package service

import (
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
)

type Registry struct {
	UserService        UserService
	ForumService       ForumService
	ForumThreadService ForumThreadService
}

func NewRegistry(log *customtypes.Logger, repository *db.Repository) *Registry {
	registry := &Registry{}

	registry.UserService = NewUserService(log, repository)
	registry.ForumService = NewForumService(log, repository)
	registry.ForumThreadService = NewForumThreadService(log, repository)

	return registry
}
