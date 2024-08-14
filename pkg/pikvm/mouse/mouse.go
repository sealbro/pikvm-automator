package mouse

type Button string

var Buttons = []Button{Left, Right}

const (
	Left  Button = "left"
	Right Button = "right"
)
