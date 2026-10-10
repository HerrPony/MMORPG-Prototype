package loot

import (
	"math/rand"
	"testing"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/storage"
)

func TestRollWolfGivesCopper(t *testing.T) {
	b := Roll(gamedata.Mobs["wolf"], 1, Looter{Class: "mage"}, rand.New(rand.NewSource(1)))
	if b.Copper < 1 {
		t.Fatal("wolf should drop copper")
	}
}

func TestQuestOnlyDrops(t *testing.T) {
	gamedata.Init()
	m := gamedata.Mobs["icewolf"]
	rng := rand.New(rand.NewSource(3))
	without := 0
	for i := 0; i < 200; i++ {
		for _, s := range Roll(m, 14, Looter{Class: "mage"}, rng).Items {
			if s.Item == "icewolf_pelt" {
				without++
			}
		}
	}
	if without != 0 {
		t.Fatal("quest items must not drop without the quest")
	}
	with := 0
	qs := map[string]*storage.QuestState{"f_wolves": {Status: "active", Counts: []int{0, 0}}}
	for i := 0; i < 200; i++ {
		for _, s := range Roll(m, 14, Looter{Class: "mage", Quests: qs}, rng).Items {
			if s.Item == "icewolf_pelt" {
				with++
			}
		}
	}
	if with == 0 {
		t.Fatal("quest item should drop with an active quest")
	}
}

func TestBossDropsClassGear(t *testing.T) {
	gamedata.Init()
	b := Roll(gamedata.Mobs["kobold_king"], 20, Looter{Class: "hunter"}, rand.New(rand.NewSource(5)))
	found := 0
	for _, s := range b.Items {
		it := gamedata.Items[s.Item]
		if it.Slot != "" {
			found++
			if len(it.Classes) > 0 && it.Classes[0] != "hunter" {
				t.Fatalf("boss dropped %s for a hunter", it.ID)
			}
		}
	}
	if found < 2 {
		t.Fatalf("expected 2 epic pieces, got %d", found)
	}
}

func TestAddToInventoryStacks(t *testing.T) {
	var inv []storage.ItemStack
	inv, left := AddToInventory(inv, "wolf_pelt", 15, 4)
	inv, left = AddToInventory(inv, "wolf_pelt", 10, 4)
	if left != 0 || len(inv) != 2 || inv[0].Count != 20 || inv[1].Count != 5 {
		t.Fatalf("inv=%v left=%d", inv, left)
	}
}
