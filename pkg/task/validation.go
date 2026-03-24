package task

import (
	"context"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
)

type ValidationEngine struct {
	task Tasker
}

func (v *ValidationEngine) Start(ctx context.Context, submitChan <-chan models.ValidationTrigger, onFeedback func(ValidationFeedback, models.ValidationTriggerSource)) (done <-chan struct{}, stop chan<- struct{}) {
	doneChan := make(chan struct{}, 1)
	stopChan := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case trigger := <-submitChan:
				feedback := v.task.Validate(ctx)
				source := trigger.Source
				if source == "" {
					source = models.ValidationTriggerUnknown
				}

				if onFeedback != nil {
					onFeedback(feedback, source)
				}

				if feedback.Success {
					time.Sleep(3 * time.Second)
					doneChan <- struct{}{}
					return
				}

			case <-stopChan:
				return

			case <-ctx.Done():
				return
			}
		}
	}()

	return doneChan, stopChan
}
