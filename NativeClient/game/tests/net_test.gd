extends Node
# Integration test: log in against a running server and dump what arrives. Usage:
#   godot --headless --path NativeClient/game res://tests/net_test.tscn -- ws://localhost:8081/ws user pass
var step := 0
func _ready() -> void:
	var args := OS.get_cmdline_user_args()
	var url: String = args[0] if args.size() > 0 else "ws://localhost:8080/ws"
	var user: String = args[1] if args.size() > 1 else "nettest"
	var pw: String = args[2] if args.size() > 2 else "secret123"
	Net.on("auth", func(m): print("auth ok: ", m.user, " admin=", m.get("admin"), " classes=", m.classOrder, " skills=", m.skills.size()); Net.send("createChar", {"name": "Nt" + str(randi() % 100000).replace("0","a").replace("1","b").replace("2","c").replace("3","d").replace("4","e").replace("5","f").replace("6","g").replace("7","h").replace("8","i").replace("9","j"), "race": "human", "class": "warrior", "look": {"gender": "m", "hair": 0, "skin": 0}}))
	Net.on("authError", func(m): print("authError: ", m.error); Net.send("register", {"user": user, "pass": pw}))
	Net.on("chars", func(m): print("chars: ", m.chars.size()); if m.chars.size() > 0: Net.send("enterWorld", {"id": m.chars[0].id}))
	Net.on("createError", func(m): print("createError: ", m.error))
	Net.on("world", func(m):
		var sc: Dictionary = m.scenery
		print("world: id=", m.id, " items=", m.defs.items.size(), " quests=", m.defs.quests.size(), " scenery keys=", sc.keys().size(), " props b64 len=", str(sc.props).length(), " now=", m.now)
		var hm := Marshalls.base64_to_raw(sc.hm)
		print("heightmap bytes=", hm.size(), " nx=", sc.nx, " nz=", sc.nz, " h[0]=", hm.decode_s16(0) / 20.0))
	var cnt := {"n": 0}
	Net.on("s", func(m):
		cnt.n += 1
		if cnt.n == 3:
			print("snapshot entities=", m.e.size(), " first=", m.e[0])
			get_tree().quit())
	var seen := {}
	Net.set_meta("seen", seen)
	Net.connected.connect(func(): print("connected"); Net.send("login", {"user": user, "pass": pw}))
	Net.closed.connect(func(): print("closed"); get_tree().quit())
	print("connecting ", url, " err=", Net.connect_to(url))
	get_tree().create_timer(25.0).timeout.connect(func(): print("TIMEOUT n=", cnt.n, " seen=", Net.get_meta("seen")); get_tree().quit())
