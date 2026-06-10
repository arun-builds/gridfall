package ws

// 0 = empty
// 1 = alive entity
// 2 = destroyed entity
// 3 = miss

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
			if cell == 1 {
				count++
			}
		}
	}
	return count
}
