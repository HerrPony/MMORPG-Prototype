class_name Ent
extends RefCounted
## Client-side record of one entity (player, mob, npc, pet, object, ship) as sent by the server, plus its scene nodes.

var id := 0
var k := ""
var n := ""
var ti := ""
var m := ""
var c := ""
var l := 1
var hp := 1
var mh := 1
var d := false
var tg := 0
var st: Array = []
var lk: Dictionary = {}
var el := false
var bs := false
var sc := 1.0
var ow := 0
var fl := false
var q := ""
var gn := false
var eq: Dictionary = {}
var pt := false
var lt := false          # lootable corpse (from the snapshot's lt list)
# target (server) position, rendered (smoothed) position
var x := 0.0
var y := 0.0
var z := 0.0
var r := 0.0
var rx := 0.0
var ry := 0.0
var rz := 0.0
var rr := 0.0
var speed := 0.0
var airborne := false
var inited := false
var eq_key := ""
var mark_state := "-"
var obj_glow := false
# scene nodes
var node: Node3D
var model: BlockyModel
var mount: BlockyModel
var pad: BlockyModel
var marker: Label3D
var sparkle: CPUParticles3D
var plate: Control
var plate_key := ""
var was_target := false
var anim_t := 0.0

func height() -> float:
	return (model.height if model else 1.8) * sc
