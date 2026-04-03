package simple_db_conn

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func GetDBConn(ctx context.Context) (*pgx.Conn, error) {
	// postgresql://user:password@localhost:5432/dbname
	connString := "postgresql://postgres:postgres@localhost:5432/postgres"
	return pgx.Connect(ctx, connString)
}
