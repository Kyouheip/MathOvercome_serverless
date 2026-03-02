package apperr

import "errors"

// 認証系
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

// リソース系
var (
	ErrNotFound   = errors.New("not found")
	ErrForbidden  = errors.New("forbidden")
	ErrOutOfRange = errors.New("index out of range")
)
