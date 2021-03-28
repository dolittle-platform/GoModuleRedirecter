package server

import (
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	ProxyHostHeader = "X-Forwarded-Host"
	GoGetQueryName  = "go-get"
)

func newHandler(configuration Configuration, logger *zap.Logger) http.Handler {
	return &handler{
		configuration: configuration,
		logger:        logger,
	}
}

type handler struct {
	configuration Configuration
	logger        *zap.Logger
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	correlation := uuid.New().String()
	url := h.getFullRequestURL(r)
	h.logger.Info("Handling request", zap.String("correlation", correlation), zap.String("url", url))

	defer h.recoverPanic(w, r, correlation)

	h.isGoGetRequest(r)
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
