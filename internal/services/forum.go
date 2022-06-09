//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/senago/technopark-dbms/internal/constants"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
	"github.com/senago/technopark-dbms/internal/model/core"
	"github.com/senago/technopark-dbms/internal/model/dto"
)

type ForumService interface {
	CreateForum(ctx context.Context, request *dto.CreateForumRequest) (*dto.Response, error)
	GetForumBySlug(ctx context.Context, request *dto.GetForumBySlugRequest) (*dto.Response, error)
	GetForumThreads(ctx context.Context, request *dto.GetForumThreadsRequest) (*dto.Response, error)
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

	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.User)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.User)}, Code: http.StatusNotFound}, nil
		}
	}
	request.User = user.Nickname

	if err := svc.db.ForumRepository.CreateForum(ctx, &core.Forum{Title: request.Title, User: request.User, Slug: request.Slug}); err != nil {
		return nil, err
	}

	forum, err := svc.db.ForumRepository.GetForumBySlug(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: forum, Code: http.StatusCreated}, nil
}

func (svc *forumServiceImpl) GetForumBySlug(ctx context.Context, request *dto.GetForumBySlugRequest) (*dto.Response, error) {
	forum, err := svc.db.ForumRepository.GetForumBySlug(ctx, request.Slug)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find forum with slug: %s", request.Slug)}, Code: http.StatusNotFound}, nil
		}
	}
	return &dto.Response{Data: forum, Code: http.StatusOK}, nil
}

func (svc *forumServiceImpl) GetForumThreads(ctx context.Context, request *dto.GetForumThreadsRequest) (*dto.Response, error) {
	if forum, err := svc.db.ForumRepository.GetForumBySlug(ctx, request.Slug); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find forum with slug: %s", request.Slug)}, Code: http.StatusNotFound}, nil
		}
	} else {
		request.Slug = forum.Slug
	}

	threads, err := svc.db.ForumRepository.GetForumThreads(ctx, request.Slug, request.Limit, request.Since, request.Desc)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: threads, Code: http.StatusOK}, nil
}

func NewForumService(log *customtypes.Logger, db *db.Repository) ForumService {
	return &forumServiceImpl{log: log, db: db}
}
