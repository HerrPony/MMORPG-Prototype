class_name Hud
extends Control
## In-game overlay: unit frames, cast bar, action bar, XP, chat, minimap, quest tracker, banners.

signal toggle_window(name: String)
signal request_logout

const STATUS_ICONS := {"slow": ["🐌", "Verlangsamt", true], "dot": ["🩸", "Schaden über Zeit", true], "buff": ["⬆️", "Gestärkt", false], "shield": ["🛡️", "Schild", false]}
const MOB_EMOJI := {"wolf": "🐺", "boar": "🐗", "spider": "🕷️", "bear": "🐻", "bandit": "🥷", "hyena": "🐕", "plains_bandit": "🥷", "bull": "🐂", "wolf_alpha": "🐺", "bandit_chief": "🥷", "water_elem": "💧", "croc": "🐊", "zombie": "🧟", "wisp": "✨", "witch": "🧙", "hag": "🧙", "zombie_giant": "🧟", "icewolf": "🐺", "yeti": "❄️", "yeti_chief": "❄️", "polarbear": "🐻‍❄️", "polarbear_big": "🐻‍❄️", "bear_big": "🐻", "frostspirit": "👻", "scorpion": "🦂", "scorpion_king": "🦂", "desert_raider": "🥷", "mummy": "🧟", "tomb_guardian": "🗿", "kobold": "👹", "kobold_shaman": "👹", "kobold_boss": "👹", "kobold_king": "👑", "skeleton": "💀", "cavespider": "🕷️", "spider_queen": "🕷️", "sheep": "🐑", "water_elem_pet": "💧", "smith": "🔨"}
const NPC_EMOJI := {"forester": "🧔", "farmer": "👩‍🌾", "scholar": "👩‍🔬", "guard": "💂", "huntsman": "🧔🏿", "merchant": "👩‍💼", "trainer_warrior": "🛡️", "trainer_hunter": "🏹", "trainer_mage": "🔮", "mayor": "👑", "priest": "🙏", "elder": "👴", "herbalist": "🌿", "fisher": "🎣", "sheikh": "🕌", "nomad": "🐪", "flightmaster": "🦅"}

var minimap: Minimap
var zone_label: Label
var clock_label: Label
var fps_label: Label
# frames
var p_frame: PanelContainer
var p_name: Label
var p_lvl: Label
var p_portrait: Label
var p_portrait_bg: Panel
var p_hp: Bar
var p_res: Bar
var p_buffs: HBoxContainer
var p_stance: Label
var t_frame: PanelContainer
var t_name: Label
var t_lvl: Label
var t_portrait: Label
var t_hp: Bar
var t_tag: Label
var t_buffs: HBoxContainer
var party_box: VBoxContainer
var pet_box: VBoxContainer
var _party_key := ""
var _pet_key := ""
var _party_bars: Dictionary = {}
var _pet_bars: Array = []
# bottom
var cast_bar: Bar
var xp_bar: Bar
var slots: Array = []
var _slot_rows: HBoxContainer
# chat
var chat_log: RichTextLabel
var chat_in: LineEdit
var chat_filter := "all"
var chat_lines: Array = []
var chat_tabs: Array = []
# misc
var tracker_body: VBoxContainer
var tracker_panel: PanelContainer
var err_label: Label
var err_t := 0.0
var banner: VBoxContainer
var banner_zone: Label
var banner_sub: Label
var banner_t := 0.0
var levelup: Label
var levelup_t := 0.0
var death_box: PanelContainer
var flight_bar: PanelContainer
var flight_dest: Label
var invite_box: PanelContainer
var invite_text: Label
var invite_t := 0.0
var player_menu: PopupMenu
var _menu_ent := 0
var toast_label: Label
var toast_t := 0.0
var fct_root: Control
var _last_clock := ""
var _q_timers: Array = []

func _ready() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	mouse_filter = Control.MOUSE_FILTER_IGNORE
	theme = UI.theme()
	_build_frames()
	_build_minimap()
	_build_bottom()
	_build_chat()
	_build_misc()
	_connect_net()

# ---------------------------------------------------------------- build

func _frame_panel(pos: Vector2, w: float) -> PanelContainer:
	var p := PanelContainer.new()
	p.position = pos
	p.custom_minimum_size.x = w
	p.add_theme_stylebox_override("panel", UI.box(Color(0.06, 0.055, 0.05, 0.82), UI.EDGE, 1, 8, 6))
	return p

func _portrait(size_px: float) -> Array:
	var holder := Control.new()
	holder.custom_minimum_size = Vector2(size_px, size_px)
	var bg := Panel.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	holder.add_child(bg)
	var l := Label.new()
	l.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	l.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	l.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	l.add_theme_font_size_override("font_size", int(size_px * 0.62))
	holder.add_child(l)
	return [holder, bg, l]

