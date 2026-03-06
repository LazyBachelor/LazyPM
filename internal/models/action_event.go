package models

import "encoding/json"

type ActionEvent struct {
	Source string `json:"source"`
	Action string `json:"action"`
	Target string `json:"target,omitempty"`
	Result string `json:"result,omitempty"`
}

func EncodeActionEvent(event ActionEvent) string {
	bytes, err := json.Marshal(event)
	if err != nil {
		return event.Action
	}
	return string(bytes)
}

func DecodeActionEvent(raw string) (ActionEvent, bool) {
	var event ActionEvent
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		return ActionEvent{}, false
	}
	if event.Source == "" && event.Action == "" && event.Target == "" && event.Result == "" {
		return ActionEvent{}, false
	}
	return event, true
}
