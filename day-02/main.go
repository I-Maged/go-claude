package main

import "fmt"

func cToF(c float64) float64 {
	return c*9/5 + 32
}

func fToC(f float64) float64 {
	return f - 32*5/9
}

func main() {
	fmt.Printf("30 Clesius = %f Fahrenheit\n", cToF(30))
	fmt.Printf("68 Fahrenheit = %f Celsius\n", fToC(68))

	fmt.Println(North, East, South, West) // North East South West
	fmt.Printf("1 KB = %.0f bytes\n", float64(KB))
	fmt.Printf("1 MB = %.0f bytes\n", float64(MB))
	fmt.Printf("1 GB = %.0f bytes\n", float64(GB))
}
