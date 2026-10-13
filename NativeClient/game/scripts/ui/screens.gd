class_name Screens
extends Control
## Login / register, character selection and character creation.

signal enter_world(char_id: String)
signal back_to_login

var preview: Preview
var mode := "login"
var chars: Array = []
var selected := ""   # character ids are hex strings, not numbers
var after_create := ""
var form := {"race": "human", "cls": "warrior", "gender": "m", "hair": 0, "skin": 0}
var status := ""
var _pending: Callable = Callable()

var login_box: Control
var chars_box: Control
var create_box: Control
var in_server: LineEdit
var in_user: LineEdit
var in_pass: LineEdit
var in_pass2: LineEdit
var login_err: Label
var login_status: Label
var btn_login: Button
var tab_login: Button
var tab_register: Button
var char_list: VBoxContainer
var char_info: Label
var btn_enter: Button
var btn_delete: Button
var btn_new: Button
var create_err: Label
var in_charname: LineEdit
var race_list: HBoxContainer
var class_list: HBoxContainer
var class_detail: VBoxContainer
var hair_list: HBoxContainer
var skin_list: HBoxContainer
var g_m: Button
var g_f: Button
var btn_back: Button
var title_label: Label

func _ready() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	mouse_filter = Control.MOUSE_FILTER_IGNORE
	_build_login()
	_build_chars()
	_build_create()
	show_screen("login")
	Net.on("auth", _on_auth)
	Net.on("authError", func(m): login_err.text = m.error)
	Net.on("chars", _on_chars)
	Net.on("created", func(m): after_create = str(m.id))
	Net.on("createError", func(m): create_err.text = m.error)
	Net.on("enterError", func(m): create_err.text = m.error; _alert(m.error))
	Net.connected.connect(_on_connected)
	Net.closed.connect(_on_closed)

func show_screen(id: String) -> void:
	login_box.visible = id == "login"
	chars_box.visible = id == "chars"
	create_box.visible = id == "create"
	visible = id != "none"
	if preview:
		preview.visible = id != "none"
		preview.center_offset = 0.0
		if id == "login":
			preview.clear()
		elif id == "chars":
			preview.center_offset = 0.0
		elif id == "create":
			preview.center_offset = 0.0

# ------------------------------------------------------------------ login

func _build_login() -> void:
	login_box = CenterContainer.new()
	login_box.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	login_box.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(login_box)
	var panel := PanelContainer.new()
	panel.custom_minimum_size = Vector2(400, 0)
	login_box.add_child(panel)
	var v := UI.vbox(8)
	panel.add_child(v)
	title_label = UI.label("ERLENWALD ONLINE", 30, UI.GOLD, HORIZONTAL_ALIGNMENT_CENTER)
	title_label.add_theme_constant_override("outline_size", 6)
	title_label.add_theme_color_override("font_outline_color", Color(0.2, 0.1, 0, 0.9))
	v.add_child(title_label)
	v.add_child(UI.label("Ein Low-Poly-MMORPG", 13, UI.MUTED, HORIZONTAL_ALIGNMENT_CENTER))
	v.add_child(UI.spacer(6))
	var tabs := UI.hbox(6)
	tab_login = UI.button("Anmelden", func(): _set_mode("login"))
	tab_register = UI.button("Neuer Account", func(): _set_mode("register"))
	tab_login.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	tab_register.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	tabs.add_child(tab_login)
	tabs.add_child(tab_register)
	v.add_child(tabs)
	in_server = _field(v, "Server", "https://game.herrpony.de", String(Settings.v["server"]))
	in_user = _field(v, "Benutzername", "", String(Settings.v["username"]))
	in_pass = _field(v, "Passwort", "", "", true)
	in_pass2 = _field(v, "Passwort wiederholen", "", "", true)
	in_pass2.get_parent().visible = false
	btn_login = UI.button("Anmelden", _submit, true)
	v.add_child(UI.spacer(4))
	v.add_child(btn_login)
	login_err = UI.label("", 13, UI.RED, HORIZONTAL_ALIGNMENT_CENTER)
	login_err.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	v.add_child(login_err)
	login_status = UI.label("", 12, UI.MUTED, HORIZONTAL_ALIGNMENT_CENTER)
	v.add_child(login_status)
	for f in [in_server, in_user, in_pass, in_pass2]:
		f.text_submitted.connect(func(_t): _submit())
	_set_mode("login")

