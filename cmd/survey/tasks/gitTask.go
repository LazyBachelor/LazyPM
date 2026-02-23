package tasks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/charmbracelet/huh"
	"github.com/go-git/go-git/v6"
)

var (
	gitTaskDescription = `You are tasked with performing a Git operation.
This task will test your ability to use Git effectively.`
)

type GitTask struct {
	setupIssue *models.Issue
	repo       *git.Repository
	done       bool

	app *service.App
}

func NewGitTask(app *service.App) *GitTask {
	return &GitTask{app: app, done: false}
}

func (t *GitTask) Config() task.Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/git-task-stats.json")
}

func (t *GitTask) Details() taskui.TaskDetails {
	return BaseDetails().WithTitle("Git Task").WithDescription(gitTaskDescription)
}

func (t *GitTask) Questions(interfaceType task.InterfaceType) taskui.Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[string]().Title("What Git Interface did you use?").
				Options(
					huh.Option[string]{Value: "cli", Key: "Command Line Interface"},
					huh.Option[string]{Value: "tui", Key: "Terminal User Interface"},
					huh.Option[string]{Value: "gui", Key: "Graphical User Interface"},
				),
		),
	)
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
		[]byte("This is a Git task. Please perform a Git operation here."),
		os.FileMode(os.O_WRONLY|os.O_CREATE))

	t.setupIssue = models.NewBaseIssue().
		WithTitle("Git Task Setup Issue").
		WithDescription(gitTaskDescription).Build()

	if err := t.app.Issues.CreateIssue(ctx, t.setupIssue, ""); err != nil {
		return err
	}

	return nil
}

func (t *GitTask) Validate(ctx context.Context) (bool, error) {

	if t.setupIssue == nil {
		return false, errors.New("setup issue not found")
	}

	if t.setupIssue.Assignee != "Me" {
		return false, fmt.Errorf("issue not assigned to self")
	}

	if t.setupIssue.Status != models.StatusInProgress {
		return false, fmt.Errorf("issue status is not in progress")
	}

	// Git related validation can be added here, such as checking for commits, branches, etc.

	return EndTaskWithTimeout(&t.done, "Task completed!", 5*time.Second)
}

func (t *GitTask) initRepo() (*git.Repository, error) {
	repo, err := git.PlainInit("./task/.git/", true)
	if err != nil {
		if errors.Is(err, git.ErrTargetDirNotEmpty) {
			repo, err = git.PlainOpen("./task/.git/")
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}
	return repo, nil
}
