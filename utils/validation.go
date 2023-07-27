package utils

import (
	"unicode"

	"github.com/go-playground/validator"
)

var Validator *validator.Validate = validator.New()

var Init = func() error {
	err := Validator.RegisterValidation("username", usernameValidation)
	if err != nil {
		return err
	}
	return Validator.RegisterValidation("password", passwordValidation)
}

func usernameValidation(fl validator.FieldLevel) bool {
	return UsernameValid(fl.Field().String())
}

func UsernameValid(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	for _, value := range username {
		if unicode.IsSpace(value) || (!unicode.IsLetter(value) && !unicode.IsNumber(value) && !unicode.IsPunct(value) && !unicode.IsSymbol(value)) {
			return false
		}
	}
	return true
}

func passwordValidation(fl validator.FieldLevel) bool {
	return PasswordValid(fl.Field().String())
}

func PasswordValid(password string) bool {
	var upper, lower, number bool
	if len(password) < 8 {
		return false
	}

	for _, value := range password {
		switch {
		case unicode.IsUpper(value):
			upper = true
		case unicode.IsLower(value):
			lower = true
		case unicode.IsNumber(value):
			number = true
		default:
			if !unicode.IsPunct(value) && !unicode.IsSymbol(value) {
				return false
			}
		}
	}

	if !upper || !lower || !number {
		return false
	}
	return true
}
