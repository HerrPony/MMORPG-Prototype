class_name Bar
extends Control
## A progress bar with centered text (health, resource, XP, cast).

var value := 0.0
var maxv := 1.0
var fill := Color("3ccf4e")
var text := ""
var font_size := 12
var back := Color(0, 0, 0, 0.6)
var edge := Color(0.84, 0.7, 0.38, 0.45)
var shine := true
var ticks := 0

func _init(w := 200.0, h := 16.0, col := Color("3ccf4e")) -> void:
	custom_minimum_size = Vector2(w, h)
	size = Vector2(w, h)
	fill = col
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func set_state(v: float, m: float, t: String) -> void:
	if v != value or m != maxv or t != text:
		value = v
		maxv = m
		text = t
		queue_redraw()

func set_color(c: Color) -> void:
	if c != fill:
		fill = c
		queue_redraw()

func _draw() -> void:
	var r := Rect2(Vector2.ZERO, size)
	draw_style_box(UI.box(back, edge, 1, 3), r)
	var f := clampf(value / maxv, 0.0, 1.0) if maxv > 0.0 else 0.0
	if f > 0.0:
		var fr := Rect2(Vector2(1, 1), Vector2((size.x - 2.0) * f, size.y - 2.0))
		draw_style_box(UI.box(fill, Color(0, 0, 0, 0), 0, 2), fr)
		if shine:
			draw_rect(Rect2(fr.position, Vector2(fr.size.x, fr.size.y * 0.42)), Color(1, 1, 1, 0.14))
	if ticks > 0:
		for i in range(1, ticks):
			var x := size.x * float(i) / ticks
			draw_line(Vector2(x, 0), Vector2(x, size.y), Color(0, 0, 0, 0.55), 1.0)
	if text != "":
		var fnt := get_theme_default_font()
		var tw := fnt.get_string_size(text, HORIZONTAL_ALIGNMENT_CENTER, size.x, font_size)
		var p := Vector2((size.x - tw.x) / 2.0, size.y / 2.0 + font_size * 0.36)
		draw_string_outline(fnt, p, text, HORIZONTAL_ALIGNMENT_LEFT, -1, font_size, 4, Color.BLACK)
		draw_string(fnt, p, text, HORIZONTAL_ALIGNMENT_LEFT, -1, font_size, Color.WHITE)
