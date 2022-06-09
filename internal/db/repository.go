package db

import (
	"github.com/senago/technopark-dbms/internal/customtypes"
)

type Repository struct {
	UserRepository  UserRepository
	ForumRepository ForumRepository
}

func NewRepository(dbConn *customtypes.DBConn) (*Repository, error) {
	var err error
	repository := new(Repository)

	repository.UserRepository, err = NewUserRepository(dbConn)
	if err != nil {
		return nil, err
	}

	repository.ForumRepository, err = NewForumRepository(dbConn)
	if err != nil {
		return nil, err
	}

	return repository, nil
}
