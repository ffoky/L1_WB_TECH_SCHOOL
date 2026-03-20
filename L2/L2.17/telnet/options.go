package telnet

import "time"

type Options struct {
	Host    string
	Port    string
	Timeout time.Duration
}
