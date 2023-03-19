package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument without commas.
func WordCount(text string) map[string]int {
	// Create a map to hold the word frequencies
	freqs := make(map[string]int)
	text = strings.ToLower(text)

	// Split the text into words
	words := splitWords(text)

	// Count the word frequencies
	for _, word := range words {
		freqs[word]++
	}

	return freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {

	content, err := ioutil.ReadFile("\Users\timel\indapluswindows\twhite-palinda-3\src\singleworker\loremipsum.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := string(content)

	// Print the word frequencies
	fmt.Printf("%#v", WordCount(string(data)))
	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}

func splitWords(text string) []string {
	// Split the text into words
	words := strings.Fields(text)

	// Remove commas from the words
	for i, word := range words {
		words[i] = strings.Trim(word, ",")
		words[i] = strings.Trim(word, ".")
	}

	return words
}