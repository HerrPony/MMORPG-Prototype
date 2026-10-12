class_name Models
extends RefCounted
## Builds the blocky low-poly models of players, NPCs, mobs and world objects.

const SKIN := [0xf1c9a5, 0xe0ac82, 0xc68a5d, 0x8d5a3b, 0x5b3a29]
const HAIR := [0x2b1b10, 0x6b4423, 0xc9a24a, 0xb0361d, 0x9aa0a6, 0x2c4ea8]

static func C(h: int, a := 1.0) -> Color:
	return Color(((h >> 16) & 255) / 255.0, ((h >> 8) & 255) / 255.0, (h & 255) / 255.0, a)

static func lighten(c: Color, k: float) -> Color:
	return Color(minf(1.0, c.r * k), minf(1.0, c.g * k), minf(1.0, c.b * k), c.a)

# ---------------------------------------------------------------- spec tables

const CLASS_LOOK := {
	"warrior": {"pants": 0x4a3b2c, "shirt": 0x8a8f9c, "hat": {"kind": "helmet", "color": 0x9aa0b0}, "mh": "sword", "oh": "shield", "body": Vector3(1.14, 1.0, 1.12)},
	"hunter": {"pants": 0x5a4a2d, "shirt": 0x4f7a3c, "hat": {"kind": "hood", "color": 0x2f4f2a}, "mh": "bow", "oh": "", "cloak": 0x2f4f2a, "body": Vector3(0.97, 1.02, 0.97), "quiver": true},
	"mage": {"pants": 0x23358a, "shirt": 0x3a4ab8, "hat": {"kind": "wizardhat", "color": 0x23358a}, "mh": "staff", "oh": "", "robe": true, "body": Vector3(0.92, 1.0, 0.92)},
}

const HUMANOIDS := {
	"forester": {"shirt": 0x5a8a4a, "hat": {"kind": "hood", "color": 0x3a5a2a}, "props": [["axe", "right"]]},
	"trainer_warrior": {"shirt": 0x9a9aa8, "props": [["sword", "right"], ["shield", "left"]], "body": Vector3(1.14, 1, 1.12)},
	"trainer_hunter": {"shirt": 0x5a8a45, "hat": {"kind": "hood", "color": 0x2f4f2a}, "props": [["bow", "left"]]},
	"trainer_mage": {"shirt": 0x7a3aa8, "robe": true, "hat": {"kind": "wizardhat", "color": 0x4a2a7a}, "props": [["staff", "right"]]},
	"merchant": {"shirt": 0xb8864a, "hat": {"kind": "straw"}},
	"smith": {"shirt": 0x8a6a4a, "props": [["pick", "right"]], "body": Vector3(1.2, 0.97, 1.2)},
	"farmer": {"shirt": 0xb8704a, "hat": {"kind": "straw"}, "props": [["pick", "right"]], "body": Vector3(1.05, 1, 1.1)},
	"scholar": {"shirt": 0x3a9a6a, "robe": true, "hat": {"kind": "wizardhat", "color": 0x1f5a3a}},
	"guard": {"shirt": 0x4a68a0, "hat": {"kind": "helmet", "color": 0x8a90a0}, "props": [["spear", "right"]], "body": Vector3(1.12, 1, 1.1)},
	"huntsman": {"shirt": 0x7a5a3a, "hat": {"kind": "hood", "color": 0x5a3a22}, "props": [["rifle", "right"]]},
	"mayor": {"shirt": 0xc8a050, "hat": {"kind": "crown"}, "robe": true, "body": Vector3(1.25, 0.95, 1.3)},
	"priest": {"shirt": 0xe8e8f4, "robe": true, "hat": {"kind": "hood", "color": 0xeeeeff}, "props": [["staff", "right"]]},
	"elder": {"shirt": 0x8a8a7a, "robe": true, "scale": 0.96, "hat": {"kind": "hood", "color": 0x6a6a5a}, "props": [["staff", "right"]]},
	"herbalist": {"shirt": 0x5a8a3a, "robe": true, "hat": {"kind": "hood", "color": 0x3a6a2a}},
	"fisher": {"shirt": 0x4a6a8a, "hat": {"kind": "straw"}},
	"sheikh": {"shirt": 0xe8e0c8, "robe": true, "hat": {"kind": "turban", "color": 0xe8e0c8}},
	"nomad": {"shirt": 0xc8a868, "robe": true, "hat": {"kind": "turban", "color": 0xb89858}},
	"flightmaster": {"shirt": 0xc8a040, "hat": {"kind": "hood", "color": 0x6a4a22}, "props": [["spear", "right"]]},
	"bandit": {"shirt": 0x5a3a3a, "hat": {"kind": "mask"}, "props": [["dagger", "right"]]},
	"bandit_chief": {"shirt": 0x7a2a2a, "scale": 1.05, "hat": {"kind": "mask", "color": 0xaa2222}, "props": [["axe", "right"]], "body": Vector3(1.2, 1, 1.15)},
	"desert_raider": {"shirt": 0xb89860, "hat": {"kind": "turban", "color": 0x8a2a2a}, "props": [["scimitar", "right"]]},
	"skeleton": {"shirt": 0xe0dccc, "skin": 0xe0dccc, "style": "skeleton", "props": [["sword", "right"]], "body": Vector3(0.8, 1, 0.8)},
	"zombie": {"shirt": 0x6a8a5a, "skin": 0x6a8a5a, "slow": true, "arms_forward": true},
	"zombie_giant": {"shirt": 0x5a7a4a, "skin": 0x5a7a4a, "slow": true, "arms_forward": true, "body": Vector3(1.6, 1, 1.45)},
	"kobold": {"shirt": 0x9a6a3a, "skin": 0x9a6a3a, "scale": 0.62, "hat": {"kind": "ears"}, "props": [["pick", "right"]], "head": 1.3, "body": Vector3(1.1, 0.95, 1.1)},
	"kobold_shaman": {"shirt": 0x7a3a9a, "skin": 0x9a6a3a, "scale": 0.62, "hat": {"kind": "ears", "color": 0x7a3a9a}, "props": [["witchstaff", "right"]], "head": 1.3},
	"kobold_boss": {"shirt": 0xb87a3a, "skin": 0x9a6a3a, "scale": 0.7, "hat": {"kind": "horns"}, "props": [["axe", "right"]], "head": 1.25, "body": Vector3(1.35, 0.95, 1.3)},
	"kobold_king": {"shirt": 0xd8a838, "skin": 0xb8803a, "scale": 0.7, "hat": {"kind": "crown"}, "props": [["sword", "right"]], "head": 1.2, "body": Vector3(1.4, 0.98, 1.35)},
	"yeti": {"shirt": 0xe8f2ff, "skin": 0xe8f2ff, "style": "fur", "scale": 1.25, "hat": {"kind": "ears", "color": 0xe8f2ff}, "head": 0.9, "body": Vector3(1.55, 1.0, 1.35)},
	"yeti_chief": {"shirt": 0xd8e8ff, "skin": 0xd8e8ff, "style": "fur", "scale": 1.25, "hat": {"kind": "horns"}, "props": [["axe", "right", 1.6]], "head": 0.9, "body": Vector3(1.7, 1.0, 1.45)},
	"witch": {"shirt": 0x6a3a8a, "robe": true, "hat": {"kind": "witchhat"}, "props": [["witchstaff", "right"]]},
	"hag": {"shirt": 0x4a5a2a, "skin": 0x7a9a5a, "robe": true, "hat": {"kind": "witchhat"}, "props": [["witchstaff", "right"]], "body": Vector3(0.9, 0.92, 0.9)},
	"mummy": {"shirt": 0xe8e0c0, "skin": 0xe8e0c0, "slow": true, "arms_forward": true},
	"tomb_guardian": {"shirt": 0xd8b038, "skin": 0xd8b038, "emissive": 0x302200, "hat": {"kind": "pharaoh"}, "props": [["spear", "right", 1.1]], "body": Vector3(1.35, 1.0, 1.3)},
	"water_elem": {"shirt": 0x2a6ad8, "skin": 0x2a6ad8, "emissive": 0x0a2a6a, "opacity": 0.72, "scale": 1.1, "body": Vector3(1.3, 1.0, 1.3)},
	"water_elem_pet": {"shirt": 0x4aaaff, "skin": 0x4aaaff, "emissive": 0x1a5a9a, "opacity": 0.72, "scale": 0.95, "body": Vector3(1.1, 1, 1.1)},
	"captain": {"shirt": 0x2a4a8a, "hat": {"kind": "captainhat"}, "props": [["dagger", "right"]]},
	"empress": {"shirt": 0xb8362a, "robe": true, "hat": {"kind": "crown"}, "body": Vector3(0.95, 1, 0.95)},
	"monk": {"shirt": 0xe0a040, "bald": true, "props": [["staff", "right"]]},
	"samurai": {"shirt": 0x8a2a2a, "hat": {"kind": "helmet", "color": 0x2a2a30}, "props": [["sword", "right"]], "body": Vector3(1.08, 1, 1.06)},
	"ape": {"shirt": 0x4a3426, "skin": 0x4a3426, "style": "fur", "body": Vector3(1.4, 0.95, 1.3), "head": 0.95, "arms_long": true},
	"poacher": {"shirt": 0x4a6a3a, "hat": {"kind": "hood", "color": 0x3a4a22}, "props": [["rifle", "right"]]},
	"lizardman": {"shirt": 0x7a5a2a, "skin": 0x4a8a4a, "style": "lizard", "props": [["spear", "right"]], "head": 1.05},
	"lizard_shaman": {"shirt": 0x3a6a9a, "skin": 0x4a8a8a, "style": "lizard", "hat": {"kind": "horns"}, "props": [["witchstaff", "right"]]},
	"ronin": {"shirt": 0x5a3a3a, "hat": {"kind": "straw"}, "props": [["scimitar", "right"]]},
	"ghost_warrior": {"shirt": 0x8aaed8, "skin": 0x8aaed8, "style": "ghost", "opacity": 0.7, "emissive": 0x102040, "hat": {"kind": "helmet", "color": 0x8aaed8}, "props": [["sword", "right"]]},
	"stone_golem": {"shirt": 0x8c8c86, "skin": 0x8c8c86, "style": "stone", "body": Vector3(1.5, 1.05, 1.4), "kind": "golem"},
	"canyon_bandit": {"shirt": 0xa84a2a, "hat": {"kind": "mask", "color": 0xc03a2a}, "props": [["axe", "right"]], "body": Vector3(1.05, 1, 1.05)},
	"cultist": {"shirt": 0x6a1a14, "robe": true, "hat": {"kind": "hood", "color": 0x4a1010}, "props": [["witchstaff", "right"]]},
	"fire_imp": {"shirt": 0x5a2410, "skin": 0x7a2a14, "style": "demon", "emissive": 0x200804, "hat": {"kind": "horns"}, "scale": 0.7, "head": 1.3, "body": Vector3(0.9, 0.9, 0.9)},
	"magma_golem": {"shirt": 0x3a2a28, "skin": 0x3a2a28, "style": "lava", "emissive": 0x3a1204, "body": Vector3(1.5, 1.05, 1.4), "kind": "golem"},
	"lava_elem": {"shirt": 0x3a2a28, "skin": 0x3a2a28, "style": "lava", "emissive": 0x601c06, "scale": 1.1, "body": Vector3(1.15, 1, 1.15), "kind": "golem"},
	"obsidian_knight": {"shirt": 0x23232e, "skin": 0x23232e, "style": "obsidian", "emissive": 0x0c0818, "hat": {"kind": "helmet", "color": 0x23232e}, "props": [["axe", "right", 1.3]], "body": Vector3(1.2, 1, 1.15)},
	"naga_warrior": {"shirt": 0x2f9a78, "skin": 0x2f9a78, "style": "lizard", "hat": {"kind": "helmet", "color": 0x2f7a68}, "props": [["spear", "right"]]},
	"naga_priest": {"shirt": 0x5a2ab8, "skin": 0x5a8a9a, "style": "lizard", "robe": true, "hat": {"kind": "turban", "color": 0x6a2ab8}, "props": [["witchstaff", "right"]]},
	"snake_emperor": {"shirt": 0xe0b040, "skin": 0x5aaa6a, "style": "lizard", "emissive": 0x302200, "robe": true, "hat": {"kind": "crown"}, "props": [["scimitar", "right", 1.3]], "body": Vector3(1.25, 1.0, 1.2)},
	"demon_lord": {"shirt": 0x2a1018, "skin": 0x5a1a1a, "style": "demon", "emissive": 0x2a0804, "hat": {"kind": "horns"}, "props": [["axe", "right", 1.8]], "head": 1.1, "body": Vector3(1.4, 1.05, 1.3)},
}
const FEMALE_NPC := ["priest", "herbalist", "witch", "hag", "scholar", "merchant", "trainer_mage", "fisher", "empress"]

