package gamedata

import (
	"encoding/base64"
	"encoding/binary"
	"math"
	"sync"
)

// The world consists of two continents (each about 900 x 900 m) separated by a sea, stored as a
// height grid and a biome grid. The server generates it once at start-up and sends both grids to the
// clients, so client and server always agree on the terrain.
const (
	// WorldScale stretches the original 600 x 600 m continents; all positions in the data tables are written for
	// the original size and moved by scaleWorld (zzz_scale.go) before anything uses them.
	WorldScale = 1.5

	WorldHalf = 300.0 * WorldScale // half size of one continent
	GridCell  = 2.0
	WorldMinX = -300.0 * WorldScale
	WorldMaxX = 1140.0 * WorldScale
	WorldMinZ = -300.0 * WorldScale
	WorldMaxZ = 300.0 * WorldScale
	GridNX    = int(1440*WorldScale/GridCell) + 1 // vertices along x
	GridNZ    = int(600*WorldScale/GridCell) + 1  // vertices along z
	DeepWater = -1.6
	// Swimming: players may leave the shallows into the sea as long as the sea floor is above SwimFloor
	// (a belt of some 30 m around the coast). The open sea between the continents stays closed, so the
	// other continent can only be reached with the ferry. Clients use the same numbers.
	SwimFloor = -6.8
	SwimLevel = -1.1 // y of a swimming player's feet (the water surface is at about -0.1)
	SeaLevel  = 0.0

	Cont2X     = 840.0 * WorldScale // center of the second continent
	ContSplitX = 420.0 * WorldScale // x beyond this belongs to the second continent
)

// Dungeon instances live far outside the continents on flat ground.
const DungeonZ0 = 600.0

type Biome uint8

const (
	BiomeSea Biome = iota
	BiomeForest
	BiomePlains
	BiomeSwamp
	BiomeSnow
	BiomeDesert
	BiomeJungle
	BiomeBlossom
	BiomeMesa
	BiomeVolcano
	BiomeLava // lava lakes and rivers inside the volcano zone (not walkable)
	numBiomes
)

type Zone struct {
	ID     Biome  `json:"id"`
	Key    string `json:"key"`
	Name   string `json:"name"`
	MinLvl int    `json:"minLvl"`
	MaxLvl int    `json:"maxLvl"`
	Color  string `json:"color"` // map color
	Cont   int    `json:"cont"`
}

var Zones = []Zone{
	{BiomeSea, "sea", "Das Sturmmeer", 0, 0, "#27567e", 0},
	{BiomeForest, "forest", "Erlenwald", 1, 6, "#3d7a38", 0},
	{BiomePlains, "plains", "Goldene Ebenen", 6, 10, "#a9b552", 0},
	{BiomeSwamp, "swamp", "Düstermoor", 10, 14, "#566a3a", 0},
	{BiomeSnow, "snow", "Eisgipfel", 14, 18, "#d8e4ee", 0},
	{BiomeDesert, "desert", "Glutwüste", 17, 20, "#d3b46a", 0},
	{BiomeJungle, "jungle", "Smaragddschungel", 20, 25, "#1f9a4c", 1},
	{BiomeBlossom, "blossom", "Kirschblütenhain", 25, 30, "#f0a8c8", 1},
	{BiomeMesa, "mesa", "Rotkluft-Canyon", 30, 35, "#c8603a", 1},
	{BiomeVolcano, "volcano", "Aschenlande", 35, 40, "#4a3a44", 1},
	{BiomeLava, "volcano", "Aschenlande", 35, 40, "#ff5a1a", 1},
}

// Continents describes the two land masses for the world map.
type Continent struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	X0   float64 `json:"x0"` // west edge of the map view
	Size float64 `json:"size"`
}

var Continents = []Continent{
	{0, "Erlenmark", -WorldHalf, 2 * WorldHalf},
	{1, "Jadereich", Cont2X - WorldHalf, 2 * WorldHalf},
}

// ContinentOf returns the continent a world x position belongs to.
func ContinentOf(x float64) int {
	if x > ContSplitX {
		return 1
	}
	return 0
}

func ZoneByBiome(b Biome) Zone { return Zones[b] }

// ---- noise ----

func hash2(ix, iz int, seed uint32) float64 {
	h := uint32(ix)*374761393 + uint32(iz)*668265263 + seed*2246822519
	h = (h ^ (h >> 13)) * 1274126177
	h ^= h >> 16
	return float64(h&0xffffff) / float64(0xffffff)
}

