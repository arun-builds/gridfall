package ws

type Hub struct {
	Room *Room
}

func NewHub() *Hub {
	room := NewRoom("global")

	return &Hub{
		Room: room,
	}
}

func (h *Hub) Run() {
	go h.Room.Run()
}
