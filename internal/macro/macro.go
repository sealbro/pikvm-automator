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
	Repeats int
	Events  []Macro
}
