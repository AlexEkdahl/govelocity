package process

import (
	"fmt"
	"os"
	"os/exec"
)

type Process struct {
	PID  int    `yaml:"pid"`
	Path string `yaml:"path"`
	Name string `yaml:"name"`
	cmd  *exec.Cmd
}

func NewProcess(name, path string) *Process {
	return &Process{
		Path: path,
		Name: name,
		cmd: &exec.Cmd{
			Path: path,
			Args: []string{name},
		},
	}
}

func (p *Process) start() error {
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("error starting process: %w", err)
	}
	p.PID = p.cmd.Process.Pid

	return nil
}

func RestoreProcess(p *Process) *Process {
	process, err := os.FindProcess(p.PID)
	if err != nil {
		return NewProcess(p.Name, p.Path)
	}

	return &Process{
		Path: p.Path,
		Name: p.Name,
		PID:  process.Pid,
		cmd: &exec.Cmd{
			Path:    p.Path,
			Args:    []string{p.Name},
			Process: process,
		},
	}
}

func (p *Process) stop() error {
	if err := p.cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("error stopping process: %s\n", err)
	}

	return nil
}

func (p *Process) String() string {
	return fmt.Sprintf("PID: %d, Path: %s, Name: %s", p.PID, p.Path, p.Name)
}
