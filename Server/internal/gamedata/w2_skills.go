package gamedata

// Skills for levels 22-40 (one new skill every 2 levels for every class).

func init() {
	add := func(s *Skill) {
		s.Price = s.Level * s.Level * 8
		Skills[s.ID] = s
	}

	// ---- Krieger ----
	add(&Skill{ID: "w_cleave", Name: "Klingentanz", Icon: "🌪️", Kind: KindAoE, Cost: 25, Cooldown: 8, Radius: 6, Power: 32, Level: 22, Style: "nova", Color: "#e8e8f0",
		Desc: "Du wirbelst deine Klinge in einem weiten Bogen und triffst alle Feinde in der Nähe."})
	add(&Skill{ID: "w_ironskin", Name: "Eisenhaut", Icon: "🧱", Kind: KindShield, Cost: 20, Cooldown: 40, Power: 120, Duration: 15, Level: 24, Style: "buff", Color: "#aab4c8",
		Desc: "Deine Haut wird hart wie Eisen: Absorbiert 15 Sek. lang Schaden (steigt mit der Stufe)."})
	add(&Skill{ID: "w_crush", Name: "Niederschmettern", Icon: "🔨", Kind: KindDamage, Cost: 30, Cooldown: 14, Range: 3.5, Power: 70, Stun: 2, Level: 26, Style: "melee", Color: "#ffcc88",
		Desc: "Ein vernichtender Hieb, der das Ziel 2 Sek. lang betäubt."})
	add(&Skill{ID: "w_bloodlust", Name: "Blutrausch", Icon: "🩸", Kind: KindBuff, Cost: 25, Cooldown: 60, Power: 0.5, Taken: 0.1, Duration: 15, Level: 28, Style: "buff", Color: "#d02030",
		Desc: "Dein Schaden steigt 15 Sek. lang um 50%, du erleidest aber 10% mehr Schaden."})
	add(&Skill{ID: "w_thunder", Name: "Donnerschlag", Icon: "⚡", Kind: KindAoE, Cost: 35, Cooldown: 14, Radius: 7, Power: 48, Slow: 0.6, Duration: 4, Level: 30, Style: "nova", Color: "#99bbff",
		Desc: "Eine Schockwelle trifft alle Feinde um dich und verlangsamt sie 4 Sek. lang."})
	add(&Skill{ID: "w_wrath", Name: "Zorn der Ahnen", Icon: "💢", Kind: KindDamage, Cost: 40, Cooldown: 10, Range: 3.5, Power: 110, Level: 32, Style: "melee", Color: "#ff9a40",
		Desc: "Ein gewaltiger Schlag, der enormen Schaden verursacht."})
	add(&Skill{ID: "w_annihilate", Name: "Vernichten", Icon: "💀", Kind: KindDamage, Cost: 45, Cooldown: 12, Range: 3.5, Power: 85, Execute: 0.35, ExecMult: 3, Level: 34, Style: "melee", Color: "#ff2020",
		Desc: "Verursacht dreifachen Schaden an Zielen mit weniger als 35% Leben."})
	add(&Skill{ID: "w_titan", Name: "Titanenstärke", Icon: "🗿", Kind: KindBuff, Cost: 30, Cooldown: 90, Power: 0.7, Duration: 20, Level: 36, Style: "buff", Color: "#ffd24a",
		Desc: "Die Kraft eines Titanen: 70% mehr Schaden für 20 Sek."})
	add(&Skill{ID: "w_ancestors", Name: "Ruf der Ahnen", Icon: "👑", Kind: KindLastStand, Cooldown: 150, Power: 0.5, Duration: 20, Level: 38, Style: "buff", Color: "#fff0a0",
		Desc: "Die Ahnen stärken dich: Dein maximales Leben steigt 20 Sek. lang um 50% und du wirst um diesen Betrag geheilt."})
	add(&Skill{ID: "w_worldbreaker", Name: "Weltenbrecher", Icon: "🌋", Kind: KindAoETarget, Cost: 60, Cooldown: 30, Range: 5, Radius: 9, Power: 170, Level: 40, Style: "nova", Color: "#ff5a1a",
		Desc: "Du spaltest den Boden: Das Ziel und alle Feinde in seiner Nähe erleiden gewaltigen Schaden."})

	// ---- Jäger ----
	add(&Skill{ID: "h_explosive", Name: "Explosivschuss", Icon: "💥", Kind: KindAoETarget, Cost: 35, Cooldown: 10, Range: 30, Radius: 6, Power: 40, Level: 22, Style: "fire", Color: "#ff8a30",
		Desc: "Ein explosiver Pfeil trifft das Ziel und alle Feinde um es herum."})
	add(&Skill{ID: "h_viper", Name: "Gift der Viper", Icon: "🐍", Kind: KindDot, Cost: 25, Cooldown: 10, Range: 30, Power: 10, DotPower: 9, Duration: 12, Level: 24, Style: "arrow", Color: "#50e050",
		Desc: "Ein Pfeil mit Viperngift: starker Schaden über 12 Sek."})
	add(&Skill{ID: "h_barrage", Name: "Salve", Icon: "🎇", Kind: KindAoETarget, Cost: 40, Cooldown: 14, Range: 30, Radius: 10, Power: 38, Level: 26, Style: "arrow", Color: "#ffe090",
		Desc: "Eine Salve von Pfeilen regnet auf einen großen Bereich um das Ziel."})
	add(&Skill{ID: "h_falcon", Name: "Falkenblick", Icon: "🦅", Kind: KindBuff, Cost: 30, Cooldown: 60, Power: 0.5, Duration: 20, Level: 28, Style: "buff", Color: "#ffe055",
		Desc: "Deine Sinne schärfen sich: 50% mehr Schaden für 20 Sek."})
	add(&Skill{ID: "h_bind", Name: "Fesselschuss", Icon: "⛓️", Kind: KindDamage, Cost: 30, Cooldown: 15, Range: 30, Power: 45, Slow: 0.8, Duration: 5, Level: 30, Style: "arrow", Color: "#aab0ff",
		Desc: "Der Pfeil fesselt das Ziel: 80% langsamer für 5 Sek."})
	add(&Skill{ID: "h_precise", Name: "Präzisionsschuss", Icon: "🎯", Kind: KindDamage, Cost: 40, Cooldown: 8, Range: 35, Power: 105, Level: 32, Style: "arrow", Color: "#fff088",
		Desc: "Ein perfekt gezielter Schuss mit enormem Schaden."})
	add(&Skill{ID: "h_arrowrain", Name: "Pfeilhagel", Icon: "🌧️", Kind: KindGround, Cast: 1.5, Cost: 50, Cooldown: 30, Range: 30, Radius: 10, Power: 22, Slow: 0.3, Duration: 8, Level: 34, Style: "arrow", Color: "#d8c090",
		Desc: "Ein Hagel aus Pfeilen geht 8 Sek. lang auf den Zielort nieder."})
	add(&Skill{ID: "h_predator", Name: "Raubtierinstinkt", Icon: "🐆", Kind: KindBuff, Cost: 35, Cooldown: 90, Power: 0.6, Duration: 15, Level: 36, Style: "buff", Color: "#ff9a40",
		Desc: "Du jagst wie ein Raubtier: 60% mehr Schaden für 15 Sek."})
	add(&Skill{ID: "h_mercy", Name: "Gnadenschuss", Icon: "☠️", Kind: KindDamage, Cost: 55, Cooldown: 20, Range: 35, Power: 120, Execute: 0.3, ExecMult: 3, Level: 38, Style: "arrow", Color: "#ff3030",
		Desc: "Dreifacher Schaden gegen Ziele mit weniger als 30% Leben."})
	add(&Skill{ID: "h_storm", Name: "Sturm der Pfeile", Icon: "🌪️", Kind: KindAoETarget, Cost: 70, Cooldown: 40, Range: 35, Radius: 10, Power: 165, Level: 40, Style: "arrow", Color: "#ffd060",
		Desc: "Ein Sturm aus Pfeilen vernichtet alles um das Ziel."})

	// ---- Magier ----
	add(&Skill{ID: "m_lance", Name: "Eislanze", Icon: "🧊", Kind: KindDamage, Cast: 1.5, Cost: 40, Cooldown: 3, Range: 30, Power: 78, Slow: 0.4, Duration: 4, Level: 22, Style: "frost", Color: "#8ae0ff",
		Desc: "Eine Lanze aus Eis durchbohrt das Ziel und verlangsamt es."})
	add(&Skill{ID: "m_flamewave", Name: "Flammenwelle", Icon: "🌊", Kind: KindAoE, Cost: 45, Cooldown: 10, Radius: 9, Power: 46, Level: 24, Style: "nova", Color: "#ff7a1a",
		Desc: "Eine Welle aus Feuer rollt um dich herum und versengt alle Feinde."})
	add(&Skill{ID: "m_barrier", Name: "Arkane Barriere", Icon: "🔰", Kind: KindShield, Cost: 50, Cooldown: 40, Power: 140, Duration: 20, Level: 26, Style: "buff", Color: "#99aaff",
		Desc: "Eine mächtige Barriere absorbiert 20 Sek. lang Schaden (steigt mit der Stufe)."})
	add(&Skill{ID: "m_firestorm", Name: "Feuersturm", Icon: "🔥", Kind: KindGround, Cast: 1.5, Cost: 60, Cooldown: 30, Range: 30, Radius: 9, Power: 20, Duration: 8, Level: 28, Style: "fire", Color: "#ff6a20",
		Desc: "Ein Feuersturm tobt 8 Sek. lang am Zielort."})
	add(&Skill{ID: "m_arcblast", Name: "Arkaner Schlag", Icon: "💫", Kind: KindAoETarget, Cast: 1.5, Cost: 55, Cooldown: 12, Range: 30, Radius: 7, Power: 58, Level: 30, Style: "bolt", Color: "#c070ff",
		Desc: "Eine Detonation arkaner Energie trifft das Ziel und seine Umgebung."})
	add(&Skill{ID: "m_meteor", Name: "Meteor", Icon: "☄️", Kind: KindDamage, Cast: 3, Cost: 70, Cooldown: 18, Range: 32, Power: 190, Level: 32, Style: "fire", Color: "#ff4a10",
		Desc: "Ein Meteor stürzt auf das Ziel (3 Sek. Zauberzeit)."})
	add(&Skill{ID: "m_deepfreeze", Name: "Tiefenfrost", Icon: "🥶", Kind: KindAoE, Cost: 55, Cooldown: 30, Radius: 9, Power: 40, Slow: 1.0, Duration: 4, Level: 34, Style: "nova", Color: "#a0e8ff",
		Desc: "Alle Feinde um dich werden 4 Sek. lang festgefroren."})
	add(&Skill{ID: "m_timewarp", Name: "Zeitverzerrung", Icon: "⏳", Kind: KindBuff, Cost: 50, Cooldown: 90, Power: 0.6, Duration: 15, Level: 36, Style: "buff", Color: "#ffe9a0",
		Desc: "Du beschleunigst deine Magie: 60% mehr Schaden für 15 Sek."})
	add(&Skill{ID: "m_arcpower", Name: "Arkane Macht", Icon: "🔮", Kind: KindBuff, Cost: 60, Cooldown: 120, Power: 0.8, Taken: 0.2, Duration: 12, Level: 38, Style: "buff", Color: "#d070ff",
		Desc: "Pure Macht: 80% mehr Schaden für 12 Sek., aber du erleidest 20% mehr Schaden."})
	add(&Skill{ID: "m_armageddon", Name: "Armageddon", Icon: "🌋", Kind: KindGround, Cast: 2.5, Cost: 90, Cooldown: 60, Range: 32, Radius: 12, Power: 44, Duration: 10, Level: 40, Style: "fire", Color: "#ff3a10",
		Desc: "Feuer regnet 10 Sek. lang vom Himmel und verwüstet den Zielort."})

	Classes["warrior"].Skills = append(Classes["warrior"].Skills, "w_cleave", "w_ironskin", "w_crush", "w_bloodlust", "w_thunder", "w_wrath", "w_annihilate", "w_titan", "w_ancestors", "w_worldbreaker")
	Classes["hunter"].Skills = append(Classes["hunter"].Skills, "h_explosive", "h_viper", "h_barrage", "h_falcon", "h_bind", "h_precise", "h_arrowrain", "h_predator", "h_mercy", "h_storm")
	Classes["mage"].Skills = append(Classes["mage"].Skills, "m_lance", "m_flamewave", "m_barrier", "m_firestorm", "m_arcblast", "m_meteor", "m_deepfreeze", "m_timewarp", "m_arcpower", "m_armageddon")
}
