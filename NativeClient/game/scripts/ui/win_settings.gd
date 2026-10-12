class_name WinSettings
extends GameWindow
## Settings: graphics (with presets), audio, controls and key bindings. Changes apply immediately and are saved.

var tabs: TabContainer
var rebinding := ""
var bind_buttons: Dictionary = {}

func _init() -> void:
	super._init("Einstellungen", Vector2(560, 520))
	win_id = "settings"
	default_pos = Vector2(340, 80)
	tabs = TabContainer.new()
	tabs.custom_minimum_size = Vector2(540, 470)
	body.add_child(tabs)
	_build_graphics()
	_build_audio()
	_build_controls()
	_build_keys()
	closed.connect(func(): rebinding = ""; Settings.save_file())

func _page(name: String) -> VBoxContainer:
	var sc := ScrollContainer.new()
	sc.name = name
	sc.horizontal_scroll_mode = ScrollContainer.SCROLL_MODE_DISABLED
	var v := UI.vbox(8)
	v.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	sc.add_child(v)
	tabs.add_child(sc)
	return v

func _row(parent: Control, label: String, ctrl: Control, hint := "") -> void:
	var h := UI.hbox(10)
	var l := UI.label(label, 14, UI.TEXT)
	l.custom_minimum_size.x = 210
	if hint != "":
		l.tooltip_text = hint
		l.mouse_filter = Control.MOUSE_FILTER_PASS
	h.add_child(l)
	ctrl.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	h.add_child(ctrl)
	parent.add_child(h)

func _check(parent: Control, label: String, key: String, hint := "") -> CheckButton:
	var c := CheckButton.new()
	c.button_pressed = bool(Settings.v[key])
	c.toggled.connect(func(on): Settings.set_value(key, on); Settings.apply_window())
	_row(parent, label, c, hint)
	return c

func _slider(parent: Control, label: String, key: String, lo: float, hi: float, step: float, fmt := "%.2f", hint := "") -> HSlider:
	var box := UI.hbox(8)
	var s := HSlider.new()
	s.min_value = lo; s.max_value = hi; s.step = step
	s.value = float(Settings.v[key])
	s.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	s.custom_minimum_size.y = 22
	var val := UI.label(fmt % s.value, 13, UI.GOLD, HORIZONTAL_ALIGNMENT_RIGHT)
	val.custom_minimum_size.x = 50
	s.value_changed.connect(func(v):
		val.text = fmt % v
		Settings.set_value(key, v if step < 1.0 or not (Settings.v[key] is int) else int(v)))
	box.add_child(s)
	box.add_child(val)
	_row(parent, label, box, hint)
	return s

func _option(parent: Control, label: String, key: String, items: Array, values: Array, hint := "") -> OptionButton:
	var o := OptionButton.new()
	for i in items.size():
		o.add_item(items[i], i)
	var cur = Settings.v[key]
	o.selected = maxi(0, values.find(cur))
	o.item_selected.connect(func(i):
		Settings.set_value(key, values[i])
		Settings.apply_window())
	_row(parent, label, o, hint)
	return o

func _build_graphics() -> void:
	var v := _page("Grafik")
	var preset := OptionButton.new()
	var names := ["Niedrig", "Mittel", "Hoch", "Ultra"]
	var ids := ["low", "medium", "high", "ultra"]
	for i in 4:
		preset.add_item(names[i], i)
	preset.selected = maxi(0, ids.find(Settings.v["preset"]))
	preset.item_selected.connect(func(i):
		Settings.apply_preset(ids[i])
		_refresh_pages())
	_row(v, "Grafik-Voreinstellung", preset, "Setzt viele Optionen auf einmal")
	v.add_child(HSeparator.new())
	_option(v, "Fenstermodus", "window_mode", ["Fenster", "Randlos", "Vollbild"], ["windowed", "borderless", "fullscreen"])
	_check(v, "Vertikale Synchronisation", "vsync")
	_option(v, "FPS-Limit", "fps_limit", ["Unbegrenzt", "30", "60", "90", "120", "144", "240"], [0, 30, 60, 90, 120, 144, 240])
	_slider(v, "Renderauflösung", "res_scale", 0.5, 1.0, 0.05, "%.2f", "Unter 1.0 rendert intern kleiner und skaliert hoch (schneller)")
	_check(v, "FSR-Hochskalierung", "fsr", "AMD FSR 1.0 beim Hochskalieren der Renderauflösung")
	_option(v, "Kantenglättung (MSAA)", "msaa", ["Aus", "2x", "4x", "8x"], [0, 1, 2, 3])
	_check(v, "FXAA", "fxaa")
	_check(v, "TAA", "taa", "Temporale Kantenglättung (kann leicht verschwimmen)")
	v.add_child(HSeparator.new())
	_option(v, "Schattenqualität", "shadows", ["Aus", "Niedrig", "Mittel", "Hoch", "Ultra"], [0, 1, 2, 3, 4])
	_slider(v, "Schattenreichweite", "shadow_distance", 30, 200, 10, "%d m")
	_slider(v, "Sichtweite", "view_distance", 100, 600, 10, "%d m")
	_slider(v, "Vegetationsdichte", "prop_density", 0.2, 1.0, 0.05, "%.2f", "Anzahl der Bäume, Büsche und Steine")
	_option(v, "Wasserqualität", "water_quality", ["Niedrig", "Mittel", "Hoch"], [0, 1, 2])
	_option(v, "Partikel", "particles", ["Aus", "Wenig", "Viel"], [0, 1, 2])
	v.add_child(HSeparator.new())
	_check(v, "Umgebungsverdeckung (SSAO)", "ssao")
	_check(v, "Indirekte Beleuchtung (SSIL)", "ssil")
	_check(v, "Bloom / Leuchten", "bloom")
	_check(v, "Nebel", "fog")
	_check(v, "Volumetrischer Nebel", "volumetric_fog", "Lichtstrahlen im Nebel (teuer)")
	_option(v, "Tonemapping", "tonemap", ["Linear", "Filmic", "ACES"], ["linear", "filmic", "aces"])
	_slider(v, "Helligkeit", "brightness", 0.6, 1.6, 0.05)
	_slider(v, "Kontrast", "contrast", 0.6, 1.6, 0.05)
	_slider(v, "Sättigung", "saturation", 0.0, 2.0, 0.05)
	_slider(v, "Sichtfeld (FOV)", "fov", 50, 110, 1, "%d°")
	v.add_child(HSeparator.new())
	_slider(v, "Oberflächen-Skalierung", "ui_scale", 0.7, 1.6, 0.05, "%.2f")
	_check(v, "Namensschilder", "nameplates")
	_check(v, "Schadenszahlen", "damage_text")
	_check(v, "FPS anzeigen", "show_fps")

