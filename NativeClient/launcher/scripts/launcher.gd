extends Control
## Erlenwald launcher: finds, downloads and installs client updates by itself, then starts the game.

const DEFAULT_UPDATE_URL := "https://herrpony.github.io/MMORPG-Prototype/updates"
const CFG_PATH := "user://launcher.cfg"
const GOLD := Color("f0cf7a")
const MUTED := Color("a89a7e")
const TEXT := Color("ece3cf")
const RECHECK_SECONDS := 600.0

var cfg := {"update_url": "", "server": "", "auto_update": true, "close_on_launch": false, "root": ""}
var root := ""
var updater: Updater
var own_version := "1.0.2"
var busy := false
var state := "checking"       # checking | updating | ready | needs_install | error | offline
var cli := {}
var recheck_t := 0.0

var notes: RichTextLabel
var status_label: Label
var version_label: Label
var latest_label: Label
var bar: ProgressBar
var bar_label: Label
var play_btn: Button
var repair_btn: Button
var settings_btn: Button
var settings_win: Window
var in_url: LineEdit
var in_server: LineEdit
var cb_auto: CheckBox
var cb_close: CheckBox
var root_label: Label

func _ready() -> void:
	own_version = str(ProjectSettings.get_setting("application/config/version", "1.0.0"))
	DisplayServer.window_set_title("Erlenwald Online - Launcher")
	_parse_cli()
	_load_config()
	_resolve_root()
	updater = Updater.new(cfg.update_url, root)
	add_child(updater)
	updater.status.connect(func(t): _set_status(t))
	updater.progress.connect(_on_progress)
	if cli.has("cli"):
		await _cli_run()
		return
	_build_ui()
	_refresh_labels()
	_check_and_update.call_deferred()
	if cli.has("shot"):
		await get_tree().create_timer(float(cli.shot_delay)).timeout
		get_viewport().get_texture().get_image().save_png(str(cli.shot))
		print("saved ", cli.shot)
		get_tree().quit()

# ---------------------------------------------------------------- config

func _parse_cli() -> void:
	var a := OS.get_cmdline_user_args()
	var i := 0
	while i < a.size():
		match a[i]:
			"--cli-update": cli["cli"] = true
			"--repair": cli["repair"] = true
			"--self-update": cli["self_update"] = true
			"--no-launch": cli["no_launch"] = true
			"--shot": cli["shot"] = a[i + 1]; cli["shot_delay"] = float(a[i + 2]); i += 2
			"--url": cli["url"] = a[i + 1]; i += 1
			"--root": cli["root"] = a[i + 1]; i += 1
			"--server": cli["server"] = a[i + 1]; i += 1
		i += 1

func _defaults_file() -> Dictionary:
	# the installer (or a server operator) can ship launcher_defaults.json next to the launcher
	var p := OS.get_executable_path().get_base_dir().path_join("launcher_defaults.json")
	if FileAccess.file_exists(p):
		var d = JSON.parse_string(FileAccess.get_file_as_string(p))
		if d is Dictionary:
			return d
	return {}

func _load_config() -> void:
	var d := _defaults_file()
	cfg.update_url = str(d.get("update_url", DEFAULT_UPDATE_URL))
	cfg.server = str(d.get("server", "localhost:8080"))
	var cf := ConfigFile.new()
	if cf.load(CFG_PATH) == OK:
		for k in cfg:
			if cf.has_section_key("launcher", k):
				cfg[k] = cf.get_value("launcher", k)
	if cli.has("url"):
		cfg.update_url = cli.url
	if cli.has("root"):
		cfg.root = cli.root
	if cli.has("server"):
		cfg.server = cli.server

func _save_config() -> void:
	var cf := ConfigFile.new()
	for k in cfg:
		cf.set_value("launcher", k, cfg[k])
	cf.save(CFG_PATH)

func _resolve_root() -> void:
	if str(cfg.root) != "":
		root = str(cfg.root)
	else:
		root = OS.get_executable_path().get_base_dir()
		# started from the editor / a build folder: keep the game next to the project instead of inside the engine folder
		if OS.has_feature("editor"):
			root = ProjectSettings.globalize_path("res://").path_join("../build/install")
	root = root.simplify_path()
	DirAccess.make_dir_recursive_absolute(root)

