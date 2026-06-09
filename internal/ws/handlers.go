package ws

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) HandleWS(
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:   uuid.New().String(),
		Conn: conn,
		Send: make(chan []byte, 256),
		Room: h.Room,
	}

	h.Room.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
