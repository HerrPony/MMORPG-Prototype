package gamedata

import (
	"slices"
	"testing"
)

func TestSecondQuestWave(t *testing.T) {
	Init()
	byLevel := map[int]int{}
	n := 0
	for id, q := range Quests {
		byLevel[q.Level]++
		if len(id) > 2 && id[:2] == "n_" {
			n++
		}
	}
	if n < 140 {
		t.Errorf("only %d new quests, want at least 140", n)
	}
	if len(Quests) < 2*140 {
		t.Errorf("%d quests in total, want at least 280", len(Quests))
	}
	for l := 1; l <= MaxLevel; l++ {
		if byLevel[l] < 5 {
			t.Errorf("level %d has only %d quests", l, byLevel[l])
		}
	}
	t.Logf("%d quests (%d new), per level: %v", len(Quests), n, byLevel)
}

// Every tier of gear must be for sale somewhere: weapons, armor and jewellery.
func TestVendorsCoverEveryTier(t *testing.T) {
	Init()
	for tier := 0; tier < tiers; tier++ {
		for _, kind := range []string{"weapon", "armor", "jewel"} {
			found := false
			for _, n := range NPCs {
				if !n.Vendor || tier < n.StockLo || tier > n.StockHi {
					continue
				}
				for _, id := range n.Stock {
					it := Items[id]
					if it.Slot == "" || it.Req != tier*2+1 {
						continue
					}
					isJewel := it.Slot == "finger" || it.Slot == "neck" || it.Slot == "back"
					isWeapon := it.Slot == "mainhand" || it.Slot == "offhand"
					if (kind == "weapon" && isWeapon) || (kind == "jewel" && isJewel) || (kind == "armor" && !isJewel && !isWeapon) {
						found = true
					}
				}
			}
			if !found {
				t.Errorf("no vendor sells tier %d %s", tier, kind)
			}
		}
	}
}

func TestEveryClassFindsGearInEveryTier(t *testing.T) {
	Init()
	for tier := 0; tier < tiers; tier++ {
		for _, class := range []string{"warrior", "hunter", "mage"} {
			slots := map[string]bool{}
			for _, id := range append(slices.Clone(gearPool[tier]), gearVariants[tier]...) {
				it := Items[id]
				if len(it.Classes) == 0 || slices.Contains(it.Classes, class) {
					slots[it.Slot] = true
				}
			}
			for _, s := range []string{"head", "chest", "legs", "mainhand", "offhand", "neck", "finger", "back"} {
				if !slots[s] {
					t.Errorf("tier %d: %s has no %s", tier, class, s)
				}
			}
		}
	}
}

func TestMobLoot(t *testing.T) {
	Init()
	for id, m := range Mobs {
		if m.Boss || m.RespawnSec >= 9000 {
			continue
		}
		if len(m.Loot) < 5 {
			t.Errorf("mob %s has only %d loot entries", id, len(m.Loot))
		}
		if m.GearChance < 0.1 {
			t.Errorf("mob %s drops gear too rarely (%.2f)", id, m.GearChance)
		}
	}
	for mob, ids := range namedDrops {
		if Mobs[mob] == nil {
			t.Errorf("named drops for unknown mob %s", mob)
		}
		for _, id := range ids {
			if Items[id] == nil {
				t.Errorf("unknown named item %s", id)
			}
		}
	}
}
