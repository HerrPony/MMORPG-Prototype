package gamedata

// Interactive objects and the two new dungeons of the second continent.

var extraObjectSpawns = []ObjectSpawn{
	{"portal_temple", 768, -112},
	{"portal_fire", 1000, -198},
}

func init() {
	for _, o := range []*ObjectDef{
		{ID: "orchid_flower", Name: "Mondorchidee", Model: "lotus", Item: "jungle_orchid", Quest: "j_orchids", Respawn: 45},
		{ID: "bamboo_stalks", Name: "Bambusstängel", Model: "wheat", Item: "bamboo_bundle", Quest: "j_bamboo", Respawn: 45},
		{ID: "incense_burner", Name: "Weihrauchschale", Model: "urn", Item: "temple_incense", Quest: "b_incense", Respawn: 45},
		{ID: "sakura_pile", Name: "Kirschblütenhaufen", Model: "lotus", Item: "sakura_petal", Quest: "b_sakura", Respawn: 45},
		{ID: "bamboo_shoots", Name: "Bambussprossen", Model: "wheat", Item: "bamboo_shoot", Quest: "k_tea", Respawn: 45},
		{ID: "clay_pit", Name: "Tonvorkommen", Model: "urn", Item: "red_clay", Quest: "m_clay", Respawn: 45},
		{ID: "mesa_crystal_node", Name: "Kluftkristall", Model: "crystal", Item: "mesa_crystal", Quest: "m_crystals", Respawn: 45},
		{ID: "ember_crystal_node", Name: "Glutkristall", Model: "crystal", Item: "ember_crystal", Quest: "v_crystals", Respawn: 45},
		{ID: "fire_flower_plant", Name: "Glutblüte", Model: "cactus", Item: "fire_flower", Quest: "v_flowers", Respawn: 45},
		{ID: "portal_temple", Name: "Portal zum Nebeltempel", Model: "portal", Portal: "dungeon:temple"},
		{ID: "portal_fire", Name: "Portal zum Feuerschlund", Model: "portal", Portal: "dungeon:fire"},
	} {
		Objects[o.ID] = o
	}
	objGroups = append(objGroups,
		objGroup{"orchid_flower", 700, 12, 10, 26},
		objGroup{"bamboo_stalks", 688, 118, 10, 28},
		objGroup{"incense_burner", 880, -62, 7, 14},
		objGroup{"sakura_pile", 850, 60, 11, 28},
		objGroup{"bamboo_shoots", 760, 52, 9, 22},
		objGroup{"clay_pit", 1010, 150, 9, 22},
		objGroup{"mesa_crystal_node", 1040, 58, 7, 20},
		objGroup{"ember_crystal_node", 1062, -108, 9, 22},
		objGroup{"fire_flower_plant", 900, -165, 9, 22},
	)

	Dungeons["temple"] = &Dungeon{
		ID: "temple", Name: "Der Versunkene Nebeltempel", MinLevel: 26, MaxLevel: 30, Style: "temple", Origin: [2]float64{0, DungeonZ0}, Entrance: [2]float64{768, -112},
		Start: [2]float64{0, 3},
		Rooms: []Rect{
			{0, 0, 26, 20, "temple"}, {0, -20, 6, 20, "corridor"}, {0, -42, 34, 24, "flooded"}, {0, -60, 6, 12, "corridor"},
			{0, -84, 36, 36, "temple"}, {0, -108, 6, 12, "corridor"}, {0, -132, 40, 36, "flooded"}, {0, -156, 6, 12, "corridor"},
			{0, -182, 38, 40, "temple"}, {0, -208, 6, 12, "corridor"}, {0, -238, 46, 48, "shrine"},
		},
		Spawns: []DungeonSpawn{
			// first flooded hall
			{"naga_warrior", -9, -38}, {"naga_warrior", 0, -35}, {"naga_warrior", 9, -38}, {"naga_priest", -4, -48}, {"naga_priest", 5, -49},
			// Steinwächter Gorun
			{"gorun", 0, -88}, {"stone_guard", -12, -74}, {"stone_guard", 12, -74},
			// second flooded hall
			{"temple_cobra", -11, -124}, {"temple_cobra", 11, -124}, {"temple_cobra", 0, -144}, {"naga_warrior", -12, -140}, {"naga_warrior", 12, -140}, {"naga_priest", 0, -133},
			// Hohepriesterin Ssyra
			{"temple_cobra", -13, -170}, {"temple_cobra", 13, -170}, {"naga_priest", -8, -186}, {"naga_priest", 8, -186}, {"ssyra", 0, -192},
			// Nagakaiser Zhul'Karat
			{"zhulkarat", 0, -248},
		},
	}
	Dungeons["fire"] = &Dungeon{
		ID: "fire", Name: "Der Feuerschlund", MinLevel: 36, MaxLevel: 40, Style: "fire", Origin: [2]float64{0, DungeonZ0}, Entrance: [2]float64{1000, -198},
		Start: [2]float64{0, 3},
		Rooms: []Rect{
			{0, 0, 26, 20, "forge"}, {0, -20, 6, 20, "corridor"}, {0, -42, 34, 24, "lava"}, {0, -60, 6, 12, "corridor"},
			{0, -84, 36, 36, "forge"}, {0, -108, 6, 12, "corridor"}, {0, -132, 40, 36, "lava"}, {0, -156, 6, 12, "corridor"},
			{0, -182, 38, 40, "obsidian"}, {0, -208, 6, 12, "corridor"}, {0, -238, 46, 48, "inferno"},
		},
		Spawns: []DungeonSpawn{
			{"ember_cultist", -9, -38}, {"ember_cultist", 0, -35}, {"ember_cultist", 9, -38}, {"flame_hound", -4, -48}, {"flame_hound", 5, -49},
			// Schmiedemeister Vulkor
			{"vulkor", 0, -88}, {"lava_golem", -12, -74}, {"lava_golem", 12, -74},
			{"flame_hound", -11, -124}, {"flame_hound", 11, -124}, {"lava_golem", 0, -144}, {"ember_cultist", -12, -140}, {"ember_cultist", 12, -140}, {"ember_cultist", 0, -133},
			// Aschehexe Kaliss
			{"lava_golem", -13, -170}, {"lava_golem", 13, -170}, {"ember_cultist", -8, -186}, {"ember_cultist", 8, -186}, {"kaliss", 0, -192},
			// Drachenfürst Ignaros
			{"ignaros", 0, -248},
		},
	}
	// the first dungeon gets the same data fields
	Dungeons["mine"].MaxLevel = 20
	Dungeons["mine"].Style = "mine"
}
