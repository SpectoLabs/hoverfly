package faker

import (
	"strings"
)

// Crypto is a faker struct for generating bitcoin data
type Crypto struct {
	Faker *Faker
}

var (
	bitcoinMin = 26
	bitcoinMax = 35
	ethLen     = 42
	ethPrefix  = "0x"
)

// Checks whether the ascii value provided is in the exclusion for bitcoin.
func (c Crypto) isInExclusionZone(ascii int) bool {
	switch ascii {
	// Ascii for uppercase letter "O", uppercase letter "I", lowercase letter "l", and the number "0"
	case
		48,
		73,
		79,
		108:
		return true
	}
	return false
}

// algorithmRange decides whether to get digit, uppercase, or lowercase. returns the ascii range to do IntBetween on
func (c Crypto)  algorithmRange() (int, int) {
	dec := c.Faker.IntBetween(0, 2)
	if dec == 0 {
		// digit
		return 48, 57
	} else if dec == 1 {
		// upper
		return 65, 90
	}
	// lower
	return 97, 122
}

// generateBicoinAddress returns a bitcoin address with a given prefix and length
func (c Crypto) generateBicoinAddress(length int, prefix string, f *Faker) string {
	address := []string{prefix}

	for i := 0; i < length; i++ {
		asciiStart, asciiEnd := c.algorithmRange()
		val := f.IntBetween(asciiStart, asciiEnd)
		if c.isInExclusionZone(val) {
			val++
		}
		address = append(address, string(rune(val)))
	}
	return strings.Join(address, "")
}

// P2PKHAddress generates a P2PKH bitcoin address.
func (c Crypto) P2PKHAddress() string {
	length := c.Faker.IntBetween(bitcoinMin, bitcoinMax)
	// subtract 1 for prefix
	return c.generateBicoinAddress(length-1, "1", c.Faker)
}

// P2PKHAddressWithLength generates a P2PKH bitcoin address with specified length.
func (c Crypto) P2PKHAddressWithLength(length int) string {
	return c.generateBicoinAddress(length-1, "1", c.Faker)
}

// P2SHAddress generates a P2SH bitcoin address.
func (c Crypto) P2SHAddress() string {
	length := c.Faker.IntBetween(bitcoinMin, bitcoinMax)
	// subtract 1 for prefix
	return c.generateBicoinAddress(length-1, "3", c.Faker)
}

// P2SHAddressWithLength generates a P2PKH bitcoin address with specified length.
func (c Crypto) P2SHAddressWithLength(length int) string {
	return c.generateBicoinAddress(length-1, "3", c.Faker)
}

// Bech32Address generates a Bech32 bitcoin address
func (c Crypto) Bech32Address() string {
	length := c.Faker.IntBetween(bitcoinMin, bitcoinMax)
	// subtract 1 for prefix
	return c.generateBicoinAddress(length-3, "bc1", c.Faker)
}

// Bech32AddressWithLength generates a Bech32 bitcoin address with specified length.
func (c Crypto) Bech32AddressWithLength(length int) string {
	return c.generateBicoinAddress(length-3, "bc1", c.Faker)
}

// BitcoinAddress returns an address of either Bech32, P2PKH, or P2SH type.
func (c Crypto) BitcoinAddress() string {
	dec := c.Faker.IntBetween(0, 2)
	if dec == 0 {
		return c.Bech32Address()
	} else if dec == 1 {
		return c.P2SHAddress()
	}
	return c.P2PKHAddress()
}

// EtheriumAddress returns a hexadecimal ethereum address of 42 characters.
func (c Crypto) EtheriumAddress() string {
	address := []string{ethPrefix}

	for i := 0; i < ethLen-2; i++ {
		asciiStart, asciiEnd := c.algorithmRange()
		val := c.Faker.IntBetween(asciiStart, asciiEnd)
		address = append(address, string(rune(val)))
	}
	return strings.Join(address, "")
}