# quadrupeds: body colour, secondary, scale, features
const QUADS := {
	"wolf": {"c": 0x777b80, "c2": 0xa8acb0, "s": 1.0, "ears": true, "tail": 0.5, "snout": true},
	"wolf_alpha": {"c": 0x5a3a30, "c2": 0x8a6a58, "s": 1.3, "ears": true, "tail": 0.55, "snout": true, "eyes": 0xff3a20},
	"icewolf": {"c": 0xc8dcea, "c2": 0xffffff, "s": 1.05, "ears": true, "tail": 0.5, "snout": true, "eyes": 0x60d0ff},
	"hyena": {"c": 0xb89868, "c2": 0xd8c498, "s": 0.95, "ears": true, "tail": 0.35, "snout": true, "spots": 0x3a2a1a, "slope": true},
	"boar": {"c": 0x6a4a34, "c2": 0x8a6a50, "s": 0.85, "tusks": true, "snout": true, "tail": 0.15, "short": true},
	"bear": {"c": 0x5a3c26, "c2": 0x7a5a40, "s": 1.35, "ears": true, "tail": 0.12, "snout": true, "stocky": true},
	"bear_big": {"c": 0x4a2e1c, "c2": 0x6a4a30, "s": 1.9, "ears": true, "tail": 0.12, "snout": true, "stocky": true},
	"polarbear": {"c": 0xeaf0f4, "c2": 0xffffff, "s": 1.4, "ears": true, "tail": 0.12, "snout": true, "stocky": true},
	"polarbear_big": {"c": 0xdfe8ef, "c2": 0xffffff, "s": 1.95, "ears": true, "tail": 0.12, "snout": true, "stocky": true},
	"bull": {"c": 0x4a3a30, "c2": 0xd8d0c0, "s": 1.2, "horns": true, "tail": 0.4, "stocky": true},
	"sheep": {"c": 0xf0f0e8, "c2": 0x3a3a3a, "s": 0.75, "wool": true, "tail": 0.15},
	"tiger": {"c": 0xd07a30, "c2": 0xf4e0c0, "s": 1.15, "ears": true, "tail": 0.7, "snout": true, "stripes": 0x1a1008},
	"whitetiger": {"c": 0xece8e0, "c2": 0xffffff, "s": 1.3, "ears": true, "tail": 0.7, "snout": true, "stripes": 0x606878, "eyes": 0x60b0ff},
	"panther": {"c": 0x1c1c22, "c2": 0x2c2c34, "s": 1.05, "ears": true, "tail": 0.75, "snout": true, "eyes": 0xe0e020},
	"kitsune": {"c": 0xd8782a, "c2": 0xffffff, "s": 0.9, "ears": true, "tail": 0.8, "snout": true, "tailtip": 0xffffff, "bigtail": true},
	"panda": {"c": 0xf4f4f0, "c2": 0x1a1a1a, "s": 1.25, "ears": true, "tail": 0.1, "snout": true, "stocky": true, "panda": true},
	"sandcat": {"c": 0xd0b078, "c2": 0xeed8a8, "s": 1.05, "ears": true, "tail": 0.7, "snout": true, "mane": 0x8a6a30},
	"bison": {"c": 0x4a3226, "c2": 0x2a1c14, "s": 1.5, "horns": true, "tail": 0.3, "stocky": true, "hump": true},
	"lizard": {"c": 0x5a8a3a, "c2": 0x9aba5a, "s": 1.0, "tail": 1.0, "snout": true, "low": true},
	"ashhound": {"c": 0x2c2220, "c2": 0x4a3a34, "s": 1.05, "ears": true, "tail": 0.5, "snout": true, "cracks": 0xff5a1a, "eyes": 0xff6a10},
}

# ---------------------------------------------------------------- entry point

## ent: Dictionary with k (kind), m (model key), c (class), lk (look), eq (equipment), id
static func make(ent: Dictionary) -> BlockyModel:
	var k: String = ent.get("k", "")
	var m: String = ent.get("m", "")
	var id := int(ent.get("id", 0))
	if k == "obj":
		return make_object(m)
	if k == "ship":
		return make_ship(m)
	if QUADS.has(m):
		return make_quad(QUADS[m], id)
	match m:
		"spider": return make_spider(0x2a2420, 0xd02020)
		"cavespider": return make_spider(0x3a3a52, 0x60e0ff)
		"spider_queen": return make_spider(0x5a1a2e, 0xff2060, 1.7)
		"scorpion": return make_scorpion(0xc8a040, 1.0)
		"scorpion_king": return make_scorpion(0x8a2a1a, 2.0)
		"croc": return make_croc()
		"wisp": return make_wisp(0x9aff4a)
		"frostspirit": return make_wisp(0x7ad8ff)
		"griffin": return make_griffin()
	if k == "player":
		return make_player(ent)
	var spec: Dictionary = (HUMANOIDS.get(m, {"shirt": 0x888888}) as Dictionary).duplicate()
	spec["female"] = m in FEMALE_NPC
	if not spec.has("skin"):
		spec["skin"] = SKIN[id % 4]
	spec["hair"] = HAIR[(id >> 1) % 6]
	return humanoid(spec)

# ---------------------------------------------------------------- players

static func tier_of(it: Dictionary) -> int:
	return clampi((int(it.get("req", 1)) - 1) / 2, 0, 19)

static func armor_color(it: Dictionary) -> Color:
	var t := tier_of(it)
	var typ: String = it.get("type", "")
	var pal: Array
	match typ:
		"Kette": pal = [0x80808a, 0x9aa8bc, 0xb8c4d4, 0xd4c470, 0x4a4a5c, 0xc8d4e4, 0xd8a860, 0x30303c]
		"Leder": pal = [0x7a5a3a, 0x8a6a3a, 0x4a6a3a, 0x3a5a4a, 0x6a4a5a, 0x2a4a5a, 0x5a3a2a, 0x2a2a30]
		"Stoff": pal = [0x3a4ab8, 0x5a3ab8, 0x2a8ab8, 0x8a2a8a, 0xb83a4a, 0xe0e0f0, 0x2a6a7a, 0x501a50]
		_: pal = [0x8a8a8a, 0xaaaaaa]
	var c := C(pal[(t / 3) % pal.size()])
	# a bit brighter with every tier, rare+ pieces glint
	return lighten(c, 0.9 + 0.012 * (t % 3) + (0.1 if it.get("rarity", "") == "epic" else 0.0))

