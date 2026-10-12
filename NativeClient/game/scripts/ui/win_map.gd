class_name WinMap
extends GameWindow
## World map per continent with fog on undiscovered zones, towns, points of interest and quest areas.

class MapCanvas extends MapView:
	var win: WinMap
	func _draw() -> void:
		win.draw_overlay(self)

var map: MapCanvas
var cont := 0
var centroids: Dictionary = {}
var pts: Array = []
var areas: Array = []
var tip_label: Label
var cont_buttons: Array = []
var world: GameWorld

func _init() -> void:
	super._init("Weltkarte", Vector2(640, 100))
	win_id = "map"
	default_pos = Vector2(250, 40)
	var tabs := UI.hbox(4)
	body.add_child(tabs)
	for i in 2:
		var idx := i
		var b := UI.button("Kontinent %d" % (i + 1), func(): cont = idx; _sync_tabs(); map.queue_redraw())
		b.toggle_mode = true
		cont_buttons.append(b)
		tabs.add_child(b)
	map = MapCanvas.new()
	map.win = self
	map.custom_minimum_size = Vector2(620, 620)
	map.size = Vector2(620, 620)
	map.mouse_filter = Control.MOUSE_FILTER_STOP
	map.gui_input.connect(_map_input)
	map.mouse_exited.connect(func(): tip_label.visible = false)
	body.add_child(map)
	var legend := UI.label("★ Hauptstadt  ■ Dorf  ⚓ Hafen  🦅 Flugpunkt  ⛺ Lager  🕳️ Höhle  💀 Dungeon  ◎ Questgebiet  ▲ Du  ● Gruppe", 11, UI.MUTED)
	legend.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	legend.custom_minimum_size.x = 600
	body.add_child(legend)
	tip_label = UI.label("", 12, Color.WHITE)
	tip_label.visible = false
	tip_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	tip_label.z_index = 5
	tip_label.add_theme_constant_override("outline_size", 5)
	tip_label.add_theme_color_override("font_outline_color", Color.BLACK)
	map.add_child(tip_label)

func setup(w: GameWorld) -> void:
	world = w
	map.setup(w)
	centroids.clear()

func open(at := Vector2(-1, -1)) -> void:
	cont = WD.cont_of(S.local.x)
	_sync_tabs()
	super.open(at)
	var sz := clampf(get_viewport_rect().size.y - 190.0, 300.0, 700.0)
	map.custom_minimum_size = Vector2(sz, sz)
	map.size = Vector2(sz, sz)
	reset_size()

func _sync_tabs() -> void:
	for i in cont_buttons.size():
		cont_buttons[i].button_pressed = i == cont

func _process(_dt: float) -> void:
	if not visible or world == null:
		return
	var c: Dictionary = S.scenery.conts[cont]
	var known: Array = []
	for b in 11:
		var kb := 9 if b == 10 else b
		known.append(1.0 if (b == 0 or S.disc.has("zone:" + WD.BIOME_KEYS[kb])) else 0.0)
	map.mat.set_shader_parameter("known", known)
	var cx: float = float(c.x0) + float(c.size) / 2.0
	var half: float = float(c.size) / 2.0
	map.set_view(Vector2(cx, WD.z0 + half), Vector2(half, half), Vector2(1, 0), Vector2(0, -1), false, true)
	map.queue_redraw()

func _centroids() -> Dictionary:
	if not centroids.is_empty():
		return centroids
	var acc := {}
	var j := 0
	while j < WD.nz:
		var i := 0
		while i < WD.nx:
			var b: int = WD.bm[j * WD.nx + i]
			if b != 0:
				if b == 10:
					b = 9
				var a: Dictionary = acc.get(b, {"x": 0.0, "z": 0.0, "n": 0})
				a.x += WD.x0 + i * WD.cell
				a.z += WD.z0 + j * WD.cell
				a.n += 1
				acc[b] = a
			i += 3
		j += 3
	for b in acc:
		centroids[b] = Vector2(acc[b].x / acc[b].n, acc[b].z / acc[b].n)
	return centroids

