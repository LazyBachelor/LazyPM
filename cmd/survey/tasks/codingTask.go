package tasks

import (
	"context"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/internal/utils"
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

func (t *CodingTask) Validate(ctx context.Context) ValidationFeedback {
	expect := utils.NewExpector()
	expect.Assert(true, "This task is always valid")
	return expect.Complete()
}