func vnoise(x, z float64, seed uint32) float64 {
	ix, iz := int(math.Floor(x)), int(math.Floor(z))
	fx, fz := x-float64(ix), z-float64(iz)
	fx, fz = fx*fx*(3-2*fx), fz*fz*(3-2*fz)
	a, b := hash2(ix, iz, seed), hash2(ix+1, iz, seed)
	c, d := hash2(ix, iz+1, seed), hash2(ix+1, iz+1, seed)
	return a + (b-a)*fx + (c-a)*fz + (a-b-c+d)*fx*fz
}

func fbm(x, z float64, oct int, seed uint32) float64 {
	sum, amp, norm := 0.0, 1.0, 0.0
	for i := 0; i < oct; i++ {
		sum += amp * vnoise(x, z, seed+uint32(i)*131)
		norm += amp
		x, z, amp = x*2.03, z*2.03, amp*0.5
	}
	return sum / norm
}

func smoothstep(a, b, x float64) float64 {
	t := math.Max(0, math.Min(1, (x-a)/(b-a)))
	return t * t * (3 - 2*t)
}

func lerp(a, b, t float64) float64 { return a + (b-a)*t }

// ---- layout ----

type anchor struct {
	b       Biome
	x, z, s float64
}

var anchors = []anchor{
	// first continent
	{BiomeForest, -10, 10, 130},
	{BiomePlains, 200, -110, 115}, {BiomePlains, 110, -170, 70},
	{BiomeSwamp, -175, 175, 105}, {BiomeSwamp, -120, 110, 45},
	{BiomeSnow, -90, -225, 110}, {BiomeSnow, 40, -250, 90}, {BiomeSnow, -210, -130, 70},
	{BiomeDesert, 195, 190, 105}, {BiomeDesert, 130, 110, 60},
	// second continent
	{BiomeJungle, 640, 60, 105}, {BiomeJungle, 700, 150, 80}, {BiomeJungle, 640, -100, 75},
	{BiomeBlossom, 830, -10, 115}, {BiomeBlossom, 800, 130, 70},
	{BiomeMesa, 1010, 90, 105}, {BiomeMesa, 950, 190, 70},
	{BiomeVolcano, 1010, -120, 105}, {BiomeVolcano, 900, -200, 75},
}

type Lake struct {
	X    float64 `json:"x"`
	Z    float64 `json:"z"`
	R    float64 `json:"r"`
	Kind string  `json:"k,omitempty"` // "" = water, "lava"
}

var Lakes = []Lake{
	{62, 72, 17, ""}, {-176, 192, 20, ""}, {-205, 140, 14, ""}, {232, 222, 13, ""}, {250, -190, 16, ""},
	{730, 160, 18, ""}, {905, 62, 13, ""}, {1045, 150, 10, ""},
	{1010, -140, 16, "lava"}, {960, -195, 9, "lava"}, {1075, -60, 9, "lava"},
}

type River struct {
	Pts   [][2]float64 `json:"pts"`
	Width float64      `json:"w"`
	Kind  string       `json:"k,omitempty"`
}

var Rivers = []River{
	{[][2]float64{{-125, -160}, {-100, -105}, {-86, -50}, {-90, 5}, {-100, 55}, {-124, 98}, {-150, 135}, {-176, 175}}, 6, ""},
	{[][2]float64{{150, -40}, {185, -52}, {225, -42}, {262, -34}, {295, -30}}, 5, ""},
	{[][2]float64{{222, -260}, {218, -215}, {248, -190}, {262, -150}, {290, -120}}, 5, ""},
	{[][2]float64{{760, -70}, {728, -20}, {716, 30}, {688, 70}, {668, 110}, {640, 150}, {612, 190}}, 6, ""},
	{[][2]float64{{928, 20}, {910, 80}, {880, 130}, {850, 190}, {840, 250}}, 4, ""},
	{[][2]float64{{1010, -140}, {1040, -112}, {1068, -84}, {1100, -44}}, 4, "lava"},
}

type Road struct {
	Pts [][2]float64
}

