package ws

import (
	"log"
)

type Room struct {
	ID string

	Player1 *Client
	Player2 *Client

	State MatchState

	Register   chan *Client
	Unregister chan *Client
}

func NewRoom(id string) *Room {
	log.Printf("creating room: %s", id)

	board1 := makeBoard(5)
	board2 := makeBoard(5)

	// Temporary hardcoded entities
	board1[0][0] = 1
	board1[2][1] = 1

	board2[1][2] = 1
	board2[3][3] = 1

	return &Room{
		ID: id,
		State: MatchState{
			CurrentTurn: "",
			Board1:      board1,
			Board2:      board2,
		},
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.Register:

			log.Printf("registering client %s", client.ID)

			if r.Player1 == nil {

				r.Player1 = client

				log.Printf(
					"assigned player1=%s",
					client.ID,
				)

			} else if r.Player2 == nil {

				r.Player2 = client

				log.Printf(
					"assigned player2=%s",
					client.ID,
				)

			} else {

				log.Printf(
					"room full, rejecting client %s",
					client.ID,
				)

				client.Send <- []byte(
					`{"type":"error","message":"room full"}`,
				)
				client.Conn.Close()

				continue
			}

			if r.State.CurrentTurn == "" {

				r.State.CurrentTurn = client.ID

				log.Printf(
					"first turn assigned to %s",
					client.ID,
				)
			}

			log.Printf(
				"room status p1=%s p2=%s currentTurn=%q",
				playerID(r.Player1),
				playerID(r.Player2),
				r.State.CurrentTurn,
			)

		case client := <-r.Unregister:

			log.Printf(
				"unregistering client %s",
				client.ID,
			)

			owned := false

			if r.Player1 == client {

				r.Player1 = nil

				owned = true
			}

			if r.Player2 == client {

				r.Player2 = nil

				owned = true
			}

			// If the player whose turn it was disconnected,
			// clear the turn.
			if owned && r.State.CurrentTurn == client.ID {

				r.State.CurrentTurn = ""

				log.Printf(
					"current turn cleared because %s disconnected",
					client.ID,
				)
			}

			// Only accepted players own a Send channel
			// that the room is responsible for closing.
			if owned {

				close(client.Send)
			}

			log.Printf(
				"room status p1=%s p2=%s currentTurn=%q",
				playerID(r.Player1),
				playerID(r.Player2),
				r.State.CurrentTurn,
			)
		}
	}
}

func playerID(c *Client) string {
	if c == nil {
		return "-"
	}

	return c.ID
}
