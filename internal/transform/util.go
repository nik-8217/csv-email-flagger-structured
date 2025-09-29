package transform

import "regexp"

func EmailRegex() *regexp.Regexp {
	return regexp.MustCompile(`(?i)\b[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}\b`)
}
