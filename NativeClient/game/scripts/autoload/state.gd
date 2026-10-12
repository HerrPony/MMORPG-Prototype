extends Node
## Shared client state: the single source of truth for UI and world code (mirror of Client/js/state.js).

signal changed(what: String)

var user := ""
var classes: Dictionary = {}
var class_order: Array = []
var races: Dictionary = {}
var race_order: Array = []
var all_skills: Dictionary = {}

# in the world
var defs: Dictionary = {}          # items, quests, skills, pois
var scenery: Dictionary = {}
var id := 0
var self_info: Dictionary = {}     # name, race, class, look
var cls: Dictionary = {}           # own class definition
var skill_ids: Array = []
var me: Dictionary = {}            # hp, mh, res, rm, xp, xn, lvl, tg, auto, d ...
var ents: Dictionary = {}          # id -> Entity (see world/entity.gd)
var inv: Array = []
var copper := 0
var inv_max := 28
var equip: Dictionary = {}
var stats: Dictionary = {}
var sets: Dictionary = {}
var vendor_npc := 0
var vendor_data: Dictionary = {}
var disc: Array = []
var fp: Array = []
var party: Dictionary = {}
var zone: Dictionary = {}
var instance = null
var quests: Array = []
var marks: Dictionary = {}
var cd: Dictionary = {}
var gcd := {"end": 0.0, "dur": 0.0}
var cast = null
var local := {"x": 0.0, "y": 0.0, "z": 0.0, "r": 0.0, "speed": 0.0, "airborne": false}
var typing := false
var in_world := false
var admin := false
var dialog_npc := 0
var map_cont := -1
var server_offset_ms := 0.0         # server clock minus local clock

func _ready() -> void:
	reset_world()

func reset_world() -> void:
	me = {"hp": 1, "mh": 1, "res": 0, "rm": 1, "xp": 0, "xn": 100, "lvl": 1, "tg": 0, "auto": false, "d": false}
	party = {"id": 0, "leader": 0, "members": [], "states": []}
	zone = {"name": "", "sub": "", "lo": 1, "hi": 6}
	cd = {}
	gcd = {"end": 0.0, "dur": 0.0}
	cast = null
	quests = []
	marks = {}
	instance = null
	inv = []
	equip = {}
	stats = {}
	sets = {}
	disc = []
	fp = []
	vendor_npc = 0
	vendor_data = {}

func now() -> float:
	return Time.get_ticks_msec() / 1000.0

func item(item_id: String) -> Dictionary:
	var items: Dictionary = defs.get("items", {})
	return items.get(item_id, {})

func quest_def(qid: String) -> Dictionary:
	var q: Dictionary = defs.get("quests", {})
	return q.get(qid, {})

func skill(sid: String) -> Dictionary:
	var sk: Dictionary = defs.get("skills", {})
	return sk.get(sid, {})
