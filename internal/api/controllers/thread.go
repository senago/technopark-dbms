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

func (c *ForumThreadController) UpdateVote(ctx *fiber.Ctx) error {
	request := &dto.UpdateVoteRequest{}
	if err := Bind(ctx, request); err != nil {
		return err
	}

	slugOrID := ctx.Params("slug_or_id")
	response, err := c.registry.ForumThreadService.UpdateVote(context.Background(), slugOrID, request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func (c *ForumThreadController) GetForumThreadDetails(ctx *fiber.Ctx) error {
	slugOrID := ctx.Params("slug_or_id")

	response, err := c.registry.ForumThreadService.GetThreadDetails(context.Background(), slugOrID)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func (c *ForumThreadController) UpdateForumThread(ctx *fiber.Ctx) error {
	request := &dto.UpdateForumThreadRequest{}
	if err := Bind(ctx, request); err != nil {
		return err
	}
	slugOrID := ctx.Params("slug_or_id")

	response, err := c.registry.ForumThreadService.UpdateForumThread(context.Background(), slugOrID, request)
	if err != nil {
		return err
	}

	return ctx.Status(response.Code).JSON(response.Data)
}

func NewForumThreadController(log *customtypes.Logger, registry *service.Registry) *ForumThreadController {
	return &ForumThreadController{log: log, registry: registry}
}
