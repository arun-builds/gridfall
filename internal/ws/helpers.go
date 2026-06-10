package ws

func (r *Room) getOpponent(client *Client) *Client {
	switch client {
	case r.Player1:
		return r.Player2

	case r.Player2:
		return r.Player1

	default:
		return nil
	}
}
