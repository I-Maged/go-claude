package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrEmptyFile    = errors.New("file is empty")
	ErrNotFound     = errors.New("file not found")
	ErrNoPermission = errors.New("permission denied")
)

type FileError struct {
	Path      string
	Operation string
	Err       error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("file error [%s] on '%s': %v", e.Operation, e.Path, e.Err)
}
func (e *FileError) Unwrap() error { return e.Err }

type FileStats struct {
	Path      string
	Lines     int
	Words     int
	Chars     int
	LinesList []string
}

func (s FileStats) String() string {
	return fmt.Sprintf("'%s': %d lines, %d words, %d chars",
		filepath.Base(s.Path), s.Lines, s.Words, s.Chars)
}

func readFileStats(path string) (*FileStats, error) {
	// Check existence
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &FileError{
				Path:      path,
				Operation: "stat",
				Err:       ErrNotFound,
			}
		}
		if os.IsPermission(err) {
			return nil, &FileError{
				Path:      path,
				Operation: "stat",
				Err:       ErrNoPermission,
			}
		}
		return nil, fmt.Errorf("readFileStats: %w", err)
	}

	// Check not a directory
	if info.IsDir() {
		return nil, &FileError{
			Path:      path,
			Operation: "read",
			Err:       fmt.Errorf("path is a directory, not a file"),
		}
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, &FileError{Path: path, Operation: "open", Err: err}
	}
	defer file.Close()

	// Read line by line
	stats := &FileStats{Path: path}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		stats.Chars += len(line) + 1 // +1 for newline
		stats.Words += len(strings.Fields(line))
		stats.LinesList = append(stats.LinesList, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, &FileError{Path: path, Operation: "scan", Err: err}
	}

	if stats.Lines == 0 {
		return nil, &FileError{
			Path:      path,
			Operation: "read",
			Err:       ErrEmptyFile,
		}
	}

	return stats, nil
}

func processFiles(paths []string) ([]*FileStats, []error) {
	results := make([]*FileStats, 0, len(paths))
	errs := make([]error, 0)

	for _, path := range paths {
		stats, err := readFileStats(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results = append(results, stats)
	}

	return results, errs
}

// classifyError returns a human-friendly category for an error
func classifyError(err error) string {
	var fe *FileError
	if errors.As(err, &fe) {
		switch {
		case errors.Is(err, ErrNotFound):
			return "NOT FOUND"
		case errors.Is(err, ErrNoPermission):
			return "PERMISSION DENIED"
		case errors.Is(err, ErrEmptyFile):
			return "EMPTY FILE"
		default:
			return "FILE ERROR"
		}
	}
	return "UNKNOWN ERROR"
}

func createTestFiles() {
	os.WriteFile("test1.txt",
		[]byte("Hello Go\nThis is line two\nAnd line three\n"), 0644)
	os.WriteFile("test2.txt",
		[]byte("Go error handling\nis explicit\nand powerful\nno exceptions needed\n"), 0644)
	os.WriteFile("empty.txt", []byte(""), 0644)
}

func cleanupTestFiles() {
	for _, f := range []string{"test1.txt", "test2.txt", "empty.txt"} {
		os.Remove(f)
	}
}
