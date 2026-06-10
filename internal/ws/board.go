package ws

const BoardSize = 8

func makeBoard(size int) [][]int {
	board := make([][]int, size)

	for i := range board {
		board[i] = make([]int, size)
	}

	return board

}

func countRemainingEntities(board [][]int) int {
	count := 0

	for _, row := range board {
		for _, cell := range row {
			if cell > 0 {
				count++
			}
		}
	}
	return count
}
