package gamedata

// More vendors, so that every settlement sells gear of its own level range and the best gear of the next range
// is always a short trip away. Stock windows are power tiers (see gear.go): tier t needs character level 2t+1.

func init() {
	v := func(id, name, title, model string, x, z, r float64, kind string, lo, hi int, greet string) {
		NPCs = append(NPCs, &NPC{ID: id, Name: name, Title: title, Model: model, X: x, Z: z, R: r, Vendor: true, StockK: kind, StockLo: lo, StockHi: hi, Greet: greet})
	}
	// Sonnenburg: the capital also sells the jewellery up to tier 6
	if n := NPCByID("pip"); n != nil {
		n.StockHi = 6
	}
	v("wiebke", "Hafenhändlerin Wiebke", "Händlerin", "merchant", 224, 46, 1.2, "weapons,armor,jewelry", 4, 6, "Was vom Schiff fällt, landet bei mir. Alles legal, versprochen.")
	v("torg", "Grubenhändler Torg", "Händler", "merchant", -166, -201, 0.4, "weapons,armor", 7, 8, "Wer in die Mine geht, sollte vernünftig ausgerüstet sein.")
	v("sven", "Pelzhändler Sven", "Schmuck & Pelze", "merchant", -104, -217, 3.9, "jewelry", 7, 9, "Warm gehalten wird nur, wer gut gekleidet ist.")
	v("pott", "Trödler Pott", "Schmuck & Trödel", "merchant", -178, 170, 2.0, "jewelry", 5, 7, "Das Moor gibt vieles her. Das meiste davon glitzert.")
	v("tarek", "Schmied Tarek", "Waffen & Rüstungen", "smith", 194, 178, 4.4, "weapons,armor", 9, 10, "Wüstenstahl wird heiß geschmiedet und kalt geführt.")
	v("xia", "Händlerin Xia", "Schmuck", "merchant", 616, 4, 3.3, "jewelry", 10, 12, "Jade, Perlen und Glück – alles zu haben.")
	v("daisuke", "Schmied Daisuke", "Waffen & Rüstungen", "smith", 692, 128, 4.7, "weapons,armor", 12, 13, "Bambus biegt sich, mein Stahl nicht.")
	v("kaito", "Juwelier Kaito", "Juwelier", "merchant", 856, -10, 4.7, "jewelry", 14, 16, "Selbst die Kaiserin trägt meine Ringe.")
	v("ume", "Händlerin Ume", "Rüstungen", "merchant", 802, 154, 4.7, "armor", 14, 15, "Seide vom Feinsten, Rüstung vom Besten.")
	v("zafira", "Juwelierin Zafira", "Schmuck", "merchant", 996, 136, 4.4, "jewelry", 15, 17, "Rote Steine, rote Preise.")
	v("ruk", "Waffenhändler Ruk", "Waffen", "smith", 984, 124, 3.5, "weapons", 16, 17, "Klingen für die Schlucht.")
	v("cinder", "Händlerin Cinder", "Schmuck & Tränke", "merchant", 932, -98, 1.2, "jewelry", 17, 19, "In der Glut werden Edelsteine geboren.")
	v("obran", "Rüstmeister Obran", "Waffen & Rüstungen", "smith", 956, -106, 3.5, "weapons,armor", 18, 19, "Wer die Aschenlande betritt, sollte das Beste tragen.")

	// food and potions of the right strength at every vendor (the stock lists are filled in by Init)
	for _, n := range NPCs {
		if !n.Vendor || n.StockK == "" {
			continue
		}
		for _, lv := range []int{2*n.StockLo + 1, 2*n.StockHi + 1} {
			n.Stock = appendUnique(n.Stock, foods[min(len(foods)-1, lv/6)].id)
		}
		if n.StockHi >= 2 {
			n.Stock = appendUnique(n.Stock, potionID(2*n.StockHi+1))
		}
	}
}

func appendUnique(l []string, id string) []string {
	for _, s := range l {
		if s == id {
			return l
		}
	}
	return append(l, id)
}

func potionID(level int) string {
	switch {
	case level >= 33:
		return "potion_hp6"
	case level >= 25:
		return "potion_hp5"
	case level >= 17:
		return "potion_hp4"
	case level >= 11:
		return "potion_hp3"
	case level >= 5:
		return "potion_hp2"
	}
	return "potion_hp"
}
