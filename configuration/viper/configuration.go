package viper

import (
	config "redirecter/configuration"
	"redirecter/modules"
	"redirecter/server"
	"strings"

	"github.com/fsnotify/fsnotify"
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

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.WatchConfig()

	return &configuration{
		modules: &modulesConfiguration{},
		server:  &serverConfiguration{},
	}, nil
}

type configuration struct {
	modules *modulesConfiguration
	server  *serverConfiguration
}

func (c *configuration) OnChange(callback func()) {
	viper.OnConfigChange(func(in fsnotify.Event) {
		callback()
	})
}

func (c *configuration) Modules() modules.Configuration {
	return c.modules
}

func (c *configuration) Server() server.Configuration {
	return c.server
}
