package game

import (
	"math"
	"slices"
	"sort"
	"time"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/loot"
	"mmorpg/internal/progression"
	"mmorpg/internal/protocol"
	"mmorpg/internal/quest"
)

func xpNext(level int) int { return progression.XPToNext(level) }

func finite(vs ...float64) bool {
	for _, v := range vs {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}

// onMove validates a client position update (client-side movement, server-side sanity checks).
func (w *World) onMove(e *Entity, m protocol.Move) {
	p := e.P
	if e.Dead || p.Flight != nil || !finite(m.X, m.Y, m.Z, m.R) {
		return
	}
	if p.Aboard != 0 {
		w.movePassenger(e, m)
		return
	}
	dt := math.Min(w.now-p.LastMove, 2.0)
	p.LastMove = w.now
	allowed := (dt+0.1)*MaxMoveSpeed + 1.5
	if math.Hypot(m.X-e.X, m.Z-e.Z) > allowed {
		w.send(e, protocol.M{"t": "correct", "x": e.X, "y": e.Y, "z": e.Z})
		return
	}
	if e.Inst != 0 {
		if !w.instWalkable(e, m.X, m.Z) {
			w.send(e, protocol.M{"t": "correct", "x": e.X, "y": e.Y, "z": e.Z})
			return
		}
		e.X, e.Z, e.Y = m.X, m.Z, clampf(m.Y, -0.3, 4)
	} else {
		if ship := w.shipAtPos(m.X, m.Z); ship != nil { // stepped onto the deck of the ferry
			w.boardShip(e, ship, m.X, m.Z)
			e.R = m.R
			return
		}
		if !gamedata.Walkable(m.X, m.Z) {
			w.send(e, protocol.M{"t": "correct", "x": e.X, "y": e.Y, "z": e.Z})
			return
		}
		e.X, e.Z = m.X, m.Z
		h := gamedata.Height(e.X, e.Z)
		if gamedata.OnDock(e.X, e.Z) {
			h = deckY
		}
		e.Y = clampf(m.Y, h-0.3, h+4)
	}
	e.R = m.R
	if p.Cast != nil && math.Hypot(e.X-p.Cast.SX, e.Z-p.Cast.SZ) > 0.5 {
		w.cancelCast(e, "Unterbrochen")
	}
}

func (w *World) tickPlayer(e *Entity) {
	p := e.P
	if e.Dead {
		return
	}
	if p.Flight != nil {
		w.tickFlight(e)
		return
	}
	// casting
	if p.Cast != nil && w.now >= p.Cast.End {
		c := p.Cast
		p.Cast = nil
		w.send(e, protocol.M{"t": "castStop"})
		sk := gamedata.Skills[c.Skill]
		t := w.ents[c.Target]
		if sk != nil && w.finishCast(e, sk, t) {
			w.execSkill(e, sk, t)
			w.hunterAuto(e, sk, t)
		}
	}
	// auto attack
	if p.Auto && p.Cast == nil {
		t := w.ents[p.Target]
		basic := w.basicSkill(e)
		if t == nil || t.Kind != KindMob || t.Dead || basic == nil {
			p.Auto = false
		} else if w.now >= p.CD[basic.ID] && dist(e, t) <= basic.Range+1.0 && p.Res >= float64(basic.Cost) && w.hasLOS(e, t, true) {
			w.execSkill(e, basic, t)
		}
	}
	inCombat := w.now-p.LastCombat < 5
	// resources
	switch p.Class.Resource {
	case "rage":
		if !inCombat {
			p.Res = math.Max(0, p.Res-3*w.dt)
		}
	default:
		p.Res = math.Min(float64(p.ResMax), p.Res+p.Class.ResRegen*w.dt)
	}
	// health regeneration
	if !inCombat && e.HP < e.MaxHP {
		p.HPAcc += (0.02*float64(e.MaxHP) + 0.5) * w.dt
		if p.HPAcc >= 1 {
			n := int(p.HPAcc)
			p.HPAcc -= float64(n)
			e.HP = min(e.MaxHP, e.HP+n)
		}
	}
}

func (w *World) basicSkill(e *Entity) *gamedata.Skill {
	for _, id := range e.P.Class.Skills {
		if s := gamedata.Skills[id]; s.Basic {
			return s
		}
	}
	return nil
}

// ---- inventory ----

func (w *World) sendInventory(e *Entity) {
	w.send(e, protocol.M{"t": "inv", "items": e.P.Char.Inventory, "copper": e.P.Char.Copper, "max": MaxInventory})
}

func (w *World) addItem(e *Entity, item string, n int) int {
	inv, left := loot.AddToInventory(e.P.Char.Inventory, item, n, MaxInventory)
	e.P.Char.Inventory = inv
	return left
}

// ---- revive ----

func (w *World) onRevive(e *Entity) {
	if !e.Dead {
		return
	}
	e.Dead = false
	w.recalc(e)
	e.HP = e.MaxHP
	if e.Inst != 0 {
		if in := w.insts[e.Inst]; in != nil {
			e.X, e.Z = in.Def.Origin[0]+in.Def.Start[0], in.Def.Origin[1]+in.Def.Start[1]
		}
	} else { // wake up at the nearest settlement
		var best *gamedata.Town
		bd := 1e9
		for _, t := range gamedata.Towns {
			if t.Kind == "site" {
				continue
			}
			if d := math.Hypot(t.X-e.X, t.Z-e.Z); d < bd {
				best, bd = t, d
			}
		}
		e.X, e.Z, e.R = best.X, best.Z+4, 0
		if best.ID == "eichenhain" {
			e.X, e.Z, e.R = gamedata.SpawnPoint.X, gamedata.SpawnPoint.Z, gamedata.SpawnPoint.R
		}
	}
	e.Y = gamedata.Height(e.X, e.Z)
	e.P.Res = float64(e.P.Class.ResStart)
	if e.P.Class.Resource != "rage" {
		e.P.Res = float64(e.P.ResMax)
	}
	e.P.ZoneKey, e.P.PoiID = "", ""
	e.P.LastMove = w.now
	w.send(e, protocol.M{"t": "correct", "x": e.X, "y": e.Y, "z": e.Z})
}

// ---- quests ----

func (w *World) sendQuests(e *Entity) {
	now := time.Now().Unix()
	list := make([]protocol.M, 0, len(e.P.Char.Quests))
	for id, st := range e.P.Char.Quests {
		if st.Status == quest.Done {
			continue
		}
		m := protocol.M{"id": id, "status": st.Status, "counts": st.Counts}
		if st.Expires > 0 {
			m["left"] = max(0, st.Expires-now)
		}
		list = append(list, m)
	}
	w.send(e, protocol.M{"t": "quests", "list": list})
}

// sendMarks tells the client which quest marker to show above each NPC.
func (w *World) sendMarks(e *Entity) {
	marks := protocol.M{}
	for _, n := range w.ents {
		if n.Kind != KindNPC {
			continue
		}
		_, m := quest.NPCState(n.N, e.P.Char.Quests, e.Level, e.P.Char.Class)
		marks[itoa(n.ID)] = m
	}
	w.send(e, protocol.M{"t": "marks", "marks": marks})
}

// questEvent credits an event (kill, deliver, talk, explore, tame) to the player's quests.
func (w *World) questEvent(e *Entity, kind, target string) {
	ch := quest.Event(e.P.Char.Quests, kind, target)
	if len(ch) == 0 {
		return
	}
	for _, qid := range ch {
		q := gamedata.Quests[qid]
		st := e.P.Char.Quests[qid]
		if kind == "deliver" && q.GiveItem != "" { // the item is handed over
			e.P.Char.Inventory = loot.RemoveFromInventory(e.P.Char.Inventory, q.GiveItem, 1)
			w.sendInventory(e)
		}
		for i, o := range q.Objectives {
			if o.Kind == kind && o.Target == target {
				w.sysChat(e, "quest", q.Title+": "+o.Text+": "+itoa(st.Counts[i])+"/"+itoa(o.Count))
			}
		}
		if st.Status == quest.Complete {
			w.sysChat(e, "quest", "Quest abgeschlossen – kehre zu "+npcName(quest.TurnInNPC(q))+" zurück!")
		}
	}
	w.send(e, protocol.M{"t": "questUpdate"})
	w.sendQuests(e)
	w.sendMarks(e)
}

func npcName(id string) string {
	if n := gamedata.NPCByID(id); n != nil {
		return n.Name
	}
	return "dem Questgeber"
}

// refreshCollect updates collect objectives after the inventory changed.
func (w *World) refreshCollect(e *Entity) {
	ch := quest.RefreshCollect(e.P.Char.Quests, e.P.Char.Inventory)
	if len(ch) == 0 {
		return
	}
	for _, qid := range ch {
		q := gamedata.Quests[qid]
		st := e.P.Char.Quests[qid]
		for i, o := range q.Objectives {
			if o.Kind == "collect" {
				w.sysChat(e, "quest", q.Title+": "+o.Text+": "+itoa(st.Counts[i])+"/"+itoa(o.Count))
			}
		}
		if st.Status == quest.Complete {
			w.sysChat(e, "quest", "Quest abgeschlossen – kehre zu "+npcName(quest.TurnInNPC(q))+" zurück!")
		}
	}
	w.sendQuests(e)
	w.sendMarks(e)
}

// tickQuestTimers fails timed quests that ran out.
func (w *World) tickQuestTimers() {
	now := time.Now().Unix()
	for _, e := range w.players {
		failed := false
		for id, st := range e.P.Char.Quests {
			if st.Status != quest.Active || st.Expires == 0 || now < st.Expires {
				continue
			}
			q := gamedata.Quests[id]
			delete(e.P.Char.Quests, id)
			if q.GiveItem != "" {
				e.P.Char.Inventory = loot.RemoveFromInventory(e.P.Char.Inventory, q.GiveItem, 99)
			}
			w.sysChat(e, "yell", "Die Zeit ist abgelaufen! Quest gescheitert: "+q.Title)
			w.send(e, protocol.M{"t": "err", "text": "Quest gescheitert: " + q.Title})
			failed = true
		}
		if failed {
			w.sendInventory(e)
			w.sendQuests(e)
			w.sendMarks(e)
		}
	}
}

func (w *World) onInteract(e *Entity, id int) {
	t := w.ents[id]
	if t == nil || e.Dead || t.Inst != e.Inst || e.P.Flight != nil {
		return
	}
	if t.Kind == KindShip {
		return // the ferry is boarded by walking onto it
	}
	if dist(e, t) > InteractDist {
		w.errMsg(e, "Du bist zu weit entfernt.")
		return
	}
	switch {
	case t.Kind == KindNPC:
		w.questEvent(e, "deliver", t.N.ID)
		w.questEvent(e, "talk", t.N.ID)
		if t.N.Flight {
			w.unlockFlightPoint(e, t)
		}
		w.sendDialog(e, t)
	case t.Kind == KindMob && t.Dead && !t.M.Loot[e.ID].Empty():
		w.sendLoot(e, t)
	case t.Kind == KindObj:
		w.useObject(e, t)
	}
}

func (w *World) useObject(e, o *Entity) {
	def := o.O.Def
	if o.O.Gone {
		return
	}
	switch {
	case def.Portal == "exit":
		w.exitDungeon(e)
		return
	case len(def.Portal) > 8 && def.Portal[:8] == "dungeon:":
		w.enterDungeon(e, def.Portal[8:])
		return
	}
	if def.Quest == "" || !quest.WantsItem(e.P.Char.Quests, def.Quest) {
		w.errMsg(e, "Dafür hast du keine Verwendung.")
		return
	}
	if left := w.addItem(e, def.Item, 1); left > 0 {
		w.errMsg(e, "Dein Inventar ist voll.")
		return
	}
	w.sysChat(e, "loot", "Du erhältst Questgegenstand: ["+gamedata.Items[def.Item].Name+"]")
	o.O.Gone = true
	o.O.RespawnAt = w.now + def.Respawn
	w.sendInventory(e)
	w.refreshCollect(e)
}

func (w *World) tickObject(e *Entity) {
	if e.O.Gone && e.O.Def.Respawn > 0 && w.now >= e.O.RespawnAt {
		e.O.Gone = false
	}
}

func (w *World) npcWith(e *Entity, pred func(n *Entity) bool) *Entity {
	for _, n := range w.ents {
		if n.Kind == KindNPC && n.Inst == e.Inst && dist(e, n) <= InteractDist && pred(n) {
			return n
		}
	}
	return nil
}

func (w *World) onAccept(e *Entity, qid string) {
	q := gamedata.Quests[qid]
	if q == nil {
		return
	}
	if w.npcWith(e, func(n *Entity) bool { return slices.Contains(n.N.Quests, qid) }) == nil {
		w.errMsg(e, "Du bist zu weit entfernt.")
		return
	}
	if !quest.Available(e.P.Char.Quests, e.Level, e.P.Char.Class, qid) {
		return
	}
	if q.GiveItem != "" && len(e.P.Char.Inventory) >= MaxInventory && loot.CountItem(e.P.Char.Inventory, q.GiveItem) == 0 {
		w.errMsg(e, "Dein Inventar ist voll.")
		return
	}
	quest.Accept(e.P.Char.Quests, e.Level, e.P.Char.Class, qid)
	w.sysChat(e, "quest", "Quest angenommen: "+q.Title)
	if q.GiveItem != "" {
		w.addItem(e, q.GiveItem, 1)
		w.sendInventory(e)
	}
	for _, sk := range q.GrantSkill {
		w.learnSkill(e, sk)
	}
	if q.TimeLimit > 0 {
		w.sysChat(e, "yell", "Du hast "+itoa(q.TimeLimit)+" Sekunden Zeit!")
	}
	quest.RefreshCollect(e.P.Char.Quests, e.P.Char.Inventory)
	w.sendQuests(e)
	w.sendMarks(e)
	w.savePlayer(e)
}

func (w *World) learnSkill(e *Entity, id string) {
	sk := gamedata.Skills[id]
	if sk == nil || slices.Contains(e.P.Char.Skills, id) {
		return
	}
	e.P.Char.Skills = append(e.P.Char.Skills, id)
	w.sysChat(e, "level", "Du hast eine neue Fertigkeit erlernt: "+sk.Name)
	w.send(e, protocol.M{"t": "skills", "skills": e.P.Char.Skills, "learned": id})
}

func (w *World) onTurnIn(e *Entity, qid string) {
	st := e.P.Char.Quests[qid]
	q := gamedata.Quests[qid]
	if st == nil || q == nil || st.Status != quest.Complete {
		return
	}
	turn := quest.TurnInNPC(q)
	if w.npcWith(e, func(n *Entity) bool { return n.N.ID == turn }) == nil {
		w.errMsg(e, "Du bist zu weit entfernt.")
		return
	}
	inv := e.P.Char.Inventory
	for _, o := range quest.CollectItems(q) {
		inv = loot.RemoveFromInventory(inv, o.Target, o.Count)
	}
	need := len(q.Reward.Items)
	if q.Reward.Gear {
		need++
	}
	if MaxInventory-len(inv) < need {
		w.errMsg(e, "Dein Inventar ist voll.")
		return
	}
	e.P.Char.Inventory = inv
	st.Status = quest.Done
	e.P.Char.Copper += q.Reward.Copper
	rewards := append([]string(nil), q.Reward.Items...)
	if q.Reward.Gear {
		rewards = append(rewards, gamedata.RandomGear(q.Level, e.P.Char.Class, w.rng))
	}
	for _, it := range rewards {
		w.addItem(e, it, 1)
		w.sysChat(e, "loot", "Du erhältst Gegenstand: ["+gamedata.Items[it].Name+"]")
	}
	for _, sk := range q.Reward.Skills {
		w.learnSkill(e, sk)
	}
	w.sysChat(e, "quest", "Quest abgeschlossen: "+q.Title)
	w.send(e, protocol.M{"t": "questDone", "quest": qid})
	quest.RefreshCollect(e.P.Char.Quests, e.P.Char.Inventory)
	w.sendInventory(e)
	w.sendQuests(e)
	w.giveXP(e, q.Reward.XP)
	w.sendMarks(e)
	w.savePlayer(e)
}

func (w *World) onAbandon(e *Entity, qid string) {
	if st := e.P.Char.Quests[qid]; st != nil && st.Status != quest.Done {
		if q := gamedata.Quests[qid]; q != nil && q.GiveItem != "" {
			e.P.Char.Inventory = loot.RemoveFromInventory(e.P.Char.Inventory, q.GiveItem, 99)
			w.sendInventory(e)
		}
		delete(e.P.Char.Quests, qid)
		w.sendQuests(e)
		w.sendMarks(e)
	}
}

// ---- NPC dialog, trainer, vendor, griffin master ----

func (w *World) sendDialog(e, n *Entity) {
	m := protocol.M{"t": "dialog", "npc": n.ID, "name": n.Name, "title": n.Title, "greet": n.N.Greet}
	offers := quest.Offers(n.N, e.P.Char.Quests, e.Level, e.P.Char.Class)
	rank := map[string]int{quest.MarkReady: 0, quest.MarkAvailable: 1, quest.MarkProgress: 2}
	sort.SliceStable(offers, func(i, j int) bool { return rank[offers[i].State] < rank[offers[j].State] })
	if len(offers) > 0 {
		list := make([]protocol.M, 0, len(offers))
		for _, o := range offers {
			list = append(list, protocol.M{"id": o.ID, "state": o.State})
		}
		m["quests"] = list
	}
	if n.N.Vendor {
		m["vendor"] = true
		m["stock"] = n.N.Stock
	}
	if n.N.Flight {
		m["flight"] = w.flightInfo(e, n)
	}
	if n.N.Ship {
		m["ferry"] = w.ferryInfo(e)
	}
	if n.N.Trainer != "" {
		cls := gamedata.Classes[n.N.Trainer]
		tr := protocol.M{"class": n.N.Trainer, "className": cls.Name, "mine": n.N.Trainer == e.P.Char.Class}
		if n.N.Trainer == e.P.Char.Class {
			var list []protocol.M
			for _, id := range cls.Skills {
				sk := gamedata.Skills[id]
				if sk.Level <= 1 || sk.Quest {
					continue
				}
				state := "locked"
				switch {
				case slices.Contains(e.P.Char.Skills, id):
					state = "known"
				case e.Level >= sk.Level:
					state = "learnable"
				}
				list = append(list, protocol.M{"id": id, "level": sk.Level, "price": sk.Price, "state": state})
			}
			tr["skills"] = list
		}
		m["trainer"] = tr
	}
	w.send(e, m)
}

func (w *World) npcNear(e *Entity, id int) *Entity {
	n := w.ents[id]
	if n == nil || n.Kind != KindNPC || e.Dead || n.Inst != e.Inst {
		return nil
	}
	if dist(e, n) > InteractDist {
		w.errMsg(e, "Du bist zu weit entfernt.")
		return nil
	}
	return n
}

func (w *World) onTrain(e *Entity, m protocol.TrainMsg) {
	n := w.npcNear(e, m.NPC)
	sk := gamedata.Skills[m.Skill]
	if n == nil || sk == nil || n.N.Trainer != e.P.Char.Class || !slices.Contains(e.P.Class.Skills, sk.ID) || sk.Quest {
		return
	}
	switch {
	case slices.Contains(e.P.Char.Skills, sk.ID):
		return
	case e.Level < sk.Level:
		w.errMsg(e, "Dafür bist du noch nicht erfahren genug.")
		return
	case e.P.Char.Copper < sk.Price:
		w.errMsg(e, "Du hast nicht genug Geld.")
		return
	}
	e.P.Char.Copper -= sk.Price
	w.learnSkill(e, sk.ID)
	w.sendInventory(e)
	w.sendDialog(e, n)
	w.savePlayer(e)
}

func (w *World) onSell(e *Entity, m protocol.SellMsg) {
	n := w.npcNear(e, m.NPC)
	if n == nil || !n.N.Vendor {
		return
	}
	inv := e.P.Char.Inventory
	if m.Idx < 0 || m.Idx >= len(inv) {
		return
	}
	s := inv[m.Idx]
	def := gamedata.Items[s.Item]
	if def == nil || def.Quest {
		w.errMsg(e, "Das kann man nicht verkaufen.")
		return
	}
	gain := def.Sell * s.Count
	e.P.Char.Inventory = append(inv[:m.Idx:m.Idx], inv[m.Idx+1:]...)
	e.P.Char.Copper += gain
	w.sysChat(e, "loot", "Du verkaufst ["+def.Name+"]"+countSuffix(s.Count)+" für "+formatCopper(gain)+".")
	w.sendInventory(e)
	w.refreshCollect(e)
}
