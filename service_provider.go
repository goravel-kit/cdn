package cdn

import (
	"github.com/goravel/framework/contracts/foundation"
)

const Binding = "cdn"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Singleton(Binding, func(app foundation.Application) (any, error) {
		return NewCdn(app.MakeConfig()), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	app.Publishes("github.com/goravel-kit/cdn", map[string]string{
		"config/cdn.go": app.ConfigPath("cdn.go"),
	})
}
