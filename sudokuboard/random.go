package sudokuboard

import (
	"math/rand"
	"time"
)

// Return a slice of randomly generated distinct integers; all values are incremented
func getRandomVals() []int {

	vals := rand.Perm(9)
	for i := range vals {
		vals[i]++
	}
	return vals

}

func seedRand() {

	rand.Seed(time.Now().Unix())
}


