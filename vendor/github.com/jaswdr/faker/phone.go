package faker

import (
	"fmt"
	"strings"
)

var (
	phoneFormats = []string{
		// International format
		"+1-{{areaCode}}-{{exchangeCode}}-####",
		"+1 ({{areaCode}}) {{exchangeCode}}-####",
		"+1-{{areaCode}}-{{exchangeCode}}-####",
		"+1.{{areaCode}}.{{exchangeCode}}.####",
		"+1{{areaCode}}{{exchangeCode}}####",
		// Standard formats
		"{{areaCode}}-{{exchangeCode}}-####",
		"({{areaCode}}) {{exchangeCode}}-####",
		"1-{{areaCode}}-{{exchangeCode}}-####",
		"{{areaCode}}.{{exchangeCode}}.####",
		"{{areaCode}}-{{exchangeCode}}-####",
		"({{areaCode}}) {{exchangeCode}}-####",
		"1-{{areaCode}}-{{exchangeCode}}-####",
		"{{areaCode}}.{{exchangeCode}}.####",
		// Extensions
		"{{areaCode}}-{{exchangeCode}}-#### x###",
		"({{areaCode}}) {{exchangeCode}}-#### x###",
		"1-{{areaCode}}-{{exchangeCode}}-#### x###",
		"{{areaCode}}.{{exchangeCode}}.#### x###",
		"{{areaCode}}-{{exchangeCode}}-#### x####",
		"({{areaCode}}) {{exchangeCode}}-#### x####",
		"1-{{areaCode}}-{{exchangeCode}}-#### x####",
		"{{areaCode}}.{{exchangeCode}}.#### x####",
		"{{areaCode}}-{{exchangeCode}}-#### x#####",
		"({{areaCode}}) {{exchangeCode}}-#### x#####",
		"1-{{areaCode}}-{{exchangeCode}}-#### x#####",
		"{{areaCode}}.{{exchangeCode}}.#### x#####"}

	tollFreeAreaCodes = []string{"800", "844", "855", "866", "877", "888"}

	tollFreeFormats = []string{ // Standard formats
		"{{tollFreeAreaCode}}-{{exchangeCode}}-####",
		"({{tollFreeAreaCode}}) {{exchangeCode}}-####",
		"1-{{tollFreeAreaCode}}-{{exchangeCode}}-####",
		"{{tollFreeAreaCode}}.{{exchangeCode}}.####"}
)

// Phone is a faker struct for Phone
type Phone struct {
	Faker *Faker
}

// AreaCode returns a fake area code for Phone
func (p Phone) AreaCode() (code string) {
	number1 := p.Faker.IntBetween(2, 9)
	number2 := p.Faker.RandomDigit()
	number3 := p.Faker.RandomDigitNot(number2)
	return fmt.Sprintf("%d%d%d", number1, number2, number3)
}

// ExchangeCode returns a fake exchange code for Phone
func (p Phone) ExchangeCode() (code string) {
	number1 := p.Faker.IntBetween(2, 9)
	number2 := p.Faker.RandomDigit()
	number3 := p.Faker.RandomDigit()

	if number2 == 1 {
		number3 = p.Faker.RandomDigitNot(1)
	}

	return fmt.Sprintf("%d%d%d", number1, number2, number3)
}

// Number returns a fake phone number for Phone
func (p Phone) Number() string {
	number := p.Faker.RandomStringElement(phoneFormats)

	// {{areaCode}}
	number = strings.Replace(number, "{{areaCode}}", p.AreaCode(), 1)

	// {{exchangeCode}}
	number = strings.Replace(number, "{{exchangeCode}}", p.ExchangeCode(), 1)

	return p.Faker.Numerify(number)
}

// TollFreeAreaCode returns a fake toll free area code for Phone
func (p Phone) TollFreeAreaCode() string {
	return p.Faker.RandomStringElement(tollFreeAreaCodes)
}

// ToolFreeNumber returns a fake tool free number for Phone
func (p Phone) ToolFreeNumber() string {
	number := p.Faker.RandomStringElement(tollFreeFormats)

	// {{tollFreeAreaCode}}
	number = strings.Replace(number, "{{tollFreeAreaCode}}", p.TollFreeAreaCode(), 1)

	// {{exchangeCode}}
	number = strings.Replace(number, "{{exchangeCode}}", p.ExchangeCode(), 1)

	return p.Faker.Numerify(number)
}

// E164Number returns a fake E164 phone number for Phone
func (p Phone) E164Number() string {
	// Excluding 0 since it cannot be any country code in E164
	firstDigit := p.Faker.IntBetween(1, 9)
	randomDigits := p.Faker.Numerify("##########")

	return fmt.Sprintf("+%d%s", firstDigit, randomDigits)
}
