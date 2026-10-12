class_name Controller
extends Node
## Local player: movement with client-side collision, third-person camera, mouse picking, target handling.

signal action(name: String)
signal player_menu(ent_id: int, pos: Vector2)

const RUN := 7.0
const BACK := 4.2
const WALK := 3.2
const JUMP := 8.2
const GRAVITY := 24.0

var camera: Camera3D
var world: GameWorld
var entities: Entities
var dungeon = null               # Dungeon while inside an instance

var cam_yaw := 0.0
var cam_pitch := 0.35
var cam_dist := 9.0
var cam_dist_cur := 9.0
var vy := 0.0
var grounded := true
var autorun := false
var lmb := false
var rmb := false
var down_pos := Vector2.ZERO
var down_acc := 0.0
var down_time := 0
var down_btn := -1
var last_send := 0
var last_sent := {"x": 0.0, "y": 0.0, "z": 0.0, "r": 0.0}
var cam_pos := Vector3.ZERO
var deck: Ent = null
var _ab_id := 0
var _ab_x := 0.0
var _ab_z := 0.0
var hover_t := 0
var cursor_over_ent := 0
var _was_grounded := true

func _ready() -> void:
	add_to_group("controller")

func spawn(x: float, y: float, z: float, r: float) -> void:
	S.local.x = x; S.local.y = y; S.local.z = z; S.local.r = r
	S.local.speed = 0.0; S.local.airborne = false
	cam_yaw = r
	vy = 0.0
	grounded = true
	cam_pos = Vector3(x, y + 4, z)

func correct(x: float, y: float, z: float) -> void:
	S.local.x = x; S.local.y = y; S.local.z = z
	vy = 0.0
	grounded = true

func _ground(x: float, z: float) -> float:
	if dungeon != null:
		return 0.0
	return WD.gh(x, z)

# ------------------------------------------------------------- input

func _input(ev: InputEvent) -> void:
	if not S.in_world:
		return
	if ev is InputEventMouseButton:
		var mb := ev as InputEventMouseButton
		if mb.pressed:
			if _over_ui():
				return
			if mb.button_index == MOUSE_BUTTON_LEFT or mb.button_index == MOUSE_BUTTON_RIGHT:
				if mb.button_index == MOUSE_BUTTON_LEFT:
					lmb = true
				else:
					rmb = true
				down_pos = mb.position
				down_acc = 0.0
				down_time = Time.get_ticks_msec()
				down_btn = mb.button_index
				if rmb and not S.me.d:
					Input.mouse_mode = Input.MOUSE_MODE_CAPTURED
				get_viewport().gui_release_focus()
			elif mb.button_index == MOUSE_BUTTON_WHEEL_UP:
				cam_dist = clampf(cam_dist * 0.9, 2.0, 22.0)
			elif mb.button_index == MOUSE_BUTTON_WHEEL_DOWN:
				cam_dist = clampf(cam_dist * 1.1, 2.0, 22.0)
		else:
			if mb.button_index == MOUSE_BUTTON_LEFT and lmb:
				lmb = false
				if down_btn == MOUSE_BUTTON_LEFT and down_acc < 5.0:
					_click(down_pos, false)
			if mb.button_index == MOUSE_BUTTON_RIGHT and rmb:
				rmb = false
				Input.mouse_mode = Input.MOUSE_MODE_VISIBLE
				if down_btn == MOUSE_BUTTON_RIGHT and down_acc < 6.0 and Time.get_ticks_msec() - down_time < 500:
					_click(down_pos, true)
	elif ev is InputEventMouseMotion:
		var mm := ev as InputEventMouseMotion
		if not (lmb or rmb):
			_hover_cursor(mm.position)
		if lmb or rmb:
			down_acc += absf(mm.relative.x) + absf(mm.relative.y)
			if rmb or down_acc > 4.0:
				var sens := 0.0045 * float(Settings.v["mouse_sens"])
				cam_yaw -= mm.relative.x * sens
				var inv := -1.0 if Settings.v["invert_y"] else 1.0
				cam_pitch = clampf(cam_pitch + mm.relative.y * sens * inv, -0.35, 1.45)

