package gamedata

// ---- Items ----

type Item struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Rarity string `json:"rarity"` // poor | common | uncommon | rare | epic
	Desc   string `json:"desc"`
	Sell   int    `json:"sell"` // copper a vendor pays
	Stack  int    `json:"stack"`

	// equipment
	Slot    string         `json:"slot,omitempty"` // head, neck, shoulder, back, chest, wrist, hands, waist, legs, feet, finger, mainhand, offhand
	Type    string         `json:"type,omitempty"` // display text: Stoff, Leder, Platte, Schwert ...
	Armor   int            `json:"armor,omitempty"`
	Weapon  int            `json:"weapon,omitempty"` // weapon damage rating
	Stats   map[string]int `json:"stats,omitempty"`  // sta, str, agi, int
	Req     int            `json:"req,omitempty"`    // required character level
	Classes []string       `json:"classes,omitempty"`
	Buy     int            `json:"buy,omitempty"`   // copper a vendor charges
	Heal    int            `json:"heal,omitempty"`  // consumable: hit points restored
	Quest   bool           `json:"quest,omitempty"` // quest item: cannot be sold
}

var Items = map[string]*Item{
	// wolves
	"wolf_pelt": {ID: "wolf_pelt", Name: "Wolfspelz", Icon: "🐺", Rarity: "common", Desc: "Ein zotteliges graues Fell.", Sell: 12, Stack: 20},
	"wolf_fang": {ID: "wolf_fang", Name: "Wolfszahn", Icon: "🦷", Rarity: "common", Desc: "Scharf und gelblich.", Sell: 8, Stack: 20},
	"wolf_meat": {ID: "wolf_meat", Name: "Zähes Wolffleisch", Icon: "🍖", Rarity: "poor", Desc: "Nicht gerade ein Festmahl.", Sell: 3, Stack: 20},
	// boars
	"boar_tusk": {ID: "boar_tusk", Name: "Eberhauer", Icon: "🐗", Rarity: "common", Desc: "Ein gebogener, gelblicher Hauer.", Sell: 10, Stack: 20},
	"boar_meat": {ID: "boar_meat", Name: "Wildschweinfleisch", Icon: "🥩", Rarity: "common", Desc: "Saftiges, frisches Fleisch.", Sell: 6, Stack: 20},
	// spiders
	"spider_silk": {ID: "spider_silk", Name: "Spinnenseide", Icon: "🕸️", Rarity: "common", Desc: "Zäh, klebrig und überraschend reißfest.", Sell: 14, Stack: 20},
	"spider_eye":  {ID: "spider_eye", Name: "Spinnenauge", Icon: "👁️", Rarity: "poor", Desc: "Starrt dich vorwurfsvoll an.", Sell: 4, Stack: 20},
	"venom_sac":   {ID: "venom_sac", Name: "Giftbeutel", Icon: "🧪", Rarity: "common", Desc: "Vorsicht, tropft.", Sell: 16, Stack: 20},
	// bandits
	"rusty_dagger": {ID: "rusty_dagger", Name: "Rostiger Dolch", Icon: "🗡️", Rarity: "poor", Desc: "Gerade noch scharf genug für Brot.", Sell: 10, Stack: 5},
	"bandit_mask":  {ID: "bandit_mask", Name: "Banditentuch", Icon: "🧣", Rarity: "common", Desc: "Ein schwarzes Tuch zum Verhüllen des Gesichts.", Sell: 15, Stack: 20},
	"stolen_purse": {ID: "stolen_purse", Name: "Gestohlener Beutel", Icon: "👛", Rarity: "common", Desc: "Gehört irgendeinem Bauern aus der Gegend.", Sell: 35, Stack: 20},
	// bears
	"bear_claw": {ID: "bear_claw", Name: "Bärenklaue", Icon: "🐾", Rarity: "common", Desc: "Riesig und messerscharf.", Sell: 30, Stack: 20},
	"bear_pelt": {ID: "bear_pelt", Name: "Bärenfell", Icon: "🐻", Rarity: "uncommon", Desc: "Dick, warm und sehr schwer.", Sell: 60, Stack: 10},
	// quest rewards
	"ranger_cloak":  {ID: "ranger_cloak", Name: "Mantel des Waldhüters", Icon: "🧥", Rarity: "uncommon", Desc: "Ein Geschenk von Förster Aldric.", Sell: 60, Stack: 1},
	"hilda_pie":     {ID: "hilda_pie", Name: "Hildas Apfelkuchen", Icon: "🥧", Rarity: "common", Desc: "Noch warm. Schmeckt nach Zuhause.", Sell: 20, Stack: 5},
	"scholar_charm": {ID: "scholar_charm", Name: "Amulett der Naturforscherin", Icon: "📿", Rarity: "uncommon", Desc: "Elara schwört, es bringe Glück bei der Feldarbeit.", Sell: 90, Stack: 1},
	"guard_badge":   {ID: "guard_badge", Name: "Abzeichen der Wache", Icon: "🎖️", Rarity: "uncommon", Desc: "Ehrenabzeichen für mutige Dorfbewohner.", Sell: 120, Stack: 1},
	"hunter_ring":   {ID: "hunter_ring", Name: "Ring des Bärenjägers", Icon: "💍", Rarity: "uncommon", Desc: "Aus der Klaue eines Bären gefertigt.", Sell: 180, Stack: 1},
}

