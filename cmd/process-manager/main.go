package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexEkdahl/govelocity/internal/process"
)

func main() {
	// Define command-line flags
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "Path to the configuration file")

	// Parse command-line flags
	flag.Parse()

	// Load the configuration
	// config, err := config.LoadConfig(configFile)
	// if err != nil {
	// 	log.Fatalf("Error loading configuration: %v", err)
	// }

	// Create a new process manager
	manager := process.NewManager()

	// Create a new notifications channel

	// Start the process manager
	if err := manager.Start(); err != nil {
		log.Fatalf("Error starting process manager: %v", err)
	}

	// Start the notifications handler

	// Handle signals sent to the process manager
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			// Stop the process manager and exit
			if err := manager.Stop(); err != nil {
				log.Fatalf("Error stopping process manager: %v", err)
			}
			os.Exit(0)
		}
	}
}
