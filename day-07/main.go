package main

import "fmt"

func main() {
	PrintGreeting(English{Name: "Ahmed"})
	PrintGreeting(Arabic{Name: "ماجد"})

	t := Temperature{Celsius: 35}
	fmt.Println(t)
}
