package gamedata

import (
	"encoding/base64"
	"math"
	"sync"
)

// World props: trees, rocks, bushes, grass ... scattered over the terrain. This is the single place that
// decides what stands where. The list is sent to the clients (rendering + their own collision) and used by
// the server for mob movement and line of sight (collide.go).

// PropKinds are the model names the client knows (Client/js/world/vegetation.js); a prop refers to one by index.
var PropKinds = []string{"pine", "oak", "dead", "snowpine", "cactus", "palm", "rock", "bush", "reed", "mushroom", "tuft",
	"flower", "bones", "jtree", "fern", "bamboo", "sakura", "mesapillar", "drybush", "charred", "obsidian", "ember"}

type Prop struct {
	K     int     // index into PropKinds
	X, Z  float64 // world position (quantised to 1/40 m)
	S     float64 // scale (1/100)
	Ry    float64 // rotation around y (radians, 8 bit)
	Solid float64 // collision radius per unit scale, 0 = walk through
}

var (
	propsOnce sync.Once
	propList  []Prop
)

// worldProps returns the (cached) prop list for a base scenery.
func worldProps(sc Scenery) []Prop {
	propsOnce.Do(func() { propList = applyPropEdits(generateProps(sc)) })
	return propList
}

const (
	propPosQ = 40.0  // position steps per metre
	propRyQ  = 255.0 // rotation steps per 2π (6.28)
)

func kindIndex(name string) int {
	for i, k := range PropKinds {
		if k == name {
			return i
		}
	}
	panic("unknown prop kind " + name)
}

// quantised grid access as the clients see it (int16 of h*20, stored as float32)
func qHeight(i, j int) float64 {
	i = max(0, min(GridNX-1, i))
	j = max(0, min(GridNZ-1, j))
	return float64(float32(float64(int16(math.Round(float64(hm[j*GridNX+i])*20))) / 20))
}

func qHeightAt(x, z float64) float64 {
	gx, gz := (x-WorldMinX)/GridCell, (z-WorldMinZ)/GridCell
	i, j := int(math.Floor(gx)), int(math.Floor(gz))
	fx, fz := gx-float64(i), gz-float64(j)
	a, b, c, d := qHeight(i, j), qHeight(i+1, j), qHeight(i, j+1), qHeight(i+1, j+1)
	return a + (b-a)*fx + (c-a)*fz + (a-b-c+d)*fx*fz
}

// cellSeed makes the random stream of one 3x3 m cell independent of all others, so a change in one place
// (terrain edit, moved spawn) does not reshuffle the vegetation of the whole world.
func cellSeed(i, j int) uint32 {
	h := uint32(i)*73856093 ^ uint32(j)*19349663 ^ 2024
	h ^= h >> 16
	h *= 0x7feb352d
	h ^= h >> 15
	h *= 0x846ca68b
	h ^= h >> 16
	return h
}

