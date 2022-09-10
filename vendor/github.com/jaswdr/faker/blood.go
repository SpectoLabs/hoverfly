package faker

// Blood is a faker struct for Blood
type Blood struct {
	Faker *Faker
}

var ( 
		bloodTypes = [] string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
)

// Name returns a Blood name for Blood
func (f Blood) Name() string {
	return f.Faker.RandomStringElement(bloodTypes)
}

