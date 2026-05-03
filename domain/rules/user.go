package rules

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
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

func validateName(name string) error {
	if len(name) < NameMinLetters {
		return derr.NameIsTooShort
	}

	if len(name) > NameMaxLetters {
		return derr.NameIsTooLong
	}

	return nil
}

func validateEmail(email string) bool {
	res, err := mail.ParseAddress(email)
	return err == nil && res.Address == email
}

func validatePassword(password string) bool {
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

	strong := letters >= PasswordMinLetters && letters <= PasswordMaxLetters && number && special
	if !strong {
		return false
	}

	return true
}

func validateBiography(biography string) error {
	if len(biography) > BiographyMaxLetters {
		return derr.BiographyIsTooLong
	}

	return nil
}

func ValidateUserInformation(userInformation entities.UserInformation) error {
	err := validateName(userInformation.Username)
	if err != nil {
		return err
	}

	err = validateBiography(userInformation.Biography)
	if err != nil {
		return err
	}

	return nil
}

func ValidateRegisterCredentials(credentials entities.UserCredentials) error {
	err := validateName(credentials.Name)
	if err != nil {
		return err
	}

	valid := validateEmail(credentials.Email)
	if !valid {
		return derr.InvalidEmail
	}

	valid = validatePassword(credentials.Password)
	if !valid {
		return derr.WeakPassword
	}

	return nil
}

func ValidateLoginCredentials(credentials entities.UserCredentials) error {
	emailValid := validateEmail(credentials.Email)
	if !emailValid {
		return derr.InvalidEmail
	}

	passwordValid := validatePassword(credentials.Password)
	if !passwordValid {
		return derr.InvalidPassword
	}

	return nil
}

var toUsernameRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func ToUsername(input string) string {
	s := strings.ToLower(input)
	s = toUsernameRegex.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	return s
}
