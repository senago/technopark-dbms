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
	ServiceRepository     ServiceRepository
}

func NewRepository(dbConn *customtypes.DBConn) (*Repository, error) {
	repository := &Repository{}

	repository.UserRepository = NewUserRepository(dbConn)
	repository.ForumRepository = NewForumRepository(dbConn)
	repository.ForumThreadRepository = NewForumThreadRepository(dbConn)
	repository.PostsRepository = NewPostsRepository(dbConn)
	repository.VotesRepository = NewVotesRepository(dbConn)
	repository.VotesRepository = NewVotesRepository(dbConn)
	repository.ServiceRepository = NewServiceRepository(dbConn)

	return repository, nil
}
