package game

import (
	"math"
	"strconv"
	"strings"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
)

// Boss abilities are listed in the mob definition as "name:arg:arg...":
//
//	heavy:<sec>              a crushing blow on the main target every <sec> seconds
//	aoe:<sec>                damages all players close to the boss
//	poison:<sec>             poisons a random player
//	summon:<hp%>:<mob>:<n>   calls <n> adds once the boss falls below <hp%>
//	enrage:<hp%>             deals 50% more damage below <hp%>
func (w *World) tickBoss(e, tgt *Entity) {
	m := e.M
	d := m.Def
	avg := float64(d.DmgMin+d.DmgMax) / 2
	pct := float64(e.HP) / float64(e.MaxHP) * 100
	for i, spec := range d.Abilities {
		p := strings.Split(spec, ":")
		num := func(k int) float64 {
			if k < len(p) {
				f, _ := strconv.ParseFloat(p[k], 64)
				return f
			}
			return 0
		}
		switch p[0] {
		case "heavy", "aoe", "poison":
			iv := num(1)
			if m.AbNext[i] == 0 {
				m.AbNext[i] = w.now + iv*0.6
			}
			if w.now < m.AbNext[i] {
				continue
			}
			m.AbNext[i] = w.now + iv
			switch p[0] {
			case "heavy":
				w.sendNear(e, protocol.M{"t": "fx", "k": "bossheavy", "s": e.ID, "d": tgt.ID})
				w.yell(e, d.Name+" holt zu einem vernichtenden Schlag aus!")
				w.damage(e, tgt, avg*2.4, false)
			case "aoe":
				w.sendNear(e, protocol.M{"t": "fx", "k": "bossnova", "s": e.ID, "d": e.ID})
				w.yell(e, d.Name+" lässt den Boden beben!")
				for _, pl := range w.players {
					if pl.Inst == e.Inst && !pl.Dead && dist(e, pl) < 14 {
						w.damage(e, pl, avg*1.2, false)
					}
				}
			case "poison":
				var cands []*Entity
				for _, pl := range w.players {
					if pl.Inst == e.Inst && !pl.Dead && dist(e, pl) < 30 {
						cands = append(cands, pl)
					}
				}
				if len(cands) > 0 {
					v := cands[w.rng.Intn(len(cands))]
					v.Dots = append(v.Dots, Dot{Src: e.ID, Per: avg * 0.4, Until: w.now + 10, Next: w.now + 1})
					w.sendNear(e, protocol.M{"t": "fx", "k": "bossnova", "s": e.ID, "d": v.ID})
					w.yell(e, d.Name+" vergiftet "+v.Name+"!")
				}
			}
		case "summon":
			if m.AbDone[i] || pct > num(1) {
				continue
			}
			m.AbDone[i] = true
			n := int(num(3))
			w.yell(e, d.Name+" ruft Verstärkung!")
			for k := 0; k < n; k++ {
				a := float64(k)/float64(n)*2*math.Pi + 0.5
				x, z := e.X+math.Cos(a)*5, e.Z+math.Sin(a)*5
				add := w.newMob(p[2], x, z, e.Inst)
				add.M.Threat[tgt.ID] = 1
				add.M.Target = tgt.ID
				m.Adds = append(m.Adds, add.ID)
			}
		case "enrage":
			if !m.Enraged && pct <= num(1) {
				m.Enraged = true
				w.yell(e, d.Name+" wird wütend!")
			}
		}
	}
}

// removeAdds removes summoned helpers when a boss resets.
func (w *World) removeAdds(e *Entity) {
	for _, id := range e.M.Adds {
		if a := w.ents[id]; a != nil && a.M != nil {
			delete(w.ents, id)
		}
	}
	e.M.Adds = nil
}

// yell shows a boss message to everyone nearby.
func (w *World) yell(e *Entity, text string) {
	w.sendNear(e, protocol.M{"t": "chat", "ch": "yell", "text": text})
}

var _ = gamedata.Mobs
