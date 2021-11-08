package service

import (
	"context"
	"os"

	"github.com/kardianos/service"
)

// Config configures the service.
type Config struct {
	Name string // service name
	Desc string // service description
	Dir  string
	Exec string
	Args []string
	Env  []string

	Stderr, Stdout string
}

type ErgoService struct {
	cancel context.CancelFunc
}

var nocontext = context.Background()

func (es *ErgoService) Start(s service.Service) error {
	_, cancel := context.WithCancel(nocontext)
	es.cancel = cancel
	return nil
}

func (es *ErgoService) Stop(s service.Service) error {
	es.cancel()
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}
