class_name Buildings
extends RefCounted
## Town buildings (houses, towers, walls, keeps, pagodas ...) as blocky meshes merged per 128 m chunk.
## Windows and lamps live in a separate "glow" mesh that lights up at night.

const CHUNK := 128.0

const STYLES := {
	"timber": {"wall": 0x9a6a3a, "roof": 0x7a2f22, "trim": 0x5a3a1c},
	"plaster": {"wall": 0xe8dcc0, "roof": 0xb04a2a, "trim": 0x6a4a2c},
	"stilt": {"wall": 0x5a4a32, "roof": 0x6a7a3a, "trim": 0x3a2e1e},
	"lodge": {"wall": 0x6a4a30, "roof": 0xdde8f0, "trim": 0x3a2a1c},
	"adobe": {"wall": 0xd8b878, "roof": 0xc89858, "trim": 0x8a6a3a},
	"bamboo": {"wall": 0xcdbc7a, "roof": 0x9c8a3c, "trim": 0x7a6a30},
	"jade": {"wall": 0xb8362a, "roof": 0x2f9a78, "trim": 0xe8c860},
	"redrock": {"wall": 0xc8693f, "roof": 0xb4583a, "trim": 0x8a4628},
	"obsidian": {"wall": 0x2c2c38, "roof": 0x14141c, "trim": 0xff6a20},
}

## One chunk worth of geometry with the building origin as local frame.
class Ctx:
	var main: MeshBuilder
	var glow: MeshBuilder
	var water: MeshBuilder
	var ox := 0.0
	var oy := 0.0
	var oz := 0.0
	func box(cx: float, cy: float, cz: float, w: float, h: float, d: float, col: int, rot := Vector3.ZERO) -> void:
		main.box(Vector3(ox + cx, oy + cy, oz + cz), Vector3(w, h, d), Models.C(col), rot, 0.05)
	func fr(cx: float, by: float, cz: float, w0: float, d0: float, w1: float, d1: float, h: float, col: int, rot := Vector3.ZERO) -> void:
		main.frustum(Vector3(ox + cx, oy + by, oz + cz), w0, d0, w1, d1, h, Models.C(col), rot, 0.05)
	func lamp(cx: float, cy: float, cz: float, w: float, h: float, d: float, col: int) -> void:
		glow.box(Vector3(ox + cx, oy + cy, oz + cz), Vector3(w, h, d), Models.C(col), Vector3.ZERO, 0.0)
	func pool(cx: float, cy: float, cz: float, w: float, h: float, d: float, col: int, a: float) -> void:
		water.box(Vector3(ox + cx, oy + cy, oz + cz), Vector3(w, h, d), Models.C(col, a), Vector3.ZERO, 0.0)

static func build(parent: Node3D) -> void:
	for c in parent.get_children():
		c.queue_free()
	var chunks: Dictionary = {}
	for b in WD.buildings:
		var key := Vector2i(floori(b.x / CHUNK), floori(b.z / CHUNK))
		if not chunks.has(key):
			var c := Ctx.new()
			c.main = MeshBuilder.new(); c.glow = MeshBuilder.new(); c.water = MeshBuilder.new()
			chunks[key] = c
		var ctx: Ctx = chunks[key]
		ctx.ox = b.x; ctx.oz = b.z
		ctx.oy = 0.0 if b.k == "dock" else (WD.height_at(b.x, b.z) if b.k == "pad" else WD.height_at(b.x, b.z) - 0.3)
		match b.k:
			"house", "hut": _house(ctx, b)
			"tower": _tower(ctx, b)
			"wall": _wall(ctx, b)
			"keep": _keep(ctx, b)
			"fountain": _fountain(ctx, b)
			"pad": _pad(ctx, b)
			"tent": _tent(ctx, b)
			"cave": _cave(ctx, b)
			"pagoda": _pagoda(ctx, b)
			"gate": _gate(ctx, b)
			"pond": _pond(ctx, b)
			"lantern": _lantern(ctx, b)
			"dock": _dock(ctx, b)
	for key in chunks:
		var c: Ctx = chunks[key]
		var mi := c.main.instance(material())
		mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_ON
		parent.add_child(mi)
		if not c.glow.verts.is_empty():
			var gi := c.glow.instance(glow_material())
			gi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
			parent.add_child(gi)
		if not c.water.verts.is_empty():
			var wi := c.water.instance(MeshBuilder.alpha_material())
			wi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
			parent.add_child(wi)

