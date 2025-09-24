package main

import (
	"fmt"
)

type hours int
type meters int

type Human struct {
	IsStarving bool
	Rested     bool
	Healthy    bool
	Energy     int
}

func (h *Human) Eat() {
	h.IsStarving = false
	h.Energy += 40
}

func (h *Human) Sleep(timeToSleep hours) error {
	if !h.Rested {
		h.Energy += 50
		h.Rested = true
		fmt.Printf("Creature of type %T rested for %v hours\n", h, timeToSleep)
		return nil
	}
	return fmt.Errorf("already rested, no need to sleep")
}

func (h *Human) Run(distance meters) error {
	energyToRun := 5 * int(distance)
	if h.Energy-energyToRun >= 0 {
		h.Energy -= energyToRun
		h.Rested = false
		fmt.Printf("Creature of type %T ran for %v meters\n", h, distance)
		return nil
	}
	return fmt.Errorf("not enough energy to run")
}

type Dog struct {
	woof bool
}

func (d *Dog) Run() {
	if d.woof {
		fmt.Printf("Dog is running and barking")
	} else {
		fmt.Printf("Dog is running silently")
	}
}

//Пример Shadowing
//func (a *Action) Run(distance meters) {
//  fmt.Printf("Some creature run fow %v meters\n", distance)
//}

type Action struct {
	Human
	//Dog
	InProcess  bool
	IsFinished bool
}

func (a *Action) Fly(distance meters) {
	fmt.Printf("Creature of type %T flew for %v meters\n", a, distance)
}

func main() {
	h := &Human{Energy: 20, Rested: false}
	a := &Action{Human: *h}

	a.Eat()
	fmt.Printf("Energy of creature %T after eat: %v \n", a, a.Energy)

	if err := a.Run(3); err != nil {
		fmt.Println("Run error:", err)
	}
	//Shadowing
	//если я напишу a.Run и раскомментирую метод Action Run()
	//то, чтобы вызвать код Human нужно использовать такой синтаксис
	//a.Human.Run(5)
	//Colliding
	//Так же если существует embedded структура Dog, то нужно так же написать
	// отдельно a.Human.Run(5)
	// и отдельно a.Dog.Run(5)

	if err := a.Sleep(5); err != nil {
		fmt.Println("Sleep error:", err)
	}

	a.Fly(10)
}
