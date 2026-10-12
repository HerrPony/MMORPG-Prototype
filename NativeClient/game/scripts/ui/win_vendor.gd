class_name WinVendor
extends GameWindow
## Vendor: list of goods, filter for usable items, sell junk.

var list: VBoxContainer
var money: RichTextLabel
var only_usable := true

func _init() -> void:
	super._init("Händler", Vector2(380, 100))
	win_id = "vendor"
	default_pos = Vector2(760, 130)
	var top := UI.hbox(8)
	var cb := CheckBox.new()
	cb.text = "Nur Nutzbares"
	cb.button_pressed = true
	cb.toggled.connect(func(on): only_usable = on; refresh())
	top.add_child(cb)
	var f := UI.fill()
	f.size_flags_vertical = Control.SIZE_SHRINK_BEGIN
	top.add_child(f)
	top.add_child(UI.button("🗑 Graues verkaufen", func(): if S.vendor_npc != 0: Net.send("sellJunk", {"npc": S.vendor_npc})))
	body.add_child(top)
	var sc := ScrollContainer.new()
	sc.custom_minimum_size = Vector2(360, 380)
	sc.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	list = UI.vbox(3)
	list.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	sc.add_child(list)
	body.add_child(sc)
	money = UI.rich("", 14)
	body.add_child(money)
	closed.connect(func(): S.vendor_npc = 0; S.vendor_data = {})

func open_for(m: Dictionary) -> void:
	S.vendor_npc = int(m.npc)
	S.vendor_data = m
	title_label.text = m.name
	open()
	refresh()

func refresh() -> void:
	var m: Dictionary = S.vendor_data
	if m.is_empty() or not visible:
		return
	UI.clear(list)
	var n := 0
	for id in m.get("stock", []):
		var it := S.item(id)
		if it.is_empty() or (only_usable and not WinBags._usable(it)):
			continue
		n += 1
		var buy := int(it.get("buy", 0))
		var card := UI.card(func(): pass)
		var h := UI.hbox(8)
		var ic := Slot.new(38)
		ic.icon = it.icon
		ic.border = UI.rarity_color(it.get("rarity", "common"))
		ic.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(ic)
		var v := UI.vbox(0)
		v.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		v.mouse_filter = Control.MOUSE_FILTER_IGNORE
		var nm := UI.label(it.name, 14, UI.rarity_color(it.get("rarity", "common")))
		nm.mouse_filter = Control.MOUSE_FILTER_IGNORE
		v.add_child(nm)
		var sub := ""
		if it.get("slot", "") != "":
			sub = "%s · %s" % [ItemTip.SLOT_LABEL.get(it.slot, ""), it.get("type", "")]
		if int(it.get("req", 0)) > 0:
			sub += " · Stufe %d" % int(it.req)
		var sl := UI.label(sub, 11, UI.MUTED)
		sl.mouse_filter = Control.MOUSE_FILTER_IGNORE
		v.add_child(sl)
		h.add_child(v)
		var price := UI.rich(UI.copper_text(buy), 13)
		price.custom_minimum_size.x = 90
		price.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(price)
		h.mouse_filter = Control.MOUSE_FILTER_IGNORE
		card.add_child(h)
		card.modulate = Color(1, 0.6, 0.6) if S.copper < buy else (Color(1, 1, 1, 0.55) if int(it.get("req", 0)) > S.me.lvl else Color.WHITE)
		var item_id: String = id
		card.gui_input.connect(func(ev: InputEvent):
			if ev is InputEventMouseButton and ev.pressed and ev.button_index == MOUSE_BUTTON_LEFT:
				Net.send("buy", {"npc": int(m.npc), "item": item_id, "count": 5 if (Input.is_key_pressed(KEY_SHIFT) and int(it.get("stack", 1)) > 1) else 1})
			elif ev is InputEventMouseMotion:
				Tooltip.follow_mouse(ItemTip.item(item_id, 1, "buy"), card))
		card.mouse_exited.connect(Tooltip.hide_tip)
		list.add_child(card)
	if n == 0:
		list.add_child(UI.label("Nichts im Angebot.", 13, UI.MUTED))
	money.text = UI.copper_text(S.copper)
