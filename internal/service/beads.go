package service

import (
	"context"
	"fmt"
	"time"

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

func (s *BeadsService) GetComments(ctx context.Context, issueID string) ([]models.Comment, error) {
	query := `SELECT id, issue_id, author, text, created_at FROM comments WHERE issue_id = ? ORDER BY created_at ASC LIMIT 100`

	rows, err := s.UnderlyingDB().QueryContext(ctx, query, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.IssueID, &c.Author, &c.Text, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *BeadsService) AddComment(ctx context.Context, issueID, author, text string) (*models.Comment, error) {
	query := `INSERT INTO comments (issue_id, author, text, created_at) VALUES (?, ?, ?, ?)`

	createdAt := time.Now()
	result, err := s.UnderlyingDB().ExecContext(ctx, query, issueID, author, text, createdAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Comment{
		ID:        id,
		IssueID:   issueID,
		Author:    author,
		Text:      text,
		CreatedAt: createdAt,
	}, nil
}
