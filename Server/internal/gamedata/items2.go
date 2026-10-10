package gamedata

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
	// vendor trash from the new zones
	good("hyena_pelt", "Hyänenfell", "🐆", 24, "Struppig und streng riechend.")
	good("croc_hide", "Krokodilleder", "🐊", 44, "Zäh wie Eisen.")
	junk("rotten_cloth", "Faulender Stofffetzen", "🧻", 18, "Besser nicht anfassen.")
	good("witch_herb", "Hexenkraut", "🌿", 50, "Duftet verdächtig.")
	good("yeti_fur", "Yetifell", "🧶", 85, "Dick, weiß und warm.")
	good("scorpion_tail", "Skorpionstachel", "🦂", 100, "Gift tropft noch davon.")
	junk("mummy_wrap", "Alte Mumienbinde", "🩹", 70, "Staubtrocken.")
	junk("kobold_candle", "Kobold-Kerze", "🕯️", 45, "Brennt mit grünlicher Flamme.")
	junk("bone_dust", "Knochenstaub", "💀", 40, "Knirscht zwischen den Fingern.")

	// quest items
	quest("letter_aldric", "Eiliger Brief", "✉️", "Ein Brief von Förster Aldric. Er muss schnell zugestellt werden.")
	quest("letter_alena", "Brief an die Bürgermeisterin", "📜", "Aldric bittet Alena um Unterstützung gegen die Wolfsplage.")
	quest("order_royal", "Königlicher Befehl", "📯", "Mit dem Siegel von Sonnenburg versehen. Dringend!")
	quest("phial_water", "Versiegelte Phiole", "⚗️", "Ein schimmernder Trank, der ins Quellheiligtum gebracht werden soll.")
	quest("shield_plate", "Schildplatte", "🔰", "Frisch geschmiedetes Eisen. Rodrik wartet darauf.")
	quest("medicine_vial", "Medizin-Phiole", "🧪", "Sylvas Heilmittel. Rasch zum Kranken bringen!")
	quest("pass_message", "Botschaft des Hauptmanns", "📩", "Eine Nachricht von Hauptmann Isgard an die Wache am Windpass.")
	quest("spice_box", "Gewürzkiste", "📦", "Teure Gewürze, die schnell geliefert werden müssen.")
	quest("pumpkin", "Kürbis", "🎃", "Ein praller Herbstkürbis.")
	quest("wheat", "Weizenbündel", "🌾", "Goldgelber Weizen.")
	quest("mill_gear", "Mühlenzahnrad", "⚙️", "Gehört in die Windmühle von Hannes.")
	quest("bull_horn", "Stierhorn", "🐂", "Ein mächtiges Horn.")
	quest("spring_water", "Quellwasser", "💧", "Rein und kühl.")
	quest("track", "Wildfährte", "👣", "Aldwin will wissen, welche Tiere hier jagen.")
	quest("crystal_water", "Wasserkristall", "💎", "Er pulsiert im Takt des Wassers.")
	quest("crystal_heart", "Herz des Wassers", "💠", "Die Quelle der Verzerrung.")
	quest("wolf_alpha_claw", "Alphaklaue", "🐾", "Beweis für den Sieg über Rotpelz.")
	quest("shield_shards", "Zerbeulte Schildreste", "🛡️", "Von Banditen erbeutet. Brauchbar für Rodrik.")
	quest("lotus", "Moorlotus", "🪷", "Heilkraut aus den Düstermooren.")
	quest("wisp_dust", "Irrlichtstaub", "✨", "Schimmert unwirklich.")
	quest("icewolf_pelt", "Eiswolfpelz", "🐺", "Weiß wie frisch gefallener Schnee.")
	quest("ice_shard", "Eissplitter", "🧊", "Eiskalt und scharfkantig.")
	quest("caravan_map", "Karawanenkarte", "🗺️", "Zeigt die Wege durch die Dünen.")
	quest("cactus_fruit", "Kaktusfrucht", "🌵", "Süß und saftig.")
	quest("ancient_shard", "Antike Scherbe", "🏺", "Aus der Zeit der Pharaonen.")
}
