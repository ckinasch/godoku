package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"html/template"
	math_rand "math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var _ = template.ErrAmbigContext

type Cell struct {
	row   int
	col   int
	Value int
	box   int
}

type Box struct {
	index int
	cells [9]*Cell
}

type Board struct {
	cells [9][9]*Cell
	boxes [9]Box
}

const (
	EASY   int = 10
	MEDIUM int = 20
	HARD   int = 30
)

var (
	board Board
	// SolveBoard() helper variable; 0, 1, 1+ for count of solutions to the board, board will only be valid if there is a single solution
	n_solutions int = 1
	difficulty  int = EASY
)

func (b *Board) PrintBoard() {
	lines := "---------------------"
	fmt.Println(lines)
	for r, row := range b.cells {
		for c := range row {
			cur := b.cells[r][c]
			fmt.Print("|")
			fmt.Printf("%d", cur.Value)

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
func crypto_seed() int64 {
	// https://stackoverflow.com/a/54491783

	var b [8]byte

	if _, err := crypto_rand.Read(b[:]); err != nil {
		panic("Cannot seed")
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}

func time_seed() int64 {
	return time.Now().Unix()
}

func seedRand() {

	// math_rand.Seed(crypto_seed())
	math_rand.Seed(time_seed())
}

// Return a slice of randomly generated distinct integers; all values are incremented
func GetRandomVals() []int {

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
		if val == board.cells[i][c.col].Value {
			return true
		}
		if val == board.cells[c.row][i].Value {
			return true
		}
	}
	return false
}

func (b *Box) ContainsValue(val int) bool {

	for i := 0; i < len(b.cells); i++ {
		if b.cells[i].Value == val {
			return true
		}
	}
	return false

}

// Populate the structs with references to cells;
// Index by Box and Board position to simplify processing and access
func (b *Board) IndexBoard() {

	for i := 0; i < 81; i++ {
		row := i / 9
		col := i % 9
		// Index the box
		b.boxes[row].index = row

		// Get the box position
		ibox := row/3*3 + col/3

		// Store reference to current created cell at this position
		cell := &Cell{row: row, col: col, box: ibox}

		// Create reference to cell in board row/col
		b.cells[row][col] = cell

		// Create reference to cell in appropriate box and position in the boxes index
		b.boxes[ibox].cells[row%3*3+col%3] = cell

	}

}

// Check if the board has been completed
func (b *Board) CheckBoard() bool {

	for i := 0; i < 81; i++ {
		row := i / 9
		col := i % 9
		if b.cells[row][col].Value == 0 {
			// Board incomplete
			return false
		}
	}
	// Board has been completed, trigger recursive exit
	return true
}

func (b *Board) RemoveCells() {
	// Reset val
	n_solutions = 1
	for {

		random_vals := GetRandomVals()
		rRow := random_vals[0] - 1
		rCol := random_vals[1] - 1

		// Get a new cell if the value is 0
		for b.cells[rRow][rCol].Value == 0 {
			random_vals = GetRandomVals()
			rRow = random_vals[0] - 1
			rCol = random_vals[1] - 1
		}

		old_val := b.cells[rRow][rCol].Value
		b.cells[rRow][rCol].Value = 0

		// Either 0 or multiple solutions; invalid board, restore the value and repeat
		if n_solutions != 1 {
			b.cells[rRow][rCol].Value = old_val
		}

		if threshold_met, _ := b.CountZeroCells(); threshold_met {
			b.SolveBoard(false)
			break
		}

	}
}

// Returns true/false if the desired number of cells at 0 is met
func (b *Board) CountZeroCells() (bool, int) {

	zero_cells := 0

	for i := 0; i < 81; i++ {
		row := i / 9
		col := i % 9
		cell := b.cells[row][col]

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
func (b *Board) SolveBoard(fill bool) bool {

	for i := 0; i < 81; i++ {
		row := i / 9
		col := i % 9
		// Refer to current cell
		cell := b.cells[row][col]

		if cell.Value == 0 {

			// Next set of values to attempt placement
			// Iterate and recursively place
			vals := GetRandomVals()
			for _, v := range vals {

				// Value is valid and can be placed
				if !cell.CheckValueConflict(v) && !b.boxes[cell.box].ContainsValue(v) {
					cell.Value = v

					// Board is completed
					if b.CheckBoard() {
						n_solutions++

						if fill {
							return true
						} else {
							break
						}

					} else {
						if b.SolveBoard(fill) {
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

func main() {

	board.IndexBoard()

	board.SolveBoard(true)
	board.RemoveCells()

	board.PrintBoard()

	app := gin.Default()

	app.LoadHTMLGlob("templates/*")

	page := gin.H{
		"title": "Godoku",
		"cells": board.cells}

	app.SetHTMLTemplate(template.Must(template.ParseFiles("templates/header.html", "templates/index.html")))

	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "homeHTML", page)
	})

	app.Run(":3399")

}

func init() {
	seedRand()

}
