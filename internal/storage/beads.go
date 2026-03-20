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

	storage.SetConfig(ctx, "status.custom", "ready_to_sprint")

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


/////////////////////////////////////////////////////////////////////////////////////
// Dependency Management API wrappers over the underlying beads storage API
/////////////////////////////////////////////////////////////////////////////////////
// AddDependency creates a dependency edge between two issues.
// It is a thin wrapper over the underlying beads storage AddDependency API.
func (s *BeadsService) AddDependency(ctx context.Context, issueID, dependsOnID string, depType models.DependencyType, actor string) error {
	dep := &beads.Dependency{
		IssueID:     issueID,
		DependsOnID: dependsOnID,
		Type:        beads.DependencyType(depType),
	}
	return s.Storage.AddDependency(ctx, dep, actor)
}

// RemoveDependency removes a dependency edge between two issues.
func (s *BeadsService) RemoveDependency(ctx context.Context, issueID, dependsOnID string, actor string) error {
	return s.Storage.RemoveDependency(ctx, issueID, dependsOnID, actor)
}

// GetDependencies returns issues that the given issue depends on.
func (s *BeadsService) GetDependencies(ctx context.Context, issueID string) ([]*models.Issue, error) {
	issues, err := s.Storage.GetDependencies(ctx, issueID)
	if err != nil {
		return nil, err
	}
	if len(issues) == 0 {
		return []*models.Issue{}, nil
	}
	// types.Issue is layout-compatible with models.Issue (alias to beads.Issue),
	// so we can return the slice directly as []*models.Issue.
	result := make([]*models.Issue, 0, len(issues))
	for _, iss := range issues {
		if iss == nil {
			continue
		}
		casted := models.Issue(*iss)
		result = append(result, &casted)
	}
	return result, nil
}

// GetDependents returns issues that depend on the given issue.
func (s *BeadsService) GetDependents(ctx context.Context, issueID string) ([]*models.Issue, error) {
	issues, err := s.Storage.GetDependents(ctx, issueID)
	if err != nil {
		return nil, err
	}
	if len(issues) == 0 {
		return []*models.Issue{}, nil
	}
	result := make([]*models.Issue, 0, len(issues))
	for _, iss := range issues {
		if iss == nil {
			continue
		}
		casted := models.Issue(*iss)
		result = append(result, &casted)
	}
	return result, nil
}
