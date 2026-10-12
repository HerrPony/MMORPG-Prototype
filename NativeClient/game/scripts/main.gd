extends Node
## Entry point: wires the network, the 3D world, the HUD and the screens together.

var preview: Preview
var ui_layer: CanvasLayer
var screens: Screens
var hud: Hud
var windows: Windows
var plates: Control
var loading: Label
var world: GameWorld
var entities: Entities
var fx: FX
var controller: Controller
var cam: Camera3D
var dungeon: Dungeon = null
var game_root: Node3D
var spawned := false
var last_combat := -99.0
var zone_t := 0.0
var map_t := 0.0
var shot_pending := false
var world_dirty := false

func _ready() -> void:
	get_window().title = "Erlenwald Online"
	get_window().min_size = Vector2i(960, 540)
	get_tree().root.content_scale_factor = float(Settings.v["ui_scale"])
	Settings.changed.connect(func(): get_tree().root.content_scale_factor = float(Settings.v["ui_scale"]))
	preview = Preview.new()
	add_child(preview)
	cam = Camera3D.new()
	cam.far = 800.0
	cam.near = 0.1
	add_child(cam)
	var canvas := CanvasLayer.new()
	canvas.layer = 1
	add_child(canvas)
	# Controls directly below a CanvasLayer get no size, so everything hangs below a root that follows the viewport
	ui_layer = canvas
	var ui_root := Control.new()
	ui_root.name = "UIRoot"
	ui_root.mouse_filter = Control.MOUSE_FILTER_IGNORE
	canvas.add_child(ui_root)
	var fit := func(): ui_root.size = get_viewport().get_visible_rect().size
	get_viewport().size_changed.connect(fit)
	fit.call()
	ui_layer = canvas
	plates = Control.new()
	plates.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	plates.mouse_filter = Control.MOUSE_FILTER_IGNORE
	ui_root.add_child(plates)
	hud = Hud.new()
	hud.visible = false
	ui_root.add_child(hud)
	windows = Windows.new()
	windows.hud = hud
	ui_root.add_child(windows)
	screens = Screens.new()
	screens.preview = preview
	screens.theme = UI.theme()
	ui_root.add_child(screens)
	loading = UI.label("Welt wird aufgebaut …", 28, UI.GOLD, HORIZONTAL_ALIGNMENT_CENTER)
	loading.theme = UI.theme()
	loading.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	loading.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	loading.visible = false
	var lbg := ColorRect.new()
	lbg.color = Color(0.03, 0.04, 0.05, 1)
	lbg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbg.mouse_filter = Control.MOUSE_FILTER_IGNORE
	lbg.name = "LoadingBg"
	lbg.visible = false
	ui_root.add_child(lbg)
	ui_root.add_child(loading)
	preview.camera.current = true
	_wire()
	if OS.get_cmdline_user_args().size() > 0:
		var dev := DevHooks.new()
		dev.main = self
		add_child(dev)

func _wire() -> void:
	# chars must be registered before the screens decide what to show
	Net.on("chars", func(_m): if S.in_world: leave_world())
	Net.on("world", _on_world)
	Net.on("s", _on_snapshot)
	Net.on("correct", func(m): if controller: controller.correct(m.x, m.y, m.z))
	Net.on("instance", _on_instance)
	Net.on("dmg", func(m):
		if fx: fx.on_dmg(m)
		if m.id == S.id:
			last_combat = S.now()
			Sound.sfx("hurt")
		elif m.get("s", 0) == S.id:
			last_combat = S.now()
			Sound.sfx("hit"))
	Net.on("fx", func(m):
		if fx: fx.on_fx(m)
		if m.get("s", 0) != S.id:
			return
		var k: String = m.k
		if k == "levelup":
			Sound.sfx("levelup")
		elif not k.ends_with(":cast"):
			var st: String = S.skill(k).get("style", "")
			if st != "":
				Sound.sfx("melee" if st == "melee" else st))
	Net.on("questDone", func(_m): Sound.sfx("quest"))
	Net.on("skills", func(_m): Sound.sfx("learn"))
	Net.on("loot", func(_m): Sound.sfx("loot"))
	Net.on("err", func(_m): Sound.sfx("error", -8.0))
	Net.on("marks", func(m):
		S.marks = m.get("marks", {})
		if entities: entities.refresh_all_markers())
	Net.on("quests", func(_m): if entities: entities.refresh_objects())
	Net.on("worldChanged", func(_m):
		world_dirty = true
		_message("Ein Admin hat die Welt verändert. Bitte melde dich neu an, die Welt wird dann neu aufgebaut.", true))
	Net.on("kicked", func(m): _message(str(m.text), true))
	Net.closed.connect(func():
		if S.in_world:
			_message("Verbindung zum Server getrennt.", true))
	screens.enter_world.connect(func(id): Net.send("enterWorld", {"id": id}))
	screens.back_to_login.connect(func():
		Net.close()
		screens.show_screen("login"))
	windows.logout_requested.connect(func(): Net.send("logout"))
	windows.quit_requested.connect(func(): get_tree().quit())
	windows.editor_requested.connect(_open_editor)
	hud.toggle_window.connect(windows.toggle)

