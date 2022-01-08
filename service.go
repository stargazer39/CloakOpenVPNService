package main

import (
	"log"
	"os"
	"os/exec"
)

type Service struct {
	proc     exec.Cmd
	name     string
	args     []string
	Running  bool
	waitChan chan bool
}

func NewService(name string, args ...string) *Service {
	return &Service{
		proc:     *exec.Command(name, args...),
		Running:  false,
		name:     name,
		args:     args,
		waitChan: make(chan bool),
	}
	// s.proc.Command(name)
}

func (s *Service) Start() error {
	if !s.Running {
		if err := s.proc.Start(); err != nil {
			return err
		}
		go func() {
			s.Running = true
			s.proc.Wait()
			s.waitChan <- true
			log.Printf("Process %s ended\n", s.name)
			s.Running = false
			s = NewService(s.name, s.args...)
		}()
	} else {
		log.Printf("Process %s already started\n", s.name)
	}
	return nil
}

func (s *Service) Kill() {
	if s.proc.Process == nil {
		return
	}
	if sErr := s.proc.Process.Signal(os.Interrupt); sErr != nil {
		log.Println(sErr)
		if sErr == os.ErrProcessDone {
			return
		}
	}
	<-s.waitChan
}
