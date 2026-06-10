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

func maskBoard(board [][]int) [][]int {
	masked := makeBoard(len(board))

	for y := range board {
		for x := range board[y] {

			cell := board[y][x]

			switch {

			case cell == MissCell:
				masked[y][x] = MissCell

			case cell < 0:
				masked[y][x] = cell

			default:
				masked[y][x] = UnknownCell
			}
		}
	}

	return masked
}

func copyBoard(board [][]int) [][]int {
	boardCopy := makeBoard(len(board))

	for y := range board {
		copy(boardCopy[y], board[y])
	}

	return boardCopy
}
