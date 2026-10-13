package gamedata

import (
	"math"
	"testing"
)

// nobody stands inside a building, a wall or a tower (the castle city of Sonnenburg is full of them)
func TestNPCsStandOutsideOfBuildings(t *testing.T) {
	Init()
	for _, n := range NPCs {
		for _, b := range Buildings() {
			switch b.Kind {
			case "pad", "dock", "ship", "lantern", "pond", "gate", "archgate", "fountain":
				continue
			}
			if math.Abs(n.X-b.X) < b.W/2+0.8 && math.Abs(n.Z-b.Z) < b.D/2+0.8 {
				t.Errorf("NPC %s stands in a %s at %.0f,%.0f", n.ID, b.Kind, b.X, b.Z)
			}
		}
	}
}

// the starting village and the other settlements stay free of mobs
func TestMobsKeepAwayFromSettlements(t *testing.T) {
	for _, s := range AllSpawns() {
		for _, tw := range Towns {
			slack := 40.0 // spawns that were pushed out of water or off a road may drift back a little
			if tw.ID == "eichenhain" {
				slack = 8
			}
			if r := safeRadius(tw); r > 0 && math.Hypot(s.X-tw.X, s.Z-tw.Z) < r-slack {
				t.Errorf("%s spawns %.0f m from %s (safe radius %.0f)", s.Mob, math.Hypot(s.X-tw.X, s.Z-tw.Z), tw.ID, r)
			}
		}
	}
	// the world is about WorldScale times as wide as before
	if WorldHalf < 440 || WorldHalf > 460 {
		t.Fatalf("unexpected world size %.0f", WorldHalf)
	}
}
