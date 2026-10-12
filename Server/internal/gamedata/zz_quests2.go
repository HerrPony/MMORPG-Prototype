package gamedata

func init() {
	// ================= Moorbruck / Düstermoor (10-14) =================
	qitem("n_swamp_tonic", "Moortonic", "🍵", "Sylvas Gebräu. Es riecht nach nassen Socken, wirkt aber Wunder.")
	qitem("n_net_repair", "Geflickte Netze", "🕸️", "Tams Netze, frisch geflickt.")
	qitem("n_ward_charm", "Schutzamulett", "🧿", "Ein einfacher Talisman gegen das Moorlicht.")

	nq("moros", "n_croc_scales", "Schuppen für die Palisade", 10,
		"Die Palisade um das Dorf braucht Verstärkung, und Krokodilschuppen sind hart wie Eisen. Sammle sechs Schuppen von den Sumpfkrokodilen, %n.",
		"Sechs Schuppen, dann hält die Palisade.", "Das wird halten. Die Moorbewohner danken dir.",
		objs(cn("x_croc_1", 6)), gear())
	nq("sylva", "n_croc_eyes", "Augen für den Trank", 10,
		"Krokodilaugen sind der Grundstoff für meinen Sehtrank. Ich brauche vier. Sei vorsichtig, die Biester beißen!",
		"Vier Augen, nicht mehr.", "Perfekt. Bald können die alten Männer wieder lesen.",
		objs(cn("x_croc_2", 4)), items("potion_hp3"))

	nq("hektor", "n_swamp_patrol", "Sumpfpatrouille", 11,
		"Die Krokodilsümpfe sind das Rückgrat unserer Verteidigung, aber niemand sieht dort mehr nach dem Rechten. Erkunde sie, %n, und erlege dabei sechs Krokodile.",
		"Die Krokodilsümpfe, sechs Krokodile.", "Gut. Jetzt weiß ich, wo wir die Wache verstärken müssen.",
		objs(ex("krokodilsuempfe"), kn("croc", 6)))
	nq("sylva", "n_tonic_run", "Moortonic für Hektor", 11,
		"Hektor hat sich bei der Wache eine Erkältung geholt. Bring ihm dieses Tonic, %n. Er wird behaupten, es schmecke widerlich, aber er wird es trinken.",
		"Hektor wartet auf sein Tonic.", "Er hat es getrunken und nur zweimal geflucht? Dann wirkt es.",
		objs(dv("hektor")), give("n_swamp_tonic"), items("food_stew"))
	nq("boggs", "n_zombie_trophies", "Beweise für die Wache", 11,
		"Die Wache zahlt eine Prämie für jedes erledigte Untote. Sie wollen Beweise: Zombie-Grabstaub, in kleinen Fläschchen. Bring mir sechs davon, ich löse sie ein und teile mit dir.",
		"Sechs Fläschchen Grabstaub.", "Das sind die Prämien. Und dein Anteil.",
		objs(cn("x_zombie_2", 6)))
	nq("moros", "n_zombie_purge", "Reinigung des Dorfrandes", 11,
		"Die Zombies kommen jede Nacht näher an die Häuser. Erlege zehn von ihnen und die Kinder können wieder schlafen, %n.",
		"Zehn Zombies, damit die Nächte ruhig bleiben.", "Es ist still geworden. Hab Dank.",
		objs(kn("zombie", 10)))

	nq("tam", "n_net_to_boggs", "Geflickte Netze", 12,
		"Boggs wollte meine Netze flicken lassen, aber Hexen haben das Dorf verhext und er findet keine Fäden. Bring sie zu Hektor, er hat noch Garn.",
		"Hektor hat das Garn.", "Genau diese Netze. Tam wird sich freuen.",
		objs(dv("hektor")), give("n_net_repair"))
	nq("hektor", "n_witch_purses", "Hexenbeutel", 12,
		"Die Moorhexen tragen Beutel mit seltsamen Zutaten. Mit etwas Glück erfahren wir, was sie planen. Bring mir fünf davon, %n.",
		"Fünf Beutel der Hexen.", "Interessant. Wurzeln, Pilze, Augen. Nichts Gutes.",
		objs(cn("x_witch_1", 5)))
	nq("sylva", "n_witch_compass", "Verlorene Orientierung", 12,
		"Die Hexen stehlen die Kompasse der Reisenden, damit sie sich im Moor verirren. Hol mir vier zurück, %n, die Reisenden werden es zu schätzen wissen.",
		"Vier Kompasse, damit sich niemand verirrt.", "Wunderbar. Die Reisenden werden es dir danken.",
		objs(cn("x_witch_2", 4)), gear())
	nq("tam", "n_wisp_hunt", "Lichter im Nebel", 12,
		"Die Irrlichter locken Fischer auf den See. Zerstreue sie, %n: Erlege sechs von ihnen und bring mir drei Schleier als Beweis.",
		"Sechs Irrlichter, drei Schleier.", "Der See ist sicherer. Ich danke dir.",
		objs(kn("wisp", 6), cn("x_wisp_2", 3)))

	nq("moros", "n_ward_amulets", "Schutzamulette", 13,
		"Die Dorfbewohner brauchen Schutz gegen das Moorlicht. Bring dieses Amulett zu Tam, er soll es in die Fischerboote hängen.",
		"Tam wartet auf die Amulette.", "Tam hat die Amulette verteilt? Dann sind wir gerüstet.",
		objs(dv("tam")), give("n_ward_charm"))
	nq("hektor", "n_hag_scout", "Erkundung des Hexenhains", 13,
		"Mutter Grimmel haust im Hexenhain. Bevor wir etwas unternehmen, will ich wissen, was dort vor sich geht. Erkunde den Ort, %n, und erlege acht Hexen.",
		"Der Hexenhain, acht Hexen.", "Wir haben ihnen einen Dämpfer verpasst. Gute Arbeit!",
		objs(ex("hexenhain"), kn("witch", 8)), gear())
	nq("boggs", "n_lake_trade", "Handelsweg über den See", 13,
		"Der Nebelsee ist der kürzeste Weg nach Süden, aber seit Wochen traut sich niemand mehr darüber. Erkunde den Nebelsee, %n, und erlege die Irrlichter, die dort lauern.",
		"Der Nebelsee, sechs Irrlichter.", "Der Weg ist frei. Meine Karawanen danken dir.",
		objs(ex("nebelsee"), kn("wisp", 6)))
	nq("pott", "n_swamp_trinkets", "Glitzerkram aus dem Moor", 13,
		"Wenn man weiß, wo man suchen muss, findet man im Moor allerlei Glitzerkram. Bring mir vier Irrlicht-Hauche, ich mache daraus Schmuck.",
		"Vier Irrlicht-Hauche.", "Wunderschön! Daraus lässt sich etwas machen.",
		objs(cn("x_wisp_1", 4)))

	// ================= Frostwacht / Eisgipfel (14-17) =================
	qitem("n_frost_charm", "Frostamulett", "❄️", "Ein Talisman aus Eissplittern, geschmiedet von Hilde.")
	qitem("n_pelt_order", "Pelzbestellung", "📜", "Eine Liste der benötigten Felle, mit dem Siegel der Wache.")
	qitem("n_ration_sack", "Verpflegungssack", "🎒", "Getrocknetes Fleisch und Hartbrot für die Späher am Pass.")
	qitem("n_ice_map", "Eiskarte", "🗺️", "Vorn hat die Pfade über den Gletscher eingezeichnet.")

	nq("olaf", "n_icewolf_bones", "Knochen für die Hunde", 14,
		"Unsere Schlittenhunde brauchen Knochen, und die Eiswölfe liefern sie reichlich. Bring mir sechs Eiswolfsknochen, %n.",
		"Sechs Knochen, die Hunde warten.", "Sie werden es dir mit Schwanzwedeln danken.",
		objs(cn("x_icewolf_1", 6)))
	nq("isgard", "n_pass_patrol", "Patrouille am Windpass", 14,
		"Der Windpass ist unsere Verbindung nach Süden. Besuche den Pass, %n, und erlege dabei acht Eiswölfe, bevor die Händler dort wieder durchkommen.",
		"Der Windpass, acht Eiswölfe.", "Die Händler können passieren. Danke!",
		objs(ex("windpass"), kn("icewolf", 8)), gear())
	nq("vorn", "n_map_for_isgard", "Karte für den Hauptmann", 14,
		"Ich habe einen sicheren Pfad über den Gletscher kartiert. Der Hauptmann sollte sie kennen. Bring sie ihm, %n, und achte auf die Eiswölfe.",
		"Isgard wartet auf die Karte.", "Isgard hat die Karte? Dann sind die Streifen künftig sicherer unterwegs.",
		objs(dv("isgard")), give("n_ice_map"))
	nq("hilde", "n_icewolf_claws", "Krallen für den Stahl", 14,
		"Eiswolfskrallen im Stahl geben der Klinge einen blauen Schimmer. Bring mir fünf, %n, und ich schmiede dir etwas Hübsches.",
		"Fünf Krallen für den blauen Stahl.", "Wunderbar! Der Stahl wird eine Pracht.",
		objs(cn("x_icewolf_2", 5)), gear())

	nq("olaf", "n_yeti_bones", "Yetiknochen", 15,
		"Die Yetis sind größer als jeder Bär und doppelt so zäh. Ihre Knochen geben die besten Schlittenkufen. Bring mir sechs, %n.",
		"Sechs Yetiknochen.", "Daraus werden die besten Schlitten der Region.",
		objs(cn("x_yeti_1", 6)))
	nq("isgard", "n_ration_run", "Verpflegung für die Späher", 15,
		"Die Späher am Pass haben seit zwei Tagen nichts gegessen. Bring ihnen diesen Sack zu Vorn, er verteilt ihn.",
		"Vorn wartet auf die Verpflegung.", "Vorn hat verteilt? Gut, dann sind die Späher versorgt.",
		objs(dv("vorn")), give("n_ration_sack"), items("food_roast"))
	nq("sven", "n_yeti_claws", "Yetikrallen für den Schmuck", 15,
		"Aus Yetikrallen mache ich die beliebtesten Anhänger des Nordens. Bring mir fünf Stück, %n, und ich zahle gut.",
		"Fünf Yetikrallen.", "Prachtvoll. Die werden sich verkaufen.",
		objs(cn("x_yeti_2", 5)))
	nq("hilde", "n_frost_charms", "Frostamulette", 15,
		"Ich habe Amulette gegen die Kälte gefertigt. Bring eines zu Olaf, er braucht es für seine Jagden im Hochland.",
		"Olaf wartet auf das Amulett.", "Olaf friert nicht mehr? Dann hat es sich gelohnt.",
		objs(dv("olaf")), give("n_frost_charm"))

	nq("olaf", "n_bear_cliffs", "Die Eisbärenklippe", 16,
		"Die Eisbären an der Klippe werden zur Plage. Erkunde die Eisbärenklippe, %n, und erlege sechs von ihnen.",
		"Die Eisbärenklippe, sechs Eisbären.", "Das war überfällig. Danke.",
		objs(ex("baerenklippe"), kn("polarbear", 6)), gear())
	nq("vorn", "n_spirit_veils", "Geisterschleier", 16,
		"Die Eisgeister hinterlassen Schleier aus eisigem Nebel. Bring mir fünf davon, ich studiere ihre Zusammensetzung.",
		"Fünf Schleier der Eisgeister.", "Faszinierend. Noch kälter, als ich dachte.",
		objs(cn("x_frostspirit_2", 5)))
	nq("isgard", "n_yeti_cull", "Yeti-Jagd", 16,
		"Die Yetis kommen in großen Gruppen näher. Erlege zehn von ihnen, bevor sie die Wache erreichen, %n.",
		"Zehn Yetis.", "Das hat die Reihen gelichtet. Danke.",
		objs(kn("yeti", 10)))
	nq("sven", "n_polar_claws", "Eisbärenkrallen", 16,
		"Eine Eisbärenkralle am Hals bringt Glück bei der Jagd, sagt man. Bring mir vier, ich fertige Anhänger.",
		"Vier Eisbärenkrallen.", "Hübsch und kräftig. Die verkaufe ich teuer.",
		objs(cn("x_polarbear_2", 4)))

	nq("torg", "n_ore_for_torg", "Erz für den Händler", 17,
		"Mein Geschäft läuft schlecht, die Banditen vergraulen die Kundschaft. Erkunde Rabenfels, %n, und finde heraus, wer dahintersteckt.",
		"Rabenfels erkunden.", "Ich habe es geahnt. Danke für die Information.",
		objs(ex("rabenfels")))
	nq("isgard", "n_ice_spirits", "Eisgeister", 17,
		"Die Eisgeister rauben den Spähern den Schlaf. Erlege sieben von ihnen, %n.",
		"Sieben Eisgeister.", "Endlich wieder Ruhe an der Wache.",
		objs(kn("frostspirit", 7)))
	nq("karim", "n_dune_scouts", "Späher in den Dünen", 17,
		"Die Skorpiondünen sind unsere Grenze. Wenn dort etwas passiert, will ich es wissen. Erkunde die Dünen, %n, und erlege acht Skorpione.",
		"Die Skorpiondünen, acht Skorpione.", "Danke. Jetzt kann ich die Karawanen beruhigen.",
		objs(ex("skorpionduenen"), kn("scorpion", 8)))
	nq("jasim", "n_scorpion_chitin", "Panzer für die Kamele", 17,
		"Meine Kamele brauchen Schutz vor den Dornen der Wüste. Skorpionchitin ist leicht und stark. Bring mir sechs Stücke, %n.",
		"Sechs Stück Chitin.", "Ausgezeichnet. Die Kamele werden stolz sein.",
		objs(cn("x_scorpion_1", 6)))

	// ================= Sandfurt / Wüste (17-20) =================
	qitem("n_water_skins", "Wasserschläuche", "💧", "Frisches Wasser für die Karawane.")
	qitem("n_desert_amulet", "Wüstenamulett", "📿", "Ein Talisman gegen Sonnenstich.")
	qitem("n_scroll_neferet", "Papyrusrolle", "📜", "Neferets Notizen zu den Gräbern.")

	nq("karim", "n_water_for_bahir", "Wasser für die Wache", 18,
		"Bahir und seine Männer haben seit gestern kein Wasser mehr. Bring ihnen diese Schläuche, %n, bevor die Sonne sie umbringt.",
		"Bahir braucht das Wasser.", "Sie haben getrunken? Dann sind sie gerettet.",
		objs(dv("bahir")), give("n_water_skins"), items("food_roast"))
	nq("bahir", "n_dune_hunt", "Skorpionplage", 18,
		"Die Skorpione verlassen die Dünen. Das ist ungewöhnlich. Erlege zehn von ihnen, bevor sie die Karawanenwege erreichen.",
		"Zehn Skorpione.", "Die Wege sind wieder frei. Danke!",
		objs(kn("scorpion", 10)))
	nq("rashid", "n_raider_purses", "Beutel der Räuber", 18,
		"Die Wüstenräuber haben mir drei Karawanen ausgeraubt. Hol mir ihre Beutel zurück, %n. Fünf Stück sollten reichen, um den Schaden zu begleichen.",
		"Fünf Beutel der Räuber.", "Das deckt zumindest einen Teil des Schadens. Hab Dank.",
		objs(cn("x_desert_raider_1", 5)), gear())
	nq("tarek", "n_scorpion_stingers", "Stachel für die Klingen", 18,
		"Skorpionstachel im Stahl macht ihn zähflüssig-giftig. Genau das, was ein Dolchschmied sich wünscht. Bring mir vier Fühler, %n.",
		"Vier Skorpionfühler.", "Perfekt. Das wird eine tödliche Klinge.",
		objs(cn("x_scorpion_2", 4)))

	nq("jasim", "n_oasis_trade", "Die Oase", 19,
		"Die Oase ist der einzige Ort mit sicherem Wasser. Erkunde die Oase, %n, und räume dort ein paar Wüstenräuber aus dem Weg.",
		"Die Oase, fünf Wüstenräuber.", "Das Wasser ist wieder sicher. Die Karawanen danken dir.",
		objs(ex("oase"), kn("desert_raider", 5)))
	nq("neferet", "n_papyrus_to_karim", "Papyrus für den Scheich", 19,
		"Meine Notizen zu den Gräbern sollten dem Scheich vorgelegt werden. Er kennt die Geschichten dazu. Bring ihm die Rolle, %n.",
		"Karim wartet auf die Notizen.", "Der Scheich hat sie gelesen? Dann ist er überzeugt.",
		objs(dv("karim")), give("n_scroll_neferet"))
	nq("neferet", "n_mummy_wraps", "Mumienstaub", 19,
		"Der Staub der Mumien enthält Spuren alter Salben. Ich möchte sie untersuchen. Sammle sechs Fläschchen Grabstaub, %n.",
		"Sechs Fläschchen Mumienstaub.", "Faszinierend. Diese Salben sind tausende Jahre alt.",
		objs(cn("x_mummy_2", 6)))
	nq("bahir", "n_raider_compass", "Verlorene Karawanen", 19,
		"Die Kompasse der verschollenen Karawanenführer sind bei den Räubern gelandet. Bring mir fünf, %n, damit wir die Familien benachrichtigen können.",
		"Fünf Kompasse der Räuber.", "Die Familien werden es erfahren. Danke.",
		objs(cn("x_desert_raider_2", 5)))

	nq("neferet", "n_tomb_survey", "Gräber der Pharaonen", 20,
		"Die Gräber sind unsicher geworden. Erkunde das Pharaonengrab, %n, und erlege dort zehn Mumien, damit ich meine Arbeit fortsetzen kann.",
		"Das Pharaonengrab, zehn Mumien.", "Jetzt kann ich in Ruhe forschen. Danke!",
		objs(ex("pharaonengrab"), kn("mummy", 10)), gear())
	nq("karim", "n_desert_amulets", "Wüstenamulette", 20,
		"Meine Männer brauchen Schutz vor der Sonne. Bring dieses Amulett zu Jasim, er soll es an die Karawanenführer verteilen.",
		"Jasim wartet auf das Amulett.", "Jasim hat es verteilt? Dann sind die Karawanen gut geschützt.",
		objs(dv("jasim")), give("n_desert_amulet"), items("potion_hp4"))
	nq("tarek", "n_mummy_bandages", "Binden für den Verband", 20,
		"Mumienbinden sind erstaunlich haltbar. Ich brauche sie für Verbände. Bring mir acht, %n, aber lass dich nicht von den Mumien einwickeln.",
		"Acht Mumien-Knochensplitter.", "Das reicht für ein ganzes Lazarett.",
		objs(cn("x_mummy_1", 8)))
}