func _build_audio() -> void:
	var v := _page("Audio")
	_slider(v, "Gesamtlautstärke", "vol_master", 0.0, 1.0, 0.01, "%.2f")
	_slider(v, "Musik", "vol_music", 0.0, 1.0, 0.01, "%.2f")
	_slider(v, "Effekte", "vol_sfx", 0.0, 1.0, 0.01, "%.2f")
	_slider(v, "Umgebung", "vol_ambient", 0.0, 1.0, 0.01, "%.2f")
	_check(v, "Musik aktiv", "music")

func _build_controls() -> void:
	var v := _page("Steuerung")
	_slider(v, "Maus-Empfindlichkeit", "mouse_sens", 0.2, 3.0, 0.05)
	_check(v, "Maus Y invertieren", "invert_y")
	_check(v, "Weiche Kamera", "camera_smooth", "Die Kamera folgt dem Charakter weich")
	_check(v, "Auto-Angriff beim Anvisieren", "autoattack_target")
	var info := UI.label("Linke Maustaste: auswählen · Rechte Maustaste: Kamera / interagieren / angreifen · Beide: vorwärts · Mausrad: Zoom\nGamepad: linker Stick laufen, rechter Stick Kamera, A springen, X interagieren, RB nächstes Ziel", 12, UI.MUTED)
	info.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	info.custom_minimum_size.x = 500
	v.add_child(info)

func _build_keys() -> void:
	var v := _page("Tasten")
	var hint := UI.label("Klicke auf eine Taste und drücke die neue Belegung. Esc bricht ab.", 12, UI.MUTED)
	v.add_child(hint)
	for a in Settings.ACTIONS:
		var b := Button.new()
		b.custom_minimum_size.x = 160
		bind_buttons[a] = b
		var act: String = a
		b.pressed.connect(func(): _start_rebind(act))
		_row(v, Settings.ACTIONS[a][0], b)
	v.add_child(UI.button("Alle zurücksetzen", func():
		Settings.binds.clear()
		Settings._setup_input()
		Settings.save_file()
		_refresh_keys()))
	_refresh_keys()

func _refresh_keys() -> void:
	for a in bind_buttons:
		bind_buttons[a].text = "…drücke eine Taste…" if rebinding == a else Settings.key_name(a)

func _start_rebind(a: String) -> void:
	rebinding = a
	_refresh_keys()

func _input(ev: InputEvent) -> void:
	if rebinding == "" or not visible:
		return
	if ev is InputEventKey and ev.pressed and not ev.echo:
		get_viewport().set_input_as_handled()
		if ev.keycode != KEY_ESCAPE:
			Settings.rebind(rebinding, ev.physical_keycode)
		rebinding = ""
		_refresh_keys()

func _refresh_pages() -> void:
	# rebuild the graphics page so that its controls show the preset values
	var cur := tabs.current_tab
	for c in tabs.get_children():
		tabs.remove_child(c)
		c.queue_free()
	bind_buttons.clear()
	_build_graphics()
	_build_audio()
	_build_controls()
	_build_keys()
	tabs.current_tab = cur
