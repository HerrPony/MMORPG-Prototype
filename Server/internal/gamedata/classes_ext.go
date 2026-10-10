package gamedata

// Skills for levels 10-20, the stance system and the class-quest skills
// (hunter pet, warrior defensive stance, mage water elemental).

func init() {
	add := func(s *Skill) { Skills[s.ID] = s }

	// ---- Krieger ----
	add(&Skill{ID: "w_stance_battle", Name: "Kampfhaltung", Icon: "⚔️", Kind: KindStance, Stance: "battle", Cooldown: 1, Level: 1, Style: "buff", Color: "#ffaa55",
		Desc: "Ausgewogene Haltung: normaler Schaden und normale Bedrohung. Du kämpfst von Beginn an in dieser Haltung."})
	add(&Skill{ID: "w_stance_def", Name: "Verteidigungshaltung", Icon: "🛡️", Kind: KindStance, Stance: "defensive", Cooldown: 1, Level: 10, Quest: true, Style: "buff", Color: "#7aa0ff",
		Desc: "Du erleidest 40% weniger Schaden und erzeugst doppelt so viel Bedrohung, richtest aber 25% weniger Schaden an."})
	add(&Skill{ID: "w_taunt", Name: "Spott", Icon: "😤", Kind: KindTaunt, Cost: 10, Cooldown: 8, Range: 30, Level: 10, Quest: true, Style: "nova", Color: "#ff6644",
		Desc: "Zwingt den Gegner, dich anzugreifen, und setzt seine Bedrohung auf deinen Wert."})
	add(&Skill{ID: "w_charge", Name: "Sturmangriff", Icon: "💨", Kind: KindCharge, Gain: 15, Cooldown: 15, Range: 25, MinRange: 7, Stun: 1.2, Level: 12, Price: 600, Style: "melee", Color: "#ffffff",
		Desc: "Stürme auf das Ziel zu und betäube es 1,2 Sek. lang. Erzeugt Wut."})
	add(&Skill{ID: "w_bash", Name: "Schildschlag", Icon: "🔰", Kind: KindDamage, Cost: 20, Cooldown: 12, Range: 3.5, Power: 14, Stun: 3, Level: 14, Price: 900, Style: "melee", Color: "#9ab8ff",
		Desc: "Ein Schlag mit dem Schild, der Schaden verursacht und das Ziel 3 Sek. betäubt."})
	add(&Skill{ID: "w_execute", Name: "Hinrichten", Icon: "☠️", Kind: KindDamage, Cost: 30, Cooldown: 8, Range: 3.5, Power: 26, Execute: 0.3, ExecMult: 3, Level: 16, Price: 1300, Style: "melee", Color: "#ff3030",
		Desc: "Verursacht dreifachen Schaden an Zielen mit weniger als 30% Leben."})
	add(&Skill{ID: "w_berserk", Name: "Berserkerwut", Icon: "😡", Kind: KindBuff, Cost: 20, Cooldown: 60, Power: 0.4, Taken: 0.2, Duration: 20, Level: 18, Price: 1800, Style: "buff", Color: "#ff2222",
		Desc: "Dein Schaden steigt 20 Sek. lang um 40%, du erleidest aber 20% mehr Schaden."})
	add(&Skill{ID: "w_laststand", Name: "Letztes Gefecht", Icon: "💪", Kind: KindLastStand, Cooldown: 120, Power: 0.3, Duration: 15, Level: 20, Price: 2500, Style: "buff", Color: "#ffd24a",
		Desc: "Erhöht dein maximales Leben 15 Sek. lang um 30% und heilt dich um diesen Betrag."})

	// ---- Jäger ----
	add(&Skill{ID: "h_tame", Name: "Tier zähmen", Icon: "🐾", Kind: KindTame, Cast: 6, Range: 25, Cooldown: 10, Level: 10, Quest: true, Style: "buff", Color: "#88dd66",
		Desc: "Zähme ein wildes Tier (höchstens deine Stufe). Du musst dich dabei 6 Sek. lang konzentrieren."})
	add(&Skill{ID: "h_callpet", Name: "Tier rufen", Icon: "🐺", Kind: KindCallPet, Cast: 3, Cooldown: 5, Level: 10, Quest: true, Style: "buff", Color: "#88dd66",
		Desc: "Ruft dein gezähmtes Tier an deine Seite oder schickt es fort. Ein gefallenes Tier wird wiederbelebt."})
	add(&Skill{ID: "h_mend", Name: "Tier heilen", Icon: "💚", Kind: KindMendPet, Cost: 25, Cooldown: 15, Power: 0.4, Level: 12, Price: 600, Style: "buff", Color: "#66ff88",
		Desc: "Heilt dein Tier sofort um 40% seiner Lebenspunkte."})
	add(&Skill{ID: "h_multi", Name: "Mehrfachschuss", Icon: "🎆", Kind: KindAoETarget, Cost: 35, Cooldown: 10, Range: 30, Radius: 8, Power: 18, Level: 14, Price: 900, Style: "arrow", Color: "#ffcc66",
		Desc: "Trifft das Ziel und alle Feinde in seiner Nähe."})
	add(&Skill{ID: "h_trap", Name: "Frostfalle", Icon: "🧊", Kind: KindAoETarget, Cost: 30, Cooldown: 25, Range: 30, Radius: 6, Power: 8, Slow: 0.7, Duration: 6, Level: 16, Price: 1300, Style: "frost", Color: "#99ddff",
		Desc: "Eine Falle verlangsamt alle Feinde um das Ziel 6 Sek. lang um 70%."})
	add(&Skill{ID: "h_rapid", Name: "Schnellfeuer", Icon: "⚡", Kind: KindBuff, Cost: 40, Cooldown: 60, Power: 0.5, Duration: 12, Level: 18, Price: 1800, Style: "buff", Color: "#ffee44",
		Desc: "Dein Schaden steigt 12 Sek. lang um 50%."})
	add(&Skill{ID: "h_kill", Name: "Todesschuss", Icon: "💀", Kind: KindDamage, Cost: 50, Cooldown: 25, Range: 35, Power: 55, Execute: 0.25, ExecMult: 2.5, Level: 20, Price: 2500, Style: "arrow", Color: "#ff4444",
		Desc: "Ein tödlicher Schuss. Gegen Ziele unter 25% Leben richtet er 2,5-fachen Schaden an."})

	// ---- Magier ----
	add(&Skill{ID: "m_elemental", Name: "Wasserelementar", Icon: "💧", Kind: KindElemental, Cast: 2, Cost: 60, Cooldown: 180, Duration: 45, Level: 10, Quest: true, Style: "frost", Color: "#4aa8ff",
		Desc: "Beschwört 45 Sek. lang einen Wasserelementar, der für dich kämpft und Feinde verlangsamt."})
	add(&Skill{ID: "m_blink", Name: "Blinzeln", Icon: "💫", Kind: KindBlink, Cost: 20, Cooldown: 15, Blink: 14, Level: 12, Price: 600, Style: "nova", Color: "#c8a0ff",
		Desc: "Teleportiert dich 14 Meter nach vorn."})
	add(&Skill{ID: "m_poly", Name: "Verwandlung", Icon: "🐑", Kind: KindPoly, Cast: 1.5, Cost: 30, Cooldown: 20, Range: 28, Duration: 10, Level: 14, Price: 900, Style: "bolt", Color: "#ff99dd",
		Desc: "Verwandelt das Ziel 10 Sek. lang in ein harmloses Schaf. Schaden beendet den Zauber."})
	add(&Skill{ID: "m_blizzard", Name: "Blizzard", Icon: "🌨️", Kind: KindGround, Cast: 1.5, Cost: 45, Cooldown: 25, Range: 30, Radius: 8, Power: 10, Slow: 0.4, Duration: 6, Level: 16, Price: 1300, Style: "frost", Color: "#aee4ff",
		Desc: "Ein Schneesturm tobt 6 Sek. lang am Zielort, verletzt und verlangsamt alle Feinde darin."})
	add(&Skill{ID: "m_frostnova", Name: "Frostnova", Icon: "🥶", Kind: KindAoE, Cost: 30, Cooldown: 25, Radius: 8, Power: 18, Slow: 1.0, Duration: 3, Level: 18, Price: 1800, Style: "nova", Color: "#7ad8ff",
		Desc: "Eine Welle aus Frost schleudert Eis um dich. Getroffene Feinde sind 3 Sek. lang festgefroren."})
	add(&Skill{ID: "m_pyro", Name: "Pyroschlag", Icon: "☄️", Kind: KindDamage, Cast: 3, Cost: 55, Cooldown: 15, Range: 30, Power: 95, Level: 20, Price: 2500, Style: "fire", Color: "#ff5a1a",
		Desc: "Ein gewaltiger Feuerschlag (3 Sek. Zauberzeit)."})

	Classes["warrior"].Skills = []string{"w_strike", "w_stance_battle", "w_heroic", "w_rend", "w_whirl", "w_shout", "w_stance_def", "w_taunt", "w_charge", "w_bash", "w_execute", "w_berserk", "w_laststand"}
	Classes["hunter"].Skills = []string{"h_shot", "h_aimed", "h_serpent", "h_cripple", "h_eagle", "h_tame", "h_callpet", "h_mend", "h_multi", "h_trap", "h_rapid", "h_kill"}
	Classes["mage"].Skills = []string{"m_bolt", "m_fireball", "m_frost", "m_nova", "m_shield", "m_elemental", "m_blink", "m_poly", "m_blizzard", "m_frostnova", "m_pyro"}
}

// StartSkills are the skills a new character knows.
func StartSkills(c *Class) []string {
	var out []string
	for _, id := range c.Skills {
		if s := Skills[id]; s.Level <= 1 && !s.Quest {
			out = append(out, id)
		}
	}
	return out
}
