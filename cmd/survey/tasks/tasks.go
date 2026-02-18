package tasks

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
)

func init() {
	task.Register("create_issue", func(svc *service.Services) task.Tasker {
		return NewCreateIssueTask(svc)
	})

	task.Register("coding_task", func(svc *service.Services) task.Tasker {
		return NewCodingTask(svc)
	})
}
