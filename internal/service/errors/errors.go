package errors

import "errors"

var (
	// User related 
	ErrUserAlreadyExists	= errors.New("user already exists")
	ErrUserNotFound 	= errors.New("user not found")

	// Password hashing
	ErrPasswordHashingFailed = errors.New("failed to hash password")

	// Password validation
	ErrPasswordTooShort	  = errors.New("password must be at least 10 characters long")
	ErrPasswordMissingDigit   = errors.New("password must contain at least one digit")
	ErrPasswordMissingUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingSymbol  = errors.New("password must contain at least one special character")
)