static var _mat: ShaderMaterial = null
## Surface of the walls and roofs (bricks, planks, shingles) - the vertex colours decide the base colour.
static func material() -> ShaderMaterial:
	if _mat == null:
		_mat = ShaderMaterial.new()
		_mat.shader = load("res://assets/shaders/building.gdshader")
	return _mat

static var _glow_mat: ShaderMaterial = null
## Windows and lamps: dim glass by day, bright light at night (the world sets the "night" parameter).
static func glow_material() -> ShaderMaterial:
	if _glow_mat == null:
		var sh: Shader = load("res://assets/shaders/window_glow.gdshader")
		_glow_mat = ShaderMaterial.new()
		_glow_mat.shader = sh
	return _glow_mat

static func _house(c: Ctx, b: Dictionary) -> void:
	var st: Dictionary = STYLES.get(b.get("s", "timber"), STYLES["timber"])
	var lift := 1.3 if b.k == "hut" else 0.0
	var w: float = b.w
	var d: float = b.d
	var h: float = b.h
	c.box(0, lift + h / 2, 0, w, h, d, st.wall)
	for sx in [-1, 1]:
		for sz in [-1, 1]:
			c.box(sx * w / 2, lift + h / 2, sz * d / 2, 0.35, h, 0.35, st.trim)
	if b.k == "hut":
		for sx in [-1, 1]:
			for sz in [-1, 1]:
				c.box(sx * w / 2.3, (lift + 0.8) / 2 - 0.4, sz * d / 2.3, 0.4, lift + 0.8, 0.4, st.trim)
	match b.get("s", "timber"):
		"adobe":
			c.box(0, lift + h + 0.25, 0, w + 0.4, 0.5, d + 0.4, st.roof)
			c.box(0, lift + h + 0.9, 0, minf(w, d) * 0.5, 0.8, minf(w, d) * 0.5, 0xe8d8a8)
			c.box(0, lift + h + 1.5, 0, minf(w, d) * 0.3, 0.5, minf(w, d) * 0.3, 0xe8d8a8)
		"redrock":
			c.box(0, lift + h + 0.3, 0, w + 0.6, 0.6, d + 0.6, st.roof)
			c.box(0, lift + h + 0.75, 0, w + 0.9, 0.35, d + 0.9, st.trim)
			for i in range(-2, 3):
				c.box(i * w / 5.0, lift + h - 0.2, 0, 0.28, 0.28, d + 1.6, 0x5a3a22)
		"jade":
			c.fr(0, lift + h, 0, w + 3.2, d + 3.2, w * 0.5, d * 0.5, 1.5, st.roof)
			c.fr(0, lift + h + 1.5, 0, w * 0.6, d * 0.6, 0.6, 0.6, 1.5, st.roof)
			for sx in [-1, 1]:
				for sz in [-1, 1]:
					c.box(sx * (w / 2 + 1.5), lift + h + 0.55, sz * (d / 2 + 1.5), 0.4, 0.4, 0.4, st.trim, Vector3(0, 0, sx * 35))
			c.box(0, lift + h + 3.4, 0, 0.2, 0.9, 0.2, st.trim)
		"obsidian":
			c.fr(0, lift + h, 0, w + 1.2, d + 1.2, 0.5, 0.5, 3.6, st.roof)
			c.lamp(0, lift + h + 0.1, 0, w + 0.2, 0.2, d + 0.2, 0xff6a20)
		"bamboo":
			var m := maxf(w, d) + 2.4
			c.fr(0, lift + h, 0, m, m, 0.6, 0.6, 3.2, st.roof)
		_:
			c.fr(0, lift + h, 0, w + 1.4, d + 1.4, w * 0.35, d * 0.35, 2.4 + (0.5 if h > 5 else 0.0), st.roof)
			c.box(-w * 0.28, lift + h + 2.2, -d * 0.2, 0.7, 1.6, 0.7, 0x6a6a64)  # chimney
	c.box(0, lift + 1.1, d / 2 + 0.05, 1.3, 2.2, 0.2, 0x3a2410)
	var wc := 0xff7a30 if b.get("s", "") == "obsidian" else 0xffe9a0
	for sx in [-1, 1]:
		c.lamp(sx * w * 0.28, lift + h * 0.6, d / 2 + 0.07, 0.9, 0.9, 0.2, wc)
		c.lamp(sx * w * 0.28, lift + h * 0.6, -d / 2 - 0.07, 0.9, 0.9, 0.2, wc)
	if b.k == "hut":
		c.box(0, 0.5, d / 2 + 1.3, 2.2, 0.3, 2.0, st.trim)

