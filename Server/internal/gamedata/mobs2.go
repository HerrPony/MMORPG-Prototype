package gamedata

import (
	"hash/fnv"
	"math"
	"math/rand"
)

// addMob fills in derived values for the generated mobs.
func addMob(m *Mob) {
	if m.AttackSec == 0 {
		m.AttackSec = 2.0
	}
	if m.Speed == 0 {
		m.Speed = 5.2
	}
	if m.Leash == 0 {
		m.Leash = 45
	}
	if m.RespawnSec == 0 {
		m.RespawnSec = 45
	}
	if m.XPBase == 0 {
		m.XPBase = 20 + 8*max(0, m.LevelMin-20) // higher level mobs are worth more
	}
	if m.Scale == 0 {
		m.Scale = 1
	}
	if m.LevelMax < m.LevelMin {
		m.LevelMax = m.LevelMin
	}
	if m.Aggressive && m.AggroRange == 0 {
		m.AggroRange = 8
	}
	hpm, dm := m.HPMult, m.DmgMult
	if m.Elite {
		if hpm == 0 {
			hpm = 3.2
		}
		if dm == 0 {
			dm = 1.5
		}
	}
	if hpm == 0 {
		hpm = 1
	}
	if dm == 0 {
		dm = 1
	}
	if m.HPBase == 0 {
		lo := m.LevelMin
		m.HPPerLevel = int(math.Round((mobHP(lo+1) - mobHP(lo)) * hpm))
		m.HPBase = int(math.Round(mobHP(lo)*hpm)) - m.HPPerLevel*(lo-1)
	}
	if m.DmgMax == 0 {
		m.DmgMin = max(1, int(mobDmg(m.LevelMin)*0.8*dm))
		m.DmgMax = max(m.DmgMin+1, int(mobDmg(m.LevelMax)*1.2*dm))
	}
	if m.CopperMax == 0 {
		m.CopperMin = 2 + 2*m.LevelMin
		m.CopperMax = 4 + 4*m.LevelMin
		if m.Elite {
			m.CopperMin *= 3
			m.CopperMax *= 3
		}
	}
	Mobs[m.ID] = m
}

func l(item string, chance float64, mn, mx int) LootEntry {
	return LootEntry{Item: item, Chance: chance, Min: mn, Max: mx}
}
func ql(item, quest string, chance float64, mn, mx int) LootEntry {
	return LootEntry{Item: item, Chance: chance, Min: mn, Max: mx, Quest: quest}
}

