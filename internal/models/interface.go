package models

import "context"

// Interface represents a user interface for interacting with the application.
// It can be any implementation that fulfills the Run method.
type Interface interface {
	Run(context.Context, Config) error
}

// InterfaceType represents the type of user interface
type InterfaceType string

// Interface represents a user interface for interacting with the application.
// Keep lower case to avoid confusion with Go's built-in interfaces.
const (
	InterfaceTypeCLI  InterfaceType = "cli"
	InterfaceTypeTUI  InterfaceType = "tui"
	InterfaceTypeWeb  InterfaceType = "web"
	InterfaceTypeREPL InterfaceType = "repl"
)
