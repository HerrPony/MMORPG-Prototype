// Package gamedata contains all static, data-driven game definitions
// (races, classes, skills, items, mobs, NPCs, quests, world layout).
package gamedata

type Race struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

var Races = map[string]*Race{
	"human": {ID: "human", Name: "Mensch", Desc: "Vielseitig und ehrgeizig. Menschen finden sich in jeder Lage zurecht."},
}

var RaceOrder = []string{"human"}

// Skill kinds.
const (
	KindDamage = "damage" // direct damage on target
	KindDot    = "dot"    // direct damage + damage over time
	KindAoE    = "aoe"    // damage to all enemies around the caster
	KindBuff   = "buff"   // self damage buff
	KindShield = "shield" // self absorb shield
)

type Skill struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Icon     string  `json:"icon"`
	Desc     string  `json:"desc"`
	Kind     string  `json:"kind"`
	Basic    bool    `json:"basic"`    // auto-attack skill
	Cost     int     `json:"cost"`     // resource cost
	Gain     int     `json:"gain"`     // resource gained on hit (rage)
	Cooldown float64 `json:"cooldown"` // seconds
	Range    float64 `json:"range"`
	Cast     float64 `json:"cast"` // cast time in seconds
	Power    float64 `json:"power"`
	DotPower float64 `json:"dotPower"` // damage per second tick
	Duration float64 `json:"duration"`
	Slow     float64 `json:"slow"`            // 0..1 movement slow on target
	Radius   float64 `json:"radius"`          // for AoE
	Color    string  `json:"color"`           // projectile/effect color
	Style    string  `json:"style"`           // visual: melee, arrow, bolt, fire, frost, nova, buff
	Level    int     `json:"level"`           // character level needed to learn it (1 = known from the start)
	Price    int     `json:"price"`           // copper the class trainer charges
	Quest    bool    `json:"quest,omitempty"` // learned through a class quest, not at a trainer

	Stance   string  `json:"stance,omitempty"`  // stance skills: battle | defensive
	Stun     float64 `json:"stun,omitempty"`    // seconds the target cannot act
	Execute  float64 `json:"execute,omitempty"` // hp fraction below which ExecMult applies
	ExecMult float64 `json:"execMult,omitempty"`
	Blink    float64 `json:"blink,omitempty"` // teleport distance forward
	MinRange float64 `json:"minRange,omitempty"`
	Taken    float64 `json:"taken,omitempty"` // extra damage taken while a buff is active
}

// Additional skill kinds.
const (
	KindStance    = "stance"
	KindTaunt     = "taunt"
	KindCharge    = "charge"
	KindLastStand = "laststand"
	KindTame      = "tame"
	KindCallPet   = "callpet"
	KindMendPet   = "mendpet"
	KindAoETarget = "aoetarget" // damage around the target
	KindElemental = "elemental"
	KindBlink     = "blink"
	KindPoly      = "polymorph"
	KindGround    = "ground" // damaging area at the target position for a while
)

// training lists {required level, price in copper} for every skill that must be learned.
// A new skill becomes available every 2 levels.
var training = map[string][2]int{
	"w_heroic": {2, 10}, "w_rend": {4, 50}, "w_whirl": {6, 150}, "w_shout": {8, 400},
	"h_aimed": {2, 10}, "h_serpent": {4, 50}, "h_cripple": {6, 150}, "h_eagle": {8, 400},
	"m_fireball": {2, 10}, "m_frost": {4, 50}, "m_nova": {6, 150}, "m_shield": {8, 400},
}

func init() {
	for id, s := range Skills {
		s.Level = 1
		if t, ok := training[id]; ok {
			s.Level, s.Price = t[0], t[1]
		}
	}
}

// BasicSkill returns the free auto-attack skill of a class.
func BasicSkill(c *Class) string {
	for _, id := range c.Skills {
		if Skills[id].Basic {
			return id
		}
	}
	return ""
}