// Roads connect the settlements.
var Roads = []Road{
	{[][2]float64{{6, 14}, {30, 8}, {62, -10}, {100, -50}, {128, -100}, {140, -120}}},      // Eichenhain - Sonnenburg (west gate)
	{[][2]float64{{190, -75}, {198, -10}, {196, 70}, {188, 120}, {190, 168}}},              // Sonnenburg - Sandfurt
	{[][2]float64{{6, 14}, {-30, 40}, {-70, 70}, {-110, 90}, {-150, 128}, {-168, 158}}},    // Eichenhain - Moorbruck
	{[][2]float64{{6, 14}, {4, -40}, {-30, -90}, {-70, -160}, {-100, -205}, {-110, -215}}}, // Eichenhain - Frostwacht
	{[][2]float64{{-110, -215}, {-140, -205}, {-175, -196}}},                               // Frostwacht - Schwarzfelsmine
	{[][2]float64{{8, 22}, {10, 36}, {14, 56}}},                                            // to the bandit camp
	{[][2]float64{{198, 20}, {214, 30}, {232, 40}}},                                        // Sandfurt road - Sturmhafen
	// second continent
	{[][2]float64{{610, 20}, {650, 14}, {700, 6}, {760, -2}, {800, -6}}}, // Lotushafen - Himmelspforte (west gate)
	{[][2]float64{{630, 24}, {660, 60}, {688, 110}}},                     // to Bambushain
	{[][2]float64{{830, 28}, {822, 70}, {814, 110}, {810, 150}}},         // Himmelspforte (south gate) - Kirschbach
	{[][2]float64{{860, -6}, {910, 20}, {960, 60}, {1000, 120}}},         // Himmelspforte (east gate) - Rotfels
	{[][2]float64{{840, -40}, {880, -70}, {910, -90}, {940, -100}}},      // Himmelspforte (north gate) - Aschenwacht
	{[][2]float64{{760, -2}, {760, -50}, {768, -110}}},                   // to the Nebeltempel
	{[][2]float64{{940, -100}, {960, -150}, {985, -190}, {1000, -198}}},  // Aschenwacht - Feuerschlund
}

func distToSeg(px, pz, ax, az, bx, bz float64) float64 {
	dx, dz := bx-ax, bz-az
	l2 := dx*dx + dz*dz
	t := 0.0
	if l2 > 0 {
		t = math.Max(0, math.Min(1, ((px-ax)*dx+(pz-az)*dz)/l2))
	}
	return math.Hypot(px-(ax+dx*t), pz-(az+dz*t))
}

func distToPoly(x, z float64, pts [][2]float64) float64 {
	d := math.Inf(1)
	for i := 0; i+1 < len(pts); i++ {
		d = math.Min(d, distToSeg(x, z, pts[i][0], pts[i][1], pts[i+1][0], pts[i+1][1]))
	}
	return d
}

// DistToRoad is the distance from a point to the nearest road centre line.
func DistToRoad(x, z float64) float64 {
	d := math.Inf(1)
	for _, r := range Roads {
		d = math.Min(d, distToPoly(x, z, r.Pts))
	}
	return d
}

// ---- grids ----

var (
	hm    [GridNX * GridNZ]float32
	bm    [GridNX * GridNZ]uint8
	genWG sync.Once
)

func gridXZ(i, j int) (float64, float64) {
	return WorldMinX + float64(i)*GridCell, WorldMinZ + float64(j)*GridCell
}

func biomeWeights(x, z float64) [numBiomes]float64 {
	wx := x + 40*(fbm(x*0.006, z*0.006, 3, 11)-0.5)
	wz := z + 40*(fbm(x*0.006+50, z*0.006+50, 3, 12)-0.5)
	var w [numBiomes]float64
	for _, a := range anchors {
		d := math.Hypot(wx-a.x, wz-a.z) / a.s
		w[a.b] += math.Exp(-d * d)
	}
	return w
}

// lift keeps dry land clearly above the water line.
func lift(h float64) float64 {
	if h < 1.4 {
		return 1.4 - (1.4-h)*0.25
	}
	return h
}

