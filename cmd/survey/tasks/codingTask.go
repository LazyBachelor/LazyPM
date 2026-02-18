package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
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
	config := BaseConfig()
	config.StatisticsStoragePath = "./.pm/coding-task-stats.json"
	return config
}

func (t *CodingTask) Details() taskui.TaskDetails {
	details := BaseDetails()
	details.Title = "Coding Task"
	details.Description = codingDescription
	return details
}

func (t *CodingTask) Questions(interfaceType task.InterfaceType) (questions taskui.Questions) {
	questions = append(questions, BaseQuestions(interfaceType)...)

	return questions
}

func (t *CodingTask) Setup(ctx context.Context) error {
	return nil
}

func (t *CodingTask) Validate(ctx context.Context) (bool, error) {
	return true, nil
}
