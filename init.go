package lambdabase

import (
	"context"

	"github.com/go-nacelle/nacelle/v2"
	"github.com/go-nacelle/service/v2"
)

func doInit(ctx context.Context, container *nacelle.ServiceContainer, handler interface{}) error {
	if err := service.Inject(ctx, container, handler); err != nil {
		return err
	}

	if initializer, ok := handler.(nacelle.Initializer); ok {
		if err := initializer.Init(ctx); err != nil {
			return err
		}
	}

	return nil
}
