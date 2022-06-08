package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/senago/technopark-dbms/internal/constants"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/dto"
	service "github.com/senago/technopark-dbms/internal/services"
)

type UserController struct {
	log      *customtypes.Logger
	registry *service.Registry
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	request := &dto.CreateUserRequest{}
	if err := Bind(ctx, request); err != nil {
		return err
	}
	request.Nickname = ctx.Params("nickname")

	response, err := c.registry.UserService.CreateUser(context.Background(), request)
	if err != nil {
		if ce, ok := err.(*constants.CodedError); ok {
			return ctx.Status(ce.Code()).JSON(response)
		}
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func NewUserController(log *customtypes.Logger, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
