// Package game is the authoritative simulation: one goroutine owns all state
// and processes network commands and the fixed-step tick.
package game

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"slices"
	"time"

	"mmorpg/internal/config"
	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
	"mmorpg/internal/storage"
)

type World struct {
	cfg   config.Config
	store storage.Store
	rng   *rand.Rand

	cmds    chan func()
	ents    map[int]*Entity
	players map[int]*Entity
	byChar  map[string]*Entity
	nextID  int
	now     float64
	tickN   int
	dt      float64

	parties   map[int]*Party
	nextParty int
	insts     map[int]*Instance
	instByKey map[string]*Instance
	nextInst  int
	grounds   []*GroundFX

	sceneryJSON json.RawMessage
	defsJSON    json.RawMessage
	lastSave    float64
	pathBudget  int // path searches left in this tick
}

func New(cfg config.Config, store storage.Store) *World {
	gamedata.Init()
	if ed, err := gamedata.LoadEdits(editsPath(cfg)); err != nil {
		log.Printf("world_edits.json: %v (ignoriert)", err)
	} else if err := gamedata.ApplyEdits(ed); err != nil {
		log.Printf("world_edits.json ungültig: %v (ignoriert)", err)
		_ = gamedata.ApplyEdits(&gamedata.Edits{})
	}
	w := &World{
		cfg:       cfg,
		store:     store,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		cmds:      make(chan func(), 2048),
		ents:      map[int]*Entity{},
		players:   map[int]*Entity{},
		byChar:    map[string]*Entity{},
		parties:   map[int]*Party{},
		insts:     map[int]*Instance{},
		instByKey: map[string]*Instance{},
		nextID:    1,
		nextInst:  1,
		nextParty: 1,
		dt:        cfg.TickRate.Seconds(),
	}
	w.refreshStatic()
	w.spawnStatic()
	return w
}

// refreshStatic rebuilds the cached scenery/definition payloads sent to clients.
func (w *World) refreshStatic() {
	w.sceneryJSON, _ = json.Marshal(gamedata.BuildScenery())
	w.defsJSON, _ = json.Marshal(protocol.M{
		"items":  gamedata.Items,
		"quests": gamedata.Quests,
		"skills": gamedata.Skills,
		"pois":   gamedata.POIs,
	})
}

func (w *World) newID() int {
	id := w.nextID
	w.nextID++
	return id
}

func (w *World) spawnStatic() {
	for _, n := range gamedata.NPCs {
		e := &Entity{ID: w.newID(), Kind: KindNPC, Name: n.Name, Title: n.Title, Model: n.Model, Level: 20,
			X: n.X, Z: n.Z, R: n.R, HP: 1000, MaxHP: 1000, N: n}
		e.Y = gamedata.Height(e.X, e.Z)
		w.ents[e.ID] = e
	}
	for _, s := range gamedata.AllSpawns() {
		w.newMob(s.Mob, s.X, s.Z, 0)
	}
	for _, o := range gamedata.ObjectSpawns {
		w.newObject(o.Obj, o.X, o.Z, 0)
	}
	w.spawnFerry()
}

func (w *World) newMob(id string, x, z float64, inst int) *Entity {
	def := gamedata.Mobs[id]
	e := &Entity{ID: w.newID(), Kind: KindMob, Name: def.Name, Model: def.Model, Inst: inst,
		M: &Mob{Def: def, SX: x, SZ: z, Hit: map[int]bool{}, Threat: map[int]float64{}}}
	w.respawnMob(e)
	w.ents[e.ID] = e
	return e
}

func (w *World) newObject(id string, x, z float64, inst int) *Entity {
	def := gamedata.Objects[id]
	e := &Entity{ID: w.newID(), Kind: KindObj, Name: def.Name, Model: def.Model, Inst: inst, X: x, Z: z, Level: 1,
		HP: 1, MaxHP: 1, O: &Obj{Def: def}}
	if inst == 0 {
		e.Y = gamedata.Height(x, z)
	}
	w.ents[e.ID] = e
	return e
}

// Do queues f to run on the world goroutine.
func (w *World) Do(f func()) { w.cmds <- f }

// Run drives the simulation until ctx is cancelled.
func (w *World) Run(ctx context.Context) {
	t := time.NewTicker(w.cfg.TickRate)
	defer t.Stop()
	for {
		select {
		case f := <-w.cmds:
			f()
		case <-t.C:
			w.tick()
		case <-ctx.Done():
			for _, p := range w.players {
				w.savePlayer(p)
			}
			return
		}
	}
}

