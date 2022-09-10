package faker

// Boolean is a faker struct for Boolean
type Boolean struct {
	Faker *Faker
}

// Bool returns a fake bool for Faker
func (b Boolean) Bool() bool {
	return b.Faker.IntBetween(0, 100) > 50
}

// BoolWithChance returns true with a given percentual chance that the value is true, otherwise returns false
func (b Boolean) BoolWithChance(chanceTrue int) bool {
	if chanceTrue <= 0 {
		return false
	} else if chanceTrue >= 100 {
		return true
	}

	return b.Faker.IntBetween(0, 100) < chanceTrue
}

// BoolInt returns a fake bool for Integer Boolean
func (b Boolean) BoolInt() int {
	return b.Faker.RandomIntElement([]int{0, 1})
}

// BoolString returns a fake bool for string Boolean
func (b Boolean) BoolString(firstArg string, secondArg string) string {
	boolean := []string{firstArg, secondArg}

	return b.Faker.RandomStringElement(boolean)
}
