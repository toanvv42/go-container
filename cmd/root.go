package cmd

import (
	"fmt"
	"os"
)

func Execute() error {
	if len(os.Args) < 2 {
		usage()
		return nil
	}

	switch os.Args[1] {
	case "run":
		return runCommand()
	case "child":
		return childCommand()
	default:
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func usage() {
	fmt.Println(`Usage: go-container <command> [args]

Commands:
  run    Run a command in a new container`)
}
