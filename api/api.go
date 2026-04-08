package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"orders/domain"
	"orders/repo"
)

/*
	тут лучше сразу готовый роутер взять chi.NewRouter()
	будем намного удобнее и безопаснее
	+ вся логика должна быть в слое сервиса, а не в хендлере, сейчас у тебя 2в1
*/

const queryTimeout = time.Second

func MainHandler(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage) {
	switch r.Method {
	case "GET":
		getAllOrders(w, r, repo)
	case "POST":
		createOrder(w, r, repo)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func MainHandlerID(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage) {
	switch r.Method {
	case "GET":
		getOrderByID(w, r, repo)
	case "PUT":
		putOrderStatus(w, r, repo)
	case "DELETE":
		deleteOrder(w, r, repo)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage) {
	defer func() { _ = r.Body.Close() }()
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var ord domain.Order
	if err = json.Unmarshal(httpRequestBody, &ord); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	ordId, err := repo.CreateOrder(ctx, &ord)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mapa := make(map[string]int)
	mapa["id"] = ordId
	responseBody, err := json.Marshal(mapa)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getOrderByID(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	ord, err := repo.GetOrderByID(ctx, id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(ord, "", "	")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func getAllOrders(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage) {
	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	ords, err := repo.GetAllOrders(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseBody, err := json.MarshalIndent(ords, "", "	")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(responseBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func putOrderStatus(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage) {
	defer func() { _ = r.Body.Close() }()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var newStatusMap map[string]string
	err = json.Unmarshal(requestBody, &newStatusMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newStatus, existed := newStatusMap["Status"]
	if !existed {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	changedOrder, err := repo.UpdateOrderStatus(ctx, id, newStatus)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := json.MarshalIndent(*changedOrder, "", "	")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func deleteOrder(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	if _, err = repo.DeleteOrder(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
