package repo_inmemory

import (
	"errors"
	"fmt"
	"orders/order"
	"sync"
)

type Repo struct {
	i       int
	Storage map[int]*order.Order
	mtx     sync.RWMutex
}

func NewRepo() *Repo {
	return &Repo{
		0,
		make(map[int]*order.Order),
		sync.RWMutex{},
	}
}

func (r *Repo) CreateOrder(o *order.Order) (int, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.Storage[r.i+1] = o
	r.i++

	return r.i, nil
}

func (r *Repo) GetOrderByID(id int) (order.Order, error) {
	var err error
	r.mtx.RLock()
	ord, exists := r.Storage[id]
	r.mtx.RUnlock()
	if !exists {
		err = errors.New(fmt.Sprintf("Заказа c id=%v не существует", id))
		return order.Order{}, err
	}
	return *ord, err
}

func (r *Repo) GetAllOrders() (map[int]*order.Order, error) {
	r.mtx.RLock()
	result := r.Storage
	resultCopy := make(map[int]*order.Order)
	for k, v := range result {
		resultCopy[k] = v
	}
	r.mtx.RUnlock()
	return resultCopy, nil
}

func (r *Repo) UpdateOrderStatus(id int, newStatus string) (*order.Order, error) {
	var err error
	r.mtx.Lock()
	ord, exists := r.Storage[id]
	if !exists {
		err = errors.New(fmt.Sprintf("Заказа c id=%v не существует", id))
	} else {
		ord.Status = newStatus
	}
	r.mtx.Unlock()
	return ord, err
}

func (r *Repo) DeleteOrder(id int) (order.Order, error) {
	var err error
	var deleted order.Order
	r.mtx.Lock()
	ord, exists := r.Storage[id]

	if !exists {
		err = errors.New(fmt.Sprintf("Заказа c id=%v не существует", id))
	} else {
		deleted = *ord
		delete(r.Storage, id)
	}
	r.mtx.Unlock()
	return deleted, err
}
