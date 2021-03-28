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

func NewResponder(resolver Resolver, writer Writer, logger *zap.Logger) Responder {
	return &responder{
		resolver: resolver,
		writer:   writer,
		logger:   logger,
	}
}

type responder struct {
	resolver Resolver
	writer   Writer
	logger   *zap.Logger
}

func (r *responder) RespondTo(w http.ResponseWriter, url string, goGet bool, ctx context.Context) error {
	repository, found, err := r.resolver.Resolve(url, ctx)
	if err != nil {
		r.logger.Error("Failed to resolve", zap.Error(err), zap.String("url", url), zap.String("correlation", correlation.CorrelationFromContext(ctx)))
		return err
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		r.writer.WriteNotFoundResponse(w, url, ctx)
		return nil
	}

	if goGet {
		w.WriteHeader(http.StatusOK)
		r.writer.WriteGoGetResponse(w, url, repository, ctx)
	} else {
		w.WriteHeader(http.StatusOK)
		r.writer.WriteUserResponse(w, url, repository, ctx)
	}

	return nil
}
