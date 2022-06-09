package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/core"
	"github.com/senago/technopark-dbms/internal/model/dto"
)

const (
	queryCheckPostParent = "SELECT thread FROM posts WHERE id = $1;"
)

type PostsRepository interface {
	CreatePosts(ctx context.Context, forum string, thread int64, posts []*dto.PostData) ([]*core.Post, error)
	CheckParentPost(ctx context.Context, parent int) (int, error)

	GetPostsFlat(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error)
	GetPostsTree(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error)
	GetPostsParentTree(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error)
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

func (repo *postsRepositoryImpl) GetPostsFlat(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error) {
	query := "SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread = $1 "

	if since != -1 {
		if desc {
			query += "AND id < $2 "
		} else {
			query += "AND id > $2 "
		}
	}

	if desc {
		query += "ORDER BY created DESC, id DESC "
	} else {
		query += "ORDER BY created ASC, id ASC "
	}

	query += fmt.Sprintf("LIMIT %d ", limit)

	var rows pgx.Rows
	var err error
	if since == -1 {
		rows, err = repo.dbConn.Query(ctx, query, id)
	} else {
		rows, err = repo.dbConn.Query(ctx, query, id, since)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*core.Post{}
	for rows.Next() {
		post := &core.Post{}
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, err
}

func (repo *postsRepositoryImpl) GetPostsTree(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error) {
	query := "SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread = $1 "

	if since != -1 {
		if desc {
			query += "and path < "
		} else {
			query += "and path > "
		}
		query += fmt.Sprintf("(SELECT path FROM posts WHERE id = %d) ", since)
	}

	if desc {
		query += "ORDER BY path desc "
	} else {
		query += "ORDER BY path asc, id "
	}

	query += fmt.Sprintf("LIMIT NULLIF(%d, 0) ", limit)

	rows, err := repo.dbConn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*core.Post{}
	for rows.Next() {
		post := &core.Post{}
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repo *postsRepositoryImpl) GetPostsParentTree(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error) {
	var rows pgx.Rows
	var err error
	if since == -1 {
		if desc {
			rows, err = repo.dbConn.Query(ctx,
				` SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path ASC, id ASC;`,
				id, limit)
		} else {
			rows, err = repo.dbConn.Query(ctx,
				`	SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 ORDER BY id ASC LIMIT $2)
					ORDER BY path ASC, id ASC;`,
				id, limit)
		}
	} else {
		if desc {
			rows, err = repo.dbConn.Query(ctx,
				` SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 AND path[1] < (SELECT path[1] FROM posts WHERE id = $2)
					ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path ASC, id ASC;`,
				id, since, limit)
		} else {
			rows, err = repo.dbConn.Query(ctx,
				` SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 AND path[1] >
					(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY path ASC, id ASC;`,
				id, since, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*core.Post{}
	for rows.Next() {
		post := &core.Post{}
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func NewPostsRepository(dbConn *customtypes.DBConn) (*postsRepositoryImpl, error) {
	return &postsRepositoryImpl{dbConn: dbConn}, nil
}
