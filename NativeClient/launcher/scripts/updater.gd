class_name Updater
extends Node
## Downloads and verifies the files of a release described by a manifest (see tools/release/make_release.py).
## Only files whose SHA-256 differs from the installed ones are fetched, so an update usually moves just the game pack.

signal status(text: String)
signal progress(done: int, total: int)

const MANIFEST_FORMAT := 1

var base_url := ""
var root := ""            # install root (contains the launcher, game/ below it)
var platform := ""
var manifest: Dictionary = {}
var installed: Dictionary = {}   # installed.json: {version, files:{path:{sha256,size,mtime}}}
var last_error := ""
var _cancel := false
var force_launcher_update := false   # tests: also update when running from the editor

func _init(url := "", install_root := "") -> void:
	base_url = url.trim_suffix("/")
	root = install_root
	platform = detect_platform()

static func detect_platform() -> String:
	match OS.get_name():
		"Windows": return "windows"
		"macOS": return "macos"
		_: return "linux"

func game_dir() -> String:
	return root.path_join("game")

func cancel() -> void:
	_cancel = true

# ------------------------------------------------------------------ version helpers

static func parse_version(v: String) -> Array:
	var out: Array = []
	for p in v.split("."):
		out.append(int(p))
	while out.size() < 3:
		out.append(0)
	return out

## Returns 1 if a > b, -1 if a < b, 0 if equal.
static func compare_versions(a: String, b: String) -> int:
	var x := parse_version(a)
	var y := parse_version(b)
	for i in 3:
		if x[i] > y[i]:
			return 1
		if x[i] < y[i]:
			return -1
	return 0

static func safe_path(p: String) -> bool:
	return p != "" and not p.begins_with("/") and not p.contains("..") and not p.contains(":") and not p.contains("\\")

# ------------------------------------------------------------------ http

func _url_for(rel: String) -> String:
	if rel.begins_with("http://") or rel.begins_with("https://"):
		return rel
	return base_url + "/" + rel

func http_get(url: String) -> Dictionary:
	var req := HTTPRequest.new()
	req.use_threads = true
	req.timeout = 25.0
	req.accept_gzip = true
	add_child(req)
	var st := {"done": false, "res": 0, "code": 0, "body": PackedByteArray()}
	req.request_completed.connect(func(result: int, code: int, _h: PackedStringArray, body: PackedByteArray):
		st.res = result; st.code = code; st.body = body; st.done = true)
	var err := req.request(url, ["Cache-Control: no-cache", "User-Agent: ErlenwaldLauncher/" + str(ProjectSettings.get_setting("application/config/version", "1.0.0"))])
	if err != OK:
		req.queue_free()
		return {"ok": false, "code": 0, "error": "Anfrage konnte nicht gestartet werden (%d)" % err}
	while not st.done:
		await get_tree().process_frame
	req.queue_free()
	if st.res != HTTPRequest.RESULT_SUCCESS:
		return {"ok": false, "code": st.code, "error": "Netzwerkfehler (%d)" % st.res}
	if st.code != 200:
		return {"ok": false, "code": st.code, "error": "Server antwortete mit HTTP %d" % st.code}
	return {"ok": true, "code": 200, "body": st.body}

## Downloads a file; `on_bytes` is called with the number of bytes received so far. Detects stalled transfers.
func http_download(url: String, dest: String, on_bytes: Callable) -> Dictionary:
	DirAccess.make_dir_recursive_absolute(dest.get_base_dir())
	var req := HTTPRequest.new()
	req.use_threads = true
	req.download_file = dest
	req.accept_gzip = false
	req.download_chunk_size = 1 << 18
	add_child(req)
	var st := {"done": false, "res": 0, "code": 0}
	req.request_completed.connect(func(result: int, code: int, _h: PackedStringArray, _b: PackedByteArray):
		st.res = result; st.code = code; st.done = true)
	var err := req.request(url, ["User-Agent: ErlenwaldLauncher"])
	if err != OK:
		req.queue_free()
		return {"ok": false, "error": "Download konnte nicht gestartet werden (%d)" % err}
	var last_bytes := 0
	var last_change := Time.get_ticks_msec()
	while not st.done:
		await get_tree().process_frame
		var b := req.get_downloaded_bytes()
		if b != last_bytes:
			last_bytes = b
			last_change = Time.get_ticks_msec()
			on_bytes.call(b)
		elif Time.get_ticks_msec() - last_change > 30000:
			req.cancel_request()
			req.queue_free()
			return {"ok": false, "error": "Zeitüberschreitung beim Download"}
		if _cancel:
			req.cancel_request()
			req.queue_free()
			return {"ok": false, "error": "Abgebrochen"}
	req.queue_free()
	if st.res != HTTPRequest.RESULT_SUCCESS:
		return {"ok": false, "error": "Netzwerkfehler (%d)" % st.res}
	if st.code != 200:
		return {"ok": false, "error": "HTTP %d" % st.code}
	return {"ok": true}

