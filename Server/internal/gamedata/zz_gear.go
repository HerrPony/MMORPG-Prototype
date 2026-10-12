package gamedata

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
)

// More equipment: alternative weapon and off-hand families for every tier, "of the Guardian/of Might" variants of
// every piece (drops only), and hand-made named items that drop from elites. The per-tier pools (gearPool) are
// extended by the families, so vendors and drops pick them up automatically.

// gearVariants holds the flavoured versions per tier; they drop from mobs but are not sold by vendors.
var gearVariants [tiers][]string

func init() {
	for t := 0; t < tiers; t++ {
		g := gearSet{prefix: tierNames[t], suffix: fmt.Sprint(t), t: float64(t), req: t*2 + 1, rarity: tierRarity[t]}
		gearPool[t] = append(gearPool[t], genFamilies(g)...)
		gearVariants[t] = genVariants(gearPool[t], t)
	}
}

// genFamilies adds the alternative weapons and off-hands of one tier.
func genFamilies(g gearSet) []string {
	var ids []string
	t := g.t
	wd := weaponDmgTier[min(tiers-1, int(math.Round(t)))]
	push := func(it *Item) { add(it); ids = append(ids, it.ID) }
	type wdef struct {
		id, name, icon, typ, class, prim string
		dmg                              float64 // weapon damage relative to the standard weapon
		primMul, staMul                  float64
	}
	for _, d := range []wdef{
		{"axe", "Streitaxt", "🪓", "Einhand-Axt", "warrior", "str", 1.04, 2.2, 0.8},
		{"mace", "Streitkolben", "🔨", "Einhand-Streitkolben", "warrior", "str", 0.94, 1.7, 1.6},
		{"xbow", "Armbrust", "🎯", "Armbrust", "hunter", "agi", 1.06, 1.9, 0.9},
		{"cbow", "Kompositbogen", "🏹", "Bogen", "hunter", "agi", 0.95, 2.4, 0.6},
		{"scepter", "Zepter", "🔱", "Zepter", "mage", "int", 0.95, 2.3, 0.7},
		{"runestaff", "Runenstab", "🪄", "Stab", "mage", "int", 1.05, 1.8, 1.2},
	} {
		it := g.item(fmt.Sprintf("%s_%s", d.id, g.suffix), d.name, d.icon, "mainhand", d.typ, 1.6)
		it.Weapon = int(math.Round(float64(wd) * d.dmg))
		it.Stats[d.prim] = 1 + int(t*d.primMul)
		it.Stats["sta"] = 1 + int(t*d.staMul)
		it.Classes = []string{d.class}
		push(it)
	}
	tower := g.item("tower_"+g.suffix, "Turmschild", "🛡️", "offhand", "Schild", 1.3)
	tower.Armor = 16 + int(t*20)
	tower.Stats["sta"] = 2 + int(t*3)
	tower.Classes = []string{"warrior"}
	push(tower)
	quiver := g.item("quiver_"+g.suffix, "Köcher", "🎒", "offhand", "Köcher", 1.2)
	quiver.Stats["agi"] = 1 + int(t*2)
	quiver.Stats["sta"] = 1 + int(t)
	quiver.Classes = []string{"hunter"}
	push(quiver)
	tome := g.item("tome_"+g.suffix, "Zauberbuch", "📖", "offhand", "Zauberbuch", 1.2)
	tome.Stats["int"] = 1 + int(t*1.8)
	tome.Stats["sta"] = 1 + int(t*1.4)
	tome.Classes = []string{"mage"}
	push(tome)
	return ids
}

// genVariants creates "der Wacht" (tougher) and "der Macht" (more damage) versions of every item of a tier.
func genVariants(base []string, t int) []string {
	var out []string
	for _, id := range base {
		b := Items[id]
		if b.Slot == "" || len(b.Stats) == 0 {
			continue
		}
		prim := ""
		for _, p := range []string{"str", "agi", "int"} {
			if b.Stats[p] > 0 {
				prim = p
			}
		}
		for _, v := range []struct {
			key, name        string
			prim, sta, armor float64
		}{
			{"w", "der Wacht", 0.7, 1.6, 1.15},
			{"m", "der Macht", 1.4, 0.55, 0.9},
		} {
			c := *b
			c.ID, c.Name = id+"_"+v.key, b.Name+" "+v.name
			c.Stats = map[string]int{}
			for k, n := range b.Stats {
				switch {
				case k == "sta":
					c.Stats[k] = max(1, int(math.Round(float64(n)*v.sta)))
				case k == prim:
					c.Stats[k] = max(1, int(math.Round(float64(n)*v.prim)))
				default:
					c.Stats[k] = n
				}
			}
			c.Armor = int(math.Round(float64(b.Armor) * v.armor))
			c.Buy = int(float64(b.Buy) * 1.25)
			c.Sell = max(1, c.Buy/4)
			c.Classes = slices.Clone(b.Classes)
			Items[c.ID] = &c
			out = append(out, c.ID)
		}
	}
	return out
}