# ---------------------------------------------------------------- UI

func _build_ui() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	var bg := ColorRect.new()
	bg.color = Color("0b1218")
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(bg)
	# a simple low-poly hero banner
	var banner := Control.new()
	banner.position = Vector2(0, 0)
	banner.size = Vector2(980, 170)
	banner.draw.connect(func(): _draw_banner(banner))
	add_child(banner)
	var title := Label.new()
	title.text = "ERLENWALD ONLINE"
	title.position = Vector2(34, 34)
	title.add_theme_font_size_override("font_size", 46)
	title.add_theme_color_override("font_color", GOLD)
	title.add_theme_constant_override("outline_size", 10)
	title.add_theme_color_override("font_outline_color", Color(0.12, 0.06, 0, 0.95))
	add_child(title)
	var sub := Label.new()
	sub.text = "Das Low-Poly-MMORPG"
	sub.position = Vector2(38, 100)
	sub.add_theme_font_size_override("font_size", 18)
	sub.add_theme_color_override("font_color", TEXT)
	sub.add_theme_constant_override("outline_size", 6)
	sub.add_theme_color_override("font_outline_color", Color(0, 0, 0, 0.9))
	add_child(sub)
	# news
	var np := PanelContainer.new()
	np.position = Vector2(24, 190)
	np.size = Vector2(610, 366)
	np.add_theme_stylebox_override("panel", _box(Color(0.05, 0.07, 0.09, 0.92), Color(0.84, 0.7, 0.38, 0.35), 1, 8, 14))
	add_child(np)
	var nv := VBoxContainer.new()
	np.add_child(nv)
	var nt := Label.new()
	nt.text = "Neuigkeiten"
	nt.add_theme_font_size_override("font_size", 20)
	nt.add_theme_color_override("font_color", GOLD)
	nv.add_child(nt)
	notes = RichTextLabel.new()
	notes.bbcode_enabled = true
	notes.size_flags_vertical = Control.SIZE_EXPAND_FILL
	notes.add_theme_font_size_override("normal_font_size", 15)
	notes.add_theme_color_override("default_color", TEXT)
	notes.text = "[color=#a89a7e]Lade Neuigkeiten …[/color]"
	nv.add_child(notes)
	# right column
	var rp := PanelContainer.new()
	rp.position = Vector2(654, 190)
	rp.size = Vector2(302, 366)
	rp.add_theme_stylebox_override("panel", _box(Color(0.05, 0.07, 0.09, 0.92), Color(0.84, 0.7, 0.38, 0.35), 1, 8, 14))
	add_child(rp)
	var rv := VBoxContainer.new()
	rv.add_theme_constant_override("separation", 8)
	rp.add_child(rv)
	version_label = _label("", 13, MUTED)
	latest_label = _label("", 13, MUTED)
	rv.add_child(version_label)
	rv.add_child(latest_label)
	var spacer := Control.new()
	spacer.size_flags_vertical = Control.SIZE_EXPAND_FILL
	rv.add_child(spacer)
	status_label = _label("", 14, TEXT)
	status_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	status_label.custom_minimum_size.y = 40
	rv.add_child(status_label)
	bar = ProgressBar.new()
	bar.custom_minimum_size.y = 18
	bar.show_percentage = false
	bar.add_theme_stylebox_override("background", _box(Color(0, 0, 0, 0.6), Color(0.84, 0.7, 0.38, 0.4), 1, 4, 0))
	bar.add_theme_stylebox_override("fill", _box(Color("cf9a33"), Color(0, 0, 0, 0), 0, 4, 0))
	rv.add_child(bar)
	bar_label = _label("", 12, MUTED)
	rv.add_child(bar_label)
	play_btn = Button.new()
	play_btn.text = "SPIELEN"
	play_btn.custom_minimum_size.y = 52
	play_btn.add_theme_font_size_override("font_size", 22)
	play_btn.add_theme_stylebox_override("normal", _box(Color("cf9a33"), Color("6a4a10"), 1, 6, 8))
	play_btn.add_theme_stylebox_override("hover", _box(Color("e4b04a"), Color("6a4a10"), 1, 6, 8))
	play_btn.add_theme_stylebox_override("pressed", _box(Color("b88626"), Color("6a4a10"), 1, 6, 8))
	play_btn.add_theme_stylebox_override("disabled", _box(Color(0.3, 0.27, 0.2), Color(0.1, 0.1, 0.1), 1, 6, 8))
	for c in ["font_color", "font_hover_color", "font_pressed_color"]:
		play_btn.add_theme_color_override(c, Color("2b1a05"))
	play_btn.add_theme_color_override("font_disabled_color", Color(0.55, 0.5, 0.4))
	play_btn.pressed.connect(_on_play)
	rv.add_child(play_btn)
	var row := HBoxContainer.new()
	row.add_theme_constant_override("separation", 6)
	rv.add_child(row)
	settings_btn = _small_button("Einstellungen", _open_settings)
	repair_btn = _small_button("Prüfen", func(): _check_and_update(true))
	repair_btn.tooltip_text = "Alle Dateien prüfen und beschädigte neu laden"
	row.add_child(settings_btn)
	row.add_child(repair_btn)
	row.add_child(_small_button("Beenden", func(): get_tree().quit()))
	_build_settings()

