package gamedata

// Mobs of the second continent (levels 20-40) and of the two new dungeons.

func init() {
	// ---- Smaragddschungel (20-25) ----
	addMob(&Mob{ID: "river_croc", Name: "Flusskrokodil", Model: "croc", LevelMin: 21, LevelMax: 22, Aggressive: true, AggroRange: 7, Tameable: true,
		Loot: []LootEntry{l("river_croc_tooth", 0.55, 1, 2), l("croc_hide", 0.3, 1, 1)}, GearChance: 0.06})
	addMob(&Mob{ID: "jungle_tiger", Name: "Dschungeltiger", Model: "tiger", LevelMin: 20, LevelMax: 21, Tameable: true,
		Loot: []LootEntry{l("tiger_pelt", 0.6, 1, 1), l("tiger_fang", 0.4, 1, 2)}, GearChance: 0.06})
	addMob(&Mob{ID: "panther", Name: "Schwarzpanther", Model: "panther", LevelMin: 21, LevelMax: 22, Aggressive: true, AggroRange: 10, Tameable: true,
		Loot: []LootEntry{l("panther_claw", 0.6, 1, 2)}, GearChance: 0.07})
	addMob(&Mob{ID: "jungle_ape", Name: "Dschungelaffe", Model: "ape", LevelMin: 22, LevelMax: 23, Aggressive: true, AggroRange: 9, Scale: 1.1,
		Loot: []LootEntry{l("ape_fur", 0.6, 1, 2)}, GearChance: 0.07})
	addMob(&Mob{ID: "poacher", Name: "Wilderer", Model: "poacher", LevelMin: 22, LevelMax: 23, Aggressive: true, AggroRange: 10,
		Loot: []LootEntry{l("poacher_trap", 0.45, 1, 1), ql("poacher_ledger", "j_poachers", 0.6, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "lizardman", Name: "Echsenkrieger", Model: "lizardman", LevelMin: 23, LevelMax: 24, Aggressive: true, AggroRange: 10,
		Loot: []LootEntry{l("lizard_scale", 0.55, 1, 2), l("snake_venom", 0.45, 1, 1), ql("lizard_idol", "j_idol", 0.5, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "lizard_shaman", Name: "Echsenschamane", Model: "lizard_shaman", LevelMin: 24, LevelMax: 24, Aggressive: true, AggroRange: 14, Range: 20, Style: "bolt", HPMult: 0.9,
		Loot: []LootEntry{l("lizard_scale", 0.4, 1, 1), l("snake_venom", 0.3, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "whitetiger", Name: "Weißer Tiger Shan", Model: "whitetiger", LevelMin: 24, Elite: true, Aggressive: true, AggroRange: 10, Scale: 1.5, RespawnSec: 300,
		Loot: []LootEntry{l("tiger_pelt", 1, 1, 2), ql("tiger_trophy", "j_whitetiger", 1, 1, 1)}, GearChance: 0.85})
	addMob(&Mob{ID: "poacher_boss", Name: "Wilderer-Anführer Kral", Model: "poacher", LevelMin: 24, Elite: true, Aggressive: true, AggroRange: 12, Scale: 1.3, RespawnSec: 300,
		Loot: []LootEntry{l("poacher_trap", 1, 1, 2)}, GearChance: 0.85})
	addMob(&Mob{ID: "lizard_chief", Name: "Echsenhäuptling Szar", Model: "lizardman", LevelMin: 25, Elite: true, Aggressive: true, AggroRange: 12, Scale: 1.45, RespawnSec: 300,
		Loot: []LootEntry{l("lizard_scale", 1, 2, 3), ql("lizard_idol", "j_idol", 1, 1, 1)}, GearChance: 0.85})
	addMob(&Mob{ID: "ape_king", Name: "Silberrücken Bokko", Model: "ape", LevelMin: 25, Elite: true, Aggressive: true, AggroRange: 10, Scale: 1.6, RespawnSec: 300,
		Loot: []LootEntry{l("ape_fur", 1, 2, 3)}, GearChance: 0.85})

	// ---- Kirschblütenhain (25-30) ----
	addMob(&Mob{ID: "fox_spirit", Name: "Fuchsgeist", Model: "kitsune", LevelMin: 25, LevelMax: 26, Aggressive: true, AggroRange: 12, Range: 18, Style: "fire", HPMult: 0.85, Tameable: true,
		Loot: []LootEntry{l("fox_tail", 0.55, 1, 1), ql("fox_pearl", "b_pearls", 0.45, 1, 1)}, GearChance: 0.08})
	addMob(&Mob{ID: "bamboo_bear", Name: "Bambusbär", Model: "panda", LevelMin: 26, LevelMax: 27, Tameable: true,
		Loot: []LootEntry{l("panda_fur", 0.55, 1, 1)}, GearChance: 0.07})
	addMob(&Mob{ID: "ronin", Name: "Verstoßener Ronin", Model: "ronin", LevelMin: 27, LevelMax: 28, Aggressive: true, AggroRange: 10,
		Loot: []LootEntry{l("ronin_sword", 0.45, 1, 1), ql("ronin_seal", "b_ronin", 0.55, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "ghost_warrior", Name: "Geisterkrieger", Model: "ghost_warrior", LevelMin: 28, LevelMax: 29, Aggressive: true, AggroRange: 10, Speed: 4.4,
		Loot: []LootEntry{l("spirit_essence", 0.5, 1, 1), ql("spirit_bell", "b_ghosts", 0.45, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "stone_golem", Name: "Steinwächter", Model: "stone_golem", LevelMin: 29, LevelMax: 30, Speed: 3.8, Scale: 1.4,
		Loot: []LootEntry{l("stone_core", 0.5, 1, 1)}, GearChance: 0.09})
	addMob(&Mob{ID: "fox_queen", Name: "Fuchskönigin Kyuuko", Model: "kitsune", LevelMin: 28, Elite: true, Aggressive: true, AggroRange: 12, Range: 20, Style: "fire", Scale: 1.7, RespawnSec: 300,
		Loot: []LootEntry{l("fox_tail", 1, 2, 3)}, GearChance: 0.9})
	addMob(&Mob{ID: "ronin_master", Name: "Ronin-Meister Takeda", Model: "ronin", LevelMin: 29, Elite: true, Aggressive: true, AggroRange: 12, Scale: 1.3, RespawnSec: 300,
		Loot: []LootEntry{l("ronin_sword", 1, 1, 2)}, GearChance: 0.9})
	addMob(&Mob{ID: "temple_colossus", Name: "Tempelkoloss", Model: "stone_golem", LevelMin: 30, Elite: true, Scale: 2.1, Speed: 3.6, RespawnSec: 300,
		Loot: []LootEntry{l("stone_core", 1, 2, 3)}, GearChance: 0.9})

	// ---- Rotkluft-Canyon (30-35) ----
	addMob(&Mob{ID: "sandcat", Name: "Sandkatze", Model: "sandcat", LevelMin: 30, LevelMax: 31, Aggressive: true, AggroRange: 10, Tameable: true,
		Loot: []LootEntry{l("sandcat_fang", 0.55, 1, 2)}, GearChance: 0.07})
	addMob(&Mob{ID: "canyon_bison", Name: "Canyon-Bison", Model: "bison", LevelMin: 31, LevelMax: 32, Tameable: true,
		Loot: []LootEntry{l("bison_horn", 0.55, 1, 1), l("boar_meat", 0.4, 1, 2)}, GearChance: 0.07})
	addMob(&Mob{ID: "mesa_scorpion", Name: "Kluftskorpion", Model: "scorpion", LevelMin: 32, LevelMax: 33, Aggressive: true, AggroRange: 9, Scale: 1.3, Tameable: true,
		Loot: []LootEntry{l("scorpion_chitin", 0.55, 1, 2)}, GearChance: 0.08})
	addMob(&Mob{ID: "canyon_bandit", Name: "Schluchtbandit", Model: "canyon_bandit", LevelMin: 33, LevelMax: 34, Aggressive: true, AggroRange: 10,
		Loot: []LootEntry{l("bandit_badge", 0.5, 1, 1), ql("mesa_map", "m_map", 0.5, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "clay_golem", Name: "Lehmgolem", Model: "stone_golem", LevelMin: 34, LevelMax: 35, Speed: 3.8, Scale: 1.5,
		Loot: []LootEntry{l("golem_heart", 0.5, 1, 1), ql("titan_shard", "m_titan", 0.5, 1, 1)}, GearChance: 0.09})
	addMob(&Mob{ID: "scorpion_matriarch", Name: "Skorpionmatriarchin", Model: "scorpion_king", LevelMin: 33, Elite: true, Aggressive: true, AggroRange: 10, Scale: 2.4, RespawnSec: 300,
		Loot: []LootEntry{l("scorpion_chitin", 1, 2, 3)}, GearChance: 0.9})
	addMob(&Mob{ID: "bandit_king", Name: "Banditenkönig Rakar", Model: "canyon_bandit", LevelMin: 34, Elite: true, Aggressive: true, AggroRange: 12, Scale: 1.35, RespawnSec: 300,
		Loot: []LootEntry{l("bandit_badge", 1, 2, 3)}, GearChance: 0.9})
	addMob(&Mob{ID: "old_titan", Name: "Der Alte Titan", Model: "stone_golem", LevelMin: 35, Elite: true, Scale: 2.6, Speed: 3.6, RespawnSec: 300,
		Loot: []LootEntry{l("golem_heart", 1, 1, 2), ql("titan_shard", "m_titan", 1, 1, 1)}, GearChance: 0.95})

	// ---- Aschenlande (35-40) ----
	addMob(&Mob{ID: "ash_hound", Name: "Aschehund", Model: "ashhound", LevelMin: 35, LevelMax: 36, Aggressive: true, AggroRange: 11, Tameable: true,
		Loot: []LootEntry{l("hound_tooth", 0.55, 1, 2)}, GearChance: 0.08})
	addMob(&Mob{ID: "fire_imp", Name: "Feuerimp", Model: "fire_imp", LevelMin: 36, LevelMax: 37, Aggressive: true, AggroRange: 13, Range: 18, Style: "fire", HPMult: 0.8,
		Loot: []LootEntry{l("imp_horn", 0.55, 1, 1)}, GearChance: 0.08})
	addMob(&Mob{ID: "cultist", Name: "Glutkultist", Model: "cultist", LevelMin: 37, LevelMax: 38, Aggressive: true, AggroRange: 13, Range: 20, Style: "fire", HPMult: 0.9,
		Loot: []LootEntry{l("cultist_scrap", 0.5, 1, 1), ql("cult_tome", "v_cultists", 0.5, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "magma_golem", Name: "Magmagolem", Model: "magma_golem", LevelMin: 38, LevelMax: 39, Aggressive: true, AggroRange: 9, Speed: 3.8, Scale: 1.5,
		Loot: []LootEntry{l("magma_core", 0.5, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "lava_elemental", Name: "Lavaelementar", Model: "lava_elem", LevelMin: 39, LevelMax: 40, Aggressive: true, AggroRange: 14, Range: 18, Style: "fire", HPMult: 0.9, Scale: 1.3,
		Loot: []LootEntry{l("obsidian_shard", 0.4, 1, 1), ql("lava_heart", "v_lava", 0.5, 1, 1)}, GearChance: 0.1})
	addMob(&Mob{ID: "obsidian_knight", Name: "Obsidianritter", Model: "obsidian_knight", LevelMin: 39, LevelMax: 40, Aggressive: true, AggroRange: 10, Scale: 1.15,
		Loot: []LootEntry{l("obsidian_shard", 0.5, 1, 2), ql("knight_emblem", "v_knights", 0.5, 1, 1)}, GearChance: 0.11})
	addMob(&Mob{ID: "flame_lord", Name: "Flammenfürst Pyrax", Model: "cultist", LevelMin: 38, Elite: true, Aggressive: true, AggroRange: 14, Range: 22, Style: "fire", Scale: 1.6, RespawnSec: 300,
		Loot: []LootEntry{l("cultist_scrap", 1, 2, 3)}, GearChance: 0.95})
	addMob(&Mob{ID: "magma_titan", Name: "Magmatitan Brokkar", Model: "magma_golem", LevelMin: 39, Elite: true, Aggressive: true, AggroRange: 10, Scale: 2.4, Speed: 3.6, RespawnSec: 300,
		Loot: []LootEntry{l("magma_core", 1, 2, 3)}, GearChance: 0.95})
	addMob(&Mob{ID: "obsidian_lord", Name: "Obsidianwächter Zeth", Model: "obsidian_knight", LevelMin: 40, Elite: true, Aggressive: true, AggroRange: 12, Scale: 1.6, RespawnSec: 300,
		Loot: []LootEntry{l("obsidian_shard", 1, 2, 3)}, GearChance: 0.95})

	// ---- Dungeon: Nebeltempel (27-30) ----
	addMob(&Mob{ID: "naga_warrior", Name: "Nagakrieger", Model: "naga_warrior", LevelMin: 27, LevelMax: 28, Elite: true, HPMult: 2.2, DmgMult: 1.2, Aggressive: true, AggroRange: 12, Leash: 80, RespawnSec: 9999,
		Loot: []LootEntry{l("naga_scale", 0.5, 1, 1)}})
	addMob(&Mob{ID: "naga_priest", Name: "Nagapriesterin", Model: "naga_priest", LevelMin: 28, LevelMax: 28, Elite: true, HPMult: 1.6, DmgMult: 1.1, Range: 20, Style: "bolt", Aggressive: true, AggroRange: 14, Leash: 80, RespawnSec: 9999,
		Loot: []LootEntry{l("temple_idol", 0.3, 1, 1)}})
	addMob(&Mob{ID: "temple_cobra", Name: "Tempelkobra", Model: "lizard", LevelMin: 28, LevelMax: 29, Elite: true, HPMult: 1.8, Aggressive: true, AggroRange: 12, Leash: 80, RespawnSec: 9999, Scale: 1.2,
		Loot: []LootEntry{l("serpent_fang", 0.4, 1, 1)}})
	addMob(&Mob{ID: "stone_guard", Name: "Tempelwächter", Model: "stone_golem", LevelMin: 29, LevelMax: 29, Elite: true, HPMult: 2.6, Aggressive: true, AggroRange: 10, Leash: 80, RespawnSec: 9999, Scale: 1.3, Speed: 3.8,
		Loot: []LootEntry{l("stone_core", 0.4, 1, 1)}})
	addMob(&Mob{ID: "naga_add", Name: "Naga-Leibwache", Model: "naga_warrior", LevelMin: 29, LevelMax: 29, HPMult: 0.8, Aggressive: true, AggroRange: 20, Leash: 80, RespawnSec: 9999, XPBase: 20})
	addMob(&Mob{ID: "gorun", Name: "Steinwächter Gorun", Model: "stone_golem", LevelMin: 28, Elite: true, Boss: true, Set: "ns", HPMult: 5.0, DmgMult: 1.4, Scale: 2.3, Speed: 3.6, Aggressive: true, AggroRange: 14, Leash: 90, RespawnSec: 9999,
		Abilities: []string{"heavy:9", "aoe:16", "enrage:25"}})
	addMob(&Mob{ID: "ssyra", Name: "Hohepriesterin Ssyra", Model: "naga_priest", LevelMin: 29, Elite: true, Boss: true, Set: "ns", HPMult: 5.2, DmgMult: 1.3, Scale: 1.8, Aggressive: true, AggroRange: 14, Leash: 90, RespawnSec: 9999,
		Range: 22, Style: "bolt", Abilities: []string{"poison:9", "summon:60:naga_add:3", "summon:30:naga_add:3"}})
	addMob(&Mob{ID: "zhulkarat", Name: "Nagakaiser Zhul'Karat", Model: "snake_emperor", LevelMin: 30, Elite: true, Boss: true, Final: true, Set: "ns", HPMult: 7.5, DmgMult: 1.5, Scale: 2.4, Aggressive: true, AggroRange: 16, Leash: 100, RespawnSec: 9999,
		Abilities: []string{"heavy:8", "aoe:12", "poison:14", "summon:65:naga_add:3", "summon:35:naga_add:4", "enrage:20"}})

	// ---- Dungeon: Feuerschlund (37-40) ----
	addMob(&Mob{ID: "ember_cultist", Name: "Aschekultist", Model: "cultist", LevelMin: 37, LevelMax: 38, Elite: true, HPMult: 1.6, DmgMult: 1.1, Range: 20, Style: "fire", Aggressive: true, AggroRange: 14, Leash: 80, RespawnSec: 9999,
		Loot: []LootEntry{l("ember_chip", 0.5, 1, 1)}})
	addMob(&Mob{ID: "flame_hound", Name: "Flammenhund", Model: "ashhound", LevelMin: 38, LevelMax: 38, Elite: true, HPMult: 1.9, Aggressive: true, AggroRange: 12, Leash: 80, RespawnSec: 9999,
		Loot: []LootEntry{l("hound_tooth", 0.4, 1, 1)}})
	addMob(&Mob{ID: "lava_golem", Name: "Lavagolem", Model: "magma_golem", LevelMin: 38, LevelMax: 39, Elite: true, HPMult: 2.5, Aggressive: true, AggroRange: 10, Leash: 80, RespawnSec: 9999, Scale: 1.3, Speed: 3.8,
		Loot: []LootEntry{l("molten_gem", 0.3, 1, 1)}})
	addMob(&Mob{ID: "imp_add", Name: "Glutimp", Model: "fire_imp", LevelMin: 39, LevelMax: 39, HPMult: 0.7, Aggressive: true, AggroRange: 20, Range: 16, Style: "fire", Leash: 80, RespawnSec: 9999, XPBase: 25})
	addMob(&Mob{ID: "vulkor", Name: "Schmiedemeister Vulkor", Model: "magma_golem", LevelMin: 38, Elite: true, Boss: true, Set: "gk", HPMult: 5.5, DmgMult: 1.4, Scale: 2.4, Speed: 3.6, Aggressive: true, AggroRange: 14, Leash: 90, RespawnSec: 9999,
		Abilities: []string{"heavy:8", "aoe:14", "summon:50:imp_add:4"}})
	addMob(&Mob{ID: "kaliss", Name: "Aschehexe Kaliss", Model: "cultist", LevelMin: 39, Elite: true, Boss: true, Set: "gk", HPMult: 5.5, DmgMult: 1.3, Scale: 1.8, Aggressive: true, AggroRange: 14, Leash: 90, RespawnSec: 9999,
		Range: 24, Style: "fire", Abilities: []string{"poison:8", "aoe:15", "summon:60:imp_add:3", "summon:30:imp_add:3"}})
	addMob(&Mob{ID: "ignaros", Name: "Drachenfürst Ignaros", Model: "demon_lord", LevelMin: 40, Elite: true, Boss: true, Final: true, Set: "gk", HPMult: 8.0, DmgMult: 1.5, Scale: 2.7, Aggressive: true, AggroRange: 16, Leash: 100, RespawnSec: 9999,
		Abilities: []string{"heavy:7", "aoe:11", "poison:13", "summon:70:imp_add:3", "summon:40:imp_add:4", "enrage:20"}})
}

// groups2 are the open world spawn groups of the second continent.
func groups2() []SpawnGroup {
	return []SpawnGroup{
		// Smaragddschungel
		{"river_croc", 716, 34, 7, 16}, {"river_croc", 668, 112, 4, 12},
		{"jungle_tiger", 680, -30, 8, 20}, {"jungle_tiger", 650, 70, 4, 14},
		{"panther", 700, -72, 7, 18},
		{"jungle_ape", 722, 66, 9, 18},
		{"poacher", 640, -100, 9, 14}, {"poacher_boss", 640, -96, 1, 1},
		{"lizardman", 650, 165, 9, 16}, {"lizard_shaman", 654, 168, 4, 12}, {"lizard_chief", 652, 172, 1, 1},
		{"whitetiger", 690, -46, 1, 1}, {"ape_king", 726, 72, 1, 1},
		// Kirschblütenhain
		{"fox_spirit", 880, -62, 8, 16}, {"fox_queen", 884, -66, 1, 1},
		{"bamboo_bear", 760, 52, 7, 20},
		{"ronin", 800, -150, 9, 14}, {"ronin_master", 804, -152, 1, 1},
		{"ghost_warrior", 905, 62, 8, 18},
		{"stone_golem", 850, 100, 6, 16}, {"temple_colossus", 778, -100, 1, 1},
		// Rotkluft-Canyon
		{"sandcat", 1040, 58, 8, 18}, {"canyon_bison", 960, 180, 7, 20},
		{"mesa_scorpion", 1000, 190, 9, 18}, {"scorpion_matriarch", 1006, 196, 1, 1},
		{"canyon_bandit", 1030, 165, 10, 14}, {"bandit_king", 1034, 163, 1, 1},
		{"clay_golem", 940, 108, 7, 16}, {"old_titan", 944, 112, 1, 1},
		// Aschenlande
		{"ash_hound", 900, -165, 9, 18}, {"fire_imp", 1062, -108, 9, 16},
		{"cultist", 880, -212, 10, 14}, {"flame_lord", 884, -208, 1, 1},
		{"magma_golem", 985, -105, 6, 16}, {"magma_titan", 990, -110, 1, 1},
		{"lava_elemental", 1052, -92, 6, 14},
		{"obsidian_knight", 960, -172, 6, 14}, {"obsidian_lord", 1062, -112, 1, 1},
	}
}
