package gamedata

func init() {
	// ================= Rotfels / Rotkluft-Canyon (31-35) =================
	qitem("n_clay_jugs", "Tonkrüge", "🏺", "Olins beste Krüge. Sie dürfen auf keinen Fall zerbrechen.")
	qitem("n_ore_crate", "Erzkiste", "📦", "Brokks Erzproben, schwer und wertvoll.")
	qitem("n_canyon_report", "Späherbericht", "📜", "Noxens Beobachtungen zu den Banditenbewegungen.")

	nq("mara", "n_cat_bones", "Sandkatzenknochen", 31,
		"Die Sandkatzen dezimieren unsere Herden. Als Beweis für die Zahlung der Prämie brauche ich Knochen: sechs Stück, %n.",
		"Sechs Sandkatzenknochen.", "Hier ist die Prämie. Gut gemacht.",
		objs(cn("x_sandcat_1", 6)))
	nq("olin", "n_jugs_to_zane", "Tonkrüge für den Händler", 31,
		"Zane wartet seit Tagen auf meine Krüge. Bring sie ihm, %n, und lass nichts zerbrechen!",
		"Zane wartet auf die Krüge.", "Zane hat sie bekommen? Gut. Ich wusste, du bist vorsichtig.",
		objs(dv("zane")), give("n_clay_jugs"), items("food_feast"))
	nq("brokk", "n_cat_claws", "Krallen für die Klingen", 31,
		"Sandkatzenkrallen, in den Stahl eingeschmolzen, geben der Klinge einen unverwechselbaren Schliff. Bring mir fünf, %n.",
		"Fünf Sandkatzenkrallen.", "Perfekt. Der Schliff wird makellos.",
		objs(cn("x_sandcat_2", 5)), gear())

	nq("nox", "n_pass_scout", "Der Kluftpass", 32,
		"Der Kluftpass ist unsere Lebensader. Erkunde ihn, %n, und erlege dabei sechs Canyon-Bisons, die den Weg blockieren.",
		"Der Kluftpass, sechs Bisons.", "Der Pass ist offen. Hab Dank.",
		objs(ex("kluftpass"), kn("canyon_bison", 6)))
	nq("zane", "n_bison_horns", "Bisonhörner", 32,
		"Aus Bisonhörnern lassen sich hervorragende Trinkgefäße machen. Bring mir fünf Knochen als Rohmaterial, %n.",
		"Fünf Bisonknochen.", "Genau, was ich brauchte. Meine Kunden werden begeistert sein.",
		objs(cn("x_canyon_bison_1", 5)))
	nq("olin", "n_bison_trample", "Gefahr für die Töpferei", 32,
		"Die Bisons trampeln meine Lehmgruben nieder. Erlege acht von ihnen, %n.",
		"Acht Bisons.", "Die Gruben sind sicher. Hab Dank.",
		objs(kn("canyon_bison", 8)), gear())

	nq("brokk", "n_scorp_chitin", "Chitin für die Rüstung", 33,
		"Kluftskorpion-Chitin ist leichter als Stahl und härter als Leder. Ideal für Rüstungen. Bring mir sechs Stück, %n.",
		"Sechs Stück Chitin.", "Hervorragend. Das wird eine leichte, harte Rüstung.",
		objs(cn("x_mesa_scorpion_1", 6)))
	nq("ruk", "n_scorp_cull", "Skorpionplage", 33,
		"Die Kluftskorpione werden zahlreicher. Erlege zehn von ihnen, %n, bevor sie das Dorf erreichen.",
		"Zehn Skorpione.", "Das hat die Zahl gedrückt. Danke!",
		objs(kn("mesa_scorpion", 10)))
	nq("mara", "n_stone_arch", "Der Steinbogen", 33,
		"Der Steinbogen ist ein wichtiger Wegpunkt. Erkunde ihn, %n, und erlege die Canyon-Banditen, die dort Zoll erheben.",
		"Der Steinbogen, sechs Banditen.", "Die Händler können wieder passieren. Danke.",
		objs(ex("steinbogen"), kn("canyon_bandit", 6)), gear())

	nq("zafira", "n_bandit_compasses", "Gestohlene Kompasse", 34,
		"Banditen haben meine Kompasse gestohlen. Es waren fünf, schwarz gebeizt. Bring sie mir zurück, %n.",
		"Fünf Banditenkompasse.", "Da sind sie! Danke.",
		objs(cn("x_canyon_bandit_2", 5)))
	nq("mara", "n_bandit_fort", "Die Banditenfestung", 34,
		"Die Festung ist das Herz der Banditen. Erkunde sie, %n, und erlege acht Banditen, damit sie uns nicht mehr belästigen.",
		"Die Banditenfestung, acht Banditen.", "Das war ein harter Schlag für sie. Danke.",
		objs(ex("banditenfestung"), kn("canyon_bandit", 8)))
	nq("olin", "n_clay_cores", "Kerne der Lehmgolems", 34,
		"Die Lehmgolems tragen warme Kerne im Inneren. Ich möchte sie im Ofen testen. Bring mir sechs Kernstücke, %n.",
		"Sechs Lehmgolem-Kernstücke.", "Wunderbar. Damit kann ich experimentieren.",
		objs(cn("x_clay_golem_2", 6)))

	nq("nox", "n_report_to_mara", "Späherbericht", 35,
		"Ich habe beobachtet, wie die Banditen ihre Reihen ordnen. Bring diesen Bericht zu Mara, %n, sie muss davon erfahren.",
		"Mara wartet auf den Bericht.", "Mara hat den Bericht? Dann kann sie handeln.",
		objs(dv("mara")), give("n_canyon_report"))
	nq("brokk", "n_golem_splitters_mesa", "Splitter für den Amboss", 35,
		"Lehmgolem-Splitter sind der härteste Zusatz für meinen Amboss. Bring mir acht, %n.",
		"Acht Lehmgolem-Splitter.", "Perfekt. Der Amboss wird unzerstörbar.",
		objs(cn("x_clay_golem_1", 8)), gear())
	nq("zane", "n_oasis_hidden", "Die verborgene Oase", 35,
		"Es heißt, in den Felsen verberge sich eine Oase. Erkunde die Verborgene Oase, %n, und erlege die Golems, die sie bewachen.",
		"Die Verborgene Oase, sechs Lehmgolems.", "Dann stimmt das Gerücht. Danke.",
		objs(ex("versteckte_oase"), kn("clay_golem", 6)))

	// ================= Aschenwacht / Aschenlande (36-40) =================
	qitem("n_ember_flask", "Glutphiole", "🧪", "Eine versiegelte Phiole voller flüssiger Glut. Vorsicht!")
	qitem("n_ash_orders", "Befehle der Kommandantin", "📜", "Mit dem Siegel von Aschenwacht versehen.")
	qitem("n_forge_coal", "Schmiedekohle", "⚫", "Rußschwarze Kohle von Ignis.")

	nq("vash", "n_hound_cull", "Aschehund-Plage", 36,
		"Die Aschehunde rotten sich am Pass zusammen. Erlege zehn von ihnen, %n, bevor sie die Wache überrennen.",
		"Zehn Aschehunde.", "Der Pass ist sicherer. Danke!",
		objs(kn("ash_hound", 10)))
	nq("cinder", "n_hound_claws", "Aschehundkrallen", 36,
		"Aschehundkrallen, in Glut gehärtet, sind der letzte Schrei. Bring mir fünf, %n, ich fertige Anhänger.",
		"Fünf Aschehundkrallen.", "Prächtig. Die gehen weg wie nichts.",
		objs(cn("x_ash_hound_2", 5)), gear())
	nq("kael", "n_flask_to_ignis", "Glutphiole für die Schmiede", 36,
		"Ignis braucht diese Phiole für ihr bestes Stück. Bring sie ihr, %n, und halte sie ruhig!",
		"Ignis wartet auf die Phiole.", "Ignis hat sie bekommen? Gut. Das Stück wird ein Meisterwerk.",
		objs(dv("ignis")), give("n_ember_flask"), items("food_dragon"))

	nq("vash", "n_ash_pass", "Der Aschepass", 37,
		"Der Aschepass ist von Imps belagert. Erkunde den Pass, %n, und erlege acht Feuerimps.",
		"Der Aschepass, acht Feuerimps.", "Der Pass ist offen. Danke.",
		objs(ex("aschepass"), kn("fire_imp", 8)), gear())
	nq("ignis", "n_imp_veils", "Schleier der Imps", 37,
		"Imp-Schleier schützen vor Hitze. Ich will sie in meine Handschuhe einarbeiten. Bring mir sechs, %n.",
		"Sechs Imp-Schleier.", "Hervorragend. Das gibt Handschuhe, die jede Hitze aushalten.",
		objs(cn("x_fire_imp_2", 6)))
	nq("rurik", "n_cultist_compasses", "Kompasse der Kultisten", 37,
		"Die Kultisten haben Kompasse, die auf ihren Schrein zeigen. Bring mir fünf, %n, damit wir ihre Wege verfolgen können.",
		"Fünf Kultistenkompasse.", "Jetzt wissen wir, wohin sie wollen. Danke!",
		objs(cn("x_cultist_2", 5)))

	nq("vash", "n_cultist_grove", "Der Hain der Glutkultisten", 38,
		"Im Hain sammeln sich die Kultisten. Erkunde ihn, %n, und erlege zehn von ihnen, bevor ihr Ritual vollendet ist.",
		"Der Kultistenhain, zehn Kultisten.", "Das Ritual ist gestört. Gut gemacht.",
		objs(ex("kultistenhain"), kn("cultist", 10)))
	nq("obran", "n_magma_splitters", "Magmasplitter", 38,
		"Magmasplitter sind die härteste Zutat für meine Rüstungen. Bring mir sechs, %n.",
		"Sechs Magmasplitter.", "Genau das, was ich brauchte.",
		objs(cn("x_magma_golem_1", 6)), gear())
	nq("karr", "n_orders_to_vash", "Befehle der Kommandantin", 38,
		"Vash wartet auf die Befehle aus dem Hauptquartier. Bring sie ihr, %n, aber lass sie nicht in falsche Hände geraten.",
		"Vash wartet auf die Befehle.", "Vash hat sie erhalten? Dann kann der Plan beginnen.",
		objs(dv("vash")), give("n_ash_orders"))

	nq("vash", "n_obsidian_field", "Das Obsidianfeld", 39,
		"Im Obsidianfeld bewegen sich Ritter, die nicht sein dürften. Erkunde das Feld, %n, und erlege sechs Obsidianritter.",
		"Das Obsidianfeld, sechs Obsidianritter.", "Das war nötig. Danke.",
		objs(ex("obsidianfeld"), kn("obsidian_knight", 6)))
	nq("kael", "n_lava_splitters", "Lavasplitter", 39,
		"Lavasplitter sind ein seltener Grundstoff für meine Tränke. Bring mir sechs, %n, aber trag Handschuhe.",
		"Sechs Lavasplitter.", "Wunderbar. Damit braue ich etwas Besonderes.",
		objs(cn("x_lava_elemental_1", 6)), gear())
	nq("ignis", "n_forge_coal_run", "Kohle für Obran", 39,
		"Obran braucht Schmiedekohle für seine Rüstungen. Bring sie ihm, %n, und achte auf die Magmagolems.",
		"Obran wartet auf die Kohle.", "Obran hat die Kohle? Dann wird seine Rüstung Perfektion.",
		objs(dv("obran")), give("n_forge_coal"), items("food_dragon"))

	nq("vash", "n_last_stand", "Letzter Widerstand", 40,
		"Die Lage spitzt sich zu. Erlege zehn Obsidianritter, %n, und sichere den Weg zum Feuerschlund.",
		"Zehn Obsidianritter.", "Das hat die Reihen gelichtet. Danke.",
		objs(kn("obsidian_knight", 10)), gear())
	nq("cinder", "n_lava_crystals", "Kristalle der Glut", 40,
		"Lavaelementare tragen Splitter aus reinem Feuer. Bring mir acht, %n, sie sind der Grundstoff für den besten Schmuck.",
		"Acht Lava-Kernstücke.", "Wunderschön. Das gibt Schmuck für Könige.",
		objs(cn("x_lava_elemental_2", 8)))
	nq("obran", "n_magma_hunt", "Magmagolem-Jagd", 40,
		"Die Magmagolems werden zahlreicher und gefährlicher. Erlege acht von ihnen, %n.",
		"Acht Magmagolems.", "Das war überfällig. Danke.",
		objs(kn("magma_golem", 8)), gear())
}
