package app

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

type Initializer interface {
	Init(path string) error
}

type InteractiveInitializer struct{}

func (i InteractiveInitializer) Init(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return nil
	}

	var initialize bool

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("PM is not initialized in this directory!").
				Description("Do you want to initialize it here?").
				Value(&initialize),
		)).WithTheme(huh.ThemeBase16()).WithAccessible(true).Run()

	if err != nil {
		return err
	}

	if !initialize {
		return fmt.Errorf("project not initialized")
	}

	return nil
}
