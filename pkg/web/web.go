package web

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/web/handler"
	"github.com/LazyBachelor/LazyPM/pkg/web/server"
)

type WebConfig = service.Config

type Web struct {
	feedbackChan chan task.ValidationFeedback
	quitChan     chan bool
}

func NewWeb() *Web {
	return &Web{}
}

//go:embed assets/*
var assets embed.FS

func (w *Web) Run(ctx context.Context, config WebConfig) error {
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	httpServer := server.NewServer(server.Server{
		Address:  config.WebAddress,
		Assets:   assets,
		Services: svc,
	})

	fmt.Printf("Starting web server on %s...\n", config.WebAddress)

	serverErr := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	if w.feedbackChan != nil {
		go func() {
			for feedback := range w.feedbackChan {
				handler.SetTaskFeedback(feedback)
			}
		}()
	}

	select {
	case <-w.quitChan:
		fmt.Println("Task completed! Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return httpServer.Shutdown(context.Background())
	}
}

func (w *Web) SetChannels(feedbackChan chan task.ValidationFeedback, quitChan chan bool) {
	w.feedbackChan = feedbackChan
	w.quitChan = quitChan
}
