package main

import (
	"fmt"
	"github.com/sealbro/pikvm-automator/internal/macro"
	"github.com/sealbro/pikvm-automator/pkg/pikvm"
	"strconv"
	"time"
	"unicode"
)

func main() {
	expressions := []string{
		"1ms",
		"2s",
		"MetaLeft",
		"MetaLeft+KeyD",
		"MetaLeft+KeyD+1ms",
		"[10](MetaLeft+KeyD+1ms)",
		"[10](MetaLeft+KeyD+1ms|10s)",
		"[10](MetaLeft+KeyD+1ms|10s)|MetaLeft|5s",
	}

	for _, exp := range expressions {
		macros := parse(exp)
		fmt.Println(macros)
	}
}

func parse(expression string) []macro.Macro {
	const (
		empty       = ""
		combinator  = '+'
		splitter    = '|'
		startGroup  = '('
		endGroup    = ')'
		startRepeat = '['
		endRepeat   = ']'
	)

	isDelay := false
	isKeyEvent := false
	//isMouseEvent := false
	isRepeat := false

	var result []macro.Macro
	var current = empty
	skipLetter := 0

	finishKeyEvent := func() {
		var tempResult []macro.Macro

		for i := len(result) - 1; i >= 0; i-- {
			if m, ok := result[i].(macro.KeyEvent); ok {
				if m.State == false {
					tempResult = append(tempResult, macro.KeyEvent{Key: m.Key, State: true})
				} else {
					break
				}
			}
		}

		if len(current) > 0 {
			result = append(result, macro.KeyEvent{Key: pikvm.Key(current), State: false})
			result = append(result, macro.KeyEvent{Key: pikvm.Key(current), State: true})
		}
		result = append(result, tempResult...)
		current = empty
		isKeyEvent = false
	}

	finishGroup := func() {
		var repeatGroup macro.Repeat
		startIndex := -1
		for i := len(result) - 1; i >= 0; i-- {
			if m, ok := result[i].(macro.Repeat); ok {
				startIndex = i + 1
				repeatGroup = m
				break
			}
		}

		if startIndex >= 0 {
			repeatGroup.Events = result[startIndex:]
			result = result[:startIndex-1]
			result = append(result, repeatGroup)
		}
	}

	for _, letter := range expression {
		if skipLetter > 0 {
			skipLetter--
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
				result = append(result, macro.Delay{Time: delay})
				current = empty
				isDelay = false
			}
			continue
		}

		if isKeyEvent {
			if letter == combinator {
				result = append(result, macro.KeyEvent{Key: pikvm.Key(current), State: false})
				current = empty
				isKeyEvent = false
			} else if letter == splitter || letter == endGroup {
				finishKeyEvent()
				current = empty
				isKeyEvent = false
			} else {
				current += string(letter)
			}
			continue
		}

		if isRepeat {
			if unicode.IsDigit(letter) {
				isRepeat = true
				current += string(letter)
			}

			if letter == endRepeat {
				atoi, _ := strconv.Atoi(current)
				result = append(result, macro.Repeat{Repeats: atoi})
				current = empty
				isRepeat = false
				skipLetter = 1 // skip '('
			}
			continue
		}

		//if isMouseEvent {
		//
		//}

		if letter == endGroup {
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

	finishKeyEvent()

	return result
}
