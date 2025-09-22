package domain

import "errors"

var (
	ErrDuplicateUser  = errors.New("user already exists")
	ErrUserCannotSave = errors.New("user cannot be saved")
	ErrPersonNotFound = errors.New("person not found")

	ErrGettingUserByEmail        = errors.New("error getting user by email")
	ErrNotFoundUserByEmail       = errors.New("user not found by email")
	ErrUserCannotFound           = errors.New("user cannot be found")
	ErrUserCannotGet             = errors.New("user cannot be retrieved")
	ErrorEmailNotVerified        = errors.New("email not verified")
	ErrVerificationTokenNotFound = errors.New("verification token not found")
	ErrTokenExpired              = errors.New("token expired")
	ErrTokenAlreadyUsed          = errors.New("token already used")

	ErrInvalidJSONFormat = errors.New("invalid JSON format")
)