// named creates a hand made item of power tier t (slightly above the generated gear of that tier).
func named(id, name, icon, slot, typ string, t float64, rarity, desc string, classes ...string) {
	w := map[string]float64{"head": 1.2, "shoulder": 1.0, "chest": 1.5, "wrist": 0.5, "hands": 0.8, "waist": 0.6, "legs": 1.3, "feet": 0.8,
		"back": 0.7, "neck": 1.1, "finger": 1.0, "mainhand": 1.6, "offhand": 1.2}[slot]
	prim := map[string]string{"warrior": "str", "hunter": "agi", "mage": "int"}
	it := &Item{ID: id, Name: name, Icon: icon, Rarity: rarity, Desc: desc, Slot: slot, Type: typ, Req: int(t)*2 + 1, Classes: classes, Stats: map[string]int{},
		Buy: int(math.Round(tierPrice(t+0.8) * w * 1.6))}
	p := "str"
	if len(classes) == 1 {
		p = prim[classes[0]]
	}
	boost := 1.18
	switch slot {
	case "mainhand":
		it.Weapon = int(math.Round(float64(weaponDmgTier[min(tiers-1, int(t))]) * boost))
		it.Stats[p] = 1 + int(t*2.2)
		it.Stats["sta"] = 1 + int(t)
	case "offhand":
		it.Stats["sta"] = 2 + int(t*2.5)
		if classes[0] == "warrior" {
			it.Armor = 14 + int(t*18)
		} else {
			it.Stats[p] = 1 + int(t*2)
		}
	case "neck", "finger":
		it.Stats["sta"] = 1 + int(t*1.4)
		if len(classes) > 1 { // usable by several classes: every class finds its main stat on it
			for _, c := range classes {
				it.Stats[prim[c]] = 1 + int(t*2.0)
			}
		} else {
			it.Stats[p] = 1 + int(t*2.3)
		}
	case "back":
		it.Armor = 3 + int(t*2.2)
		it.Stats["sta"] = 1 + int(t*1.4)
		if len(classes) > 1 {
			for _, c := range classes {
				it.Stats[prim[c]] = 1 + int(t*1.1)
			}
		} else {
			it.Stats[p] = 1 + int(t*1.2)
		}
	default:
		k := map[string]float64{"warrior": 16, "hunter": 8, "mage": 3}[classes[0]]
		it.Armor = int(math.Round(k * w * (1 + 1.2*t) * boost))
		it.Stats[p] = max(1, int(math.Round(w*(1+1.6*t)*boost)))
		it.Stats["sta"] = max(1, int(math.Round(w*(1.2+1.4*t)*boost)))
	}
	it.Sell = max(1, it.Buy/4)
	add(it)
}

// namedDrops: item ids that an elite may drop on top of the regular loot (mob id -> items, 30 % each).
var namedDrops = map[string][]string{}

func drops(mob string, ids ...string) { namedDrops[mob] = append(namedDrops[mob], ids...) }

