package gamedata

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
)

// Equipment slots in display order. "finger" items go into finger1 or finger2.
var EquipSlots = []string{"head", "neck", "shoulder", "back", "chest", "wrist", "hands", "waist", "legs", "feet", "finger1", "finger2", "mainhand", "offhand"}

// PrimaryStat is the stat that raises a class' damage.
var PrimaryStat = map[string]string{"warrior": "str", "hunter": "agi", "mage": "int"}

const (
	tiers      = 10
	statHP     = 5
	statMana   = 4
	statDamage = 0.01  // damage bonus per point of the primary stat
	weaponDmg  = 0.012 // damage bonus per weapon damage point
)

var tierNames = [tiers]string{"Rekruten", "Wächter", "Veteranen", "Helden", "Gardisten", "Kämpfer", "Eliten", "Meister", "Champions", "Legenden"}
var tierRarity = [tiers]string{"common", "uncommon", "uncommon", "rare", "rare", "rare", "rare", "epic", "epic", "epic"}
var weaponDmgTier = [tiers]int{4, 7, 11, 16, 22, 29, 37, 46, 56, 68}

func tierPrice(t float64) float64 { return 12 * math.Pow(1.85, t) }

type armorKind struct {
	id, label, class, prim string
	base                   int
}

var armorKinds = []armorKind{
	{"cloth", "Stoff", "mage", "int", 3},
	{"leather", "Leder", "hunter", "agi", 8},
	{"mail", "Kette", "warrior", "str", 16},
}

type slotDef struct {
	id    string
	w     float64
	icon  string
	nouns map[string]string
}

var armorSlots = []slotDef{
	{"head", 1.2, "🪖", map[string]string{"cloth": "Magierhut", "leather": "Lederkappe", "mail": "Kettenhaube"}},
	{"shoulder", 1.0, "🦺", map[string]string{"cloth": "Schulterstoff", "leather": "Lederschulter", "mail": "Kettenschulter"}},
	{"chest", 1.5, "👕", map[string]string{"cloth": "Robe", "leather": "Lederwams", "mail": "Kettenhemd"}},
	{"wrist", 0.5, "⌚", map[string]string{"cloth": "Stoffbänder", "leather": "Lederbänder", "mail": "Kettenbänder"}},
	{"hands", 0.8, "🧤", map[string]string{"cloth": "Stoffhandschuhe", "leather": "Lederhandschuhe", "mail": "Kettenhandschuhe"}},
	{"waist", 0.6, "🎗️", map[string]string{"cloth": "Schärpe", "leather": "Ledergürtel", "mail": "Kettengürtel"}},
	{"legs", 1.3, "👖", map[string]string{"cloth": "Stoffhose", "leather": "Lederhose", "mail": "Kettenbeinlinge"}},
	{"feet", 0.8, "🥾", map[string]string{"cloth": "Stoffschuhe", "leather": "Lederstiefel", "mail": "Kettenstiefel"}},
}

// pools of generated gear by tier, used for random drops and vendors
var gearPool [tiers][]string

// dungeonPool holds the epic dungeon set (Schwarzfels).
var dungeonPool []string

func add(it *Item) {
	if it.Stack == 0 {
		it.Stack = 1
	}
	if it.Buy > 0 && it.Sell == 0 {
		it.Sell = max(1, it.Buy/4)
	}
	Items[it.ID] = it
}

// gearSet describes one tier of generated equipment. t is the (fractional) power tier.
type gearSet struct {
	prefix string
	suffix string // id suffix
	t      float64
	req    int
	rarity string
}

func (g gearSet) item(id, name, icon, slot, typ string, w float64) *Item {
	return &Item{ID: id, Name: g.prefix + "-" + name, Icon: icon, Rarity: g.rarity, Slot: slot, Type: typ, Req: g.req,
		Buy: int(math.Round(tierPrice(g.t) * w)), Stats: map[string]int{}}
}