func _build_frames() -> void:
	p_frame = _frame_panel(Vector2(16, 16), 270)
	add_child(p_frame)
	var h := UI.hbox(8)
	p_frame.add_child(h)
	var pr := _portrait(60)
	h.add_child(pr[0])
	p_portrait_bg = pr[1]
	p_portrait = pr[2]
	var v := UI.vbox(3)
	v.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	h.add_child(v)
	var top := UI.hbox(6)
	p_name = UI.label("", 15, Color.WHITE)
	p_name.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	p_lvl = UI.label("1", 14, UI.GOLD)
	top.add_child(p_name)
	top.add_child(p_lvl)
	v.add_child(top)
	p_hp = Bar.new(190, 17, Color("3ccf4e"))
	p_hp.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	v.add_child(p_hp)
	p_res = Bar.new(190, 14, Color("e24a3a"))
	p_res.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	v.add_child(p_res)
	p_stance = UI.label("", 11, UI.MUTED)
	v.add_child(p_stance)
	p_buffs = UI.hbox(2)
	v.add_child(p_buffs)

	t_frame = _frame_panel(Vector2(300, 16), 270)
	t_frame.visible = false
	add_child(t_frame)
	var th := UI.hbox(8)
	t_frame.add_child(th)
	var tv := UI.vbox(3)
	tv.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	th.add_child(tv)
	var tt := UI.hbox(6)
	t_name = UI.label("", 15, Color.WHITE)
	t_name.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	t_tag = UI.label("", 11, Color("ffb040"))
	t_lvl = UI.label("", 14, UI.GOLD)
	tt.add_child(t_name)
	tt.add_child(t_tag)
	tt.add_child(t_lvl)
	tv.add_child(tt)
	t_hp = Bar.new(190, 17, Color("e24a3a"))
	t_hp.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	tv.add_child(t_hp)
	t_buffs = UI.hbox(2)
	tv.add_child(t_buffs)
	var tp := _portrait(52)
	th.add_child(tp[0])
	t_portrait = tp[2]
	(tp[1] as Panel).add_theme_stylebox_override("panel", UI.box(Color(0.15, 0.13, 0.11), UI.EDGE, 1, 26))

	party_box = UI.vbox(4)
	party_box.position = Vector2(16, 112)
	add_child(party_box)
	pet_box = UI.vbox(4)
	pet_box.position = Vector2(300, 82)
	add_child(pet_box)

func _build_minimap() -> void:
	var holder := Control.new()
	holder.set_anchors_and_offsets_preset(Control.PRESET_TOP_RIGHT)
	holder.offset_left = -204
	holder.offset_right = -16
	holder.offset_top = 16
	holder.offset_bottom = 300
	holder.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(holder)
	minimap = Minimap.new()
	minimap.position = Vector2(6, 26)
	holder.add_child(minimap)
	zone_label = UI.label("Erlenwald", 14, UI.GOLD, HORIZONTAL_ALIGNMENT_CENTER)
	zone_label.custom_minimum_size = Vector2(188, 0)
	zone_label.position = Vector2(0, 0)
	zone_label.add_theme_constant_override("outline_size", 5)
	zone_label.add_theme_color_override("font_outline_color", Color.BLACK)
	holder.add_child(zone_label)
	clock_label = UI.label("", 13, UI.TEXT, HORIZONTAL_ALIGNMENT_CENTER)
	clock_label.custom_minimum_size = Vector2(188, 0)
	clock_label.position = Vector2(0, 206)
	clock_label.add_theme_constant_override("outline_size", 4)
	clock_label.add_theme_color_override("font_outline_color", Color.BLACK)
	holder.add_child(clock_label)
	fps_label = UI.label("", 12, Color("9fe89f"), HORIZONTAL_ALIGNMENT_RIGHT)
	fps_label.set_anchors_and_offsets_preset(Control.PRESET_TOP_RIGHT)
	fps_label.position = Vector2(-80, 0)
	add_child(fps_label)
	# quest tracker
	tracker_panel = PanelContainer.new()
	tracker_panel.set_anchors_and_offsets_preset(Control.PRESET_TOP_RIGHT)
	tracker_panel.offset_left = -270
	tracker_panel.offset_right = -16
	tracker_panel.offset_top = 236
	tracker_panel.add_theme_stylebox_override("panel", UI.box(Color(0.05, 0.045, 0.04, 0.55), UI.EDGE_SOFT, 1, 6, 8))
	tracker_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(tracker_panel)
	var tv := UI.vbox(2)
	tracker_panel.add_child(tv)
	tv.add_child(UI.label("Quests", 13, UI.GOLD))
	tracker_body = UI.vbox(1)
	tv.add_child(tracker_body)

