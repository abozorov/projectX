package errs

import (
	"errors"
)

var (
	// ServerError
	ErrSomethingWentWrong = errors.New("something went wrong")

	// BadRequest
	ErrBadRequest        = errors.New("bad request")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrBadRequestBody    = errors.New("bad request body")
	ErrBadRequestQuery   = errors.New("bad request query")
	ErrInvalidUserId     = errors.New("invalid user id")

	// Too Many Requests
	ErrTooManyRequests = errors.New("too Many Requests")

	// NotFound
	ErrNotFound       = errors.New("not found")
	ErrUserNotFound   = errors.New("user not found")
	ErrUserIDNotFound = errors.New("user ID Not Found")
)
