//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"

	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/db"
	"github.com/senago/technopark-dbms/internal/model/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
}

type userServiceImpl struct {
	log *customtypes.Logger
	db  *db.Repository
}

func (svc *userServiceImpl) CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	return &dto.CreateUserResponse{Nickname: request.Nickname, Fullname: request.Fullname, About: request.About, Email: request.Email}, nil
}

func NewUserService(log *customtypes.Logger, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