func _field(parent: Control, label: String, placeholder: String, value: String, secret := false) -> LineEdit:
	var b := UI.vbox(2)
	b.add_child(UI.label(label, 12, UI.MUTED))
	var e := LineEdit.new()
	e.placeholder_text = placeholder
	e.text = value
	e.secret = secret
	e.max_length = 64
	b.add_child(e)
	parent.add_child(b)
	return e

func _set_mode(m: String) -> void:
	mode = m
	in_pass2.get_parent().visible = m == "register"
	btn_login.text = "Anmelden" if m == "login" else "Account erstellen"
	tab_login.disabled = m == "login"
	tab_register.disabled = m == "register"
	login_err.text = ""

func _server_url() -> String:
	var s := in_server.text.strip_edges()
	if s == "":
		s = "https://game.herrpony.de"
	if s.begins_with("ws://") or s.begins_with("wss://"):
		return s if s.ends_with("/ws") else s.trim_suffix("/") + "/ws"
	if s.begins_with("http://"):
		s = s.substr(7)
	elif s.begins_with("https://"):
		return "wss://" + s.substr(8).trim_suffix("/") + "/ws"
	return "ws://" + s.trim_suffix("/") + "/ws"

func _submit() -> void:
	var user := in_user.text.strip_edges()
	var pw := in_pass.text
	if user == "" or pw == "":
		login_err.text = "Bitte Benutzername und Passwort eingeben."
		return
	if mode == "register" and pw != in_pass2.text:
		login_err.text = "Die Passwörter stimmen nicht überein."
		return
	login_err.text = ""
	Settings.v["server"] = in_server.text.strip_edges()
	Settings.v["username"] = user
	Settings.save_file()
	var send_it := func(): Net.send(mode, {"user": user, "pass": pw})
	if Net.is_open:
		send_it.call()
	else:
		_pending = send_it
		login_status.text = "Verbinde mit %s …" % _server_url()
		btn_login.disabled = true
		var err := Net.connect_to(_server_url())
		if err != OK:
			login_status.text = ""
			login_err.text = "Ungültige Serveradresse."
			btn_login.disabled = false

func _on_connected() -> void:
	login_status.text = "Verbunden."
	btn_login.disabled = false
	if _pending.is_valid():
		var p := _pending
		_pending = Callable()
		p.call()

func _on_closed() -> void:
	btn_login.disabled = false
	if not S.in_world and visible and login_box.visible:
		login_status.text = ""
		login_err.text = "Server nicht erreichbar. Läuft der Server und stimmt die Adresse?"
	_pending = Callable()

func _on_auth(m: Dictionary) -> void:
	S.user = m.user
	S.admin = bool(m.get("admin", false))
	S.classes = m.classes
	S.class_order = m.classOrder
	S.races = m.races
	S.race_order = m.raceOrder
	S.all_skills = m.skills
	login_err.text = ""
	in_pass.text = ""
	in_pass2.text = ""
	login_status.text = ""

func _alert(t: String) -> void:
	var d := AcceptDialog.new()
	d.dialog_text = t
	d.title = "Hinweis"
	add_child(d)
	d.popup_centered()
	d.confirmed.connect(d.queue_free)
	d.canceled.connect(d.queue_free)

# ------------------------------------------------------------------ characters

