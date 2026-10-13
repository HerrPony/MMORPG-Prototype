extends Node3D
# Logs in, builds the world from the server's scenery and takes a screenshot:
#   godot --path . res://tests/world_shot.tscn -- ws://localhost:8084/ws user pass hour x z yaw_deg pitch_deg out.png [zone]
var world: GameWorld
var cam: Camera3D
var args: PackedStringArray

func _ready() -> void:
	args = OS.get_cmdline_user_args()
	DayNight.forced = float(args[3])
	Net.on("auth", func(m): S.user = m.user; Net.send("createChar", {"name": "Sh" + str(randi() % 100000).replace("0","a").replace("1","b").replace("2","c").replace("3","d").replace("4","e").replace("5","f").replace("6","g").replace("7","h").replace("8","i").replace("9","j"), "race": "human", "class": "warrior", "look": {"gender": "m", "hair": 0, "skin": 0}}))
	Net.on("authError", func(m): Net.send("register", {"user": args[1], "pass": args[2]}))
	Net.on("chars", func(m): if m.chars.size() > 0: Net.send("enterWorld", {"id": m.chars[0].id}))
	Net.on("world", _on_world)
	Net.connected.connect(func(): Net.send("login", {"user": args[1], "pass": args[2]}))
	Net.connect_to(args[0])

func _on_world(m: Dictionary) -> void:
	S.scenery = m.scenery
	S.defs = m.defs
	var t0 := Time.get_ticks_msec()
	world = GameWorld.new()
	add_child(world)
	world.build()
	print("world built in ", Time.get_ticks_msec() - t0, " ms; props=", WD.p_kind.size(), " buildings=", WD.buildings.size())
	cam = Camera3D.new()
	add_child(cam)
	cam.fov = 70
	var x := float(args[4]); var z := float(args[5])
	var y := WD.height_at(x, z)
	cam.position = Vector3(x, y + (float(args[10]) if args.size() > 10 else 3.0), z)
	cam.rotation = Vector3(deg_to_rad(float(args[7])), deg_to_rad(float(args[6])), 0)
	cam.make_current()
	world.set_zone(args[9] if args.size() > 9 else WD.biome_key(x, z))
	for i in 120:
		world.update(0.1)
	world.follow(cam.position)
	await get_tree().create_timer(2.0).timeout
	for i in 20:
		world.update(0.05)
		world.follow(cam.position)
		await get_tree().process_frame
	get_viewport().get_texture().get_image().save_png(args[8])
	print("saved ", args[8])
	get_tree().quit()
