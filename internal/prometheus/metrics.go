package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Определяем метрики как глобальные переменные
var (
	// COUNTER: количество созданных заказов
	OrdersCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "app_orders_created_total",
			Help: "Total number of orders created",
		},
	)

	// GAUGE: количество активных заказов (без удаленных)
	OrdersActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_orders_active",
			Help: "Number of active orders (excluding deleted ones)",
		},
	)

	// COUNTER: количество обработанных HTTP-запросов
	RequestsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "Total number of handled requests",
		},
	)

	//HISTOGRAM: гистограмма времени обработки запросов
	RequestDuratuion = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "app_request_duratuion",
			Help:    "Duration of proccessing request",
			Buckets: []float64{600, 1200, 1800, 2400, 3000},
		})
)
