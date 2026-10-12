class_name BlockyModel
extends Node3D
## A model made of blocky parts with simple procedural animation (walking, attacking, casting, dying).

var kind := "human"          # human | quad | spider | scorpion | croc | wisp | golem | object | ship | griffin
var height := 1.85
var parts: Dictionary = {}   # name -> Node3D pivot
var phase := 0.0
var attack_t := -1.0
var casting := false
var death := 0.0             # 0 alive .. 1 lying down
var speed_scale := 1.0       # slow walkers (zombies, golems)
var arms_forward := false
var flying := false
var bob := 0.0
var wing_speed := 9.0
var glow_nodes: Array = []   # nodes that pulse (crystals, wisps)
var t := 0.0
var player_class := ""

func attack() -> void:
	attack_t = 0.0

func cast(on: bool) -> void:
	casting = on

func part(n: String) -> Node3D:
	return parts.get(n)

func update(dt: float, st: Dictionary) -> void:
	t += dt
	var speed: float = st.get("speed", 0.0)
	var dead: bool = st.get("dead", false)
	flying = st.get("flying", false)
	death = move_toward(death, 1.0 if dead else 0.0, dt * 2.5)
	if attack_t >= 0.0:
		attack_t += dt
		if attack_t > 0.45:
			attack_t = -1.0
	phase += dt * (3.0 + speed * 9.0) * speed_scale * (1.0 if speed > 0.02 else 0.0)
	var amp: float = clampf(speed * 1.6, 0.0, 1.0) * 0.8
	var airborne: bool = st.get("airborne", false)
	match kind:
		"human", "golem": _anim_human(dt, amp, airborne)
		"quad": _anim_quad(dt, amp, airborne)
		"spider", "scorpion": _anim_crawler(dt, amp)
		"croc": _anim_croc(dt, amp)
		"wisp": _anim_wisp(dt)
		"griffin": _anim_griffin(dt, st)
		"ship": _anim_ship(dt)
		"object": _anim_object(dt)
	# dying: tip over
	if kind != "wisp" and kind != "object" and kind != "ship":
		rotation.x = death * (PI / 2.0) * (-1.0 if kind == "quad" else 1.0) * 0.98
		position.y = -death * (0.1 if kind == "quad" else 0.5) * 0.0

func _swing(n: String, a: float) -> void:
	var p: Node3D = parts.get(n)
	if p:
		p.rotation.x = a

func _anim_human(dt: float, amp: float, airborne: bool) -> void:
	var s := sin(phase)
	var lift := 0.0
	if airborne:
		_swing("leg_l", -0.5); _swing("leg_r", 0.3)
		lift = 0.0
	else:
		_swing("leg_l", s * amp)
		_swing("leg_r", -s * amp)
	var arm_l := -s * amp * 0.9
	var arm_r := s * amp * 0.9
	if arms_forward:
		arm_l = -1.35 + s * amp * 0.1
		arm_r = -1.35 - s * amp * 0.1
	if casting:
		arm_r = -2.4 + sin(t * 9.0) * 0.08
		arm_l = -2.2 + sin(t * 9.0 + 1.0) * 0.08
	if attack_t >= 0.0:
		var k := attack_t / 0.45
		var sw := sin(k * PI)
		arm_r = -2.1 * sw + 0.2
		var tw: Node3D = parts.get("torso")
		if tw:
			tw.rotation.y = -0.5 * sw
	else:
		var tw2: Node3D = parts.get("torso")
		if tw2:
			tw2.rotation.y = lerpf(tw2.rotation.y, 0.0, minf(1.0, dt * 14.0))
	_swing("arm_l", arm_l)
	_swing("arm_r", arm_r)
	var torso: Node3D = parts.get("torso")
	var head: Node3D = parts.get("head")
	var idle := sin(t * 1.6) * 0.012
	if torso:
		torso.position.y = 0.75 + (absf(s) * amp * 0.05 if not airborne else 0.0) + idle
	if head:
		head.position.y = 1.45 + (absf(s) * amp * 0.05 if not airborne else 0.0) + idle
		head.rotation.x = sin(t * 0.9) * 0.02
	for n in ["arm_l", "arm_r"]:
		var a: Node3D = parts.get(n)
		if a:
			a.position.y = 1.42 + (absf(s) * amp * 0.05 if not airborne else 0.0) + idle
	var cl: Node3D = parts.get("cloak")
	if cl:
		cl.rotation.x = 0.08 + amp * 0.35 + sin(t * 3.0) * 0.02

