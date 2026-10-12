class_name Entities
extends Node3D
## Creates, updates and removes the scene nodes of all entities from the server snapshots and handles picking.

var plates: Control            # 2D overlay for nameplates (set by main)
var camera: Camera3D
var ring: MeshInstance3D
var glow_tex: GradientTexture2D
var mark_font_color := Color("ffd01a")
var t_total := 0.0

func _ready() -> void:
	ring = MeshInstance3D.new()
	var tm := TorusMesh.new()
	tm.inner_radius = 0.8
	tm.outer_radius = 0.95
	tm.rings = 32
	tm.ring_segments = 6
	ring.mesh = tm
	ring.scale = Vector3(1, 0.12, 1)
	var mat := StandardMaterial3D.new()
	mat.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
	mat.albedo_color = Color(1, 1, 0, 0.85)
	mat.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
	mat.no_depth_test = true
	mat.render_priority = 10
	ring.material_override = mat
	ring.visible = false
	ring.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	add_child(ring)
	glow_tex = GradientTexture2D.new()
	glow_tex.fill = GradientTexture2D.FILL_RADIAL
	glow_tex.fill_from = Vector2(0.5, 0.5)
	glow_tex.fill_to = Vector2(1.0, 0.5)
	glow_tex.width = 64
	glow_tex.height = 64
	var g := Gradient.new()
	g.colors = PackedColorArray([Color(1, 1, 1, 1), Color(1, 0.92, 0.6, 0.8), Color(1, 0.8, 0.25, 0)])
	g.offsets = PackedFloat32Array([0.0, 0.3, 1.0])
	glow_tex.gradient = g

# ---------------------------------------------------------------- snapshots

func apply(list: Array, lootable: Array) -> void:
	var seen := {}
	var lt := {}
	for i in lootable:
		lt[int(i)] = true
	for v in list:
		var id := int(v.id)
		seen[id] = true
		var e: Ent = S.ents.get(id)
		var fresh := e == null
		if fresh:
			e = _create(v)
		var was_dead := e.d
		var was_lt := e.lt
		var model_changed: bool = fresh or v.get("m", "") != e.m
		# optional fields are omitted by the server when empty
		e.fl = false; e.d = false; e.tg = 0; e.st = []; e.el = false; e.bs = false; e.ow = 0; e.q = ""; e.gn = false; e.pt = false; e.ti = ""; e.c = ""; e.sc = 1.0
		var eqv = v.get("eq")
		e.eq = eqv if eqv is Dictionary else {}
		var lkv = v.get("lk")
		e.lk = lkv if lkv is Dictionary else {}
		e.k = v.get("k", ""); e.n = v.get("n", ""); e.m = v.get("m", ""); e.c = v.get("c", ""); e.ti = v.get("ti", "")
		e.l = int(v.get("l", 1)); e.hp = int(v.get("hp", 1)); e.mh = int(v.get("mh", 1))
		e.d = bool(v.get("d", false)); e.tg = int(v.get("tg", 0)); e.st = v.get("st", [])
		e.el = bool(v.get("el", false)); e.bs = bool(v.get("bs", false)); e.sc = float(v.get("sc", 1.0))
		e.ow = int(v.get("ow", 0)); e.fl = bool(v.get("fl", false)); e.q = v.get("q", ""); e.gn = bool(v.get("gn", false)); e.pt = bool(v.get("pt", false))
		e.lt = lt.has(id)
		e.x = float(v.x); e.y = float(v.y); e.z = float(v.z); e.r = float(v.r)
		if fresh:
			e.rx = e.x; e.ry = e.y; e.rz = e.z; e.rr = e.r
		if model_changed:
			_rebuild_model(e)
		e.node.scale = Vector3.ONE * e.sc
		if e.k == "player":
			var key := JSON.stringify(e.eq) + JSON.stringify(e.lk)
			if key != e.eq_key:
				e.eq_key = key
				if not fresh:
					_rebuild_model(e)
		_sync_mount(e)
		if id == S.id and not e.inited:
			e.inited = true
			S.local.x = e.x; S.local.y = e.y; S.local.z = e.z; S.local.r = e.r
		if e.lt != was_lt:
			_set_sparkle(e, e.lt)
		if e.k == "npc":
			_refresh_marker(e)
		if e.k == "obj":
			_refresh_object(e)
		if e.d and not was_dead and e.k == "mob":
			pass
		_refresh_plate(e)
	for id in S.ents.keys():
		if not seen.has(id):
			remove(S.ents[id])

