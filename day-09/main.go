package main

import (
	customerror "day-09/customError"
	"day-09/login"
	"errors"
	"fmt"
)

func main() {
	login.NewUserStore().TestLogin()
	customerror.TestValidationErrors()

	fmt.Println("=== File Reader ===")
	createTestFiles()
	defer cleanupTestFiles()

	files := []string{
		"test1.txt",
		"test2.txt",
		"empty.txt",
		"missing.txt", // does not exist
		".",           // a directory
	}

	fmt.Println("=== Processing files ===")
	results, errs := processFiles(files)

	// Print successful results
	if len(results) > 0 {
		fmt.Println("Successful reads:")
		for _, s := range results {
			fmt.Printf("  ✓ %s\n", s)
		}
	}

	// Print errors with classification
	if len(errs) > 0 {
		fmt.Printf("\nErrors (%d):\n", len(errs))
		for _, err := range errs {
			fmt.Printf("  ✗ [%s] %v\n", classifyError(err), err)
		}
	}

	// Demonstrate errors.Is through wrapping
	fmt.Println("\n=== errors.Is through the chain ===")
	_, err := readFileStats("missing.txt")
	fmt.Println("errors.Is(err, ErrNotFound):", errors.Is(err, ErrNotFound))

	var fe *FileError
	if errors.As(err, &fe) {
		fmt.Printf("FileError.Path:      %s\n", fe.Path)
		fmt.Printf("FileError.Operation: %s\n", fe.Operation)
		fmt.Printf("Underlying error:    %v\n", fe.Err)
	}

	// Demonstrate empty file error
	fmt.Println("\n=== Empty file detection ===")
	_, err = readFileStats("empty.txt")
	fmt.Println("errors.Is(err, ErrEmptyFile):", errors.Is(err, ErrEmptyFile))
	fmt.Println("Error message:", err)

	// Show first file's content preview
	if len(results) > 0 {
		fmt.Printf("\n=== Preview: %s ===\n", results[0].Path)
		for i, line := range results[0].LinesList {
			fmt.Printf("  %2d: %s\n", i+1, line)
		}
	}
}