// ---- Mobs ----

type LootEntry struct {
	Item   string
	Chance float64 // 0..1
	Min    int
	Max    int
	Quest  string // only drops while this quest is active
}

type Mob struct {
	ID         string
	Name       string
	Model      string
	LevelMin   int
	LevelMax   int
	HPBase     int
	HPPerLevel int
	DmgMin     int
	DmgMax     int
	AttackSec  float64
	Speed      float64 // run speed (m/s)
	Aggressive bool    // attacks players on sight
	AggroRange float64
	Leash      float64
	XPBase     int
	CopperMin  int
	CopperMax  int
	Loot       []LootEntry
	GearChance float64 // chance to drop a random piece of equipment
	RespawnSec float64

	Elite     bool
	Boss      bool    // dungeon boss: personal epic loot, long corpse time
	HPMult    float64 // multiplies the level based hit points (new mobs)
	DmgMult   float64
	Range     float64 // > 0: ranged attacker
	Style     string  // ranged effect style (bolt, frost, ...)
	Scale     float64 // model scale
	Tameable  bool
	Abilities []string // boss abilities, see game/boss.go
	Faction   string
	Set       string // dungeon gear set a boss drops from (bf, ns, gk)
	Final     bool   // last boss of a dungeon: two set items
}

var Mobs = map[string]*Mob{
	"wolf": {
		ID: "wolf", Name: "Waldwolf", Model: "wolf", LevelMin: 1, LevelMax: 2,
		HPBase: 45, HPPerLevel: 20, DmgMin: 3, DmgMax: 5, AttackSec: 2.0, Speed: 5.2,
		Leash: 45, XPBase: 20, CopperMin: 1, CopperMax: 6, RespawnSec: 40,
		Loot: []LootEntry{{Item: "wolf_pelt", Chance: 0.65, Min: 1, Max: 1}, {Item: "wolf_fang", Chance: 0.40, Min: 1, Max: 2}, {Item: "wolf_meat", Chance: 0.45, Min: 1, Max: 1}},
	},
	"boar": {
		ID: "boar", Name: "Wildschwein", Model: "boar", LevelMin: 2, LevelMax: 3,
		HPBase: 70, HPPerLevel: 25, DmgMin: 4, DmgMax: 7, AttackSec: 2.2, Speed: 5.0,
		Leash: 45, XPBase: 25, CopperMin: 2, CopperMax: 8, RespawnSec: 40,
		Loot: []LootEntry{{Item: "boar_tusk", Chance: 0.45, Min: 1, Max: 1}, {Item: "boar_meat", Chance: 0.70, Min: 1, Max: 2}},
	},
	"spider": {
		ID: "spider", Name: "Riesenspinne", Model: "spider", LevelMin: 3, LevelMax: 4,
		HPBase: 80, HPPerLevel: 28, DmgMin: 5, DmgMax: 8, AttackSec: 1.8, Speed: 5.6,
		Aggressive: true, AggroRange: 9, Leash: 40, XPBase: 30, CopperMin: 3, CopperMax: 10, RespawnSec: 45,
		Loot: []LootEntry{{Item: "spider_silk", Chance: 0.65, Min: 1, Max: 2}, {Item: "spider_eye", Chance: 0.40, Min: 1, Max: 1}, {Item: "venom_sac", Chance: 0.25, Min: 1, Max: 1}},
	},
	"bandit": {
		ID: "bandit", Name: "Straßenräuber", Model: "bandit", LevelMin: 4, LevelMax: 5,
		HPBase: 110, HPPerLevel: 35, DmgMin: 7, DmgMax: 11, AttackSec: 2.0, Speed: 5.2,
		Aggressive: true, AggroRange: 8, Leash: 40, XPBase: 38, CopperMin: 8, CopperMax: 25, RespawnSec: 50,
		Loot: []LootEntry{{Item: "rusty_dagger", Chance: 0.45, Min: 1, Max: 1}, {Item: "bandit_mask", Chance: 0.35, Min: 1, Max: 1}, {Item: "stolen_purse", Chance: 0.40, Min: 1, Max: 1}},
	},
	"bear": {
		ID: "bear", Name: "Grizzlybär", Model: "bear", LevelMin: 5, LevelMax: 6,
		HPBase: 190, HPPerLevel: 45, DmgMin: 10, DmgMax: 15, AttackSec: 2.6, Speed: 5.0,
		Leash: 50, XPBase: 45, CopperMin: 6, CopperMax: 18, RespawnSec: 60,
		Loot: []LootEntry{{Item: "bear_claw", Chance: 0.60, Min: 1, Max: 2}, {Item: "bear_pelt", Chance: 0.35, Min: 1, Max: 1}},
	},
}

