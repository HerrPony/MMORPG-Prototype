package game

import (
	"math"
	"slices"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/progression"
	"mmorpg/internal/protocol"
)

const globalCooldown = 1.0

var resNames = map[string]string{"rage": "Wut", "focus": "Fokus", "mana": "Mana"}

func (w *World) tickEffects(e *Entity) {
	if e.HPBoostEnd > 0 && w.now >= e.HPBoostEnd {
		e.HPBoostEnd = 0
		if e.P != nil {
			w.recalc(e)
			e.HP = min(e.HP, e.MaxHP)
		}
	}
	if len(e.Dots) == 0 {
		return
	}
	if e.Dead {
		e.Dots = nil
		return
	}
	var keep []Dot
	for _, d := range e.Dots {
		if w.now >= d.Next {
			d.Next += 1
			w.damage(w.ents[d.Src], e, d.Per, false)
		}
		if e.Dead {
			e.Dots = nil
			return
		}
		if w.now < d.Until {
			keep = append(keep, d)
		}
	}
	e.Dots = keep
}

// ---- helpers for finding enemies ----

func (w *World) enemiesNear(inst int, x, z, r float64) []*Entity {
	var out []*Entity
	for _, o := range w.ents {
		if o.Kind == KindMob && !o.Dead && o.Inst == inst && math.Hypot(o.X-x, o.Z-z) <= r {
			out = append(out, o)
		}
	}
	return out
}

func needsTarget(sk *gamedata.Skill) bool {
	switch sk.Kind {
	case gamedata.KindDamage, gamedata.KindDot, gamedata.KindAoETarget, gamedata.KindGround, gamedata.KindPoly,
		gamedata.KindTaunt, gamedata.KindCharge, gamedata.KindTame:
		return true
	}
	return false
}

// checkCast returns an error text if the skill cannot be used right now.
func (w *World) checkCast(e *Entity, sk *gamedata.Skill, t *Entity) string {
	p := e.P
	if p.Res < float64(sk.Cost) {
		return "Nicht genug " + resNames[p.Class.Resource] + "."
	}
	if needsTarget(sk) {
		if t == nil || t.Kind != KindMob || t.Dead || t.Inst != e.Inst {
			return "Du brauchst ein gültiges Ziel."
		}
		d := dist(e, t)
		if d > sk.Range+1.0 {
			return "Das Ziel ist außer Reichweite."
		}
		if sk.MinRange > 0 && d < sk.MinRange {
			return "Das Ziel ist zu nah."
		}
		if !w.hasLOS(e, t, true) {
			return "Kein Sichtkontakt zum Ziel."
		}
	}
	switch sk.Kind {
	case gamedata.KindAoE:
		if len(w.enemiesNear(e.Inst, e.X, e.Z, sk.Radius)) == 0 {
			return "Keine Feinde in der Nähe."
		}
	case gamedata.KindTame:
		switch {
		case !t.M.Def.Tameable:
			return "Dieses Wesen kann nicht gezähmt werden."
		case t.M.Def.Elite || t.M.Def.Boss:
			return "Elite-Gegner können nicht gezähmt werden."
		case t.Level > e.Level:
			return "Dieses Tier ist zu stark für dich."
		case p.Char.Pet != nil && p.PetEnt != 0:
			return "Du hast bereits ein Tier. Entlasse es zuerst."
		}
	case gamedata.KindCallPet:
		if p.Char.Pet == nil {
			return "Du hast noch kein gezähmtes Tier."
		}
	case gamedata.KindMendPet:
		if w.ents[p.PetEnt] == nil {
			return "Du hast kein Tier an deiner Seite."
		}
	}
	return ""
}

