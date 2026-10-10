package game

import (
	"math"
	"slices"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/loot"
	"mmorpg/internal/progression"
	"mmorpg/internal/protocol"
)

// killCredit returns every player that gets credit for a kill: the ones that hit the mob
// plus their party members who are close by.
func (w *World) killCredit(mob *Entity) []*Entity {
	var out []*Entity
	seen := map[int]bool{}
	add := func(p *Entity) {
		if p != nil && !p.Dead && !seen[p.ID] && p.Inst == mob.Inst {
			seen[p.ID] = true
			out = append(out, p)
		}
	}
	for id := range mob.M.Hit {
		hitter := w.players[id]
		if hitter == nil {
			continue
		}
		add(hitter)
		if party := w.partyOf(hitter); party != nil {
			for _, mid := range party.Members {
				if m := w.players[mid]; m != nil && dist(m, mob) < 100 {
					add(m)
				}
			}
		}
	}
	return out
}

func (w *World) mobDied(t *Entity) {
	m := t.M
	d := m.Def
	m.Target = 0
	m.Threat = map[int]float64{}
	m.CorpseEnd = w.now + 30
	if d.Boss {
		m.CorpseEnd = w.now + 180
	}
	m.RespawnAt = w.now + d.RespawnSec
	m.Loot = map[int]*loot.Bag{}
	w.removeAdds(t)
	credit := w.killCredit(t)
	for _, pl := range credit {
		if pl.P.Target == t.ID {
			pl.P.Auto = false
		}
		// group experience: shared between the party members that get credit
		n := 0
		if party := w.partyOf(pl); party != nil {
			for _, o := range credit {
				if slices.Contains(party.Members, o.ID) {
					n++
				}
			}
		}
		xp := progression.KillXP(d.XPBase, t.Level, pl.Level)
		switch {
		case d.Boss:
			xp = xp * 4
		case d.Elite:
			xp = xp * 5 / 2
		}
		if n > 1 {
			xp = int(float64(xp) * (1 + 0.15*float64(n-1)) / float64(n))
		}
		w.sysChat(pl, "combat", t.Name+" stirbt.")
		w.giveXP(pl, xp)
		w.questEvent(pl, "kill", d.ID)
		bag := loot.Roll(d, t.Level, loot.Looter{Class: pl.P.Char.Class, Quests: pl.P.Char.Quests}, w.rng)
		if !bag.Empty() {
			m.Loot[pl.ID] = bag
		}
	}
	if d.Boss && len(credit) > 0 {
		for _, p := range w.players {
			if p.Inst == t.Inst {
				w.sysChat(p, "yell", d.Name+" wurde besiegt!")
			}
		}
	}
}

// ---- experience ----

func (w *World) giveXP(e *Entity, amount int) {
	p := e.P
	if amount <= 0 || e.Level >= progression.MaxLevel {
		return
	}
	w.sysChat(e, "xp", "Du erhältst "+itoa(amount)+" Erfahrung.")
	lvl, xp, ups := progression.AddXP(e.Level, p.Char.XP, amount)
	p.Char.XP = xp
	if ups == 0 {
		return
	}
	e.Level, p.Char.Level = lvl, lvl
	w.recalc(e)
	e.HP = e.MaxHP
	if p.Class.Resource != "rage" {
		p.Res = float64(p.ResMax)
	}
	w.sendEquipment(e)
	w.send(e, protocol.M{"t": "lvl", "lvl": lvl})
	w.sysChat(e, "level", "Glückwunsch! Du hast Stufe "+itoa(lvl)+" erreicht!")
	if lvl == 10 {
		w.sysChat(e, "level", "Du kannst jetzt deine Klassenquest in Sonnenburg beginnen!")
	}
	w.sendNear(e, protocol.M{"t": "fx", "k": "levelup", "s": e.ID, "d": e.ID})
	w.sendMarks(e) // newly available quests
	w.savePlayer(e)
}

// ---- loot ----

func (w *World) sendLoot(e, t *Entity) {
	b := t.M.Loot[e.ID]
	items := make([]protocol.M, 0, len(b.Items))
	for _, it := range b.Items {
		items = append(items, protocol.M{"i": it.Item, "c": it.Count})
	}
	w.send(e, protocol.M{"t": "loot", "id": t.ID, "name": t.Name, "items": items, "copper": b.Copper})
}

func (w *World) onLoot(e *Entity, m protocol.LootMsg) {
	t := w.ents[m.ID]
	if t == nil || t.M == nil || !t.Dead || e.Dead || t.M.Loot[e.ID].Empty() {
		w.send(e, protocol.M{"t": "lootClose", "id": m.ID})
		return
	}
	if dist(e, t) > InteractDist {
		w.errMsg(e, "Du bist zu weit entfernt.")
		return
	}
	b := t.M.Loot[e.ID]
	takeCopper := func() {
		if b.Copper > 0 {
			e.P.Char.Copper += b.Copper
			w.sysChat(e, "loot", "Du erhältst "+formatCopper(b.Copper)+".")
			b.Copper = 0
		}
	}
	takeItem := func(i int) bool {
		it := b.Items[i]
		left := w.addItem(e, it.Item, it.Count)
		got := it.Count - left
		if got > 0 {
			w.sysChat(e, "loot", "Du erhältst Beute: ["+gamedata.Items[it.Item].Name+"]"+countSuffix(got))
		}
		if left > 0 {
			b.Items[i].Count = left
			w.errMsg(e, "Dein Inventar ist voll.")
			return false
		}
		return true
	}
	switch {
	case m.Idx == -2:
		takeCopper()
	case m.Idx == -1:
		takeCopper()
		kept := b.Items[:0:0]
		for i := range b.Items {
			if !takeItem(i) {
				kept = append(kept, b.Items[i])
			}
		}
		b.Items = kept
	case m.Idx >= 0 && m.Idx < len(b.Items):
		if takeItem(m.Idx) {
			b.Items = append(b.Items[:m.Idx], b.Items[m.Idx+1:]...)
		}
	}
	w.sendInventory(e)
	w.refreshCollect(e)
	if b.Empty() {
		delete(t.M.Loot, e.ID)
		if len(t.M.Loot) == 0 {
			t.M.CorpseEnd = math.Min(t.M.CorpseEnd, w.now+3)
		}
		w.send(e, protocol.M{"t": "lootClose", "id": t.ID})
	} else {
		w.sendLoot(e, t)
	}
}
