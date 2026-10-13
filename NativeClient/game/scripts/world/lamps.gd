class_name Lamps
extends Node3D
## Street lamps and house lights that come on by themselves at night.
## Every lamp is a glowing mesh with a soft halo; the nearest few also get a real OmniLight3D (with shadows) from a pool,
## so the light really falls on the ground and the walls and the posts and houses cast shadows.

# lights in the pool and how many of them (the nearest) cast shadows, per shadow quality 0..4
const POOL := [6, 8, 12, 16, 24]
const SHADOWS := [0, 0, 3, 5, 8]
const REACH := 64.0          # lights further away than this are only glowing meshes
const POST_H := 3.3

var sources: Array = []      # {p, col, e, r, ph, halo}
var lights: Array = []       # OmniLight3D pool
var light_src: Array = []    # source index of each light (-1 = none)
var light_fade: Array = []   # distance fade of each light
var head_mat: ShaderMaterial
var halo_mat: ShaderMaterial
var halo_mm: MultiMeshInstance3D
var meshes: Node3D
var t := 0.0
var on := 0.0
var _timer := 0.0
var _cfg := -1

func build() -> void:
	_cfg = -1
	for c in get_children():
		c.queue_free()
	sources.clear()
	lights.clear()
	light_src.clear()
	light_fade.clear()
	head_mat = ShaderMaterial.new()
	head_mat.shader = load("res://assets/shaders/lamp_head.gdshader")
	halo_mat = ShaderMaterial.new()
	halo_mat.shader = load("res://assets/shaders/lamp_halo.gdshader")
	meshes = Node3D.new()
	meshes.name = "Posts"
	add_child(meshes)
	_build_posts()
	_collect_buildings()
	_build_halos()
	configure()

## (Re)creates the light pool for the current shadow quality.
func configure() -> void:
	var q := clampi(int(Settings.v["shadows"]), 0, 4)
	if q == _cfg:
		return
	_cfg = q
	for l in lights:
		(l as Node).queue_free()
	lights.clear()
	light_src.clear()
	light_fade.clear()
	for i in POOL[q]:
		var l := OmniLight3D.new()
		l.visible = false
		l.omni_attenuation = 1.5
		l.shadow_bias = 0.04
		l.shadow_normal_bias = 1.2
		l.shadow_blur = 1.4
		l.light_specular = 0.7
		add_child(l)
		lights.append(l)
		light_src.append(-1)
		light_fade.append(0.0)

static func style_colour(style: String) -> Color:
	match style:
		"jade":
			return Color(1.0, 0.55, 0.28)
		"obsidian":
			return Color(1.0, 0.4, 0.14)
		"bamboo":
			return Color(1.0, 0.7, 0.34)
	return Color(1.0, 0.72, 0.4)

func _add_source(p: Vector3, col: Color, e: float, r: float, halo: float) -> void:
	sources.append({"p": p, "col": col, "e": e, "r": r, "ph": randf() * TAU, "halo": halo})

## Lamp posts: iron pole, glass lantern and cap, merged per 128 m chunk.
func _build_posts() -> void:
	var chunks := {}
	for l in WD.lamps:
		var key := Vector2i(floori(l.x / 128.0), floori(l.z / 128.0))
		if not chunks.has(key):
			chunks[key] = [MeshBuilder.new(), MeshBuilder.new()]
		var main: MeshBuilder = chunks[key][0]
		var glow: MeshBuilder = chunks[key][1]
		var y := WD.height_at(l.x, l.z) - 0.15
		var style: String = l.get("s", "")
		var jade := style == "jade"
		main.box(Vector3(l.x, y + 0.2, l.z), Vector3(0.55, 0.4, 0.55), Models.C(0x5a5650), Vector3.ZERO, 0.03)
		main.box(Vector3(l.x, y + POST_H / 2.0, l.z), Vector3(0.16, POST_H, 0.16), Models.C(0x2e2a28), Vector3.ZERO, 0.02)
		main.box(Vector3(l.x, y + POST_H - 0.25, l.z), Vector3(0.34, 0.1, 0.34), Models.C(0x2e2a28), Vector3.ZERO, 0.0)
		glow.box(Vector3(l.x, y + POST_H + 0.3, l.z), Vector3(0.5, 0.6, 0.5), Models.C(0xff7a40 if jade else 0xffc068), Vector3.ZERO, 0.0)
		main.frustum(Vector3(l.x, y + POST_H + 0.58, l.z), 0.78, 0.78, 0.14, 0.14, 0.4, Models.C(0xc0302a if jade else 0x2a2420), Vector3.ZERO, 0.0)
		main.box(Vector3(l.x, y + POST_H + 1.04, l.z), Vector3(0.12, 0.12, 0.12), Models.C(0x2e2a28), Vector3.ZERO, 0.0)
		_add_source(Vector3(l.x, y + POST_H + 0.3, l.z), style_colour(style), 2.6, 16.0, 1.0)
	for key in chunks:
		var main2: MeshBuilder = chunks[key][0]
		var glow2: MeshBuilder = chunks[key][1]
		var mi := main2.instance(Buildings.material())
		mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_ON
		meshes.add_child(mi)
		var gi := glow2.instance(head_mat)
		gi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
		meshes.add_child(gi)

