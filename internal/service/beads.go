package service

import (
	"context"
	"fmt"

	"github.com/steveyegge/beads"
)

type Service struct {
	beads.Storage
}

func NewService(ctx context.Context, storage beads.Storage, prefix string) (*Service, error) {
	issue_prefix, err := storage.GetConfig(ctx, "issue_prefix")
	if err != nil || issue_prefix == "" {
		if err := storage.SetConfig(ctx, "issue_prefix", prefix); err != nil {
			return nil, fmt.Errorf("failed to set issue_prefix: %w", err)
		}
		fmt.Println("Initialized with prefix:", prefix)
	}

	return &Service{
		Storage:    storage,
	}, nil
}
