package builtin

import (
	"fmt"
	"io"
	"os"
)

func PrintWorkingDirectory(writer io.Writer) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	_, _ = fmt.Fprintln(writer, wd)
}
