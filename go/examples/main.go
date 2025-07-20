// TOPAY-Z512 Examples Runner
//
// This file serves as the entry point to run all the examples in the TOPAY-Z512 Go implementation.
// It provides a simple menu to select which example to run.

package examples

import (
	"fmt"
	"os"
	"strconv"
)

// RunMenu displays a menu to select which example to run
func RunMenu() {
	if len(os.Args) > 1 {
		// If an argument is provided, run the specified example
		exampleNum, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Error: Invalid example number '%s'\n", os.Args[1])
			printUsage()
			os.Exit(1)
		}

		runExample(exampleNum)
		return
	}

	// If no argument is provided, show the menu
	showMenu()
}

func showMenu() {
	fmt.Println("TOPAY-Z512 Examples")
	fmt.Println("=================")
	fmt.Println("Please select an example to run:")
	fmt.Println("1. Hash Example")
	fmt.Println("2. Key Pair Example")
	fmt.Println("3. Private to Public Key Conversion Example")
	fmt.Println("0. Exit")
	fmt.Print("\nEnter your choice: ")

	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil {
		fmt.Println("\nInvalid input. Please enter a number.")
		return
	}

	fmt.Println() // Add a newline for better formatting
	runExample(choice)
}

func runExample(choice int) {
	switch choice {
	case 0:
		fmt.Println("Exiting...")
		os.Exit(0)
	case 1:
		Run()
	case 2:
		RunKeypair()
	case 3:
		RunPrivateToPublic()
	default:
		fmt.Printf("Invalid choice: %d\n", choice)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("\nUsage:")
	fmt.Println("  go run main.go [example_number]")
	fmt.Println("\nExample numbers:")
	fmt.Println("  1 - Hash Example")
	fmt.Println("  2 - Key Pair Example")
	fmt.Println("  3 - Private to Public Key Conversion Example")
}