func ui_root_node() -> Control:
	return ui_layer.get_node("UIRoot")

func _message(text: String, to_login: bool) -> void:
	var d := AcceptDialog.new()
	d.title = "Erlenwald Online"
	d.dialog_text = text
	ui_root_node().add_child(d)
	d.popup_centered()
	d.confirmed.connect(d.queue_free)
	if to_login:
		leave_world()
		Net.close()
		screens.show_screen("login")

func _open_editor() -> void:
	var url := Net.url.replace("ws://", "http://").replace("wss://", "https://").trim_suffix("/ws")
	url += "/editor.html?x=%d&z=%d" % [int(S.local.x), int(S.local.z)]
	OS.shell_open(url)

# ---------------------------------------------------------------- entering and leaving the world

func _on_world(m: Dictionary) -> void:
	if m.has("now"):
		S.server_offset_ms = float(m.now) - Time.get_unix_time_from_system() * 1000.0
	S.reset_world()
	S.id = int(m.id)
	S.self_info = m.get("self", {})
	S.cls = m.get("cls", {})
	S.skill_ids = m.get("skills", [])
	S.defs = m.defs
	S.scenery = m.scenery
	S.disc = m.get("disc") if m.get("disc") is Array else []
	S.fp = m.get("fp") if m.get("fp") is Array else []
	screens.visible = false
	preview.visible = false
	ui_root_node().get_node("LoadingBg").visible = true
	loading.visible = true
	await get_tree().process_frame
	await get_tree().process_frame
	if world != null and world_dirty:
		_destroy_game()
	if world == null:
		_build_game()
	_leave_dungeon()
	entities.remove_all()
	fx.clear_all()
	spawned = false
	world.visible = true
	cam.make_current()
	map_setup()
	windows.reset()
	hud.visible = true
	hud.render_tracker()
	hud.add_chat("[color=#ffe27a]Tippe /hilfe für die Steuerung.[/color]")
	S.in_world = true
	ui_root_node().get_node("LoadingBg").visible = false
	loading.visible = false
	Sound.set_mode("world")

func _build_game() -> void:
	game_root = Node3D.new()
	game_root.name = "Game"
	add_child(game_root)
	world = GameWorld.new()
	game_root.add_child(world)
	world.build()
	entities = Entities.new()
	game_root.add_child(entities)
	entities.plates = plates
	entities.camera = cam
	fx = FX.new()
	game_root.add_child(fx)
	fx.ents = entities
	fx.hud = hud
	controller = Controller.new()
	add_child(controller)
	controller.camera = cam
	controller.world = world
	controller.entities = entities
	controller.action.connect(_on_action)
	controller.player_menu.connect(hud.open_player_menu)
	cam.fov = float(Settings.v["fov"])

func _destroy_game() -> void:
	world_dirty = false
	if controller:
		controller.queue_free()
		controller = null
	if game_root:
		game_root.queue_free()
		game_root = null
	world = null
	entities = null
	fx = null
	dungeon = null
	_map_ready = false
	S.ents.clear()

var _map_ready := false
func map_setup() -> void:
	if _map_ready:
		return
	_map_ready = true
	hud.minimap.setup(world)
	windows.map.setup(world)

