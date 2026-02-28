package storage

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"

	"github.com/steveyegge/beads"
)

type BeadsService struct {
	beads.Storage
}

func NewBeadsIssueStorage(ctx context.Context, storage beads.Storage, prefix string) (*BeadsService, error) {
	issue_prefix, err := storage.GetConfig(ctx, "issue_prefix")
	if err != nil || issue_prefix == "" {
		if err := storage.SetConfig(ctx, "issue_prefix", prefix); err != nil {
			return nil, fmt.Errorf("failed to set issue_prefix: %w", err)
		}
	}

	return &BeadsService{
		Storage: storage,
	}, nil
}

func (s *BeadsService) AllIssues(ctx context.Context) ([]models.Issue, error) {
	issuesPtr, err := s.Storage.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return nil, err
	}

	if len(issuesPtr) == 0 {
		return []models.Issue{}, nil
	}

	issues := models.IssuesPtrToIssues(issuesPtr)

	return issues, nil
}

func (s *BeadsService) DeleteIssues() error {

	var deleteIssues = "DELETE FROM issues;"

	if _, err := s.UnderlyingDB().Exec(deleteIssues); err != nil {
		return err
	}
	return nil
}
