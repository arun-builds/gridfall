package ws

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Room *Room
}

func (c *Client) ReadPump() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			break
		}

		c.Room.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		err := c.Conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		if err != nil {
			return
		}
	}
}
