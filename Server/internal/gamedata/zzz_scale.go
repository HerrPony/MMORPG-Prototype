package gamedata

import (
	"math"
	"sync"
)

// The data tables of the world (towns, places, NPCs, roads, rivers ...) were written for continents of
// 600 x 600 m. The continents are WorldScale times larger now: scaleWorld moves every position, so the places
// lie further apart. Things that belong to a place (the buildings of a camp, the NPCs of a village, the mobs
// around a camp) keep their offset to it, so a village is as big as before - only the distance between the
// villages grows. This file sorts after all other data files, so its init runs last.

var (
	scaleOnce sync.Once
	// placeAnchors are the original and the new centre of every town and place
	placeAnchors []placeAnchor
)

type placeAnchor struct {
	ox, oz, nx, nz, r float64
	town              bool
	kind              string
}

// safeRadius is how far hostile places keep away from a settlement (the first village is the starting area).
func safeRadius(t *Town) float64 {
	switch {
	case t.ID == "eichenhain":
		return 130
	case t.Kind == "capital":
		return 120
	case t.Kind == "village":
		return 80
	}
	return 0
}

// keepAway moves a point radially away from the settlements until it is outside of their safe radius.
func keepAway(x, z float64) (float64, float64) {
	for _, t := range Towns {
		r := safeRadius(t)
		if d := math.Hypot(x-t.X, z-t.Z); d < r && r > 0 {
			if d < 1 {
				d, x = 1, x+1
			}
			// the points that are pushed spread over a band, not a single ring
			band := math.Mod(math.Abs(x*7.31+z*13.17), 1) * 45
			x, z = t.X+(x-t.X)/math.Max(d, 1)*(r+band), t.Z+(z-t.Z)/math.Max(d, 1)*(r+band)
		}
	}
	return x, z
}

func init() { scaleOnce.Do(scaleWorld) }

// moveToWorld maps an original position to the new world: next to a place it keeps the offset to the place,
// elsewhere it is stretched.
func moveToWorld(x, z float64) (float64, float64) {
	best, bd := -1, 1e18
	for i, a := range placeAnchors {
		if d := math.Hypot(x-a.ox, z-a.oz) / a.r; d < 1 && d < bd {
			best, bd = i, d
		}
	}
	if best >= 0 {
		a := placeAnchors[best]
		return a.nx + (x - a.ox), a.nz + (z - a.oz)
	}
	return x * WorldScale, z * WorldScale
}

func scaleWorld() {
	if worldGenerated {
		panic("gamedata: the terrain was generated before the data tables were scaled")
	}
	const S = WorldScale
	for _, t := range Towns {
		placeAnchors = append(placeAnchors, placeAnchor{t.X, t.Z, t.X * S, t.Z * S, 80, true, ""})
	}
	poiAnchor := map[*POI]int{}
	for _, p := range POIs {
		if !p.Town {
			poiAnchor[p] = len(placeAnchors)
			placeAnchors = append(placeAnchors, placeAnchor{p.X, p.Z, p.X * S, p.Z * S, 45, false, p.Kind})
		}
	}
	// hostile places (camps, wolf hills ...) keep their distance to the settlements: they move outwards as a whole,
	// with everything that belongs to them
	for i := range placeAnchors {
		a := &placeAnchors[i]
		if !a.town && a.kind != "lake" { // lakes stay where the water is
			for _, t := range Towns {
				if r := safeRadius(t); r > 0 {
					tx, tz := t.X*S, t.Z*S
					if d := math.Hypot(a.nx-tx, a.nz-tz); d < r {
						a.nx, a.nz = tx+(a.nx-tx)/d*r, tz+(a.nz-tz)/d*r
					}
				}
			}
		}
	}
	for _, n := range NPCs {
		n.X, n.Z = moveToWorld(n.X, n.Z)
	}
	for i := range Spawns {
		Spawns[i].X, Spawns[i].Z = moveToWorld(Spawns[i].X, Spawns[i].Z)
	}
	for i := range extraObjectSpawns {
		extraObjectSpawns[i].X, extraObjectSpawns[i].Z = moveToWorld(extraObjectSpawns[i].X, extraObjectSpawns[i].Z)
	}
	for i := range objGroups {
		objGroups[i].x, objGroups[i].z = moveToWorld(objGroups[i].x, objGroups[i].z)
	}
	for _, d := range Dungeons {
		d.Entrance[0], d.Entrance[1] = moveToWorld(d.Entrance[0], d.Entrance[1])
	}
	for _, p := range POIs {
		if i, ok := poiAnchor[p]; ok {
			p.X, p.Z = placeAnchors[i].nx, placeAnchors[i].nz
			if p.Kind == "wild" {
				p.R *= S
			}
		} else {
			p.X, p.Z = p.X*S, p.Z*S
		}
	}
	for i := range anchors {
		anchors[i].x, anchors[i].z, anchors[i].s = anchors[i].x*S, anchors[i].z*S, anchors[i].s*S
	}
	for i := range Lakes {
		Lakes[i].X, Lakes[i].Z, Lakes[i].R = Lakes[i].X*S, Lakes[i].Z*S, Lakes[i].R*1.2
	}
	for i := range Rivers {
		for j := range Rivers[i].Pts {
			Rivers[i].Pts[j][0] *= S
			Rivers[i].Pts[j][1] *= S
		}
	}
	for i := range Roads {
		pts := Roads[i].Pts
		for j := range pts {
			pts[j][0] *= S
			pts[j][1] *= S
		}
		// a road that ends at a place that moved leads to its new position
		for _, end := range []int{0, len(pts) - 1} {
			for _, a := range placeAnchors {
				if !a.town && math.Hypot(pts[end][0]-a.ox*S, pts[end][1]-a.oz*S) < 40*S {
					pts[end] = [2]float64{a.nx, a.nz}
				}
			}
		}
	}
	portalMineX, portalMineZ = moveToWorld(-175, -198)
	for _, t := range Towns {
		t.X, t.Z = t.X*S, t.Z*S
	}
}

