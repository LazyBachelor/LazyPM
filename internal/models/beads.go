package models

import "github.com/steveyegge/beads"

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
