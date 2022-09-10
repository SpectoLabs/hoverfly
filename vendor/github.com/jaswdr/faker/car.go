package faker

import "strconv"

var (
	carMakers            = []string{"Acura", "Alfa Romeo", "Audi", "BMW", "Bentley", "Buick", "Cadillac", "Chevrolet", "Chrysler", "Dodge", "Fiat", "Ford", "GMC", "Genesis", "Honda", "Hyundai", "Infiniti", "Jaguar", "Jeep", "Kia", "Land Rover", "Lexus", "Lincoln", "Maserati", "Mazda", "Mercedes-Benz", "Mini", "Mitsubishi", "Nissan", "Polestar", "Porsche", "Ram", "Saab", "Smart", "Subaru", "Tesla", "Toyota", "Volkswagen", "Volvo"}
	carModels            = []string{"Q3", "Malibu", "Escalade ESV", "Corvette", "RLX", "Silverado 2500 HD Crew Cab", "3 Series", "Pacifica", "Colorado Crew Cab", "X3", "TLX", "Silverado 3500 HD Crew Cab", "7 Series", "Fusion", "Envision", "SQ5", "R8", "Traverse", "MDX", "QX80", "Encore", "Sierra 2500 HD Crew Cab", "Insight", "XT6", "XT5", "XT4", "Enclave", "Q5", "Santa Fe", "EcoSport", "Escape", "Mustang", "Sonata", "Edge", "Camaro", "Kona Electric", "Equinox", "Sierra 3500 HD Crew Cab", "Gladiator", "X7", "CT6-V", "A7", "Blazer", "F150 SuperCrew Cab", "Suburban", "Civic", "Compass", "Escalade", "Voyager", "Accord Hybrid", "Terrain", "Spark", "Sierra 1500 Crew Cab", "NEXO", "Veloster", "Silverado 1500 Crew Cab", "G70", "CT5", "Odyssey", "Elantra GT", "RDX", "Yukon XL", "Ranger SuperCab", "Expedition MAX", "Kona", "QX50", "Durango", "Yukon", "Palisade", "Ridgeline", "Cherokee", "Bolt EV", "Expedition", "Elantra", "Passport", "Charger", "Accord", "QX60", "Venue", "Pilot", "Grand Cherokee", "Tahoe", "Acadia", "Impala", "CR-V", "X5", "Q60", "Ranger SuperCrew", "Trax", "Ioniq Plug-in Hybrid", "E-PACE", "Tucson", "Explorer", "HR-V", "I-PACE", "Q50", "G80", "F-PACE", "Renegade", "Accent"}
	carCategories        = []string{"SUV", "Sedan", "Coupe", "Convertible", "Hatchback", "Pickup", "Van", "Minivan", "Wagon"}
	carFuelTypes         = []string{"Bio Gas", "Diesel", "Eletric", "Ethanol", "Hybrid", "Petrol"}
	carTransmissionGears = []string{"Automatic", "CVT", "Eletronic", "Manual", "Semi-auto", "Tiptronic"}
	carSeries            = []string{"KР", "ВI", "ВO", "АA", "EА", "BА", "PЕ", "НA", "IB", "KА", "KK", "OМ", "АM", "TА", "HI", "ОA", "CК", "PВ", "КC", "CА", "TЕ", "XА", "XО", "XМ", "MА", "МK", "MО"}
	countryRegions       = []string{"АK", "АB", "АC", "AЕ", "AН", "АM", "АO", "АP", "АT", "AА", "АI", "BА", "ВB", "BС", "ВE", "BН", "ВI", "BК", "СH", "ВM", "ВO", "АX", "ВT", "ВX", "CА", "CВ", "СE"}
)

// Car is a faker struct for Car
type Car struct {
	Faker *Faker
}

// Maker will return a random car maker
func (c Car) Maker() string {
	return c.Faker.RandomStringElement(carMakers)
}

// Model will return a random car model
func (c Car) Model() string {
	return c.Faker.RandomStringElement(carModels)
}

// Category will return a random car category
func (c Car) Category() string {
	return c.Faker.RandomStringElement(carCategories)
}

// FuelType will return a random car fuel type
func (c Car) FuelType() string {
	return c.Faker.RandomStringElement(carFuelTypes)
}

// TransmissionGear will return a random car transmission gear
func (c Car) TransmissionGear() string {
	return c.Faker.RandomStringElement(carTransmissionGears)
}

// Plate will return a random car plate
func (c Car) Plate() string {
	return c.Faker.RandomStringElement(countryRegions) + strconv.Itoa(c.Faker.RandomNumber(4)) + c.Faker.RandomStringElement(carSeries)
}
