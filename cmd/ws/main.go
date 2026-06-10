package main

import (
	"log"
	"net/http"

	"github.com/arun-builds/gridfall/internal/ws"
)

func main() {
	hub := ws.NewHub()

	http.HandleFunc("/ws", hub.HandleWS)

	log.Println("WS server listening on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
