//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"

	"github.com/senago/technopark-dbms/internal/constants"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
	"github.com/senago/technopark-dbms/internal/model/core"
	"github.com/senago/technopark-dbms/internal/model/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, request *dto.CreateUserRequest) (interface{}, error)
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

func NewUserService(log *customtypes.Logger, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