// ---- NPCs ----

type NPC struct {
	ID      string
	Name    string
	Title   string
	Model   string
	X, Z    float64
	R       float64
	Quests  []string
	Trainer string   // class id this NPC trains ("" = no trainer)
	Vendor  bool     // buys loot for copper
	Stock   []string // item ids the vendor sells
	StockK  string   // weapons | armor | jewelry | potions: generated goods
	StockLo int
	StockHi int
	Flight  bool     // griffin master
	Ship    bool     // ferry captain (sea route between the harbors)
	TurnIn  []string // quests this NPC only accepts hand-ins for
	Greet   string
}

var NPCs = []*NPC{
	{ID: "forester", Name: "Förster Aldric", Title: "Waldhüter", Model: "forester", X: 6, Z: 9, R: 3.6,
		Quests: []string{"wolf_plague", "wolf_pelts"}, Greet: "Der Wald ist friedlich heute… meistens."},
	{ID: "brandt", Name: "Hauptmann Brandt", Title: "Kriegerlehrer", Model: "trainer_warrior", X: -9, Z: 14, R: 1.2,
		Trainer: "warrior", Greet: "Stärke kommt nicht von allein. Zeig, was du kannst!"},
	{ID: "lyra", Name: "Lyra Falkenauge", Title: "Jägerlehrerin", Model: "trainer_hunter", X: -13, Z: 18, R: 1.0,
		Trainer: "hunter", Greet: "Ruhig atmen, Ziel erfassen, loslassen."},
	{ID: "orlan", Name: "Magister Orlan", Title: "Magierlehrer", Model: "trainer_mage", X: -5, Z: 19, R: 0.8,
		Trainer: "mage", Greet: "Die arkanen Künste verlangen Disziplin."},
	{ID: "marta", Name: "Händlerin Marta", Title: "Händlerin", Model: "merchant", X: 11, Z: 15, R: 4.0,
		Vendor: true, StockK: "jewelry", StockLo: 0, StockHi: 1, Stock: []string{"potion_hp", "potion_hp2"}, Greet: "Was du im Wald findest, kaufe ich dir gern ab!"},
	{ID: "garrick", Name: "Schmied Garrick", Title: "Waffen & Rüstungen", Model: "smith", X: 15, Z: 11, R: 4.2,
		Vendor: true, StockK: "weapons,armor", StockLo: 0, StockHi: 1, Greet: "Gutes Eisen, ehrliche Preise."},
	{ID: "hilda", Name: "Bäuerin Hilda", Title: "Bäuerin", Model: "farmer", X: -26, Z: 8, R: 4.4,
		Quests: []string{"boar_trouble", "boar_feast"}, Greet: "Die Ernte ist dieses Jahr gut. Wenn nur die Wildschweine nicht wären."},
	{ID: "elara", Name: "Elara Moosfeld", Title: "Naturforscherin", Model: "scholar", X: 8, Z: -22, R: 3.4,
		Quests: []string{"spider_silk", "spider_nest"}, Greet: "Faszinierend, diese Waldbewohner. Aus sicherer Entfernung."},
	{ID: "gerd", Name: "Wachposten Gerd", Title: "Dorfwache", Model: "guard", X: 14, Z: 28, R: 0.4,
		Quests: []string{"bandit_camp", "stolen_goods"}, Greet: "Halt die Augen offen. Die Straßen sind nicht mehr sicher."},
	{ID: "torben", Name: "Torben Eisenbart", Title: "Bärenjäger", Model: "huntsman", X: 56, Z: 18, R: 2.0,
		Quests: []string{"bear_hunt", "bear_claws"}, Greet: "Die Bären hier sind groß. Sehr groß."},
}

