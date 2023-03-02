package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Push []Command `yaml:"push"`
}

type Command struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

func main() {
	configFile := flag.String("config", "config.yml", "path to configuration file")
	event := flag.String("event", "push", "GitHub webhook event type")
	// url := flag.String("url", "", "GitHub webhook URL")
	// secret := flag.String("secret", "", "GitHub webhook secret key")
	addCommand := flag.Bool("add", false, "add a new command")
	flag.Parse()

	config, err := readConfig(*configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	if *addCommand {
		name := flag.Arg(0)
		if name == "" {
			fmt.Println("Error: command name is required")
			return
		}
		command := strings.Join(flag.Args()[1:], " ")
		if command == "" {
			fmt.Println("Error: command is required")
			return
		}

		commands, ok := config[*event]
		if !ok {
			commands = []Command{}
		}

		commands = append(commands, Command{Name: name, Command: command})
		config[*event] = commands

		err := writeConfig(*configFile, config)
		if err != nil {
			fmt.Printf("Error writing config file: %s\n", err)
			return
		}

		fmt.Printf("Added command '%s' for event '%s'\n", name, *event)
		return
	}

	commands, ok := config[*event]
	if !ok {
		fmt.Printf("No commands defined for event %s\n", *event)
		return
	}

	for _, cmd := range commands {
		err := runCommand(cmd.Command)
		if err != nil {
			fmt.Printf("Error running command %s: %s\n", cmd.Name, err)
			return
		}
		fmt.Printf("Command %s executed successfully\n", cmd.Name)
	}
}

func readConfig(filename string) (map[string][]Command, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	commands := map[string][]Command{
		"push": config.Push,
	}
	return commands, nil
}

func writeConfig(filename string, config map[string][]Command) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := Config{Push: config["push"]}
	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func runCommand(command string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Run()
}
