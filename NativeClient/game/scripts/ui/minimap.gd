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

func _draw() -> void:
	var c := size / 2.0
	var sc := px_scale()
	if S.instance != null:
		# dungeon: draw the room rectangles
		draw_circle(c, size.x / 2.0, Color("16141a"))
		var org: Array = S.instance.origin
		for rm in S.instance.rooms:
			var a := to_px(Vector2(rm.x + org[0] - rm.w / 2.0, rm.z + org[1] - rm.d / 2.0))
			var b := to_px(Vector2(rm.x + org[0] + rm.w / 2.0, rm.z + org[1] + rm.d / 2.0))
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
				draw_circle(p, rr, Color(1, 0.78, 0.16, 0.15))
				dashed_circle(p, rr, Color(1, 0.84, 0.25, 0.85))
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
	# rim
	draw_arc(c, size.x / 2.0 - 1.0, 0, TAU, 64, Color(0.84, 0.7, 0.38, 0.9), 2.5, true)
