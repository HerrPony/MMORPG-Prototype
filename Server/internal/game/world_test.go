package game

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"mmorpg/internal/config"
	"mmorpg/internal/gamedata"
	"mmorpg/internal/protocol"
	"mmorpg/internal/storage"
)

type fakeOut struct {
	msgs []map[string]any
}

func (f *fakeOut) Send(v any) {
	b, _ := json.Marshal(v)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	f.msgs = append(f.msgs, m)
}

func (f *fakeOut) last(t string) map[string]any {
	for i := len(f.msgs) - 1; i >= 0; i-- {
		if f.msgs[i]["t"] == t {
			return f.msgs[i]
		}
	}
	return nil
}

func newTestWorld(t *testing.T) *World {
	t.Helper()
	store, err := storage.OpenJSON(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return New(config.Config{TickRate: 50 * time.Millisecond, SnapshotEach: 2, AutosaveEach: time.Hour}, store)
}

func addPlayer(t *testing.T, w *World, name, class string, level int) (*Entity, *fakeOut) {
	t.Helper()
	c := &storage.Character{ID: name, AccountID: "a", Name: name, Class: class, Race: "human", Level: level}
	out := &fakeOut{}
	e := w.enter(c, out)
	return e, out
}

func run(w *World, secs float64) {
	for i := 0; i < int(secs/w.dt); i++ {
		w.tick()
	}
}

func findMob(w *World, id string) *Entity {
	for _, e := range w.ents {
		if e.M != nil && e.M.Def.ID == id && !e.Dead && e.Inst == 0 {
			return e
		}
	}
	return nil
}

func TestEnterAndSnapshot(t *testing.T) {
	w := newTestWorld(t)
	e, out := addPlayer(t, w, "Anna", "mage", 1)
	if out.last("world") == nil {
		t.Fatal("no world message")
	}
	if len(e.P.Char.Skills) != 1 {
		t.Fatalf("a new mage knows only the basic attack, got %v", e.P.Char.Skills)
	}
	run(w, 1)
	if s := out.last("s"); s == nil {
		t.Fatal("no snapshot")
	}
}

func TestKillWolfGivesXPAndPersonalLoot(t *testing.T) {
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "mage", 1)
	wolf := findMob(w, "wolf")
	e.X, e.Z = clearShotPos(t, wolf, 8)
	e.Y = gamedata.Height(e.X, e.Z)
	e.P.Target = wolf.ID
	for i := 0; i < 400 && !wolf.Dead; i++ {
		w.tryCast(e, "m_bolt")
		run(w, 1)
		if e.Dead {
			e.Dead = false
			e.HP = e.MaxHP
		}
	}
	if !wolf.Dead || e.P.Char.XP == 0 {
		t.Fatalf("wolf dead=%v xp=%d", wolf.Dead, e.P.Char.XP)
	}
}

