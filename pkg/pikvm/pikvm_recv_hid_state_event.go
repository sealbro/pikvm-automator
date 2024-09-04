package pikvm

type HIDStateEvent struct {
	Online    bool `json:"online"`
	Busy      bool `json:"busy"`
	Connected any  `json:"connected"`
	Keyboard  struct {
		Online bool `json:"online"`
		Leds   struct {
			Caps   bool `json:"caps"`
			Scroll bool `json:"scroll"`
			Num    bool `json:"num"`
		} `json:"leds"`
		Outputs Outputs `json:"outputs"`
	} `json:"keyboard"`
	Mouse struct {
		Outputs  Outputs `json:"outputs"`
		Online   bool    `json:"online"`
		Absolute bool    `json:"absolute"`
	} `json:"mouse"`
	Jiggler struct {
		Enabled  bool  `json:"enabled"`
		Active   bool  `json:"active"`
		Interval int64 `json:"interval"`
	} `json:"jiggler"`
}

type Outputs struct {
	Available []any  `json:"available"`
	Active    string `json:"active"`
}
