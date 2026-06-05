package handlers

import (
	"errors"
	"net/http"

	"github.com/abozorov/projectX/package/errs"
)

func errDistributor(err error, w http.ResponseWriter) {

	switch {
	// http.StatusNotFound
	case errors.Is(err, errs.ErrNotFound):
		http.Error(w, errs.ErrNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, errs.ErrUserNotFound):
		http.Error(w, errs.ErrUserNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, errs.ErrUserIDNotFound):
		http.Error(w, errs.ErrUserIDNotFound.Error(), http.StatusNotFound)

	// http.StatusBadRequest
	case errors.Is(err, errs.ErrBadRequest):
		http.Error(w, errs.ErrBadRequest.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrUserAlreadyExists):
		http.Error(w, errs.ErrUserAlreadyExists.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrBadRequestBody):
		http.Error(w, errs.ErrBadRequestBody.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrBadRequestQuery):
		http.Error(w, errs.ErrBadRequestQuery.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrInvalidUserId):
		http.Error(w, errs.ErrInvalidUserId.Error(), http.StatusBadRequest)

	// http.StatusTooManyRequests
	case errors.Is(err, errs.ErrTooManyRequests):
		http.Error(w, errs.ErrTooManyRequests.Error(), http.StatusTooManyRequests)

	// http.StatusGatewayTimeout
	case errors.Is(err, errs.ErrTimeoutExceeded):
		http.Error(w, errs.ErrTimeoutExceeded.Error(), http.StatusGatewayTimeout)

	// http.StatusInternalServerError
	default:
		http.Error(w, errs.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)

	}
}
