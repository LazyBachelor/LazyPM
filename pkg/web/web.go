package web

import (
	"context"
	"embed"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web/server"
)

type WebConfig = service.Config

type Web struct{}

func NewWeb() *Web {
	return &Web{}
}

//go:embed assets/*
var assets embed.FS

func (w Web) Run(ctx context.Context, config WebConfig) error {
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	server := server.NewServer(server.Server{
		Address:  config.WebAddress,
		Assets:   assets,
		Services: svc,
	})

	fmt.Printf("Starting web server on %s...\n", config.WebAddress)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
