package process

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Manager is a process manager that manages a collection of processes.
type Manager struct {
	processes map[string]*Process
	lock      sync.RWMutex
}

// NewManager creates a new process manager.
func NewManager() *Manager {
	return &Manager{
		processes: make(map[string]*Process),
	}
}

// AddProcess adds a process to the process manager.
func (m *Manager) AddProcess(p *Process) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.processes[p.Name]; ok {
		return fmt.Errorf("process '%s' already exists", p.Name)
	}

	m.Start()
	m.processes[p.Name] = p
	return nil
}

// RemoveProcess removes a process from the process manager.
func (m *Manager) RemoveProcess(name string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.processes[name]; !ok {
		return fmt.Errorf("process '%s' does not exist", name)
	}

	delete(m.processes, name)
	return nil
}

// Start starts all the processes managed by the process manager.
func (m *Manager) Start() error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Start each process
	for _, p := range m.processes {
		if p.Autostart {
			if err := p.Start(); err != nil {
				return fmt.Errorf("failed to start process '%s': %v", p.Name, err)
			}
		}
	}

	// Wait for signals to stop the process manager
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	// Stop each process
	for _, p := range m.processes {
		if p.IsRunning() {
			if err := p.Stop(); err != nil {
				return fmt.Errorf("failed to stop process '%s': %v", p.Name, err)
			}
		}
	}

	return nil
}

// Status returns the status of the process with the given name.
func (m *Manager) Status(name string) Status {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if p, ok := m.processes[name]; ok {
		return p.Status()
	}

	return StatusUnknown
}

// Statuses returns the statuses of all the processes managed by the process manager.
func (m *Manager) Statuses() map[string]Status {
	m.lock.RLock()
	defer m.lock.RUnlock()

	statuses := make(map[string]Status)
	for name, p := range m.processes {
		statuses[name] = p.Status()
	}

	return statuses
}

// Stop stops all the processes managed by the process manager.
func (m *Manager) Stop() error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Stop each process
	for _, p := range m.processes {
		if p.IsRunning() {
			if err := p.Stop(); err != nil {
				return fmt.Errorf("failed to stop process '%s': %v", p.Name, err)
			}
		}
	}

	return nil
}

// Wait waits for all the processes managed by the process manager to exit.
func (m *Manager) Wait() {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, p := range m.processes {
		p.Wait()
	}
}

// Restart restarts all the processes managed by the process manager.
func (m *Manager) Restart() error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Stop each process
	for _, p := range m.processes {
		if p.IsRunning() {
			if err := p.Restart(); err != nil {
				return fmt.Errorf("failed to restart process '%s': %v", p.Name, err)
			}
		}
	}

	return nil
}

// StopProcess stops a single process managed by the process manager.
func (m *Manager) StopProcess(name string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if p, ok := m.processes[name]; ok {
		return p.Stop()
	}

	return fmt.Errorf("process '%s' not found", name)
}

// StartProcess starts a single process managed by the process manager.
func (m *Manager) StartProcess(name string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if p, ok := m.processes[name]; ok {
		return p.Start()
	}

	return fmt.Errorf("process '%s' not found", name)
}

// StopAll stops all processes managed by the process manager.
func (m *Manager) StopAll() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, p := range m.processes {
		if p.IsRunning() {
			_ = p.Stop()
		}
	}
}

// Get returns the process with the given name, or an error if it doesn't exist.
func (m *Manager) Get(name string) (*Process, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	p, ok := m.processes[name]
	if !ok {
		return nil, fmt.Errorf("process '%s' not found", name)
	}
	return p, nil
}

// GetProcess returns the process with the given name, or an error if it doesn't exist.
func (m *Manager) GetProcess(name string) (*Process, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	p, ok := m.processes[name]
	if !ok {
		return nil, fmt.Errorf("process '%s' not found", name)
	}
	return p, nil
}

// ListProcesses returns a slice of all the processes managed by the process manager.
func (m *Manager) ListProcesses() []*Process {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var processes []*Process
	for _, p := range m.processes {
		processes = append(processes, p)
	}
	return processes
}

// RestartProcess restarts a single process managed by the process manager.
func (m *Manager) RestartProcess(name string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if p, ok := m.processes[name]; ok {
		return p.Restart()
	}

	return fmt.Errorf("process '%s' not found", name)
}
