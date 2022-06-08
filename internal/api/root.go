package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/senago/technopark-dbms/internal/api/controllers"
	"github.com/senago/technopark-dbms/internal/customtypes"
)

type APIService struct {
	log    *customtypes.Logger
	router *fiber.App
}

func (svc *APIService) Serve(addr string) {
	svc.log.Fatal(svc.router.Listen(addr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	return svc.router.Shutdown()
}

func NewAPIService(log *customtypes.Logger, dbConn *customtypes.DBConn) (*APIService, error) {
	svc := &APIService{
		log:    log,
		router: fiber.New(),
	}

	controllersRegistry := controllers.NewRegistry(log, dbConn)

	svc.router.Use(recover.New())

	svc.router.Post("/user/:nickname/create", controllersRegistry.UserController.CreateUser)

	return svc, nil
}
