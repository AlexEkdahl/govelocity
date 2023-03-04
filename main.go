package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlexEkdahl/govelocity/internal/cli"
	"github.com/AlexEkdahl/govelocity/internal/process"
	"github.com/AlexEkdahl/govelocity/internal/storage"
)

const dbFile = "velocity.db"

func main() {
	args := os.Args
	action := args[1]

	dbPath := filepath.Join(".", dbFile)
	s, err := storage.OpenDatabase(dbPath)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	pm := process.NewProcessManager(s)
	c := cli.NewCli()

	h, ok := c.Handler[action]
	if !ok {
		fmt.Println("comand does not exist: ", action)
		os.Exit(1)
	}

	if err := h.Validation(args, h.Length); err != nil {
		panic(err)
	}

	if err := h.Cmd(pm, args); err != nil {
		panic(err)
	}
}
