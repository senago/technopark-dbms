package customtypes

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Logger = zap.SugaredLogger
type DBConn = pgxpool.Pool
