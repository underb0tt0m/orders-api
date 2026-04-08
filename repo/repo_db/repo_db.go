package repo_db

import (
	"context"

	"orders/domain"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{
		conn: conn,
	}
}

func (r *Repo) CreateOrder(ctx context.Context, o *domain.Order) (int, error) {
	stmt := `
INSERT INTO orders(name, count, status)
VALUES ($1, $2, $3)
RETURNING id;
`
	var id int
	if err := r.conn.QueryRow(
		ctx,
		stmt,
		o.Name,
		o.Count,
		o.Status,
	).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r *Repo) GetOrderByID(ctx context.Context, id int) (domain.Order, error) {
	stmt := `
SELECT name, count, status
FROM orders
WHERE id = $1;
`
	var (
		name, status string
		count        int
	)

	if err := r.conn.QueryRow(
		ctx,
		stmt,
		id,
	).Scan(
		&name,
		&count,
		&status,
	); err != nil {
		return domain.Order{}, err
	}

	return domain.Order{
		Name:   name,
		Count:  count,
		Status: status,
	}, nil
}

func (r *Repo) GetAllOrders(ctx context.Context) (map[int]*domain.Order, error) {
	stmt := `
SELECT id, name, count, status
FROM orders;
`
	orders := make(map[int]*domain.Order)

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id, count    int
		name, status string
	)

	for rows.Next() {
		if serr := rows.Scan(&id, &name, &count, &status); serr != nil {
			return nil, serr
		}

		orders[id] = &domain.Order{
			Name:   name,
			Count:  count,
			Status: status,
		}
	}

	return orders, nil

}

func (r *Repo) UpdateOrderStatus(ctx context.Context, id int, newStatus string) (*domain.Order, error) {
	stmt := `
UPDATE orders
SET status = $1
WHERE id = $2
RETURNING name, count, status;
`
	var (
		name, status string
		count        int
	)

	if err := r.conn.QueryRow(ctx, stmt, newStatus, id).Scan(&name, &count, &status); err != nil {
		return nil, err
	}

	return &domain.Order{
		Name:   name,
		Count:  count,
		Status: status,
	}, nil
}

func (r *Repo) DeleteOrder(ctx context.Context, id int) (domain.Order, error) {
	stmt := `
DELETE FROM orders
WHERE id = $1
RETURNING name, count, status;
`
	var (
		name, status string
		count        int
	)

	if err := r.conn.QueryRow(ctx, stmt, id).Scan(&name, &count, &status); err != nil {
		return domain.Order{}, err
	}

	return domain.Order{
		Name:   name,
		Count:  count,
		Status: status,
	}, nil
}
