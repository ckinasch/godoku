package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var _ = errors.New
var _ = fmt.Println
var _ = rand.Perm

type Cell struct {
	row   int
	col   int
	value int
	box   Box
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
			fmt.Print(cur.value)

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

// Return a slice of randomly generated distinct integers; all values are incremented
func GetRandomVals() []int {

	vals := rand.Perm(9)
	for i := range vals {
		vals[i]++
	}
	return vals

}

func PopulateBoard(b *Board) {

	var _ = GetRandomVals()

	for r, row := range b.cells {
		for c := range row {
			cur := &b.cells[r][c]
			cur.value = r
		}
	}

}

func main() {

	PopulateBoard(&b)

	PrintBoard(&b)

}

func init() {
	rand.Seed(time.Now().Unix())
}