func leave_world() -> void:
	if not S.in_world:
		return
	S.in_world = false
	Sound.set_mode("menu")
	_leave_dungeon()
	if entities:
		entities.remove_all()
		entities.ring.visible = false
	if fx:
		fx.clear_all()
	if world:
		world.visible = false
	S.me.tg = 0
	hud.visible = false
	windows.reset()
	Input.mouse_mode = Input.MOUSE_MODE_VISIBLE
	preview.camera.make_current()
	preview.visible = true

func _on_snapshot(m: Dictionary) -> void:
	if not S.in_world or entities == null:
		return
	entities.apply(m.e if m.get("e") is Array else [], m.get("lt") if m.get("lt") is Array else [])
	var tg: int = S.me.tg
	for k in m.me:
		S.me[k] = m.me[k]
	S.me.tg = tg
	S.party.states = m.get("pt") if m.get("pt") is Array else []
	var own: Ent = S.ents.get(S.id)
	if own and not spawned:
		spawned = true
		controller.spawn(own.x, own.y, own.z, own.r)
		world.set_zone(WD.biome_key(own.x, own.z))
		cam.make_current()

func _on_instance(m: Dictionary) -> void:
	if world == null:
		return
	if m.get("on", false):
		_leave_dungeon()
		S.instance = {"name": m.get("name", ""), "rooms": m.get("rooms", []), "origin": m.get("origin", [0, 0])}
		dungeon = Dungeon.new()
		game_root.add_child(dungeon)
		dungeon.build({"rooms": m.get("rooms", []), "origin": m.get("origin", [0, 0]), "style": m.get("style", "mine")})
		world.set_active(false)
		cam.environment = dungeon.env
		controller.dungeon = dungeon
		Sound.set_mode("dungeon")
	else:
		S.instance = null
		_leave_dungeon()
		Sound.set_mode("world")
	entities.remove_all()

func _leave_dungeon() -> void:
	if dungeon:
		dungeon.queue_free()
		dungeon = null
	if world:
		world.set_active(true)
	if cam:
		cam.environment = null
	if controller:
		controller.dungeon = null
	S.instance = null

# ---------------------------------------------------------------- input actions

func _on_action(a: String) -> void:
	match a:
		"toggle_bags": windows.toggle("bags")
		"toggle_quests": windows.toggle("quests")
		"toggle_char": windows.toggle("char")
		"toggle_spells": windows.toggle("spells")
		"toggle_map": windows.toggle("map")
		"toggle_party": windows.toggle("party")
		"menu":
			if windows.close_top():
				return
			if S.me.tg != 0:
				controller.set_target(0)
				return
			windows.toggle("menu")
		"chat": hud.open_chat("")
		"chat_cmd": hud.open_chat("/")
		"mute": Settings.set_value("music", not Settings.v["music"])
		"screenshot": _screenshot()
		"fullscreen":
			Settings.v["window_mode"] = "windowed" if Settings.v["window_mode"] != "windowed" else "borderless"
			Settings.apply_window()
		_:
			if a.begins_with("slot_"):
				hud.cast_slot(int(a.substr(5)))

func _screenshot() -> void:
	DirAccess.make_dir_recursive_absolute("user://screenshots")
	var img := get_viewport().get_texture().get_image()
	var path := "user://screenshots/erlenwald_%s.png" % Time.get_datetime_string_from_system().replace(":", "-")
	img.save_png(path)
	hud.toast("Screenshot gespeichert")

func _unhandled_key_input(ev: InputEvent) -> void:
	# F11 / F12 also work outside the game
	if S.in_world:
		return
	if ev.is_action_pressed("fullscreen"):
		_on_action("fullscreen")

# ---------------------------------------------------------------- frame loop

func _process(dt: float) -> void:
	dt = minf(dt, 0.05)
	if not S.in_world or world == null or not spawned:
		return
	controller.update(dt)
	entities.update(dt)
	fx.update(dt)
	if dungeon == null:
		world.update(dt)
		zone_t += dt
		if zone_t > 1.5:
			zone_t = 0.0
			var key := WD.biome_key(S.local.x, S.local.z)
			world.set_zone(key)
			Sound.zone = key
	Sound.combat = S.now() - last_combat < 6.0 and not S.me.d
	hud.update(dt)
	hud.minimap.update_view(controller.cam_yaw)
	windows.update_distance()
	if windows.quests.visible:
		map_t += dt
		if map_t > 1.0:
			map_t = 0.0
			windows.quests.refresh()
