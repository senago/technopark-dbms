package db

import (
	"context"

	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/core"
)

const (
	queryCreateForumThread = "INSERT INTO threads (title, author, forum, message, slug, created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, author, forum, message, votes, slug, created;"

	querGetForumThreadByID    = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1;"
	queryGetForumThreadBySlug = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE slug = $1;"

	queryUpdateForumThreadByID = "UPDATE threads SET title = $2, message = $3 WHERE id = $1 RETURNING id, title, author, forum, message, votes, slug, created;"
)

type ForumThreadRepository interface {
	CreateForumThread(ctx context.Context, thread *core.Thread) (*core.Thread, error)

	GetForumThreadByID(ctx context.Context, id int64) (*core.Thread, error)
	GetForumThreadBySlug(ctx context.Context, slug string) (*core.Thread, error)

	UpdateForumThreadByID(ctx context.Context, id int64, title string, message string) (*core.Thread, error)
}

type forumThreadRepositoryImpl struct {
	dbConn *customtypes.DBConn
}

func (repo *forumThreadRepositoryImpl) CreateForumThread(ctx context.Context, thread *core.Thread) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx, queryCreateForumThread, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).
		Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return t, err
}

func (repo *forumThreadRepositoryImpl) GetForumThreadByID(ctx context.Context, id int64) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx, querGetForumThreadByID, id).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return t, wrapErr(err)
}

func (repo *forumThreadRepositoryImpl) GetForumThreadBySlug(ctx context.Context, slug string) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx, queryGetForumThreadBySlug, slug).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return t, wrapErr(err)
}

func (repo *forumThreadRepositoryImpl) UpdateForumThreadByID(ctx context.Context, id int64, title string, message string) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx, queryUpdateForumThreadByID, id, title, message).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return t, wrapErr(err)
}

func NewForumThreadRepository(dbConn *customtypes.DBConn) *forumThreadRepositoryImpl {
	return &forumThreadRepositoryImpl{dbConn: dbConn}
}