func init() {
	W, H, M := "warrior", "hunter", "mage"
	// ---- Erlenmark (tier 0-3) ----
	named("grimmbart_pelt", "Fellmantel des Grimmbart", "🧥", "back", "Fell", 2, "rare", "Ein Mantel aus dem Pelz des alten Bären. Er riecht noch nach Wald.", W, H, M)
	named("dorn_blade", "Dorns Räuberklinge", "🗡️", "mainhand", "Einhand-Schwert", 3, "rare", "Gestohlen, geschliffen und nie gereinigt.", W)
	named("dorn_hood", "Kapuze des Räuberhauptmanns", "🪖", "head", "Leder", 3, "rare", "Wer sie trägt, wird selten erkannt.", H)
	named("stream_pendant", "Anhänger des Wasserherzens", "📿", "neck", "Amulett", 3, "rare", "Ein Tropfen des verzerrten Elementars, gefangen in Kristall.", M)
	named("alpha_ring", "Ring des Rudelführers", "💍", "finger", "Ring", 4, "rare", "Wer ihn trägt, hört das Rudel heulen.", H)
	// ---- Goldene Ebenen / Düstermoor (tier 4-6) ----
	named("chief_greaves", "Beinschienen des Hauptmanns", "👖", "legs", "Kette", 4, "rare", "Schwer, laut und unerschütterlich.", W)
	named("grimmel_robe", "Robe der Mutter Grimmel", "👕", "chest", "Stoff", 6, "rare", "Genäht aus Moorgras und alten Flüchen.", M)
	named("grimmel_staff", "Krötenstab", "🪄", "mainhand", "Stab", 6, "rare", "Quakt leise, wenn man ihn zu fest greift.", M)
	named("colossus_belt", "Gürtel des Moorkolosses", "🎗️", "waist", "Kette", 6, "rare", "Zusammengehalten von Schlamm und Willenskraft.", W)
	named("moor_cloak", "Schleier der Irrlichter", "🧥", "back", "Stoff", 6, "rare", "Funkelt, wenn niemand hinsieht.", M, H)
	// ---- Eisgipfel (tier 7-8) ----
	named("mother_fur", "Mantel der Eisbärenmutter", "🧥", "back", "Fell", 8, "rare", "Dick genug für eine Nacht auf dem Gletscher.", W, H, M)
	named("grimmhorn_axe", "Grimmhorns Eisbeil", "🪓", "mainhand", "Einhand-Axt", 8, "epic", "Der Griff ist nie warm geworden.", W)
	named("grimmhorn_bow", "Frostsehne", "🏹", "mainhand", "Bogen", 8, "epic", "Die Sehne singt im Wind.", H)
	named("yeti_boots", "Yetistiefel", "🥾", "feet", "Leder", 8, "rare", "Sie hinterlassen große Spuren.", H)
	// ---- Sandfurt / Wüste (tier 9-10) ----
	named("king_stinger", "Stachel des Skorpionkönigs", "🗡️", "mainhand", "Dolch", 9, "epic", "Noch immer tropft Gift herab.", W)
	named("king_scepter", "Zepter der Dünen", "🔱", "mainhand", "Zepter", 9, "epic", "Sand rieselt aus dem Kopfstück.", M)
	named("guardian_helm", "Helm des Grabwächters", "🪖", "head", "Kette", 10, "epic", "Hohl, alt und voller Echo.", W)
	named("guardian_amulet", "Skarabäus der Pharaonen", "📿", "neck", "Amulett", 10, "epic", "Ein goldener Käfer, der leise summt.", H, M)
	// ---- Smaragddschungel (tier 10-12) ----
	named("tiger_cloak", "Mantel des Weißen Tigers", "🧥", "back", "Fell", 11, "epic", "Weiß mit goldenen Streifen.", W, H, M)
	named("shan_bow", "Shans Reißzahn", "🏹", "mainhand", "Bogen", 11, "epic", "Aus dem Zahn eines Tigers geschnitzt.", H)
	named("kral_vest", "Weste des Wilderer-Anführers", "👕", "chest", "Leder", 11, "rare", "Viele Taschen, alle voller Fallen.", H)
	named("bokko_bracers", "Armschienen des Silberrückens", "⌚", "wrist", "Kette", 12, "epic", "Mächtig wie die Arme, die sie trugen.", W)
	named("szar_staff", "Idolstab von Szar", "🪄", "mainhand", "Stab", 12, "epic", "Das Götzenbild darauf blinzelt.", M)
	// ---- Kirschblütenhain (tier 12-14) ----
	named("kyuuko_tail", "Kyuukos Neunschweif", "🧣", "back", "Stoff", 13, "epic", "Neun seidene Schweife, die sich von allein bewegen.", M, H)
	named("kyuuko_ring", "Fuchsperlen-Ring", "💍", "finger", "Ring", 13, "epic", "Die Perle lächelt.", M)
	named("takeda_katana", "Takedas Klinge", "🗡️", "mainhand", "Einhand-Schwert", 14, "epic", "Eine Klinge, die kein Wort verliert.", W)
	named("takeda_armor", "Rüstung des Ronin-Meisters", "👕", "chest", "Kette", 14, "epic", "Zerkratzt, aber unbesiegt.", W)
	named("colossus_core", "Kern des Tempelkolosses", "🔮", "offhand", "Zauberkugel", 14, "epic", "Pulsiert im Takt eines uralten Herzens.", M)
	// ---- Rotkluft-Canyon (tier 15-17) ----
	named("matriarch_gloves", "Handschuhe der Matriarchin", "🧤", "hands", "Leder", 16, "epic", "Kratzfest und leicht giftig.", H)
	named("rakar_blade", "Rakars Schluchtschneider", "🗡️", "mainhand", "Einhand-Schwert", 16, "epic", "Wurde nie ehrlich erworben.", W)
	named("rakar_crown", "Banditenkrone", "🪖", "head", "Stoff", 16, "epic", "Zu viele Edelsteine für zu wenig Verstand.", M)
	named("titan_fist", "Faust des Alten Titanen", "🔨", "mainhand", "Einhand-Streitkolben", 17, "epic", "Ein Stück Fels, das nicht aufhören will zu wachsen.", W)
	named("titan_core", "Titanenherz", "💎", "neck", "Amulett", 17, "epic", "Es schlägt. Langsam.", W, H, M)
	// ---- Aschenlande (tier 18-19) ----
	named("pyrax_staff", "Pyrax' Flammenstab", "🪄", "mainhand", "Stab", 18, "epic", "Heiß genug, um Handschuhe zu schmelzen.", M)
	named("pyrax_mantle", "Schulterumhang des Flammenfürsten", "🦺", "shoulder", "Stoff", 18, "epic", "Glimmt an den Rändern.", M)
	named("brokkar_hammer", "Brokkars Magmahammer", "🔨", "mainhand", "Einhand-Streitkolben", 19, "epic", "Jeder Schlag lässt den Boden beben.", W)
	named("brokkar_plate", "Magmaplatten", "👕", "chest", "Kette", 19, "epic", "Schmieden aus Vulkangestein.", W)
	named("zeth_visor", "Visier des Obsidianwächters", "🪖", "head", "Kette", 19, "epic", "Schwarz wie die Nacht, scharf wie Glas.", W)
	named("zeth_bow", "Obsidiansplitter", "🏹", "mainhand", "Bogen", 19, "epic", "Jeder Pfeil schneidet die Luft.", H)
	named("zeth_orb", "Auge des Obsidianwächters", "🔮", "offhand", "Zauberkugel", 19, "epic", "Es beobachtet dich zurück.", M)

	drops("bear_elder", "grimmbart_pelt")
	drops("bandit_chief", "dorn_blade", "dorn_hood", "chief_greaves")
	drops("water_corrupt", "stream_pendant")
	drops("wolf_alpha", "alpha_ring")
	drops("hag", "grimmel_robe", "grimmel_staff")
	drops("bog_colossus", "colossus_belt", "moor_cloak")
	drops("ice_mother", "mother_fur")
	drops("yeti_chief", "grimmhorn_axe", "grimmhorn_bow", "yeti_boots")
	drops("scorpion_king", "king_stinger", "king_scepter")
	drops("tomb_guardian", "guardian_helm", "guardian_amulet")
	drops("whitetiger", "tiger_cloak", "shan_bow")
	drops("poacher_boss", "kral_vest")
	drops("ape_king", "bokko_bracers")
	drops("lizard_chief", "szar_staff")
	drops("fox_queen", "kyuuko_tail", "kyuuko_ring")
	drops("ronin_master", "takeda_katana", "takeda_armor")
	drops("temple_colossus", "colossus_core")
	drops("scorpion_matriarch", "matriarch_gloves")
	drops("bandit_king", "rakar_blade", "rakar_crown")
	drops("old_titan", "titan_fist", "titan_core")
	drops("flame_lord", "pyrax_staff", "pyrax_mantle")
	drops("magma_titan", "brokkar_hammer", "brokkar_plate")
	drops("obsidian_lord", "zeth_visor", "zeth_bow", "zeth_orb")
}

// NamedDrop picks one of the named items a mob may drop, usable by the class ("" if none).
func NamedDrop(mob, class string, rng *rand.Rand) string {
	var usable []string
	for _, id := range namedDrops[mob] {
		it := Items[id]
		if len(it.Classes) == 0 || class == "" || slices.Contains(it.Classes, class) {
			usable = append(usable, id)
		}
	}
	if len(usable) == 0 {
		return ""
	}
	return usable[rng.Intn(len(usable))]
}
