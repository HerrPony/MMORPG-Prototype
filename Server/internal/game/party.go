package game

import (
	"slices"
	"strings"

	"mmorpg/internal/protocol"
)

// Party is a group of up to five players.
type Party struct {
	ID      int
	Leader  int
	Members []int // entity ids
}

func (w *World) partyOf(e *Entity) *Party {
	if e == nil || e.P == nil || e.P.PartyID == 0 {
		return nil
	}
	return w.parties[e.P.PartyID]
}

// partyStates is the per-tick info shown in the party frames.
func (w *World) partyStates(p *Party) []protocol.M {
	out := make([]protocol.M, 0, len(p.Members))
	for _, id := range p.Members {
		m := w.players[id]
		if m == nil {
			continue
		}
		out = append(out, protocol.M{"id": id, "hp": m.HP, "mh": m.MaxHP, "res": int(m.P.Res), "rm": m.P.ResMax,
			"d": m.Dead, "x": round2(m.X), "z": round2(m.Z), "i": m.Inst, "l": m.Level})
	}
	return out
}

// sendParty tells every member about the current roster.
func (w *World) sendParty(p *Party) {
	list := make([]protocol.M, 0, len(p.Members))
	for _, id := range p.Members {
		if m := w.players[id]; m != nil {
			list = append(list, protocol.M{"id": id, "n": m.Name, "c": m.P.Char.Class, "l": m.Level})
		}
	}
	for _, id := range p.Members {
		if m := w.players[id]; m != nil {
			w.send(m, protocol.M{"t": "party", "id": p.ID, "leader": p.Leader, "members": list})
		}
	}
}

func (w *World) partyChat(from *Entity, text string) {
	p := w.partyOf(from)
	if p == nil {
		w.errMsg(from, "Du bist in keiner Gruppe.")
		return
	}
	for _, id := range p.Members {
		if m := w.players[id]; m != nil {
			w.send(m, protocol.M{"t": "chat", "ch": "party", "from": from.Name, "id": from.ID, "text": text})
		}
	}
}

func (w *World) onInvite(e *Entity, m protocol.InviteMsg) {
	var t *Entity
	if m.ID != 0 {
		t = w.players[m.ID]
	} else {
		for _, p := range w.players {
			if strings.EqualFold(p.Name, strings.TrimSpace(m.Name)) {
				t = p
			}
		}
	}
	switch {
	case t == nil || t == e:
		w.errMsg(e, "Dieser Spieler ist nicht erreichbar.")
		return
	case t.P.PartyID != 0 && t.P.PartyID == e.P.PartyID:
		w.errMsg(e, t.Name+" ist bereits in deiner Gruppe.")
		return
	case t.P.PartyID != 0:
		w.errMsg(e, t.Name+" ist schon in einer Gruppe.")
		return
	}
	if p := w.partyOf(e); p != nil {
		if p.Leader != e.ID {
			w.errMsg(e, "Nur der Gruppenleiter kann einladen.")
			return
		}
		if len(p.Members) >= MaxParty {
			w.errMsg(e, "Deine Gruppe ist voll (maximal 5 Spieler).")
			return
		}
	}
	t.P.Invite = &Invite{From: e.ID, Expire: w.now + 30}
	w.send(t, protocol.M{"t": "partyInvite", "from": e.Name, "id": e.ID})
	w.sysChat(e, "sys", "Du hast "+t.Name+" in die Gruppe eingeladen.")
}

func (w *World) onInviteReply(e *Entity, accept bool) {
	inv := e.P.Invite
	e.P.Invite = nil
	if inv == nil || w.now > inv.Expire {
		return
	}
	from := w.players[inv.From]
	if from == nil {
		return
	}
	if !accept {
		w.sysChat(from, "sys", e.Name+" hat die Einladung abgelehnt.")
		return
	}
	p := w.partyOf(from)
	if p == nil {
		p = &Party{ID: w.nextParty, Leader: from.ID, Members: []int{from.ID}}
		w.nextParty++
		w.parties[p.ID] = p
		from.P.PartyID = p.ID
	}
	if len(p.Members) >= MaxParty {
		w.errMsg(e, "Die Gruppe ist voll.")
		return
	}
	p.Members = append(p.Members, e.ID)
	e.P.PartyID = p.ID
	for _, id := range p.Members {
		if m := w.players[id]; m != nil {
			w.sysChat(m, "sys", e.Name+" ist der Gruppe beigetreten.")
		}
	}
	w.sendParty(p)
}

// leaveParty removes e from its party (announce = tell the others).
func (w *World) leaveParty(e *Entity, announce bool) {
	p := w.partyOf(e)
	if p == nil {
		return
	}
	p.Members = slices.DeleteFunc(p.Members, func(id int) bool { return id == e.ID })
	e.P.PartyID = 0
	w.send(e, protocol.M{"t": "party", "id": 0, "leader": 0, "members": []any{}})
	if len(p.Members) <= 1 {
		for _, id := range p.Members {
			if m := w.players[id]; m != nil {
				m.P.PartyID = 0
				w.send(m, protocol.M{"t": "party", "id": 0, "leader": 0, "members": []any{}})
				w.sysChat(m, "sys", "Die Gruppe wurde aufgelöst.")
			}
		}
		delete(w.parties, p.ID)
		return
	}
	if p.Leader == e.ID {
		p.Leader = p.Members[0]
	}
	for _, id := range p.Members {
		if m := w.players[id]; m != nil {
			w.sysChat(m, "sys", e.Name+" hat die Gruppe verlassen.")
		}
	}
	w.sendParty(p)
}

func (w *World) onKick(e *Entity, id int) {
	p := w.partyOf(e)
	if p == nil || p.Leader != e.ID || id == e.ID {
		return
	}
	if t := w.players[id]; t != nil && t.P.PartyID == p.ID {
		w.sysChat(t, "sys", "Du wurdest aus der Gruppe geworfen.")
		w.leaveParty(t, true)
	}
}

func (w *World) onPromote(e *Entity, id int) {
	p := w.partyOf(e)
	if p == nil || p.Leader != e.ID || !slices.Contains(p.Members, id) {
		return
	}
	p.Leader = id
	w.sendParty(p)
}

func (w *World) tickParties() {
	for _, pe := range w.players {
		if pe.P.Invite != nil && w.now > pe.P.Invite.Expire {
			pe.P.Invite = nil
		}
	}
}
