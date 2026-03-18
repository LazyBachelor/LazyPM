package web

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/browser"
	"github.com/LazyBachelor/LazyPM/pkg/web/handler"
	"github.com/LazyBachelor/LazyPM/pkg/web/server"
)

type Config = models.Config
type ValidationFeedback = models.ValidationFeedback

type Web struct {
	feedbackChan chan ValidationFeedback
	quitChan     chan bool
	submitChan   chan<- struct{}
}

func New() *Web {
	return &Web{}
}

//go:embed assets/*
var assets embed.FS

type tuiModel struct {
	width    int
	height   int
	address  string
	quitChan chan struct{}
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			select {
			case <-m.quitChan:
			default:
				close(m.quitChan)
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m tuiModel) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, fmt.Sprintf(
		"Web server running at %s\n\nPress q, esc, or Ctrl+C to stop the task and server.\n",
		m.address,
	))
}

func (w *Web) Run(ctx context.Context, config Config) error {
	app, cleanup, err := app.New(ctx, config)
	if err != nil {
		return err
	}
	defer cleanup()

	httpServer := server.NewServer(server.Server{
		Address: config.WebAddress,
		Assets:  assets,
		App:     app,
	})

	address := config.WebAddress
	if !strings.Contains(address, "http") {
		address = "http://localhost" + address
	}

	serverErr := make(chan error, 1)
	uiQuitChan := make(chan struct{})

	var screen *tea.Program
	var screenDone chan struct{}

	if w.quitChan != nil {
		screen = tea.NewProgram(
			tuiModel{
				address:  address,
				quitChan: uiQuitChan,
			},
			tea.WithAltScreen(),
		)

		screenDone = make(chan struct{})

		go func() {
			defer close(screenDone)
			screen.Run()
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	fmt.Println("Starting server at", address+"...")
	go func() {
		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
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

	if w.submitChan != nil {
		handler.SetSubmitChan(w.submitChan)
	}

	var shutdownMsg string

	if w.quitChan != nil {

		if err := browser.Open(address); err != nil {
			fmt.Printf("failed to open browser: %s\n", err)
		}
		select {
		case <-w.quitChan:
			shutdownMsg = "Task completed! Ending task and server..."
		case <-uiQuitChan:
			shutdownMsg = "Quit key pressed. Ending task and server..."
		case sig := <-sigChan:
			shutdownMsg = fmt.Sprintf("Received signal %s. Ending task and server...", sig)
		case <-ctx.Done():
			shutdownMsg = "Context cancelled. Ending task and server..."
		case err := <-serverErr:
			if screen != nil {
				screen.Quit()
				<-screenDone
			}
			return err
		}
	} else {
		select {
		case sig := <-sigChan:
			shutdownMsg = fmt.Sprintf("Received signal %s. Closing server...", sig)
		case <-ctx.Done():
			shutdownMsg = "Context cancelled. Closing server..."
		case err := <-serverErr:
			return err
		}
	}

	if screen != nil {
		screen.Quit()
		<-screenDone
	}

	fmt.Println(shutdownMsg)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(shutdownCtx)
}

func (w *Web) SetChannels(feedbackChan chan ValidationFeedback, quitChan chan bool) {
	w.feedbackChan = feedbackChan
	w.quitChan = quitChan
}

func (w *Web) SetSubmitChan(submitChan chan<- struct{}) {
	w.submitChan = submitChan
}
