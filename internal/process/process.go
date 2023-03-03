package process

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Process represents a process that can be managed by a process manager.
type Process struct {
	Name          string
	Command       string
	Args          []string
	Env           []string
	Dir           string
	Autostart     bool
	Autorestart   bool
	StartRetries  int
	StartDelay    time.Duration
	StartTimeout  time.Duration
	StopSignal    string
	StopTimeout   time.Duration
	RestartPolicy RestartPolicy

	cmd        *exec.Cmd
	startTime  time.Time
	stopTime   time.Time
	exitStatus error
	lock       sync.RWMutex
	Pid        int
}

// Status represents the status of a process.
type Status string

const (
	StatusUnknown  Status = "unknown"
	StatusStarting Status = "starting"
	StatusRunning  Status = "running"
	StatusStopped  Status = "stopped"
	StatusFailed   Status = "failed"
)

// RestartPolicy represents the policy for restarting a process.
type RestartPolicy string

const (
	RestartAlways    RestartPolicy = "always"
	RestartOnFailure RestartPolicy = "on-failure"
	RestartNever     RestartPolicy = "never"
)

// Start starts the process.
func (p *Process) Start() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	// Check if the process is already running
	if p.IsRunning() {
		p.lock.RUnlock()
		return errors.New("process already running")
	}

	// Create the command
	p.cmd = exec.Command(p.Command, p.Args...)
	p.cmd.Dir = p.Dir
	p.cmd.Env = p.Env

	// Create pipes for capturing output
	stdoutPipe, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := p.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	err = p.cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start process: %v", err)
	}
	p.Pid = p.cmd.Process.Pid // set the PID
	p.startTime = time.Now()

	// Read the output of the process
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, io.MultiReader(stdoutPipe, stderrPipe)); err != nil {
			p.lock.Lock()
			p.exitStatus = fmt.Errorf("failed to read process output: %v", err)
			p.lock.Unlock()
			return
		}
		p.lock.Lock()
		p.exitStatus = fmt.Errorf("process exited with status: %v", p.cmd.Wait())
		p.lock.Unlock()
	}()

	// Wait for the process to start
	status := StatusStarting
	ticker := time.NewTicker(p.StartDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := p.checkStatus()
			if err == nil {
				return nil
			}
			if status != StatusRunning && status != StatusStarting {
				p.exitStatus = err
				return err
			}
			status = StatusRunning
			ticker.Stop()
		case <-time.After(p.StartTimeout):
			p.Stop()
			p.exitStatus = errors.New("timeout waiting for process to start")
			return p.exitStatus
		}
	}
}

// Wait waits for the process to exit.
func (p *Process) Wait() {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Wait()
	}
}

// Stop stops the process.
func (p *Process) Stop() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	// Check if the process is running
	if !p.IsRunning() {
		return errors.New("process not running")
	}

	// Send the stop signal to the process
	signal := getSignal(p.StopSignal)
	if signal == syscall.SIGKILL {
		p.cmd.Process.Kill()
	} else {
		p.cmd.Process.Signal(signal)
	}

	// Wait for the process to stop
	select {
	case <-p.done():
		p.stopTime = time.Now()
		return nil
	case <-time.After(p.StopTimeout):
		// If we're still waiting, kill the process
		p.cmd.Process.Kill()
		p.exitStatus = errors.New("timeout waiting for process to stop")
		return p.exitStatus
	}
}

// Restart restarts the process.
func (p *Process) Restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	return p.Start()
}

// IsRunning returns true if the process is running.
func (p *Process) IsRunning() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.cmd == nil || p.cmd.Process == nil {
		return false
	}

	if p.cmd.ProcessState != nil && p.cmd.ProcessState.Exited() {
		return false
	}

	return true
}

// Status returns the status of the process.
func (p *Process) Status() Status {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.cmd == nil {
		return StatusUnknown
	}

	if p.exitStatus != nil {
		return StatusFailed
	}

	if !p.IsRunning() {
		return StatusStopped
	}

	return StatusRunning
}

// Output returns the combined stdout and stderr output of the process.
func (p *Process) Output() (string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if !p.IsRunning() {
		return "", errors.New("process not running")
	}

	var buf bytes.Buffer
	stdoutPipe, err := p.cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := p.cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe: %v", err)
	}
	_, stdoutWriter := io.Pipe()
	_, stderrWriter := io.Pipe()

	go func() {
		io.Copy(io.MultiWriter(&buf, stdoutWriter), stdoutPipe)
		stdoutWriter.Close()
	}()
	go func() {
		io.Copy(io.MultiWriter(&buf, stderrWriter), stderrPipe)
		stderrWriter.Close()
	}()

	if err := p.cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start process: %v", err)
	}
	p.startTime = time.Now()

	if err := p.cmd.Wait(); err != nil {
		return "", fmt.Errorf("process exited with status code: %d", p.cmd.ProcessState.ExitCode())
	}

	return strings.TrimSpace(buf.String()), nil
}

// ExitStatus returns the exit status of the process.
func (p *Process) ExitStatus() error {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.exitStatus
}

func (p *Process) done() chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		p.cmd.Wait()
	}()
	return done
}

func (p *Process) checkStatus() error {
	if p.cmd.ProcessState != nil && p.cmd.ProcessState.Exited() {
		if p.cmd.ProcessState.Success() {
			return nil
		}
		if p.RestartPolicy == RestartAlways || (p.RestartPolicy == RestartOnFailure && p.cmd.ProcessState.ExitCode() != 0) {
			return fmt.Errorf("process '%s' exited with error: %v", p.Name, p.cmd.ProcessState)
		}
		return fmt.Errorf("process '%s' exited with status code: %d", p.Name, p.cmd.ProcessState.ExitCode())
	}
	return fmt.Errorf("process '%s' is still starting", p.Name)
}

func NewProcess(cmd []string, name string) *Process {
	a := exec.Command(name, cmd...)
	return &Process{
		cmd:  a,
		Name: name,
	}
}

func getSignal(signal string) syscall.Signal {
	switch signal {
	case "SIGINT":
		return syscall.SIGINT
	case "SIGQUIT":
		return syscall.SIGQUIT
	case "SIGTERM":
		return syscall.SIGTERM
	case "SIGKILL":
		return syscall.SIGKILL
	default:
		return syscall.SIGTERM
	}
}
