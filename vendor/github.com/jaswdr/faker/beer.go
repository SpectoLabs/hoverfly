package faker

import "strconv"

var (
	beerNames  = []string{"Dale’s Pale Ale", "Breckenridge Vanilla Porter", "Brooklyn Brewery Lager", "Surly Brewing Darkness", "New Belgium Fat Tire", "Gigantic IPA", "NoDa Hop Drop n Roll", "Sam Adams Boston Lager", "Green Flash Palate Wrecker", "Dogfish Head 90 Minute IPA", "Pipeworks Citra", "Widmer Brothers Hefeweizen", "The Bruery Saison Rue", "Foothills Brewing Sexual Chocolate", "Avery Uncle Jacob’s Stout", "The Alchemist Focal Banger", "Hill Farmstead Abner", "Westbrook Gose", "Firestone Walker Union Jack IPA", "Highland Cold Mountain Winter Ale", "Sierra Nevada Pale Ale", "Allagash White", "Anchor Steam Beer", "Alpine Duet IPA", "Russian River Supplication", "Craftsman Cabernale", "Bell’s Two Hearted", "Deschutes Black Butte Porter", "Half Acre Daisy Cutter", "Smuttynose Finest Kind IPA", "Hair of the Dog Adam", "AleSmith Horny Devil", "21st Amendment Bitter American", "Stone IPA", "Tröegs Nugget Nectar", "Ballast Point Sculpin", "Upslope Brown Ale", "Rogue Shakespeare Oatmeal Stout", "Saint Arnold Fancy Lawnmower", "DC Brau On the Wings Of Armageddon", "Haymarket Angry Birds Rye IPA", "Capital Autumnal Fire Doppelbock", "Fullsteam Carver", "Green Flash Hop Head Red", "Russian River Blind Pig IPA", "Revolution Anti-Hero IPA", "Bell’s Hop Slam", "Great Lakes Edmund Fitzgerald Porter", "Jolly Pumpkin La Roja", "Toppling Goliath PseudoSue", "Lagunitas Brown Shugga", "Avery Rumpkin", "Firestone Walker Velvet Merkin", "Boulevard Tank 7", "Founders Red’s Rye", "Schlafly Pumpkin Ale", "Perennial Artisan Ales Abraxas Imperial Stout", "Three Floyds Zombie Dust", "Wicked Weed Serenity", "Stone Enjoy By… IPA", "The Bruery Sans Pagaie", "Firestone Walker Wookey Jack", "Cascade Brewing Apricot Ale", "Odell 90 Shilling Ale", "Left Hand Milk Stout Nitro", "Kern River Brewing Citra DIPA", "New Holland Dragon’s Milk", "Jester King Boxer’s Revenge", "Funky Buddha Maple Bacon Coffee Porter", "New Glarus Brewing Serendipity", "Westbrook Mexican Cake", "Alpine Great", "Terrapin Wake n Bake", "Wild Heaven Eschaton", "Ten Fidy", "New Belgium Lips of Faith La Folie", "Jai Alai IPA", "Founders Breakfast Stout", "Allagash Curieux", "Lagunitas IPA", "Great Divide Yeti", "Hill Farmstead Everett Porter", "Samuel Adams Utopias", "Troegs Mad Elf", "Dark Horse Plead the 5th", "Clown Shoes Undead Party Crasher", "Brewery Ommegang Three Philosophers", "North Coast Old Rasputin", "Avery Mephistopheles Stout", "Goose Island Bourbon County Stout", "Sierra Nevada Bigfoot Barleywine-Style Ale", "Firestone Walker Parabola", "Victory Prima Pils", "Lost Abbey Duck Duck Gooze", "Cigar City Hanaphu Imperial Stout", "Founders KBS (Kentucky Breakfast Stout)", "The Alchemist Heady Topper", "Russian River Pliny the Elder", "Floyds Dark Lord"}
	beerStyles = []string{"Altbier", "Amber ale", "Barley wine", "Berliner Weisse", "Bière de Garde", "Bitter", "Blonde Ale", "Bock", "Brown ale", "California Common/Steam Beer", "Cream Ale", "Dortmunder Export", "Doppelbock", "Dunkel", "Dunkelweizen", "Eisbock", "Flanders red ale", "Golden/Summer ale", "Gose", "Gueuze", "Hefeweizen", "Helles", "India pale ale", "Kölsch", "Lambic", "Light ale", "Maibock/Helles bock", "Malt liquor", "Mild", "Oktoberfestbier/Märzenbier", "Old ale", "Oud bruin", "Pale ale", "Pilsener/Pilsner/Pils", "Porter", "Red ale", "Roggenbier", "Saison", "Scotch ale", "Stout", "Schwarzbier", "Vienna lager", "Witbier", "Weissbier", "Weizenbock", "Fruit beer", "Herb and spiced beer", "Honey beer", "Rye Beer", "Smoked beer", "Vegetable beer", "Wild beer", "Wood-aged beer"}
	beerHops   = []string{"Admiral Hops", "Agnus Hops", "Ahtanum Hops", "AlphAroma Hops", "Amarillo Hops", "Amethyst Hops", "Apollo Hops", "Aramis Hops", "Atlas Hops", "Aurora Hops", "Beata Hops", "Belma Hops", "Bitter Gold Hops", "Boadicea Hops", "Bobek Hops", "Bouclier Hops", "Bramling Cross Hops", "Bravo Hops", "Brewers Gold Hops", "British Kent Goldings Hops", "Bullion Hops", "Calicross Hops", "California Cluster Hops", "Calypso Hops", "Cascade Hops", "Cashmere Hops", "Cekin Hops", "Celeia Hops", "Centennial Hops", "Challenger Hops", "Chelan Hops", "Chinook Hops", "Cicero Hops", "Citra Hops", "Cluster Hops", "Cobb’s Golding Hops", "Columbia Hops", "Columbus Hops", "Comet Hops", "Crystal Hops", "Dana Hops", "Delta Hops", "Dr. Rudi Hops", "Early Green Hops", "El Dorado Hops", "Ella Hops", "Endeavour Hops", "Equinox Hops", "Eroica Hops", "Falconer's Flight Hops", "First Gold Hops", "Flyer Hops", "Fuggle Hops", "Galaxy Hops", "Galena Hops", "Glacier Hops", "Golding Hops", "Green Bullet Hops", "Hallertau Hops", "HBC 342 Experimental Hops", "HBC 472 Experimental Hops", "Helga Hops", "Herald Hops", "Herkules Hops", "Hersbrucker Hops", "Horizon Hops", "Huell Melon Hops", "Idaho 7 Hops", "Idaho Gem™ Hops", "Jester Hops", "Junga Hops", "Kazbek Hops", "Kohatu Hops", "Liberty Hops", "Lubelski Hops", "Magnum Hops", "Mandarina Bavaria Hops", "Mathon Hops", "Marynka Hops", "Medusa™ Hops", "Meridian Hops", "Merkur Hops", "Millennium Hops", "Mittelfruh Hops", "Mosaic Hops", "Motueka Hops", "Mt. Hood Hops", "Mt. Rainier Hops", "Multihead Hops", "Nelson Sauvin Hops", "Neo1 Hops", "Newport Hops", "Northdown Hops", "Northern Brewer Hops", "Nugget Hops", "Opal Hops", "Orbit Hops", "Orion Hops", "Outeniqua Hops", "Pacific Gem Hops", "Pacific Jade Hops", "Pacific Sunrise Hops", "Pacifica Hops", "Palisade Hops", "Perle Hops", "Phoenix Hops", "Pilgrim Hops", "Pilot Hops", "Pioneer Hops", "Polaris Hops", "Premiant Hops", "Pride of Ringwood Hops", "Progress Hops", "Rakau Hops", "Riwaka Hops", "Saaz Hops", "Sabro Hops / Ron Mexico", "Santiam Hops", "Saphir Hops", "Satus Hops", "Select Hops", "Serebrianka Hops", "Simcoe Hops", "Sladek Hops", "Smaragd Hops", "Sonnet Hops", "Sorachi Ace Hops", "Southern Brewer Hops", "Southern Cross Hops", "Southern Promise Hops", "Southern Star Hops", "Sovereign Hops", "Spault Hops", "Spaulter Select Hops", "Sterling Hops", "Strata Hops", "Strickelbract Hops", "Strisselspault Hops", "Styrian Gold Hops", "Styrian Golding Hops", "Summer Hops", "Summit Hops", "Super Galena Hops", "Super Pride Hops", "Sussex Hops", "Sybilla Hops", "Sylva Hops", "Tahoma Hops", "Tardif de Burgogne Hops", "Target Hops", "Taurus Hops", "Teamaker Hops", "Tettnanger Hops", "Tillicum Hops", "Topaz Hops", "Tradition Hops", "Triple Pearl Hops", "Triskel Hops", "Ultra Hops", "Universal Hops", "Vanguard Hops", "Victoria Hops", "Vic Secret Hops", "Viking Hops", "Vital Hops", "Vojvodina Hops", "Wai-iti Hops", "Waimea Hops", "Wakatu Hops", "Warrior Hops", "Whitbread Goldings Hops", "Willamette Hops", "Yakima Cluster Hops", "Yakima Gold Hops", "Zappa Hops", "Zenith Hops", "Zythos Hops"}
	beerMalts  = []string{"Pale Malt", "Wheat Malt", "Rye Malt", "Vienna Malt", "Munich Malt", "Carapils", "Caramel/Crystal 10", "Caramel/Crystal 40", "Caramel/Crystal 60", "Caramel/Crystal 120", "Victory Malt", "Special Roast", "Chocolate Malt", "Roasted Barley", "Black Barley", "Black Patent", "German Pale", "Weizen", "Wiener", "Munchener", "Crystal", "Carafa Special", "German Acidulated", "German Melanoidin", "Belgian Pilsner", "Aromatic", "Belgian Special", "Biscuit Malt", "CaraVienne", "CaraMunich", "Rauchmalt 2", "Honey Malt", "Peated Malt"}
)

// Beer is a faker struct for Beer
type Beer struct {
	Faker *Faker
}

// Name will return a random beer name
func (b Beer) Name() string {
	return b.Faker.RandomStringElement(beerNames)
}

// Style will return a random beer style
func (b Beer) Style() string {
	return b.Faker.RandomStringElement(beerStyles)
}

// Hop will return a random beer hop
func (b Beer) Hop() string {
	return b.Faker.RandomStringElement(beerHops)
}

// Malt will return a random beer malt
func (b Beer) Malt() string {
	return b.Faker.RandomStringElement(beerMalts)
}

// Alcohol will return a random beer alcohol level between 2.0 and 10.0
func (b Beer) Alcohol() string {
	return strconv.FormatFloat(b.Faker.RandomFloat(2, 2.0, 10.0), 'f', 1, 64) + "%"
}

// Ibu will return a random beer ibu value between 10 and 100
func (b Beer) Ibu() string {
	return strconv.Itoa(b.Faker.IntBetween(10, 100)) + " IBU"
}

// Blg will return a random beer blg between 5.0 and 20.0
func (b Beer) Blg() string {
	return strconv.FormatFloat(b.Faker.RandomFloat(2, 5.0, 20.0), 'f', 1, 64) + "°Blg"
}
