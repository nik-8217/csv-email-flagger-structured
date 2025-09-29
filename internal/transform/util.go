package transform

import (
	"regexp"
	"strings"
)

func EmailRegex() *regexp.Regexp {
	// Email regex that matches valid email addresses
	// This regex ensures the email ends with a valid TLD (2+ chars) and doesn't end with a dot
	return regexp.MustCompile(`(?i)\b[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}(?:\.[A-Z]{2,})*\b`)
}

// IsValidEmail checks if a string is a valid email address
// This function provides additional validation beyond the regex
func IsValidEmail(email string) bool {
	// First check with regex
	if !emailRe.MatchString(email) {
		return false
	}

	// Additional check: email should not end with a dot
	email = strings.TrimSpace(email)
	if strings.HasSuffix(email, ".") {
		return false
	}

	return true
}
