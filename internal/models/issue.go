package models

type IssueBuilder struct {
	ID          string
	Title       string
	Description string
	Status      Status
	IssueType   IssueType
	Priority    int
}

func NewIssueBuilder() *IssueBuilder {
	return &IssueBuilder{}
}

func NewBaseIssue() *IssueBuilder {
	return NewIssueBuilder().
		WithID("pm-abc").
		WithTitle("Basic Issue").
		WithDescription("Basic Description").
		WithStatus(StatusOpen).
		WithIssueType(TypeTask)

}

func (b *IssueBuilder) WithID(id string) *IssueBuilder {
	b.ID = id
	return b
}

func (b *IssueBuilder) WithTitle(title string) *IssueBuilder {
	b.Title = title
	return b
}

func (b *IssueBuilder) WithDescription(description string) *IssueBuilder {
	b.Description = description
	return b
}

func (b *IssueBuilder) WithStatus(status Status) *IssueBuilder {
	b.Status = status
	return b
}

func (b *IssueBuilder) WithIssueType(issueType IssueType) *IssueBuilder {
	b.IssueType = issueType
	return b
}

func (b *IssueBuilder) WithPriority(priority int) *IssueBuilder {
	b.Priority = priority
	return b
}

func (b IssueBuilder) Build() Issue {
	return Issue{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Status:      b.Status,
		IssueType:   b.IssueType,
		Priority:    b.Priority,
	}
}