func _build_bottom() -> void:
	var bottom := VBoxContainer.new()
	bottom.set_anchors_and_offsets_preset(Control.PRESET_CENTER_BOTTOM)
	bottom.grow_horizontal = Control.GROW_DIRECTION_BOTH
	bottom.grow_vertical = Control.GROW_DIRECTION_BEGIN
	bottom.offset_bottom = -10
	bottom.add_theme_constant_override("separation", 5)
	bottom.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(bottom)
	cast_bar = Bar.new(300, 20, Color("e8b830"))
	cast_bar.font_size = 13
	cast_bar.visible = false
	cast_bar.size_flags_horizontal = Control.SIZE_SHRINK_CENTER
	bottom.add_child(cast_bar)
	_slot_rows = UI.hbox(4)
	_slot_rows.alignment = BoxContainer.ALIGNMENT_CENTER
	bottom.add_child(_slot_rows)
	var keys := ["1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=", "[", "]"]
	for i in 14:
		var s := Slot.new(46)
		s.key_text = keys[i]
		var idx := i
		s.clicked.connect(func(b): if b == MOUSE_BUTTON_LEFT: cast_slot(idx))
		s.hovered.connect(func(on): _skill_tip(idx, s, on))
		slots.append(s)
		_slot_rows.add_child(s)
	xp_bar = Bar.new(14 * 50.0, 12, Color("8a4ad8"))
	xp_bar.font_size = 10
	xp_bar.size_flags_horizontal = Control.SIZE_SHRINK_CENTER
	bottom.add_child(xp_bar)
	# micro menu
	var micro := UI.hbox(4)
	micro.set_anchors_and_offsets_preset(Control.PRESET_BOTTOM_RIGHT)
	micro.grow_horizontal = Control.GROW_DIRECTION_BEGIN
	micro.grow_vertical = Control.GROW_DIRECTION_BEGIN
	micro.offset_right = -16
	micro.offset_bottom = -14
	add_child(micro)
	for it in [["🎒", "Taschen (B)", "bags"], ["🧍", "Charakter (C)", "char"], ["📖", "Zauberbuch (P)", "spells"], ["📜", "Questlog (L)", "quests"], ["🗺️", "Weltkarte (M)", "map"], ["👥", "Gruppe (O)", "party"], ["⚙️", "Menü (Esc)", "menu"]]:
		var b := Button.new()
		b.text = it[0]
		b.tooltip_text = it[1]
		b.custom_minimum_size = Vector2(40, 40)
		b.add_theme_font_size_override("font_size", 20)
		var wn: String = it[2]
		b.pressed.connect(func(): toggle_window.emit(wn))
		micro.add_child(b)

func _build_chat() -> void:
	var panel := VBoxContainer.new()
	panel.set_anchors_and_offsets_preset(Control.PRESET_BOTTOM_LEFT)
	panel.grow_vertical = Control.GROW_DIRECTION_BEGIN
	panel.offset_left = 16
	panel.offset_bottom = -14
	panel.custom_minimum_size = Vector2(430, 0)
	panel.add_theme_constant_override("separation", 3)
	panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(panel)
	var tabs := UI.hbox(3)
	for f in [["all", "Allgemein"], ["combat", "Kampf"], ["loot", "Beute"]]:
		var b := Button.new()
		b.text = f[1]
		b.toggle_mode = true
		b.button_pressed = f[0] == "all"
		b.add_theme_font_size_override("font_size", 12)
		var fid: String = f[0]
		b.pressed.connect(func(): chat_filter = fid; _chat_refresh(); for t in chat_tabs: t.button_pressed = t == b)
		chat_tabs.append(b)
		tabs.add_child(b)
	panel.add_child(tabs)
	var bg := PanelContainer.new()
	bg.add_theme_stylebox_override("panel", UI.box(Color(0, 0, 0, 0.4), Color(0, 0, 0, 0), 0, 6, 6))
	bg.mouse_filter = Control.MOUSE_FILTER_IGNORE
	panel.add_child(bg)
	chat_log = RichTextLabel.new()
	chat_log.bbcode_enabled = true
	chat_log.scroll_following = true
	chat_log.custom_minimum_size = Vector2(0, 150)
	chat_log.mouse_filter = Control.MOUSE_FILTER_PASS
	chat_log.add_theme_font_size_override("normal_font_size", 13)
	chat_log.add_theme_font_size_override("bold_font_size", 13)
	bg.add_child(chat_log)
	chat_in = LineEdit.new()
	chat_in.placeholder_text = "Nachricht… (Enter)"
	chat_in.max_length = 200
	chat_in.focus_entered.connect(func(): S.typing = true)
	chat_in.focus_exited.connect(func(): S.typing = false)
	chat_in.text_submitted.connect(_chat_submit)
	chat_in.gui_input.connect(func(ev: InputEvent):
		if ev is InputEventKey and ev.pressed and ev.keycode == KEY_ESCAPE:
			chat_in.text = ""
			chat_in.release_focus()
			get_viewport().set_input_as_handled())
	panel.add_child(chat_in)

