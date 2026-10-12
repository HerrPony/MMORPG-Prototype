extends Node
# Loads (and so compiles) every script of the project with the autoloads in place; errors show up in the engine log.
func _ready() -> void:
	var bad := 0
	for dir in ["res://scripts", "res://tests"]:
		bad += _scan(dir)
	print("load_all: " + ("FAILED (%d)" % bad if bad > 0 else "ok"))
	get_tree().quit(1 if bad > 0 else 0)

func _scan(path: String) -> int:
	var bad := 0
	var d := DirAccess.open(path)
	if d == null:
		return 0
	for f in d.get_files():
		if f.ends_with(".gd") and f != "load_all.gd":
			var s := ResourceLoader.load(path + "/" + f, "", ResourceLoader.CACHE_MODE_IGNORE)
			if s == null or not (s as GDScript).can_instantiate():
				print("cannot load ", path, "/", f)
				bad += 1
	for sub in d.get_directories():
		bad += _scan(path + "/" + sub)
	return bad
