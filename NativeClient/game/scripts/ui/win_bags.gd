class_name WinBags
extends GameWindow
## Inventory: slot grid with tooltips, use/equip/sell, drag and drop with the equipment.

var grid: GridContainer
var money: RichTextLabel
var slot_nodes: Array = []

func _init() -> void:
	super._init("Taschen", Vector2(340, 100))
	win_id = "bags"
	default_pos = Vector2(1180, 330)
	grid = GridContainer.new()
	grid.columns = 7
	grid.add_theme_constant_override("h_separation", 4)
	grid.add_theme_constant_override("v_separation", 4)
	body.add_child(grid)
	money = UI.rich("", 14)
	money.fit_content = true
	body.add_child(money)

func refresh() -> void:
	if not visible and slot_nodes.size() > 0:
		return
	UI.clear(grid)
	slot_nodes.clear()
	for i in S.inv_max:
		var s := Slot.new(44)
		var st = S.inv[i] if i < S.inv.size() else null
		if st != null:
			var d := S.item(st.i)
			s.icon = d.get("icon", "?")
			s.count = int(st.c)
			s.border = UI.rarity_color(d.get("rarity", "common"))
			if d.get("slot", "") != "" and (int(d.get("req", 0)) > S.me.lvl or not _usable(d)):
				s.dim = true
			var idx := i
			s.hovered.connect(func(on): _tip(on, s, st))
			s.clicked.connect(func(b): if b == MOUSE_BUTTON_RIGHT: use_item(idx, false))
			s.double_clicked.connect(func(): use_item(idx, true))
			s.drag_data = {"from": "bag", "idx": i}
		else:
			s.border = Color(0.84, 0.7, 0.38, 0.18)
		s.accept_drop = func(data): return data is Dictionary and data.get("from") == "eq"
		s.dropped.connect(func(data): Net.send("unequip", {"slot": data.slot}))
		grid.add_child(s)
		slot_nodes.append(s)
	money.text = UI.copper_text(S.copper)
	reset_size()

func _tip(on: bool, s: Slot, st: Dictionary) -> void:
	if not on:
		Tooltip.hide_tip()
		return
	Tooltip.follow_mouse(ItemTip.item(st.i, int(st.c), "bag"), s)

static func _usable(it: Dictionary) -> bool:
	var classes = it.get("classes", [])
	return not (classes is Array) or classes.is_empty() or classes.has(S.self_info.get("class", ""))

static func use_item(idx: int, no_sell: bool) -> void:
	if idx >= S.inv.size():
		return
	var st: Dictionary = S.inv[idx]
	var it := S.item(st.i)
	Tooltip.hide_tip()
	if S.vendor_npc != 0 and not no_sell:
		Net.send("sell", {"npc": S.vendor_npc, "idx": idx})
	elif it.get("slot", "") != "":
		Net.send("equip", {"idx": idx})
	elif it.get("heal", 0) != 0:
		Net.send("use", {"idx": idx})

