package parser

import (
	"fmt"
	"io"
	"os"
	"strings"

	"L2.15/builtin"
	"L2.15/exec"
)

type ParsedCommand struct {
	cmd     string
	inFile  string
	outFile string
}

func Run(s string, writer io.Writer) {
	parts := splitCmdsByConds(s)

	lastOk := true
	skip := false
	for _, p := range parts {
		switch p {
		case "&&":
			skip = !lastOk
		case "||":
			skip = lastOk
		default:
			if !skip {
				lastOk = parse(p, writer)
			}
			skip = false
		}
	}
}

func parse(s string, writer io.Writer) bool {
	if strings.Contains(s, "|") {
		return exec.RunPipeline(s, writer)
	}

	parsed := parseRedirects(s)

	reader, closeIn := openInFile(parsed.inFile)
	defer closeIn()
	if reader == nil {
		return false
	}

	writer, closeOut := createOutFile(parsed.outFile, writer)
	defer closeOut()
	if writer == nil {
		return false
	}

	args := strings.Fields(parsed.cmd)
	if len(args) == 0 {
		return true
	}

	return executeCommand(args, writer, reader)
}

func openInFile(path string) (io.Reader, func()) {
	if path == "" {
		return os.Stdin, func() {}
	}
	openFile, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, func() {}
	}
	return openFile, func() { _ = openFile.Close() }
}

func createOutFile(path string, fallback io.Writer) (io.Writer, func()) {
	if path == "" {
		return fallback, func() {}
	}
	outFile, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, func() {}
	}
	return outFile, func() { _ = outFile.Close() }
}

func executeCommand(args []string, writer io.Writer, reader io.Reader) bool {
	switch args[0] {
	case "echo":
		builtin.Echo(args[1:], writer)
	case "pwd":
		builtin.PrintWorkingDirectory(writer)
	case "cd":
		if len(args) > 1 {
			builtin.ChangeDirectory(args[1])
		}
	case "kill":
		builtin.Kill(args[1:])
	case "ps":
		builtin.ProcessStatus(writer)
	default:
		return exec.RunFromStdin(args, writer, reader)
	}
	return true
}

func cutRedirect(s string, symbol string) (before, after string) {
	idx := strings.Index(s, symbol)
	if idx < 0 {
		return s, ""
	}
	return strings.TrimSpace(s[:idx]), strings.TrimSpace(s[idx+1:])
}

func parseRedirects(s string) *ParsedCommand {
	res := &ParsedCommand{cmd: s}
	res.cmd, res.outFile = cutRedirect(res.cmd, ">")
	res.cmd, res.inFile = cutRedirect(res.cmd, "<")
	return res
}

func splitCmdsByConds(s string) []string {
	s = strings.ReplaceAll(s, "&&", "\n&&\n")
	s = strings.ReplaceAll(s, "||", "\n||\n")

	var parts []string
	for _, p := range strings.Split(s, "\n") {
		p = strings.TrimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}
