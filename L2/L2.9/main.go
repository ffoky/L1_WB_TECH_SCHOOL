package main

import (
	"L2.9/unpack"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	s := "45"
	u, err := unpack.Unpack(s)
	if err != nil {
		slog.Error("error", "err", err)
		os.Exit(1)
	}
	fmt.Println(u)
}
