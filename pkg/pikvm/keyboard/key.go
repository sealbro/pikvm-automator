package keyboard

// Key see official docs https://github.com/pikvm/kvmd/blob/d7675cdf231896851c492e1e0e719d2ccdc28082/hid/arduino/lib/drivers/usb-keymap.h
type Key string

var Keys = []Key{
	KeyA,
	KeyB,
	KeyC,
	KeyD,
	KeyE,
	KeyF,
	KeyG,
	KeyH,
	KeyI,
	KeyJ,
	KeyK,
	KeyL,
	KeyM,
	KeyN,
	KeyO,
	KeyP,
	KeyQ,
	KeyR,
	KeyS,
	KeyT,
	KeyU,
	KeyV,
	KeyW,
	KeyX,
	KeyY,
	KeyZ,
	Digit1,
	Digit2,
	Digit3,
	Digit4,
	Digit5,
	Digit6,
	Digit7,
	Digit8,
	Digit9,
	Digit0,
	Enter,
	Escape,
	Backspace,
	Tab,
	Space,
	Minus,
	Equal,
	BracketLeft,
	BracketRight,
	Backslash,
	Semicolon,
	Quote,
	Backquote,
	Comma,
	Period,
	Slash,
	CapsLock,
	F1,
	F2,
	F3,
	F4,
	F5,
	F6,
	F7,
	F8,
	F9,
	F10,
	F11,
	F12,
	PrintScreen,
	Insert,
	Home,
	PageUp,
	Delete,
	End,
	PageDown,
	ArrowRight,
	ArrowLeft,
	ArrowDown,
	ArrowUp,
	ControlLeft,
	ShiftLeft,
	AltLeft,
	MetaLeft,
	ControlRight,
	ShiftRight,
	AltRight,
	MetaRight,
	Pause,
	ScrollLock,
	NumLock,
	ContextMenu,
	NumpadDivide,
	NumpadMultiply,
	NumpadSubtract,
	NumpadAdd,
	NumpadEnter,
	Numpad1,
	Numpad2,
	Numpad3,
	Numpad4,
	Numpad5,
	Numpad6,
	Numpad7,
	Numpad8,
	Numpad9,
	Numpad0,
	NumpadDecimal,
	Power,
	IntlBackslash,
	IntlYen,
	IntlRo,
	KanaMode,
	Convert,
	NonConvert,
}

const (
	KeyA           Key = "KeyA"
	KeyB           Key = "KeyB"
	KeyC           Key = "KeyC"
	KeyD           Key = "KeyD"
	KeyE           Key = "KeyE"
	KeyF           Key = "KeyF"
	KeyG           Key = "KeyG"
	KeyH           Key = "KeyH"
	KeyI           Key = "KeyI"
	KeyJ           Key = "KeyJ"
	KeyK           Key = "KeyK"
	KeyL           Key = "KeyL"
	KeyM           Key = "KeyM"
	KeyN           Key = "KeyN"
	KeyO           Key = "KeyO"
	KeyP           Key = "KeyP"
	KeyQ           Key = "KeyQ"
	KeyR           Key = "KeyR"
	KeyS           Key = "KeyS"
	KeyT           Key = "KeyT"
	KeyU           Key = "KeyU"
	KeyV           Key = "KeyV"
	KeyW           Key = "KeyW"
	KeyX           Key = "KeyX"
	KeyY           Key = "KeyY"
	KeyZ           Key = "KeyZ"
	Digit1         Key = "Digit1"
	Digit2         Key = "Digit2"
	Digit3         Key = "Digit3"
	Digit4         Key = "Digit4"
	Digit5         Key = "Digit5"
	Digit6         Key = "Digit6"
	Digit7         Key = "Digit7"
	Digit8         Key = "Digit8"
	Digit9         Key = "Digit9"
	Digit0         Key = "Digit0"
	Enter          Key = "Enter"
	Escape         Key = "Escape"
	Backspace      Key = "Backspace"
	Tab            Key = "Tab"
	Space          Key = "Space"
	Minus          Key = "Minus"
	Equal          Key = "Equal"
	BracketLeft    Key = "BracketLeft"
	BracketRight   Key = "BracketRight"
	Backslash      Key = "Backslash"
	Semicolon      Key = "Semicolon"
	Quote          Key = "Quote"
	Backquote      Key = "Backquote"
	Comma          Key = "Comma"
	Period         Key = "Period"
	Slash          Key = "Slash"
	CapsLock       Key = "CapsLock"
	F1             Key = "F1"
	F2             Key = "F2"
	F3             Key = "F3"
	F4             Key = "F4"
	F5             Key = "F5"
	F6             Key = "F6"
	F7             Key = "F7"
	F8             Key = "F8"
	F9             Key = "F9"
	F10            Key = "F10"
	F11            Key = "F11"
	F12            Key = "F12"
	PrintScreen    Key = "PrintScreen"
	Insert         Key = "Insert"
	Home           Key = "Home"
	PageUp         Key = "PageUp"
	Delete         Key = "Delete"
	End            Key = "End"
	PageDown       Key = "PageDown"
	ArrowRight     Key = "ArrowRight"
	ArrowLeft      Key = "ArrowLeft"
	ArrowDown      Key = "ArrowDown"
	ArrowUp        Key = "ArrowUp"
	ControlLeft    Key = "ControlLeft"
	ShiftLeft      Key = "ShiftLeft"
	AltLeft        Key = "AltLeft"
	MetaLeft       Key = "MetaLeft"
	ControlRight   Key = "ControlRight"
	ShiftRight     Key = "ShiftRight"
	AltRight       Key = "AltRight"
	MetaRight      Key = "MetaRight"
	Pause          Key = "Pause"
	ScrollLock     Key = "ScrollLock"
	NumLock        Key = "NumLock"
	ContextMenu    Key = "ContextMenu"
	NumpadDivide   Key = "NumpadDivide"
	NumpadMultiply Key = "NumpadMultiply"
	NumpadSubtract Key = "NumpadSubtract"
	NumpadAdd      Key = "NumpadAdd"
	NumpadEnter    Key = "NumpadEnter"
	Numpad1        Key = "Numpad1"
	Numpad2        Key = "Numpad2"
	Numpad3        Key = "Numpad3"
	Numpad4        Key = "Numpad4"
	Numpad5        Key = "Numpad5"
	Numpad6        Key = "Numpad6"
	Numpad7        Key = "Numpad7"
	Numpad8        Key = "Numpad8"
	Numpad9        Key = "Numpad9"
	Numpad0        Key = "Numpad0"
	NumpadDecimal  Key = "NumpadDecimal"
	Power          Key = "Power"
	IntlBackslash  Key = "IntlBackslash"
	IntlYen        Key = "IntlYen"
	IntlRo         Key = "IntlRo"
	KanaMode       Key = "KanaMode"
	Convert        Key = "Convert"
	NonConvert     Key = "NonConvert"
)
