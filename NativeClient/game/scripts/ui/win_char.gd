class_name WinChar
extends GameWindow
## Character sheet: paper doll with equipment slots, stats and the equipment set manager.

const LEFT := ["head", "neck", "shoulder", "back", "chest", "wrist"]
const RIGHT := ["hands", "waist", "legs", "feet", "finger1", "finger2"]
const WEAP := ["mainhand", "offhand"]
const PRIMARY := {"warrior": "str", "hunter": "agi", "mage": "int"}
const RES_NAME := {"rage": "Wut", "focus": "Fokus", "mana": "Mana"}

var content: VBoxContainer
var set_name: LineEdit

func _init() -> void:
	super._init("Charakter", Vector2(430, 100))
	win_id = "char"
	default_pos = Vector2(120, 120)
	content = UI.vbox(8)
	body.add_child(content)

func refresh() -> void:
	if not visible:
		return
	UI.clear(content)
	var me := S.me
	var c := S.cls
	var st := S.stats
	var cls_id: String = S.self_info.get("class", "")
	var head := UI.vbox(0)
	head.add_child(UI.label(S.self_info.get("name", ""), 20, UI.CLASS_COLORS.get(cls_id, Color.WHITE)))
	head.add_child(UI.label("Stufe %d %s %s" % [me.lvl, S.races.get(S.self_info.get("race", ""), {}).get("name", ""), c.get("name", "")], 13, UI.MUTED))
	content.add_child(head)
	var paper := UI.hbox(10)
	paper.alignment = BoxContainer.ALIGNMENT_CENTER
	paper.add_child(_col(LEFT))
	var mid := UI.vbox(8)
	mid.alignment = BoxContainer.ALIGNMENT_CENTER
	var doll := UI.label("👩" if S.self_info.get("look", {}).get("gender", "m") == "f" else "🧑", 80, Color.WHITE, HORIZONTAL_ALIGNMENT_CENTER)
	doll.custom_minimum_size = Vector2(110, 120)
	mid.add_child(doll)
	var weps := UI.hbox(6)
	weps.alignment = BoxContainer.ALIGNMENT_CENTER
	for sl in WEAP:
		weps.add_child(_slot(sl))
	mid.add_child(weps)
	paper.add_child(mid)
	paper.add_child(_col(RIGHT))
	content.add_child(paper)
	# stats
	var prim: String = PRIMARY.get(cls_id, "")
	var grid := GridContainer.new()
	grid.columns = 2
	grid.add_theme_constant_override("h_separation", 30)
	grid.add_theme_constant_override("v_separation", 1)
	var rows: Array = [
		["Lebenspunkte", "%d / %d" % [me.hp, me.mh], false], [RES_NAME.get(c.get("resource", ""), "Ressource"), "%d / %d" % [me.res, me.rm], false],
	]
	for k in ["str", "agi", "int", "sta"]:
		rows.append([ItemTip.STAT_LABEL[k], str(int(st.get(k, 0))), k == prim])
	rows.append(["Rüstung", "%d (−%d%% Schaden)" % [int(st.get("armor", 0)), int(st.get("red", 0))], false])
	rows.append(["Waffenschaden", str(int(st.get("weapon", 0))), false])
	rows.append(["Schadensbonus", "+%d%%" % int(st.get("dmg", 0)), false])
	rows.append(["Erfahrung", ("%d / %d" % [me.xp, me.xn]) if me.xn > 0 else "Max", false])
	for r in rows:
		var l := UI.label(r[0], 13, UI.GOLD if r[2] else UI.MUTED)
		l.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		grid.add_child(l)
		var v := UI.label(r[1], 13, UI.GOLD if r[2] else UI.TEXT, HORIZONTAL_ALIGNMENT_RIGHT)
		grid.add_child(v)
	grid.add_child(UI.label("Geld", 13, UI.MUTED))
	var money := UI.rich(UI.copper_text(S.copper), 13)
	money.fit_content = true
	money.custom_minimum_size.x = 140
	grid.add_child(money)
	content.add_child(grid)
	# equipment manager
	content.add_child(HSeparator.new())
	content.add_child(UI.label("Ausrüstungsmanager", 14, UI.GOLD))
	var names: Array = S.sets.keys()
	names.sort()
	for n in names:
		var set: Dictionary = S.sets[n]
		var row := UI.hbox(4)
		var nm := UI.label(n, 13, UI.TEXT)
		nm.custom_minimum_size.x = 90
		row.add_child(nm)
		var icons := ""
		var cnt := 0
		for slot_id in set:
			if cnt >= 8:
				break
			icons += S.item(set[slot_id]).get("icon", "")
			cnt += 1
		var il := UI.label(icons, 15)
		il.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		row.add_child(il)
		var nn: String = n
		row.add_child(UI.button("Anlegen", func(): Net.send("equipSet", {"name": nn})))
		var ow := UI.button("Überschreiben", func(): Net.send("saveSet", {"name": nn}))
		ow.tooltip_text = "Aktuelle Ausrüstung in diesem Set speichern"
		row.add_child(ow)
		row.add_child(UI.button("✕", func(): Net.send("deleteSet", {"name": nn})))
		content.add_child(row)
	if names.is_empty():
		content.add_child(UI.label("Noch keine Sets. Lege deine Ausrüstung an und speichere sie.", 12, UI.MUTED))
	var add := UI.hbox(6)
	set_name = LineEdit.new()
	set_name.max_length = 16
	set_name.placeholder_text = "Setname"
	set_name.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	set_name.focus_entered.connect(func(): S.typing = true)
	set_name.focus_exited.connect(func(): S.typing = false)
	set_name.text_submitted.connect(func(_t): _save())
	add.add_child(set_name)
	add.add_child(UI.button("Aktuelle Ausrüstung speichern", _save))
	content.add_child(add)
	reset_size()

func _save() -> void:
	var n := set_name.text.strip_edges()
	if n == "":
		return
	Net.send("saveSet", {"name": n})

func _col(slots: Array) -> VBoxContainer:
	var v := UI.vbox(5)
	for s in slots:
		v.add_child(_slot(s))
	return v

func _slot(slot: String) -> Slot:
	var s := Slot.new(46)
	var id: String = S.equip.get(slot, "")
	var d := S.item(id) if id != "" else {}
	if not d.is_empty():
		s.icon = d.icon
		s.border = UI.rarity_color(d.get("rarity", "common"))
		s.hovered.connect(func(on):
			if on:
				Tooltip.follow_mouse(ItemTip.item(id, 1, "eq"), s)
			else:
				Tooltip.hide_tip())
		s.clicked.connect(func(b): if b == MOUSE_BUTTON_RIGHT: Tooltip.hide_tip(); Net.send("unequip", {"slot": slot}))
		s.drag_data = {"from": "eq", "slot": slot}
	else:
		s.placeholder = ItemTip.SLOT_ICON.get(slot, "")
		s.border = Color(0.84, 0.7, 0.38, 0.22)
		s.hovered.connect(func(on):
			if on:
				Tooltip.follow_mouse("[color=#999999]%s[/color]" % ItemTip.SLOT_LABEL.get(slot, slot), s)
			else:
				Tooltip.hide_tip())
	s.accept_drop = func(data): return data is Dictionary and data.get("from") == "bag"
	s.dropped.connect(func(data): Net.send("equip", {"idx": int(data.idx), "slot": slot}))
	return s
