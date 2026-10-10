// Package quest contains pure quest-state logic.
package quest

import (
	"time"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/storage"
)

const (
	Active   = "active"
	Complete = "complete"
	Done     = "done"
)

// Marker states shown above NPCs / in dialogs.
const (
	MarkNone      = ""
	MarkAvailable = "available" // yellow !
	MarkProgress  = "progress"  // grey ?
	MarkReady     = "ready"     // yellow ?
)

// TurnInNPC is the NPC a finished quest is handed in to.
func TurnInNPC(q *gamedata.Quest) string {
	if q.Hand != "" {
		return q.Hand
	}
	return q.Giver
}

// Available reports whether the player may accept the quest.
func Available(quests map[string]*storage.QuestState, level int, class, id string) bool {
	q := gamedata.Quests[id]
	if q == nil || level < q.MinLevel {
		return false
	}
	if q.Class != "" && q.Class != class {
		return false
	}
	if _, ok := quests[id]; ok {
		return false
	}
	if q.Requires != "" {
		if r, ok := quests[q.Requires]; !ok || r.Status != Done {
			return false
		}
	}
	return true
}

// Accept starts a quest. Returns false if not allowed.
func Accept(quests map[string]*storage.QuestState, level int, class, id string) bool {
	if !Available(quests, level, class, id) {
		return false
	}
	q := gamedata.Quests[id]
	st := &storage.QuestState{Status: Active, Counts: make([]int, len(q.Objectives))}
	if q.TimeLimit > 0 {
		st.Expires = time.Now().Unix() + int64(q.TimeLimit)
	}
	if len(q.Objectives) == 0 {
		st.Status = Complete
	}
	quests[id] = st
	return true
}

// Event registers a kill/deliver/talk/explore/tame event and returns the ids of changed quests.
func Event(quests map[string]*storage.QuestState, kind, target string) (changed []string) {
	for id, st := range quests {
		if st.Status != Active {
			continue
		}
		q := gamedata.Quests[id]
		if q == nil {
			continue
		}
		hit := false
		for i, o := range q.Objectives {
			if o.Kind == kind && o.Target == target && st.Counts[i] < o.Count {
				st.Counts[i]++
				hit = true
			}
		}
		if hit {
			if Finished(q, st) {
				st.Status = Complete
				st.Expires = 0
			}
			changed = append(changed, id)
		}
	}
	return
}

func Finished(q *gamedata.Quest, st *storage.QuestState) bool {
	for i, o := range q.Objectives {
		if st.Counts[i] < o.Count {
			return false
		}
	}
	return true
}

// WantsItem reports whether an active quest still needs the given item (for quest-only drops).
func WantsItem(quests map[string]*storage.QuestState, questID string) bool {
	st, ok := quests[questID]
	return ok && st.Status == Active
}

// RefreshCollect recounts "collect" objectives from the inventory and returns changed quests.
func RefreshCollect(quests map[string]*storage.QuestState, inv []storage.ItemStack) (changed []string) {
	for id, st := range quests {
		if st.Status == Done {
			continue
		}
		q := gamedata.Quests[id]
		if q == nil {
			continue
		}
		diff := false
		for i, o := range q.Objectives {
			if o.Kind != "collect" {
				continue
			}
			n := 0
			for _, s := range inv {
				if s.Item == o.Target {
					n += s.Count
				}
			}
			n = min(n, o.Count)
			if st.Counts[i] != n {
				st.Counts[i] = n
				diff = true
			}
		}
		if Finished(q, st) {
			if st.Status != Complete {
				st.Expires = 0
			}
			st.Status = Complete
		} else {
			st.Status = Active
		}
		if diff {
			changed = append(changed, id)
		}
	}
	return
}

// CollectItems lists the items a quest consumes on hand-in.
func CollectItems(q *gamedata.Quest) (items []gamedata.Objective) {
	for _, o := range q.Objectives {
		if o.Kind == "collect" {
			items = append(items, o)
		}
	}
	return
}

// Offer is one quest an NPC can currently do something with.
type Offer struct {
	ID    string
	State string
}

var rank = map[string]int{MarkNone: 0, MarkProgress: 1, MarkAvailable: 2, MarkReady: 3}

// Offers lists every quest an NPC has for the player (new, in progress, ready to hand in).
func Offers(npc *gamedata.NPC, quests map[string]*storage.QuestState, level int, class string) []Offer {
	var out []Offer
	consider := func(id string, giver bool) {
		q := gamedata.Quests[id]
		if q == nil {
			return
		}
		st, ok := quests[id]
		turnIn := TurnInNPC(q) == npc.ID
		switch {
		case !ok:
			if giver && Available(quests, level, class, id) {
				out = append(out, Offer{id, MarkAvailable})
			}
		case st.Status == Active && turnIn:
			out = append(out, Offer{id, MarkProgress})
		case st.Status == Complete && turnIn:
			out = append(out, Offer{id, MarkReady})
		}
	}
	for _, id := range npc.Quests {
		consider(id, true)
	}
	for _, id := range npc.TurnIn {
		consider(id, false)
	}
	return out
}

// NPCState returns the best quest and marker an NPC currently offers a player.
func NPCState(npc *gamedata.NPC, quests map[string]*storage.QuestState, level int, class string) (questID, mark string) {
	for _, o := range Offers(npc, quests, level, class) {
		if rank[o.State] > rank[mark] {
			questID, mark = o.ID, o.State
		}
	}
	return
}
