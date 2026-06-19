package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type DownloadResult struct {
	URL      string
	Bytes    int
	Duration time.Duration
	Err      error
}

func (r *DownloadResult) String() string {
	if r.Err != nil {
		return fmt.Sprintf("✗ %-30s FAILED: %v", r.URL, r.Err)
	}
	return fmt.Sprintf("✓ %-30s %6d bytes  in %v", r.URL, r.Bytes, r.Duration.Round(time.Millisecond))
}

func simulateDownload(url string) DownloadResult {
	start := time.Now()

	// Random latency between 100-500ms
	latency := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(latency)

	// Simulate occasional failure (1 in 6 chance)
	if rand.Intn(6) == 0 {
		return DownloadResult{
			URL:      url,
			Duration: time.Since(start),
			Err:      fmt.Errorf("connection timeout"),
		}
	}

	// Simulate file size
	bytes := 1000 + rand.Intn(50000)

	return DownloadResult{
		URL:      url,
		Bytes:    bytes,
		Duration: time.Since(start),
	}
}

func downloadAll(urls []string) []DownloadResult {
	result := make([]DownloadResult, len(urls))

	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go func(index int, u string) {
			defer wg.Done()
			result[index] = simulateDownload(u)
		}(i, url)
	}

	wg.Wait()
	return result
}

func summarize(result []DownloadResult, totalWallTime time.Duration) {
	var success, failed, totalBytes int
	var totalIndvidualTime time.Duration

	for _, r := range result {
		if r.Err != nil {
			failed++
		} else {
			success++
			totalBytes += r.Bytes
		}
		totalIndvidualTime += r.Duration
	}

	fmt.Println()
	fmt.Println("=== Summary ===")
	fmt.Printf("Total downloads:      %d\n", len(result))
	fmt.Printf("Succeeded:            %d\n", success)
	fmt.Printf("Failed:               %d\n", failed)
	fmt.Printf("Total bytes:          %d\n", totalBytes)
	fmt.Printf("Wall clock time:      %v\n", totalWallTime.Round(time.Millisecond))
	fmt.Printf("Sum of individual times: %v (if run sequentially)\n", totalIndvidualTime.Round(time.Millisecond))
	speedup := float64(totalIndvidualTime) / float64(totalWallTime)
	fmt.Printf("Speedup from concurrency: %.1fx\n", speedup)
}

func TestFileDownloader() {
	rand.Seed(time.Now().UnixNano())

	urls := []string{
		"cdn.example.com/file1.zip",
		"cdn.example.com/file2.zip",
		"cdn.example.com/file3.zip",
		"cdn.example.com/file4.zip",
		"cdn.example.com/file5.zip",
		"cdn.example.com/file6.zip",
		"cdn.example.com/file7.zip",
		"cdn.example.com/file8.zip",
	}

	fmt.Println("=== Concurrent File Downloader ===")
	fmt.Printf("Downloading %d files concurrently...\n\n", len(urls))

	start := time.Now()
	results := downloadAll(urls)
	elapsed := time.Since(start)

	// Print results in original order
	for _, r := range results {
		fmt.Println(r)
	}

	summarize(results, elapsed)

}
