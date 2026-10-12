package gamedata

// Second quest wave: roughly as many quests again, spread over all levels (ids start with "n_").

// nq registers a quest of the second wave (rewards a bit smaller than the original ones: there are twice as many now).
func nq(giver, id, title string, lvl int, text, progress, complete string, objs []Objective, opts ...qopt) {
	qd(giver, id, title, lvl, 0.8, text, progress, complete, objs, opts...)
}

// cn is a collect objective whose text is derived from the item name.
func cn(item string, n int) Objective {
	name := item
	if it := Items[item]; it != nil {
		name = it.Name
	}
	return coll(item, n, name+" gesammelt")
}

// kn is a kill objective whose text is derived from the mob name.
func kn(mob string, n int) Objective {
	name := mob
	if m := Mobs[mob]; m != nil {
		name = m.Name
	}
	return kill(mob, n, name+" besiegt")
}

func tk(npc string) Objective {
	name := npc
	if n := NPCByID(npc); n != nil {
		name = n.Name
	}
	return talkTo(npc, "Mit "+name+" gesprochen")
}

func dv(npc string) Objective {
	name := npc
	if n := NPCByID(npc); n != nil {
		name = n.Name
	}
	return deliver(npc, "Lieferung an "+name+" übergeben")
}

func ex(poi string) Objective {
	name := poi
	if p := POIByID(poi); p != nil {
		name = p.Name
	}
	return explore(poi, name+" erkundet")
}

// qitem defines a quest item (cannot be sold).
func qitem(id, name, icon, desc string) {
	Items[id] = &Item{ID: id, Name: name, Icon: icon, Rarity: "common", Desc: desc, Stack: 20, Quest: true}
}

func objs(o ...Objective) []Objective { return o }

