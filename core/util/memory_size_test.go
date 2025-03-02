package util

import "testing"
import . "github.com/onsi/gomega"

func TestMemorySize_Set(t *testing.T) {
	RegisterTestingT(t) // Register Gomega for the test

	t.Run("valid inputs", func(t *testing.T) {
		var ms MemorySize

		// Test inputs with valid values
		err := ms.Set("128KB")
		Expect(err).To(BeNil())
		Expect(ms).To(Equal(MemorySize(128 * 1024)))

		err = ms.Set("2MB")
		Expect(err).To(BeNil())
		Expect(ms).To(Equal(MemorySize(2 * 1024 * 1024)))

		err = ms.Set("1GB")
		Expect(err).To(BeNil())
		Expect(ms).To(Equal(MemorySize(1 * 1024 * 1024 * 1024)))

		err = ms.Set("1024") // No suffix, treat as bytes
		Expect(err).To(BeNil())
		Expect(ms).To(Equal(MemorySize(1024)))

		err = ms.Set(" 64MB ") // Test with leading/trailing spaces
		Expect(err).To(BeNil())
		Expect(ms).To(Equal(MemorySize(64 * 1024 * 1024)))
	})

	t.Run("invalid inputs", func(t *testing.T) {
		var ms MemorySize

		// Test inputs with invalid values
		Expect(ms.Set("10XYZ")).To(Not(BeNil())) // Unknown unit
		Expect(ms.Set("ABC")).To(Not(BeNil()))   // Non-numeric input
		Expect(ms.Set("")).To(Not(BeNil()))      // Empty input
		Expect(ms.Set("-5MB")).To(Not(BeNil()))  // Negative value
	})

	t.Run("boundary cases", func(t *testing.T) {
		var ms MemorySize

		// Test extremely large numbers
		err := ms.Set("1099511627776GB") // 1 PB (petabyte, very large)
		Expect(err).To(BeNil())

		// Overflow handling (in practice, you'd want to handle overflow explicitly)
		err = ms.Set("9223372036854775808") // Larger than int64 max
		Expect(err).To(Not(BeNil()))
	})
}

func TestMemorySize_AsBytes(t *testing.T) {
	RegisterTestingT(t) // Register Gomega for the test

	var ms MemorySize

	// Set a value and check its string representation
	ms = 128 * 1024
	Expect(ms.ToBytes()).To(Equal(131072))

	ms = 2 * 1024 * 1024
	Expect(ms.ToBytes()).To(Equal(2097152))
}
