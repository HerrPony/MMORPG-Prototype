package gamedata

import (
	"fmt"
	"math"
	"strings"
)

// Wider loot tables: every mob drops a few more trade goods (generated from its name and kind), food and
// potions fitting its level, and open-world mobs drop equipment more often.

// mobKind: stem used in item names and the kind of body parts the mob leaves behind.
var mobKind = map[string][2]string{
	"wolf": {"Wolfs", "beast"}, "boar": {"Keiler", "beast"}, "spider": {"Spinnen", "insect"}, "bandit": {"Räuber", "human"}, "bear": {"Bären", "beast"},
	"hyena": {"Hyänen", "beast"}, "plains_bandit": {"Wegelagerer", "human"}, "bull": {"Stier", "beast"}, "croc": {"Krokodil", "reptile"},
	"zombie": {"Zombie", "undead"}, "witch": {"Hexen", "human"}, "wisp": {"Irrlicht", "spirit"}, "icewolf": {"Eiswolfs", "beast"}, "yeti": {"Yeti", "beast"},
	"polarbear": {"Eisbären", "beast"}, "frostspirit": {"Eisgeist", "spirit"}, "scorpion": {"Skorpion", "insect"}, "desert_raider": {"Wüstenräuber", "human"},
	"mummy": {"Mumien", "undead"}, "jungle_tiger": {"Tiger", "beast"}, "panther": {"Panther", "beast"}, "river_croc": {"Flusskrokodil", "reptile"},
	"jungle_ape": {"Affen", "beast"}, "poacher": {"Wilderer", "human"}, "lizardman": {"Echsen", "reptile"}, "lizard_shaman": {"Schamanen", "reptile"},
	"fox_spirit": {"Fuchs", "spirit"}, "bamboo_bear": {"Bambusbär", "beast"}, "ronin": {"Ronin", "human"}, "ghost_warrior": {"Geisterkrieger", "spirit"},
	"stone_golem": {"Steinwächter", "elemental"}, "sandcat": {"Sandkatzen", "beast"}, "canyon_bison": {"Bison", "beast"}, "mesa_scorpion": {"Kluftskorpion", "insect"},
	"canyon_bandit": {"Schluchtbandit", "human"}, "clay_golem": {"Lehmgolem", "elemental"}, "ash_hound": {"Aschehund", "beast"}, "fire_imp": {"Imp", "spirit"},
	"cultist": {"Kultisten", "human"}, "magma_golem": {"Magma", "elemental"}, "lava_elemental": {"Lava", "elemental"}, "obsidian_knight": {"Obsidianritter", "undead"},
	// elites of the open world
	"bear_elder": {"Grimmbart", "beast"}, "bandit_chief": {"Hauptmann", "human"}, "water_corrupt": {"Wasser", "elemental"}, "wolf_alpha": {"Alphawolf", "beast"},
	"hag": {"Grimmel", "human"}, "bog_colossus": {"Moorkoloss", "elemental"}, "ice_mother": {"Eismutter", "beast"}, "yeti_chief": {"Häuptlings", "beast"},
	"scorpion_king": {"Skorpionkönig", "insect"}, "tomb_guardian": {"Grabwächter", "undead"}, "whitetiger": {"Weißtiger", "beast"}, "poacher_boss": {"Anführer", "human"},
	"ape_king": {"Silberrücken", "beast"}, "lizard_chief": {"Häuptlings", "reptile"}, "fox_queen": {"Fuchskönigin", "spirit"}, "ronin_master": {"Meister", "human"},
	"temple_colossus": {"Kolossal", "elemental"}, "scorpion_matriarch": {"Matriarchen", "insect"}, "bandit_king": {"Banditenkönig", "human"}, "old_titan": {"Titan", "elemental"},
	"flame_lord": {"Flammenfürst", "spirit"}, "magma_titan": {"Magmatitan", "elemental"}, "obsidian_lord": {"Obsidianwächter", "undead"},
}

type part struct {
	suffix, icon, rarity string
	chance, value        float64
}