func generateProps(sc Scenery) []Prop {
	var seed uint32
	r := func() float64 { seed = seed*1664525 + 1013904223; return float64(seed) / 4294967296 }
	const step = 3.0
	blocked := func(x, z float64) bool {
		for _, c := range sc.Clears {
			if math.Hypot(x-c.X, z-c.Z) < c.R+1.2 {
				return true
			}
		}
		for _, rd := range sc.Roads {
			if distToPoly(x, z, rd) < 3.2 {
				return true
			}
		}
		for _, b := range sc.Buildings {
			if b.Kind != "wall" && math.Abs(x-b.X) < b.W/2+2 && math.Abs(z-b.Z) < b.D/2+2 {
				return true
			}
		}
		for _, b := range sc.House {
			if math.Abs(x-b.X) < b.W/2+3 && math.Abs(z-b.Z) < b.D/2+3 {
				return true
			}
		}
		return false
	}
	var out []Prop
	add := func(kind string, x, z, scale, solid float64) {
		ry := r() * 6.28
		qx, qz := math.Round((x-WorldMinX)*propPosQ), math.Round((z-WorldMinZ)*propPosQ)
		out = append(out, Prop{
			K: kindIndex(kind), X: WorldMinX + qx/propPosQ, Z: WorldMinZ + qz/propPosQ,
			S:  math.Min(2.55, math.Round(scale*100)/100),
			Ry: math.Round(ry/6.28*propRyQ) / propRyQ * 6.28, Solid: math.Round(solid*100) / 100,
		})
	}
	xEnd := WorldMinX + float64(GridNX-1)*GridCell - 4
	zEnd := WorldMinZ + float64(GridNZ-1)*GridCell - 4
	for j := 0; WorldMinZ+4+float64(j)*step < zEnd; j++ {
		z := WorldMinZ + 4 + float64(j)*step
		for i := 0; WorldMinX+4+float64(i)*step < xEnd; i++ {
			x := WorldMinX + 4 + float64(i)*step
			seed = cellSeed(i, j)
			px, pz := x+(r()-0.5)*step, z+(r()-0.5)*step
			b := int(qBiome(px, pz))
			if b == 0 || b == 10 {
				continue
			}
			h := qHeightAt(px, pz)
			if h < 0.35 {
				if b == 3 && h > -1.0 && h < 0.35 && r() < 0.12 && !blocked(px, pz) {
					add("reed", px, pz, 0.9+r()*0.8, 0)
				}
				continue
			}
			if r() > 0.5 { // quick reject: most cells stay empty
				continue
			}
			if blocked(px, pz) {
				continue
			}
			slope := math.Abs(qHeightAt(px+2, pz)-h) + math.Abs(qHeightAt(px, pz+2)-h)
			s := 0.8 + r()*0.9
			u := r()
			switch b {
			case 1: // forest
				if u < 0.55 {
					k := "oak"
					if r() < 0.72 {
						k = "pine"
					}
					add(k, px, pz, s, 0.5)
				} else if u < 0.62 {
					add("bush", px, pz, s, 0.4)
				} else if u < 0.66 {
					add("rock", px, pz, s*0.8, 0.6)
				} else if u < 0.90 {
					add("tuft", px, pz, s, 0)
				} else {
					add("flower", px, pz, s, 0)
				}
			case 2: // plains
				if u < 0.08 {
					add("oak", px, pz, s, 0.5)
				} else if u < 0.12 {
					add("bush", px, pz, s, 0.4)
				} else if u < 0.40 {
					add("tuft", px, pz, s, 0)
				} else if u < 0.55 {
					add("flower", px, pz, s, 0)
				}
			case 3: // swamp
				if u < 0.22 {
					add("dead", px, pz, s, 0.4)
				} else if u < 0.32 {
					add("mushroom", px, pz, s, 0)
				} else if u < 0.38 {
					add("bush", px, pz, s, 0.4)
				} else if u < 0.55 {
					add("tuft", px, pz, s, 0)
				}
			case 4: // snow
				if h > 24 || slope > 3.2 {
					if u < 0.25 {
						add("rock", px, pz, s*1.4, 0.9)
					}
				} else if u < 0.30 {
					add("snowpine", px, pz, s, 0.5)
				} else if u < 0.34 {
					add("rock", px, pz, s*1.1, 0.8)
				}
			case 5: // desert
				if u < 0.06 {
					add("cactus", px, pz, s, 0.5)
				} else if u < 0.085 {
					add("rock", px, pz, s, 0.7)
				} else if u < 0.10 && h < 1.8 {
					add("palm", px, pz, s, 0.4)
				} else if u < 0.105 {
					add("bones", px, pz, s, 0)
				} else if u < 0.15 {
					add("bush", px, pz, s*0.7, 0)
				}
			case 6: // jungle: dense and tall
				if u < 0.46 {
					add("jtree", px, pz, s*1.05, 0.7)
				} else if u < 0.51 && h < 2.2 {
					add("palm", px, pz, s, 0.4)
				} else if u < 0.68 {
					add("fern", px, pz, s, 0)
				} else if u < 0.76 {
					add("bamboo", px, pz, s*0.9, 0.4)
				} else if u < 0.82 {
					add("bush", px, pz, s, 0.4)
				} else if u < 0.88 {
					add("flower", px, pz, s, 0)
				} else if u < 0.9 {
					add("rock", px, pz, s, 0.6)
				}
			case 7: // cherry blossom groves
				if u < 0.2 {
					add("sakura", px, pz, s, 0.5)
				} else if u < 0.28 {
					add("bamboo", px, pz, s, 0.4)
				} else if u < 0.32 {
					add("bush", px, pz, s*0.8, 0.4)
				} else if u < 0.34 {
					add("rock", px, pz, s*0.8, 0.6)
				} else if u < 0.58 {
					add("tuft", px, pz, s, 0)
				} else if u < 0.72 {
					add("flower", px, pz, s, 0)
				}
			case 8: // red canyons
				if u < 0.045 && slope < 2.5 {
					add("mesapillar", px, pz, s*1.2, 1.3)
				} else if u < 0.12 {
					add("rock", px, pz, s*1.1, 0.8)
				} else if u < 0.16 {
					add("cactus", px, pz, s, 0.5)
				} else if u < 0.26 {
					add("drybush", px, pz, s, 0)
				} else if u < 0.28 {
					add("bones", px, pz, s, 0)
				}
			case 9: // volcanic ash
				if u < 0.12 && slope < 3 {
					add("charred", px, pz, s, 0.4)
				} else if u < 0.25 {
					add("rock", px, pz, s*1.2, 0.8)
				} else if u < 0.32 {
					add("obsidian", px, pz, s, 0.7)
				} else if u < 0.36 {
					add("ember", px, pz, s, 0)
				}
			}
		}
	}
	return out
}

