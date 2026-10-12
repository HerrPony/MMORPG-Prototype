class_name Dungeon
extends Node3D
## Procedural dungeon interior built from the rooms the server sends for an instance (mine, temple, fire cavern).

const WALL_H := 9.0
const STYLES := {
	"mine": {"bg": 0x07070b, "amb": 0x6a6a88, "rock": 0x4a443e, "floor": 0x3a3632, "ceil": 0x1c1a18, "wallH": [0.08, 0.1], "flame": 0xffaa33, "light": 0xff9a44, "player": 0xffd9a0},
	"temple": {"bg": 0x04100e, "amb": 0x5a8a7a, "rock": 0x4a6a5a, "floor": 0x2e463c, "ceil": 0x142a22, "wallH": [0.42, 0.14], "flame": 0x60ffb8, "light": 0x50ffc0, "player": 0xc8fff0},
	"fire": {"bg": 0x120504, "amb": 0x8a5a4a, "rock": 0x2c1e1e, "floor": 0x241618, "ceil": 0x120a0a, "wallH": [0.0, 0.1], "flame": 0xff5a1a, "light": 0xff5a20, "player": 0xffb888},
}
const FLOORS := {"throne": 0x2a2a3a, "web": 0x2e2e2a, "shrine": 0x24483c, "flooded": 0x28403a, "inferno": 0x1c1012, "obsidian": 0x16161e, "forge": 0x2a1c1c}

var env: Environment
var rects: Array = []
var flames: Array = []
var lights: Array = []
var anims: Array = []
var player_light: OmniLight3D
var rng_state := 77
var style: Dictionary
var rooms: Array = []
var _mats: Dictionary = {}

func _rnd() -> float:
	rng_state = (rng_state * 1664525 + 1013904223) & 0xFFFFFFFF
	return float(rng_state) / 4294967296.0

func mat(c: int, unshaded := false, alpha := 1.0, both := false) -> StandardMaterial3D:
	var key := "%d|%s|%f|%s" % [c, str(unshaded), alpha, str(both)]
	if _mats.has(key):
		return _mats[key]
	var m := StandardMaterial3D.new()
	m.albedo_color = Color(Models.C(c), alpha)
	m.roughness = 1.0
	if unshaded:
		m.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
	if alpha < 1.0:
		m.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
	if both:
		m.cull_mode = BaseMaterial3D.CULL_DISABLED
	_mats[key] = m
	return m

func add_mesh(mesh: Mesh, pos: Vector3, m: Material, rot := Vector3.ZERO) -> MeshInstance3D:
	var mi := MeshInstance3D.new()
	mi.mesh = mesh
	mi.position = pos
	mi.rotation = rot
	mi.material_override = m
	mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	add_child(mi)
	return mi

func box(sz: Vector3) -> BoxMesh:
	var b := BoxMesh.new()
	b.size = sz
	return b

func cyl(top: float, bottom: float, h: float, seg := 8) -> CylinderMesh:
	var c := CylinderMesh.new()
	c.top_radius = top
	c.bottom_radius = bottom
	c.height = h
	c.radial_segments = seg
	return c

func sphere(r: float, seg := 8) -> SphereMesh:
	var s := SphereMesh.new()
	s.radius = r
	s.height = r * 2.0
	s.radial_segments = seg
	s.rings = maxi(3, seg / 2)
	return s

func octa(r: float) -> SphereMesh:
	var s := SphereMesh.new()
	s.radius = r
	s.height = r * 2.0
	s.radial_segments = 4
	s.rings = 2
	return s

func disc(r: float, y_h := 0.04) -> CylinderMesh:
	return cyl(r, r, y_h, 20)

func inside(x: float, z: float, skip: int) -> bool:
	for i in rooms.size():
		if i == skip:
			continue
		var q: Dictionary = rooms[i]
		if x > q.x - q.w / 2.0 + 0.01 and x < q.x + q.w / 2.0 - 0.01 and z > q.z - q.d / 2.0 - 0.01 and z < q.z + q.d / 2.0 + 0.01:
			return true
	return false

func light(color: int, energy: float, rng_m: float, pos: Vector3) -> OmniLight3D:
	var l := OmniLight3D.new()
	l.light_color = Models.C(color)
	l.light_energy = energy
	l.omni_range = rng_m
	l.position = pos
	l.set_meta("base", energy)
	add_child(l)
	lights.append(l)
	return l

func torch(x: float, z: float) -> void:
	add_mesh(cyl(0.1, 0.12, 2.2, 5), Vector3(x, 1.1, z), mat(0x4a3020))
	flames.append(add_mesh(sphere(0.28, 6), Vector3(x, 2.4, z), mat(style.flame, true)))

