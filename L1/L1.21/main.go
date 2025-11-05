package main

import "fmt"

type exactInterface interface {
	DoingExactThing()
}

type someStruct struct{}

func (s *someStruct) DoingSomething() {
	fmt.Println("someStruct doing Something")
}

type Adapter struct {
	*someStruct
}

func (a *Adapter) DoingExactThing() {
	a.DoingSomething()
}

func client(e exactInterface) {
	e.DoingExactThing()
}

func main() {
	existing := &someStruct{}
	adapter := &Adapter{someStruct: existing}

	client(adapter)
}
