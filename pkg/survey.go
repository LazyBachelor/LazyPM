package pkg

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
)

type SurveyConfig = service.Config

func Run(ctx context.Context, config SurveyConfig) error {
	_, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}
	defer cleanup()

	return nil
}
