package telnet

import (
	"bytes"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEcho(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()

	go func() {
		defer server.Close()
		buf := make([]byte, 1024)
		for {
			n, err := server.Read(buf)
			if err != nil {
				return
			}
			_, _ = server.Write(buf[:n])
		}
	}()

	pr, pw := io.Pipe()
	output := &bytes.Buffer{}

	go func() {
		_, _ = pw.Write([]byte("hello\nworld\n"))
		time.Sleep(100 * time.Millisecond)
		pw.Close()
	}()

	handleConn(client, pr, output)

	assert.Equal(t, "hello\nworld", strings.TrimSpace(output.String()))
}

func TestServerClose(t *testing.T) {
	server, client := net.Pipe()

	go func() {
		_, _ = server.Write([]byte("bye\n"))
		server.Close()
	}()

	output := &bytes.Buffer{}
	handleConn(client, &bytes.Buffer{}, output)

	assert.Equal(t, "bye", strings.TrimSpace(output.String()))
}

func TestDialError(t *testing.T) {
	err := Run(Options{
		Host:    "invalid.host.that.does.not.exist",
		Port:    "9999",
		Timeout: 1 * time.Second,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dial tcp")
}
