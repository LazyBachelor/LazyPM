package forms

import (
	"embed"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
)

type Introduction struct {
	Title       string
	Description string
	About       struct {
		Title       string
		Description string
	}
	Disclaimer struct {
		Title       string
		Description string
	}
}

func NewIntroduction(fs embed.FS) (Introduction, error) {
	var intro Introduction
	if _, err := toml.DecodeFS(fs, "assets/intro.toml", &intro); err != nil {
		return Introduction{}, err
	}

	return intro, nil
}

func (i Introduction) Run() error {
	fmt.Print("\033[H\033[2J")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(i.Title).
				Description(i.Description),
		),
		huh.NewGroup(
			huh.NewNote().
				Title(i.About.Title).
				Description(i.About.Description),
		),
		huh.NewGroup(
			huh.NewNote().
				Title(i.Disclaimer.Title).
				Description(i.Disclaimer.Description),
		),
	).WithLayout(huh.LayoutStack).
		WithTheme(huh.ThemeBase16())

	return form.Run()
}
