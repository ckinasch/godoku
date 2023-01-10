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
	vals  [9]int
}

type Board struct {
	cells [9][9]Cell
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
			// fmt.Printf("%d:%d", cur.value, cur.box)
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

func (b *Board) PopulateBox(box int) {
	vals := GetRandomVals()
	i := 0
	for r, row := range b.cells {
		for c := range row {
			cur := &b.cells[r][c]

			if cur.box == box {

				valueConflict := cur.CheckValueConflict(vals[i])

				switch {
				case !valueConflict:
					cur.value = vals[i]
					vals = removeSliceZero(vals)
					fmt.Println(vals)

				case i+1 == len(vals):
					fmt.Println("No solution", cur, vals[i])

				case valueConflict:
					i++
				}

			}
		}
	}
}

func (b *Board) PopulateBoard() {

	// Populate diagonal boxes first to increase chance of success of brute force
	box_order := [9]int{0, 4, 8, 1, 2, 3, 5, 6, 7}

	for r, row := range b.cells {
		for c := range row {
			box := r/3*3 + c/3
			cur := &b.cells[r][c]
			cur.row = r
			cur.col = c
			cur.box = box
		}
	}

	for _, v := range box_order {
		b.PopulateBox(v)
	}

}

func main() {

	board.PopulateBoard()

	board.PrintBoard()

}

func init() {
}
