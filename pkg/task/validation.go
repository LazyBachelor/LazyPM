package task

import (
	"context"
	"time"
)

type ValidationEngine struct {
	task Tasker
}

func (v *ValidationEngine) Start(ctx context.Context, onFeedback func(ValidationFeedback)) (done <-chan struct{}, stop chan<- struct{}) {
	doneChan := make(chan struct{}, 1)
	stopChan := make(chan struct{}, 1)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				feedback := v.task.Validate(ctx)

				if onFeedback != nil {
					onFeedback(feedback)
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