func (w *World) tryCast(e *Entity, id string) {
	p := e.P
	sk := gamedata.Skills[id]
	if sk == nil || !slices.Contains(p.Char.Skills, id) {
		return
	}
	switch {
	case e.Dead:
		w.errMsg(e, "Du bist tot.")
		return
	case p.Flight != nil:
		w.errMsg(e, "Du kannst im Flug nicht zaubern.")
		return
	}
	if sk.Basic && p.Auto { // pressing the basic attack again toggles auto-attack off
		p.Auto = false
		return
	}
	if p.Cast != nil {
		w.errMsg(e, "Du wirkst bereits.")
		return
	}
	if w.now < p.CD[id] || (!sk.Basic && sk.Kind != gamedata.KindStance && w.now < p.GCD) {
		w.errMsg(e, "Noch nicht bereit.")
		return
	}
	t := w.ents[p.Target]
	if msg := w.checkCast(e, sk, t); msg != "" {
		w.errMsg(e, msg)
		return
	}
	if sk.Cast > 0 {
		p.Cast = &CastState{Skill: id, End: w.now + sk.Cast, Target: p.Target, SX: e.X, SZ: e.Z}
		w.send(e, protocol.M{"t": "cast", "skill": id, "dur": sk.Cast})
		w.sendNear(e, protocol.M{"t": "fx", "k": id + ":cast", "s": e.ID, "d": p.Target})
		return
	}
	w.execSkill(e, sk, t)
	if sk.Basic {
		p.Auto = true
	}
	w.hunterAuto(e, sk, t)
}

// hunterAuto: hunters start their auto shot with any offensive skill used on a mob.
func (w *World) hunterAuto(e *Entity, sk *gamedata.Skill, t *Entity) {
	if e.P.Class.ID != "hunter" || sk.Basic || (sk.Power <= 0 && sk.DotPower <= 0) {
		return
	}
	if t != nil && t.Kind == KindMob && !t.Dead {
		e.P.Auto = true
	}
}

// finishCast re-checks a skill when its cast time ran out.
func (w *World) finishCast(e *Entity, sk *gamedata.Skill, t *Entity) bool {
	if e.Dead {
		return false
	}
	if msg := w.checkCast(e, sk, t); msg != "" {
		w.errMsg(e, msg)
		return false
	}
	return true
}

func (w *World) cancelCast(e *Entity, reason string) {
	e.P.Cast = nil
	w.send(e, protocol.M{"t": "castStop"})
	if reason != "" {
		w.errMsg(e, reason)
	}
}

func (w *World) calcDamage(e *Entity, power float64) float64 {
	gear := 1.0
	if e.P != nil {
		gear = e.P.Stats.damageMult(e.P.Char.Class)
		if e.P.Stance == "defensive" {
			gear *= 0.75
		}
	}
	return power * progression.DamageScale(e.Level) * e.dmgMult(w.now) * gear * (0.9 + 0.2*w.rng.Float64())
}

func (w *World) correctPlayer(e *Entity) {
	e.Y = gamedata.FootY(e.X, e.Z)
	e.P.LastMove = w.now
	w.send(e, protocol.M{"t": "correct", "x": e.X, "y": e.Y, "z": e.Z})
}

