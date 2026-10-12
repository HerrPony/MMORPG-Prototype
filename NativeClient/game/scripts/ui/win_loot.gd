class_name WinLoot
extends GameWindow
## Loot window of a corpse or container.

var list: VBoxContainer
var loot_id := 0
var data: Dictionary = {}

func _init() -> void:
	super._init("Beute", Vector2(280, 100))
	win_id = "loot"
	default_pos = Vector2(700, 280)
	list = UI.vbox(3)
	body.add_child(list)
	body.add_child(UI.button("Alles nehmen", func(): if loot_id != 0: Net.send("loot", {"id": loot_id, "idx": -1}), true))
	closed.connect(func(): loot_id = 0)

func show_loot(m: Dictionary) -> void:
	loot_id = int(m.id)
	data = m
	title_label.text = m.name
	UI.clear(list)
	var idx := 0
	for it in m.get("items", []):
		var d := S.item(it.i)
		var i2 := idx
		var card := UI.card(func(): Net.send("loot", {"id": int(m.id), "idx": i2}))
		var h := UI.hbox(8)
		var ic := Slot.new(34)
		ic.icon = d.get("icon", "?")
		ic.border = UI.rarity_color(d.get("rarity", "common"))
		ic.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(ic)
		var nm := UI.label(d.get("name", it.i) + (" x%d" % int(it.c) if int(it.c) > 1 else ""), 14, UI.rarity_color(d.get("rarity", "common")))
		nm.mouse_filter = Control.MOUSE_FILTER_IGNORE
		nm.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		h.add_child(nm)
		h.mouse_filter = Control.MOUSE_FILTER_IGNORE
		card.add_child(h)
		var iid: String = it.i
		var cnt := int(it.c)
		card.gui_input.connect(func(ev: InputEvent):
			if ev is InputEventMouseMotion:
				Tooltip.follow_mouse(ItemTip.item(iid, cnt), card))
		card.mouse_exited.connect(Tooltip.hide_tip)
		list.add_child(card)
		idx += 1
	if int(m.get("copper", 0)) > 0:
		var card2 := UI.card(func(): Net.send("loot", {"id": int(m.id), "idx": -2}))
		var rt := UI.rich("💰 " + UI.copper_text(int(m.copper)), 14)
		rt.mouse_filter = Control.MOUSE_FILTER_IGNORE
		card2.add_child(rt)
		list.add_child(card2)
	if not visible:
		open()
	reset_size()
