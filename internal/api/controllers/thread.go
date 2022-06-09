package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/dto"
	service "github.com/senago/technopark-dbms/internal/services"
)

type ForumThreadController struct {
	log      *customtypes.Logger
	registry *service.Registry
}

func (c *ForumThreadController) CreateForumThread(ctx *fiber.Ctx) error {
	request := &dto.CreateForumThreadRequest{Forum: ctx.Params("slug")}
	if err := Bind(ctx, request); err != nil {
		return err
	}

	response, err := c.registry.ForumThreadService.CreateForumThread(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func NewForumThreadController(log *customtypes.Logger, registry *service.Registry) *ForumThreadController {
	return &ForumThreadController{log: log, registry: registry}
}
