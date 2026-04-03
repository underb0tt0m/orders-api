package db_ops

import (
	"context"
	"orders/order"

	"github.com/jackc/pgx/v5"
)

func InsertOrder(ctx context.Context, conn *pgx.Conn, o *order.Order) (int, error) {
	stmt := `
INSERT INTO orders(name, count, status)
VALUES ($1, $2, $3)
RETURNING id;
`
	var id int
	err := conn.QueryRow(ctx, stmt, o.Name, o.Count, o.Status).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func SelectOrderByID(ctx context.Context, conn *pgx.Conn, id int) (order.Order, error) {
	stmt := `
SELECT name, count, status
FROM orders
WHERE id = $1;
`
	var name string
	var count int
	var status string
	err := conn.QueryRow(ctx, stmt, id).Scan(&name, &count, &status)
	return order.Order{
		Name:   name,
		Count:  count,
		Status: status,
	}, err
}

func GetAllOrders(ctx context.Context, conn *pgx.Conn) (map[int]*order.Order, error) {
	stmt := `
SELECT id, name, count, status
FROM orders;
`

	rows, err := conn.Query(ctx, stmt)
	orders := make(map[int]*order.Order)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	var id int
	var name string
	var count int
	var status string
	for rows.Next() {
		err = rows.Scan(&id, &name, &count, &status)
		if err != nil {
			return orders, err
		}
		orders[id] = order.NewOrder(name, count, status)
	}
	return orders, nil
}

func UpdateOrderStatus(ctx context.Context, conn *pgx.Conn, newStatus string, id int) (order.Order, error) {
	stmt := `
UPDATE orders
SET status = $1
WHERE id = $2
RETURNING name, count, status;
`
	var name string
	var count int
	var status string
	err := conn.QueryRow(ctx, stmt, newStatus, id).Scan(&name, &count, &status)
	return order.Order{
		Name:   name,
		Count:  count,
		Status: status,
	}, err
}

func DeleteOrderByID(ctx context.Context, conn *pgx.Conn, id int) (order.Order, error) {
	stmt := `
DELETE FROM orders
WHERE id = $1
RETURNING name, count, status;
`
	var name string
	var count int
	var status string
	err := conn.QueryRow(ctx, stmt, id).Scan(&name, &count, &status)
	return order.Order{
		Name:   name,
		Count:  count,
		Status: status,
	}, err
}
