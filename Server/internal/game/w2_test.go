package game

import (
	"math"
	"math/rand"
	"strings"
	"testing"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/loot"
	"mmorpg/internal/progression"
	"mmorpg/internal/protocol"
)

func shipEntity(w *World) *Entity {
	for _, e := range w.ents {
		if e.Ship != nil {
			return e
		}
	}
	return nil
}

// walk is a series of small steps to a world position, like a client would send them
func walk(w *World, e *Entity, x, z float64) {
	for i := 0; i < 400; i++ {
		dx, dz := x-e.X, z-e.Z
		d := math.Hypot(dx, dz)
		if d < 0.05 {
			return
		}
		step := math.Min(d, 1.2)
		w.tick()
		w.onMove(e, protocol.Move{X: e.X + dx/d*step, Y: e.Y, Z: e.Z + dz/d*step, R: math.Atan2(dx, dz)})
	}
}

// the ferry lies alongside the pier on a timetable: players just walk onto it, move around on deck and walk off again
func TestFerryTimetable(t *testing.T) {
	w := newTestWorld(t)
	ship := shipEntity(w)
	if ship == nil {
		t.Fatal("no ferry")
	}
	e, _ := addPlayer(t, w, "Anna", "mage", 20)
	from := ship.Ship.From
	e.X, e.Z = gamedata.PadPos(from)
	e.Y = gamedata.Height(e.X, e.Z)
	if !gamedata.Walkable(e.X, e.Z) {
		t.Fatal("pad not walkable")
	}
	// along the pier and across the rail onto the deck
	walk(w, e, ship.X-3, from.Z)
	if e.P.Aboard != 0 {
		t.Fatal("still on the pier")
	}
	walk(w, e, ship.X-2, ship.Z-1)
	if e.P.Aboard != ship.ID {
		t.Fatalf("could not walk onto the ship (at %.1f,%.1f, ship %.1f,%.1f)", e.X, e.Z, ship.X, ship.Z)
	}
	// free movement on the deck
	w.onMove(e, protocol.Move{AB: true, LX: 0.5, LZ: -3})
	if math.Abs(e.P.DeckZ+3) > 0.01 {
		t.Fatal("moving on deck did not work")
	}
	// the ship leaves; the passenger moves along and stays on deck
	sx := ship.X
	for i := 0; i < 20*200 && !ship.Ship.Sailing; i++ {
		w.tick()
	}
	if !ship.Ship.Sailing {
		t.Fatal("the ferry never left")
	}
	for i := 0; i < 20*200 && ship.Ship.Sailing; i++ {
		w.tick()
		w.onMove(e, protocol.Move{AB: true, LX: float64(i%3) - 1, LZ: float64(i%9) - 4})
		if math.Hypot(e.X-ship.X, e.Z-ship.Z) > 8 || e.Y < 1 || e.P.Aboard != ship.ID {
			t.Fatalf("passenger fell off the ship: %.1f,%.1f vs %.1f,%.1f", e.X, e.Z, ship.X, ship.Z)
		}
		// walking off at sea is not possible
		if i == 100 {
			w.onMove(e, protocol.Move{X: e.X + 3, Z: e.Z + 3, Y: 1})
			if e.P.Aboard == 0 {
				t.Fatal("must not leave the ship at sea")
			}
		}
	}
	if math.Abs(ship.X-sx) < 100 || gamedata.ContinentOf(ship.X) == gamedata.ContinentOf(sx) {
		t.Fatalf("the ferry did not cross the sea (%.0f -> %.0f)", sx, ship.X)
	}
	// docked: walk over the pier to the shore
	to := ship.Ship.From
	walk(w, e, ship.X+3, to.Z)
	if e.P.Aboard != 0 {
		t.Fatalf("could not step onto the pier (at %.1f,%.1f)", e.X, e.Z)
	}
	px, pz := gamedata.PadPos(to)
	walk(w, e, px, pz)
	if dist(e, &Entity{X: px, Z: pz}) > 1 || !gamedata.Walkable(e.X, e.Z) || gamedata.ContinentOf(e.X) != 1 {
		t.Fatalf("did not reach the shore: %.0f,%.0f", e.X, e.Z)
	}
	// and the ferry comes back
	first := ship.Ship.From
	run(w, 130)
	if ship.Ship.From == first && !ship.Ship.Sailing {
		t.Fatal("the ferry should have sailed back")
	}
}
func TestGriffinStaysOnContinent(t *testing.T) {
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "mage", 25)
	suki := w.npcByID("suki")
	e.X, e.Z = suki.X, suki.Z+2
	e.P.Char.Copper = 5000
	e.P.Char.FlightPoints = []string{"sonnenburg", "himmelspforte"}
	w.onFlight(e, protocol.FlightMsg{NPC: suki.ID, Dest: "sonnenburg"})
	if e.P.Flight != nil {
		t.Fatal("griffins cannot cross the sea")
	}
	w.onFlight(e, protocol.FlightMsg{NPC: suki.ID, Dest: "himmelspforte"})
	if e.P.Flight == nil {
		t.Fatal("flight to the capital of the same continent should start")
	}
}

