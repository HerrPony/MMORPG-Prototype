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
	ox, oz := e.X, e.Z
	e.X += dx / d * step
	e.Z += dz / d * step
	if e.Inst == 0 { // trees, rocks and buildings are solid for mobs and pets too
		rad := 0.5
		if e.M != nil {
			rad = 0.35 + 0.25*math.Max(1, e.M.Def.Scale)
		}
		e.X, e.Z = gamedata.Collide(e.X, e.Z, rad)
		if step > 0.01 && math.Hypot(e.X-ox, e.Z-oz) < step*0.3 {
			e.Stuck += w.dt
		} else {
			e.Stuck = 0
		}
	}
	e.R = math.Atan2(dx, dz)
	e.Y = gamedata.Height(e.X, e.Z)
	return d - step
}

// walkTo moves e towards (tx,tz) like moveToward, but walks around obstacles: with a clear line it goes straight,
// otherwise it follows an A* route over the navigation grid (recomputed when the target moves or the route
// gets old). Dungeons have no grid and use the straight line. Returns the distance left to the current waypoint.
func (w *World) walkTo(e *Entity, tx, tz, speed float64) float64 {
	n := &e.Nav
	if e.Inst != 0 || math.Hypot(tx-e.X, tz-e.Z) < 1.5 || gamedata.NavLineFree(e.X, e.Z, tx, tz) {
		n.Path = nil
		return w.moveToward(e, tx, tz, speed)
	}
	stale := n.Path == nil || math.Hypot(tx-n.GX, tz-n.GZ) > 3 || w.now-n.At > 2.5
	if stale && w.now >= n.Fail && w.pathBudget > 0 {
		w.pathBudget--
		path, _, ok := gamedata.FindPath(e.X, e.Z, tx, tz, 20000)
		if ok {
			n.Path, n.GX, n.GZ, n.At = path, tx, tz, w.now
		} else {
			n.Path, n.Fail = nil, w.now+1 // unreachable: try again in a second, meanwhile slide along
		}
	}
	for len(n.Path) > 0 && math.Hypot(n.Path[0][0]-e.X, n.Path[0][1]-e.Z) < 0.7 {
		n.Path = n.Path[1:]
	}
	if len(n.Path) == 0 {
		return w.moveToward(e, tx, tz, speed)
	}
	return w.moveToward(e, n.Path[0][0], n.Path[0][1], speed)
}

// hasLOS reports whether a has a clear line to b (open world only). Trees and rocks block ranged attacks,
// buildings block everything.
func (w *World) hasLOS(a, b *Entity, trees bool) bool {
	if a.Inst != 0 || b.Inst != 0 {
		return true
	}
	return !gamedata.LineBlocked(a.X, a.Z, b.X, b.Z, trees)
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

// packRadius is how close an aggressive mob has to stand to join a fight that starts next to it.
const packRadius = 6.5

// pullPack makes aggressive mobs standing close to e join the fight.
func (w *World) pullPack(e *Entity, targetID int) {
	if !e.M.Def.Aggressive && !e.M.Def.Elite {
		return
	}
	for _, o := range w.ents {
		if o == e || o.M == nil || o.Dead || o.Inst != e.Inst || len(o.M.Threat) > 0 || !o.M.Def.Aggressive {
			continue
		}
		if math.Hypot(o.X-e.X, o.Z-e.Z) < packRadius {
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
			if p.Inst == e.Inst && !p.Dead && p.P.Flight == nil && dist(e, p) < d.AggroRange && w.hasLOS(e, p, false) {
				w.addThreat(e, p.ID, 1)
				break
			}
		}
	}
	t := w.pickTarget(e)
	if t == nil {
		switch {
		case m.Returning:
			w.walkTo(e, m.SX, m.SZ, d.Speed*1.4)
			if math.Hypot(e.X-m.SX, e.Z-m.SZ) < 0.5 || e.Stuck > 2 {
				if e.Stuck > 2 { // wedged behind an obstacle: beam home
					e.X, e.Z = m.SX, m.SZ
					e.Y = gamedata.Height(e.X, e.Z)
				}
				e.Stuck = 0
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
	if e.Stuck > 5 { // cannot reach the target (cliff, wall): give up
		e.Stuck = 0
		w.evade(e)
		return
	}
	sp := d.Speed * e.speedMult(w.now)
	dd := dist(e, t)
	los := w.hasLOS(e, t, d.Range > 0)
	dmgMult := 1.0
	if m.Enraged {
		dmgMult = 1.5
	}
	if d.Range > 0 { // ranged attackers keep their distance
		if dd > d.Range*0.9 || !los {
			w.walkTo(e, t.X, t.Z, sp)
		} else {
			e.R = math.Atan2(t.X-e.X, t.Z-e.Z)
		}
		if dd <= d.Range && los && w.now >= m.NextAtk {
			m.NextAtk = w.now + d.AttackSec*1.1
			w.sendNear(e, protocol.M{"t": "fx", "k": "mobshot", "s": e.ID, "d": t.ID, "style": d.Style})
			w.damage(e, t, float64(d.DmgMin+w.rng.Intn(d.DmgMax-d.DmgMin+1))*dmgMult, false)
		}
	} else {
		if dd > 2.2*math.Max(1, d.Scale*0.8) {
			w.walkTo(e, t.X, t.Z, sp)
		} else {
			e.R = math.Atan2(t.X-e.X, t.Z-e.Z)
		}
		if dd <= 2.8*math.Max(1, d.Scale*0.8) && los && w.now >= m.NextAtk {
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
		r := 1 + w.rng.Float64()*2.5
		m.WX, m.WZ = m.SX+math.Cos(a)*r, m.SZ+math.Sin(a)*r
		if gamedata.Height(m.WX, m.WZ) > 0.2 && gamedata.NavLineFree(e.X, e.Z, m.WX, m.WZ) {
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
