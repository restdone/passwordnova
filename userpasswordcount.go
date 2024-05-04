package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Parse command line flags
	threshold := flag.Int("t", 0, "Threshold for username occurrences")
	flag.Parse()

	// Open the file
	file, err := os.Open("password_trim.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Initialize a map to store username occurrences
	occurrences := make(map[string]int)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}
		username := parts[0]
		occurrences[username]++
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

	// Print usernames that occur more than the threshold
	for username, count := range occurrences {
		if count > *threshold {
			fmt.Printf("%s occurs %d times\n", username, count)
		}
	}
}
