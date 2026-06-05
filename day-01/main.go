package main

import (
	"fmt"
	"time"
)

func main() {
	myName := "Maged"
	today := time.Now().Format("Monday, January 2, 2006")
	now := time.Now()
	hour := now.Hour()
	greeting := getGreeting(hour)

	fmt.Printf("%s, %s!\n", greeting, myName)
	fmt.Printf("today is %s\n", today)
	fmt.Printf("The yeas is %d\n", time.Now().Year())
}