func _build_chars() -> void:
	chars_box = Control.new()
	chars_box.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	chars_box.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(chars_box)
	var left := PanelContainer.new()
	left.set_anchors_and_offsets_preset(Control.PRESET_LEFT_WIDE)
	left.offset_left = 24
	left.offset_top = 24
	left.offset_bottom = -24
	left.offset_right = 344
	chars_box.add_child(left)
	var lv := UI.vbox(8)
	left.add_child(lv)
	lv.add_child(UI.label("Charaktere", 22, UI.GOLD))
	char_list = UI.vbox(6)
	var sc := ScrollContainer.new()
	sc.size_flags_vertical = Control.SIZE_EXPAND_FILL
	sc.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	char_list.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	sc.add_child(char_list)
	lv.add_child(sc)
	btn_new = UI.button("＋ Neuer Charakter", func(): _open_create())
	lv.add_child(btn_new)
	lv.add_child(UI.button("Abmelden", func(): back_to_login.emit()))
	var bottom := VBoxContainer.new()
	bottom.set_anchors_and_offsets_preset(Control.PRESET_CENTER_BOTTOM)
	bottom.grow_horizontal = Control.GROW_DIRECTION_BOTH
	bottom.grow_vertical = Control.GROW_DIRECTION_BEGIN
	bottom.offset_bottom = -40
	bottom.alignment = BoxContainer.ALIGNMENT_END
	bottom.add_theme_constant_override("separation", 8)
	chars_box.add_child(bottom)
	char_info = UI.label("", 20, UI.TEXT, HORIZONTAL_ALIGNMENT_CENTER)
	char_info.add_theme_constant_override("outline_size", 6)
	char_info.add_theme_color_override("font_outline_color", Color(0, 0, 0, 0.9))
	bottom.add_child(char_info)
	var row := UI.hbox(10)
	row.alignment = BoxContainer.ALIGNMENT_CENTER
	btn_enter = UI.button("Welt betreten", _enter_selected, true)
	btn_enter.custom_minimum_size.x = 220
	btn_delete = UI.button("Löschen", _delete_selected)
	row.add_child(btn_enter)
	row.add_child(btn_delete)
	bottom.add_child(row)

func _on_chars(m: Dictionary) -> void:
	chars = m.get("chars", [])
	if after_create != "":
		selected = after_create
		after_create = ""
	var found := false
	for c in chars:
		if str(c.id) == selected:
			found = true
	if not found:
		selected = str(chars[0].id) if chars.size() > 0 else ""
	if chars.is_empty():
		_open_create()
		return
	show_screen("chars")
	_render_chars()

func _render_chars() -> void:
	UI.clear(char_list)
	for c in chars:
		var id := str(c.id)
		var card := UI.card(func(): selected = id; _render_chars(), Vector2(0, 58))
		UI.set_selected(card, id == selected)
		var h := UI.hbox(10)
		var ico := UI.label(UI.CLASS_ICONS.get(c["class"], "?"), 26, Color.WHITE, HORIZONTAL_ALIGNMENT_CENTER)
		ico.custom_minimum_size.x = 40
		ico.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(ico)
		var v := UI.vbox(0)
		var nm := UI.label(c.name, 17, UI.CLASS_COLORS.get(c["class"], Color.WHITE))
		nm.mouse_filter = Control.MOUSE_FILTER_IGNORE
		var sub := UI.label("Stufe %d %s %s" % [int(c.level), S.races[c.race].name, S.classes[c["class"]].name], 12, UI.MUTED)
		sub.mouse_filter = Control.MOUSE_FILTER_IGNORE
		v.add_child(nm)
		v.add_child(sub)
		v.mouse_filter = Control.MOUSE_FILTER_IGNORE
		h.add_child(v)
		h.mouse_filter = Control.MOUSE_FILTER_IGNORE
		card.add_child(h)
		card.gui_input.connect(func(ev: InputEvent):
			if ev is InputEventMouseButton and ev.double_click and ev.button_index == MOUSE_BUTTON_LEFT:
				enter_world.emit(id))
		char_list.add_child(card)
	btn_new.visible = chars.size() < 5
	var sel = null
	for c in chars:
		if str(c.id) == selected:
			sel = c
	btn_enter.disabled = sel == null
	btn_delete.disabled = sel == null
	if sel:
		char_info.text = "%s – Stufe %d %s %s" % [sel.name, int(sel.level), S.races[sel.race].name, S.classes[sel["class"]].name]
		char_info.add_theme_color_override("font_color", UI.CLASS_COLORS.get(sel["class"], Color.WHITE))
		preview.set_style(sel["class"], sel.look)
	else:
		char_info.text = ""
		preview.clear()

func _enter_selected() -> void:
	if selected != "":
		enter_world.emit(selected)

