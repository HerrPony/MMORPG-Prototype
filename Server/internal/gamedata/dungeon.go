package gamedata

// Dungeons are instanced areas far outside the continent (z > DungeonZ0).

type Rect struct {
	X     float64 `json:"x"` // center
	Z     float64 `json:"z"`
	W     float64 `json:"w"`
	D     float64 `json:"d"`
	Theme string  `json:"theme"` // hall | mine | web | throne | corridor
}

type DungeonSpawn struct {
	Mob  string
	X, Z float64
}

type Dungeon struct {
	ID       string
	Name     string
	MinLevel int
	Origin   [2]float64 // world position of the dungeon's local (0,0)
	Entrance [2]float64 // world position of the entrance portal in the open world
	Rooms    []Rect     // local coordinates
	Spawns   []DungeonSpawn
	Start    [2]float64 // local arrival position
}

var Dungeons = map[string]*Dungeon{
	"mine": {
		ID: "mine", Name: "Die Schwarzfelsmine", MinLevel: 15, Origin: [2]float64{0, DungeonZ0}, Entrance: [2]float64{-175, -198},
		Start: [2]float64{0, 3},
		Rooms: []Rect{
			{0, 0, 26, 20, "hall"}, {0, -20, 6, 20, "corridor"}, {0, -42, 32, 24, "mine"}, {0, -60, 6, 12, "corridor"},
			{0, -84, 34, 36, "mine"}, {0, -108, 6, 12, "corridor"}, {0, -132, 40, 36, "mine"}, {0, -156, 6, 12, "corridor"},
			{0, -182, 38, 40, "web"}, {0, -208, 6, 12, "corridor"}, {0, -238, 46, 48, "throne"},
		},
		Spawns: []DungeonSpawn{
			// first cave
			{"kobold", -9, -38}, {"kobold", 0, -35}, {"kobold", 9, -38}, {"kobold_shaman", -4, -48}, {"kobold_shaman", 5, -49},
			// Vorarbeiter Grimm
			{"foreman", 0, -88}, {"kobold", -12, -74}, {"kobold", 12, -74},
			// second cave
			{"skeleton_miner", -11, -124}, {"skeleton_miner", 11, -124}, {"kobold", -12, -140}, {"kobold", 12, -140}, {"kobold_shaman", 0, -133}, {"skeleton_miner", 0, -144},
			// Spinnenkönigin Vessa
			{"cave_spider", -13, -170}, {"cave_spider", 13, -170}, {"cave_spider", 0, -176}, {"spider_queen", 0, -192},
			// Kobold-König Skarn
			{"kobold_king", 0, -248},
		},
	},
}

// Walkable reports whether a local dungeon position lies inside the rooms (with a small margin).
func (d *Dungeon) Walkable(x, z float64) bool {
	lx, lz := x-d.Origin[0], z-d.Origin[1]
	for _, r := range d.Rooms {
		if lx >= r.X-r.W/2+0.5 && lx <= r.X+r.W/2-0.5 && lz >= r.Z-r.D/2-0.5 && lz <= r.Z+r.D/2+0.5 {
			return true
		}
	}
	return false
}

// Clamp pushes a world position back into the walkable area.
func (d *Dungeon) Clamp(x, z float64) (float64, float64) {
	if d.Walkable(x, z) {
		return x, z
	}
	lx, lz := x-d.Origin[0], z-d.Origin[1]
	bx, bz, best := lx, lz, 1e18
	for _, r := range d.Rooms {
		cx := max(r.X-r.W/2+0.6, min(r.X+r.W/2-0.6, lx))
		cz := max(r.Z-r.D/2-0.4, min(r.Z+r.D/2+0.4, lz))
		if dd := (cx-lx)*(cx-lx) + (cz-lz)*(cz-lz); dd < best {
			best, bx, bz = dd, cx, cz
		}
	}
	return bx + d.Origin[0], bz + d.Origin[1]
}