func qBiome(x, z float64) Biome {
	i := max(0, min(GridNX-1, int(math.Round((x-WorldMinX)/GridCell))))
	j := max(0, min(GridNZ-1, int(math.Round((z-WorldMinZ)/GridCell))))
	return Biome(bm[j*GridNX+i])
}

// PropsB64 packs the props for the clients: 8 bytes each = kind, x u16, z u16 (1/40 m), scale u8 (1/100),
// rotation u8, solid radius u8 (1/100).
func PropsB64(props []Prop) string {
	buf := make([]byte, 0, len(props)*8)
	for _, p := range props {
		qx, qz := uint16(math.Round((p.X-WorldMinX)*propPosQ)), uint16(math.Round((p.Z-WorldMinZ)*propPosQ))
		buf = append(buf, byte(p.K), byte(qx), byte(qx>>8), byte(qz), byte(qz>>8),
			byte(math.Round(p.S*100)), byte(math.Round(p.Ry/6.28*propRyQ)), byte(math.Round(p.Solid*100)))
	}
	return base64.StdEncoding.EncodeToString(buf)
}

// ---- manual prop edits (map editor) ----

// PropEdit places one prop by hand.
type PropEdit struct {
	Kind string  `json:"k"`
	X    float64 `json:"x"`
	Z    float64 `json:"z"`
	S    float64 `json:"s"`
	Ry   float64 `json:"ry"`
}

// defaultSolid is the collision radius per unit scale of hand placed props.
var defaultSolid = map[string]float64{"pine": 0.5, "oak": 0.5, "dead": 0.4, "snowpine": 0.5, "cactus": 0.5, "palm": 0.4,
	"rock": 0.7, "bush": 0.4, "jtree": 0.7, "bamboo": 0.4, "sakura": 0.5, "mesapillar": 1.3, "charred": 0.4, "obsidian": 0.7}

var (
	propAdd []PropEdit
	propDel [][2]float64
)

func propKey(x, z float64) [2]int {
	return [2]int{int(math.Round((x - WorldMinX) * propPosQ)), int(math.Round((z - WorldMinZ) * propPosQ))}
}

func quantProp(k int, x, z, s, ry, solid float64) Prop {
	qx, qz := math.Round((x-WorldMinX)*propPosQ), math.Round((z-WorldMinZ)*propPosQ)
	ry = math.Mod(math.Mod(ry, 6.28)+6.28, 6.28)
	return Prop{K: k, X: WorldMinX + qx/propPosQ, Z: WorldMinZ + qz/propPosQ, S: math.Min(2.55, math.Max(0.05, math.Round(s*100)/100)),
		Ry: math.Round(ry/6.28*propRyQ) / propRyQ * 6.28, Solid: solid}
}

func applyPropEdits(gen []Prop) []Prop {
	if len(propDel) > 0 {
		del := map[[2]int]bool{}
		for _, d := range propDel {
			del[propKey(d[0], d[1])] = true
		}
		out := gen[:0:0]
		for _, p := range gen {
			if !del[propKey(p.X, p.Z)] {
				out = append(out, p)
			}
		}
		gen = out
	}
	for _, a := range propAdd {
		if k := kindIndex2(a.Kind); k >= 0 {
			gen = append(gen, quantProp(k, a.X, a.Z, a.S, a.Ry, defaultSolid[a.Kind]))
		}
	}
	return gen
}

func kindIndex2(name string) int {
	for i, k := range PropKinds {
		if k == name {
			return i
		}
	}
	return -1
}

// ResetProps forgets the generated props and colliders (after the terrain, spawns or buildings changed).
func resetProps() {
	resetNav()
	propsOnce = sync.Once{}
	collOnce = sync.Once{}
	propList = nil
	circles, boxes, boxCircle = nil, nil, nil
}
