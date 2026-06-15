package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/abozorov/projectX/pkg/errs"
	"github.com/lib/pq"
)

// sql.error -> pkg.errs
func distributor(err error) error {
	if err == nil {
		return err
	}

	// context / context
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotFound

	case errors.Is(err, context.DeadlineExceeded):
		return errs.ErrTimeoutExceeded

	case errors.Is(err, context.Canceled):
		return errs.ErrTimeoutExceeded
	}

	// postgres
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {

		// unique_violation
		case "23505":
			return errs.ErrUserAlreadyExists

		// foreign_key_violation
		case "23503":
			return errs.ErrBadRequest

		// not_null_violation
		case "23502":
			return errs.ErrBadRequestBody

		// check_violation
		case "23514":
			return errs.ErrBadRequest

		// invalid_text_representation
		case "22P02":
			return errs.ErrBadRequest

		default:
			return errs.ErrSomethingWentWrong
		}
	}

	return errs.ErrSomethingWentWrong
}
