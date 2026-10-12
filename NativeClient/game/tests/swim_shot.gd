extends Node3D
# Renders a swimming player on the water shader: godot --path . res://tests/swim_shot.tscn -- out.png [moving 0|1]
func _ready() -> void:
	var out := OS.get_cmdline_user_args()[0] if OS.get_cmdline_user_args().size() > 0 else "/tmp/swim.png"
	var moving := OS.get_cmdline_user_args().size() > 1 and OS.get_cmdline_user_args()[1] == "1"
	var env := Environment.new()
	env.background_mode = Environment.BG_COLOR
	env.background_color = Color(0.6, 0.78, 0.95)
	env.ambient_light_source = Environment.AMBIENT_SOURCE_COLOR
	env.ambient_light_color = Color(0.7, 0.75, 0.8)
	var we := WorldEnvironment.new()
	we.environment = env
	add_child(we)
	var sun := DirectionalLight3D.new()
	sun.rotation = Vector3(deg_to_rad(-50), deg_to_rad(30), 0)
	add_child(sun)
	var plane := MeshInstance3D.new()
	var pm := PlaneMesh.new()
	pm.size = Vector2(200, 200)
	pm.subdivide_width = 100
	pm.subdivide_depth = 100
	plane.mesh = pm
	var mat := ShaderMaterial.new()
	mat.shader = load("res://assets/shaders/water.gdshader")
	plane.material_override = mat
	plane.position.y = -0.12
	add_child(plane)
	var floor_mi := MeshInstance3D.new()
	var fm := PlaneMesh.new()
	fm.size = Vector2(200, 200)
	floor_mi.mesh = fm
	floor_mi.position.y = -6.0
	add_child(floor_mi)
	var m: BlockyModel = Models.make({"k": "player", "c": "warrior", "lk": {"gender": "m", "hair": 1, "skin": 1}, "eq": {}, "id": 1})
	m.position = Vector3(0, -1.1, 0)
	add_child(m)
	var cam := Camera3D.new()
	add_child(cam)
	cam.position = Vector3(2.5, 1.4, 5.0)
	cam.look_at(Vector3(0, -0.3, 0))
	cam.make_current()
	for i in 90:
		m.update(0.05, {"speed": 0.6 if moving else 0.0, "swim": true})
		await get_tree().process_frame
	get_viewport().get_texture().get_image().save_png(out)
	print("saved ", out)
	get_tree().quit()
