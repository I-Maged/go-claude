package validate

import (
	"fmt"
	"strings"
)

func NotEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", field)
	}
	return nil
}

func Positive(field string, value float64) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive, got %.2f", field, value)
	}
	return nil
}

func NonNegative(field string, value float64) error {
	if value < 0 {
		return fmt.Errorf("%s cannot be negative, got %.2f", field, value)
	}
	return nil
}
