class_name FX
extends Node3D
## Combat and spell effects (projectiles, bursts, rings, beams) and floating damage numbers.

var ents: Entities
var hud: Hud
var list: Array = []
var tex: GradientTexture2D
var streak_tex: GradientTexture2D

func _ready() -> void:
	tex = GradientTexture2D.new()
	tex.fill = GradientTexture2D.FILL_RADIAL
	tex.fill_from = Vector2(0.5, 0.5)
	tex.fill_to = Vector2(1.0, 0.5)
	tex.width = 64
	tex.height = 64
	var g := Gradient.new()
	g.colors = PackedColorArray([Color(1, 1, 1, 1), Color(1, 1, 1, 0.55), Color(1, 1, 1, 0)])
	g.offsets = PackedFloat32Array([0.0, 0.4, 1.0])
	tex.gradient = g
	streak_tex = GradientTexture2D.new()
	streak_tex.fill = GradientTexture2D.FILL_LINEAR
	streak_tex.fill_from = Vector2(0, 0.5)
	streak_tex.fill_to = Vector2(1, 0.5)
	streak_tex.width = 64
	streak_tex.height = 8
	var g2 := Gradient.new()
	g2.colors = PackedColorArray([Color(1, 1, 1, 0), Color(1, 1, 1, 1), Color(1, 1, 1, 0)])
	g2.offsets = PackedFloat32Array([0.0, 0.5, 1.0])
	streak_tex.gradient = g2

func clear_all() -> void:
	for c in get_children():
		c.queue_free()
	list.clear()

func _mat(color: Color, t: Texture2D = null, billboard := true, opacity := 1.0) -> StandardMaterial3D:
	var m := StandardMaterial3D.new()
	m.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
	m.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
	m.blend_mode = BaseMaterial3D.BLEND_MODE_ADD
	m.cull_mode = BaseMaterial3D.CULL_DISABLED
	m.albedo_color = Color(color.r, color.g, color.b, opacity)
	if t:
		m.albedo_texture = t
	if billboard:
		m.billboard_mode = BaseMaterial3D.BILLBOARD_ENABLED
	return m

func sprite(color: Color, size: float, opacity := 1.0) -> MeshInstance3D:
	var mi := MeshInstance3D.new()
	var q := QuadMesh.new()
	q.size = Vector2.ONE
	mi.mesh = q
	mi.material_override = _mat(color, tex, true, opacity)
	mi.scale = Vector3.ONE * size
	mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	add_child(mi)
	return mi

func _set_alpha(mi: MeshInstance3D, a: float) -> void:
	var m := mi.material_override as StandardMaterial3D
	m.albedo_color.a = a

func burst(at: Vector3, color: Color, size := 1.4, dur := 0.35) -> void:
	var s := sprite(color, size * 0.3)
	s.position = at
	var st := {"t": 0.0}
	list.append(func(dt: float) -> bool:
		st.t += dt
		var k: float = st.t / dur
		s.scale = Vector3.ONE * size * (0.3 + k)
		_set_alpha(s, maxf(0.0, 1.0 - k))
		if k >= 1.0:
			s.queue_free()
			return false
		return true)

func ring(at: Vector3, color: Color, r0: float, r1: float, dur: float, vertical := 0.0) -> void:
	var mi := MeshInstance3D.new()
	var tm := TorusMesh.new()
	tm.inner_radius = 0.9
	tm.outer_radius = 1.0
	tm.rings = 32
	tm.ring_segments = 4
	mi.mesh = tm
	mi.material_override = _mat(color, null, false, 0.9)
	mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	mi.position = at
	add_child(mi)
	var st := {"t": 0.0}
	list.append(func(dt: float) -> bool:
		st.t += dt
		var k := minf(1.0, st.t / dur)
		var r := r0 + (r1 - r0) * k
		mi.scale = Vector3(r, 0.1, r)
		mi.position.y = at.y + vertical * k
		_set_alpha(mi, 0.9 * (1.0 - k))
		if k >= 1.0:
			mi.queue_free()
			return false
		return true)

