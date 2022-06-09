package controllers

import (
	"context"
	"encoding/json"
	"strconv"

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

func (c *PostsController) GetPosts(ctx *fiber.Ctx) error {
	slugOrID := ctx.Params("slug_or_id")
	sort := ctx.Query("sort", "flat")
	since, _ := strconv.ParseInt(ctx.Query("since", "-1"), 10, 64)
	desc, _ := strconv.ParseBool(ctx.Query("desc"))
	limit, _ := strconv.ParseInt(ctx.Query("limit", "100"), 10, 64)

	response, err := c.registry.PostsService.GetPosts(context.Background(), slugOrID, sort, since, desc, limit)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func (c *PostsController) GetPostDetails(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
	request := &dto.GetPostDetailsRequest{ID: id, Related: ctx.Query("related")}

	response, err := c.registry.PostsService.GetPostDetails(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func NewPostsController(log *customtypes.Logger, registry *service.Registry) *PostsController {
	return &PostsController{log: log, registry: registry}
}
