package main

import (
	"log"
	"net/http"

	controller "github.com/elnur0000/fish-backend/pkg/controller"
	"github.com/elnur0000/fish-backend/pkg/transports"
)

func main() {
	http.HandleFunc("/main.wasm", transports.CORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/wasm")
		http.ServeFile(w, r, "bin/main.wasm")
	}))

	gamePool := controller.NewGamePool(controller.GamePoolOptions{
		Transport: &transports.WebRTCTransport,
	})

	go gamePool.Start()

	log.Println("Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
