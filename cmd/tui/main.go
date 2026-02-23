package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
)

func main() {
	if err := tui.NewTui().Run(context.Background(), service.BaseConfig); err != nil {
		return
	}
}
