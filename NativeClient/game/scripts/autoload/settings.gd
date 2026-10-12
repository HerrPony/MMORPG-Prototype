extends Node
## User settings (graphics, audio, controls) stored in user://settings.cfg. Applied at start and live.

signal changed

const PATH := "user://settings.cfg"

const DEFAULTS := {
	"server": "localhost:8080", "username": "",
	# graphics
	"preset": "high", "window_mode": "windowed", "vsync": true, "fps_limit": 0, "res_scale": 1.0, "msaa": 2, "taa": false, "fxaa": false, "fsr": false,
	"shadows": 3, "shadow_distance": 90, "view_distance": 300, "ssao": true, "ssil": false, "bloom": true, "fog": true, "volumetric_fog": false,
	"prop_density": 1.0, "particles": 2, "water_quality": 2, "tonemap": "filmic", "brightness": 1.0, "contrast": 1.0, "saturation": 1.05, "fov": 65.0,
	"ui_scale": 1.0, "nameplates": true, "damage_text": true, "show_fps": false,
	# audio
	"vol_master": 0.7, "vol_music": 0.7, "vol_sfx": 0.8, "vol_ambient": 0.6, "music": true,
	# controls
	"mouse_sens": 1.0, "invert_y": false, "camera_smooth": true, "autoattack_target": true,
}

var v: Dictionary = {}
var binds: Dictionary = {}

const ACTIONS := {
	"move_forward": ["Vorwärts", KEY_W], "move_back": ["Rückwärts", KEY_S], "move_left": ["Links", KEY_A], "move_right": ["Rechts", KEY_D],
	"jump": ["Springen", KEY_SPACE], "walk": ["Gehen (halten)", KEY_SHIFT], "autorun": ["Automatisch laufen", KEY_R],
	"target_next": ["Nächstes Ziel", KEY_TAB], "interact": ["Interagieren", KEY_F],
	"toggle_bags": ["Taschen", KEY_B], "toggle_quests": ["Questlog", KEY_L], "toggle_char": ["Charakter", KEY_C], "toggle_spells": ["Zauberbuch", KEY_P],
	"toggle_map": ["Weltkarte", KEY_M], "toggle_party": ["Gruppe", KEY_O], "chat": ["Chat öffnen", KEY_ENTER], "chat_cmd": ["Befehl eingeben", KEY_SLASH],
	"menu": ["Menü / Abbrechen", KEY_ESCAPE], "mute": ["Musik an/aus", KEY_N], "screenshot": ["Screenshot", KEY_F12], "fullscreen": ["Vollbild", KEY_F11],
}
const SLOT_KEYS := [KEY_1, KEY_2, KEY_3, KEY_4, KEY_5, KEY_6, KEY_7, KEY_8, KEY_9, KEY_0, KEY_MINUS, KEY_EQUAL, KEY_BRACKETLEFT, KEY_BRACKETRIGHT]

func _ready() -> void:
	v = DEFAULTS.duplicate()
	load_file()
	_setup_input()
	apply_window()

func load_file() -> void:
	var cf := ConfigFile.new()
	if cf.load(PATH) != OK:
		return
	for k in DEFAULTS:
		if cf.has_section_key("settings", k):
			v[k] = cf.get_value("settings", k)
	if cf.has_section("binds"):
		for a in cf.get_section_keys("binds"):
			binds[a] = int(cf.get_value("binds", a))

func save_file() -> void:
	var cf := ConfigFile.new()
	for k in v:
		cf.set_value("settings", k, v[k])
	for a in binds:
		cf.set_value("binds", a, binds[a])
	cf.save(PATH)

func set_value(key: String, val) -> void:
	v[key] = val
	changed.emit()

func _setup_input() -> void:
	for a in ACTIONS:
		if not InputMap.has_action(a):
			InputMap.add_action(a)
		InputMap.action_erase_events(a)
		var key: int = binds.get(a, ACTIONS[a][1])
		var ev := InputEventKey.new()
		ev.physical_keycode = key as Key
		InputMap.action_add_event(a, ev)
	for i in SLOT_KEYS.size():
		var name := "slot_%d" % i
		if not InputMap.has_action(name):
			InputMap.add_action(name)
		InputMap.action_erase_events(name)
		var ev := InputEventKey.new()
		ev.physical_keycode = SLOT_KEYS[i] as Key
		InputMap.action_add_event(name, ev)
	# gamepad
	var pad := {"jump": JOY_BUTTON_A, "interact": JOY_BUTTON_X, "target_next": JOY_BUTTON_RIGHT_SHOULDER, "menu": JOY_BUTTON_START}
	for a in pad:
		var jb := InputEventJoypadButton.new()
		jb.button_index = pad[a] as JoyButton
		InputMap.action_add_event(a, jb)

func rebind(action: String, key: int) -> void:
	binds[action] = key
	_setup_input()
	save_file()

func key_name(action: String) -> String:
	var key: int = binds.get(action, ACTIONS[action][1])
	return OS.get_keycode_string(key as Key)

func apply_window() -> void:
	match v["window_mode"]:
		"fullscreen": DisplayServer.window_set_mode(DisplayServer.WINDOW_MODE_EXCLUSIVE_FULLSCREEN)
		"borderless": DisplayServer.window_set_mode(DisplayServer.WINDOW_MODE_FULLSCREEN)
		_: DisplayServer.window_set_mode(DisplayServer.WINDOW_MODE_WINDOWED)
	DisplayServer.window_set_vsync_mode(DisplayServer.VSYNC_ENABLED if v["vsync"] else DisplayServer.VSYNC_DISABLED)
	Engine.max_fps = int(v["fps_limit"])
	if get_tree() and get_tree().root:
		get_tree().root.content_scale_factor = float(v["ui_scale"])

# graphics presets: they set many of the individual options at once
const PRESETS := {
	"low": {"msaa": 0, "fxaa": true, "shadows": 1, "shadow_distance": 50, "view_distance": 150, "ssao": false, "bloom": false, "volumetric_fog": false, "prop_density": 0.5, "particles": 0, "water_quality": 0, "res_scale": 0.8},
	"medium": {"msaa": 0, "fxaa": true, "shadows": 2, "shadow_distance": 70, "view_distance": 220, "ssao": false, "bloom": true, "volumetric_fog": false, "prop_density": 0.75, "particles": 1, "water_quality": 1, "res_scale": 1.0},
	"high": {"msaa": 2, "fxaa": false, "shadows": 3, "shadow_distance": 90, "view_distance": 300, "ssao": true, "bloom": true, "volumetric_fog": false, "prop_density": 1.0, "particles": 2, "water_quality": 2, "res_scale": 1.0},
	"ultra": {"msaa": 3, "fxaa": false, "shadows": 4, "shadow_distance": 140, "view_distance": 450, "ssao": true, "ssil": true, "bloom": true, "volumetric_fog": true, "prop_density": 1.0, "particles": 2, "water_quality": 2, "res_scale": 1.0},
}

func apply_preset(name: String) -> void:
	if not PRESETS.has(name):
		return
	for k in PRESETS[name]:
		v[k] = PRESETS[name][k]
	v["preset"] = name
	changed.emit()