func rawHeight(b Biome, x, z float64) float64 {
	switch b {
	case BiomeForest:
		return lift(3.4 + (fbm(x*0.014, z*0.014, 4, 21)-0.5)*9)
	case BiomePlains:
		return lift(2.6 + (fbm(x*0.009, z*0.009, 3, 22)-0.5)*5)
	case BiomeSwamp:
		n := fbm(x*0.03, z*0.03, 3, 23)
		h := 0.5 + (n-0.5)*2.2
		if n < 0.38 { // shallow pools
			h = lerp(h, -0.9, smoothstep(0.38, 0.28, n))
		}
		return h
	case BiomeSnow:
		r := 1 - math.Abs(fbm(x*0.011, z*0.011, 4, 24)*2-1) // ridged noise
		return 5 + r*r*34 + fbm(x*0.05, z*0.05, 2, 25)*3
	case BiomeDesert:
		dunes := math.Sin(x*0.045+fbm(x*0.02, z*0.02, 2, 26)*6) * math.Cos(z*0.03) * 1.4
		return lift(2.8 + (fbm(x*0.012, z*0.012, 3, 27)-0.5)*5 + dunes)
	case BiomeJungle: // steep, overgrown hills
		return lift(4.5 + (fbm(x*0.02, z*0.02, 4, 41)-0.5)*16 + math.Sin(x*0.09)*math.Cos(z*0.07)*0.8)
	case BiomeBlossom: // gentle terraced slopes
		h := 4.0 + (fbm(x*0.012, z*0.012, 3, 42)-0.5)*9
		step := math.Floor(h / 1.6)
		f := h/1.6 - step
		return lift(1.6*(step+smoothstep(0.55, 0.95, f)) + 0.5)
	case BiomeMesa: // flat-topped plateaus split by canyons
		v := fbm(x*0.008, z*0.008, 3, 43)
		h := 3 + 7*smoothstep(0.40, 0.62, v) + 7*smoothstep(0.58, 0.8, v)
		return lift(h + fbm(x*0.06, z*0.06, 2, 44)*1.2)
	case BiomeVolcano: // jagged ash mountains around a great cone
		r := 1 - math.Abs(fbm(x*0.016, z*0.016, 4, 45)*2-1)
		cone := 30 * math.Exp(-math.Pow(math.Hypot(x-1010*WorldScale, z+140*WorldScale)/(70*WorldScale), 2))
		return 4 + r*r*16 + fbm(x*0.06, z*0.06, 2, 46)*2.5 + cone
	}
	return -6
}

// contMask is 1 on land and 0 at sea; each continent has its own shape.
func contMask(x, z float64) float64 {
	d1 := math.Pow(math.Pow(math.Abs(x), 4)+math.Pow(math.Abs(z), 4), 0.25)/WorldHalf + (fbm(x*0.02, z*0.02, 3, 31)-0.5)*0.12
	m1 := 1 - smoothstep(0.82, 0.97, d1)
	dx := math.Abs(x - Cont2X)
	d2 := math.Pow(math.Pow(dx, 3)+math.Pow(math.Abs(z), 3), 1.0/3)/WorldHalf + (fbm(x*0.018+70, z*0.018, 3, 32)-0.5)*0.18
	m2 := 1 - smoothstep(0.80, 0.96, d2)
	return math.Max(m1, m2)
}

func generate() {
	worldGenerated = true
	for j := 0; j < GridNZ; j++ {
		for i := 0; i < GridNX; i++ {
			x, z := gridXZ(i, j)
			w := biomeWeights(x, z)
			best, bw := BiomeForest, -1.0
			sum := 0.0
			for b := BiomeForest; b <= BiomeVolcano; b++ {
				sum += w[b]
				if w[b] > bw {
					best, bw = b, w[b]
				}
			}
			h := -7.0
			if sum > 1e-12 {
				h = 0
				for b := BiomeForest; b <= BiomeVolcano; b++ {
					h += w[b] / sum * rawHeight(b, x, z)
				}
			}
			m := contMask(x, z)
			h = lerp(-7, h, m)
			// lakes
			lava := false
			for _, l := range Lakes {
				dd := math.Hypot(x-l.X, z-l.Z)
				if l.Kind == "lava" {
					if dd < l.R*1.4 {
						h = lerp(h, 2.0, 1-smoothstep(l.R*0.9, l.R*1.4, dd)) // crater rim
					}
					if dd < l.R {
						lava = true
					}
					continue
				}
				if dd < l.R*1.5 {
					h = lerp(h, -1.1, 1-smoothstep(l.R*0.55, l.R*1.1, dd))
				}
			}
			// rivers
			for _, r := range Rivers {
				dd := distToPoly(x, z, r.Pts)
				if r.Kind == "lava" {
					if dd < r.Width*0.7 {
						lava = true
					}
					continue
				}
				if dd < r.Width*1.3 {
					h = lerp(h, math.Min(h, -1.0), 1-smoothstep(r.Width*0.3, r.Width, dd))
				}
			}
			// settlements are flattened
			for _, t := range Towns {
				dd := math.Hypot(x-t.X, z-t.Z)
				if dd < t.BlendR {
					h = lerp(t.H, h, smoothstep(t.FlatR, t.BlendR, dd))
					if dd < t.FlatR {
						lava = false
					}
				}
			}
			idx := j*GridNX + i
			if lava && h > -2 {
				h = -0.05 // lava surface: just above the water plane so the sea does not cover it
				best = BiomeLava
			}
			hm[idx] = float32(h)
			if h < -2.5 {
				best = BiomeSea
			}
			bm[idx] = uint8(best)
		}
	}
}