func init() {
	for _, id := range []string{"wolf", "boar", "bear", "spider"} {
		Mobs[id].Tameable = true
	}
	// ---- Goldene Ebenen (6-10) ----
	addMob(&Mob{ID: "hyena", Name: "Steppenhyäne", Model: "hyena", LevelMin: 6, LevelMax: 7, Aggressive: true, AggroRange: 9, Tameable: true,
		Loot: []LootEntry{l("hyena_pelt", 0.6, 1, 1), l("wolf_fang", 0.3, 1, 1)}, GearChance: 0.06})
	addMob(&Mob{ID: "plains_bandit", Name: "Wegelagerer", Model: "bandit", LevelMin: 7, LevelMax: 8, Aggressive: true,
		Loot: []LootEntry{l("rusty_dagger", 0.4, 1, 1), l("bandit_mask", 0.3, 1, 1), ql("mill_gear", "p_mill", 0.7, 1, 1), ql("shield_shards", "w_trial", 0.5, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "bull", Name: "Wildstier", Model: "bull", LevelMin: 8, LevelMax: 9, Tameable: true,
		Loot: []LootEntry{l("boar_meat", 0.5, 1, 2), ql("bull_horn", "p_bulls", 0.7, 1, 1)}, GearChance: 0.06})
	addMob(&Mob{ID: "wolf_alpha", Name: "Alpha-Wolf Rotpelz", Model: "wolf_alpha", LevelMin: 9, LevelMax: 9, Tameable: true, Scale: 1.25, HPMult: 1.6, RespawnSec: 120,
		Loot: []LootEntry{l("wolf_pelt", 1, 1, 2), l("wolf_fang", 1, 1, 2)}})
	addMob(&Mob{ID: "bandit_chief", Name: "Räuberhauptmann Dorn", Model: "bandit_chief", LevelMin: 10, Elite: true, Aggressive: true, AggroRange: 10, Scale: 1.25, RespawnSec: 240,
		Loot: []LootEntry{l("stolen_purse", 1, 1, 2)}, GearChance: 0.7})
	addMob(&Mob{ID: "water_corrupt", Name: "Verzerrter Wasserelementar", Model: "water_elem", LevelMin: 10, Elite: true, Range: 20, Style: "frost", Scale: 1.3, RespawnSec: 240,
		Loot: []LootEntry{ql("crystal_heart", "m_heart", 1, 1, 1)}, GearChance: 0.7})

	// ---- Düstermoor (10-14) ----
	addMob(&Mob{ID: "croc", Name: "Sumpfkrokodil", Model: "croc", LevelMin: 10, LevelMax: 11, Aggressive: true, AggroRange: 7, Tameable: true,
		Loot: []LootEntry{l("croc_hide", 0.55, 1, 1), l("wolf_meat", 0.3, 1, 1)}, GearChance: 0.07})
	addMob(&Mob{ID: "zombie", Name: "Moorzombie", Model: "zombie", LevelMin: 11, LevelMax: 12, Aggressive: true, AggroRange: 8, Speed: 3.6,
		Loot: []LootEntry{l("rotten_cloth", 0.5, 1, 1)}, GearChance: 0.08})
	addMob(&Mob{ID: "wisp", Name: "Irrlicht", Model: "wisp", LevelMin: 12, LevelMax: 13, Aggressive: true, AggroRange: 14, Range: 18, Style: "bolt", HPMult: 0.8,
		Loot: []LootEntry{ql("wisp_dust", "s_wisps", 0.8, 1, 1)}, GearChance: 0.07})
	addMob(&Mob{ID: "witch", Name: "Moorhexe", Model: "witch", LevelMin: 12, LevelMax: 13, Aggressive: true, AggroRange: 12, Range: 20, Style: "fire", HPMult: 0.9,
		Loot: []LootEntry{l("witch_herb", 0.5, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "hag", Name: "Mutter Grimmel", Model: "hag", LevelMin: 13, Elite: true, Aggressive: true, AggroRange: 12, Range: 20, Style: "fire", Scale: 1.3, RespawnSec: 240,
		Loot: []LootEntry{l("witch_herb", 1, 2, 3)}, GearChance: 0.8})
	addMob(&Mob{ID: "bog_colossus", Name: "Alter Moorkoloss", Model: "zombie_giant", LevelMin: 14, Elite: true, Scale: 2.0, Speed: 4, RespawnSec: 300,
		Loot: []LootEntry{l("rotten_cloth", 1, 2, 3)}, GearChance: 0.8})

	// ---- Eisgipfel (14-18) ----
	addMob(&Mob{ID: "icewolf", Name: "Eiswolf", Model: "icewolf", LevelMin: 14, LevelMax: 15, Tameable: true,
		Loot: []LootEntry{ql("icewolf_pelt", "f_wolves", 0.75, 1, 1), l("wolf_fang", 0.3, 1, 1)}, GearChance: 0.07})
	addMob(&Mob{ID: "yeti", Name: "Yeti", Model: "yeti", LevelMin: 15, LevelMax: 16, Aggressive: true, AggroRange: 8,
		Loot: []LootEntry{l("yeti_fur", 0.5, 1, 1)}, GearChance: 0.09})
	addMob(&Mob{ID: "polarbear", Name: "Eisbär", Model: "polarbear", LevelMin: 16, LevelMax: 17, Tameable: true,
		Loot: []LootEntry{l("bear_claw", 0.5, 1, 2), l("bear_pelt", 0.35, 1, 1)}, GearChance: 0.08})
	addMob(&Mob{ID: "frostspirit", Name: "Eisgeist", Model: "frostspirit", LevelMin: 16, LevelMax: 17, Aggressive: true, AggroRange: 14, Range: 18, Style: "frost", HPMult: 0.8,
		Loot: []LootEntry{ql("ice_shard", "f_spirits", 0.8, 1, 1)}, GearChance: 0.08})
	addMob(&Mob{ID: "yeti_chief", Name: "Yeti-Häuptling Grimmhorn", Model: "yeti_chief", LevelMin: 17, Elite: true, Aggressive: true, AggroRange: 10, Scale: 1.5, RespawnSec: 300,
		Loot: []LootEntry{l("yeti_fur", 1, 2, 3)}, GearChance: 0.85})
	addMob(&Mob{ID: "ice_mother", Name: "Eisbärenmutter", Model: "polarbear_big", LevelMin: 16, Elite: true, Scale: 1.5, RespawnSec: 300,
		Loot: []LootEntry{l("bear_pelt", 1, 1, 2)}, GearChance: 0.8})

	// ---- Glutwüste (17-20) ----
	addMob(&Mob{ID: "scorpion", Name: "Sandskorpion", Model: "scorpion", LevelMin: 17, LevelMax: 18, Aggressive: true, AggroRange: 8, Tameable: true,
		Loot: []LootEntry{l("scorpion_tail", 0.5, 1, 1)}, GearChance: 0.08})
	addMob(&Mob{ID: "desert_raider", Name: "Wüstenräuber", Model: "desert_raider", LevelMin: 18, LevelMax: 19, Aggressive: true, AggroRange: 9,
		Loot: []LootEntry{ql("caravan_map", "g_raiders", 0.7, 1, 1), l("bandit_mask", 0.3, 1, 1)}, GearChance: 0.11})
	addMob(&Mob{ID: "mummy", Name: "Mumie", Model: "mummy", LevelMin: 19, LevelMax: 20, Aggressive: true, AggroRange: 8, Speed: 4,
		Loot: []LootEntry{l("mummy_wrap", 0.5, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "scorpion_king", Name: "Skorpionkönig", Model: "scorpion_king", LevelMin: 19, Elite: true, Aggressive: true, AggroRange: 10, Scale: 2.4, RespawnSec: 300,
		Loot: []LootEntry{l("scorpion_tail", 1, 2, 3)}, GearChance: 0.85})
	addMob(&Mob{ID: "tomb_guardian", Name: "Grabwächter", Model: "tomb_guardian", LevelMin: 20, Elite: true, Aggressive: true, AggroRange: 9, Scale: 1.6, Speed: 4, RespawnSec: 300,
		Loot: []LootEntry{l("mummy_wrap", 1, 2, 3)}, GearChance: 0.9})

	// ---- Open world elite for the early game ----
	addMob(&Mob{ID: "bear_elder", Name: "Alter Grimmbart", Model: "bear_big", LevelMin: 6, Elite: true, Scale: 1.35, RespawnSec: 240,
		Loot: []LootEntry{l("bear_claw", 1, 2, 2)}, GearChance: 0.7})

	// ---- Dungeon: Schwarzfelsmine (16-20) ----
	addMob(&Mob{ID: "kobold", Name: "Minenkobold", Model: "kobold", LevelMin: 16, LevelMax: 17, Elite: true, HPMult: 2.0, DmgMult: 1.2, Aggressive: true, AggroRange: 12, Leash: 80, RespawnSec: 9999,
		Loot: []LootEntry{l("kobold_candle", 0.5, 1, 1)}, GearChance: 0.0})
	addMob(&Mob{ID: "kobold_shaman", Name: "Kobold-Schamane", Model: "kobold_shaman", LevelMin: 16, LevelMax: 17, Elite: true, HPMult: 1.6, DmgMult: 1.1, Range: 20, Style: "fire", Aggressive: true, AggroRange: 14, Leash: 80, RespawnSec: 9999,
		Loot: []LootEntry{l("kobold_candle", 0.5, 1, 1)}})
	addMob(&Mob{ID: "skeleton_miner", Name: "Skelettbergmann", Model: "skeleton", LevelMin: 17, LevelMax: 17, Elite: true, HPMult: 2.2, Aggressive: true, AggroRange: 12, Leash: 80, RespawnSec: 9999,
		Loot: []LootEntry{l("bone_dust", 0.5, 1, 1)}})
	addMob(&Mob{ID: "cave_spider", Name: "Höhlenspinne", Model: "cavespider", LevelMin: 17, LevelMax: 18, Elite: true, HPMult: 1.8, Aggressive: true, AggroRange: 12, Leash: 80, RespawnSec: 9999,
		Loot: []LootEntry{l("spider_silk", 0.5, 1, 2)}})
	addMob(&Mob{ID: "spiderling", Name: "Spinnenbrut", Model: "cavespider", LevelMin: 18, LevelMax: 18, HPMult: 0.35, DmgMult: 0.6, Scale: 0.5, Aggressive: true, AggroRange: 20, Leash: 80, RespawnSec: 9999, XPBase: 4})
	addMob(&Mob{ID: "kobold_add", Name: "Kobold-Leibwache", Model: "kobold", LevelMin: 19, LevelMax: 19, HPMult: 0.8, Aggressive: true, AggroRange: 20, Leash: 80, RespawnSec: 9999, XPBase: 10})
	addMob(&Mob{ID: "foreman", Name: "Vorarbeiter Grimm", Model: "kobold_boss", LevelMin: 18, Elite: true, Boss: true, HPMult: 4.5, DmgMult: 1.4, Scale: 1.9, Aggressive: true, AggroRange: 14, Leash: 90, RespawnSec: 9999,
		Abilities: []string{"heavy:9", "summon:50:kobold_add:2"}})
	addMob(&Mob{ID: "spider_queen", Name: "Spinnenkönigin Vessa", Model: "spider_queen", LevelMin: 19, Elite: true, Boss: true, HPMult: 5.0, DmgMult: 1.3, Scale: 2.6, Aggressive: true, AggroRange: 14, Leash: 90, RespawnSec: 9999,
		Abilities: []string{"poison:10", "summon:60:spiderling:4", "summon:30:spiderling:4"}})
	addMob(&Mob{ID: "kobold_king", Name: "Kobold-König Skarn", Model: "kobold_king", LevelMin: 20, Elite: true, Boss: true, HPMult: 7.0, DmgMult: 1.5, Scale: 2.3, Aggressive: true, AggroRange: 16, Leash: 100, RespawnSec: 9999,
		Abilities: []string{"heavy:8", "aoe:12", "summon:70:kobold_add:3", "summon:40:kobold_add:3", "enrage:20"}})
}

// ---- spawn groups ----

type SpawnGroup struct {
	Mob    string
	X, Z   float64
	N      int
	Radius float64
}

func groups() []SpawnGroup {
	return append(groups1(), groups2()...)
}

func groups1() []SpawnGroup {
	return []SpawnGroup{
		// Erlenwald elite
		{"bear_elder", 70, 40, 1, 2},
		// Goldene Ebenen
		{"hyena", 150, -62, 7, 14}, {"hyena", 105, -150, 6, 14},
		{"plains_bandit", 258, -68, 8, 12},
		{"bull", 120, -190, 5, 14}, {"bull", 100, -105, 5, 14},
		{"wolf_alpha", 135, -218, 1, 1},
		{"hyena", 130, -214, 4, 10},
		{"bandit_chief", 262, -64, 1, 1},
		{"water_corrupt", 178, -234, 1, 1},
		// Düstermoor
		{"croc", -130, 212, 8, 18}, {"zombie", -95, 150, 8, 16}, {"zombie", -220, 130, 5, 12},
		{"wisp", -176, 205, 6, 16}, {"witch", -225, 92, 7, 12},
		{"hag", -228, 96, 1, 1}, {"bog_colossus", -92, 178, 1, 1},
		// Eisgipfel
		{"icewolf", -60, -172, 7, 16}, {"icewolf", -150, -170, 5, 14},
		{"yeti", -40, -248, 6, 14}, {"yeti", -140, -250, 5, 14},
		{"polarbear", 30, -255, 5, 16},
		{"frostspirit", -70, -240, 6, 16},
		{"yeti_chief", -38, -254, 1, 1}, {"ice_mother", 36, -260, 1, 1},
		// Glutwüste
		{"scorpion", 245, 142, 9, 20}, {"scorpion", 150, 120, 6, 14},
		{"desert_raider", 225, 175, 7, 14},
		{"mummy", 140, 235, 9, 16},
		{"scorpion_king", 250, 148, 1, 1}, {"tomb_guardian", 138, 242, 1, 1},
	}
}

var generatedSpawns []Spawn

// BuildSpawns places the mob groups on valid ground (deterministic).
func buildSpawns() {
	for _, g := range groups() {
		g.X, g.Z = moveToWorld(g.X, g.Z)
		h := fnv.New32a()
		h.Write([]byte(g.Mob))
		rng := rand.New(rand.NewSource(int64(h.Sum32()) + int64(g.X*7) + int64(g.Z*13)))
		placed := 0
		for try := 0; try < 800 && placed < g.N; try++ {
			x, z := g.X, g.Z
			if g.N > 1 || g.Radius > 2 || try > 20 {
				a := rng.Float64() * 2 * math.Pi
				rad := g.Radius + float64(try/40)*3 // widen the search if the spot is blocked
				r := math.Sqrt(rng.Float64()) * rad
				x, z = g.X+math.Cos(a)*r, g.Z+math.Sin(a)*r
			}
			if Height(x, z) < 0.8 || nearNPC(x, z, 7) || inTownCore(x, z) || DistToRoad(x, z) < roadClear || BiomeAt(x, z) == BiomeLava {
				continue
			}
			tooClose := false
			gap := packSpacing
			if g.N == 1 { // single bosses/elites sit in their camp
				gap = 4
			}
			for _, s := range generatedSpawns {
				if math.Hypot(s.X-x, s.Z-z) < gap {
					tooClose = true
					break
				}
			}
			if tooClose {
				continue
			}
			generatedSpawns = append(generatedSpawns, Spawn{g.Mob, x, z})
			placed++
		}
	}
}

// packSpacing is the minimum distance between two spawns, so that you can pull one mob without its neighbours
// (aggro range 8, group pull radius 6.5, wandering radius 2.5).
const packSpacing = 14.0

// roadClear is the distance mobs keep from roads when they spawn.
const roadClear = 5.0

// offRoad moves a point away from the nearest road until it is clear of it (keeps points that already are).
func offRoad(x, z float64) (float64, float64) {
	if DistToRoad(x, z) >= roadClear {
		return x, z
	}
	for r := 2.0; r <= 24; r += 2 {
		for a := 0.0; a < 2*math.Pi; a += math.Pi / 8 {
			nx, nz := x+math.Cos(a)*r, z+math.Sin(a)*r
			if DistToRoad(nx, nz) >= roadClear && Height(nx, nz) >= 0.8 && BiomeAt(nx, nz) != BiomeLava {
				return nx, nz
			}
		}
	}
	return x, z
}

func nearNPC(x, z, r float64) bool {
	for _, n := range NPCs {
		if math.Hypot(n.X-x, n.Z-z) < r {
			return true
		}
	}
	return false
}

func inTownCore(x, z float64) bool {
	for _, t := range Towns {
		if math.Hypot(t.X-x, t.Z-z) < t.FlatR*0.8 {
			return true
		}
	}
	return false
}
