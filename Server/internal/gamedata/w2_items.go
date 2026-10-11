package gamedata

// Items of the second continent: vendor goods, quest items and dungeon loot.

func init() {
	junk := func(id, name, icon string, sell int, desc string) {
		Items[id] = &Item{ID: id, Name: name, Icon: icon, Rarity: "poor", Desc: desc, Sell: sell, Stack: 20}
	}
	good := func(id, name, icon string, sell int, desc string) {
		Items[id] = &Item{ID: id, Name: name, Icon: icon, Rarity: "common", Desc: desc, Sell: sell, Stack: 20}
	}
	quest := func(id, name, icon, desc string) {
		Items[id] = &Item{ID: id, Name: name, Icon: icon, Rarity: "common", Desc: desc, Stack: 20, Quest: true}
	}

	// ---- Smaragddschungel ----
	good("tiger_pelt", "Tigerfell", "🐅", 160, "Gestreift und erstaunlich weich.")
	good("tiger_fang", "Tigerzahn", "🦷", 120, "Länger als dein Finger.")
	good("panther_claw", "Pantherkralle", "🐾", 150, "Schwarz und messerscharf.")
	junk("ape_fur", "Affenfell", "🧶", 90, "Riecht nach Dschungel.")
	junk("poacher_trap", "Wildererfalle", "🪤", 110, "Rostig, aber noch funktionstüchtig.")
	good("lizard_scale", "Echsenschuppe", "🦎", 170, "Schimmert in tausend Grüntönen.")
	good("river_croc_tooth", "Krokodilzahn", "🐊", 140, "Ein beeindruckender Zahn.")
	good("snake_venom", "Schlangengift", "🧪", 160, "Eine Phiole mit schimmerndem Gift.")
	quest("jungle_orchid", "Mondorchidee", "🌺", "Blüht nur im Schatten der großen Bäume.")
	quest("tiger_trophy", "Tigerkralle", "🐾", "Beweis für den Sieg über den Weißen Tiger.")
	quest("poacher_ledger", "Wilderer-Kassenbuch", "📒", "Voller Namen und Summen.")
	quest("lizard_idol", "Echsenidol", "🗿", "Ein kleines, grün glasiertes Götzenbild.")
	quest("letter_wen", "Hafenbrief", "✉️", "Ein versiegelter Brief für das Festland.")
	quest("tea_box", "Teekiste", "📦", "Duftet nach Jasmin. Muss frisch ankommen!")
	quest("ship_manifest", "Frachtbrief", "📜", "Eine Liste der Ladung, mit dem Siegel des Hafenmeisters.")
	quest("bamboo_bundle", "Bambusbündel", "🎋", "Frisch geschlagen.")

	// ---- Kirschblütenhain ----
	good("fox_tail", "Fuchsschweif", "🦊", 230, "Er leuchtet leicht im Dunkeln.")
	good("panda_fur", "Pandafell", "🐼", 250, "Schwarz-weiß und flauschig.")
	junk("ronin_sword", "Verbogene Katana", "⚔️", 260, "Hat schon bessere Tage gesehen.")
	good("stone_core", "Steinkern", "🪨", 280, "Warm, als schliefe etwas darin.")
	good("spirit_essence", "Geisteressenz", "👻", 270, "Ein kalter, leuchtender Nebel in einem Glas.")
	quest("bamboo_shoot", "Bambussprosse", "🌱", "Zart und frisch, ein Leckerbissen.")
	quest("sakura_petal", "Kirschblüte", "🌸", "Rosa Blüten, vom Wind getragen.")
	quest("temple_incense", "Tempelweihrauch", "🪔", "Duftet nach Zedern und Ruhe.")
	quest("fox_pearl", "Fuchsperle", "🔮", "Ein Juwel, das Füchse sammeln.")
	quest("ronin_seal", "Ronin-Siegel", "🈴", "Das Zeichen eines Verstoßenen.")
	quest("spirit_bell", "Geisterglocke", "🔔", "Sie klingelt, ohne dass jemand sie berührt.")
	quest("spirit_scroll", "Geisterrolle", "📜", "Beschrieben mit silberner Tinte. Dringend!")
	quest("sake_flask", "Reiswein-Flasche", "🍶", "Ein Geschenk an die Wachen. Nicht schütteln!")

	// ---- Rotkluft-Canyon ----
	good("sandcat_fang", "Sandkatzenzahn", "🐈", 340, "Gelblich und scharf.")
	good("bison_horn", "Bisonhorn", "🐃", 370, "Dick wie ein Oberarm.")
	junk("bandit_badge", "Schluchtbanditen-Abzeichen", "🎖️", 330, "Ein rostiges Blechzeichen.")
	good("golem_heart", "Golemherz", "💎", 420, "Pulsiert langsam im Takt der Erde.")
	good("scorpion_chitin", "Skorpionpanzer", "🦂", 350, "Hart wie Stahl.")
	quest("red_clay", "Roter Ton", "🧱", "Feucht und rot, ideal für Töpfer.")
	quest("mesa_crystal", "Kluftkristall", "💠", "Fängt das Sonnenlicht ein.")
	quest("mesa_map", "Schluchtkarte", "🗺️", "Zeigt geheime Wege durch den Canyon.")
	quest("titan_shard", "Titanenscherbe", "🪨", "Ein Stück eines uralten Riesen.")
	quest("ore_sample", "Erzprobe", "🔶", "Eine Probe für den Schmied.")
	quest("pottery_crate", "Tonwaren-Kiste", "🏺", "Zerbrechlich! Schnell und vorsichtig liefern.")

	// ---- Aschenlande ----
	good("hound_tooth", "Glutzahn", "🔥", 480, "Noch warm.")
	good("imp_horn", "Imphorn", "😈", 460, "Kleine, rußige Hörner.")
	good("magma_core", "Magmakern", "🌋", 560, "Brennt sich durch deine Tasche, wenn du ihn zu lange hältst.")
	junk("cultist_scrap", "Kultistenrobe-Fetzen", "🧵", 440, "Verbrannt und eingerissen.")
	good("obsidian_shard", "Obsidiansplitter", "⚫", 520, "Messerscharf und glänzend schwarz.")
	quest("ember_crystal", "Glutkristall", "🔥", "Heiß, aber nicht verbrennend.")
	quest("fire_flower", "Glutblüte", "🏵️", "Blüht nur auf Lavaboden.")
	quest("cult_tome", "Kultbuch", "📕", "Seiten aus verkohltem Leder, voller Runen.")
	quest("lava_heart", "Lavaherz", "❤️‍🔥", "Das schlagende Herz eines Lavaelementars.")
	quest("knight_emblem", "Obsidian-Wappen", "🛡️", "Das Wappen eines gefallenen Ritters.")
	quest("fire_vial", "Glutphiole", "⚗️", "Eine Phiole mit flüssigem Feuer. Vorsicht!")
	quest("ash_report", "Aschebericht", "📄", "Ein Bericht des Spähers an die Kommandantin.")

	// ---- Dungeons ----
	junk("naga_scale", "Nagaschuppe", "🐍", 300, "Glatt und kühl.")
	good("temple_idol", "Tempelidol", "🗿", 380, "Ein Schlangengott in Jade.")
	junk("ember_chip", "Glutsplitter", "🔸", 520, "Glimmt noch leicht.")
	good("serpent_fang", "Schlangenzahn", "🦷", 400, "Nicht berühren!")
	good("molten_gem", "Geschmolzener Edelstein", "💎", 600, "Flüssiges Feuer in Kristallform.")

	// quest rewards with real stats
	set := func(id, name, icon, slot, typ string, armor int, stats map[string]int, req int, classes ...string) {
		Items[id] = &Item{ID: id, Name: name, Icon: icon, Rarity: "rare", Desc: "Ein Geschenk für verdiente Helden.", Sell: 400 + req*40, Stack: 1,
			Slot: slot, Type: typ, Armor: armor, Stats: stats, Req: req, Classes: classes}
	}
	set("jade_amulet", "Jadeamulett", "📿", "neck", "Amulett", 0, map[string]int{"sta": 14, "int": 10, "str": 10, "agi": 10}, 22)
	set("lotus_ring", "Lotusring", "💍", "finger", "Ring", 0, map[string]int{"sta": 18, "agi": 14, "int": 14, "str": 14}, 26)
	set("fox_cloak", "Fuchsmantel", "🧥", "back", "Stoff", 40, map[string]int{"sta": 20, "agi": 12, "int": 12}, 28)
	set("titan_belt", "Titanengürtel", "🎗️", "waist", "Kette", 90, map[string]int{"sta": 26, "str": 18}, 32, "warrior")
	set("ember_pendant", "Glutanhänger", "📿", "neck", "Amulett", 0, map[string]int{"sta": 28, "int": 22, "agi": 22, "str": 22}, 36)
	set("dragon_ring", "Drachenring", "💍", "finger", "Ring", 0, map[string]int{"sta": 34, "int": 26, "agi": 26, "str": 26}, 38)
}
