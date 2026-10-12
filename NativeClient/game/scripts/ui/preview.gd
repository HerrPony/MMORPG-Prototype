class_name Preview
extends Node3D
## Stage for the character selection / creation screens: a platform in a pine forest with the rotating character.

var camera: Camera3D
var model: BlockyModel
var yaw := 0.3
var drag := false
var atk_t := 0.0
var center_offset := 0.0   # horizontal shift so that the character stays between the UI panels

func _ready() -> void:
	var env := Environment.new()
	env.background_mode = Environment.BG_COLOR
	env.background_color = Color("10202a")
	env.fog_enabled = true
	env.fog_light_color = Color("10202a")
	env.fog_density = 0.035
	env.ambient_light_source = Environment.AMBIENT_SOURCE_COLOR
	env.ambient_light_color = Color("7090b8")
	env.ambient_light_energy = 0.7
	env.tonemap_mode = Environment.TONE_MAPPER_FILMIC
	env.glow_enabled = true
	env.glow_intensity = 0.5
	var key := DirectionalLight3D.new()
	key.light_color = Color("fff0d0")
	key.light_energy = 1.6
	key.shadow_enabled = true
	key.rotation_degrees = Vector3(-40, 35, 0)
	add_child(key)
	var rim := DirectionalLight3D.new()
	rim.light_color = Color("6aa0ff")
	rim.light_energy = 0.8
	rim.rotation_degrees = Vector3(-25, 215, 0)
	add_child(rim)
	camera = Camera3D.new()
	camera.environment = env   # (a WorldEnvironment node would also affect the game world)
	camera.fov = 35
	camera.position = Vector3(0, 1.5, 4.6)
	add_child(camera)
	camera.look_at(Vector3(0, 1.0, 0))
	var ground := MeshInstance3D.new()
	var gm := CylinderMesh.new()
	gm.top_radius = 14; gm.bottom_radius = 14; gm.height = 0.2; gm.radial_segments = 24
	ground.mesh = gm
	ground.position.y = -0.15
	ground.material_override = _mat(Color("27452a"))
	add_child(ground)
	var plate := MeshInstance3D.new()
	var pm := CylinderMesh.new()
	pm.top_radius = 1.1; pm.bottom_radius = 1.2; pm.height = 0.12; pm.radial_segments = 12
	plate.mesh = pm
	plate.position.y = -0.04
	plate.material_override = _mat(Color("5a4a36"))
	add_child(plate)
	var trunk := _mat(Color("4a3322"))
	var leaf := _mat(Color("1f4a2a"))
	for i in 26:
		var a := float(i) / 26.0 * TAU
		var d := 7.0 + (i % 3) * 2.2
		var g := Node3D.new()
		var t := MeshInstance3D.new()
		var tm := CylinderMesh.new()
		tm.top_radius = 0.2; tm.bottom_radius = 0.3; tm.height = 3.0; tm.radial_segments = 5
		t.mesh = tm; t.position.y = 1.5; t.material_override = trunk
		g.add_child(t)
		for k in 3:
			var c := MeshInstance3D.new()
			var cm := CylinderMesh.new()
			cm.top_radius = 0.0; cm.bottom_radius = 1.8 - k * 0.45; cm.height = 2.4; cm.radial_segments = 6
			c.mesh = cm; c.position.y = 3.2 + k * 1.4; c.material_override = leaf
			g.add_child(c)
		g.position = Vector3(cos(a) * d, 0, sin(a) * d - 2.0)
		add_child(g)

func _mat(c: Color) -> StandardMaterial3D:
	var m := StandardMaterial3D.new()
	m.albedo_color = c
	m.roughness = 1.0
	return m

func set_style(cls: String, look: Dictionary) -> void:
	if model:
		model.queue_free()
	model = Models.make_player({"c": cls, "lk": look, "eq": {}, "id": 1})
	model.rotation.y = yaw
	add_child(model)

func clear() -> void:
	if model:
		model.queue_free()
		model = null

func _process(dt: float) -> void:
	if not visible:
		return
	var vs := get_viewport().get_visible_rect().size
	# shift the view so the model appears in the middle of the free space between the panels
	camera.h_offset = center_offset
	if model:
		if not drag:
			yaw += dt * 0.5
		model.rotation.y = yaw
		atk_t += dt
		if atk_t > 3.2:
			atk_t = 0.0
			model.attack()
		model.update(dt, {"speed": 0.0, "airborne": false, "dead": false})

func _input(ev: InputEvent) -> void:
	if not visible:
		return
	if ev is InputEventMouseMotion and drag:
		yaw += (ev as InputEventMouseMotion).relative.x * 0.01
	elif ev is InputEventMouseButton and (ev as InputEventMouseButton).button_index == MOUSE_BUTTON_LEFT:
		if not (ev as InputEventMouseButton).pressed:
			drag = false
		else:
			var h := get_viewport().gui_get_hovered_control()
			drag = h == null or h.mouse_filter == Control.MOUSE_FILTER_IGNORE
