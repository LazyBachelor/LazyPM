package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func main() {
	if err := web.New().Run(context.Background(), models.BaseConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
