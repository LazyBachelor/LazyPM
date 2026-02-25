package tasks

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/charmbracelet/huh"
)

const codingDescription = `You are tasked with writing a simple function.

This task will test your ability to write clean, working code.

Your task:
1. Review the requirements below
2. Write a function that takes two integers and returns their sum
3. The function should be named "Add"
4. The function should be part of the "coding" package
5. Save your code to the code.txt file

Requirements:
- Function name: Add
- Parameters: two integers
- Return value: integer (sum of the two inputs)
- Package: coding`

var textFileContent = codingDescription + `
Please write your code below this line!
############################################################
`

type CodingTask struct {
	done bool
	app  *service.App
}

func NewCodingTask(app *service.App) *CodingTask {
	return &CodingTask{app: app, done: false}
}

func (t *CodingTask) Config() task.Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/coding-task-stats.json")
}

func (t *CodingTask) Details() taskui.TaskDetails {
	return BaseDetails().WithTitle("Coding Task").WithDescription(codingDescription)
}

func (t *CodingTask) Questions(interfaceType task.InterfaceType) (questions taskui.Questions) {
	return BaseQuestions(interfaceType).
		With(
			ReplQuestion(interfaceType,
				huh.NewConfirm().Title("Question only for REPL interface")),
		).
		With(
			WebQuestion(interfaceType,
				huh.NewInput().Title("Question only for Web interface")),
		).
		With(
			TUIQuestion(interfaceType,
				huh.NewConfirm().Title("Question only for TUI interface")),
		).
		With(
			Question(
				huh.NewConfirm().Title("One last question for all interfaces!")),
		)
}

func (t *CodingTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	if err := os.WriteFile("./code.txt", []byte(textFileContent), 0644); err != nil {
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

	return EndTaskWithTimeout(&t.done, "Task completed!", 5*time.Second)
}
