package customerror

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field   string
	Value   any
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' value '%v' — %s", e.Field, e.Value, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %d not found", e.Resource, e.ID)
}

type HTTPError struct {
	Code    int
	Message string
	Err     error // wrapped underlying error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("HTTP %d: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

func (e *HTTPError) Unwrap() error { return e.Err }

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{
			Field:   "age",
			Value:   age,
			Message: "must be non-negative",
		}
	}
	if age > 150 {
		return &ValidationError{
			Field:   "age",
			Value:   age,
			Message: "unrealistically large",
		}
	}
	return nil
}

func findUser(id int) error {
	if id != 1 {
		underlying := &NotFoundError{Resource: "user", ID: id}
		return &HTTPError{Code: 404, Message: "user lookup failed", Err: underlying}
	}
	return nil
}

func TestValidationErrors() {
	fmt.Println("=== Custom Validation Errors ===")
	// errors.As — extracts the concrete type from the chain
	err := validateAge(-5)
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Printf("Field:   %s\n", ve.Field)
		fmt.Printf("Value:   %v\n", ve.Value)
		fmt.Printf("Message: %s\n", ve.Message)
	}

	fmt.Println()

	err = findUser(99)
	fmt.Println(err) // HTTP 404: user lookup failed: user with id 99 not found

	// errors.As walks the chain to find *NotFoundError
	var nfe *NotFoundError
	if errors.As(err, &nfe) {
		fmt.Printf("Resource: %s\n", nfe.Resource)
		fmt.Printf("ID:       %d\n", nfe.ID)
	}

	// errors.As finds *HTTPError directly
	var he *HTTPError
	if errors.As(err, &he) {
		fmt.Printf("HTTP Code: %d\n", he.Code)
	}
}
