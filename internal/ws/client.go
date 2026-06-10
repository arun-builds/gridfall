package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
	Room *Room
}

func (c *Client) ReadPump() {
	defer func() {
		log.Printf("ws client disconnected from room %s", c.Room.ID)
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			log.Printf("ws read error: %v", err)
			break
		}

		var event Event

		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("client=%s invalid json", c.ID)
			continue
		}

		log.Printf(
			"client=%s event=%s",
			c.ID,
			event.Type,
		)

		switch event.Type {

		case "attack":
			c.handleAttack(event)
		default:
			c.Send <- []byte(`{"type":"error","message":"unknown event"}`)
		}
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
			log.Printf("ws write error: %v", err)
			return
		}
	}
}

func (c *Client) handleAttack(event Event) {

	var payload AttackPayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {

		c.Send <- []byte(`{"type":"error","message":"invalid attack payload"}`)

		return
	}

	if c.Room.State.CurrentTurn != c.ID {

		c.Send <- []byte(`{"type":"error","message":"not your turn"}`)

		return
	}

	opponent := c.Room.Opponent(c)

	if opponent == nil {

		c.Send <- []byte(`{"type":"error","message":"waiting for opponent"}`)

		return
	}

	log.Printf(
		"attack accepted from %s at (%d,%d)",
		c.ID,
		payload.X,
		payload.Y,
	)

	c.Room.State.CurrentTurn = opponent.ID

	c.Send <- []byte(
		fmt.Sprintf(
			`{
				"type":"attack_accepted",
				"x":%d,
				"y":%d
			}`,
			payload.X,
			payload.Y,
		),
	)

	opponent.Send <- []byte(
		fmt.Sprintf(
			`{
				"type":"opponent_attacked",
				"x":%d,
				"y":%d
			}`,
			payload.X,
			payload.Y,
		),
	)

	log.Printf("turn switched to %s", opponent.ID)
}