// genSet creates a complete set (armor, weapons, offhands, cloak, rings, amulets) and returns the ids.
func genSet(g gearSet) []string {
	var ids []string
	push := func(it *Item) { add(it); ids = append(ids, it.ID) }
	t := g.t
	for _, k := range armorKinds {
		for _, s := range armorSlots {
			it := g.item(fmt.Sprintf("%s_%s_%s", k.id, s.id, g.suffix), s.nouns[k.id], s.icon, s.id, k.label, s.w)
			if k.id == "cloth" && s.id == "head" {
				it.Icon = "🎩"
			}
			it.Armor = int(math.Round(float64(k.base) * s.w * (1 + 1.2*t)))
			it.Stats[k.prim] = max(1, int(math.Round(s.w*(1+1.6*t))))
			it.Stats["sta"] = max(1, int(math.Round(s.w*(1.2+1.4*t))))
			it.Classes = []string{k.class}
			push(it)
		}
	}
	type wdef struct{ id, name, icon, slot, typ, class, prim string }
	wd := weaponDmgTier[min(tiers-1, int(math.Round(t)))]
	for _, d := range []wdef{
		{"sword", "Schwert", "🗡️", "mainhand", "Einhand-Schwert", "warrior", "str"},
		{"bow", "Jagdbogen", "🏹", "mainhand", "Bogen", "hunter", "agi"},
		{"staff", "Stab", "🪄", "mainhand", "Stab", "mage", "int"},
	} {
		it := g.item(fmt.Sprintf("%s_%s", d.id, g.suffix), d.name, d.icon, d.slot, d.typ, 1.6)
		it.Weapon = wd
		it.Stats[d.prim] = 1 + int(t*2)
		if t > 0 {
			it.Stats["sta"] = int(t)
		}
		it.Classes = []string{d.class}
		push(it)
	}
	sh := g.item("shield_"+g.suffix, "Schild", "🛡️", "offhand", "Schild", 1.2)
	sh.Armor = 12 + int(t*14)
	sh.Stats["sta"] = 1 + int(t*2)
	sh.Classes = []string{"warrior"}
	push(sh)
	orb := g.item("orb_"+g.suffix, "Zauberkugel", "🔮", "offhand", "Zauberkugel", 1.2)
	orb.Stats["int"] = 1 + int(t*2)
	orb.Stats["sta"] = 1 + int(t)
	orb.Classes = []string{"mage"}
	push(orb)
	cl := g.item("back_"+g.suffix, "Umhang", "🧥", "back", "Stoff", 0.7)
	cl.Armor = 2 + int(t*2)
	cl.Stats["sta"] = 1 + int(t)
	push(cl)
	statNames := map[string]string{"str": "der Stärke", "agi": "der Gewandtheit", "int": "der Weisheit"}
	for _, p := range []string{"str", "agi", "int"} {
		ring := g.item(fmt.Sprintf("ring_%s_%s", p, g.suffix), "Ring "+statNames[p], "💍", "finger", "Ring", 1.0)
		ring.Stats[p] = 1 + int(t*2)
		ring.Stats["sta"] = 1 + int(t)
		neck := g.item(fmt.Sprintf("neck_%s_%s", p, g.suffix), "Amulett "+statNames[p], "📿", "neck", "Amulett", 1.1)
		neck.Stats[p] = 1 + int(t*2)
		neck.Stats["sta"] = 1 + int(t)
		push(ring)
		push(neck)
	}
	return ids
}

