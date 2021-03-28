package modules

type ModuleToRepositoryMappings map[string]Repository

type Configuration interface {
	TemplateGoGetPath() string
	TemplateUserPath() string
	TemplateNotFoundPath() string
	Mappings() ModuleToRepositoryMappings
}
