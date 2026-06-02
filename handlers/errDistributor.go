package handlers

import (
	"errors"
	"net/http"

	"github.com/abozorov/projectX/package/errs"
)

func errDistributor(err error, w http.ResponseWriter) {
	if errors.Is(err, errs.ErrNotFound) ||
		errors.Is(err, errs.ErrUserNotFound) ||
		errors.Is(err, errs.ErrUserIDNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else if errors.Is(err, errs.ErrBadRequest) ||
		errors.Is(err, errs.ErrUserAlreadyExists) ||
		errors.Is(err, errs.ErrBadRequestBody) ||
		errors.Is(err, errs.ErrBadRequestQuery) ||
		errors.Is(err, errs.ErrInvalidUserId) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		http.Error(w, errs.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
	}
}
