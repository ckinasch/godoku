package main

import (
	"errors"
	"fmt"
	"math/rand"
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

func main() {

}
