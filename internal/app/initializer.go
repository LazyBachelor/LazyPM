package app

import (
	"fmt"
	"os"

	"charm.land/huh/v2"
)

type Initializer interface {
	Init(path string) error
}

type InteractiveInitializer struct{}

func (i InteractiveInitializer) Init(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
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