static func make_player(ent: Dictionary) -> BlockyModel:
	var cls: String = ent.get("c", "warrior")
	var look: Dictionary = ent.get("lk", {}) if ent.get("lk") is Dictionary else {}
	var eq: Dictionary = ent.get("eq", {}) if ent.get("eq") is Dictionary else {}
	var base: Dictionary = CLASS_LOOK.get(cls, CLASS_LOOK["warrior"])
	var spec := {
		"skin": SKIN[clampi(int(look.get("skin", 0)), 0, SKIN.size() - 1)], "hair": HAIR[clampi(int(look.get("hair", 0)), 0, HAIR.size() - 1)],
		"female": look.get("gender", "m") == "f", "shirt": base["shirt"], "pants": base["pants"], "body": base["body"], "robe": base.get("robe", false),
		"hat": base.get("hat", {}), "props": [], "player_class": cls,
	}
	if base.has("cloak"):
		spec["cloak"] = base["cloak"]
	if base.get("quiver", false):
		spec["quiver"] = true
	# equipment overrides
	var dress := func(slot: String) -> Dictionary:
		var id: String = eq.get(slot, "")
		return S.item(id) if id != "" else {}
	var chest: Dictionary = dress.call("chest")
	if not chest.is_empty():
		spec["shirt"] = armor_color(chest).to_html(false).hex_to_int() if false else _color_int(armor_color(chest))
		spec["robe"] = chest.get("type", "") == "Stoff"
	var legs: Dictionary = dress.call("legs")
	if not legs.is_empty():
		spec["pants"] = _color_int(armor_color(legs))
	var feet: Dictionary = dress.call("feet")
	if not feet.is_empty():
		spec["boots"] = _color_int(armor_color(feet))
	var hands: Dictionary = dress.call("hands")
	if not hands.is_empty():
		spec["glove"] = _color_int(armor_color(hands))
	var waist: Dictionary = dress.call("waist")
	if not waist.is_empty():
		spec["belt"] = _color_int(armor_color(waist))
	var head: Dictionary = dress.call("head")
	if not head.is_empty():
		var typ: String = head.get("type", "")
		spec["hat"] = {"kind": "helmet" if typ == "Kette" else "hood" if typ == "Leder" else "wizardhat", "color": _color_int(armor_color(head))}
	var sh: Dictionary = dress.call("shoulder")
	if not sh.is_empty():
		spec["shoulders"] = _color_int(armor_color(sh))
	elif cls == "warrior":
		spec["shoulders"] = 0x8a8f9c
	var back: Dictionary = dress.call("back")
	if not back.is_empty():
		spec["cloak"] = _color_int(armor_color({"type": "Leder", "req": back.get("req", 1), "rarity": back.get("rarity", "")}))
	var props: Array = []
	var mh: Dictionary = dress.call("mainhand")
	var oh: Dictionary = dress.call("offhand")
	var mhk := _weapon_kind(mh.get("type", ""), base.get("mh", ""))
	if mhk != "":
		props.append([mhk, "left" if mhk in ["bow", "crossbow"] else "right", 1.0, _color_int(_metal(mh))])
	var ohk := _weapon_kind(oh.get("type", ""), base.get("oh", "")) if not oh.is_empty() else base.get("oh", "")
	if ohk != "" and ohk != "quiver":
		props.append([ohk, "left", 1.0, _color_int(_metal(oh)) if not oh.is_empty() else 0x9aa0b0])
	if ohk == "quiver":
		spec["quiver"] = true
	spec["props"] = props
	var mdl := humanoid(spec)
	mdl.player_class = cls
	return mdl

static func _color_int(c: Color) -> int:
	return (int(c.r * 255.0) << 16) | (int(c.g * 255.0) << 8) | int(c.b * 255.0)

static func _metal(it: Dictionary) -> Color:
	if it.is_empty():
		return C(0xb0b8c4)
	var pal := [0x8a8a92, 0xb0b8c4, 0xc8d4e4, 0xe0d070, 0x6a8aff, 0xff7a4a, 0x8aeaff, 0xc070ff]
	var c := C(pal[(tier_of(it) / 3) % pal.size()])
	return lighten(c, 1.0 if it.get("rarity", "") != "epic" else 1.15)

static func _weapon_kind(typ: String, fallback: String) -> String:
	if typ == "":
		return fallback
	if "Schwert" in typ: return "sword"
	if "Axt" in typ: return "axe"
	if "Streitkolben" in typ: return "mace"
	if typ == "Armbrust": return "crossbow"
	if "Bogen" in typ: return "bow"
	if typ == "Zepter": return "scepter"
	if "Stab" in typ: return "staff"
	if "Dolch" in typ: return "dagger"
	if typ == "Schild": return "shield"
	if typ == "Zauberkugel": return "orb"
	if typ == "Zauberbuch": return "tome"
	if typ == "Köcher": return "quiver"
	return fallback

# ---------------------------------------------------------------- humanoid

static func _part(parent: Node3D, name: String, pos: Vector3, mb: MeshBuilder, mat: Material = null) -> Node3D:
	var p := Node3D.new()
	p.name = name
	p.position = pos
	parent.add_child(p)
	if mb != null and not mb.verts.is_empty():
		p.add_child(mb.instance(mat))
	return p

## A bow standing along Y and bulging towards +X, with a string and a nocked arrow that BlockyModel._set_draw moves.
static func _bow_node() -> Node3D:
	var root := Node3D.new()
	var mb := MeshBuilder.new()
	var wood := C(0x6b4a22)
	var arc: Array[Vector2] = []
	for k in 9:
		var y := -0.5 + k / 8.0
		arc.append(Vector2(0.17 * (1.0 - pow(y / 0.5, 2.0)), y))
	for k in 8:
		var a := arc[k]
		var b := arc[k + 1]
		var mid := (a + b) / 2.0
		mb.box(Vector3(mid.x, mid.y, 0), Vector3(0.07, (b - a).length() + 0.03, 0.07), wood, Vector3(0, 0, rad_to_deg(-atan2(b.x - a.x, b.y - a.y))), 0.0)
	mb.box(Vector3(0.17, 0, 0), Vector3(0.1, 0.2, 0.1), C(0x3a2412), Vector3.ZERO, 0.0)
	root.add_child(mb.instance())
	var halves: Array = []
	for i in 2:
		var sb := MeshBuilder.new()
		sb.box(Vector3.ZERO, Vector3(0.012, 1.0, 0.012), C(0xdddddd), Vector3.ZERO, 0.0)
		var h := sb.instance()
		root.add_child(h)
		halves.append(h)
	var ab := MeshBuilder.new()
	ab.box(Vector3(0.35, 0, 0), Vector3(0.7, 0.014, 0.014), C(0xd8d0b8), Vector3.ZERO, 0.0)
	ab.box(Vector3(0.72, 0, 0), Vector3(0.08, 0.05, 0.03), C(0xc0c4cc), Vector3.ZERO, 0.0)
	ab.box(Vector3(0.04, 0, 0), Vector3(0.1, 0.05, 0.012), C(0xc83030), Vector3.ZERO, 0.0)
	var nock := Node3D.new()
	nock.add_child(ab.instance())
	root.add_child(nock)
	root.set_meta("halves", halves)
	root.set_meta("nock", nock)
	return root

