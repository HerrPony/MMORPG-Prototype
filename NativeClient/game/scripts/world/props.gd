class_name Props
extends RefCounted
## Trees, rocks, bushes ...: blocky meshes instanced in 96 m chunks. The server decides where they stand (WD.p_*).

const CHUNK := 96.0
const SMALL := ["tuft", "flower", "reed", "mushroom", "bones", "fern", "ember", "drybush"]

static var _meshes: Dictionary = {}

static func hsl(h: float, s: float, l: float) -> Color:
	var q := l * (1.0 + s) if l < 0.5 else l + s - l * s
	var p := 2.0 * l - q
	var f := func(t: float) -> float:
		t = fposmod(t, 1.0)
		if t < 1.0 / 6.0: return p + (q - p) * 6.0 * t
		if t < 0.5: return q
		if t < 2.0 / 3.0: return p + (q - p) * (2.0 / 3.0 - t) * 6.0
		return p
	return Color(f.call(h + 1.0 / 3.0), f.call(h), f.call(h - 1.0 / 3.0))

static func leaf_tint(biome: int, r: float) -> Color:
	match biome:
		6: return hsl(0.33 + r * 0.05, 0.65, 0.2 + r * 0.09)
		7: return hsl(0.93 + r * 0.05, 0.6, 0.74 + r * 0.1)
		8: return hsl(0.1 + r * 0.03, 0.35, 0.28 + r * 0.06)
		9: return hsl(0.04, 0.2, 0.1 + r * 0.05)
		1: return hsl(0.31 + r * 0.03, 0.45, 0.17 + r * 0.08)
		2: return hsl(0.24 + r * 0.04, 0.5, 0.26 + r * 0.08)
		3: return hsl(0.22 + r * 0.04, 0.35, 0.16 + r * 0.06)
		4: return hsl(0.3, 0.3, 0.3)
		5: return hsl(0.18, 0.4, 0.3)
	return hsl(0.3, 0.4, 0.2)

static func tint_for(kind: String, biome: int, r: float) -> Color:
	match kind:
		"rock":
			match biome:
				4: return hsl(0.08, 0.05, 0.63 + r * 0.1)
				5: return hsl(0.1, 0.35, 0.5)
				8: return hsl(0.04, 0.55, 0.3 + r * 0.12)
				9: return hsl(0.02, 0.2, 0.08 + r * 0.06)
				7: return hsl(0.08, 0.1, 0.5)
			return hsl(0.08, 0.05, 0.38 + r * 0.15)
		"tuft", "fern":
			return hsl(0.14, 0.4, 0.45) if biome == 5 else hsl(0.26 + r * 0.07, 0.5, 0.25 + r * 0.15)
		"flower":
			return [Color(1, 0.42, 0.63), Color(1, 0.88, 0.29), Color(1, 1, 1), Color(0.63, 0.48, 1)][int(r * 3.99)]
	return leaf_tint(biome, r)

const T := 0.5   # alpha marking tintable parts
static func _lf(c := Color(1, 1, 1)) -> Color:
	return Color(c.r, c.g, c.b, T)