func TestPartyAndDungeon(t *testing.T) {
	w := newTestWorld(t)
	a, _ := addPlayer(t, w, "Anna", "warrior", 18)
	b, outB := addPlayer(t, w, "Bert", "mage", 18)
	w.onInvite(a, protocol.InviteMsg{ID: b.ID})
	w.onInviteReply(b, true)
	if a.P.PartyID == 0 || a.P.PartyID != b.P.PartyID {
		t.Fatal("party not formed")
	}
	if outB.last("party") == nil {
		t.Fatal("no party message")
	}
	for _, n := range []string{"C1", "D1", "E1"} {
		p, _ := addPlayer(t, w, n, "hunter", 18)
		w.onInvite(a, protocol.InviteMsg{ID: p.ID})
		w.onInviteReply(p, true)
	}
	x, _ := addPlayer(t, w, "Extra", "hunter", 18)
	w.onInvite(a, protocol.InviteMsg{ID: x.ID})
	w.onInviteReply(x, true)
	if x.P.PartyID != 0 {
		t.Fatal("a party has at most 5 members")
	}
	// enter the dungeon together
	for _, id := range w.parties[a.P.PartyID].Members {
		m := w.players[id]
		w.enterDungeon(m, "mine")
		if m.Inst == 0 {
			t.Fatalf("%s did not enter", m.Name)
		}
	}
	if a.Inst != b.Inst {
		t.Fatal("party members must share one instance")
	}
	if len(w.insts) != 1 {
		t.Fatalf("expected one instance, got %d", len(w.insts))
	}
	// the foreman fights back
	var boss *Entity
	for _, e := range w.ents {
		if e.M != nil && e.M.Def.ID == "foreman" {
			boss = e
		}
	}
	if boss == nil {
		t.Fatal("no foreman")
	}
	a.X, a.Z = boss.X, boss.Z+6
	a.MaxHP, a.HP = 100000, 100000
	w.damage(a, boss, 1, false)
	run(w, 20)
	if boss.M.Target == 0 && !a.Dead {
		t.Fatal("boss should be fighting")
	}
	boss.HP = boss.MaxHP / 3 // trigger the add summon
	run(w, 3)
	adds := 0
	for _, e := range w.ents {
		if e.M != nil && e.M.Def.ID == "kobold_add" {
			adds++
		}
	}
	if adds == 0 {
		t.Fatal("boss should have summoned helpers")
	}
	// dying in the dungeon revives at the entrance
	a.Dead = true
	w.onRevive(a)
	if a.Inst == 0 {
		t.Fatal("revive must stay inside the dungeon")
	}
	w.exitDungeon(a)
	if a.Inst != 0 {
		t.Fatal("exit failed")
	}
}

func TestGriffinFlight(t *testing.T) {
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "mage", 6)
	kira := w.npcByID("kira")
	e.X, e.Z = kira.X, kira.Z+2
	w.onInteract(e, kira.ID) // unlocks Eichenhain
	bran := w.npcByID("bran")
	e.P.Char.FlightPoints = append(e.P.Char.FlightPoints, "sonnenburg")
	e.P.Char.Copper = 1000
	w.onFlight(e, protocol.FlightMsg{NPC: kira.ID, Dest: "sonnenburg"})
	if e.P.Flight == nil {
		t.Fatal("flight did not start")
	}
	if e.P.Char.Copper >= 1000 {
		t.Fatal("flight must cost money")
	}
	run(w, 15)
	if e.P.Flight != nil {
		t.Fatal("flight should be over")
	}
	px, pz := gamedata.PadPos(gamedata.TownByID("sonnenburg"))
	if dist(e, &Entity{X: px, Z: pz}) > 3 {
		t.Fatalf("landed at %.0f,%.0f, want %.0f,%.0f", e.X, e.Z, px, pz)
	}
	_ = bran
}

func TestTimedQuestExpires(t *testing.T) {
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "mage", 3)
	e.P.Char.Quests["wolf_plague"] = &storage.QuestState{Status: "done", Counts: []int{6}}
	e.X, e.Z = 6, 9
	w.onAccept(e, "e_letter")
	st := e.P.Char.Quests["e_letter"]
	if st == nil || len(e.P.Char.Inventory) != 1 {
		t.Fatal("quest not accepted or no letter given")
	}
	st.Expires = time.Now().Unix() - 1
	w.tickQuestTimers()
	if e.P.Char.Quests["e_letter"] != nil || len(e.P.Char.Inventory) != 0 {
		t.Fatal("expired quest must be removed together with its item")
	}
}

func TestWarriorStanceAndTaunt(t *testing.T) {
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "warrior", 12)
	e.P.Char.Skills = append(e.P.Char.Skills, "w_stance_def", "w_taunt")
	w.tryCast(e, "w_stance_def")
	if e.P.Stance != "defensive" || e.takenMult(w.now) >= 1 {
		t.Fatal("defensive stance should reduce damage taken")
	}
	wolf := findMob(w, "wolf")
	e.X, e.Z = clearShotPos(t, wolf, 5)
	e.P.Target = wolf.ID
	e.P.Res = 50
	w.now += 5
	w.tryCast(e, "w_taunt")
	if wolf.M.Target != e.ID {
		t.Fatal("taunt should force the mob onto the warrior")
	}
}

