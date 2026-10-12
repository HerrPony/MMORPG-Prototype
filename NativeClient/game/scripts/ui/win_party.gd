class_name WinParty
extends GameWindow
## Party window: member list, kick and leave.

var content: VBoxContainer

func _init() -> void:
	super._init("Gruppe", Vector2(300, 100))
	win_id = "party"
	default_pos = Vector2(120, 330)
	content = UI.vbox(5)
	body.add_child(content)

func refresh() -> void:
	if not visible:
		return
	UI.clear(content)
	var members: Array = S.party.members
	if members.is_empty():
		var l := UI.label("Du bist in keiner Gruppe. Rechtsklick auf einen Spieler, um ihn einzuladen (oder /einladen Name). Gruppen haben bis zu 5 Spieler.", 13, UI.MUTED)
		l.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
		l.custom_minimum_size.x = 280
		content.add_child(l)
		reset_size()
		return
	for m in members:
		var h := UI.hbox(8)
		var ic := UI.label(UI.CLASS_ICONS.get(m.c, "?"), 20, Color.WHITE, HORIZONTAL_ALIGNMENT_CENTER)
		ic.custom_minimum_size.x = 30
		h.add_child(ic)
		var nm := UI.label(("⭐ " if int(m.id) == int(S.party.leader) else "") + "%s  (Stufe %d)" % [m.n, int(m.l)], 14, UI.CLASS_COLORS.get(m.c, Color.WHITE))
		nm.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		h.add_child(nm)
		if int(S.party.leader) == S.id and int(m.id) != S.id:
			var id := int(m.id)
			h.add_child(UI.button("✕", func(): Net.send("kick", {"id": id})))
		content.add_child(h)
	content.add_child(UI.button("Gruppe verlassen", func(): Net.send("leaveParty")))
	reset_size()
