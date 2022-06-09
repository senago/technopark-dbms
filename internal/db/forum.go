package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/core"
)

const (
	queryCreateForum = `INSERT INTO forums (title, "user", slug) VALUES ($1, $2, $3);`

	queryGetForumBySlug  = `SELECT title, "user", slug, posts, threads FROM forums WHERE slug = $1;`
	queryGetForumThreads = "SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM threads as t LEFT JOIN forum f on t.forum = f.slug WHERE f.slug = $1;"
)

type ForumRepository interface {
	CreateForum(ctx context.Context, forum *core.Forum) error

	GetForumBySlug(ctx context.Context, slug string) (*core.Forum, error)
	GetForumUsers(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.User, error)
	GetForumThreads(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.Thread, error)
}

type forumRepositoryImpl struct {
	dbConn *customtypes.DBConn
}

func (repo *forumRepositoryImpl) CreateForum(ctx context.Context, forum *core.Forum) error {
	_, err := repo.dbConn.Exec(ctx, queryCreateForum, &forum.Title, &forum.User, &forum.Slug)
	return err
}

func (repo *forumRepositoryImpl) GetForumBySlug(ctx context.Context, slug string) (*core.Forum, error) {
	forum := &core.Forum{}
	err := repo.dbConn.QueryRow(ctx, queryGetForumBySlug, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	return forum, wrapErr(err)
}

func (repo *forumRepositoryImpl) GetForumUsers(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.User, error) {
	query := constructGetForumUsersQuery(limit, since, desc)
	rows, err := repo.dbConn.Query(ctx, query, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &core.User{}
	users := []*core.User{}
	for rows.Next() {
		if err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *forumRepositoryImpl) GetForumThreads(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.Thread, error) {
	var rows pgx.Rows
	var err error

	query := queryGetForumThreads

	queryOrderBy := "ORDER BY t.created "
	if desc {
		queryOrderBy += "DESC"
	}
	if limit > 0 {
		queryOrderBy += fmt.Sprintf(" LIMIT %d", limit)
	}

	if since != "" {
		querySince := " AND t.created >= $2 "
		if since != "" && desc {
			querySince = " and t.created <= $2 "
		} else if since != "" && !desc {
			querySince = " and t.created >= $2 "
		}

		query = query + querySince + queryOrderBy
		rows, err = repo.dbConn.Query(ctx, query, slug, since)
	} else {
		query = query + queryOrderBy
		rows, err = repo.dbConn.Query(ctx, query, slug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	thread := &core.Thread{}
	threads := []*core.Thread{}
	for rows.Next() {
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}

func NewForumRepository(dbConn *customtypes.DBConn) (*forumRepositoryImpl, error) {
	return &forumRepositoryImpl{dbConn: dbConn}, nil
}

func constructGetForumUsersQuery(limit int64, since string, desc bool) string {
	query := "SELECT u.nickname, u.fullname, u.about, u.email from forum_users u where u.forum = $1"

	if desc && since != "" {
		query += fmt.Sprintf(" and u.nickname < '%s'", since)
	} else if since != "" {
		query += fmt.Sprintf(" and u.nickname > '%s'", since)
	}

	query += " ORDER BY u.nickname"
	if desc {
		query += " desc"
	}
	query += fmt.Sprintf(" LIMIT %d", limit)

	return query
}