# ------------------------------------------------------------------ manifest

## Fetches and validates the manifest of this platform. Returns true on success (see `manifest` / `last_error`).
func fetch_manifest() -> bool:
	var url := _url_for("manifest-%s.json" % platform)
	status.emit("Suche nach Updates …")
	var r: Dictionary = await http_get(url)
	if not r.ok:
		last_error = r.error
		return false
	var parsed = JSON.parse_string((r.body as PackedByteArray).get_string_from_utf8())
	if not (parsed is Dictionary):
		last_error = "Ungültiges Manifest."
		return false
	var m: Dictionary = parsed
	if int(m.get("format", 0)) != MANIFEST_FORMAT or not m.has("files") or not m.has("version"):
		last_error = "Manifest hat ein unbekanntes Format."
		return false
	for f in m.files:
		if not safe_path(str(f.get("path", ""))) or str(f.get("sha256", "")).length() != 64:
			last_error = "Manifest enthält ungültige Einträge."
			return false
	for f in m.get("launcher_files", []):
		if not safe_path(str(f.get("path", ""))) or str(f.get("sha256", "")).length() != 64:
			last_error = "Manifest enthält ungültige Einträge."
			return false
	manifest = m
	return true

func object_url(f: Dictionary) -> String:
	if f.has("url"):
		return _url_for(str(f.url))
	var h: String = f.sha256
	var tmpl: String = manifest.get("objects", "objects/{h2}/{h}")
	return _url_for(tmpl.replace("{h2}", h.substr(0, 2)).replace("{h}", h))

# ------------------------------------------------------------------ installed state

func load_installed() -> void:
	installed = {"version": "", "files": {}}
	var p := game_dir().path_join("installed.json")
	if FileAccess.file_exists(p):
		var d = JSON.parse_string(FileAccess.get_file_as_string(p))
		if d is Dictionary and d.has("files"):
			installed = d

func save_installed() -> void:
	DirAccess.make_dir_recursive_absolute(game_dir())
	var f := FileAccess.open(game_dir().path_join("installed.json"), FileAccess.WRITE)
	if f:
		f.store_string(JSON.stringify(installed, "  "))

func installed_version() -> String:
	return str(installed.get("version", ""))

## True if the file on disk is the one the manifest wants. `deep` hashes the file, otherwise size and mtime are trusted.
func file_ok(dir: String, f: Dictionary, deep: bool) -> bool:
	var path := dir.path_join(f.path)
	if not FileAccess.file_exists(path):
		return false
	var sz := _file_size(path)
	if sz != int(f.get("size", -1)):
		return false
	if not deep:
		var rec = installed.get("files", {}).get(f.path)
		if rec is Dictionary and rec.get("sha256", "") == f.sha256 and int(rec.get("mtime", -1)) == int(FileAccess.get_modified_time(path)):
			return true
	return FileAccess.get_sha256(path) == f.sha256

func _file_size(path: String) -> int:
	var f := FileAccess.open(path, FileAccess.READ)
	if f == null:
		return -1
	var n := f.get_length()
	f.close()
	return n

func needs_update(deep := false) -> Array:
	var todo: Array = []
	for f in manifest.get("files", []):
		if not file_ok(game_dir(), f, deep):
			todo.append(f)
	return todo

# ------------------------------------------------------------------ installing

