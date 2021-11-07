package service

import (
	"os"
	"os/exec"

	"github.com/ergoapi/log"
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
	*Config
	Cmd     *exec.Cmd
	Service service.Service
	Exit    chan struct{}
}

func (es *ErgoService) Start(s service.Service) error {
	fullexec, err := exec.LookPath(es.Exec)
	if err != nil {
		return err
	}
	es.Cmd = exec.Command(fullexec, es.Args...)
	es.Cmd.Dir = es.Dir
	es.Cmd.Env = append(os.Environ(), es.Env...)
	go es.run()
	return nil
}

func (es *ErgoService) run() {
	eslog := log.GetInstance()
	eslog.Infof("start %v", es.Name)
	defer func() {
		if service.Interactive() {
			es.Stop(es.Service)
		} else {
			es.Service.Stop()
		}
	}()

	if es.Stderr != "" {
		f, err := os.OpenFile(es.Stderr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			eslog.Warnf("Failed to open std err %q: %v", es.Stderr, err)
			return
		}
		defer f.Close()
		es.Cmd.Stderr = f
	}
	if es.Stdout != "" {
		f, err := os.OpenFile(es.Stdout, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			eslog.Warnf("Failed to open std out %q: %v", es.Stdout, err)
			return
		}
		defer f.Close()
		es.Cmd.Stdout = f
	}

	err := es.Cmd.Run()
	if err != nil {
		eslog.Warnf("Error running: %v", err)
	}
}

func (es *ErgoService) Stop(s service.Service) error {
	eslog := log.GetInstance()
	close(es.Exit)
	eslog.Infof("Stopping %v", es.Name)
	if es.Cmd.Process != nil {
		es.Cmd.Process.Kill()
	}
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func New(conf *Config) (service.Service, error) {
	config := &service.Config{
		Name:             conf.Name,
		DisplayName:      conf.Name,
		Description:      conf.Desc,
		Arguments:        conf.Args,
		WorkingDirectory: conf.Dir,
		Executable:       conf.Exec,
	}
	m := new(ErgoService)
	return service.New(m, config)
}
