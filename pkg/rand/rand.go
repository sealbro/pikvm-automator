package rand

import (
	"math/rand/v2"
	"time"
)

type Rand struct {
	rnd *rand.Rand
}

func New() *Rand {
	return &Rand{
		rnd: rand.New(&seedSource{}),
	}
}

func (r *Rand) Range(min, max int) int {
	return rand.IntN(max-min) + min
}

type seedSource struct {
	rand.Source
}

func (s seedSource) Uint64() uint64 {
	return uint64(time.Now().UnixNano())
}