static func _tower(c: Ctx, b: Dictionary) -> void:
	var w: float = b.w
	var h: float = b.h
	c.fr(0, 0, 0, w + 1.0, w + 1.0, w, w, h, 0x9a9690)
	c.box(0, h - 0.4, 0, w + 0.8, 0.8, w + 0.8, 0x8a8680)
	for i in 4:
		var a := i * PI / 2.0
		c.box(cos(a) * (w / 2 + 0.2), h + 0.4, sin(a) * (w / 2 + 0.2), 1.2, 0.8, 1.2, 0x8a8680)
	c.fr(0, h, 0, w + 1.8, w + 1.8, 0.2, 0.2, 4.0, 0x3a5a9a)
	c.lamp(0, h * 0.6, w / 2 + 0.55, 0.8, 1.4, 0.2, 0xffe9a0)
	c.lamp(w / 2 + 0.55, h * 0.8, 0, 0.2, 1.4, 0.8, 0xffe9a0)

static func _wall(c: Ctx, b: Dictionary) -> void:
	var w: float = b.w
	var d: float = b.d
	var h: float = b.h
	c.box(0, h / 2, 0, w, h, d, 0xaaa69e)
	var long := w > d
	var n := int((w if long else d) / 3.0)
	for i in n:
		var off := -((w if long else d) / 2.0) + 1.5 + i * 3.0
		c.box(off if long else 0.0, h + 0.6, 0.0 if long else off, 1.4 if long else d, 1.2, d if long else 1.4, 0xaaa69e)

static func _keep(c: Ctx, b: Dictionary) -> void:
	var s: String = b.get("s", "")
	var stone := 0xd8b878 if s == "adobe" else (0xc8693f if s == "redrock" else 0xcfc9bc)
	var roof := 0xc89858 if s == "adobe" else (0xa84e2e if s == "redrock" else 0x3a5a9a)
	var w: float = b.w
	var d: float = b.d
	var h: float = b.h
	c.box(0, h / 2, 0, w, h, d, stone)
	c.fr(0, h, 0, w + 1.0, d + 1.0, 1.0, 1.0, 5.0, roof)
	if s != "adobe" and s != "redrock":
		for sx in [-1, 1]:
			for sz in [-1, 1]:
				c.fr(sx * w / 2, 0, sz * d / 2, 4.6, 4.6, 4.2, 4.2, h + 4, stone)
				c.fr(sx * w / 2, h + 4, sz * d / 2, 5.4, 5.4, 0.3, 0.3, 4.0, 0xc8a24a)
	c.box(0, 3, d / 2 + 0.1, 4, 6, 0.4, 0x3a2410)
	for i in range(-2, 3):
		c.box(i * 4.0, h - 3, d / 2 + 0.12, 1.2, 3, 0.1, 0xc8a24a if i % 2 else 0x3a5a9a)
	for i in range(-2, 3):
		c.lamp(i * 5.0, h * 0.62, d / 2 + 0.1, 1.0, 1.8, 0.1, 0xffe9a0)

