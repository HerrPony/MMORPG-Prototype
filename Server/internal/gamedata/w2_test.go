package gamedata

import "testing"

func TestWorldPlacements(t *testing.T) {
	Init()
	for _, n := range NPCs {
		if !Walkable(n.X, n.Z) {
			t.Errorf("NPC %s stands on unwalkable ground (%.0f,%.0f h=%.1f)", n.ID, n.X, n.Z, Height(n.X, n.Z))
		}
	}
	for _, tw := range Towns {
		if tw.Kind == "site" {
			continue
		}
		x, z := PadPos(tw)
		if !Walkable(x, z) {
			t.Errorf("landing place of %s is not walkable", tw.ID)
		}
	}
	for _, s := range AllSpawns() {
		if !Walkable(s.X, s.Z) {
			t.Errorf("mob %s spawns on unwalkable ground (%.0f,%.0f)", s.Mob, s.X, s.Z)
		}
		if d := DistToRoad(s.X, s.Z); d < roadClear-0.5 {
			t.Errorf("mob %s spawns on a road (%.0f,%.0f, %.1f m from it)", s.Mob, s.X, s.Z, d)
		}
	}
	for _, o := range ObjectSpawns {
		if !Walkable(o.X, o.Z) {
			t.Errorf("object %s placed on unwalkable ground (%.0f,%.0f)", o.Obj, o.X, o.Z)
		}
	}
	for id, d := range Dungeons {
		for _, s := range d.Spawns {
			if !d.Walkable(d.Origin[0]+s.X, d.Origin[1]+s.Z) {
				t.Errorf("dungeon %s: %s spawns outside the rooms", id, s.Mob)
			}
			if Mobs[s.Mob] == nil {
				t.Errorf("dungeon %s: unknown mob %s", id, s.Mob)
			}
		}
		if !Walkable(d.Entrance[0], d.Entrance[1]) {
			t.Errorf("dungeon %s: entrance not walkable", id)
		}
	}
}

// every class can learn a new skill every 2 levels up to level 40 (level 10 comes from the class quest)
func TestSkillsUpTo40(t *testing.T) {
	for _, c := range Classes {
		have := map[int]bool{}
		for _, id := range c.Skills {
			if s := Skills[id]; !s.Quest && s.Level > 1 {
				have[s.Level] = true
			}
		}
		for lvl := 2; lvl <= MaxLevel; lvl += 2 {
			if lvl == 10 {
				continue
			}
			if !have[lvl] {
				t.Errorf("class %s has no skill for level %d", c.ID, lvl)
			}
		}
	}
}

func TestQuestAreas(t *testing.T) {
	Init()
	n := 0
	for id, q := range Quests {
		for i, o := range q.Objectives {
			if len(q.Areas) <= i || len(q.Areas[i]) == 0 {
				t.Errorf("quest %s objective %d (%s %s) has no map area", id, i, o.Kind, o.Target)
			}
		}
		if q.HandIn == nil {
			t.Errorf("quest %s has no hand-in area", id)
		}
		n++
	}
	t.Logf("%d quests checked", n)
}

func TestLevelRange(t *testing.T) {
	Init()
	best := 0
	for _, q := range Quests {
		best = max(best, q.Level)
	}
	if best < 40 {
		t.Errorf("highest quest level is %d, want 40", best)
	}
	lv := map[int]int{}
	for _, q := range Quests {
		lv[q.Level]++
	}
	for l := 1; l <= 40; l++ {
		if lv[l] == 0 {
			t.Errorf("no quest for level %d", l)
		}
	}
}
