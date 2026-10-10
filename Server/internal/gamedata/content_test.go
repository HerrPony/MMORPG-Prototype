package gamedata

import (
	"testing"
)

// The quest data is large; this test catches typos and unreachable objectives.
func TestQuestData(t *testing.T) {
	Init()
	spawned := map[string]int{}
	for _, s := range AllSpawns() {
		spawned[s.Mob]++
	}
	for _, d := range Dungeons {
		for _, s := range d.Spawns {
			spawned[s.Mob]++
		}
	}
	obtainable := map[string]bool{}
	for _, m := range Mobs {
		for _, l := range m.Loot {
			obtainable[l.Item] = true
		}
	}
	for _, o := range Objects {
		if o.Item != "" {
			obtainable[o.Item] = true
		}
	}
	objSpawned := map[string]int{}
	for _, s := range ObjectSpawns {
		objSpawned[s.Obj]++
	}
	for id, q := range Quests {
		if q.Giver == "" || NPCByID(q.Giver) == nil {
			t.Errorf("quest %s has no valid giver (%q)", id, q.Giver)
		}
		if q.Requires != "" && Quests[q.Requires] == nil {
			t.Errorf("quest %s requires unknown quest %s", id, q.Requires)
		}
		if q.GiveItem != "" && Items[q.GiveItem] == nil {
			t.Errorf("quest %s gives unknown item %s", id, q.GiveItem)
		}
		if q.Hand != "" && NPCByID(q.Hand) == nil {
			t.Errorf("quest %s hand-in NPC %s unknown", id, q.Hand)
		}
		for _, it := range q.Reward.Items {
			if Items[it] == nil {
				t.Errorf("quest %s rewards unknown item %s", id, it)
			}
		}
		for _, sk := range append(append([]string{}, q.GrantSkill...), q.Reward.Skills...) {
			if Skills[sk] == nil {
				t.Errorf("quest %s: unknown skill %s", id, sk)
			}
		}
		if q.TimeLimit > 0 && q.GiveItem == "" {
			t.Logf("note: timed quest %s gives no item", id)
		}
		for _, o := range q.Objectives {
			switch o.Kind {
			case "kill", "tame":
				if Mobs[o.Target] == nil {
					t.Errorf("quest %s: unknown mob %s", id, o.Target)
				} else if spawned[o.Target] == 0 {
					t.Errorf("quest %s: mob %s never spawns", id, o.Target)
				}
			case "collect":
				if Items[o.Target] == nil {
					t.Errorf("quest %s: unknown item %s", id, o.Target)
				} else if !obtainable[o.Target] {
					t.Errorf("quest %s: item %s cannot be obtained", id, o.Target)
				}
			case "deliver", "talk":
				if NPCByID(o.Target) == nil {
					t.Errorf("quest %s: unknown NPC %s", id, o.Target)
				}
			case "explore":
				if POIByID(o.Target) == nil {
					t.Errorf("quest %s: unknown POI %s", id, o.Target)
				}
			default:
				t.Errorf("quest %s: unknown objective kind %s", id, o.Kind)
			}
		}
	}
	// quest-gated loot and objects must name real quests
	for _, m := range Mobs {
		for _, l := range m.Loot {
			if Items[l.Item] == nil {
				t.Errorf("mob %s drops unknown item %s", m.ID, l.Item)
			}
			if l.Quest != "" && Quests[l.Quest] == nil {
				t.Errorf("mob %s: loot gated by unknown quest %s", m.ID, l.Quest)
			}
		}
	}
	for id, o := range Objects {
		if o.Quest != "" && Quests[o.Quest] == nil {
			t.Errorf("object %s: unknown quest %s", id, o.Quest)
		}
		if o.Item != "" && objSpawned[id] == 0 {
			t.Errorf("object %s is never placed", id)
		}
	}
	for _, g := range groups() {
		if Mobs[g.Mob] == nil {
			t.Errorf("group of unknown mob %s", g.Mob)
		}
	}
	for _, n := range NPCs {
		for _, id := range n.Quests {
			if Quests[id] == nil {
				t.Errorf("NPC %s offers unknown quest %s", n.ID, id)
			}
		}
		for _, id := range n.Stock {
			if Items[id] == nil {
				t.Errorf("NPC %s sells unknown item %s", n.ID, id)
			}
		}
		if n.Vendor && len(n.Stock) == 0 {
			t.Errorf("vendor %s has no stock", n.ID)
		}
	}
	t.Logf("%d quests, %d mobs, %d NPCs, %d items, %d spawns", len(Quests), len(Mobs), len(NPCs), len(Items), len(AllSpawns()))
}

func TestMobCountsPerGroup(t *testing.T) {
	Init()
	want := map[string]int{}
	for _, g := range groups() {
		want[g.Mob] += g.N
	}
	got := map[string]int{}
	for _, s := range generatedSpawns {
		got[s.Mob]++
	}
	for m, n := range want {
		if got[m] < n {
			t.Errorf("mob %s: placed %d of %d", m, got[m], n)
		}
	}
}

func TestSkillLevelsEvery2(t *testing.T) {
	for _, c := range Classes {
		seen := map[int]bool{}
		for _, id := range c.Skills {
			s := Skills[id]
			if s == nil {
				t.Fatalf("class %s: unknown skill %s", c.ID, id)
			}
			if !s.Quest && s.Level > 1 {
				if seen[s.Level] || s.Level%2 != 0 {
					t.Errorf("class %s: skill %s at odd or duplicate level %d", c.ID, id, s.Level)
				}
				seen[s.Level] = true
			}
		}
	}
}