static func humanoid(spec: Dictionary) -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = spec.get("kind", "human")
	if mdl.kind == "golem":
		mdl.kind = "golem"
		mdl.speed_scale = 0.7
	mdl.speed_scale = 0.55 if spec.get("slow", false) else mdl.speed_scale
	mdl.arms_forward = spec.get("arms_forward", false)
	var sc: float = spec.get("scale", 1.0) * (0.95 if spec.get("female", false) else 1.0)
	var body: Vector3 = spec.get("body", Vector3.ONE)
	var rig := Node3D.new()
	rig.name = "rig"
	rig.scale = body * sc
	mdl.add_child(rig)
	var skin := C(spec.get("skin", SKIN[0]))
	var hair := C(spec.get("hair", HAIR[0]))
	var shirt := C(spec.get("shirt", 0x888888))
	var pants := C(spec.get("pants", _color_int(shirt.darkened(0.35))))
	var boots := C(spec.get("boots", _color_int(pants.darkened(0.35))))
	var glove := C(spec["glove"]) if spec.has("glove") else skin
	var style: String = spec.get("style", "")
	var emissive: int = spec.get("emissive", 0)
	var opacity: float = spec.get("opacity", 1.0)
	var female: bool = spec.get("female", false)
	if style == "fur":
		shirt = shirt.darkened(0.0)
	if style == "ghost" or opacity < 1.0:
		skin.a = opacity; shirt.a = opacity; pants.a = opacity; boots.a = opacity; glove.a = opacity
	var tw := 0.5 if not female else 0.44
	var arm_w := 0.2 if style != "fur" else 0.28
	var arm_len := 0.68 if not spec.get("arms_long", false) else 0.9
	# legs
	for side in [-1, 1]:
		var mb := MeshBuilder.new()
		mb.box(Vector3(0, -0.375, 0), Vector3(0.23, 0.75, 0.25), pants)
		mb.box(Vector3(0, -0.66, 0.02), Vector3(0.25, 0.18, 0.3), boots)
		if style == "skeleton":
			mb = MeshBuilder.new()
			mb.box(Vector3(0, -0.375, 0), Vector3(0.1, 0.75, 0.1), skin)
		_part(rig, "leg_l" if side < 0 else "leg_r", Vector3(side * 0.13, 0.75, 0), mb, MeshBuilder.alpha_material() if opacity < 1.0 else null)
	# torso
	var tb := MeshBuilder.new()
	if style == "skeleton":
		tb.box(Vector3(0, 0.35, 0), Vector3(0.36, 0.6, 0.16), skin)
		for i in 4:
			tb.box(Vector3(0, 0.15 + i * 0.13, 0.06), Vector3(0.4, 0.04, 0.12), skin.darkened(0.15))
	else:
		tb.box(Vector3(0, 0.35, 0), Vector3(tw, 0.7, 0.28), shirt)
		if female:
			tb.box(Vector3(0, 0.0, 0), Vector3(0.5, 0.18, 0.3), shirt.darkened(0.05))
		if spec.get("robe", false):
			tb.frustum(Vector3(0, -0.55, 0), 0.72, 0.5, 0.52, 0.3, 0.62, shirt.darkened(0.08))
		# belt + collar
		tb.box(Vector3(0, 0.05, 0), Vector3(tw + 0.02, 0.07, 0.3), C(spec["belt"]) if spec.has("belt") else shirt.darkened(0.35))
		tb.box(Vector3(0, 0.05, 0.15), Vector3(0.1, 0.08, 0.02), C(0xd8c060))
		if style == "stone" or style == "lava" or style == "obsidian":
			for i in 3:
				tb.box(Vector3(-0.1 + i * 0.1, 0.4, 0.145), Vector3(0.07, 0.3, 0.02), C(0xff7a1a if style == "lava" else 0x6a6a74))
	if spec.has("shoulders"):
		var sc2 := C(spec["shoulders"])
		for side in [-1, 1]:
			tb.box(Vector3(side * (tw / 2 + 0.07), 0.68, 0), Vector3(0.26, 0.14, 0.3), sc2)
			tb.box(Vector3(side * (tw / 2 + 0.1), 0.74, 0), Vector3(0.14, 0.1, 0.2), lighten(sc2, 1.15))
	if spec.get("quiver", false):
		var q := _prop("quiver", 0x6a4a2a)
		tb.append(q, Vector3(-0.1, 0.35, -0.2))
	var torso := _part(rig, "torso", Vector3(0, 0.75, 0), tb, MeshBuilder.alpha_material() if opacity < 1.0 else null)
	if spec.has("cloak"):
		var cb := MeshBuilder.new()
		cb.frustum(Vector3(0, -0.85, -0.17), 0.52, 0.04, 0.46, 0.05, 1.0, C(spec["cloak"]), Vector3.ZERO, 0.03)
		var cl := _part(torso, "cloak", Vector3(0, 0.7, 0), cb)
		cl.position = Vector3(0, 0.68, -0.14)
		(cl.get_child(0) as Node3D).position = Vector3(0, 0, 0)
		mdl.parts["cloak"] = cl
	# head
	var hb := MeshBuilder.new()
	var hs: float = spec.get("head", 1.0)
	hb.box(Vector3(0, 0.2, 0), Vector3(0.4, 0.4, 0.4), skin, Vector3.ZERO, 0.02)
	if style == "lizard":
		hb.box(Vector3(0, 0.15, 0.3), Vector3(0.3, 0.2, 0.25), skin.darkened(0.05))
		hb.box(Vector3(0, 0.3, -0.05), Vector3(0.1, 0.25, 0.4), skin.darkened(0.2))
	elif style == "demon":
		hb.box(Vector3(0, 0.1, 0.22), Vector3(0.3, 0.12, 0.08), skin.darkened(0.3))
	elif style == "fur":
		hb.box(Vector3(0, 0.2, 0.0), Vector3(0.44, 0.44, 0.44), shirt, Vector3.ZERO, 0.06)
		hb.box(Vector3(0, 0.1, 0.22), Vector3(0.26, 0.2, 0.1), skin.darkened(0.25))
	elif style == "skeleton":
		hb.box(Vector3(0, 0.1, 0.2), Vector3(0.3, 0.12, 0.05), skin.darkened(0.2))
	# hair
	if not spec.get("bald", false) and style == "" and spec.get("hat", {}).get("kind", "") not in ["helmet", "hood", "wizardhat", "witchhat", "turban", "captainhat", "straw", "pharaoh"]:
		hb.box(Vector3(0, 0.395, -0.01), Vector3(0.42, 0.1, 0.42), hair)
		hb.box(Vector3(0, 0.22, -0.2), Vector3(0.42, 0.4, 0.05), hair)
		if female:
			hb.box(Vector3(0, 0.0, -0.22), Vector3(0.36, 0.5, 0.06), hair)
	# face
	var eye := C(0x1a1a1a)
	var glow_eye := emissive != 0 or style == "demon" or style == "ghost" or style == "lava"
	if glow_eye:
		eye = C(0xffd040 if style != "ghost" else 0xc0ffff)
	if style != "skeleton":
		hb.box(Vector3(-0.1, 0.24, 0.205), Vector3(0.07, 0.07, 0.02), eye, Vector3.ZERO, 0.0)
		hb.box(Vector3(0.1, 0.24, 0.205), Vector3(0.07, 0.07, 0.02), eye, Vector3.ZERO, 0.0)
		hb.box(Vector3(0, 0.12, 0.205), Vector3(0.1, 0.03, 0.02), skin.darkened(0.3), Vector3.ZERO, 0.0)
	else:
		hb.box(Vector3(-0.1, 0.24, 0.2), Vector3(0.08, 0.1, 0.03), C(0x20180c), Vector3.ZERO, 0.0)
		hb.box(Vector3(0.1, 0.24, 0.2), Vector3(0.08, 0.1, 0.03), C(0x20180c), Vector3.ZERO, 0.0)
	var hat: Dictionary = spec.get("hat", {})
	var hat_h := 0.0
	if hat.has("kind"):
		hat_h = _hat(hb, hat["kind"], C(hat["color"]) if hat.has("color") else shirt.darkened(0.2))
	var head := _part(rig, "head", Vector3(0, 1.45, 0), hb, MeshBuilder.alpha_material() if opacity < 1.0 else null)
	head.scale = Vector3.ONE * hs
	# arms
	var props: Array = spec.get("props", [])
	for side in [-1, 1]:
		var ab := MeshBuilder.new()
		var sleeve := shirt.darkened(0.05) if style != "skeleton" else skin
		if style == "skeleton":
			ab.box(Vector3(0, -arm_len / 2, 0), Vector3(0.08, arm_len, 0.08), skin)
		else:
			ab.box(Vector3(0, -arm_len / 2 + 0.02, 0), Vector3(arm_w, arm_len - 0.1, arm_w), sleeve)
			ab.box(Vector3(0, -arm_len + 0.06, 0), Vector3(arm_w - 0.02, 0.14, arm_w - 0.02), glove)
		var ap := _part(rig, "arm_l" if side < 0 else "arm_r", Vector3(side * (tw / 2 + arm_w / 2 + 0.01), 1.42, 0), ab, MeshBuilder.alpha_material() if opacity < 1.0 else null)
		# held items
		for p in props:
			var slot: String = p[1]
			if (slot == "right" and side > 0) or (slot == "left" and side < 0):
				var kind: String = p[0]
				var psc: float = p[2] if p.size() > 2 else 1.0
				var pcol: int = p[3] if p.size() > 3 else 0
				if kind == "bow":
					# the arm points forward: hand-local +Z is up and -Y is forward, so the bow stands upright, belly towards the target
					var bh := Node3D.new()
					bh.position = Vector3(0, -arm_len + 0.06, 0.02)
					bh.basis = Basis(Vector3(0, -1, 0), Vector3(0, 0, 1), Vector3(-1, 0, 0)).scaled(Vector3.ONE * psc * 1.25)
					var bow_node := _bow_node()
					bh.add_child(bow_node)
					ap.add_child(bh)
					mdl.bow = bow_node
					mdl.set_meta("bow_delay", 0.3)   # the arrow leaves the string this long after the shot started
					continue
				var pm := _prop(kind, pcol)
				var holder := Node3D.new()
				holder.position = Vector3(0, -arm_len + 0.06, 0.05)
				if kind in ["shield", "tome", "orb"]:
					holder.position = Vector3(-0.1 if side < 0 else 0.1, -arm_len + 0.25, 0.08)
					holder.rotation_degrees = Vector3(0, 0, 0)
				elif kind in ["bow", "crossbow"]:
					holder.rotation_degrees = Vector3(-80, 0, 0)
					holder.position = Vector3(0, -arm_len + 0.06, 0.08)
				else:
					holder.rotation_degrees = Vector3(-75, 0, 0)
				holder.scale = Vector3.ONE * psc
				holder.add_child(pm.instance(MeshBuilder.glow_material() if kind in ["orb"] else null))
				ap.add_child(holder)
	# tails etc
	if style == "lizard" or style == "demon":
		var tb2 := MeshBuilder.new()
		tb2.frustum(Vector3(0, 0, 0), 0.22, 0.22, 0.1, 0.1, 0.0001, skin)
		var tl := MeshBuilder.new()
		tl.box(Vector3(0, 0, -0.35), Vector3(0.2, 0.2, 0.7), skin.darkened(0.05))
		tl.box(Vector3(0, 0, -0.85), Vector3(0.12, 0.12, 0.4), skin.darkened(0.1))
		_part(rig, "tail", Vector3(0, 0.78, -0.12), tl)
	if style == "demon":
		var wb := MeshBuilder.new()
		for side in [-1, 1]:
			wb.box(Vector3(side * 0.45, 0.1, -0.1), Vector3(0.7, 0.05, 0.35), skin.darkened(0.2), Vector3(0, side * 20, side * -25))
		_part(rig, "wings", Vector3(0, 1.1, -0.18), wb)
	var is_big: bool = mdl.kind == "golem"
	mdl.height = (1.85 + hat_h) * sc * body.y * (1.0 if not is_big else 1.0)
	mdl.player_class = spec.get("player_class", "")
	if emissive != 0 and style != "ghost":
		pass
	return mdl