## Brings the game up to date. Returns true on success. `deep` re-hashes every file (repair).
func update_game(deep := false) -> bool:
	_cancel = false
	load_installed()
	var todo := needs_update(deep)
	var total := 0
	for f in todo:
		total += int(f.size)
	var done_before := 0
	var own_exe := OS.get_executable_path()
	var own_hash := ""
	if FileAccess.file_exists(own_exe):
		own_hash = ""  # computed lazily
	var idx := 0
	for f in todo:
		idx += 1
		if _cancel:
			last_error = "Abgebrochen"
			return false
		status.emit("Lade %s (%d/%d) …" % [str(f.path).get_file(), idx, todo.size()])
		var dest := game_dir().path_join(f.path)
		var part := root.path_join(".partial").path_join(str(f.sha256))
		DirAccess.make_dir_recursive_absolute(part.get_base_dir())
		var got := false
		# the engine binary of the launcher is the same file as the game's one: copy it instead of downloading it again
		if str(f.path).get_extension() in ["exe", "x86_64"] and int(f.size) == _file_size(own_exe):
			if own_hash == "":
				own_hash = FileAccess.get_sha256(own_exe)
			if own_hash == f.sha256:
				DirAccess.make_dir_recursive_absolute(dest.get_base_dir())
				if DirAccess.copy_absolute(own_exe, part) == OK:
					got = true
					progress.emit(done_before + int(f.size), total)
		var attempt := 0
		while not got and attempt < 3:
			attempt += 1
			var base := done_before
			var r: Dictionary = await http_download(object_url(f), part, func(b: int): progress.emit(base + b, total))
			if not r.ok:
				last_error = "%s: %s" % [f.path, r.error]
				if _cancel:
					return false
				await get_tree().create_timer(1.5 * attempt).timeout
				continue
			if FileAccess.get_sha256(part) != f.sha256:
				last_error = "%s: Prüfsumme stimmt nicht" % f.path
				DirAccess.remove_absolute(part)
				continue
			got = true
		if not got:
			return false
		if not _install_file(part, dest, f):
			return false
		done_before += int(f.size)
		progress.emit(done_before, total)
	# remove files that were part of an earlier version but are not wanted anymore
	var wanted := {}
	for f in manifest.files:
		wanted[f.path] = true
	for p in installed.get("files", {}).keys():
		if not wanted.has(p) and safe_path(p):
			DirAccess.remove_absolute(game_dir().path_join(p))
	var rec := {}
	for f in manifest.files:
		var path := game_dir().path_join(f.path)
		rec[f.path] = {"sha256": f.sha256, "size": int(f.size), "mtime": int(FileAccess.get_modified_time(path))}
	installed = {"version": manifest.version, "files": rec, "installed_at": Time.get_datetime_string_from_system()}
	save_installed()
	DirAccess.remove_absolute(root.path_join(".partial"))
	return true

## Records the current files as the installed version (used when the installer already bundled them).
func adopt_existing() -> void:
	var rec := {}
	for f in manifest.files:
		var path := game_dir().path_join(f.path)
		rec[f.path] = {"sha256": f.sha256, "size": int(f.size), "mtime": int(FileAccess.get_modified_time(path))}
	installed = {"version": manifest.version, "files": rec, "installed_at": Time.get_datetime_string_from_system()}
	save_installed()

func _install_file(part: String, dest: String, f: Dictionary) -> bool:
	DirAccess.make_dir_recursive_absolute(dest.get_base_dir())
	if FileAccess.file_exists(dest):
		# a running executable cannot be overwritten on Windows but it can be renamed away
		var old := dest + ".old"
		DirAccess.remove_absolute(old)
		if DirAccess.rename_absolute(dest, old) != OK and DirAccess.remove_absolute(dest) != OK:
			last_error = "%s kann nicht ersetzt werden (läuft das Spiel noch?)" % f.path
			return false
		DirAccess.remove_absolute(old)
	if DirAccess.rename_absolute(part, dest) != OK:
		if DirAccess.copy_absolute(part, dest) != OK:
			last_error = "%s konnte nicht installiert werden" % f.path
			return false
		DirAccess.remove_absolute(part)
	if f.get("exec", false) and OS.get_name() != "Windows":
		OS.execute("chmod", ["+x", dest])
	return true

