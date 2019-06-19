package lambdabase

import (
	"github.com/go-nacelle/nacelle"
)

func doInit(config nacelle.Config, container nacelle.ServiceContainer, handler interface{}) error {
	if err := container.Inject(handler); err != nil {
		return err
	}

	if initializer, ok := handler.(nacelle.Initializer); ok {
		if err := initializer.Init(config); err != nil {
			return err
		}
	}

	return nil
}
