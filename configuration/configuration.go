package configuration

import (
	"redirecter/server"
)

type Configuration interface {
	OnChange(callback func())

	Server() server.Configuration
}
