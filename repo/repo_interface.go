package repo

import (
	"context"
	"orders/domain"

	"go.uber.org/zap"
)

type OrderStorage interface {
	CreateOrder(ctx context.Context, o *domain.Order, logger *zap.Logger) (int, error)
	GetOrderByID(ctx context.Context, id int, logger *zap.Logger) (domain.Order, error)
	GetAllOrders(ctx context.Context, logger *zap.Logger) (map[int]*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id int, newStatus string, logger *zap.Logger) (*domain.Order, error)
	DeleteOrder(ctx context.Context, id int, logger *zap.Logger) (domain.Order, error)
}
