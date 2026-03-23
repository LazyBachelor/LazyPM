package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"

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

	storage.UnderlyingDB().Exec(`
	CREATE TABLE IF NOT EXISTS sprints (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		issues TEXT,
		sprint_num INTEGER UNIQUE,
		is_backlog BOOLEAN DEFAULT 0
	);
	`)

	backlogNum, err := getBacklogSprintNum(storage)
	if err != nil {
		_, err = storage.UnderlyingDB().Exec(
			"INSERT INTO sprints (name, issues, sprint_num, is_backlog) VALUES (?, ?, 0, 1)",
			"backlog", "[]",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create backlog sprint: %w", err)
		}
		storage.SetConfig(ctx, "backlog_sprint", "0")
	} else {
		storage.SetConfig(ctx, "backlog_sprint", fmt.Sprintf("%d", backlogNum))
	}

	return &BeadsService{
		Storage: storage,
	}, nil
}

func (s *BeadsService) CreateIssue(ctx context.Context, issue *models.Issue, actor string) error {
	if err := s.Storage.CreateIssue(ctx, issue, actor); err != nil {
		return err
	}

	backlogNum, err := s.GetBacklogSprint(ctx)
	if err != nil {
		return nil
	}

	if err := s.AddIssueToSprint(ctx, issue.ID, backlogNum); err != nil {
		return nil
	}

	return nil
}

func (s *BeadsService) CreateIssues(ctx context.Context, issues []*models.Issue, actor string) error {
	if err := s.Storage.CreateIssues(ctx, issues, actor); err != nil {
		return err
	}

	backlogNum, err := s.GetBacklogSprint(ctx)
	if err != nil {
		return nil
	}

	for _, issue := range issues {
		if err := s.AddIssueToSprint(ctx, issue.ID, backlogNum); err != nil {
			continue
		}
	}

	return nil
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

	var deleteIssues = `DELETE FROM issues;
	DELETE FROM sprints;`

	if _, err := s.UnderlyingDB().Exec(deleteIssues); err != nil {
		return err
	}
	return nil
}

func (s *BeadsService) AddSprint(ctx context.Context) (int, error) {
	var addSprint = "INSERT INTO sprints (sprint_num, issues) VALUES ((SELECT IFNULL(MAX(sprint_num), 0) + 1 FROM sprints), '[]');"

	r, err := s.UnderlyingDB().Exec(addSprint)
	if err != nil {
		return 0, fmt.Errorf("failed to add sprint: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	var sprintNum int
	err = s.UnderlyingDB().QueryRow("SELECT sprint_num FROM sprints WHERE id = ?", id).Scan(&sprintNum)
	if err != nil {
		return 0, fmt.Errorf("failed to get sprint_num: %w", err)
	}

	return sprintNum, nil
}

func (s *BeadsService) RemoveSprint(ctx context.Context, sprintNum int) error {
	result, err := s.UnderlyingDB().Exec("DELETE FROM sprints WHERE sprint_num = ?", sprintNum)
	if err != nil {
		return fmt.Errorf("failed to remove sprint: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sprint %d not found", sprintNum)
	}

	return nil
}

func (s *BeadsService) GetSprints(ctx context.Context) ([]int, error) {
	rows, err := s.UnderlyingDB().Query("SELECT sprint_num FROM sprints WHERE is_backlog = 0 ORDER BY sprint_num")
	if err != nil {
		return nil, fmt.Errorf("failed to get sprints: %w", err)
	}
	defer rows.Close()

	var sprints []int
	for rows.Next() {
		var sprintNum int
		if err := rows.Scan(&sprintNum); err != nil {
			return nil, fmt.Errorf("failed to scan sprint: %w", err)
		}
		sprints = append(sprints, sprintNum)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sprints: %w", err)
	}

	return sprints, nil
}

func (s *BeadsService) GetIssuesBySprint(ctx context.Context, sprintNum int) ([]*models.Issue, error) {
	var issuesJSON string
	err := s.UnderlyingDB().QueryRow("SELECT issues FROM sprints WHERE sprint_num = ?", sprintNum).Scan(&issuesJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Issue{}, nil
		}
		return nil, fmt.Errorf("failed to get sprint issues: %w", err)
	}

	var issueIDs []string
	if err := json.Unmarshal([]byte(issuesJSON), &issueIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal issues: %w", err)
	}

	if len(issueIDs) == 0 {
		return []*models.Issue{}, nil
	}

	var issues []*models.Issue
	for _, id := range issueIDs {
		issue, err := s.Storage.GetIssue(ctx, id)
		if err != nil {
			// Skip issues that don't exist or can't be retrieved
			continue
		}
		issues = append(issues, issue)
	}

	return issues, nil
}

func (s *BeadsService) AddIssueToSprint(ctx context.Context, issueID string, sprintNum int) error {
	var issuesJSON string
	err := s.UnderlyingDB().QueryRow("SELECT issues FROM sprints WHERE sprint_num = ?", sprintNum).Scan(&issuesJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("sprint %d not found", sprintNum)
		}
		return fmt.Errorf("failed to get sprint: %w", err)
	}

	var issueIDs []string
	if err := json.Unmarshal([]byte(issuesJSON), &issueIDs); err != nil {
		return fmt.Errorf("failed to unmarshal issues: %w", err)
	}

	if slices.Contains(issueIDs, issueID) {
		return nil
	}

	issueIDs = append(issueIDs, issueID)

	updatedJSON, err := json.Marshal(issueIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal issues: %w", err)
	}

	_, err = s.UnderlyingDB().Exec("UPDATE sprints SET issues = ? WHERE sprint_num = ?", string(updatedJSON), sprintNum)
	if err != nil {
		return fmt.Errorf("failed to update sprint: %w", err)
	}

	return nil
}