// execSkill performs a validated skill: spends resources, applies effects, broadcasts visuals.
func (w *World) execSkill(e *Entity, sk *gamedata.Skill, t *Entity) {
	p := e.P
	p.Res = math.Max(0, p.Res-float64(sk.Cost))
	p.CD[sk.ID] = w.now + sk.Cooldown
	gcd := 0.0
	if !sk.Basic && sk.Kind != gamedata.KindStance {
		gcd = globalCooldown
		p.GCD = w.now + gcd
	}
	w.send(e, protocol.M{"t": "cd", "skill": sk.ID, "dur": sk.Cooldown, "gcd": gcd})

	dstID := e.ID
	if t != nil && needsTarget(sk) {
		dstID = t.ID
	}
	if sk.Kind != gamedata.KindTame && sk.Kind != gamedata.KindCallPet && sk.Kind != gamedata.KindGround {
		w.sendNear(e, protocol.M{"t": "fx", "k": sk.ID, "s": e.ID, "d": dstID})
	}

	switch sk.Kind {
	case gamedata.KindDamage:
		amt := w.calcDamage(e, sk.Power)
		if sk.Execute > 0 && float64(t.HP) <= sk.Execute*float64(t.MaxHP) {
			amt *= sk.ExecMult
		}
		w.damage(e, t, amt, true)
		w.applySlow(t, sk)
		w.applyStun(t, sk.Stun)
	case gamedata.KindDot:
		w.damage(e, t, w.calcDamage(e, sk.Power), true)
		if !t.Dead {
			t.Dots = append(t.Dots, Dot{Src: e.ID, Per: w.calcDamage(e, sk.DotPower), Until: w.now + sk.Duration, Next: w.now + 1})
		}
	case gamedata.KindAoE:
		for _, o := range w.enemiesNear(e.Inst, e.X, e.Z, sk.Radius) {
			w.damage(e, o, w.calcDamage(e, sk.Power), true)
			w.applySlow(o, sk)
		}
	case gamedata.KindAoETarget:
		for _, o := range w.enemiesNear(e.Inst, t.X, t.Z, sk.Radius) {
			w.damage(e, o, w.calcDamage(e, sk.Power), true)
			w.applySlow(o, sk)
		}
	case gamedata.KindBuff:
		e.BuffMult = 1 + sk.Power
		e.BuffTaken = sk.Taken
		e.BuffUntil = w.now + sk.Duration
		w.sysChat(e, "combat", sk.Name+" aktiviert.")
	case gamedata.KindShield:
		e.Shield = sk.Power * progression.DamageScale(e.Level)
		e.ShieldUntil = w.now + sk.Duration
		w.sysChat(e, "combat", sk.Name+" absorbiert "+itoa(int(e.Shield))+" Schaden.")
	case gamedata.KindLastStand:
		e.HPBoostEnd = w.now + sk.Duration
		old := e.MaxHP
		e.MaxHP = int(float64(old) * (1 + sk.Power))
		e.HP += e.MaxHP - old
		w.sysChat(e, "combat", sk.Name+": maximales Leben um "+itoa(int(sk.Power*100))+"% erhöht.")
	case gamedata.KindStance:
		p.Stance = sk.Stance
		w.sysChat(e, "combat", sk.Name+" eingenommen.")
	case gamedata.KindTaunt:
		w.taunt(e, t)
		w.sysChat(e, "combat", "Du verspottest "+t.Name+".")
	case gamedata.KindCharge:
		// land next to the target and stun it
		a := math.Atan2(e.X-t.X, e.Z-t.Z)
		e.X, e.Z = t.X+math.Sin(a)*2.2, t.Z+math.Cos(a)*2.2
		e.R = a + math.Pi
		w.correctPlayer(e)
		w.applyStun(t, sk.Stun)
		w.addThreat(t, e.ID, 30)
		if p.Class.Resource == "rage" {
			p.Res = math.Min(float64(p.ResMax), p.Res+float64(sk.Gain))
		}
	case gamedata.KindBlink:
		ox, oz := e.X, e.Z
		for d := sk.Blink; d > 1; d -= 1 {
			nx, nz := ox+math.Sin(e.R)*d, oz+math.Cos(e.R)*d
			if e.Inst == 0 && gamedata.Walkable(nx, nz) || e.Inst != 0 && w.instWalkable(e, nx, nz) {
				e.X, e.Z = nx, nz
				break
			}
		}
		w.correctPlayer(e)
	case gamedata.KindPoly:
		t.PolyUntil = w.now + sk.Duration
		t.StunUntil = t.PolyUntil
		t.M.Target = 0
		t.M.Threat = map[int]float64{}
	case gamedata.KindGround:
		w.addGround(e, t.X, t.Z, sk)
	case gamedata.KindTame:
		w.tame(e, t)
	case gamedata.KindCallPet:
		w.callPet(e)
	case gamedata.KindMendPet:
		if pet := w.ents[p.PetEnt]; pet != nil {
			pet.HP = min(pet.MaxHP, pet.HP+int(sk.Power*float64(pet.MaxHP)))
		}
	case gamedata.KindElemental:
		w.summonElemental(e, sk)
	}
	if sk.Gain > 0 && p.Class.Resource == "rage" && sk.Kind != gamedata.KindCharge {
		p.Res = math.Min(float64(p.ResMax), p.Res+float64(sk.Gain))
	}
	switch sk.Kind {
	case gamedata.KindBuff, gamedata.KindShield, gamedata.KindStance, gamedata.KindLastStand, gamedata.KindBlink, gamedata.KindMendPet, gamedata.KindCallPet:
	default:
		p.LastCombat = w.now
	}
}

