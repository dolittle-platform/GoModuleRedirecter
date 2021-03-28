package modules

type ModuleToRepositoryMappings map[string]Repository

type Configuration interface {
	Documentation() string
	TemplateGoGetPath() string
	TemplateUserPath() string
	TemplateNotFoundPath() string
	Mappings() ModuleToRepositoryMappings
}