func _draw_banner(c: Control) -> void:
	# sky gradient and low-poly mountains and pines
	var top := Color("1c3b5a")
	var bot := Color("d79a5a")
	for i in 34:
		var t := float(i) / 33.0
		c.draw_rect(Rect2(0, i * 5.0, 980, 5.2), top.lerp(bot, t * t))
	var rng := RandomNumberGenerator.new()
	rng.seed = 7
	for layer in 3:
		var col := Color("2a3f55").lerp(Color("0f1c14"), float(layer) / 2.0)
		var pts := PackedVector2Array([Vector2(0, 170)])
		var x := 0.0
		while x <= 1000.0:
			pts.append(Vector2(x, 100.0 + layer * 22.0 - rng.randf() * (60.0 - layer * 14.0)))
			x += 38.0 + rng.randf() * 30.0
		pts.append(Vector2(1000, 170))
		c.draw_colored_polygon(pts, col)
	for i in 28:
		var px := rng.randf() * 980.0
		var h := 24.0 + rng.randf() * 34.0
		var base := 168.0
		for k in 3:
			var w := (h * 0.42) * (1.0 - k * 0.22)
			var y := base - k * h * 0.28
			c.draw_colored_polygon(PackedVector2Array([Vector2(px - w, y), Vector2(px + w, y), Vector2(px, y - h * 0.5)]), Color("0d2418"))
	c.draw_rect(Rect2(0, 0, 980, 170), Color(0, 0, 0, 0.18))

func _box(bg: Color, border: Color, bw: int, radius: int, pad: int) -> StyleBoxFlat:
	var s := StyleBoxFlat.new()
	s.bg_color = bg
	s.border_color = border
	s.set_border_width_all(bw)
	s.set_corner_radius_all(radius)
	s.set_content_margin_all(pad)
	return s

func _label(t: String, size_px: int, col: Color) -> Label:
	var l := Label.new()
	l.text = t
	l.add_theme_font_size_override("font_size", size_px)
	l.add_theme_color_override("font_color", col)
	return l

func _small_button(t: String, cb: Callable) -> Button:
	var b := Button.new()
	b.text = t
	b.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	b.add_theme_stylebox_override("normal", _box(Color(1, 1, 1, 0.06), Color(0.84, 0.7, 0.38, 0.4), 1, 4, 6))
	b.add_theme_stylebox_override("hover", _box(Color(1, 1, 1, 0.12), Color("f0cf7a"), 1, 4, 6))
	b.add_theme_stylebox_override("pressed", _box(Color(0.94, 0.81, 0.48, 0.2), Color("f0cf7a"), 1, 4, 6))
	b.add_theme_stylebox_override("disabled", _box(Color(1, 1, 1, 0.02), Color(0.84, 0.7, 0.38, 0.15), 1, 4, 6))
	b.pressed.connect(cb)
	return b

