package modules

import (
	"context"
	"math"
	"redirecter/configuration/changes"
	"redirecter/server/correlation"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

var majorVersionPattern = regexp.MustCompile(`^v[0-9]+`)

type Match struct {
	Exact        bool
	MajorVersion string
	Repository   *Repository
}

type Resolver interface {
	Resolve(url string, ctx context.Context) (*Match, bool, error)
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

func (r *resolver) Resolve(url string, ctx context.Context) (*Match, bool, error) {
	r.logger.Info("Resolving package", zap.String("url", url), zap.String("correlation", correlation.CorrelationFromContext(ctx)))

	var bestMatch *Match
	bestMatchRemainingLength := math.MaxInt32
	for mappingURL, repository := range r.modules {
		match, matches, remainingLength := r.matchRequestedURLToMapping(url, mappingURL)
		if matches && remainingLength < bestMatchRemainingLength {
			match.Repository = &repository
			bestMatch = match
			bestMatchRemainingLength = remainingLength
		}

	}

	if bestMatch != nil {
		r.logger.Info("Found package", zap.String("url", url), zap.String("type", bestMatch.Repository.Type), zap.String("repository", bestMatch.Repository.Source), zap.String("correlation", correlation.CorrelationFromContext(ctx)))
		return bestMatch, true, nil
	}

	return nil, false, nil
}

func (r *resolver) matchRequestedURLToMapping(url, mapping string) (*Match, bool, int) {
	mapping = strings.TrimSuffix(mapping, "/")

	if !strings.HasPrefix(url, mapping) {
		return nil, false, 0
	}

	remainder := strings.TrimPrefix(strings.TrimPrefix(url, mapping), "/")
	majorVersion := majorVersionPattern.FindString(remainder)

	remainingLength := len(remainder) - len(majorVersion)

	return &Match{
		Exact:        remainingLength == 0,
		MajorVersion: majorVersion,
	}, true, remainingLength

}
