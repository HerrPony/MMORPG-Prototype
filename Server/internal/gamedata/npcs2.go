package gamedata

func init() {
	NPCs = append(NPCs,
		// ---- Eichenhain ----
		&NPC{ID: "kira", Name: "Greifenmeisterin Kira", Title: "Flugmeisterin", Model: "flightmaster", X: 20, Z: 24, R: 3.9, Flight: true,
			Greet: "Mit meinen Greifen bist du schneller, als jeder Wind. Für ein paar Münzen, versteht sich."},

		// ---- Sonnenburg (Hauptstadt) ----
		&NPC{ID: "bran", Name: "Greifenmeister Bran", Title: "Flugmeister", Model: "flightmaster", X: 216, Z: -106, R: 3.9, Flight: true,
			Greet: "Willkommen in Sonnenburg! Meine Greife bringen dich in jede Stadt, die du schon kennst."},
		&NPC{ID: "alena", Name: "Bürgermeisterin Alena", Title: "Herrin von Sonnenburg", Model: "mayor", X: 190, Z: -134, R: 0,
			Greet: "Sonnenburg heißt jeden Helden willkommen."},
		&NPC{ID: "varus", Name: "Hauptmann Varus", Title: "Wache von Sonnenburg", Model: "guard", X: 205, Z: -82, R: 0.2,
			Greet: "Die Ebenen sind unsicher geworden. Bleib auf den Straßen."},
		&NPC{ID: "marius", Name: "Gladiator Marius", Title: "Kriegerlehrer", Model: "trainer_warrior", X: 172, Z: -110, R: 1.57, Trainer: "warrior",
			Greet: "Wer in der Hauptstadt kämpft, kämpft für Ruhm."},
		&NPC{ID: "ceres", Name: "Meisterschützin Ceres", Title: "Jägerlehrerin", Model: "trainer_hunter", X: 172, Z: -102, R: 1.57, Trainer: "hunter",
			Greet: "Ein ruhiger Atem, ein sicherer Pfeil."},
		&NPC{ID: "valen", Name: "Erzmagier Valen", Title: "Magierlehrer", Model: "trainer_mage", X: 172, Z: -94, R: 1.57, Trainer: "mage",
			Greet: "Magie ist Wissen, und Wissen ist Macht."},
		&NPC{ID: "rodrik", Name: "Veteran Rodrik", Title: "Meister der Krieger", Model: "trainer_warrior", X: 158, Z: -110, R: 1.57,
			Greet: "Ein Krieger, der standhält, gewinnt jede Schlacht."},
		&NPC{ID: "aldwin", Name: "Meisterjäger Aldwin", Title: "Meister der Jäger", Model: "huntsman", X: 158, Z: -102, R: 1.57,
			Greet: "Die Wildnis hat ihre eigenen Gesetze. Ich kenne sie alle."},
		&NPC{ID: "ilvane", Name: "Erzmagierin Ilvane", Title: "Meisterin der Magier", Model: "trainer_mage", X: 158, Z: -94, R: 1.57,
			Greet: "Das Wasser vergisst nie, wer es um Hilfe bittet."},
		&NPC{ID: "thoren", Name: "Waffenmeister Thoren", Title: "Waffenhändler", Model: "smith", X: 208, Z: -110, R: 4.7, Vendor: true,
			StockK: "weapons", StockLo: 2, StockHi: 6, Greet: "Scharfes Eisen für scharfe Köpfe."},
		&NPC{ID: "isolde", Name: "Rüstmeisterin Isolde", Title: "Rüstungshändlerin", Model: "merchant", X: 208, Z: -102, R: 4.7, Vendor: true,
			StockK: "armor", StockLo: 2, StockHi: 6, Greet: "Die beste Rüstung der Ebenen."},
		&NPC{ID: "pip", Name: "Alchemist Pip", Title: "Tränke & Schmuck", Model: "merchant", X: 200, Z: -92, R: 3.1, Vendor: true,
			StockK: "jewelry", StockLo: 2, StockHi: 5, Stock: []string{"potion_hp2", "potion_hp3"}, Greet: "Tränke, Tränke, Tränke! Und etwas Glitzerkram."},
		&NPC{ID: "odo", Name: "Bauer Odo", Title: "Weizenhof", Model: "farmer", X: 140, Z: -60, R: 1.0,
			Greet: "Gutes Land, aber schwere Arbeit."},
		&NPC{ID: "hannes", Name: "Müller Hannes", Title: "Windmühle", Model: "farmer", X: 100, Z: -160, R: 2.0,
			Greet: "Ohne Mühle kein Brot."},
		&NPC{ID: "selene", Name: "Jägerin Selene", Title: "Jägerhain", Model: "huntsman", X: 130, Z: -212, R: 1.0,
			Greet: "Still. Du verscheuchst das Wild."},
		&NPC{ID: "ilyana", Name: "Priesterin Ilyana", Title: "Hüterin der Quelle", Model: "priest", X: 176, Z: -222, R: 2.0,
			Greet: "Die Quelle ist heilig. Behandle sie mit Respekt."},

		// ---- Moorbruck ----
		&NPC{ID: "mira", Name: "Greifenmeisterin Mira", Title: "Flugmeisterin", Model: "flightmaster", X: -154, Z: 180, R: 4.0, Flight: true,
			Greet: "Der Himmel über dem Moor ist klarer als der Boden darunter."},
		&NPC{ID: "moros", Name: "Dorfältester Moros", Title: "Ältester von Moorbruck", Model: "elder", X: -172, Z: 160, R: 0.3,
			Greet: "Moorbruck ist klein, aber stolz."},
		&NPC{ID: "sylva", Name: "Kräuterfrau Sylva", Title: "Heilkundige", Model: "herbalist", X: -186, Z: 176, R: 1.2,
			Greet: "Jedes Kraut hat seinen Zweck."},
		&NPC{ID: "hektor", Name: "Wachmann Hektor", Title: "Moorwacht", Model: "guard", X: -157, Z: 158, R: 5.0,
			Greet: "Nachts ist es hier nicht ruhig."},
		&NPC{ID: "tam", Name: "Fischer Tam", Title: "Nebelsee", Model: "fisher", X: -152, Z: 202, R: 4.0,
			Greet: "Der Nebelsee gibt viel her. Wenn man ihn lässt."},
		&NPC{ID: "boggs", Name: "Händler Boggs", Title: "Händler", Model: "merchant", X: -182, Z: 164, R: 1.0, Vendor: true,
			StockK: "weapons,armor", StockLo: 5, StockHi: 7, Stock: []string{"potion_hp2", "potion_hp3"}, Greet: "Gute Ware, krumme Preise."},

		// ---- Frostwacht ----
		&NPC{ID: "sigrid", Name: "Greifenmeisterin Sigrid", Title: "Flugmeisterin", Model: "flightmaster", X: -94, Z: -205, R: 4.0, Flight: true,
			Greet: "Bei diesem Wind fliegt kein Greif gern. Aber sie tun es."},
		&NPC{ID: "isgard", Name: "Hauptmann Isgard", Title: "Frostwacht-Garde", Model: "guard", X: -112, Z: -207, R: 0.3,
			Greet: "Die Berge sind unser Schutz und unser Feind."},
		&NPC{ID: "olaf", Name: "Jäger Olaf", Title: "Bergjäger", Model: "huntsman", X: -121, Z: -224, R: 0.8,
			Greet: "Im Schnee hinterlassen alle Spuren."},
		&NPC{ID: "hilde", Name: "Schmiedin Hilde", Title: "Waffen & Rüstungen", Model: "smith", X: -100, Z: -224, R: 3.5, Vendor: true,
			StockK: "weapons,armor", StockLo: 7, StockHi: 9, Stock: []string{"potion_hp3", "potion_hp4"}, Greet: "Frostgeschmiedet hält länger."},
		&NPC{ID: "vorn", Name: "Gelehrter Vorn", Title: "Eisforscher", Model: "scholar", X: -60, Z: -172, R: 1.0,
			Greet: "Fascinierend! Die Kälte erzeugt geistähnliche Strukturen."},
		&NPC{ID: "dolf", Name: "Grubenaufseher Dolf", Title: "Schwarzfelsmine", Model: "guard", X: -166, Z: -190, R: 4.5,
			Greet: "Die alte Mine gehört den Kobolden. Wir wollen sie zurück."},

		// ---- Sandfurt ----
		&NPC{ID: "zara", Name: "Greifenmeisterin Zara", Title: "Flugmeisterin", Model: "flightmaster", X: 206, Z: 200, R: 3.9, Flight: true,
			Greet: "Sand im Gefieder, aber die Greife lieben die Hitze."},
		&NPC{ID: "karim", Name: "Scheich Karim", Title: "Herr von Sandfurt", Model: "sheikh", X: 188, Z: 182, R: 0.3,
			Greet: "Die Wüste vergibt keine Fehler, Fremder."},
		&NPC{ID: "jasim", Name: "Kameltreiber Jasim", Title: "Karawane", Model: "nomad", X: 176, Z: 196, R: 1.2,
			Greet: "Meine Kamele sind die schnellsten im ganzen Land."},
		&NPC{ID: "bahir", Name: "Wache Bahir", Title: "Dünenwacht", Model: "guard", X: 240, Z: 138, R: 3.5,
			Greet: "Skorpione, Räuber, Sonne. Such dir etwas aus."},
		&NPC{ID: "neferet", Name: "Archäologin Neferet", Title: "Grabforscherin", Model: "scholar", X: 135, Z: 230, R: 0.5,
			Greet: "Dieses Grab ist tausende Jahre alt. Und nicht ganz tot."},
		&NPC{ID: "rashid", Name: "Händler Rashid", Title: "Händler", Model: "merchant", X: 200, Z: 184, R: 4.4, Vendor: true,
			StockK: "weapons,armor,jewelry", StockLo: 8, StockHi: 9, Stock: []string{"potion_hp3", "potion_hp4"}, Greet: "Für dich, mein Freund, ein besonderer Preis!"},
	)
}

// NPCByID returns an NPC by id.
func NPCByID(id string) *NPC {
	for _, n := range NPCs {
		if n.ID == id {
			return n
		}
	}
	return nil
}