static func _fountain(c: Ctx, b: Dictionary) -> void:
	c.box(0, 0.5, 0, 8.4, 1.0, 8.4, 0xc4bfb4)
	c.box(0, 0.4, 0, 9.0, 0.8, 5.0, 0xc4bfb4)
	c.box(0, 0.4, 0, 5.0, 0.8, 9.0, 0xc4bfb4)
	c.pool(0, 1.0, 0, 7.4, 0.12, 7.4, 0x4aa8d8, 0.8)
	c.box(0, 1.9, 0, 1.2, 3.0, 1.2, 0xc4bfb4)
	c.box(0, 3.3, 0, 3.0, 0.5, 3.0, 0xc4bfb4)
	c.pool(0, 3.62, 0, 2.4, 0.1, 2.4, 0x4aa8d8, 0.8)

static func _pad(c: Ctx, b: Dictionary) -> void:
	var w: float = b.w
	c.box(0, 0.25, 0, w, 0.5, w, 0xa8a298)
	c.box(0, 0.25, 0, w + 0.6, 0.5, w - 1.0, 0xa8a298)
	c.box(0, 0.56, 0, w - 1.6, 0.1, w - 1.6, 0xc8a24a)
	c.box(0, 0.62, 0, w - 2.4, 0.1, w - 2.4, 0xa8a298)
	for a in [0.8, 2.4]:
		c.box(cos(a) * 4, 2.0, sin(a) * 4, 0.5, 4.0, 0.5, 0x6a4a2a)
		c.box(cos(a) * 4 + 0.7, 3.4, sin(a) * 4, 1.4, 0.9, 0.05, 0xc8a24a)
	c.box(0, 0.8, 0, 3.0, 0.5, 3.0, 0x8a6a3a)

static func _tent(c: Ctx, b: Dictionary) -> void:
	var cols := [0x8a3a2a, 0x6a5a3a, 0x4a4a5a]
	var col: int = cols[int(absf(b.x * 7.0 + b.z)) % 3]
	var w: float = b.w * 1.4
	c.fr(0, 0, 0, w, w, 0.3, 0.3, b.h, col)
	c.box(0, 1.0, b.w * 0.4, 1.2, 2.0, 0.1, 0x2a1a10, Vector3(-17, 0, 0))

static func _cave(c: Ctx, b: Dictionary) -> void:
	var kind: String = b.get("s", "mine")
	var rock := 0x24181c if kind == "fire" else (0x6a8a78 if kind == "temple" else 0x5a5650)
	var w: float = b.w
	var d: float = b.d
	var h: float = b.h
	c.box(0, h * 0.35, -d * 0.3, w * 1.2, h * 0.8, d * 1.0, rock, Vector3(0, 8, 0))
	c.box(w * 0.2, h * 0.65, -d * 0.4, w * 0.8, h * 0.6, d * 0.8, rock, Vector3(0, -20, 6))
	c.box(-w * 0.25, h * 0.5, -d * 0.3, w * 0.7, h * 0.7, d * 0.9, rock, Vector3(0, 25, -5))
	var front := d * 0.62
	c.box(0, h * 0.28, front, w * 0.34, h * 0.55, 0.3, 0x050505)
	for sx in [-1, 1]:
		c.box(sx * w * 0.2, h * 0.31, front + 0.1, 0.9, h * 0.62, 0.9, 0x7a746a)
	c.box(0, h * 0.62, front + 0.1, w * 0.5, 0.9, 1.0, 0x7a746a)
	c.lamp(0, 3.0, front + 1.5, 0.5, 0.7, 0.5, 0x60ffb0 if kind == "temple" else 0xff9040)
	if kind == "temple":
		for sx in [-1, 1]:
			c.box(sx * w * 0.3, h * 0.4, front + 0.6, 1.2, h * 0.8, 1.2, 0x9ac8b0)
		c.box(0, h * 0.82, front + 0.6, w * 0.8, 0.8, 1.2, 0x2f9a78)
	elif kind == "fire":
		c.lamp(0, h * 0.25, front + 0.2, w * 0.3, h * 0.45, 0.2, 0xff5a1a)
		for i in 5:
			c.fr((i - 2) * w * 0.16, h * 0.8, front - 0.5, 1.0, 1.0, 0.0, 0.0, 2.4, 0x14141c)