func pillar(x: float, z: float, rad: float, color: int) -> void:
	add_mesh(cyl(rad * 0.9, rad, WALL_H, 8), Vector3(x, WALL_H / 2.0, z), mat(color))

func pool(x: float, z: float, rad: float, color: int, with_light := false) -> void:
	var m := add_mesh(disc(rad), Vector3(x, 0.04, z), mat(color, true))
	add_mesh(cyl(rad + 0.7, rad + 0.7, 0.03, 20), Vector3(x, 0.02, z), mat(0x1a1214))
	var base := Models.C(color)
	anims.append(func(t: float): (m.material_override as StandardMaterial3D).albedo_color = base * (0.85 + sin(t * 2.0 + x) * 0.15))
	if with_light:
		light(color, 1.0, 16.0, Vector3(x, 2, z))

func throne(q: Dictionary, seat_c: int, back_c: int, seat_sz: Vector3, back_sz: Vector3, off_z: float) -> Node3D:
	var th := Node3D.new()
	th.position = Vector3(q.x, 0, q.z - q.d / 2.0 + off_z)
	add_child(th)
	var s := MeshInstance3D.new(); s.mesh = box(seat_sz); s.position.y = seat_sz.y / 2.0; s.material_override = mat(seat_c); th.add_child(s)
	var b := MeshInstance3D.new(); b.mesh = box(back_sz); b.position = Vector3(0, back_sz.y / 2.0 + 1.0, -seat_sz.z / 2.0 + 0.1); b.material_override = mat(back_c); th.add_child(b)
	return th

func rune_ring(q: Dictionary, color: int, base_alpha: float) -> void:
	var r := MeshInstance3D.new()
	var tm := TorusMesh.new()
	tm.inner_radius = 8.0
	tm.outer_radius = 9.0
	tm.rings = 40
	tm.ring_segments = 3
	r.mesh = tm
	r.scale = Vector3(1, 0.02, 1)
	r.position = Vector3(q.x, 0.05, q.z + 2.0)
	r.material_override = mat(color, true, base_alpha, true)
	r.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	add_child(r)
	anims.append(func(t: float): (r.material_override as StandardMaterial3D).albedo_color.a = base_alpha - 0.15 + sin(t * 2.5) * 0.2)

func scatter(q: Dictionary, margin: float) -> Vector3:
	return Vector3(q.x + (_rnd() - 0.5) * (q.w - margin), 0, q.z + (_rnd() - 0.5) * (q.d - margin))

func bones(q: Dictionary, n: int, s: float) -> void:
	for k in n:
		var p := scatter(q, 4)
		var g := Node3D.new()
		g.position = p
		g.rotation.y = _rnd() * 6.0
		add_child(g)
		for j in 3:
			var b := MeshInstance3D.new()
			b.mesh = box(Vector3(0.5, 0.08, 0.1) * s)
			b.position = Vector3((j - 1) * 0.2, 0.06, (j - 1) * 0.15) * s
			b.rotation.y = j * 1.1
			b.material_override = mat(0xd8d4c4)
			g.add_child(b)
		var skull := MeshInstance3D.new()
		skull.mesh = box(Vector3(0.22, 0.2, 0.24) * s)
		skull.position = Vector3(0.1, 0.12, 0.1) * s
		skull.material_override = mat(0xe8e4d4)
		g.add_child(skull)

