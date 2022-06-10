package controllers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
)

type ServiceController struct {
	log *customtypes.Logger
	db  *db.Repository
}

func (c *ServiceController) Status(ctx *fiber.Ctx) error {
	response, err := c.db.ServiceRepository.Status(context.Background())
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(response)
}

func (c *ServiceController) Delete(ctx *fiber.Ctx) error {
	err := c.db.ServiceRepository.Delete(context.Background())
	if err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}

func NewServiceController(log *customtypes.Logger, db *db.Repository) *ServiceController {
	return &ServiceController{log: log, db: db}
}
