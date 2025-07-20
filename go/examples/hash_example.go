// Example demonstrating TOPAY-Z512 hashing functionality
//
// This example demonstrates how to use the TOPAY-Z512 hash functionality.
// It uses the current time as input to show dynamic hashing.

package examples

import (
	"fmt"
	"strconv"
	"time"

	"github.com/topay-foundation/topayz512/pkg/hash"
)

// Run demonstrates the hash functionality
func Run() {
	fmt.Println("TOPAY-Z512 Hash Example with Time-based Input")

	// Get current time as input
	now := time.Now()
	timeStr := fmt.Sprintf("TOPAY-Z512 time: %d seconds", now.Unix())
	data := []byte(timeStr)
	hashValue := hash.New(data)

	fmt.Printf("Input: %s\n", string(data))
	fmt.Printf("Hash: %s\n", hashValue)
	fmt.Printf("Hash size: %d bytes\n\n", len(hashValue))

	// Hash combination with time components
	timeMillis := strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)
	data1 := []byte(fmt.Sprintf("TOPAY-%s", timeMillis))
	data2 := []byte(fmt.Sprintf("Z512-%d", now.UnixNano()%1000000))

	combinedHash := hash.Combine(data1, data2)

	fmt.Printf("Input 1: %s\n", string(data1))
	fmt.Printf("Input 2: %s\n", string(data2))
	fmt.Printf("Combined Hash: %s\n", combinedHash)

	// Concatenated hash (different from combined hash)
	concatenated := append([]byte{}, data1...)
	concatenated = append(concatenated, data2...)
	concatHash := hash.New(concatenated)

	fmt.Printf("Concatenated Hash: %s\n", concatHash)
	fmt.Printf("Are they equal? %t\n\n", combinedHash == concatHash)

	// Convenience functions
	hashBytes := hash.Sum512(data)
	fmt.Printf("Hash bytes (first 8): %v\n", hashBytes[:8])

	combinedBytes := hash.SumCombine(data1, data2)
	fmt.Printf("Combined hash bytes (first 8): %v\n", combinedBytes[:8])

	// Hex conversion
	hex := hashValue.String()
	fmt.Printf("\nHex string: %s\n", hex)

	hashFromHex, err := hash.FromHex(hex)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Converted back from hex: %s\n", hashFromHex)
		fmt.Printf("Equal to original: %t\n", hashValue == hashFromHex)
	}
}