type Class struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Desc       string   `json:"desc"`
	Role       string   `json:"role"`
	Resource   string   `json:"resource"` // rage | focus | mana
	ResColor   string   `json:"resColor"`
	ResMax     int      `json:"resMax"`
	ResStart   int      `json:"resStart"`
	ResRegen   float64  `json:"resRegen"` // per second (out of combat decay for rage is handled in game)
	HP         int      `json:"hp"`
	HPPerLevel int      `json:"hpPerLevel"`
	Color      string   `json:"color"`
	Skills     []string `json:"skills"`
}

var Classes = map[string]*Class{
	"warrior": {
		ID: "warrior", Name: "Krieger", Role: "Nahkampf · Tank",
		Desc:     "Unerschütterlicher Frontkämpfer. Baut im Kampf Wut auf und setzt sie für mächtige Schläge ein.",
		Resource: "rage", ResColor: "#c83030", ResMax: 100, ResStart: 0, ResRegen: 0,
		HP: 130, HPPerLevel: 20, Color: "#c79c6e",
		Skills: []string{"w_strike", "w_heroic", "w_rend", "w_whirl", "w_shout"},
	},
	"hunter": {
		ID: "hunter", Name: "Jäger", Role: "Fernkampf · Schaden",
		Desc:     "Meisterschütze der Wildnis. Bekämpft Feinde aus sicherer Distanz mit Pfeilen und Giften.",
		Resource: "focus", ResColor: "#e8a030", ResMax: 100, ResStart: 100, ResRegen: 5,
		HP: 105, HPPerLevel: 15, Color: "#abd473",
		Skills: []string{"h_shot", "h_aimed", "h_serpent", "h_cripple", "h_eagle"},
	},
	"mage": {
		ID: "mage", Name: "Magier", Role: "Fernkampf · Zauberer",
		Desc:     "Gelehrter der arkanen Künste. Entfesselt Feuer und Frost, schützt sich mit magischen Schilden.",
		Resource: "mana", ResColor: "#3070e0", ResMax: 120, ResStart: 120, ResRegen: 4,
		HP: 85, HPPerLevel: 11, Color: "#69ccf0",
		Skills: []string{"m_bolt", "m_fireball", "m_frost", "m_nova", "m_shield"},
	},
}

var ClassOrder = []string{"warrior", "hunter", "mage"}

