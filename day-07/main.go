package main

import "fmt"

func main() {
	PrintGreeting(English{Name: "Ahmed"})
	PrintGreeting(Arabic{Name: "ماجد"})

	t := Temperature{Celsius: 35}
	fmt.Println(t)

	shapes := []Shape{
		Circle{Radius: 2.5},
		Circle{Radius: 5},
		Rectangle{Width: 5, Height: 4},
		Rectangle{Width: 6, Height: 6},
	}

	for _, s := range shapes {
		printShapeInfo(s)
	}
	fmt.Printf("\nTotal area of all shapes: %.2f\n", totalArea(shapes))
	for _, s := range shapes {
		describe(s)
	}
}
