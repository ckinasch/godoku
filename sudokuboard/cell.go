package sudokuboard

type Cell struct {
	row   int
	col   int
	Value int
	box   int
}

func (c *Cell) checkValueConflict(val int) bool {

	for i := 0; i < nRowsCols; i++ {
		if val == board.Cells[i][c.col].Value {
			return true
		}
		if val == board.Cells[c.row][i].Value {
			return true
		}
	}
	return false
}
