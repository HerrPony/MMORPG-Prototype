class_name Windows
extends Control
## Owns all game windows and routes server messages to them.

signal logout_requested
signal quit_requested
signal editor_requested

var bags: WinBags
var char_win: WinChar
var spells: WinSpells
var quests: WinQuests
var dialog: WinDialog
var vendor: WinVendor
var loot: WinLoot
var party: WinParty
var map: WinMap
var menu: WinMenu
var settings: WinSettings
var hud: Hud
var tooltip: Tooltip
var all: Array = []

func _ready() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	mouse_filter = Control.MOUSE_FILTER_IGNORE
	theme = UI.theme()
	bags = WinBags.new(); char_win = WinChar.new(); spells = WinSpells.new(); quests = WinQuests.new()
	dialog = WinDialog.new(); vendor = WinVendor.new(); loot = WinLoot.new(); party = WinParty.new()
	map = WinMap.new(); menu = WinMenu.new(); settings = WinSettings.new()
	all = [bags, char_win, spells, quests, dialog, vendor, loot, party, map, menu, settings]
	for w in all:
		add_child(w)
	tooltip = Tooltip.new()
	add_child(tooltip)
	dialog.open_vendor.connect(func(m): vendor.open_for(m); bags.open(); bags.refresh())
	menu.open_settings.connect(func(): settings.open())
	menu.logout.connect(func(): logout_requested.emit())
	menu.quit_game.connect(func(): quit_requested.emit())
	menu.open_editor.connect(func(): editor_requested.emit())
	bags.visibility_changed.connect(func(): if bags.visible: bags.refresh())
	char_win.visibility_changed.connect(func(): if char_win.visible: char_win.refresh())
	spells.visibility_changed.connect(func(): if spells.visible: spells.refresh())
	quests.visibility_changed.connect(func(): if quests.visible: quests.refresh())
	party.visibility_changed.connect(func(): if party.visible: party.refresh())
	Net.on("dialog", func(m): dialog.show_dialog(m))
	Net.on("loot", func(m): loot.show_loot(m))
	Net.on("lootClose", func(m):
		if loot.loot_id == int(m.id):
			loot.close()
		Tooltip.hide_tip())
	Net.on("inv", func(m):
		S.inv = m.get("items") if m.get("items") is Array else []
		S.copper = int(m.copper)
		S.inv_max = int(m.max)
		bags.refresh() if bags.visible else null
		char_win.refresh()
		if dialog.visible:
			dialog.render()
		vendor.refresh())
	Net.on("equip", func(m): S.equip = m.get("eq") if m.get("eq") is Dictionary else {}; char_win.refresh())
	Net.on("stats", func(m): S.stats = m; char_win.refresh())
	Net.on("sets", func(m): S.sets = m.get("sets") if m.get("sets") is Dictionary else {}; char_win.refresh())
	Net.on("skills", func(m):
		S.skill_ids = m.skills
		spells.refresh()
		var d := S.skill(m.learned)
		if not d.is_empty() and hud:
			hud.toast("Neue Fertigkeit: " + d.name))
	Net.on("quests", func(_m): quests.refresh())
	Net.on("questDone", func(_m): dialog.close())
	Net.on("party", func(_m): party.refresh())
	Net.on("disc", func(m): S.disc = m.ids if m.ids is Array else [])
	Net.on("fp", func(m): S.fp = m.ids if m.ids is Array else [])

func toggle(name: String) -> void:
	var w: GameWindow = {"bags": bags, "char": char_win, "spells": spells, "quests": quests, "map": map, "party": party, "menu": menu}.get(name)
	if w:
		w.toggle()

func close_top() -> bool:
	for w in [loot, vendor, dialog, map, party, spells, quests, bags, char_win, settings, menu]:
		if w.visible:
			w.close()
			return true
	return false

func reset() -> void:
	for w in all:
		w.visible = false
	S.vendor_npc = 0
	S.dialog_npc = 0
	loot.loot_id = 0

func update_distance() -> void:
	var chk := func(w: GameWindow, id: int) -> void:
		if not w.visible or id == 0:
			return
		var e: Ent = S.ents.get(id)
		if e == null or Vector2(e.rx - S.local.x, e.rz - S.local.z).length() > 12.0:
			w.close()
	chk.call(dialog, S.dialog_npc)
	chk.call(loot, loot.loot_id)
	chk.call(vendor, S.vendor_npc)

func any_open() -> bool:
	for w in all:
		if w.visible:
			return true
	return false
