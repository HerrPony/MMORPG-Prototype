package gamedata

import (
	"math"
	"sort"
)

// QArea is a circular map area where a quest objective can be completed.
type QArea struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
	K string  `json:"k"` // mob | obj | npc | poi | dungeon
}

type pt struct{ x, z float64 }

// cluster groups points that lie close together and returns one circle per group (largest groups first).
func cluster(pts []pt, kind string, maxN int) []QArea {
	type cl struct {
		sx, sz float64
		n      int
		pts    []pt
	}
	var cs []*cl
	for _, p := range pts {
		var hit *cl
		for _, c := range cs {
			if math.Hypot(c.sx/float64(c.n)-p.x, c.sz/float64(c.n)-p.z) < 45 {
				hit = c
				break
			}
		}
		if hit == nil {
			hit = &cl{}
			cs = append(cs, hit)
		}
		hit.sx += p.x
		hit.sz += p.z
		hit.n++
		hit.pts = append(hit.pts, p)
	}
	sort.SliceStable(cs, func(i, j int) bool { return cs[i].n > cs[j].n })
	var out []QArea
	for i, c := range cs {
		if i >= maxN {
			break
		}
		cx, cz := c.sx/float64(c.n), c.sz/float64(c.n)
		r := 0.0
		for _, p := range c.pts {
			r = math.Max(r, math.Hypot(p.x-cx, p.z-cz))
		}
		out = append(out, QArea{X: math.Round(cx), Z: math.Round(cz), R: math.Round(math.Max(14, r+10)), K: kind})
	}
	return out
}

func mobAreas(mob string) []QArea {
	var pts []pt
	for _, s := range spawnsAll {
		if s.Mob == mob {
			pts = append(pts, pt{s.X, s.Z})
		}
	}
	out := cluster(pts, "mob", 3)
	for _, d := range Dungeons {
		for _, s := range d.Spawns {
			if s.Mob == mob {
				out = append(out, QArea{X: d.Entrance[0], Z: d.Entrance[1], R: 14, K: "dungeon"})
				break
			}
		}
	}
	return out
}

func objectAreas(item string) []QArea {
	var pts []pt
	for _, s := range ObjectSpawns {
		if def := Objects[s.Obj]; def != nil && def.Item == item {
			pts = append(pts, pt{s.X, s.Z})
		}
	}
	return cluster(pts, "obj", 2)
}

func npcArea(id string) *QArea {
	if n := NPCByID(id); n != nil {
		return &QArea{X: n.X, Z: n.Z, R: 10, K: "npc"}
	}
	return nil
}

// ObjectiveAreas returns the map areas of one quest objective.
func ObjectiveAreas(o Objective) []QArea {
	switch o.Kind {
	case "kill", "tame":
		return mobAreas(o.Target)
	case "collect":
		out := objectAreas(o.Target)
		if len(out) > 0 {
			return out
		}
		for _, m := range Mobs {
			for _, l := range m.Loot {
				if l.Item == o.Target {
					out = append(out, mobAreas(m.ID)...)
					break
				}
			}
		}
		return out
	case "deliver", "talk":
		if a := npcArea(o.Target); a != nil {
			return []QArea{*a}
		}
	case "explore":
		if p := POIByID(o.Target); p != nil {
			return []QArea{{X: p.X, Z: p.Z, R: p.R, K: "poi"}}
		}
	}
	return nil
}

func computeQuestAreas() {
	for _, q := range Quests {
		q.Areas = nil
		for _, o := range q.Objectives {
			q.Areas = append(q.Areas, ObjectiveAreas(o))
		}
		hand := q.Hand
		if hand == "" {
			hand = q.Giver
		}
		q.HandIn = npcArea(hand)
	}
}
