package calculator

import "fmt"

func main() {
	// Build shapes — using constructors with error handling
	builders := []func() (Shape, error){
		func() (Shape, error) { return NewCircle(5) },
		func() (Shape, error) { return NewRectangle(10, 4) },
		func() (Shape, error) { return NewTriangle(3, 4, 5) },
		func() (Shape, error) { return NewCircle(2.5) },
		func() (Shape, error) { return NewRectangle(6, 6) },
		func() (Shape, error) { return NewTriangle(5, 5, 5) },
	}

	shapes := make([]Shape, 0, len(builders))
	for _, build := range builders {
		s, err := build()
		if err != nil {
			fmt.Println("Skipping invalid shape:", err)
			continue
		}
		shapes = append(shapes, s)
	}

	fmt.Println("=== Shape Calculator ===")
	fmt.Println("All shapes:")
	printTable(shapes)

	fmt.Printf("\nLargest shape: %s (area=%.2f)\n",
		largest(shapes).String(), largest(shapes).Area())

	fmt.Println("\nSorted by area (smallest first):")
	sorted := sortByArea(shapes)
	for i, s := range sorted {
		fmt.Printf("  %d. %s — %.2f\n", i+1, s.String(), s.Area())
	}

	fmt.Println("\nAll shapes scaled by 2x:")
	scaled := scaleAll(shapes, 2)
	printTable(scaled)

	// Test error cases
	fmt.Println("\n--- Validation errors ---")
	if _, err := NewCircle(-5); err != nil {
		fmt.Println("Circle error:", err)
	}
	if _, err := NewRectangle(0, 10); err != nil {
		fmt.Println("Rectangle error:", err)
	}
	if _, err := NewTriangle(1, 2, 100); err != nil {
		fmt.Println("Triangle error:", err)
	}

	// Type assertion demo
	fmt.Println("\n--- Type assertions ---")
	for _, s := range shapes {
		switch v := s.(type) {
		case Circle:
			fmt.Printf("Circle with radius %.2f\n", v.Radius)
		case Rectangle:
			if v.Width == v.Height {
				fmt.Printf("Square with side %.2f\n", v.Width)
			} else {
				fmt.Printf("Rectangle %0.fx%.0f\n", v.Width, v.Height)
			}
		case Triangle:
			fmt.Printf("Triangle with perimeter %.2f\n", v.Perimeter())
		}
	}
}