var Skills = map[string]*Skill{
	// ---- Krieger ----
	"w_strike": {ID: "w_strike", Name: "Hieb", Icon: "⚔️", Kind: KindDamage, Basic: true, Gain: 12, Cooldown: 1.6, Range: 3.5, Power: 7, Style: "melee", Color: "#ffffff",
		Desc: "Ein schneller Nahkampfhieb. Erzeugt Wut. Wird automatisch wiederholt."},
	"w_heroic": {ID: "w_heroic", Name: "Heldenhafter Stoß", Icon: "🗡️", Kind: KindDamage, Cost: 15, Gain: 0, Cooldown: 4, Range: 3.5, Power: 22, Style: "melee", Color: "#ffcc66",
		Desc: "Ein kraftvoller Stoß, der hohen Schaden verursacht."},
	"w_rend": {ID: "w_rend", Name: "Zerfleischen", Icon: "🩸", Kind: KindDot, Cost: 20, Cooldown: 8, Range: 3.5, Power: 6, DotPower: 4, Duration: 8, Style: "melee", Color: "#d02020",
		Desc: "Verursacht Schaden und lässt das Ziel 8 Sek. lang bluten."},
	"w_whirl": {ID: "w_whirl", Name: "Wirbelwind", Icon: "🌀", Kind: KindAoE, Cost: 25, Cooldown: 10, Radius: 5, Power: 14, Style: "nova", Color: "#dddddd",
		Desc: "Du wirbelst herum und triffst alle Feinde in der Nähe."},
	"w_shout": {ID: "w_shout", Name: "Kriegsschrei", Icon: "📣", Kind: KindBuff, Cost: 10, Cooldown: 30, Power: 0.3, Duration: 20, Style: "buff", Color: "#ff8844",
		Desc: "Erhöht deinen verursachten Schaden 20 Sek. lang um 30%."},

	// ---- Jäger ----
	"h_shot": {ID: "h_shot", Name: "Schuss", Icon: "🏹", Kind: KindDamage, Basic: true, Cooldown: 1.8, Range: 30, Power: 7, Style: "arrow", Color: "#d8c090",
		Desc: "Ein einfacher Schuss mit dem Bogen. Wird automatisch wiederholt."},
	"h_aimed": {ID: "h_aimed", Name: "Zielschuss", Icon: "🎯", Kind: KindDamage, Cost: 25, Cooldown: 5, Range: 35, Power: 24, Style: "arrow", Color: "#ffee88",
		Desc: "Ein sorgfältig gezielter Schuss, der hohen Schaden verursacht."},
	"h_serpent": {ID: "h_serpent", Name: "Schlangenbiss", Icon: "🐍", Kind: KindDot, Cost: 15, Cooldown: 8, Range: 30, Power: 3, DotPower: 3.5, Duration: 8, Style: "arrow", Color: "#60d040",
		Desc: "Ein vergifteter Pfeil, der das Ziel 8 Sek. lang vergiftet."},
	"h_cripple": {ID: "h_cripple", Name: "Lähmender Schuss", Icon: "🕸️", Kind: KindDamage, Cost: 15, Cooldown: 10, Range: 30, Power: 8, Slow: 0.5, Duration: 8, Style: "arrow", Color: "#9999ff",
		Desc: "Trifft das Ziel und verlangsamt es 8 Sek. lang um 50%."},
	"h_eagle": {ID: "h_eagle", Name: "Adlerauge", Icon: "🦅", Kind: KindBuff, Cost: 10, Cooldown: 30, Power: 0.3, Duration: 20, Style: "buff", Color: "#ffdd55",
		Desc: "Schärft deine Sinne: 30% mehr Schaden für 20 Sek."},

	// ---- Magier ----
	"m_bolt": {ID: "m_bolt", Name: "Arkanes Geschoss", Icon: "✨", Kind: KindDamage, Basic: true, Cooldown: 1.8, Range: 28, Power: 6, Style: "bolt", Color: "#c070ff",
		Desc: "Ein kleines Geschoss reiner Magie. Wird automatisch wiederholt."},
	"m_fireball": {ID: "m_fireball", Name: "Feuerball", Icon: "🔥", Kind: KindDamage, Cost: 18, Cooldown: 0, Cast: 1.5, Range: 30, Power: 34, Style: "fire", Color: "#ff7a1a",
		Desc: "Schleudert einen Feuerball (1,5 Sek. Zauberzeit)."},
	"m_frost": {ID: "m_frost", Name: "Frostblitz", Icon: "❄️", Kind: KindDamage, Cost: 14, Cooldown: 3, Cast: 1.0, Range: 30, Power: 20, Slow: 0.4, Duration: 5, Style: "frost", Color: "#7ad8ff",
		Desc: "Verursacht Frostschaden und verlangsamt das Ziel 5 Sek. lang."},
	"m_nova": {ID: "m_nova", Name: "Arkane Explosion", Icon: "💥", Kind: KindAoE, Cost: 30, Cooldown: 8, Radius: 8, Power: 22, Style: "nova", Color: "#c070ff",
		Desc: "Eine Explosion arkaner Energie trifft alle Feinde in der Nähe."},
	"m_shield": {ID: "m_shield", Name: "Arkaner Schild", Icon: "🛡️", Kind: KindShield, Cost: 25, Cooldown: 30, Power: 45, Duration: 15, Style: "buff", Color: "#8899ff",
		Desc: "Absorbiert Schaden (steigt mit der Stufe) für 15 Sek."},
}
