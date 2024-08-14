package screen

import "math"

type Position struct {
	ratioX int
	ratioY int
	width  int
	height int
}

func NewFullHD() *Position {
	return NewScreen(1920, 1080)
}

//func NewQuadHD() *Position {
//	return NewScreen(2560, 1440)
//}

//func NewUltraHD() *Position {
//	return NewScreen(3840, 2160)
//}

func NewScreen(width, height int) *Position {
	return &Position{
		width:  width,
		height: height,
	}
}

func (p *Position) ToPiKvmPoints(inX, inY int) (int16, int16) {
	const (
		oX = math.MinInt16
		oY = math.MinInt16
		mX = math.MaxInt16
		mY = math.MaxInt16
	)

	x := int16(mapRange(float64(inX), 0, float64(p.width), float64(oX), float64(mX)))
	y := int16(mapRange(float64(inY), 0, float64(p.height), float64(oY), float64(mY)))

	return x, y
}

func mapRange(value, inMin, inMax, outMin, outMax float64) float64 {
	return (value-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}
