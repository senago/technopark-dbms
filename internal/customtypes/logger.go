package customtypes

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Logger = zap.SugaredLogger
type DBConn = pgxpool.Pool
