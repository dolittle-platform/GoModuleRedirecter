package viper

import (
	"redirecter/modules"

	"github.com/spf13/viper"
)

const (
	modulesTemplatesGoGetKey    = "modules.templates.go-get"
	modulesTemplatesUserKey     = "modules.templates.user"
	modulesTemplatesNotFoundKey = "modules.templates.not-found"
	modulesMappingsKey          = "modules.mappings"

	defaultModulesTemplatesGoGet    = "/var/lib/redirecter/go-get.html"
	defaultModulesTemplatesUser     = "/var/lib/redirecter/user.html"
	defaultModulesTemplatesNotFound = "/var/lib/redirecter/not-found.html"
)

var (
	defaultModulesMappings = make(map[string]modules.Repository)
)

type modulesConfiguration struct{}

type moduleMappingEntry struct {
	Module string
	Type   string
	Source string
}

func (c *modulesConfiguration) TemplateGoGetPath() string {
	if path := viper.GetString(modulesTemplatesGoGetKey); path != "" {
		return path
	}
	return defaultModulesTemplatesGoGet
}

func (c *modulesConfiguration) TemplateUserPath() string {
	if path := viper.GetString(modulesTemplatesUserKey); path != "" {
		return path
	}
	return defaultModulesTemplatesUser
}

func (c *modulesConfiguration) TemplateNotFoundPath() string {
	if path := viper.GetString(modulesTemplatesNotFoundKey); path != "" {
		return path
	}
	return defaultModulesTemplatesNotFound
}

func (c *modulesConfiguration) Mappings() modules.ModuleToRepositoryMappings {
	mappingEntries := []moduleMappingEntry{}
	if err := viper.UnmarshalKey(modulesMappingsKey, &mappingEntries); err != nil {
		return defaultModulesMappings
	}

	mappings := make(map[string]modules.Repository)
	for _, mappingEntry := range mappingEntries {
		mappings[mappingEntry.Module] = modules.Repository{
			Type:   mappingEntry.Type,
			Source: mappingEntry.Source,
		}
	}
	return mappings
}
