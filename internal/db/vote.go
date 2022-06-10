package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/senago/technopark-dbms/internal/customtypes"
	"github.com/senago/technopark-dbms/internal/model/core"
)

const (
	queryCreateVote = "INSERT INTO votes (nickname, thread, voice) VALUES ($1, $2, $3);"
	queryVoteExists = "SELECT voice from votes where nickname = $1 and thread = $2;"
	queryUpdateVote = "UPDATE votes SET voice = $3 WHERE thread = $1 and nickname = $2 and voice != $3;"
)

type VotesRepository interface {
	CreateVote(ctx context.Context, vote *core.Vote) error
	VoteExists(ctx context.Context, nickname string, threadID int64) (bool, error)
	UpdateVote(ctx context.Context, threadID int64, nickname string, voice int64) (bool, error)
}

type votesRepositoryImpl struct {
	dbConn *customtypes.DBConn
}

func (repo *votesRepositoryImpl) CreateVote(ctx context.Context, vote *core.Vote) error {
	_, err := repo.dbConn.Exec(ctx, queryCreateVote, vote.Nickname, vote.ThreadID, vote.Voice)
	return wrapErr(err)
}

func (repo *votesRepositoryImpl) VoteExists(ctx context.Context, nickname string, threadID int64) (bool, error) {
	voice := 0
	err := repo.dbConn.QueryRow(ctx, queryVoteExists, nickname, threadID).Scan(&voice)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *votesRepositoryImpl) UpdateVote(ctx context.Context, threadID int64, nickname string, voice int64) (bool, error) {
	res, err := repo.dbConn.Exec(ctx, queryUpdateVote, threadID, nickname, voice)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() == 1, nil
}

func NewVotesRepository(dbConn *customtypes.DBConn) *votesRepositoryImpl {
	return &votesRepositoryImpl{dbConn: dbConn}
}
