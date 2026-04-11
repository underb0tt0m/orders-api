package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"orders/domain"
	"orders/repo"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const queryTimeout = time.Second

func MainHandler(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage, logger *zap.Logger) {
	switch r.Method {
	case "GET":
		getAllOrders(w, r, repo, logger)
	case "POST":
		createOrder(w, r, repo, logger)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func MainHandlerID(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage, logger *zap.Logger) {
	switch r.Method {
	case "GET":
		getOrderByID(w, r, repo, logger)
	case "PUT":
		putOrderStatus(w, r, repo, logger)
	case "DELETE":
		deleteOrder(w, r, repo, logger)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage, logger *zap.Logger) {
	logger.Info(
		"Handle request",
		zap.String("host", r.Host),
		zap.String("method", r.Method),
	)

	httpRequestBody, err := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()
	if err != nil {
		logger.Error(
			"Can't read request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ord domain.Order
	if err = json.Unmarshal(httpRequestBody, &ord); err != nil {
		logger.Error(
			"Can't deserialize request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	ordId, err := repo.CreateOrder(ctx, &ord)
	if err != nil {
		logger.Error("Can't write data in DB", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mapa := make(map[string]int)
	mapa["id"] = ordId

	responseBody, err := json.Marshal(mapa)
	if err != nil {
		logger.Error(
			"Can't serialize response body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	if _, err = w.Write(responseBody); err != nil {
		logger.Error(
			"Can't write request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Request are successfully handled")
}

func getOrderByID(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage, logger *zap.Logger) {
	logger.Info(
		"Handle request",
		zap.String("host", r.Host),
		zap.String("method", r.Method),
	)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Warn(
			"Can't convert user data to ID",
			zap.String("data", r.PathValue("id")),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	ord, err := repo.GetOrderByID(ctx, id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		logger.Info(
			"Row with ID are not found",
			zap.Int("ID", id),
		)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		logger.Info(
			"Can't read data from DB",
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalIndent(ord, "", "	")
	if err != nil {
		logger.Error(
			"Can't serialize response body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(data); err != nil {
		logger.Error(
			"Can't write request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Request are successfully handled")
}

func getAllOrders(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage, logger *zap.Logger) {
	logger.Info(
		"Handle request",
		zap.String("host", r.Host),
		zap.String("method", r.Method),
	)

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()
	ords, err := repo.GetAllOrders(ctx)
	if err != nil {
		logger.Info(
			"Can't read data from DB",
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.MarshalIndent(ords, "", "	")
	if err != nil {
		logger.Error(
			"Can't serialize response body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(responseBody); err != nil {
		logger.Error(
			"Can't write request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Request are successfully handled")
}

func putOrderStatus(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage, logger *zap.Logger) {
	logger.Info(
		"Handle request",
		zap.String("host", r.Host),
		zap.String("method", r.Method),
	)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Warn(
			"Can't convert user data to ID",
			zap.String("data", r.PathValue("id")),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()
	if err != nil {
		logger.Error(
			"Can't read request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var newStatusMap map[string]string
	if err = json.Unmarshal(
		requestBody,
		&newStatusMap,
	); err != nil {
		logger.Error(
			"Can't deserialize request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newStatus, existed := newStatusMap["Status"]
	if !existed {
		logger.Warn(
			"Put request without 'status' key",
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	changedOrder, err := repo.UpdateOrderStatus(ctx, id, newStatus)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		logger.Info(
			"Row with ID are not found",
			zap.Int("ID", id),
		)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		logger.Info(
			"Can't update data in DB",
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.MarshalIndent(*changedOrder, "", "	")
	if err != nil {
		logger.Error(
			"Can't serialize request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(data); err != nil {
		logger.Error(
			"Can't write request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Request are successfully handled")
}

func deleteOrder(w http.ResponseWriter, r *http.Request, repo repo.OrderStorage, logger *zap.Logger) {
	logger.Info(
		"Handle request",
		zap.String("host", r.Host),
		zap.String("method", r.Method),
	)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Warn(
			"Can't convert user data to ID",
			zap.String("data", r.PathValue("id")),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()
	_, err = repo.DeleteOrder(ctx, id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		logger.Info(
			"Row with ID are not found",
			zap.Int("ID", id),
		)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		logger.Info(
			"Can't delete data from DB",
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("Request are successfully handled")
}
