package db

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/senago/technopark-dbms/internal/constants"
)

func wrapErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return constants.ErrDBNotFound
	}

	return err
}