## Pointing hand over NPCs, objects and lootable corpses, crosshair over attackable creatures.
func _hover_cursor(pos: Vector2) -> void:
	var t := Time.get_ticks_msec()
	if t - hover_t < 80:
		return
	hover_t = t
	var shape := Input.CURSOR_ARROW
	if not _over_ui():
		var e: Ent = S.ents.get(entities.pick(pos))
		if e:
			if e.k == "npc" or e.k == "obj" or (e.k == "mob" and e.d and e.lt):
				shape = Input.CURSOR_POINTING_HAND
			elif e.k == "mob" and not e.d:
				shape = Input.CURSOR_CROSS
	Input.set_default_cursor_shape(shape)

func _over_ui() -> bool:
	var h := get_viewport().gui_get_hovered_control()
	return h != null and h.mouse_filter != Control.MOUSE_FILTER_IGNORE

func _unhandled_key_input(ev: InputEvent) -> void:
	if not S.in_world or S.typing:
		return
	var k := ev as InputEventKey
	if k == null or not k.pressed or k.echo:
		return
	for a in ["toggle_bags", "toggle_quests", "toggle_char", "toggle_spells", "toggle_map", "toggle_party", "menu", "mute", "screenshot", "fullscreen"]:
		if ev.is_action_pressed(a):
			action.emit(a)
			get_viewport().set_input_as_handled()
			return
	if k.unicode == 47:   # "/" on any keyboard layout
		action.emit("chat_cmd")
		get_viewport().set_input_as_handled()
	elif ev.is_action_pressed("chat"):
		action.emit("chat")
		get_viewport().set_input_as_handled()
	elif ev.is_action_pressed("chat_cmd"):
		action.emit("chat_cmd")
		get_viewport().set_input_as_handled()
	elif ev.is_action_pressed("autorun"):
		autorun = not autorun
	elif ev.is_action_pressed("interact"):
		interact_target()
	elif ev.is_action_pressed("target_next"):
		target_next(k.shift_pressed)
		get_viewport().set_input_as_handled()
	else:
		for i in Settings.SLOT_KEYS.size():
			if ev.is_action_pressed("slot_%d" % i):
				action.emit("slot_%d" % i)
				return

func _click(pos: Vector2, right: bool) -> void:
	var id := entities.pick(pos)
	if id == 0:
		return
	var e: Ent = S.ents.get(id)
	if e == null:
		return
	if S.me.tg != id:
		set_target(id)
	if not right:
		return
	if e.k == "npc" or e.k == "obj" or (e.k == "mob" and e.d):
		Net.send("interact", {"id": id})
	elif e.k == "player" and e.id != S.id:
		player_menu.emit(id, pos)
	elif e.k == "mob" and not e.d:
		action.emit("slot_0")

func interact_target() -> void:
	var e: Ent = S.ents.get(S.me.tg)
	if e and (e.k == "npc" or e.k == "obj" or (e.k == "mob" and e.d)):
		Net.send("interact", {"id": e.id})

func set_target(id: int) -> void:
	S.me.tg = id
	Net.send("target", {"id": id})

func target_next(back: bool) -> void:
	var list: Array = []
	for e: Ent in S.ents.values():
		if e.k == "mob" and not e.d and Vector2(e.rx - S.local.x, e.rz - S.local.z).length() < 45.0:
			list.append(e)
	if list.is_empty():
		return
	list.sort_custom(func(a: Ent, b: Ent): return Vector2(a.rx - S.local.x, a.rz - S.local.z).length() < Vector2(b.rx - S.local.x, b.rz - S.local.z).length())
	var i := -1
	for j in list.size():
		if list[j].id == S.me.tg:
			i = j
	if back:
		i = list.size() - 1 if i <= 0 else i - 1
	else:
		i = (i + 1) % list.size()
	set_target(list[i].id)

func face_target() -> void:
	var t: Ent = S.ents.get(S.me.tg)
	if t and t.k == "mob":
		S.local.r = atan2(t.rx - S.local.x, t.rz - S.local.z)

# ------------------------------------------------------------- per frame

