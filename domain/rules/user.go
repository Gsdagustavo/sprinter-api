package rules

import (
	"net/mail"
	"unicode"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// Password rules
const (
	PasswordMinLetters = 8
	PasswordMinNumbers = 1
	PasswordMinSpecial = 1
	PasswordMaxLetters = 32
)

// Name/Username rules
const (
	NameMinLetters = 3
	NameMaxLetters = 32
)

const (
	BiographyMinLetters = 0
	BiographyMaxLetters = 100
)

func ValidateEmail(email string) bool {
	res, err := mail.ParseAddress(email)
	return err == nil && res.Address == email
}

func ValidatePassword(password string) bool {
	var letters int
	var number bool
	var special bool

	for _, char := range password {
		letters++

		switch {
		case unicode.IsNumber(char):
			number = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special = true
		}
	}

	return letters >= PasswordMinLetters && letters <= PasswordMaxLetters && number && special
}

func ValidateName(name string) bool {
	if len(name) < NameMinLetters || len(name) > NameMaxLetters {
		return false
	}

	return true
}

func ValidateBiography(biography string) bool {
	if len(biography) > BiographyMaxLetters {
		return false
	}

	return true
}

func ValidateCredentials(credentials entities.UserCredentials) bool {
	return ValidateEmail(credentials.Email) && ValidatePassword(credentials.Password)
}