func build(info: Dictionary) -> void:
	style = STYLES.get(info.get("style", "mine"), STYLES.mine)
	env = Environment.new()
	env.background_mode = Environment.BG_COLOR
	env.background_color = Models.C(style.bg)
	env.ambient_light_source = Environment.AMBIENT_SOURCE_COLOR
	env.ambient_light_color = Models.C(style.amb)
	env.ambient_light_energy = 0.55
	env.fog_enabled = true
	env.fog_light_color = Models.C(style.bg)
	env.fog_density = 0.0
	env.fog_depth_begin = 8.0
	env.fog_depth_end = 55.0
	env.fog_mode = Environment.FOG_MODE_DEPTH
	env.tonemap_mode = Environment.TONE_MAPPER_FILMIC
	env.glow_enabled = true
	env.glow_intensity = 0.7
	var org: Array = info.origin
	for r in info.rooms:
		var q: Dictionary = r.duplicate()
		q.x = float(r.x) + float(org[0])
		q.z = float(r.z) + float(org[1])
		rooms.append(q)
	var wall_pieces: Array = []
	for i in rooms.size():
		var q: Dictionary = rooms[i]
		var fl := add_mesh(box(Vector3(q.w, 0.4, q.d)), Vector3(q.x, -0.2, q.z), mat(FLOORS.get(q.get("theme", ""), style.floor)))
		fl.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
		var ce := add_mesh(box(Vector3(q.w, 0.2, q.d)), Vector3(q.x, WALL_H + 0.1, q.z), mat(style.ceil))
		var x0: float = q.x - q.w / 2.0
		var x1: float = q.x + q.w / 2.0
		var z0: float = q.z - q.d / 2.0
		var z1: float = q.z + q.d / 2.0
		var x: float = x0 + 1.0
		while x < x1:
			for zz in [[z0, -1.0], [z1, 1.0]]:
				if not inside(x, zz[0] + zz[1] * 0.5, i) and not inside(x, zz[0] - zz[1] * 0.5, i):
					wall_pieces.append([x, zz[0] + zz[1] * 0.6, 2.1, 1.2])
			x += 2.0
		var z: float = z0 + 1.0
		while z < z1:
			for xx in [[x0, -1.0], [x1, 1.0]]:
				if not inside(xx[0] + xx[1] * 0.5, z, i) and not inside(xx[0] - xx[1] * 0.5, z, i):
					wall_pieces.append([xx[0] + xx[1] * 0.6, z, 1.2, 2.1])
			z += 2.0
	# walls as one MultiMesh with per-piece colour
	var mm := MultiMesh.new()
	mm.transform_format = MultiMesh.TRANSFORM_3D
	mm.use_colors = true
	var bm := BoxMesh.new()
	bm.size = Vector3.ONE
	mm.mesh = bm
	mm.instance_count = wall_pieces.size()
	for i in wall_pieces.size():
		var p: Array = wall_pieces[i]
		var h := WALL_H * (0.8 + _rnd() * 0.3)
		mm.set_instance_transform(i, Transform3D(Basis.from_scale(Vector3(p[2], h, p[3])), Vector3(p[0], h / 2.0, p[1])))
		var c := Color.from_hsv(style.wallH[0], style.wallH[1], 0.22 + _rnd() * 0.1 + 0.2)
		mm.set_instance_color(i, c)
	var wm := MultiMeshInstance3D.new()
	wm.multimesh = mm
	var wmat := StandardMaterial3D.new()
	wmat.vertex_color_use_as_albedo = true
	wmat.roughness = 1.0
	wm.material_override = wmat
	wm.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	add_child(wm)
	# rooms
	for i in rooms.size():
		var q: Dictionary = rooms[i]
		var n := maxi(2, int(q.d / 14.0))
		for k in n:
			var tz: float = q.z - q.d / 2.0 + (k + 0.5) * (q.d / n)
			if q.w > 8:
				torch(q.x - q.w / 2.0 + 1.2, tz)
				torch(q.x + q.w / 2.0 - 1.2, tz)
		if q.w > 8 and i % 2 == 0:
			light(style.light, 1.6, 38.0, Vector3(q.x, 5, q.z))
		_decorate(q)
	player_light = OmniLight3D.new()
	player_light.light_color = Models.C(style.player)
	player_light.light_energy = 1.2
	player_light.omni_range = 30.0
	add_child(player_light)
	for q in rooms:
		rects.append({"x0": q.x - q.w / 2.0 + 0.6, "x1": q.x + q.w / 2.0 - 0.6, "z0": q.z - q.d / 2.0 - 0.4, "z1": q.z + q.d / 2.0 + 0.4})

