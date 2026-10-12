class_name ActionFrame
extends Control
## Ornate bronze frame of the action bar with a medallion at each end (drawn, no image assets).

const CAP := 44

func _init() -> void:
	mouse_filter = Control.MOUSE_FILTER_STOP

func _draw() -> void:
	var w := size.x
	var h := size.y
	var bronze := Color("8a6a3a")
	var light := Color("d8b068")
	var dark := Color("2a1c10")
	# body with a vertical gradient
	var top := 8.0
	for i in 12:
		var t := float(i) / 11.0
		draw_rect(Rect2(CAP - 6, top + t * (h - top) * 0.92, w - CAP * 2 + 12, (h - top) / 11.0 + 1.0), Color(0.13, 0.09, 0.06, 0.96).lerp(Color(0.05, 0.035, 0.025, 0.98), t))
	draw_line(Vector2(CAP - 6, top), Vector2(w - CAP + 6, top), light, 2.0)
	draw_line(Vector2(CAP - 6, top + 2), Vector2(w - CAP + 6, top + 2), dark, 2.0)
	# end medallions
	for side in [0, 1]:
		var cx := 24.0 if side == 0 else w - 24.0
		var cy := h - 22.0
		var dir := 1.0 if side == 0 else -1.0
		# wings
		var wing := PackedVector2Array([Vector2(cx - dir * 6, cy - 22), Vector2(cx + dir * 30, cy - 10), Vector2(cx + dir * 24, cy + 8), Vector2(cx - dir * 6, cy + 20)])
		draw_colored_polygon(wing, Color("5a3e20"))
		draw_polyline(PackedVector2Array([wing[0], wing[1], wing[2], wing[3]]), light, 1.5)
		draw_circle(Vector2(cx, cy), 21.0, dark)
		draw_circle(Vector2(cx, cy), 19.0, bronze)
		draw_circle(Vector2(cx, cy), 14.0, Color("3a2814"))
		draw_arc(Vector2(cx, cy), 14.0, 0, TAU, 28, light, 1.5, true)
		draw_circle(Vector2(cx, cy), 6.0, light)
		draw_circle(Vector2(cx, cy), 3.5, Color("6a1a10"))
