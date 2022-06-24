package main

import (
	"log"
	"net/http"

	"github.com/elnur0000/fish-backend/pkg/controllers"
	"github.com/elnur0000/fish-backend/pkg/transports"
)

func main() {
	http.HandleFunc("/main.wasm", transports.CORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/wasm")
		http.ServeFile(w, r, "bin/main.wasm")
	}))

	g := controllers.NewGameController(controllers.GameOptions{
		Transport: &transports.WebRTCTransport,
	})

	go g.StartGameLoop()

	log.Println("Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
