package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type APIService struct {
	log    *zap.SugaredLogger
	router *fiber.App
}

func (svc *APIService) Serve(addr string) {
	svc.log.Fatal(svc.router.Listen(addr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	return svc.router.Shutdown()
}

func NewAPIService(log *zap.SugaredLogger, dbPool *pgxpool.Pool) (*APIService, error) {
	svc := &APIService{
		log:    log,
		router: fiber.New(),
	}

	return svc, nil
}