# ------------------------------------------------------------------ launcher self-update

## True if one of the launcher files of the manifest differs from the installed one. Decided by content, not by
## version numbers, so a launcher that is already current never restarts itself in a loop. Big files (the engine
## binary) are compared by size only.
func launcher_update_available(_own_version := "") -> bool:
	if OS.has_feature("editor") and not force_launcher_update:
		return false
	for f in manifest.get("launcher_files", []):
		var path := root.path_join(f.path)
		if not FileAccess.file_exists(path):
			return true
		var sz := _file_size(path)
		if sz != int(f.size):
			return true
		if sz < 16 * 1048576 and FileAccess.get_sha256(path) != f.sha256:
			return true
	return false

## Downloads the new launcher files into <root>/.update and starts a helper that swaps them in after this process exits.
func stage_launcher_update() -> bool:
	var dir := root.path_join(".update")
	DirAccess.make_dir_recursive_absolute(dir)
	var files: Array = manifest.launcher_files
	var total := 0
	for f in files:
		total += int(f.size)
	var done := 0
	for f in files:
		status.emit("Lade Launcher-Update …")
		var dest := dir.path_join(f.path)
		# unchanged files (for instance the engine binary) are not downloaded again
		if FileAccess.file_exists(root.path_join(f.path)) and FileAccess.get_sha256(root.path_join(f.path)) == f.sha256:
			continue
		var base := done
		var r: Dictionary = await http_download(object_url(f), dest, func(b: int): progress.emit(base + b, total))
		if not r.ok or FileAccess.get_sha256(dest) != f.sha256:
			last_error = "Launcher-Update fehlgeschlagen: " + str(r.get("error", "Prüfsumme"))
			DirAccess.remove_absolute(dest)
			return false
		done += int(f.size)
	return true

## Writes and starts the swap helper. The caller must quit right after this returns true.
func apply_staged_launcher(restart: bool) -> bool:
	var dir := root.path_join(".update")
	var files: Array = manifest.launcher_files
	var exe := OS.get_executable_path()
	if OS.get_name() == "Windows":
		var cmd := "@echo off\r\nping -n 3 127.0.0.1 >nul\r\ncd /d \"%s\"\r\n" % root.replace("/", "\\")
		for f in files:
			var src := dir.path_join(f.path)
			if FileAccess.file_exists(src):
				cmd += "move /y \".update\\%s\" \"%s\" >nul\r\n" % [str(f.path).replace("/", "\\"), str(f.path).replace("/", "\\")]
		if restart:
			cmd += "start \"\" \"%s\"\r\n" % exe.replace("/", "\\")
		cmd += "rmdir /s /q .update\r\n(goto) 2>nul & del \"%~f0\"\r\n"
		var sp := root.path_join("update_launcher.cmd")
		var fa := FileAccess.open(sp, FileAccess.WRITE)
		if fa == null:
			last_error = "Update-Skript konnte nicht geschrieben werden"
			return false
		fa.store_string(cmd)
		fa.close()
		OS.create_process("cmd.exe", ["/c", sp.replace("/", "\\")])
	else:
		var sh := "#!/bin/sh\nsleep 1\ncd \"%s\" || exit 1\n" % root
		for f in files:
			var src := dir.path_join(f.path)
			if FileAccess.file_exists(src):
				sh += "mv -f \".update/%s\" \"%s\"\n" % [f.path, f.path]
				if f.get("exec", false):
					sh += "chmod +x \"%s\"\n" % f.path
		sh += "rm -rf .update\n"
		if restart:
			sh += "\"%s\" &\n" % exe
		sh += "rm -- \"$0\"\n"
		var sp2 := root.path_join("update_launcher.sh")
		var fb := FileAccess.open(sp2, FileAccess.WRITE)
		if fb == null:
			last_error = "Update-Skript konnte nicht geschrieben werden"
			return false
		fb.store_string(sh)
		fb.close()
		OS.execute("chmod", ["+x", sp2])
		OS.create_process("/bin/sh", [sp2])
	return true
