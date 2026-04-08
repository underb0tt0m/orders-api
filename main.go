package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"orders/api"
	"orders/db/db_conn/simple_db_conn"
	"orders/repo/repo_db"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signalChan
		log.Println("got exit signal, cancel context")
		cancel()
	}()

	conn, err := simple_db_conn.GetDBConn(ctx)
	if err != nil {
		log.Fatalf("can't create conn: %v", err)
	}
	defer conn.Close(ctx)

	/*
		pgx conn -> repo -> service -> handler
		каждый слой отдельно, service не должен зависеть от pgx никак
	*/

	myRepo := repo_db.New(conn)

	http.HandleFunc("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.MainHandlerID(w, r, myRepo)
	})
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		api.MainHandler(w, r, myRepo)
	})

	if err := http.ListenAndServe(":9091", nil); err != nil {
		log.Fatalf("can't start server: %v", err)
	}
}
