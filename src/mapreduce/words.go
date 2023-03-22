package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument without commas or periods.
func WordCount(text string) map[string]int {
	// Create a map to hold the word frequencies
	freqs := make(map[string]int)
	ch := make(chan map[string]int)
	text = strings.ToLower(text)
	words := strings.Fields(text)
	size := 1000
	length := len(words)
	var wg sync.WaitGroup

	for i, j := 0, size; i < length; i, j = j, j+size {
		if j > length {
			j = length
		}
		wg.Add(1)
		go func(words []string) {
			subRoutineFreqs := make(map[string]int)
			for _, word := range words {
				word = strings.Trim(word, ",")
				word = strings.Trim(word, ".")
				subRoutineFreqs[word]++
			}
			ch <- subRoutineFreqs
			wg.Done()
		}(words[i:j])
	}
	wg.Wait()
	close(ch)

	for subRoutine := range ch {
		for word, value := range subRoutine {
			freqs[word] += value
		}
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
	// read in DataFile as a string called data
	data, err := ioutil.ReadFile(DataFile)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", WordCount(string(data)))

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
