package ws

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

		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}

		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.Send <- message:

				default:
					delete(r.Clients, client)
					close(client.Send)
				}
			}
		}
	}
}