func _build_settings() -> void:
	settings_win = Window.new()
	settings_win.title = "Launcher-Einstellungen"
	settings_win.size = Vector2i(560, 360)
	settings_win.visible = false
	settings_win.exclusive = true
	settings_win.close_requested.connect(func(): settings_win.hide())
	add_child(settings_win)
	var bgc := ColorRect.new()
	bgc.color = Color("10181f")
	bgc.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	settings_win.add_child(bgc)
	var m := MarginContainer.new()
	m.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	for s in ["left", "right", "top", "bottom"]:
		m.add_theme_constant_override("margin_" + s, 18)
	settings_win.add_child(m)
	var v := VBoxContainer.new()
	v.add_theme_constant_override("separation", 8)
	m.add_child(v)
	v.add_child(_label("Update-Server (Adresse der Updates)", 13, MUTED))
	in_url = LineEdit.new()
	v.add_child(in_url)
	v.add_child(_label("Spielserver (wird dem Spiel vorgegeben)", 13, MUTED))
	in_server = LineEdit.new()
	v.add_child(in_server)
	cb_auto = CheckBox.new()
	cb_auto.text = "Updates automatisch herunterladen und installieren"
	v.add_child(cb_auto)
	cb_close = CheckBox.new()
	cb_close.text = "Launcher nach dem Spielstart schließen"
	v.add_child(cb_close)
	root_label = _label("", 12, MUTED)
	root_label.autowrap_mode = TextServer.AUTOWRAP_ARBITRARY
	v.add_child(root_label)
	var f := Control.new()
	f.size_flags_vertical = Control.SIZE_EXPAND_FILL
	v.add_child(f)
	var row := HBoxContainer.new()
	row.add_theme_constant_override("separation", 8)
	v.add_child(row)
	row.add_child(_small_button("Installationsordner öffnen", func(): OS.shell_open(root)))
	row.add_child(_small_button("Speichern", _save_settings))

func _open_settings() -> void:
	in_url.text = cfg.update_url
	in_server.text = cfg.server
	cb_auto.button_pressed = cfg.auto_update
	cb_close.button_pressed = cfg.close_on_launch
	root_label.text = "Installationsordner: " + root
	settings_win.popup_centered()

func _save_settings() -> void:
	cfg.update_url = in_url.text.strip_edges().trim_suffix("/")
	cfg.server = in_server.text.strip_edges()
	cfg.auto_update = cb_auto.button_pressed
	cfg.close_on_launch = cb_close.button_pressed
	_save_config()
	updater.base_url = cfg.update_url
	settings_win.hide()
	_check_and_update()

# ---------------------------------------------------------------- state

func _set_status(t: String) -> void:
	if status_label:
		status_label.text = t

func _on_progress(done: int, total: int) -> void:
	if bar:
		bar.max_value = maxf(1.0, total)
		bar.value = done
		bar_label.text = "%s / %s" % [_mb(done), _mb(total)]

func _mb(b: int) -> String:
	return "%.1f MB" % (float(b) / 1048576.0)

func _refresh_labels() -> void:
	updater.load_installed()
	var inst := updater.installed_version()
	version_label.text = "Installiert: " + (inst if inst != "" else "nicht installiert") + "    Launcher " + own_version
	latest_label.text = "Aktuell: " + (str(updater.manifest.version) if updater.manifest.has("version") else "…")

func _set_state(s: String) -> void:
	state = s
	var can_play := false
	var txt := "SPIELEN"
	match s:
		"ready": can_play = true
		"needs_install": can_play = not busy; txt = "INSTALLIEREN"
		"error": can_play = true; txt = "ERNEUT VERSUCHEN"
		"offline": can_play = FileAccess.file_exists(_game_exe()); txt = "SPIELEN"
		_: can_play = false
	play_btn.disabled = not can_play
	play_btn.text = txt
	repair_btn.disabled = busy
	settings_btn.disabled = busy

func _game_exe() -> String:
	var name := "ErlenwaldOnline.exe" if OS.get_name() == "Windows" else "ErlenwaldOnline.x86_64"
	if updater.manifest.has("launch"):
		name = str(updater.manifest.launch)
	else:
		var pth := root.path_join("game/launch.txt")
		if FileAccess.file_exists(pth):
			name = FileAccess.get_file_as_string(pth).strip_edges()
	return root.path_join("game").path_join(name)

func _notes_bb(m: Dictionary) -> String:
	var t := str(m.get("notes", ""))
	if t == "":
		return "[color=#a89a7e]Keine Neuigkeiten.[/color]"
	return t

func _process(dt: float) -> void:
	if cli.has("cli") or busy or state == "updating":
		return
	recheck_t += dt
	if recheck_t > RECHECK_SECONDS:
		recheck_t = 0.0
		_check_and_update()

