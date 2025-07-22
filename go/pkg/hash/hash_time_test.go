package hash

import (
	"testing"
	"time"
)

func TestNewWithTime(t *testing.T) {
	// Generate two hashes with a small time gap
	hash1 := NewWithTime()
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps
	hash2 := NewWithTime()

	// The hashes should be different due to different timestamps
	if hash1 == hash2 {
		t.Errorf("Expected different hashes for different timestamps")
	}

	// Verify the hash size is correct
	if len(hash1) != HashSizeBytes {
		t.Errorf("Expected hash size %d, got %d", HashSizeBytes, len(hash1))
	}
}

func TestSum512WithTime(t *testing.T) {
	// Generate two hashes with a small time gap
	hash1 := Sum512WithTime()
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps
	hash2 := Sum512WithTime()

	// The hashes should be different due to different timestamps
	if hash1 == hash2 {
		t.Errorf("Expected different hashes for different timestamps")
	}

	// Verify the hash size is correct
	if len(hash1) != HashSizeBytes {
		t.Errorf("Expected hash size %d, got %d", HashSizeBytes, len(hash1))
	}
}