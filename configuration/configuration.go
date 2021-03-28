package configuration

import (
	"redirecter/modules"
	"redirecter/server"
)

type Configuration interface {
	OnChange(callback func())

	Modules() modules.Configuration
	Server() server.Configuration
}
