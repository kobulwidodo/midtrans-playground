package domain

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrEmailConflict  = errors.New("email has already taken")
	ErrNotFound       = errors.New("not found")
	ErrPassNotMatch   = errors.New("credential doesnt match to our system")
)
