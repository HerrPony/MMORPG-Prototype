class_name WD
extends RefCounted
## World data decoded from the server's scenery message: terrain grid, biomes, props, buildings, collision.
## Static state (there is only one world).

static var ready := false
static var nx := 0
static var nz := 0
static var cell := 2.0
static var x0 := -300.0
static var z0 := -300.0
static var hm := PackedFloat32Array()
static var bm := PackedByteArray()
static var roads: Array = []
static var rivers: Array = []
static var lakes: Array = []
static var towns: Array = []
static var pois: Array = []
static var buildings: Array = []
static var house: Array = []
static var bridges: Array = []
static var zones: Array = []
static var conts: Array = []
static var prop_kinds: Array = []
# props: parallel arrays
static var p_kind := PackedInt32Array()
static var p_x := PackedFloat32Array()
static var p_z := PackedFloat32Array()
static var p_s := PackedFloat32Array()
static var p_ry := PackedFloat32Array()
static var p_solid := PackedFloat32Array()
static var tree_grid: Dictionary = {}    # Vector2i (8 m cells) -> PackedFloat32Array [x, z, r, x, z, r ...]
static var solids: Array = []            # building collision boxes {x,z,w,d}
static var docks: Array = []

const BIOME_KEYS := ["sea", "forest", "plains", "swamp", "snow", "desert", "jungle", "blossom", "mesa", "volcano", "volcano"]
const NO_SOLID := ["pad", "dock", "ship", "lantern", "pond", "gate"]

static func load_scenery(sc: Dictionary) -> void:
	nx = int(sc.nx); nz = int(sc.nz); cell = float(sc.cell); x0 = float(sc.x0); z0 = float(sc.z0)
	var raw := Marshalls.base64_to_raw(sc.hm)
	hm = PackedFloat32Array()
	hm.resize(nx * nz)
	for i in nx * nz:
		hm[i] = raw.decode_s16(i * 2) / 20.0
	bm = Marshalls.base64_to_raw(sc.bm)
	roads = sc.get("roads", [])
	rivers = sc.get("rivers", [])
	lakes = sc.get("lakes", [])
	towns = sc.get("towns", [])
	pois = sc.get("pois", [])
	buildings = sc.get("buildings", [])
	house = sc.get("house", [])
	bridges = sc.get("bridges", [])
	zones = sc.get("zones", [])
	conts = sc.get("conts", [])
	prop_kinds = sc.get("pk", [])
	_decode_props(sc.get("props", ""))
	docks = buildings.filter(func(b): return b.k == "dock")
	solids = []
	for b in buildings:
		if b.k in NO_SOLID:
			continue
		var mul := 0.9 if b.k == "tent" else 1.0
		solids.append({"x": b.x, "z": b.z, "w": b.w * mul, "d": b.d * mul})
	ready = true

static func _decode_props(b64: String) -> void:
	var bin := Marshalls.base64_to_raw(b64)
	var n := bin.size() / 8
	p_kind.resize(n); p_x.resize(n); p_z.resize(n); p_s.resize(n); p_ry.resize(n); p_solid.resize(n)
	tree_grid.clear()
	for i in n:
		var o := i * 8
		p_kind[i] = bin[o]
		p_x[i] = x0 + (bin[o + 1] | (bin[o + 2] << 8)) / 40.0
		p_z[i] = z0 + (bin[o + 3] | (bin[o + 4] << 8)) / 40.0
		p_s[i] = bin[o + 5] / 100.0
		p_ry[i] = bin[o + 6] / 255.0 * 6.28
		p_solid[i] = bin[o + 7] / 100.0
		if p_solid[i] > 0.0:
			var key := Vector2i(floori(p_x[i] / 8.0), floori(p_z[i] / 8.0))
			if not tree_grid.has(key):
				tree_grid[key] = PackedFloat32Array()
			var arr: PackedFloat32Array = tree_grid[key]
			arr.append(p_x[i]); arr.append(p_z[i]); arr.append(p_solid[i] * p_s[i])
			tree_grid[key] = arr

static func _grid(i: int, j: int) -> float:
	i = clampi(i, 0, nx - 1)
	j = clampi(j, 0, nz - 1)
	return hm[j * nx + i]

## Terrain height (bilinear on the grid, like the server); dungeons far to the south are flat.
static func height_at(x: float, z: float) -> float:
	if z > 500.0 or not ready:
		return 0.0
	var gx := (x - x0) / cell
	var gz := (z - z0) / cell
	var i := floori(gx)
	var j := floori(gz)
	var fx := gx - i
	var fz := gz - j
	var a := _grid(i, j)
	var b := _grid(i + 1, j)
	var c := _grid(i, j + 1)
	var d := _grid(i + 1, j + 1)
	return a + (b - a) * fx + (c - a) * fz + (a - b - c + d) * fx * fz

