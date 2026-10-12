package gamedata

import "testing"

func TestNPCNotInBuilding(t *testing.T) {
	initColliders()
	for _, n := range NPCs {
		for _, b := range boxes {
			if n.X > b.x-b.hw && n.X < b.x+b.hw && n.Z > b.z-b.hd && n.Z < b.z+b.hd {
				t.Errorf("NPC %s inside a building at %.0f,%.0f", n.ID, n.X, n.Z)
			}
		}
		if NavBlocked(n.X, n.Z) {
			t.Logf("note: NPC %s stands in a blocked nav cell (%.0f,%.0f)", n.ID, n.X, n.Z)
		}
	}
}
