class_name ItemTip
extends RefCounted
## BBCode tooltips for items and skills.

const SLOT_LABEL := {
	"head": "Kopf", "neck": "Hals", "shoulder": "Schultern", "back": "Rücken", "chest": "Brust", "wrist": "Handgelenke", "hands": "Hände", "waist": "Taille",
	"legs": "Beine", "feet": "Füße", "finger": "Finger", "finger1": "Finger", "finger2": "Finger", "mainhand": "Waffenhand", "offhand": "Nebenhand",
}
const SLOT_ICON := {"head": "🪖", "neck": "📿", "shoulder": "🦺", "back": "🧥", "chest": "👕", "wrist": "⌚", "hands": "🧤", "waist": "🎗️", "legs": "👖", "feet": "🥾", "finger1": "💍", "finger2": "💍", "mainhand": "⚔️", "offhand": "🛡️"}
const STAT_LABEL := {"str": "Stärke", "agi": "Beweglichkeit", "int": "Intelligenz", "sta": "Ausdauer"}
const RARITY_LABEL := {"poor": "Schlecht", "common": "Gewöhnlich", "uncommon": "Ungewöhnlich", "rare": "Selten", "epic": "Episch"}
const RES_NAME := {"rage": "Wut", "focus": "Fokus", "mana": "Mana"}

static func col(rarity: String) -> String:
	return "#" + UI.rarity_color(rarity).to_html(false)

static func item(id: String, count := 1, mode := "") -> String:
	var it := S.item(id)
	if it.is_empty():
		return ""
	var rar: String = it.get("rarity", "common")
	var h := "[b][color=%s]%s[/color][/b]\n" % [col(rar), it.name]
	if it.get("slot", "") != "":
		h += "%s  [color=#a89a7e]%s[/color]\n" % [SLOT_LABEL.get(it.slot, it.slot), it.get("type", "")]
	else:
		h += "[color=#a89a7e]%s[/color]\n" % RARITY_LABEL.get(rar, "")
	if it.get("armor", 0) > 0:
		h += "%d Rüstung\n" % int(it.armor)
	if it.get("weapon", 0) > 0:
		h += "%d Waffenschaden\n" % int(it.weapon)
	var stats: Dictionary = it.get("stats", {}) if it.get("stats") is Dictionary else {}
	for k in ["str", "agi", "int", "sta"]:
		if stats.get(k, 0) != 0:
			h += "[color=#5bd65b]+%d %s[/color]\n" % [int(stats[k]), STAT_LABEL[k]]
	var classes: Array = it.get("classes", []) if it.get("classes") is Array else []
	if classes.size() > 0:
		var names: Array = []
		for c in classes:
			names.append(S.classes[c].name if S.classes.has(c) else c)
		var ok: bool = classes.has(S.self_info.get("class", ""))
		h += "%sKlasse: %s%s\n" % ["" if ok else "[color=#ff5a4a]", ", ".join(names), "" if ok else "[/color]"]
	if int(it.get("req", 0)) > 0:
		var ok2: bool = S.me.lvl >= int(it.req)
		h += "%sBenötigt Stufe %d%s\n" % ["" if ok2 else "[color=#ff5a4a]", int(it.req), "" if ok2 else "[/color]"]
	if it.get("desc", "") != "":
		h += "[color=#e8d98a]%s[/color]\n" % it.desc
	if mode == "buy":
		h += "\nPreis: " + UI.copper_text(int(it.get("buy", 0)))
	elif int(it.get("sell", 0)) > 0:
		h += "\nVerkaufspreis: " + UI.copper_text(int(it.sell) * maxi(1, count))
	if mode == "bag":
		if S.vendor_npc != 0:
			h += "\n[color=#8a8a7a]Rechtsklick: verkaufen · Doppelklick: anlegen/benutzen[/color]"
		elif it.get("slot", "") != "":
			h += "\n[color=#8a8a7a]Rechtsklick: anlegen[/color]"
		elif it.get("heal", 0) != 0:
			h += "\n[color=#8a8a7a]Rechtsklick: benutzen[/color]"
	return h.strip_edges(false, true)

static func skill(sk: Dictionary) -> String:
	var h := "[b][color=#f0cf7a]%s[/color][/b]\n" % sk.name
	var meta: Array = []
	var res: String = RES_NAME.get(S.cls.get("resource", ""), "")
	if sk.get("cost", 0) > 0:
		meta.append("%d %s" % [int(sk.cost), res])
	if sk.get("range", 0) > 0:
		meta.append("%d m Reichweite" % int(sk.range))
	if sk.get("cast", 0) > 0:
		meta.append("%s Sek. Zauberzeit" % str(sk.cast))
	elif not sk.get("basic", false):
		meta.append("Sofort")
	if sk.get("cooldown", 0) > 0:
		meta.append("%s Sek. Abklingzeit" % str(sk.cooldown))
	h += "[color=#a89a7e]%s[/color]\n%s" % [" · ".join(meta), sk.get("desc", "")]
	return h
