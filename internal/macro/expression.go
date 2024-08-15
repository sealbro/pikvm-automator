package macro

import (
	"github.com/sealbro/pikvm-automator/pkg/pikvm/keyboard"
	"github.com/sealbro/pikvm-automator/pkg/pikvm/mouse"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Key [\w\d]+
// MouseMove @\d+\'\d+
// Delay \d+[sm]
// Combinator \+
// Splitter \|
// StartGroup \(
// EndGroup \)
// StartRepeat \[
// EndRepeat \]

const (
	empty         = ""
	startGroup    = '('
	endGroup      = ')'
	startRepeat   = '['
	endRepeat     = ']'
	mousePrefix   = '@'
	mouseSplitter = '\''
	combinator    = '+'
	splitter      = '|'
)

type Expression struct {
	events []Macro
	exp    string
}

func New(exp string) *Expression {
	return &Expression{
		exp: exp,
	}
}

func (e *Expression) Parse() Group {
	if len(e.events) > 0 {
		return Group{Events: append(make([]Macro, 0), e.events...)}
	}

	isDelay := false
	isKeyEvent := false
	isMouseEvent := false
	isRepeat := false

	var result []Macro
	var current = empty
	skipLetter := 0

	finishGroup := func() {
		var repeatGroup Repeat
		startIndex := -1
		for i := len(result) - 1; i >= 0; i-- {
			if m, ok := result[i].(Repeat); ok && len(m.Events) == 0 {
				startIndex = i + 1
				repeatGroup = m
				break
			}
		}

		if startIndex >= 0 {
			repeatGroup.Events = append(repeatGroup.Events, result[startIndex:]...)
			result = append(make([]Macro, 0), result[:startIndex-1]...)
			result = append(result, repeatGroup)
		}
	}

	finishBind := func() {
		var tempResult []Macro

		var startIndex = -1
		for i := len(result) - 1; i >= 0; i-- {
			if m, ok := result[i].(MouseClickEvent); ok {
				if m.State {
					startIndex = i
					tempResult = append(tempResult, MouseClickEvent{Button: m.Button, State: false})
				} else {
					break
				}
			}
			if m, ok := result[i].(KeyPressEvent); ok {
				if m.State {
					startIndex = i
					tempResult = append(tempResult, KeyPressEvent{Key: m.Key, State: false})
				} else {
					break
				}
			}
		}

		if len(current) > 0 && slices.Contains(keyboard.Keys, keyboard.Key(current)) {
			result = append(result, KeyPressEvent{Key: keyboard.Key(current), State: true})
			result = append(result, KeyPressEvent{Key: keyboard.Key(current), State: false})
		}

		if len(tempResult) > 0 {
			var bindMacro Bind
			expectedBind := result[startIndex:]
			if len(expectedBind) > 1 {
				bindMacro.Events = append(bindMacro.Events, expectedBind...)
				bindMacro.Events = append(bindMacro.Events, tempResult...)

				result = append(make([]Macro, 0), result[:startIndex]...)
				result = append(result, bindMacro)
			} else {
				result = append(result, tempResult...)
			}
		}

		current = empty
		isMouseEvent = false
		isKeyEvent = false
	}

	finishMouseMoveEvent := func() {
		if len(current) > 0 && strings.Contains(current, string(mouseSplitter)) {
			split := strings.Split(current, string(mouseSplitter))
			xtoi, _ := strconv.Atoi(split[0])
			ytoi, _ := strconv.Atoi(split[1])
			result = append(result, MouseMoveEvent{X: xtoi, Y: ytoi})
			current = empty
		}

		isMouseEvent = false
	}

	for _, letter := range e.exp {
		if skipLetter > 0 {
			skipLetter--
			continue
		}

		if letter == splitter {
			finishMouseMoveEvent()
			finishBind()
			continue
		}

		if isDelay {
			if unicode.IsDigit(letter) {
				current += string(letter)
			} else {
				atoi, _ := strconv.Atoi(current)
				var delay time.Duration
				if letter == 's' {
					delay = time.Second * time.Duration(atoi)
				} else if letter == 'm' {
					delay = time.Millisecond * time.Duration(atoi)
					skipLetter = 1 // skip 's'
				}
				result = append(result, Delay{Time: delay})
				current = empty
				isDelay = false
			}
			continue
		}

		if isRepeat {
			if unicode.IsDigit(letter) {
				isRepeat = true
				current += string(letter)
			} else if letter == endRepeat {
				atoi, _ := strconv.Atoi(current)
				result = append(result, Repeat{Repeats: atoi})
				current = empty
				isRepeat = false
				skipLetter = 1 // skip '(' - startGroup
			}
			continue
		}

		if isKeyEvent {
			if letter == combinator {
				result = append(result, KeyPressEvent{Key: keyboard.Key(current), State: true})
				current = empty
				isKeyEvent = false
			} else {
				current += string(letter)
			}
			continue
		}

		if isMouseEvent {
			if letter == combinator {
				finishMouseMoveEvent()
				current = empty
				isMouseEvent = false
			} else if letter == 'l' {
				result = append(result, MouseClickEvent{Button: mouse.Left, State: true})
				isMouseEvent = false
				skipLetter = 3 // skip 'eft'
			} else if letter == 'r' {
				result = append(result, MouseClickEvent{Button: mouse.Right, State: true})
				isMouseEvent = false
				skipLetter = 4 // skip 'ight'
			} else if unicode.IsDigit(letter) || letter == mouseSplitter {
				isMouseEvent = true
				current += string(letter)
			}
			continue
		}

		if letter == mousePrefix {
			isMouseEvent = true
			continue
		}

		if letter == endGroup {
			finishMouseMoveEvent()
			finishBind()
			finishGroup()
			continue
		}

		if letter == startRepeat {
			isRepeat = true
			continue
		}

		if unicode.IsDigit(letter) {
			isDelay = true
			current = string(letter)
			continue
		}

		if unicode.IsUpper(letter) {
			isKeyEvent = true
			current = string(letter)
			continue
		}
	}

	finishMouseMoveEvent()
	finishBind()

	e.events = result
	return Group{Events: append(make([]Macro, 0), e.events...)}
}

func (e *Expression) String() string {
	return compile(e.events)
}

func compile(events []Macro) string {
	if len(events) == 0 {
		return empty
	}

	sb := strings.Builder{}

	for i, m := range events {
		switch v := m.(type) {
		case Delay:
			milliseconds := v.Time.Milliseconds()
			if milliseconds%1000 == 0 {
				sb.WriteString(strconv.Itoa(int(v.Time.Seconds())))
				sb.WriteString("s")
				break
			} else {
				sb.WriteString(strconv.Itoa(int(milliseconds)))
				sb.WriteString("ms")
			}
		case KeyPressEvent:
			if v.State {
				sb.WriteString(string(v.Key))
			} else {
				// skip splitter
				continue
			}
		case MouseMoveEvent:
			sb.WriteString(string(mousePrefix))
			sb.WriteString(strconv.Itoa(v.X))
			sb.WriteString(string(mouseSplitter))
			sb.WriteString(strconv.Itoa(v.Y))
		case MouseClickEvent:
			if v.State {
				sb.WriteString(string(mousePrefix))
				sb.WriteString(string(v.Button))
			} else {
				// skip splitter
				continue
			}
		case Repeat:
			sb.WriteString(string(startRepeat))
			sb.WriteString(strconv.Itoa(v.Repeats))
			sb.WriteString(string(endRepeat))
			sb.WriteString(string(startGroup))
			sb.WriteString(compile(v.Events))
			sb.WriteString(string(endGroup))
		case Bind:
			sb.WriteString(strings.ReplaceAll(compile(v.Events), string(splitter), string(combinator)))
		}

		if i < len(events)-1 {
			sb.WriteString(string(splitter))
		}
	}

	s := sb.String()
	if s[len(s)-1] == splitter {
		s = s[:len(s)-1]
	}
	return s
}
