package repo_db

import (
	"context"
	"orders/db/db_ops"
	"orders/order"
	"sync"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	mtx  sync.RWMutex
	conn *pgx.Conn
	ctx  context.Context
}

func NewRepo(conn *pgx.Conn, ctx context.Context) *Repo {
	return &Repo{
		sync.RWMutex{},
		conn,
		ctx}
}

func (r *Repo) CreateOrder(o *order.Order) (int, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	id, err := db_ops.InsertOrder(r.ctx, r.conn, o)
	return id, err
}

func (r *Repo) GetOrderByID(id int) (order.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	ord, err := db_ops.SelectOrderByID(r.ctx, r.conn, id)
	if err != nil {
		return order.Order{}, err
	}
	return ord, err
}

func (r *Repo) GetAllOrders() (map[int]*order.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	ords, err := db_ops.GetAllOrders(r.ctx, r.conn)
	return ords, err
}

func (r *Repo) UpdateOrderStatus(id int, newStatus string) (*order.Order, error) {
	var err error
	r.mtx.Lock()
	ord, err := db_ops.UpdateOrderStatus(r.ctx, r.conn, newStatus, id)
	r.mtx.Unlock()
	return &ord, err
}

func (r *Repo) DeleteOrder(id int) (order.Order, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	ord, err := db_ops.DeleteOrderByID(r.ctx, r.conn, id)
	return ord, err
}
