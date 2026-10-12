package game

import (
	"math"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
	"mmorpg/internal/storage"
)

func petHP(level int, kind string) int {
	hp := 40 + 22*level + level*level
	if kind == "elemental" {
		return int(float64(hp) * 0.8)
	}
	return int(float64(hp) * 0.9)
}

func (w *World) spawnPet(owner *Entity, kind, mobID, name, model string, expire float64) *Entity {
	a := w.rng.Float64() * 2 * math.Pi
	e := &Entity{ID: w.newID(), Kind: KindPet, Name: name, Model: model, Level: owner.Level, Inst: owner.Inst,
		X: owner.X + math.Cos(a)*2, Z: owner.Z + math.Sin(a)*2, R: owner.R,
		Pet: &Pet{Owner: owner.ID, Kind: kind, MobID: mobID, ExpireAt: expire}}
	e.Y = gamedata.Height(e.X, e.Z)
	e.MaxHP = petHP(e.Level, kind)
	e.HP = e.MaxHP
	w.ents[e.ID] = e
	return e
}

// ---- hunter beast ----

func (w *World) tame(owner, t *Entity) {
	def := t.M.Def
	owner.P.Char.Pet = &storage.PetData{Mob: def.ID, Name: def.Name}
	w.sysChat(owner, "quest", "Du hast "+def.Name+" gezähmt!")
	w.sendNear(owner, protocol.M{"t": "fx", "k": "levelup", "s": t.ID, "d": t.ID})
	// the wild animal disappears until it respawns
	t.Dead = true
	t.M.Gone = true
	t.M.RespawnAt = w.now + def.RespawnSec
	t.M.Threat = map[int]float64{}
	w.questEvent(owner, "tame", def.ID)
	w.summonBeast(owner)
}

func (w *World) summonBeast(owner *Entity) {
	pd := owner.P.Char.Pet
	if pd == nil || owner.P.PetEnt != 0 {
		return
	}
	def := gamedata.Mobs[pd.Mob]
	pd.Dead = false
	e := w.spawnPet(owner, "beast", pd.Mob, pd.Name, def.Model, 0)
	owner.P.PetEnt = e.ID
	w.sendNear(owner, protocol.M{"t": "fx", "k": "pet_summon", "s": e.ID, "d": e.ID})
}

func (w *World) callPet(owner *Entity) {
	if pe := w.ents[owner.P.PetEnt]; pe != nil { // dismiss
		delete(w.ents, pe.ID)
		owner.P.PetEnt = 0
		w.sysChat(owner, "combat", "Du schickst dein Tier fort.")
		return
	}
	if owner.P.Char.Pet != nil {
		w.summonBeast(owner)
		w.sysChat(owner, "combat", "Du rufst "+owner.P.Char.Pet.Name+".")
	}
}

// ---- mage elemental ----

func (w *World) summonElemental(owner *Entity, sk *gamedata.Skill) {
	if old := w.ents[owner.P.ElemEnt]; old != nil {
		delete(w.ents, old.ID)
	}
	e := w.spawnPet(owner, "elemental", "", "Wasserelementar", "water_elem_pet", w.now+sk.Duration)
	owner.P.ElemEnt = e.ID
	w.sendNear(owner, protocol.M{"t": "fx", "k": "pet_summon", "s": e.ID, "d": e.ID})
	w.sysChat(owner, "combat", "Ein Wasserelementar erscheint an deiner Seite.")
}

func (w *World) dismissPets(owner *Entity) {
	for _, id := range []int{owner.P.PetEnt, owner.P.ElemEnt} {
		delete(w.ents, id)
	}
	owner.P.PetEnt, owner.P.ElemEnt = 0, 0
}

func (w *World) petDied(t *Entity) {
	owner := w.players[t.Pet.Owner]
	if owner != nil {
		if t.Pet.Kind == "beast" {
			owner.P.PetEnt = 0
			if owner.P.Char.Pet != nil {
				owner.P.Char.Pet.Dead = true
			}
			w.sysChat(owner, "combat", t.Name+" ist gefallen. Rufe dein Tier erneut, um es wiederzubeleben.")
		} else {
			owner.P.ElemEnt = 0
		}
	}
	delete(w.ents, t.ID)
}

