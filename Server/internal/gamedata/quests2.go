package gamedata

// ---- builder helpers ----

type qopt func(*Quest)

func tag(t string) qopt        { return func(q *Quest) { q.Tag = t } }
func timed(sec int) qopt       { return func(q *Quest) { q.TimeLimit = sec } }
func give(item string) qopt    { return func(q *Quest) { q.GiveItem = item } }
func hand(npc string) qopt     { return func(q *Quest) { q.Hand = npc } }
func needs(id string) qopt     { return func(q *Quest) { q.Requires = id } }
func forClass(c string) qopt   { return func(q *Quest) { q.Class = c } }
func grant(s ...string) qopt   { return func(q *Quest) { q.GrantSkill = s } }
func learns(s ...string) qopt  { return func(q *Quest) { q.Reward.Skills = s } }
func items(ids ...string) qopt { return func(q *Quest) { q.Reward.Items = ids } }
func gear() qopt               { return func(q *Quest) { q.Reward.Gear = true } }

func kill(target string, n int, text string) Objective {
	return Objective{Kind: "kill", Target: target, Count: n, Text: text}
}
func coll(item string, n int, text string) Objective {
	return Objective{Kind: "collect", Target: item, Count: n, Text: text}
}
func deliver(npc, text string) Objective {
	return Objective{Kind: "deliver", Target: npc, Count: 1, Text: text}
}
func talkTo(npc, text string) Objective {
	return Objective{Kind: "talk", Target: npc, Count: 1, Text: text}
}
func explore(poi, text string) Objective {
	return Objective{Kind: "explore", Target: poi, Count: 1, Text: text}
}
func tame(mob, text string) Objective {
	return Objective{Kind: "tame", Target: mob, Count: 1, Text: text}
}

// qd registers a quest. w scales the XP/money reward (1 = a normal quest of that level).
func qd(giver, id, title string, lvl int, w float64, text, progress, complete string, objs []Objective, opts ...qopt) {
	q := &Quest{ID: id, Title: title, MinLevel: lvl, Level: lvl, Text: text, Progress: progress, Complete: complete,
		Objectives: objs, Giver: giver}
	q.Reward = Reward{XP: questXP(lvl, w), Copper: questCopper(lvl, w)}
	for _, o := range opts {
		o(q)
	}
	Quests[id] = q
	if n := NPCByID(giver); n != nil {
		n.Quests = append(n.Quests, id)
	}
	if q.Hand != "" && q.Hand != giver {
		if n := NPCByID(q.Hand); n != nil {
			n.TurnIn = append(n.TurnIn, id)
		}
	}
}

