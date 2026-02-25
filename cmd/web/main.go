package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func main() {
	if err := web.NewWeb().Run(context.Background(), service.BaseConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
