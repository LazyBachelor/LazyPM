package models

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/muesli/reflow/truncate"
	"github.com/steveyegge/beads"
)

type (
	Issue            = beads.Issue
	Status           = beads.Status
	IssueType        = beads.IssueType
	Dependency       = beads.Dependency
	DependencyType   = beads.DependencyType
	Label            = beads.Label
	Comment          = beads.Comment
	Event            = beads.Event
	EventType        = beads.EventType
	BlockedIssue     = beads.BlockedIssue
	TreeNode         = beads.TreeNode
	IssueFilter      = beads.IssueFilter
	WorkFilter       = beads.WorkFilter
	StaleFilter      = beads.StaleFilter
	DependencyCounts = beads.DependencyCounts
	IssueWithCounts  = beads.IssueWithCounts
	SortPolicy       = beads.SortPolicy
	EpicStatus       = beads.EpicStatus
)

// Status constants
const (
	StatusOpen       = beads.StatusOpen
	StatusInProgress = beads.StatusInProgress
	StatusBlocked    = beads.StatusBlocked
	StatusDeferred   = beads.StatusDeferred
	StatusClosed     = beads.StatusClosed
)

// IssueType constants
const (
	TypeBug     = beads.TypeBug
	TypeFeature = beads.TypeFeature
	TypeTask    = beads.TypeTask
	TypeEpic    = beads.TypeEpic
	TypeChore   = beads.TypeChore
)

// DependencyType constants
const (
	DepBlocks            = beads.DepBlocks
	DepRelated           = beads.DepRelated
	DepParentChild       = beads.DepParentChild
	DepDiscoveredFrom    = beads.DepDiscoveredFrom
	DepConditionalBlocks = beads.DepConditionalBlocks // B runs only if A fails (bd-kzda)
)

// SortPolicy constants
const (
	SortPolicyHybrid   = beads.SortPolicyHybrid
	SortPolicyPriority = beads.SortPolicyPriority
	SortPolicyOldest   = beads.SortPolicyOldest
)

// EventType constants
const (
	EventCreated           = beads.EventCreated
	EventUpdated           = beads.EventUpdated
	EventStatusChanged     = beads.EventStatusChanged
	EventCommented         = beads.EventCommented
	EventClosed            = beads.EventClosed
	EventReopened          = beads.EventReopened
	EventDependencyAdded   = beads.EventDependencyAdded
	EventDependencyRemoved = beads.EventDependencyRemoved
	EventLabelAdded        = beads.EventLabelAdded
	EventLabelRemoved      = beads.EventLabelRemoved
	EventCompacted         = beads.EventCompacted
)

func IssuesPtrToIssues(issuePtr []*Issue) []Issue {
	issues := make([]Issue, 0, len(issuePtr))
	for _, issue := range issuePtr {
		issues = append(issues, *issue)
	}
	return issues
}

func FormatIssueRow(issue Issue) string {
	return fmt.Sprintf(
		"%s\t%s\t%s\t%s\t%s\t%d",
		truncate.String(issue.ID, 5),
		truncate.StringWithTail(issue.Title, 25, "..."),
		truncate.StringWithTail(issue.Description, 40, "..."),
		issue.Status,
		issue.IssueType,
		issue.Priority,
	)
}

func PrintIssues(issues []Issue) {
	w := tabwriter.NewWriter(os.Stdout, 8, 10, 5, ' ', 0)

	fmt.Fprintln(w, "ID\tTITLE\tDESCRIPTION\tSTATUS\tTYPE\tPRIORITY")

	for _, issue := range issues {
		fmt.Fprintln(w, FormatIssueRow(issue))
	}

	w.Flush()
}
