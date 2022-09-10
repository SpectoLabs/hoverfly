package faker

import "strings"

// BinaryString is the faker struct for BinaryString
type BinaryString struct {
	faker *Faker
}

// BinaryString returns a random binarystring of given input length
func (bn BinaryString) BinaryString(length int) string {
	var bs strings.Builder
	for i := 0; i < length; i++ {
		bs.WriteString(bn.faker.RandomStringElement([]string{"0", "1"}))
	}
	return bs.String()
}
