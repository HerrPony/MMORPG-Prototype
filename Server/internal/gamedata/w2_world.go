package gamedata

// The second continent "Jadereich": places and decorations.

func init() {
	POIs = append(POIs,
		// Smaragddschungel (20-25)
		&POI{ID: "sturmhafen", Name: "Sturmhafen", X: 232, Z: 40, R: 36, Kind: "town", Town: true},
		&POI{ID: "lotushafen", Name: "Lotushafen", X: 612, Z: 20, R: 34, Kind: "town", Town: true},
		&POI{ID: "tigerpfad", Name: "Tigerpfad", X: 680, Z: -30, R: 24, Kind: "wild"},
		&POI{ID: "affenhain", Name: "Affenhain", X: 722, Z: 66, R: 24, Kind: "wild"},
		&POI{ID: "bambushain", Name: "Bambushain", X: 688, Z: 118, R: 32, Kind: "town", Town: true},
		&POI{ID: "lotusteich", Name: "Lotusteich", X: 730, Z: 160, R: 22, Kind: "lake"},
		&POI{ID: "wildererlager", Name: "Wildererlager", X: 640, Z: -100, R: 20, Kind: "camp"},
		&POI{ID: "echsenruinen", Name: "Echsenruinen", X: 650, Z: 165, R: 22, Kind: "ruin"},
		// Kirschblütenhain (25-30)
		&POI{ID: "himmelspforte", Name: "Himmelspforte", X: 830, Z: -10, R: 78, Kind: "town", Town: true},
		&POI{ID: "bambuswald", Name: "Bambuswald", X: 760, Z: 52, R: 24, Kind: "wild"},
		&POI{ID: "fuchsschrein", Name: "Fuchsschrein", X: 880, Z: -62, R: 18, Kind: "ruin"},
		&POI{ID: "spiegelteich", Name: "Spiegelteich", X: 905, Z: 62, R: 20, Kind: "lake"},
		&POI{ID: "kirschbach", Name: "Kirschbach", X: 810, Z: 150, R: 32, Kind: "town", Town: true},
		&POI{ID: "roninlager", Name: "Ronin-Lager", X: 800, Z: -150, R: 22, Kind: "camp"},
		&POI{ID: "nebeltempel", Name: "Nebeltempel", X: 768, Z: -112, R: 16, Kind: "cave"},
		// Rotkluft-Canyon (30-35)
		&POI{ID: "rotfels_dorf", Name: "Rotfels", X: 1000, Z: 120, R: 32, Kind: "town", Town: true},
		&POI{ID: "kluftpass", Name: "Kluftpass", X: 960, Z: 180, R: 24, Kind: "wild"},
		&POI{ID: "steinbogen", Name: "Steinbogen", X: 1040, Z: 58, R: 22, Kind: "wild"},
		&POI{ID: "titanengrab", Name: "Titanengrab", X: 940, Z: 108, R: 20, Kind: "ruin"},
		&POI{ID: "banditenfestung", Name: "Banditenfestung", X: 1030, Z: 165, R: 22, Kind: "camp"},
		&POI{ID: "versteckte_oase", Name: "Verborgene Oase", X: 1045, Z: 150, R: 14, Kind: "lake"},
		// Aschenlande (35-40)
		&POI{ID: "aschenwacht", Name: "Aschenwacht", X: 940, Z: -100, R: 32, Kind: "town", Town: true},
		&POI{ID: "aschepass", Name: "Aschepass", X: 900, Z: -165, R: 24, Kind: "wild"},
		&POI{ID: "kultistenhain", Name: "Hain der Glutkultisten", X: 880, Z: -212, R: 20, Kind: "camp"},
		&POI{ID: "schlund", Name: "Der Schlund", X: 1010, Z: -140, R: 26, Kind: "wild"},
		&POI{ID: "obsidianfeld", Name: "Obsidianfeld", X: 1062, Z: -108, R: 22, Kind: "wild"},
		&POI{ID: "feuerschlund", Name: "Feuerschlund", X: 1000, Z: -198, R: 16, Kind: "cave"},
	)
}

func extraBuildings2() []Building {
	return []Building{
		// wildererlager
		{Kind: "tent", X: 636, Z: -96, W: 5, D: 5, H: 3.5}, {Kind: "tent", X: 646, Z: -104, W: 5, D: 5, H: 3.5}, {Kind: "tent", X: 640, Z: -88, W: 5, D: 5, H: 3.5},
		// echsenruinen
		{Kind: "tower", X: 646, Z: 160, W: 3, D: 3, H: 7}, {Kind: "tower", X: 658, Z: 170, W: 3, D: 3, H: 5}, {Kind: "tower", X: 652, Z: 176, W: 3, D: 3, H: 9},
		// fuchsschrein
		{Kind: "gate", Style: "jade", X: 880, Z: -58, W: 9, D: 3, H: 7},
		{Kind: "pagoda", Style: "jade", X: 880, Z: -70, W: 8, D: 8, H: 9},
		// ronin camp
		{Kind: "tent", X: 796, Z: -146, W: 6, D: 6, H: 4}, {Kind: "tent", X: 806, Z: -156, W: 6, D: 6, H: 4}, {Kind: "tent", X: 802, Z: -138, W: 6, D: 6, H: 4},
		// nebeltempel entrance (a sunken temple door in the hillside)
		{Kind: "cave", Style: "temple", X: 768, Z: -118, W: 14, D: 6, H: 10},
		// titanengrab
		{Kind: "keep", Style: "redrock", X: 940, Z: 112, W: 16, D: 12, H: 7},
		// banditenfestung
		{Kind: "tent", X: 1026, Z: 162, W: 6, D: 6, H: 4}, {Kind: "tent", X: 1038, Z: 170, W: 6, D: 6, H: 4}, {Kind: "tent", X: 1034, Z: 156, W: 6, D: 6, H: 4},
		{Kind: "tower", X: 1020, Z: 172, W: 4, D: 4, H: 8}, {Kind: "tower", X: 1042, Z: 158, W: 4, D: 4, H: 8},
		// kultistenhain
		{Kind: "tent", X: 876, Z: -208, W: 6, D: 6, H: 4}, {Kind: "tent", X: 886, Z: -216, W: 6, D: 6, H: 4}, {Kind: "hut", Style: "obsidian", X: 884, Z: -204, W: 6, D: 6, H: 4},
		// feuerschlund entrance
		{Kind: "cave", Style: "fire", X: 1000, Z: -204, W: 14, D: 6, H: 10},
	}
}
