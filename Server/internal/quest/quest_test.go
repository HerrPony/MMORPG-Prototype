package quest

import (
	"testing"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/storage"
)

func TestChainRequiresPrevious(t *testing.T) {
	qs := map[string]*storage.QuestState{}
	if Accept(qs, 5, "mage", "wolf_pelts") {
		t.Fatal("follow-up quest must need the first one")
	}
	Accept(qs, 1, "mage", "wolf_plague")
	qs["wolf_plague"].Status = Done
	if !Accept(qs, 2, "mage", "wolf_pelts") {
		t.Fatal("follow-up should unlock after hand-in")
	}
}

func TestClassQuestOnlyForClass(t *testing.T) {
	qs := map[string]*storage.QuestState{}
	if Accept(qs, 10, "mage", "w_trial") {
		t.Fatal("warrior quest must not be available to a mage")
	}
	if !Accept(qs, 10, "warrior", "w_trial") {
		t.Fatal("warrior quest should be available to a warrior")
	}
}

func TestCollectTracksInventory(t *testing.T) {
	qs := map[string]*storage.QuestState{}
	qs["wolf_plague"] = &storage.QuestState{Status: Done, Counts: []int{6}}
	Accept(qs, 2, "mage", "wolf_pelts")
	RefreshCollect(qs, []storage.ItemStack{{Item: "wolf_pelt", Count: 4}})
	if qs["wolf_pelts"].Status != Complete {
		t.Fatal("4 pelts should complete the quest")
	}
	RefreshCollect(qs, []storage.ItemStack{{Item: "wolf_pelt", Count: 1}})
	if qs["wolf_pelts"].Status != Active {
		t.Fatal("selling pelts should reopen the quest")
	}
}

func TestTimedDeliveryQuest(t *testing.T) {
	qs := map[string]*storage.QuestState{}
	qs["wolf_plague"] = &storage.QuestState{Status: Done, Counts: []int{6}}
	if !Accept(qs, 3, "mage", "e_letter") || qs["e_letter"].Expires == 0 {
		t.Fatal("timed quest should get an expiry")
	}
	if ch := Event(qs, "deliver", "torben"); len(ch) != 1 || qs["e_letter"].Status != Complete {
		t.Fatal("delivering should complete the quest")
	}
	if qs["e_letter"].Expires != 0 {
		t.Fatal("completed quest must not expire")
	}
}

func TestHandInNPC(t *testing.T) {
	gamedata.Init()
	qs := map[string]*storage.QuestState{"wolf_plague": {Status: Done, Counts: []int{6}}}
	Accept(qs, 5, "mage", "e_to_capital")
	if _, m := NPCState(gamedata.NPCByID("alena"), qs, 5, "mage"); m != MarkProgress {
		t.Fatalf("hand-in NPC should show a progress marker, got %q", m)
	}
	Event(qs, "deliver", "alena")
	if _, m := NPCState(gamedata.NPCByID("alena"), qs, 5, "mage"); m != MarkReady {
		t.Fatalf("got %q", m)
	}
}

func TestNPCMarkerPriority(t *testing.T) {
	npc := gamedata.NPCs[0]
	_, m := NPCState(npc, map[string]*storage.QuestState{}, 1, "mage")
	if m != MarkAvailable {
		t.Fatalf("got %q", m)
	}
}
