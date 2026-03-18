package builtin

import (
	"fmt"
	"io"
	"strings"
)

func Echo(args []string, writer io.Writer) {
	_, _ = fmt.Fprintln(writer, strings.Join(args, " "))
}
