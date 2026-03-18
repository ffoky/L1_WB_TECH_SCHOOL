package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"L2.15/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer func() { _ = writer.Flush() }()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	for {
		currentWorkingDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		_, _ = fmt.Fprintf(writer, "%s> ", currentWorkingDir)
		_ = writer.Flush()

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				_, _ = fmt.Fprintln(writer)
				_, _ = fmt.Fprintln(writer, "logout")
				return
			}
			fmt.Fprintln(os.Stderr, err)
		}
		line = os.ExpandEnv(line)
		parser.Run(line, writer)
		_ = writer.Flush()
	}
}