func _anim_quad(dt: float, amp: float, airborne: bool) -> void:
	var s := sin(phase * 0.85)
	_swing("leg_fl", s * amp)
	_swing("leg_br", s * amp)
	_swing("leg_fr", -s * amp)
	_swing("leg_bl", -s * amp)
	var tail: Node3D = parts.get("tail")
	if tail:
		tail.rotation.y = sin(t * 4.0) * 0.25 * (0.4 + amp)
	var head: Node3D = parts.get("head")
	var body: Node3D = parts.get("body")
	if head:
		head.rotation.x = sin(t * 1.3) * 0.04
		if attack_t >= 0.0:
			head.rotation.x = 0.5 * sin(attack_t / 0.45 * PI)
	if body:
		body.position.z = 0.0
		if attack_t >= 0.0:
			body.position.z = 0.35 * sin(attack_t / 0.45 * PI)
		body.position.y = absf(s) * amp * 0.03

func _anim_crawler(dt: float, amp: float) -> void:
	for i in 8:
		var l: Node3D = parts.get("leg%d" % i)
		if l:
			var side := 1.0 if i < 4 else -1.0
			l.rotation.z = side * (0.15 + sin(phase * 1.4 + i * 1.7) * 0.25 * amp)
			l.rotation.y = sin(phase * 1.4 + i * 1.7 + 1.0) * 0.2 * amp
	var tail: Node3D = parts.get("tail")
	if tail:
		tail.rotation.x = sin(t * 2.0) * 0.1 - 0.1
		if attack_t >= 0.0:
			tail.rotation.x = -0.9 * sin(attack_t / 0.45 * PI)
	var body: Node3D = parts.get("body")
	if body and attack_t >= 0.0:
		body.position.z = 0.3 * sin(attack_t / 0.45 * PI)

func _anim_croc(dt: float, amp: float) -> void:
	var tail: Node3D = parts.get("tail")
	if tail:
		tail.rotation.y = sin(phase * 0.8) * 0.35 * (0.4 + amp)
	var jaw: Node3D = parts.get("jaw")
	if jaw:
		jaw.rotation.x = 0.0 if attack_t < 0.0 else 0.6 * sin(attack_t / 0.45 * PI)
	for n in ["leg_fl", "leg_br"]:
		_swing(n, sin(phase * 0.8) * amp)
	for n in ["leg_fr", "leg_bl"]:
		_swing(n, -sin(phase * 0.8) * amp)

func _anim_wisp(dt: float) -> void:
	position.y = 1.0 + sin(t * 2.2) * 0.25
	for n in glow_nodes:
		if n is Node3D:
			n.rotation.y += dt * 1.5
			n.rotation.x += dt * 0.9

func _anim_griffin(dt: float, st: Dictionary) -> void:
	var fly: bool = st.get("flying", false)
	var w := sin(t * wing_speed) * (0.9 if fly else 0.08)
	var wl: Node3D = parts.get("wing_l")
	var wr: Node3D = parts.get("wing_r")
	if wl:
		wl.rotation.z = 0.3 + w
	if wr:
		wr.rotation.z = -0.3 - w
	var tail: Node3D = parts.get("tail")
	if tail:
		tail.rotation.x = sin(t * 2.0) * 0.1

func _anim_ship(dt: float) -> void:
	rotation.z = sin(t * 0.9) * 0.025
	var sail: Node3D = parts.get("sail")
	if sail:
		sail.rotation.y = sin(t * 0.6) * 0.04

func _anim_object(dt: float) -> void:
	for n in glow_nodes:
		if n is MeshInstance3D:
			n.scale = Vector3.ONE * (1.0 + sin(t * 3.0) * 0.04)
	var sw: Node3D = parts.get("swirl")
	if sw:
		sw.rotation.z += dt * 1.2
