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
	cells  [9]*Cell
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

// func (b *Board) _PopulateBox(box int) error {
	
// 	vals := GetRandomVals()
// 	i := 0
// 	for r, row := range b.cells {
// 		for c := range row {
// 			cur := &b.cells[r][c]

// 			if cur.box == box {

// 				valueConflict := cur.CheckValueConflict(vals[i])
// 				fmt.Println(i, len(vals))

// 				switch {
// 				case !valueConflict:
// 					cur.value = vals[i]
// 					vals = removeSliceZero(vals)
// 					fmt.Println(vals)

// 				case i+1 >= len(vals):

// 					return fmt.Errorf("no solution: %v %v", cur, vals[i])

// 				case valueConflict:
// 					i++
// 					// default:
// 					// return fmt.Errorf("unhandled exception: i=%v, vals=%v", cur, vals[i])
// 				}

// 			}
// 		}
// 	}
// 	return nil
// }

func (b *Box) PopulateBox(){
	for i, v := range GetRandomVals(){
		b.cells[i].value = v
	}

}

func (b *Board) PopulateBoard() {

	// Populate diagonal boxes first to increase chance of success of brute force
	// box_order := [9]int{0, 4, 8, 1, 2, 3, 5, 6, 7}

	for i := 0; i < 9; i++ {
		// Index the box
		b.boxes[i].index = i
		for j := 0; j < 9; j++ {
			
			// Get the boxes position using division truncation trick
			box := i/3*3 + j/3

			// Store reference to current created cell at this position
			cell := &Cell{row: i, col: j, box: box}

			// Create reference to cell in board row/col
			b.cells[i][j] = cell

			// Create reference to cell in appropriate box and position in the boxes index
			b.boxes[i/3*3+j/3].cells[i%3*3+j%3] = cell

		}
	}

	b.boxes[0].PopulateBox()

	// for _, v := range b.boxes[0].cells{
		// b.cells[v.row][v.col].value = v.value
	// }






	// for i := 0; i < 9; i++ {
	// boxn := box_order[i]
	// err := b.PopulateBox(boxn)
	// if err != nil {
	// fmt.Println(err)
	// i--
	// b.PrintBoard()

	// }
	// }

	// b.PopulateBox(0)
	// b.PopulateBox(4)
	// b.PopulateBox(8)

	// for {
	// 	err := b.PopulateBox(1)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		b.PrintBoard()
	// 	}
	// }

}

func main() {

	board.PopulateBoard()

	board.PrintBoard()


}

func init() {
}
