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

func ValidateUserInformation(userInformation entities.AccountInformation) error {
	err := ValidateName(userInformation.Username)
	if err != nil {
		return err
	}

	err = ValidateBiography(userInformation.Biography)
	if err != nil {
		return err
	}

	return nil
}

func ValidateEmail(email string) bool {
	res, err := mail.ParseAddress(email)
	return err == nil && res.Address == email
}

func ValidatePassword(password string) error {
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
		return derr.WeakPassword
	}

	return nil
}

func ValidateName(name string) error {
	if len(name) < NameMinLetters {
		return derr.NameIsTooShort
	}

	if len(name) > NameMaxLetters {
		return derr.NameIsTooLong
	}

	return nil
}

func ValidateBiography(biography string) error {
	if len(biography) > BiographyMaxLetters {
		return derr.BiographyIsTooLong
	}

	return nil
}

func ValidateCredentials(credentials entities.UserCredentials) error {
	if valid := ValidateEmail(credentials.Email); !valid {
		return derr.InvalidEmail
	}

	if err := ValidatePassword(credentials.Password); err != nil {
		return derr.WeakPassword
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
