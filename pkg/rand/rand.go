package rand

import (
	rand2 "math/rand/v2"
	"time"
)

type Rand struct {
	rnd *rand2.Rand
}

func New() *Rand {
	return &Rand{
		rnd: rand2.New(&seedSource{}),
	}
}

func (r *Rand) Range(min, max int) int {
	return rand2.IntN(max-min) + min
}

type seedSource struct {
	rand2.Source
}

func (s seedSource) Uint64() uint64 {
	return uint64(time.Now().UnixNano())
}
