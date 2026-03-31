package cmd

import (
	"fmt"
	"os"

	"github.com/toanvv/go-container/pkg/container"
)

func runCommand() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: go-container run <command> [args]")
	}

	return container.Run(os.Args[2], os.Args[3:]...)
}

// childCommand is invoked inside the new namespaces (internal, not user-facing).
func childCommand() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("child requires a command")
	}

	return container.Child(os.Args[2], os.Args[3:]...)
}
