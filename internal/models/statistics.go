package models

import (
	"time"

	"github.com/google/uuid"
)

type InterfaceType string

const (
	InterfaceTypeCLI InterfaceType = "CLI"
	InterfaceTypeWeb InterfaceType = "Web"
)

type Statistics struct {
	ID        uuid.UUID     `json:"id"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`

	InterfaceType  InterfaceType `json:"interface_type"`
	TasksCompleted int           `json:"tasks_completed"`
}