static func _pagoda(c: Ctx, b: Dictionary) -> void:
	var st: Dictionary = STYLES.get(b.get("s", "jade"), STYLES["jade"])
	var w: float = b.w
	var d: float = b.d
	c.box(0, 0.6, 0, w + 5, 1.2, d + 5, 0xb8a898)
	c.box(0, 0.3, d / 2 + 3.4, w * 0.5, 0.6, 3, 0xc8b8a8)
	var tiers := maxi(3, roundi(b.h / 4.4))
	var y := 1.2
	for i in tiers:
		var hgt := 4.6 if i == 0 else 3.4
		c.box(0, y + hgt / 2, 0, w * 0.82, hgt, d * 0.82, st.wall)
		for sx in [-1, 1]:
			for sz in [-1, 1]:
				c.box(sx * w * 0.41, y + hgt / 2, sz * d * 0.41, 0.6, hgt, 0.6, st.trim)
		c.lamp(0, y + hgt * 0.55, d * 0.41 + 0.07, w * 0.3, 1.1, 0.2, 0xffe9a0)
		c.lamp(0, y + hgt * 0.55, -d * 0.41 - 0.07, w * 0.3, 1.1, 0.2, 0xffe9a0)
		c.fr(0, y + hgt, 0, w + 3.2, d + 3.2, w * 0.4, d * 0.4, 1.8, st.roof)
		for sx in [-1, 1]:
			for sz in [-1, 1]:
				c.box(sx * (w / 2 + 1.5), y + hgt + 0.6, sz * (d / 2 + 1.5), 0.5, 0.5, 0.5, st.trim, Vector3(0, 0, sx * 40))
		y += hgt + 1.1
		w *= 0.8
		d *= 0.8
	c.box(0, y + 1.4, 0, 0.4, 3.6, 0.4, st.trim)
	c.lamp(0, y + 3.4, 0, 0.8, 0.8, 0.8, 0xe8c040)

static func _gate(c: Ctx, b: Dictionary) -> void:
	# torii: the beam spans the street (x for rot 0 faces z...)
	var along_x: bool = not (b.get("r", 0) == 1)
	var w: float = b.w
	var h: float = b.h
	var sgn := 1.0
	for s in [-1, 1]:
		if along_x:
			c.box(0, h / 2, s * w / 2.4, 1.0, h, 1.0, 0xc0302a)
		else:
			c.box(s * w / 2.4, h / 2, 0, 1.0, h, 1.0, 0xc0302a)
	if along_x:
		c.box(0, h + 0.3, 0, 1.0, 0.7, w + 4, 0x201a1a)
		c.box(0, h - 1.4, 0, 0.8, 0.5, w + 1.2, 0xc0302a)
		for s in [-1, 1]:
			c.box(0, h + 0.7, s * (w / 2 + 2.2), 1.0, 0.5, 1.4, 0x201a1a, Vector3(s * 20, 0, 0))
	else:
		c.box(0, h + 0.3, 0, w + 4, 0.7, 1.0, 0x201a1a)
		c.box(0, h - 1.4, 0, w + 1.2, 0.5, 0.8, 0xc0302a)
		for s in [-1, 1]:
			c.box(s * (w / 2 + 2.2), h + 0.7, 0, 1.4, 0.5, 1.0, 0x201a1a, Vector3(0, 0, -s * 20))

static func _pond(c: Ctx, b: Dictionary) -> void:
	var w: float = b.w
	var d: float = b.d
	c.box(0, 0.2, 0, w + 1.2, 0.6, d + 1.2, 0xb8a898)
	c.pool(0, 0.55, 0, w, 0.2, d, 0x3aa8b8, 0.85)
	for i in 7:
		c.box(((i * 37) % 17 - 8) / 8.0 * w * 0.4, 0.68, ((i * 53) % 13 - 6) / 6.0 * d * 0.38, 1.0, 0.06, 1.0, 0x3f8a3a)
	c.box(0, 1.2, 0, 3, 0.4, d + 4, 0xc0302a)
	c.box(0, 1.7, 0, 3, 0.4, d * 0.5, 0xc0302a)
	for s in [-1, 1]:
		c.box(s * 1.4, 1.8, 0, 0.2, 0.8, d + 4, 0xe8c860)

