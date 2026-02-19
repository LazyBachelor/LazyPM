package tasks

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
	"github.com/charmbracelet/huh"
)

const (
	InterfaceTUI  task.InterfaceType = "tui"
	InterfaceREPL task.InterfaceType = "repl"
	InterfaceWeb  task.InterfaceType = "web"
)

func InterfaceToType(it task.Interface) task.InterfaceType {
	switch it.(type) {
	case *repl.REPL:
		return InterfaceREPL
	case *tui.Tui:
		return InterfaceTUI
	case *web.Web:
		return InterfaceWeb
	default:
		return task.InterfaceType("unknown")
	}
}

func BaseDetails() taskui.TaskDetails {
	return taskui.TaskDetails{
		Title:          "Base Task",
		Description:    "This is a base task.",
		TimeToComplete: "10m",
		Difficulty:     "Easy",
	}
}

func BaseConfig() task.TaskConfig {
	return task.TaskConfig{
		IssuePrefix:           "pm",
		WebAddress:            ":8080",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
	}
}

func ClearIssues(svc *service.Services) error {
	return svc.Beads.DeleteIssues()
}

func BaseQuestions(interfaceType task.InterfaceType) taskui.Questions {
	var taskRating int
	return taskui.Questions{
		huh.NewGroup(
			huh.NewConfirm().
				Title("Did you complete the task?"),
		),
		huh.NewGroup(
			huh.NewSelect[int]().Value(&taskRating).
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

func Question(fields ...huh.Field) *huh.Group {
	return huh.NewGroup(fields...)
}

func ReplQuestion(interfaceType task.InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceREPL {
		return nil
	}
	return huh.NewGroup(fields...)
}

func WebQuestion(interfaceType task.InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceWeb {
		return nil
	}
	return huh.NewGroup(fields...)
}

func TUIQuestion(interfaceType task.InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceTUI {
		return nil
	}
	return huh.NewGroup(fields...)
}
