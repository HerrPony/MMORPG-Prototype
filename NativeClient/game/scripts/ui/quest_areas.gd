class_name QAreas
extends RefCounted
## Map areas of the active quests: where to go to finish the open objectives, or where to hand the quest in.

## Returns [{x, z, r, k, qid, title, text, turn_in}]
static func get_all() -> Array:
	var out: Array = []
	for q in S.quests:
		if q.status != "active" and q.status != "complete":
			continue
		var def := S.quest_def(q.id)
		if def.is_empty():
			continue
		if q.status == "complete":
			var h = def.get("handIn")
			if h is Dictionary:
				var a: Dictionary = (h as Dictionary).duplicate()
				a["qid"] = q.id; a["title"] = def.title; a["text"] = "Quest abgeben"; a["turn_in"] = true
				out.append(a)
			continue
		var objs: Array = def.get("objectives", [])
		var areas: Array = def.get("areas", [])
		for i in objs.size():
			var o: Dictionary = objs[i]
			var counts: Array = q.get("counts", [])
			if (int(counts[i]) if i < counts.size() else 0) >= int(o.count):
				continue
			if i < areas.size() and areas[i] is Array:
				for a0 in areas[i]:
					var a: Dictionary = (a0 as Dictionary).duplicate()
					a["qid"] = q.id; a["title"] = def.title; a["text"] = o.text; a["turn_in"] = false
					out.append(a)
	return out
