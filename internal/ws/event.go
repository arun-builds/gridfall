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

type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
