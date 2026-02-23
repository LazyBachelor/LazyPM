package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func main() {
	if err := web.NewWeb().Run(context.Background(), service.BaseConfig); err != nil {
		return
	}
}