func projectile(from: Vector3, dst_id: int, style: String, color_s: String) -> void:
	var col := Color(color_s) if color_s != "" else Color("c070ff")
	var node := Node3D.new()
	add_child(node)
	node.position = from
	if style == "arrow":
		var shaft := MeshInstance3D.new()
		var bm := BoxMesh.new()
		bm.size = Vector3(0.04, 0.04, 0.7)
		shaft.mesh = bm
		var sm := StandardMaterial3D.new()
		sm.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
		sm.albedo_color = Color("e8d8a8")
		shaft.material_override = sm
		node.add_child(shaft)
		var tip := MeshInstance3D.new()
		var tm := BoxMesh.new()
		tm.size = Vector3(0.1, 0.1, 0.16)
		tip.mesh = tm
		tip.position.z = 0.4
		var tmat := StandardMaterial3D.new()
		tmat.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
		tmat.albedo_color = Color("cccccc")
		tip.material_override = tmat
		node.add_child(tip)
	else:
		var core := MeshInstance3D.new()
		var sp := SphereMesh.new()
		var r := 0.3 if style == "fire" else 0.2
		sp.radius = r
		sp.height = r * 2.0
		sp.radial_segments = 8
		sp.rings = 4
		core.mesh = sp
		var cm := StandardMaterial3D.new()
		cm.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
		cm.albedo_color = col
		core.material_override = cm
		node.add_child(core)
		var halo := sprite(col, 1.6 if style == "fire" else 1.1)
		remove_child(halo)
		node.add_child(halo)
	var speed := 42.0 if style == "arrow" else (22.0 if style == "fire" else 28.0)
	list.append(func(dt: float) -> bool:
		var to = ents.world_pos(dst_id, 0.6)
		if to == null:
			node.queue_free()
			return false
		var d: Vector3 = to - node.position
		var step := speed * dt
		if d.length() <= step + 0.3:
			node.queue_free()
			burst(to, col, 2.2 if style == "fire" else 1.4)
			return false
		d = d.normalized()
		node.position += d * step
		if style == "arrow":
			node.look_at(to, Vector3.UP)
		return true)

func slash(at: Vector3, color: Color) -> void:
	var mi := MeshInstance3D.new()
	var q := QuadMesh.new()
	q.size = Vector2(1.7, 0.2)
	mi.mesh = q
	mi.material_override = _mat(color, streak_tex, false, 1.0)
	mi.position = at
	mi.rotation = Vector3(randf() * 0.8, randf() * TAU, randf() * 3.0)
	mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	add_child(mi)
	var st := {"t": 0.0}
	list.append(func(dt: float) -> bool:
		st.t += dt
		mi.scale = Vector3.ONE * (1.0 + st.t * 4.0)
		_set_alpha(mi, 1.0 - st.t / 0.25)
		if st.t > 0.25:
			mi.queue_free()
			return false
		return true)
	burst(at, color, 1.0, 0.25)

func rising(id: int, color: Color, dur := 1.0) -> void:
	if int(Settings.v["particles"]) == 0:
		return
	for i in 10:
		var s := sprite(color, 0.4, 0.9)
		var a := float(i) / 10.0 * TAU
		var st := {"t": -i * 0.05}
		list.append(func(dt: float) -> bool:
			st.t += dt
			if st.t < 0.0:
				return true
			var p = ents.world_pos(id, 0.0)
			if p == null:
				s.queue_free()
				return false
			s.position = Vector3(p.x + cos(a + st.t * 4.0) * 0.5, p.y + st.t * 2.2, p.z + sin(a + st.t * 4.0) * 0.5)
			_set_alpha(s, 0.9 * maxf(0.0, 1.0 - st.t / dur))
			if st.t > dur:
				s.queue_free()
				return false
			return true)

