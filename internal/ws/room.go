package ws

import (
	"log"
)

type Room struct {
	ID         string
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.Register:
			r.Clients[client] = true
			log.Printf("ws client connected id: %s", client.ID)

		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
				log.Printf("ws client disconnected id: %s", client.ID)
			}

		case message := <-r.Broadcast:
			log.Printf("room %s: broadcasting message (%d bytes) to %d clients", r.ID, len(message), len(r.Clients))
			for client := range r.Clients {
				select {
				case client.Send <- message:

				default:
					log.Printf("room %s: dropping slow client", r.ID)
					delete(r.Clients, client)
					close(client.Send)
				}
			}
		}
	}
}
