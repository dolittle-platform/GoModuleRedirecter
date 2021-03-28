package configuration

import (
	"redirecter/configuration/changes"
	"redirecter/server"

	"go.uber.org/zap"
)

type Container struct {
	Notifier changes.ConfigurationChangeNotifier

	Server server.Server
}

func NewContainer(config Configuration) (*Container, error) {
	logger, _ := zap.NewDevelopment()
	container := Container{}

	container.Notifier = changes.NewConfigurationChangeNotifier(logger)
	config.OnChange(container.Notifier.TriggerChanged)

	container.Server = server.NewServer(config.Server(), container.Notifier, logger)

	return &container, nil
}
