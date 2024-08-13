package macro

import (
	"github.com/sealbro/pikvm-automator/pkg/pikvm"
	"time"
)

type Macro interface {
}

type Delay struct {
	Time time.Duration
}

type KeyEvent struct {
	Key   pikvm.Key
	State bool
}

type MouseEvent struct {
	X int
	Y int
}

type Repeat struct {
	Group
	Repeats int
}

type Bind struct {
	Group
}

type Group struct {
	Events []Macro
}

func (g *Group) TotalDelay() time.Duration {
	var total time.Duration
	for _, m := range g.Events {
		switch v := m.(type) {
		case Delay:
			total += v.Time
		case Repeat:
			total += time.Duration(v.Repeats) * v.TotalDelay()
		case Bind:
			total += v.TotalDelay()
		case Group:
			total += v.TotalDelay()
		}
	}
	return total
}
