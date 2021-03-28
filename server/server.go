package server

import (
	"fmt"

	"go.uber.org/zap"
)

type Server interface {
	Run() error
}

func NewServer(configuration Configuration, logger *zap.Logger) Server {
	return &server{}
}

type server struct {
}

func (s *server) Run() error {
	fmt.Println("Running server")
	return nil
}
