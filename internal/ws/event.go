package ws

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type PlaceEntitiesPayload struct {
	Placements []Placement `json:"placements"`
}

type Placement struct {
	Entity int `json:"entity"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type AttackPayload struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type AttackResult struct {
	Type            string `json:"type"`
	Result          string `json:"result"`
	X               int    `json:"x,omitempty"`
	Y               int    `json:"y,omitempty"`
	DestroyedEntity string `json:"destroyed_entity,omitempty"`
	Winner          string `json:"winner,omitempty"`
}

type GameStateResponse struct {
	Type         string     `json:"type"`
	YourID       string     `json:"your_id"`
	Phase        MatchPhase `json:"phase"`
	CurrentTurn  string     `json:"current_turn"`
	YourBoard    [][]int    `json:"your_board"`
	OpponentView [][]int    `json:"opponent_view"`
}

type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