static func biome_at(x: float, z: float) -> int:
	if not ready or z > 500.0:
		return 1
	var i := clampi(roundi((x - x0) / cell), 0, nx - 1)
	var j := clampi(roundi((z - z0) / cell), 0, nz - 1)
	return bm[j * nx + i]

static func biome_key(x: float, z: float) -> String:
	return BIOME_KEYS[biome_at(x, z)]

static func cont_of(x: float) -> int:
	return 1 if x > 420.0 else 0

static func on_dock(x: float, z: float) -> bool:
	for b in docks:
		if absf(x - b.x) <= b.w / 2.0 and absf(z - b.z) <= b.d / 2.0:
			return true
	return false

## Walking height: piers are walkways above the water.
static func gh(x: float, z: float) -> float:
	return 1.2 if on_dock(x, z) else height_at(x, z)

static func dist_to_seg(px: float, pz: float, ax: float, az: float, bx: float, bz: float) -> float:
	var dx := bx - ax
	var dz := bz - az
	var l2 := dx * dx + dz * dz
	var t := 0.0
	if l2 > 0.0:
		t = clampf(((px - ax) * dx + (pz - az) * dz) / l2, 0.0, 1.0)
	return Vector2(px - (ax + dx * t), pz - (az + dz * t)).length()

## Collision with trees, rocks and buildings, plus water depth and steep slopes (same rules as the web client / server).
static func collide(x: float, z: float, ox: float, oz: float, has_old := true, radius := 0.4) -> Vector2:
	var gx := floori(x / 8.0)
	var gz := floori(z / 8.0)
	for ix in range(gx - 1, gx + 2):
		for iz in range(gz - 1, gz + 2):
			var cellv = tree_grid.get(Vector2i(ix, iz))
			if cellv == null:
				continue
			var arr: PackedFloat32Array = cellv
			var k := 0
			while k < arr.size():
				var dx := x - arr[k]
				var dz := z - arr[k + 1]
				var d := sqrt(dx * dx + dz * dz)
				var mn := arr[k + 2] + radius
				if d < mn and d > 1e-4:
					x = arr[k] + dx / d * mn
					z = arr[k + 1] + dz / d * mn
				k += 3
	for b in house:
		var hx: float = b.w / 2.0 + radius
		var hz: float = b.d / 2.0 + radius
		var dx2: float = x - b.x
		var dz2: float = z - b.z
		if absf(dx2) < hx and absf(dz2) < hz:
			var px := hx - absf(dx2)
			var pz := hz - absf(dz2)
			if px < pz:
				x = b.x + signf(dx2 if dx2 != 0.0 else 1.0) * hx
			else:
				z = b.z + signf(dz2 if dz2 != 0.0 else 1.0) * hz
	for b in solids:
		if absf(x - b.x) >= b.w / 2.0 + 3.0 or absf(z - b.z) >= b.d / 2.0 + 3.0:
			continue
		var hx2: float = b.w / 2.0 + radius
		var hz2: float = b.d / 2.0 + radius
		var dx3: float = x - b.x
		var dz3: float = z - b.z
		if absf(dx3) < hx2 and absf(dz3) < hz2:
			var px2 := hx2 - absf(dx3)
			var pz2 := hz2 - absf(dz3)
			if px2 < pz2:
				x = b.x + signf(dx3 if dx3 != 0.0 else 1.0) * hx2
			else:
				z = b.z + signf(dz3 if dz3 != 0.0 else 1.0) * hz2
	var xmax := x0 + (nx - 1) * cell - 6.0
	var zmax := z0 + (nz - 1) * cell - 6.0
	x = clampf(x, x0 + 6.0, xmax)
	z = clampf(z, z0 + 6.0, zmax)
	if has_old:
		var h := gh(x, z)
		if h < -1.5:
			return Vector2(ox, oz)
		var step := Vector2(x - ox, z - oz).length()
		if step > 1e-4 and (h - gh(ox, oz)) / step > 1.35:
			if (gh(x, oz) - gh(ox, oz)) / maxf(absf(x - ox), 1e-4) <= 1.35 and gh(x, oz) > -1.5:
				return Vector2(x, oz)
			if (gh(ox, z) - gh(ox, oz)) / maxf(absf(z - oz), 1e-4) <= 1.35 and gh(ox, z) > -1.5:
				return Vector2(ox, z)
			return Vector2(ox, oz)
	return Vector2(x, z)