# ---------------------------------------------------------------- hats and held items

static func _hat(mb: MeshBuilder, kind: String, col: Color) -> float:
	match kind:
		"helmet":
			mb.box(Vector3(0, 0.25, 0), Vector3(0.46, 0.34, 0.46), col, Vector3.ZERO, 0.02)
			mb.box(Vector3(0, 0.44, 0), Vector3(0.44, 0.08, 0.44), lighten(col, 1.1))
			mb.box(Vector3(0, 0.2, 0.235), Vector3(0.3, 0.07, 0.02), C(0x101018))
			mb.box(Vector3(0, 0.58, -0.02), Vector3(0.06, 0.18, 0.3), C(0xb83030))
			return 0.3
		"hood":
			mb.box(Vector3(0, 0.26, -0.02), Vector3(0.46, 0.36, 0.46), col)
			mb.box(Vector3(0, 0.1, -0.2), Vector3(0.46, 0.5, 0.12), col.darkened(0.1))
			mb.box(Vector3(0, 0.46, 0), Vector3(0.36, 0.08, 0.36), col.darkened(0.12))
			return 0.15
		"wizardhat":
			mb.box(Vector3(0, 0.41, 0), Vector3(0.7, 0.04, 0.7), col)
			mb.frustum(Vector3(0, 0.42, 0), 0.4, 0.4, 0.2, 0.2, 0.3, col)
			mb.frustum(Vector3(0, 0.72, 0), 0.2, 0.2, 0.06, 0.06, 0.28, col.darkened(0.05), Vector3.ZERO, 0.04, Vector2(0.02, 0.0))
			mb.box(Vector3(0, 0.45, 0), Vector3(0.42, 0.05, 0.42), C(0xd8c060))
			return 0.65
		"witchhat":
			mb.box(Vector3(0, 0.41, 0), Vector3(0.75, 0.04, 0.75), C(0x1a1424))
			mb.frustum(Vector3(0, 0.42, 0), 0.42, 0.42, 0.2, 0.2, 0.35, C(0x1a1424))
			mb.frustum(Vector3(0, 0.77, 0), 0.2, 0.2, 0.05, 0.05, 0.3, C(0x1a1424), Vector3.ZERO, 0.04, Vector2(0.06, 0))
			mb.box(Vector3(0, 0.46, 0), Vector3(0.43, 0.06, 0.43), C(0xb040b0))
			return 0.75
		"straw":
			mb.box(Vector3(0, 0.41, 0), Vector3(0.85, 0.03, 0.85), C(0xe0c070))
			mb.box(Vector3(0, 0.5, 0), Vector3(0.4, 0.14, 0.4), C(0xe8cc80))
			return 0.25
		"crown":
			mb.box(Vector3(0, 0.44, 0), Vector3(0.44, 0.1, 0.44), C(0xf0c030))
			for i in 4:
				var a := i * PI / 2.0
				mb.box(Vector3(cos(a) * 0.19, 0.55, sin(a) * 0.19), Vector3(0.08, 0.14, 0.08), C(0xf4d050))
			mb.box(Vector3(0, 0.46, 0.225), Vector3(0.08, 0.08, 0.02), C(0xd02030))
			return 0.3
		"horns":
			for side in [-1, 1]:
				mb.box(Vector3(side * 0.2, 0.46, 0), Vector3(0.08, 0.2, 0.08), C(0xe8e0c8), Vector3(0, 0, side * -25))
				mb.box(Vector3(side * 0.28, 0.6, 0), Vector3(0.06, 0.16, 0.06), C(0xf0e8d0), Vector3(0, 0, side * -45))
			return 0.3
		"turban":
			mb.box(Vector3(0, 0.46, 0), Vector3(0.48, 0.18, 0.48), col)
			mb.box(Vector3(0, 0.6, 0), Vector3(0.38, 0.12, 0.38), col.darkened(0.08))
			mb.box(Vector3(0, 0.45, 0.245), Vector3(0.1, 0.1, 0.02), C(0xd02030))
			return 0.3
		"mask":
			mb.box(Vector3(0, 0.1, 0.215), Vector3(0.4, 0.16, 0.04), col if col != Color() else C(0x20202a))
			mb.box(Vector3(0, 0.42, 0), Vector3(0.44, 0.06, 0.44), col.darkened(0.2))
			return 0.1
		"ears":
			for side in [-1, 1]:
				mb.box(Vector3(side * 0.32, 0.3, 0), Vector3(0.18, 0.1, 0.06), col, Vector3(0, 0, side * 20))
			return 0.1
		"pharaoh":
			mb.box(Vector3(0, 0.3, 0), Vector3(0.46, 0.3, 0.46), C(0xf0c030), Vector3.ZERO, 0.02)
			mb.box(Vector3(-0.26, 0.05, 0), Vector3(0.1, 0.5, 0.3), C(0x2a50a8))
			mb.box(Vector3(0.26, 0.05, 0), Vector3(0.1, 0.5, 0.3), C(0x2a50a8))
			mb.box(Vector3(0, 0.52, 0), Vector3(0.1, 0.12, 0.1), C(0x20a050))
			return 0.3
		"captainhat":
			mb.box(Vector3(0, 0.44, 0), Vector3(0.6, 0.08, 0.5), C(0x1a2a50))
			mb.box(Vector3(0, 0.52, 0), Vector3(0.44, 0.14, 0.4), C(0x1a2a50))
			mb.box(Vector3(0, 0.5, 0.22), Vector3(0.1, 0.08, 0.02), C(0xf0c030))
			return 0.25
	return 0.0

static func _prop(kind: String, tint: int) -> MeshBuilder:
	var mb := MeshBuilder.new()
	var steel := C(tint) if tint != 0 else C(0xb0b8c4)
	var wood := C(0x6a4a2a)
	var gold := C(0xd8c060)
	match kind:
		"sword":
			mb.box(Vector3(0, 0.09, 0), Vector3(0.05, 0.2, 0.05), wood)
			mb.box(Vector3(0, 0.2, 0), Vector3(0.28, 0.05, 0.08), gold)
			mb.box(Vector3(0, 0.67, 0), Vector3(0.08, 0.9, 0.025), steel, Vector3.ZERO, 0.02)
			mb.frustum(Vector3(0, 1.1, 0), 0.08, 0.025, 0.0, 0.0, 0.14, steel)
		"dagger":
			mb.box(Vector3(0, 0.06, 0), Vector3(0.04, 0.14, 0.04), wood)
			mb.box(Vector3(0, 0.15, 0), Vector3(0.16, 0.04, 0.06), gold)
			mb.box(Vector3(0, 0.38, 0), Vector3(0.06, 0.42, 0.02), steel)
		"scimitar":
			mb.box(Vector3(0, 0.08, 0), Vector3(0.05, 0.18, 0.05), wood)
			mb.box(Vector3(0, 0.19, 0), Vector3(0.2, 0.04, 0.06), gold)
			mb.box(Vector3(0, 0.55, 0), Vector3(0.08, 0.7, 0.02), steel, Vector3(0, 0, -6))
			mb.box(Vector3(0.05, 0.95, 0), Vector3(0.08, 0.24, 0.02), steel, Vector3(0, 0, -22))
		"axe":
			mb.box(Vector3(0, 0.4, 0), Vector3(0.06, 0.9, 0.06), wood)
			mb.frustum(Vector3(0.06, 0.62, -0.025), 0.02, 0.05, 0.3, 0.05, 0.0001, steel)
			mb.box(Vector3(0.14, 0.72, 0), Vector3(0.28, 0.24, 0.04), steel, Vector3(0, 0, 0))
			mb.box(Vector3(-0.05, 0.72, 0), Vector3(0.08, 0.2, 0.05), steel.darkened(0.2))
		"mace":
			mb.box(Vector3(0, 0.35, 0), Vector3(0.06, 0.75, 0.06), wood)
			mb.box(Vector3(0, 0.82, 0), Vector3(0.2, 0.2, 0.2), steel)
			for s in [-1, 1]:
				mb.box(Vector3(s * 0.12, 0.82, 0), Vector3(0.06, 0.08, 0.08), steel.darkened(0.2))
				mb.box(Vector3(0, 0.82, s * 0.12), Vector3(0.08, 0.08, 0.06), steel.darkened(0.2))
		"spear":
			mb.box(Vector3(0, 0.7, 0), Vector3(0.05, 1.7, 0.05), wood)
			mb.frustum(Vector3(0, 1.5, 0), 0.1, 0.03, 0.0, 0.0, 0.3, steel)
		"staff", "scepter":
			mb.box(Vector3(0, 0.65, 0), Vector3(0.06, 1.5, 0.06), C(0x5a3a1e))
			var orb := C(tint) if tint != 0 else C(0x60c8ff)
			if kind == "scepter":
				mb.box(Vector3(0, 1.45, 0), Vector3(0.14, 0.14, 0.14), gold, Vector3(45, 45, 0))
			mb.box(Vector3(0, 1.5, 0), Vector3(0.16, 0.16, 0.16), Color(orb.r, orb.g, orb.b), Vector3(45, 0, 45))
			for s in [-1, 1]:
				mb.box(Vector3(s * 0.08, 1.42, 0), Vector3(0.04, 0.2, 0.04), C(0x5a3a1e), Vector3(0, 0, s * -20))
		"witchstaff":
			mb.box(Vector3(0, 0.6, 0), Vector3(0.06, 1.4, 0.06), C(0x3a2a20), Vector3(0, 0, 4))
			mb.box(Vector3(0.05, 1.35, 0), Vector3(0.06, 0.3, 0.06), C(0x3a2a20), Vector3(0, 0, -30))
			mb.box(Vector3(0.14, 1.52, 0), Vector3(0.16, 0.16, 0.16), C(0x70ff50), Vector3(45, 0, 45))
		"bow":
			var n := 6
			for i in n:
				var a := (float(i) / (n - 1) - 0.5) * 1.6
				mb.box(Vector3(0, a * 0.5 + 0.1, sin(a * 1.2) * -0.02 + (1.0 - absf(a)) * 0.0 + 0.08 - absf(a) * 0.12), Vector3(0.04, 0.2, 0.04), wood, Vector3(a * 28, 0, 0))
			mb.box(Vector3(0, 0.1, -0.04), Vector3(0.01, 0.82, 0.01), C(0xe0e0d0))
		"crossbow":
			mb.box(Vector3(0, 0.1, 0), Vector3(0.06, 0.5, 0.06), wood)
			mb.box(Vector3(0, 0.3, 0), Vector3(0.5, 0.05, 0.05), steel.darkened(0.2))
			mb.box(Vector3(0, 0.3, 0.02), Vector3(0.48, 0.01, 0.01), C(0xe0e0d0))
		"rifle":
			mb.box(Vector3(0, 0.45, 0), Vector3(0.06, 1.0, 0.06), C(0x2a2a30))
			mb.box(Vector3(0, 0.0, 0), Vector3(0.08, 0.3, 0.1), wood)
		"pick":
			mb.box(Vector3(0, 0.4, 0), Vector3(0.06, 0.85, 0.06), wood)
			mb.box(Vector3(0, 0.82, 0), Vector3(0.5, 0.07, 0.06), steel.darkened(0.2))
			mb.frustum(Vector3(0.25, 0.82, 0), 0.0001, 0.06, 0.08, 0.06, 0.0001, steel)
		"shield":
			mb.box(Vector3(0, 0, 0), Vector3(0.07, 0.62, 0.5), steel.darkened(0.05))
			mb.box(Vector3(0.05, 0, 0), Vector3(0.06, 0.26, 0.26), gold)
			mb.box(Vector3(0.05, 0.24, 0), Vector3(0.04, 0.06, 0.4), steel.darkened(0.3))
		"orb":
			mb.box(Vector3(0, 0, 0), Vector3(0.22, 0.22, 0.22), C(0x60b0ff), Vector3(45, 0, 45))
			mb.box(Vector3(0, -0.15, 0), Vector3(0.14, 0.06, 0.14), gold)
		"tome":
			mb.box(Vector3(0, 0, 0), Vector3(0.08, 0.34, 0.26), C(0x5a2a8a))
			mb.box(Vector3(0.0, 0, 0.0), Vector3(0.1, 0.3, 0.22), C(0xe8dcb0))
			mb.box(Vector3(0.05, 0, 0), Vector3(0.02, 0.1, 0.1), gold)
		"quiver":
			mb.box(Vector3(0, 0, 0), Vector3(0.14, 0.55, 0.14), C(0x6a4a2a), Vector3(0, 0, 18))
			for i in 3:
				mb.box(Vector3(-0.04 + i * 0.04, 0.34, 0), Vector3(0.02, 0.2, 0.02), C(0xe0e0d0), Vector3(0, 0, 18))
	return mb

