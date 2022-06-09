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

type ForumThreadService interface {
	CreateForumThread(ctx context.Context, request *dto.CreateForumThreadRequest) (*dto.Response, error)
}

type forumThreadServiceImpl struct {
	log *customtypes.Logger
	db  *db.Repository
}

func (svc *forumThreadServiceImpl) CreateForumThread(ctx context.Context, request *dto.CreateForumThreadRequest) (*dto.Response, error) {
	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.Author)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Author)}, Code: http.StatusNotFound}, nil
		}
	}
	request.Author = user.Nickname

	if forum, err := svc.db.ForumRepository.GetForumBySlug(ctx, request.Forum); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", request.Forum)}, Code: http.StatusNotFound}, nil
		}
	} else {
		request.Forum = forum.Slug
	}

	if request.Slug != "" {
		if thread, err := svc.db.ForumThreadRepository.GetForumThreadBySlug(ctx, request.Slug); err != nil {
			if !errors.Is(err, constants.ErrDBNotFound) {
				return nil, err
			}
		} else {
			return &dto.Response{Data: thread, Code: http.StatusConflict}, nil
		}
	}

	reqThread := &core.Thread{Forum: request.Forum, Title: request.Title, Author: request.Author, Message: request.Message, Slug: request.Slug, Created: request.Created}
	thread, err := svc.db.ForumThreadRepository.CreateForumThread(ctx, reqThread)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: thread, Code: http.StatusCreated}, nil
}

func NewForumThreadService(log *customtypes.Logger, db *db.Repository) ForumThreadService {
	return &forumThreadServiceImpl{log: log, db: db}
}
