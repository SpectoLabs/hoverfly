package faker

import (
	"strings"
	"time"
)

var (
	musicNameFormats = []string{
		"{{adverb}} {{verb}} {{noun}} {{ending}}",
		"{{verb}} {{adjective}} {{noun}}",
		"{{adjective}} {{noun}}",
		"{{adjective}} {{noun}} {{ending}}",
	}

	musicGenres = []string{
		"2-step", "4-beat", "Acid breaks", "Acid house", "Acid jazz", "Acid rock", "Acid techno", "Acid trance", "Aggrotech",
		"Alternative dance", "Alternative metal", "Alternative rock", "Ambient dub", "Ambient house", "Ambient techno", "Ambient", "Anarcho punk", "Anti-folk",
		"Art punk", "Art rock", "Asian Underground", "Avant-garde jazz", "Baggy", "Balearic Beat", "Baltimore Club", "Bassline", "Beat music",
		"Bebop", "Big beat", "Bitpop", "Black metal", "Boogie-woogie", "Boogie", "Bossa nova", "Bouncy house", "Bouncy techno",
		"Breakbeat hardcore", "Breakbeat", "Breakcore", "Breakstep", "British dance", "Britpop", "Broken beat", "Bubblegum dance", "Canterbury scene",
		"Cape jazz", "Celtic metal", "Celtic punk", "Celtic", "Chamber jazz", "Chicago house", "Chill out", "Chillwave", "Chinese rock",
		"Chiptune", "Christian metal", "Christian punk", "Christian rock", "Classic trance", "Coldwave", "Contemporary folk", "Continental Jazz", "Cool jazz",
		"Cosmic disco", "Cowpunk", "Crossover jazz", "Crossover thrash", "Crunk", "Crust punk", "Crustgrind", "Cybergrind", "D-beat",
		"Dance-pop", "Dance-punk", "Dance-rock", "Dark ambient", "Dark cabaret", "Dark electro", "Dark psytrance", "Dark Wave", "Darkcore",
		"Darkside jungle", "Darkstep", "Death industrial", "Death metal", "Deathcore", "Deathrock", "Deep house", "Desert rock", "Detroit techno",
		"Digital hardcore", "Disco house", "Disco polo", "Disco", "Diva house", "Dixieland", "Djent", "Doom metal", "Doomcore",
		"Downtempo", "Dream house", "Dream pop", "Dream trance", "Drone metal", "Drone", "Drum and bass", "Drumfunk", "Drumstep",
		"Dub", "Dubstep", "Dubstyle", "Dubtronica", "Dunedin Sound", "Dutch house", "EDM", "Electro backbeat", "Electro house",
		"Electro-grime", "Electro-industrial", "Electro", "Electroacoustic", "Electroclash", "Electronic art music", "Electronic rock", "Electronica", "Electronicore",
		"Electropop", "Electropunk", "Emo", "Epic doom", "Ethereal wave", "Ethnic electronica", "Euro disco", "Eurobeat", "Eurodance",
		"European free jazz", "Europop", "Experimental rock", "Filk", "Florida breaks", "Folk metal", "Folk punk", "Folk rock", "Folk",
		"Folktronica", "Freak folk", "Freakbeat", "Free tekno", "Freestyle house", "Freestyle", "French house", "Full on", "Funeral doom",
		"Funk metal", "Funky house", "Funky", "Futurepop", "Gabber", "Garage punk", "Garage rock", "Ghetto house", "Ghettotech",
		"Glam metal", "Glam rock", "Glitch", "Goregrind", "Gothic metal", "Gothic rock", "Grime", "Grindcore", "Groove metal",
		"Grunge", "Happy hardcore", "Hard bop", "Hard NRG", "Hard rock", "Hard trance", "Hardbag", "Hardcore punk", "Hardcore/Hard dance",
		"Hardstep", "Hardstyle", "Heavy metal", "Hi-NRG", "Hip house", "Horror punk", "House", "IDM", "Illbient",
		"Indie folk", "Indie pop", "Indie rock", "Indietronica", "Industrial folk", "Industrial metal", "Industrial rock", "Industrial", "Intelligent drum and bass",
		"Italo dance", "Italo disco", "Italo house", "Japanoise", "Jazz blues", "Jazz fusion", "Jazz rap", "Jazz rock", "Jazz-funk",
		"Jump-Up", "Jumpstyle", "Krautrock", "Laptronica", "Latin house", "Latin jazz", "Liquid funk", "Livetronica", "Lowercase",
		"Lo-fi", "Madchester", "Mainstream jazz", "Makina", "Math rock", "Mathcore", "Medieval metal", "Melodic death metal", "Metalcore",
		"Minimal house/Microhouse", "Minimal", "Modal jazz", "Moombahton", "Neo-bop jazz", "Neo-psychedelia", "Neo-swing", "Neofolk", "Neurofunk",
		"New Beat", "New jack swing", "New prog", "New rave", "New wave", "New-age", "Nintendocore", "No wave", "Noise pop",
		"Noise rock", "Noise", "Noisegrind", "Nortec", "Novelty ragtime", "Nu jazz", "Nu metal", "Nu skool breaks", "Nu-disco",
		"Oldschool jungle", "Orchestral jazz", "Orchestral Uplifting", "Paisley Underground", "Pop punk", "Pop rock", "Post-bop", "Post-Britpop", "Post-disco",
		"Post-grunge", "Post-hardcore", "Post-metal", "Post-punk revival", "Post-punk", "Post-rock", "Power electronics", "Power metal", "Power noise",
		"Power pop", "Powerviolence", "Progressive breaks", "Progressive drum & bass", "Progressive folk", "Progressive house", "Progressive metal", "Progressive rock", "Progressive techno",
		"Progressive", "Psybreaks", "Psychedelic folk", "Psychedelic rock", "Psychedelic trance", "Psychobilly", "Psyprog", "Punk jazz", "Punk rock",
		"Raga rock", "Ragga-jungle", "Raggacore", "Ragtime", "Rap metal", "Rap rock", "Rapcore", "Riot grrrl", "Rock and roll",
		"Rock in Opposition", "Sadcore", "Sambass", "Screamo", "Shibuya-kei", "Shoegaze", "Ska jazz", "Ska punk", "Skate punk",
		"Skweee", "Slowcore", "Sludge metal", "Smooth jazz", "Soft rock", "Soul jazz", "Sound art", "Southern rock", "Space disco",
		"Space house", "Space rock", "Speed garage", "Speed metal", "Speedcore", "Stoner rock", "Straight-ahead jazz", "Street punk", "Stride jazz",
		"Sufi rock", "Sung poetry", "Suomisaundi", "Surf rock", "Swing house", "Swing", "Symphonic metal", "Synthcore", "Synthpop",
		"Synthpunk", "Tech house", "Tech trance", "Technical death metal", "Techno-DNB", "Techno-folk", "Techno", "Technopop", "Techstep",
		"Tecno brega", "Terrorcore", "Third stream", "Thrash metal", "Thrashcore", "Toytown Techno", "Trad jazz", "Traditional doom", "Trance",
		"Trap", "Tribal house", "Trip hop", "Turbofolk", "Twee Pop", "Uplifting trance", "Vaporwave", "Viking metal", "Vocal house",
		"Vocal jazz", "Vocal trance", "West Coast jazz", "Western", "Witch House/Drag", "World fusion", "Worldbeat", "Yacht rock", "Yorkshire Bleeps and Bass"}

	musicNameAdverbs = []string{"Abnormally", "Absentmindedly", "Accidentally", "Acidly", "Actually", "Adventurously", "Afterwards", "Almost", "Always", "Angrily", "Annually", "Anxiously", "Arrogantly", "Awkwardly", "Badly",
		"Bashfully", "Beautifully", "Bitterly", "Bleakly", "Blindly", "Blissfully", "Boastfully", "Boldly", "Bravely", "Briefly", "Brightly", "Briskly", "Broadly", "Busily", "Calmly", "Carefully", "Carelessly", "Cautiously",
		"Certainly", "Cheerfully", "Clearly", "Cleverly", "Closely", "Coaxingly", "Colorfully", "Commonly", "Continually", "Coolly", "Correctly", "Courageously", "Crossly", "Cruelly", "Curiously", "Daily", "Daintily", "Dearly",
		"Deceivingly", "Deeply", "Defiantly", "Deliberately", "Delightfully", "Diligently", "Dimly", "Doubtfully", "Dreamily", "Easily", "Elegantly", "Energetically", "Enormously", "Enthusiastically", "Equally", "Especially",
		"Even", "Evenly", "Eventually", "Exactly", "Excitedly", "Extremely", "Fairly", "Faithfully", "Famously", "Far", "Fast", "Fatally", "Ferociously", "Fervently", "Fiercely", "Fondly", "Foolishly", "Fortunately", "Frankly",
		"Frantically", "Freely", "Frenetically", "Frightfully", "Fully", "Furiously", "Generally", "Generously", "Gently", "Gladly", "Gleefully", "Gracefully", "Gratefully", "Greatly", "Greedily", "Happily", "Hastily", "Healthily",
		"Heavily", "Helpfully", "Helplessly", "Highly", "Honestly", "Hopelessly", "Hourly", "Hungrily", "Immediately", "Innocently", "Inquisitively", "Instantly", "Intensely", "Intently", "Interestingly", "Inwardly", "Irritably",
		"Jaggedly", "Jealously", "Jovially", "Joyfully", "Joyously", "Jubilantly", "Judgementally", "Justly", "Kiddingly", "Kindheartedly", "Kindly", "Knottily", "Knowingly", "Knowledgeably", "Lazily", "Less", "Lightly", "Likely",
		"Limply", "Lively", "Loftily", "Longingly", "Loosely", "Loudly", "Lovingly", "Loyally", "Madly", "Majestically", "Meaningfully", "Mechanically", "Merrily", "Miserably", "Mockingly", "Monthly", "More", "Mortally", "Mostly",
		"Mysteriously", "Naturally", "Nearly", "Neatly", "Needily", "Nervously", "Never", "Nicely", "Noisily", "Not", "Obediently", "Obnoxiously", "Oddly", "Offensively", "Officially", "Often", "Only", "Openly", "Optimistically",
		"Overconfidently", "Owlishly", "Painfully", "Partially", "Patiently", "Perfectly", "Physically", "Playfully", "Politely", "Poorly", "Positively", "Potentially", "Powerfully", "Promptly", "Properly", "Questionably",
		"Quickly", "Quietly", "Quirkily", "Quizzically", "Rapidly", "Rarely", "Readily", "Really", "Reassuringly", "Recklessly", "Regularly", "Reluctantly", "Repeatedly", "Reproachfully", "Restfully", "Righteously", "Rightfully",
		"Rigidly", "Roughly", "Rudely", "Sadly", "Safely", "Scarcely", "Scarily", "Searchingly", "Sedately", "Seemingly", "Seldom", "Selfishly", "Separately", "Seriously", "Shakily", "Sharply", "Sheepishly", "Shrilly", "Shyly",
		"Silently", "Sleepily", "Slowly", "Smoothly", "Softly", "Solemnly", "Solidly", "Sometimes", "Soon", "Speedily", "Stealthily", "Sternly", "Strictly", "Successfully", "Suddenly", "Surprisingly", "Suspiciously", "Sweetly",
		"Swiftly", "Sympathetically", "Tenderly", "Tensely", "Terribly", "Thankfully", "Thoroughly", "Thoughtfully", "Tightly", "Tomorrow", "Too", "Tremendously", "Triumphantly", "Truly", "Truthfully", "Ultimately", "Unabashedly",
		"Unaccountably", "Unbearably", "Unethically", "Unexpectedly", "Unfortunately", "Unimpressively", "Unnaturally", "Unnecessarily", "Upbeat", "Upliftingly", "Upright", "Upside-down", "Upward", "Upwardly", "Urgently",
		"Usefully", "Uselessly", "Usually", "Utterly", "Vacantly", "Vaguely", "Vainly", "Valiantly", "Vastly", "Verbally", "Very", "Viciously", "Victoriously", "Violently", "Vivaciously", "Voluntarily", "Warmly", "Weakly",
		"Wearily", "Well", "Wetly", "Wholly", "Wildly", "Willfully", "Wisely", "Woefully", "Wonderfully", "Worriedly", "Wrongly", "Yawningly", "Yearly", "Yearningly", "Yesterday", "Yieldingly", "Youthfully", "Zealously",
		"Zestfully", "Zestily"}

	musicNameVerbs = []string{"Abiding", "Accelerating", "Accepting", "Accomplishing", "Achieving", "Acquiring", "Acting", "Activating", "Adapting", "Adding", "Addressing", "Administering", "Admiring", "Admitting", "Adopting", "Advising", "Affording", "Agreeing", "Alerting", "Allowing", "Altering", "Amusing",
		"Analyzing", "Announcing", "Annoying", "Answering", "Anticipating", "Apologizing", "Appearing", "Applauding", "Approving", "Arguing", "Arranging", "Arresting", "Arriving", "Asking", "Assembling", "Assisting", "Attaching", "Attacking", "Attracting", "Avoiding", "Awaking", "Backing",
		"Baking", "Balancing", "Banging", "Baring", "Bathing", "Bating", "Battling", "Beaming", "Bearing", "Beating", "Becoming", "Beging", "Beginning", "Behaving", "Beholding", "Belonging", "Bending", "Beting", "Biding", "Binding", "Biting", "Bleaching", "Bleeding", "Blessing", "Blinding",
		"Blinking", "Blowing", "Blushing", "Boasting", "Boiling", "Bolting", "Bombing", "Booking", "Boring", "Borrowing", "Bouncing", "Bowing", "Boxing", "Braking", "Branching", "Breaking", "Breathing", "Breeding", "Briefing", "Bring", "Broadcasting", "Bruising", "Brushing", "Bubbling",
		"Budgeting", "Building", "Bumping", "Burning", "Bursting", "Burying", "Busting", "Buying", "Buzzing", "Calculating", "Calling", "Camping", "Caring", "Carrying", "Carving", "Casting", "Cataloging", "Catching", "Causing", "Challenging", "Changing", "Charging", "Charting", "Chasing",
		"Cheating", "Checking", "Cheering", "Chewing", "Choking", "Choosing", "Chopping", "Claiming", "Clapping", "Clarifying", "Classifying", "Cleaning", "Clearing", "Cling", "Clipping", "Closing", "Clothing", "Coaching", "Coiling", "Collecting", "Coloring", "Combing", "Coming", "Commanding",
		"Communicating", "Comparing", "Competing", "Compiling", "Complaining", "Completing", "Composing", "Computing", "Conceiving", "Concentrating", "Conceptualizing", "Concerning", "Concluding", "Conducting", "Confessing", "Confronting", "Confusing", "Connecting", "Conserving", "Considering",
		"Consisting", "Consolidating", "Constructing", "Consulting", "Containing", "Continuing", "Contracting", "Controlling", "Converting", "Coordinating", "Copying", "Correcting", "Correlating", "Costing", "Coughing", "Counseling", "Counting", "Covering", "Cracking", "Crashing", "Crawling",
		"Creating", "Creeping", "Critiquing", "Crossing", "Crushing", "Crying", "Curing", "Curling", "Curving", "Cutting", "Cycling", "Damaging", "Dancing", "Daring", "Dealing", "Decaying", "Deceiving", "Deciding", "Decorating", "Defining", "Delaying", "Delegating", "Delighting", "Delivering",
		"Demonstrating", "Depending", "Describing", "Deserting", "Deserving", "Designing", "Destroying", "Detailing", "Detecting", "Determining", "Developing", "Devising", "Diagnosing", "Diging", "Directing", "Disagreeing", "Disappearing", "Disapproving", "Disarming", "Discovering", "Dispensing",
		"Displaying", "Disproving", "Dissecting", "Distributing", "Diverting", "Dividing", "Diving", "Doing", "Doubling", "Doubting", "Drafting", "Dragging", "Draining", "Dramatizing", "Drawing", "Dreaming", "Dressing", "Drinking", "Dripping", "Driving", "Dropping", "Drowning", "Drumming", "Drying",
		"Dusting", "Dwelling", "Earning", "Eating", "Editing", "Educating", "Eliminating", "Embarrassing", "Employing", "Emptying", "Encouraging", "Ending", "Enduring", "Enforcing", "Engineering", "Enhancing", "Enjoying", "Enlisting", "Ensuring", "Entering", "Entertaining", "Escaping", "Establishing",
		"Estimating", "Evaluating", "Examining", "Exceeding", "Exciting", "Excusing", "Executing", "Exercising", "Exhibiting", "Existing", "Expanding", "Expecting", "Expediting", "Experimenting", "Explaining", "Exploding", "Expressing", "Extending", "Extracting", "Facilitating", "Facing", "Fading",
		"Failing", "Fancying", "Fastening", "Faxing", "Fearing", "Feeding", "Feeling", "Fencing", "Fetching", "Fighting", "Filing", "Filling", "Filming", "Finalizing", "Financing", "Finding", "Firing", "Fitting", "Fixing", "Flapping", "Flashing", "Fleeing", "Fling", "Floating", "Flooding", "Flowering",
		"Flowing", "Flying", "Folding", "Following", "Fooling", "Forbidding", "Forcing", "Forecasting", "Foregoing", "Foreseeing", "Foretelling", "Forgetting", "Forgiving", "Forming", "Formulating", "Forsaking", "Framing", "Freezing", "Frightening", "Frying", "Gathering", "Gazing", "Generating",
		"Geting", "Giving", "Glowing", "Gluing", "Going", "Governing", "Grabbing", "Graduating", "Grating", "Greasing", "Greeting", "Grinding", "Grinning", "Griping", "Groaning", "Growing", "Guaranteeing", "Guarding", "Guessing", "Guiding", "Hammering", "Handing", "Handling", "Handwriting", "Hanging",
		"Happening", "Harassing", "Harming", "Hating", "Haunting", "Heading", "Healing", "Heaping", "Hearing", "Heating", "Helping", "Hiding", "Hitting", "Holding", "Hooking", "Hoping", "Hovering", "Huging", "Humming", "Hunting", "Hurrying", "Hurting", "Hypothesizing", "Identifying", "Ignoring",
		"Illustrating", "Imagining", "Implementing", "Impressing", "Improving", "Improvising", "Including", "Increasing", "Inducing", "Influencing", "Informing", "Initiating", "Injecting", "Injuring", "Inlaying", "Innovating", "Inputing", "Inspecting", "Inspiring", "Installing", "Instituting",
		"Instructing", "Insuring", "Integrating", "Intending", "Intensifying", "Interesting", "Interfering", "Interlaying", "Interpreting", "Interrupting", "Interviewing", "Introducing", "Inventing", "Inventorying", "Investigating", "Inviting", "Irritating", "Itching", "Jailing", "Jamming", "Jogging",
		"Joining", "Joking", "Judging", "Juggling", "Jumping", "Justifying", "Keeping", "Kicking", "Killing", "Kissing", "Kneeling", "Knocking", "Knotting", "Knowing", "Labeling", "Landing", "Lasting", "Laughing", "Launching", "Laying", "Leading", "Leaning", "Leaping", "Learning", "Leaving", "Lecturing",
		"Lending", "Letting", "Leveling", "Licensing", "Licking", "Lifting", "Lightening", "Lighting", "Liking", "Listening", "Listing", "Living", "Loading", "Locating", "Locking", "Logging", "Longing", "Looking", "Losing", "Loving", "Maintaining", "Making", "Managing", "Manipulating", "Manufacturing",
		"Mapping", "Marching", "Marketing", "Marking", "Marrying", "Matching", "Mating", "Mattering", "Meaning", "Measuring", "Meddling", "Mediating", "Meeting", "Melting", "Memorizing", "Mending", "Mentoring", "Milking", "Mining", "Misleading", "Missing", "Misspelling", "Mistaking", "Misunderstanding",
		"Mixing", "Moaning", "Modeling", "Modifying", "Monitoring", "Motivating", "Mourning", "Moving", "Mowing", "Muddling", "Muging", "Multiplying", "Murdering", "Nailing", "Naming", "Navigating", "Needing", "Negotiating", "Nesting", "Normalizing", "Noticing", "Noting", "Numbering", "Obeying",
		"Objecting", "Observing", "Obtaining", "Occurring", "Offending", "Offering", "Officiating", "Opening", "Operating", "Ordering", "Organizing", "Originating", "Overcoming", "Overdoing", "Overflowing", "Overhearing", "Overtaking", "Overthrowing", "Owing", "Owning", "Packing", "Paddling",
		"Painting", "Parking", "Participating", "Parting", "Passing", "Pasting", "Pausing", "Paying", "Pecking", "Pedaling", "Peeling", "Peeping", "Perceiving", "Perfecting", "Performing", "Permitting", "Persuading", "Phoning", "Photographing", "Picking", "Piloting", "Pinching", "Pioneering", "Placing",
		"Planing", "Planting", "Playing", "Pleading", "Pleasing", "Plugging", "Pointing", "Poking", "Polishing", "Popping", "Possessing", "Posting", "Pouring", "Practicing", "Praising", "Praying", "Preaching", "Preceding", "Predicting", "Preparing", "Prescribing", "Presenting", "Preserving", "Presiding",
		"Pressing", "Pretending", "Preventing", "Pricking", "Printing", "Processing", "Procuring", "Producing", "Professing", "Programing", "Progressing", "Projecting", "Promising", "Promoting", "Proofreading", "Proposing", "Protecting", "Providing", "Proving", "Publicizing", "Pulling", "Pumping",
		"Punching", "Puncturing", "Punishing", "Purchasing", "Pushing", "Putting", "Qualifying", "Questioning", "Queuing", "Quitting", "Racing", "Radiating", "Raining", "Raising", "Ranking", "Rating", "Reaching", "Reading", "Realigning", "Realizing", "Reasoning", "Receiving", "Recognizing",
		"Recommending", "Reconciling", "Recording", "Recruiting", "Reducing", "Referring", "Reflecting", "Refusing", "Regretting", "Regulating", "Rehabilitating", "Reigning", "Reinforcing", "Rejecting", "Rejoicing", "Relating", "Relaxing", "Releasing", "Relying", "Remaining", "Remembering",
		"Reminding", "Removing", "Rendering", "Reorganizing", "Repairing", "Repeating", "Replacing", "Replying", "Reporting", "Representing", "Reproducing", "Requesting", "Rescuing", "Researching", "Resolving", "Responding", "Restoring", "Restructuring", "Retiring", "Retrieving", "Returning",
		"Reviewing", "Revising", "Rhyming", "Riding", "Ring", "Rinsing", "Rising", "Risking", "Robbing", "Rocking", "Rolling", "Rotting", "Rubbing", "Ruining", "Ruling", "Running", "Rushing", "Sacking", "Sailing", "Satisfying", "Saving", "Sawing", "Saying", "Scaring", "Scattering", "Scheduling", "Scolding",
		"Scorching", "Scraping", "Scratching", "Screaming", "Screwing", "Scribbling", "Scrubbing", "Sealing", "Searching", "Securing", "Seeking", "Seing", "Selecting", "Selling", "Sending", "Sensing", "Separating", "Servicing", "Serving", "Setting", "Settling", "Sewing", "Shading", "Shaking", "Shaping",
		"Sharing", "Shaving", "Shearing", "Sheltering", "Shining", "Shivering", "Shocking", "Shoing", "Shooting", "Shopping", "Showing", "Shrinking", "Shutting", "Sighing", "Signaling", "Signing", "Simplifying", "Sing", "Singing", "Sinking", "Sipping", "Sitting", "Sketching", "Skiing", "Skipping",
		"Slapping", "Slaying", "Sleeping", "Sliding", "Sling", "Slinking", "Slipping", "Slowing", "Smashing", "Smelling", "Smiling", "Smoking", "Snatching", "Sneaking", "Sneezing", "Sniffing", "Snoring", "Snowing", "Soaking", "Solving", "Soothing", "Soothsaying", "Sorting", "Sounding", "Sowing", "Sparing",
		"Sparking", "Sparkling", "Speaking", "Specifying", "Speeding", "Spelling", "Spending", "Spilling", "Spinning", "Spiting", "Splitting", "Spoiling", "Spotting", "Spraying", "Spreading", "Spring", "Sprouting", "Squashing", "Squeaking", "Squealing", "Squeezing", "Staining", "Stamping", "Standing",
		"Staring", "Starting", "Staying", "Stealing", "Steering", "Stepping", "Sticking", "Stimulating", "Sting", "Stinking", "Stirring", "Stitching", "Stoping", "Storing", "Strapping", "Streamlining", "Strengthening", "Stretching", "Striding", "Striking", "String", "Striping", "Striving", "Stroking",
		"Structuring", "Studying", "Stuffing", "Subtracting", "Succeeding", "Sucking", "Suffering", "Suggesting", "Suiting", "Summarizing", "Supervising", "Supplying", "Supporting", "Supposing", "Surprising", "Surrounding", "Suspecting", "Suspending", "Swearing", "Sweating", "Sweeping", "Swelling",
		"Swimming", "Swing", "Switching", "Symbolizing", "Synthesizing", "Systemizing", "Tabulating", "Taking", "Talking", "Taming", "Taping", "Targeting", "Tasting", "Teaching", "Tearing", "Teasing", "Telephoning", "Telling", "Tempting", "Terrifying", "Testing", "Thanking", "Thawing", "Thinking",
		"Thriving", "Throwing", "Thrusting", "Ticking", "Tickling", "Timing", "Tipping", "Tiring", "Touching", "Touring", "Towing", "Tracing", "Trading", "Training", "Transcribing", "Transferring", "Transforming", "Translating", "Transporting", "Trapping", "Traveling", "Treading", "Treating",
		"Trembling", "Tricking", "Tripping", "Trotting", "Troubleshooting", "Troubling", "Trusting", "Trying", "Tugging", "Tumbling", "Turning", "Tutoring", "Twisting", "Typing", "Undergoing", "Understanding", "Undertaking", "Undressing", "Unfastening", "Unifying", "Uniting", "Unlocking", "Unpacking",
		"Untidying", "Updating", "Upgrading", "Upholding", "Upsetting", "Using", "Utilizing", "Vanishing", "Verbalizing", "Verifying", "Vexing", "Visiting", "Wailing", "Waiting", "Waking", "Walking", "Wandering", "Wanting", "Warming", "Warning", "Washing", "Wasting", "Watching", "Watering", "Waving",
		"Wearing", "Weaving", "Wedding", "Weeping", "Weighing", "Welcoming", "Wending", "Wetting", "Whining", "Whipping", "Whirling", "Whispering", "Whistling", "Winding", "Wining", "Winking", "Wiping", "Wishing", "Withdrawing", "Withholding", "Withstanding", "Wobbling", "Wondering", "Working",
		"Worrying", "Wrapping", "Wrecking", "Wrestling", "Wriggling", "Writing", "Yawning", "Yelling"}

	musicNameAdjectives = []string{"Acceptable", "Alcoholic", "Apathetic", "Barbarous", "Bashful", "Bawdy", "Beautiful", "Befitting", "Belligerent", "Beneficial", "Bent", "Berserk", "Best", "Better", "Bewildered", "Big", "Billowy", "Bite-sized", "Bitter", "Bizarre", "Black", "Black-and-white", "Bloody", "Blue", "Blue-Brown",
		"Cheap", "Coherent", "Crabby", "Damaged", "Defiant", "Direful", "Dull", "Elegant", "Evanescent", "Evasive", "Even", "Excellent", "Excited", "Exciting", "Exclusive", "Exotic", "Expensive", "Extra-large", "Extra-small", "Exuberant", "Exultant", "Fabulous", "Faded", "Faint", "Fair", "Faithful", "Fallacious",
		"False", "Familiar", "Famous", "Fanatical", "Fancy", "Fantastic", "Far", "Far-Five", "Frail", "Gabby", "Good", "Grumpy", "Guarded", "Guiltless", "Gullible", "Gusty", "Guttural", "Habitual", "Half", "Hallowed", "Halting", "Handsome", "Handsomely", "Handy", "Hanging", "Hapless", "Happy", "Hard", "Hard-to-find", "Harmonious",
		"Harsh", "Hateful", "Heady", "Healthy", "Heartbreaking", "Heavenly", "Heavy", "Hellish", "Helpful", "Helpless", "Hesitant", "Hideous", "High", "High-Hurt", "Hushed", "Husky", "Hypnotic", "Hysterical", "Icky", "Icy", "Idiotic", "Ignorant", "Ill", "Ill-fated", "Ill-Infamous", "Jolly", "Lame", "Limping", "Literate", "Little",
		"Lively", "Living", "Lonely", "Long", "Long-Madly", "Measly", "Moaning", "Near", "Nonstop", "Obtainable", "Oceanic", "Odd", "Offbeat", "Old", "Old-Overt", "Perpetual", "Possessive", "Puffy", "Racial", "Remarkable", "Rough", "Scattered", "Scientific", "Scintillating", "Scrawny", "Screeching", "Second", "Second-Shut",
		"Smart", "Spiteful", "Sticky", "Super", "Tart", "Tasteful", "Tasteless", "Tasty", "Tawdry", "Tearful", "Tedious", "Teeny", "Teeny-Thoughtful", "Trite", "Undesirable", "Uppity", "Victorious", "Watery", "Weak", "Wealthy", "Weary", "Well-groomed", "Well-made", "Well-off", "Well-to-do", "Wet", "Whimsical", "Whispering",
		"White", "Whole", "Wholesale", "Wicked", "Wide", "Wide-Wretched", "Wrong", "Wry", "Xenophobic", "Yellow", "Yielding", "Young", "Youthful", "Yummy", "Zany", "Zealous", "Zesty", "Zippy", "Zonked"}

	musicNameNouns = []string{"Accounts", "Achievers", "Acoustics", "Actions", "Activities", "Actors", "Acts", "Additions", "Adjustments", "Advertisements", "Advices", "Aftermaths", "Afternoons", "Agreements", "Airplanes", "Airports", "Airs", "Alarms", "Amounts", "Amplifiers", "Amusements", "Angles", "Animals", "Answers", "Apparatus", "Apples", "Appliances", "Approvals", "Arguments",
		"Arithmetics", "Arms", "Arts", "Attacks", "Attempts", "Attentions", "Attractions", "Babys", "Backs", "Badges", "Bags", "Baits", "Balances", "Balloons", "Balls", "Bananas", "Bands", "Barrels", "Bars", "Baseballs", "Bases", "Basins", "Basketballs", "Baskets", "Baths", "Bats", "Battles", "Beads", "Beams", "Beans", "Bears", "Beasts", "Bedrooms", "Beds", "Beefs", "Bees", "Beetles",
		"Beggars", "Beginners", "Behaviors", "Beliefs", "Believes", "Bell", "Bikes", "Birds", "Birthdays", "Births", "Bites", "Bits", "Blades", "Blankets", "Bloods", "Blows", "Boards", "Boats", "Bodys", "Bombs", "Bones", "Books", "Boots", "Borders", "Bottles", "Boxes", "Boys", "Brains", "Brakes", "Branches", "Breads", "Breakfasts", "Breaths", "Bricks", "Bridges", "Brothers", "Brushes",
		"Bubbles", "Buckets", "Buildings", "Bulbs", "Bullets", "Buns", "Burns", "Bursts", "Bushes", "Buttons", "Cabbages", "Cables", "Cakes", "Calculators", "Calendars", "Cameras", "Camps", "Cannons", "Cans", "Canvass", "Caps", "Captions", "Cards", "Cares", "Carpenters", "Carriages", "Cars", "Cartoons", "Carts", "Casts", "Cats", "Causes", "Caves", "Cellars", "Cents", "Chains", "Chairs",
		"Chalks", "Chances", "Changes", "Channels", "Chickens", "Chins", "Churches", "Circles", "Citys", "Clams", "Clocks", "Cloths", "Clouds", "Clovers", "Clubs", "Coasts", "Coats", "Cobwebs", "Coils", "Collars", "Colors", "Combs", "Comforts", "Committees", "Comparisons", "Competitions", "Computers", "Conditions", "Connections", "Controls", "Cooks", "Copies", "Coppers", "Cords",
		"Corks", "Corns", "Coughs", "Covers", "Cows", "Crackers", "Cracks", "Crates", "Crayons", "Creams", "Creators", "Creatures", "Credits", "Cribs", "Crimes", "Crooks", "Crowds", "Crowns", "Crows", "Cups", "Currents", "Curtains", "Curves", "Cushions", "Cymbals", "Dads", "Daughters", "Days", "Deaths", "Debts", "Decisions", "Deers", "Degrees", "Designers", "Designs", "Desires", "Desks",
		"Destructions", "Details", "Developments", "Digestions", "Dimes", "Dinners", "Directions", "Dirts", "Discoveries", "Discussions", "Diseases", "Disgusts", "Disks", "Distances", "Distributions", "Divisions", "Docks", "Doctors", "Dogs", "Dolls", "Donkeys", "Doors", "Downtowns", "Dragons", "Drains", "Drawers", "Dreams", "Drinks", "Drives", "Drivings", "Drops", "Drugs",
		"Drums", "Ducks", "Dusts", "Ears", "Earthquakes", "Earths", "Edges", "Educations", "Effects", "Eggnogs", "Eggs", "Elbows", "Empires", "Ends", "Engines", "Equipments", "Errors", "Events", "Examples", "Exchanges", "Existences", "Expansions", "Experiences", "Experts", "Explosions", "Eyes", "Faces", "Facts", "Falls", "Fangs", "Fans", "Farmers", "Farms", "Fathers", "Faucets",
		"Fears", "Feasts", "Feathers", "Feelings", "Fictions", "Fields", "Fifths", "Fights", "Fingers", "Fires", "Flags", "Flames", "Flavors", "Flights", "Flocks", "Floors", "Flowers", "Fogs", "Folds", "Foods", "Foots", "Forces", "Forks", "Forms", "Frames", "Frictions", "Friends", "Frogs", "Fronts", "Fruits", "Fuels", "Furnitures", "Futures", "Galaxies", "Galleys", "Games", "Gardens",
		"Gates", "Giraffes", "Girls", "Gloves", "Glues", "Goats", "Golds", "Good-byes", "Gooses", "Governments", "Governors", "Grades", "Grains", "Grandfathers", "Grandmothers", "Grapes", "Grills", "Grips", "Grounds", "Groups", "Growths", "Guides", "Guitars", "Guns", "Haircuts", "Hairs", "Halls", "Hammers", "Handles", "Hands", "Harbors", "Harmonies", "Hates", "Hats", "Heads",
		"Healers", "Heals", "Healths", "Hearings", "Hearts", "Heats", "Heavens", "Helps", "Hens", "Hills", "Holes", "Holidays", "Homes", "Honeys", "Hooks", "Hopes", "Horns", "Horses", "Hoses", "Hospitals", "Hots", "Hours", "Houses", "Humors", "Hydrants", "Ices", "Icicles", "Ideas", "Impulses", "Incomes", "Increases", "Inks", "Inputs", "Insects", "Instruments", "Insurances", "Interests",
		"Inventions", "Invoices", "Irons", "Islands", "Jails", "Jams", "Jars", "Jeans", "Jellyfishs", "Jellys", "Jewels", "Joins", "Jokes", "Journeys", "Judges", "Juices", "Jumps", "Kettles", "Keyboards", "Keys", "Kicks", "Kingdoms", "Kites", "Kittens", "Knees", "Knots", "Knowledges", "Laborers", "Laces", "Lakes", "Lamps", "Lands", "Languages", "Laughs", "Lawyers", "Leads", "Leathers",
		"Legs", "Letters", "Lettuces", "Levels", "Lifts", "Lights", "Limits", "Linens", "Lines", "Lips", "Liquids", "Lists", "Lizards", "Lockets", "Locks", "Locusts", "Looks", "Losses", "Lovers", "Loves", "Lumber", "Lunches", "Machines", "Magic", "Maids", "Managers", "Maps", "Marbles", "Markets", "Marks", "Masks", "Matchs", "Meals", "Measures", "Meats", "Medicines", "Meetings", "Memories",
		"Men", "Metals", "Mice", "Milks", "Minds", "Mines", "Mints", "Minutes", "Mists", "Mittens", "Moms", "Money", "Monkeys", "Months", "Moons", "Mornings", "Mothers", "Motions", "Mountains", "Mouses", "Moves", "Muscles", "Music", "Nails", "Names", "Nations", "Necks", "Needles", "Needs", "Nerves", "Nests", "Nets", "Nightclubs", "Nights", "Ninjas", "Noises", "Noses", "Notebooks", "Notes",
		"Numbers", "Nuts", "Observations", "Oceans", "Offers", "Offices", "Oils", "Operations", "Opinions", "Oranges", "Orders", "Organizations", "Ornaments", "Outputs", "Ovens", "Owls", "Owners", "Pages", "Pails", "Pains", "Paints", "Pajamas", "Pancakes", "Pans", "Papers", "Parcels", "Parents", "Parks", "Partners", "Parts", "Passengers", "Pencils", "Pens", "Peoples", "Pests", "Phones",
		"Pickles", "Pictures", "Pies", "Pigs", "Pipes", "Piranhas", "Places", "Planes", "Planes", "Plants", "Plastics", "Plates", "Plays", "Pleasures", "Plots", "Plugs", "Pockets", "Points", "Poisons", "Police", "Positions", "Pots", "Powders", "Powers", "Prices", "Prints", "Prisons", "Processes", "Produces", "Profits", "Properties", "Protests", "Pulls", "Pumps", "Punishments",
		"Purposes", "Quarters", "Queens", "Questions", "Quicksands", "Quiets", "Quills", "Quilts", "Quinces", "Quivers", "Rabbits", "Rails", "Railways", "Rains", "Rainstorms", "Rakes", "Ranges", "Rates", "Rats", "Rays", "Reactions", "Reasons", "Receipts", "Records", "Regrets", "Relations", "Religions", "Representatives", "Requests", "Respects", "Rests", "Rewards", "Rhythms",
		"Rices", "Riddles", "Rifles", "Rings", "Rivers", "Roads", "Robins", "Robots", "Rocks", "Rooms", "Roots", "Ropes", "Roses", "Routes", "Rubs", "Rules", "Runs", "Sacks", "Sails", "Salts", "Samurais", "Sanctums", "Sands", "Scales", "Scenes", "Scents", "Schools", "Sciences", "Scissors", "Screens", "Screws", "Seas", "Seats", "Seeds", "Selections", "Senses", "Servants", "Shades", "Shakes",
		"Shames", "Shapes", "Sheep", "Sheets", "Ships", "Shirts", "Shocks", "Shoes", "Shops", "Shows", "Sides", "Sidewalks", "Signs", "Silks", "Silvers", "Sinks", "Sisters", "Sizes", "Skates", "Skins", "Skirts", "Skys", "Slaves", "Sleeps", "Sleets", "Slips", "Slopes", "Smashs", "Smells", "Smiles", "Smokes", "Snakes", "Societies", "Socks", "Sodas", "Sofas", "Soldiers", "Songs", "Sons", "Sorts",
		"Souls", "Sounds", "Soups", "Spaces", "Spades", "Sparks", "Speakers", "Spiders", "Sponges", "Spoons", "Spots", "Springs", "Spys", "Squares", "Squirrels", "Stages", "Stamps", "Stars", "Starts", "Statements", "States", "Stations", "Steams", "Steels", "Stems", "Steps", "Stews", "Sticks", "Stomachs", "Stones", "Stops", "Stores", "Storys", "Stoves", "Strangers", "Straws", "Streams",
		"Streets", "Strings", "Structures", "Substances", "Sugars", "Suggestions", "Suits", "Summers", "Suns", "Supports", "Surprises", "Sweaters", "Swims", "Swings", "Systems", "Tables", "Tails", "Talks", "Tanks", "Tastes", "Teachings", "Teams", "Tempers", "Tents", "Tests", "Textures", "Theories", "Things", "Thoughts", "Threads", "Thrills", "Throats", "Thrones", "Thumbs",
		"Thunders", "Tickets", "Tigers", "Times", "Tins", "Tips", "Titles", "Toads", "Toes", "Tomatoes", "Tongues", "Tops", "Touchs", "Towns", "Toys", "Trades", "Trails", "Trains", "Tramps", "Transports", "Trays", "Treatments", "Trees", "Trees", "Tricks", "Triggers", "Trips", "Troubles", "Trucks", "Tubs", "Turkeys", "Turns", "Twigs", "Twists", "Umbrellas", "Units", "Universes", "Vacations",
		"Valleys", "Values", "Vans", "Vases", "Vegetables", "Veils", "Veins", "Verses", "Vessels", "Vests", "Views", "Villages", "Visitors", "Voices", "Volcanos", "Volleyballs", "Voyages", "Walks", "Walls", "Wars", "Waves", "Ways", "Weeks", "Weights", "Wheels", "Whips", "Whistles", "Windows", "Winds", "Wines", "Wings", "Winters", "Wires", "Wishes", "Women", "Words", "Works", "Worlds", "Worms",
		"Worshippers", "Wounds", "Wrenches", "Wrens", "Wrists", "Writers", "Writings", "Yards", "Years", "Zebras", "Zippers"}

	musicNameEndings = []string{"In Paradise", "On The Floor", "In Throwing Distance", "Still Standing", "In The Jungle", "On My Mind", "On Fire", "In This Club", "In This World", "On The Go", "At The Party", "In The Mix", "Filtered Past",
		"In The Fire", "In The Future", "In His Heart", "Above The Darkness", "Under It All", "From My Mind", "It Worked", "In The Thick Of It", "In A Minute", "From Somewhere", "Far Too Bright", "From Distant Love",
		"Remember That Night", "All Our Money", "She Had It", "He Was Crazy", "Mama Always Said", "Having A Memory", "Treating You Well", "Knowing The Ropes", "In An Instant", "After Forever", "Under Attack",
		"Above The Clouds", "In The Story", "Under You", "On The Fence", "In Crazy Times", "In Case", "In My Day", "Gone Forever", "Under Heaven", "In Your Eyes", "Of Nothing", "Of Life", "Always Shining", "Watching That Thing",
		"Speaking The Truth", "Sitting Down", "Singing Loudly", "Over The Top", "Against All Odds", "Floating Away", "Plugging It In", "Breaking Dreams", "Jumping The Wall", "Developing Strangely", "Growing Big", "Taking A Bow",
		"Under Pressure", "Falling In Love", "Deciding To Go", "In Here", "Finding Out", "Packing It Up", "Changing The Channel", "Then He Ate It", "Trying Harder", "Rolling It Up", "In Here", "Making A Pass", "Taking A Break",
		"This Time Around", "Panic Setting In", "Bass Drop", "Of The Sun", "After We Ate", "Spinning In Circles", "Shaking That Thing", "Killing It", "Dancing Strangely", "Taking Pictures", "Feeling Weird", "Lighting Up",
		"Trying It", "Hitting The Field", "Quitting The Team", "On The Hunt", "In Her Heart", "Loving Your Neighbor", "Wasting Time", "Trickling Down", "After I Finish", "Fishing Around", "Loading The Program", "In Binary Language",
		"Growing Like A Weed", "Fighting For Freedom", "In A Room", "In Another Future", "Before Time Began", "Deciding Against It", "From Wrong", "Finding Where It Was", "When We Were Happy", "In The Middle", "At The Bottom",
		"Started Up", "Kicked In", "Running Out", "In A Flash", "It Started Again", "Starting Over", "Flashing Lights", "In Their Language", "Taking Over", "Yelling At Me", "Barking Orders", "Riding On A Wave", "Spinning In Circles",
		"Bringing It Up", "In Your House", "In Cold Blood", "Getting Better", "It Will Change", "Like It Happened", "Behind The Tree", "After We Eat", "Of Power", "Running Around", "Loving It", "Lost Everything", "Kill It Now",
		"In Life", "In Here", "Of Beauty", "Drawing The Line", "Driving The Road", "Lighting The Fire", "On The Beach", "By The Riverside", "In My Truck", "In My Belly", "In Class", "Turning Up", "In England", "Shooting Stars",
		"On Fire", "Putting It Out There", "After Dinner", "From Heaven", "Putting It On", "On Pause", "Delayed Reaction", "On The News", "Walking Around", "Writing It Down", "Ripping It Up", "Under It", "Coming Up Short",
		"Making It Happen", "Planting A Seed", "Helping Them Out", "Sitting Down", "Standing Up", "Breathing Heavy", "Asking Nicely", "Having A Hot Meal", "From The Past", "Asking Nicely", "At Work", "Below The Shelf", "Off The Roof",
		"Dancing All Night", "Blasting The Hillside", "Through The Window", "Leaving Town", "Breaking Down", "Making Noise", "Getting A Ride", "On My Brain", "Of My Dreams", "Loving It", "On My Watch", "From The Girl", "Running Off",
		"Starting Over", "Zooming In"}
)