func (w *World) tickPet(e *Entity) {
	pt := e.Pet
	owner := w.players[pt.Owner]
	if owner == nil || owner.Dead && pt.Kind == "elemental" {
		delete(w.ents, e.ID)
		return
	}
	if pt.ExpireAt > 0 && w.now >= pt.ExpireAt {
		w.sendNear(e, protocol.M{"t": "fx", "k": "pet_summon", "s": e.ID, "d": e.ID})
		owner.P.ElemEnt = 0
		delete(w.ents, e.ID)
		return
	}
	if e.Inst != owner.Inst || dist(e, owner) > 70 { // keep up with the owner
		e.Inst = owner.Inst
		e.X, e.Z = owner.X+2, owner.Z+2
		e.Y = gamedata.Height(e.X, e.Z)
	}
	if e.stunned(w.now) {
		return
	}
	// choose a target: the owner's target, or whatever fights the owner
	var tgt *Entity
	if t := w.ents[owner.P.Target]; t != nil && t.M != nil && !t.Dead && t.Inst == e.Inst && dist(t, owner) < 45 &&
		(owner.P.Auto || w.now-owner.P.LastCombat < 6) {
		tgt = t
	}
	if tgt == nil {
		if cur := w.ents[pt.Target]; cur != nil && cur.M != nil && !cur.Dead && cur.M.Threat[owner.ID] > 0 && dist(cur, e) < 40 {
			tgt = cur
		} else {
			for _, o := range w.ents {
				if o.M != nil && !o.Dead && o.Inst == e.Inst && o.M.Threat[owner.ID] > 0 && dist(o, e) < 35 {
					tgt = o
					break
				}
			}
		}
	}
	if tgt != nil {
		pt.Target = tgt.ID
		dd := dist(e, tgt)
		if pt.Kind == "elemental" {
			if dd > 17 {
				w.walkTo(e, tgt.X, tgt.Z, 7.5)
			} else {
				e.R = math.Atan2(tgt.X-e.X, tgt.Z-e.Z)
			}
			if dd <= 20 && w.now >= pt.NextAtk {
				pt.NextAtk = w.now + 2.4
				w.sendNear(e, protocol.M{"t": "fx", "k": "mobshot", "s": e.ID, "d": tgt.ID, "style": "frost"})
				w.damage(e, tgt, (8+float64(e.Level)*3.4)*(0.9+0.2*w.rng.Float64()), false)
				tgt.SlowMult, tgt.SlowUntil = 0.3, w.now+3
			}
		} else {
			if dd > 2.4 {
				w.walkTo(e, tgt.X, tgt.Z, 8)
			} else {
				e.R = math.Atan2(tgt.X-e.X, tgt.Z-e.Z)
			}
			if dd <= 3.2 && w.now >= pt.NextAtk {
				pt.NextAtk = w.now + 1.8
				w.sendNear(e, protocol.M{"t": "fx", "k": "bite", "s": e.ID, "d": tgt.ID})
				w.damage(e, tgt, (4+float64(e.Level)*2.3)*(0.9+0.2*w.rng.Float64()), false)
			}
		}
		return
	}
	pt.Target = 0
	// follow
	d := dist(e, owner)
	switch {
	case d > 14:
		w.walkTo(e, owner.X, owner.Z, 11)
	case d > 4.5:
		w.walkTo(e, owner.X, owner.Z, 7.5)
	}
	// slow regeneration out of combat
	if w.tickN%20 == 0 && e.HP < e.MaxHP {
		e.HP = min(e.MaxHP, e.HP+max(1, e.MaxHP/50))
	}
}

// ---- ground effects (Blizzard) ----

type GroundFX struct {
	Src       int
	Inst      int
	X, Z, R   float64
	End, Next float64
	Power     float64
	Slow      float64
}

func (w *World) addGround(src *Entity, x, z float64, sk *gamedata.Skill) {
	w.grounds = append(w.grounds, &GroundFX{Src: src.ID, Inst: src.Inst, X: x, Z: z, R: sk.Radius, End: w.now + sk.Duration,
		Next: w.now + 0.5, Power: sk.Power, Slow: sk.Slow})
	w.sendNear(src, protocol.M{"t": "fx", "k": "ground", "s": src.ID, "d": src.ID, "x": x, "z": z, "r": sk.Radius, "dur": sk.Duration, "color": sk.Color})
}

func (w *World) tickGrounds() {
	if len(w.grounds) == 0 {
		return
	}
	keep := w.grounds[:0]
	for _, g := range w.grounds {
		src := w.ents[g.Src]
		if w.now >= g.End || src == nil {
			continue
		}
		if w.now >= g.Next {
			g.Next += 1
			for _, o := range w.enemiesNear(g.Inst, g.X, g.Z, g.R) {
				w.damage(src, o, w.calcDamage(src, g.Power), false)
				if g.Slow > 0 && !o.Dead {
					o.SlowMult, o.SlowUntil = g.Slow, w.now+1.5
				}
			}
		}
		keep = append(keep, g)
	}
	w.grounds = keep
}