func init() {
	for t := 0; t < tiers; t++ {
		gearPool[t] = genSet(gearSet{prefix: tierNames[t], suffix: fmt.Sprint(t), t: float64(t), req: t*2 + 1, rarity: tierRarity[t]})
	}
	// the dungeon set is a bit stronger than the open world tier of the same level
	dungeonPool = genSet(gearSet{prefix: "Schwarzfels", suffix: "bf", t: 8.6, req: 17, rarity: "epic"})
	for _, id := range dungeonPool {
		Items[id].Desc = "Aus den Tiefen der Schwarzfelsmine."
	}

	// --- consumables ---
	add(&Item{ID: "potion_hp", Name: "Kleiner Heiltrank", Icon: "🧪", Rarity: "common", Desc: "Stellt sofort 80 Lebenspunkte wieder her.", Heal: 80, Stack: 10, Buy: 25, Sell: 6})
	add(&Item{ID: "potion_hp2", Name: "Heiltrank", Icon: "⚗️", Rarity: "uncommon", Desc: "Stellt sofort 200 Lebenspunkte wieder her.", Heal: 200, Stack: 10, Buy: 90, Sell: 22, Req: 5})
	add(&Item{ID: "potion_hp3", Name: "Großer Heiltrank", Icon: "🍶", Rarity: "rare", Desc: "Stellt sofort 450 Lebenspunkte wieder her.", Heal: 450, Stack: 10, Buy: 260, Sell: 65, Req: 11})
	add(&Item{ID: "potion_hp4", Name: "Meisterheiltrank", Icon: "🏺", Rarity: "rare", Desc: "Stellt sofort 800 Lebenspunkte wieder her.", Heal: 800, Stack: 10, Buy: 700, Sell: 175, Req: 17})

	// --- quest reward items get real stats ---
	set := func(id, slot, typ string, armor int, stats map[string]int, req int, classes ...string) {
		it := Items[id]
		it.Slot, it.Type, it.Armor, it.Stats, it.Req, it.Classes = slot, typ, armor, stats, req, classes
	}
	set("ranger_cloak", "back", "Stoff", 6, map[string]int{"sta": 2, "agi": 1}, 1)
	set("scholar_charm", "neck", "Amulett", 0, map[string]int{"int": 3, "sta": 2}, 3)
	set("guard_badge", "neck", "Amulett", 0, map[string]int{"str": 4, "sta": 3}, 4)
	set("hunter_ring", "finger", "Ring", 0, map[string]int{"agi": 5, "sta": 3}, 5)

	for id, c := range map[string]float64{"wolf": 0.05, "boar": 0.06, "spider": 0.08, "bandit": 0.12, "bear": 0.12} {
		Mobs[id].GearChance = c
	}
}

// RandomGear picks a piece of equipment suitable for a mob of the given level and
// usable by the given class (class-less jewelry and cloaks count for everyone).
func RandomGear(level int, class string, rng *rand.Rand) string {
	t := min(tiers-1, max(0, (level-1)/2))
	var usable []string
	for _, id := range gearPool[t] {
		it := Items[id]
		if len(it.Classes) == 0 || class == "" || slices.Contains(it.Classes, class) {
			usable = append(usable, id)
		}
	}
	return usable[rng.Intn(len(usable))]
}

// DungeonGear picks a piece of the dungeon set usable by the class.
func DungeonGear(class string, rng *rand.Rand) string {
	var usable []string
	for _, id := range dungeonPool {
		it := Items[id]
		if len(it.Classes) == 0 || slices.Contains(it.Classes, class) {
			usable = append(usable, id)
		}
	}
	return usable[rng.Intn(len(usable))]
}

// stockFor lists vendor goods between two power tiers. Kind: "weapons", "armor", "jewelry".
func stockFor(kind string, tmin, tmax int) []string {
	var out []string
	for t := tmin; t <= tmax && t < tiers; t++ {
		for _, id := range gearPool[t] {
			it := Items[id]
			isJewel := it.Slot == "finger" || it.Slot == "neck" || it.Slot == "back"
			isWeapon := it.Slot == "mainhand" || it.Slot == "offhand"
			switch kind {
			case "weapons":
				if isWeapon {
					out = append(out, id)
				}
			case "armor":
				if !isJewel && !isWeapon {
					out = append(out, id)
				}
			case "jewelry":
				if isJewel {
					out = append(out, id)
				}
			}
		}
	}
	return out
}
