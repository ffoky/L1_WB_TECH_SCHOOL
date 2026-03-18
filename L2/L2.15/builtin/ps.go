package builtin

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ProcessStatus(writer io.Writer) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ps:", err)
		return
	}

	_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", "PID", "STATE", "NAME")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		data, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
		if err != nil {
			continue
		}

		name := ""
		state := ""
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "Name:") {
				name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			}
			if strings.HasPrefix(line, "State:") {
				state = strings.TrimSpace(strings.TrimPrefix(line, "State:"))
			}
		}

		_, _ = fmt.Fprintf(writer, "%d\t%s\t%s\n", pid, state, name)
	}
}
