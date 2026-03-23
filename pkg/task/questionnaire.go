package task

import (
	"charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

type Questions = models.Questions

type QuestionnaireModel struct {
	Questions
	form          *huh.Form
	keys          []string
	width, height int
	userQuit      bool
}

func NewQuestionnaireModel(questions Questions, keys []string) *QuestionnaireModel {
	form := huh.NewForm(questions...).
		WithTheme(style.HuhCenterTheme()).
		WithLayout(huh.LayoutGrid(1, 1)).
		WithWidth(80)

	return &QuestionnaireModel{
		Questions: questions,
		form:      form,
		keys:      keys,
	}
}

func (q *QuestionnaireModel) Init() tea.Cmd {
	return q.form.Init()
}

func (q *QuestionnaireModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		q.SetSize(msg.Width, msg.Height)
		// Update form width to match window
		q.form.WithWidth(msg.Width)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			q.userQuit = true
			return q, tea.Quit
		}
	}

	form, cmd := q.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		q.form = f
	}

	if q.form.State == huh.StateCompleted {
		return q, tea.Quit
	}

	return q, cmd
}

func (q *QuestionnaireModel) View() tea.View {
	form := lipgloss.NewStyle().
		Width(q.width).Align(lipgloss.Center).
		Render(q.form.View())

	content := lipgloss.Place(
		q.width, q.height, lipgloss.Center, lipgloss.Center, form,
	)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (q *QuestionnaireModel) SetSize(width, height int) {
	q.width, q.height = width, height
}

func (q QuestionnaireModel) GetUserQuit() bool {
	return q.userQuit
}

func (q QuestionnaireModel) GetCompleted() bool {
	return q.form != nil && q.form.State == huh.StateCompleted
}

func (q QuestionnaireModel) GetAnswers() map[string]any {
	if q.form == nil || len(q.keys) == 0 {
		return nil
	}

	answers := make(map[string]any)
	for _, key := range q.keys {
		if key == "" {
			continue
		}
		answers[key] = q.form.Get(key)
	}
	if len(answers) == 0 {
		return nil
	}

	return answers
}
