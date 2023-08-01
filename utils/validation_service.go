package utils

import (
	"unicode"

	"github.com/go-playground/validator"
)

type IValidationService interface {
	GetValidator() *validator.Validate
	UsernameValid(username string) bool
	PasswordValid(password string) bool
}

type ValidationService struct {
	Validator *validator.Validate
}

func NewValidationService() *ValidationService {
	validator := validator.New()
	vs := ValidationService{
		Validator: validator,
	}

	vs.Validator.RegisterValidation("username", vs.usernameValidation)
	vs.Validator.RegisterValidation("password", vs.passwordValidation)
	return &vs
}

func (v *ValidationService) GetValidator() *validator.Validate {
	return v.Validator
}

func (v *ValidationService) usernameValidation(fl validator.FieldLevel) bool {
	return v.UsernameValid(fl.Field().String())
}

func (v *ValidationService) UsernameValid(username string) bool {
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

func (v *ValidationService) passwordValidation(fl validator.FieldLevel) bool {
	return v.PasswordValid(fl.Field().String())
}

func (v *ValidationService) PasswordValid(password string) bool {
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
