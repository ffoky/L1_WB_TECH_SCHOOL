package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

func typeDefinitor(v interface{}) string{
	switch reflect.TypeOf(v).Kind() {
	case reflect.Int:
		return "int"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "bool"
	case reflect.Chan:
		return "chan"
	}
	return ""
}

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var channel chan int
	var v int
	var s string
	var f bool
	fmt.Fprintln(writer, typeDefinitor(channel))
	fmt.Fprintln(writer, typeDefinitor(v))
	fmt.Fprintln(writer, typeDefinitor(s))
	fmt.Fprintln(writer, typeDefinitor(f))
}