func _create(v: Dictionary) -> Ent:
	var e := Ent.new()
	e.id = int(v.id)
	e.node = Node3D.new()
	add_child(e.node)
	S.ents[e.id] = e
	e.k = v.get("k", "")
	e.m = v.get("m", "")
	return e

func _rebuild_model(e: Ent) -> void:
	if e.model:
		e.model.queue_free()
		e.model = null
	var src := {"id": e.id, "k": e.k, "m": e.m, "c": e.c, "lk": e.lk, "eq": e.eq}
	e.model = Models.make(src)
	e.node.add_child(e.model)
	if e.m == "flightmaster" and e.pad == null:
		e.pad = Models.make_griffin()
		e.pad.position = Vector3(2.8, 0, 0.5)
		e.pad.rotation.y = -0.6
		e.node.add_child(e.pad)
	if e.mount:
		e.model.position = Vector3(0, e.mount.get_meta("seat_y", 0.0), e.mount.get_meta("seat_z", 0.0))

func _sync_mount(e: Ent) -> void:
	if e.k != "player":
		return
	if e.fl and e.mount == null:
		e.mount = Models.make_griffin()
		e.node.add_child(e.mount)
		e.model.position = Vector3(0, float(e.mount.get_meta("seat_y", 1.1)), float(e.mount.get_meta("seat_z", 0.1)))
	elif not e.fl and e.mount != null:
		e.mount.queue_free()
		e.mount = null
		e.model.position = Vector3.ZERO

func remove(e: Ent) -> void:
	if e.node:
		e.node.queue_free()
	if e.plate:
		e.plate.queue_free()
	if e.sparkle:
		e.sparkle.queue_free()
	S.ents.erase(e.id)
	if S.me.tg == e.id:
		S.me.tg = 0

func remove_all() -> void:
	for e in S.ents.values():
		remove(e)

func attack(id: int) -> void:
	var e: Ent = S.ents.get(id)
	if e and e.model:
		e.model.attack()

func cast_anim(id: int, on: bool) -> void:
	var e: Ent = S.ents.get(id)
	if e and e.model:
		e.model.cast(on)

# ---------------------------------------------------------------- markers

func _refresh_marker(e: Ent) -> void:
	var stt: String = S.marks.get(e.id, "")
	if e.mark_state == stt:
		return
	e.mark_state = stt
	if e.marker:
		e.marker.queue_free()
		e.marker = null
	if stt == "":
		return
	var l := Label3D.new()
	l.text = "!" if stt == "available" else "?"
	l.font_size = 140
	l.outline_size = 36
	l.pixel_size = 0.0075
	l.billboard = BaseMaterial3D.BILLBOARD_ENABLED
	l.modulate = Color("9a9a9a") if stt == "progress" else Color("ffd01a")
	l.outline_modulate = Color("2a1500")
	l.no_depth_test = false
	l.shaded = false
	l.double_sided = true
	l.position.y = e.model.height + 0.75
	e.marker = l
	e.node.add_child(l)

func refresh_all_markers() -> void:
	for e in S.ents.values():
		if e.k == "npc":
			_refresh_marker(e)