func TestMageElementalAndHunterPet(t *testing.T) {
	w := newTestWorld(t)
	m, _ := addPlayer(t, w, "Mia", "mage", 12)
	m.P.Char.Skills = append(m.P.Char.Skills, "m_elemental")
	m.P.Res = 100
	w.tryCast(m, "m_elemental")
	run(w, 3)
	if m.P.ElemEnt == 0 || w.ents[m.P.ElemEnt] == nil {
		t.Fatal("elemental was not summoned")
	}
	run(w, 50)
	if w.ents[m.P.ElemEnt] != nil {
		t.Fatal("elemental should have expired")
	}

	h, _ := addPlayer(t, w, "Hans", "hunter", 12)
	h.P.Char.Skills = append(h.P.Char.Skills, "h_tame", "h_callpet")
	wolf := findMob(w, "wolf")
	h.X, h.Z = clearShotPos(t, wolf, 6)
	h.P.Target = wolf.ID
	w.tryCast(h, "h_tame")
	run(w, 7)
	if h.P.Char.Pet == nil || h.P.PetEnt == 0 {
		t.Fatalf("taming failed: pet=%v ent=%d", h.P.Char.Pet, h.P.PetEnt)
	}
	pet := w.ents[h.P.PetEnt]
	pet.X, pet.Z = h.X+30, h.Z
	run(w, 2)
	if dist(pet, h) > 25 {
		t.Fatal("pet should follow its owner")
	}
}

func TestObjectsGiveQuestItems(t *testing.T) {
	w := newTestWorld(t)
	e, _ := addPlayer(t, w, "Anna", "mage", 2)
	var pumpkin *Entity
	for _, o := range w.ents {
		if o.O != nil && o.O.Def.ID == "pumpkin_patch" {
			pumpkin = o
			break
		}
	}
	e.X, e.Z = pumpkin.X, pumpkin.Z+1
	w.useObject(e, pumpkin)
	if len(e.P.Char.Inventory) != 0 {
		t.Fatal("object must not give items without the quest")
	}
	e.P.Char.Quests["e_pumpkin"] = &storage.QuestState{Status: "active", Counts: []int{0}}
	w.useObject(e, pumpkin)
	if len(e.P.Char.Inventory) != 1 || !pumpkin.O.Gone {
		t.Fatal("object should give the quest item and disappear")
	}
	run(w, 50)
	if pumpkin.O.Gone {
		t.Fatal("object should respawn")
	}
}

func (w *World) npcByID(id string) *Entity {
	for _, e := range w.ents {
		if e.N != nil && e.N.ID == id {
			return e
		}
	}
	return nil
}

func TestMobsRespectObstaclesAndSpacing(t *testing.T) {
	sp := gamedata.AllSpawns()
	tooClose := 0
	for i, a := range sp {
		for _, b := range sp[i+1:] {
			if math.Hypot(a.X-b.X, a.Z-b.Z) < 3.9 {
				tooClose++
			}
		}
	}
	if tooClose > 0 {
		t.Fatalf("%d spawn pairs closer than 3.9 m", tooClose)
	}
	// a body pushed into a building ends up outside of it
	for _, b := range gamedata.BuildScenery().House {
		x, z := gamedata.Collide(b.X, b.Z, 0.4)
		if math.Abs(x-b.X) < b.W/2 && math.Abs(z-b.Z) < b.D/2 {
			t.Fatalf("still inside wall at %.1f,%.1f", b.X, b.Z)
		}
	}
	// no line of sight through the cabin wall (north wall at z=-4)
	if !gamedata.LineBlocked(0, -8, 0, 0, false) {
		t.Fatal("shot through the house wall")
	}
}
