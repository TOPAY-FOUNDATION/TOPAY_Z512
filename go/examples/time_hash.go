package examples

import (
	"fmt"
	"time"

	"github.com/topay-foundation/topayz512/pkg/hash"
)

// RunTimeHash runs the time-based hash example
func RunTimeHash() {
	fmt.Println("TOPAY-Z512 Time-based Hash Example")

	// Create a hash using the current time
	timeHash := hash.NewWithTime()
	fmt.Printf("Time-based hash: %s\n", timeHash)

	// Create another hash after a short delay
	fmt.Println("\nWaiting for 1 second...")
	time.Sleep(1 * time.Second)
	timeHash2 := hash.NewWithTime()
	fmt.Printf("Time-based hash after 1 second: %s\n", timeHash2)

	// Demonstrate that the hashes are different
	fmt.Printf("\nHashes are different: %v\n", timeHash != timeHash2)

	// Using the convenience function
	fmt.Printf("\nTime-based hash bytes length: %d\n", len(hash.Sum512WithTime()))

	fmt.Println("\nThis example demonstrates how to create hashes using the current time as input.")
	fmt.Println("This is useful for generating random-like hashes when no specific input is available.")
	fmt.Println("Each hash will be different because it's based on the current timestamp.")
}
