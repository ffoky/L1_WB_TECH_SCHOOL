package builtin

import (
	"fmt"
	"os"
)

func ChangeDirectory(dir string) {
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintln(os.Stderr, "shell: cd:", err)
	}
}
