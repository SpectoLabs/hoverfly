package faker

var (
	gamerTags = []string{"EatBullets", "PR0_GGRAM3D", "CollateralDamage",
		"TheSickness", "Shoot2Kill", "Overkill", "Killspree", "MindlessKilling", "Born2Kill",
		"TheZodiac", "ZodiacKiller", "Osamaisback", "OsamasGhost", "T3rr0r1st", "ToySoldier",
		"MilitaryMan", "DeathSquad", "Veteranofdeath", "Angelofdeath", "Ebola", "MustardGas",
		"Knuckles", "KnuckleBreaker", "KnuckleDuster", "BloodyKnuckles", "JackTheRipper", "TedBundyHandsome",
		"Necromancer", "SmilingSadist", "ManicLaughter", "Tearsofjoy", "ShowMeUrguts", "KnifeInGutsOut",
		"Talklesswinmore", "Guillotine", "Decapitator", "TheExecutor", "BigKnives", "SharpKnives",
		"LocalBackStabber", "BodyParts", "BodySnatcher", "TheButcher", "meat", "ChopChop", "ChopSuey",
		"TheZealot", "VagaBond", "LoneAssailant", "9mm", "SemiAutomatic", "101WaysToMeetYourMaker",
		"SayHi2God", "Welcome2Hell", "HellNBack", "Dudemister", "MiseryInducing", "SmashDtrash",
		"TakinOutThaTrash", "StreetSweeper", "TheBully", "Getoutofmyway", "NoMercy4TheWeak", "Sl4ught3r",
		"HappyKilling", "HappyPurgeDay", "HappyPurging", "RiotStarter", "CantStop", "CantStopWontstop",
		"SweetPoison", "SimplyTheBest", "PuppyDrowner", "EatYourHeartOut", "RipYourHeartOut", "BloodDrainer",
		"AcidAttack", "AcidFace", "PetrolBomb", "Molotov", "TequilaSunrise", "TeKillaSunrise", "LocalGrimReaper",
		"SoulTaker", "DreamHaunter", "Grave", "YSoSerious", "Revenge", "Avenged", "BestServedCold", "HitNRUN",
		"Fastandfurious", "MrBlond", "TheKingIsDead", "TheNihilist", "Bad2TheBone", "OneShot", "SmokinAces", "DownInSmoke", "NoFun4U"}
)

// Gamer is a faker struct for Gamer
type Gamer struct {
	Faker *Faker
}

// Tag returns a fake gamer tag for Gamer
func (g Gamer) Tag() string {
	return g.Faker.RandomStringElement(gamerTags)
}
