package controllers

import (
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
	service "github.com/senago/technopark-dbms/internal/services"
)

type Registry struct {
	UserController        *UserController
	ForumController       *ForumController
	ForumThreadController *ForumThreadController
}

func NewRegistry(log *customtypes.Logger, dbConn *customtypes.DBConn) *Registry {
	repository, err := db.NewRepository(dbConn)
	if err != nil {
		log.Fatal(err)
	}
	serviceRegistry := service.NewRegistry(log, repository)

	registry := &Registry{}

	registry.UserController = NewUserController(log, serviceRegistry)
	registry.ForumController = NewForumController(log, serviceRegistry)
	registry.ForumThreadController = NewForumThreadController(log, serviceRegistry)

	return registry
}