func update(dt: float) -> void:
	var L: Dictionary = S.local
	if S.me.get("fl", false):
		_update_camera(dt)
		return
	var dead: bool = S.me.d
	var typing: bool = S.typing
	var fw := 0.0
	var rt := 0.0
	if not dead and not typing:
		fw = Input.get_action_strength("move_forward") - Input.get_action_strength("move_back")
		rt = Input.get_action_strength("move_right") - Input.get_action_strength("move_left")
		if lmb and rmb or autorun:
			fw += 1.0
		# gamepad: left stick moves, right stick turns the camera
		var jx := Input.get_joy_axis(0, JOY_AXIS_LEFT_X)
		var jy := -Input.get_joy_axis(0, JOY_AXIS_LEFT_Y)
		if absf(jx) > 0.2: rt += jx
		if absf(jy) > 0.2: fw += jy
		var cx := Input.get_joy_axis(0, JOY_AXIS_RIGHT_X)
		var cy := Input.get_joy_axis(0, JOY_AXIS_RIGHT_Y)
		if absf(cx) > 0.2: cam_yaw -= cx * dt * 2.4
		if absf(cy) > 0.2: cam_pitch = clampf(cam_pitch + cy * dt * 1.8, -0.35, 1.45)
	if Input.is_action_pressed("move_forward") or Input.is_action_pressed("move_back"):
		autorun = false
	fw = clampf(fw, -1.0, 1.0)
	rt = clampf(rt, -1.0, 1.0)
	var fx := sin(cam_yaw)
	var fz := cos(cam_yaw)
	var dx := fx * fw + -fz * rt
	var dz := fz * fw + fx * rt
	var length := sqrt(dx * dx + dz * dz)
	var speed := 0.0
	if length > 0.001:
		dx /= length
		dz /= length
		speed = BACK if (fw < 0.0 and rt == 0.0) else (WALK if Input.is_action_pressed("walk") else RUN)
		if fw < 0.0 and rt != 0.0:
			speed = RUN * 0.8
		var stick := Vector2(rt, fw).length()
		if stick < 0.6 and not Input.is_action_pressed("move_forward") and not Input.is_action_pressed("move_back") and not Input.is_action_pressed("move_left") and not Input.is_action_pressed("move_right") and not autorun:
			speed = WALK
	if not dead:
		var strafe := rmb and length > 0.001 and rt != 0.0 and fw >= 0.0
		if rmb and not strafe:
			L.r = cam_yaw
		elif length > 0.001:
			var want := atan2(dx, dz)
			L.r += wrapf(want - L.r, -PI, PI) * minf(1.0, dt * 14.0)
	# ships: the deck moves with the ship
	var ship: Ent = S.ents.get(S.me.get("ab", 0)) if int(S.me.get("ab", 0)) != 0 else null
	if ship:
		if _ab_id == ship.id:
			L.x += ship.rx - _ab_x
			L.z += ship.rz - _ab_z
		_ab_id = ship.id
		_ab_x = ship.rx
		_ab_z = ship.rz
	else:
		_ab_id = 0
	var nx: float = L.x + dx * speed * dt
	var nz: float = L.z + dz * speed * dt
	var on_deck: Ent = null
	if dungeon == null:
		var cand: Ent = ship if ship else _near_ship(L.x, L.z)
		if cand and ((ship and cand.fl) or _on_deck(cand, nx, nz, 0.3)):
			on_deck = cand
	deck = on_deck
	if on_deck:
		var c := _clamp_deck(on_deck, nx, nz, 0.0 if on_deck.fl else 0.3)
		nx = c.x; nz = c.y
	elif speed > 0.0:
		var c2: Vector2 = dungeon.collide(nx, nz) if dungeon != null else WD.collide(nx, nz, L.x, L.z)
		nx = c2.x; nz = c2.y
	var moved := Vector2(nx - L.x, nz - L.z).length() / maxf(dt, 0.0001)
	L.x = nx
	L.z = nz
	L.speed = moved / RUN
	# vertical
	var ground: float = (on_deck.ry + 1.2) if on_deck else _ground(L.x, L.z)
	if grounded:
		if Input.is_action_pressed("jump") and not dead and not typing:
			vy = JUMP
			grounded = false
		elif L.y - ground > 0.35:
			grounded = false
			vy = 0.0
		else:
			L.y = ground
	if not grounded:
		vy -= GRAVITY * dt
		L.y += vy * dt
		if L.y <= ground and vy <= 0.0:
			L.y = ground
			vy = 0.0
			grounded = true
	L.airborne = not grounded
	Sound.footsteps(dt, L.speed, grounded and not dead, L.x, L.z)
	if grounded != _was_grounded:
		if grounded and not _was_grounded:
			Sound.sfx("land", -8.0)
		elif not grounded and vy > 0.0:
			Sound.sfx("jump", -6.0)
		_was_grounded = grounded
	_send_move()
	_update_camera(dt)