func level_up(id: int) -> void:
	var p = ents.world_pos(id, 0.0)
	if p == null:
		return
	var mi := MeshInstance3D.new()
	var cm := CylinderMesh.new()
	cm.top_radius = 0.9
	cm.bottom_radius = 0.9
	cm.height = 7.0
	cm.radial_segments = 12
	cm.cap_top = false
	cm.cap_bottom = false
	mi.mesh = cm
	mi.material_override = _mat(Color("ffdd55"), null, false, 0.7)
	mi.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	add_child(mi)
	var st := {"t": 0.0}
	list.append(func(dt: float) -> bool:
		st.t += dt
		var q = ents.world_pos(id, 0.0)
		if q != null:
			mi.position = Vector3(q.x, q.y + 3.5, q.z)
		_set_alpha(mi, 0.7 * maxf(0.0, 1.0 - st.t / 1.8))
		mi.scale = Vector3(1.0 + st.t * 0.2, 1.0, 1.0 + st.t * 0.2)
		if st.t > 1.8:
			mi.queue_free()
			return false
		return true)
	rising(id, Color("ffe27a"), 1.6)
	ring(Vector3(p.x, p.y + 0.1, p.z), Color("ffe27a"), 0.5, 3.5, 1.0, 1.5)

# ---------------------------------------------------------------- network events

func _later(sec: float, fn: Callable) -> void:
	get_tree().create_timer(sec).timeout.connect(fn)

