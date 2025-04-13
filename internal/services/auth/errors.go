package auth

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrIncorrectData = errors.New("incorrect password or email")
	ErrAlreadyExists = errors.New("user already exists")
	ErrIncorrectRole = errors.New("incorrect access level")
	ErrBadRequest    = errors.New("bad request")
)
