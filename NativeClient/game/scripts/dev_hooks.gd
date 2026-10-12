class_name DevHooks
extends Node
## Developer-only automation used by the tests: screenshots and scripted play (see NativeClient/README.md).
## Enabled with user args: --autologin host user pass, --autoenter, --shot out.png seconds, --open a,b, --play

var main: Node

func _ready() -> void:
	var a := OS.get_cmdline_user_args()
	if a.has("--preset"):
		Settings.apply_preset(a[a.find("--preset") + 1])
	if a.has("--autologin"):
		_autologin.call_deferred(a)
	if a.has("--autoenter"):
		Net.on("chars", _auto_enter)
	if a.has("--shot"):
		_shot.call_deferred(a)
	if a.has("--play"):
		_play.call_deferred(a)
	if a.has("--play2"):
		_play2.call_deferred(a)

func _autologin(a: PackedStringArray) -> void:
	var i := a.find("--autologin")
	main.screens.in_server.text = a[i + 1]
	main.screens.in_user.text = a[i + 2]
	main.screens.in_pass.text = a[i + 3]
	if a.has("--autoenter") or a.has("--autoenter-register"):
		Net.on("authError", func(_m): Net.send("register", {"user": main.screens.in_user.text, "pass": main.screens.in_pass.text}))
	main.screens._submit()

func _auto_enter(m: Dictionary) -> void:
	if S.in_world:
		return
	if m.chars.is_empty():
		var nm := "Tester"
		for ch in str(randi() % 100000):
			nm += "abcdefghij"[int(ch)]
		Net.send("createChar", {"name": nm, "race": "human", "class": "warrior", "look": {"gender": "m", "hair": 0, "skin": 0}})
	else:
		Net.send("enterWorld", {"id": m.chars[0].id})

func _wait_spawn() -> void:
	while not (S.in_world and main.spawned):
		await get_tree().process_frame

func _shot(a: PackedStringArray) -> void:
	var i := a.find("--shot")
	var path: String = a[i + 1]
	var delay := float(a[i + 2])
	var opens: PackedStringArray = []
	if a.has("--open"):
		opens = a[a.find("--open") + 1].split(",")
	if a.has("--hour"):
		DayNight.forced = float(a[a.find("--hour") + 1])
	var chats: PackedStringArray = []
	if a.has("--chat"):
		chats = a[a.find("--chat") + 1].split(";")
	var sent := false
	var interact_s := ""
	if a.has("--interact"):
		interact_s = a[a.find("--interact") + 1]
	var t := 0.0
	while t < delay:
		await get_tree().process_frame
		t += get_process_delta_time()
		if S.in_world and main.spawned and not sent and chats.size() > 0 and t > 3.0:
			sent = true
			for c in chats:
				Net.send("chat", {"text": c})
		if S.in_world and main.spawned and interact_s != "" and t > 6.0:
			var best: Ent = null
			var bd := 1e9
			for e: Ent in S.ents.values():
				if (e.m.contains(interact_s) or e.n.contains(interact_s)) and e.id != S.id:
					var d := Vector2(e.rx - S.local.x, e.rz - S.local.z).length()
					if d < bd:
						bd = d
						best = e
			if best:
				print("interact with ", best.n, " (", best.m, ") at ", bd, " m")
				main.controller.set_target(best.id)
				Net.send("interact", {"id": best.id})
			interact_s = ""
		if S.in_world and main.spawned and opens.size() > 0 and t > delay - 2.0:
			for o in opens:
				main.windows.toggle(o)
			opens = []
	get_viewport().get_texture().get_image().save_png(path)
	print("saved ", path)
	get_tree().quit()

