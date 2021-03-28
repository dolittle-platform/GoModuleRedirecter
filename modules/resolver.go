package modules

import (
	"context"
	"redirecter/configuration/changes"
	"redirecter/server/correlation"

	"go.uber.org/zap"
)

type Resolver interface {
	Resolve(url string, ctx context.Context) (*Repository, bool, error)
}

func NewResolver(configuration Configuration, notifier changes.ConfigurationChangeNotifier, logger *zap.Logger) Resolver {
	r := &resolver{
		configuration: configuration,
		logger:        logger,
	}
	r.reloadModulesFromConfiguration()
	notifier.RegisterCallback("resolver", r.reloadModulesFromConfiguration)
	return r
}

type resolver struct {
	configuration Configuration
	modules       ModuleToRepositoryMappings
	logger        *zap.Logger
}

func (r *resolver) reloadModulesFromConfiguration() error {
	r.modules = r.configuration.Mappings()
	r.logger.Info("Resolver configured with module mappings")
	for url, repository := range r.modules {
		r.logger.Info("Mapping", zap.String("url", url), zap.String("type", repository.Type), zap.String("repository", repository.Source))
	}
	return nil
}

func (r *resolver) Resolve(url string, ctx context.Context) (*Repository, bool, error) {
	r.logger.Info("Resolving package", zap.String("url", url), zap.String("correlation", correlation.CorrelationFromContext(ctx)))

	for mappingURL, repository := range r.modules {
		if url == mappingURL {
			r.logger.Info("Found package", zap.String("url", url), zap.String("type", repository.Type), zap.String("repository", repository.Source), zap.String("correlation", correlation.CorrelationFromContext(ctx)))
			return &repository, true, nil
		}
	}

	return nil, false, nil
}
