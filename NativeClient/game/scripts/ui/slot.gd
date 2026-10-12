class_name Slot
extends Control
## An icon slot (inventory, equipment, action bar, loot). Shows an emoji icon, border, count, key hint and cooldown sweep.

signal clicked(button: int)
signal double_clicked
signal hovered(on: bool)
signal dropped(data: Variant)

var icon := ""
var count := 0
var border := Color(0.84, 0.7, 0.38, 0.35)
var key_text := ""
var cost_text := ""
var cd_frac := 0.0           # 0..1 of the cooldown that is left
var cd_text := ""
var dim := false             # unusable (red tint)
var no_res := false
var active := false
var placeholder := ""
var flash := 0.0
var drag_data: Variant = null
var accept_drop: Callable = Callable()
var icon_size := 28
var _hover := false

func _init(sz := 46) -> void:
	custom_minimum_size = Vector2(sz, sz)
	size = Vector2(sz, sz)
	mouse_filter = Control.MOUSE_FILTER_STOP
	icon_size = int(sz * 0.62)

func _notification(what: int) -> void:
	if what == NOTIFICATION_MOUSE_ENTER:
		_hover = true
		hovered.emit(true)
		queue_redraw()
	elif what == NOTIFICATION_MOUSE_EXIT:
		_hover = false
		hovered.emit(false)
		queue_redraw()

func _gui_input(ev: InputEvent) -> void:
	if ev is InputEventMouseButton and ev.pressed:
		var mb := ev as InputEventMouseButton
		if mb.double_click and mb.button_index == MOUSE_BUTTON_LEFT:
			double_clicked.emit()
		clicked.emit(mb.button_index)

func _draw() -> void:
	var r := Rect2(Vector2.ZERO, size)
	draw_style_box(UI.box(Color(0, 0, 0, 0.55), border, 2 if icon != "" else 1, 5), r)
	var f := get_theme_default_font()
	if icon != "":
		var a := 0.45 if dim else 1.0
		if no_res:
			a = 0.55
		var tw := f.get_string_size(icon, HORIZONTAL_ALIGNMENT_LEFT, -1, icon_size)
		draw_string(f, Vector2((size.x - tw.x) / 2.0, size.y / 2.0 + icon_size * 0.36), icon, HORIZONTAL_ALIGNMENT_LEFT, -1, icon_size, Color(1, 1, 1, a))
	elif placeholder != "":
		var tw2 := f.get_string_size(placeholder, HORIZONTAL_ALIGNMENT_LEFT, -1, icon_size)
		draw_string(f, Vector2((size.x - tw2.x) / 2.0, size.y / 2.0 + icon_size * 0.36), placeholder, HORIZONTAL_ALIGNMENT_LEFT, -1, icon_size, Color(1, 1, 1, 0.16))
	if dim and icon != "":
		draw_rect(r.grow(-2), Color(0.6, 0, 0, 0.28))
	if cd_frac > 0.0:
		# clock sweep: the dark part shrinks clockwise as the cooldown runs out
		var c := size / 2.0
		var rad := size.x
		var pts := PackedVector2Array([c])
		var a0 := -PI / 2.0 + (1.0 - cd_frac) * TAU
		var a1 := -PI / 2.0 + TAU
		var n := 24
		for i in n + 1:
			var ang := lerpf(a0, a1, float(i) / n)
			var p := c + Vector2(cos(ang), sin(ang)) * rad
			p = p.clamp(Vector2(2, 2), size - Vector2(2, 2))
			pts.append(p)
		draw_colored_polygon(pts, Color(0, 0, 0, 0.72))
		if cd_text != "":
			var tw3 := f.get_string_size(cd_text, HORIZONTAL_ALIGNMENT_LEFT, -1, 16)
			draw_string_outline(f, Vector2((size.x - tw3.x) / 2.0, size.y / 2.0 + 6), cd_text, HORIZONTAL_ALIGNMENT_LEFT, -1, 16, 4, Color.BLACK)
			draw_string(f, Vector2((size.x - tw3.x) / 2.0, size.y / 2.0 + 6), cd_text, HORIZONTAL_ALIGNMENT_LEFT, -1, 16, Color.WHITE)
	if count > 1:
		var t := str(count)
		var tw4 := f.get_string_size(t, HORIZONTAL_ALIGNMENT_LEFT, -1, 12)
		draw_string_outline(f, Vector2(size.x - tw4.x - 4, size.y - 4), t, HORIZONTAL_ALIGNMENT_LEFT, -1, 12, 4, Color.BLACK)
		draw_string(f, Vector2(size.x - tw4.x - 4, size.y - 4), t, HORIZONTAL_ALIGNMENT_LEFT, -1, 12, Color.WHITE)
	if key_text != "":
		draw_string_outline(f, Vector2(4, 12), key_text, HORIZONTAL_ALIGNMENT_LEFT, -1, 11, 4, Color.BLACK)
		draw_string(f, Vector2(4, 12), key_text, HORIZONTAL_ALIGNMENT_LEFT, -1, 11, Color("d8cba4"))
	if cost_text != "":
		var tw5 := f.get_string_size(cost_text, HORIZONTAL_ALIGNMENT_LEFT, -1, 11)
		draw_string_outline(f, Vector2(size.x - tw5.x - 4, 12), cost_text, HORIZONTAL_ALIGNMENT_LEFT, -1, 11, 4, Color.BLACK)
		draw_string(f, Vector2(size.x - tw5.x - 4, 12), cost_text, HORIZONTAL_ALIGNMENT_LEFT, -1, 11, Color("8ad0ff"))
	if active:
		draw_rect(r.grow(-1), Color("ffd84a"), false, 2.0)
	if _hover and icon != "":
		draw_rect(r.grow(-2), Color(1, 1, 1, 0.12))
	if flash > 0.0:
		draw_rect(r, Color(1, 0.9, 0.4, flash * 0.6))

func _process(dt: float) -> void:
	if flash > 0.0:
		flash = maxf(0.0, flash - dt * 2.0)
		queue_redraw()

func _get_drag_data(_pos: Vector2) -> Variant:
	if drag_data == null:
		return null
	var p := Label.new()
	p.text = icon
	p.add_theme_font_size_override("font_size", 32)
	set_drag_preview(p)
	Tooltip.hide_tip()
	return drag_data

func _can_drop_data(_pos: Vector2, data: Variant) -> bool:
	return accept_drop.is_valid() and bool(accept_drop.call(data))

func _drop_data(_pos: Vector2, data: Variant) -> void:
	dropped.emit(data)
