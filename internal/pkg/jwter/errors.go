package jwter

import "errors"

var (
	ErrNoID                    = errors.New("no id in payload")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrTokenExpired            = errors.New("token expired")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
)