# ---------------------------------------------------------------- animals

static func make_quad(spec: Dictionary, id: int) -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = "quad"
	var s: float = spec.get("s", 1.0)
	var c := C(spec["c"])
	var c2 := C(spec["c2"])
	var low: bool = spec.get("low", false)
	var stocky: bool = spec.get("stocky", false)
	var leg_h := 0.18 if low else (0.42 if not stocky else 0.4)
	var bw := 0.5 if not stocky else 0.7
	var bh := 0.4 if not stocky else 0.55
	var bl := 1.0 if not spec.get("short", false) else 0.8
	var root := Node3D.new()
	root.scale = Vector3.ONE * s
	mdl.add_child(root)
	var bb := MeshBuilder.new()
	bb.box(Vector3(0, leg_h + bh / 2, 0), Vector3(bw, bh, bl), c)
	bb.box(Vector3(0, leg_h + bh / 2 - bh * 0.3, 0), Vector3(bw + 0.02, bh * 0.45, bl * 0.8), c2)
	if spec.get("wool", false):
		bb.box(Vector3(0, leg_h + bh / 2 + 0.05, 0), Vector3(bw + 0.14, bh + 0.1, bl + 0.1), c, Vector3.ZERO, 0.08)
	if spec.get("hump", false):
		bb.box(Vector3(0, leg_h + bh + 0.1, 0.15), Vector3(bw, 0.3, 0.5), c.darkened(0.1))
	if spec.get("slope", false):
		bb.box(Vector3(0, leg_h + bh / 2 + 0.06, 0.1), Vector3(bw, bh, bl * 0.6), c.darkened(0.08))
	if spec.has("mane"):
		bb.box(Vector3(0, leg_h + bh / 2 + 0.05, bl / 2 - 0.1), Vector3(bw + 0.14, bh + 0.14, 0.34), C(spec["mane"]))
	if spec.has("spots"):
		var r := RandomNumberGenerator.new()
		r.seed = id
		for i in 7:
			bb.box(Vector3(r.randf_range(-0.2, 0.2), leg_h + bh + 0.005, r.randf_range(-0.4, 0.4)), Vector3(0.1, 0.02, 0.1), C(spec["spots"]), Vector3.ZERO, 0.0)
	if spec.has("stripes"):
		for i in 5:
			bb.box(Vector3(0, leg_h + bh / 2, -0.38 + i * 0.19), Vector3(bw + 0.03, bh + 0.03, 0.06), C(spec["stripes"]), Vector3.ZERO, 0.0)
	if spec.get("panda", false):
		bb.box(Vector3(0, leg_h + bh / 2, -0.28), Vector3(bw + 0.03, bh + 0.03, 0.34), c2)
		bb.box(Vector3(0, leg_h + bh / 2, 0.35), Vector3(bw + 0.03, bh + 0.03, 0.3), c2)
	if spec.has("cracks"):
		for i in 4:
			bb.box(Vector3(-0.15 + i * 0.1, leg_h + bh + 0.005, -0.3 + i * 0.2), Vector3(0.05, 0.02, 0.2), C(spec["cracks"]), Vector3(0, 20 * i, 0), 0.0)
	var body := _part(root, "body", Vector3.ZERO, bb)
	mdl.parts["body"] = body
	# head
	var hb := MeshBuilder.new()
	var hsz := 0.34 if not stocky else 0.42
	hb.box(Vector3(0, 0, 0), Vector3(hsz, hsz, hsz), c)
	if spec.get("snout", false):
		hb.box(Vector3(0, -hsz * 0.15, hsz * 0.6), Vector3(hsz * 0.62, hsz * 0.45, hsz * 0.5), c2 if not spec.get("stocky", false) else c.lightened(0.1) if false else c2)
		hb.box(Vector3(0, -hsz * 0.05, hsz * 0.86), Vector3(0.08, 0.07, 0.04), C(0x1a1414), Vector3.ZERO, 0.0)
	var ecol := C(spec["eyes"]) if spec.has("eyes") else C(0x101010)
	hb.box(Vector3(-hsz * 0.28, hsz * 0.12, hsz * 0.5), Vector3(0.06, 0.06, 0.03), ecol, Vector3.ZERO, 0.0)
	hb.box(Vector3(hsz * 0.28, hsz * 0.12, hsz * 0.5), Vector3(0.06, 0.06, 0.03), ecol, Vector3.ZERO, 0.0)
	if spec.get("ears", false):
		for side in [-1, 1]:
			hb.box(Vector3(side * hsz * 0.34, hsz * 0.55, -0.02), Vector3(0.1, 0.14, 0.06), c.darkened(0.1))
	if spec.get("tusks", false):
		for side in [-1, 1]:
			hb.box(Vector3(side * 0.1, -hsz * 0.28, hsz * 0.78), Vector3(0.04, 0.12, 0.04), C(0xf0e8d0), Vector3(-30, 0, side * 15), 0.0)
	if spec.get("horns", false):
		for side in [-1, 1]:
			hb.box(Vector3(side * hsz * 0.55, hsz * 0.4, 0), Vector3(0.22, 0.07, 0.07), C(0xe0d8c0), Vector3(0, 0, side * 25))
			hb.box(Vector3(side * hsz * 0.78, hsz * 0.6, 0), Vector3(0.07, 0.18, 0.07), C(0xe8e0c8), Vector3(0, 0, side * -10))
	var head := _part(root, "head", Vector3(0, leg_h + bh * 0.85, bl / 2 + hsz * 0.25), hb)
	mdl.parts["head"] = head
	# legs
	var lw := 0.16 if not stocky else 0.22
	for leg in [["leg_fl", -1, 1], ["leg_fr", 1, 1], ["leg_bl", -1, -1], ["leg_br", 1, -1]]:
		var lb := MeshBuilder.new()
		lb.box(Vector3(0, -leg_h / 2, 0), Vector3(lw, leg_h, lw), c.darkened(0.05))
		lb.box(Vector3(0, -leg_h + 0.04, 0.02), Vector3(lw + 0.01, 0.08, lw + 0.04), c.darkened(0.3))
		_part(root, leg[0], Vector3(leg[1] * (bw / 2 - lw / 2 - 0.01), leg_h, leg[2] * (bl / 2 - lw / 2 - 0.04)), lb)
		mdl.parts[leg[0]] = root.get_child(root.get_child_count() - 1)
	# tail
	var tl: float = spec.get("tail", 0.4)
	if tl > 0.0:
		var tb := MeshBuilder.new()
		var tcol := C(spec["tailtip"]) if spec.has("tailtip") else c
		tb.box(Vector3(0, 0, -tl / 2), Vector3(0.1, 0.1, tl), c)
		if spec.has("tailtip"):
			tb.box(Vector3(0, 0, -tl + 0.04), Vector3(0.12, 0.12, 0.16), tcol)
		if spec.get("bigtail", false):
			tb.box(Vector3(0, 0.05, -tl / 2), Vector3(0.22, 0.22, tl), c)
		_part(root, "tail", Vector3(0, leg_h + bh * 0.8, -bl / 2), tb)
		mdl.parts["tail"] = root.get_child(root.get_child_count() - 1)
	mdl.height = (leg_h + bh + 0.3) * s
	return mdl

