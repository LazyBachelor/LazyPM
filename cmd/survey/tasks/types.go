package tasks

import (
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

const (
	InterfaceTUI task.InterfaceType = "tui"
	InterfaceCLI task.InterfaceType = "repl"
	InterfaceWeb task.InterfaceType = "web"
)

func InterfaceToType(it task.Interface) task.InterfaceType {
	switch it.(type) {
	case *repl.REPL:
		return InterfaceCLI
	case *tui.Tui:
		return InterfaceTUI
	case *web.Web:
		return InterfaceWeb
	default:
		return "unknown"
	}
}