func (s *BeadsService) RemoveIssueFromSprint(ctx context.Context, issueID string, sprintNum int) error {
	var issuesJSON string
	err := s.UnderlyingDB().QueryRow("SELECT issues FROM sprints WHERE sprint_num = ?", sprintNum).Scan(&issuesJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("sprint %d not found", sprintNum)
		}
		return fmt.Errorf("failed to get sprint: %w", err)
	}

	var issueIDs []string
	if err := json.Unmarshal([]byte(issuesJSON), &issueIDs); err != nil {
		return fmt.Errorf("failed to unmarshal issues: %w", err)
	}

	found := false
	var updatedIDs []string
	for _, id := range issueIDs {
		if id != issueID {
			updatedIDs = append(updatedIDs, id)
		} else {
			found = true
		}
	}

	if !found {
		return nil
	}

	updatedJSON, err := json.Marshal(updatedIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal issues: %w", err)
	}

	_, err = s.UnderlyingDB().Exec("UPDATE sprints SET issues = ? WHERE sprint_num = ?", string(updatedJSON), sprintNum)
	if err != nil {
		return fmt.Errorf("failed to update sprint: %w", err)
	}

	return nil
}

func (s *BeadsService) GetBacklogSprint(ctx context.Context) (int, error) {
	sprintNum, err := getBacklogSprintNum(s.Storage)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("backlog sprint not found")
		}
		return 0, fmt.Errorf("failed to get backlog sprint: %w", err)
	}
	return sprintNum, nil
}

// GetIssuesNotInAnySprint returns issues that are only in the backlog
func (s *BeadsService) GetIssuesNotInAnySprint(ctx context.Context) ([]*models.Issue, error) {
	allIssues, err := s.Storage.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all issues: %w", err)
	}

	rows, err := s.UnderlyingDB().Query("SELECT sprint_num, issues FROM sprints WHERE is_backlog = 0")
	if err != nil {
		return nil, fmt.Errorf("failed to get sprint issues: %w", err)
	}
	defer rows.Close()

	issuesInSprints := make(map[string]bool)
	for rows.Next() {
		var sprintNum int
		var issuesJSON string
		if err := rows.Scan(&sprintNum, &issuesJSON); err != nil {
			continue
		}
		var issueIDs []string
		if err := json.Unmarshal([]byte(issuesJSON), &issueIDs); err != nil {
			continue
		}
		for _, id := range issueIDs {
			issuesInSprints[id] = true
		}
	}

	var backlogIssues []*models.Issue
	for _, issue := range allIssues {
		if !issuesInSprints[issue.ID] {
			backlogIssues = append(backlogIssues, issue)
		}
	}

	return backlogIssues, nil
}

func getBacklogSprintNum(storage beads.Storage) (int, error) {
	var sprintNum int
	err := storage.UnderlyingDB().QueryRow("SELECT sprint_num FROM sprints WHERE is_backlog = 1 LIMIT 1").Scan(&sprintNum)
	if err != nil {
		return 0, err
	}
	return sprintNum, nil
}