static func make_spider(col: int, eyes: int, scale := 1.0) -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = "spider"
	var root := Node3D.new()
	root.scale = Vector3.ONE * scale
	mdl.add_child(root)
	var c := C(col)
	var bb := MeshBuilder.new()
	bb.box(Vector3(0, 0.45, 0.15), Vector3(0.6, 0.4, 0.6), c)
	bb.box(Vector3(0, 0.55, -0.55), Vector3(0.8, 0.55, 0.9), c.darkened(0.15))
	bb.box(Vector3(0, 0.62, -0.6), Vector3(0.3, 0.08, 0.5), C(eyes).darkened(0.3))
	for i in 4:
		bb.box(Vector3(-0.12 + (i % 2) * 0.24, 0.52 + (i / 2) * 0.1, 0.46), Vector3(0.09, 0.09, 0.04), C(eyes), Vector3.ZERO, 0.0)
	var body := _part(root, "body", Vector3.ZERO, bb)
	mdl.parts["body"] = body
	for i in 8:
		var side := 1.0 if i < 4 else -1.0
		var k := i % 4
		var lb := MeshBuilder.new()
		lb.box(Vector3(side * 0.35, 0.1, 0), Vector3(0.7, 0.07, 0.07), c.lightened(0.05) if false else c, Vector3(0, 0, side * -30), 0.0)
		lb.box(Vector3(side * 0.78, -0.18, 0), Vector3(0.07, 0.5, 0.07), c.darkened(0.1), Vector3.ZERO, 0.0)
		var lp := _part(root, "leg%d" % i, Vector3(side * 0.25, 0.5, 0.35 - k * 0.28), lb)
		lp.rotation.y = side * (-0.5 + k * 0.33)
		mdl.parts["leg%d" % i] = lp
	mdl.height = 1.0 * scale
	return mdl

static func make_scorpion(col: int, scale: float) -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = "scorpion"
	var root := Node3D.new()
	root.scale = Vector3.ONE * scale
	mdl.add_child(root)
	var c := C(col)
	var bb := MeshBuilder.new()
	bb.box(Vector3(0, 0.3, 0), Vector3(0.6, 0.24, 0.9), c)
	bb.box(Vector3(0, 0.34, 0.5), Vector3(0.3, 0.2, 0.3), c.darkened(0.1))
	for side in [-1, 1]:
		bb.box(Vector3(side * 0.42, 0.3, 0.7), Vector3(0.12, 0.1, 0.5), c.darkened(0.1), Vector3(0, side * -20, 0))
		bb.box(Vector3(side * 0.5, 0.3, 1.0), Vector3(0.3, 0.12, 0.24), c.darkened(0.2), Vector3(0, side * 25, 0))
		bb.box(Vector3(side * 0.12, 0.42, 0.62), Vector3(0.06, 0.06, 0.03), C(0x101010), Vector3.ZERO, 0.0)
	var body := _part(root, "body", Vector3.ZERO, bb)
	mdl.parts["body"] = body
	var tb := MeshBuilder.new()
	for i in 4:
		tb.box(Vector3(0, 0.1 + i * 0.16, -0.12 - i * 0.04 + (0.0 if i < 3 else 0.14)), Vector3(0.14 - i * 0.015, 0.2, 0.14 - i * 0.015), c.darkened(0.05 * i), Vector3(-15 * i, 0, 0))
	tb.box(Vector3(0, 0.78, 0.12), Vector3(0.05, 0.08, 0.16), C(0x2a1a14), Vector3(-60, 0, 0))
	var tail := _part(root, "tail", Vector3(0, 0.35, -0.45), tb)
	mdl.parts["tail"] = tail
	for i in 8:
		var side := 1.0 if i < 4 else -1.0
		var k := i % 4
		var lb := MeshBuilder.new()
		lb.box(Vector3(side * 0.28, -0.04, 0), Vector3(0.5, 0.05, 0.05), c.darkened(0.2), Vector3(0, 0, side * -25), 0.0)
		var lp := _part(root, "leg%d" % i, Vector3(side * 0.25, 0.28, 0.3 - k * 0.25), lb)
		mdl.parts["leg%d" % i] = lp
	mdl.height = 0.9 * scale
	return mdl

static func make_croc() -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = "croc"
	var root := Node3D.new()
	mdl.add_child(root)
	var c := C(0x4a6a32)
	var bb := MeshBuilder.new()
	bb.box(Vector3(0, 0.3, 0), Vector3(0.7, 0.34, 1.5), c)
	bb.box(Vector3(0, 0.5, 0), Vector3(0.5, 0.06, 1.3), c.darkened(0.15))
	bb.box(Vector3(0, 0.22, 0), Vector3(0.64, 0.2, 1.4), C(0xb8c070))
	var body := _part(root, "body", Vector3.ZERO, bb)
	mdl.parts["body"] = body
	var hb := MeshBuilder.new()
	hb.box(Vector3(0, 0.0, 0.5), Vector3(0.46, 0.2, 1.0), c)
	for i in 5:
		hb.box(Vector3(-0.2, -0.12, 0.2 + i * 0.18), Vector3(0.04, 0.07, 0.04), C(0xf0f0e0), Vector3.ZERO, 0.0)
		hb.box(Vector3(0.2, -0.12, 0.2 + i * 0.18), Vector3(0.04, 0.07, 0.04), C(0xf0f0e0), Vector3.ZERO, 0.0)
	hb.box(Vector3(-0.14, 0.14, 0.1), Vector3(0.1, 0.1, 0.1), C(0xe0d020), Vector3.ZERO, 0.0)
	hb.box(Vector3(0.14, 0.14, 0.1), Vector3(0.1, 0.1, 0.1), C(0xe0d020), Vector3.ZERO, 0.0)
	_part(root, "head", Vector3(0, 0.34, 0.75), hb)
	var jb := MeshBuilder.new()
	jb.box(Vector3(0, -0.08, 0.4), Vector3(0.42, 0.08, 0.8), C(0xb8c070))
	var jaw := _part(root, "jaw", Vector3(0, 0.3, 0.8), jb)
	mdl.parts["jaw"] = jaw
	var tb := MeshBuilder.new()
	tb.frustum(Vector3(0, 0, 0), 0.5, 0.5, 0.18, 0.18, 0.001, c)
	tb.box(Vector3(0, 0, -0.5), Vector3(0.4, 0.28, 1.0), c.darkened(0.05))
	tb.box(Vector3(0, 0, -1.2), Vector3(0.24, 0.2, 0.6), c.darkened(0.1))
	var tail := _part(root, "tail", Vector3(0, 0.3, -0.7), tb)
	mdl.parts["tail"] = tail
	for leg in [["leg_fl", -1, 1], ["leg_fr", 1, 1], ["leg_bl", -1, -1], ["leg_br", 1, -1]]:
		var lb := MeshBuilder.new()
		lb.box(Vector3(leg[1] * 0.18, -0.1, 0), Vector3(0.2, 0.2, 0.2), c.darkened(0.1))
		var lp := _part(root, leg[0], Vector3(leg[1] * 0.28, 0.2, leg[2] * 0.5), lb)
		mdl.parts[leg[0]] = lp
	mdl.height = 0.7
	return mdl

static func make_wisp(col: int) -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = "wisp"
	var c := C(col)
	var core := MeshBuilder.new()
	core.box(Vector3.ZERO, Vector3(0.34, 0.34, 0.34), Color(c.r, c.g, c.b), Vector3(45, 0, 45), 0.0)
	var m1 := core.instance(MeshBuilder.glow_material())
	mdl.add_child(m1)
	var shell := MeshBuilder.new()
	shell.box(Vector3.ZERO, Vector3(0.6, 0.6, 0.6), Color(c.r, c.g, c.b, 0.3), Vector3(0, 45, 20), 0.0)
	var m2 := shell.instance(MeshBuilder.alpha_material())
	mdl.add_child(m2)
	var light := OmniLight3D.new()
	light.light_color = c
	light.omni_range = 7.0
	light.light_energy = 1.2
	mdl.add_child(light)
	mdl.glow_nodes = [m1, m2]
	mdl.height = 1.6
	return mdl

