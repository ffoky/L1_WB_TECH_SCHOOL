package telnet

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func Run(opts Options) error {
	address := net.JoinHostPort(opts.Host, opts.Port)
	conn, err := net.DialTimeout("tcp", address, opts.Timeout)
	if err != nil {
		return fmt.Errorf("dial tcp: %w", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}(conn)

	handleConn(conn, os.Stdin, os.Stdout)
	return nil
}

func handleConn(conn net.Conn, input io.Reader, output io.Writer) {
	var wg sync.WaitGroup
	wg.Go(func() {
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}(conn)
		if err := writeToConn(conn, input); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	})
	wg.Go(func() {
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}(conn)
		if err := readFromConn(conn, output); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	})
	wg.Wait()
}

func writeToConn(conn net.Conn, input io.Reader) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		b := scanner.Bytes()
		b = append(b, '\n')
		_, err := conn.Write(b)
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func readFromConn(conn net.Conn, output io.Writer) error {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		_, err := fmt.Fprintln(output, scanner.Text())
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}
