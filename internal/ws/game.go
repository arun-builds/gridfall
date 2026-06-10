package ws

import (
	"encoding/json"
	"fmt"
	"log"
)

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

	var board [][]int
	if c.Room.Player1 == c {
		board = c.Room.State.Board2
	} else {
		board = c.Room.State.Board1
	}

	if opponent == nil {

		c.Send <- []byte(`{"type":"error","message":"waiting for opponent"}`)

		return
	}

	if payload.X < 0 ||
		payload.X >= len(board) ||
		payload.Y < 0 ||
		payload.Y >= len(board[0]) {

		c.Send <- []byte(
			`{"type":"error","message":"invalid coordinates"}`,
		)

		return
	}

	log.Printf(
		"attack accepted from %s at (%d,%d)",
		c.ID,
		payload.X,
		payload.Y,
	)

	result := ""

	switch board[payload.X][payload.Y] {

	case 0:
		board[payload.X][payload.Y] = 3
		result = "miss"

	case 1:
		board[payload.X][payload.Y] = 2
		result = "hit"

	case 2, 3:
		c.Send <- []byte(
			`{"type":"error","message":"cell already attacked"}`,
		)

		return
	}

	log.Printf(
		"%s attacked (%d,%d): %s",
		c.ID,
		payload.X,
		payload.Y,
		result,
	)

	remaining := countRemainingEntities(board)

	log.Printf(
		"remaining entities: %d",
		remaining,
	)

	if remaining == 0 {

		log.Printf(
			"game over winner=%s",
			c.ID,
		)

		c.Send <- []byte(
			fmt.Sprintf(
				`{
					"type":"game_over",
					"winner":"%s"
				}`,
				c.ID,
			),
		)

		opponent.Send <- []byte(
			fmt.Sprintf(
				`{
					"type":"game_over",
					"winner":"%s"
				}`,
				c.ID,
			),
		)

		return
	}

	c.Room.State.CurrentTurn = opponent.ID

	c.Send <- []byte(
		fmt.Sprintf(
			`{
				"type":"attack_result",
				"result":"%s",
				"x":%d,
				"y":%d
			}`,
			result,
			payload.X,
			payload.Y,
		),
	)

	opponent.Send <- []byte(
		fmt.Sprintf(
			`{
				"type":"opponent_attacked",
				"result":"%s",
				"x":%d,
				"y":%d
			}`,
			result,
			payload.X,
			payload.Y,
		),
	)

	log.Printf("turn switched to %s", opponent.ID)
}
