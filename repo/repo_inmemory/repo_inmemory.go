package repo_inmemory

import (
	"fmt"
	"sync"

	"orders/domain"
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

func (r *Repo) CreateOrder(o *domain.Order) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Storage[r.i+1] = o
	r.i++

	return r.i, nil
}

func (r *Repo) GetOrderByID(id int) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ord, exists := r.Storage[id]
	if !exists {
		return domain.Order{}, fmt.Errorf("заказа c id=%v не существует", id)
	}

	return *ord, nil
}

func (r *Repo) GetAllOrders() (map[int]*domain.Order, error) {
	r.mu.RLock()
	result := r.Storage
	resultCopy := make(map[int]*domain.Order)
	for k, v := range result {
		resultCopy[k] = v
	}
	r.mu.RUnlock()
	return resultCopy, nil
}

func (r *Repo) UpdateOrderStatus(id int, newStatus string) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ord, exists := r.Storage[id]
	if !exists {
		return nil, fmt.Errorf("заказа c id=%v не существует", id)
	}
	r.Storage[id].Status = newStatus

	return ord, nil
}

func (r *Repo) DeleteOrder(id int) (domain.Order, error) {
	var deleted domain.Order
	r.mu.Lock()
	defer r.mu.Unlock()

	ord, exists := r.Storage[id]
	if !exists {
		return domain.Order{}, fmt.Errorf("заказа c id=%v не существует", id)
	}

	deleted = *ord
	delete(r.Storage, id)

	return deleted, nil
}