func _to_deck(ship: Ent, x: float, z: float) -> Vector2:
	var s := sin(ship.rr)
	var c := cos(ship.rr)
	var wx := x - ship.rx
	var wz := z - ship.rz
	return Vector2(wx * c - wz * s, wx * s + wz * c)

func _on_deck(ship: Ent, x: float, z: float, margin: float) -> bool:
	var l := _to_deck(ship, x, z)
	return absf(l.x) <= 1.5 + margin and absf(l.y) <= 6.0 + margin

func _near_ship(x: float, z: float) -> Ent:
	for e: Ent in S.ents.values():
		if e.k == "ship" and Vector2(e.rx - x, e.rz - z).length() < 25.0:
			return e
	return null

func _clamp_deck(ship: Ent, x: float, z: float, margin: float) -> Vector2:
	var l := _to_deck(ship, x, z)
	var lx := clampf(l.x, -1.5 - margin, 1.5 + margin)
	var lz := clampf(l.y, -6.0 - margin, 6.0 + margin)
	var s := sin(ship.rr)
	var c := cos(ship.rr)
	return Vector2(ship.rx + lx * c + lz * s, ship.rz - lx * s + lz * c)

func _send_move() -> void:
	var L: Dictionary = S.local
	var t := Time.get_ticks_msec()
	if t - last_send < 50:
		return
	var s := last_sent
	var changed: bool = absf(L.x - s.x) > 0.01 or absf(L.z - s.z) > 0.01 or absf(L.y - s.y) > 0.01 or absf(L.r - s.r) > 0.02
	if not changed and t - last_send < 600:
		return
	last_send = t
	s.x = L.x; s.y = L.y; s.z = L.z; s.r = L.r
	var msg := {"x": snappedf(L.x, 0.01), "y": snappedf(L.y, 0.01), "z": snappedf(L.z, 0.01), "r": snappedf(L.r, 0.001)}
	if int(S.me.get("ab", 0)) != 0 and deck != null:
		var l := _to_deck(deck, L.x, L.z)
		msg["ab"] = true
		msg["lx"] = snappedf(l.x, 0.01)
		msg["lz"] = snappedf(l.y, 0.01)
	Net.send("move", msg)

func _update_camera(dt: float) -> void:
	var L: Dictionary = S.local
	var tx: float = L.x
	var ty: float = L.y + 1.7
	var tz: float = L.z
	var cp := cos(cam_pitch)
	var sp := sin(cam_pitch)
	# pull the camera in when terrain or walls are in the way
	var want := cam_dist
	var min_d := 1.2
	var d := want
	if dungeon == null:
		var steps := 10
		for i in range(steps, 0, -1):
			var dd := want * float(i) / steps
			var px := tx - sin(cam_yaw) * cp * dd
			var pz := tz - cos(cam_yaw) * cp * dd
			var py := ty + sp * dd
			if py > WD.height_at(px, pz) + 0.45:
				d = dd
				break
			d = min_d
	cam_dist_cur = d if d < cam_dist_cur else lerpf(cam_dist_cur, d, minf(1.0, dt * 6.0))
	var cx := tx - sin(cam_yaw) * cp * cam_dist_cur
	var cz := tz - cos(cam_yaw) * cp * cam_dist_cur
	var cy := ty + sp * cam_dist_cur
	if dungeon != null:
		var c: Vector2 = dungeon.collide(cx, cz)
		cx = c.x; cz = c.y
	else:
		cy = maxf(cy, WD.height_at(cx, cz) + 0.5)
	var target := Vector3(cx, cy, cz)
	if Settings.v["camera_smooth"] and cam_pos.distance_to(target) < 6.0:
		cam_pos = cam_pos.lerp(target, minf(1.0, dt * 28.0))
	else:
		cam_pos = target
	camera.position = cam_pos
	camera.look_at(Vector3(tx, ty + 0.1, tz), Vector3.UP)
	camera.fov = float(Settings.v["fov"])
	if dungeon != null:
		dungeon.update(Time.get_ticks_msec() / 1000.0, L.x, L.z)
	elif world:
		world.follow(Vector3(tx, L.y, tz))
