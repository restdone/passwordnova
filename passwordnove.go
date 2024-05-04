package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Define a struct to represent the JSON response
type Response struct {
	Count int      `json:"count"`
	Lines []string `json:"lines"`
}

func main() {
	// Define command-line flags
	usernameFile := flag.String("u", "", "Path to the file containing usernames")
	filterString := flag.String("d", "", "String to filter the lines")
	outputFile := flag.String("o", "", "Path to the output file")
	flag.Parse()

	// Check if the username file and output file flags are provided
	if *usernameFile == "" || *outputFile == "" {
		fmt.Println("Please provide a path to the file containing usernames using the -u flag and specify the output file using the -o flag")
		fmt.Println("Use -d to filter what to search in result")
		return
	}

	// Open the file containing usernames
	file, err := os.Open(*usernameFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read usernames from the file into a slice
	var usernames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		usernames = append(usernames, scanner.Text())
	}

	// Define the base URL without the parameter
	baseURL := "https://api.proxynova.com/comb?query="

	// Create a WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Create a channel to collect results from goroutines
	results := make(chan string)

	// Spawn multiple goroutines to make requests concurrently
	for _, username := range usernames {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()

			// Construct the URL with the parameter value
			url := baseURL + username

			// Make the HTTP request
			response, err := http.Get(url)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer response.Body.Close()

			// Check if the response is successful
			if response.StatusCode == http.StatusOK {
				// Decode JSON response
				var resp Response
				err := json.NewDecoder(response.Body).Decode(&resp)
				if err != nil {
					fmt.Println("Error decoding JSON:", err)
					return
				}

				// Filter results based on the filter string if provided
				if *filterString != "" {
					filteredLines := []string{}
					for _, line := range resp.Lines {
						if strings.Contains(line, *filterString) {
							filteredLines = append(filteredLines, line)
						}
					}
					resp.Lines = filteredLines
				}

				// Process and send response to results channel
				for _, line := range resp.Lines {
					results <- line
				}
			}
		}(username)
	}

	// Close the results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Open the output file for writing
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	// Write results to the output file
	writer := bufio.NewWriter(outFile)
	for line := range results {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}
	writer.Flush()

	fmt.Println("Results have been exported to", *outputFile)
}