func (w *World) tick() {
	w.now += w.dt
	w.tickN++
	w.pathBudget = 6
	for _, e := range w.ents {
		w.tickEffects(e)
		switch e.Kind {
		case KindMob:
			w.tickMob(e)
		case KindPlayer:
			w.tickPlayer(e)
		case KindPet:
			w.tickPet(e)
		case KindObj:
			w.tickObject(e)
		case KindShip:
			w.tickFerry(e)
		}
	}
	w.tickGrounds()
	if w.tickN%10 == 0 {
		w.tickZones()
	}
	if w.tickN%20 == 0 {
		w.tickQuestTimers()
		w.tickParties()
		w.tickInstances()
	}
	if w.tickN%w.cfg.SnapshotEach == 0 {
		w.sendSnapshots()
	}
	if w.now-w.lastSave >= w.cfg.AutosaveEach.Seconds() {
		w.lastSave = w.now
		for _, p := range w.players {
			w.savePlayer(p)
		}
	}
}

// ---- messaging ----

func (w *World) send(e *Entity, m protocol.M) {
	if e.P != nil {
		e.P.Out.Send(m)
	}
}

func (w *World) broadcast(m protocol.M) {
	for _, p := range w.players {
		p.P.Out.Send(m)
	}
}

// sendNear delivers a message to all players that can see src (same instance, close by).
func (w *World) sendNear(src *Entity, m protocol.M) {
	for _, p := range w.players {
		if p.Inst == src.Inst && math.Abs(p.X-src.X) < InterestDist && math.Abs(p.Z-src.Z) < InterestDist {
			p.P.Out.Send(m)
		}
	}
}

func (w *World) sysChat(e *Entity, ch, text string) {
	w.send(e, protocol.M{"t": "chat", "ch": ch, "text": text})
}

func (w *World) errMsg(e *Entity, text string) {
	w.send(e, protocol.M{"t": "err", "text": text})
}

func (w *World) broadcastExcept(ex *Entity, m protocol.M) {
	for _, p := range w.players {
		if p != ex {
			p.P.Out.Send(m)
		}
	}
}

// ---- snapshots ----

type entView struct {
	ID   int               `json:"id"`
	K    string            `json:"k"`
	N    string            `json:"n"`
	Ti   string            `json:"ti,omitempty"`
	M    string            `json:"m"`
	C    string            `json:"c,omitempty"`
	L    int               `json:"l"`
	X    float64           `json:"x"`
	Y    float64           `json:"y"`
	Z    float64           `json:"z"`
	R    float64           `json:"r"`
	HP   int               `json:"hp"`
	MH   int               `json:"mh"`
	D    bool              `json:"d,omitempty"`
	Tg   int               `json:"tg,omitempty"`
	St   []string          `json:"st,omitempty"`
	Look *storage.Look     `json:"lk,omitempty"`
	El   bool              `json:"el,omitempty"`
	Bs   bool              `json:"bs,omitempty"`
	Sc   float64           `json:"sc,omitempty"`
	Ow   int               `json:"ow,omitempty"`
	Fl   bool              `json:"fl,omitempty"`
	Q    string            `json:"q,omitempty"`
	Gone bool              `json:"gn,omitempty"`
	Eq   map[string]string `json:"eq,omitempty"`
	Pt   bool              `json:"pt,omitempty"`
}

var visualSlots = []string{"mainhand", "offhand", "head", "shoulder", "chest", "back", "legs", "feet", "hands", "waist"}

