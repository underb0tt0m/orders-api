package repo_inmemory

import (
	"context"
	"fmt"
	"orders/domain"
	"sync"
)

type Repo struct {
	i       int
	Storage map[int]*domain.Order
	mu      sync.RWMutex
}

func New() *Repo {
	return &Repo{
		i:       0,
		Storage: make(map[int]*domain.Order),
		mu:      sync.RWMutex{},
	}
}

func (r *Repo) CreateOrder(ctx context.Context, o *domain.Order) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Storage[r.i+1] = o
	r.i++
	return r.i, nil
}

func (r *Repo) GetOrderByID(ctx context.Context, id int) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ord, exists := r.Storage[id]

	if !exists {
		return domain.Order{}, fmt.Errorf("заказа c id=%v не существует", id)
	}
	return *ord, nil
}

func (r *Repo) GetAllOrders(ctx context.Context) (map[int]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := r.Storage
	resultCopy := make(map[int]*domain.Order)
	for k, v := range result {
		resultCopy[k] = v
	}
	return resultCopy, nil
}

func (r *Repo) UpdateOrderStatus(ctx context.Context, id int, newStatus string) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ord, exists := r.Storage[id]
	if !exists {
		return &domain.Order{}, fmt.Errorf("заказа c id=%v не существует", id)
	}
	ord.Status = newStatus
	return ord, nil
}

func (r *Repo) DeleteOrder(ctx context.Context, id int) (domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ord, exists := r.Storage[id]

	if !exists {
		return domain.Order{}, fmt.Errorf("заказа c id=%v не существует", id)
	}
	deleted := *ord
	delete(r.Storage, id)
	return deleted, nil
}