// ---- Quests ----

type Objective struct {
	Kind   string `json:"kind"` // "kill" | "collect"
	Target string `json:"target"`
	Count  int    `json:"count"`
	Text   string `json:"text"`
}

type Reward struct {
	XP     int      `json:"xp"`
	Copper int      `json:"copper"`
	Items  []string `json:"items"`
	Gear   bool     `json:"gear,omitempty"`   // a random piece of gear for your class
	Skills []string `json:"skills,omitempty"` // skills learned
}

type Quest struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	MinLevel   int         `json:"minLevel"`
	Requires   string      `json:"requires,omitempty"` // quest that must be handed in first
	Text       string      `json:"text"`
	Progress   string      `json:"progress"`
	Complete   string      `json:"complete"`
	Objectives []Objective `json:"objectives"`
	Reward     Reward      `json:"reward"`

	Level      int      `json:"level"`
	Zone       string   `json:"zone,omitempty"`
	Tag        string   `json:"tag,omitempty"`       // Gruppe | Elite | Zeit | Klasse | Dungeon
	TimeLimit  int      `json:"timeLimit,omitempty"` // seconds
	GiveItem   string   `json:"giveItem,omitempty"`  // quest item handed over on accept
	Class      string   `json:"class,omitempty"`
	GrantSkill []string `json:"grantSkill,omitempty"` // skills granted on accept
	Giver      string   `json:"giver,omitempty"`
	Hand       string   `json:"hand,omitempty"` // NPC to hand in to (empty = giver)

	// map areas (filled in by Init): where the objectives can be completed and where to hand in
	Areas  [][]QArea `json:"areas,omitempty"`
	HandIn *QArea    `json:"handIn,omitempty"`
}

