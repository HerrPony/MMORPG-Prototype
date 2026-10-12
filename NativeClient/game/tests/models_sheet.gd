extends Node3D
# Renders every model into a contact sheet: godot --path . res://tests/models_sheet.tscn -- out.png
func _ready() -> void:
	var out := OS.get_cmdline_user_args()[0] if OS.get_cmdline_user_args().size() > 0 else "/tmp/sheet.png"
	var keys: Array = []
	keys.append_array(Models.HUMANOIDS.keys())
	keys.append_array(Models.QUADS.keys())
	keys.append_array(["spider", "cavespider", "spider_queen", "scorpion", "scorpion_king", "croc", "wisp", "frostspirit", "griffin"])
	var objs := ["pumpkin", "wheat", "spring", "tracks", "crystal", "lotus", "cactus", "urn", "portal"]
	var cols := 12
	var i := 0
	var rows := 0
	for k in keys:
		var e: Dictionary = {"k": "mob", "m": k, "id": i}
		var m := Models.make(e)
		m.position = Vector3((i % cols) * 3.2, 0, (i / cols) * 3.6)
		add_child(m)
		i += 1
	for c in ["warrior", "hunter", "mage"]:
		for g in ["m", "f"]:
			var m := Models.make({"k": "player", "c": c, "lk": {"gender": g, "hair": i % 6, "skin": i % 5}, "eq": {}, "id": i})
			m.position = Vector3((i % cols) * 3.2, 0, (i / cols) * 3.6)
			add_child(m)
			i += 1
	for o in objs:
		var m := Models.make({"k": "obj", "m": o})
		m.position = Vector3((i % cols) * 3.2, 0, (i / cols) * 3.6)
		add_child(m)
		i += 1
	for s in ["ship_cog", "ship_junk"]:
		var m := Models.make({"k": "ship", "m": s})
		m.position = Vector3((i % cols) * 3.2, 0, (i / cols) * 3.6)
		m.scale = Vector3.ONE * 0.3
		add_child(m)
		i += 1
	rows = (i - 1) / cols + 1
	var cam := Camera3D.new()
	add_child(cam)
	cam.position = Vector3(cols * 3.2 / 2.0 - 1.6, 14, rows * 3.6 + 10)
	cam.look_at(Vector3(cols * 3.2 / 2.0 - 1.6, 0.5, rows * 3.6 / 2.0))
	cam.fov = 60
	var l := DirectionalLight3D.new()
	add_child(l)
	l.rotation_degrees = Vector3(-50, 30, 0)
	var env := WorldEnvironment.new()
	var e2 := Environment.new()
	e2.background_mode = Environment.BG_COLOR
	e2.background_color = Color(0.45, 0.6, 0.75)
	e2.ambient_light_source = Environment.AMBIENT_SOURCE_COLOR
	e2.ambient_light_color = Color(0.7, 0.75, 0.85)
	e2.ambient_light_energy = 0.6
	env.environment = e2
	add_child(env)
	get_viewport().size = Vector2i(1900, 1000)
	await get_tree().create_timer(1.5).timeout
	for m in get_children():
		if m is BlockyModel:
			m.update(0.1, {"speed": 0.0})
	await get_tree().create_timer(0.3).timeout
	get_viewport().get_texture().get_image().save_png(out)
	print("saved ", out, " models=", i)
	get_tree().quit()
