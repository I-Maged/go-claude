package main

import "fmt"

type Greeter interface {
	Greet() string
}

type English struct {
	Name string
}

type Arabic struct {
	Name string
}

func (e English) Greet() string {
	return "Hello, " + e.Name
}
func (a Arabic) Greet() string {
	return "مرحبا،  " + a.Name
}

func PrintGreeting(g Greeter) {
	fmt.Println(g.Greet())
}