static func mesh_for(kind: String) -> ArrayMesh:
	if _meshes.has(kind):
		return _meshes[kind]
	var mb := MeshBuilder.new()
	var trunk := Models.C(0x5a3d22)
	var white := Color(1, 1, 1)
	match kind:
		"pine", "snowpine":
			mb.box(Vector3(0, 1.0, 0), Vector3(0.36, 2.0, 0.36), trunk)
			mb.frustum(Vector3(0, 1.3, 0), 2.6, 2.6, 1.5, 1.5, 1.1, _lf(white * 0.92))
			mb.frustum(Vector3(0, 2.2, 0), 2.0, 2.0, 1.0, 1.0, 1.0, _lf(white))
			mb.frustum(Vector3(0, 3.0, 0), 1.4, 1.4, 0.5, 0.5, 0.9, _lf(white * 1.05))
			mb.box(Vector3(0, 4.1, 0), Vector3(0.4, 0.5, 0.4), _lf(white * 1.1))
			if kind == "snowpine":
				mb.frustum(Vector3(0, 1.95, 0), 1.6, 1.6, 1.52, 1.52, 0.16, Models.C(0xf4f8ff))
				mb.frustum(Vector3(0, 2.9, 0), 1.2, 1.2, 1.02, 1.02, 0.16, Models.C(0xf4f8ff))
				mb.box(Vector3(0, 4.2, 0), Vector3(0.44, 0.3, 0.44), Models.C(0xf4f8ff))
		"oak", "sakura":
			mb.box(Vector3(0, 1.2, 0), Vector3(0.52, 2.4, 0.52), Models.C(0x6b4a2a))
			mb.box(Vector3(0, 3.3, 0), Vector3(2.7, 2.2, 2.7), _lf(white), Vector3.ZERO, 0.08)
			mb.box(Vector3(0.9, 2.6, 0.5), Vector3(1.7, 1.5, 1.7), _lf(white * 0.9), Vector3.ZERO, 0.08)
			mb.box(Vector3(-0.8, 3.9, -0.5), Vector3(1.8, 1.4, 1.8), _lf(white * 1.08), Vector3.ZERO, 0.08)
			mb.box(Vector3(0.2, 4.6, 0.1), Vector3(1.2, 0.8, 1.2), _lf(white * 1.12), Vector3.ZERO, 0.08)
		"dead", "charred":
			var tc := Models.C(0x3a3226) if kind == "dead" else Models.C(0x1c1816)
			mb.box(Vector3(0, 2.0, 0), Vector3(0.36, 4.0, 0.36), tc)
			mb.box(Vector3(0.7, 3.3, 0), Vector3(1.4, 0.16, 0.16), tc, Vector3(0, 0, 40))
			mb.box(Vector3(-0.6, 2.6, 0.1), Vector3(1.1, 0.14, 0.14), tc, Vector3(0, 0, -45))
			mb.box(Vector3(0.1, 1.7, 0.5), Vector3(0.12, 0.12, 1.0), tc, Vector3(30, 0, 0))
		"cactus":
			var g := Models.C(0x3a8a3a)
			mb.box(Vector3(0, 1.0, 0), Vector3(0.5, 2.0, 0.5), g)
			mb.box(Vector3(0.5, 1.2, 0), Vector3(0.6, 0.3, 0.3), g)
			mb.box(Vector3(0.75, 1.65, 0), Vector3(0.3, 0.8, 0.3), g)
			mb.box(Vector3(-0.45, 0.8, 0), Vector3(0.5, 0.28, 0.28), g)
			mb.box(Vector3(-0.65, 1.2, 0), Vector3(0.28, 0.7, 0.28), g)
			mb.box(Vector3(0, 2.05, 0), Vector3(0.2, 0.14, 0.2), Models.C(0xe05080))
		"palm":
			mb.box(Vector3(0, 0.9, 0), Vector3(0.34, 1.8, 0.34), Models.C(0x8a6a3a))
			mb.box(Vector3(0.2, 2.4, 0), Vector3(0.3, 1.4, 0.3), Models.C(0x8a6a3a), Vector3(0, 0, -12))
			mb.box(Vector3(0.45, 3.5, 0), Vector3(0.26, 1.0, 0.26), Models.C(0x8a6a3a), Vector3(0, 0, -20))
			for i in 6:
				var a := i * TAU / 6.0
				mb.box(Vector3(0.6 + cos(a) * 0.9, 3.9 - 0.25, sin(a) * 0.9), Vector3(1.9, 0.07, 0.5), _lf(white), Vector3(0, -a * 57.3, 18))
		"rock":
			mb.box(Vector3(0, 0.35, 0), Vector3(1.2, 0.7, 1.0), _lf(white), Vector3(0, 20, 0), 0.1)
			mb.box(Vector3(0.45, 0.25, 0.3), Vector3(0.7, 0.5, 0.7), _lf(white * 0.88), Vector3(10, 50, 0), 0.1)
			mb.box(Vector3(-0.35, 0.55, -0.2), Vector3(0.6, 0.6, 0.6), _lf(white * 1.1), Vector3(0, 70, 15), 0.1)
		"bush", "drybush":
			var col := _lf(white) if kind == "bush" else Models.C(0x8a6a38)
			mb.box(Vector3(0, 0.4, 0), Vector3(1.1, 0.8, 1.1), col, Vector3.ZERO, 0.08)
			mb.box(Vector3(0.35, 0.3, 0.25), Vector3(0.7, 0.6, 0.7), _lf(white * 0.9) if kind == "bush" else Models.C(0x7a5a30), Vector3(0, 30, 0), 0.08)
		"reed":
			for i in 4:
				mb.box(Vector3(-0.2 + i * 0.14, 0.7, (i % 2) * 0.1), Vector3(0.05, 1.4, 0.05), Models.C(0x6a7a3a), Vector3(0, 0, -6 + i * 4))
		"mushroom":
			mb.box(Vector3(0, 0.15, 0), Vector3(0.12, 0.3, 0.12), Models.C(0xe8e0d0))
			mb.box(Vector3(0, 0.35, 0), Vector3(0.36, 0.14, 0.36), Models.C(0xc04030))
			mb.box(Vector3(0.06, 0.43, 0.02), Vector3(0.08, 0.02, 0.08), Models.C(0xf4f0e0))
		"tuft":
			for i in 3:
				mb.box(Vector3(0, 0.2, 0), Vector3(0.05, 0.4, 0.2), _lf(white), Vector3(0, i * 60, 0))
		"flower":
			mb.box(Vector3(0, 0.17, 0), Vector3(0.04, 0.34, 0.04), Models.C(0x4a8a3a))
			mb.box(Vector3(0, 0.37, 0), Vector3(0.16, 0.12, 0.16), _lf(white), Vector3(0, 30, 0))
		"bones":
			mb.box(Vector3(0, 0.06, 0), Vector3(0.7, 0.08, 0.1), Models.C(0xe8e4d4), Vector3(0, 20, 0))
			mb.box(Vector3(0.1, 0.06, 0.2), Vector3(0.5, 0.08, 0.1), Models.C(0xe8e4d4), Vector3(0, -40, 0))
			mb.box(Vector3(-0.2, 0.12, -0.15), Vector3(0.2, 0.16, 0.2), Models.C(0xf0ecdc))
		"jtree":
			mb.box(Vector3(0, 2.8, 0), Vector3(0.7, 5.6, 0.7), Models.C(0x5a4026))
			mb.box(Vector3(0, 6.0, 0), Vector3(4.2, 1.6, 4.2), _lf(white * 0.9), Vector3.ZERO, 0.08)
			mb.box(Vector3(0.4, 7.2, -0.3), Vector3(3.2, 1.4, 3.2), _lf(white), Vector3.ZERO, 0.08)
			mb.box(Vector3(-0.3, 8.2, 0.3), Vector3(2.0, 1.0, 2.0), _lf(white * 1.1), Vector3.ZERO, 0.08)
			mb.box(Vector3(1.4, 4.4, 0.4), Vector3(0.14, 3.0, 0.14), Models.C(0x3a8a2a))
		"fern":
			for i in 5:
				mb.box(Vector3(cos(i * 1.26) * 0.5, 0.45, sin(i * 1.26) * 0.5), Vector3(0.9, 0.05, 0.3), _lf(white), Vector3(0, -i * 72, 28))
		"bamboo":
			for i in 6:
				var hgt := 6.0 + (i % 3) * 1.4
				var px := cos(i * 1.9) * 0.6
				var pz := sin(i * 1.9) * 0.6
				mb.box(Vector3(px, hgt / 2.0, pz), Vector3(0.12, hgt, 0.12), Models.C(0x8cb24a))
				for k in 4:
					mb.box(Vector3(px, 1.2 + k * 1.5, pz), Vector3(0.16, 0.05, 0.16), Models.C(0x6a9234))
		"mesapillar":
			var bands := [0xc8693f, 0xb4583a, 0xd88a50, 0x8a4628]
			var y := 0.0
			for i in 6:
				var w := 3.2 - i * 0.3
				mb.box(Vector3(0, y + 0.8, 0), Vector3(w, 1.6, w * 0.9), Models.C(bands[i % 4]), Vector3(0, i * 9, 0), 0.05)
				y += 1.6
		"obsidian":
			for i in 3:
				mb.frustum(Vector3((i - 1) * 0.55, 0, i * 0.2), 0.7, 0.7, 0.0, 0.0, 3.2 - i * 0.7, Models.C(0x14141c), Vector3(0, 0, (i - 1) * 8))
		"ember":
			mb.box(Vector3(0, 0.15, 0), Vector3(0.3, 0.3, 0.3), Models.C(0xff7a20), Vector3(30, 20, 0))
			mb.box(Vector3(0.35, 0.1, 0.2), Vector3(0.2, 0.2, 0.2), Models.C(0xffaa40), Vector3(10, 40, 0))
		_:
			mb.box(Vector3(0, 0.3, 0), Vector3(0.6, 0.6, 0.6), Models.C(0xa0a040))
	var m := mb.get_mesh()
	_meshes[kind] = m
	return m