## Scripted play: walk, target a mob and fight it; prints PASS/FAIL lines.
func _play(a: PackedStringArray) -> void:
	await _wait_spawn()
	var res := {"ok": true}
	var check := func(name: String, cond: bool) -> void:
		print(("PASS " if cond else "FAIL ") + name)
		if not cond:
			res.ok = false
	Net.on("err", func(m): print("server error: ", m.text))
	var srv := {"tg": -1}
	Net.on("s", func(m): srv.tg = int(m.me.get("tg", 0)))
	await get_tree().create_timer(1.0).timeout
	var x0: float = S.local.x
	var z0: float = S.local.z
	Input.action_press("move_forward")
	await get_tree().create_timer(2.0).timeout
	Input.action_release("move_forward")
	var moved := Vector2(S.local.x - x0, S.local.z - z0).length()
	check.call("walk moved %.1f m" % moved, moved > 5.0)
	await get_tree().create_timer(0.6).timeout
	var own: Ent = S.ents.get(S.id)
	check.call("server position follows (%.1f vs %.1f)" % [own.x, S.local.x], absf(own.x - S.local.x) < 3.0 and absf(own.z - S.local.z) < 3.0)
	# find a mob and teleport near it via the server's admin command is not available: walk towards it
	var mob: Ent = null
	var best := 1e9
	for e: Ent in S.ents.values():
		if e.k == "mob" and not e.d:
			var d := Vector2(e.rx - S.local.x, e.rz - S.local.z).length()
			if d < best:
				best = d
				mob = e
	check.call("a mob is in range (%.0f m)" % best, mob != null)
	if mob:
		main.controller.set_target(mob.id)
		await get_tree().create_timer(0.3).timeout
		check.call("target set (server says %d)" % srv.tg, S.me.tg == mob.id and srv.tg == mob.id)
		var hp0 := mob.hp
		var seen := {"dmg": 0}
		Net.on("dmg", func(_m): seen.dmg += 1)
		# run towards it
		var t := 0.0
		while t < 12.0 and Vector2(mob.rx - S.local.x, mob.rz - S.local.z).length() > 2.2:
			S.local.r = atan2(mob.rx - S.local.x, mob.rz - S.local.z)
			main.controller.cam_yaw = S.local.r
			Input.action_press("move_forward")
			await get_tree().process_frame
			t += get_process_delta_time()
		Input.action_release("move_forward")
		print("skills: ", S.skill_ids, " dist=", Vector2(mob.rx - S.local.x, mob.rz - S.local.z).length(), " target=", S.me.tg, " mob=", mob.m, " lvl ", mob.l)
		for i in 12:
			main.hud.cast_slot(0)
			print("cast ", i, " dist=", Vector2(mob.rx - S.local.x, mob.rz - S.local.z).length(), " hp=", mob.hp, " frame=", Engine.get_frames_drawn())
			await get_tree().create_timer(0.5).timeout
		check.call("damage events received (%d)" % seen.dmg, seen.dmg > 0)
		check.call("mob lost hp (%d -> %d)" % [hp0, mob.hp], mob.hp < hp0 or mob.d)
	print("music dir: ", Sound._music_dir, " track: ", Sound._track)
	check.call("music plays a fight track after combat", Sound._music_dir == "" or Sound._track.ends_with("_fight"))
	print("PLAY RESULT: " + ("OK" if res.ok else "FAILED"))
	get_tree().quit()


# ---------------------------------------------------------------- scripted play 2: dialogs, vendor, gear, quests, trainer, flight

func _until(cond: Callable, timeout := 6.0) -> bool:
	var t := 0.0
	while t < timeout:
		if cond.call():
			return true
		await get_tree().process_frame
		t += get_process_delta_time()
	return cond.call()

func _tp(x: float, z: float) -> void:
	Net.send("chat", {"text": "/dev tp %d %d" % [int(x), int(z)]})
	await _until(func(): return Vector2(S.local.x - x, S.local.z - z).length() < 3.0, 4.0)
	await get_tree().create_timer(0.8).timeout

func _interact(who: String) -> bool:
	for e: Ent in S.ents.values():
		if e.id != S.id and (e.m == who or e.n.contains(who)):
			main.controller.set_target(e.id)
			Net.send("interact", {"id": e.id})
			return true
	return false