func _refresh_object(e: Ent) -> void:
	var want := false
	if e.q != "":
		for q in S.quests:
			if q.id == e.q and q.status == "active":
				want = true
	if want == e.obj_glow:
		return
	e.obj_glow = want
	_set_sparkle(e, want or e.lt)

func refresh_objects() -> void:
	for e in S.ents.values():
		if e.k == "obj":
			_refresh_object(e)

func _set_sparkle(e: Ent, on: bool) -> void:
	if not on:
		if e.sparkle:
			e.sparkle.queue_free()
			e.sparkle = null
		return
	if e.sparkle:
		return
	var p := CPUParticles3D.new()
	p.amount = 18
	p.lifetime = 1.5
	p.direction = Vector3.UP
	p.spread = 25.0
	p.initial_velocity_min = 0.6
	p.initial_velocity_max = 1.1
	p.gravity = Vector3.ZERO
	p.emission_shape = CPUParticles3D.EMISSION_SHAPE_SPHERE
	p.emission_sphere_radius = 0.6
	p.scale_amount_min = 0.5
	p.scale_amount_max = 1.0
	var qm := QuadMesh.new()
	qm.size = Vector2(0.4, 0.4)
	var mat := StandardMaterial3D.new()
	mat.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
	mat.billboard_mode = BaseMaterial3D.BILLBOARD_ENABLED
	mat.blend_mode = BaseMaterial3D.BLEND_MODE_ADD
	mat.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
	mat.albedo_texture = glow_tex
	mat.albedo_color = Color(1, 0.93, 0.66)
	mat.vertex_color_use_as_albedo = true
	mat.no_depth_test = false
	qm.material = mat
	p.mesh = qm
	p.position.y = 0.3
	p.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	e.node.add_child(p)
	e.sparkle = p

# ---------------------------------------------------------------- nameplates

func _name_color(e: Ent) -> Color:
	if e.k == "player":
		return UI.CLASS_COLORS.get(e.c, Color.WHITE)
	if e.k == "npc":
		return Color("ffe27a")
	if e.k == "pet":
		return Color("9fe89f")
	if e.k == "ship":
		return Color("9fd8ff")
	if e.k == "obj":
		return Color("ffe14a") if e.obj_glow else Color("cccccc")
	return Color("888888") if e.d else Color("ffe14a")

func _refresh_plate(e: Ent) -> void:
	if plates == null:
		return
	var is_t: bool = S.me.tg == e.id
	var show_hp: bool = is_t and not e.d and (e.k == "mob" or e.k == "pet")
	var key := "%s|%s|%s|%d|%s|%s|%s|%d|%d" % [e.n, e.ti, str(e.d), e.l, str(show_hp), str(e.el), str(e.bs), e.hp if show_hp else 0, S.me.lvl]
	if e.plate == null:
		var p := Control.new()
		p.mouse_filter = Control.MOUSE_FILTER_IGNORE
		p.size = Vector2(160, 44)
		var nm := Label.new()
		nm.name = "Name"
		nm.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		nm.size = Vector2(160, 20)
		nm.position = Vector2(0, 0)
		nm.add_theme_font_size_override("font_size", 14)
		nm.add_theme_constant_override("outline_size", 5)
		nm.add_theme_color_override("font_outline_color", Color(0, 0, 0, 0.9))
		nm.mouse_filter = Control.MOUSE_FILTER_IGNORE
		p.add_child(nm)
		var sub := Label.new()
		sub.name = "Sub"
		sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		sub.size = Vector2(160, 16)
		sub.position = Vector2(0, 17)
		sub.add_theme_font_size_override("font_size", 11)
		sub.add_theme_constant_override("outline_size", 4)
		sub.add_theme_color_override("font_outline_color", Color(0, 0, 0, 0.9))
		sub.mouse_filter = Control.MOUSE_FILTER_IGNORE
		p.add_child(sub)
		var bar := ColorRect.new()
		bar.name = "Bar"
		bar.color = Color(0, 0, 0, 0.7)
		bar.size = Vector2(70, 7)
		bar.position = Vector2(45, 33)
		bar.mouse_filter = Control.MOUSE_FILTER_IGNORE
		var fill := ColorRect.new()
		fill.name = "Fill"
		fill.color = Color("d63a2c")
		fill.position = Vector2(1, 1)
		fill.size = Vector2(68, 5)
		fill.mouse_filter = Control.MOUSE_FILTER_IGNORE
		bar.add_child(fill)
		p.add_child(bar)
		plates.add_child(p)
		e.plate = p
		e.plate_key = ""
	if key == e.plate_key and not show_hp:
		return
	e.plate_key = key
	var nm: Label = e.plate.get_node("Name")
	var sub: Label = e.plate.get_node("Sub")
	var bar: ColorRect = e.plate.get_node("Bar")
	var lvl := ""
	if (e.k == "mob" and not e.d) or e.k == "pet":
		lvl = " " + str(e.l)
	nm.text = e.n + lvl
	nm.add_theme_color_override("font_color", _name_color(e))
	var subt := ""
	if e.k == "npc":
		subt = "<%s>" % e.ti
	elif e.k == "pet":
		subt = "Begleiter"
	elif e.k == "ship":
		subt = e.ti
	elif e.k == "mob" and e.el:
		subt = "Boss" if e.bs else "Elite"
	sub.text = subt
	sub.add_theme_color_override("font_color", Color("ffb040") if e.el else Color("c8c8c8"))
	bar.visible = show_hp
	if show_hp:
		(bar.get_node("Fill") as ColorRect).size.x = 68.0 * clampf(float(e.hp) / maxf(1.0, float(e.mh)), 0.0, 1.0)

