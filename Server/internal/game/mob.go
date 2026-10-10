package game

import (
	"math"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
)

func (w *World) respawnMob(e *Entity) {
	m := e.M
	d := m.Def
	e.Level = d.LevelMin
	if d.LevelMax > d.LevelMin {
		e.Level += w.rng.Intn(d.LevelMax - d.LevelMin + 1)
	}
	e.MaxHP = d.HPBase + d.HPPerLevel*(e.Level-1)
	e.HP = e.MaxHP
	e.Dead = false
	e.Dots = nil
	e.SlowUntil, e.StunUntil, e.PolyUntil = 0, 0, 0
	m.Gone = false
	m.Target = 0
	m.Returning = false
	m.Loot = nil
	m.Hit = map[int]bool{}
	m.Threat = map[int]float64{}
	m.Wandering = false
	m.Enraged = false
	m.AbNext = make([]float64, len(d.Abilities))
	m.AbDone = make([]bool, len(d.Abilities))
	m.WanderWait = w.rng.Float64() * 4
	if e.Inst == 0 && d.RespawnSec < 9000 {
		a := w.rng.Float64() * 2 * math.Pi
		r := w.rng.Float64() * 1.5
		e.X, e.Z = m.SX+math.Cos(a)*r, m.SZ+math.Sin(a)*r
	} else {
		e.X, e.Z = m.SX, m.SZ
	}
	e.Y = gamedata.Height(e.X, e.Z)
	e.R = w.rng.Float64() * 2 * math.Pi
}

func (w *World) moveToward(e *Entity, tx, tz, speed float64) float64 {
	dx, dz := tx-e.X, tz-e.Z
	d := math.Hypot(dx, dz)
	if d < 0.001 {
		return 0
	}
	step := math.Min(d, speed*w.dt)
	e.X += dx / d * step
	e.Z += dz / d * step
	e.R = math.Atan2(dx, dz)
	e.Y = gamedata.Height(e.X, e.Z)
	return d - step
}

// addThreat raises src's threat on a mob; the first threat pulls nearby pack members.
func (w *World) addThreat(mobE *Entity, srcID int, amt float64) {
	m := mobE.M
	if m == nil {
		return
	}
	first := len(m.Threat) == 0
	m.Threat[srcID] += amt
	m.Returning = false
	if m.Target == 0 {
		m.Target = srcID
	}
	if first {
		w.pullPack(mobE, srcID)
	}
}

// pullPack makes aggressive mobs standing close to e join the fight.
func (w *World) pullPack(e *Entity, targetID int) {
	if !e.M.Def.Aggressive && !e.M.Def.Elite {
		return
	}
	for _, o := range w.ents {
		if o == e || o.M == nil || o.Dead || o.Inst != e.Inst || len(o.M.Threat) > 0 || !o.M.Def.Aggressive {
			continue
		}
		if math.Hypot(o.X-e.X, o.Z-e.Z) < 11 {
			o.M.Threat[targetID] = 1
			o.M.Target = targetID
		}
	}
}

func (w *World) taunt(src, mobE *Entity) {
	m := mobE.M
	top := 0.0
	for _, v := range m.Threat {
		top = math.Max(top, v)
	}
	m.Threat[src.ID] = top*1.1 + 10
	m.Target = src.ID
	m.TauntBy = src.ID
	m.TauntUntil = w.now + 3
	m.Returning = false
}

func (w *World) pickTarget(e *Entity) *Entity {
	m := e.M
	var best *Entity
	bestT := 0.0
	for id, th := range m.Threat {
		t := w.ents[id]
		if t == nil || t.Dead || t.Inst != e.Inst || (t.P != nil && t.P.Flight != nil) ||
			math.Hypot(t.X-e.X, t.Z-e.Z) > m.Def.Leash*1.5 {
			delete(m.Threat, id)
			continue
		}
		if th > bestT {
			best, bestT = t, th
		}
	}
	if w.now < m.TauntUntil {
		if t := w.ents[m.TauntBy]; t != nil && !t.Dead && t.Inst == e.Inst {
			m.Target = t.ID
			return t
		}
	}
	if cur := w.ents[m.Target]; cur != nil && best != nil && cur != best && m.Threat[cur.ID] >= bestT/1.1 {
		best = cur
	}
	if best != nil {
		m.Target = best.ID
	} else {
		m.Target = 0
	}
	return best
}

