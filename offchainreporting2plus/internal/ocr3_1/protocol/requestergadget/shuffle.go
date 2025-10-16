package requestergadget

import (
	"math"
	"math/rand"

	"github.com/RoSpaceDev/libocr/commontypes"
)

func shuffle(ranks []commontypes.OracleID) []commontypes.OracleID {

	return softmaxShuffle(ranks, math.Ln2)
}

func softmaxShuffle[T any](ranks []T, alpha float64) []T {
	n := len(ranks)
	// precompute weights
	weights := make([]float64, n)
	for i := range n {
		weights[i] = math.Exp(-alpha * float64(i))
	}

	shuffled := make([]T, 0, n)
	used := make([]bool, n)

	for len(shuffled) < n {
		// compute total weight of unused
		total := 0.0
		for i, w := range weights {
			if !used[i] {
				total += w
			}
		}

		// sample
		r := rand.Float64() * total
		acc := 0.0
		for i, w := range weights {
			if used[i] {
				continue
			}
			acc += w
			if r <= acc {
				used[i] = true
				shuffled = append(shuffled, ranks[i])
				break
			}
		}
	}

	return shuffled
}