# ---------------------------------------------------------------- per frame

func update(dt: float) -> void:
	t_total += dt
	var k := 1.0 - exp(-dt * 11.0)
	var vp_size := get_viewport().get_visible_rect().size
	var show_plates: bool = Settings.v["nameplates"]
	for e: Ent in S.ents.values():
		var is_self: bool = e.id == S.id
		if (S.me.tg == e.id) != e.was_target:
			e.was_target = S.me.tg == e.id
			e.plate_key = ""
			_refresh_plate(e)
		if is_self and not S.me.get("fl", false):
			e.speed = S.local.speed
			e.airborne = S.local.airborne
			e.rx = S.local.x; e.ry = S.local.y; e.rz = S.local.z
			e.rr = lerp_angle(e.rr, S.local.r, minf(1.0, dt * 14.0))
		else:
			var ox := e.rx
			var oz := e.rz
			e.rx += (e.x - e.rx) * k
			e.rz += (e.z - e.rz) * k
			e.ry += (e.y - e.ry) * k
			var sp := Vector2(e.rx - ox, e.rz - oz).length() / maxf(dt, 0.001)
			e.speed += (sp / (5.0 if e.k == "mob" else 7.0) - e.speed) * minf(1.0, dt * 8.0)
			e.airborne = e.y - e.ry > 0.6 and not e.fl
			e.rr = lerp_angle(e.rr, e.r, minf(1.0, dt * 10.0))
			if is_self and S.me.get("fl", false):
				S.local.x = e.rx; S.local.y = e.ry; S.local.z = e.rz; S.local.r = e.rr; S.local.speed = 0.0
		e.node.position = Vector3(e.rx, e.ry, e.rz)
		e.node.rotation.y = e.rr
		var stt := {"speed": 0.0 if e.fl else e.speed, "airborne": e.airborne, "dead": e.d, "flying": e.fl}
		e.model.update(dt, stt)
		if e.mount:
			e.mount.update(dt, stt)
		if e.pad:
			e.pad.update(dt, {"speed": 0.0, "flying": false})
		if e.marker:
			e.marker.position.y = e.model.height + 0.75 + sin(t_total * 3.0) * 0.08
		# nameplate
		if e.plate:
			var far := Vector2(e.rx - S.local.x, e.rz - S.local.z).length() > (32.0 if e.k == "mob" else 14.0 if e.k == "obj" else 55.0)
			var hide: bool = (not show_plates) or far or (e.k == "mob" and e.d) or (is_self and S.me.get("fl", false)) or camera == null
			if not hide:
				var wp := Vector3(e.rx, e.ry + (e.model.height + (0.2 if e.marker else 0.25)) * e.sc + (2.0 if e.fl else 0.0), e.rz)
				if camera.is_position_behind(wp):
					hide = true
				else:
					var sp2 := camera.unproject_position(wp)
					e.plate.position = sp2 - Vector2(80, 30)
			e.plate.visible = not hide
	# selection ring
	var t: Ent = S.ents.get(S.me.tg)
	if t and t.node:
		ring.visible = true
		ring.position = Vector3(t.rx, t.ry + 0.08, t.rz)
		var col := Color("40ff40")
		if t.k == "mob":
			col = Color("888888") if t.d else Color("ff9a20") if t.el else Color("ffe000")
		elif t.k == "obj":
			col = Color("40c8ff")
		(ring.material_override as StandardMaterial3D).albedo_color = Color(col.r, col.g, col.b, 0.85)
		ring.rotation.y += dt * 0.8
		var s := (0.85 if t.model.kind == "human" else maxf(0.9, t.model.height * 0.75)) * t.sc
		ring.scale = Vector3(s, 0.12, s)
	else:
		ring.visible = false

