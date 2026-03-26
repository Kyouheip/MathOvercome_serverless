package apperr

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrForbidden  = errors.New("forbidden")
	ErrOutOfRange = errors.New("index out of range")
)
