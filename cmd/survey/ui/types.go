package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
)

type TaskDetails struct {
	Title          string
	Description    string
	TimeToComplete string
	Difficulty     string
}

type TaskModel struct {
	TaskDetails
	keys          TaskHelpKeys
	help          help.Model
	width, height int
}

type Questions []*huh.Group

type QuestionnaireModel struct {
	Questions
	form          *huh.Form
	width, height int
}
