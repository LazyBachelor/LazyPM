package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/charmbracelet/huh"
)

const codingDescription = `You are tasked with writing a simple function.

Write a function that takes two integers and returns their sum.
The function should be named "Add" and be part of the "coding" package.`

type CodingTask struct {
	svc *service.Services
}

func NewCodingTask(svc *service.Services) *CodingTask {
	return &CodingTask{svc: svc}
}

func (t *CodingTask) Config() task.TaskConfig {
	return task.TaskConfig{
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/coding-stats.json",
		WebAddress:            ":8080",
	}
}

func (t *CodingTask) Details() taskui.TaskDetails {
	return taskui.TaskDetails{
		Title:          "Coding Task",
		Description:    codingDescription,
		TimeToComplete: "10m",
		Difficulty:     "Easy",
	}
}

func (t *CodingTask) Questions(interfaceType task.InterfaceType) taskui.Questions {
	return taskui.Questions{
		huh.NewGroup(huh.NewConfirm().Title("Did you complete the coding task?")),
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(
					huh.NewOption("Very easy", 1),
					huh.NewOption("Easy", 2),
					huh.NewOption("Moderate", 3),
					huh.NewOption("Hard", 4),
				).
				Title("How difficult was the task?"),
		),
	}
}

func (t *CodingTask) Setup(ctx context.Context) error {
	return nil
}

func (t *CodingTask) Validate(ctx context.Context) (bool, error) {
	return true, nil
}
