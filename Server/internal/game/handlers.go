package game

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode/utf8"

	"mmorpg/internal/protocol"
)

func itoa(n int) string { return strconv.Itoa(n) }

// Dispatch routes an in-world client message to its handler (runs on the world goroutine).
func (w *World) Dispatch(id int, t string, raw []byte) {
	w.Do(func() {
		e := w.players[id]
		if e == nil {
			return
		}
		w.handle(e, t, raw)
	})
}

func decode[T any](raw []byte) (v T, ok bool) {
	return v, json.Unmarshal(raw, &v) == nil
}

func (w *World) handle(e *Entity, t string, raw []byte) {
	switch t {
	case protocol.CMove:
		if m, ok := decode[protocol.Move](raw); ok {
			w.onMove(e, m)
		}
	case protocol.CTarget:
		if m, ok := decode[protocol.TargetMsg](raw); ok {
			w.onTarget(e, m.ID)
		}
	case protocol.CCast:
		if m, ok := decode[protocol.CastMsg](raw); ok {
			w.tryCast(e, m.Skill)
		}
	case protocol.CAuto:
		if m, ok := decode[protocol.AutoMsg](raw); ok && !m.On {
			e.P.Auto = false
		}
	case protocol.CInteract:
		if m, ok := decode[protocol.InteractMsg](raw); ok {
			w.onInteract(e, m.ID)
		}
	case protocol.CAccept:
		if m, ok := decode[protocol.QuestMsg](raw); ok {
			w.onAccept(e, m.Quest)
		}
	case protocol.CTurnIn:
		if m, ok := decode[protocol.QuestMsg](raw); ok {
			w.onTurnIn(e, m.Quest)
		}
	case protocol.CAbandon:
		if m, ok := decode[protocol.QuestMsg](raw); ok {
			w.onAbandon(e, m.Quest)
		}
	case protocol.CLoot:
		if m, ok := decode[protocol.LootMsg](raw); ok {
			w.onLoot(e, m)
		}
	case protocol.CChat:
		if m, ok := decode[protocol.ChatMsg](raw); ok {
			w.onChat(e, m.Text)
		}
	case protocol.CTrain:
		if m, ok := decode[protocol.TrainMsg](raw); ok {
			w.onTrain(e, m)
		}
	case protocol.CSell:
		if m, ok := decode[protocol.SellMsg](raw); ok {
			w.onSell(e, m)
		}
	case protocol.CBuy:
		if m, ok := decode[protocol.BuyMsg](raw); ok {
			w.onBuy(e, m)
		}
	case protocol.CSellJunk:
		if m, ok := decode[protocol.SellJunkMsg](raw); ok {
			w.onSellJunk(e, m.NPC)
		}
	case protocol.CEquip:
		if m, ok := decode[protocol.EquipMsg](raw); ok {
			w.onEquip(e, m)
		}
	case protocol.CUnequip:
		if m, ok := decode[protocol.UnequipMsg](raw); ok {
			w.onUnequip(e, m.Slot)
		}
	case protocol.CUse:
		if m, ok := decode[protocol.UseMsg](raw); ok {
			w.onUse(e, m.Idx)
		}
	case protocol.CSaveSet:
		if m, ok := decode[protocol.SetMsg](raw); ok {
			w.onSaveSet(e, m.Name)
		}
	case protocol.CEquipSet:
		if m, ok := decode[protocol.SetMsg](raw); ok {
			w.onEquipSet(e, m.Name)
		}
	case protocol.CDelSet:
		if m, ok := decode[protocol.SetMsg](raw); ok {
			w.onDeleteSet(e, m.Name)
		}
	case protocol.CInvite:
		if m, ok := decode[protocol.InviteMsg](raw); ok {
			w.onInvite(e, m)
		}
	case protocol.CInviteRe:
		if m, ok := decode[protocol.InviteReplyMsg](raw); ok {
			w.onInviteReply(e, m.Accept)
		}
	case protocol.CLeaveP:
		w.leaveParty(e, true)
	case protocol.CKick:
		if m, ok := decode[protocol.IDMsg](raw); ok {
			w.onKick(e, m.ID)
		}
	case protocol.CPromote:
		if m, ok := decode[protocol.IDMsg](raw); ok {
			w.onPromote(e, m.ID)
		}
	case protocol.CFlight:
		if m, ok := decode[protocol.FlightMsg](raw); ok {
			w.onFlight(e, m)
		}
	case protocol.CRevive:
		w.onRevive(e)
	}
}

func (w *World) onTarget(e *Entity, id int) {
	p := e.P
	if id == p.Target {
		return
	}
	p.Auto = false
	if id == 0 {
		p.Target = 0
		return
	}
	t := w.ents[id]
	if t == nil || t.Inst != e.Inst || (t.M != nil && t.M.Gone) || (t.O != nil && t.O.Gone) {
		return
	}
	p.Target = id
}

func (w *World) onChat(e *Entity, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	if utf8.RuneCountInString(text) > 200 {
		text = string([]rune(text)[:200])
	}
	if w.cfg.Dev && strings.HasPrefix(text, "/dev ") {
		w.devCommand(e, strings.Fields(text[5:]))
		return
	}
	switch {
	case strings.HasPrefix(text, "/p "):
		w.partyChat(e, text[3:])
	case text == "/verlassen" || text == "/leave":
		w.leaveParty(e, true)
	case strings.HasPrefix(text, "/einladen "):
		w.onInvite(e, protocol.InviteMsg{Name: text[len("/einladen "):]})
	default:
		w.broadcast(protocol.M{"t": "chat", "ch": "say", "from": e.Name, "id": e.ID, "text": text})
	}
}

func countSuffix(n int) string {
	if n > 1 {
		return " x" + itoa(n)
	}
	return ""
}

func formatCopper(c int) string {
	g, s, k := c/10000, (c/100)%100, c%100
	var parts []string
	if g > 0 {
		parts = append(parts, itoa(g)+" Gold")
	}
	if s > 0 {
		parts = append(parts, itoa(s)+" Silber")
	}
	if k > 0 || len(parts) == 0 {
		parts = append(parts, itoa(k)+" Kupfer")
	}
	return strings.Join(parts, " ")
}
