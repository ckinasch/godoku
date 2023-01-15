package sudokuboard

type Box struct {
	index int
	cells [nRowsCols]*Cell
}

func (b *Box) containsValue(val int) bool {

	for i := 0; i < len(b.cells); i++ {
		if b.cells[i].Value == val {
			return true
		}
	}
	return false

}
