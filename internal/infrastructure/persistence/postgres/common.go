package postgres

import (
	"errors"

	"github.com/go-pg/pg/v9"
)

func isNoRowsError(err error) bool {
	return errors.Is(err, pg.ErrNoRows)
}

func isMultiRowsError(err error) bool {
	return errors.Is(err, pg.ErrMultiRows)
}

func getConstraint(err error) string {
	pgErr, ok := err.(pg.Error)
	if ok && pgErr.IntegrityViolation() {
		return pgErr.Field(byte('n'))
	}

	return ""
}
