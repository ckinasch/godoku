package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	_ "errors"
	"fmt"
	math_rand "math/rand"
	_ "time"
)

type Cell struct {
	row   int
	col   int
	value int
	box   int
}

// TODO: vals as reference to Cells?
// on instantiate step, point to vals
type Box struct {
	index int
	cells [9]*Cell
}

type Board struct {
	cells [9][9]*Cell
	boxes [9]Box
}

var board Board

func (b *Board) PrintBoard() {
	lines := "---------------------"
	fmt.Println(lines)
	for r, row := range b.cells {
		for c := range row {
			cur := b.cells[r][c]
			fmt.Print("|")
			fmt.Printf("%d", cur.value)

			if (c+1)%3 == 0 {
				fmt.Print("|")
			}
		}
		if (r+1)%3 == 0 {
			fmt.Println("")
			fmt.Println(lines)

		} else {
			fmt.Println("")
		}
	}
}

// Simultaneous calls for new random slices may overlap; time based seed is insufficient;
// Cryptographically secure approach to ensure unique slices
func seedRand() {
	// https://stackoverflow.com/a/54491783

	var b [8]byte
	_, err := crypto_rand.Read(b[:])

	if err != nil {
		panic("Cannot seed")
	}

	seedval := int64(binary.LittleEndian.Uint64(b[:]))
	math_rand.Seed(seedval)
}

// Return a slice of randomly generated distinct integers; all values are incremented
func GetRandomVals() []int {

	seedRand()

	vals := math_rand.Perm(9)
	for i := range vals {
		vals[i]++
	}
	return vals

}

// Replace the first element of a slice with the last one and return a slice with the tail cut off
func removeSliceZero(s []int) []int {
	s[0] = s[len(s)-1]
	return s[:len(s)-1]
}

func (c *Cell) CheckValueConflict(val int) bool {

	for i := 0; i < 9; i++ {
		if val == board.cells[i][c.col].value {
			return true
		}
		if val == board.cells[c.row][i].value {
			return true
		}
	}
	return false
}

func (b *Box) ContainsValue(val int) bool {

	for i := 0; i < len(b.cells); i++ {
		if b.cells[i].value == val {
			return true
		}
	}
	return false

}

func (b *Board) IndexBoard() {

	// Populate diagonal boxes first to increase chance of success of brute force
	// box_order := [9]int{0, 4, 8, 1, 2, 3, 5, 6, 7}

	for i := 0; i < 81; i++ {
		row := i / 9
		col := i % 9
		// Index the box
		b.boxes[row].index = row

		// Get the box position
		box := row/3*3 + col/3

		// Store reference to current created cell at this position
		cell := &Cell{row: row, col: col, box: box}

		// Create reference to cell in board row/col
		b.cells[row][col] = cell

		// Create reference to cell in appropriate box and position in the boxes index
		b.boxes[row/3*3+col/3].cells[row%3*3+col%3] = cell

	}

}

// Check if the board has been completed
func (b *Board) CheckBoard() bool {

	for i := 0; i < 81; i++ {
		row := i / 9
		col := i % 9
		if b.cells[row][col].value == 0 {
			// Board incomplete
			return false
		}
	}
	// Board has been completed, trigger recursive exit
	return true
}

// Recursive backtracking function
// Begins placing a random sequence and placing valid values
// Backtraces as far as required to create a valid solution
func (b *Board) FillBoard() (bool, error) {

	for i := 0; i < 81; i++ {
		row := i / 9
		col := i % 9
		// Refer to current cell
		cell := b.cells[row][col]

		if cell.value == 0 {

			// Next set of values to attempt placement
			// Iterate and recursively place
			vals := GetRandomVals()
			for _, v := range vals {

				// Value is valid and can be placed
				if !cell.CheckValueConflict(v) && !b.boxes[cell.box].ContainsValue(v) {
					cell.value = v

					// Board is completed
					if b.CheckBoard() {
						return true, nil

					} else {
						if next, _ := b.FillBoard(); next {
							return true, nil
						}
					}
				}
			}
			// No solution from here; reset the current cell's value and break loop
			b.cells[row][col].value = 0
			break
		}
	}

	// Trace back to previous cell and attempt a different value
	return false, nil

}

func main() {

	board.IndexBoard()

	board.FillBoard()
	board.PrintBoard()

}

func init() {
}
