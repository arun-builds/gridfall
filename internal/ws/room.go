package ws

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
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

	board1 := makeBoard(BoardSize)
	board2 := makeBoard(BoardSize)

	return &Room{
		ID: id,

		State: MatchState{
			Phase:       PhaseWaiting,
			CurrentTurn: "",
			Board1:      board1,
			Board2:      board2,
		},
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// TODO:
// Room.Run(), Attack(), and PlaceEntities()
// currently mutate shared Room state from different goroutines.
//
// This is acceptable for the current single-room prototype.
// Migrate game actions through Room.Events to establish
// single ownership before introducing reconnects,
// matchmaking, or multiple rooms.

func (r *Room) Run() {
	for {
		select {

		case client := <-r.Register:

			log.Printf("registering client %s", client.ID)

			if r.Player1 == nil {
				r.Player1 = client
				log.Printf("assigned player1=%s", client.ID)
			} else if r.Player2 == nil {
				r.Player2 = client
				log.Printf("assigned player2=%s", client.ID)
				r.State.Phase = PhasePlacement
				log.Printf("placement phase started")
			} else {

				log.Printf(
					"room full, rejecting client %s",
					client.ID,
				)

				client.sendJSON(ErrorResponse{
					Type:    "error",
					Message: "room full",
				})
				client.Conn.Close()

				continue
			}

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

func (r *Room) PlaceEntities(
	client *Client,
	placements []Placement,
) error {

	if r.State.Phase != PhasePlacement {
		return errors.New("placement phase is not active")
	}

	var (
		board     [][]int
		readyFlag *bool
	)

	switch client {
	case r.Player1:
		board = r.State.Board1
		readyFlag = &r.State.Player1Ready

	case r.Player2:
		board = r.State.Board2
		readyFlag = &r.State.Player2Ready

	default:
		return errors.New("client does not belong to this room")
	}

	if *readyFlag {
		return errors.New("fleet already placed")
	}

	required := map[int]bool{
		Scout:      false,
		Battleship: false,
		Mage:       false,
		Assassin:   false,
	}

	if len(placements) != len(required) {
		return errors.New("invalid fleet size")
	}

	occupied := make(map[string]bool)

	// Validation pass
	for _, p := range placements {

		if _, ok := required[p.Entity]; !ok {
			return fmt.Errorf("invalid entity: %d", p.Entity)
		}

		if required[p.Entity] {
			return fmt.Errorf("duplicate entity: %d", p.Entity)
		}

		required[p.Entity] = true

		if p.X < 0 || p.X >= BoardSize ||
			p.Y < 0 || p.Y >= BoardSize {

			return fmt.Errorf(
				"invalid coordinates (%d,%d)",
				p.X,
				p.Y,
			)
		}

		key := fmt.Sprintf("%d:%d", p.X, p.Y)

		if occupied[key] {
			return errors.New("overlapping placements")
		}

		occupied[key] = true
	}

	// Populate board only after validation succeeds
	for _, p := range placements {
		board[p.Y][p.X] = p.Entity
	}

	*readyFlag = true

	client.sendJSON(map[string]any{
		"type": "placement_success",
	})

	opponent := r.getOpponent(client)

	if opponent != nil {
		opponent.sendJSON(map[string]any{
			"type": "opponent_ready",
		})
	}

	if r.State.Player1Ready && r.State.Player2Ready {

		r.State.Phase = PhaseBattle

		if rand.Intn(2) == 0 {
			r.State.CurrentTurn = r.Player1.ID
		} else {
			r.State.CurrentTurn = r.Player2.ID
		}

		battleMsg := map[string]any{
			"type":         "battle_started",
			"current_turn": r.State.CurrentTurn,
		}

		r.Player1.sendJSON(battleMsg)
		r.Player2.sendJSON(battleMsg)
	}

	return nil
}

func (r *Room) Attack(
	client *Client,
	payload AttackPayload,
) error {

	if r.State.Phase != PhaseBattle {
		return errors.New("cannot attack during placement")

	}

	if r.State.CurrentTurn != client.ID {
		return errors.New("not your turn")
	}

	opponent := client.Room.getOpponent(client)
	if opponent == nil {
		return errors.New("waiting for opponent")
	}

	var board [][]int
	if client.Room.Player1 == client {
		board = client.Room.State.Board2
	} else {
		board = client.Room.State.Board1
	}

	if payload.X < 0 ||
		payload.X >= len(board) ||
		payload.Y < 0 ||
		payload.Y >= len(board[0]) {
		return errors.New("invalid coordinates")
	}

	log.Printf(
		"attack accepted from %s at (%d,%d)",
		client.ID,
		payload.X,
		payload.Y,
	)

	result := ""

	cell := board[payload.Y][payload.X]
	destroyedEntity := ""

	switch cell {
	case EmptyCell:
		board[payload.Y][payload.X] = MissCell
		result = "miss"
	case Scout:
		board[payload.Y][payload.X] = DestroyedScout
		result = "hit"
		destroyedEntity = "scout"

	case Battleship:
		board[payload.Y][payload.X] = DestroyedBattleship
		result = "hit"
		destroyedEntity = "battleship"

	case Mage:
		board[payload.Y][payload.X] = DestroyedMage
		result = "hit"
		destroyedEntity = "mage"

	case Assassin:
		board[payload.Y][payload.X] = DestroyedAssassin
		result = "hit"
		destroyedEntity = "assassin"

	default:

		return errors.New("cell already attacked")
	}

	log.Printf(
		"%s attacked (%d,%d): %s",
		client.ID,
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
			client.ID,
		)

		r.State.Phase = PhaseGameOver

		msg := AttackResult{Type: "game_over", Winner: client.ID}
		client.sendJSON(msg)
		opponent.sendJSON(msg)

		return nil
	}

	client.Room.State.CurrentTurn = opponent.ID

	client.sendJSON(AttackResult{
		Type:            "attack_result",
		Result:          result,
		X:               payload.X,
		Y:               payload.Y,
		DestroyedEntity: destroyedEntity,
	})

	opponent.sendJSON(AttackResult{
		Type:   "opponent_attacked",
		Result: result,
		X:      payload.X,
		Y:      payload.Y,
	})

	log.Printf("turn switched to %s", opponent.ID)

	return nil
}
