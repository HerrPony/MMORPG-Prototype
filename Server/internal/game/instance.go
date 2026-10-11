package game

import (
	"fmt"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
)

// Instance is one copy of a dungeon for a party (or a single player).
type Instance struct {
	ID         int
	Def        *gamedata.Dungeon
	Key        string
	EmptySince float64
}

func (w *World) instanceKey(e *Entity, dungeonID string) string {
	if e.P.PartyID != 0 {
		return fmt.Sprintf("%s:party:%d", dungeonID, e.P.PartyID)
	}
	return fmt.Sprintf("%s:solo:%s", dungeonID, e.P.Char.ID)
}

func (w *World) instWalkable(e *Entity, x, z float64) bool {
	if in := w.insts[e.Inst]; in != nil {
		return in.Def.Walkable(x, z)
	}
	return true
}

func (w *World) createInstance(def *gamedata.Dungeon, key string) *Instance {
	in := &Instance{ID: w.nextInst, Def: def, Key: key}
	w.nextInst++
	w.insts[in.ID] = in
	w.instByKey[key] = in
	ox, oz := def.Origin[0], def.Origin[1]
	for _, s := range def.Spawns {
		w.newMob(s.Mob, ox+s.X, oz+s.Z, in.ID)
	}
	exit := w.newObject("portal_exit", ox+def.Start[0]-9, oz+def.Start[1]-1, in.ID)
	exit.Y = 0
	return in
}

func (w *World) sendInstanceInfo(e *Entity, in *Instance) {
	if in == nil {
		w.send(e, protocol.M{"t": "instance", "on": false})
		return
	}
	w.send(e, protocol.M{"t": "instance", "on": true, "name": in.Def.Name, "rooms": in.Def.Rooms, "style": in.Def.Style,
		"origin": in.Def.Origin, "start": in.Def.Start})
}

func (w *World) enterDungeon(e *Entity, dungeonID string) {
	def := gamedata.Dungeons[dungeonID]
	if def == nil || e.Inst != 0 {
		return
	}
	if e.Level < def.MinLevel {
		w.errMsg(e, "Du benötigst mindestens Stufe "+itoa(def.MinLevel)+" für diesen Dungeon.")
		return
	}
	key := w.instanceKey(e, dungeonID)
	in := w.instByKey[key]
	if in == nil {
		in = w.createInstance(def, key)
	}
	e.P.Char.ReturnX, e.P.Char.ReturnZ = e.X, e.Z
	if e.P.Cast != nil {
		w.cancelCast(e, "")
	}
	e.P.Auto = false
	e.P.Target = 0
	e.Inst = in.ID
	e.X, e.Z = def.Origin[0]+def.Start[0], def.Origin[1]+def.Start[1]
	e.R = 3.14159
	in.EmptySince = 0
	w.sendInstanceInfo(e, in)
	w.correctPlayer(e)
	w.sysChat(e, "sys", "Du betrittst: "+def.Name)
	w.send(e, protocol.M{"t": "zone", "zone": def.Name, "sub": "", "lo": def.MinLevel, "hi": def.MaxLevel, "instance": true})
}

func (w *World) exitDungeon(e *Entity) {
	if e.Inst == 0 {
		return
	}
	e.Inst = 0
	e.P.Target = 0
	e.X, e.Z = e.P.Char.ReturnX+3, e.P.Char.ReturnZ+3
	if !gamedata.Walkable(e.X, e.Z) {
		e.X, e.Z = e.P.Char.ReturnX, e.P.Char.ReturnZ
	}
	w.sendInstanceInfo(e, nil)
	w.correctPlayer(e)
	e.P.ZoneKey, e.P.PoiID = "", "" // announce the zone again
}

// tickInstances removes dungeon copies that have been empty for a while.
func (w *World) tickInstances() {
	for id, in := range w.insts {
		n := 0
		for _, p := range w.players {
			if p.Inst == id {
				n++
			}
		}
		if n > 0 {
			in.EmptySince = 0
			continue
		}
		if in.EmptySince == 0 {
			in.EmptySince = w.now
		}
		if w.now-in.EmptySince > 120 {
			for eid, e := range w.ents {
				if e.Inst == id {
					delete(w.ents, eid)
				}
			}
			delete(w.insts, id)
			delete(w.instByKey, in.Key)
		}
	}
}
