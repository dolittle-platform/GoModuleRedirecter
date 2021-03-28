package server

import (
	"net/http"
	"redirecter/modules"
	"redirecter/server/correlation"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	ProxyHostHeader = "X-Forwarded-Host"
	GoGetQueryName  = "go-get"
)

func newHandler(configuration Configuration, responder modules.Responder, logger *zap.Logger) http.Handler {
	return &handler{
		configuration: configuration,
		responder:     responder,
		logger:        logger,
	}
}

type handler struct {
	configuration Configuration
	responder     modules.Responder
	logger        *zap.Logger
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	corr := uuid.New().String()
	url := h.getFullRequestURL(r)
	h.logger.Info("Handling request", zap.String("correlation", corr), zap.String("url", url))

	defer h.recoverPanic(w, r, corr)

	ctx := correlation.ContextWithCorrelation(r.Context(), corr)

	isGoGetRequest := h.isGoGetRequest(r)

	if err := h.responder.RespondTo(w, url, isGoGetRequest, ctx); err != nil {
		h.logger.Error("Failed to respond to request", zap.Error(err), zap.String("correlation", corr), zap.String("url", url))
	}
}

func (h *handler) getFullRequestURL(r *http.Request) string {
	host := r.Host
	if forwardedHost := r.Header.Get(ProxyHostHeader); forwardedHost != "" {
		host = forwardedHost
	}
	return host + r.URL.Path
}

func (h *handler) isGoGetRequest(r *http.Request) bool {
	return r.URL.Query().Get(GoGetQueryName) == "1"
}

func (h *handler) recoverPanic(w http.ResponseWriter, r *http.Request, correlation string) {
	if err := recover(); err != nil {
		h.logger.Error("Recovered from request panic", zap.String("correlation", correlation), zap.Reflect("error", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
