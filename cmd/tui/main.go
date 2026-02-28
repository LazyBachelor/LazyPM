package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
)

func main() {
	if err := tui.NewTui().Run(context.Background(), models.BaseConfig); err != nil {
		return
	}
}
