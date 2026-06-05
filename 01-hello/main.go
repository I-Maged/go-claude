package main

import (
	"fmt"
	"time"
)

func main() {
	myName := "Maged"
	today := time.Now().Format("Monday, January 2, 2006")
	fmt.Println("Hello, " + myName + "!")
	fmt.Printf("today is %s\n", today)
	fmt.Printf("The yeas is %d\n", time.Now().Year())
}
