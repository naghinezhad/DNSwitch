package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

// get user choice
func GetUserChoice(maxChoice int) int {
	for {
		var choice string
		fmt.Printf("Please enter your choice (1-%d): ", maxChoice)
		if _, err := fmt.Scanln(&choice); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		choice = strings.TrimSpace(choice)
		if choice == "" {
			fmt.Println("Input cannot be empty. Please enter a valid number.")
			continue
		}

		index, err := strconv.Atoi(choice)
		if err == nil && index >= 1 && index <= maxChoice {
			return index
		}
		fmt.Printf("Invalid input '%s'. Please enter a number between 1 and %d.\n", choice, maxChoice)
	}
}
