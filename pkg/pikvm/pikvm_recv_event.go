package pikvm

import (
	"encoding/json"
	"fmt"
)

// RecvState see https://docs.pikvm.org/api/#websocket-events
type RecvState string

const (
	GpioModelState    RecvState = "gpio_model_state"
	InfoFanState      RecvState = "info_fan_state"
	InfoExtrasState   RecvState = "info_extras_state"
	InfoHardwareState RecvState = "info_hw_state"
	InfoMetaState     RecvState = "info_meta_state"
	InfoSystemState   RecvState = "info_system_state"
	InfoWolState      RecvState = "wol_state"
	GpioState         RecvState = "gpio_state"
	HidState          RecvState = "hid_state"
	AtxState          RecvState = "atx_state"
	MsdState          RecvState = "msd_state"
	StreamState       RecvState = "streamer_state"
	LoopState         RecvState = "loop"
	PingState         RecvState = "ping"
	PongState         RecvState = "pong"
)

type PiKVMRecvEvent struct {
	EventType RecvState `json:"event_type"`
	Event     any       `json:"event"`
}

func (ce *PiKVMRecvEvent) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var eventType RecvState

	err = json.Unmarshal(*objMap["event_type"], &eventType)
	if err != nil {
		return fmt.Errorf("event_type is not a string: %w", err)
	}

	ce.EventType = eventType

	switch eventType {
	case HidState:
		state := HIDStateEvent{}
		err = json.Unmarshal(*objMap["event"], &state)
		if err != nil {
			return fmt.Errorf("hid_state event is not a HIDStateEvent: %w", err)
		}
		ce.Event = state
	default:
		_ = json.Unmarshal(*objMap["event"], &ce.Event)
	}

	return nil
}
