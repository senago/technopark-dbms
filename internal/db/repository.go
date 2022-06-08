package db

import (
	"github.com/senago/technopark-dbms/internal/customtypes"
)

type Repository struct {
}

func NewRepository(dbConn *customtypes.DBConn) (*Repository, error) {
	repository := new(Repository)

	return repository, nil
}
