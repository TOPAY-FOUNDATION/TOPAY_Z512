// TOPAY-Z512 Time-based Hash Example
//
// This example demonstrates how to generate hashes using the current time
// with the TOPAY-Z512 library.
package examples

import (
	"fmt"
	"time"

	"github.com/topay-foundation/topayz512/pkg/hash"
)

// RunTimeHashExample demonstrates time-based hash generation
func RunTimeHashExample() {
	fmt.Println("TOPAY-Z512 Time-based Hash Example")

	// Create a hash using the current time
	timeHash := hash.NewWithTime()
	fmt.Printf("Time-based hash: %s\n", timeHash)

	// Create another hash after a short delay
	time.Sleep(1 * time.Second)
	timeHash2 := hash.NewWithTime()
	fmt.Printf("Time-based hash after 1 second: %s\n", timeHash2)

	// Demonstrate that the hashes are different
	fmt.Printf("Hashes are different: %v\n", timeHash != timeHash2)

	// Using the convenience function
	timeHashBytes := hash.Sum512WithTime()
	fmt.Printf("Time-based hash bytes: %v\n", timeHashBytes)
	fmt.Printf("Time-based hash bytes length: %d\n", len(timeHashBytes))
}