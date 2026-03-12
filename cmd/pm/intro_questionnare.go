package main

import (
	"github.com/LazyBachelor/LazyPM/pkg/task"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type IntroQuestionnaire struct{}

func newIntroQuestionnaire() *IntroQuestionnaire {
	return &IntroQuestionnaire{}
}

func (iq *IntroQuestionnaire) Run() (map[string]any, error) {
	model := task.NewQuestionnaireModel(iq.Questions(), iq.Keys())
	app := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := app.Run(); err != nil {
		return nil, err
	}

	return model.GetAnswers(), nil
}

func (iq *IntroQuestionnaire) Questions() task.Questions {
	return task.Questions{
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which age group do you belong to?").
				Description("This helps us understand the background of participants.").
				Options(
					huh.NewOption("Under 18", "under_18"),
					huh.NewOption("18–24", "18_24"),
					huh.NewOption("25–34", "25_34"),
					huh.NewOption("35–44", "35_44"),
					huh.NewOption("45+", "45_plus"),
				).
				Key("age_group"),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Are you a student, employed, or both?").
				Description("This helps us understand your current situation.").
				Options(
					huh.NewOption("Student", "student"),
					huh.NewOption("Employed", "employed"),
					huh.NewOption("Both student and employed", "both"),
					huh.NewOption("Neither", "neither"),
				).
				Key("occupation_status"),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How far along are you in your education?").
				Description("Select the option that best matches your current level.").
				Options(
					huh.NewOption("Primary school", "primary"),
					huh.NewOption("Secondary school", "secondary"),
					huh.NewOption("Bachelor's degree", "bachelor"),
					huh.NewOption("Master's degree", "master"),
					huh.NewOption("PhD / Doctorate", "phd"),
					huh.NewOption("Other", "other"),
				).
				Key("education_level"),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How would you describe your experience with the command line?").
				Description("This helps us tailor the questions to your experience level.").
				Options(
					huh.NewOption("No experience", "none"),
					huh.NewOption("Some experience", "some"),
					huh.NewOption("Extensive experience", "extensive"),
				).
				Key("cli_experience"),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How often do you use the command line?").
				Description("Select the option that best describes your usage.").
				Options(
					huh.NewOption("Never", "never"),
					huh.NewOption("Rarely", "rarely"),
					huh.NewOption("Weekly", "weekly"),
					huh.NewOption("Several times a week", "multiple_weekly"),
					huh.NewOption("Daily", "daily"),
				).
				Key("cli_frequency"),
		),
	}
}

func (iq *IntroQuestionnaire) Keys() []string {
	return []string{"age_group", "occupation_status", "education_level", "cli_experience", "cli_frequency"}
}
