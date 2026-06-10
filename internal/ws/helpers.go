package ws

func (r *Room) Opponent(c *Client) *Client {

	if r.Player1 == c {
		return r.Player2
	}

	if r.Player2 == c {
		return r.Player1
	}

	return nil
}
