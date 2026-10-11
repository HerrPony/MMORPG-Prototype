package gamedata

import (
	"hash/fnv"
	"math"
	"math/rand"
	"strings"
	"sync"
)

// ObjectDef is an interactive world object (herbs, chests, portals).
type ObjectDef struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Model   string  `json:"model"`
	Item    string  `json:"item,omitempty"`   // item received when used
	Quest   string  `json:"quest,omitempty"`  // quest that must be active for the object to be usable
	Respawn float64 `json:"-"`                // seconds until it returns
	Portal  string  `json:"portal,omitempty"` // "dungeon:<id>" or "exit"
}

var Objects = map[string]*ObjectDef{
	"pumpkin_patch": {ID: "pumpkin_patch", Name: "Kürbis", Model: "pumpkin", Item: "pumpkin", Quest: "e_pumpkin", Respawn: 45},
	"wheat_sheaf":   {ID: "wheat_sheaf", Name: "Weizenbündel", Model: "wheat", Item: "wheat", Quest: "p_wheat", Respawn: 45},
	"spring_pool":   {ID: "spring_pool", Name: "Heilige Quelle", Model: "spring", Item: "spring_water", Quest: "p_holy", Respawn: 45},
	"track_marks":   {ID: "track_marks", Name: "Wildfährte", Model: "tracks", Item: "track", Quest: "h_tracks", Respawn: 45},
	"crystal_node":  {ID: "crystal_node", Name: "Wasserkristall", Model: "crystal", Item: "crystal_water", Quest: "m_crystals", Respawn: 45},
	"lotus_flower":  {ID: "lotus_flower", Name: "Moorlotus", Model: "lotus", Item: "lotus", Quest: "s_herbs", Respawn: 45},
	"cactus_fruit":  {ID: "cactus_fruit", Name: "Kaktusfrucht", Model: "cactus", Item: "cactus_fruit", Quest: "g_cactus", Respawn: 45},
	"ancient_urn":   {ID: "ancient_urn", Name: "Antike Scherbe", Model: "urn", Item: "ancient_shard", Quest: "g_shards", Respawn: 45},
	"portal_mine":   {ID: "portal_mine", Name: "Portal zur Schwarzfelsmine", Model: "portal", Portal: "dungeon:mine"},
	"portal_exit":   {ID: "portal_exit", Name: "Ausgang", Model: "portal", Portal: "exit"},
}

type ObjectSpawn struct {
	Obj  string
	X, Z float64
}

var ObjectSpawns []ObjectSpawn

type objGroup struct {
	obj    string
	x, z   float64
	n      int
	radius float64
}

var objGroups = []objGroup{
	{"pumpkin_patch", -24, 18, 10, 10},
	{"wheat_sheaf", 140, -62, 12, 14},
	{"track_marks", 130, -212, 6, 14},
	{"spring_pool", 178, -228, 8, 9},
	{"crystal_node", 178, -228, 6, 12},
	{"lotus_flower", -185, 150, 12, 24},
	{"cactus_fruit", 175, 150, 9, 22},
	{"ancient_urn", 140, 238, 7, 10},
}

var (
	initOnce  sync.Once
	spawnsAll []Spawn
)

// Init finishes the data set: places spawn groups, objects and vendor stock.
// It must run before the world is created and is safe to call repeatedly.
func Init() {
	initOnce.Do(func() {
		buildSpawns()
		spawnsAll = append([]Spawn(nil), Spawns...)
		for i := range spawnsAll { // the hand placed spawns also keep off the roads
			spawnsAll[i].X, spawnsAll[i].Z = offRoad(spawnsAll[i].X, spawnsAll[i].Z)
		}
		spawnsAll = append(spawnsAll, generatedSpawns...)
		for _, g := range objGroups {
			h := fnv.New32a()
			h.Write([]byte(g.obj))
			rng := rand.New(rand.NewSource(int64(h.Sum32()) + int64(g.x*3)))
			placed := 0
			for try := 0; try < 500 && placed < g.n; try++ {
				a := rng.Float64() * 2 * math.Pi
				r := math.Sqrt(rng.Float64()) * g.radius
				x, z := g.x+math.Cos(a)*r, g.z+math.Sin(a)*r
				if g.obj == "lotus_flower" {
					if h := Height(x, z); h < -0.9 || h > 1.0 { // shallow water
						continue
					}
				} else if Height(x, z) < 0.8 {
					continue
				}
				if nearNPC(x, z, 4) || !Walkable(x, z) {
					continue
				}
				ObjectSpawns = append(ObjectSpawns, ObjectSpawn{g.obj, x, z})
				placed++
			}
		}
		ObjectSpawns = append(ObjectSpawns, ObjectSpawn{"portal_mine", -175, -198})
		ObjectSpawns = append(ObjectSpawns, extraObjectSpawns...)
		computeQuestAreas()
		// vendor stock
		for _, n := range NPCs {
			if !n.Vendor || n.StockK == "" {
				continue
			}
			for _, k := range strings.Split(n.StockK, ",") {
				n.Stock = append(n.Stock, stockFor(k, n.StockLo, n.StockHi)...)
			}
		}
	})
}

// AllSpawns lists every open world spawn.
func AllSpawns() []Spawn {
	Init()
	return spawnsAll
}
