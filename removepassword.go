package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"regexp"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./removepassword [username]")
		return
	}

	username := os.Args[1]

	file, err := os.OpenFile("password_trim.txt", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var found bool
	var lines []string
	var removeIndex int

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == username || parts[1] == username {
			found = true
			lines = append(lines, line)
		}
	}

	if !found {
		fmt.Println("No matching passwords found for username:", username)
		return
	}

	if len(lines) == 0 {
		fmt.Println("No passwords found for", username)
		return
	}

	fmt.Println("Passwords found for", username+":")
	for i, line := range lines {
		fmt.Printf("[%d] %s\n", i+1, line)
	}

	fmt.Print("Remove: ")
	if _, err := fmt.Scan(&removeIndex); err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	if removeIndex < 1 || removeIndex > len(lines) {
		fmt.Println("Invalid choice")
		return
	}

	// Print the chosen line
	fmt.Println("Chosen line:", lines[removeIndex-1])
	var  stringToRemove = lines[removeIndex-1]

	// Append \
	// Define a regular expression pattern to match special characters
	regex := regexp.MustCompile(`([[\]{}()*+?.,:\\^$|#\s])`)
	// Replace special characters with their escaped versions
	escapedStr := regex.ReplaceAllString(stringToRemove, `\$1`)

	fmt.Println("Original string:", stringToRemove)
	fmt.Println("Escaped string:", escapedStr)

	cmd := "sed"
	args := []string{"-i", "/"+escapedStr+"/d", "password_trim.txt"}

	out, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		fmt.Println("Error executing sed command:", err)
		return
	}


	fmt.Println(string(out))
	fmt.Println("Line containing", stringToRemove, "removed successfully.")
}
