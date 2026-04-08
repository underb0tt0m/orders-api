package repo

import (
	"context"

	"orders/domain"
)

type OrderStorage interface {
	CreateOrder(ctx context.Context, o *domain.Order) (int, error)
	GetOrderByID(ctx context.Context, id int) (domain.Order, error)
	GetAllOrders(ctx context.Context) (map[int]*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id int, newStatus string) (*domain.Order, error)
	DeleteOrder(ctx context.Context, id int) (domain.Order, error)
}