func on_fx(m: Dictionary) -> void:
	var k: String = m.k
	var skill := S.skill(k.replace(":cast", ""))
	var src: Ent = S.ents.get(int(m.get("s", 0)))
	var sid := int(m.get("s", 0))
	var did := int(m.get("d", 0))
	if k == "levelup":
		level_up(sid)
		return
	if k == "mobshot":
		ents.attack(sid)
		var from = ents.world_pos(sid, 0.9)
		if from != null:
			var cols := {"frost": "#7ad8ff", "fire": "#ff7a1a", "bolt": "#9aff4a"}
			var st: String = m.get("style", "")
			projectile(from, did, "frost" if st == "frost" else ("fire" if st == "fire" else "bolt"), cols.get(st, "#c070ff"))
		return
	if k == "pet_summon":
		var p = ents.world_pos(sid, 0.3)
		if p != null:
			ring(p, Color("66ccff"), 0.5, 2.2, 0.6)
			rising(sid, Color("66ccff"), 0.8)
		return
	if k == "heal":
		rising(sid, Color("66ff88"), 1.0)
		return
	if k == "bossheavy":
		var p2 = ents.world_pos(did, 0.0)
		if p2 != null:
			ents.attack(sid)
			ring(p2, Color("ff4020"), 0.5, 3.0, 0.45)
			burst(p2 + Vector3(0, 1, 0), Color("ff6030"), 2.6, 0.4)
		return
	if k == "bossnova":
		var p3 = ents.world_pos(sid, 0.0)
		if p3 != null:
			ring(p3 + Vector3(0, 0.1, 0), Color("ff5030"), 0.5, 14.0, 0.8)
		return
	if k == "ground":
		var col := Color(m.get("color", "#aee4ff"))
		var by: float = src.ry if src else 0.0
		var disc := MeshInstance3D.new()
		var cyl := CylinderMesh.new()
		cyl.top_radius = float(m.r)
		cyl.bottom_radius = float(m.r)
		cyl.height = 0.05
		cyl.radial_segments = 24
		disc.mesh = cyl
		disc.material_override = _mat(col, null, false, 0.25)
		disc.position = Vector3(m.x, by + 0.15, m.z)
		disc.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
		add_child(disc)
		var st2 := {"t": 0.0}
		var dur := float(m.dur)
		var rad := float(m.r)
		var px := float(m.x)
		var pz := float(m.z)
		list.append(func(dt: float) -> bool:
			st2.t += dt
			if int(Settings.v["particles"]) > 0:
				for i in 3:
					var q := sprite(col, 0.5, 0.8)
					q.position = Vector3(px + (randf() - 0.5) * rad * 2.0, by + 4.0 + randf() * 2.0, pz + (randf() - 0.5) * rad * 2.0)
					var s2 := {"t": 0.0}
					list.append(func(d2: float) -> bool:
						s2.t += d2
						q.position.y -= d2 * 5.0
						_set_alpha(q, 0.8 * (1.0 - s2.t / 0.9))
						if s2.t > 0.9:
							q.queue_free()
							return false
						return true)
			_set_alpha(disc, 0.25 * maxf(0.0, 1.0 - st2.t / dur))
			if st2.t > dur:
				disc.queue_free()
				return false
			return true)
		return
	if k == "bite":
		ents.attack(sid)
		var to = ents.world_pos(did, 0.6)
		if to != null:
			_later(0.2, func(): burst(to, Color("ff3030"), 1.1, 0.3))
		return
	if skill.is_empty():
		return
	var color := Color(skill.get("color", "#ffffff"))
	if k.ends_with(":cast"):
		ents.cast_anim(sid, true)
		_later(float(skill.get("cast", 1)), func(): ents.cast_anim(sid, false))
		var p4 = ents.world_pos(sid, 0.9)
		if p4 == null:
			return
		var s := sprite(color, 0.2, 0.9)
		var st3 := {"t": 0.0}
		var castt := maxf(0.8, float(skill.get("cast", 1)))
		list.append(func(dt: float) -> bool:
			st3.t += dt
			var q = ents.world_pos(sid, 0.9)
			var e: Ent = S.ents.get(sid)
			if q == null or e == null:
				s.queue_free()
				return false
			s.position = q + Vector3(sin(e.rr) * 0.6, 0, cos(e.rr) * 0.6)
			s.scale = Vector3.ONE * (0.3 + minf(st3.t, 1.5) * 0.6)
			if st3.t > castt:
				s.queue_free()
				return false
			return true)
		return
	if src:
		ents.attack(sid)
	var from2 = ents.world_pos(sid, 0.85)
	var to2 = ents.world_pos(did, 0.6)
	match skill.get("style", ""):
		"melee":
			if to2 != null:
				_later(0.12, func(): slash(to2, color))
		"nova":
			if from2 != null:
				ring(from2 - Vector3(0, 0.7, 0), color, 0.5, float(skill.get("radius", 5)), 0.5)
		"buff":
			rising(sid, color, 1.1)
		_:
			if from2 != null and to2 != null:
				var delay := 0.0
				if skill.get("style", "") == "arrow" and src and src.model and src.model.has_meta("bow_delay"):
					delay = float(src.model.get_meta("bow_delay"))
				if delay > 0.0:
					_later(delay, func():
						var f3 = ents.world_pos(sid, 0.85)
						if f3 != null:
							projectile(f3, did, skill.get("style", ""), skill.get("color", "")))
				else:
					projectile(from2, did, skill.get("style", ""), skill.get("color", ""))

func on_dmg(m: Dictionary) -> void:
	var sp = ents.screen_pos(int(m.id), 0.85)
	if sp == null or hud == null:
		return
	var cls := "other"
	if int(m.get("s", 0)) == S.id:
		cls = "out"
	elif int(m.id) == S.id:
		cls = "in"
	var crit := bool(m.get("c", false))
	if int(m.get("a", 0)) > 0:
		hud.floating_text(sp, str(int(m.a)) + ("!" if crit else ""), cls, crit)
	if int(m.get("ab", 0)) > 0:
		hud.floating_text(sp + Vector2(0, -18), ("(%d absorbiert)" % int(m.ab)) if int(m.get("a", 0)) > 0 else "Absorbiert", "abs")

func update(dt: float) -> void:
	var cur := list
	list = []
	var keep: Array = []
	for f in cur:
		if f.call(dt):
			keep.append(f)
	list = keep + list
