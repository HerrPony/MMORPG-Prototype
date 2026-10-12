class_name WinSpells
extends GameWindow
## Spell book: all skills of the class, learned and still locked ones.

var list: VBoxContainer

func _init() -> void:
	super._init("Zauberbuch", Vector2(460, 100))
	win_id = "spells"
	default_pos = Vector2(300, 90)
	var sc := ScrollContainer.new()
	sc.custom_minimum_size = Vector2(440, 480)
	sc.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	list = UI.vbox(4)
	list.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	sc.add_child(list)
	body.add_child(sc)

func refresh() -> void:
	if not visible:
		return
	UI.clear(list)
	var known := {}
	for id in S.skill_ids:
		known[id] = true
	for id in S.cls.get("skills", []):
		var d := S.skill(id)
		if d.is_empty():
			continue
		var card := PanelContainer.new()
		card.add_theme_stylebox_override("panel", UI.box(Color(1, 1, 1, 0.04 if known.has(id) else 0.015), UI.EDGE_SOFT, 1, 5, 6))
		card.modulate.a = 1.0 if known.has(id) else 0.55
		var h := UI.hbox(10)
		var ic := Slot.new(42)
		ic.icon = d.icon
		ic.border = Color(0.84, 0.7, 0.38, 0.6)
		ic.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(ic)
		var v := UI.vbox(1)
		v.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		v.add_child(UI.label(d.name, 15, UI.GOLD if known.has(id) else UI.MUTED))
		var meta: Array = []
		if d.get("cost", 0) > 0: meta.append("%d %s" % [int(d.cost), ItemTip.RES_NAME.get(S.cls.get("resource", ""), "")])
		if d.get("range", 0) > 0: meta.append("%d m" % int(d.range))
		if d.get("cast", 0) > 0: meta.append("%s Sek. Zauberzeit" % str(d.cast))
		if d.get("cooldown", 0) > 0: meta.append("%s Sek. Abklingzeit" % str(d.cooldown))
		v.add_child(UI.label(" · ".join(meta), 11, UI.MUTED))
		var ds := UI.label(d.get("desc", ""), 12, UI.TEXT)
		ds.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		ds.custom_minimum_size.x = 300
		v.add_child(ds)
		h.add_child(v)
		if not known.has(id):
			var st := ("Klassenquest ab Stufe %d" % int(d.level)) if d.get("quest", false) else ("Stufe %d · beim Klassenlehrer" % int(d.level))
			var sl := UI.label(st, 11, Color("c8a050"))
			sl.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
			sl.custom_minimum_size.x = 100
			h.add_child(sl)
		card.add_child(h)
		list.add_child(card)
	reset_size()
