class_name WinDialog
extends GameWindow
## NPC dialog: quests, class trainer, griffin flights, ferry information. Vendors open their own window.

signal open_vendor(m: Dictionary)

var m: Dictionary = {}
var view := "greet"
var quest_id := ""
var tabs: HBoxContainer
var content: VBoxContainer

func _init() -> void:
	super._init("", Vector2(470, 100))
	win_id = "dialog"
	default_pos = Vector2(420, 120)
	tabs = UI.hbox(4)
	body.add_child(tabs)
	var sc := ScrollContainer.new()
	sc.custom_minimum_size = Vector2(450, 340)
	sc.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	sc.size_flags_vertical = Control.SIZE_EXPAND_FILL
	content = UI.vbox(7)
	content.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	sc.add_child(content)
	body.add_child(sc)
	closed.connect(func(): S.dialog_npc = 0)

func _views() -> Array:
	var v: Array = []
	var quests: Array = m.get("quests", [])
	if quests.size() > 0:
		var ready := false
		for q in quests:
			if q.state == "ready":
				ready = true
		v.append(["quest", "❓ Quests" if ready else "❗ Quests"])
	if m.get("trainer") != null:
		v.append(["trainer", "📚 Training"])
	if m.get("flight") != null:
		v.append(["flight", "🦅 Greifenflug"])
	var fy = m.get("ferry")
	if fy is Dictionary and fy.get("ship", false):
		v.append(["ferry", "⛵ Fähre"])
	return v

func show_dialog(msg: Dictionary) -> void:
	var same: bool = visible and int(m.get("npc", -1)) == int(msg.npc)
	var prev := view if same else ""
	if not same:
		quest_id = ""
	m = msg
	S.dialog_npc = int(msg.npc)
	var vs := _views()
	if bool(m.get("vendor", false)) and vs.is_empty():
		close()
		open_vendor.emit(m)
		return
	title_label.text = m.name + ((" – " + m.title) if m.get("title", "") != "" else "")
	view = "greet"
	for v in vs:
		if v[0] == prev:
			view = prev
	if view == "greet" and vs.size() > 0:
		view = vs[0][0]
	if view == "quest" and (m.get("quests", []) as Array).size() == 1:
		quest_id = m.quests[0].id
	var still := false
	for q in m.get("quests", []):
		if q.id == quest_id:
			still = true
	if not still:
		quest_id = ""
	open()
	render()

func render() -> void:
	if m.is_empty():
		return
	UI.clear(tabs)
	UI.clear(content)
	var vs := _views()
	if vs.size() > 1 or (bool(m.get("vendor", false)) and vs.size() > 0):
		if m.get("vendor", false):
			tabs.add_child(UI.button("💰 Handel", func(): open_vendor.emit(m)))
		for v in vs:
			var b := UI.button(v[1], Callable())
			var vid: String = v[0]
			b.toggle_mode = true
			b.button_pressed = vid == view
			b.pressed.connect(func(): view = vid; render())
			tabs.add_child(b)
	match view:
		"quest": _render_quest()
		"trainer": _render_trainer()
		"flight": _render_flight()
		"ferry": _render_ferry()
		_:
			content.add_child(_text(str(m.get("greet", "…"))))
	reset_size()

func _text(t: String, size_px := 14) -> RichTextLabel:
	var r := UI.rich(t, size_px)
	r.custom_minimum_size.x = 430
	return r

func _render_quest() -> void:
	var quests: Array = m.quests
	if quest_id == "":
		content.add_child(_text(str(m.get("greet", ""))))
		for o in quests:
			var q := S.quest_def(o.id)
			var oid: String = o.id
			var card := UI.card(func(): quest_id = oid; render())
			var h := UI.hbox(8)
			var mark := {"available": "❗", "ready": "❓", "progress": "❔"}
			var ic := UI.label(mark.get(o.state, ""), 22)
			ic.mouse_filter = Control.MOUSE_FILTER_IGNORE
			h.add_child(ic)
			var v := UI.vbox(0)
			v.mouse_filter = Control.MOUSE_FILTER_IGNORE
			var t := UI.label(q.title, 15, UI.GOLD)
			t.mouse_filter = Control.MOUSE_FILTER_IGNORE
			v.add_child(t)
			var s := UI.label("Stufe %d%s" % [int(q.level), (" · " + q.tag) if q.get("tag", "") != "" else ""], 11, UI.MUTED)
			s.mouse_filter = Control.MOUSE_FILTER_IGNORE
			v.add_child(s)
			h.add_child(v)
			h.mouse_filter = Control.MOUSE_FILTER_IGNORE
			card.add_child(h)
			content.add_child(card)
		return
	var offer := {}
	for o in quests:
		if o.id == quest_id:
			offer = o
	var q := S.quest_def(quest_id)
	var st := {}
	for x in S.quests:
		if x.id == quest_id:
			st = x
	var name: String = S.self_info.get("name", "")
	content.add_child(UI.label(q.title + (("  [" + q.tag + "]") if q.get("tag", "") != "" else ""), 19, UI.GOLD))
	var meta := "Stufe %d" % int(q.level)
	if int(q.get("timeLimit", 0)) > 0:
		meta += " · Zeitlimit: %d Sek." % int(q.timeLimit)
	content.add_child(UI.label(meta, 12, UI.MUTED))
	match offer.state:
		"available":
			content.add_child(_text(str(q.text).replace("%n", name)))
			content.add_child(UI.label("Aufgabe", 14, UI.GOLD))
			content.add_child(_text(WinQuests.objectives_bb(q, null), 13))
			content.add_child(UI.label("Belohnung", 14, UI.GOLD))
			content.add_child(_text(WinQuests.reward_bb(q), 13))
			content.add_child(UI.button("Quest annehmen", func(): Net.send("acceptQuest", {"quest": q.id}); close(), true))
		"progress":
			content.add_child(_text(str(q.progress).replace("%n", name)))
			content.add_child(UI.label("Aufgabe", 14, UI.GOLD))
			content.add_child(_text(WinQuests.objectives_bb(q, st.get("counts", [])), 13))
		_:
			content.add_child(_text(str(q.complete).replace("%n", name)))
			content.add_child(UI.label("Belohnung", 14, UI.GOLD))
			content.add_child(_text(WinQuests.reward_bb(q), 13))
			content.add_child(UI.button("Quest abschließen", func(): Net.send("turnInQuest", {"quest": q.id}); close(), true))
	if quests.size() > 1:
		content.add_child(UI.button("‹ Alle Quests", func(): quest_id = ""; render()))

