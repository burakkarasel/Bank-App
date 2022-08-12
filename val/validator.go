package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

// ValidateString checks if a given string has correct characters
func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)

	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

// ValidateUsername checks if a given username matches its validations
func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits or underscore")
	}
	return nil
}

// ValidatePassword check for a given password has correct characters count
func ValidatePassword(value string) error {
	return ValidateString(value, 8, 100)
}

// ValidateEmail checks if given email is valid or not
func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}

// ValidateFullName checks for a given full name input if it meets the validations
func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}
