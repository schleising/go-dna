package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Sturcture to contain the name and counts for a file
type histogramContainer struct {
	name      string
	histogram map[rune]int
}

func readFile(filename string, callNum int, ch chan histogramContainer) {
	// Open the file
	file, err := os.Open(filename)

	// Check for no error
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Failed to open file")
	}

	// Ensure that defer is called when the function exist,
	// Using defer ensures the function is called even in the event of an error
	defer file.Close()

	// Create a new Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a new histogramConatiner for this file
	var histogram histogramContainer
	histogram.name = filename
	histogram.histogram = make(map[rune]int)

	// Read the file line by line using Scan()
	for scanner.Scan() {
		// Discard the line if it starts with a >
		if !strings.HasPrefix(scanner.Text(), ">") {
			// Loop through each line character (as a rune), use the rune as a key for the map
			for _, c := range scanner.Text() {
				histogram.histogram[c]++
			}
		}
	}

	// Print the totals for this file
	fmt.Println(histogram.name)
	for base, count := range histogram.histogram {
		fmt.Printf("Base %s\tCount %v\n", string(base), count)
	}

	// Send the histogramContainer back to the initiating function
	// This is safe even though histogram is defined in this function
	// as the garbage collector will recognise that it still has an active reference
	ch <- histogram
}

func main() {
	// Get the list of filenames matching *.fa in data
	matches, err := filepath.Glob("data/*.fa")

	// Check there was no error
	if err != nil {
		fmt.Println(err)
	} else {
		// Make a channel to send the results back from the goroutines
		// The channel type is histogramContainer as this is the type it will send/receive
		ch := make(chan histogramContainer)

		// Loop through each file printing the filename and spawning a
		// gorouting do read and parse the files
		for i, match := range matches {
			// The go keyword makes this run in a new goroutine,
			// the ch variable is the channel used to send back the results
			go readFile(match, i, ch)
		}

		// Create a new histogramContainer for the totals
		var totals histogramContainer
		totals.name = "Totals"

		// Need to use make here to initialise the map, otherwise adding data will cause a panic
		totals.histogram = make(map[rune]int)

		// Loop through waiting for results from the goroutines
		for i := 0; i < len(matches); i++ {
			// This is a blocking wait, waiting for data to be sent on channel ch
			fileCounts := <-ch

			// Add the latest results to the totals container
			for base, count := range fileCounts.histogram {
				totals.histogram[base] += count
			}
		}

		// Print the totals
		fmt.Println(totals.name)
		for base, count := range totals.histogram {
			fmt.Printf("Base %s\tCount %v\n", string(base), count)
		}
	}
}
