package main

import (
	"beadstest/cmd/cli/commands"
	"beadstest/internal/service"
	"context"
	"fmt"
	"os"

	"github.com/steveyegge/beads"
)

func main() {
	ctx := context.Background()

	store, err := beads.NewSQLiteStorage(ctx, "./db.db")
	handleError(err)
	defer store.Close()

	svc, err := service.NewService(ctx, store, "pm")
	handleError(err)
	defer svc.Close()

	handleError(commands.Execute(svc))
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
