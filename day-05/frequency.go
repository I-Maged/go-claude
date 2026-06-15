package main

import (
	"fmt"
	"sort"
	"strings"
)

func wordFrequency(text string) map[string]int {
	freq := make(map[string]int)
	words := strings.Fields(strings.ToLower(text))

	for _, word := range words {
		word = strings.Trim(word, ".,!?\"'()[]{}:;")
		if word == "" {
			continue
		}
		freq[word]++
	}
	return freq
}

func topN(freq map[string]int, n int) []string {
	// Collect all words into a slice so we can sort them
	type wordCount struct {
		word  string
		count int
	}

	pairs := make([]wordCount, 0, len(freq))
	for word, count := range freq {
		pairs = append(pairs, wordCount{word, count})
	}

	// Sort by count descending, then alphabetically for ties
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].word < pairs[j].word
	})

	// Extract just the words
	if n > len(pairs) {
		n = len(pairs)
	}
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = pairs[i].word
	}
	return result
}

// uniqueWords returns sorted list of unique words
func uniqueWords(freq map[string]int) []string {
	words := make([]string, 0, len(freq))
	for word := range freq {
		words = append(words, word)
	}
	sort.Strings(words)
	return words
}

func printFreqTable(freq map[string]int, words []string) {
	fmt.Printf("%-20s %s\n", "WORD", "COUNT")
	fmt.Println(strings.Repeat("-", 28))
	for _, word := range words {
		bar := strings.Repeat("█", freq[word])
		fmt.Printf("%-20s %d  %s\n", word, freq[word], bar)
	}
}
