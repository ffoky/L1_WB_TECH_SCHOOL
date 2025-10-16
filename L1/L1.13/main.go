package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	a:=1
	b:=2
	fmt.Fprintf(writer, "a = %d\nb = %d\n",a,b)
	a=a ^ b
	b = a ^ b
	a=a ^ b
	fmt.Fprintf(writer, "a = %d\nb = %d\n",a,b)
}
