package game

import (
	"math"
	"slices"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
)

// tickZones announces zone / place changes and records discoveries.
func (w *World) tickZones() {
	for _, e := range w.players {
		if e.Inst != 0 {
			continue
		}
		p := e.P
		b := gamedata.BiomeAt(e.X, e.Z)
		zone := gamedata.ZoneByBiome(b)
		poi := gamedata.PoiAt(e.X, e.Z)
		pid, pname := "", ""
		if poi != nil {
			pid, pname = poi.ID, poi.Name
		}
		if zone.Key == p.ZoneKey && pid == p.PoiID {
			continue
		}
		p.ZoneKey, p.PoiID = zone.Key, pid
		w.send(e, protocol.M{"t": "zone", "zone": zone.Name, "sub": pname, "lo": zone.MinLvl, "hi": zone.MaxLvl, "sea": b == gamedata.BiomeSea})
		if b != gamedata.BiomeSea {
			w.discover(e, "zone:"+zone.Key, zone.Name, 0)
		}
		if poi != nil {
			w.discover(e, poi.ID, poi.Name, 15+e.Level*4)
			w.questEvent(e, "explore", poi.ID)
		}
	}
}

func (w *World) discover(e *Entity, id, name string, xp int) {
	c := e.P.Char
	if slices.Contains(c.Discovered, id) {
		return
	}
	c.Discovered = append(c.Discovered, id)
	w.send(e, protocol.M{"t": "disc", "ids": c.Discovered, "new": id, "name": name})
	if xp > 0 {
		w.sysChat(e, "level", "Neu entdeckt: "+name)
		w.giveXP(e, xp)
	}
}

// ---- griffin flights ----

// townOfNPC finds the settlement a griffin master belongs to.
func townOfNPC(n *gamedata.NPC) *gamedata.Town {
	var best *gamedata.Town
	bd := 1e9
	for _, t := range gamedata.Towns {
		if t.Kind == "site" {
			continue
		}
		if d := math.Hypot(t.X-n.X, t.Z-n.Z); d < bd {
			best, bd = t, d
		}
	}
	return best
}

func flightPrice(a, b *gamedata.Town) int {
	ax, az := gamedata.PadPos(a)
	bx, bz := gamedata.PadPos(b)
	return max(15, int(math.Hypot(ax-bx, az-bz)*0.4))
}

func (w *World) flightInfo(e *Entity, n *Entity) protocol.M {
	here := townOfNPC(n.N)
	var list []protocol.M
	for _, id := range e.P.Char.FlightPoints {
		t := gamedata.TownByID(id)
		if t == nil || t.ID == here.ID {
			continue
		}
		list = append(list, protocol.M{"id": t.ID, "name": t.Name, "price": flightPrice(here, t)})
	}
	return protocol.M{"here": here.ID, "hereName": here.Name, "points": list}
}

// unlockFlightPoint remembers a town's griffin master after the player talked to them.
func (w *World) unlockFlightPoint(e, n *Entity) {
	t := townOfNPC(n.N)
	if t == nil || slices.Contains(e.P.Char.FlightPoints, t.ID) {
		return
	}
	e.P.Char.FlightPoints = append(e.P.Char.FlightPoints, t.ID)
	w.sysChat(e, "level", "Flugpunkt freigeschaltet: "+t.Name)
	w.send(e, protocol.M{"t": "fp", "ids": e.P.Char.FlightPoints})
}

func (w *World) onFlight(e *Entity, m protocol.FlightMsg) {
	n := w.npcNear(e, m.NPC)
	if n == nil || !n.N.Flight || e.P.Flight != nil || e.Inst != 0 {
		return
	}
	here := townOfNPC(n.N)
	dest := gamedata.TownByID(m.Dest)
	if dest == nil || dest.ID == here.ID || !slices.Contains(e.P.Char.FlightPoints, dest.ID) {
		w.errMsg(e, "Dorthin kennst du noch keinen Flugpunkt.")
		return
	}
	price := flightPrice(here, dest)
	if e.P.Char.Copper < price {
		w.errMsg(e, "Du hast nicht genug Geld.")
		return
	}
	e.P.Char.Copper -= price
	w.sendInventory(e)
	if e.P.Cast != nil {
		w.cancelCast(e, "")
	}
	e.P.Auto = false
	dx, dz := gamedata.PadPos(dest)
	d := math.Hypot(dx-e.X, dz-e.Z)
	e.P.Flight = &Flight{Dest: dest.ID, SX: e.X, SZ: e.Z, DX: dx, DZ: dz, Dist: d, Speed: 34}
	w.sysChat(e, "sys", "Du fliegst nach "+dest.Name+".")
	w.send(e, protocol.M{"t": "flightStart", "dest": dest.Name, "dur": d / 34})
}

func (w *World) tickFlight(e *Entity) {
	f := e.P.Flight
	f.Done += f.Speed * w.dt
	k := math.Min(1, f.Done/math.Max(f.Dist, 1))
	e.X = f.SX + (f.DX-f.SX)*k
	e.Z = f.SZ + (f.DZ-f.SZ)*k
	e.R = math.Atan2(f.DX-f.SX, f.DZ-f.SZ)
	ramp := math.Min(math.Min(f.Done, f.Dist-f.Done)/3, 1)
	ground := gamedata.Height(e.X, e.Z)
	e.Y = ground + 14*math.Max(0, ramp) + 8*math.Max(0, ramp)*math.Sin(k*math.Pi) + math.Max(0, 6-ground*0)
	if k >= 1 {
		e.P.Flight = nil
		e.Y = gamedata.Height(e.X, e.Z)
		e.P.LastMove = w.now
		w.send(e, protocol.M{"t": "correct", "x": e.X, "y": e.Y, "z": e.Z})
		w.send(e, protocol.M{"t": "flightEnd"})
		w.sysChat(e, "sys", "Du bist angekommen.")
	}
}
