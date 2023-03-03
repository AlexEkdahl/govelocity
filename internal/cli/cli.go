package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlexEkdahl/govelocity/internal/process"
)

// HandlerFunc is a function that handles a command.
type HandlerFunc func(manager *process.Manager, args []string, name string) error

// Command represents a command that can be executed by the CLI.
type Command struct {
	Name        string
	Description string
	Usage       string
	Handler     HandlerFunc
}

// CLI represents a command-line interface.
type CLI struct {
	manager  *process.Manager
	commands []*Command
}

// New creates a new CLI instance.
func New(manager *process.Manager) *CLI {
	return &CLI{
		manager: manager,
	}
}

// RegisterCommand registers a new command.
func (c *CLI) RegisterCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

// Run runs the CLI with the given arguments.
func (c *CLI) Run(command string, args []string, name string) error {
	// Find the command to execute
	var cmd *Command
	fmt.Println("c.commands", c.commands)
	for _, c := range c.commands {
		fmt.Println("c.Name", c.Name)
		fmt.Println("command", command)
		if c.Name == command {
			cmd = c
			break
		}
	}
	if cmd == nil {
		return fmt.Errorf("unknown command: %s", command)
	}

	// Execute the command
	return cmd.Handler(c.manager, args, name)
}

// Help prints help information for the CLI.
func (c *CLI) Help() string {
	var sb strings.Builder
	sb.WriteString("Usage: process-managerctl [options] <command> [arguments]\n\n")
	sb.WriteString("Options:\n")
	sb.WriteString("  --help: Print the help message\n")
	sb.WriteString("  --name: Name of the process\n\n")
	sb.WriteString("Commands:\n")
	for _, cmd := range c.commands {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", cmd.Name, cmd.Description))
		sb.WriteString(fmt.Sprintf("    Usage: %s\n\n", cmd.Usage))
	}
	return sb.String()
}

// StartCommandHandler starts a process.
func StartCommandHandler(manager *process.Manager, args []string, name string) error {
	if len(args) == 0 {
		return errors.New("missing command to start")
	}

	p, err := manager.Get(name)
	if err != nil {
		return err
	}

	p.Command = args[0]
	p.Args = args[1:]
	return p.Start()
}

// StopCommandHandler stops a process.
func StopCommandHandler(manager *process.Manager, args []string, name string) error {
	p, err := manager.Get(name)
	if err != nil {
		return err
	}

	return p.Stop()
}

// RestartCommandHandler restarts a process.
func RestartCommandHandler(manager *process.Manager, args []string, name string) error {
	p, err := manager.Get(name)
	if err != nil {
		return err
	}

	return p.Restart()
}

// StatusCommandHandler prints the status of a process.
func StatusCommandHandler(manager *process.Manager, args []string, name string) error {
	p, err := manager.Get(name)
	if err != nil {
		return err
	}

	status := p.Status()
	fmt.Printf("%s (%d): %s\n", name, p.Pid, status)
	return nil
}

// StartCommandHandler handles the "start" command.
func (c *CLI) StartCommandHandler(args []string, name string) error {
	if len(args) < 1 {
		return errors.New("missing command argument")
	}

	process, err := c.manager.GetProcess(name)
	if err != nil {
		return fmt.Errorf("process not found: %s", name)
	}

	if err := process.Start(); err != nil {
		return fmt.Errorf("failed to start process: %v", err)
	}

	fmt.Printf("Process started: %s (PID: %d)\n", process.Name, process.Pid)

	return nil
}

// StopCommandHandler handles the "stop" command.
func (c *CLI) StopCommandHandler(args []string, name string) error {
	if len(args) < 1 {
		return errors.New("missing command argument")
	}

	process, err := c.manager.GetProcess(name)
	if err != nil {
		return fmt.Errorf("process not found: %s", name)
	}

	if err := process.Stop(); err != nil {
		return fmt.Errorf("failed to stop process: %v", err)
	}

	fmt.Printf("Process stopped: %s (PID: %d)\n", process.Name, process.Pid)

	return nil
}

// RestartCommandHandler handles the "restart" command.
func (c *CLI) RestartCommandHandler(args []string, name string) error {
	if len(args) < 1 {
		return errors.New("missing command argument")
	}

	process, err := c.manager.GetProcess(name)
	if err != nil {
		return fmt.Errorf("process not found: %s", name)
	}

	if err := process.Restart(); err != nil {
		return fmt.Errorf("failed to restart process: %v", err)
	}

	fmt.Printf("Process restarted: %s (PID: %d)\n", process.Name, process.Pid)

	return nil
}

// StatusCommandHandler handles the "status" command.
func (c *CLI) StatusCommandHandler(args []string, name string) error {
	if len(args) > 1 {
		return errors.New("too many arguments")
	}

	processes := c.manager.ListProcesses()
	if len(processes) == 0 {
		fmt.Println("No processes found")
		return nil
	}

	if name == "" {
		// Print status for all processes
		for _, process := range processes {
			fmt.Printf("%s: %s\n", process.Name, process.Status())
		}
	} else {
		// Print status for a specific process
		process, err := c.manager.GetProcess(name)
		if err != nil {
			return fmt.Errorf("process not found: %s", name)
		}
		fmt.Printf("%s: %s\n", process.Name, process.Status())
	}

	return nil
}
