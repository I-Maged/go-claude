package main

// Simple constant
const Pi = 3.14159
const AppName = "MyApp"

// iota — starts at 0, increments by 1 each line
type Direction int

const (
	North Direction = iota // 0
	East                   // 1
	South                  // 2
	West                   // 3
)

// iota can do math
type ByteSize float64

const (
	KB ByteSize = 1024 << (10 * iota) // iota=0 → 1024 << 0  = 1024
	MB                                // iota=1 → 1024 << 10 = 1,048,576
	GB                                // iota=2 → 1024 << 20
)

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

// func tryIota() {
// 	fmt.Println(North, East, South, West) // North East South West
// 	fmt.Printf("1 KB = %.0f bytes\n", float64(KB))
// 	fmt.Printf("1 MB = %.0f bytes\n", float64(MB))
// 	fmt.Printf("1 GB = %.0f bytes\n", float64(GB))
// }
