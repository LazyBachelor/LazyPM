package models

import "sort"

func PriorityLabel(p int) string {
	if p >= 2 {
		return "E"
	}
	if p >= 1 {
		return "H"
	}
	return "N"
}

func StateLabel(s Status) string {
	if s == StatusClosed {
		return "Closed"
	}
	return "Open"
}

func TypeLabel(t IssueType) string {
	switch t {
	case TypeTask:
		return "Task"
	case TypeBug:
		return "Bug"
	case TypeFeature:
		return "Story"
	case TypeEpic:
		return "Epic"
	default:
		return string(t)
	}
}

func SortIssuesByPriority(issues []Issue, idDesc bool) {
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Priority != issues[j].Priority {
			return issues[i].Priority > issues[j].Priority
		}
		if idDesc {
			return issues[i].ID > issues[j].ID
		}
		return issues[i].ID < issues[j].ID
	})
}

func RecentIssuesByID(issues []Issue, limit int) []Issue {
	if limit <= 0 || len(issues) == 0 {
		return nil
	}
	sorted := make([]Issue, len(issues))
	copy(sorted, issues)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID > sorted[j].ID })
	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

func CountByStatus(issues []Issue) (open, closed int) {
	for _, issue := range issues {
		if issue.Status == StatusClosed {
			closed++
		} else {
			open++
		}
	}
	return open, closed
}

func FindIssueByID(issues []Issue, id string) *Issue {
	for i := range issues {
		if issues[i].ID == id {
			return &issues[i]
		}
	}
	return nil
}

func FilterIssues(issues []Issue, filterPri, filterState string) []Issue {
	var filtered []Issue
	for _, issue := range issues {
		isClosed := issue.Status == StatusClosed
		if filterState == "closed" && !isClosed {
			continue
		}
		if filterState == "open" && isClosed {
			continue
		}
		if filterPri != "" && PriorityLabel(issue.Priority) != filterPri {
			continue
		}
		filtered = append(filtered, issue)
	}
	return filtered
}

func GroupIssuesByStatus(issues []Issue) (todo, inProgress, done []Issue) {
	for _, issue := range issues {
		switch issue.Status {
		case StatusOpen, StatusBlocked, StatusDeferred:
			todo = append(todo, issue)
		case StatusInProgress:
			inProgress = append(inProgress, issue)
		case StatusClosed:
			done = append(done, issue)
		}
	}
	return todo, inProgress, done
}
