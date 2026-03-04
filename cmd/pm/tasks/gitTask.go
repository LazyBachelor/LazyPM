package tasks

import (
	"context"
	"errors"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
	"github.com/charmbracelet/huh"
	"github.com/go-git/go-git/v6"
)

const gitTaskDescription = `You are tasked with performing a Git operation.

This task will test your ability to use Git effectively within a project management workflow.

Your task:
1. Initialize or open a Git repository
2. Review the current repository status
3. Make appropriate changes to complete the task
4. Commit your changes with a meaningful message
5. Update the task status to reflect completion

The repository has been initialized in ./task/.git/ for you to work with.`

type GitTask struct {
	app        *App
	done       bool
	repo       *git.Repository
	setupIssue *Issue
}

func NewGitTask(app *App) *GitTask {
	return &GitTask{app: app, done: false}
}

func (t *GitTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/git-task-stats.json")
}

func (t *GitTask) Details() TaskDetails {
	return BaseDetails().WithTitle("Git Task").WithDescription(gitTaskDescription)
}

func (t *GitTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[string]().Title("What Git Interface did you use?").
				Key("git_interface_used").
				Options(
					huh.Option[string]{Value: "cli", Key: "Command Line Interface"},
					huh.Option[string]{Value: "tui", Key: "Terminal User Interface"},
					huh.Option[string]{Value: "gui", Key: "Graphical User Interface"},
				),
		),
	)
}

func (t *GitTask) QuestionnaireKeys(interfaceType InterfaceType) []string {
	_ = interfaceType
	return []string{"task_completed", "task_difficulty", "git_interface_used"}
}

func (t *GitTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	var err error
	t.repo, err = t.initRepo()
	if err != nil {
		return err
	}

	os.WriteFile("./task/README.md",
		[]byte("This is a Git task. Please perform a Git operation here."), 0o644)

	t.setupIssue = NewIssueBuilder().
		WithTitle("Git Task Setup Issue").
		WithDescription(gitTaskDescription).
		WithIssueType(models.TypeTask).
		Build()

	if err := t.app.Issues.CreateIssue(ctx, t.setupIssue, ""); err != nil {
		return err
	}

	return nil
}

func (t *GitTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	return expect.Complete()
}

func (t *GitTask) initRepo() (*git.Repository, error) {
	repo, err := git.PlainInit("./task/.git/", true)
	if err != nil {
		if errors.Is(err, git.ErrTargetDirNotEmpty) {
			repo, err = git.PlainOpen("./task/.git/")
			if err != nil {
				return nil, err
			}
			return repo, nil
		}
		return nil, err
	}
	return repo, nil
}