func (w *World) view(e *Entity) entView {
	v := entView{ID: e.ID, K: e.Kind, N: e.Name, Ti: e.Title, M: e.Model, L: e.Level,
		X: round2(e.X), Y: round2(e.Y), Z: round2(e.Z), R: round2(e.R), HP: e.HP, MH: e.MaxHP, D: e.Dead}
	if e.Kind == KindNPC || e.Kind == KindObj || e.Kind == KindShip {
		v.HP, v.MH = 1, 1
	}
	if e.P != nil {
		v.C = e.P.Char.Class
		lk := e.P.Char.Look
		v.Look = &lk
		v.Tg = e.P.Target
		v.Fl = e.P.Flight != nil
		for _, s := range visualSlots {
			if id := e.P.Char.Equipment[s]; id != "" {
				if v.Eq == nil {
					v.Eq = map[string]string{}
				}
				v.Eq[s] = id
			}
		}
		v.Pt = e.P.PartyID != 0
	}
	if e.Ship != nil {
		v.Fl = e.Ship.Sailing
	}
	if e.M != nil {
		v.Tg = e.M.Target
		v.El = e.M.Def.Elite
		v.Bs = e.M.Def.Boss
		if e.M.Def.Scale != 1 {
			v.Sc = e.M.Def.Scale
		}
		v.Gone = e.M.Gone
		if w.now < e.PolyUntil {
			v.M = "sheep"
		}
	}
	if e.Pet != nil {
		v.Ow = e.Pet.Owner
	}
	if e.O != nil {
		v.Q = e.O.Def.Quest
		v.Gone = e.O.Gone
	}
	if w.now < e.SlowUntil {
		v.St = append(v.St, "slow")
	}
	if len(e.Dots) > 0 {
		v.St = append(v.St, "dot")
	}
	if w.now < e.BuffUntil {
		v.St = append(v.St, "buff")
	}
	if w.now < e.ShieldUntil && e.Shield > 0 {
		v.St = append(v.St, "shield")
	}
	if w.now < e.StunUntil {
		v.St = append(v.St, "stun")
	}
	if e.P != nil && e.P.Stance == "defensive" {
		v.St = append(v.St, "defstance")
	}
	return v
}

func (w *World) sendSnapshots() {
	if len(w.players) == 0 {
		return
	}
	type rec struct {
		e   *Entity
		raw []byte
	}
	recs := make([]rec, 0, len(w.ents))
	for _, e := range w.ents {
		if (e.M != nil && e.M.Gone) || (e.O != nil && e.O.Gone) {
			continue
		}
		b, _ := json.Marshal(w.view(e))
		recs = append(recs, rec{e, b})
	}
	for _, pe := range w.players {
		p := pe.P
		buf := make([]byte, 0, 16384)
		buf = append(buf, '[')
		first := true
		party := w.partyOf(pe)
		for _, r := range recs {
			e := r.e
			if e.Inst != pe.Inst {
				continue
			}
			near := math.Abs(e.X-pe.X) < InterestDist && math.Abs(e.Z-pe.Z) < InterestDist
			if !near && !(party != nil && e.Kind == KindPlayer && slices.Contains(party.Members, e.ID)) {
				continue
			}
			if !first {
				buf = append(buf, ',')
			}
			first = false
			buf = append(buf, r.raw...)
		}
		buf = append(buf, ']')
		var lt []int
		for _, e := range w.ents {
			if e.M != nil && e.Dead && !e.M.Gone && e.Inst == pe.Inst && !e.M.Loot[pe.ID].Empty() {
				lt = append(lt, e.ID)
			}
		}
		me := protocol.M{
			"hp": pe.HP, "mh": pe.MaxHP, "res": int(p.Res), "rm": p.ResMax,
			"xp": p.Char.XP, "xn": xpNext(pe.Level), "lvl": pe.Level,
			"tg": p.Target, "auto": p.Auto, "d": pe.Dead, "fl": p.Flight != nil, "ab": p.Aboard,
		}
		if p.Class.ID == "warrior" {
			me["stance"] = p.Stance
		}
		if p.Char.Pet != nil || p.ElemEnt != 0 {
			pet := protocol.M{"name": ""}
			if p.Char.Pet != nil {
				pet["name"], pet["mob"], pet["dead"] = p.Char.Pet.Name, p.Char.Pet.Mob, p.Char.Pet.Dead
			}
			if pt := w.ents[p.PetEnt]; pt != nil {
				pet["id"], pet["hp"], pet["mh"] = pt.ID, pt.HP, pt.MaxHP
			}
			if el := w.ents[p.ElemEnt]; el != nil {
				pet["eid"], pet["ehp"], pet["emh"] = el.ID, el.HP, el.MaxHP
			}
			me["pet"] = pet
		}
		msg := protocol.M{"t": "s", "e": json.RawMessage(buf), "me": me, "lt": lt}
		if party != nil {
			msg["pt"] = w.partyStates(party)
		}
		p.Out.Send(msg)
	}
}

// ---- enter / leave ----

