package fake

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type creditCard struct {
	vendor   string
	length   int
	prefixes []int
}

func (c creditCard) RandomPrefix() int {
	return c.prefixes[r.Intn(len(c.prefixes))]
}

var (
	creditCards = map[string]creditCard{
		"visa":       {"VISA", 16, []int{4539, 4556, 4916, 4532, 4929, 40240071, 4485, 4716, 4}},
		"mastercard": {"MasterCard", 16, []int{51, 52, 53, 54, 55}},
		"amex":       {"American Express", 15, []int{34, 37}},
		"discover":   {"Discover", 16, []int{6011}},
	}
	creditCardsKeys = make([]string, len(creditCards))
)

func init() {
	n := 0
	for key := range creditCards {
		creditCardsKeys[n] = key
		n++
	}
	sort.Strings(creditCardsKeys)
}

// CreditCardType returns one of the following credit values:
// VISA, MasterCard, American Express and Discover
func CreditCardType() string {
	n := len(creditCards)
	var vendors []string
	for _, cc := range creditCards {
		vendors = append(vendors, cc.vendor)
	}

	return vendors[r.Intn(n)]
}

// CreditCardNum generated credit card number according to the card number rules
func CreditCardNum(vendor string) string {
	if vendor != "" {
		vendor = strings.ToLower(vendor)
	} else {
		vendor = creditCardsKeys[r.Intn(len(creditCardsKeys))]
	}
	card, ok := creditCards[vendor]
	if !ok {
		panic(fmt.Sprintf("unsupported vendor %q", vendor))
	}

	prefix := strconv.Itoa(card.RandomPrefix())
	num := []rune(prefix)
	for i := 0; i < card.length-len(prefix)-1; i++ {
		num = append(num, rune(strconv.Itoa(r.Intn(10))[0]))
	}
	num = append(num, creditCardNumChecksum(num))

	return string(num)
}

func creditCardNumChecksum(num []rune) rune {
	// See: https://en.wikipedia.org/wiki/Luhn_algorithm
	sum := 0
	pos := 0
	for i := len(num) - 1; i >= 0; i-- {
		n := int(num[i] - '0')
		if pos%2 == 0 {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		pos++
	}

	// https://en.wikipedia.org/wiki/Talk:Luhn_algorithm#Formula_error
	checksum := 10 - (sum%10)%10
	return rune(strconv.Itoa(checksum)[0])
}