var Quests = map[string]*Quest{
	"wolf_plague": {
		ID: "wolf_plague", Title: "Die Wolfsplage", MinLevel: 1,
		Text: "Gut, dass du vorbeikommst, %n! Die Wölfe im Wald sind in letzter Zeit viel zu zahlreich geworden. " +
			"Sie reißen mein Wild und wagen sich immer näher an mein Haus heran.\n\n" +
			"Geh ein Stück nach Osten, dort streunen sie herum. Sie greifen nicht von selbst an, aber wenn du sie reizt, wehren sie sich. " +
			"Erlege sechs von ihnen, dann kehrt hier wieder Ruhe ein.",
		Progress:   "Die Wölfe sind noch immer da draußen. Hast du schon sechs erlegt?",
		Complete:   "Sechs Wölfe weniger – der Wald dankt dir, %n! Nimm diesen Mantel als Zeichen meiner Dankbarkeit.",
		Objectives: []Objective{{Kind: "kill", Target: "wolf", Count: 6, Text: "Waldwölfe getötet"}},
		Reward:     Reward{XP: 150, Copper: 75, Items: []string{"ranger_cloak"}},
	},
	"wolf_pelts": {
		ID: "wolf_pelts", Title: "Pelze für den Winter", MinLevel: 2, Requires: "wolf_plague",
		Text: "Du hast dich als tüchtig erwiesen, %n. Der Winter kommt früh, und meine Decken sind löchrig.\n\n" +
			"Bring mir vier gute Wolfspelze, dann nähe ich mir etwas Warmes daraus. Die Wölfe im Osten tragen genug davon.",
		Progress:   "Vier Wolfspelze, mehr brauche ich nicht. Die Wölfe im Osten haben sie.",
		Complete:   "Prächtige Felle! Damit kommen wir warm durch den Winter. Hier, für deine Mühe.",
		Objectives: []Objective{{Kind: "collect", Target: "wolf_pelt", Count: 4, Text: "Wolfspelze gesammelt"}},
		Reward:     Reward{XP: 130, Copper: 60},
	},
	"boar_trouble": {
		ID: "boar_trouble", Title: "Schweinerei auf dem Feld", MinLevel: 2,
		Text: "Ach, %n, gut dass du kommst! Die Wildschweine aus dem Westwald wühlen mein ganzes Feld um.\n\n" +
			"Sie sind nicht angriffslustig, aber mit ihren Hauern ist nicht zu spaßen. Erlege sechs von ihnen, dann habe ich endlich wieder Ruhe.",
		Progress:   "Die Schweine wühlen immer noch! Hast du sechs erwischt?",
		Complete:   "Endlich Ruhe auf dem Feld! Nimm einen von meinen Kuchen mit – frisch gebacken.",
		Objectives: []Objective{{Kind: "kill", Target: "boar", Count: 6, Text: "Wildschweine getötet"}},
		Reward:     Reward{XP: 220, Copper: 90, Items: []string{"hilda_pie"}},
	},
	"boar_feast": {
		ID: "boar_feast", Title: "Ein Festmahl", MinLevel: 2, Requires: "boar_trouble",
		Text: "Wo wir schon dabei sind: Am Sonntag feiert das Dorf Erntedank und ich habe nichts Rechtes zu braten.\n\n" +
			"Bring mir vier Stücke Wildschweinfleisch, %n, und du bist eingeladen.",
		Progress:   "Vier Stücke Wildschweinfleisch brauche ich noch.",
		Complete:   "Das wird ein Festschmaus! Hier ist dein Anteil.",
		Objectives: []Objective{{Kind: "collect", Target: "boar_meat", Count: 4, Text: "Wildschweinfleisch gesammelt"}},
		Reward:     Reward{XP: 200, Copper: 80},
	},
	"spider_silk": {
		ID: "spider_silk", Title: "Seide aus dem Norden", MinLevel: 3,
		Text: "Ah, ein Abenteurer! Ich erforsche die Riesenspinnen im Norden – leider wage ich mich nicht nah genug heran.\n\n" +
			"Bring mir fünf Büschel Spinnenseide. Aber Vorsicht, %n: Diese Spinnen greifen jeden an, der ihnen zu nahe kommt.",
		Progress:   "Fünf Büschel Spinnenseide, bitte. Die Spinnen leben weiter im Norden.",
		Complete:   "Hervorragend! Die Fasern sind ganz außergewöhnlich. Hier, deine Belohnung.",
		Objectives: []Objective{{Kind: "collect", Target: "spider_silk", Count: 5, Text: "Spinnenseide gesammelt"}},
		Reward:     Reward{XP: 260, Copper: 110},
	},
	"spider_nest": {
		ID: "spider_nest", Title: "Das Spinnennest", MinLevel: 3, Requires: "spider_silk",
		Text: "Die Spinnen werden immer mehr, %n. Wenn wir nichts unternehmen, erreichen sie bald das Dorf.\n\n" +
			"Erlege sechs von ihnen und dezimiere damit das Nest.",
		Progress:   "Das Nest ist noch nicht dezimiert. Erlege sechs Riesenspinnen.",
		Complete:   "Das Nest ist erst mal ruhig. Nimm dieses Amulett, es hat mir immer Glück gebracht.",
		Objectives: []Objective{{Kind: "kill", Target: "spider", Count: 6, Text: "Riesenspinnen getötet"}},
		Reward:     Reward{XP: 380, Copper: 140, Items: []string{"scholar_charm"}},
	},
	"bandit_camp": {
		ID: "bandit_camp", Title: "Das Banditenlager", MinLevel: 4,
		Text: "Straßenräuber haben südlich des Dorfes ein Lager aufgeschlagen und überfallen unsere Händler.\n\n" +
			"Sie greifen sofort an, wenn sie dich sehen, also sei vorbereitet, %n. Töte sechs von ihnen.",
		Progress:   "Die Straßen sind noch nicht sicher. Sechs Banditen, vergiss das nicht.",
		Complete:   "Gute Arbeit! Die Händler können wieder aufatmen. Dieses Abzeichen steht dir zu.",
		Objectives: []Objective{{Kind: "kill", Target: "bandit", Count: 6, Text: "Straßenräuber getötet"}},
		Reward:     Reward{XP: 520, Copper: 220, Items: []string{"guard_badge"}},
	},
	"stolen_goods": {
		ID: "stolen_goods", Title: "Gestohlenes Gut", MinLevel: 4, Requires: "bandit_camp",
		Text:       "Die Banditen tragen die Beutel ihrer Opfer bei sich. Bring mir drei Gestohlene Beutel, %n, dann gebe ich sie den Besitzern zurück.",
		Progress:   "Drei Gestohlene Beutel fehlen noch. Die Banditen haben sie.",
		Complete:   "Die Besitzer werden sich freuen! Hier ist ein Finderlohn.",
		Objectives: []Objective{{Kind: "collect", Target: "stolen_purse", Count: 3, Text: "Gestohlene Beutel gesammelt"}},
		Reward:     Reward{XP: 450, Copper: 200},
	},
	"bear_hunt": {
		ID: "bear_hunt", Title: "Bärenjagd", MinLevel: 5,
		Text: "Ich bin Torben, und ich jage Bären seit dreißig Jahren. Dieses Jahr sind es zu viele.\n\n" +
			"Sie greifen nicht von sich aus an – aber wehe, du reizt sie. Töte vier Grizzlys, %n, wenn du dich traust.",
		Progress:   "Vier Bären, mein Freund. Sie streifen weiter östlich umher.",
		Complete:   "Vier Bären! Du bist ein echter Jäger.",
		Objectives: []Objective{{Kind: "kill", Target: "bear", Count: 4, Text: "Grizzlybären getötet"}},
		Reward:     Reward{XP: 700, Copper: 300},
	},
	"bear_claws": {
		ID: "bear_claws", Title: "Klauen als Trophäe", MinLevel: 5, Requires: "bear_hunt",
		Text:       "Aus Bärenklauen fertige ich Schmuck. Bring mir vier davon und ich mache dir einen Ring, %n.",
		Progress:   "Vier Bärenklauen, dann bekommst du deinen Ring.",
		Complete:   "Hervorragende Klauen! Hier, dein Ring – trag ihn mit Stolz.",
		Objectives: []Objective{{Kind: "collect", Target: "bear_claw", Count: 4, Text: "Bärenklauen gesammelt"}},
		Reward:     Reward{XP: 650, Copper: 350, Items: []string{"hunter_ring"}},
	},
}

