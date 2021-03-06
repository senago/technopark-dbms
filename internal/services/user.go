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

type UserService interface {
	CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.Response, error)
	GetUserProfile(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.Response, error)
	UpdateUserProfile(ctx context.Context, request *dto.UpdateUserProfileRequest) (*dto.Response, error)
}

type userServiceImpl struct {
	log *customtypes.Logger
	db  *db.Repository
}

func (svc *userServiceImpl) CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.Response, error) {
	if users, err := svc.db.UserRepository.GetUsersByEmailOrNickname(ctx, request.Email, request.Nickname); err != nil {
		return nil, err
	} else if len(users) > 0 {
		return &dto.Response{Data: users, Code: http.StatusConflict}, nil
	}

	user := &core.User{Nickname: request.Nickname, Fullname: request.Fullname, About: request.About, Email: request.Email}
	if err := svc.db.UserRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &dto.Response{Data: user, Code: http.StatusCreated}, nil
}

func (svc *userServiceImpl) GetUserProfile(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.Response, error) {
	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.Nickname)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Nickname)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}
	return &dto.Response{Data: user, Code: http.StatusOK}, nil
}

func (svc *userServiceImpl) UpdateUserProfile(ctx context.Context, request *dto.UpdateUserProfileRequest) (*dto.Response, error) {
	if len(request.Email) > 0 {
		if user, err := svc.db.UserRepository.GetUserByEmail(ctx, request.Email); err != nil {
			if !errors.Is(err, constants.ErrDBNotFound) {
				return nil, err
			}
		} else if user.Nickname != request.Nickname {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("This email is already registered by user: %s", user.Nickname)}, Code: http.StatusConflict}, nil
		}
	}

	user := &core.User{Nickname: request.Nickname, Fullname: request.Fullname, About: request.About, Email: request.Email}
	updatedUser, err := svc.db.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Nickname)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}
	return &dto.Response{Data: updatedUser, Code: http.StatusOK}, nil
}

func NewUserService(log *customtypes.Logger, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
