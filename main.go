package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"orders/api"
	"orders/db/db_conn/simple_db_conn"
	"orders/repo/repo_inmemory"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signalChan
		log.Println("got exit signal, exit context")
		cancel()
	}()

	conn, err := simple_db_conn.GetDBConn(ctx)
	if err != nil {
		log.Fatalf("can't create conn: %v", err)
	}
	defer func() { _ = conn.Close(ctx) }()

	myRepo := repo_inmemory.New()

	http.HandleFunc("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.MainHandlerID(w, r, myRepo)
	})
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		api.MainHandler(w, r, myRepo)
	})

	if err := http.ListenAndServe(":9091", nil); err != nil {
		fmt.Println("Возникла ошибка")
	}
}
