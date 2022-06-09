package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/dto"
	service "github.com/senago/technopark-dbms/internal/services"
)

type ForumController struct {
	log      *customtypes.Logger
	registry *service.Registry
}

func (c *ForumController) CreateForum(ctx *fiber.Ctx) error {
	request := &dto.CreateForumRequest{}
	if err := Bind(ctx, request); err != nil {
		return err
	}

	response, err := c.registry.ForumService.CreateForum(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func (c *ForumController) GetForumBySlug(ctx *fiber.Ctx) error {
	request := &dto.GetForumBySlugRequest{Slug: ctx.Params("slug")}

	response, err := c.registry.ForumService.GetForumBySlug(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func NewForumController(log *customtypes.Logger, registry *service.Registry) *ForumController {
	return &ForumController{log: log, registry: registry}
}
