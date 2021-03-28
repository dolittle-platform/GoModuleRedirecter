package viper

import (
	"github.com/spf13/viper"
)

const (
	servePortKey = "serve.port"

	defaultServePort = 8080
)

type serverConfiguration struct{}

func (c *serverConfiguration) Port() int {
	port := viper.GetInt(servePortKey)
	if port == 0 {
		return defaultServePort
	}
	return port
}