func _check_and_update(deep := false) -> void:
	if busy:
		return
	busy = true
	recheck_t = 0.0
	bar.value = 0
	bar_label.text = ""
	_set_state("checking")
	var ok: bool = await updater.fetch_manifest()
	if not ok:
		busy = false
		var has_game := FileAccess.file_exists(_game_exe())
		_set_status("Keine Verbindung zum Update-Server.\n" + updater.last_error if has_game else "Update-Server nicht erreichbar:\n" + updater.last_error)
		_set_state("offline" if has_game else "error")
		_refresh_labels()
		return
	notes.text = _notes_bb(updater.manifest)
	# launcher update first
	if updater.launcher_update_available(own_version):
		_set_status("Launcher-Update %s wird geladen …" % updater.manifest.launcher_version)
		_set_state("updating")
		if await updater.stage_launcher_update():
			_set_status("Starte den Launcher neu …")
			if updater.apply_staged_launcher(true):
				await get_tree().create_timer(0.4).timeout
				get_tree().quit()
				return
		busy = false
		_set_status("Launcher-Update fehlgeschlagen: " + updater.last_error)
		_set_state("error")
		return
	updater.load_installed()
	var todo := updater.needs_update(deep)
	_refresh_labels()
	if todo.is_empty():
		busy = false
		if updater.installed_version() != str(updater.manifest.version):
			updater.adopt_existing()
			_refresh_labels()
		_set_status("Das Spiel ist auf dem neuesten Stand.")
		bar.value = bar.max_value
		_set_state("ready")
		return
	if updater.installed_version() == "":
		busy = false
		_set_status("Das Spiel ist noch nicht installiert (%s)." % _mb(_total(todo)))
		_set_state("needs_install")
		if cfg.auto_update:
			_run_update(deep)
		return
	if cfg.auto_update or deep:
		_run_update(deep)
	else:
		busy = false
		_set_status("Ein Update ist verfügbar (%s)." % _mb(_total(todo)))
		_set_state("needs_install")
		play_btn.text = "AKTUALISIEREN"

func _total(todo: Array) -> int:
	var t := 0
	for f in todo:
		t += int(f.size)
	return t

func _run_update(deep := false) -> void:
	busy = true
	_set_state("updating")
	var ok: bool = await updater.update_game(deep)
	busy = false
	_refresh_labels()
	if ok:
		_set_status("Installation abgeschlossen.")
		_set_state("ready")
	else:
		_set_status("Fehler: " + updater.last_error)
		_set_state("error")

func _on_play() -> void:
	match state:
		"needs_install":
			_run_update()
		"error":
			_check_and_update()
		_:
			_launch()

func _launch() -> void:
	var exe := _game_exe()
	if not FileAccess.file_exists(exe):
		_set_status("Spieldatei fehlt: " + exe)
		return
	var args: PackedStringArray = ["--", "--server", str(cfg.server)]
	var pid := OS.create_process(exe, args)
	if pid <= 0:
		_set_status("Das Spiel konnte nicht gestartet werden.")
		return
	_set_status("Spiel gestartet …")
	if cfg.close_on_launch:
		await get_tree().create_timer(1.0).timeout
		get_tree().quit()

# ---------------------------------------------------------------- command line mode (tests, scripted installs)

func _cli_run() -> void:
	updater.force_launcher_update = cli.has("self_update")
	updater.status.connect(func(t): print(t))
	var ok: bool = await updater.fetch_manifest()
	if not ok:
		print("RESULT: error manifest: ", updater.last_error)
		get_tree().quit(2)
		return
	if updater.launcher_update_available(own_version):
		print("launcher update available: ", updater.manifest.launcher_version)
		if cli.has("self_update"):
			if await updater.stage_launcher_update() and updater.apply_staged_launcher(false):
				print("RESULT: launcher update staged")
				get_tree().quit(0)
			else:
				print("RESULT: error ", updater.last_error)
				get_tree().quit(3)
			return
	ok = await updater.update_game(cli.has("repair"))
	print("RESULT: ", ("ok version=" + str(updater.installed_version())) if ok else "error " + updater.last_error)
	get_tree().quit(0 if ok else 3)
