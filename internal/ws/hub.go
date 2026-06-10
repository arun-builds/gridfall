package ws

import "sync"

type Hub struct {
	mu    sync.Mutex
	Rooms map[string]*Room
}

func NewHub() *Hub {

	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetOrCreateRoom(id string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	room, ok := h.Rooms[id]
	if ok {
		return room
	}

	room = NewRoom(id)
	room.OnEmpty = h.DeleteRoom
	h.Rooms[id] = room

	go room.Run()

	return room
}

func (h *Hub) DeleteRoom(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.Rooms, id)
}
