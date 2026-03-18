package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"L2.15/parser"
)

func TestShellEcho(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("echo hello world", &buf)
	assert.Equal(t, "hello world\n", buf.String())
}

func TestShellPwd(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("pwd", &buf)

	expected, _ := os.Getwd()
	assert.Equal(t, expected+"\n", buf.String())
}

func TestShellCdAndPwd(t *testing.T) {
	original, _ := os.Getwd()
	defer func() { _ = os.Chdir(original) }()

	var buf bytes.Buffer
	parser.Run("cd /tmp", &buf)

	buf.Reset()
	parser.Run("pwd", &buf)
	assert.Equal(t, "/tmp\n", buf.String())
}

func TestShellExternalCommand(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("echo test_external", &buf)
	assert.Equal(t, "test_external\n", buf.String())
}

func TestShellPipeline(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("echo hello world | wc -w", &buf)
	assert.Equal(t, "2\n", strings.TrimSpace(buf.String())+"\n")
}

func TestShellRedirectOutput(t *testing.T) {
	tmpFile := "/tmp/shell_test_out.txt"
	defer func() { _ = os.Remove(tmpFile) }()

	var buf bytes.Buffer
	parser.Run("echo redirect_test > "+tmpFile, &buf)

	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	assert.Equal(t, "redirect_test\n", string(data))
}

func TestShellRedirectInput(t *testing.T) {
	tmpFile := "/tmp/shell_test_in.txt"
	err := os.WriteFile(tmpFile, []byte("line1\nline2\nline3\n"), 0o644)
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile) }()

	var buf bytes.Buffer
	parser.Run("wc -l < "+tmpFile, &buf)
	assert.Equal(t, "3\n", strings.TrimSpace(buf.String())+"\n")
}

func TestShellAndOperatorSuccess(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("echo first && echo second", &buf)

	output := buf.String()
	assert.Contains(t, output, "first")
	assert.Contains(t, output, "second")
}

func TestShellAndOperatorFail(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("__no_such_cmd__ && echo should_not_appear", &buf)
	assert.NotContains(t, buf.String(), "should_not_appear")
}

func TestShellOrOperatorFail(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("__no_such_cmd__ || echo fallback", &buf)
	assert.Contains(t, buf.String(), "fallback")
}

func TestShellOrOperatorSuccess(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("echo ok || echo should_not_appear", &buf)
	assert.NotContains(t, buf.String(), "should_not_appear")
}

func TestShellEmptyInput(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("", &buf)
	assert.Empty(t, buf.String())
}

func TestShellPs(t *testing.T) {
	var buf bytes.Buffer
	parser.Run("ps", &buf)

	output := buf.String()
	assert.Contains(t, output, "PID")
}
