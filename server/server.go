package server

import (
	"context"
	"fmt"
	"net/http"
	"redirecter/configuration/changes"
	"redirecter/modules"
	"time"

	"go.uber.org/zap"
)

type Server interface {
	Run() error
}

func NewServer(configuration Configuration, notifier changes.ConfigurationChangeNotifier, responder modules.Responder, logger *zap.Logger) Server {
	return &server{
		configuration:    configuration,
		notifier:         notifier,
		handler:          newHandler(configuration, responder, logger),
		logger:           logger,
		shutdownComplete: make(chan struct{}),
	}
}

type server struct {
	configuration    Configuration
	notifier         changes.ConfigurationChangeNotifier
	handler          http.Handler
	logger           *zap.Logger
	httpServer       *http.Server
	shutdownComplete chan struct{}
}

func (s *server) Run() error {
	if err := s.notifier.RegisterCallback("server", s.handleConfigurationChanged); err != nil {
		return err
	}
	return s.loop()
}

func (s *server) loop() error {
	for {
		if err := s.run(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("http server failed", zap.Error(err))

			select {
			case <-s.shutdownComplete:
			case <-time.After(2 * time.Second):
			}
		}
	}
}

func (s *server) run() error {
	s.httpServer = &http.Server{}

	s.logger.Info("Starting server", zap.Int("port", s.configuration.Port()))
	s.httpServer.Addr = fmt.Sprintf(":%d", s.configuration.Port())
	s.httpServer.Handler = s.handler

	return s.httpServer.ListenAndServe()
}

func (s *server) handleConfigurationChanged() error {
	err := s.httpServer.Shutdown(context.Background())
	select {
	case s.shutdownComplete <- struct{}{}:
	default:
	}
	return err
}
