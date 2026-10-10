// Package loot rolls loot tables.
package loot

import (
	"math/rand"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/quest"
	"mmorpg/internal/storage"
)

type Bag struct {
	Items  []storage.ItemStack
	Copper int
}

func (b *Bag) Empty() bool { return b == nil || (len(b.Items) == 0 && b.Copper == 0) }

// Looter describes who the loot is rolled for (loot is personal).
type Looter struct {
	Class  string
	Quests map[string]*storage.QuestState
}

// Roll creates the loot of a killed mob for one player.
func Roll(m *gamedata.Mob, mobLevel int, who Looter, rng *rand.Rand) *Bag {
	b := &Bag{}
	if rng.Float64() < m.GearChance {
		b.Items = append(b.Items, storage.ItemStack{Item: gamedata.RandomGear(mobLevel, who.Class, rng), Count: 1})
	}
	if m.Boss {
		n := 1
		if m.ID == "kobold_king" {
			n = 2
		}
		for i := 0; i < n; i++ {
			b.Items = append(b.Items, storage.ItemStack{Item: gamedata.DungeonGear(who.Class, rng), Count: 1})
		}
		b.Items = append(b.Items, storage.ItemStack{Item: "potion_hp3", Count: 2})
	}
	if m.CopperMax > 0 {
		b.Copper = m.CopperMin + rng.Intn(m.CopperMax-m.CopperMin+1)
	}
	for _, e := range m.Loot {
		if e.Quest != "" && !quest.WantsItem(who.Quests, e.Quest) {
			continue
		}
		if rng.Float64() < e.Chance {
			n := e.Min
			if e.Max > e.Min {
				n += rng.Intn(e.Max - e.Min + 1)
			}
			b.Items = append(b.Items, storage.ItemStack{Item: e.Item, Count: n})
		}
	}
	return b
}

// RemoveFromInventory removes count of item (across stacks) and returns the new inventory.
func RemoveFromInventory(inv []storage.ItemStack, item string, count int) []storage.ItemStack {
	out := make([]storage.ItemStack, 0, len(inv))
	for _, s := range inv {
		if s.Item == item && count > 0 {
			take := min(count, s.Count)
			s.Count -= take
			count -= take
		}
		if s.Count > 0 {
			out = append(out, s)
		}
	}
	return out
}

// CountItem returns how many of an item the inventory holds.
func CountItem(inv []storage.ItemStack, item string) int {
	n := 0
	for _, s := range inv {
		if s.Item == item {
			n += s.Count
		}
	}
	return n
}

// AddToInventory merges stacks into inv (respecting stack sizes and slot limit).
// It returns the new inventory and the amount that did not fit.
func AddToInventory(inv []storage.ItemStack, item string, count, maxSlots int) ([]storage.ItemStack, int) {
	def := gamedata.Items[item]
	stack := 1
	if def != nil {
		stack = def.Stack
	}
	for i := range inv {
		if inv[i].Item == item && inv[i].Count < stack {
			add := min(count, stack-inv[i].Count)
			inv[i].Count += add
			count -= add
		}
	}
	for count > 0 && len(inv) < maxSlots {
		add := min(count, stack)
		inv = append(inv, storage.ItemStack{Item: item, Count: add})
		count -= add
	}
	return inv, count
}
