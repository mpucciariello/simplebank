package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile("^[a-z0-9_]+$").MatchString
	isValidFullName = regexp.MustCompile("^[a-zA-Z]+$").MatchString
)

func ValidateLength(s string, min, max int) error {
	n := len(s)
	if n < min || n > max {
		return fmt.Errorf("invalid string length: must have between %d and %d characters", min, max)
	}

	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateLength(username, 5, 20); err != nil {
		return err
	}

	if !isValidUsername(username) {
		return fmt.Errorf("invalid characters. username must only contain letters, digits or underscore")
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if err := ValidateLength(fullName, 2, 25); err != nil {
		return err
	}

	if !isValidFullName(fullName) {
		return fmt.Errorf("invalid characters. fullname must only contain letters, digits or underscore")
	}
	return nil
}

func ValidatePassword(pwd string) error {
	if err := ValidateLength(pwd, 8, 45); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(email string) error {
	if err := ValidateLength(email, 8, 25); err != nil {
		return err
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}
