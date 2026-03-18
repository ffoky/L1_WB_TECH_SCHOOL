package builtin

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"одно слово", []string{"hello"}, "hello\n"},
		{"несколько слов", []string{"hello", "world"}, "hello world\n"},
		{"пустые аргументы", []string{}, "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Echo(tt.args, &buf)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestPrintWorkingDirectory(t *testing.T) {
	var buf bytes.Buffer
	PrintWorkingDirectory(&buf)

	expected, _ := os.Getwd()
	assert.Equal(t, expected+"\n", buf.String())
}

func TestChangeDirectory(t *testing.T) {
	original, _ := os.Getwd()
	defer func() { _ = os.Chdir(original) }()

	ChangeDirectory("/tmp")
	current, _ := os.Getwd()
	assert.Equal(t, "/tmp", current)
}

func TestChangeDirectoryNotExists(t *testing.T) {
	original, _ := os.Getwd()
	defer func() { _ = os.Chdir(original) }()

	ChangeDirectory("/not/existing/path")
	current, _ := os.Getwd()
	assert.Equal(t, original, current)
}

func TestProcessStatus(t *testing.T) {
	var buf bytes.Buffer
	ProcessStatus(&buf)

	output := buf.String()
	assert.Contains(t, output, "PID")
	assert.Contains(t, output, "STATE")
	assert.Contains(t, output, "NAME")
}
