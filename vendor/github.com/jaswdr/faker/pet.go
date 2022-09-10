package faker

var (
	dogNames = []string{"Alfie", "Archie", "Bailey", "Banjo", "Barney", "Baxter", "Bear", "Beau", "Bella", "Benji", "Bentley", "Billie", "Billy", "Bonnie", "Bruce", "Bruno", "Buddy", "Buster", "Charlie", "Chester", "Chilli", "Chloe", "Cleo", "Coco", "Cookie", "Cooper", "Daisy", "Dexter", "Diesel", "Duke", "Ella", "Ellie", "Frankie", "George", "Gus", "Harley", "Harry", "Harvey", "Henry", "Holly", "Honey", "Hugo", "Jack", "Jasper", "Jax", "Jessie", "Jet", "Leo", "Lexi", "Lilly", "Lily", "Loki", "Lola", "Louie", "Louis", "Lucky", "Lucy", "Lulu", "Luna", "Maggie", "Marley", "Max", "Mia", "Millie", "Milly", "Milo", "Missy", "Molly", "Monty", "Murphy", "Nala", "Ollie", "Oscar", "Penny", "Pepper", "Pippa", "Poppy", "Ralph", "Rex", "Rocky", "Rosie", "Roxy", "Ruby", "Rusty", "Sam", "Sasha", "Scout", "Shadow", "Simba", "Sophie", "Stella", "Teddy", "Tilly", "Toby", "Willow", "Winston", "Zeus", "Ziggy", "Zoe"}
	catNames = []string{"Bella", "Tigger", "Chloe", "Shadow", "Luna", "Oreo", "Oliver", "Kitty", "Lucy", "Molly", "Jasper", "Smokey", "Gizmo", "Simba", "Tiger", "Charlie", "Angel", "Jack", "Lily", "Peanut", "Toby", "Baby", "Loki", "Midnight", "Milo", "Princess", "Sophie", "Harley", "Max", "Missy", "Rocky", "Zoe", "CoCo", "Misty", "Nala", "Oscar", "Pepper", "Sasha", "Buddy", "Pumpkin", "Kiki", "Mittens", "Bailey", "Callie", "Lucky", "Patches", "Simon", "Garfield", "George", "Maggie", "Sammy", "Sebastian", "Boots", "Cali", "Felix", "Lilly", "Phoebe", "Sassy", "Tucker", "Bandit", "Dexter", "Fiona", "Jake", "Precious", "Romeo", "Snickers", "Socks", "Daisy", "Gracie", "Lola", "Sadie", "Sox", "Casper", "Fluffy", "Marley", "Minnie", "Sweetie", "Ziggy", "Belle", "Blackie", "Chester", "Frankie", "Ginger", "Muffin", "Murphy", "Rusty", "Scooter", "Batman", "Boo", "Cleo", "Izzy", "Jasmine", "Mimi", "Sugar", "Cupcake", "Dusty", "Leo", "Noodle", "Panda", "Peaches"}
)

// Pet is a faker struct for Pet
type Pet struct {
	Faker *Faker
}

// Dog returns a fake dog name for App
func (p Pet) Dog() string {
	return p.Faker.RandomStringElement(dogNames)
}

// Cat returns a fake cat name for App
func (p Pet) Cat() string {
	return p.Faker.RandomStringElement(catNames)
}

// Name returns a fake pet name for App
func (p Pet) Name() string {
	petNames := []string{}
	petNames = append(petNames, catNames...)
	petNames = append(petNames, dogNames...)
	return p.Faker.RandomStringElement(petNames)
}