func TestNewDungeons(t *testing.T) {
	for id, name := range map[string]string{"temple": "zhulkarat", "fire": "ignaros"} {
		w := newTestWorld(t)
		def := gamedata.Dungeons[id]
		e, _ := addPlayer(t, w, "Anna", "warrior", def.MinLevel)
		w.enterDungeon(e, id)
		if e.Inst == 0 {
			t.Fatalf("%s: did not enter", id)
		}
		var boss *Entity
		bosses := 0
		for _, m := range w.ents {
			if m.M != nil && m.Inst == e.Inst && m.M.Def.Boss {
				bosses++
				if m.M.Def.ID == name {
					boss = m
				}
			}
		}
		if boss == nil || bosses != 3 {
			t.Fatalf("%s: want 2 minibosses and the end boss, got %d bosses (final found: %v)", id, bosses, boss != nil)
		}
		// the end boss drops two pieces of the dungeon set
		bag := loot.Roll(boss.M.Def, boss.Level, loot.Looter{Class: "hunter"}, rand.New(rand.NewSource(3)))
		sets := 0
		for _, it := range bag.Items {
			if strings.HasSuffix(it.Item, "_"+boss.M.Def.Set) {
				sets++
			}
		}
		if sets != 2 {
			t.Fatalf("%s: expected two set items, got %v", id, bag.Items)
		}
		w.exitDungeon(e)
		if e.Inst != 0 {
			t.Fatal("exit failed")
		}
	}
	// the level requirement is enforced
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "mage", 20)
	w.enterDungeon(e, "fire")
	if e.Inst != 0 {
		t.Fatal("a level 20 character must not enter the Feuerschlund")
	}
}

func TestLevelCap40(t *testing.T) {
	lvl, xp := 1, 0
	for i := 0; i < 1000 && lvl < progression.MaxLevel; i++ {
		lvl, xp, _ = progression.AddXP(lvl, xp, 5000)
	}
	if lvl != 40 || progression.MaxLevel != 40 {
		t.Fatalf("level cap should be 40, reached %d", lvl)
	}
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "warrior", 40)
	w.giveXP(e, 1000000)
	if e.Level != 40 || e.P.Char.XP != 0 {
		t.Fatal("no experience beyond the cap")
	}
	// high level skills can be trained at the new capital's trainer
	goro := w.npcByID("goro")
	e.X, e.Z = goro.X+1, goro.Z
	e.P.Char.Copper = 1000000
	w.onTrain(e, protocol.TrainMsg{NPC: goro.ID, Skill: "w_worldbreaker"})
	found := false
	for _, s := range e.P.Char.Skills {
		if s == "w_worldbreaker" {
			found = true
		}
	}
	if !found {
		t.Fatal("level 40 skill could not be learned")
	}
}

func TestHunterSkillStartsAutoShot(t *testing.T) {
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "hunter", 5)
	wolf := findMob(w, "wolf")
	e.X, e.Z = wolf.X-10, wolf.Z
	e.Y = gamedata.Height(e.X, e.Z)
	e.P.Target = wolf.ID
	e.P.Char.Skills = append(e.P.Char.Skills, "h_aimed")
	w.tryCast(e, "h_aimed")
	run(w, 3)
	if !e.P.Auto {
		t.Fatal("an offensive skill should start the hunter's auto shot")
	}
	m, _ := addPlayer(t, w, "Bert", "mage", 5)
	m.X, m.Z, m.Y = e.X, e.Z, e.Y
	m.P.Target = wolf.ID
	w.tryCast(m, "m_bolt")
	run(w, 3)
	if m.P.Auto {
		t.Fatal("only hunters start auto attacks with skills")
	}
}
