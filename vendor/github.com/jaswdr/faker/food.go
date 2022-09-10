package faker

var (
	fruits     = []string{"Abiu", "Açaí", "Acerola", "Ackee", "Apple", "Apricot", "Avocado", "Banana", "Bilberry", "Blackberry", "Blackcurrant", "Black sapote", "Blueberry", "Boysenberry", "Breadfruit", "Buddha's hand", "Cactus pear", "Cempedak", "Crab apple", "Currant", "Cherry", "Cherimoya", "Chico fruit", "Cloudberry", "Coco De Mer", "Coconut", "Cranberry", "Damson", "Date", "Dragonfruit", "Durian", "Egg Fruit", "Elderberry", "Feijoa", "Fig", "Goji berry", "Gooseberry", "Grape", "Grewia asiatica (phalsa or falsa)", "Grapefruit", "Guava", "Honeyberry", "Huckleberry", "Jabuticaba", "Jackfruit", "Jambul", "Japanese plum", "Jostaberry", "Jujube", "Juniper berry", "Kiwano", "Kiwifruit", "Kumquat", "Lemon", "Lime", "Loganberry", "Loquat", "Longan", "Lulo", "Lychee", "Mamey Apple", "Mamey Sapote", "Mango", "Mangosteen", "Marionberry", "Melon", "Miracle fruit", "Monstera deliciosa", "Mulberry", "Nance", "Nectarine", "Orange", "Papaya", "Passionfruit", "Peach", "Pear", "Persimmon", "Plantain", "Plum", "Pineapple", "Pineberry", "Plumcot", "Pomegranate", "Pomelo", "Purple mangosteen", "Quince", "Raspberry", "Rambutan", "Redcurrant", "Rose apple", "Salal", "Salak", "Satsuma", "Shine Muscat or Vitis Vinifera", "Soursop", "Star apple", "Star fruit", "Strawberry", "Surinam cherry", "Tamarillo", "Tamarind", "Tangelo", "Tayberry", "Tomato", "Ugli fruit", "White currant", "White sapote", "Yuzu"}
	vegetables = []string{"Amaranth Leaves", "Arrowroot", "Artichoke", "Arugula", "Asparagus", "Bamboo Shoots", "Green Beans", "Beets", "Belgian Endive", "Bitter Melon", "Bok Choy", "Broadbeans", "Broccoli", "Broccoli Rabe", "Brussel Sprouts", "Green Cabbage ", "Red Cabbage", "Carrot", "Cassava", "Cauliflower", "Celeriac", "Celery", "Chayote", "Chicory", "Collards", "Corn", "Crookneck", "Cucumber", "Daikon", "Dandelion Greens", "Soybeans Edamame", "Eggplant", "Fennel", "Fiddleheads", "Ginger Root", "Horseradish", "Jicama", "Kale", "Kohlrabi", "Leeks", "Iceberg Lettuce", "Leaf Lettuce", "Romaine Lettuce", "Mushrooms", "Mustard Greens", "Okra", "Red Onion", "Parsnip", "Green Peas", "Green Pepper", "Sweet Red Pepper", "Red Potato", "White Potato", "Yellow Potato", "Pumpkin", "Radicchio", "Radishes", "Rutabaga", "Salsify", "Shallots", "Snow Peas", "Sorrel", "Spaghetti Squash", "Spinach", "Squash, Butternut", "Sugar Snap Peas", "Sweet Potato", "Swiss Chard", "Tomatillo", "Tomato", "Turnip", "Watercress", "Yam Root", "Zucchini"}
)

// Food is a faker struct for Food
type Food struct {
	Faker *Faker
}

// Fruit returns a fake fruit for Food
func (f Food) Fruit() string {
	return f.Faker.RandomStringElement(fruits)
}

// Vegetable returns a fake fruit for Food
func (f Food) Vegetable() string {
	return f.Faker.RandomStringElement(vegetables)
}
