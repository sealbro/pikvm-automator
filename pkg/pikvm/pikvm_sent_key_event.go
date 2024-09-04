package pikvm

import (
	"github.com/sealbro/pikvm-automator/pkg/pikvm/keyboard"
	"github.com/sealbro/pikvm-automator/pkg/pikvm/mouse"
)

// https://github.com/jtyers/pikvm-mouse-wiggle/blob/main/pikvm_mouse_wiggle/wiggle.py
type EventType string

const (
	Ping        EventType = "ping"
	Pong        EventType = "pong"
	Keyboard    EventType = "key"
	MouseMove   EventType = "mouse_move"
	MouseButton EventType = "mouse_button"
)

type PiKvmEvent struct {
	EventType EventType `json:"event_type"`
	Event     any       `json:"event"`
}

type KeyboardEvent struct {
	Key   keyboard.Key `json:"key"`
	State bool         `json:"state"`
}

type MouseMoveEvent struct {
	To MousePoint `json:"to"`
}

type MousePoint struct {
	X int16 `json:"x"`
	Y int16 `json:"y"`
}

type MouseButtonEvent struct {
	Button mouse.Button `json:"button"`
	State  bool         `json:"state"`
}