func _build_misc() -> void:
	fct_root = Control.new()
	fct_root.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	fct_root.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(fct_root)
	err_label = UI.label("", 18, Color("ff6a5a"), HORIZONTAL_ALIGNMENT_CENTER)
	err_label.set_anchors_and_offsets_preset(Control.PRESET_CENTER_TOP)
	err_label.grow_horizontal = Control.GROW_DIRECTION_BOTH
	err_label.offset_top = 150
	err_label.add_theme_constant_override("outline_size", 6)
	err_label.add_theme_color_override("font_outline_color", Color.BLACK)
	err_label.modulate.a = 0.0
	add_child(err_label)
	toast_label = UI.label("", 16, UI.GOLD, HORIZONTAL_ALIGNMENT_CENTER)
	toast_label.set_anchors_and_offsets_preset(Control.PRESET_CENTER_TOP)
	toast_label.grow_horizontal = Control.GROW_DIRECTION_BOTH
	toast_label.offset_top = 190
	toast_label.add_theme_constant_override("outline_size", 6)
	toast_label.add_theme_color_override("font_outline_color", Color.BLACK)
	toast_label.modulate.a = 0.0
	add_child(toast_label)
	banner = UI.vbox(0)
	banner.set_anchors_and_offsets_preset(Control.PRESET_CENTER_TOP)
	banner.grow_horizontal = Control.GROW_DIRECTION_BOTH
	banner.offset_top = 90
	banner.modulate.a = 0.0
	banner.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(banner)
	banner_zone = UI.label("", 40, Color("f6e2a0"), HORIZONTAL_ALIGNMENT_CENTER)
	banner_zone.add_theme_constant_override("outline_size", 10)
	banner_zone.add_theme_color_override("font_outline_color", Color(0.1, 0.05, 0, 0.95))
	banner_sub = UI.label("", 18, Color("d8cba4"), HORIZONTAL_ALIGNMENT_CENTER)
	banner_sub.add_theme_constant_override("outline_size", 6)
	banner_sub.add_theme_color_override("font_outline_color", Color.BLACK)
	banner.add_child(banner_zone)
	banner.add_child(banner_sub)
	levelup = UI.label("", 54, Color("ffeb3b"), HORIZONTAL_ALIGNMENT_CENTER)
	levelup.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	levelup.grow_horizontal = Control.GROW_DIRECTION_BOTH
	levelup.offset_top = -160
	levelup.add_theme_constant_override("outline_size", 12)
	levelup.add_theme_color_override("font_outline_color", Color(0.3, 0.15, 0, 1))
	levelup.modulate.a = 0.0
	add_child(levelup)
	death_box = PanelContainer.new()
	death_box.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	death_box.grow_horizontal = Control.GROW_DIRECTION_BOTH
	death_box.offset_top = -120
	death_box.visible = false
	add_child(death_box)
	var dv := UI.vbox(8)
	death_box.add_child(dv)
	dv.add_child(UI.label("Du bist gestorben", 24, Color("ff5a4a"), HORIZONTAL_ALIGNMENT_CENTER))
	dv.add_child(UI.button("Wiederbeleben", func(): Net.send("revive"), true))
	flight_bar = PanelContainer.new()
	flight_bar.set_anchors_and_offsets_preset(Control.PRESET_CENTER_TOP)
	flight_bar.grow_horizontal = Control.GROW_DIRECTION_BOTH
	flight_bar.offset_top = 70
	flight_bar.visible = false
	add_child(flight_bar)
	flight_dest = UI.label("", 16, UI.GOLD, HORIZONTAL_ALIGNMENT_CENTER)
	flight_bar.add_child(flight_dest)
	invite_box = PanelContainer.new()
	invite_box.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	invite_box.grow_horizontal = Control.GROW_DIRECTION_BOTH
	invite_box.offset_top = -200
	invite_box.visible = false
	add_child(invite_box)
	var iv := UI.vbox(8)
	invite_box.add_child(iv)
	invite_text = UI.label("", 15, UI.TEXT, HORIZONTAL_ALIGNMENT_CENTER)
	iv.add_child(invite_text)
	var ir := UI.hbox(8)
	ir.alignment = BoxContainer.ALIGNMENT_CENTER
	ir.add_child(UI.button("Annehmen", func(): Net.send("inviteReply", {"accept": true}); invite_box.visible = false, true))
	ir.add_child(UI.button("Ablehnen", func(): Net.send("inviteReply", {"accept": false}); invite_box.visible = false))
	iv.add_child(ir)
	player_menu = PopupMenu.new()
	player_menu.id_pressed.connect(_on_player_menu)
	add_child(player_menu)

# ---------------------------------------------------------------- net events

func _connect_net() -> void:
	Net.on("err", func(m): show_error(m.text))
	Net.on("chat", _on_chat)
	Net.on("lvl", func(m): levelup.text = "STUFE %d!" % int(m.lvl); levelup_t = 3.0)
	Net.on("cast", func(m): S.cast = {"skill": m.skill, "start": S.now(), "dur": float(m.dur)})
	Net.on("castStop", func(_m): S.cast = null)
	Net.on("cd", func(m):
		var t := S.now()
		if float(m.get("dur", 0)) > 0.0:
			S.cd[m.skill] = {"end": t + float(m.dur), "dur": float(m.dur)}
		if float(m.get("gcd", 0)) > 0.0:
			S.gcd = {"end": t + float(m.gcd), "dur": float(m.gcd)})
	Net.on("skills", func(m):
		var i: int = (m.skills as Array).find(m.learned)
		if i >= 0 and i < slots.size():
			slots[i].flash = 1.0)
	Net.on("zone", _on_zone)
	Net.on("partyInvite", func(m): invite_text.text = "%s lädt dich in seine Gruppe ein." % m.from; invite_box.visible = true; invite_t = 30.0)
	Net.on("flightStart", func(m): flight_dest.text = ("⛵ " if m.get("ship", false) else "🦅 ") + str(m.dest); flight_bar.visible = true)
	Net.on("flightEnd", func(_m): flight_bar.visible = false)
	Net.on("quests", func(m):
		var now := S.now()
		S.quests = []
		for q in m.list:
			var d: Dictionary = q.duplicate()
			d["end"] = (now + float(q.left)) if q.has("left") else 0.0
			S.quests.append(d)
		render_tracker())
	Net.on("party", func(m): S.party = {"id": m.id, "leader": m.leader, "members": m.members, "states": S.party.states})

