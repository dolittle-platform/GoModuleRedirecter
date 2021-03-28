package configuration

import (
	"redirecter/configuration/changes"
	"redirecter/modules"
	"redirecter/server"

	"go.uber.org/zap"
)

type Container struct {
	Notifier changes.ConfigurationChangeNotifier

	Resolver  modules.Resolver
	Writer    modules.Writer
	Responder modules.Responder

	Server server.Server
}

func NewContainer(config Configuration) (*Container, error) {
	logger, _ := zap.NewDevelopment()
	container := Container{}

	container.Notifier = changes.NewConfigurationChangeNotifier(logger)
	config.OnChange(container.Notifier.TriggerChanged)

	container.Resolver = modules.NewResolver(config.Modules(), container.Notifier, logger)
	container.Writer = modules.NewWriter(config.Modules(), container.Notifier, logger)
	container.Responder = modules.NewResponder(config.Modules(), container.Resolver, container.Writer, logger)

	container.Server = server.NewServer(config.Server(), container.Notifier, container.Responder, logger)

	return &container, nil
}
