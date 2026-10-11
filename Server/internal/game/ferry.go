package game

import (
	"fmt"
	"math"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
)

// The ferry sails back and forth between the two harbors on a fixed timetable.
// It lies alongside the pier, so players simply walk onto its deck (no interaction needed), can move
// around freely on deck while it sails and walk off again over the pier at the other side.
// The hull is symmetric and the ship never turns: it just sails bow or stern first.
const (
	ferrySpeed   = 9.0         // m/s
	ferryDocked  = 45.0        // seconds at each pier
	ferryHeading = math.Pi / 2 // the ship always lies along the east-west axis
	deckHalfX    = 1.5         // deck half width / length in ship coordinates
	deckHalfZ    = 6.0
	deckMargin   = 0.35 // tolerance when stepping on or off
	deckY        = 1.2
)

// Ferry is the state of the ship (stored in the ship entity).
type Ferry struct {
	From, To *gamedata.Town // current leg: the harbor it lies at / sails to
	Sailing  bool
	Until    float64 // docked: departure time
	Start    float64 // sailing: time the leg started
	Dur      float64
	Warned   bool
	SX, SZ   float64
	DX, DZ   float64
}

func (w *World) spawnFerry() {
	var a, b *gamedata.Town
	for _, t := range gamedata.Towns {
		if !t.Harbor {
			continue
		}
		if a == nil {
			a = t
		} else if t.Cont != a.Cont && b == nil {
			b = t
		}
	}
	if a == nil || b == nil {
		return
	}
	e := &Entity{ID: w.newID(), Kind: KindShip, Name: "Fähre", Model: "ship_cog", Level: 1, HP: 1, MaxHP: 1, R: ferryHeading}
	e.X, e.Z = gamedata.Berth(a)
	e.Ship = &Ferry{From: a, To: b, Until: w.now + 30}
	w.ents[e.ID] = e
}

// deckToWorld converts a deck position (ship coordinates) into a world position.
func deckToWorld(ship *Entity, lx, lz float64) (float64, float64) {
	s, c := math.Sincos(ship.R)
	return ship.X + lx*c + lz*s, ship.Z - lx*s + lz*c
}

// worldToDeck is the inverse of deckToWorld.
func worldToDeck(ship *Entity, x, z float64) (float64, float64) {
	s, c := math.Sincos(ship.R)
	dx, dz := x-ship.X, z-ship.Z
	return dx*c - dz*s, dx*s + dz*c
}

func onDeck(ship *Entity, x, z, margin float64) bool {
	lx, lz := worldToDeck(ship, x, z)
	return math.Abs(lx) <= deckHalfX+margin && math.Abs(lz) <= deckHalfZ+margin
}

// shipAtPos returns a ferry whose deck covers the world position.
func (w *World) shipAtPos(x, z float64) *Entity {
	for _, s := range w.ents {
		if s.Ship != nil && onDeck(s, x, z, deckMargin) {
			return s
		}
	}
	return nil
}

func (w *World) tickFerry(e *Entity) {
	f := e.Ship
	if !f.Sailing {
		left := f.Until - w.now
		e.Title = fmt.Sprintf("nach %s · Abfahrt in %d s", f.To.Name, int(math.Ceil(math.Max(0, left))))
		if left <= 10 && !f.Warned {
			f.Warned = true
			w.tellPassengers(e, "Das Schiff legt in 10 Sekunden nach "+f.To.Name+" ab!")
		}
		if left <= 0 {
			f.Sailing = true
			f.SX, f.SZ = e.X, e.Z
			f.DX, f.DZ = gamedata.Berth(f.To)
			f.Dur = math.Hypot(f.DX-f.SX, f.DZ-f.SZ) / ferrySpeed
			f.Start = w.now
			w.tellPassengers(e, "Die Fähre legt ab. Nächster Halt: "+f.To.Name+".")
		}
	} else {
		e.Title = "auf See nach " + f.To.Name
		k := math.Min(1, (w.now-f.Start)/math.Max(f.Dur, 1))
		e.X = f.SX + (f.DX-f.SX)*k
		e.Z = f.SZ + (f.DZ-f.SZ)*k
		if k >= 1 {
			f.From, f.To = f.To, f.From
			f.Sailing = false
			f.Until = w.now + ferryDocked
			f.Warned = false
			w.tellPassengers(e, "Das Schiff hat in "+f.From.Name+" angelegt. Du kannst über den Steg von Bord gehen.")
		}
	}
	// carry the passengers
	for _, p := range w.players {
		if p.P.Aboard == e.ID {
			p.X, p.Z = deckToWorld(e, p.P.DeckX, p.P.DeckZ)
			p.Y = deckY
		}
	}
}

