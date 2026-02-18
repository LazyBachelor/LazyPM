package tasks

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func init() {
	task.Register("create_issue", func(svc *service.Services) task.Tasker {
		return NewCreateIssueTask(svc)
	})

	task.Register("coding_task", func(svc *service.Services) task.Tasker {
		return NewCodingTask(svc)
	})
}

func InterfaceToType(it task.Interface) task.InterfaceType {
	switch it.(type) {
	case *repl.REPL:
		return task.InterfaceCLI
	case *tui.Tui:
		return task.InterfaceTUI
	case *web.Web:
		return task.InterfaceWeb
	default:
		return task.InterfaceType("unknown")
	}
}
