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
	"github.com/geobuff/api/validation"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	UserID                  int       `json:"userId"`
	AvatarId                int       `json:"avatarId"`
	AvatarName              string    `json:"avatarName"`
	AvatarDescription       string    `json:"avatarDescription"`
	AvatarPrimaryImageUrl   string    `json:"avatarPrimaryImageUrl"`
	AvatarSecondaryImageUrl string    `json:"avatarSecondaryImageUrl"`
	Username                string    `json:"username"`
	Email                   string    `json:"email"`
	CountryCode             string    `json:"countryCode"`
	XP                      int       `json:"xp"`
	IsAdmin                 bool      `json:"isAdmin"`
	IsPremium               bool      `json:"isPremium"`
	Joined                  time.Time `json:"joined"`
	jwt.StandardClaims
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterDto struct {
	AvatarId    int    `json:"avatarId" validate:"required"`
	Username    string `json:"username" validate:"required,username"`
	Email       string `json:"email" validate:"required,email"`
	CountryCode string `json:"countryCode" validate:"required"`
	Password    string `json:"password" validate:"required,password"`
}

type RefreshTokenDto struct {
	Token string `json:"token"`
}

type PasswordResetDto struct {
	Email string `json:"email"`
}

type ResetTokenUpdateDto struct {
	UserID   int    `json:"userId"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

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

	user, err := repo.GetAuthUserUsingEmail(loginDto.Email)
	if err == sql.ErrNoRows {
		http.Error(writer, "Invalid email. Please try again.", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDto.Password)); err != nil {
		http.Error(writer, "Invalid password. Please try again.", http.StatusBadRequest)
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

	err = validation.Validator.Struct(registerDto)
	if err != nil {
		http.Error(writer, "There was a validation error on sign up. Please ensure all fields are filled in correctly and try again.", http.StatusBadRequest)
		return
	}

	usernameExists, err := repo.UsernameExists(registerDto.Username)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if usernameExists {
		http.Error(writer, "Username already in use. Please choose another and try again.", http.StatusBadRequest)
		return
	}

	emailExists, err := repo.EmailExists(registerDto.Email)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if emailExists {
		http.Error(writer, "Email already in use. Please choose another and try again.", http.StatusBadRequest)
		return
	}

	passwordHash, err := hashPassword([]byte(registerDto.Password))
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	newUser := repo.User{
		AvatarId:     registerDto.AvatarId,
		Username:     registerDto.Username,
		Email:        registerDto.Email,
		PasswordHash: passwordHash,
		CountryCode:  registerDto.CountryCode,
	}

	if err := repo.InsertUser(newUser); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func RefreshToken(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var refreshTokenDto RefreshTokenDto
	err = json.Unmarshal(requestBody, &refreshTokenDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	claims, err := getClaims(refreshTokenDto.Token)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := repo.GetAuthUserUsingEmail(claims.Email)
	if err == sql.ErrNoRows {
		http.Error(writer, fmt.Sprintf("User with email %s does not exist.", claims.Email), http.StatusBadRequest)
		return
	} else if err != nil {
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

func EmailExists(writer http.ResponseWriter, request *http.Request) {
	exists, err := repo.EmailExists(mux.Vars(request)["email"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(exists)
}

func UsernameExists(writer http.ResponseWriter, request *http.Request) {
	exists, err := repo.UsernameExists(mux.Vars(request)["username"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(exists)
}

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

	user, err := repo.GetAuthUserUsingEmail(passwordResetDto.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, fmt.Sprintf("User with email %s does not exist.", passwordResetDto.Email), http.StatusBadRequest)
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
	_, err = email.SendResetToken(passwordResetDto.Email, resetLink)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func ResetTokenValid(writer http.ResponseWriter, request *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(request)["userId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := repo.GetAuthUser(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, fmt.Sprintf("User with id %d does not exist.", userID), http.StatusBadRequest)
			return
		}

		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	valid := isResetTokenValid(user.PasswordResetToken, mux.Vars(request)["token"], user.PasswordResetExpiry)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(valid)
}

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

	user, err := repo.GetAuthUser(resetTokenUpdateDto.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, fmt.Sprintf("User with id %d does not exist.", resetTokenUpdateDto.UserID), http.StatusBadRequest)
			return
		}

		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if !isResetTokenValid(user.PasswordResetToken, resetTokenUpdateDto.Token, user.PasswordResetExpiry) {
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

var isResetTokenValid = func(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool {
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

var getToken = func(request *http.Request) (string, error) {
	header := request.Header.Get("Authorization")
	if len(header) < 8 {
		return "", errors.New("token missing or invalid length")
	}
	return header[7:], nil
}

var getClaims = func(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH_SIGNING_KEY")), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

var buildToken = func(user repo.AuthUserDto) (string, error) {
	claims := CustomClaims{
		UserID:                  user.ID,
		AvatarId:                user.AvatarId,
		AvatarName:              user.AvatarName,
		AvatarDescription:       user.AvatarDescription,
		AvatarPrimaryImageUrl:   user.AvatarPrimaryImageUrl,
		AvatarSecondaryImageUrl: user.AvatarSecondaryImageUrl,
		Username:                user.Username,
		Email:                   user.Email,
		CountryCode:             user.CountryCode,
		XP:                      user.XP,
		IsAdmin:                 user.IsAdmin,
		IsPremium:               user.IsPremium,
		Joined:                  user.Joined,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
			Issuer:    os.Getenv("AUTH_ISSUER"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("AUTH_SIGNING_KEY")))
}

var hashPassword = func(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
