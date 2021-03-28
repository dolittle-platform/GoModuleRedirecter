package configuration

import (
	"redirecter/server"

	"go.uber.org/zap"
)

type Container struct {
	Server server.Server
}

func NewContainer(config Configuration) (*Container, error) {
	logger, _ := zap.NewDevelopment()
	container := Container{}

	container.Server = server.NewServer(config.Server(), logger)

	return &container, nil
}