func _play2(_a: PackedStringArray) -> void:
	await _wait_spawn()
	var res := {"ok": true}
	var check := func(name: String, cond: bool) -> void:
		print(("PASS " if cond else "FAIL ") + name)
		if not cond:
			res.ok = false
	var errs: Array = []
	Net.on("err", func(m): errs.append(m.text); print("server error: ", m.text))
	var W: Windows = main.windows
	Net.on("inv", func(m): print("inv message: ", (m.get("items", []) as Array).size(), " items, copper ", m.get("copper"), " max ", m.get("max")))
	Net.send("chat", {"text": "/dev level 12"})
	Net.send("chat", {"text": "/dev gold 200"})
	await get_tree().create_timer(1.0).timeout
	# --- vendor
	await _tp(11, 13)
	await get_tree().create_timer(1.0).timeout
	check.call("merchant found", _interact("merchant"))
	check.call("dialog opens", await _until(func(): return W.dialog.visible))
	check.call("dialog lists quests", (W.dialog.m.get("quests", []) as Array).size() > 0)
	W.vendor.open_for(W.dialog.m)
	check.call("vendor shows goods", W.vendor.list.get_child_count() > 0)
	var stock: Array = W.dialog.m.get("stock", [])
	var bought := ""
	for id in stock:
		var it := S.item(id)
		if int(it.get("buy", 99999999)) <= S.copper and not (it.get("req", 0) > S.me.lvl):
			bought = id
			break
	check.call("something affordable", bought != "")
	if bought != "":
		Net.send("buy", {"npc": W.dialog.m.npc, "item": bought, "count": 1})
		check.call("inventory updated", await _until(func(): return S.inv.size() > 0))
		W.bags.open()
		W.bags.refresh()
		check.call("bag slot shows item", W.bags.slot_nodes.size() > 0 and (W.bags.slot_nodes[0] as Slot).icon != "")
		var idx := 0
		var it2 := S.item(S.inv[0].i)
		if it2.get("slot", "") != "":
			Net.send("equip", {"idx": 0})
			check.call("item equipped", await _until(func(): return S.equip.size() > 0))
		W.char_win.open()
		W.char_win.refresh()
		check.call("character window builds", W.char_win.content.get_child_count() > 3)
	# --- quest
	var qid := ""
	for o in W.dialog.m.get("quests", []):
		if o.state == "available":
			qid = o.id
			break
	check.call("quest offered", qid != "")
	if qid != "":
		W.dialog.quest_id = qid
		W.dialog.view = "quest"
		W.dialog.render()
		check.call("quest text rendered", W.dialog.content.get_child_count() > 3)
		Net.send("acceptQuest", {"quest": qid})
		check.call("quest accepted", await _until(func(): return S.quests.size() == 1))
		await get_tree().create_timer(0.3).timeout
		check.call("tracker shows quest", main.hud.tracker_body.get_child_count() >= 2)
		W.quests.open()
		W.quests.refresh()
		check.call("quest log builds", W.quests.detail.get_child_count() > 3)
	W.close_top()
	W.close_top()
	# --- trainer
	await _tp(-9, 13)
	await get_tree().create_timer(1.0).timeout
	check.call("trainer found", _interact("trainer_warrior"))
	check.call("trainer dialog opens", await _until(func(): return W.dialog.visible and W.dialog.m.has("trainer")))
	W.dialog.view = "trainer"
	W.dialog.render()
	check.call("trainer lists skills", W.dialog.content.get_child_count() > 2)
	var learn := ""
	for sk in W.dialog.m.trainer.get("skills", []):
		if sk.state == "learnable":
			learn = sk.id
			break
	if learn != "":
		Net.send("train", {"npc": W.dialog.m.npc, "skill": learn})
		check.call("skill learned", await _until(func(): return S.skill_ids.has(learn)))
		W.spells.open()
		W.spells.refresh()
		check.call("spell book builds", W.spells.list.get_child_count() > 3)
	W.dialog.close()
	# --- map and party windows
	W.map.open()
	await get_tree().create_timer(0.5).timeout
	W.map.close()
	W.party.open()
	W.party.refresh()
	W.party.close()
	# --- griffin flight
	Net.send("chat", {"text": "/dev unlock"})
	await get_tree().create_timer(0.5).timeout
	await _tp(20, 22)
	await get_tree().create_timer(1.0).timeout
	check.call("griffin master found", _interact("flightmaster"))
	check.call("flight dialog opens", await _until(func(): return W.dialog.visible and W.dialog.view == "flight"))
	var pts: Array = W.dialog.m.flight.get("points", [])
	print("flight points: ", pts.size())
	if pts.size() > 0:
		var started := {"v": false}
		Net.on("flightStart", func(_m): started.v = true)
		Net.send("flight", {"npc": W.dialog.m.npc, "dest": pts[0].id})
		check.call("flight starts", await _until(func(): return started.v))
		check.call("player is flying", await _until(func(): return S.me.get("fl", false)))
		var p0 := Vector2(S.local.x, S.local.z)
		await get_tree().create_timer(3.0).timeout
		check.call("flight moves the player", Vector2(S.local.x, S.local.z).distance_to(p0) > 5.0)
		var own: Ent = S.ents.get(S.id)
		check.call("griffin mount attached", own != null and own.mount != null)
	print("server errors: ", errs.size())
	print("PLAY2 RESULT: " + ("OK" if res.ok else "FAILED"))
	get_tree().quit()
