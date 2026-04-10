package simple_db_conn

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetDBConn(ctx context.Context) (*pgx.Conn, error) {
	// postgresql://user:password@localhost:5432/dbname
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgresql://postgres:postgres@localhost:5432/postgres"
	}
	return pgx.Connect(ctx, connString)
}