static var _mat: ShaderMaterial = null
static func material() -> ShaderMaterial:
	if _mat == null:
		_mat = ShaderMaterial.new()
		_mat.shader = load("res://assets/shaders/foliage.gdshader")
	return _mat

## Builds all prop chunks below parent. density (0..1) thins out the small decorative props.
static func build(parent: Node3D, density: float, view_dist: float, shadows: bool) -> void:
	for c in parent.get_children():
		c.queue_free()
	var groups: Dictionary = {}   # Vector3i(kind, cx, cz) -> Array of indices
	var rnd := RandomNumberGenerator.new()
	rnd.seed = 99
	for i in WD.p_kind.size():
		var kind: String = WD.prop_kinds[WD.p_kind[i]]
		if density < 1.0 and kind in SMALL and (i * 2654435761) % 1000 / 1000.0 > density:
			continue
		var key := Vector3i(WD.p_kind[i], floori(WD.p_x[i] / CHUNK), floori(WD.p_z[i] / CHUNK))
		if not groups.has(key):
			groups[key] = []
		groups[key].append(i)
	var mat := material()
	for key in groups:
		var idxs: Array = groups[key]
		var kind: String = WD.prop_kinds[key.x]
		var mm := MultiMesh.new()
		mm.transform_format = MultiMesh.TRANSFORM_3D
		mm.use_custom_data = true
		mm.mesh = mesh_for(kind)
		mm.instance_count = idxs.size()
		var minp := Vector3(1e9, 1e9, 1e9)
		var maxp := Vector3(-1e9, -1e9, -1e9)
		for n in idxs.size():
			var i: int = idxs[n]
			var x: float = WD.p_x[i]
			var z: float = WD.p_z[i]
			var y := WD.height_at(x, z) - 0.1
			var s: float = WD.p_s[i]
			var basis := Basis(Vector3.UP, WD.p_ry[i]).scaled(Vector3(s, s * (0.9 + (n % 3) * 0.1), s))
			mm.set_instance_transform(n, Transform3D(basis, Vector3(x, y, z)))
			var r := rnd.randf()
			var t := tint_for(kind, WD.biome_at(x, z), r)
			mm.set_instance_custom_data(n, Color(t.r, t.g, t.b, 1.0))
			minp = minp.min(Vector3(x, y, z)); maxp = maxp.max(Vector3(x, y, z))
		mm.custom_aabb = AABB(minp - Vector3(6, 2, 6), (maxp - minp) + Vector3(12, 14, 12))
		var mi := MultiMeshInstance3D.new()
		mi.multimesh = mm
		mi.material_override = mat
		mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_ON if (shadows and not (kind in SMALL)) else GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
		mi.visibility_range_end = view_dist if not (kind in SMALL) else minf(view_dist, 70.0)
		parent.add_child(mi)
