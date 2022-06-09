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

type UserService interface {
	CreateUser(ctx context.Context, request *dto.CreateUserRequest) (interface{}, error)
	GetUserProfile(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.GetUserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, request *dto.UpdateUserProfileRequest) (interface{}, error)
}

type userServiceImpl struct {
	log *customtypes.Logger
	db  *db.Repository
}

func (svc *userServiceImpl) CreateUser(ctx context.Context, request *dto.CreateUserRequest) (interface{}, error) {
	user := &core.User{Nickname: request.Nickname, Fullname: request.Fullname, About: request.About, Email: request.Email}

	if users, err := svc.db.UserRepository.GetUsersByEmailOrNickname(ctx, user.Email, user.Nickname); err != nil {
		return nil, err
	} else if len(users) > 0 {
		return users, constants.ErrUserAlreadyExists
	}

	if err := svc.db.UserRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *userServiceImpl) GetUserProfile(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.GetUserProfileResponse, error) {
	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.Nickname)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return nil, constants.NewCodedError(fmt.Sprintf("Can't find user by nickname: %s", request.Nickname), http.StatusNotFound)
		}
		return nil, err
	}
	return user, nil
}

func (svc *userServiceImpl) UpdateUserProfile(ctx context.Context, request *dto.UpdateUserProfileRequest) (interface{}, error) {
	user := &core.User{Nickname: request.Nickname, Fullname: request.Fullname, About: request.About, Email: request.Email}
	updatedUser, err := svc.db.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return nil, constants.NewCodedError(fmt.Sprintf("Can't find user by nickname: %s", request.Nickname), http.StatusNotFound)
		} else {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return nil, constants.NewCodedError(fmt.Sprintf("This email is already registered by user: %s", user.Nickname), http.StatusConflict)
			}
		}
		return nil, err
	}
	return updatedUser, nil
}

func NewUserService(log *customtypes.Logger, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
