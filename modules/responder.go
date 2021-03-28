package modules

import (
	"context"
	"net/http"
	"redirecter/server/correlation"

	"go.uber.org/zap"
)

type Responder interface {
	RespondTo(w http.ResponseWriter, url string, goGet bool, ctx context.Context) error
}

func NewResponder(configuration Configuration, resolver Resolver, writer Writer, logger *zap.Logger) Responder {
	return &responder{
		configuration: configuration,
		resolver:      resolver,
		writer:        writer,
		logger:        logger,
	}
}

type responder struct {
	configuration Configuration
	resolver      Resolver
	writer        Writer
	logger        *zap.Logger
}

func (r *responder) RespondTo(w http.ResponseWriter, url string, goGet bool, ctx context.Context) error {
	match, found, err := r.resolver.Resolve(url, ctx)
	if err != nil {
		r.logger.Error("Failed to resolve", zap.Error(err), zap.String("url", url), zap.String("correlation", correlation.CorrelationFromContext(ctx)))
		return err
	}

	if !found || (goGet && !match.Exact) {
		w.WriteHeader(http.StatusNotFound)
		r.writer.WriteNotFoundResponse(w, url, ctx)
		return nil
	}

	if goGet {
		w.WriteHeader(http.StatusOK)
		r.writer.WriteGoGetResponse(w, url, match, ctx)
	} else {
		w.Header().Set("Location", r.configuration.Documentation()+url)
		w.WriteHeader(http.StatusFound)
		r.writer.WriteUserResponse(w, url, match, ctx)
	}

	return nil
}
