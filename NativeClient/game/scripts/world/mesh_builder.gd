class_name MeshBuilder
extends RefCounted
## Collects boxes/frusta with baked per-face shading into one ArrayMesh (the blocky "voxel" look).

var verts := PackedVector3Array()
var norms := PackedVector3Array()
var cols := PackedColorArray()
var idx := PackedInt32Array()

static var _mat: StandardMaterial3D = null
static var _mat_glow: StandardMaterial3D = null
static var _mat_alpha: StandardMaterial3D = null

static func material() -> StandardMaterial3D:
	if _mat == null:
		_mat = StandardMaterial3D.new()
		_mat.vertex_color_use_as_albedo = true
		_mat.vertex_color_is_srgb = true
		_mat.roughness = 0.92
		_mat.metallic_specular = 0.15
	return _mat

## Vertex colours are used as emission too: windows, eyes, lava, crystals.
static func glow_material() -> StandardMaterial3D:
	if _mat_glow == null:
		_mat_glow = StandardMaterial3D.new()
		_mat_glow.vertex_color_use_as_albedo = true
		_mat_glow.vertex_color_is_srgb = true
		_mat_glow.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
		_mat_glow.albedo_color = Color(1.1, 1.1, 1.1)
	return _mat_glow

static func alpha_material() -> StandardMaterial3D:
	if _mat_alpha == null:
		_mat_alpha = StandardMaterial3D.new()
		_mat_alpha.vertex_color_use_as_albedo = true
		_mat_alpha.vertex_color_is_srgb = true
		_mat_alpha.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
		_mat_alpha.roughness = 0.6
		_mat_alpha.cull_mode = BaseMaterial3D.CULL_DISABLED
	return _mat_alpha

const FACE_SHADE := {"top": 1.0, "bottom": 0.5, "x": 0.82, "z": 0.7}

## A frustum with a rectangular base (w0 x d0) and top (w1 x d1); bottom centre at c. A box has w0==w1 and d0==d1.
func frustum(c: Vector3, w0: float, d0: float, w1: float, d1: float, h: float, color: Color, rot: Vector3 = Vector3.ZERO, jitter := 0.04, shear := Vector2.ZERO) -> void:
	var b := Basis.from_euler(rot * (PI / 180.0)) if rot != Vector3.ZERO else Basis.IDENTITY
	var corners := [
		Vector3(-w0 / 2, 0, -d0 / 2), Vector3(w0 / 2, 0, -d0 / 2), Vector3(w0 / 2, 0, d0 / 2), Vector3(-w0 / 2, 0, d0 / 2),
		Vector3(-w1 / 2 + shear.x, h, -d1 / 2 + shear.y), Vector3(w1 / 2 + shear.x, h, -d1 / 2 + shear.y), Vector3(w1 / 2 + shear.x, h, d1 / 2 + shear.y), Vector3(-w1 / 2 + shear.x, h, d1 / 2 + shear.y),
	]
	var j := 1.0 + (randf() - 0.5) * 2.0 * jitter
	var faces := [
		[[0, 1, 2, 3], "bottom"], [[4, 7, 6, 5], "top"], [[0, 4, 5, 1], "z"], [[2, 6, 7, 3], "z"], [[3, 7, 4, 0], "x"], [[1, 5, 6, 2], "x"],
	]
	var centre := c + b * Vector3(0, h / 2, 0)
	for f in faces:
		var ids: Array = f[0]
		var p: Array = []
		for i in ids:
			p.append(c + b * corners[i])
		var nc: Vector3 = (p[2] - p[0]).cross(p[3] - p[1]).normalized() # counter-clockwise normal from the diagonals
		var mid: Vector3 = (p[0] + p[1] + p[2] + p[3]) / 4.0 - centre
		var out: Vector3 = nc if nc.dot(mid) >= 0 else -nc
		if nc.dot(mid) > 0:
			p.reverse() # Godot wants clockwise front faces
		var s: float = FACE_SHADE[f[1]]
		var col := Color(color.r * s * j, color.g * s * j, color.b * s * j, color.a)
		var base := verts.size()
		for q in p:
			verts.append(q)
			norms.append(out)
			cols.append(col)
		idx.append_array([base, base + 1, base + 2, base, base + 2, base + 3])

## A box centred at c.
func box(c: Vector3, size: Vector3, color: Color, rot: Vector3 = Vector3.ZERO, jitter := 0.04) -> void:
	if rot == Vector3.ZERO:
		frustum(c - Vector3(0, size.y / 2, 0), size.x, size.z, size.x, size.z, size.y, color, Vector3.ZERO, jitter)
	else:
		_box_rot(c, size, color, rot, jitter)

func _box_rot(c: Vector3, size: Vector3, color: Color, rot: Vector3, jitter: float) -> void:
	# rotate about the box centre: build around the origin and offset
	var b := Basis.from_euler(rot * (PI / 180.0))
	var start := verts.size()
	frustum(Vector3(0, -size.y / 2, 0), size.x, size.z, size.x, size.z, size.y, color, Vector3.ZERO, jitter)
	for i in range(start, verts.size()):
		verts[i] = c + b * verts[i]
		norms[i] = b * norms[i]

func get_mesh(count_only := false) -> ArrayMesh:
	var mesh := ArrayMesh.new()
	if verts.is_empty():
		return mesh
	var arr := []
	arr.resize(Mesh.ARRAY_MAX)
	arr[Mesh.ARRAY_VERTEX] = verts
	arr[Mesh.ARRAY_NORMAL] = norms
	arr[Mesh.ARRAY_COLOR] = cols
	arr[Mesh.ARRAY_INDEX] = idx
	mesh.add_surface_from_arrays(Mesh.PRIMITIVE_TRIANGLES, arr)
	return mesh

## Adds the contents of another builder, moved by an offset.
func append(o: MeshBuilder, offset: Vector3 = Vector3.ZERO) -> void:
	var base := verts.size()
	for i in o.verts.size():
		verts.append(o.verts[i] + offset)
		norms.append(o.norms[i])
		cols.append(o.cols[i])
	for i in o.idx:
		idx.append(base + i)

func instance(mat: Material = null) -> MeshInstance3D:
	var mi := MeshInstance3D.new()
	mi.mesh = get_mesh()
	mi.material_override = mat if mat else material()
	return mi
