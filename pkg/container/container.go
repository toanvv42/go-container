//go:build linux

package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Run starts a new containerized process.
// Phase 1: sets up namespaces, chroot, and /proc mount.
func Run(command string, args ...string) error {
	fmt.Printf("Running %s %v\n", command, args)

	cmd := exec.Command("/proc/self/exe", append([]string{"child", command}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
	}

	return cmd.Run()
}

// Child is called inside the new namespaces to set up the container environment.
func Child(command string, args ...string) error {
	fmt.Printf("Running container process %s %v\n", command, args)

	// TODO Phase 1.2: set hostname
	// TODO Phase 1.5: pivot_root into rootfs
	// TODO Phase 1.6: mount /proc

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