func init() {
	// default zone/level info for the original quests
	for id, q := range Quests {
		q.Level = q.MinLevel
		q.Giver = ""
		_ = id
	}
	for _, n := range NPCs {
		for _, id := range n.Quests {
			if q := Quests[id]; q != nil {
				q.Giver = n.ID
			}
		}
	}

	// ================= Eichenhain / Erlenwald (1-6) =================
	qd("forester", "e_scout", "Der Kundschafter", 1, 0.8,
		"Ich kenne den Wald wie meine Westentasche, aber die Augen werden schlechter. Sieh dir für mich drei Orte an: den Wolfshügel im Osten, die Keilerwiese im Westen und den Spiegelsee im Südosten, %n. Einfach hinlaufen, dann weiß ich, dass dort alles in Ordnung ist.",
		"Du musst noch alle drei Orte besucht haben.", "Gut gemacht! Jetzt weiß ich, wo ich als Nächstes nach dem Rechten sehen muss.",
		[]Objective{explore("wolfshuegel", "Wolfshügel erkundet"), explore("keilerwiese", "Keilerwiese erkundet"), explore("spiegelsee", "Spiegelsee erkundet")})
	qd("forester", "e_letter", "Eilbrief für Torben", 3, 0.9,
		"Torben im Osten muss unbedingt noch heute von den Wölfen erfahren. Nimm diesen Brief und lauf, so schnell du kannst, %n! Du hast nur 40 Sekunden, bevor der Wind das Siegel löst.",
		"Torben wartet auf den Brief! Beeile dich!", "Der Brief ist da? Gerade noch rechtzeitig. Danke!",
		[]Objective{deliver("torben", "Brief an Torben überbracht")}, give("letter_aldric"), timed(40), needs("wolf_plague"), tag("Zeit"), hand("forester"))
	qd("hilda", "e_pumpkin", "Verstreute Ernte", 2, 0.9,
		"Der Sturm hat meine Kürbisse über das ganze Feld verteilt! Sammle acht davon ein, die überall zwischen den Beeten liegen, %n. Du erkennst sie an den orangen Kugeln im Gras.",
		"Mir fehlen noch Kürbisse. Sie liegen verstreut neben den Beeten.", "Meine schönen Kürbisse! Hier, für die Mühe.",
		[]Objective{coll("pumpkin", 8, "Kürbisse gesammelt")})
	qd("torben", "e_bear_elder", "Der Schrecken der Klamm", 6, 1.8,
		"Seit Wochen terrorisiert ein riesiger Bär die Bärenklamm: Alter Grimmbart. Er ist weit größer und zäher als jeder Bär, den du bisher gesehen hast. Allein ist er kaum zu schaffen, %n – nimm einen Freund mit oder bring viele Heiltränke.",
		"Grimmbart lebt noch! Er haust in der Bärenklamm.", "Du hast Grimmbart erlegt! Das hat noch niemand allein geschafft.",
		[]Objective{kill("bear_elder", 1, "Alter Grimmbart besiegt")}, needs("bear_hunt"), tag("Elite"), gear())
	qd("forester", "e_to_capital", "Eine Reise nach Sonnenburg", 5, 0.7,
		"Du hast dich hier bewährt, %n. In der Hauptstadt Sonnenburg, weit im Nordosten, gibt es mehr Arbeit für Leute wie dich. Bring Bürgermeisterin Alena diesen Brief. Folge einfach der Straße nach Osten. Oder sprich mit Kira bei den Greifen, aber das kostet etwas.",
		"Alena wartet auf deinen Brief. Folge der Straße nach Osten.", "Aldric schickt dich? Dann bist du uns willkommen!",
		[]Objective{deliver("alena", "Brief an Alena überbracht")}, give("letter_alena"), needs("wolf_plague"), hand("alena"))

	// ================= Sonnenburg / Goldene Ebenen (6-10) =================
	qd("alena", "p_welcome", "Willkommen in Sonnenburg", 6, 0.5,
		"Sonnenburg ist die Hauptstadt des Reiches, %n. Hauptmann Varus am Südtor kann dir sagen, wo es gefährlich ist. Sprich bitte mit ihm.",
		"Hauptmann Varus steht am Südtor.", "Varus lässt grüßen? Dann kennst du dich jetzt aus.",
		[]Objective{talkTo("varus", "Mit Hauptmann Varus gesprochen")}, hand("alena"))
	qd("alena", "p_flight", "Auf Greifenschwingen", 6, 0.5,
		"Unsere Greife bringen dich zu jeder Stadt, die du schon einmal besucht hast, %n. Sprich mit Flugmeister Bran am Landeplatz im Osten der Stadt. Er merkt sich dann, dass du Sonnenburg kennst, und du kannst von jetzt an von jedem Landeplatz hierher fliegen.",
		"Der Landeplatz liegt im Osten der Stadt, bei den Greifen.", "Bran hat dich notiert! Von nun an stehen dir die Lüfte offen.",
		[]Objective{talkTo("bran", "Mit Flugmeister Bran gesprochen")}, hand("alena"))
	qd("odo", "p_hyenas", "Hyänenplage", 6, 1.0,
		"Steppenhyänen reißen mein Vieh, %n. Sie jagen im Rudel und greifen sofort an. Erlege acht von ihnen, bevor sie noch mehr Schaden anrichten.",
		"Die Hyänen sind noch da draußen. Acht Stück!", "Endlich Ruhe auf der Weide! Nimm das hier.",
		[]Objective{kill("hyena", 8, "Steppenhyänen getötet")})
	qd("odo", "p_wheat", "Die Weizenernte", 6, 0.9,
		"Die Ernte ist reif, aber mir fehlen Hände. Sammle acht Weizenbündel von den Feldern um meinen Hof, %n.",
		"Weizenbündel liegen auf den Feldern bei meinem Hof.", "Prächtig! Daraus wird das beste Brot der Ebenen.",
		[]Objective{coll("wheat", 8, "Weizenbündel gesammelt")})
	qd("hannes", "p_mill", "Die Mühle steht still", 7, 1.1,
		"Wegelagerer haben die Zahnräder meiner Windmühle gestohlen! Hol dir vier davon zurück, indem du die Räuber erledigst – sechs von ihnen sollten es sein, %n.",
		"Ohne die Zahnräder steht die Mühle still. Die Wegelagerer treiben sich im Osten herum.", "Meine Mühle läuft wieder! Wie kann ich dir danken?",
		[]Objective{kill("plains_bandit", 6, "Wegelagerer getötet"), coll("mill_gear", 4, "Mühlenzahnräder gesammelt")})
	qd("alena", "p_order", "Königlicher Befehl", 7, 0.9,
		"Dieser Befehl muss sofort zu Müller Hannes auf dem Windmühlenhügel im Südwesten! Du hast 70 Sekunden, %n. Beeil dich, es geht um die Versorgung der Stadt!",
		"Hannes erwartet den Befehl am Windmühlenhügel. Schnell!", "Der Befehl ist angekommen. Gute Arbeit, Läufer!",
		[]Objective{deliver("hannes", "Befehl an Hannes überbracht")}, give("order_royal"), timed(70), tag("Zeit"), hand("alena"))
	qd("hannes", "p_bulls", "Wilde Stiere", 8, 1.0,
		"Wildstiere zertrampeln meine Felder. Töte sechs und bring mir vier Hörner, damit ich beweisen kann, dass es getan ist.",
		"Die Stiere stehen draußen auf den Weiden.", "Gewaltige Hörner! Das hast du gut gemacht.",
		[]Objective{kill("bull", 6, "Wildstiere getötet"), coll("bull_horn", 4, "Stierhörner gesammelt")}, needs("p_mill"))
	qd("selene", "p_pack", "Das Rudel am Jägerhain", 8, 1.0,
		"Die Hyänen am Jägerhain werden dreister, %n. Töte sechs von ihnen, damit das Wild sich wieder erholt.",
		"Noch sind die Hyänen zu zahlreich.", "Das Rudel ist gebrochen. Danke dir.",
		[]Objective{kill("hyena", 6, "Hyänen am Jägerhain getötet")})
	qd("ilyana", "p_holy", "Das Wasser der Quelle", 8, 1.0,
		"Das Quellheiligtum verliert seine Kraft. Besuche die Quelle und sammle sechs Gefäße heiliges Wasser, damit ich es reinigen kann.",
		"Das Quellwasser sprudelt an den Teichen um das Heiligtum.", "Das reine Wasser wird die Quelle stärken.",
		[]Objective{explore("quellheiligtum", "Quellheiligtum besucht"), coll("spring_water", 6, "Quellwasser gesammelt")})
	qd("varus", "p_raiders", "Die Räuberfestung", 9, 1.8,
		"Die Räuber haben im Osten eine Festung errichtet und überfallen die Händler. Ihr Anführer, Hauptmann Dorn, ist ein harter Brocken. Töte acht seiner Männer und dann Dorn selbst, %n. Allein ist das gefährlich: Such dir Verstärkung.",
		"Die Räuberfestung liegt im Osten. Dorn ist noch nicht tot.", "Die Straßen sind sicherer. Sonnenburg dankt dir!",
		[]Objective{kill("plains_bandit", 8, "Räuber getötet"), kill("bandit_chief", 1, "Hauptmann Dorn besiegt")}, tag("Elite"), gear())
	qd("ilyana", "p_elemental", "Die Verzerrung", 10, 1.8,
		"Etwas Dunkles hat den Wasserelementar am Heiligtum verdorben. Er greift jeden an, der sich nähert. Bring ihn zur Ruhe, %n – aber unterschätze ihn nicht.",
		"Der verzerrte Elementar wacht am Heiligtum.", "Das Wasser ist wieder rein. Du hast mir sehr geholfen.",
		[]Objective{kill("water_corrupt", 1, "Verzerrter Wasserelementar besiegt")}, tag("Elite"), gear(), needs("p_holy"))
	qd("alena", "p_swamp", "Der Ruf der Moore", 10, 0.6,
		"Im Südwesten, im Düstermoor, braucht man Hilfe. Dorfältester Moros von Moorbruck sucht tapfere Abenteurer. Sprich mit ihm, %n. Die Straße führt vom Eichenhain nach Südwesten, oder du fliegst von Sonnenburg nach Eichenhain und gehst von dort.",
		"Moros wartet in Moorbruck im Südwesten.", "Alena schickt dich? Wir können jede Hilfe gebrauchen.",
		[]Objective{talkTo("moros", "Mit Moros gesprochen")}, hand("moros"))

	// ---- Klassenquests (Stufe 10) ----
	qd("rodrik", "w_trial", "Die Probe der Standhaftigkeit", 10, 2.0,
		"Ein Krieger, der nicht standhält, ist nur ein Mann mit einem Schwert, %n. Ich will sehen, wie du dich schlägst. Die Räuber im Osten tragen die Schildreste ihrer Opfer. Töte acht Wegelagerer und bring mir vier der Reste.",
		"Die Räuber liegen im Osten. Vier Schildreste, acht Wegelagerer.", "Gut. Du hast den Kampfgeist. Aber Mut allein reicht nicht.",
		[]Objective{kill("plains_bandit", 8, "Wegelagerer getötet"), coll("shield_shards", 4, "Schildreste gesammelt")}, forClass("warrior"), tag("Klasse"))
	qd("rodrik", "w_forge", "Das Eisenholzschild", 10, 1.5,
		"Aus den Resten lässt sich ein Schild schmieden. Bring diese Platte so schnell du kannst zu Waffenmeister Thoren, solange das Eisen noch glüht. Du hast 90 Sekunden! Danach kehre zu mir zurück.",
		"Thoren wartet auf die Platte. Die Zeit läuft!", "Das Schild wird bald fertig sein. Eine letzte Aufgabe bleibt.",
		[]Objective{deliver("thoren", "Schildplatte zu Thoren gebracht")}, give("shield_plate"), timed(90), needs("w_trial"), forClass("warrior"), tag("Klasse"), hand("rodrik"))
	qd("rodrik", "w_stand", "Der Wall von Sonnenburg", 10, 3.0,
		"Die letzte Prüfung: Hauptmann Dorn, der Anführer der Räuber. Er schlägt hart, aber ein echter Krieger hält jeden Schlag aus. Besiege ihn, und ich lehre dich die Verteidigungshaltung und den Spott.",
		"Dorn haust in der Räuberfestung. Er ist hartnäckig.", "Du hast es geschafft! Hier, ich lehre dich die Verteidigungshaltung und den Spott. Wer andere beschützen kann, ist ein wahrer Krieger.",
		[]Objective{kill("bandit_chief", 1, "Hauptmann Dorn besiegt")}, needs("w_forge"), forClass("warrior"), tag("Klasse"), learns("w_stance_def", "w_taunt"), gear())
	qd("aldwin", "h_tracks", "Fährtenleser", 10, 1.5,
		"Ein Jäger muss die Wildnis lesen können, %n. Geh zum Jägerhain und sammle drei Wildfährten ein, die dort überall im Gras zu finden sind.",
		"Die Fährten liegen im Jägerhain, im Nordwesten der Ebenen.", "Gut! Aus diesen Fährten kann ich viel lesen.",
		[]Objective{explore("jaegerhain", "Jägerhain erkundet"), coll("track", 3, "Wildfährten gesammelt")}, forClass("hunter"), tag("Klasse"))
	qd("aldwin", "h_tame", "Der Alpha-Wolf", 10, 2.5,
		"Ein Jäger ist nie allein. Bei Selene im Jägerhain streift ein Alpha-Wolf namens Rotpelz umher. Ich lehre dich Tier zähmen: Wähle ihn aus, wirke den Zauber und bleib 6 Sekunden ruhig stehen. Wenn er dich lässt, gehört er dir.",
		"Rotpelz streift im Jägerhain umher. Wähle ihn aus und wirke Tier zähmen.", "Du hast ihn gezähmt! Rotpelz ist jetzt dein Gefährte.",
		[]Objective{tame("wolf_alpha", "Alpha-Wolf Rotpelz gezähmt")}, needs("h_tracks"), forClass("hunter"), tag("Klasse"), grant("h_tame"))
	qd("aldwin", "h_bond", "Das Band der Wildnis", 10, 2.0,
		"Ein Tier und sein Jäger kämpfen als eines. Zieh mit Rotpelz los und erlege sechs Hyänen. Danach lehre ich dich, dein Tier jederzeit zu rufen.",
		"Ruf dein Tier und töte die Hyänen gemeinsam mit ihm.", "Ihr seid ein gutes Team! Hier, ich zeige dir, wie du dein Tier jederzeit rufen kannst.",
		[]Objective{kill("hyena", 6, "Hyänen getötet")}, needs("h_tame"), forClass("hunter"), tag("Klasse"), learns("h_callpet"), gear())
	qd("ilvane", "m_crystals", "Quelle der Macht", 10, 1.5,
		"Das Wasser ist das älteste der Elemente, %n. Im Quellheiligtum kristallisiert sich seine Macht. Geh dorthin und sammle vier Wasserkristalle.",
		"Die Kristalle liegen am Quellheiligtum im Norden.", "Wunderschön. Ich spüre ihre Kraft schon von hier.",
		[]Objective{coll("crystal_water", 4, "Wasserkristalle gesammelt")}, forClass("mage"), tag("Klasse"))
	qd("ilvane", "m_phial", "Die eilige Phiole", 10, 1.5,
		"Diese Phiole enthält die Essenz der Kristalle. Bring sie zu Priesterin Ilyana am Quellheiligtum, bevor sie ihre Kraft verliert. 100 Sekunden, %n! Dann kehre zu mir zurück.",
		"Ilyana erwartet die Phiole am Quellheiligtum.", "Ilyana hat die Essenz erhalten. Nun fehlt nur noch der letzte Schritt.",
		[]Objective{deliver("ilyana", "Phiole zu Ilyana gebracht")}, give("phial_water"), timed(100), needs("m_crystals"), forClass("mage"), tag("Klasse"), hand("ilvane"))
	qd("ilvane", "m_heart", "Das Herz des Wassers", 10, 3.0,
		"Der verzerrte Elementar am Heiligtum trägt das Herz des Wassers in sich. Bezwinge ihn und bring mir das Herz. Dann lehre ich dich, einen Wasserelementar zu beschwören.",
		"Der verzerrte Elementar wacht am Quellheiligtum.", "Das Herz des Wassers! Hier, ich lehre dich, einen Wasserelementar zu beschwören. Er wird dir treu dienen.",
		[]Objective{kill("water_corrupt", 1, "Verzerrter Wasserelementar besiegt"), coll("crystal_heart", 1, "Herz des Wassers erbeutet")}, needs("m_phial"), forClass("mage"), tag("Klasse"), learns("m_elemental"), gear())

	// ================= Moorbruck / Düstermoor (10-14) =================
	qd("moros", "s_croc", "Sumpfkrokodile", 10, 1.0,
		"Die Krokodile sind aggressiver geworden und gefährden unsere Fischer. Töte acht von ihnen, %n. Sie lauern in den Sümpfen im Süden.",
		"Die Krokodile lauern in den Sümpfen im Süden.", "Danke! Jetzt trauen sich die Fischer wieder ans Wasser.",
		[]Objective{kill("croc", 8, "Sumpfkrokodile getötet")})
	qd("sylva", "s_herbs", "Heilkräuter", 10, 1.0,
		"Mir gehen die Heilkräuter aus. Sammle acht Moorlotusblüten, die überall in den Sümpfen wachsen, %n. Sie schwimmen im flachen Wasser.",
		"Moorlotus wächst überall im Moor, besonders im flachen Wasser.", "Wunderbare Blüten! Daraus braue ich Heiltränke.",
		[]Objective{coll("lotus", 8, "Moorlotus gesammelt")})
	qd("hektor", "s_zombies", "Der Zombie-Ansturm", 11, 1.1,
		"Nachts steigen Zombies aus dem Moor. Töte zehn von ihnen, damit wir wieder ruhig schlafen können, %n.",
		"Die Zombies wandern östlich und westlich des Dorfes.", "Das war dringend nötig. Danke!",
		[]Objective{kill("zombie", 10, "Moorzombies getötet")})
	qd("sylva", "s_medicine", "Medizin für Tam", 11, 0.9,
		"Fischer Tam am Nebelsee ist krank! Nimm diese Medizin und bring sie ihm, aber schnell: Du hast nur 45 Sekunden, sonst verliert sie ihre Wirkung.",
		"Tam wartet am Nebelsee im Süden. Beeil dich!", "Tam ist wieder gesund? Das ist gut zu hören.",
		[]Objective{deliver("tam", "Medizin zu Tam gebracht")}, give("medicine_vial"), timed(45), tag("Zeit"), hand("sylva"), needs("s_herbs"))
	qd("tam", "s_wisps", "Irrlichter", 12, 1.1,
		"Irrlichter locken Fischer ins tiefe Wasser. Erlege acht von ihnen und bring mir fünf Beutel Irrlichtstaub als Beweis, %n.",
		"Irrlichter schweben über dem Nebelsee.", "Der Staub glitzert wunderschön. Danke dir.",
		[]Objective{kill("wisp", 8, "Irrlichter getötet"), coll("wisp_dust", 5, "Irrlichtstaub gesammelt")})
	qd("hektor", "s_witches", "Die Hexen von Hexenhain", 12, 1.2,
		"Im Hexenhain im Westen brauen Hexen finstere Tränke. Sieh dir ihr Lager an und töte acht von ihnen, %n.",
		"Das Hexenhain liegt im Westen des Moores.", "Die Hexen sind in die Schranken gewiesen.",
		[]Objective{explore("hexenhain", "Hexenhain erkundet"), kill("witch", 8, "Moorhexen getötet")})
	qd("hektor", "s_hag", "Mutter Grimmel", 13, 2.0,
		"Die Anführerin der Hexen, Mutter Grimmel, verflucht das ganze Moor. Sie ist stark und wirft Feuerbälle. Besiege sie, %n, am besten mit Verstärkung.",
		"Grimmel haust im Hexenhain.", "Das Moor atmet auf. Nimm diesen Dank.",
		[]Objective{kill("hag", 1, "Mutter Grimmel besiegt")}, needs("s_witches"), tag("Elite"), gear())
	qd("moros", "s_tower", "Der Schiefe Turm", 13, 2.0,
		"Im Osten des Moores steht ein verfallener Turm. Dort hat sich ein riesiger Moorkoloss festgesetzt. Sieh dir den Turm an und besiege das Ungetüm, %n.",
		"Der Schiefe Turm liegt östlich von Moorbruck.", "Der Koloss ist besiegt! Moorbruck wird dich nie vergessen.",
		[]Objective{explore("schiefer_turm", "Schiefen Turm erkundet"), kill("bog_colossus", 1, "Alter Moorkoloss besiegt")}, tag("Elite"), gear(), needs("s_croc"))
	qd("moros", "s_north", "Nach Frostwacht", 14, 0.6,
		"Im hohen Norden liegt Frostwacht. Hauptmann Isgard braucht Hilfe gegen Yetis und Schlimmeres. Sprich mit ihm, %n. Von Moorbruck kannst du nach Eichenhain fliegen und der Nordstraße folgen.",
		"Isgard wartet in Frostwacht im Norden.", "Moros schickt dich? Dann sieht es nicht gut aus.",
		[]Objective{talkTo("isgard", "Mit Hauptmann Isgard gesprochen")}, hand("isgard"))

	// ================= Frostwacht / Eisgipfel (14-18) =================
	qd("olaf", "f_wolves", "Eiswölfe", 14, 1.0,
		"Die Eiswölfe reißen unser Vieh. Töte acht und bring mir sechs Pelze. Daraus nähen wir Winterkleidung, %n.",
		"Die Eiswölfe streifen durch den Schnee westlich und südlich der Stadt.", "Gute Pelze! Das reicht für den ganzen Winter.",
		[]Objective{kill("icewolf", 8, "Eiswölfe getötet"), coll("icewolf_pelt", 6, "Eiswolfpelze gesammelt")})
	qd("isgard", "f_yeti", "Yeti-Gefahr", 15, 1.1,
		"Die Yetis sind in letzter Zeit viel zu nahe herangekommen. Töte acht von ihnen, %n.",
		"Die Yetis lauern in den Bergen nördlich der Stadt.", "Die Yetis werden uns eine Weile in Ruhe lassen.",
		[]Objective{kill("yeti", 8, "Yetis getötet")}, needs("f_wolves"))
	qd("isgard", "f_message", "Botschaft für den Windpass", 15, 0.9,
		"Diese Nachricht muss sofort zu Gelehrtem Vorn am Windpass im Südosten. Du hast nur 70 Sekunden, %n, sonst ist der Pass zugeschneit!",
		"Vorn wartet am Windpass im Südosten. Beeil dich!", "Die Nachricht ist angekommen. Danke!",
		[]Objective{deliver("vorn", "Botschaft bei Vorn abgeliefert")}, give("pass_message"), timed(70), tag("Zeit"), hand("isgard"))
	qd("vorn", "f_spirits", "Eisgeister", 16, 1.1,
		"Die Kälte hat Geister geformt. Faszinierend, aber gefährlich! Erlege acht und sammle fünf Eissplitter für meine Studien, %n.",
		"Die Eisgeister schweben nördlich der Stadt.", "Hervorragende Proben! Das wird die Wissenschaft voranbringen.",
		[]Objective{kill("frostspirit", 8, "Eisgeister getötet"), coll("ice_shard", 5, "Eissplitter gesammelt")})
	qd("olaf", "f_bears", "Eisbärenjagd", 16, 1.1,
		"Die Eisbären an der Klippe werden zu viele. Erlege sechs davon, %n.",
		"Die Eisbären hausen an der Eisbärenklippe.", "Das hilft uns viel. Danke!",
		[]Objective{kill("polarbear", 6, "Eisbären getötet")}, needs("f_wolves"))
	qd("olaf", "f_mother", "Die Mutter der Klippe", 16, 2.0,
		"Die Eisbärenmutter ist größer und wütender als alle anderen. Besiege sie, %n. Vorsicht: Sie ist ein Elite-Gegner.",
		"Die Eisbärenmutter wacht an der Klippe.", "Das war eine große Tat!",
		[]Objective{kill("ice_mother", 1, "Eisbärenmutter besiegt")}, needs("f_bears"), tag("Elite"), gear())
	qd("isgard", "f_chief", "Häuptling Grimmhorn", 17, 2.2,
		"Der Yeti-Häuptling Grimmhorn führt die Angriffe an. Besiege ihn in der Yetihöhle, %n. Das schaffst du kaum allein.",
		"Grimmhorn haust in der Yetihöhle im Norden.", "Grimmhorn ist gefallen! Frostwacht ist gerettet.",
		[]Objective{explore("yetihoehle", "Yetihöhle erkundet"), kill("yeti_chief", 1, "Yeti-Häuptling Grimmhorn besiegt")}, needs("f_yeti"), tag("Elite"), gear())
	qd("isgard", "f_desert", "Reise in die Wüste", 17, 0.6,
		"Im Südosten, in der Glutwüste, liegt Sandfurt. Scheich Karim sucht Verstärkung. Sprich mit ihm, %n.",
		"Karim wartet in Sandfurt im Südosten.", "Isgard schickt dich? Willkommen, Fremder.",
		[]Objective{talkTo("karim", "Mit Scheich Karim gesprochen")}, hand("karim"))

	// ---- Dungeon: Schwarzfelsmine ----
	qd("dolf", "d_kobolds", "Die Schwarzfelsmine", 16, 1.6,
		"Kobolde haben die alte Schwarzfelsmine besetzt. Wir brauchen eine Gruppe von bis zu fünf Helden, %n. Betritt die Mine durch das Portal und töte zehn Kobolde. Allein hast du dort keine Chance!",
		"Die Mine liegt gleich hinter dem Portal bei Rabenfels.", "Zehn Kobolde weniger. Gut gemacht!",
		[]Objective{explore("rabenfels", "Rabenfels erkundet"), kill("kobold", 10, "Minenkobolde getötet")}, tag("Dungeon"))
	qd("dolf", "d_foreman", "Vorarbeiter Grimm", 17, 2.0,
		"Der Vorarbeiter treibt die Kobolde an und ist ein gewaltiger Brocken. Besiege Vorarbeiter Grimm in der Mine, %n.",
		"Grimm wartet im Herzen der ersten Höhle.", "Grimm ist erledigt! Das war kein leichter Kampf.",
		[]Objective{kill("foreman", 1, "Vorarbeiter Grimm besiegt")}, needs("d_kobolds"), tag("Dungeon"), gear())
	qd("dolf", "d_queen", "Die Spinnenkönigin", 18, 2.2,
		"Tief in der Mine haust Vessa, die Spinnenkönigin. Ihr Gift ist tödlich und ihre Brut zahllos. Besiege sie, %n!",
		"Vessa haust in einer der tieferen Höhlen.", "Vessa ist besiegt! Die Mine ist ein Stück sicherer.",
		[]Objective{kill("spider_queen", 1, "Spinnenkönigin Vessa besiegt")}, needs("d_foreman"), tag("Dungeon"), gear())
	qd("dolf", "d_king", "Der Kobold-König", 19, 3.0,
		"Der Kobold-König Skarn thront am Ende der Mine. Besiege ihn, %n, und die Schwarzfelsmine gehört wieder uns. Bring deine ganze Gruppe mit!",
		"Skarn wartet in der tiefsten Halle.", "Skarn ist besiegt! Die Schwarzfelsmine ist frei. Du bist ein Held!",
		[]Objective{kill("kobold_king", 1, "Kobold-König Skarn besiegt")}, needs("d_queen"), tag("Dungeon"), gear())

	// ================= Sandfurt / Glutwüste (17-20) =================
	qd("karim", "g_scorpions", "Die Skorpionplage", 17, 1.1,
		"Skorpione überrennen unsere Karawanenwege. Töte zehn von ihnen, %n. Hüte dich vor ihrem Gift.",
		"Die Skorpione wimmeln in den Dünen östlich der Stadt.", "Endlich wieder sichere Wege!",
		[]Objective{kill("scorpion", 10, "Sandskorpione getötet")})
	qd("jasim", "g_caravan", "Eilige Gewürze", 17, 0.9,
		"Diese Gewürzkiste muss sofort zu Wache Bahir an den Dünen im Osten! Die Karawane wartet nicht, %n. Du hast 60 Sekunden!",
		"Bahir wartet an den Skorpiondünen im Osten. Lauf!", "Die Gewürze sind rechtzeitig angekommen. Perfekt!",
		[]Objective{deliver("bahir", "Gewürzkiste bei Bahir abgegeben")}, give("spice_box"), timed(60), tag("Zeit"), hand("jasim"))
	qd("bahir", "g_raiders", "Wüstenräuber", 18, 1.2,
		"Räuber überfallen unsere Karawanen. Töte acht von ihnen und bring mir vier Karawanenkarten zurück, die sie gestohlen haben, %n.",
		"Die Wüstenräuber lagern südwestlich der Dünen.", "Die Karten! Jetzt können die Karawanen wieder sicher reisen.",
		[]Objective{kill("desert_raider", 8, "Wüstenräuber getötet"), coll("caravan_map", 4, "Karawanenkarten gesammelt")})
	qd("jasim", "g_cactus", "Süße Früchte", 18, 0.9,
		"Meine Kamele lieben Kaktusfrüchte. Sammle sechs davon in der Wüste für mich, %n.",
		"Kakteen wachsen überall in der Wüste, vor allem im Süden.", "Meine Kamele werden sich freuen!",
		[]Objective{coll("cactus_fruit", 6, "Kaktusfrüchte gesammelt")})
	qd("bahir", "g_king", "Der Skorpionkönig", 19, 2.2,
		"Der Skorpionkönig ist riesig und tödlich. Er haust tief in den Dünen. Besiege ihn, %n, aber nicht allein.",
		"Der Skorpionkönig lauert in den Dünen.", "Der König ist tot! Ein großer Sieg für Sandfurt.",
		[]Objective{kill("scorpion_king", 1, "Skorpionkönig besiegt")}, needs("g_raiders"), tag("Elite"), gear())
	qd("neferet", "g_shards", "Scherben der Vergangenheit", 19, 1.2,
		"Das Pharaonengrab birgt Geheimnisse. Sammle vier antike Scherben in der Nähe des Eingangs, %n.",
		"Die Scherben liegen im Sand beim Grab.", "Fantastisch! Diese Scherben erzählen Geschichten.",
		[]Objective{explore("pharaonengrab", "Pharaonengrab erkundet"), coll("ancient_shard", 4, "Antike Scherben gesammelt")})
	qd("neferet", "g_mummies", "Der Fluch der Mumien", 19, 1.3,
		"Die Mumien im Grab sind erwacht. Töte zehn von ihnen, %n, damit ich in Ruhe forschen kann.",
		"Die Mumien wandern rund um das Grab.", "Die Mumien sind zur Ruhe gekommen.",
		[]Objective{kill("mummy", 10, "Mumien getötet")}, needs("g_shards"))
	qd("neferet", "g_guardian", "Der Grabwächter", 20, 2.5,
		"Der Wächter schützt das Grab seit Jahrtausenden. Besiege ihn, %n. Es wird kein einfacher Kampf.",
		"Der Wächter steht vor dem Grab.", "Der Grabwächter ist gefallen. Das Grab liegt offen!",
		[]Objective{kill("tomb_guardian", 1, "Grabwächter besiegt")}, needs("g_mummies"), tag("Elite"), gear())
	qd("karim", "g_final", "Held der Wüste", 20, 3.0,
		"Du hast Sandfurt von Skorpionen, Räubern und Mumien befreit, %n. Dafür schulde ich dir mehr, als ich zahlen kann. Nimm dies als Zeichen unserer Dankbarkeit.",
		"Besiege zuerst den Grabwächter.", "Möge die Sonne immer auf deinem Weg scheinen, Held der Wüste!",
		nil, needs("g_guardian"), gear(), items("hunter_ring"))
}