func discovered(id: String) -> bool:
	return S.disc.has(id)

func draw_overlay(m: MapCanvas) -> void:
	var sc: Dictionary = S.scenery
	var size := m.size.x
	var cen := _centroids()
	var big := size > 0.0
	# zone labels
	for z in sc.zones:
		if int(z.id) == 0 or int(z.id) == 10 or int(z.cont) != cont:
			continue
		var ce = cen.get(int(z.id))
		if ce == null:
			continue
		var known := discovered("zone:" + str(z.key))
		var p := m.to_px(ce)
		var col := Color(1, 0.94, 0.78, 0.92) if known else Color(0.7, 0.7, 0.75, 0.55)
		m.text_out(p, str(z.name) if known else "???", 20, col, HORIZONTAL_ALIGNMENT_CENTER, 300)
		m.text_out(p + Vector2(0, 17), "Stufe %d–%d" % [int(z.minLvl), int(z.maxLvl)], 12, col, HORIZONTAL_ALIGNMENT_CENTER, 200)
	# quest areas
	areas = []
	var tt := Time.get_ticks_msec() / 1000.0
	for a in QAreas.get_all():
		if WD.cont_of(a.x) != cont:
			continue
		var p2 := m.to_px(Vector2(a.x, a.z))
		var rr := maxf(8.0, a.r * m.px_scale())
		m.draw_circle(p2, rr, Color(1, 0.78, 0.16, 0.2))
		m.dashed_circle(p2, rr, Color(1, 0.84, 0.25, 0.7 + sin(tt * 3.0) * 0.2), 2.5)
		m.text_out(p2 + Vector2(0, 6), "?" if (a.turn_in or a.k == "npc") else "!", 18, Color("ffd01a"))
		areas.append(a)
	# towns and POIs
	pts = []
	for t in sc.towns:
		if t.kind == "site" or WD.cont_of(t.x) != cont:
			continue
		var kn := discovered(t.id)
		pts.append({"x": t.x, "z": t.z, "name": t.name if kn else "???", "type": "harbor" if t.get("harbor", false) else t.kind, "known": kn, "id": t.id, "fp": S.fp.has(t.id)})
	for p3 in sc.pois:
		if p3.get("town", false) or WD.cont_of(p3.x) != cont:
			continue
		if not discovered(p3.id):
			continue
		pts.append({"x": p3.x, "z": p3.z, "name": p3.name, "type": p3.kind, "known": true, "id": p3.id})
	for dg in [["rabenfels", "Schwarzfelsmine (Dungeon)"], ["nebeltempel", "Nebeltempel (Dungeon)"], ["feuerschlund", "Feuerschlund (Dungeon)"]]:
		for p4 in sc.pois:
			if p4.id == dg[0] and WD.cont_of(p4.x) == cont and discovered(dg[0]):
				pts.append({"x": p4.x, "z": p4.z, "name": dg[1], "type": "dungeon", "known": true, "id": "dungeon"})
	for p in pts:
		var px := m.to_px(Vector2(p.x, p.z))
		var is_cap: bool = p.type == "capital"
		var is_town: bool = p.type == "village" or p.type == "harbor"
		if is_cap or is_town:
			var col2 := (Color("ffd86a") if is_cap else Color("ffe9a8")) if p.known else Color("777788")
			if is_cap:
				var star := PackedVector2Array()
				for i in 5:
					var a1 := -PI / 2.0 + i * 1.2566
					star.append(px + Vector2(cos(a1), sin(a1)) * 9.0)
					star.append(px + Vector2(cos(a1 + 0.628), sin(a1 + 0.628)) * 4.0)
				m.draw_colored_polygon(star, col2)
				star.append(star[0])
				m.draw_polyline(star, Color.BLACK, 1.5)
			else:
				m.draw_rect(Rect2(px - Vector2(5, 5), Vector2(10, 10)), col2)
				m.draw_rect(Rect2(px - Vector2(5, 5), Vector2(10, 10)), Color.BLACK, false, 1.5)
			if p.get("fp", false):
				m.text_out(px + Vector2(14, 5), "🦅", 15)
			if p.type == "harbor":
				m.text_out(px + Vector2(14 if px.x < size / 2.0 else -14, 5), "⚓", 15)
			m.text_out(px + Vector2(0, -12), p.name, 15 if is_cap else 13, Color.WHITE if p.known else Color("9999aa"), HORIZONTAL_ALIGNMENT_CENTER, 220)
		else:
			var icon: String = {"camp": "⛺", "cave": "🕳️", "lake": "💧", "ruin": "🏛️", "wild": "📍", "dungeon": "💀"}.get(p.type, "•")
			m.text_out(px + Vector2(0, 5), icon, 15)
			m.text_out(px + Vector2(0, 20), p.name, 11, Color("eeeeee"), HORIZONTAL_ALIGNMENT_CENTER, 200)
	# ferry
	var harbors: Array = sc.towns.filter(func(t): return t.get("harbor", false))
	var here = null
	var other = null
	for h in harbors:
		if WD.cont_of(h.x) == cont:
			here = h
		else:
			other = h
	if here != null and other != null and discovered(here.id):
		var a0 := m.to_px(Vector2(here.x, here.z))
		var a1 := Vector2(size if cont == 0 else 0.0, a0.y + 0.2 * (m.to_px(Vector2(other.x, other.z)).y - a0.y))
		var n := 22
		for i in n:
			if i % 2 == 0:
				m.draw_line(a0.lerp(a1, float(i) / n), a0.lerp(a1, float(i + 1) / n), Color(1, 1, 1, 0.65), 2.0)
		m.text_out(Vector2(size - 70 if cont == 0 else 70, a0.y + 28), "⛵ Überfahrt", 12, Color("cfe6ff"))
	# party and the player
	for st in S.party.states:
		if int(st.id) == S.id or int(st.get("i", 0)) != 0 or WD.cont_of(st.x) != cont:
			continue
		var pp := m.to_px(Vector2(st.x, st.z))
		m.draw_circle(pp, 6.0, Color.BLACK)
		m.draw_circle(pp, 5.0, Color("6aa8ff"))
	if S.instance == null and WD.cont_of(S.local.x) == cont:
		var me := m.to_px(Vector2(S.local.x, S.local.z))
		m.arrow(me, m.dir_angle(sin(S.local.r), cos(S.local.r)), Color("ff4a4a"), 1.0, Color.WHITE)
		var ring := fposmod(Time.get_ticks_msec() / 600.0, 1.0)
		m.draw_arc(me, 8.0 + ring * 14.0, 0, TAU, 28, Color(1, 0.3, 0.3, 1.0 - ring), 2.0, true)
	m.draw_rect(Rect2(Vector2.ZERO, m.size), Color("6d5a30"), false, 4.0)

func _map_input(ev: InputEvent) -> void:
	if not (ev is InputEventMouseMotion):
		return
	var pos: Vector2 = (ev as InputEventMouseMotion).position
	var c: Dictionary = S.scenery.conts[cont]
	var half: float = float(c.size) / 2.0
	var w := Vector2(float(c.x0) + half + (pos.x / map.size.x * 2.0 - 1.0) * half, WD.z0 + half + (pos.y / map.size.y * 2.0 - 1.0) * half)
	var best = null
	var bd := 18.0
	for p in pts:
		var d := Vector2(p.x - w.x, p.z - w.y).length()
		if d < bd:
			bd = d
			best = p
	var text := ""
	if best != null:
		text = best.name + (" · Flugpunkt" if best.get("fp", false) else (" · Hafen" if best.type == "harbor" else ""))
	for a in areas:
		if Vector2(a.x - w.x, a.z - w.y).length() <= maxf(10.0, a.r):
			text += ("\n" if text != "" else "") + "Quest: %s – %s" % [a.title, a.text]
			break
	tip_label.visible = text != ""
	tip_label.text = text
	tip_label.position = pos + Vector2(12, 12)