func show_error(t: String) -> void:
	err_label.text = t
	err_t = 1.8

func toast(t: String) -> void:
	toast_label.text = t
	toast_t = 2.6

func _on_zone(m: Dictionary) -> void:
	var changed: bool = m.zone != S.zone.name
	S.zone = {"name": m.zone, "sub": m.get("sub", ""), "lo": m.get("lo", 0), "hi": m.get("hi", 0)}
	zone_label.text = ("%s – %s" % [m.zone, m.sub]) if m.get("sub", "") != "" else str(m.zone)
	if changed:
		var sub: String = m.get("sub", "")
		if sub == "":
			sub = "Offene See" if m.get("sea", false) else ("Stufe %d–%d" % [int(m.lo), int(m.hi)] if int(m.get("lo", 0)) > 0 else "")
		_show_banner(m.zone, sub)
	elif m.get("sub", "") != "":
		_show_banner(m.sub, m.zone)

func _show_banner(zone: String, sub: String) -> void:
	banner_zone.text = zone
	banner_sub.text = sub
	banner_t = 4.2

# ---------------------------------------------------------------- chat

func _on_chat(m: Dictionary) -> void:
	var ch: String = m.ch
	var text: String = m.text
	var color := "ffe27a"
	var kind := "all"
	match ch:
		"say":
			add_chat("[b][color=#ffffdd][%s][/color][/b]: %s" % [m.from, _esc(text)], "all")
			return
		"party":
			add_chat("[color=#7ab8ff][b][Gruppe] %s[/b]: %s[/color]" % [m.from, _esc(text)], "all")
			return
		"yell":
			add_chat("[color=#ff8a6a][b][%s][/b]: %s[/color]" % [m.get("from", "?"), _esc(text)], "all")
			return
		"combat": color = "ff9b6a"; kind = "combat"
		"xp": color = "b98cff"; kind = "combat"
		"loot": color = "5bd65b"; kind = "loot"
		"quest": color = "ffd650"; kind = "loot"
		"level": color = "ffeb3b"; kind = "loot"
	add_chat("[color=#%s]%s[/color]" % [color, _esc(text)], kind)

func _esc(s: String) -> String:
	return s.replace("[", "[lb]")

func add_chat(bb: String, kind := "all") -> void:
	chat_lines.append([bb, kind])
	if chat_lines.size() > 200:
		chat_lines.pop_front()
	if chat_filter == "all" or chat_filter == kind:
		chat_log.append_text(bb + "\n")

func _chat_refresh() -> void:
	chat_log.clear()
	for l in chat_lines:
		if chat_filter == "all" or chat_filter == l[1]:
			chat_log.append_text(l[0] + "\n")

func open_chat(prefix: String) -> void:
	chat_in.text = prefix
	chat_in.grab_focus()
	chat_in.caret_column = prefix.length()

func _chat_submit(t: String) -> void:
	chat_in.text = ""
	chat_in.release_focus()
	t = t.strip_edges()
	if t == "":
		return
	if t == "/hilfe" or t == "/help":
		add_chat("[color=#ffe27a]Steuerung: WASD bewegen · Leertaste springen · Rechtsklick halten: Kamera · Linksklick: auswählen · Rechtsklick: interagieren/angreifen/looten · Tab: Ziel wechseln · 1-5: Fertigkeiten · B/L/C/P: Fenster · R: Autorun · F: mit Ziel interagieren[/color]")
		return
	Net.send("chat", {"text": t.substr(3) if t.begins_with("/s ") else t})

# ---------------------------------------------------------------- action bar

func skill_at(i: int) -> Dictionary:
	if i >= S.skill_ids.size():
		return {}
	return S.skill(S.skill_ids[i])

func cast_slot(i: int) -> void:
	var sk := skill_at(i)
	if sk.is_empty() or S.me.d:
		return
	if ["damage", "dot", "aoetarget", "ground", "polymorph", "tame", "taunt", "charge"].has(sk.kind):
		for c in get_tree().get_nodes_in_group("controller"):
			c.face_target()
	Net.send("cast", {"skill": sk.id})
	slots[i].flash = 0.5

func _skill_tip(i: int, s: Slot, on: bool) -> void:
	if not on:
		Tooltip.hide_tip()
		return
	var sk := skill_at(i)
	if sk.is_empty():
		return
	Tooltip.show_at(ItemTip.skill(sk), s.get_global_rect().position + Vector2(0, -6), true)

