package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func Bind(ctx *fiber.Ctx, out interface{}) error {
	if err := ctx.BodyParser(out); err != nil {
		return err
	}

	if err := ctx.QueryParser(out); err != nil {
		return err
	}

	if err := ctx.ReqHeaderParser(out); err != nil {
		return err
	}

	return nil
}
