package db

import (
	"context"

	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/core"
)

const (
	queryDeleteAllTables           = "TRUNCATE TABLE users, forums, threads, posts, forum_users, votes CASCADE;"
	queryCountForumPostThreadUsers = "SELECT (SELECT count(*) FROM users) AS user, (SELECT count(*) FROM forums) AS forum, (SELECT count(*) FROM threads) AS thread, (SELECT count(*) FROM posts) AS post;"
)

type ServiceRepository interface {
	Status(ctx context.Context) (*core.ServiceInfo, error)
	Delete(ctx context.Context) error
}

type serviceRepositoryImpl struct {
	dbConn *customtypes.DBConn
}

func (repo *serviceRepositoryImpl) Status(ctx context.Context) (*core.ServiceInfo, error) {
	res := &core.ServiceInfo{}
	err := repo.dbConn.QueryRow(ctx, queryCountForumPostThreadUsers).Scan(&res.User, &res.Forum, &res.Thread, &res.Post)
	return res, err
}

func (repo *serviceRepositoryImpl) Delete(ctx context.Context) error {
	_, err := repo.dbConn.Exec(ctx, queryDeleteAllTables)
	return err
}

func NewServiceRepository(dbConn *customtypes.DBConn) *serviceRepositoryImpl {
	return &serviceRepositoryImpl{dbConn: dbConn}
}