// the portal of the mine in the open world (the original position was -175, -198)
var portalMineX, portalMineZ float64

// landOK reports whether a point is dry, walkable ground.
func landOK(x, z float64) bool {
	return Height(x, z) >= 0.9 && Walkable(x, z) && BiomeAt(x, z) != BiomeLava
}

// snapLand returns the nearest dry ground to a point (the point itself if it is fine already).
func snapLand(x, z float64) (float64, float64) {
	if landOK(x, z) {
		return x, z
	}
	for r := 2.0; r <= 80; r += 2 {
		for a := 0.0; a < 2*math.Pi; a += math.Pi / 8 {
			nx, nz := x+math.Cos(a)*r, z+math.Sin(a)*r
			if landOK(nx, nz) {
				return nx, nz
			}
		}
	}
	return x, z
}

// snapEntities moves NPCs, places and portals that landed in water or lava (the terrain is made from noise that
// does not move with the places) onto the nearest dry ground.
func snapEntities() {
	for _, n := range NPCs {
		n.X, n.Z = snapLand(n.X, n.Z)
	}
	for _, p := range POIs {
		if p.Kind != "lake" {
			p.X, p.Z = snapLand(p.X, p.Z)
		}
	}
	for _, d := range Dungeons {
		d.Entrance[0], d.Entrance[1] = snapLand(d.Entrance[0], d.Entrance[1])
	}
	portalMineX, portalMineZ = snapLand(portalMineX, portalMineZ)
	for i := range extraObjectSpawns {
		extraObjectSpawns[i].X, extraObjectSpawns[i].Z = snapLand(extraObjectSpawns[i].X, extraObjectSpawns[i].Z)
	}
}

// separateSpawns moves spawns that ended up on top of each other (after being pushed away from a settlement)
// apart, onto free dry ground.
func separateSpawns(sp []Spawn) {
	tooClose := func(i int, x, z float64) bool {
		for j := range sp {
			if j != i && math.Hypot(sp[j].X-x, sp[j].Z-z) < 4.5 {
				return true
			}
		}
		return false
	}
	for i := range sp {
		if !tooClose(i, sp[i].X, sp[i].Z) {
			continue
		}
	search:
		for r := 4.0; r <= 40; r += 2 {
			for a := 0.0; a < 2*math.Pi; a += math.Pi / 6 {
				x, z := sp[i].X+math.Cos(a+r)*r, sp[i].Z+math.Sin(a+r)*r
				if landOK(x, z) && DistToRoad(x, z) >= roadClear && !tooClose(i, x, z) {
					sp[i].X, sp[i].Z = x, z
					break search
				}
			}
		}
	}
}