# ---------------------------------------------------------------- tracker

func render_tracker() -> void:
	UI.clear(tracker_body)
	_q_timers.clear()
	for q in S.quests:
		var def := S.quest_def(q.id)
		if def.is_empty():
			continue
		var t := UI.label(def.title + ("  (fertig)" if q.status == "complete" else ""), 13, Color("ffd650"))
		t.add_theme_constant_override("outline_size", 3)
		t.add_theme_color_override("font_outline_color", Color.BLACK)
		tracker_body.add_child(t)
		if float(q.get("end", 0)) > 0.0 and q.status == "active":
			var tl := UI.label("", 12, Color("ffd650"))
			tracker_body.add_child(tl)
			_q_timers.append([tl, float(q.end)])
		var counts: Array = q.get("counts", [])
		for i in def.objectives.size():
			var o: Dictionary = def.objectives[i]
			var n: int = int(counts[i]) if i < counts.size() else 0
			var l := UI.label("- %s: %d/%d" % [o.text, n, int(o.count)], 12, Color("7ee27e") if n >= int(o.count) else Color("d8d8d8"))
			l.add_theme_constant_override("outline_size", 3)
			l.add_theme_color_override("font_outline_color", Color.BLACK)
			l.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
			l.custom_minimum_size.x = 220
			tracker_body.add_child(l)
	if S.quests.is_empty():
		tracker_body.add_child(UI.label("Keine aktiven Quests", 12, Color("999999")))

# ---------------------------------------------------------------- player menu

func open_player_menu(ent_id: int, pos: Vector2) -> void:
	var e: Ent = S.ents.get(ent_id)
	if e == null:
		return
	_menu_ent = ent_id
	player_menu.clear()
	var mine := false
	for m in S.party.members:
		if int(m.id) == ent_id:
			mine = true
	var leader: bool = int(S.party.leader) == S.id
	player_menu.add_item("%s  (Stufe %d)" % [e.n, e.l], 99)
	player_menu.set_item_disabled(0, true)
	if not mine:
		player_menu.add_item("👥 In Gruppe einladen", 1)
	if mine and leader:
		player_menu.add_item("⭐ Zum Anführer machen", 2)
		player_menu.add_item("🚪 Aus Gruppe werfen", 3)
	player_menu.add_item("💬 Anflüstern", 4)
	player_menu.position = Vector2i(get_window().position) + Vector2i(pos)
	player_menu.popup(Rect2i(Vector2i(get_viewport().get_mouse_position()) + Vector2i(get_window().position), Vector2i(10, 10)))

func _on_player_menu(id: int) -> void:
	match id:
		1: Net.send("invite", {"id": _menu_ent})
		2: Net.send("promote", {"id": _menu_ent})
		3: Net.send("kick", {"id": _menu_ent})
		4: open_chat("/p ")

# ---------------------------------------------------------------- floating text

func floating_text(screen: Vector2, text: String, cls: String, crit := false) -> void:
	if not Settings.v["damage_text"]:
		return
	var l := Label.new()
	l.text = text
	l.mouse_filter = Control.MOUSE_FILTER_IGNORE
	var col := Color.WHITE
	var sz := 22
	match cls:
		"out": col = Color("ffe14a") if not crit else Color("ffb020"); sz = 28 if crit else 22
		"in": col = Color("ff4a3a"); sz = 24
		"abs": col = Color("8ad0ff"); sz = 16
		_: col = Color("e8e8e8"); sz = 18
	l.add_theme_font_size_override("font_size", sz)
	l.add_theme_color_override("font_color", col)
	l.add_theme_constant_override("outline_size", 7)
	l.add_theme_color_override("font_outline_color", Color(0, 0, 0, 0.95))
	l.position = screen + Vector2(randf_range(-26, 26), 0) - Vector2(30, 14)
	fct_root.add_child(l)
	var tw := create_tween()
	tw.set_parallel(true)
	tw.tween_property(l, "position:y", l.position.y - 70.0, 1.15).set_ease(Tween.EASE_OUT)
	tw.tween_property(l, "modulate:a", 0.0, 0.45).set_delay(0.7)
	if crit:
		l.scale = Vector2(0.6, 0.6)
		l.pivot_offset = Vector2(30, 14)
		tw.tween_property(l, "scale", Vector2(1.25, 1.25), 0.18)
	tw.chain().tween_callback(l.queue_free)

# ---------------------------------------------------------------- per frame

func _portrait_emoji(e: Ent) -> String:
	if e.k == "player":
		return "👩" if e.lk.get("gender", "m") == "f" else "🧑"
	if MOB_EMOJI.has(e.m):
		return MOB_EMOJI[e.m]
	if NPC_EMOJI.has(e.m):
		return NPC_EMOJI[e.m]
	return "❓"

func _set_buffs(box: HBoxContainer, st: Array) -> void:
	var key := ",".join(PackedStringArray(st))
	if box.has_meta("k") and box.get_meta("k") == key:
		return
	box.set_meta("k", key)
	UI.clear(box)
	for s in st:
		if not STATUS_ICONS.has(s):
			continue
		var d: Array = STATUS_ICONS[s]
		var l := UI.label(d[0], 14)
		l.tooltip_text = d[1]
		l.mouse_filter = Control.MOUSE_FILTER_PASS
		box.add_child(l)

