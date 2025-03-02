package util

import (
	"fmt"
	"strconv"
	"strings"
)

// MemorySize is a custom type for parsing memory sizes (e.g., "128KB", "2MB")
type MemorySize int64

// String returns the string representation of the memory size
func (m *MemorySize) String() string {
	return fmt.Sprintf("%d bytes", *m)
}

// ToBytes returns the memory size as an int64 in bytes.
func (m *MemorySize) ToBytes() int64 {
	return int64(*m)
}

// Set parses a string like "128KB" or "2MB" and converts it to bytes
func (m *MemorySize) Set(value string) error {
	multiplier := int64(1)

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

	size, err := strconv.ParseInt(value, 10, 64)
	if err != nil || size < 0 {
		return fmt.Errorf("invalid memory size: %s", value)
	}

	*m = MemorySize(size * multiplier)
	return nil
}

