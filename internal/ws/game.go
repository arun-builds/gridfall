package ws

import (
	"encoding/json"
	"log"
)

func (c *Client) sendJSON(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return
	}
	c.Send <- b
}

func (c *Client) handleAttackEvent(event Event) {
	var payload AttackPayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		c.sendJSON(ErrorResponse{
			Type:    "error",
			Message: "invalid attack payload",
		})
		return
	}

	if err := c.Room.Attack(c, payload); err != nil {
		c.sendJSON(ErrorResponse{
			Type:    "error",
			Message: err.Error(),
		})
	}
}

func (c *Client) handlePlacementEvent(event Event) {
	var payload PlaceEntitiesPayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		c.sendJSON(ErrorResponse{
			Type:    "error",
			Message: "invalid placement payload",
		})
		return
	}

	if err := c.Room.PlaceEntities(c, payload.Placements); err != nil {
		c.sendJSON(ErrorResponse{
			Type:    "error",
			Message: err.Error(),
		})
		return
	}
}

func (c *Client) handleGetStateEvent() {
	state := c.Room.GetState(c)
	c.sendJSON(state)
}