func init() {
	// ================= Eichenhain / Erlenwald (1-5) =================
	qitem("n_basket", "Korb mit Brot", "🧺", "Frisch gebacken. Hilda wollte, dass er noch warm ankommt.")
	qitem("n_ore_sample", "Eisenprobe", "🪨", "Ein grauer Brocken, schwerer, als er aussieht.")
	qitem("n_herbs_bundle", "Kräuterbündel", "🌿", "Elaras Kräuter für die Händlerin.")
	qitem("n_signal_horn", "Signalhorn", "📯", "Ein kleines Horn der Wache. Ein kräftiger Ton und alle wissen Bescheid.")

	nq("hilda", "n_bread_basket", "Brot für die Händlerin", 1,
		"Marta hat seit dem Morgen nichts gegessen, die Arme, und der Laden läuft heute wie verrückt. Bring ihr diesen Korb, solange das Brot noch warm ist, %n!",
		"Marta wartet bestimmt schon auf ihr Brot.", "Marta hat sich gefreut? Ich wusste es. Hier, für deine Mühe.",
		objs(dv("marta")), give("n_basket"), items("food_bread"))
	nq("garrick", "n_wolf_bones", "Knochen für die Schmiede", 1,
		"Wolfsknochen, sagst du? Nein, ich esse sie nicht. Aber gemahlen geben sie einen hervorragenden Zusatz für meinen Stahl. Bring mir fünf davon, %n, und ich zahle fair.",
		"Fünf Knochen, nicht weniger. Sonst wird der Stahl brüchig.", "Genau, was ich brauchte! Mein Stahl wird es dir danken.",
		objs(cn("x_wolf_1", 5)))
	nq("brandt", "n_first_drill", "Die erste Übung", 1,
		"Jeder Kämpfer fängt klein an. Meine Schüler trainieren an Strohpuppen, aber wahre Erfahrung bekommst du nur im Feld. Erlege vier Wölfe im Osten, dann sehen wir weiter.",
		"Vier Wölfe, %n. Kein Problem für einen echten Kämpfer.", "Gut gemacht. Du hast den ersten Schritt getan.",
		objs(kn("wolf", 4)))
	nq("marta", "n_wolf_trophies", "Wolfskrallen als Schmuck", 1,
		"Die Kundschaft verlangt nach etwas Wildem. Wolfskrallen, zu Anhängern verarbeitet, verkaufen sich wie warme Semmeln. Ich brauche vier davon.",
		"Vier Krallen, mehr nicht. Frisch von der Jagd wären sie am besten.", "Wunderbar! Daraus mache ich etwas Hübsches.",
		objs(cn("x_wolf_2", 4)))

	nq("elara", "n_herb_delivery", "Kräuter für die Händlerin", 2,
		"Ich habe Heilkräuter gesammelt, aber meine Beine tragen mich nicht mehr bis ins Dorf. Marta braucht sie für ihre Tränke. Würdest du sie hinbringen, %n?",
		"Marta braucht die Kräuter dringend.", "Marta war zufrieden? Gut. Du bist ein zuverlässiger Läufer.",
		objs(dv("marta")), give("n_herbs_bundle"), items("food_bread"))
	nq("hilda", "n_boar_bones", "Keilerknochen für den Hofhund", 2,
		"Unser Hofhund ist alt, aber er liebt Knochen. Die Keiler im Westen haben mehr als genug davon. Bring mir fünf, %n, und ich backe dir einen Kuchen dazu.",
		"Der Hund sabbert schon. Hast du die Knochen?", "Ach, wie er sich freut! Hier ist dein Kuchen.",
		objs(cn("x_boar_1", 5)), items("hilda_pie"))
	nq("torben", "n_boar_charge", "Keilerplage", 2,
		"Die Keiler auf der Westwiese wühlen mir die Felder um. Wenn die Wache nichts unternimmt, tue ich es selbst. Nein, besser du, %n: Erlege acht von ihnen.",
		"Acht Keiler, dann bleibt mein Feld ganz.", "Endlich Ruhe auf dem Feld. Hab Dank!",
		objs(kn("boar", 8)))
	nq("lyra", "n_hunters_practice", "Zielübung für Jäger", 2,
		"Ruhig atmen, Ziel erfassen, loslassen. Doch Theorie nützt nichts ohne Praxis: Erlege drei Wildschweine mit sicherem Schuss, %n, und bring mir ihre Hauer zum Beweis.",
		"Die Hauer, %n. Ohne Beweis kein Lob.", "Saubere Arbeit. Du hast ein Auge dafür.",
		objs(kn("boar", 3), cn("boar_tusk", 3)))

	nq("elara", "n_web_collection", "Netze zählen", 3,
		"Ich dokumentiere die Spinnenpopulation im Norden. Dafür brauche ich Fühler, an ihnen sieht man das Alter. Sammle fünf von diesen Riesenspinnen, aber sei vorsichtig!",
		"Fünf Fühler, bitte. Und halte dich von den Netzen fern.", "Perfekt! Das erleichtert meine Arbeit enorm.",
		objs(cn("x_spider_2", 5)))
	nq("orlan", "n_spider_venom", "Gift für die Magie", 3,
		"Spinnengift ist ein hervorragender Ausgangsstoff für arkane Tinkturen. Bring mir drei Giftbeutel, aber trag sie sorgfältig. Ich brauche sie unversehrt.",
		"Die Giftbeutel dürfen nicht platzen.", "Ausgezeichnet. Du hast ein Händchen für sensibles Material.",
		objs(cn("venom_sac", 3)))
	nq("gerd", "n_watch_the_road", "Die Straße sichern", 3,
		"Die Banditen werden dreister. Sie lauern an der Südstraße auf Händler. Besuche das Banditenlager im Süden und sieh dir an, wie groß ihre Zahl ist, %n.",
		"Wie viele sind es im Lager?", "Verstanden. Danke für die Aufklärung.",
		objs(ex("banditenlager")))
	nq("forester", "n_spider_wood", "Der Spinnenwald", 3,
		"Im Norden ist der Wald in Spinnweben gehüllt. Selbst ich traue mich nicht mehr hinein. Erkunde den Spinnenwald, %n, und erlege dort fünf Spinnen, damit ich sehe, ob es besser wird.",
		"Der Spinnenwald, fünf Spinnen. Du schaffst das.", "Etwas lichter schon. Gut gemacht!",
		objs(ex("spinnenwald"), kn("spider", 5)))

	nq("gerd", "n_bandit_masks", "Maskenball", 4,
		"Ich möchte die Gesichter hinter den Tüchern sehen. Bring mir vier Banditentücher, dann kann ich die Räuber den Fahndungsplakaten zuordnen.",
		"Vier Tücher, dann weiß ich mehr.", "Interessant. Die kenne ich alle noch aus der Stadt.",
		objs(cn("bandit_mask", 4)))
	nq("garrick", "n_ore_run", "Eisenprobe nach Sonnenburg", 4,
		"Meine Eisenprobe muss zu Thoren in Sonnenburg. Ich vertraue sie nur jemandem an, der Beine hat und sie auch benutzt. Das Hauptstadttor liegt weit im Osten, %n.",
		"Thoren wartet auf die Probe.", "Thoren sagt, die Probe ist gut? Natürlich ist sie das. Hier, für den weiten Weg.",
		objs(dv("thoren")), give("n_ore_sample"), gear())
	nq("brandt", "n_bandit_drill", "Kampfübung gegen Banditen", 4,
		"Zeit für etwas Echtes. Die Banditen sind kein Strohpuppen-Gegner: Sie schlagen zurück. Erlege fünf von ihnen und hol dir den Beweis in Form ihrer Beutel.",
		"Fünf Banditen, fünf Beutel.", "Gute Arbeit. Du hast den Mut eines Kriegers.",
		objs(kn("bandit", 5), cn("x_bandit_1", 5)))
	nq("marta", "n_stolen_compass", "Der gestohlene Kompass", 4,
		"Ein Räuber hat meinem Cousin den Kompass abgenommen. Er trug ihn immer bei sich, ein Familienerbstück. Bring mir drei Kompasse der Banditen zurück, einer davon ist es bestimmt.",
		"Einer der Kompasse gehört meinem Cousin.", "Da ist er! Und noch zwei dazu. Danke!",
		objs(cn("x_bandit_2", 3)))

	nq("torben", "n_bear_bones", "Bärenknochen für Torben", 5,
		"Ein Bär, der nichts mehr anrichten kann, ist ein guter Bär. Aber ich brauche mehr: Fünf Knochen für meine Schmiede, die rechtschaffene Wut des Feuers braucht.",
		"Fünf Knochen von den Bären.", "Gut. Mein Feuer wird es zu schätzen wissen.",
		objs(cn("x_bear_1", 5)))
	nq("gerd", "n_horn_of_alarm", "Das Signalhorn", 5,
		"Wir brauchen ein neues Signalhorn für die Wache am Bärenpass. Bring es zu Torben, er ist der Einzige, der das Horn richtig stimmen kann.",
		"Torben wird das Horn stimmen.", "Torben hat das Horn gestimmt? Dann sind wir bereit.",
		objs(dv("torben")), give("n_signal_horn"))
	nq("forester", "n_bear_gorge", "Die Bärenklamm", 5,
		"Die Bären ziehen sich in die Klamm im Osten zurück. Erkunde die Bärenklamm, %n, und finde heraus, wie viele es sind. Wenn du schon dort bist, erlege ruhig drei.",
		"Die Bärenklamm, drei Bären.", "Dann ist der Wald ruhiger. Danke!",
		objs(ex("baerenklamm"), kn("bear", 3)), gear())
	nq("hilda", "n_pie_for_the_hunter", "Kuchen für Lyra", 5,
		"Lyra übt von morgens bis abends und vergisst dabei zu essen. Bring ihr diesen Apfelkuchen, %n, aber lass sie nicht merken, dass er von mir ist. Sie ist stolz.",
		"Lyra wartet nicht, aber der Kuchen schon.", "Lyra hat den Kuchen genommen? Und nicht gemerkt, von wem? Ausgezeichnet.",
		objs(dv("lyra")), give("hilda_pie"), items("food_cheese"))

	// ================= Weizenhof / Sonnenburg (6-9) =================
	qitem("n_flour_sack", "Mehlsack", "🛍️", "Feinstes Weizenmehl vom Weizenhof.")
	qitem("n_blessing", "Segensband", "🎀", "Ein schlichtes Band, von Ilyana gesegnet.")

	nq("odo", "n_hyena_bones", "Hyänenknochen", 6,
		"Meine Schafe sind verschwunden, und überall liegen Hyänenspuren. Erlege die Räuber, %n. Als Beweis genügen mir fünf Knochen.",
		"Fünf Knochen, dann bin ich beruhigt.", "Die Schafe werden dir danken, wenn sie sprechen könnten.",
		objs(cn("x_hyena_1", 5)))
	nq("alena", "n_capital_errands", "Besorgungen für die Bürgermeisterin", 6,
		"Sonnenburg hat viele Bewohner und noch mehr Wünsche. Sprich mit Thoren, dem Waffenmeister, und frag ihn nach seinem Bedarf an Rohstoffen.",
		"Thoren wartet auf dich.", "Rohstoffe also. Ich werde mich darum kümmern.",
		objs(tk("thoren")))
	nq("varus", "n_patrol_duty", "Patrouille", 6,
		"Die Ebenen sind groß, die Wache klein. Erkunde den Weizenhof und die Windmühle für mich, %n, damit ich weiß, ob dort alles sicher ist.",
		"Der Weizenhof und der Windmühlenhügel.", "Danke, das beruhigt mich. Ich schicke bald eine Streife.",
		objs(ex("weizenhof"), ex("windmuehlenhuegel")))

	nq("hannes", "n_flour_delivery", "Mehl für die Bürgermeisterin", 7,
		"Für das Fest braucht Alena frisches Mehl. Bring es in die Stadt, %n, aber stolpere nicht über die Wegelagerer.",
		"Das Mehl muss vor dem Fest ankommen.", "Alena hat das Mehl bekommen? Dann wird das Fest ein Erfolg.",
		objs(dv("alena")), give("n_flour_sack"), items("food_cheese"))
	nq("odo", "n_bandit_purses", "Beutel der Wegelagerer", 7,
		"Die Wegelagerer plündern meine Händler. Hol mir ihre Beutel zurück, %n: fünf Stück sollten genügen, um die Verluste zu begleichen.",
		"Fünf Beutel, dann ist der Schaden ersetzt.", "Daran erkenne ich meine Lieferungen. Hab Dank!",
		objs(cn("x_plains_bandit_1", 5)))
	nq("selene", "n_hyena_hunt", "Die Jagd auf die Rudelführer", 7,
		"Hyänen jagen in Rudeln, und jedes Rudel hat seinen Anführer. Erlege acht von ihnen und bring mir drei Krallen, dann weiß ich, wie stark die Rudel sind.",
		"Acht Hyänen, drei Krallen.", "Du hast die Rudel geschwächt. Gut so.",
		objs(kn("hyena", 8), cn("x_hyena_2", 3)), gear())
	nq("hannes", "n_bull_bones", "Stierknochen für den Mühlstein", 7,
		"Der Mühlstein ist abgenutzt und das Metall nicht mehr hart genug. Mit Knochenasche geht es besser. Bring mir fünf Stierknochen, %n.",
		"Fünf Knochen für den Mühlstein.", "Das wird den Stein wieder schärfen.",
		objs(cn("x_bull_1", 5)))

	nq("ilyana", "n_blessing_band", "Das Segensband", 8,
		"Die Quelle verleiht Segen, wenn man ihr mit Respekt begegnet. Bring dieses Band zu Aldwin, damit er auf der Jagd Glück hat.",
		"Aldwin wartet auf den Segen.", "Aldwin trägt das Band? Möge die Quelle ihm gewogen sein.",
		objs(dv("aldwin")), give("n_blessing"), gear())
	nq("varus", "n_cleanse_the_road", "Die Handelsstraße", 8,
		"Die Straße zur Räuberfestung ist die wichtigste Handelsroute. Erkunde die Räuberfestung, %n, und erlege dabei acht Wegelagerer.",
		"Die Räuberfestung, acht Wegelagerer.", "Das war überfällig. Die Händler werden aufatmen.",
		objs(ex("raeuberfestung"), kn("plains_bandit", 8)))
	nq("odo", "n_bull_run", "Stierhatz", 8,
		"Die Wildstiere zertrampeln meine Weiden. Erlege acht, %n, und die Weide bleibt in diesem Jahr heil.",
		"Acht Stiere, dann hat die Weide Ruhe.", "Endlich Ruhe. Hier, dein Lohn.",
		objs(kn("bull", 8)))

	nq("alena", "n_diplomacy", "Diplomatische Mission", 9,
		"Die Moorbewohner haben um Hilfe gebeten, aber ich traue dem Boten nicht. Sprich mit Rodrik, er kennt die Gerüchte.",
		"Rodrik weiß mehr.", "Verstehe. Die Lage ist ernster, als ich dachte.",
		objs(tk("rodrik")))
	nq("selene", "n_alpha_trail", "Die Spur des Alphas", 9,
		"Rotpelz, der Alpha-Wolf, hat meine Fallen überlistet. Erkunde den Jägerhain und sieh dir seine Spuren an, %n. Danach erlege fünf Wölfe seines Rudels.",
		"Der Jägerhain, fünf Wölfe.", "Sein Rudel ist geschwächt. Jetzt kann ich ihm das Handwerk legen.",
		objs(ex("jaegerhain"), kn("wolf", 5)))
	nq("aldwin", "n_feral_feast", "Fleisch für die Falken", 9,
		"Meine Falken haben Hunger, und die Beute ist rar. Bring mir sechs Stücke Wildschweinfleisch, damit sie bei Kräften bleiben.",
		"Sechs Stücke Fleisch für die Falken.", "Sie stürzen sich schon darauf. Danke, %n.",
		objs(cn("boar_meat", 6)), items("potion_hp2"))
	nq("varus", "n_chief_prep", "Vorbereitung auf den Schlag", 9,
		"Bevor ich den Räuberhauptmann angehe, will ich seine Handlanger geschwächt wissen. Erlege sechs Wegelagerer und bring mir drei ihrer Kompasse, %n. Dann wissen wir, wohin sie ziehen.",
		"Sechs Wegelagerer, drei Kompasse.", "Jetzt wissen wir, wo wir ansetzen müssen. Gut gemacht.",
		objs(kn("plains_bandit", 6), cn("x_plains_bandit_2", 3)), gear())
	nq("ilyana", "n_pilgrim_path", "Der Pilgerweg", 9,
		"Pilger wollen zur Quelle, aber Hyänen lauern am Weg. Geh zum Quellheiligtum, %n, und erlege unterwegs sechs Hyänen. Die Quelle wird es dir danken.",
		"Das Quellheiligtum, sechs Hyänen.", "Die Pilger können wieder ohne Furcht kommen. Möge die Quelle dich segnen.",
		objs(ex("quellheiligtum"), kn("hyena", 6)), items("potion_hp2"))
}
