package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func main() {
	web := web.NewWeb()

	config := service.Config{
		WebAddress:            "localhost:8080",
		BeadsDBPath:           "./.pm/db.db",
		IssuePrefix:           "pm",
		StatisticsStoragePath: "./.pm/stats.json",
	}

	if err := web.Run(context.Background(), config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
