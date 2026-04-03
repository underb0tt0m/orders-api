package db_ops

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTables(ctx context.Context, conn *pgx.Conn) error {
	stmt := `
CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50),
	count INTEGER,
	status VARCHAR(10)
);`
	_, err := conn.Exec(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}
