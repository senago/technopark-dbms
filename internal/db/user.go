package db

import (
	"context"

	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/core"
)

const (
	queryCreateUser = "INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4);"

	queryGetUserByEmail            = "SELECT nickname, fullname, about, email FROM users where email = $1;"
	queryGetUserByNickname         = "SELECT nickname, fullname, about, email FROM users where nickname = $1;"
	queryGetUsersByEmailOrNickname = "SELECT nickname, fullname, about, email FROM users WHERE email = $1 OR nickname = $2;"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *core.User) error

	GetUserByEmail(ctx context.Context, email string) (*core.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (*core.User, error)
	GetUsersByEmailOrNickname(ctx context.Context, email, nickname string) ([]*core.User, error)
}

type userRepositoryImpl struct {
	dbConn *customtypes.DBConn
}

func (repo *userRepositoryImpl) CreateUser(ctx context.Context, user *core.User) error {
	_, err := repo.dbConn.Exec(ctx, queryCreateUser, user.Nickname, user.Fullname, user.About, user.Email)
	return err
}

func (repo *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	user := &core.User{}
	err := repo.dbConn.QueryRow(ctx, queryGetUserByEmail, email).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return user, err
}

func (repo *userRepositoryImpl) GetUserByNickname(ctx context.Context, nickname string) (*core.User, error) {
	user := &core.User{}
	err := repo.dbConn.QueryRow(ctx, queryGetUserByNickname, nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return user, err
}

func (repo *userRepositoryImpl) GetUsersByEmailOrNickname(ctx context.Context, email, nickname string) ([]*core.User, error) {
	rows, err := repo.dbConn.Query(ctx, queryGetUsersByEmailOrNickname, email, nickname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*core.User{}
	for rows.Next() {
		user := &core.User{}
		if err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// NewUserRepository creates a new instance of userRepositoryImpl
func NewUserRepository(dbConn *customtypes.DBConn) (*userRepositoryImpl, error) {
	return &userRepositoryImpl{dbConn: dbConn}, nil
}