func (w *World) applySlow(t *Entity, sk *gamedata.Skill) {
	if sk.Slow > 0 && !t.Dead {
		t.SlowMult = sk.Slow
		t.SlowUntil = w.now + sk.Duration
	}
}

func (w *World) applyStun(t *Entity, secs float64) {
	if secs > 0 && !t.Dead && !t.M.Def.Boss {
		t.StunUntil = math.Max(t.StunUntil, w.now+secs)
	}
}

// ownerOf returns the player credited for damage dealt by src (pets credit their owner).
func (w *World) ownerOf(src *Entity) *Entity {
	if src == nil {
		return nil
	}
	if src.P != nil {
		return src
	}
	if src.Pet != nil {
		return w.players[src.Pet.Owner]
	}
	return nil
}

// damage applies damage from src (may be nil for expired DoT sources) to t.
func (w *World) damage(src, t *Entity, amt float64, canCrit bool) {
	if t.Dead {
		return
	}
	if t.P != nil && t.P.Flight != nil {
		return
	}
	if t.P != nil && src != nil { // armor soaks damage taken by players
		amt *= 1 - armorReduction(t.P.Stats.Armor, src.Level)
	}
	if t.P != nil || t.Pet != nil {
		amt *= t.takenMult(w.now)
	}
	crit := canCrit && w.rng.Float64() < 0.1
	if crit {
		amt *= 1.8
	}
	a := max(1, int(math.Round(amt)))
	absorbed := 0
	if w.now < t.ShieldUntil && t.Shield > 0 {
		absorbed = min(a, int(t.Shield))
		t.Shield -= float64(absorbed)
		a -= absorbed
	}
	t.HP -= a
	srcID := 0
	if src != nil {
		srcID = src.ID
	}
	w.sendNear(t, protocol.M{"t": "dmg", "id": t.ID, "a": a, "ab": absorbed, "c": crit, "s": srcID})

	owner := w.ownerOf(src)
	if owner != nil {
		owner.P.LastCombat = w.now
		if src == owner {
			suffix := ""
			if crit {
				suffix = " (kritisch)"
			}
			w.sysChat(owner, "combat", "Du triffst "+t.Name+" für "+itoa(a)+" Schaden"+suffix+".")
		} else {
			w.sysChat(owner, "combat", src.Name+" trifft "+t.Name+" für "+itoa(a)+" Schaden.")
		}
	}
	if t.P != nil {
		t.P.LastCombat = w.now
		if t.P.Class.Resource == "rage" {
			t.P.Res = math.Min(float64(t.P.ResMax), t.P.Res+float64(a)*0.5+1)
		}
		who := "Etwas"
		if src != nil {
			who = src.Name
		}
		w.sysChat(t, "combat", who+" trifft dich für "+itoa(a)+" Schaden.")
	}
	if t.M != nil && src != nil {
		if owner != nil {
			t.M.Hit[owner.ID] = true
		}
		mult := 1.0
		if src.P != nil && src.P.Stance == "defensive" {
			mult = 2.0
		}
		w.addThreat(t, src.ID, float64(a)*mult+1)
		if t.PolyUntil > w.now { // damage breaks polymorph
			t.PolyUntil, t.StunUntil = 0, 0
		}
	}
	if t.HP <= 0 {
		t.HP = 0
		w.die(t, src)
	}
}

func (w *World) die(t, killer *Entity) {
	t.Dead = true
	t.Dots = nil
	switch {
	case t.P != nil:
		t.P.Auto = false
		if t.P.Cast != nil {
			w.cancelCast(t, "")
		}
		w.sysChat(t, "combat", "Du bist gestorben.")
	case t.Pet != nil:
		w.petDied(t)
	case t.M != nil:
		w.mobDied(t)
	}
}
