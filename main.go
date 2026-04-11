package main

import (
	"context"
	"net/http"
	"orders/api"
	"orders/db/db_conn/simple_db_conn"
	"orders/repo/repo_db"
	"orders/zapLogger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	Logger, closeLogger, err := zapLogger.Create()
	if err != nil {
		panic(err)
	}
	Logger.Info("Initialize Logger")
	defer func() {
		Logger.Info("Close Logger")
		Logger.Sync()
		closeLogger()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signalChan
		Logger.Info("Got exit signal, exit context")
		cancel()
	}()

	Logger.Info("Create db connection")
	conn, err := simple_db_conn.GetDBConn(ctx)
	if err != nil {
		Logger.Fatal("Can't create connection", zap.Error(err))
	}
	defer func() { _ = conn.Close(ctx) }()

	Logger.Info("Initialize repository")
	myRepo := repo_db.New(conn)

	Logger.Info("Register handlers")
	http.HandleFunc("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.MainHandlerID(w, r, myRepo, Logger)
	})
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		api.MainHandler(w, r, myRepo, Logger)
	})

	server := &http.Server{Addr: ":9091"}

	go func() {
		Logger.Info("Start HTTP-server")
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Fatal("Can't start HTTP-server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	Logger.Info("Shutting down...")
	server.Shutdown(ctx)
	Logger.Info("Server stopped")
}
