package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
	String() string
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width, Height float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}
func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.2f)", c.Radius)
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%.2f x %.2f)", r.Width, r.Height)
}

func printShapeInfo(s Shape) {
	fmt.Printf("%-30s  area=%-10.2f  perimeter=%.2f\n",
		s.String(), s.Area(), s.Perimeter())
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func describe(s Shape) {
	// Type switch — like a regular switch but on the type
	switch v := s.(type) {
	case Circle:
		fmt.Printf("Circle with radius %.2f, area %.2f\n", v.Radius, v.Area())
	case Rectangle:
		fmt.Printf("Rectangle %0.fx%.0f, area %.2f\n", v.Width, v.Height, v.Area())
	default:
		fmt.Printf("Unknown shape: %T\n", v)
	}
}
