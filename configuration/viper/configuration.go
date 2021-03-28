package viper

import (
	config "redirecter/configuration"
	"redirecter/server"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func NewViperConfiguration(configPath string) (config.Configuration, error) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".redirecter")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &configuration{
		server: &serverConfiguration{},
	}, nil
}

type configuration struct {
	server *serverConfiguration
}

func (c *configuration) Server() server.Configuration {
	return c.server
}
