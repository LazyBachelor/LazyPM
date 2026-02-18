package tasks

import (
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
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
		WebAddress:            ":8080",
	}
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

func AppendGroup(questions *taskui.Questions, group *huh.Group) taskui.Questions {
	*questions = append(*questions, group)
	return *questions
}

func AppendQuestion(questions *taskui.Questions, interfaceType task.InterfaceType, field ...huh.Field) taskui.Questions {
	*questions = append(*questions, huh.NewGroup(field...))
	return *questions
}

func AppendReplQuestion(questions *taskui.Questions, interfaceType task.InterfaceType, field ...huh.Field) taskui.Questions {
	if interfaceType != InterfaceREPL {
		return *questions
	}
	*questions = append(*questions, huh.NewGroup(field...))
	return *questions
}

func AppendWebQuestion(questions *taskui.Questions, interfaceType task.InterfaceType, field ...huh.Field) taskui.Questions {
	if interfaceType != InterfaceWeb {
		return *questions
	}
	*questions = append(*questions, huh.NewGroup(field...))
	return *questions
}

func AppendTUIQuestion(questions *taskui.Questions, interfaceType task.InterfaceType, field ...huh.Field) taskui.Questions {
	if interfaceType != InterfaceTUI {
		return *questions
	}
	*questions = append(*questions, huh.NewGroup(field...))
	return *questions
}
