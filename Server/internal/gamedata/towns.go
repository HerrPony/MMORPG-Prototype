package gamedata

import (
	"math"
	"math/rand"
	"sync"
)

type Town struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Kind   string  `json:"kind"` // capital | village | site
	X      float64 `json:"x"`
	Z      float64 `json:"z"`
	FlatR  float64 `json:"-"`
	BlendR float64 `json:"-"`
	H      float64 `json:"-"`
	Style  string  `json:"style"` // building style: timber, plaster, stilt, lodge, adobe, bamboo, jade, redrock, obsidian
	Harbor bool    `json:"harbor,omitempty"`
	Cont   int     `json:"cont"`
}

var Towns = []*Town{
	{ID: "eichenhain", Name: "Eichenhain", Kind: "village", X: 6, Z: 12, FlatR: 20, BlendR: 40, H: 1.8, Style: "timber"},
	{ID: "sonnenburg", Name: "Sonnenburg", Kind: "capital", X: 190, Z: -120, FlatR: 62, BlendR: 95, H: 3.0, Style: "plaster"},
	{ID: "moorbruck", Name: "Moorbruck", Kind: "village", X: -170, Z: 170, FlatR: 24, BlendR: 44, H: 1.4, Style: "stilt"},
	{ID: "frostwacht", Name: "Frostwacht", Kind: "village", X: -110, Z: -215, FlatR: 24, BlendR: 46, H: 9, Style: "lodge"},
	{ID: "sandfurt", Name: "Sandfurt", Kind: "village", X: 190, Z: 190, FlatR: 24, BlendR: 44, H: 2.6, Style: "adobe"},
	{ID: "rabenfels", Name: "Rabenfels", Kind: "site", X: -175, Z: -196, FlatR: 11, BlendR: 22, H: 11, Style: "lodge"},
	{ID: "sturmhafen", Name: "Sturmhafen", Kind: "village", X: 232, Z: 40, FlatR: 18, BlendR: 40, H: 1.6, Style: "timber", Harbor: true},
	// second continent
	{ID: "lotushafen", Name: "Lotushafen", Kind: "village", X: 612, Z: 20, FlatR: 20, BlendR: 40, H: 1.8, Style: "bamboo", Harbor: true, Cont: 1},
	{ID: "himmelspforte", Name: "Himmelspforte", Kind: "capital", X: 830, Z: -10, FlatR: 66, BlendR: 100, H: 5.0, Style: "jade", Cont: 1},
	{ID: "bambushain", Name: "Bambushain", Kind: "village", X: 688, Z: 118, FlatR: 24, BlendR: 44, H: 3.2, Style: "bamboo", Cont: 1},
	{ID: "kirschbach", Name: "Kirschbach", Kind: "village", X: 810, Z: 150, FlatR: 24, BlendR: 44, H: 3.6, Style: "jade", Cont: 1},
	{ID: "rotfels_dorf", Name: "Rotfels", Kind: "village", X: 1000, Z: 120, FlatR: 24, BlendR: 44, H: 6.5, Style: "redrock", Cont: 1},
	{ID: "aschenwacht", Name: "Aschenwacht", Kind: "village", X: 940, Z: -100, FlatR: 24, BlendR: 46, H: 7, Style: "obsidian", Cont: 1},
	{ID: "nebeltempel", Name: "Nebeltempel", Kind: "site", X: 768, Z: -112, FlatR: 12, BlendR: 24, H: 4.5, Style: "jade", Cont: 1},
	{ID: "feuerschlund", Name: "Feuerschlund", Kind: "site", X: 1000, Z: -198, FlatR: 12, BlendR: 24, H: 8, Style: "obsidian", Cont: 1},
}

func TownByID(id string) *Town {
	for _, t := range Towns {
		if t.ID == id {
			return t
		}
	}
	return nil
}