func update(dt: float) -> void:
	var me := S.me
	var ent: Ent = S.ents.get(S.id)
	# player frame
	var cc: Color = UI.CLASS_COLORS.get(S.self_info.get("class", ""), Color.WHITE)
	p_name.text = S.self_info.get("name", "")
	p_name.add_theme_color_override("font_color", cc)
	p_lvl.text = str(int(me.lvl))
	p_portrait.text = "👩" if S.self_info.get("look", {}).get("gender", "m") == "f" else "🧑"
	p_portrait_bg.add_theme_stylebox_override("panel", UI.box(Color(cc.r * 0.5, cc.g * 0.5, cc.b * 0.5, 0.9), cc, 2, 30))
	p_hp.set_state(me.hp, me.mh, "%d / %d" % [me.hp, me.mh])
	p_res.set_color(Color(S.cls.get("resColor", "#e24a3a")))
	p_res.set_state(me.res, me.rm, "%d / %d" % [me.res, me.rm])
	_set_buffs(p_buffs, ent.st if ent else [])
	if S.cls.get("id", "") == "warrior":
		p_stance.text = "🛡️ Verteidigungshaltung" if me.get("stance", "") == "defensive" else "⚔️ Kampfhaltung"
	else:
		p_stance.text = ""
	# target frame
	var t: Ent = S.ents.get(me.tg)
	if t:
		t_frame.visible = true
		t_name.text = t.n
		var tc := Color("ffe14a")
		if t.k == "player":
			tc = UI.CLASS_COLORS.get(t.c, Color.WHITE)
		elif t.k == "npc":
			tc = Color("ffe27a")
		t_name.add_theme_color_override("font_color", tc)
		t_lvl.text = "!" if t.k == "npc" else str(int(t.l))
		t_portrait.text = _portrait_emoji(t)
		t_hp.set_color(Color("e24a3a") if t.k == "mob" else Color("3ccf4e"))
		if t.d:
			t_hp.set_color(Color("777777"))
		t_hp.set_state(t.hp, t.mh, "" if t.k == "npc" else ("Tot" if t.d else "%d / %d" % [t.hp, t.mh]))
		t_tag.text = "Boss" if t.bs else ("Elite" if t.el else "")
		_set_buffs(t_buffs, t.st)
	else:
		t_frame.visible = false
	# xp
	xp_bar.set_state(me.xp, maxf(1.0, me.xn), ("Erfahrung: %d / %d" % [me.xp, me.xn]) if me.xn > 0 else "Höchststufe")
	if me.xn <= 0:
		xp_bar.value = 1.0
		xp_bar.maxv = 1.0
	# cast bar
	if S.cast != null:
		var k: float = (S.now() - S.cast.start) / S.cast.dur
		if k >= 1.0:
			S.cast = null
		else:
			cast_bar.visible = true
			cast_bar.set_state(k, 1.0, S.skill(S.cast.skill).get("name", ""))
	if S.cast == null:
		cast_bar.visible = false
	death_box.visible = me.d
	# action bar
	_update_actionbar()
	_update_social()
	# quest timers
	for qt in _q_timers:
		var left := maxi(0, ceili(qt[1] - S.now()))
		(qt[0] as Label).text = "⏱ %d Sek." % left
		(qt[0] as Label).add_theme_color_override("font_color", Color("ff5a4a") if left < 10 else Color("ffd650"))
	# clock + fps
	var gh := DayNight.hours()
	var hh := int(floor(gh))
	var mm := int(floor((gh - hh) * 60.0))
	var ck := ("☀ " if float(DayNight.info().elev) > -0.05 else "☾ ") + "%02d:%02d" % [hh, mm]
	if ck != _last_clock:
		_last_clock = ck
		clock_label.text = ck
	fps_label.text = ("%d FPS" % Engine.get_frames_per_second()) if Settings.v["show_fps"] else ""
	# timers for transient labels
	if err_t > 0.0:
		err_t -= dt
		err_label.modulate.a = clampf(err_t / 0.3, 0.0, 1.0)
	if toast_t > 0.0:
		toast_t -= dt
		toast_label.modulate.a = clampf(toast_t / 0.4, 0.0, 1.0)
	if banner_t > 0.0:
		banner_t -= dt
		banner.modulate.a = clampf(minf(banner_t / 0.8, (4.2 - banner_t) / 0.5), 0.0, 1.0)
	else:
		banner.modulate.a = 0.0
	if levelup_t > 0.0:
		levelup_t -= dt
		levelup.modulate.a = clampf(levelup_t / 0.8, 0.0, 1.0)
		levelup.scale = Vector2.ONE * (1.0 + (3.0 - levelup_t) * 0.04)
		levelup.pivot_offset = levelup.size / 2.0
	else:
		levelup.modulate.a = 0.0
	if invite_t > 0.0:
		invite_t -= dt
		if invite_t <= 0.0:
			invite_box.visible = false

