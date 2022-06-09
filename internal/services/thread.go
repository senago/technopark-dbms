//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/senago/technopark-dbms/internal/constants"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
	"github.com/senago/technopark-dbms/internal/model/core"
	"github.com/senago/technopark-dbms/internal/model/dto"
)

type ForumThreadService interface {
	CreateForumThread(ctx context.Context, request *dto.CreateForumThreadRequest) (*dto.Response, error)
	UpdateVote(ctx context.Context, slugOrID string, request *dto.UpdateVoteRequest) (*dto.Response, error)
	GetThreadDetails(ctx context.Context, slugOrID string) (*dto.Response, error)
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

func (svc *forumThreadServiceImpl) UpdateVote(ctx context.Context, slugOrID string, request *dto.UpdateVoteRequest) (*dto.Response, error) {
	var id int
	var err error
	id, err = strconv.Atoi(slugOrID)

	var thread *core.Thread
	if err != nil {
		if thread, err = svc.db.ForumThreadRepository.GetForumThreadBySlug(ctx, slugOrID); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
		} else {
			id = int(thread.ID)
		}
	} else {
		if thread, err = svc.db.ForumThreadRepository.GetForumThreadByID(ctx, int64(id)); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
			}
		}
	}

	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.Nickname)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Nickname)}, Code: http.StatusNotFound}, nil
		}
	}
	request.Nickname = user.Nickname

	exists, err := svc.db.VotesRepository.VoteExists(ctx, request.Nickname, thread.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		if ok, err := svc.db.VotesRepository.UpdateVote(ctx, thread.ID, request.Nickname, request.Voice); err != nil {
			return nil, err
		} else if ok {
			thread.Votes += request.Voice * 2
		}
	} else {
		newVote := &core.Vote{
			Nickname: request.Nickname,
			ThreadID: thread.ID,
			Voice:    request.Voice,
		}

		if err := svc.db.VotesRepository.CreateVote(ctx, newVote); err != nil {
			return nil, err
		}

		thread.Votes += request.Voice
	}

	return &dto.Response{Data: thread, Code: http.StatusOK}, nil
}

func (svc *forumThreadServiceImpl) GetThreadDetails(ctx context.Context, slugOrID string) (*dto.Response, error) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		if thread, err := svc.db.ForumThreadRepository.GetForumThreadBySlug(ctx, slugOrID); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
			return nil, err
		} else {
			return &dto.Response{Data: thread, Code: http.StatusOK}, nil
		}
	}

	thread, err := svc.db.ForumThreadRepository.GetForumThreadByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
		}
	}

	return &dto.Response{Data: thread, Code: http.StatusOK}, nil
}

func NewForumThreadService(log *customtypes.Logger, db *db.Repository) ForumThreadService {
	return &forumThreadServiceImpl{log: log, db: db}
}