func _decorate(q: Dictionary) -> void:
	match q.get("theme", ""):
		"hall":
			for sx in [-1, 1]:
				for sz in [-1, 1]:
					pillar(q.x + sx * (q.w / 2.0 - 3.0), q.z + sz * (q.d / 2.0 - 3.0), 0.8, 0x6a645a)
		"mine":
			for k in 12:
				var p := scatter(q, 6)
				add_mesh(cyl(0.0, 0.5 + _rnd() * 0.5, 1.5 + _rnd(), 5), p + Vector3(0, 0.8, 0), mat(0x58504a))
			for k in 6:
				var p2 := scatter(q, 4)
				add_mesh(octa(0.5 + _rnd() * 0.5), p2 + Vector3(0, 0.8, 0), mat([0x4acfff, 0xff6aff, 0x6aff9a][k % 3], true))
			for k in 4:
				var p3 := scatter(q, 6)
				add_mesh(box(Vector3(1.6, 0.9, 1.1)), p3 + Vector3(0, 0.5, 0), mat(0x5a4030), Vector3(0, _rnd() * 3.0, 0))
		"web":
			for k in 14:
				var w := add_mesh(box(Vector3(4 + _rnd() * 4, 0.02, 4 + _rnd() * 3)), Vector3(q.x + (_rnd() - 0.5) * (q.w - 2), 3 + _rnd() * 4, q.z + (_rnd() - 0.5) * (q.d - 2)), mat(0xdddde8, true, 0.3, true), Vector3(_rnd() * 1.2, _rnd() * 3.0, _rnd() * 1.2))
			for k in 8:
				var p4 := scatter(q, 6)
				var co := add_mesh(sphere(0.7, 8), p4 + Vector3(0, 1, 0), mat(0xc8c8b8))
				co.scale.y = 1.6
			bones(q, 8, 1.0)
		"throne":
			throne(q, 0x6a5a2a, 0x8a7a3a, Vector3(3, 1.4, 3), Vector3(3, 5, 0.8), 4.0)
			rune_ring(q, 0xc03030, 0.6)
			for sx in [-1, 1]:
				for k in 3:
					pillar(q.x + sx * (q.w / 2.0 - 3.0), q.z - q.d / 2.0 + 8 + k * 14, 0.95, 0x4a4640)
					torch(q.x + sx * (q.w / 2.0 - 4.2), q.z - q.d / 2.0 + 8 + k * 14)
			light(0xff5a30, 2.0, 42.0, Vector3(q.x, 6, q.z))
			bones(q, 14, 1.2)
		"temple":
			for sx in [-1, 1]:
				for sz in [-1, 1]:
					pillar(q.x + sx * (q.w / 2.0 - 3.0), q.z + sz * (q.d / 2.0 - 3.0), 0.9, 0x78b8a0)
			for k in 4:
				var p5 := scatter(q, 8)
				var rn := add_mesh(cyl(1.5, 1.5, 0.03, 6), p5 + Vector3(0, 0.05, 0), mat(0x60ffb8, true, 0.5, true))
		"flooded":
			var w2 := add_mesh(box(Vector3(q.w - 1.5, 0.05, q.d - 1.5)), Vector3(q.x, 0.12, q.z), mat(0x2a9a9a, false, 0.55))
			(w2.material_override as StandardMaterial3D).roughness = 0.2
			anims.append(func(t: float): w2.position.y = 0.12 + sin(t * 1.3) * 0.03)
			for k in 7:
				var p6 := scatter(q, 6)
				add_mesh(cyl(0.7, 0.8, 1.5 + _rnd() * 4.0, 7), p6 + Vector3(0, 1, 0), mat(0x78a898))
			for k in 6:
				var p7 := scatter(q, 4)
				var cr := add_mesh(octa(0.5 + _rnd() * 0.4), p7 + Vector3(0, 0.9, 0), mat(0x50ffd0, true))
				anims.append(func(t: float): cr.rotation.y = t + k)
			for k in 10:
				var p8 := scatter(q, 4)
				add_mesh(cyl(0.6, 0.6, 0.04, 8), p8 + Vector3(0, 0.16, 0), mat(0x3f8a3a))
			light(0x40ffd0, 1.2, 30.0, Vector3(q.x, 4, q.z))
		"shrine":
			var th := throne(q, 0x2f9a78, 0x3fb890, Vector3(4, 1.2, 4), Vector3(5, 6, 0.8), 5.0)
			for sx in [-1, 1]:
				var cobra := MeshInstance3D.new()
				cobra.mesh = cyl(0.0, 0.8, 4.5, 6)
				cobra.position = Vector3(sx * 3.4, 2.3, -1.5)
				cobra.material_override = mat(0x2a8a60)
				th.add_child(cobra)
			rune_ring(q, 0x30ffa0, 0.6)
			for sx in [-1, 1]:
				for k in 3:
					pillar(q.x + sx * (q.w / 2.0 - 3.0), q.z - q.d / 2.0 + 8 + k * 14, 0.95, 0x78b8a0)
					torch(q.x + sx * (q.w / 2.0 - 4.2), q.z - q.d / 2.0 + 8 + k * 14)
			light(0x30ffa0, 2.0, 42.0, Vector3(q.x, 6, q.z))
		"forge":
			for sx in [-1, 1]:
				for sz in [-1, 1]:
					pillar(q.x + sx * (q.w / 2.0 - 3.0), q.z + sz * (q.d / 2.0 - 3.0), 0.9, 0x3a2a2a)
			for sx in [-1, 1]:
				add_mesh(box(Vector3(2, 0.12, q.d - 8)), Vector3(q.x + sx * (q.w / 2.0 - 5.5), 0.06, q.z), mat(0xff5a1a, true))
			for k in 3:
				var an := Node3D.new()
				an.position = Vector3(q.x + (k - 1) * 7, 0, q.z + (_rnd() - 0.5) * (q.d - 12))
				add_child(an)
				var base := MeshInstance3D.new(); base.mesh = box(Vector3(1.4, 1.1, 2)); base.position.y = 0.55; base.material_override = mat(0x2a2a30); an.add_child(base)
				var top := MeshInstance3D.new(); top.mesh = box(Vector3(2.2, 0.4, 1.0)); top.position.y = 1.3; top.material_override = mat(0x3a3a44); an.add_child(top)
			light(0xff5a20, 1.7, 34.0, Vector3(q.x, 4, q.z))
		"lava":
			pool(q.x - q.w / 2.0 + 5, q.z - q.d / 4.0, 4.2, 0xff5a1a, true)
			pool(q.x + q.w / 2.0 - 5, q.z + q.d / 4.0, 4.8, 0xff5a1a)
			pool(q.x + q.w / 2.0 - 6, q.z - q.d / 3.0, 3.2, 0xff7a20)
			for k in 8:
				var p9 := scatter(q, 6)
				add_mesh(cyl(0.0, 0.5, 2 + _rnd() * 2.4, 5), p9 + Vector3(0, 1.2, 0), mat(0x14141c))
		"obsidian":
			for k in 12:
				var p10 := scatter(q, 4)
				var cr2 := add_mesh(octa(0.7 + _rnd() * 0.8), p10 + Vector3(0, 1.2, 0), mat(0x14141c))
				cr2.scale.y = 2.0
			for k in 6:
				var p11 := scatter(q, 6)
				var cr3 := add_mesh(octa(0.4), p11 + Vector3(0, 0.7, 0), mat(0xff6a20, true))
				anims.append(func(t: float): cr3.rotation.y = t * 1.5 + k)
			for sx in [-1, 1]:
				for k in 2:
					pillar(q.x + sx * (q.w / 2.0 - 3.0), q.z - q.d / 4.0 + k * 16, 0.9, 0x1a1a24)
			pool(q.x, q.z + q.d / 2.0 - 5, 3.4, 0xff5a1a)
		"inferno":
			var th2 := throne(q, 0x1a1a24, 0x20202c, Vector3(3.4, 1.4, 3.4), Vector3(3.4, 5.5, 0.8), 4.0)
			for sx in [-1, 1]:
				var horn := MeshInstance3D.new()
				horn.mesh = cyl(0.0, 0.45, 3.2, 5)
				horn.position = Vector3(sx * 1.7, 7.2, -1.6)
				horn.rotation.z = -sx * 0.4
				horn.material_override = mat(0x14141c)
				th2.add_child(horn)
			rune_ring(q, 0xff3a10, 0.7)
			for sx in [-1, 1]:
				pool(q.x + sx * (q.w / 2.0 - 4.0), q.z - q.d / 6.0, 4.2, 0xff5a1a)
				pool(q.x + sx * (q.w / 2.0 - 4.0), q.z + q.d / 4.0, 3.6, 0xff6a20)
				for k in 3:
					pillar(q.x + sx * (q.w / 2.0 - 3.0), q.z - q.d / 2.0 + 8 + k * 14, 0.95, 0x1a1a24)
					torch(q.x + sx * (q.w / 2.0 - 4.2), q.z - q.d / 2.0 + 8 + k * 14)
			light(0xff4a20, 2.2, 46.0, Vector3(q.x, 6, q.z))

## Keeps a position inside the nearest walkable room rectangle.
func collide(x: float, z: float) -> Vector2:
	for q in rects:
		if x >= q.x0 and x <= q.x1 and z >= q.z0 and z <= q.z1:
			return Vector2(x, z)
	var best := Vector2(x, z)
	var bd := INF
	for q in rects:
		var cx := clampf(x, q.x0, q.x1)
		var cz := clampf(z, q.z0, q.z1)
		var dd := (cx - x) * (cx - x) + (cz - z) * (cz - z)
		if dd < bd:
			bd = dd
			best = Vector2(cx, cz)
	return best

func update(t: float, px: float, pz: float) -> void:
	player_light.position = Vector3(px, 3.5, pz)
	for i in flames.size():
		(flames[i] as Node3D).scale = Vector3.ONE * (0.85 + sin(t * 9.0 + i) * 0.15)
	for i in lights.size():
		var l: OmniLight3D = lights[i]
		l.light_energy = float(l.get_meta("base")) * (0.9 + sin(t * 7.0 + i * 2.0) * 0.1)
	for a in anims:
		a.call(t)