static func _lantern(c: Ctx, b: Dictionary) -> void:
	var h: float = b.h
	c.box(0, h / 2, 0, 0.3, h, 0.3, 0x3a2a22)
	c.lamp(0, h + 0.3, 0, 0.7, 0.8, 0.7, 0xffb040)
	c.fr(0, h + 0.7, 0, 0.95, 0.95, 0.1, 0.1, 0.55, 0xc0302a)

static func _dock(c: Ctx, b: Dictionary) -> void:
	var w: float = b.w
	var d: float = b.d
	c.box(0, 1.03, 0, w, 0.35, d, 0x8a6a40)
	var x := -w / 2 + 1.5
	while x <= w / 2:
		for z in [-1, 1]:
			c.box(x, -0.6, z * (d / 2 - 0.3), 0.5, 4.2, 0.5, 0x4a3220)
		x += 4.0

## The cabin at the world origin (server sends its wall boxes).
static func build_cabin(parent: Node3D) -> void:
	var mb := MeshBuilder.new()
	var h0 := WD.height_at(0, 0)
	for b in WD.house:
		mb.box(Vector3(b.x, h0 + b.h / 2, b.z), Vector3(b.w, b.h, b.d), Models.C(0x8a5a2e), Vector3.ZERO, 0.04)
	mb.box(Vector3(0, h0 + 0.1, 0), Vector3(10.6, 0.5, 8.6), Models.C(0x8a8a84))
	mb.box(Vector3(0, h0 + 0.36, 0), Vector3(9.6, 0.1, 7.6), Models.C(0x6b4a28))
	mb.frustum(Vector3(0, h0 + 3.0, 0), 11.6, 9.6, 1.4, 1.4, 2.6, Models.C(0x8a2f22))
	mb.box(Vector3(-2.8, h0 + 4.4, -1.6), Vector3(0.9, 2.6, 0.9), Models.C(0x8a8a84))
	mb.box(Vector3(0, h0 + 2.5, 4), Vector3(2.1, 0.2, 0.5), Models.C(0x5a3a1c))
	for sx in [-1, 1]:
		mb.box(Vector3(sx * 1.0, h0 + 1.25, 4), Vector3(0.2, 2.5, 0.5), Models.C(0x5a3a1c))
	var gb := MeshBuilder.new()
	for p in [[-3.2, 4.0, 0], [3.2, 4.0, 0], [-5.0, 0.0, 1], [5.0, 0.0, 1]]:
		gb.box(Vector3(p[0], h0 + 1.8, p[1]), Vector3(1.1 if p[2] == 0 else 0.5, 0.9, 0.5 if p[2] == 0 else 1.1), Models.C(0xffe9a0), Vector3.ZERO, 0.0)
	var mi := mb.instance()
	mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_ON
	parent.add_child(mi)
	parent.add_child(gb.instance(glow_material()))

static func build_bridges(parent: Node3D) -> void:
	var mb := MeshBuilder.new()
	for br in WD.bridges:
		var h := maxf(0.1, WD.height_at(br.x, br.z))
		var inner := MeshBuilder.new()
		inner.box(Vector3(0, 0.55, 0), Vector3(4, 0.3, br.len), Models.C(0x7a5a32))
		for sx in [-1, 1]:
			inner.box(Vector3(sx * 1.9, 1.1, 0), Vector3(0.2, 0.9, br.len), Models.C(0x7a5a32))
			for sz in [-1, 1]:
				inner.box(Vector3(sx * 1.9, 0.9, sz * br.len / 2.2), Vector3(0.3, 1.8, 0.3), Models.C(0x7a5a32))
		var rotated := Basis(Vector3.UP, br.rot)
		var base := mb.verts.size()
		for i in inner.verts.size():
			mb.verts.append(Vector3(br.x, h, br.z) + rotated * inner.verts[i])
			mb.norms.append(rotated * inner.norms[i])
			mb.cols.append(inner.cols[i])
		for i in inner.idx:
			mb.idx.append(base + i)
	var mi := mb.instance()
	mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_ON
	parent.add_child(mi)
