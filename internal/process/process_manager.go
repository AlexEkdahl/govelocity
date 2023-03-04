package process

import (
	"fmt"
)

type Storer interface {
	Open() error
	Close() error
	CreateTable() error
	InsertProcess(p *Process) error
	GetProcesses() ([]*Process, error)
	RemoveProcess(pid int) error
}

type ProcessManager struct {
	processes map[string]*Process
	storer    Storer
}

func NewProcessManager(db Storer) *ProcessManager {
	manager := &ProcessManager{
		processes: restoreProcesses(db),
		storer:    db,
	}

	return manager
}

func restoreProcesses(db Storer) map[string]*Process {
	processes := make(map[string]*Process)

	// Retrieve processes from database
	processList, err := db.GetProcesses()
	if err != nil {
		return processes
	}

	// Convert process list to map
	for _, p := range processList {
		processes[p.Name] = RestoreProcess(p)
	}

	return processes
}

func (pm *ProcessManager) Add(name, path string) error {
	p := NewProcess(name, path)

	if err := p.start(); err != nil {
		return fmt.Errorf("error starting process: %w", err)
	}

	pm.processes[name] = p

	if err := pm.storer.InsertProcess(p); err != nil {
		return fmt.Errorf("error saving process: %w", err)
	}

	return nil
}

func (pm *ProcessManager) Start(name string) error {
	p := pm.processes[name]
	if err := p.start(); err != nil {
		return err
	}
	return nil
}

func (pm *ProcessManager) List() {
	// TODO
	fmt.Println("pm.processes", pm.processes)
}

func (pm *ProcessManager) Remove(name string) error {
	p, ok := pm.processes[name]
	if !ok {
		return fmt.Errorf("process with name %s not found", name)
	}

	// Kill the process if it is still running
	if p.cmd.ProcessState == nil {
		if err := p.stop(); err != nil {
			fmt.Println("err", err)
			// return err
		}
	}

	// Remove the process from the map
	delete(pm.processes, name)

	// Remove the process from the database
	if err := pm.storer.RemoveProcess(p.PID); err != nil {
		return err
	}

	return nil
}

func (pm *ProcessManager) Stop(name string) error {
	p := pm.processes[name]
	if err := p.stop(); err != nil {
		return err
	}

	return nil
}
