package progression

import "testing"

func TestAddXPLevelsUp(t *testing.T) {
	lvl, xp, ups := AddXP(1, 90, 20) // needs 100 for level 2
	if lvl != 2 || xp != 10 || ups != 1 {
		t.Fatalf("got lvl=%d xp=%d ups=%d", lvl, xp, ups)
	}
	lvl, _, ups = AddXP(1, 0, 100+240)
	if lvl != 3 || ups != 2 {
		t.Fatalf("multi level-up: lvl=%d ups=%d", lvl, ups)
	}
}

func TestKillXPScaling(t *testing.T) {
	if KillXP(20, 1, 1) != 30 {
		t.Fatal("same level should give base+10*lvl")
	}
	if KillXP(20, 1, 10) >= KillXP(20, 1, 1) {
		t.Fatal("low level mobs should give less xp")
	}
}
