package models

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"unicode"

	"charm.land/lipgloss/v2"
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

func IssueString(issue Issue) string {
	labelStyle := lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("5"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("6"))

	boxStyle := lipgloss.NewStyle().Padding(1)

	line := func(label, value string) string {
		return labelStyle.Render(label) + ": " + value + "\t"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		line("Title", titleStyle.Render(issue.Title)+"\t"+line("Assignee", issue.Assignee)),
		line("ID", issue.ID+"\t"+line("Type", string(issue.IssueType))+line("Status", string(issue.Status))+line("Priority", fmt.Sprintf("%d", issue.Priority))),
		line("Description", "\n"+issue.Description),
	)

	return boxStyle.Render(content)
}

func IssuesPtrToIssues(issuePtr []*Issue) []Issue {
	issues := make([]Issue, 0, len(issuePtr))
	for _, issuePtr := range issuePtr {
		if issuePtr != nil {
			issues = append(issues, *issuePtr)
		}
	}
	return issues
}

func sanitizeCell(s string) string {
	s = strings.Map(func(r rune) rune {
		switch r {
		case '\n', '\r', '\t':
			return ' '
		default:
			return r
		}
	}, s)

	return strings.Join(strings.FieldsFunc(s, unicode.IsSpace), " ")
}

func FormatIssueRow(issue Issue) string {
	id := truncate.String(sanitizeCell(issue.ID), 10)
	title := truncate.StringWithTail(sanitizeCell(issue.Title), 35, "...")
	description := truncate.StringWithTail(sanitizeCell(issue.Description), 40, "...")
	status := sanitizeCell(string(issue.Status))
	issueType := sanitizeCell(string(issue.IssueType))

	return fmt.Sprintf(
		"%s\t%s\t%s\t%s\t%s\t%d",
		id,
		title,
		description,
		status,
		issueType,
		issue.Priority,
	)
}

func PrintIssues(issues []Issue) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "ID\tTITLE\tDESCRIPTION\tSTATUS\tTYPE\tPRIORITY")

	for _, issue := range issues {
		fmt.Fprintln(w, FormatIssueRow(issue))
	}

	_ = w.Flush()
}
