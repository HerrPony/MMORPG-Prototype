package gamedata

func init() {
	// ================= Lotushafen / Bambushain / Smaragddschungel (21-24) =================
	qitem("n_tea_chest", "Teekiste für den Hafen", "📦", "Feinster Jasmintee, sorgfältig verpackt.")
	qitem("n_fish_basket", "Korb mit Silberfischen", "🐟", "Frisch gefangen und gut gekühlt.")
	qitem("n_jade_seal", "Jadesiegel", "🔰", "Das Siegel der Hafenmeisterin.")
	qitem("n_herb_pouch", "Heilkräuterbeutel", "🌿", "Hanas Kräuter für die Wachen.")

	nq("lin", "n_tea_to_xia", "Tee für die Händlerin", 21,
		"Xia verkauft den besten Jasmintee der Küste, aber er muss frisch ankommen. Bring die Kiste zu ihr, %n, ohne sie fallen zu lassen.",
		"Xia wartet auf ihren Tee.", "Die Kiste ist heil geblieben? Hervorragend, Xia wird zufrieden sein.",
		objs(dv("xia")), give("n_tea_chest"), items("food_rice"))
	nq("mei", "n_croc_scales_jade", "Schuppen für die Netze", 21,
		"Krokodilschuppen im Netz halten die Fische besser. Bring mir sechs von den Flusskrokodilen, %n, aber trau dem Wasser nicht.",
		"Sechs Flusskrokodil-Schuppen.", "Genau, was ich brauchte. Meine Netze werden es dir danken.",
		objs(cn("x_river_croc_1", 6)))
	nq("jiro", "n_tiger_bones", "Tigerknochen", 21,
		"Die Tiger im Dschungel werden zur Gefahr. Bring mir sechs Knochen als Beweis für ihre Zahl, %n.",
		"Sechs Tigerknochen.", "Das sind mehr, als ich dachte. Wir müssen vorsichtiger sein.",
		objs(cn("x_jungle_tiger_1", 6)))
	nq("bao", "n_tiger_claws", "Krallen für die Kunden", 21,
		"Tigerkrallen sind in der Hauptstadt hochbegehrt. Bring mir vier, %n, und wir teilen den Gewinn.",
		"Vier Tigerkrallen.", "Schöne Stücke! Das wird sich lohnen.",
		objs(cn("x_jungle_tiger_2", 4)), gear())

	nq("mei", "n_fish_run", "Fisch für das Dorf", 22,
		"Bambushain wartet auf frischen Fisch, und die Träger sind krank. Bring den Korb zu Yuna, %n, bevor der Fisch verdirbt.",
		"Yuna wartet auf den Fisch.", "Yuna hat den Fisch erhalten? Gut. Das Dorf wird ein Festmahl feiern.",
		objs(dv("yuna")), give("n_fish_basket"), items("food_rice"))
	nq("jiro", "n_tiger_path", "Der Tigerpfad", 22,
		"Der Tigerpfad ist die Hauptroute zum Dschungel. Erkunde ihn, %n, und erlege acht Tiger, damit die Händler wieder passieren können.",
		"Der Tigerpfad, acht Tiger.", "Danke. Der Pfad ist jetzt sicher.",
		objs(ex("tigerpfad"), kn("jungle_tiger", 8)), gear())
	nq("xia", "n_panther_claws", "Pantherkrallen für den Schmuck", 22,
		"Schwarze Pantherkrallen, in Jade gefasst, sind der letzte Schrei. Bring mir sechs, %n.",
		"Sechs Pantherkrallen.", "Wunderschön. Daraus wird feinster Schmuck.",
		objs(cn("x_panther_2", 6)))
	nq("taro", "n_panther_hunt", "Schwarze Schatten", 22,
		"Die Panther schleichen sich nachts ins Dorf. Erlege sieben von ihnen, bevor sie einen Wächter erwischen, %n.",
		"Sieben Panther.", "Die Nächte sind wieder ruhig. Danke.",
		objs(kn("panther", 7)))

	nq("yuna", "n_ape_grove", "Der Affenhain", 23,
		"Die Affen im Hain werden dreist und stehlen unsere Vorräte. Erkunde den Affenhain, %n, und erlege acht von ihnen.",
		"Der Affenhain, acht Affen.", "Danke. Wir können wieder ohne Sorge ernten.",
		objs(ex("affenhain"), kn("jungle_ape", 8)))
	nq("taro", "n_poacher_ledgers", "Wildererbeutel", 23,
		"Die Wilderer tragen Beutel voller gestohlener Münzen. Bring mir sechs, %n, damit wir die Opfer entschädigen können.",
		"Sechs Wildererbeutel.", "Damit können wir einigen Familien helfen. Gut gemacht.",
		objs(cn("x_poacher_1", 6)), gear())
	nq("hana", "n_herbs_for_guards", "Kräuter für die Wachen", 23,
		"Die Wachen im Lager sind verletzt und brauchen Kräuter. Bring ihnen diesen Beutel, %n, aber geh nicht durch den Dschungel.",
		"Jiro wartet auf die Kräuter.", "Jiro hat die Kräuter bekommen? Dann sind die Männer versorgt.",
		objs(dv("jiro")), give("n_herb_pouch"), items("potion_hp4"))
	nq("kenzo", "n_ape_bones", "Affenknochen", 23,
		"Aus Affenknochen schnitze ich Pfeifen, die sich in den Dörfern gut verkaufen. Bring mir sechs, %n.",
		"Sechs Affenknochen.", "Schöne Stücke. Die werden hübsche Pfeifen.",
		objs(cn("x_jungle_ape_1", 6)))

	nq("yuna", "n_poacher_camp", "Das Wildererlager", 24,
		"Im Wildererlager sammeln sich die Schurken. Erkunde es, %n, und erlege zehn Wilderer, damit der Dschungel Ruhe findet.",
		"Das Wildererlager, zehn Wilderer.", "Dem Dschungel geht es jetzt besser. Danke.",
		objs(ex("wildererlager"), kn("poacher", 10)))
	nq("hana", "n_lizard_scales", "Schuppen für die Medizin", 24,
		"Echsenschuppen sind ein Grundstoff für meine Tinkturen. Bring mir sechs, %n, ich bezahle gut.",
		"Sechs Echsenschuppen.", "Hervorragend. Meine Patienten werden es dir danken.",
		objs(cn("x_lizardman_1", 6)), gear())
	nq("taro", "n_lizard_ruins", "Die Echsenruinen", 24,
		"In den Ruinen leben Echsenmenschen. Erlege acht von ihnen, %n, damit sich die Dörfer sicher fühlen.",
		"Acht Echsenkrieger.", "Das hilft uns enorm.",
		objs(kn("lizardman", 8)))
	nq("jiro", "n_lotus_pond", "Der Lotusteich", 24,
		"Am Lotusteich sind seltsame Dinge geschehen. Erkunde ihn, %n, und erlege die Echsenschamanen, die dort ihr Unwesen treiben.",
		"Der Lotusteich, vier Schamanen.", "Wir wissen jetzt, wer dahintersteckt. Danke.",
		objs(ex("lotusteich"), kn("lizard_shaman", 4)))

	// ================= Himmelspforte / Kirschblütenhain (25-30) =================
	qitem("n_imperial_letter", "Kaiserlicher Brief", "✉️", "Mit dem Siegel der Kaiserin versehen.")
	qitem("n_tea_set", "Teeservice", "🍵", "Zerbrechlich! Aikos bestes Porzellan.")
	qitem("n_ember_lantern", "Geisterlaterne", "🏮", "Eine Laterne mit blauem Licht. Sie wärmt nicht, schützt aber.")
	qitem("n_blade_oil", "Klingenöl", "🧴", "Genjis Spezialöl für Klingen.")

	nq("meilin", "n_letter_to_akane", "Der kaiserliche Brief", 25,
		"Ein Brief an Generalin Akane, von höchster Wichtigkeit. Überbringe ihn, %n, und lass ihn nicht in falsche Hände geraten.",
		"Akane wartet auf den Brief.", "Akane hat ihn gelesen? Dann ist alles in Bewegung.",
		objs(dv("akane")), give("n_imperial_letter"), items("food_feast"))
	nq("akane", "n_fox_tails", "Fuchsschweife", 25,
		"Die Fuchsgeister stehlen Reisenden die Orientierung. Sammle vier Hauche ihrer Magie, %n, damit unsere Magier ihre Schwächen studieren können.",
		"Vier Fuchs-Hauche.", "Das reicht für die Forschung. Danke.",
		objs(cn("x_fox_spirit_1", 4)), gear())
	nq("genji", "n_incense_run", "Duft der Ruhe", 25,
		"Meine Mönche brauchen Weihrauch für die Zeremonie, aber der Weg zum Schrein ist zu gefährlich. Erkunde den Schrein, %n, und erlege dort sechs Fuchsgeister.",
		"Der Fuchsschrein, sechs Fuchsgeister.", "Nun können wir die Zeremonie vorbereiten.",
		objs(ex("fuchsschrein"), kn("fox_spirit", 6)))
	nq("sen", "n_fox_veils", "Schleier für die Seide", 25,
		"Fuchsschleier lassen Seide schimmern wie Mondlicht. Bring mir vier, %n, ich bezahle in Gold.",
		"Vier Fuchs-Schleier.", "Prachtvoll! Die Seide wird eine Sensation.",
		objs(cn("x_fox_spirit_2", 4)))

	nq("aiko", "n_tea_set_delivery", "Teeservice für die Kaiserin", 26,
		"Das beste Porzellan für die Kaiserin. Es darf auf keinen Fall zerbrechen, %n. Trage es mit beiden Händen.",
		"Mei-Lin wartet auf das Teeservice.", "Mei-Lin war entzückt? Ich wusste, du bist vorsichtig.",
		objs(dv("meilin")), give("n_tea_set"), items("food_feast"))
	nq("hiro", "n_panda_hunt", "Bambusfresser", 26,
		"Die Pandabären zertrampeln meine Bambusfelder. Erlege sechs von ihnen, %n, aber tu es schnell: Sie sind größer, als sie aussehen.",
		"Sechs Bambusbären.", "Der Bambus kann nachwachsen. Hab Dank.",
		objs(kn("bamboo_bear", 6)))
	nq("rin", "n_bear_claws_jade", "Krallen für den Schmuck", 26,
		"Bambusbärenkrallen sind ein beliebtes Mitbringsel aus Kirschbach. Bring mir fünf, %n.",
		"Fünf Bambusbärkrallen.", "Prächtig. Die verkaufen sich von selbst.",
		objs(cn("x_bamboo_bear_2", 5)), gear())
	nq("ume", "n_silk_trail", "Die Seidenstraße", 26,
		"Seide kommt aus Kirschbach, aber die Straße ist unsicher. Erkunde Kirschbach, %n, und sieh dir die Lage an.",
		"Kirschbach erkunden.", "Dann ist es besser, als ich dachte. Danke.",
		objs(ex("kirschbach")))

	nq("akane", "n_ronin_camp", "Das Ronin-Lager", 27,
		"Die Ronin sammeln sich im Lager im Norden. Erkunde es, %n, und erlege zehn von ihnen, bevor sie zur Hauptstadt marschieren.",
		"Das Ronin-Lager, zehn Ronin.", "Das war überfällig. Die Stadt ist sicherer.",
		objs(ex("roninlager"), kn("ronin", 10)))
	nq("goro", "n_golem_cores", "Kerne für die Schmiede", 27,
		"Steinwächterkerne speichern Hitze lange. Ideal für meine Schmiede. Bring mir fünf Kernstücke, %n.",
		"Fünf Steinwächter-Kernstücke.", "Wunderbar! Mein Feuer wird nie wieder ausgehen.",
		objs(cn("x_stone_golem_2", 5)))
	nq("takumi", "n_ghost_essence", "Hauch der Geister", 27,
		"Meine Forschungen über die Geisterkrieger stocken. Bring mir sechs Hauche, %n, damit ich ihre Natur ergründen kann.",
		"Sechs Hauche der Geisterkrieger.", "Faszinierend. Das ist sehr aufschlussreich.",
		objs(cn("x_ghost_warrior_1", 6)))
	nq("hiro", "n_lantern_to_rin", "Laterne für die Füchsin", 27,
		"Rin will in den Wald, aber ohne Laterne findet sie nicht zurück. Bring ihr diese hier, %n, sie schützt vor den Geistern.",
		"Rin wartet auf die Laterne.", "Rin hat sie genommen? Dann bin ich beruhigt.",
		objs(dv("rin")), give("n_ember_lantern"))

	nq("akane", "n_ronin_seals_jade", "Siegel der Verstoßenen", 28,
		"Jeder Ronin trägt ein Zeichen seines früheren Herrn. Sammle sechs Beutel von ihnen, %n, damit wir die Verbindungen aufdecken.",
		"Sechs Ronin-Beutel.", "Genau die Hinweise, die ich brauchte.",
		objs(cn("x_ronin_1", 6)), gear())
	nq("genji", "n_blade_oil_run", "Öl für den Meister", 28,
		"Großmeister Ryu braucht Öl für seine Klinge. Bring es ihm, %n, und sag ihm, dass ich seinen Gruß schicke.",
		"Ryu wartet auf das Öl.", "Ryu hat das Öl erhalten? Das freut mich.",
		objs(dv("ryu")), give("n_blade_oil"))
	nq("rin", "n_mirror_pond", "Der Spiegelteich", 28,
		"Am Spiegelteich spuken Geisterkrieger. Erkunde den Teich, %n, und erlege acht von ihnen, damit der Ort wieder friedlich wird.",
		"Der Spiegelteich, acht Geisterkrieger.", "Der Ort ist wieder friedlich. Danke.",
		objs(ex("spiegelteich"), kn("ghost_warrior", 8)))
	nq("miyu", "n_temple_approach", "Der Weg zum Tempel", 28,
		"Bevor wir in den Nebeltempel eindringen, müssen die Wege gesichert sein. Erlege sechs Steinwächter, %n, die den Zugang blockieren.",
		"Sechs Steinwächter.", "Der Weg ist frei. Danke.",
		objs(kn("stone_golem", 6)))

	nq("takumi", "n_golem_splitters", "Splitter der Wächter", 29,
		"Die Splitter der Steinwächter tragen Inschriften, die ich entziffern will. Bring mir sechs, %n.",
		"Sechs Steinwächter-Splitter.", "Faszinierend. Da steht tatsächlich etwas drauf.",
		objs(cn("x_stone_golem_1", 6)))
	nq("akane", "n_ronin_trail", "Auf der Spur der Verstoßenen", 29,
		"Ein Ronin-Meister führt die anderen an. Erlege acht Ronin, um seine Reihen zu schwächen, %n.",
		"Acht Ronin.", "Gut. Jetzt wird er nervös.",
		objs(kn("ronin", 8)))
	nq("miyu", "n_ghost_shrouds", "Schleier der Toten", 29,
		"Die Schleier der Geisterkrieger geben mir Aufschluss über die Zeit, in der sie lebten. Bring mir fünf, %n.",
		"Fünf Schleier der Geisterkrieger.", "Das sind wertvolle Hinweise. Danke.",
		objs(cn("x_ghost_warrior_2", 5)), gear())
	nq("goro", "n_stone_hunt", "Steinwächter-Jagd", 29,
		"Die Steinwächter marschieren neuerdings in die Dörfer. Erlege acht von ihnen, %n.",
		"Acht Steinwächter.", "Gut. Die Dörfer sind sicher.",
		objs(kn("stone_golem", 8)))

	nq("meilin", "n_empress_envoy", "Gesandte der Kaiserin", 30,
		"Die Lage an der Grenze zum Rotkluft-Canyon ist angespannt. Sprich mit Mara, %n, und überbringe ihr die Grüße der Kaiserin.",
		"Mara wartet auf deine Nachricht.", "Das sind gute Neuigkeiten. Danke.",
		objs(tk("mara")), items("food_feast"))
	nq("takumi", "n_temple_scouting", "Der Nebeltempel", 30,
		"Wir müssen wissen, was im Nebeltempel vor sich geht. Erkunde ihn, %n, und erlege dabei sechs Steinwächter.",
		"Der Nebeltempel, sechs Steinwächter.", "Dann ist der Tempel noch unberührt. Danke.",
		objs(ex("nebeltempel"), kn("stone_golem", 6)))
	nq("miyu", "n_ancient_inscriptions", "Alte Inschriften", 30,
		"Die Inschriften am Tempel geben Hinweise auf den Nagakaiser. Erkunde den Tempel und bring mir drei Steinwächterkerne, %n.",
		"Der Nebeltempel, drei Kernstücke.", "Das bestätigt meine Vermutung. Danke.",
		objs(ex("nebeltempel"), cn("x_stone_golem_2", 3)), gear())
	nq("akane", "n_defend_gate", "Die Verteidigung der Pforte", 30,
		"Geisterkrieger rücken auf Himmelspforte vor. Erlege zehn von ihnen, %n, bevor sie die Mauern erreichen.",
		"Zehn Geisterkrieger.", "Die Stadt ist gerettet. Danke!",
		objs(kn("ghost_warrior", 10)))
}
