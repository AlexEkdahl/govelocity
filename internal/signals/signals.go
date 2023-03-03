package signals

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexEkdahl/govelocity/internal/process"
)

// Signals represents the signals received by the process manager.
type Signals struct {
	manager *process.Manager
}

// New creates a new signals object with the specified process manager.
func New(manager *process.Manager) *Signals {
	return &Signals{manager: manager}
}

// Start starts listening for signals and handles them.
func (s *Signals) Start() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for sig := range sigChan {
		log.Printf("Received signal: %v", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			go s.manager.StopAll()
		}
	}
}
