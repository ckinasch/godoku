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

// Populate boxes 0, 4, 8; they are diagonal and therefore independant.
// This completes as much of the board as possible before attempting to brute force the remaining boxes;
// (theoretically) increasing the chance of success
func SeedInitialBoxes(b *Board) {

	// 2d array of random integers used to seed boxes 0, 4, 8 (diagonal)
	vals := [3][]int{GetRandomVals(), GetRandomVals(), GetRandomVals()}

	// Complementary Array of indices for incrementing the vals array when val has been used
	i := [3]int{0, 0, 0}

	for r, row := range b.cells {
		for c := range row {
			cur := &b.cells[r][c]
			if cur.box == 0 {
				cur.value = vals[0][i[0]]
				i[0]++
			}
			if cur.box == 4 {
				cur.value = vals[1][i[1]]
				i[1]++
			}
			if cur.box == 8 {
				cur.value = vals[2][i[2]]
				i[2]++
			}
		}
	}

}

func PopulateBoard(b *Board) {

	for r, row := range b.cells {
		for c := range row {
			box := r/3*3 + c/3
			cur := &b.cells[r][c]
			cur.box = box
		}
	}

	SeedInitialBoxes(b)

}

func main() {

	PopulateBoard(&b)

	PrintBoard(&b)

}

func init() {
}