// POI is a named sub area; entering it shows its name and uncovers it on the map.
type POI struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	X     float64 `json:"x"`
	Z     float64 `json:"z"`
	R     float64 `json:"r"`
	Kind  string  `json:"kind"` // town | camp | cave | lake | ruin | wild
	Town  bool    `json:"town,omitempty"`
	Level int     `json:"level,omitempty"`
}

var POIs = []*POI{
	// Erlenwald
	{ID: "eichenhain", Name: "Eichenhain", X: 6, Z: 12, R: 32, Kind: "town", Town: true},
	{ID: "wolfshuegel", Name: "Wolfshügel", X: 38, Z: 2, R: 24, Kind: "wild"},
	{ID: "keilerwiese", Name: "Keilerwiese", X: -38, Z: 8, R: 22, Kind: "wild"},
	{ID: "spinnenwald", Name: "Spinnenwald", X: 6, Z: -48, R: 26, Kind: "wild"},
	{ID: "banditenlager", Name: "Banditenlager", X: 14, Z: 48, R: 20, Kind: "camp"},
	{ID: "baerenklamm", Name: "Bärenklamm", X: 66, Z: 36, R: 22, Kind: "cave"},
	{ID: "spiegelsee", Name: "Spiegelsee", X: 62, Z: 72, R: 22, Kind: "lake"},
	// Goldene Ebenen
	{ID: "sonnenburg", Name: "Sonnenburg", X: 190, Z: -120, R: 70, Kind: "town", Town: true},
	{ID: "weizenhof", Name: "Weizenhof", X: 140, Z: -60, R: 24, Kind: "wild"},
	{ID: "windmuehlenhuegel", Name: "Windmühlenhügel", X: 100, Z: -160, R: 22, Kind: "wild"},
	{ID: "jaegerhain", Name: "Jägerhain", X: 130, Z: -212, R: 22, Kind: "wild"},
	{ID: "quellheiligtum", Name: "Quellheiligtum", X: 178, Z: -228, R: 18, Kind: "ruin"},
	{ID: "raeuberfestung", Name: "Räuberfestung", X: 258, Z: -70, R: 22, Kind: "camp"},
	// Düstermoor
	{ID: "moorbruck", Name: "Moorbruck", X: -170, Z: 170, R: 32, Kind: "town", Town: true},
	{ID: "krokodilsuempfe", Name: "Krokodilsümpfe", X: -130, Z: 215, R: 26, Kind: "wild"},
	{ID: "hexenhain", Name: "Hexenhain", X: -225, Z: 92, R: 22, Kind: "camp"},
	{ID: "nebelsee", Name: "Nebelsee", X: -176, Z: 202, R: 24, Kind: "lake"},
	{ID: "schiefer_turm", Name: "Schiefer Turm", X: -92, Z: 175, R: 18, Kind: "ruin"},
	// Eisgipfel
	{ID: "frostwacht", Name: "Frostwacht", X: -110, Z: -215, R: 32, Kind: "town", Town: true},
	{ID: "yetihoehle", Name: "Yetihöhle", X: -40, Z: -250, R: 24, Kind: "cave"},
	{ID: "windpass", Name: "Windpass", X: -60, Z: -172, R: 22, Kind: "wild"},
	{ID: "rabenfels", Name: "Rabenfels", X: -175, Z: -196, R: 16, Kind: "cave", Town: false},
	{ID: "baerenklippe", Name: "Eisbärenklippe", X: 30, Z: -256, R: 24, Kind: "wild"},
	// Glutwüste
	{ID: "sandfurt", Name: "Sandfurt", X: 190, Z: 190, R: 32, Kind: "town", Town: true},
	{ID: "skorpionduenen", Name: "Skorpiondünen", X: 245, Z: 140, R: 26, Kind: "wild"},
	{ID: "pharaonengrab", Name: "Pharaonengrab", X: 140, Z: 238, R: 22, Kind: "ruin"},
	{ID: "oase", Name: "Die Oase", X: 232, Z: 222, R: 18, Kind: "lake"},
}

