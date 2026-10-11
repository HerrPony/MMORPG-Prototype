package gamedata

// Quests for levels 20-40: Sturmhafen, Smaragddschungel, Kirschblütenhain, Rotkluft-Canyon, Aschenlande and the two new dungeons.

func init() {
	// ================= Überfahrt =================
	qd("olga", "c_sail", "Die Überfahrt", 20, 0.6,
		"Du hast die Wüste bezwungen, %n. Aber die Welt endet nicht an unserer Küste: Jenseits des Sturmmeers liegt Jadereich. Kapitän Brenner bringt dich hinüber. Melde dich bei Kapitän Wen in Lotushafen, er erwartet Abenteurer wie dich.",
		"Kapitän Wen erwartet dich in Lotushafen. Die Fähre legt bei Brenner ab.", "Willkommen in Jadereich! Hier gibt es viel zu tun.",
		[]Objective{talkTo("wen", "Mit Kapitän Wen gesprochen")}, hand("wen"), needs("g_guardian"))

	// ================= Lotushafen / Smaragddschungel (20-25) =================
	qd("lin", "j_manifest", "Der Frachtbrief", 20, 0.9,
		"Die Fracht der letzten Fähre muss bestätigt werden, und Dorfälteste Yuna in Bambushain wartet auf das Siegel. Nimm den Frachtbrief und lauf, %n! Du hast 70 Sekunden, dann ist der Handel geplatzt.",
		"Yuna wartet in Bambushain im Südosten. Beeil dich!", "Der Frachtbrief! Pünktlich. Danke dir.",
		[]Objective{deliver("yuna", "Frachtbrief bei Yuna abgegeben")}, give("ship_manifest"), timed(70), tag("Zeit"), hand("yuna"))
	qd("mei", "j_crocs", "Flusskrokodile", 21, 1.0,
		"Die Krokodile im Fluss sind eine Plage. Erlege acht von ihnen und bring mir fünf Zähne, %n. Aus den Zähnen fertige ich Fischerhaken.",
		"Die Krokodile lauern im Fluss nördlich und östlich des Hafens.", "Prächtige Zähne! Die Fische werden Schlange stehen.",
		[]Objective{kill("river_croc", 8, "Flusskrokodile getötet"), coll("river_croc_tooth", 5, "Krokodilzähne gesammelt")})
	qd("mei", "j_orchids", "Mondorchideen", 21, 0.9,
		"Die Heilerin in Bambushain braucht Mondorchideen. Sie wachsen im Schatten am Fluss. Sammle acht davon, %n.",
		"Die Orchideen wachsen in den Schatten der großen Bäume östlich des Hafens.", "Wunderschön! Hana wird sich freuen.",
		[]Objective{coll("jungle_orchid", 8, "Mondorchideen gesammelt")})
	qd("jiro", "j_tigers", "Die Tiger am Pfad", 21, 1.0,
		"Der Tigerpfad ist für Händler lebensgefährlich geworden. Erlege acht Dschungeltiger, %n. Sie greifen nicht an, solange man sie in Ruhe lässt, doch bei einem Angriff wehren sie sich.",
		"Die Tiger streifen auf dem Tigerpfad im Norden herum.", "Die Händler können wieder atmen. Ich danke dir.",
		[]Objective{kill("jungle_tiger", 8, "Dschungeltiger getötet")})
	qd("jiro", "j_panthers", "Schatten im Dickicht", 22, 1.0,
		"Schwarzpanther lauern im Unterholz und greifen Reisende an. Töte sechs von ihnen, %n.",
		"Die Panther jagen nördlich des Tigerpfads.", "Sechs Panther weniger. Gute Arbeit!",
		[]Objective{kill("panther", 6, "Schwarzpanther getötet")}, needs("j_tigers"))
	qd("lin", "j_to_bambushain", "Auf nach Bambushain", 22, 0.5,
		"Mehr Arbeit wartet im Landesinneren, %n. Geh nach Bambushain und sprich mit Dorfältester Yuna.",
		"Yuna wartet in Bambushain im Südosten.", "Lin schickt dich? Wir können Hilfe gut gebrauchen.",
		[]Objective{talkTo("yuna", "Mit Yuna gesprochen")}, hand("yuna"), needs("j_manifest"))
	qd("taro", "j_tame_tiger", "Raubkatze zähmen", 22, 1.5,
		"Ein Jäger braucht einen starken Begleiter. Bring einen Dschungeltiger dazu, dir zu folgen, %n: Wähle einen aus und wirke Tier zähmen. Du darfst dabei nicht abgelenkt werden.",
		"Wähle einen Dschungeltiger aus und wirke Tier zähmen.", "Ein Tiger an deiner Seite! Das macht dich zum echten Jäger.",
		[]Objective{tame("jungle_tiger", "Dschungeltiger gezähmt")}, forClass("hunter"), tag("Klasse"), gear())

	// ---- Bambushain ----
	qd("taro", "j_apes", "Plünderer im Affenhain", 22, 1.0,
		"Die Dschungelaffen plündern unsere Felder. Erlege acht von ihnen und bring mir fünf Felle, %n. Aus den Fellen mache ich Decken.",
		"Die Affen leben im Affenhain, nördlich des Dorfes.", "Gute Felle! Die Felder sind sicher.",
		[]Objective{kill("jungle_ape", 8, "Dschungelaffen getötet"), coll("ape_fur", 5, "Affenfelle gesammelt")})
	qd("yuna", "j_bamboo", "Frischer Bambus", 23, 0.9,
		"Wir brauchen dringend Bambus für die neuen Häuser. Sammle acht Bündel in der Umgebung des Dorfes, %n.",
		"Der Bambus wächst überall rund um Bambushain.", "Hervorragend! Daraus bauen wir ein neues Haus.",
		[]Objective{coll("bamboo_bundle", 8, "Bambusbündel gesammelt")})
	qd("yuna", "j_tea", "Teelieferung", 23, 0.9,
		"Händler Bao in Lotushafen wartet auf diese Kiste Tee. Sie muss frisch ankommen! Du hast 75 Sekunden, %n.",
		"Bao wartet in Lotushafen im Westen. Lauf!", "Der Tee ist noch frisch. Wunderbar!",
		[]Objective{deliver("bao", "Teekiste bei Bao abgegeben")}, give("tea_box"), timed(75), tag("Zeit"), hand("bao"))
	qd("taro", "j_poachers", "Die Wilderer", 23, 1.1,
		"Wilderer haben im Nordwesten ein Lager aufgeschlagen und bedrohen unsere Tiere. Erlege acht von ihnen und bring mir ihr Kassenbuch, %n.",
		"Das Wildererlager liegt im Nordwesten des Dschungels.", "Das Kassenbuch! Jetzt wissen wir, für wen sie arbeiten.",
		[]Objective{kill("poacher", 8, "Wilderer getötet"), coll("poacher_ledger", 1, "Wilderer-Kassenbuch gesammelt")})
	qd("hana", "j_snakes", "Schlangengift", 24, 1.1,
		"Für meine Gegengifte brauche ich frisches Schlangengift. Die Echsenkrieger in den Ruinen tragen es bei sich. Erlege acht und sammle vier Phiolen, %n.",
		"Die Echsenkrieger hausen in den Ruinen im Süden.", "Das reicht für viele Gegengifte. Danke!",
		[]Objective{kill("lizardman", 8, "Echsenkrieger getötet"), coll("snake_venom", 4, "Schlangengift gesammelt")})
	qd("taro", "j_poachchief", "Wilderer-Anführer Kral", 24, 2.0,
		"Kral führt die Wilderer an. Er ist stark und gut bewaffnet. Besiege ihn, %n, am besten mit Verstärkung.",
		"Kral befehligt das Wildererlager im Nordwesten.", "Kral ist gefallen! Der Dschungel atmet auf.",
		[]Objective{kill("poacher_boss", 1, "Wilderer-Anführer Kral besiegt")}, needs("j_poachers"), tag("Elite"), gear())
	qd("yuna", "j_ruins", "Die Echsenruinen", 24, 1.2,
		"Im Süden liegen uralte Ruinen, in denen sich Echsenschamanen versammeln. Sieh dir die Ruinen an und erlege vier Schamanen, %n.",
		"Die Echsenruinen liegen im Süden von Bambushain.", "Die Schamanen sind vertrieben. Vorerst.",
		[]Objective{explore("echsenruinen", "Echsenruinen erkundet"), kill("lizard_shaman", 4, "Echsenschamanen getötet")}, needs("j_bamboo"))
	qd("yuna", "j_idol", "Das Echsenidol", 25, 2.0,
		"Der Häuptling der Echsen, Szar, trägt ein heiliges Idol. Solange er es hat, werden seine Krieger nicht aufgeben. Besiege ihn und bring mir das Idol, %n.",
		"Szar haust in den Echsenruinen im Süden.", "Das Idol! Nun können wir den Echsen Frieden anbieten.",
		[]Objective{kill("lizard_chief", 1, "Echsenhäuptling Szar besiegt"), coll("lizard_idol", 1, "Echsenidol erbeutet")}, needs("j_ruins"), tag("Elite"), gear(), items("jade_amulet"))
	qd("taro", "j_whitetiger", "Der Weiße Tiger", 24, 2.0,
		"Ein riesiger weißer Tiger, Shan, terrorisiert den Tigerpfad. Töte ihn und bring mir seine Kralle als Beweis, %n.",
		"Shan streift auf dem Tigerpfad im Norden umher.", "Shan ist besiegt. Ein mächtiges Tier.",
		[]Objective{kill("whitetiger", 1, "Weißer Tiger Shan besiegt"), coll("tiger_trophy", 1, "Tigerkralle erbeutet")}, needs("j_tigers"), tag("Elite"), gear())
	qd("yuna", "j_kong", "Der Silberrücken", 25, 2.0,
		"Der alte Silberrücken Bokko beherrscht den Affenhain. Seine Wutanfälle zerstören alles. Besiege ihn, %n.",
		"Bokko lebt im Affenhain nördlich des Dorfes.", "Bokko ist besiegt! Der Hain ist wieder friedlich.",
		[]Objective{kill("ape_king", 1, "Silberrücken Bokko besiegt")}, needs("j_apes"), tag("Elite"), gear())
	qd("yuna", "j_to_capital", "Zur Himmelspforte", 25, 0.6,
		"Du hast dich im Dschungel bewährt, %n. Die Kaiserin in der Himmelspforte sucht tapfere Helden. Folge der Straße nach Osten, oder nimm einen Greifen. Sprich mit Kaiserin Mei-Lin.",
		"Mei-Lin wartet in der Himmelspforte im Osten.", "Yuna schickt dich? Dann bist du uns willkommen!",
		[]Objective{talkTo("meilin", "Mit Kaiserin Mei-Lin gesprochen")}, hand("meilin"), needs("j_ruins"))

	// ================= Himmelspforte / Kirschblütenhain (25-30) =================
	qd("meilin", "b_welcome", "Willkommen in Jadereich", 25, 0.5,
		"Die Himmelspforte öffnet ihre Tore für dich, %n. Generalin Akane kann dir sagen, wo meine Hilfe am nötigsten ist.",
		"Generalin Akane steht an der Westseite des Platzes.", "Meine Kaiserin lässt dich grüßen? Dann hör mir zu.",
		[]Objective{talkTo("akane", "Mit Generalin Akane gesprochen")}, hand("akane"))
	qd("meilin", "b_flight", "Auf Greifenschwingen", 25, 0.5,
		"Greife tragen dich zu jedem Dorf, das du kennst. Sprich mit Greifenmeister Haru am Landeplatz im Südosten der Stadt, %n, damit er sich dich merkt.",
		"Der Landeplatz liegt im Südosten der Himmelspforte.", "Haru hat dich notiert. Der Himmel gehört dir!",
		[]Objective{talkTo("haru", "Mit Greifenmeister Haru gesprochen")}, hand("meilin"))
	qd("genji", "b_incense", "Weihrauch für den Schrein", 25, 1.0,
		"Der Fuchsschrein im Norden braucht neuen Weihrauch. Besuche den Schrein und sammle sechs Schalen Weihrauch, %n.",
		"Die Weihrauchschalen stehen rund um den Fuchsschrein im Norden.", "Der Schrein ist wieder rein. Danke dir.",
		[]Objective{explore("fuchsschrein", "Fuchsschrein erkundet"), coll("temple_incense", 6, "Tempelweihrauch gesammelt")})
	qd("genji", "b_scroll", "Die Geisterrolle", 26, 0.9,
		"Diese Rolle muss sofort zu Füchsin Rin nach Kirschbach. Sie enthält einen Bann. Du hast 90 Sekunden, %n, bevor die Tinte verblasst.",
		"Rin wartet in Kirschbach im Süden. Lauf!", "Die Rolle! Gerade noch rechtzeitig.",
		[]Objective{deliver("rin", "Geisterrolle bei Rin abgegeben")}, give("spirit_scroll"), timed(90), tag("Zeit"), hand("rin"))
	qd("akane", "b_foxes", "Fuchsgeister", 26, 1.1,
		"Fuchsgeister treiben um den Schrein ihr Unwesen und werfen Feuerbälle auf Pilger. Töte acht von ihnen, %n.",
		"Die Fuchsgeister streifen um den Fuchsschrein im Norden.", "Die Pilger sind sicher. Gut gemacht.",
		[]Objective{kill("fox_spirit", 8, "Fuchsgeister getötet")})
	qd("takumi", "b_pearls", "Fuchsperlen", 26, 1.1,
		"Füchse sammeln magische Perlen, die für die Forschung unbezahlbar sind. Bring mir drei, %n. Die Geister tragen sie.",
		"Die Fuchsgeister um den Schrein tragen die Perlen.", "Faszinierend! Diese Perlen leuchten von innen.",
		[]Objective{coll("fox_pearl", 3, "Fuchsperlen gesammelt")}, needs("b_foxes"))
	qd("takumi", "b_sakura", "Kirschblüten", 27, 0.9,
		"Zum Frühlingsfest braucht der Palast Kirschblüten. Sammle zehn Haufen davon in den Hainen, %n.",
		"Die Blüten liegen unter den Bäumen im Hain südlich der Stadt.", "Prächtig! Das wird ein schönes Fest.",
		[]Objective{coll("sakura_petal", 10, "Kirschblüten gesammelt")})
	qd("akane", "b_ronin", "Die Verstoßenen", 27, 1.2,
		"Ronin, Samurai ohne Herren, plündern im Nordosten. Erlege acht von ihnen und bring mir vier Siegel, %n.",
		"Das Ronin-Lager liegt im Norden der Stadt.", "Die Siegel! Ich werde die Namen prüfen.",
		[]Objective{kill("ronin", 8, "Verstoßene Ronin getötet"), coll("ronin_seal", 4, "Ronin-Siegel gesammelt")})
	qd("goro", "b_warrior_trial", "Die Prüfung des Berges", 27, 2.0,
		"Ein Krieger muss beweisen, dass er unbeweglich wie ein Berg ist. Besiege sechs Steinwächter und bring mir drei ihrer Kerne, %n.",
		"Die Steinwächter stehen südlich der Stadt.", "Du hast den Berg bezwungen. Du bist würdig.",
		[]Objective{kill("stone_golem", 6, "Steinwächter getötet"), coll("stone_core", 3, "Steinkerne gesammelt")}, forClass("warrior"), tag("Klasse"), gear())
	qd("ryu", "b_mage_trial", "Die Prüfung der Geister", 27, 2.0,
		"Ein Magier muss die Geister beherrschen, %n. Bezwinge sechs Geisterkrieger und bring mir vier Essenzen ihrer Seele.",
		"Die Geisterkrieger spuken am Spiegelteich im Südosten.", "Die Geister haben sich gebeugt. Ein würdiger Meister.",
		[]Objective{kill("ghost_warrior", 6, "Geisterkrieger getötet"), coll("spirit_essence", 4, "Geisteressenz gesammelt")}, forClass("mage"), tag("Klasse"), gear())
	qd("sora", "b_tame_panda", "Der Bambusbär", 27, 1.5,
		"Pandas sind die sanftesten Begleiter, aber im Kampf unerschütterlich. Zähme einen Bambusbär, %n.",
		"Bambusbären leben im Bambuswald westlich der Stadt. Wirke Tier zähmen.", "Ein Panda! Er wird dich beschützen.",
		[]Objective{tame("bamboo_bear", "Bambusbär gezähmt")}, forClass("hunter"), tag("Klasse"), gear())
	qd("genji", "b_ghosts", "Ruhelose Geister", 28, 1.2,
		"Am Spiegelteich spuken die Geister gefallener Krieger. Erlöse acht von ihnen und bring mir drei Geisterglocken, %n.",
		"Die Geisterkrieger spuken am Spiegelteich im Südosten.", "Sie sind erlöst. Ihre Glocken schweigen endlich.",
		[]Objective{kill("ghost_warrior", 8, "Geisterkrieger getötet"), coll("spirit_bell", 3, "Geisterglocken gesammelt")}, needs("b_incense"))
	qd("akane", "b_foxqueen", "Fuchskönigin Kyuuko", 28, 2.2,
		"Die Fuchskönigin befehligt alle Geister des Schreins. Besiege sie, %n. Sie ist schlau und mächtig; bring Verstärkung mit.",
		"Kyuuko wacht am Fuchsschrein im Norden.", "Kyuuko ist besiegt! Der Schrein ist befreit.",
		[]Objective{kill("fox_queen", 1, "Fuchskönigin Kyuuko besiegt")}, needs("b_foxes"), tag("Elite"), gear(), items("fox_cloak"))
	qd("takumi", "b_golems", "Erwachte Wächter", 29, 1.2,
		"Die Steinwächter der alten Tempel sind erwacht und marschieren Richtung Stadt. Zerstöre sechs und sammle vier Kerne, %n.",
		"Die Steinwächter stehen südlich der Stadt, auf dem Weg nach Kirschbach.", "Die Kerne bergen uralte Magie. Danke!",
		[]Objective{kill("stone_golem", 6, "Steinwächter getötet"), coll("stone_core", 4, "Steinkerne gesammelt")})
	qd("akane", "b_roninmaster", "Ronin-Meister Takeda", 29, 2.2,
		"Takeda führt die Ronin an und ist ein Meister des Schwertes. Besiege ihn im Ronin-Lager, %n. Das schaffst du kaum allein.",
		"Takeda befehligt das Ronin-Lager im Norden.", "Takeda ist gefallen. Das Reich ist dir zu Dank verpflichtet.",
		[]Objective{kill("ronin_master", 1, "Ronin-Meister Takeda besiegt")}, needs("b_ronin"), tag("Elite"), gear())
	qd("takumi", "b_colossus", "Der Tempelkoloss", 30, 2.2,
		"Beim Nebeltempel erhebt sich ein gewaltiger Koloss aus Stein. Besiege ihn, %n, bevor er die Straßen erreicht.",
		"Der Tempelkoloss steht vor dem Nebeltempel im Nordwesten.", "Der Koloss liegt in Trümmern. Der Tempel ist zugänglich.",
		[]Objective{explore("nebeltempel", "Nebeltempel erkundet"), kill("temple_colossus", 1, "Tempelkoloss besiegt")}, needs("b_golems"), tag("Elite"), gear())
	qd("akane", "b_to_mesa", "Der Weg in die Schlucht", 30, 0.6,
		"Im Südosten beginnt der Rotkluft-Canyon. Sheriff Mara in Rotfels braucht Hilfe gegen Banditen und Schlimmeres. Sprich mit ihr, %n.",
		"Mara wartet in Rotfels im Südosten.", "Akane schickt dich? Gut, wir können jeden gebrauchen.",
		[]Objective{talkTo("mara", "Mit Sheriff Mara gesprochen")}, hand("mara"), needs("b_roninmaster"))

	// ---- Kirschbach ----
	qd("aiko", "k_tea", "Zarte Sprossen", 26, 0.9,
		"Für meinen besten Tee brauche ich junge Bambussprossen. Sammle acht im Bambuswald, %n.",
		"Der Bambuswald liegt westlich von Kirschbach.", "Wunderbar! Dieser Tee wird besonders.",
		[]Objective{coll("bamboo_shoot", 8, "Bambussprossen gesammelt")})
	qd("hiro", "k_pandas", "Bären im Bambus", 26, 1.1,
		"Die Bambusbären fressen meinen ganzen Bambus. Erlege sechs von ihnen und bring mir vier Felle, %n.",
		"Die Bären leben im Bambuswald westlich des Dorfes.", "Mein Bambus ist gerettet!",
		[]Objective{kill("bamboo_bear", 6, "Bambusbären getötet"), coll("panda_fur", 4, "Pandafelle gesammelt")})
	qd("rin", "k_foxhunt", "Die Fuchsplage", 27, 1.1,
		"Die Fuchsgeister kommen immer näher an unser Dorf. Töte acht von ihnen, %n.",
		"Die Fuchsgeister streifen um den Fuchsschrein nördlich der Hauptstadt.", "Endlich Ruhe in den Gärten.",
		[]Objective{kill("fox_spirit", 8, "Fuchsgeister getötet")})
	qd("rin", "k_ghostbells", "Das Läuten am Teich", 28, 1.1,
		"Nachts läuten Geisterglocken am Spiegelteich und lassen niemanden schlafen. Erlöse sechs Geisterkrieger, %n.",
		"Der Spiegelteich liegt östlich von Kirschbach.", "Endlich Stille! Danke dir.",
		[]Objective{explore("spiegelteich", "Spiegelteich erkundet"), kill("ghost_warrior", 6, "Geisterkrieger getötet")})

	// ---- Dungeon: Nebeltempel ----
	qd("miyu", "d2_explore", "Der Versunkene Tempel", 27, 1.8,
		"Hinter dem Portal liegt der Nebeltempel, halb im Wasser versunken. Nagas haben ihn besetzt. Wir brauchen eine Gruppe von bis zu fünf Helden, %n. Töte zehn Nagakrieger.",
		"Das Portal zum Tempel steht vor dem Nebeltempel im Nordwesten.", "Zehn Nagakrieger weniger. Das war erst der Anfang.",
		[]Objective{explore("nebeltempel", "Nebeltempel erkundet"), kill("naga_warrior", 10, "Nagakrieger getötet")}, tag("Dungeon"))
	qd("miyu", "d2_golem", "Steinwächter Gorun", 28, 2.2,
		"Im Herzen des Tempels steht Gorun, ein Wächter aus Fels. Besiege ihn, %n.",
		"Gorun steht in der zweiten Halle.", "Gorun ist gefallen. Das war ein harter Kampf.",
		[]Objective{kill("gorun", 1, "Steinwächter Gorun besiegt")}, needs("d2_explore"), tag("Dungeon"), gear())
	qd("miyu", "d2_priestess", "Hohepriesterin Ssyra", 29, 2.4,
		"Die Priesterin der Nagas ruft Verstärkung und vergiftet ihre Feinde. Besiege Ssyra, %n.",
		"Ssyra wartet in der vierten Halle.", "Ssyra ist besiegt. Nur noch der Kaiser.",
		[]Objective{kill("ssyra", 1, "Hohepriesterin Ssyra besiegt")}, needs("d2_golem"), tag("Dungeon"), gear())
	qd("miyu", "d2_emperor", "Der Nagakaiser", 30, 3.2,
		"Der Nagakaiser Zhul'Karat thront am Ende des Tempels. Besiege ihn, %n, und der Tempel gehört wieder uns. Bring deine ganze Gruppe mit!",
		"Zhul'Karat wartet im Schrein am Ende des Tempels.", "Der Kaiser ist gefallen! Du bist ein Held von Jadereich!",
		[]Objective{kill("zhulkarat", 1, "Nagakaiser Zhul'Karat besiegt")}, needs("d2_priestess"), tag("Dungeon"), gear(), items("lotus_ring"))

	// ================= Rotfels / Rotkluft-Canyon (30-35) =================
	qd("mara", "m_cats", "Sandkatzen", 30, 1.0,
		"Sandkatzen reißen unser Vieh und greifen Reisende an. Erlege acht von ihnen, %n.",
		"Die Katzen jagen beim Steinbogen im Norden.", "Die Herden sind sicher. Danke!",
		[]Objective{kill("sandcat", 8, "Sandkatzen getötet")})
	qd("olin", "m_clay", "Roter Ton", 30, 0.9,
		"Für die Töpferei brauche ich roten Ton aus den Vorkommen in den Canyonwänden. Sammle acht Klumpen, %n.",
		"Die Tonvorkommen liegen östlich des Dorfes in den Canyonwänden.", "Prächtiger Ton! Das gibt feine Töpfe.",
		[]Objective{coll("red_clay", 8, "Roter Ton gesammelt")})
	qd("olin", "m_crates", "Zerbrechliche Fracht", 31, 0.9,
		"Diese Kiste mit Tonwaren muss heil zu Späher Nox am Steinbogen. Du hast 55 Sekunden, %n, aber lauf vorsichtig!",
		"Nox wartet am Steinbogen im Norden. Beeil dich!", "Nichts gebrochen! Wunderbar.",
		[]Objective{deliver("nox", "Tonwaren bei Nox abgegeben")}, give("pottery_crate"), timed(55), tag("Zeit"), hand("nox"))
	qd("mara", "m_bison", "Die Bisonherde", 31, 1.0,
		"Die Bisons zertrampeln unsere Pfade. Erlege sechs und bring mir vier Hörner, %n.",
		"Die Bisons grasen am Kluftpass im Süden.", "Mächtige Hörner! Danke dir.",
		[]Objective{kill("canyon_bison", 6, "Canyon-Bisons getötet"), coll("bison_horn", 4, "Bisonhörner gesammelt")})
	qd("nox", "m_tame_cat", "Die Sandkatze", 31, 1.5,
		"Ein Jäger braucht Schnelligkeit. Zähme eine Sandkatze, %n. Wähle sie aus und wirke Tier zähmen.",
		"Sandkatzen jagen beim Steinbogen. Wirke Tier zähmen.", "Eine Sandkatze! Sie wird dir treu folgen.",
		[]Objective{tame("sandcat", "Sandkatze gezähmt")}, forClass("hunter"), tag("Klasse"), gear())
	qd("brokk", "m_scorp", "Skorpionpanzer", 32, 1.1,
		"Aus Skorpionpanzern schmiede ich die besten Rüstungen. Erlege acht Kluftskorpione und bring mir fünf Panzer, %n.",
		"Die Kluftskorpione hausen im Süden der Schlucht.", "Hervorragende Panzer! Du bekommst einen Anteil der Rüstung.",
		[]Objective{kill("mesa_scorpion", 8, "Kluftskorpione getötet"), coll("scorpion_chitin", 5, "Skorpionpanzer gesammelt")})
	qd("olin", "m_crystals", "Kluftkristalle", 32, 1.0,
		"Die Kristalle am Steinbogen leuchten im Sonnenlicht. Der Schmied braucht sie für seine Esse. Sammle sechs, %n.",
		"Die Kristalle liegen beim Steinbogen im Norden.", "Sie funkeln wunderbar. Danke!",
		[]Objective{coll("mesa_crystal", 6, "Kluftkristalle gesammelt")}, needs("m_clay"))
	qd("brokk", "m_matriarch", "Die Skorpionmatriarchin", 33, 2.2,
		"Tief im Süden lauert die Mutter aller Skorpione. Besiege sie, %n. Ihr Panzer ist kaum zu durchdringen.",
		"Die Matriarchin lauert im Süden der Schlucht.", "Die Matriarchin ist tot! Ein Fest für den Schmied.",
		[]Objective{kill("scorpion_matriarch", 1, "Skorpionmatriarchin besiegt")}, needs("m_scorp"), tag("Elite"), gear())
	qd("mara", "m_map", "Die Schluchtkarte", 33, 1.2,
		"Die Schluchtbanditen haben eine Karte mit geheimen Pfaden gestohlen. Erlege acht von ihnen und bring mir drei Kartenstücke zurück, %n.",
		"Die Banditenfestung liegt im Südosten der Schlucht.", "Die Karte ist vollständig! Wir können wieder planen.",
		[]Objective{kill("canyon_bandit", 8, "Schluchtbanditen getötet"), coll("mesa_map", 3, "Kartenstücke gesammelt")})
	qd("mara", "m_banditking", "Banditenkönig Rakar", 34, 2.2,
		"Rakar führt die Banditen an und plündert jede Karawane. Besiege ihn in seiner Festung, %n. Allein hast du keine Chance.",
		"Rakar haust in der Banditenfestung im Südosten.", "Rakar ist gefallen! Die Karawanen sind frei.",
		[]Objective{kill("bandit_king", 1, "Banditenkönig Rakar besiegt")}, needs("m_map"), tag("Elite"), gear(), items("titan_belt"))
	qd("brokk", "m_golems", "Lehmgolems", 34, 1.2,
		"Beim Titanengrab sind Golems aus Lehm erwacht. Zerstöre sechs und bring mir drei Herzen, %n.",
		"Das Titanengrab liegt westlich des Dorfes.", "Diese Herzen sind wertvoll. Danke!",
		[]Objective{explore("titanengrab", "Titanengrab erkundet"), kill("clay_golem", 6, "Lehmgolems getötet"), coll("golem_heart", 3, "Golemherzen gesammelt")})
	qd("olin", "m_titan", "Der Alte Titan", 35, 2.4,
		"Im Titanengrab liegt der Alte Titan, ein Riese aus Fels. Besiege ihn und bring mir eine Scherbe, %n. Nur mit einer Gruppe!",
		"Der Titan steht im Titanengrab westlich des Dorfes.", "Der Titan ist gefallen! Die Erde wird sich beruhigen.",
		[]Objective{kill("old_titan", 1, "Der Alte Titan besiegt"), coll("titan_shard", 1, "Titanenscherbe erbeutet")}, needs("m_golems"), tag("Elite"), gear())
	qd("mara", "m_to_volcano", "Zu den Aschenlanden", 35, 0.6,
		"Im Norden brennt der Himmel. Kommandantin Vash in der Aschenwacht braucht Verstärkung. Sprich mit ihr, %n.",
		"Vash wartet in der Aschenwacht im Norden.", "Mara schickt dich? Dann ist es ernst.",
		[]Objective{talkTo("vash", "Mit Kommandantin Vash gesprochen")}, hand("vash"), needs("m_banditking"))

	// ================= Aschenwacht / Aschenlande (35-40) =================
	qd("vash", "v_hounds", "Aschehunde", 35, 1.1,
		"Aschehunde streifen in Rudeln um den Pass und greifen Patrouillen an. Erlege acht von ihnen und bring mir fünf Zähne, %n.",
		"Die Hunde streifen am Aschepass im Westen.", "Die Patrouillen sind sicher. Gut gemacht.",
		[]Objective{kill("ash_hound", 8, "Aschehunde getötet"), coll("hound_tooth", 5, "Glutzähne gesammelt")})
	qd("rurik", "v_report", "Der Aschebericht", 36, 0.9,
		"Dieser Bericht muss sofort zur Kommandantin Vash. Die Kultisten bewegen sich! Du hast 60 Sekunden, %n.",
		"Vash wartet in der Aschenwacht im Osten.", "Der Bericht! Rechtzeitig. Danke.",
		[]Objective{deliver("vash", "Aschebericht bei Vash abgegeben")}, give("ash_report"), timed(60), tag("Zeit"), hand("vash"))
	qd("kael", "v_imps", "Feuerimps", 36, 1.1,
		"Die Imps schleudern Feuer auf alles, was sich bewegt. Erlege acht und bring mir fünf Hörner, %n.",
		"Die Imps hausen im Obsidianfeld im Osten.", "Gute Hörner! Das gibt starke Tränke.",
		[]Objective{kill("fire_imp", 8, "Feuerimps getötet"), coll("imp_horn", 5, "Imphörner gesammelt")})
	qd("kael", "v_flowers", "Glutblüten", 36, 1.0,
		"Glutblüten wachsen nur auf heißem Boden und bilden die Grundlage meiner besten Tränke. Sammle acht, %n.",
		"Glutblüten wachsen am Aschepass im Westen.", "Perfekt! Diese Blüten glühen noch.",
		[]Objective{coll("fire_flower", 8, "Glutblüten gesammelt")})
	qd("rurik", "v_tame_hound", "Der Aschehund", 36, 1.5,
		"Ein Jäger kann auch in der Asche überleben. Zähme einen Aschehund, %n. Wirke Tier zähmen und halte still.",
		"Aschehunde streifen am Aschepass. Wirke Tier zähmen.", "Ein Aschehund! Er wird dir nicht von der Seite weichen.",
		[]Objective{tame("ash_hound", "Aschehund gezähmt")}, forClass("hunter"), tag("Klasse"), gear())
	qd("vash", "v_cultists", "Die Glutkultisten", 37, 1.2,
		"Die Glutkultisten beschwören Feuer aus dem Schlund. Erlege neun von ihnen und bring mir drei Kultbücher, %n.",
		"Das Kultistenlager liegt im Nordwesten der Aschenlande.", "Die Bücher verraten uns viel. Gute Arbeit.",
		[]Objective{kill("cultist", 9, "Glutkultisten getötet"), coll("cult_tome", 3, "Kultbücher gesammelt")})
	qd("ignis", "v_crystals", "Glutkristalle", 37, 1.0,
		"Mein Ofen braucht Glutkristalle. Sammle acht im Obsidianfeld, %n. Pass auf die Imps auf.",
		"Die Kristalle liegen im Obsidianfeld im Osten.", "Hervorragend! Das gibt eine Hitze!",
		[]Objective{coll("ember_crystal", 8, "Glutkristalle gesammelt")})
	qd("vash", "v_ember", "Glut für den Schlund", 38, 0.9,
		"Diese Phiole muss zu Feldwebel Karr am Feuerschlund. Sie ist heiß, und der Weg ist weit. Du hast 75 Sekunden, %n.",
		"Karr wartet am Feuerschlund im Süden.", "Die Phiole! Das wird die Wache stärken.",
		[]Objective{deliver("karr", "Glutphiole bei Karr abgegeben")}, give("fire_vial"), timed(75), tag("Zeit"), hand("karr"))
	qd("vash", "v_flamelord", "Flammenfürst Pyrax", 38, 2.4,
		"Pyrax führt die Kultisten an. Sein Feuer verbrennt Rüstung und Haut. Besiege ihn im Kultistenlager, %n, und nimm Verstärkung mit.",
		"Pyrax befehligt das Kultistenlager im Nordwesten.", "Pyrax ist gefallen! Die Aschenwacht dankt dir.",
		[]Objective{kill("flame_lord", 1, "Flammenfürst Pyrax besiegt")}, needs("v_cultists"), tag("Elite"), gear())
	qd("ignis", "v_magma", "Magmagolems", 38, 1.2,
		"Magmagolems blockieren die Straße zum Schlund. Zerstöre sechs und bring mir vier Kerne, %n.",
		"Die Golems wandern am Fuß des Vulkans, im Zentrum der Aschenlande.", "Diese Kerne sind perfekt für meine Schmiede!",
		[]Objective{kill("magma_golem", 6, "Magmagolems getötet"), coll("magma_core", 4, "Magmakerne gesammelt")}, needs("v_crystals"))
	qd("kael", "v_lava", "Lavaelementare", 39, 1.2,
		"Aus der Lava steigen Elementare. Besiege sechs und bring mir zwei Lavaherzen, %n.",
		"Die Elementare steigen aus den Lavaseen im Osten.", "Lavaherzen! Wer hätte gedacht, dass ich sie je bekomme.",
		[]Objective{kill("lava_elemental", 6, "Lavaelementare getötet"), coll("lava_heart", 2, "Lavaherzen gesammelt")}, needs("v_flowers"))
	qd("vash", "v_knights", "Obsidianritter", 39, 1.3,
		"Verfluchte Ritter in Obsidianrüstung patrouillieren zwischen Pass und Schlund. Besiege sechs und bring mir drei Wappen, %n.",
		"Die Ritter patrouillieren südlich der Aschenwacht.", "Die Wappen! Wir werden ihre Herkunft klären.",
		[]Objective{kill("obsidian_knight", 6, "Obsidianritter getötet"), coll("knight_emblem", 3, "Obsidian-Wappen gesammelt")})
	qd("ignis", "v_titan", "Magmatitan Brokkar", 39, 2.4,
		"Der Titan aus Magma erhebt sich aus den Lavaseen. Besiege ihn, %n. Nur mit einer starken Gruppe!",
		"Brokkar steht im Zentrum der Aschenlande, in der Nähe des Schlunds.", "Brokkar ist gefallen! Die Erde bebt nicht mehr.",
		[]Objective{kill("magma_titan", 1, "Magmatitan Brokkar besiegt")}, needs("v_magma"), tag("Elite"), gear())
	qd("vash", "v_obsidianlord", "Obsidianwächter Zeth", 40, 2.6,
		"Zeth ist der Anführer der verfluchten Ritter. Besiege ihn im Obsidianfeld, %n, und die Aschenlande sind frei.",
		"Zeth wacht im Obsidianfeld im Osten.", "Zeth ist gefallen! Du hast die Aschenlande befreit.",
		[]Objective{kill("obsidian_lord", 1, "Obsidianwächter Zeth besiegt")}, needs("v_knights"), tag("Elite"), gear(), items("ember_pendant"))

	// ---- Dungeon: Feuerschlund ----
	qd("karr", "d3_explore", "Der Feuerschlund", 37, 1.8,
		"Aus dem Portal steigt Rauch. Der Schlund ist von Kultisten besetzt. Wir brauchen eine Gruppe von bis zu fünf Helden, %n. Töte zehn Aschekultisten.",
		"Das Portal zum Feuerschlund liegt im Süden der Aschenlande.", "Zehn Kultisten weniger. Der Weg ist frei.",
		[]Objective{explore("feuerschlund", "Feuerschlund erkundet"), kill("ember_cultist", 10, "Aschekultisten getötet")}, tag("Dungeon"))
	qd("karr", "d3_smith", "Schmiedemeister Vulkor", 38, 2.4,
		"In der großen Schmiede wacht Vulkor, ein Golem aus lebendigem Feuer. Besiege ihn, %n.",
		"Vulkor steht in der zweiten Halle.", "Vulkor ist besiegt. Die Esse erlischt.",
		[]Objective{kill("vulkor", 1, "Schmiedemeister Vulkor besiegt")}, needs("d3_explore"), tag("Dungeon"), gear())
	qd("karr", "d3_witch", "Aschehexe Kaliss", 39, 2.6,
		"Kaliss beschwört Feuer und Imps. Besiege sie, %n, bevor sie das Portal öffnet.",
		"Kaliss wartet in der vierten Halle.", "Kaliss ist gefallen. Nur noch der Drachenfürst.",
		[]Objective{kill("kaliss", 1, "Aschehexe Kaliss besiegt")}, needs("d3_smith"), tag("Dungeon"), gear())
	qd("karr", "d3_dragon", "Drachenfürst Ignaros", 40, 3.4,
		"Ignaros, der Drachenfürst, thront am Grund des Schlunds und speit Feuer. Besiege ihn, %n. Bring deine ganze Gruppe mit! Das ist die letzte Prüfung.",
		"Ignaros wartet im Inferno am Ende des Schlunds.", "Der Drachenfürst ist besiegt! Du bist ein Held der Aschenlande!",
		[]Objective{kill("ignaros", 1, "Drachenfürst Ignaros besiegt")}, needs("d3_witch"), tag("Dungeon"), gear(), items("dragon_ring"))
}
