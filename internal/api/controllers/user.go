package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/dto"
	service "github.com/senago/technopark-dbms/internal/services"
)

type UserController struct {
	log      *customtypes.Logger
	registry *service.Registry
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	request := &dto.CreateUserRequest{Nickname: ctx.Params("nickname")}
	if err := ctx.BodyParser(request); err != nil {
		return err
	}

	response, err := c.registry.UserService.CreateUser(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func (c *UserController) GetUserProfile(ctx *fiber.Ctx) error {
	request := &dto.GetUserProfileRequest{Nickname: ctx.Params("nickname")}

	response, err := c.registry.UserService.GetUserProfile(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func (c *UserController) UpdateUserProfile(ctx *fiber.Ctx) error {
	request := &dto.UpdateUserProfileRequest{Nickname: ctx.Params("nickname")}
	if err := ctx.BodyParser(request); err != nil {
		return err
	}

	response, err := c.registry.UserService.UpdateUserProfile(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func NewUserController(log *customtypes.Logger, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