func POIByID(id string) *POI {
	for _, p := range POIs {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// PoiAt returns the smallest POI containing the point (nil if none).
func PoiAt(x, z float64) *POI {
	var best *POI
	for _, p := range POIs {
		if math.Hypot(x-p.X, z-p.Z) <= p.R && (best == nil || p.R < best.R) {
			best = p
		}
	}
	return best
}

// ---- buildings ----

type Building struct {
	Kind  string  `json:"k"` // house, wall, tower, keep, fountain, stall, tent, pad, hut
	Style string  `json:"s,omitempty"`
	X     float64 `json:"x"`
	Z     float64 `json:"z"`
	W     float64 `json:"w"`
	D     float64 `json:"d"`
	H     float64 `json:"h"`
	Rot   int     `json:"r,omitempty"` // 1 = the front faces east/west (W and D swapped)
}

var (
	buildingsOnce sync.Once
	buildingList  []Building
)

// jadeCapital is the second capital: an open terraced city without walls, built around a five-tier palace.
func jadeCapital(t *Town) []Building {
	var b []Building
	cx, cz := t.X, t.Z
	b = append(b, Building{Kind: "pagoda", Style: "jade", X: cx, Z: cz - 30, W: 26, D: 22, H: 22})
	for _, p := range [][2]float64{{-34, -34}, {34, -34}, {-38, 30}, {38, 30}} {
		b = append(b, Building{Kind: "pagoda", Style: "jade", X: cx + p[0], Z: cz + p[1], W: 11, D: 11, H: 13})
	}
	b = append(b, Building{Kind: "pond", X: cx, Z: cz + 4, W: 24, D: 16, H: 0.5})
	for _, g := range [][3]float64{{-52, 0, 0}, {52, 0, 0}, {0, 48, 1}, {0, -46, 1}} {
		b = append(b, Building{Kind: "gate", Style: "jade", X: cx + g[0], Z: cz + g[1], W: 12, D: 3, H: 9, Rot: int(g[2])})
	}
	for i := -4; i <= 4; i++ {
		if i == 0 {
			continue
		}
		d := float64(i) * 11
		b = append(b, Building{Kind: "lantern", X: cx + d, Z: cz - 8, W: 1, D: 1, H: 3.2}, Building{Kind: "lantern", X: cx + d, Z: cz + 8, W: 1, D: 1, H: 3.2})
		b = append(b, Building{Kind: "lantern", X: cx - 8, Z: cz + d, W: 1, D: 1, H: 3.2}, Building{Kind: "lantern", X: cx + 8, Z: cz + d, W: 1, D: 1, H: 3.2})
	}
	return b
}

func capitalBuildings(t *Town) []Building {
	if t.Style == "jade" {
		return jadeCapital(t)
	}
	var b []Building
	cx, cz := t.X, t.Z
	x0, x1, z0, z1 := cx-50, cx+50, cz-45, cz+45
	const th, wh = 3.0, 7.0
	seg := func(ax, az, bx, bz float64) {
		b = append(b, Building{Kind: "wall", X: (ax + bx) / 2, Z: (az + bz) / 2, W: math.Abs(bx-ax) + th, D: math.Abs(bz-az) + th, H: wh})
	}
	seg(x0, z0, x1, z0)   // north
	seg(x1, z0, x1, z1)   // east
	seg(x0, z0, x0, cz-7) // west, upper part
	seg(x0, cz+7, x0, z1) // west, lower part
	seg(x0, z1, cx-7, z1) // south, left
	seg(cx+7, z1, x1, z1) // south, right
	for _, p := range [][2]float64{{x0, z0}, {x1, z0}, {x0, z1}, {x1, z1}, {x0, cz - 8}, {x0, cz + 8}, {cx - 8, z1}, {cx + 8, z1}} {
		b = append(b, Building{Kind: "tower", X: p[0], Z: p[1], W: 7, D: 7, H: 13})
	}
	b = append(b, Building{Kind: "keep", Style: "plaster", X: cx, Z: cz - 30, W: 28, D: 18, H: 17})
	b = append(b, Building{Kind: "fountain", X: cx, Z: cz + 6, W: 8, D: 8, H: 1.6})
	return b
}

func buildTown(t *Town, rng *rand.Rand, avoid [][3]float64) []Building {
	var out []Building
	var rx0, rx1, rz0, rz1 float64
	count, plazaR := 0, 0.0
	switch t.Kind {
	case "capital":
		out = capitalBuildings(t)
		rx0, rx1, rz0, rz1 = t.X-44, t.X+44, t.Z-38, t.Z+38
		count, plazaR = 26, 14
		if t.Style == "jade" {
			rx0, rx1, rz0, rz1 = t.X-50, t.X+50, t.Z-44, t.Z+44
			count, plazaR = 30, 16
		}
	case "village":
		rx0, rx1, rz0, rz1 = t.X-34, t.X+34, t.Z-34, t.Z+34
		count, plazaR = 9, 12
	default:
		return nil
	}
	hutKind := "house"
	if t.Style == "stilt" || t.Style == "bamboo" {
		hutKind = "hut"
	}
	placed := 0
	for try := 0; try < 4000 && placed < count; try++ {
		x := rx0 + rng.Float64()*(rx1-rx0)
		z := rz0 + rng.Float64()*(rz1-rz0)
		w, d := 7+rng.Float64()*4, 6+rng.Float64()*3.5
		rot := 0
		if rng.Intn(2) == 0 {
			rot = 1
			w, d = d, w
		}
		h := 4.2 + rng.Float64()*2
		r := math.Hypot(w, d)/2 + 1.5
		if math.Hypot(x-t.X, z-t.Z) < plazaR+r {
			continue
		}
		if t.Kind == "capital" && t.Style == "jade" {
			if math.Abs(z-t.Z) < 12+d/2 || math.Abs(x-t.X) < 12+w/2 { // cross shaped avenues
				continue
			}
			if math.Abs(x-t.X) < 20+w/2 && z < t.Z-10 { // palace court
				continue
			}
		} else if t.Kind == "capital" {
			if math.Abs(z-t.Z) < 11+d/2 && x < t.X { // road from the west gate
				continue
			}
			if math.Abs(x-t.X) < 11+w/2 && z > t.Z { // road from the south gate
				continue
			}
			if math.Abs(x-t.X) < 22+w/2 && z < t.Z-12 { // keep and courtyard
				continue
			}
			if x < t.X-44 || x > t.X+44 || z < t.Z-38 || z > t.Z+38 {
				continue
			}
		}
		ok := true
		for _, a := range avoid {
			if math.Hypot(x-a[0], z-a[1]) < a[2]+r {
				ok = false
				break
			}
		}
		for _, o := range out {
			if o.Kind == "wall" || o.Kind == "tower" {
				continue
			}
			if math.Hypot(x-o.X, z-o.Z) < r+math.Hypot(o.W, o.D)/2+1 {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		out = append(out, Building{Kind: hutKind, Style: t.Style, X: x, Z: z, W: w, D: d, H: h, Rot: rot})
		placed++
	}
	if t.Harbor { // a pier into the sea; the ferry berths next to it
		dir := 1.0
		if t.X > ContSplitX {
			dir = -1
		}
		out = append(out, Building{Kind: "dock", Style: t.Style, X: t.X + dir*26, Z: t.Z, W: 28, D: 5, H: 1.2})
		return out
	}
	// the griffin landing pad
	out = append(out, Building{Kind: "pad", X: t.X + padOffset(t)[0], Z: t.Z + padOffset(t)[1], W: 9, D: 9, H: 0.4})
	return out
}

func padOffset(t *Town) [2]float64 {
	switch t.ID {
	case "sonnenburg":
		return [2]float64{26, 14}
	case "eichenhain":
		return [2]float64{14, 12}
	case "sturmhafen":
		return [2]float64{15, 0}
	case "lotushafen":
		return [2]float64{-15, 0}
	case "himmelspforte":
		return [2]float64{30, 20}
	}
	return [2]float64{16, 10}
}

// Berth is the place on the water where the ferry lies at the pier of a harbor.
func Berth(t *Town) (float64, float64) {
	dir := 1.0
	if t.X > ContSplitX {
		dir = -1
	}
	z := t.Z + 4
	for d := 10.0; d < 90; d++ { // the first deep water in front of the pier
		if Height(t.X+dir*d, z) < -1.2 {
			return t.X + dir*(d+5), z
		}
	}
	return t.X + dir*40, z
}

// PadPos is where the griffin master of a town stands.
func PadPos(t *Town) (float64, float64) {
	o := padOffset(t)
	return t.X + o[0], t.Z + o[1]
}

func extraBuildings() []Building {
	return []Building{
		// banditenlager (tents)
		{Kind: "tent", X: 8, Z: 46, W: 5, D: 5, H: 3.5}, {Kind: "tent", X: 20, Z: 52, W: 5, D: 5, H: 3.5}, {Kind: "tent", X: 14, Z: 40, W: 5, D: 5, H: 3.5},
		// raeuberfestung (palisade tents)
		{Kind: "tent", X: 250, Z: -66, W: 6, D: 6, H: 4}, {Kind: "tent", X: 262, Z: -74, W: 6, D: 6, H: 4}, {Kind: "tent", X: 258, Z: -56, W: 6, D: 6, H: 4},
		// hexenhain huts
		{Kind: "hut", Style: "stilt", X: -222, Z: 88, W: 5, D: 5, H: 3}, {Kind: "hut", Style: "stilt", X: -232, Z: 98, W: 5, D: 5, H: 3},
		// pharaonengrab
		{Kind: "keep", Style: "adobe", X: 140, Z: 242, W: 20, D: 14, H: 8},
		// quellheiligtum
		{Kind: "tower", X: 172, Z: -232, W: 3, D: 3, H: 6}, {Kind: "tower", X: 184, Z: -224, W: 3, D: 3, H: 4},
		// schiefer turm
		{Kind: "tower", X: -92, Z: 175, W: 6, D: 6, H: 15},
		// dungeon entrance
		{Kind: "cave", X: -175, Z: -202, W: 12, D: 6, H: 9},
	}
}

// Buildings returns all static buildings (generated once, deterministic).
func Buildings() []Building {
	buildingsOnce.Do(func() {
		rng := rand.New(rand.NewSource(4242))
		for _, t := range Towns {
			var avoid [][3]float64
			for _, n := range NPCs {
				if math.Hypot(n.X-t.X, n.Z-t.Z) < 70 {
					avoid = append(avoid, [3]float64{n.X, n.Z, 5})
				}
			}
			for _, r := range Roads {
				_ = r
			}
			buildingList = append(buildingList, buildTown(t, rng, avoid)...)
		}
		buildingList = append(buildingList, extraBuildings()...)
		buildingList = append(buildingList, extraBuildings2()...)
	})
	return buildingList
}

// ---- scenery payload ----

type Box struct {
	X float64 `json:"x"` // center
	Z float64 `json:"z"`
	W float64 `json:"w"` // size along x
	D float64 `json:"d"` // size along z
	H float64 `json:"h"`
}

type Clearing struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
}

type Bridge struct {
	X   float64 `json:"x"`
	Z   float64 `json:"z"`
	Rot float64 `json:"rot"`
	Len float64 `json:"len"`
}

type Scenery struct {
	X0        float64        `json:"x0"` // world position of grid column 0
	Z0        float64        `json:"z0"`
	Cell      float64        `json:"cell"`
	NX        int            `json:"nx"`
	NZ        int            `json:"nz"`
	Conts     []Continent    `json:"conts"`
	Heights   string         `json:"hm"`
	Biomes    string         `json:"bm"`
	Zones     []Zone         `json:"zones"`
	Towns     []*Town        `json:"towns"`
	POIs      []*POI         `json:"pois"`
	Buildings []Building     `json:"buildings"`
	Lamps     []Lamp         `json:"lamps"`
	House     []Box          `json:"house"` // the cabin at the world origin (with door)
	Roads     [][][2]float64 `json:"roads"`
	Rivers    []River        `json:"rivers"`
	Lakes     []Lake         `json:"lakes"`
	Bridges   []Bridge       `json:"bridges"`
	Clears    []Clearing     `json:"clears"`
	PropKinds []string       `json:"pk"`
	Props     string         `json:"props"` // packed, see PropsB64
}

func segIntersect(a, b, c, d [2]float64) (float64, float64, bool) {
	r := [2]float64{b[0] - a[0], b[1] - a[1]}
	s := [2]float64{d[0] - c[0], d[1] - c[1]}
	den := r[0]*s[1] - r[1]*s[0]
	if math.Abs(den) < 1e-9 {
		return 0, 0, false
	}
	t := ((c[0]-a[0])*s[1] - (c[1]-a[1])*s[0]) / den
	u := ((c[0]-a[0])*r[1] - (c[1]-a[1])*r[0]) / den
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return 0, 0, false
	}
	return a[0] + t*r[0], a[1] + t*r[1], true
}

func bridges() []Bridge {
	var out []Bridge
	for _, rd := range Roads {
		for i := 0; i+1 < len(rd.Pts); i++ {
			for _, rv := range Rivers {
				for k := 0; k+1 < len(rv.Pts); k++ {
					if x, z, ok := segIntersect(rd.Pts[i], rd.Pts[i+1], rv.Pts[k], rv.Pts[k+1]); ok {
						out = append(out, Bridge{X: x, Z: z, Rot: math.Atan2(rd.Pts[i+1][0]-rd.Pts[i][0], rd.Pts[i+1][1]-rd.Pts[i][1]), Len: rv.Width*2 + 6})
					}
				}
			}
		}
	}
	return out
}

// BuildScenery returns the static world layout sent to clients.
func BuildScenery() Scenery {
	sc := baseScenery()
	sc.PropKinds = PropKinds
	sc.Props = PropsB64(worldProps(sc))
	return sc
}

// baseScenery is everything except the props (which are generated from it).
func baseScenery() Scenery {
	ensureWorld()
	const hw, hd, t, wallH = 5.0, 4.0, 0.4, 3.0
	const doorW = 1.8
	sideW := (2*hw - doorW) / 2
	house := []Box{
		{X: 0, Z: -hd, W: 2 * hw, D: t, H: wallH},
		{X: -hw, Z: 0, W: t, D: 2 * hd, H: wallH},
		{X: hw, Z: 0, W: t, D: 2 * hd, H: wallH},
		{X: -hw + sideW/2, Z: hd, W: sideW, D: t, H: wallH},
		{X: hw - sideW/2, Z: hd, W: sideW, D: t, H: wallH},
	}
	var clears []Clearing
	for _, tw := range Towns {
		clears = append(clears, Clearing{tw.X, tw.Z, tw.FlatR * 0.95})
	}
	for _, n := range NPCs {
		clears = append(clears, Clearing{n.X, n.Z, 3.5})
	}
	for _, s := range AllSpawns() {
		clears = append(clears, Clearing{s.X, s.Z, 3})
	}
	for _, o := range ObjectSpawns {
		clears = append(clears, Clearing{o.X, o.Z, 2.5})
	}
	roads := make([][][2]float64, len(Roads))
	for i, r := range Roads {
		roads[i] = r.Pts
	}
	return Scenery{
		X0: WorldMinX, Z0: WorldMinZ, Cell: GridCell, NX: GridNX, NZ: GridNZ, Conts: Continents, Heights: HeightmapB64(), Biomes: BiomemapB64(),
		Zones: Zones, Towns: Towns, POIs: POIs, Buildings: Buildings(), Lamps: Lamps(), House: house,
		Roads: roads, Rivers: Rivers, Lakes: Lakes, Bridges: bridges(), Clears: clears,
	}
}
