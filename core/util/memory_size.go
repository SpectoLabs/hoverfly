package util

import (
	"fmt"
	"strconv"
	"strings"
)

// MemorySize is a custom type for parsing memory sizes (e.g., "128KB", "2MB")
type MemorySize int

// String returns the string representation of the memory size
func (m *MemorySize) String() string {
	return fmt.Sprintf("%d bytes", *m)
}

// ToBytes returns the memory size as an int in bytes.
func (m *MemorySize) ToBytes() int {
	return int(*m)
}

// Set parses a string like "128KB" or "2MB" and converts it to bytes
func (m *MemorySize) Set(value string) error {
	multiplier := 1

	value = strings.ToUpper(strings.TrimSpace(value))

	switch {
	case strings.HasSuffix(value, "KB"):
		multiplier = 1024
		value = strings.TrimSuffix(value, "KB")
	case strings.HasSuffix(value, "MB"):
		multiplier = 1024 * 1024
		value = strings.TrimSuffix(value, "MB")
	case strings.HasSuffix(value, "GB"):
		multiplier = 1024 * 1024 * 1024
		value = strings.TrimSuffix(value, "GB")
	}

	size, err := strconv.Atoi(value)

	if err != nil {
		return fmt.Errorf("invalid memory size: %s", value)
	}

	if size <= 0 {
		return fmt.Errorf("memory size must be greater than 0")
	}

	*m = MemorySize(size * multiplier)
	return nil
}