func _update_actionbar() -> void:
	var t := S.now()
	for i in slots.size():
		var s: Slot = slots[i]
		var sk := skill_at(i)
		if sk.is_empty():
			if s.icon != "":
				s.icon = ""
				s.cost_text = ""
				s.cd_frac = 0.0
				s.border = Color(0.84, 0.7, 0.38, 0.2)
				s.queue_redraw()
			continue
		s.icon = sk.icon
		s.border = Color(0.84, 0.7, 0.38, 0.55)
		s.cost_text = str(int(sk.cost)) if int(sk.get("cost", 0)) > 0 else ""
		var rem := 0.0
		var dur := 1.0
		var cd = S.cd.get(sk.id)
		if cd != null and cd.end > t:
			rem = cd.end - t
			dur = cd.dur
		if not sk.get("basic", false) and S.gcd.end > t and S.gcd.end - t > rem:
			rem = S.gcd.end - t
			dur = S.gcd.dur
		s.cd_frac = clampf(rem / dur, 0.0, 1.0) if rem > 0.0 else 0.0
		s.cd_text = str(int(ceil(rem))) if rem > 1.0 and dur > 1.2 else ""
		s.no_res = S.me.res < int(sk.get("cost", 0))
		s.active = (sk.get("basic", false) and S.me.get("auto", false)) or (sk.kind == "stance" and S.me.get("stance", "") == sk.get("stance", "x"))
		s.queue_redraw()

func _update_social() -> void:
	# party frames
	var others: Array = []
	for m in S.party.members:
		if int(m.id) != S.id:
			others.append(m)
	var key := ""
	for m in others:
		key += str(m.id) + ","
	if key != _party_key:
		_party_key = key
		UI.clear(party_box)
		_party_bars.clear()
		for m in others:
			var id := int(m.id)
			var card := UI.card(func(): Net.send("target", {"id": id}); S.me.tg = id)
			card.add_theme_stylebox_override("panel", UI.box(Color(0.06, 0.055, 0.05, 0.78), UI.EDGE_SOFT, 1, 5, 4))
			var v := UI.vbox(2)
			var nm := UI.label(("⭐ " if int(m.id) == int(S.party.leader) else "") + m.n, 13, UI.CLASS_COLORS.get(m.c, Color.WHITE))
			nm.mouse_filter = Control.MOUSE_FILTER_IGNORE
			v.add_child(nm)
			var hp := Bar.new(150, 10, Color("3ccf4e"))
			var rs := Bar.new(150, 6, Color("4a8ae0"))
			v.add_child(hp)
			v.add_child(rs)
			card.add_child(v)
			card.set_meta("v", v)
			party_box.add_child(card)
			_party_bars[id] = [card, hp, rs]
	for st in S.party.states:
		var b = _party_bars.get(int(st.id))
		if b == null:
			continue
		b[1].set_state(st.hp, maxf(1, st.mh), "")
		b[2].set_state(st.res, maxf(1, st.rm), "")
		b[0].modulate.a = 0.45 if int(st.get("i", 0)) != (1 if S.instance != null else 0) else (0.6 if st.get("d", false) else 1.0)
	# pet frames
	var pet = S.me.get("pet")
	var pk := ""
	if pet is Dictionary:
		pk = "%s|%s|%s|%s" % [str(pet.get("id", 0)), str(pet.get("eid", 0)), str(pet.get("name", "")), str(pet.get("dead", false))]
	if pk != _pet_key:
		_pet_key = pk
		UI.clear(pet_box)
		_pet_bars.clear()
		if pet is Dictionary:
			if str(pet.get("name", "")) != "":
				var txt := "🐾 " + str(pet.name) + (" (gefallen)" if pet.get("dead", false) else ("" if pet.get("id", 0) else " (nicht gerufen)"))
				_add_pet_frame(txt, int(pet.get("id", 0)), "hp")
			if pet.get("eid", 0):
				_add_pet_frame("💧 Wasserelementar", int(pet.eid), "e")
	if pet is Dictionary:
		for pb in _pet_bars:
			if pb[1] == "hp":
				pb[0].set_state(pet.get("hp", 0), maxf(1, pet.get("mh", 1)), "")
			else:
				pb[0].set_state(pet.get("ehp", 0), maxf(1, pet.get("emh", 1)), "")

func _add_pet_frame(txt: String, id: int, kind: String) -> void:
	var card := UI.card(func(): if id != 0: Net.send("target", {"id": id}); S.me.tg = id)
	card.add_theme_stylebox_override("panel", UI.box(Color(0.06, 0.055, 0.05, 0.78), UI.EDGE_SOFT, 1, 5, 4))
	var v := UI.vbox(2)
	var nm := UI.label(txt, 12, Color("9fe89f"))
	nm.mouse_filter = Control.MOUSE_FILTER_IGNORE
	v.add_child(nm)
	var hp := Bar.new(150, 8, Color("3ccf4e"))
	v.add_child(hp)
	card.add_child(v)
	pet_box.add_child(card)
	_pet_bars.append([hp, kind])
