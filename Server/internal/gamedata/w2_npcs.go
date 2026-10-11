package gamedata

// NPCs of the harbor Sturmhafen and of the second continent.

func init() {
	NPCs = append(NPCs,
		// ---- Sturmhafen (Erlenmark) ----
		&NPC{ID: "brenner", Name: "Kapitän Brenner", Title: "Fährmeister", Model: "captain", X: 247, Z: 40, R: 1.57, Ship: true,
			Greet: "Die Überfahrt nach Jadereich dauert nicht lang. Sofern der Sturm es zulässt."},
		&NPC{ID: "olga", Name: "Hafenmeisterin Olga", Title: "Hafenmeisterin", Model: "guard", X: 238, Z: 50, R: 0.5,
			Greet: "Jeden Tag laufen Schiffe aus. Manche kommen sogar zurück."},

		// ---- Lotushafen (Jadereich, Dschungel) ----
		&NPC{ID: "wen", Name: "Kapitän Wen", Title: "Fährmeister", Model: "captain", X: 597, Z: 20, R: -1.57, Ship: true,
			Greet: "Willkommen in Jadereich! Für ein paar Münzen bringe ich dich zurück nach Sturmhafen."},
		&NPC{ID: "lin", Name: "Hafenmeisterin Lin", Title: "Hafenmeisterin von Lotushafen", Model: "elder", X: 606, Z: 12, R: 0.5,
			Greet: "Jede Kiste, die hier ankommt, geht durch meine Hände."},
		&NPC{ID: "mei", Name: "Fischerin Mei", Title: "Fischerin", Model: "fisher", X: 620, Z: 32, R: 1.0,
			Greet: "Die Flüsse hier sind voller Leben. Und voller Zähne."},
		&NPC{ID: "bao", Name: "Händler Bao", Title: "Händler", Model: "merchant", X: 620, Z: 8, R: 3.3, Vendor: true,
			StockK: "weapons,armor", StockLo: 10, StockHi: 11, Stock: []string{"potion_hp4", "potion_hp5"}, Greet: "Für Reisende vom Festland nur das Beste!"},
		&NPC{ID: "jiro", Name: "Wache Jiro", Title: "Hafenwache", Model: "samurai", X: 628, Z: 24, R: 1.57,
			Greet: "Der Dschungel beginnt gleich hinter dem Tor. Geh nicht ohne Waffe."},

		// ---- Bambushain ----
		&NPC{ID: "suki", Name: "Greifenmeisterin Suki", Title: "Flugmeisterin", Model: "flightmaster", X: 704, Z: 128, R: 3.9, Flight: true,
			Greet: "Auf Greifenflügeln siehst du den ganzen Dschungel auf einmal."},
		&NPC{ID: "yuna", Name: "Dorfälteste Yuna", Title: "Älteste von Bambushain", Model: "elder", X: 686, Z: 112, R: 0.4,
			Greet: "Bambushain steht, solange der Bambus wächst."},
		&NPC{ID: "taro", Name: "Jäger Taro", Title: "Dschungeljäger", Model: "huntsman", X: 672, Z: 126, R: 0.8,
			Greet: "Im Dschungel jagt man nicht. Man wird gejagt."},
		&NPC{ID: "hana", Name: "Heilerin Hana", Title: "Kräuterheilerin", Model: "herbalist", X: 696, Z: 104, R: 1.2, Vendor: true,
			Stock: []string{"potion_hp4", "potion_hp5"}, Greet: "Die Natur heilt, wenn man sie lässt."},
		&NPC{ID: "kenzo", Name: "Händler Kenzo", Title: "Händler", Model: "merchant", X: 680, Z: 134, R: 4.7, Vendor: true,
			StockK: "weapons,armor,jewelry", StockLo: 11, StockHi: 12, Stock: []string{"potion_hp5"}, Greet: "Bambus, Seide, Stahl – alles da."},

		// ---- Himmelspforte (Hauptstadt von Jadereich) ----
		&NPC{ID: "haru", Name: "Greifenmeister Haru", Title: "Flugmeister", Model: "flightmaster", X: 860, Z: 10, R: 3.9, Flight: true,
			Greet: "Von der Himmelspforte aus erreichst du jeden Winkel des Reiches."},
		&NPC{ID: "meilin", Name: "Kaiserin Mei-Lin", Title: "Herrscherin von Jadereich", Model: "empress", X: 830, Z: -24, R: 0,
			Greet: "Jadereich heißt jeden willkommen, der Frieden bringt."},
		&NPC{ID: "akane", Name: "Generalin Akane", Title: "Oberbefehlshaberin", Model: "samurai", X: 806, Z: -2, R: 1.57,
			Greet: "Das Reich braucht Klingen, keine Worte."},
		&NPC{ID: "goro", Name: "Meister Goro", Title: "Kriegerlehrer", Model: "trainer_warrior", X: 800, Z: -26, R: 1.57, Trainer: "warrior",
			Greet: "Ein Krieger kämpft wie der Berg: unbeweglich und unaufhaltsam."},
		&NPC{ID: "sora", Name: "Meisterin Sora", Title: "Jägerlehrerin", Model: "trainer_hunter", X: 800, Z: -19, R: 1.57, Trainer: "hunter",
			Greet: "Der Pfeil kennt sein Ziel, bevor er den Bogen verlässt."},
		&NPC{ID: "ryu", Name: "Großmeister Ryu", Title: "Magierlehrer", Model: "trainer_mage", X: 800, Z: -12, R: 1.57, Trainer: "mage",
			Greet: "Die Elemente gehorchen dem, der sie versteht."},
		&NPC{ID: "daichi", Name: "Waffenmeister Daichi", Title: "Waffenschmied", Model: "smith", X: 862, Z: -26, R: 4.7, Vendor: true,
			StockK: "weapons", StockLo: 12, StockHi: 14, Greet: "Jede Klinge hat ihre Seele."},
		&NPC{ID: "hotaru", Name: "Rüstmeisterin Hotaru", Title: "Rüstungsschmiedin", Model: "merchant", X: 862, Z: -18, R: 4.7, Vendor: true,
			StockK: "armor", StockLo: 12, StockHi: 14, Greet: "Seide, Leder und Stahl."},
		&NPC{ID: "yumi", Name: "Alchemistin Yumi", Title: "Tränke & Schmuck", Model: "merchant", X: 852, Z: -32, R: 3.1, Vendor: true,
			StockK: "jewelry", StockLo: 12, StockHi: 15, Stock: []string{"potion_hp5", "potion_hp6"}, Greet: "Tränke aus tausend Kräutern."},
		&NPC{ID: "genji", Name: "Mönch Genji", Title: "Hüter des Schreins", Model: "monk", X: 846, Z: -26, R: 3.5,
			Greet: "Stille ist die lauteste Antwort."},
		&NPC{ID: "takumi", Name: "Gelehrter Takumi", Title: "Kaiserlicher Gelehrter", Model: "scholar", X: 814, Z: -26, R: 1.2,
			Greet: "Wissen ist das wertvollste Gut des Reiches."},
		&NPC{ID: "miyu", Name: "Archäologin Miyu", Title: "Nebeltempel", Model: "scholar", X: 762, Z: -104, R: 0.6,
			Greet: "Der Tempel liegt tief unter dem Wasser. Und er ist nicht leer."},

		// ---- Kirschbach ----
		&NPC{ID: "nori", Name: "Greifenmeisterin Nori", Title: "Flugmeisterin", Model: "flightmaster", X: 826, Z: 160, R: 3.9, Flight: true,
			Greet: "Kirschblüten und Greifenfedern, ein hübscher Anblick."},
		&NPC{ID: "hiro", Name: "Bauer Hiro", Title: "Bambusbauer", Model: "farmer", X: 808, Z: 140, R: 1.0,
			Greet: "Der Bambus wächst gut. Die Bären auch."},
		&NPC{ID: "aiko", Name: "Teemeisterin Aiko", Title: "Teemeisterin", Model: "herbalist", X: 798, Z: 158, R: 1.2,
			Greet: "Tee ist Geduld im Becher."},
		&NPC{ID: "rin", Name: "Füchsin Rin", Title: "Fuchsjägerin", Model: "huntsman", X: 818, Z: 138, R: 0.8,
			Greet: "Füchse sind klüger, als sie aussehen."},
		&NPC{ID: "sen", Name: "Händler Sen", Title: "Händler", Model: "merchant", X: 796, Z: 148, R: 4.7, Vendor: true,
			StockK: "weapons,armor", StockLo: 13, StockHi: 14, Stock: []string{"potion_hp5"}, Greet: "Frische Ware aus der Hauptstadt!"},

		// ---- Rotfels (Dorf) ----
		&NPC{ID: "dax", Name: "Greifenmeister Dax", Title: "Flugmeister", Model: "flightmaster", X: 1016, Z: 130, R: 3.9, Flight: true,
			Greet: "Die Canyons sind tückisch. Aus der Luft siehst du sie besser."},
		&NPC{ID: "mara", Name: "Sheriff Mara", Title: "Sheriff von Rotfels", Model: "guard", X: 998, Z: 112, R: 0.3,
			Greet: "In der Schlucht gelten meine Gesetze."},
		&NPC{ID: "brokk", Name: "Schmied Brokk", Title: "Waffen & Rüstungen", Model: "smith", X: 988, Z: 128, R: 3.5, Vendor: true,
			StockK: "weapons,armor", StockLo: 15, StockHi: 16, Greet: "Roter Stahl für rote Felsen."},
		&NPC{ID: "olin", Name: "Schürfer Olin", Title: "Minenarbeiter", Model: "farmer", X: 1010, Z: 108, R: 1.0,
			Greet: "Die Erde gibt viel her, wenn man sie nicht stört."},
		&NPC{ID: "zane", Name: "Händler Zane", Title: "Händler", Model: "merchant", X: 992, Z: 134, R: 4.4, Vendor: true,
			StockK: "jewelry", StockLo: 15, StockHi: 16, Stock: []string{"potion_hp5", "potion_hp6"}, Greet: "Schmuck, Tränke, Nützliches."},
		&NPC{ID: "nox", Name: "Späher Nox", Title: "Steinbogen-Posten", Model: "huntsman", X: 1038, Z: 62, R: 3.0,
			Greet: "Von hier sehe ich jeden, der die Schlucht betritt."},

		// ---- Aschenwacht ----
		&NPC{ID: "ember", Name: "Greifenmeisterin Ember", Title: "Flugmeisterin", Model: "flightmaster", X: 956, Z: -90, R: 3.9, Flight: true,
			Greet: "Die Asche macht meinen Greifen nichts aus. Dir wohl schon."},
		&NPC{ID: "vash", Name: "Kommandantin Vash", Title: "Kommandantin der Aschenwacht", Model: "samurai", X: 938, Z: -108, R: 0.3,
			Greet: "Die Aschenwacht hält die Linie. Immer."},
		&NPC{ID: "kael", Name: "Alchemist Kael", Title: "Alchemist", Model: "scholar", X: 928, Z: -94, R: 1.2, Vendor: true,
			StockK: "jewelry", StockLo: 17, StockHi: 19, Stock: []string{"potion_hp6"}, Greet: "Feuer kann heilen und töten. Ich verkaufe beides."},
		&NPC{ID: "ignis", Name: "Schmiedin Ignis", Title: "Waffen & Rüstungen", Model: "smith", X: 952, Z: -110, R: 3.5, Vendor: true,
			StockK: "weapons,armor", StockLo: 17, StockHi: 19, Greet: "Im Vulkanfeuer geschmiedet. Es gibt nichts Besseres."},
		&NPC{ID: "rurik", Name: "Späher Rurik", Title: "Aschepass-Wache", Model: "huntsman", X: 898, Z: -160, R: 1.0,
			Greet: "Ich sehe die Kultisten kommen, lange bevor sie mich sehen."},
		&NPC{ID: "karr", Name: "Feldwebel Karr", Title: "Feuerschlund-Wache", Model: "guard", X: 996, Z: -190, R: 0.4,
			Greet: "Im Schlund brennt es seit Jahrhunderten. Jemand hält es am Brennen."},
	)
}
