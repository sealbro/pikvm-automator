package pikvm

type PiKVMRecvEvent struct {
	EventType string `json:"event_type"`
	Event     any    `json:"event"`
}
