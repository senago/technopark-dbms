package service

import (
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
)

type Registry struct {
	UserService  UserService
	ForumService ForumService
}

func NewRegistry(log *customtypes.Logger, repository *db.Repository) *Registry {
	registry := new(Registry)

	registry.UserService = NewUserService(log, repository)
	registry.ForumService = NewForumService(log, repository)

	return registry
}