func _delete_selected() -> void:
	var sel = null
	for c in chars:
		if str(c.id) == selected:
			sel = c
	if sel == null:
		return
	var d := ConfirmationDialog.new()
	d.title = "Charakter löschen"
	d.dialog_text = "Charakter \"%s\" wirklich endgültig löschen?" % sel.name
	d.ok_button_text = "Löschen"
	d.cancel_button_text = "Abbrechen"
	add_child(d)
	d.confirmed.connect(func(): Net.send("deleteChar", {"id": str(sel.id)}); d.queue_free())
	d.canceled.connect(d.queue_free)
	d.popup_centered()

# ------------------------------------------------------------------ creation

func _build_create() -> void:
	create_box = Control.new()
	create_box.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	create_box.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(create_box)
	var left := PanelContainer.new()
	left.set_anchors_and_offsets_preset(Control.PRESET_LEFT_WIDE)
	left.offset_left = 24; left.offset_top = 24; left.offset_bottom = -24; left.offset_right = 420
	create_box.add_child(left)
	var lscroll := ScrollContainer.new()
	lscroll.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	left.add_child(lscroll)
	var lv := UI.vbox(8)
	lv.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	lscroll.add_child(lv)
	lv.add_child(UI.label("Charakter erstellen", 22, UI.GOLD))
	lv.add_child(UI.label("Volk", 14, UI.MUTED))
	race_list = UI.hbox(6)
	lv.add_child(race_list)
	lv.add_child(UI.label("Klasse", 14, UI.MUTED))
	class_list = UI.hbox(6)
	lv.add_child(class_list)
	lv.add_child(UI.label("Aussehen", 14, UI.MUTED))
	var gr := UI.hbox(6)
	g_m = UI.button("♂ Männlich", func(): form.gender = "m"; _render_create())
	g_f = UI.button("♀ Weiblich", func(): form.gender = "f"; _render_create())
	gr.add_child(g_m)
	gr.add_child(g_f)
	lv.add_child(gr)
	lv.add_child(UI.label("Haarfarbe", 12, UI.MUTED))
	hair_list = UI.hbox(6)
	lv.add_child(hair_list)
	lv.add_child(UI.label("Hautfarbe", 12, UI.MUTED))
	skin_list = UI.hbox(6)
	lv.add_child(skin_list)
	lv.add_child(UI.label("Name", 14, UI.MUTED))
	in_charname = LineEdit.new()
	in_charname.max_length = 16
	in_charname.placeholder_text = "Charaktername"
	in_charname.text_submitted.connect(func(_t): _create())
	lv.add_child(in_charname)
	var br := UI.hbox(8)
	btn_back = UI.button("Zurück", func(): show_screen("chars"); _render_chars())
	br.add_child(btn_back)
	var create := UI.button("Erstellen", _create, true)
	create.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	br.add_child(create)
	lv.add_child(br)
	create_err = UI.label("", 13, UI.RED)
	create_err.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	lv.add_child(create_err)
	# class details on the right
	var right := PanelContainer.new()
	right.set_anchors_and_offsets_preset(Control.PRESET_RIGHT_WIDE)
	right.offset_left = -400; right.offset_right = -24; right.offset_top = 24; right.offset_bottom = -24
	create_box.add_child(right)
	var rs := ScrollContainer.new()
	rs.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	right.add_child(rs)
	class_detail = UI.vbox(6)
	class_detail.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	rs.add_child(class_detail)

func _open_create() -> void:
	create_err.text = ""
	in_charname.text = ""
	btn_back.visible = chars.size() > 0
	show_screen("create")
	_render_create()

