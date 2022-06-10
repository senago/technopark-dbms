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

type PostsService interface {
	CreatePosts(ctx context.Context, slugOrID string, posts []*dto.PostData) (*dto.Response, error)

	GetPosts(ctx context.Context, slugOrID string, sort string, since int64, desc bool, limit int64) (*dto.Response, error)
	GetPostDetails(ctx context.Context, request *dto.GetPostDetailsRequest) (*dto.Response, error)

	UpdatePost(ctx context.Context, request *dto.UpdatePostRequest) (*dto.Response, error)
}

type postsServiceImpl struct {
	log *customtypes.Logger
	db  *db.Repository
}

func (svc *postsServiceImpl) CreatePosts(ctx context.Context, slugOrID string, posts []*dto.PostData) (*dto.Response, error) {
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

	if len(posts) == 0 {
		return &dto.Response{Data: []struct{}{}, Code: http.StatusCreated}, nil
	}

	if posts[0].Parent != 0 {
		parentThreadID, err := svc.db.PostsRepository.CheckParentPost(ctx, int(posts[0].Parent))
		if err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.Response{Data: dto.ErrorResponse{Message: "Parent post was created in another thread"}, Code: http.StatusConflict}, nil
			}
		}

		if parentThreadID != id {
			return &dto.Response{Data: dto.ErrorResponse{Message: "Parent post was created in another thread"}, Code: http.StatusConflict}, nil
		}
	}

	if _, err := svc.db.UserRepository.GetUserByNickname(ctx, posts[0].Author); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", posts[0].Author)}, Code: http.StatusNotFound}, nil
		}
	}

	insertedPosts, err := svc.db.PostsRepository.CreatePosts(ctx, thread.Forum, int64(id), posts)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: insertedPosts, Code: http.StatusCreated}, nil
}

func (svc *postsServiceImpl) GetPosts(ctx context.Context, slugOrID string, sort string, since int64, desc bool, limit int64) (*dto.Response, error) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		if thread, err := svc.db.ForumThreadRepository.GetForumThreadBySlug(ctx, slugOrID); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
		} else {
			id = int(thread.ID)
		}
	}

	if _, err := svc.db.ForumThreadRepository.GetForumThreadByID(ctx, int64(id)); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
		}
	}

	var posts []*core.Post
	switch sort {
	case "flat":
		posts, err = svc.db.PostsRepository.GetPostsFlat(ctx, id, since, desc, limit)
	case "tree":
		posts, err = svc.db.PostsRepository.GetPostsTree(ctx, id, since, desc, limit)
	case "parent_tree":
		posts, err = svc.db.PostsRepository.GetPostsParentTree(ctx, id, since, desc, limit)
	default:
		posts, err = svc.db.PostsRepository.GetPostsFlat(ctx, id, since, desc, limit)
	}
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: posts, Code: http.StatusOK}, nil
}

func (svc *postsServiceImpl) GetPostDetails(ctx context.Context, request *dto.GetPostDetailsRequest) (*dto.Response, error) {
	post, err := svc.db.PostsRepository.GetPostByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find post by id: %d", request.ID)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}

	postDetails, err := svc.db.PostsRepository.GetPostDetails(ctx, request.ID, request.Related)
	if err != nil {
		return nil, err
	}
	postDetails.Post = post

	return &dto.Response{Data: postDetails, Code: http.StatusOK}, nil
}

func (svc *postsServiceImpl) UpdatePost(ctx context.Context, request *dto.UpdatePostRequest) (*dto.Response, error) {
	post, err := svc.db.PostsRepository.GetPostByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find post by id: %d", request.ID)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}

	if len(request.Message) == 0 || request.Message == post.Message {
		return &dto.Response{Data: post, Code: http.StatusOK}, nil
	}

	updatedPost, err := svc.db.PostsRepository.UpdatePost(ctx, request.ID, request.Message)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: updatedPost, Code: http.StatusOK}, nil
}

func NewPostsService(log *customtypes.Logger, db *db.Repository) PostsService {
	return &postsServiceImpl{log: log, db: db}
}
