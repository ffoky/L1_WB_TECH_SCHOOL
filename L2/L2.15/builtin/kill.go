package builtin

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func Kill(args []string) {
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill:", err)
		return
	}

	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill:", err)
	}
}