func _render_create() -> void:
	UI.clear(race_list)
	for id in S.race_order:
		var r: Dictionary = S.races[id]
		var cid: String = id
		var c := UI.card(func(): form.race = cid; _render_create(), Vector2(0, 0))
		UI.set_selected(c, form.race == cid)
		c.tooltip_text = r.desc
		var v := UI.vbox(2)
		var nm := UI.label(r.name, 14, UI.TEXT, HORIZONTAL_ALIGNMENT_CENTER)
		nm.mouse_filter = Control.MOUSE_FILTER_IGNORE
		v.add_child(nm)
		var ds := UI.label(r.desc, 10, UI.MUTED, HORIZONTAL_ALIGNMENT_CENTER)
		ds.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		ds.custom_minimum_size.x = 96
		ds.mouse_filter = Control.MOUSE_FILTER_IGNORE
		v.add_child(ds)
		v.mouse_filter = Control.MOUSE_FILTER_IGNORE
		c.add_child(v)
		c.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		race_list.add_child(c)
	UI.clear(class_list)
	for id in S.class_order:
		var cl: Dictionary = S.classes[id]
		var cid2: String = id
		var c2 := UI.card(func(): form.cls = cid2; _render_create())
		UI.set_selected(c2, form.cls == cid2)
		var v2 := UI.vbox(2)
		var ic := UI.label(UI.CLASS_ICONS.get(id, "?"), 28, Color.WHITE, HORIZONTAL_ALIGNMENT_CENTER)
		ic.mouse_filter = Control.MOUSE_FILTER_IGNORE
		var nm2 := UI.label(cl.name, 14, UI.CLASS_COLORS.get(id, Color.WHITE), HORIZONTAL_ALIGNMENT_CENTER)
		nm2.mouse_filter = Control.MOUSE_FILTER_IGNORE
		var rl := UI.label(cl.role, 10, UI.MUTED, HORIZONTAL_ALIGNMENT_CENTER)
		rl.mouse_filter = Control.MOUSE_FILTER_IGNORE
		v2.add_child(ic); v2.add_child(nm2); v2.add_child(rl)
		v2.mouse_filter = Control.MOUSE_FILTER_IGNORE
		c2.add_child(v2)
		c2.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		class_list.add_child(c2)
	_swatches(hair_list, Models.HAIR, "hair")
	_swatches(skin_list, Models.SKIN, "skin")
	g_m.disabled = form.gender == "m"
	g_f.disabled = form.gender == "f"
	# class details
	UI.clear(class_detail)
	var c: Dictionary = S.classes[form.cls]
	class_detail.add_child(UI.label(c.name, 24, UI.CLASS_COLORS.get(form.cls, Color.WHITE)))
	var res := {"rage": "Wut", "focus": "Fokus", "mana": "Mana"}
	class_detail.add_child(UI.label("%s · Ressource: %s" % [c.role, res.get(c.resource, "")], 12, UI.MUTED))
	var ds2 := UI.label(c.desc, 13, Color("cdbf9f"))
	ds2.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	class_detail.add_child(ds2)
	class_detail.add_child(UI.label("Fertigkeiten", 16, UI.GOLD))
	for sid in c.skills:
		var s: Dictionary = S.all_skills[sid]
		var row := UI.hbox(8)
		var ico := UI.label(s.icon, 24, Color.WHITE, HORIZONTAL_ALIGNMENT_CENTER)
		ico.custom_minimum_size = Vector2(36, 0)
		row.add_child(ico)
		var col := UI.vbox(0)
		col.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		var lvtxt := ("Stufe %d · beim Lehrer" % int(s.level)) if int(s.level) > 1 else "von Beginn an"
		col.add_child(UI.label("%s  (%s)" % [s.name, lvtxt], 13, UI.TEXT))
		var d := UI.label(s.desc, 11, UI.MUTED)
		d.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		col.add_child(d)
		row.add_child(col)
		class_detail.add_child(row)
	preview.set_style(form.cls, {"gender": form.gender, "hair": form.hair, "skin": form.skin})

func _swatches(host: HBoxContainer, arr: Array, key: String) -> void:
	UI.clear(host)
	for i in arr.size():
		var idx := i
		var b := Button.new()
		b.custom_minimum_size = Vector2(32, 32)
		var col := Models.C(arr[i])
		var sel: bool = form[key] == i
		b.add_theme_stylebox_override("normal", UI.box(col, UI.GOLD if sel else Color(0, 0, 0, 0.5), 3 if sel else 1, 4))
		b.add_theme_stylebox_override("hover", UI.box(col.lightened(0.15), UI.GOLD, 2, 4))
		b.add_theme_stylebox_override("pressed", UI.box(col, UI.GOLD, 3, 4))
		b.pressed.connect(func(): form[key] = idx; _render_create())
		host.add_child(b)

func _create() -> void:
	Net.send("createChar", {"name": in_charname.text.strip_edges(), "race": form.race, "class": form.cls, "look": {"gender": form.gender, "hair": form.hair, "skin": form.skin}})