func (w *World) tickMob(e *Entity) {
	m := e.M
	d := m.Def
	if e.Dead {
		if !m.Gone && w.now >= m.CorpseEnd {
			m.Gone = true
		}
		if d.RespawnSec >= 9000 { // dungeon mobs do not return
			if m.Gone {
				delete(w.ents, e.ID)
			}
			return
		}
		if w.now >= m.RespawnAt {
			w.respawnMob(e)
		}
		return
	}
	if e.stunned(w.now) {
		return
	}
	// aggressive mobs pick a target on sight
	if d.Aggressive && len(m.Threat) == 0 && !m.Returning {
		for _, p := range w.players {
			if p.Inst == e.Inst && !p.Dead && p.P.Flight == nil && dist(e, p) < d.AggroRange {
				w.addThreat(e, p.ID, 1)
				break
			}
		}
	}
	t := w.pickTarget(e)
	if t == nil {
		switch {
		case m.Returning:
			if w.moveToward(e, m.SX, m.SZ, d.Speed*1.4) < 0.5 {
				m.Returning = false
			}
		case len(m.Threat) > 0:
			w.evade(e)
		default:
			w.idleMob(e)
		}
		return
	}
	if e.Inst == 0 && math.Hypot(e.X-m.SX, e.Z-m.SZ) > d.Leash {
		w.evade(e)
		return
	}
	sp := d.Speed * e.speedMult(w.now)
	dd := dist(e, t)
	dmgMult := 1.0
	if m.Enraged {
		dmgMult = 1.5
	}
	if d.Range > 0 { // ranged attackers keep their distance
		if dd > d.Range*0.9 {
			w.moveToward(e, t.X, t.Z, sp)
		} else {
			e.R = math.Atan2(t.X-e.X, t.Z-e.Z)
		}
		if dd <= d.Range && w.now >= m.NextAtk {
			m.NextAtk = w.now + d.AttackSec*1.1
			w.sendNear(e, protocol.M{"t": "fx", "k": "mobshot", "s": e.ID, "d": t.ID, "style": d.Style})
			w.damage(e, t, float64(d.DmgMin+w.rng.Intn(d.DmgMax-d.DmgMin+1))*dmgMult, false)
		}
	} else {
		if dd > 2.2*math.Max(1, d.Scale*0.8) {
			w.moveToward(e, t.X, t.Z, sp)
		} else {
			e.R = math.Atan2(t.X-e.X, t.Z-e.Z)
		}
		if dd <= 2.8*math.Max(1, d.Scale*0.8) && w.now >= m.NextAtk {
			m.NextAtk = w.now + d.AttackSec
			w.sendNear(e, protocol.M{"t": "fx", "k": "bite", "s": e.ID, "d": t.ID})
			w.damage(e, t, float64(d.DmgMin+w.rng.Intn(d.DmgMax-d.DmgMin+1))*dmgMult, false)
		}
	}
	if len(d.Abilities) > 0 {
		w.tickBoss(e, t)
	}
}

func (w *World) idleMob(e *Entity) {
	m := e.M
	if e.Inst != 0 {
		return
	}
	if m.Wandering {
		if w.moveToward(e, m.WX, m.WZ, 1.6) < 0.2 {
			m.Wandering = false
			m.WanderWait = 3 + w.rng.Float64()*6
		}
		return
	}
	m.WanderWait -= w.dt
	if m.WanderWait <= 0 {
		a := w.rng.Float64() * 2 * math.Pi
		r := 1 + w.rng.Float64()*5
		m.WX, m.WZ = m.SX+math.Cos(a)*r, m.SZ+math.Sin(a)*r
		if gamedata.Height(m.WX, m.WZ) > 0.2 {
			m.Wandering = true
		} else {
			m.WanderWait = 2
		}
	}
}

// evade makes a mob give up, heal and walk home.
func (w *World) evade(e *Entity) {
	m := e.M
	m.Target = 0
	m.Returning = true
	m.Hit = map[int]bool{}
	m.Threat = map[int]float64{}
	m.Enraged = false
	m.AbNext = make([]float64, len(m.Def.Abilities))
	m.AbDone = make([]bool, len(m.Def.Abilities))
	e.HP = e.MaxHP
	e.Dots = nil
	e.SlowUntil = 0
	w.removeAdds(e)
}
