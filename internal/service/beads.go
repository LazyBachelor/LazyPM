package service

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"

	"github.com/steveyegge/beads"
)

type BeadsService struct {
	beads.Storage
}

func NewBeadsService(ctx context.Context, storage beads.Storage, prefix string) (*BeadsService, error) {
	issue_prefix, err := storage.GetConfig(ctx, "issue_prefix")
	if err != nil || issue_prefix == "" {
		if err := storage.SetConfig(ctx, "issue_prefix", prefix); err != nil {
			return nil, fmt.Errorf("failed to set issue_prefix: %w", err)
		}
		fmt.Println("Initialized with prefix:", prefix)
	}

	return &BeadsService{
		Storage: storage,
	}, nil
}

func (s *BeadsService) AllIssues(ctx context.Context) ([]models.Issue, error) {
	issuesPtr, err := s.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return nil, err
	}
	return models.IssuesPtrToIssues(issuesPtr), nil
}

// SearchIssues returns issues matching the query (title, description, id, assignee) and optional filter.
// The query searches across title, description, and id (LIKE %query%).
// Assignee exact match is also searched when query is non-empty.
func (s *BeadsService) SearchIssues(ctx context.Context, query string, filter models.IssueFilter) ([]*models.Issue, error) {
	seen := make(map[string]bool)
	var merged []*models.Issue

	// Search by query (title, description, id)
	issuesPtr, err := s.Storage.SearchIssues(ctx, query, filter)
	if err != nil {
		return nil, err
	}
	for _, issue := range issuesPtr {
		if issue != nil && !seen[issue.ID] {
			seen[issue.ID] = true
			merged = append(merged, issue)
		}
	}

	// Also search by assignee when query is non-empty (exact match)
	if query != "" && filter.Assignee == nil {
		assigneeFilter := filter
		assigneeFilter.Assignee = &query
		assigneeFilter.TitleSearch = ""
		assigneeFilter.DescriptionContains = ""
		assigneePtr, err := s.Storage.SearchIssues(ctx, "", assigneeFilter)
		if err != nil {
			return nil, err
		}
		for _, issue := range assigneePtr {
			if issue != nil && !seen[issue.ID] {
				seen[issue.ID] = true
				merged = append(merged, issue)
			}
		}
	}

	return merged, nil
}

func (s *BeadsService) DeleteIssues() error {

	var deleteIssues = "DELETE FROM issues;"

	if _, err := s.UnderlyingDB().Exec(deleteIssues); err != nil {
		return err
	}
	return nil
}
