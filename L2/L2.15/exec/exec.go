package exec

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func RunFromStdin(args []string, writer io.Writer, in io.Reader) bool {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = in
	cmd.Stdout = writer
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	return true
}

func RunPipeline(line string, writer io.Writer) bool {
	strs := strings.Split(line, "|")

	cmds := makeCommands(strs)

	connectIO(cmds)

	cmds[0].Stdin = os.Stdin
	cmds[len(cmds)-1].Stdout = writer

	for _, cmd := range cmds {
		cmd.Stderr = os.Stderr
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return false
		}
	}

	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return false
		}
	}
	return true
}

func makeCommands(strs []string) []*exec.Cmd {
	var cmds []*exec.Cmd
	for _, part := range strs {
		args := strings.Fields(strings.TrimSpace(part))
		if len(args) == 0 {
			continue
		}
		cmds = append(cmds, exec.Command(args[0], args[1:]...))
	}
	return cmds
}

func connectIO(cmds []*exec.Cmd) {
	for i := 0; i < len(cmds)-1; i++ {
		pipe, err := cmds[i].StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		cmds[i+1].Stdin = pipe
	}
}