// Music is a faker struct for Music
type Music struct {
	Faker *Faker
}

// Name returns a music name for Music
func (f Music) Name() string {
	var (
		adverb    = f.Faker.RandomStringElement(musicNameAdverbs)
		verb      = f.Faker.RandomStringElement(musicNameVerbs)
		adjective = f.Faker.RandomStringElement(musicNameAdjectives)
		noun      = f.Faker.RandomStringElement(musicNameNouns)
		ending    = f.Faker.RandomStringElement(musicNameEndings)
		name      = musicNameFormats[f.Faker.IntBetween(0, len(musicNameFormats)-1)]
	)

	// {{adverb}}
	name = strings.Replace(name, "{{adverb}}", adverb, 1)

	// {{verb}}
	name = strings.Replace(name, "{{verb}}", verb, 1)

	// {{adjective}}
	name = strings.Replace(name, "{{adjective}}", adjective, 1)

	// {{noun}}
	name = strings.Replace(name, "{{noun}}", noun, 1)

	// {{ending}}
	name = strings.Replace(name, "{{ending}}", ending, 1)

	return name
}

// Author returns the authors name for Music
func (f Music) Author() string {
	p := New().Person()
	return p.Name()
}

// Genre returns genre for Music
func (f Music) Genre() string {
	return f.Faker.RandomStringElement(musicGenres)
}

// Length returns how long the song takes for Music
func (f Music) Length() time.Duration {
	r := f.Faker.IntBetween(128, 512) * int(time.Second)
	return time.Duration(r)
}
