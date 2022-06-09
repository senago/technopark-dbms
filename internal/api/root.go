package api

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/senago/technopark-dbms/internal/api/controllers"
	"github.com/senago/technopark-dbms/internal/customtypes"
)

type APIService struct {
	log    *customtypes.Logger
	router *fiber.App
}

func (svc *APIService) Serve(addr string) {
	svc.log.Fatal(svc.router.Listen(addr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	return svc.router.Shutdown()
}

func NewAPIService(log *customtypes.Logger, dbConn *customtypes.DBConn) (*APIService, error) {
	svc := &APIService{
		log: log,
		router: fiber.New(fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		}),
	}

	controllersRegistry := controllers.NewRegistry(log, dbConn)

	// TODO: Remove log for better performance
	api := svc.router.Group("/api", recover.New(), logger.New())

	api.Post("/user/:nickname/create", controllersRegistry.UserController.CreateUser)
	api.Get("/user/:nickname/profile", controllersRegistry.UserController.GetUserProfile)
	api.Post("/user/:nickname/profile", controllersRegistry.UserController.UpdateUserProfile)

	api.Post("/forum/create", controllersRegistry.ForumController.CreateForum)
	api.Get("/forum/:slug/details", controllersRegistry.ForumController.GetForumBySlug)
	api.Get("/forum/:slug/threads", controllersRegistry.ForumController.GetForumThreads)
	api.Get("/forum/:slug/users", controllersRegistry.ForumController.GetForumUsers)

	api.Post("/forum/:slug/create", controllersRegistry.ForumThreadController.CreateForumThread)

	api.Post("/thread/:slug_or_id/create", controllersRegistry.PostsController.CreatePosts)
	api.Post("/thread/:slug_or_id/vote", controllersRegistry.ForumThreadController.UpdateVote)
	api.Get("/thread/:slug_or_id/details", controllersRegistry.ForumThreadController.GetForumThreadDetails)
	api.Get("/thread/:slug_or_id/posts", controllersRegistry.PostsController.GetPosts)
	api.Post("/thread/:slug_or_id/details", controllersRegistry.ForumThreadController.UpdateForumThread)

	return svc, nil
}
