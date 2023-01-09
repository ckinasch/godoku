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

type Box struct {
	index int
	vals  [9]int
}

type Board struct {
	cells [9][9]Cell
	boxes [9]Box
}

var b Board

func PrintBoard(b *Board) {
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

func PopulateBox(box int, b *Board) {
	vals := GetRandomVals()
	i := 0
	for r, row := range b.cells {
		for c := range row {
			cur := &b.cells[r][c]

			if cur.box == box {

				// Check if val can be placed
				// retry until it can be placed

				// Place if possible
				cur.value = vals[i]
				i++

			}
		}
	}
}

func PopulateBoard(b *Board) {

	// Populate diagonal boxes first to increase chance of success of brute force
	seedBoxes := [3]int{0, 4, 8}
	// Populate remaining boxes
	remainingBoxes := [6]int{1, 2, 3, 5, 6, 7}

	for r, row := range b.cells {
		for c := range row {
			box := r/3*3 + c/3
			cur := &b.cells[r][c]
			cur.box = box
		}
	}

	for _, v := range seedBoxes {
		PopulateBox(v, b)
	}

	for _, v := range remainingBoxes {
		PopulateBox(v, b)
	}

}

func main() {

	PopulateBoard(&b)

	PrintBoard(&b)

}

func init() {
}