// Enter puts a character into the world. Called from the gateway goroutine.
func (w *World) Enter(c *storage.Character, out Sender) *Entity {
	ch := make(chan *Entity, 1)
	w.Do(func() { ch <- w.enter(c, out) })
	return <-ch
}

// Leave removes a player (disconnect).
func (w *World) Leave(id int) {
	w.Do(func() { w.leave(id) })
}

func (w *World) enter(c *storage.Character, out Sender) *Entity {
	if old, ok := w.byChar[c.ID]; ok { // same character logged in elsewhere
		w.send(old, protocol.M{"t": "kicked", "text": "Du hast dich von einem anderen Ort aus angemeldet."})
		w.leave(old.ID)
	}
	class := gamedata.Classes[c.Class]
	if c.Level < 1 {
		c.Level = 1
	}
	if c.Quests == nil {
		c.Quests = map[string]*storage.QuestState{}
	}
	for _, id := range slices.Backward(gamedata.StartSkills(class)) { // new characters know the basic skills
		if !slices.Contains(c.Skills, id) {
			c.Skills = append([]string{id}, c.Skills...)
		}
	}
	e := &Entity{ID: w.newID(), Kind: KindPlayer, Name: c.Name, Model: c.Race, Level: c.Level,
		X: c.X, Y: c.Y, Z: c.Z, R: c.R}
	if !gamedata.Swimmable(c.X, c.Z) || (math.Abs(c.X) < 0.001 && math.Abs(c.Z) < 0.001) {
		e.X, e.Z, e.R = gamedata.SpawnPoint.X, gamedata.SpawnPoint.Z, gamedata.SpawnPoint.R
	}
	e.Y = gamedata.FootY(e.X, e.Z)
	e.P = &Player{Char: c, Out: out, Class: class, CD: map[string]float64{}, LastMove: w.now,
		Marks: map[int]string{}, Stance: "battle"}
	w.recalc(e)
	e.HP = e.MaxHP
	e.P.Res = float64(class.ResStart)
	if class.Resource != "rage" {
		e.P.Res = float64(e.P.ResMax)
	}
	w.ents[e.ID] = e
	w.players[e.ID] = e
	w.byChar[c.ID] = e

	w.send(e, protocol.M{
		"t": "world", "id": e.ID,
		"self": protocol.M{"name": c.Name, "race": c.Race, "class": c.Class, "look": c.Look},
		"cls":  class, "skills": c.Skills, "defs": w.defsJSON, "scenery": w.sceneryJSON, "now": time.Now().UnixMilli(),
		"disc": c.Discovered, "fp": c.FlightPoints,
	})
	w.sendInventory(e)
	w.sendEquipment(e)
	w.sendSets(e)
	w.sendQuests(e)
	w.sendMarks(e)
	w.sysChat(e, "sys", "Willkommen in der Welt, "+c.Name+"!")
	w.broadcastExcept(e, protocol.M{"t": "chat", "ch": "sys", "text": c.Name + " hat die Welt betreten."})
	return e
}

func (w *World) leave(id int) {
	e, ok := w.players[id]
	if !ok {
		return
	}
	w.leaveParty(e, false)
	w.dismissPets(e)
	w.savePlayer(e)
	delete(w.players, id)
	delete(w.ents, id)
	if w.byChar[e.P.Char.ID] == e {
		delete(w.byChar, e.P.Char.ID)
	}
	w.broadcast(protocol.M{"t": "chat", "ch": "sys", "text": e.Name + " hat die Welt verlassen."})
}

func (w *World) savePlayer(e *Entity) {
	c := e.P.Char
	if ship := w.ents[e.P.Aboard]; e.P.Aboard != 0 && ship != nil && ship.Ship != nil { // on the ferry: remember the pier it left from
		c.X, c.Z = gamedata.PadPos(ship.Ship.From)
		c.Y = gamedata.Height(c.X, c.Z)
	} else if e.Inst != 0 { // inside a dungeon: remember the entrance
		c.X, c.Z = c.ReturnX, c.ReturnZ
		c.Y = gamedata.Height(c.X, c.Z)
	} else {
		c.X, c.Y, c.Z, c.R = e.X, e.Y, e.Z, e.R
	}
	c.Level = e.Level
	if err := w.store.SaveCharacter(c); err != nil {
		log.Printf("save %s: %v", c.Name, err)
	}
}
