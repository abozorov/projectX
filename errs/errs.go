package errs

import (
	"errors"
)

var (
	// BadRequest
	ErrBadRequest        = errors.New("ErrBadRequest")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrBadRequestBody    = errors.New("bad request body")
	ErrBadRequestQuery   = errors.New("bad request query")

	// NotFound
	ErrUserNotFound   = errors.New("user not found")
	ErrUserIDNotFound = errors.New("user ID Not Found")
)
