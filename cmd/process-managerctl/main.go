package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlexEkdahl/govelocity/internal/cli"
	"github.com/AlexEkdahl/govelocity/internal/process"
)

func main() {
	// Define command-line flags
	var (
		name string
	)

	flag.StringVar(&name, "name", "status", "Name of the process to operate on")

	// Parse command-line flags
	flag.Parse()

	// Create a new process manager
	manager := process.NewManager()

	// Initialize the CLI with the process manager
	c := cli.New(manager)
	c.RegisterCommand(&cli.Command{
		Name:        "start",
		Description: "Start a process",
		Usage:       "start <command> [args...]",
		Handler:     handler,
	})

	c.RegisterCommand(&cli.Command{
		Name:        "stop",
		Description: "Stop a process",
		Usage:       "stop <name>",
		Handler:     stop,
	})

	c.RegisterCommand(&cli.Command{
		Name:        "restart",
		Description: "Restart a process",
		Usage:       "restart <name>",
		Handler:     restart,
	})
	c.RegisterCommand(&cli.Command{
		Name:        "status",
		Description: "Get the status of a process",
		Usage:       "status <name>",
		Handler:     status,
	})
	c.RegisterCommand(&cli.Command{
		Name:        "list",
		Description: "List all processes",
		Usage:       "list",
		Handler:     list,
	})
	c.RegisterCommand(&cli.Command{
		Name:        "add",
		Description: "List all processes",
		Usage:       "add",
		Handler:     add,
	})

	// Parse the command and execute it
	if len(flag.Args()) < 1 {
		fmt.Println("manager")
		// Print the help message if no command is specified
		fmt.Println(c.Help())
		os.Exit(0)
	} else {
		args := flag.Args()
		fmt.Println("args", args)
		fmt.Println("name", name)
		if err := c.Run(name, args, name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func handler(manager *process.Manager, args []string, name string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing command to start")
	}
	// manager.AddProcess(name)
	err := manager.StartProcess(name)
	if err != nil {
		return err
	}

	fmt.Printf("Started process %s \n", name)

	return nil
}

func add(manager *process.Manager, args []string, name string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing command to start")
	}

	p := process.NewProcess(args, name)
	err := manager.AddProcess(p)
	if err != nil {
		return err
	}

	fmt.Printf("Started process %s \n", name)

	return nil
}

func stop(manager *process.Manager, args []string, name string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing process name to stop")
	}

	p, err := manager.GetProcess(name)
	if err != nil {
		return err
	}

	if err := p.Stop(); err != nil {
		return err
	}

	fmt.Printf("Stopped process %s (pid %d)\n", p.Name, p.Pid)

	return nil
}

func restart(manager *process.Manager, args []string, name string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing process name to restart")
	}

	p, err := manager.GetProcess(name)
	if err != nil {
		return err
	}

	if err := p.Restart(); err != nil {
		return err
	}

	fmt.Printf("Restarted process %s (pid %d)\n", p.Name, p.Pid)
	return nil
}

func status(manager *process.Manager, args []string, name string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing process name to get status")
	}

	p, err := manager.GetProcess(name)
	if err != nil {
		return err
	}

	status := p.Status()
	fmt.Printf("Process %s (pid %d) is %s\n", p.Name, p.Pid, status)

	return nil
}

func list(manager *process.Manager, args []string, name string) error {
	processes := manager.ListProcesses()

	if len(processes) == 0 {
		fmt.Println("No processes currently running")
		return nil
	}

	fmt.Println("Running processes:")
	for _, p := range processes {
		fmt.Printf("  %s (pid %d): %s\n", p.Name, p.Pid, p.Status)
	}

	return nil
}
