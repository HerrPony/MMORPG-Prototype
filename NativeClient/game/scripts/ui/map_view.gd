class_name MapView
extends Control
## Base of the minimap and the world map: a GPU-rendered top-down map plus helpers to draw markers on top of it.

var mat: ShaderMaterial
var surface: ColorRect
var center := Vector2.ZERO
var view := Vector2(50, 50)
var right := Vector2(1, 0)
var fwd := Vector2(0, -1)

func _init() -> void:
	mat = ShaderMaterial.new()
	mat.shader = load("res://assets/shaders/map.gdshader")
	# the shader lives on a child that is drawn behind this control, so markers drawn by _draw are not affected by it
	surface = ColorRect.new()
	surface.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	surface.mouse_filter = Control.MOUSE_FILTER_IGNORE
	surface.show_behind_parent = true
	surface.material = mat
	add_child(surface)

func setup(world: GameWorld) -> void:
	mat.set_shader_parameter("hm_tex", world.tex_hm)
	mat.set_shader_parameter("bm_tex", world.tex_bm)
	mat.set_shader_parameter("map_tex", world.tex_map)
	mat.set_shader_parameter("grid", Vector2(WD.nx, WD.nz))
	mat.set_shader_parameter("origin", Vector2(WD.x0, WD.z0))
	mat.set_shader_parameter("cell", WD.cell)

func set_view(c: Vector2, v: Vector2, r: Vector2, f: Vector2, circle: bool, fog := false) -> void:
	center = c; view = v; right = r; fwd = f
	mat.set_shader_parameter("center", c)
	mat.set_shader_parameter("view", v)
	mat.set_shader_parameter("right", r)
	mat.set_shader_parameter("fwd", f)
	mat.set_shader_parameter("circle", circle)
	mat.set_shader_parameter("use_fog", fog)

## World position (x, z) -> pixel position in this control.
func to_px(w: Vector2) -> Vector2:
	var d := w - center
	var sx := d.dot(right)
	var sy := -d.dot(fwd)
	return size / 2.0 + Vector2(sx / view.x, sy / view.y) * size / 2.0

func px_scale() -> float:
	return size.x / 2.0 / view.x

## Direction (world) -> rotation angle of an upward pointing marker.
func dir_angle(dx: float, dz: float) -> float:
	var sx := dx * right.x + dz * right.y
	var sy := -(dx * fwd.x + dz * fwd.y)
	return atan2(sx, -sy)

func arrow(at: Vector2, ang: float, col: Color, scale := 1.0, outline := Color.BLACK) -> void:
	var pts := PackedVector2Array()
	for p in [Vector2(0, -11), Vector2(8, 9), Vector2(0, 4), Vector2(-8, 9)]:
		pts.append(at + (p * scale).rotated(ang))
	draw_colored_polygon(pts, col)
	pts.append(pts[0])
	draw_polyline(pts, outline, 1.5, true)

func dashed_circle(c: Vector2, r: float, col: Color, width := 2.0) -> void:
	var n := maxi(16, int(r / 2.0))
	for i in n:
		if i % 2 == 0:
			draw_arc(c, r, TAU * i / n, TAU * (i + 1) / n, 4, col, width, true)

func text_out(pos: Vector2, text: String, size_px := 14, col := Color.WHITE, align := HORIZONTAL_ALIGNMENT_CENTER, width := 200.0) -> void:
	var f := get_theme_default_font()
	var p := pos
	if align == HORIZONTAL_ALIGNMENT_CENTER:
		p.x -= width / 2.0
	draw_string_outline(f, p, text, align, width, size_px, 5, Color(0, 0, 0, 0.9))
	draw_string(f, p, text, align, width, size_px, col)