var worldGenerated bool // set when the terrain is generated; the data tables must be scaled before (zzz_scale.go)

func ensureWorld() { genWG.Do(generate) }

func gridAt(i, j int) float64 {
	i = max(0, min(GridNX-1, i))
	j = max(0, min(GridNZ-1, j))
	return float64(hm[j*GridNX+i])
}

// Height is the terrain height at a world position.
func Height(x, z float64) float64 {
	if z > DungeonZ0-100 {
		return 0
	}
	ensureWorld()
	gx, gz := (x-WorldMinX)/GridCell, (z-WorldMinZ)/GridCell
	i, j := int(math.Floor(gx)), int(math.Floor(gz))
	fx, fz := gx-float64(i), gz-float64(j)
	a, b := gridAt(i, j), gridAt(i+1, j)
	c, d := gridAt(i, j+1), gridAt(i+1, j+1)
	return a + (b-a)*fx + (c-a)*fz + (a-b-c+d)*fx*fz
}

// BiomeAt returns the biome at a world position.
func BiomeAt(x, z float64) Biome {
	if z > DungeonZ0-100 {
		return BiomeForest
	}
	ensureWorld()
	i := max(0, min(GridNX-1, int(math.Round((x-WorldMinX)/GridCell))))
	j := max(0, min(GridNZ-1, int(math.Round((z-WorldMinZ)/GridCell))))
	return Biome(bm[j*GridNX+i])
}

// OnDock reports whether a position lies on the pier of a harbor (a wooden walkway over the water).
func OnDock(x, z float64) bool {
	for _, t := range Towns {
		if !t.Harbor {
			continue
		}
		dir := 1.0
		if t.X > ContSplitX {
			dir = -1
		}
		off, w := pierSpan(t)
		cx := t.X + dir*off
		if math.Abs(x-cx) <= w/2 && math.Abs(z-t.Z) <= 2.5 {
			return true
		}
	}
	return false
}

// Walkable reports whether a player may stand at x,z in the open world.
func Walkable(x, z float64) bool {
	if OnDock(x, z) {
		return true
	}
	return x > WorldMinX+4 && x < WorldMaxX-4 && math.Abs(z) < WorldHalf-4 && Height(x, z) > DeepWater && BiomeAt(x, z) != BiomeLava
}

// Swimmable reports whether a player may be at x,z in the water: everything walkable plus the coastal
// belt of the sea. The sea floor drops to -7 away from the coasts, which keeps the continents apart.
func Swimmable(x, z float64) bool {
	if Walkable(x, z) {
		return true
	}
	return x > WorldMinX+4 && x < WorldMaxX-4 && math.Abs(z) < WorldHalf-4 && Height(x, z) > SwimFloor && BiomeAt(x, z) != BiomeLava
}

// FootY is the height a player's feet are at: the terrain, but not deeper than the swimming level.
func FootY(x, z float64) float64 { return math.Max(Height(x, z), SwimLevel) }

// HeightmapB64 / BiomemapB64 are the grids as sent to clients (int16 of h*20, uint8 biome).
func HeightmapB64() string {
	ensureWorld()
	buf := make([]byte, GridNX*GridNZ*2)
	for i, v := range hm {
		binary.LittleEndian.PutUint16(buf[i*2:], uint16(int16(math.Round(float64(v)*20))))
	}
	return base64.StdEncoding.EncodeToString(buf)
}

func BiomemapB64() string {
	ensureWorld()
	return base64.StdEncoding.EncodeToString(bm[:])
}
