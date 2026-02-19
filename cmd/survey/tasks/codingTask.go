package tasks

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
)

const codingDescription = `You are tasked with writing a simple function.

Write a function that takes two integers and returns their sum.
The function should be named "Add" and be part of the "coding" package.`

var textFileContent = codingDescription + `
Please write your code below this line!
############################################################
`

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
	if err := t.svc.DeleteIssues(); err != nil {
		return err
	}

	content := textFileContent
	
	if err := os.WriteFile("./code.txt", []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func (t *CodingTask) Validate(ctx context.Context) (bool, error) {

	file, err := os.ReadFile("./code.txt")
	if err != nil {
		return false, err
	}

	if string(file) == "" {
		return false, fmt.Errorf("the file is empty")
	}

	code, ok := strings.CutPrefix(string(file), textFileContent)
	if !ok {
		return false, fmt.Errorf("the file content is not in the expected format")
	}

	code = strings.TrimSpace(code)

	if !strings.Contains(code, "package coding") {
		return false, fmt.Errorf("the code does not belong to the 'coding' package")
	}

	if !strings.Contains(code, "func Add") {
		return false, fmt.Errorf("the function 'Add' is not defined")
	}

	if !strings.Contains(code, "return") {
		return false, fmt.Errorf("the function does not contain a return statement")
	}

	return true, nil
}