## Light sources of the buildings: lit windows in front of houses, lanterns, tower windows.
func _collect_buildings() -> void:
	for b in WD.buildings:
		var h: float = WD.height_at(b.x, b.z) - 0.3
		var style: String = b.get("s", "")
		match b.k:
			"house", "hut":
				var lift := 1.3 if b.k == "hut" else 0.0
				var col := Color(1.0, 0.42, 0.14) if style == "obsidian" else Color(1.0, 0.8, 0.48)
				_add_source(Vector3(b.x, h + lift + b.h * 0.62, b.z + b.d / 2.0 + 1.2), col, 1.7, 11.0, 0.0)
			"lantern":
				_add_source(Vector3(b.x, h + b.h + 0.3, b.z), Color(1.0, 0.62, 0.28), 2.4, 14.0, 1.0)
			"tower":
				_add_source(Vector3(b.x, h + b.h * 0.6, b.z + b.w / 2.0 + 1.0), Color(1.0, 0.78, 0.46), 1.8, 13.0, 0.0)

## One MultiMesh with a camera facing halo for every lamp.
func _build_halos() -> void:
	var idx: Array = []
	for i in sources.size():
		if sources[i].halo > 0.0:
			idx.append(i)
	var mm := MultiMesh.new()
	mm.transform_format = MultiMesh.TRANSFORM_3D
	mm.use_custom_data = true
	var q := QuadMesh.new()
	q.size = Vector2.ONE
	mm.mesh = q
	mm.instance_count = idx.size()
	for n in idx.size():
		var s: Dictionary = sources[idx[n]]
		mm.set_instance_transform(n, Transform3D(Basis.IDENTITY, s.p))
		mm.set_instance_custom_data(n, Color(s.ph / TAU, 3.4, s.halo, 0.0))
	halo_mm = MultiMeshInstance3D.new()
	halo_mm.multimesh = mm
	halo_mm.material_override = halo_mat
	halo_mm.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	halo_mm.extra_cull_margin = 400.0
	halo_mm.custom_aabb = AABB(Vector3(WD.x0 - 50.0, -20, WD.z0 - 50.0), Vector3((WD.nx - 1) * WD.cell + 100.0, 120, (WD.nz - 1) * WD.cell + 100.0))
	add_child(halo_mm)

## day: 0 night .. 1 day. Called every frame by the world.
func update(dt: float, cam_pos: Vector3, day: float) -> void:
	t += dt
	on = clampf((0.62 - day) / 0.4, 0.0, 1.0)
	head_mat.set_shader_parameter("night", 1.0 - day)
	halo_mat.set_shader_parameter("on", on)
	halo_mm.visible = on > 0.01
	if on <= 0.01:
		for l in lights:
			(l as OmniLight3D).visible = false
		_timer = 0.0
		return
	_timer -= dt
	if _timer <= 0.0:
		_timer = 0.25
		_assign(cam_pos)
	for i in lights.size():
		var s: int = light_src[i]
		var l: OmniLight3D = lights[i]
		if s < 0:
			l.visible = false
			continue
		var src: Dictionary = sources[s]
		var flick := 1.0 + 0.07 * sin(t * 9.0 + src.ph) + 0.04 * sin(t * 23.0 + src.ph * 2.0)
		l.light_energy = src.e * on * float(light_fade[i]) * flick
		l.visible = l.light_energy > 0.02

## Gives the pool lights to the sources nearest to the camera. A light keeps its source for as long as that is among
## the nearest ones, so nothing jumps around; the nearest few get the shadows.
func _assign(cam_pos: Vector3) -> void:
	var near: Array = []
	for i in sources.size():
		var d: float = (sources[i].p as Vector3).distance_to(cam_pos)
		if d < REACH + 8.0:
			near.append([d, i])
	near.sort_custom(func(a, b): return a[0] < b[0])
	var n := mini(near.size(), lights.size())
	var rank := {}      # source index -> [rank, distance]
	for r in n:
		rank[near[r][1]] = [r, near[r][0]]
	var taken := {}
	var free: Array = []
	for i in lights.size():
		var s: int = light_src[i]
		if s >= 0 and rank.has(s):
			taken[s] = true
		else:
			if s >= 0:
				light_src[i] = -1
			free.append(i)
	for r in n:
		var s: int = near[r][1]
		if taken.has(s) or free.is_empty():
			continue
		var i: int = free.pop_back()
		var src: Dictionary = sources[s]
		var l: OmniLight3D = lights[i]
		light_src[i] = s
		l.position = src.p
		l.light_color = src.col
		l.omni_range = src.r
		taken[s] = true
	var q := clampi(int(Settings.v["shadows"]), 0, 4)
	var compat := RenderingServer.get_current_rendering_method() == "gl_compatibility"
	for i in lights.size():
		var s: int = light_src[i]
		var l: OmniLight3D = lights[i]
		if s < 0:
			continue
		var info: Array = rank[s]
		light_fade[i] = 1.0 - DayNight.smooth(REACH - 22.0, REACH + 8.0, info[1])
		l.shadow_enabled = info[0] < SHADOWS[q]
		# dual paraboloid shadows are cheaper but only exist in the Forward+/Mobile renderers
		l.omni_shadow_mode = OmniLight3D.SHADOW_CUBE if (info[0] < 3 or compat) else OmniLight3D.SHADOW_DUAL_PARABOLOID
