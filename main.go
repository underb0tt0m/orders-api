package main

import (
	"context"
	"fmt"
	"net/http"
	"orders/api"
	"orders/db/db_conn/simple_db_conn"
	"orders/repo/repo_db"
)

func main() {

	ctx := context.Background()
	conn, err := simple_db_conn.GetDBConn(ctx)
	if err != nil {
		panic(err)
	}

	myRepo := repo_db.NewRepo(conn, ctx)

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
