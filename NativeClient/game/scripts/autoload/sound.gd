extends Node
## Music and sound effects. Music tracks are pre-rendered OGG files (see tools/music) that are looked up next to the
## executable (installed by the launcher), in user://music or in the build folder during development; sound effects
## are part of the game resources. Missing files are simply skipped.

const SFX_DIR := "res://assets/sfx/"
const MUSIC_DIRS := ["music", "../build/audio/music"]
const FADE := 2.0

var zone := "forest"
var mode := "menu"
var combat := false
var _music: Array[AudioStreamPlayer] = []
var _cur := 0
var _track := ""
var _track_time := 0.0
var _streams: Dictionary = {}
var _sfx_cache: Dictionary = {}
var _sfx_pool: Array[AudioStreamPlayer] = []
var _music_dir := ""
var _last_pick := ""
var _step_t := 0.0
var _bus_ready := false

func _ready() -> void:
	process_mode = Node.PROCESS_MODE_ALWAYS
	_setup_buses()
	for i in 2:
		var p := AudioStreamPlayer.new()
		p.bus = "Music"
		p.volume_db = -80.0
		add_child(p)
		_music.append(p)
	for i in 12:
		var p2 := AudioStreamPlayer.new()
		p2.bus = "SFX"
		add_child(p2)
		_sfx_pool.append(p2)
	_music_dir = _find_music_dir()
	Settings.changed.connect(apply_volumes)
	apply_volumes()

func _setup_buses() -> void:
	for n in ["Music", "SFX", "Ambient"]:
		if AudioServer.get_bus_index(n) < 0:
			AudioServer.add_bus()
			var idx := AudioServer.bus_count - 1
			AudioServer.set_bus_name(idx, n)
			AudioServer.set_bus_send(idx, "Master")
	_bus_ready = true

func apply_volumes() -> void:
	if not _bus_ready:
		return
	var set_bus := func(name: String, v: float) -> void:
		var i := AudioServer.get_bus_index(name)
		AudioServer.set_bus_volume_db(i, linear_to_db(maxf(0.0001, v)))
		AudioServer.set_bus_mute(i, v <= 0.001)
	set_bus.call("Master", float(Settings.v["vol_master"]))
	set_bus.call("Music", float(Settings.v["vol_music"]) if Settings.v["music"] else 0.0)
	set_bus.call("SFX", float(Settings.v["vol_sfx"]))
	set_bus.call("Ambient", float(Settings.v["vol_ambient"]))

func _find_music_dir() -> String:
	var cands: Array = []
	var env := OS.get_environment("ERLENWALD_MUSIC")
	if env != "":
		cands.append(env)
	var exe_dir := OS.get_executable_path().get_base_dir()
	for d in MUSIC_DIRS:
		cands.append(exe_dir.path_join(d))
		cands.append(ProjectSettings.globalize_path("res://").path_join(d))
	cands.append(ProjectSettings.globalize_path("user://music"))
	for c in cands:
		if DirAccess.dir_exists_absolute(c):
			var d := DirAccess.open(c)
			for f in d.get_files():
				if f.ends_with(".ogg"):
					return c
	return ""

func set_mode(m: String) -> void:
	mode = m

func _wanted() -> String:
	if mode == "menu":
		return "menu"
	if mode == "dungeon":
		return "dungeon_fight" if combat else "dungeon"
	var z := zone if zone != "" else "forest"
	if combat:
		return z + "_fight"
	return z + ("_night" if float(DayNight.info().night) > 0.55 else "_day")

func _load_stream(name: String) -> AudioStream:
	if _streams.has(name):
		return _streams[name]
	var s: AudioStream = null
	if _music_dir != "":
		var path := _music_dir.path_join(name + ".ogg")
		if FileAccess.file_exists(path):
			s = AudioStreamOggVorbis.load_from_file(path)
	_streams[name] = s
	return s

func _process(dt: float) -> void:
	if _music_dir == "":
		return
	_track_time += dt
	var want := _wanted()
	var cur := _music[_cur]
	var fight_switch := want.ends_with("_fight") != _track.ends_with("_fight") and _track != ""
	if want != _track:
		# switch at once for combat, otherwise only after the current track has played for a while
		if _track == "" or fight_switch or _track_time > 20.0 or mode == "menu" or _track.begins_with("menu") or (mode == "dungeon") != _track.begins_with("dungeon"):
			_play(want)
	elif not cur.playing and _track != "":
		_play(want)
	elif cur.playing and cur.stream and cur.get_playback_position() > cur.stream.get_length() - FADE - 1.0 and _track_time > 10.0:
		_play(want)
	# fade the two players
	for i in 2:
		var p := _music[i]
		var target := 0.0 if i == _cur else -80.0
		if p.playing or p.volume_db > -79.0:
			p.volume_db = move_toward(p.volume_db, target, dt * (80.0 / FADE) * 1.2 if target < p.volume_db else dt * 40.0)
			if i != _cur and p.volume_db <= -79.0:
				p.stop()

func _play(name: String) -> void:
	var s := _load_stream(name)
	if s == null:
		if not name.ends_with("_fight") and _track == "":
			_track = name
		return
	_cur = 1 - _cur
	var p := _music[_cur]
	p.stream = s
	p.volume_db = -60.0
	p.play()
	_track = name
	_track_time = 0.0

# ---------------------------------------------------------------- effects

func sfx(name: String, volume_db := 0.0, pitch := 1.0) -> void:
	var s: AudioStream = _sfx_cache.get(name)
	if s == null and not _sfx_cache.has(name):
		var path := SFX_DIR + name + ".ogg"
		s = load(path) if ResourceLoader.exists(path) else null
		_sfx_cache[name] = s
	if s == null:
		return
	for p in _sfx_pool:
		if not p.playing:
			p.stream = s
			p.volume_db = volume_db
			p.pitch_scale = pitch * randf_range(0.95, 1.05)
			p.play()
			return

## Called every frame with the player's speed so that footsteps match the ground.
func footsteps(dt: float, speed: float, grounded: bool, x: float, z: float) -> void:
	if speed < 0.15 or not grounded:
		_step_t = 0.0
		return
	_step_t -= dt
	if _step_t <= 0.0:
		_step_t = 0.62 - 0.28 * minf(speed, 1.0)
		var b := WD.biome_at(x, z)
		var name := "step"
		match b:
			4: name = "step_snow"
			5, 8: name = "step_sand"
			9, 10: name = "step_stone"
		if WD.gh(x, z) < -0.3:
			name = "splash"
		sfx(name, -6.0, 1.0)
