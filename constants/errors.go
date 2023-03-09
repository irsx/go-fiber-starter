package constants

import "errors"

var (
	ErrInvalidAuth  = errors.New("username or password is wrong")
	ErrGUIDRequired = errors.New("guid params is required")
	ErrEmailExist   = errors.New("email address already exist")
)