# ---------------------------------------------------------------- picking

## Returns the id of the entity under the screen position (0 if none).
func pick(pos: Vector2) -> int:
	if camera == null:
		return 0
	var o := camera.project_ray_origin(pos)
	var dvec := camera.project_ray_normal(pos)
	var best := 0
	var best_t := 1e9
	for e: Ent in S.ents.values():
		if e.k == "obj" and e.gn:
			continue
		var rad := (6.0 if e.k == "ship" else 0.75) * maxf(1.0, e.sc * 0.8)
		var a := Vector3(e.rx, e.ry, e.rz)
		var b := Vector3(e.rx, e.ry + maxf(1.0, e.height()), e.rz)
		var res := _ray_seg(o, dvec, a, b)
		if res.x < rad and res.y > 0.0 and res.y < best_t and res.y < 120.0:
			best_t = res.y
			best = e.id
	return best

## Distance between a ray and a segment (x) and the ray parameter of the closest point (y).
func _ray_seg(o: Vector3, d: Vector3, a: Vector3, b: Vector3) -> Vector2:
	var u := d
	var v := b - a
	var w := o - a
	var uu := u.dot(u)
	var uv := u.dot(v)
	var vv := v.dot(v)
	var uw := u.dot(w)
	var vw := v.dot(w)
	var den := uu * vv - uv * uv
	var sc_ := 0.0
	var tc := 0.0
	if den < 0.0001:
		tc = clampf(vw / maxf(vv, 0.0001), 0.0, 1.0)
		sc_ = 0.0
	else:
		tc = clampf((uu * vw - uv * uw) / den, 0.0, 1.0)
		sc_ = (uv * tc - uw) / uu
	var pc := a + v * tc
	var ps := o + u * maxf(sc_, 0.0)
	return Vector2((pc - ps).length(), sc_)

func screen_pos(id: int, h_factor: float) -> Variant:
	var e: Ent = S.ents.get(id)
	if e == null or camera == null:
		return null
	var wp := Vector3(e.rx, e.ry + e.height() * h_factor, e.rz)
	if camera.is_position_behind(wp):
		return null
	return camera.unproject_position(wp)

func world_pos(id: int, f := 0.6) -> Variant:
	var e: Ent = S.ents.get(id)
	if e == null:
		return null
	return Vector3(e.rx, e.ry + e.height() * f, e.rz)
