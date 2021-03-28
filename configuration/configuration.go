package configuration

import (
	"redirecter/server"
)

type Configuration interface {
	Server() server.Configuration
}
