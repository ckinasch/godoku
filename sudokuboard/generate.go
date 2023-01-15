package sudokuboard

const (
	// Alter this value for different sized boards: 2x2; 3x3; 4x4 etc (processing takes longer the greater the number)
	gridSize int = 3
	// Number of cells to remove when creating the board
	diffEasy   int = 10
	diffMedium int = 20
	diffHard   int = 30
)

var (
	board Board
	// solveBoard() helper variable; 0, 1, 1+ for count of solutions to the board, board will only be valid if there is a single solution
	n_solutions int = 1
	difficulty  int = diffEasy
)

func GetBoard() Board {

	board.indexBoard()

	board.solveBoard(true)

	board.removeCells()

	return board

}

func init() {
	seedRand()

}
