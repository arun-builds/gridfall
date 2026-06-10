package ws

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type AttackPayload struct {
	X int `json:"x"`
	Y int `json:"y"`
}
