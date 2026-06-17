package calculator

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Shape interface {
	Area() float64
	Perimeter() float64
	String() string
}

type Scalable interface {
	Scale(factor float64) Shape
}

type Circle struct {
	Radius float64
}

func NewCircle(r float64) (Circle, error) {
	if r <= 0 {
		return Circle{}, fmt.Errorf("radius must be positive, got %.2f", r)
	}

	return Circle{Radius: r}, nil
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

func (c Circle) Scale(f float64) Shape {
	return Circle{Radius: c.Radius * f}
}

// Rectangle

type Rectangle struct {
	Height float64
	Width  float64
}

func NewRectangle(h, w float64) (Rectangle, error) {
	if h <= 0 || w <= 0 {
		return Rectangle{}, fmt.Errorf("dimensions must be positive, got %.2fx%.2f", w, h)
	}

	return Rectangle{Height: h, Width: w}, nil
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
func (r Rectangle) Scale(f float64) Shape {
	return Rectangle{Width: r.Width * f, Height: r.Height * f}
}

// Triangle
type Triangle struct {
	A, B, C float64
}

func NewTriangle(a, b, c float64) (Triangle, error) {
	if a <= 0 || b <= 0 || c <= 0 {
		return Triangle{}, fmt.Errorf("sides must be positive")
	}
	if a+b <= c || b+c <= a || a+c <= b {
		return Triangle{}, fmt.Errorf("invalid triangle: sides %.2f, %.2f, %.2f", a, b, c)
	}
	return Triangle{A: a, B: b, C: c}, nil
}

func (t Triangle) Area() float64 {
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}
func (t Triangle) String() string {
	return fmt.Sprintf("Triangle(%.2f, %.2f, %.2f)", t.A, t.B, t.C)
}
func (t Triangle) Scale(f float64) Shape {
	return Triangle{A: t.A * f, B: t.B * f, C: t.C * f}
}

// Calculator functions

func totalArea(shapes []Shape) float64 {
	var total float64

	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func totalPerimeter(shapes []Shape) float64 {
	total := 0.0

	for _, s := range shapes {
		total += s.Area()
	}

	return total
}

func largest(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}

	largest := shapes[0]

	for _, s := range shapes {
		if s.Area() > largest.Area() {
			largest = s
		}
	}

	return largest
}

func sortByArea(shapes []Shape) []Shape {
	sorted := make([]Shape, len(shapes))

	copy(sorted, shapes)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Area() < sorted[j].Area()
	})

	return sorted
}

func scaleAll(shapes []Shape, factor float64) []Shape {
	scaled := make([]Shape, len(shapes))

	for i, s := range shapes {
		if sc, ok := s.(Scalable); ok {
			scaled[i] = sc.Scale(factor)
		} else {
			scaled[i] = s
		}
	}

	return scaled
}

func printTable(shapes []Shape) {
	fmt.Printf("%-35s  %10s  %12s\n", "SHAPE", "AREA", "PERIMETER")
	fmt.Println(strings.Repeat("-", 62))

	for _, s := range shapes {
		fmt.Printf("%-35s  %10.2f  %12.2f\n", s.String(), s.Area(), s.Perimeter())
	}

	fmt.Println(strings.Repeat("-", 62))
	fmt.Printf("%-35s  %10.2f  %12.2f\n", "TOTAL", totalArea(shapes), totalPerimeter(shapes))
}
