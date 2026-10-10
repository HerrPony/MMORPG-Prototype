package gamedata

import (
	"encoding/base64"
	"encoding/binary"
	"math"
	"sync"
)

// The world is a 600 x 600 m continent stored as a height grid and a biome grid.
// The server generates it once at start-up and sends both grids to the clients,
// so client and server always agree on the terrain.
const (
	WorldHalf = 300.0
	GridCell  = 2.0
	GridN     = 301 // vertices per side
	DeepWater = -1.6
	SeaLevel  = 0.0
)

// Dungeon instances live far outside the continent on flat ground.
const DungeonZ0 = 600.0

type Biome uint8

const (
	BiomeSea Biome = iota
	BiomeForest
	BiomePlains
	BiomeSwamp
	BiomeSnow
	BiomeDesert
)

type Zone struct {
	ID     Biome  `json:"id"`
	Key    string `json:"key"`
	Name   string `json:"name"`
	MinLvl int    `json:"minLvl"`
	MaxLvl int    `json:"maxLvl"`
	Color  string `json:"color"` // map color
}

var Zones = []Zone{
	{BiomeSea, "sea", "Das Sturmmeer", 0, 0, "#27567e"},
	{BiomeForest, "forest", "Erlenwald", 1, 6, "#3d7a38"},
	{BiomePlains, "plains", "Goldene Ebenen", 6, 10, "#a9b552"},
	{BiomeSwamp, "swamp", "Düstermoor", 10, 14, "#566a3a"},
	{BiomeSnow, "snow", "Eisgipfel", 14, 18, "#d8e4ee"},
	{BiomeDesert, "desert", "Glutwüste", 17, 20, "#d3b46a"},
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
	{BiomeForest, -10, 10, 130},
	{BiomePlains, 200, -110, 115}, {BiomePlains, 110, -170, 70},
	{BiomeSwamp, -175, 175, 105}, {BiomeSwamp, -120, 110, 45},
	{BiomeSnow, -90, -225, 110}, {BiomeSnow, 40, -250, 90}, {BiomeSnow, -210, -130, 70},
	{BiomeDesert, 195, 190, 105}, {BiomeDesert, 130, 110, 60},
}

type Lake struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
}

var Lakes = []Lake{{62, 72, 17}, {-176, 192, 20}, {-205, 140, 14}, {232, 222, 13}, {250, -190, 16}}

type River struct {
	Pts   [][2]float64 `json:"pts"`
	Width float64      `json:"w"`
}

var Rivers = []River{
	{[][2]float64{{-125, -160}, {-100, -105}, {-86, -50}, {-90, 5}, {-100, 55}, {-124, 98}, {-150, 135}, {-176, 175}}, 6},
	{[][2]float64{{150, -40}, {185, -52}, {225, -42}, {262, -34}, {295, -30}}, 5},
	{[][2]float64{{222, -260}, {218, -215}, {248, -190}, {262, -150}, {290, -120}}, 5},
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

// ---- grids ----

var (
	hm    [GridN * GridN]float32
	bm    [GridN * GridN]uint8
	genWG sync.Once
)

func gridXZ(i, j int) (float64, float64) {
	return float64(i)*GridCell - WorldHalf, float64(j)*GridCell - WorldHalf
}

func biomeWeights(x, z float64) [6]float64 {
	wx := x + 40*(fbm(x*0.006, z*0.006, 3, 11)-0.5)
	wz := z + 40*(fbm(x*0.006+50, z*0.006+50, 3, 12)-0.5)
	var w [6]float64
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
	}
	return -6
}

func generate() {
	for j := 0; j < GridN; j++ {
		for i := 0; i < GridN; i++ {
			x, z := gridXZ(i, j)
			w := biomeWeights(x, z)
			best, bw := BiomeForest, -1.0
			sum := 0.0
			for b := BiomeForest; b <= BiomeDesert; b++ {
				sum += w[b]
				if w[b] > bw {
					best, bw = b, w[b]
				}
			}
			h := 0.0
			for b := BiomeForest; b <= BiomeDesert; b++ {
				h += w[b] / sum * rawHeight(b, x, z)
			}
			// continent mask with a noisy coast line
			d := math.Pow(math.Pow(math.Abs(x), 4)+math.Pow(math.Abs(z), 4), 0.25)/WorldHalf + (fbm(x*0.02, z*0.02, 3, 31)-0.5)*0.12
			m := 1 - smoothstep(0.82, 0.97, d)
			h = lerp(-7, h, m)
			// lakes
			for _, l := range Lakes {
				dd := math.Hypot(x-l.X, z-l.Z)
				if dd < l.R*1.5 {
					h = lerp(h, -1.1, 1-smoothstep(l.R*0.55, l.R*1.1, dd))
				}
			}
			// rivers
			for _, r := range Rivers {
				dd := distToPoly(x, z, r.Pts)
				if dd < r.Width*1.3 {
					h = lerp(h, math.Min(h, -1.0), 1-smoothstep(r.Width*0.3, r.Width, dd))
				}
			}
			// settlements are flattened
			for _, t := range Towns {
				dd := math.Hypot(x-t.X, z-t.Z)
				if dd < t.BlendR {
					h = lerp(t.H, h, smoothstep(t.FlatR, t.BlendR, dd))
				}
			}
			idx := j*GridN + i
			hm[idx] = float32(h)
			if h < -2.5 {
				best = BiomeSea
			}
			bm[idx] = uint8(best)
		}
	}
}

func ensureWorld() { genWG.Do(generate) }

func gridAt(i, j int) float64 {
	i = max(0, min(GridN-1, i))
	j = max(0, min(GridN-1, j))
	return float64(hm[j*GridN+i])
}

// Height is the terrain height at a world position.
func Height(x, z float64) float64 {
	if z > DungeonZ0-100 {
		return 0
	}
	ensureWorld()
	gx, gz := (x+WorldHalf)/GridCell, (z+WorldHalf)/GridCell
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
	i := max(0, min(GridN-1, int(math.Round((x+WorldHalf)/GridCell))))
	j := max(0, min(GridN-1, int(math.Round((z+WorldHalf)/GridCell))))
	return Biome(bm[j*GridN+i])
}

// Walkable reports whether a player may stand at x,z in the open world.
func Walkable(x, z float64) bool {
	return math.Abs(x) < WorldHalf-4 && math.Abs(z) < WorldHalf-4 && Height(x, z) > DeepWater
}

// HeightmapB64 / BiomemapB64 are the grids as sent to clients (int16 of h*20, uint8 biome).
func HeightmapB64() string {
	ensureWorld()
	buf := make([]byte, GridN*GridN*2)
	for i, v := range hm {
		binary.LittleEndian.PutUint16(buf[i*2:], uint16(int16(math.Round(float64(v)*20))))
	}
	return base64.StdEncoding.EncodeToString(buf)
}

func BiomemapB64() string {
	ensureWorld()
	return base64.StdEncoding.EncodeToString(bm[:])
}
