package models

import (
	"time"
)

type Statistics struct {
	ID        int           `json:"id"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`

	InterfaceType  InterfaceType `json:"interface_type"`
	TasksCompleted int           `json:"tasks_completed"`
	ButtonClicks   ButtonClicks  `json:"button_clicks"`
}

type ButtonClicks struct {
	Clicks int `json:"clicks"`
}
