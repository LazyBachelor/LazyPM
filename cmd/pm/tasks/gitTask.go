package tasks

import (
	"context"
	"errors"
	"os"
	"strings"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
	"github.com/go-git/go-git/v6"
)

const gitTaskDescription = `You are tasked with performing a Git operation.

This task will test your ability to use Git effectively within a project management workflow. Your goal is to modify a file in a Git repository and commit the change.

Your task:
1. Set the Issue status to "In Progress" when you are ready to start.
2. A folder called "task" is created in the project directory when you start this task. Open it.
3. Inside the folder you will find README.md. Edit this file and add something to it (e.g. your name, a short note, or a new line). The file must be different from its original content.
4. Commit your change:
   - Open a terminal and change into the task folder.
   - Run "git add ." to stage the changes.
   - Run "git commit -m "updated codebase"" to commit.
5. Set the Issue status to "Closed" when done.

The repository is in the task folder (./task).`

const gitTaskReadmeContent = "This is a Git task. Add your name or a short note below this line to complete it, then commit your change."

type GitTask struct {
	app        *App
	done       bool
	repo       *git.Repository
	setupIssue *Issue
}

var gitTaskInProgress = false

func NewGitTask(app *App) *GitTask {
	return &GitTask{app: app, done: false}
}

func (t *GitTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/git-task-stats.json")
}

func (t *GitTask) Details(interfaceType InterfaceType) TaskDetails {
	return BaseDetails(interfaceType).
		WithTitle("Git Task").
		WithDescription(gitTaskDescription).
		WithDifficulty("Hard").WithTimeToComplete("5m")
}

func (t *GitTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).
		With(
			huh.NewGroup(
				huh.NewSelect[string]().Key("git_interface_used").
					Title("What Git Interface did you use for the task?").
					Description("If you used multiple interfaces, select the one you relied on most for this task.").
					Options(
						huh.Option[string]{Value: "cli", Key: "Command Line Interface"},
						huh.Option[string]{Value: "tui", Key: "Terminal User Interface"},
						huh.Option[string]{Value: "gui", Key: "Graphical User Interface"},
					),
			),
		).With(
		huh.NewGroup(
			huh.NewSelect[string]().Key("git_interface_normally").
				Title("What Git Interface do you normally use?").
				Description("If you regularly use multiple interfaces, select the one you prefer most.").
				Options(
					huh.Option[string]{Value: "cli", Key: "Command Line Interface"},
					huh.Option[string]{Value: "tui", Key: "Terminal User Interface"},
					huh.Option[string]{Value: "gui", Key: "Graphical User Interface"},
				),
		))
}

func (t *GitTask) QuestionnaireKeys(_ InterfaceType) []string {
	return BaseKeys().With("git_interface_used", "git_interface_normally")
}

func (t *GitTask) Setup(ctx context.Context) error {
	gitTaskInProgress = false
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	_ = os.RemoveAll("./task")

	var err error
	t.repo, err = t.initRepo()
	if err != nil {
		return err
	}

	_ = os.WriteFile("./task/README.md", []byte(gitTaskReadmeContent), 0o644)
	_ = os.WriteFile("./task/.gitattributes", []byte("* text=auto\n"), 0o644)

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

	issues, err := FetchIssues(ctx, t.app, t.setupIssue)
	if err != nil {
		return expect.ValidationFeedback
	}

	_ = issues

	issue := t.setupIssue

	if issue.Status == models.StatusInProgress || gitTaskInProgress {
		gitTaskInProgress = true
	} else {
		expect.Fail("The issue should be marked as In Progress while working on the Git task.")
		return expect.ValidationFeedback
	}

	// Ensure the repository is available.
	if t.repo == nil {
		t.repo, err = t.initRepo()
		if err != nil {
			expect.Fail("Failed to open the Git repository: " + err.Error())
			return expect.ValidationFeedback
		}
	}

	headRef, err := t.repo.Head()
	if err != nil {
		expect.Fail("No commits found in the Git repository. Please commit your changes.")
		return expect.ValidationFeedback
	}

	commit, err := t.repo.CommitObject(headRef.Hash())
	if err != nil {
		expect.Fail("Failed to read the latest commit: " + err.Error())
		return expect.ValidationFeedback
	}

	expect.NotEmptyAndEqual(strings.TrimSpace(commit.Message), "updated codebase", "Git commit message")

	tree, err := commit.Tree()
	if err != nil {
		expect.Fail("Failed to read commit tree: " + err.Error())
		return expect.ValidationFeedback
	}

	readmeFile, err := tree.File("README.md")
	if err != nil {
		expect.Fail("README.md should be part of the committed changes.")
		return expect.ValidationFeedback
	}

	readmeContent, err := readmeFile.Contents()
	if err != nil {
		expect.Fail("Failed to read README.md from the commit: " + err.Error())
		return expect.ValidationFeedback
	}

	expect.Assert(readmeContent != gitTaskReadmeContent,
		"You should modify README.md content before committing (make appropriate changes to complete the task).")

	if wt, err := t.repo.Worktree(); err == nil {
		if status, err := wt.Status(); err == nil {
			expect.Assert(status.IsClean(), "The working tree should be clean after committing (no unstaged changes).")
		}
	}

	expect.Assert(gitTaskInProgress && issue.Status == models.StatusClosed,
		"The issue should be marked as Closed after completing the Git task.")

	return expect.Complete()
}

func (t *GitTask) initRepo() (*git.Repository, error) {
	repo, err := git.PlainInit("./task", false)
	if err != nil {
		if errors.Is(err, git.ErrTargetDirNotEmpty) {
			repo, err = git.PlainOpen("./task")
			if err != nil {
				return nil, err
			}
			return repo, nil
		}
		return nil, err
	}
	return repo, nil
}
