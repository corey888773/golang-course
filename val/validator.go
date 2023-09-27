package val

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z]+([\s]?[a-zA-Z]*)*$`).MatchString
)

func ValidateString(str string, minLen int, maxLen int) error {
	if len(str) > maxLen ||
		len(str) < minLen {
		return fmt.Errorf("must contain %d-%d characters", minLen, maxLen)
	}

	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(username) {
		return errors.New("username must contain only small/big case letters, digits and underscore")
	}
	return nil
}

func ValidatePassword(password string) error {
	if err := ValidateString(password, 6, 100); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 6, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

func ValidateFullName(fullName string) error {
	if err := ValidateString(fullName, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(fullName) {
		return errors.New("full name must contain only small/big case letters and spaces")
	}
	return nil
}