var kindParts = map[string][3]part{
	"beast":     {{"knochen", "🦴", "common", 0.35, 1}, {"kralle", "🐾", "common", 0.30, 1.3}, {"herz", "🫀", "uncommon", 0.06, 4}},
	"insect":    {{"-Chitin", "🪲", "common", 0.35, 1}, {"-Fühler", "🦗", "common", 0.30, 1.2}, {"-Giftdrüse", "🧪", "uncommon", 0.06, 4}},
	"human":     {{"-Beutel", "👛", "common", 0.35, 1.1}, {"-Kompass", "🧭", "common", 0.28, 1.4}, {"-Schmuckstück", "💎", "uncommon", 0.06, 4.5}},
	"undead":    {{"-Knochensplitter", "💀", "common", 0.35, 1}, {"-Grabstaub", "⚱️", "common", 0.30, 1.2}, {"-Seelenfunke", "✨", "uncommon", 0.06, 4.5}},
	"spirit":    {{"-Hauch", "💨", "common", 0.35, 1}, {"-Schleier", "🕯️", "common", 0.28, 1.3}, {"-Essenzkristall", "🔮", "uncommon", 0.06, 4.5}},
	"elemental": {{"-Splitter", "🪨", "common", 0.35, 1}, {"-Kernstück", "💠", "common", 0.28, 1.3}, {"-Kristall", "💎", "uncommon", 0.06, 4.5}},
	"reptile":   {{"-Schuppe", "🦎", "common", 0.35, 1}, {"-Auge", "👁️", "common", 0.28, 1.2}, {"-Giftzahn", "🦷", "uncommon", 0.06, 4}},
}

// food and potions by level
var foods = []struct {
	id, name, icon string
	heal, req      int
}{
	{"food_bread", "Frisches Brot", "🍞", 60, 0}, {"food_cheese", "Würziger Käse", "🧀", 130, 6}, {"food_stew", "Herzhafter Eintopf", "🍲", 280, 12},
	{"food_roast", "Knuspriger Braten", "🍗", 540, 18}, {"food_rice", "Gewürzreis", "🍚", 900, 24}, {"food_feast", "Festtagsbraten", "🍖", 1500, 30}, {"food_dragon", "Glutgebratenes", "🥘", 2400, 36},
}

func init() {
	for i, f := range foods {
		Items[f.id] = &Item{ID: f.id, Name: f.name, Icon: f.icon, Rarity: "common", Desc: fmt.Sprintf("Stellt sofort %d Lebenspunkte wieder her.", f.heal), Heal: f.heal, Stack: 20,
			Buy: 6 + i*i*14 + i*20, Sell: 2 + i*i*3, Req: f.req}
	}
	foodFor := func(level int) string { return foods[min(len(foods)-1, level/6)].id }
	potionFor := func(level int) string {
		switch {
		case level >= 33:
			return "potion_hp6"
		case level >= 25:
			return "potion_hp5"
		case level >= 17:
			return "potion_hp4"
		case level >= 11:
			return "potion_hp3"
		case level >= 5:
			return "potion_hp2"
		}
		return "potion_hp"
	}
	for id, m := range Mobs {
		if m.Boss || m.RespawnSec >= 9000 { // dungeon mobs have their own loot
			continue
		}
		lvl := m.LevelMin
		if k, ok := mobKind[id]; ok {
			parts := kindParts[k[1]]
			for i, p := range parts {
				itemID := fmt.Sprintf("x_%s_%d", id, i+1)
				name := k[0] + p.suffix
				name = strings.ToUpper(name[:1]) + name[1:]
				sell := int(math.Round((5 + float64(lvl)*4.2) * p.value))
				Items[itemID] = &Item{ID: itemID, Name: name, Icon: p.icon, Rarity: p.rarity, Desc: "Ein Fundstück von " + m.Name + ". Der Händler zahlt dafür.", Sell: sell, Stack: 20}
				ch := p.chance
				if m.Elite {
					ch = math.Min(1, ch*1.8)
				}
				m.Loot = append(m.Loot, LootEntry{Item: itemID, Chance: ch, Min: 1, Max: 1 + b2i(p.rarity == "common" && !m.Elite)})
			}
		}
		m.Loot = append(m.Loot, LootEntry{Item: foodFor(lvl), Chance: 0.10, Min: 1, Max: 2}, LootEntry{Item: potionFor(lvl), Chance: 0.05, Min: 1, Max: 1})
		// equipment drops more often: about one piece in every 6 kills, elites almost always
		if !m.Elite {
			m.GearChance = math.Max(m.GearChance, 0.10+0.002*float64(lvl))
		} else {
			m.GearChance = math.Max(m.GearChance, 0.8)
		}
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
