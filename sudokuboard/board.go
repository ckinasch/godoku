package sudokuboard

import "fmt"

type Board struct {
	Cells [nRowsCols][nRowsCols]*Cell
	boxes [nRowsCols]Box
}

const (

	nRowsCols  int = gridSize * gridSize
	totalCells     = nRowsCols * nRowsCols
)

func (b *Board) PrintBoard() {
	lines := "---------------------"
	fmt.Println(lines, lines)
	for r, row := range b.Cells {
		for c := range row {
			cur := b.Cells[r][c]
			fmt.Print(" | ")
			fmt.Printf("%d", cur.Value)

			if (c+1)%gridSize == 0 {
				fmt.Print(" |")
			}
		}
		if (r+1)%gridSize == 0 {
			fmt.Println("")
			fmt.Println(lines, lines)

		} else {
			fmt.Println("")
		}
	}
}

// Check if the board has been completed
func (b *Board) checkBoard() bool {

	for i := 0; i < totalCells; i++ {
		row := i / nRowsCols
		col := i % nRowsCols
		if b.Cells[row][col].Value == 0 {
			// Board incomplete
			return false
		}
	}
	// Board has been completed, trigger recursive exit
	return true
}

// Populate the structs with references to cells;
// Index by Box and Board position to simplify processing and access
func (b *Board) indexBoard() {

	for i := 0; i < totalCells; i++ {
		row := i / nRowsCols
		col := i % nRowsCols
		// Index the box
		b.boxes[row].index = row

		// Get the box position
		ibox := row/gridSize*gridSize + col/gridSize

		// Store reference to current created cell at this position
		cell := &Cell{row: row, col: col, box: ibox}

		// Create reference to cell in board row/col
		b.Cells[row][col] = cell

		// Create reference to cell in appropriate box and position in the boxes index
		b.boxes[ibox].cells[row%gridSize*gridSize+col%gridSize] = cell

	}

}

func (b *Board) removeCells() {
	// Reset val
	n_solutions = 1
	for {

		random_vals := getRandomVals()
		rRow := random_vals[0] - 1
		rCol := random_vals[1] - 1

		// Get a new cell if the value is 0
		for b.Cells[rRow][rCol].Value == 0 {
			random_vals = getRandomVals()
			rRow = random_vals[0] - 1
			rCol = random_vals[1] - 1
		}

		old_val := b.Cells[rRow][rCol].Value
		b.Cells[rRow][rCol].Value = 0

		// Either 0 or multiple solutions; invalid board, restore the value and repeat
		if n_solutions != 1 {
			b.Cells[rRow][rCol].Value = old_val
		}

		if threshold_met, _ := b.countZeroCells(); threshold_met {
			b.solveBoard(false)
			break
		}

	}
}

// Returns true/false if the desired number of cells at 0 is met
func (b *Board) countZeroCells() (bool, int) {

	zero_cells := 0

	for i := 0; i < totalCells; i++ {
		row := i / nRowsCols
		col := i % nRowsCols
		cell := b.Cells[row][col]

		if cell.Value == 0 {
			zero_cells++
		}
	}

	return zero_cells >= difficulty, zero_cells

}

// Recursive backtracking function
// Begins placing a random sequence and placing valid values
// Backtraces as far as required to create a valid solution
// "fill" parameter should be
// : TRUE when trying to generate a board
// : FALSE when checking if the board is still valid and uniquely solvable
func (b *Board) solveBoard(fill bool) bool {

	for i := 0; i < totalCells; i++ {
		row := i / nRowsCols
		col := i % nRowsCols
		// Refer to current cell
		cell := b.Cells[row][col]

		if cell.Value == 0 {

			// Next set of values to attempt placement
			// Iterate and recursively place
			vals := getRandomVals()
			for _, v := range vals {

				// Value is valid and can be placed
				if !cell.checkValueConflict(v) && !b.boxes[cell.box].containsValue(v) {
					cell.Value = v

					// Board is completed
					if b.checkBoard() {
						n_solutions++

						if fill {
							return true
						} else {
							break
						}

					} else {
						if b.solveBoard(fill) {
							return true
						}
					}
				}
			}
			// No solution from here; reset the current cell's value and break loop
			cell.Value = 0
			break
		}
	}

	// Trace back to previous chain and attempt a different value
	return false

}
