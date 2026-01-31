package main

import (
	"beadstest/internal/models"
	"beadstest/internal/storage"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	statStore := storage.NewJsonStorage("./.pm/test.json", &models.Statistics{
		ID:        uuid.New(),
		StartTime: time.Now(),
	})

	if err := statStore.Init(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
