package faker

var (
	appNames = []string{"App Your Service", "Appcentric", "Appcare", "Develapp", "Fingertip Freedom", "Winning Widgets", "Tap Into Apps", "Download Developers", "Download Digital", "Tool Kit Digital", "Tech Happy", "Appy Digital", "Handheld Help", "In Your Palm", "For Your Palm", "Fit For Fingertips", "Fingertrip", "Tap To Begin", "Tap Into Digital", "Download Dev", "Touchpoint", "Trained For Tech", "Digitize Design", "About Apps", "A Is For App", "Handheld Digital", "Apprecicreate", "Apptitude", "Appreciate Apps", "Appster", "Digiapp", "Good Apptitude", "App Association", "Appetite", "Take-With-You Tech", "App Tech", "Appetite", "Strong Appetite", "App Natural", "Develop Digital", "Digital Daredevil", "If You Build It", "Pocket Pro", "Iconic Inc.", "Icon Inc.", "Fingertip Tech", "Dare To Design", "Pocket Pros", "Digit Widget", "Build Better", "Dual Develop", "Amazing Apps", "Action Apps", "Application Station", "App Innovation", "Fun Apps", "Fantappstic", "App Command", "Strike Apps", "App Force", "Creative Applications", "App Fly", "Sure Apps", "App Door", "App Tray", "App Sure", "Rocket Apps", "App Place", "App Cafe", "Trippy Apps", "Appkey", "App Home", "Hot Apps", "App Focus", "App Possible", "App Leader", "Whip App", "App Works", "Good Apps", "Easy Apps", "App Source", "App Stage", "App Inspire", "Fire Apps", "App Flower", "App Dog", "Advance Apps", "Chatter Apps", "App Dream", "Bold Apps", "Boss Apps", "App Joy", "App Bullet", "App Cracker", "True Apps", "Feather Apps", "Real Apps", "App Whimsy", "Jewel Apps", "Image Apps", "Rifle Apps", "Next App", "Mobile Vibes", "Candy App", "Setup App", "Personality App", "Essential Web", "VitalApp", "Interact Mobile", "HelloWeb", "Network Moment", "Major Connection", "Billing Mobile", "SmartApp", "NoteWork", "Web Influence", "PowerPhone", "Chief Network", "Connection App", "WebTools", "Gamepad", "Mobile Stick", "Know The App", "WebChecker", "PassApp", "RobotSoft", "SmartCloud", "MobileHelp", "WebDesk", "EasyClick", "WeBox", "AppCan", "Smartum", "Smartio", "Smarter Web", "GrandMobile", "Technet", "RoboVoice", "TabletSoft", "E-APPy", "SkyApp", "WebMap", "BoostApp", "UserMobile", "CheapMobile", "WireSmart", "SwipeApp", "LiveBox", "WebGroup", "LinkApp", "OneClick", "MeetAll", "MomyApp", "Moboapp Developers", "Uniworld Games", "Raptor Games", "Gamers Republic", "Atomik Games", "Javatron Games", "Ultrasonic Apps", "Graviton Games", "Virtualsphere Mobile App Developers", "Javanation", "Telesoft Mobile App Developers", "Phantom Labs", "Rededge Creations", "Virtualyard Tech", "Coderant It Solutions", "Loopsoft Developers", "Clever Co App Developers", "Primeroyal App Creations", "Cyberville Tech", "Angularis Mobile App Developers", "Oceanfloat Technology", "Intelli-Ware Creations", "Aster Mobile App Developers", "Venus Hub", "Pilot Softwares", "Dominio Software Consult", "Customs Software Developers", "Cellarstars Mobile Developers", "Helevate Games", "Tetrabyte", "Monolith Games", "Selvo Games", "Metreality Games", "Hovertec Games", "Helicion Games", "Play Monkey Studios", "Digisphere Developers", "Revolt Games", "Sabre Games", "Savagechimp Games", "Spidermokey Concept", "Gravitones Games", "Clique18 Concepts", "Blackguard Gamea", "WireSmart", "SwipeApp", "LiveBox", "WebGroup", "LinkApp", "OneClick", "MeetAll", "MomyApp"}
)

// App is a faker struct for App
type App struct {
	Faker *Faker
}

// Name returns a fake app name for App
func (a App) Name() string {
	return a.Faker.RandomStringElement(appNames)
}

// Version returns a fake app version for App
func (a App) Version() string {
	return a.Faker.Numerify("v#.#.#")
}
