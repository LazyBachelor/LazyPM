package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli"
)

func main() {
	if err := cli.NewCli().Run(context.Background(), service.BaseConfig); err != nil {
		return
	}
}
