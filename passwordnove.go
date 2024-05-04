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
	trimFlag := flag.Bool("t", false, "Flag to trim domain from email addresses")
	flag.Parse()

	// Check if the username file flag is provided
	if *usernameFile == "" {
		fmt.Println("Please provide a path to the file containing usernames using the -u flag")
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

	// Create a map to store unique passwords for each username
	passwordMap := make(map[string]map[string]bool)
	var passwordMutex sync.Mutex

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
					// Split the line by ":" to extract username and password
					parts := strings.Split(line, ":")
					if len(parts) == 2 {
						username := parts[0]
						password := parts[1]

						// Skip line if password is empty
						if password == "" {
							continue
						}

						// Check if username exists in the passwordMap
						passwordMutex.Lock()
						if _, ok := passwordMap[username]; !ok {
							passwordMap[username] = make(map[string]bool)
						}

						// Check if password is unique for the username
						if !passwordMap[username][password] {
							passwordMap[username][password] = true
							results <- line
						}
						passwordMutex.Unlock()
					}
				}
			}
		}(username)
	}

	// Close the results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Define output file name
	outputFile := "passwordnova_result.txt"

	// Open the output file for writing
	outFile, err := os.Create(outputFile)
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

	fmt.Println("Results have been exported to", outputFile)

	// Check if the -t flag is provided
	if *trimFlag {
		// Read passwordnova_result.txt and generate password_trim.txt with domain removed
		trimFile, err := os.Open(outputFile)
		if err != nil {
			fmt.Println("Error opening result file:", err)
			return
		}
		defer trimFile.Close()

		// Create password_trim.txt for writing
		trimOutFile, err := os.Create("password_trim.txt")
		if err != nil {
			fmt.Println("Error creating trim output file:", err)
			return
		}
		defer trimOutFile.Close()

		// Create a map to store unique passwords for each username
		trimmedPasswordMap := make(map[string]map[string]bool)
		var trimmedPasswordMutex sync.Mutex

		// Read lines from passwordnova_result.txt and write to password_trim.txt with domain removed
		scanner := bufio.NewScanner(trimFile)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				username := parts[0]
				password := parts[1]

				// Skip line if password is empty
				if password == "" {
					continue
				}

				// Check if username exists in the trimmedPasswordMap
				trimmedPasswordMutex.Lock()
				if _, ok := trimmedPasswordMap[username]; !ok {
					trimmedPasswordMap[username] = make(map[string]bool)
				}

				// Check if password is unique for the username
				if !trimmedPasswordMap[username][password] {
					trimmedPasswordMap[username][password] = true
					idx := strings.Index(username, "@")
					if idx != -1 {
						username = username[:idx]
					}
					trimmedLine := username + ":" + password
					fmt.Fprintln(trimOutFile, trimmedLine)
				}
				trimmedPasswordMutex.Unlock()
			}
		}

		fmt.Println("Trimmed results have been exported to password_trim.txt")
	}
}
