package game

import (
	"math"
	"slices"
	"strings"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/loot"
	"mmorpg/internal/progression"
	"mmorpg/internal/protocol"
	"mmorpg/internal/storage"
)

// Stats are the totals contributed by equipment.
type Stats struct {
	Sta, Str, Agi, Int int
	Armor, Weapon      int
}

func (s Stats) primary(class string) int {
	switch gamedata.PrimaryStat[class] {
	case "str":
		return s.Str
	case "agi":
		return s.Agi
	}
	return s.Int
}

// damageMult is the damage bonus from the primary stat and weapon.
func (s Stats) damageMult(class string) float64 {
	return 1 + 0.01*float64(s.primary(class)) + 0.012*float64(s.Weapon)
}

// armorReduction is the share of physical damage absorbed against an attacker of the given level.
func armorReduction(armor, attackerLevel int) float64 {
	a := float64(armor)
	return math.Min(0.75, a/(a+85+20*float64(attackerLevel)))
}

func (w *World) equipStats(e *Entity) Stats {
	var s Stats
	for _, id := range e.P.Char.Equipment {
		it := gamedata.Items[id]
		if it == nil {
			continue
		}
		s.Sta += it.Stats["sta"]
		s.Str += it.Stats["str"]
		s.Agi += it.Stats["agi"]
		s.Int += it.Stats["int"]
		s.Armor += it.Armor
		s.Weapon += it.Weapon
	}
	return s
}

// recalc recomputes derived values (max HP/resource) from level and equipment.
func (w *World) recalc(e *Entity) {
	p := e.P
	if p.Char.Equipment == nil {
		p.Char.Equipment = map[string]string{}
	}
	p.Stats = w.equipStats(e)
	oldMax := e.MaxHP
	e.MaxHP = progression.MaxHP(p.Class, e.Level) + p.Stats.Sta*5
	if oldMax > 0 && !e.Dead {
		e.HP = max(1, min(e.MaxHP, int(math.Round(float64(e.HP)*float64(e.MaxHP)/float64(oldMax)))))
	}
	if e.HPBoostEnd > w.now { // "Letztes Gefecht"
		e.MaxHP = int(float64(e.MaxHP) * 1.3)
	}
	p.ResMax = progression.MaxResource(p.Class, e.Level)
	if p.Class.Resource == "mana" {
		p.ResMax += p.Stats.Int * 4
	}
	p.Res = math.Min(p.Res, float64(p.ResMax))
}

func (w *World) sendEquipment(e *Entity) {
	w.send(e, protocol.M{"t": "equip", "eq": e.P.Char.Equipment})
	s := e.P.Stats
	w.send(e, protocol.M{"t": "stats", "sta": s.Sta, "str": s.Str, "agi": s.Agi, "int": s.Int, "armor": s.Armor,
		"weapon": s.Weapon, "dmg": int(math.Round((s.damageMult(e.P.Char.Class) - 1) * 100)),
		"red": int(math.Round(armorReduction(s.Armor, e.Level) * 100))})
}

func (w *World) sendSets(e *Entity) {
	w.send(e, protocol.M{"t": "sets", "sets": e.P.Char.Sets})
}

// slotFor maps an item to a free/valid equipment slot key.
func (w *World) slotFor(e *Entity, it *gamedata.Item) string {
	if it.Slot != "finger" {
		return it.Slot
	}
	if e.P.Char.Equipment["finger1"] != "" && e.P.Char.Equipment["finger2"] == "" {
		return "finger2"
	}
	return "finger1"
}

func (w *World) canEquip(e *Entity, it *gamedata.Item) string {
	switch {
	case e.Dead:
		return "Du bist tot."
	case it == nil || it.Slot == "":
		return "Das kann man nicht anlegen."
	case e.Level < it.Req:
		return "Du benötigst Stufe " + itoa(it.Req) + "."
	case len(it.Classes) > 0 && !slices.Contains(it.Classes, e.P.Char.Class):
		return "Das kann deine Klasse nicht tragen."
	}
	return ""
}

// equipFromBag equips the inventory stack at idx into slot ("" = default slot).
func (w *World) equipFromBag(e *Entity, idx int, slot string) bool {
	c := e.P.Char
	if idx < 0 || idx >= len(c.Inventory) {
		return false
	}
	it := gamedata.Items[c.Inventory[idx].Item]
	if msg := w.canEquip(e, it); msg != "" {
		w.errMsg(e, msg)
		return false
	}
	valid := slot == "finger1" || slot == "finger2"
	if slot == "" || (it.Slot == "finger") != valid || (!valid && slot != it.Slot) {
		slot = w.slotFor(e, it)
	}
	old := c.Equipment[slot]
	c.Inventory = loot.RemoveFromInventory(c.Inventory, it.ID, 1)
	c.Equipment[slot] = it.ID
	if old != "" {
		c.Inventory, _ = loot.AddToInventory(c.Inventory, old, 1, MaxInventory)
	}
	return true
}

func (w *World) afterEquipChange(e *Entity) {
	w.recalc(e)
	w.sendInventory(e)
	w.sendEquipment(e)
}

func (w *World) onEquip(e *Entity, m protocol.EquipMsg) {
	if w.equipFromBag(e, m.Idx, m.Slot) {
		w.afterEquipChange(e)
	}
}

