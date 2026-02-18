package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/storage"
	"github.com/charmbracelet/huh"

	"github.com/google/uuid"
	"github.com/steveyegge/beads"
)

type Config struct {
	RootCmd               string
	WebAddress            string
	BeadsDBPath           string
	IssuePrefix           string
	StatisticsStoragePath string
}

type Services struct {
	Config     Config
	DB         *sql.DB
	Beads      *BeadsService
	Statistics *StatisticsService
}

func NewServices(ctx context.Context, config Config) (*Services, func(), error) {
	var cleanupFuncs []func()

	if !initialized(config.BeadsDBPath) {
		fmt.Println("PM is not initialized")
		os.Exit(0)
	}

	db, err := sql.Open("sqlite3", config.BeadsDBPath)
	if err != nil {
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, func() { db.Close() })

	store, err := beads.NewSQLiteStorage(ctx, config.BeadsDBPath)
	if err != nil {
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, func() { store.Close() })

	beadsSvc, err := NewBeadsService(ctx, store, config.IssuePrefix)
	if err != nil {
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, func() { beadsSvc.Close() })

	statStore := storage.NewJsonStorage(config.StatisticsStoragePath, &models.Statistics{
		ID:        uuid.New(),
		StartTime: time.Now(),
	})

	statSvc, err := NewStatisticsService(statStore)
	if err != nil {
		return nil, nil, err
	}

	return &Services{
		DB:         db,
		Beads:      beadsSvc,
		Statistics: statSvc,
		Config:     config,
	}, func() { runCleanup(cleanupFuncs) }, nil

}

func runCleanup(funcs []func()) {
	for _, fn := range funcs {
		fn()
	}
}

func initialized(beadsPath string) bool {
	_, err := os.Stat(beadsPath)

	if os.IsNotExist(err) {
		var initialize bool

		huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Title("PM is not initialized in this directory!").
					Description("Do you want to initialize it here?").
					Value(&initialize),
			),
		).WithTheme(huh.ThemeBase16()).WithAccessible(true).Run()

		if !initialize {
			return false
		}
	}
	return true
}

// Used for reset/testing/admin operations
func (s Services) DeleteIssues() error {
	var deleteIssues = "DELETE FROM issues;"

	if _, err := s.DB.Exec(deleteIssues); err != nil {
		return err
	}
	return nil
}
