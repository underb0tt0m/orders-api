package main

import (
	"fmt"
	"net/http"
	"orders/api"
	"orders/repo"
)

func main() {
	myRepo := repo.NewRepo()
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
