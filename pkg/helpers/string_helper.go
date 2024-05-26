package helpers

import (
	"errors"
	"regexp"
	"strings"
)

func IsEmpty(str string) bool {
	str = strings.TrimSpace(str)
	if str == "" || len(str) <= 0 {
		return true
	}
	return false
}

func IdentifiesColumnValue(email, phone string) (string, string) {
	// Regular expression patterns for email and phone number validation
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	phonePattern := `^\+?[1-9]\d{1,14}$`

	// Check if the input matches the email pattern
	if matched, _ := regexp.MatchString(emailPattern, email); matched {
		return "email", email
	}

	// Check if the input matches the phone number pattern
	if matched, _ := regexp.MatchString(phonePattern, phone); matched {
		return "phone", phone
	}
	return "", ""
}

func EmailOrPhone(input string) string {
	// Regular expression patterns for email and phone number validation
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	phonePattern := `^\+?[1-9]\d{1,14}$`

	// Check if the input matches the email pattern
	if matched, _ := regexp.MatchString(emailPattern, input); matched {
		return "email"
	}

	// Check if the input matches the phone number pattern
	if matched, _ := regexp.MatchString(phonePattern, input); matched {
		return "phone"
	}
	return "unknown"
}

func ValidateEmailOrPhone(email, phone string) error {
	// Check that either Email or Phone is present, but not both
	if (email != "" && phone != "") || (email == "" && phone == "") {
		return errors.New("either Email or Phone should be present")
	}
	// Check that at least one field is non-empty
	if email == "" && phone == "" {
		return errors.New("at least one of Email or Phone should be non-empty")
	}
	// Validate email format
	if email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return errors.New("invalid email format")
		}
	}
	// Validate phone number length
	if phone != "" && len(phone) != 10 {
		phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
		if !phoneRegex.MatchString(phone) {
			return errors.New("phone number should be 10 digits")
		}
	}
	return nil
}
