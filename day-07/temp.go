package main

import "fmt"

type Stringer interface {
	String() string
}

type Temperature struct {
	Celsius float64
}

func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°C (%.1f°F)", t.Celsius, t.Celsius*9/5+32)
}
