package cli

import (
	"fmt"

	"github.com/AlexEkdahl/govelocity/internal/process"
)

type CliInterface interface {
	StartProcess(path string, name string) (int, error)
	StopProcess(pid int) error
	RestartProcess(pid int) error
	DeleteProcess(pid int) error
	ListProcesses()
}

type Cli struct {
	Handler map[string]Command
}

type Command struct {
	Action     string
	Use        string
	Short      string
	Length     int
	Validation func(args []string, n int) error
	Cmd        func(pm *process.ProcessManager, n []string) error
}

func createHandler() map[string]Command {
	m := make(map[string]Command)

	startCmd := Command{
		Action:     "start",
		Use:        "start <name> <path>",
		Short:      "Add and start a new process",
		Length:     4,
		Cmd:        startProcess,
		Validation: validateArgs,
	}

	stopCmd := Command{
		Action:     "stop",
		Use:        "start <name>",
		Short:      "Stop a process",
		Length:     3,
		Cmd:        stopProcess,
		Validation: validateArgs,
	}

	m[startCmd.Action] = startCmd
	m[stopCmd.Action] = stopCmd

	return m
}

func validateArgs(args []string, n int) error {
	if len(args) != n {
		return fmt.Errorf("expected exactly %d arguments, but got %d", n, len(args))
	}
	return nil
}

func NewCli() *Cli {
	return &Cli{
		Handler: createHandler(),
	}
}

func startProcess(pm *process.ProcessManager, args []string) error {
	name := args[2]
	path := args[3]
	if err := pm.Add(name, path); err != nil {
		return err
	}
	return nil
}

func stopProcess(pm *process.ProcessManager, args []string) error {
	name := args[2]
	if err := pm.Stop(name); err != nil {
		return err
	}
	return nil
}