func (w *World) onUnequip(e *Entity, slot string) {
	c := e.P.Char
	id := c.Equipment[slot]
	if id == "" || e.Dead {
		return
	}
	inv, left := loot.AddToInventory(c.Inventory, id, 1, MaxInventory)
	if left > 0 {
		w.errMsg(e, "Dein Inventar ist voll.")
		return
	}
	c.Inventory = inv
	delete(c.Equipment, slot)
	w.afterEquipChange(e)
}

func (w *World) onUse(e *Entity, idx int) {
	c := e.P.Char
	if idx < 0 || idx >= len(c.Inventory) || e.Dead {
		return
	}
	it := gamedata.Items[c.Inventory[idx].Item]
	if it == nil || it.Heal == 0 {
		return
	}
	switch {
	case e.Level < it.Req:
		w.errMsg(e, "Du benötigst Stufe "+itoa(it.Req)+".")
		return
	case w.now < e.P.PotCD:
		w.errMsg(e, "Noch nicht bereit.")
		return
	case e.HP >= e.MaxHP:
		w.errMsg(e, "Du hast bereits volle Gesundheit.")
		return
	}
	e.P.PotCD = w.now + 20
	healed := min(it.Heal, e.MaxHP-e.HP)
	e.HP += healed
	c.Inventory = loot.RemoveFromInventory(c.Inventory, it.ID, 1)
	w.sysChat(e, "combat", "Du benutzt ["+it.Name+"] und heilst "+itoa(healed)+" Lebenspunkte.")
	w.broadcast(protocol.M{"t": "fx", "k": "heal", "s": e.ID, "d": e.ID})
	w.sendInventory(e)
}

func (w *World) onBuy(e *Entity, m protocol.BuyMsg) {
	n := w.npcNear(e, m.NPC)
	if n == nil || !n.N.Vendor || !slices.Contains(n.N.Stock, m.Item) {
		return
	}
	it := gamedata.Items[m.Item]
	cnt := max(1, min(m.Count, 20))
	if it.Stack == 1 {
		cnt = 1
	}
	cost := it.Buy * cnt
	if e.P.Char.Copper < cost {
		w.errMsg(e, "Du hast nicht genug Geld.")
		return
	}
	inv, left := loot.AddToInventory(e.P.Char.Inventory, it.ID, cnt, MaxInventory)
	if left > 0 {
		w.errMsg(e, "Dein Inventar ist voll.")
		return
	}
	e.P.Char.Inventory = inv
	e.P.Char.Copper -= cost
	w.sysChat(e, "loot", "Du kaufst ["+it.Name+"]"+countSuffix(cnt)+" für "+formatCopper(cost)+".")
	w.sendInventory(e)
	w.refreshCollect(e)
}

func (w *World) onSellJunk(e *Entity, npc int) {
	n := w.npcNear(e, npc)
	if n == nil || !n.N.Vendor {
		return
	}
	c := e.P.Char
	var keep []storage.ItemStack
	total := 0
	for _, s := range c.Inventory {
		if it := gamedata.Items[s.Item]; it != nil && it.Rarity == "poor" {
			total += it.Sell * s.Count
			continue
		}
		keep = append(keep, s)
	}
	if total == 0 {
		w.errMsg(e, "Du hast keinen Plunder zum Verkaufen.")
		return
	}
	c.Inventory = keep
	c.Copper += total
	w.sysChat(e, "loot", "Du verkaufst Plunder für "+formatCopper(total)+".")
	w.sendInventory(e)
	w.refreshCollect(e)
}

// ---- equipment sets (equipment manager) ----

func (w *World) onSaveSet(e *Entity, name string) {
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > 16 {
		w.errMsg(e, "Ungültiger Setname.")
		return
	}
	c := e.P.Char
	if c.Sets == nil {
		c.Sets = map[string]map[string]string{}
	}
	if _, exists := c.Sets[name]; !exists && len(c.Sets) >= 10 {
		w.errMsg(e, "Maximal 10 Ausrüstungssets.")
		return
	}
	set := map[string]string{}
	for k, v := range c.Equipment {
		set[k] = v
	}
	c.Sets[name] = set
	w.sysChat(e, "sys", "Ausrüstungsset \""+name+"\" gespeichert.")
	w.sendSets(e)
}

func (w *World) onEquipSet(e *Entity, name string) {
	c := e.P.Char
	set, ok := c.Sets[name]
	if !ok || e.Dead {
		return
	}
	changed, missing := false, 0
	for _, slot := range gamedata.EquipSlots {
		id := set[slot]
		if id == "" || c.Equipment[slot] == id {
			continue
		}
		idx := slices.IndexFunc(c.Inventory, func(s storage.ItemStack) bool { return s.Item == id })
		if idx < 0 {
			missing++
			continue
		}
		if w.equipFromBag(e, idx, slot) {
			changed = true
		}
	}
	if missing > 0 {
		w.sysChat(e, "sys", itoa(missing)+" Gegenstände des Sets fehlen in deinen Taschen.")
	}
	if changed {
		w.afterEquipChange(e)
	}
}

func (w *World) onDeleteSet(e *Entity, name string) {
	delete(e.P.Char.Sets, name)
	w.sendSets(e)
}
