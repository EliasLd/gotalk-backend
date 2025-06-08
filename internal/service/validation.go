package service

import (
	"regexp"
	"github.com/EliasLd/gotalk-backend/internal/service/errors"
)

var (
	digitRegex	= regexp.MustCompile(`[0-9]`)
	upperRegex	= regexp.MustCompile(`[A-Z]`)
	lowerRegex	= regexp.MustCompile(`[a-z]`)
	symbolRegex	= regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`)
)	

func ValidatePassword(pw string) error {
	if len(pw) < 10 {
		return errors.ErrPasswordTooShort
	}
	if !digitRegex.MatchString(pw) {
		return errors.ErrPasswordMissingDigit
	}
	if !upperRegex.MatchString(pw) {
		return errors.ErrPasswordMissingUpper
	}
	if !lowerRegex.MatchString(pw) {
		return errors.ErrPasswordMissingLower
	}
	if !symbolRegex.MatchString(pw) {
		return errors.ErrPasswordMissingSymbol
	}
	return nil
}

