package taskui

import "github.com/charmbracelet/huh"

type TaskDetails struct {
	Title          string
	Description    string
	TimeToComplete string
	Difficulty     string
}

func (td TaskDetails) WithTitle(title string) TaskDetails {
	td.Title = title
	return td
}

func (td TaskDetails) WithDescription(description string) TaskDetails {
	td.Description = description
	return td
}

func (td TaskDetails) WithTimeToComplete(time string) TaskDetails {
	td.TimeToComplete = time
	return td
}

func (td TaskDetails) WithDifficulty(difficulty string) TaskDetails {
	td.Difficulty = difficulty
	return td
}

type Questions []*huh.Group

func (q Questions) With(group *huh.Group) Questions {
	if group == nil {
		return q
	}
	return append(q, group)
}
