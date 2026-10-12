class_name WinMenu
extends GameWindow
## Game menu (Esc).

signal open_settings
signal logout
signal quit_game
signal open_editor

func _init() -> void:
	super._init("Menü", Vector2(280, 100))
	win_id = "menu"
	default_pos = Vector2(660, 200)
	var v := UI.vbox(8)
	body.add_child(v)
	v.add_child(UI.button("Zurück zum Spiel", close))
	v.add_child(UI.button("⚙ Einstellungen", func(): close(); open_settings.emit()))
	var help := UI.button("Steuerung anzeigen", Callable())
	var help_label := UI.rich("", 12)
	help_label.visible = false
	help.pressed.connect(func():
		help_label.visible = not help_label.visible
		help_label.text = _controls_text()
		reset_size())
	v.add_child(help)
	v.add_child(help_label)
	var ed := UI.button("🗺 Karten-Editor (Admin)", func(): open_editor.emit())
	ed.name = "EditorBtn"
	v.add_child(ed)
	v.add_child(UI.button("Abmelden / Charakterwahl", func(): close(); logout.emit()))
	v.add_child(UI.button("Spiel beenden", func(): quit_game.emit()))

func open(at := Vector2(-1, -1)) -> void:
	(body.get_child(0).get_node("EditorBtn") as Control).visible = S.admin
	super.open(at)

func _controls_text() -> String:
	var rows := [
		["W A S D", "Bewegen"], ["Leertaste", "Springen"], ["Rechtsklick halten", "Kamera drehen (Charakter blickt mit)"], ["Linksklick halten", "Kamera drehen"],
		["Linksklick", "NPC / Spieler / Gegner auswählen"], ["Rechtsklick", "NPC ansprechen · Gegner angreifen · Leiche looten"], ["Beide Maustasten", "Vorwärts laufen"],
		["Mausrad", "Zoom"], ["Tab", "Nächsten Gegner anvisieren"], ["1 – 0", "Fertigkeiten"], ["F", "Mit Ziel interagieren"], ["R", "Automatisch laufen"],
		["Shift", "Gehen (gedrückt halten)"], ["Enter  /  /", "Chat"], ["B L C P M O", "Taschen · Questlog · Charakter · Zauberbuch · Karte · Gruppe"],
		["Esc", "Fenster schließen / Ziel abwählen / Menü"], ["F11", "Vollbild"], ["F12", "Screenshot"],
	]
	var out := ""
	for r in rows:
		out += "[b]%s[/b]  %s\n" % r
	return out
