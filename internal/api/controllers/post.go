package controllers

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/dto"
	service "github.com/senago/technopark-dbms/internal/services"
)

type PostsController struct {
	log      *customtypes.Logger
	registry *service.Registry
}

func (c *PostsController) CreatePosts(ctx *fiber.Ctx) error {
	posts := []*dto.PostData{}
	if err := json.Unmarshal(ctx.Body(), &posts); err != nil {
		return err
	}

	slugOrID := ctx.Params("slug_or_id")
	response, err := c.registry.PostsService.CreatePosts(context.Background(), slugOrID, posts)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func NewPostsController(log *customtypes.Logger, registry *service.Registry) *PostsController {
	return &PostsController{log: log, registry: registry}
}
