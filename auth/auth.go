package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/geobuff/api/config"
	"github.com/geobuff/api/repo"

	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	UserID      int    `json:"userId"`
	Username    string `json:"username"`
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
		return config.Values.Auth.SigningKey, nil
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
		CountryCode: user.CountryCode,
		XP:          user.XP,
		IsAdmin:     user.IsAdmin,
		IsPremium:   user.IsPremium,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.Values.Auth.SigningKey)
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
