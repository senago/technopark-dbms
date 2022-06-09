package db

import (
	"github.com/senago/technopark-dbms/internal/customtypes"
)

type Repository struct {
	UserRepository        UserRepository
	ForumRepository       ForumRepository
	ForumThreadRepository ForumThreadRepository
	PostsRepository       PostsRepository
	VotesRepository       VotesRepository
}

func NewRepository(dbConn *customtypes.DBConn) (*Repository, error) {
	var err error
	repository := &Repository{}

	repository.UserRepository, err = NewUserRepository(dbConn)
	if err != nil {
		return nil, err
	}

	repository.ForumRepository, err = NewForumRepository(dbConn)
	if err != nil {
		return nil, err
	}

	repository.ForumThreadRepository, err = NewForumThreadRepository(dbConn)
	if err != nil {
		return nil, err
	}

	repository.PostsRepository, err = NewPostsRepository(dbConn)
	if err != nil {
		return nil, err
	}

	repository.VotesRepository, err = NewVotesRepository(dbConn)
	if err != nil {
		return nil, err
	}

	return repository, nil
}
