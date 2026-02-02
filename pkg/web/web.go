package web

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web/handler"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
	"github.com/LazyBachelor/LazyPM/pkg/web/server"
	"context"
	"embed"
	"fmt"

	"github.com/a-h/templ"
)

type WebConfig = service.Config

//go:embed assets/*
var assets embed.FS

func Run(ctx context.Context, config WebConfig) error {
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	handler := handler.NewIssuesHandler(svc)

	server := server.NewServer(server.Server{
		Address:  config.WebAddress,
		Assets:   assets,
		Services: svc,
		Routes: []server.Route{
			{Pattern: "/", Handeler: templ.Handler(routes.Index())},
			{Pattern: "/issues", Handeler: handler.HandleGetAllIssues()},
		},
	})

	fmt.Printf("Starting web server on %s...\n", config.WebAddress)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
