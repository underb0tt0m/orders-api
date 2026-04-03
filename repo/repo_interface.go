package repo

import "orders/order"

type Repo interface {
	CreateOrder(o *order.Order) (int, error)
	GetOrderByID(id int) (order.Order, error)
	GetAllOrders() (map[int]*order.Order, error)
	UpdateOrderStatus(id int, newStatus string) (*order.Order, error)
	DeleteOrder(id int) (order.Order, error)
}
