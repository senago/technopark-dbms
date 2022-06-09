package api

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/senago/technopark-dbms/internal/api/controllers"
	"github.com/senago/technopark-dbms/internal/constants"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/dto"
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
		log: log,
		router: fiber.New(fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				if e, ok := err.(*fiber.Error); ok {
					code = e.Code
				} else if e, ok := err.(*constants.CodedError); ok {
					code = e.Code()
				}
				return ctx.Status(code).JSON(&dto.ErrorResponse{Message: err.Error()})
			},
		}),
	}

	controllersRegistry := controllers.NewRegistry(log, dbConn)

	// TODO: Remove log for better performance
	api := svc.router.Group("/api", recover.New(), logger.New())

	api.Post("/user/:nickname/create", controllersRegistry.UserController.CreateUser)
	api.Get("/user/:nickname/profile", controllersRegistry.UserController.GetUserProfile)
	api.Post("/user/:nickname/profile", controllersRegistry.UserController.UpdateUserProfile)

	return svc, nil
}
