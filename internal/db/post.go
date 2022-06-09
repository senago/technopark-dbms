package db

import (
	"context"
	"fmt"
	"time"

	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/core"
	"github.com/senago/technopark-dbms/internal/model/dto"
)

const (
	queryCheckPostParent = "SELECT threads FROM post WHERE id = $1;"
)

type PostsRepository interface {
	CreatePosts(ctx context.Context, forum string, thread int64, posts []*dto.PostData) ([]*core.Post, error)
	CheckParentPost(ctx context.Context, parent int) (int, error)
}

type postsRepositoryImpl struct {
	dbConn *customtypes.DBConn
}

func (repo *postsRepositoryImpl) CreatePosts(ctx context.Context, forum string, thread int64, posts []*dto.PostData) ([]*core.Post, error) {
	query := "INSERT INTO posts (parent, author, message, forum, thread, created) VALUES "

	queryArgs := make([]interface{}, 0, 0)
	newPosts := make([]*core.Post, 0, len(posts))
	insertTime := time.Now()
	for i, post := range posts {
		p := &core.Post{Parent: post.Parent, Author: post.Author, Message: post.Message, Forum: forum, Thread: thread, Created: insertTime}
		newPosts = append(newPosts, p)

		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)

		queryArgs = append(queryArgs, post.Parent, post.Author, post.Message, forum, thread, insertTime)
	}

	query = query[:len(query)-1] // get rid of last comma
	query += " RETURNING id;"

	rows, err := repo.dbConn.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		if err = rows.Scan(&newPosts[i].ID); err != nil {
			return nil, err
		}
	}

	return newPosts, nil
}

func (repo *postsRepositoryImpl) CheckParentPost(ctx context.Context, parent int) (int, error) {
	var threadID int
	err := repo.dbConn.QueryRow(ctx, queryCheckPostParent, parent).Scan(&threadID)
	return threadID, err
}

func NewPostsRepository(dbConn *customtypes.DBConn) (*postsRepositoryImpl, error) {
	return &postsRepositoryImpl{dbConn: dbConn}, nil
}
