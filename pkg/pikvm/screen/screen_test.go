package screen

import (
	"fmt"
	"testing"
)

func TestScreenMovements(t *testing.T) {
	testCases := []struct {
		inX    int
		inY    int
		piKvmX int16
		piKvmY int16
	}{
		{0, 0, -32768, -32768},
		{1920, 1080, 32767, 32767},
		{960, 540, 0, 0},
		{480, 270, -16384, -16384},
		{1440, 810, 16383, 16383},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase), func(t *testing.T) {
			p := NewFullHD()
			x, y := p.ToPiKvmPoints(testCase.inX, testCase.inY)
			if x != testCase.piKvmX || y != testCase.piKvmY {
				t.Errorf("expected %d, %d, got %d, %d", testCase.piKvmX, testCase.piKvmY, x, y)
			}
		})
	}
}