static func make_griffin() -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = "griffin"
	var c := C(0xc8923a)
	var wing_c := C(0xe8dcc0)
	var bb := MeshBuilder.new()
	bb.box(Vector3(0, 0.9, 0), Vector3(0.7, 0.6, 1.5), c)
	bb.box(Vector3(0, 0.78, 0.0), Vector3(0.72, 0.3, 1.2), c.lightened(0.15) if false else C(0xd8a85a))
	var body := _part(mdl, "body", Vector3.ZERO, bb)
	var hb := MeshBuilder.new()
	hb.box(Vector3(0, 0, 0), Vector3(0.4, 0.4, 0.4), wing_c)
	hb.box(Vector3(0, -0.08, 0.3), Vector3(0.16, 0.14, 0.34), C(0xf0c020))
	hb.box(Vector3(-0.12, 0.08, 0.2), Vector3(0.06, 0.06, 0.03), C(0x101010), Vector3.ZERO, 0.0)
	hb.box(Vector3(0.12, 0.08, 0.2), Vector3(0.06, 0.06, 0.03), C(0x101010), Vector3.ZERO, 0.0)
	_part(mdl, "head", Vector3(0, 1.45, 0.85), hb)
	for side in [-1, 1]:
		var wb := MeshBuilder.new()
		wb.box(Vector3(side * 0.7, 0, 0), Vector3(1.4, 0.06, 0.7), wing_c)
		wb.box(Vector3(side * 1.3, 0, -0.2), Vector3(0.7, 0.05, 0.5), wing_c.darkened(0.1))
		mdl.parts["wing_l" if side < 0 else "wing_r"] = _part(mdl, "wing_l" if side < 0 else "wing_r", Vector3(side * 0.3, 1.2, 0), wb)
	for leg in [[-1, 1], [1, 1], [-1, -1], [1, -1]]:
		var lb := MeshBuilder.new()
		lb.box(Vector3(0, -0.35, 0), Vector3(0.18, 0.7, 0.18), c.darkened(0.1))
		_part(mdl, "leg", Vector3(leg[0] * 0.25, 0.7, leg[1] * 0.5), lb)
	var tb := MeshBuilder.new()
	tb.box(Vector3(0, 0, -0.4), Vector3(0.14, 0.14, 0.8), c)
	tb.box(Vector3(0, 0, -0.85), Vector3(0.24, 0.24, 0.2), c.darkened(0.2))
	mdl.parts["tail"] = _part(mdl, "tail", Vector3(0, 1.0, -0.75), tb)
	mdl.parts["body"] = body
	mdl.height = 2.4
	mdl.set_meta("seat_y", 1.55)
	mdl.set_meta("seat_z", -0.05)
	return mdl

# ---------------------------------------------------------------- world objects

static func make_object(m: String) -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = "object"
	var mb := MeshBuilder.new()
	var glow := MeshBuilder.new()
	match m:
		"pumpkin":
			mb.box(Vector3(0, 0.22, 0), Vector3(0.5, 0.44, 0.5), C(0xe08020), Vector3.ZERO, 0.08)
			mb.box(Vector3(0, 0.48, 0), Vector3(0.08, 0.14, 0.08), C(0x4a6a20))
			mdl.height = 0.7
		"wheat":
			var r := RandomNumberGenerator.new()
			r.seed = 7
			for i in 12:
				mb.box(Vector3(r.randf_range(-0.25, 0.25), 0.35, r.randf_range(-0.25, 0.25)), Vector3(0.05, 0.7, 0.05), C(0xe8c040), Vector3(r.randf_range(-10, 10), 0, r.randf_range(-10, 10)))
			mb.box(Vector3(0, 0.3, 0), Vector3(0.5, 0.05, 0.05), C(0xa07a20))
			mdl.height = 0.9
		"spring":
			mb.box(Vector3(0, 0.1, 0), Vector3(1.4, 0.2, 1.4), C(0x8a8a86))
			glow.box(Vector3(0, 0.22, 0), Vector3(1.1, 0.06, 1.1), C(0x60c8ff, 0.8))
			for i in 6:
				var a := i * TAU / 6.0
				mb.box(Vector3(cos(a) * 0.8, 0.25, sin(a) * 0.8), Vector3(0.3, 0.3, 0.3), C(0x9a9a94), Vector3(0, i * 20, 0))
			mdl.height = 0.6
		"tracks":
			for i in 5:
				mb.box(Vector3(-0.2 + (i % 2) * 0.4, 0.02, -0.7 + i * 0.35), Vector3(0.22, 0.03, 0.28), C(0x3a2a1c), Vector3(0, 15 * i, 0), 0.0)
			mdl.height = 0.3
		"crystal":
			for i in 3:
				glow.box(Vector3(-0.2 + i * 0.2, 0.4 + i * 0.1, 0.0), Vector3(0.16, 0.8 + i * 0.2, 0.16), C(0x60b0ff, 0.9), Vector3(0, 0, -15 + i * 15), 0.0)
			mb.box(Vector3(0, 0.1, 0), Vector3(0.7, 0.2, 0.5), C(0x5a5a62))
			mdl.height = 1.2
		"lotus":
			mb.box(Vector3(0, 0.03, 0), Vector3(0.8, 0.04, 0.8), C(0x3a8a3a))
			for i in 8:
				var a := i * TAU / 8.0
				mb.box(Vector3(cos(a) * 0.14, 0.12, sin(a) * 0.14), Vector3(0.1, 0.2, 0.04), C(0xf080b0), Vector3(0, -a * 57.3, 30))
			glow.box(Vector3(0, 0.12, 0), Vector3(0.08, 0.08, 0.08), C(0xffe060))
			mdl.height = 0.5
		"cactus":
			mb.box(Vector3(0, 0.5, 0), Vector3(0.3, 1.0, 0.3), C(0x3a8a3a))
			mb.box(Vector3(0.28, 0.65, 0), Vector3(0.3, 0.14, 0.14), C(0x3a8a3a))
			mb.box(Vector3(0.4, 0.8, 0), Vector3(0.14, 0.3, 0.14), C(0x3a8a3a))
			mb.box(Vector3(0, 1.02, 0), Vector3(0.12, 0.08, 0.12), C(0xe05080))
			mdl.height = 1.1
		"urn":
			mb.frustum(Vector3(0, 0, 0), 0.3, 0.3, 0.5, 0.5, 0.35, C(0xb87a3a))
			mb.box(Vector3(0, 0.5, 0), Vector3(0.52, 0.3, 0.52), C(0xc88a46))
			mb.frustum(Vector3(0, 0.65, 0), 0.5, 0.5, 0.3, 0.3, 0.15, C(0xb87a3a))
			mb.box(Vector3(0, 0.5, 0.265), Vector3(0.3, 0.06, 0.02), C(0x4a2a14))
			mdl.height = 0.9
		"portal":
			for side in [-1, 1]:
				mb.box(Vector3(side * 1.0, 1.4, 0), Vector3(0.4, 2.8, 0.5), C(0x4a4a56))
			mb.box(Vector3(0, 2.9, 0), Vector3(2.4, 0.4, 0.5), C(0x5a5a68))
			glow.box(Vector3(0, 1.4, 0), Vector3(1.6, 2.6, 0.06), C(0x9a60ff, 0.7), Vector3.ZERO, 0.0)
			glow.box(Vector3(0, 1.4, 0), Vector3(1.0, 1.9, 0.08), C(0xd0a0ff, 0.8), Vector3.ZERO, 0.0)
			mdl.height = 3.2
		_:
			mb.box(Vector3(0, 0.3, 0), Vector3(0.6, 0.6, 0.6), C(0xa0a040))
			mdl.height = 0.8
	mdl.add_child(mb.instance())
	if not glow.verts.is_empty():
		var gi := glow.instance(MeshBuilder.alpha_material() if m in ["spring", "crystal", "portal"] else MeshBuilder.glow_material())
		mdl.add_child(gi)
		mdl.glow_nodes.append(gi)
		if m == "portal":
			mdl.parts["swirl"] = gi
	return mdl

static func make_ship(m: String) -> BlockyModel:
	var mdl := BlockyModel.new()
	mdl.kind = "ship"
	var hull := C(0x6a4a2a)
	var mb := MeshBuilder.new()
	# hull: 3 m wide, 12 m long; tapered bow
	mb.box(Vector3(0, 0.6, 0), Vector3(3.0, 1.3, 9.0), hull)
	mb.frustum(Vector3(0, 0, 4.5), 3.0, 0.1, 3.0, 0.1, 0.0001, hull)
	mb.frustum(Vector3(0, 0.0, 4.5), 3.0, 0.1, 0.4, 0.1, 1.3, hull, Vector3.ZERO, 0.04, Vector2.ZERO)
	mb.box(Vector3(0, 1.3, -5.0), Vector3(2.6, 1.8, 2.0), hull.darkened(0.1))
	mb.box(Vector3(0, 1.3, 0), Vector3(2.8, 0.1, 8.6), C(0x9a7a4a))
	for side in [-1, 1]:
		mb.box(Vector3(side * 1.45, 1.6, 0), Vector3(0.12, 0.6, 8.8), hull.darkened(0.15))
	mb.box(Vector3(0, 4.0, 0.5), Vector3(0.3, 6.0, 0.3), C(0x4a2a14))
	mb.box(Vector3(0, 6.3, 0.5), Vector3(4.2, 0.18, 0.18), C(0x4a2a14))
	var sb := MeshBuilder.new()
	var sail := C(0xece4d0) if m != "ship_junk" else C(0xc84a2a)
	sb.box(Vector3(0, 0, 0), Vector3(3.8, 3.8, 0.1), sail, Vector3.ZERO, 0.02)
	if m == "ship_junk":
		for i in 3:
			sb.box(Vector3(0, -1.2 + i * 1.2, 0.07), Vector3(3.8, 0.06, 0.04), C(0x4a2a14), Vector3.ZERO, 0.0)
	mdl.add_child(mb.instance())
	var sailn := _part(mdl, "sail", Vector3(0, 4.3, 0.55), sb)
	mdl.parts["sail"] = sailn
	mdl.height = 7.0
	return mdl