func (w *World) tellPassengers(ship *Entity, text string) {
	for _, p := range w.players {
		if p.P.Aboard == ship.ID {
			w.sysChat(p, "yell", text)
		}
	}
}

// boardShip is called when a player stepped onto the deck of a ferry.
func (w *World) boardShip(e, ship *Entity, x, z float64) {
	p := e.P
	p.Aboard = ship.ID
	p.DeckX, p.DeckZ = worldToDeck(ship, x, z)
	p.DeckX = clampf(p.DeckX, -deckHalfX, deckHalfX)
	p.DeckZ = clampf(p.DeckZ, -deckHalfZ, deckHalfZ)
	e.X, e.Z = deckToWorld(ship, p.DeckX, p.DeckZ)
	e.Y = deckY
	if p.Cast != nil {
		w.cancelCast(e, "")
	}
	f := ship.Ship
	if f.Sailing {
		w.sysChat(e, "sys", "Du bist an Bord. Die Fähre fährt nach "+f.To.Name+".")
	} else {
		w.sysChat(e, "sys", fmt.Sprintf("Du bist an Bord. Die Fähre fährt in %d Sekunden nach %s.", int(math.Max(0, math.Ceil(f.Until-w.now))), f.To.Name))
	}
}

// movePassenger handles a move message of a player who rides on a ferry.
func (w *World) movePassenger(e *Entity, m protocol.Move) {
	p := e.P
	ship := w.ents[p.Aboard]
	if ship == nil || ship.Ship == nil {
		p.Aboard = 0
		return
	}
	f := ship.Ship
	if m.AB { // walking on the deck (coordinates relative to the ship)
		if !finite(m.LX, m.LZ) {
			return
		}
		p.DeckX = clampf(m.LX, -deckHalfX-deckMargin, deckHalfX+deckMargin)
		p.DeckZ = clampf(m.LZ, -deckHalfZ-deckMargin, deckHalfZ+deckMargin)
		e.X, e.Z = deckToWorld(ship, p.DeckX, p.DeckZ)
		e.Y = deckY
		e.R = m.R
		p.LastMove = w.now
		return
	}
	// a world position: stepping off the ship onto the pier or the shore (only while it lies at a pier)
	if !f.Sailing && !onDeck(ship, m.X, m.Z, 0.05) && gamedata.Walkable(m.X, m.Z) && math.Hypot(m.X-e.X, m.Z-e.Z) < 6 {
		p.Aboard = 0
		e.X, e.Z, e.R = m.X, m.Z, m.R
		e.Y = clampf(m.Y, 0, 4)
		p.LastMove = w.now
		p.ZoneKey, p.PoiID = "", ""
		w.sysChat(e, "sys", "Du gehst in "+f.From.Name+" von Bord.")
	}
}

// ferryInfo is shown in the dialog of the harbor captains.
func (w *World) ferryInfo(e *Entity) protocol.M {
	for _, s := range w.ents {
		if s.Ship == nil {
			continue
		}
		f := s.Ship
		status := ""
		if f.Sailing {
			status = "Die Fähre ist gerade unterwegs nach " + f.To.Name + "."
		} else {
			status = fmt.Sprintf("Die Fähre liegt in %s und fährt in %d Sekunden nach %s ab.", f.From.Name, int(math.Max(0, math.Ceil(f.Until-w.now))), f.To.Name)
		}
		return protocol.M{"ship": s.ID, "status": status, "docked": !f.Sailing}
	}
	return protocol.M{}
}
