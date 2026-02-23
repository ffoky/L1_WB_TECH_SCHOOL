package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

const defaultAddress = "0.beevik-ntp.pool.ntp.org"

func currentTime(address string) (time.Time, error) {
	t, err := ntp.Time(address)
	if err != nil {
		return time.Time{}, fmt.Errorf("get time from server: %w", err)
	}
	return t, nil
}

func main() {
	time, err := currentTime(defaultAddress)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(time.Format("Mon, Jan 2 2006 15:04:05 MST"))
}
