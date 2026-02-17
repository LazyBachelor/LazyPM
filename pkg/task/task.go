package task

import (
	"context"
	"fmt"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/service"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type Tasker interface {
	Init() *Task
	Config() TaskConfig

	Details() taskui.TaskDetails
	QuestionsFunc() QuestionsFunc

	InterfaceType() Interface

	ValidateFunc(context.Context) (ok bool, errorMsg error)
	DbStateFunc(context.Context) error
}

type Task struct {
	svc    *service.Services
	Config TaskConfig

	Interface     Interface
	InterfaceType InterfaceType

	details       taskui.TaskDetails
	questionsFunc QuestionsFunc

	validateFunc ValidateFunc
	dbStateFunc  DbStateFunc

	feedbackChan chan ValidationFeedback
	doneChan     chan bool
	quitChan     chan bool
}

func NewTask(svc *service.Services, details taskui.TaskDetails, questionsFunc QuestionsFunc) *Task {
	return &Task{
		details:       details,
		questionsFunc: questionsFunc,
		svc:           svc,
	}
}

func (t *Task) IntroduceTask() error {
	if t.details == (taskui.TaskDetails{}) {
		return fmt.Errorf("details is not set")
	}

	detailsScreen := taskui.NewTaskModel(t.details)

	model, err := tea.NewProgram(detailsScreen, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		return ErrUserQuit
	}
	return nil
}

func (t *Task) StartInterface(ctx context.Context, cfg TaskConfig) error {
	if t.Interface == nil {
		return fmt.Errorf("interfaceType is not set")
	}

	return t.Interface.Run(ctx, cfg)
}

func (t *Task) Initialize(ctx context.Context) error {
	if t.dbStateFunc == nil {
		return fmt.Errorf("dbStateFunc is not set")
	}
	return t.dbStateFunc(ctx)
}

func (t *Task) Validate(ctx context.Context) (bool, error) {
	if t.validateFunc == nil {
		return false, fmt.Errorf("validateFunc is not set")
	}
	return t.validateFunc(ctx)
}

func (t *Task) StartQuestionnaire() error {
	if t.questionsFunc == nil {
		return fmt.Errorf("questionsFunc is not set")
	}

	questions := t.questionsFunc(t.InterfaceType)
	if questions == nil {
		return fmt.Errorf("questions is nil")
	}

	questionare := taskui.NewQuestionnaireModel(questions)
	model, err := tea.NewProgram(questionare, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		return ErrUserQuit
	}
	return nil
}

func (t *Task) SetConfigFunc(fn ConfigFunc) {
	t.Config = fn()
}

func (t *Task) SetInterface(interfaceType Interface) {
	t.Interface = interfaceType
}

func (t *Task) SetInterfaceType(interfaceType InterfaceType) {
	t.InterfaceType = interfaceType
}

func (t *Task) SetDbStateFunc(fn DbStateFunc) {
	t.dbStateFunc = fn
}

func (t *Task) SetValidateFunc(fn ValidateFunc) {
	t.validateFunc = fn
}

func (t *Task) SetChannels(feedbackChan chan ValidationFeedback, doneChan chan bool, quitChan chan bool) {
	t.feedbackChan = feedbackChan
	t.doneChan = doneChan
	t.quitChan = quitChan
}

func (t *Task) StartValidationLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ok, err := t.Validate(ctx)
			feedback := ValidationFeedback{
				Timestamp: time.Now(),
			}
			if ok {
				feedback.Success = true
				feedback.Message = "Task completed successfully!"
				t.feedbackChan <- feedback
				t.doneChan <- true
				return
			} else {
				feedback.Success = false
				if err != nil {
					feedback.Message = err.Error()
				} else {
					feedback.Message = "Task not yet complete"
				}
				t.feedbackChan <- feedback
			}
		case <-t.quitChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
