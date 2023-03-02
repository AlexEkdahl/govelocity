package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Push []Command `yaml:"push"`
}

type Command struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

func main() {
	configFile := "config.yml"
	event := "push"
	config, err := readConfig(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}
	commands, ok := config[event]
	if !ok {
		fmt.Printf("No commands defined for event %s\n", event)
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

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = parseYAML(data, &config)
	if err != nil {
		return nil, err
	}

	commands := map[string][]Command{
		"push": config.Push,
	}
	return commands, nil
}

func parseYAML(data []byte, v interface{}) error {
	lines := strings.Split(string(data), "\n")
	return parseYAMLLines(lines, v, 0)
}

func parseYAMLLines(lines []string, v interface{}, indent int) error {
	m := make(map[string]interface{})
	currentIndent := -1
	currentKey := ""
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		leadingSpaces := 0
		for _, c := range line {
			if c == ' ' {
				leadingSpaces++
			} else {
				break
			}
		}
		if leadingSpaces%2 != 0 {
			return fmt.Errorf("invalid indentation on line: %s", line)
		}
		indentation := leadingSpaces / 2
		if currentIndent == -1 {
			currentIndent = indentation
		} else if indentation == currentIndent {
			currentKey = strings.TrimSpace(line)
		} else if indentation > currentIndent {
			err := parseYAMLLines(lines, m, currentIndent+1)
			if err != nil {
				return err
			}
			currentIndent = indentation
			currentKey = strings.TrimSpace(line)
		} else if indentation < currentIndent {
			m[currentKey] = v
			return nil
		}
	}
	if currentIndent != -1 {
		m[currentKey] = v
	}
	return nil
}

func runCommand(command string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Run()
}