func _render_ferry() -> void:
	var f: Dictionary = m.ferry
	content.add_child(_text("[b]%s[/b]\n\nDie Fähre pendelt im festen Takt zwischen Sturmhafen und Lotushafen und liegt dabei längsseits am Steg. Geh einfach über den Steg auf das Schiff, bewege dich frei an Deck und steige im Zielhafen wieder über den Steg aus." % f.status))

func _render_flight() -> void:
	var f: Dictionary = m.flight
	if not f.get("ship", false):
		content.add_child(_text("Du bist in [b]%s[/b]. Dein Flugpunkt hier ist freigeschaltet. Dein Geld: %s" % [f.hereName, UI.copper_text(S.copper)]))
	var points: Array = f.get("points", [])
	if points.is_empty():
		content.add_child(_text("[color=#a89a7e]Du kennst noch keinen anderen Flugpunkt. Sprich mit den Greifenmeistern in anderen Städten, um sie freizuschalten. Greifen fliegen nur auf demselben Kontinent.[/color]", 13))
	for p in points:
		var pid: String = p.id
		var price := int(p.price)
		var card := UI.card(func():
			if S.copper >= price:
				Net.send("flight", {"npc": int(m.npc), "dest": pid})
				close()
			else:
				Net._dispatch("err", {"text": "Du hast nicht genug Geld."}))
		var h := UI.hbox(8)
		var ic := UI.label("⛵" if f.get("ship", false) else "🦅", 22)
		ic.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(ic)
		var nm := UI.label(p.name, 15, UI.TEXT)
		nm.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		nm.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(nm)
		var pr := UI.rich(UI.copper_text(price), 13)
		pr.custom_minimum_size.x = 100
		pr.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(pr)
		h.mouse_filter = Control.MOUSE_FILTER_IGNORE
		card.add_child(h)
		content.add_child(card)

func _render_trainer() -> void:
	var t: Dictionary = m.trainer
	if not t.get("mine", false):
		content.add_child(_text("Ich bilde nur %s aus. Such dir einen Lehrer deiner eigenen Klasse." % t.className))
		return
	content.add_child(_text("Dein Geld: " + UI.copper_text(S.copper)))
	for s in t.get("skills", []):
		var d := S.skill(s.id)
		var card := PanelContainer.new()
		card.add_theme_stylebox_override("panel", UI.box(Color(1, 1, 1, 0.04), UI.EDGE_SOFT, 1, 5, 6))
		if s.state == "locked":
			card.modulate.a = 0.55
		var h := UI.hbox(8)
		var ic := Slot.new(40)
		ic.icon = d.icon
		ic.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(ic)
		var v := UI.vbox(0)
		v.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		v.add_child(UI.label(d.name, 15, UI.GOLD))
		var ds := UI.label("Stufe %d · %s" % [int(s.level), d.desc], 11, UI.MUTED)
		ds.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		ds.custom_minimum_size.x = 230
		v.add_child(ds)
		h.add_child(v)
		if s.state == "known":
			h.add_child(UI.label("Gelernt", 13, Color("7ee27e")))
		elif s.state == "locked":
			h.add_child(UI.label("Ab Stufe %d" % int(s.level), 12, Color("c8a050")))
		else:
			var col := UI.vbox(2)
			var pr := UI.rich(UI.copper_text(int(s.price)), 12)
			pr.custom_minimum_size.x = 80
			col.add_child(pr)
			var b := UI.button("Lernen", Callable())
			b.disabled = S.copper < int(s.price)
			var sid: String = s.id
			b.pressed.connect(func(): Net.send("train", {"npc": int(m.npc), "skill": sid}))
			col.add_child(b)
			h.add_child(col)
		card.add_child(h)
		content.add_child(card)
	var has_quest_skill := false
	for id in S.cls.get("skills", []):
		if S.skill(id).get("quest", false):
			has_quest_skill = true
	if has_quest_skill:
		content.add_child(UI.label("Besondere Fertigkeiten (Stufe 10) erlernst du durch deine Klassenquest in Sonnenburg.", 12, UI.MUTED))
