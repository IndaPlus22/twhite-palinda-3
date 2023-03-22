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
	n := 8 // number of goroutines
	chunk := len(words) / n // split words into n chunks (arbitrary number
	length := len(words)
	var wg sync.WaitGroup

	// split words into chunks and run each chunk in a goroutine
	for i, j := 0, chunk; i < length; i, j = j, j+chunk {
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

	// close channel when all goroutines are done
	go func ()  {
		wg.Wait()
		close(ch)
	}()

	// read from channel
	for subRoutines := range ch {
		for word, value := range subRoutines {
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