// ---- Spawns ----

type Spawn struct {
	Mob string  `json:"mob"`
	X   float64 `json:"x"`
	Z   float64 `json:"z"`
}

var Spawns = []Spawn{
	// wolves (east)
	{"wolf", 30, 2}, {"wolf", 36, 8}, {"wolf", 33, -6}, {"wolf", 42, 0}, {"wolf", 40, 12},
	{"wolf", 46, -8}, {"wolf", 28, 16}, {"wolf", 38, -14}, {"wolf", 50, 6}, {"wolf", 34, 20},
	// boars (west)
	{"boar", -34, 4}, {"boar", -40, 10}, {"boar", -36, 18}, {"boar", -44, 2}, {"boar", -32, 24}, {"boar", -42, 20}, {"boar", -38, -6},
	// spiders (north)
	{"spider", 10, -40}, {"spider", 18, -46}, {"spider", 4, -48}, {"spider", -6, -42}, {"spider", 14, -54}, {"spider", -2, -56},
	// bandits (south camp)
	{"bandit", 10, 42}, {"bandit", 18, 48}, {"bandit", 6, 50}, {"bandit", 22, 40}, {"bandit", 14, 56}, {"bandit", 2, 44},
	// bears (far east)
	{"bear", 60, 30}, {"bear", 66, 40}, {"bear", 58, 48}, {"bear", 70, 22}, {"bear", 72, 36},
}

// SpawnPoint is where new characters start (forest, next to the house).
var SpawnPoint = struct{ X, Z, R float64 }{X: 1, Z: 10, R: 3.14159}
