//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/senago/technopark-dbms/internal/constants"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
	"github.com/senago/technopark-dbms/internal/model/core"
	"github.com/senago/technopark-dbms/internal/model/dto"
)

type ForumService interface {
	CreateForum(ctx context.Context, request *dto.CreateForumRequest) (*dto.Response, error)
}

type forumServiceImpl struct {
	log *customtypes.Logger
	db  *db.Repository
}

func (svc *forumServiceImpl) CreateForum(ctx context.Context, request *dto.CreateForumRequest) (*dto.Response, error) {
	if forum, err := svc.db.ForumRepository.GetForumBySlug(ctx, request.Slug); err != nil {
		if !errors.Is(err, constants.ErrDBNotFound) {
			return nil, err
		}
	} else {
		return &dto.Response{Data: forum, Code: http.StatusConflict}, nil
	}

	forum := &core.Forum{Title: request.Title, User: request.User, Slug: request.Slug}
	if err := svc.db.ForumRepository.CreateForum(ctx, forum); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.User)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}

	forum, err := svc.db.ForumRepository.GetForumBySlug(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: forum, Code: http.StatusCreated}, nil
}

func NewForumService(log *customtypes.Logger, db *db.Repository) ForumService {
	return &forumServiceImpl{log: log, db: db}
}
