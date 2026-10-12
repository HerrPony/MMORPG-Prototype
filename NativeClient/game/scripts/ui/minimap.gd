class_name Minimap
extends MapView
## Round rotating minimap with quest areas, party members, NPC markers and creatures.

const VIEW := 52.0
var cam_yaw := 0.0
var north: Label
var zone_label: Label

func _init() -> void:
	super._init()
	custom_minimum_size = Vector2(176, 176)
	size = Vector2(176, 176)
	mouse_filter = Control.MOUSE_FILTER_PASS
	north = Label.new()
	north.text = "N"
	north.add_theme_color_override("font_color", Color("ff6a5a"))
	north.add_theme_font_size_override("font_size", 14)
	north.add_theme_constant_override("outline_size", 4)
	north.add_theme_color_override("font_outline_color", Color.BLACK)
	north.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(north)

func update_view(yaw: float) -> void:
	cam_yaw = yaw
	var f := Vector2(sin(yaw), cos(yaw))
	var r := Vector2(-cos(yaw), sin(yaw))
	set_view(Vector2(S.local.x, S.local.z), Vector2(VIEW, VIEW), r, f, true)
	# north is world -z
	var nd := Vector2(0, -1)
	var ang := Vector2(nd.dot(r), -nd.dot(f)).angle()
	north.position = size / 2.0 + Vector2(cos(ang), sin(ang)) * (size.x / 2.0 - 10.0) - Vector2(5, 10)
	queue_redraw()

## Quest area circle (fill + dashed outline) limited to the round map, so it never spills over the rim.
func _area_clipped(p: Vector2, rr: float, c: Vector2, limit: float) -> void:
	var n := 64
	var area := PackedVector2Array()
	var disc := PackedVector2Array()
	for i in n:
		var a := TAU * i / n
		area.append(p + Vector2(cos(a), sin(a)) * rr)
		disc.append(c + Vector2(cos(a), sin(a)) * limit)
	for poly in Geometry2D.intersect_polygons(area, disc):
		if poly.size() >= 3 and not Geometry2D.triangulate_polygon(poly).is_empty():
			draw_colored_polygon(poly, Color(1, 0.78, 0.16, 0.15))
	# dashed outline: only the dashes that lie inside the map
	var dashes := maxi(16, int(rr / 2.0))
	var col := Color(1, 0.84, 0.25, 0.85)
	for i in range(0, dashes, 2):
		var a0 := TAU * i / dashes
		var a1 := TAU * (i + 1) / dashes
		var pts := PackedVector2Array()
		for k in 5:
			var q := p + Vector2(cos(lerpf(a0, a1, k / 4.0)), sin(lerpf(a0, a1, k / 4.0))) * rr
			if (q - c).length() <= limit:
				pts.append(q)
			elif pts.size() >= 2:
				draw_polyline(pts, col, 2.0, true)
				pts = PackedVector2Array()
			else:
				pts = PackedVector2Array()
		if pts.size() >= 2:
			draw_polyline(pts, col, 2.0, true)

func _draw() -> void:
	var c := size / 2.0
	var sc := px_scale()
	if S.instance != null:
		# dungeon: draw the room rectangles
		draw_circle(c, size.x / 2.0, Color("16141a"))
		var org: Array = S.instance.origin
		for rm in S.instance.rooms:
			# rooms are axis aligned in the world; the map is rotated, so draw them as polygons
			var corners := [Vector2(rm.x + org[0] - rm.w / 2.0, rm.z + org[1] - rm.d / 2.0), Vector2(rm.x + org[0] + rm.w / 2.0, rm.z + org[1] - rm.d / 2.0), Vector2(rm.x + org[0] + rm.w / 2.0, rm.z + org[1] + rm.d / 2.0), Vector2(rm.x + org[0] - rm.w / 2.0, rm.z + org[1] + rm.d / 2.0)]
			var pts := PackedVector2Array()
			for cn in corners:
				pts.append(to_px(cn))
			draw_colored_polygon(pts, Color("4a443e"))
	var t := Time.get_ticks_msec() / 1000.0
	# quest areas
	if S.instance == null:
		for a in QAreas.get_all():
			var p := to_px(Vector2(a.x, a.z))
			var rr: float = a.r * sc
			var d := (p - c).length()
			if d < size.x / 2.0 + rr:
				_area_clipped(p, rr, c, size.x / 2.0 - 3.0)
				if a.k == "npc":
					text_out(p + Vector2(0, 6), "?", 18, Color("ffd01a"))
			if d > size.x / 2.0 - 6.0:
				var ang := (p - c).angle()
				var rad := size.x / 2.0 - 9.0
				var at := c + Vector2(cos(ang), sin(ang)) * rad
				var al := 0.65 + sin(t * 4.0) * 0.3
				var pts2 := PackedVector2Array([at + Vector2(0, -7).rotated(ang + PI / 2), at + Vector2(6, 5).rotated(ang + PI / 2), at + Vector2(-6, 5).rotated(ang + PI / 2)])
				draw_colored_polygon(pts2, Color(1, 0.82, 0.1, al))
	for e: Ent in S.ents.values():
		if e.id == S.id or e.k == "obj":
			continue
		var p := to_px(Vector2(e.rx, e.rz))
		if (p - c).length() > size.x / 2.0 + 6.0:
			continue
		match e.k:
			"npc":
				draw_circle(p, 4.0, Color("ffe14a"))
				var m: String = S.marks.get(e.id, "")
				if m != "":
					text_out(p + Vector2(0, -2), "!" if m == "available" else "?", 18, Color("aaaaaa") if m == "progress" else Color("ffd01a"))
			"mob":
				var col := Color("666666")
				if not e.d:
					col = Color("ff8a20") if e.el else Color("ffd23a")
				elif e.lt:
					col = Color("ffd24a")
				var rad2 := 3.0 if e.d else (5.5 if e.el else 4.0)
				draw_circle(p, rad2 + 1.0, Color.BLACK)
				draw_circle(p, rad2, col)
			"ship":
				text_out(p + Vector2(0, 6), "⛵", 17, Color.WHITE)
			"pet":
				draw_circle(p, 3.0, Color("7ad87a"))
			_:
				var in_party := false
				for m2 in S.party.members:
					if int(m2.id) == e.id:
						in_party = true
				draw_circle(p, 5.5, Color.BLACK)
				draw_circle(p, 4.5, Color("6aa8ff") if in_party else Color("e8e8e8"))
	arrow(c, dir_angle(sin(S.local.r), cos(S.local.r)), Color.WHITE)
	# bronze rim with a notch at the top
	var r := size.x / 2.0
	draw_arc(c, r - 3.0, 0, TAU, 72, Color("1c1209"), 8.0, true)
	draw_arc(c, r - 3.0, 0, TAU, 72, Color("8a6a3a"), 6.0, true)
	draw_arc(c, r - 5.0, 0, TAU, 72, Color("d8b068"), 1.5, true)
	draw_arc(c, r - 0.5, 0, TAU, 72, Color("2a1c10"), 1.5, true)
	draw_colored_polygon(PackedVector2Array([Vector2(c.x - 9, 5), Vector2(c.x + 9, 5), Vector2(c.x, -6)]), Color("8a6a3a"))
	draw_polyline(PackedVector2Array([Vector2(c.x - 9, 5), Vector2(c.x, -6), Vector2(c.x + 9, 5)]), Color("d8b068"), 1.5)
