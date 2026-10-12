class_name WinQuests
extends GameWindow
## Quest log: list of active quests on the left, details on the right.

var list: VBoxContainer
var detail: VBoxContainer
var selected := ""

func _init() -> void:
	super._init("Questlog", Vector2(640, 420))
	win_id = "quests"
	default_pos = Vector2(260, 110)
	var h := UI.hbox(10)
	h.size_flags_vertical = Control.SIZE_EXPAND_FILL
	body.add_child(h)
	var ls := ScrollContainer.new()
	ls.custom_minimum_size = Vector2(210, 400)
	ls.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	list = UI.vbox(3)
	list.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	ls.add_child(list)
	h.add_child(ls)
	var ds := ScrollContainer.new()
	ds.custom_minimum_size = Vector2(400, 400)
	ds.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	ds.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	detail = UI.vbox(6)
	detail.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	ds.add_child(detail)
	h.add_child(ds)

static func reward_bb(q: Dictionary) -> String:
	var r: Dictionary = q.get("reward", {})
	var parts: Array = []
	if int(r.get("xp", 0)) > 0:
		parts.append("🌟 %d Erfahrung" % int(r.xp))
	if int(r.get("copper", 0)) > 0:
		parts.append(UI.copper_text(int(r.copper)))
	for it in r.get("items", []):
		var d := S.item(it)
		parts.append("[color=%s]%s %s[/color]" % [ItemTip.col(d.get("rarity", "common")), d.get("icon", ""), d.get("name", it)])
	return "   ".join(parts)

static func objectives_bb(q: Dictionary, counts) -> String:
	var out := ""
	for i in q.objectives.size():
		var o: Dictionary = q.objectives[i]
		if counts != null:
			var n: int = int(counts[i]) if i < counts.size() else 0
			out += "• %s: %d/%d\n" % [o.text, n, int(o.count)]
		else:
			out += "• %s (%d)\n" % [o.text, int(o.count)]
	return out.strip_edges()

func refresh() -> void:
	if not visible:
		return
	UI.clear(list)
	UI.clear(detail)
	if S.quests.is_empty():
		detail.add_child(UI.label("Du hast keine aktiven Quests.", 14, UI.MUTED))
		return
	var found := false
	for q in S.quests:
		if q.id == selected:
			found = true
	if not found:
		selected = S.quests[0].id
	for q in S.quests:
		var d := S.quest_def(q.id)
		var qid: String = q.id
		var card := UI.card(func(): selected = qid; refresh())
		UI.set_selected(card, qid == selected)
		var l := UI.label(d.title, 13, UI.GOLD if q.status != "complete" else Color("7ee27e"))
		l.mouse_filter = Control.MOUSE_FILTER_IGNORE
		l.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		l.custom_minimum_size.x = 170
		card.add_child(l)
		list.add_child(card)
	var q: Dictionary = {}
	for x in S.quests:
		if x.id == selected:
			q = x
	var d := S.quest_def(selected)
	var name: String = S.self_info.get("name", "")
	detail.add_child(UI.label(d.title + (("   [" + d.tag + "]") if d.get("tag", "") != "" else ""), 18, UI.GOLD))
	var meta := "Stufe %d" % int(d.level)
	if float(q.get("end", 0)) > 0.0:
		meta += " · Zeit übrig: %d Sek." % maxi(0, ceili(float(q.end) - S.now()))
	detail.add_child(UI.label(meta, 12, UI.MUTED))
	var txt := UI.rich(str(d.complete if q.status == "complete" else d.text).replace("%n", name), 13)
	txt.custom_minimum_size.x = 380
	detail.add_child(txt)
	detail.add_child(UI.label("Ziele", 14, UI.GOLD))
	var ob := UI.rich(objectives_bb(d, q.get("counts", [])), 13)
	ob.custom_minimum_size.x = 380
	detail.add_child(ob)
	detail.add_child(UI.label("Belohnung", 14, UI.GOLD))
	var rw := UI.rich(reward_bb(d), 13)
	rw.custom_minimum_size.x = 380
	detail.add_child(rw)
	detail.add_child(UI.button("Quest abbrechen", func(): Net.send("abandonQuest", {"quest": selected})))
