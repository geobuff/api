package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/geobuff/api/email"
	"github.com/geobuff/api/repo"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	UserID      int    `json:"userId"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	CountryCode string `json:"countryCode"`
	XP          int    `json:"xp"`
	IsAdmin     bool   `json:"isAdmin"`
	IsPremium   bool   `json:"isPremium"`
	jwt.StandardClaims
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterDto struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	CountryCode string `json:"countryCode"`
	Password    string `json:"password"`
}

type PasswordResetDto struct {
	Email string `json:"email"`
}

type ResetTokenUpdateDto struct {
	UserID   int    `json:"userId"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

// Login verifies a user before returning a token with relevant information.
func Login(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var loginDto LoginDto
	err = json.Unmarshal(requestBody, &loginDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := repo.GetUserUsingEmail(loginDto.Email)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	match, err := passwordsMatch([]byte(user.PasswordHash), []byte(loginDto.Password))
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if !match {
		http.Error(writer, fmt.Sprint("Invalid password. Please try again.\n", err), http.StatusBadRequest)
		return
	}

	token, err := buildToken(user)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(token)
}

// Register creates a new user in the database and returns a token holding relevant information.
func Register(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var registerDto RegisterDto
	err = json.Unmarshal(requestBody, &registerDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	usernameExists, err := repo.UsernameExists(registerDto.Username)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if usernameExists {
		http.Error(writer, fmt.Sprint("Username already in use. Please choose another and try again.\n", err), http.StatusBadRequest)
		return
	}

	emailExists, err := repo.EmailExists(registerDto.Email)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if emailExists {
		http.Error(writer, fmt.Sprint("Email already in use. Please choose another and try again.\n", err), http.StatusBadRequest)
		return
	}

	passwordHash, err := hashPassword([]byte(registerDto.Password))
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	newUser := repo.User{
		Username:     registerDto.Username,
		Email:        registerDto.Email,
		PasswordHash: passwordHash,
		CountryCode:  registerDto.CountryCode,
	}

	user, err := repo.InsertUser(newUser)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	token, err := buildToken(user)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(token)
}

// SendResetToken sends a password reset email if the user is valid.
func SendResetToken(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var passwordResetDto PasswordResetDto
	err = json.Unmarshal(requestBody, &passwordResetDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := repo.GetUserUsingEmail(passwordResetDto.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, fmt.Sprintf("User with email %s does not exist.\n", passwordResetDto.Email), http.StatusBadRequest)
			return
		}

		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	guid := uuid.New().String()
	expiryDate := time.Now().AddDate(0, 0, 1)
	err = repo.SetPasswordResetValues(user.ID, guid, expiryDate)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	resetLink := fmt.Sprintf("%s/reset-password/%d/%s", os.Getenv("SITE_URL"), user.ID, guid)
	err = email.SendResetToken(passwordResetDto.Email, resetLink)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

// ResetTokenValid checks if a token matches the user's database token.
func ResetTokenValid(writer http.ResponseWriter, request *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(request)["userId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := repo.GetUser(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, fmt.Sprintf("User with id %d does not exist.\n", userID), http.StatusBadRequest)
			return
		}

		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	valid := resetTokenValid(user.PasswordResetToken, mux.Vars(request)["token"], user.PasswordResetExpiry)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(valid)
}

// UpdatePasswordUsingToken updates the users password using the email reset token.
func UpdatePasswordUsingToken(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var resetTokenUpdateDto ResetTokenUpdateDto
	err = json.Unmarshal(requestBody, &resetTokenUpdateDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := repo.GetUser(resetTokenUpdateDto.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, fmt.Sprintf("User with id %d does not exist.\n", resetTokenUpdateDto.UserID), http.StatusBadRequest)
			return
		}

		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if !resetTokenValid(user.PasswordResetToken, resetTokenUpdateDto.Token, user.PasswordResetExpiry) {
		http.Error(writer, "Password reset token is not valid.", http.StatusBadRequest)
		return
	}

	passwordHash, err := hashPassword([]byte(resetTokenUpdateDto.Password))
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	err = repo.ResetPassword(user.ID, passwordHash)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	user.PasswordResetToken = sql.NullString{}
	user.PasswordResetExpiry = sql.NullTime{}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}

func resetTokenValid(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool {
	return userToken.Valid && expiry.Valid && userToken.String == requestToken && expiry.Time.Sub(time.Now()) > 0
}

var ValidUser = func(request *http.Request, id int) (int, error) {
	token, err := getToken(request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	claims, err := getClaims(token)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if claims.UserID != id && !claims.IsAdmin {
		return http.StatusUnauthorized, errors.New("invalid permissions to make request")
	}

	return http.StatusOK, nil
}

var IsAdmin = func(request *http.Request) (int, error) {
	token, err := getToken(request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	claims, err := getClaims(token)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !claims.IsAdmin {
		return http.StatusUnauthorized, errors.New("invalid permissions to make request")
	}

	return http.StatusOK, nil
}

func getToken(request *http.Request) (string, error) {
	header := request.Header.Get("Authorization")
	if len(header) < 8 {
		return "", errors.New("token missing or invalid length")
	}
	return header[7:], nil
}

func getClaims(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH_SIGNING_KEY")), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func buildToken(user repo.User) (string, error) {
	claims := CustomClaims{
		UserID:      user.ID,
		Username:    user.Username,
		Email:       user.Email,
		CountryCode: user.CountryCode,
		XP:          user.XP,
		IsAdmin:     user.IsAdmin,
		IsPremium:   user.IsPremium,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(0, 0, 3).Unix(),
			Issuer:    os.Getenv("AUTH_ISSUER"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("AUTH_SIGNING_KEY")))
}

func hashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func passwordsMatch(hashedPassword, plainPassword []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, plainPassword); err != nil {
		return false, err
	}
	return true, nil
}
