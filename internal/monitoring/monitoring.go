package monitoring

import (
	"context"
	"time"

	"github.com/AlexEkdahl/govelocity/internal/process"
)

// Monitor is a process monitor that periodically checks the status of a process and reports any changes.
type Monitor struct {
	ProcessManager *process.Manager
	ProcessName    string
	Interval       time.Duration
	Handler        func(status process.Status)
	stopCh         chan struct{}
}

// Start starts the monitor and begins monitoring the process.
func (m *Monitor) Start() {
	m.stopCh = make(chan struct{})
	go m.monitorLoop()
}

// Stop stops the monitor.
func (m *Monitor) Stop() {
	close(m.stopCh)
}

// monitorLoop is the main loop for the monitor.
func (m *Monitor) monitorLoop() {
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	status := m.ProcessManager.Status(m.ProcessName)
	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			newStatus := m.ProcessManager.Status(m.ProcessName)
			if newStatus != status {
				status = newStatus
				m.Handler(status)
			}
		case <-ctx.Done():
			return
		}
	}
}
