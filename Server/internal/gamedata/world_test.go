package gamedata

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestWorldMap(t *testing.T) {
	ensureWorld()
	glyph := map[Biome]byte{BiomeSea: '~', BiomeForest: 'f', BiomePlains: 'p', BiomeSwamp: 's', BiomeSnow: 'S', BiomeDesert: 'd'}
	var sb strings.Builder
	for z := -WorldHalf; z <= WorldHalf; z += 12 {
		for x := -WorldHalf; x <= WorldHalf; x += 6 {
			c := glyph[BiomeAt(x, z)]
			if Height(x, z) < 0 && c != '~' {
				c = ','
			}
			for _, tw := range Towns {
				if math.Hypot(x-tw.X, z-tw.Z) < 8 {
					c = 'T'
				}
			}
			sb.WriteByte(c)
		}
		sb.WriteByte('\n')
	}
	fmt.Println(sb.String())
}

func TestPlacesAreOnLand(t *testing.T) {
	for _, n := range NPCs {
		if h := Height(n.X, n.Z); h < 0.4 {
			t.Errorf("NPC %s stands in water (h=%.2f at %.0f,%.0f)", n.ID, h, n.X, n.Z)
		}
	}
	for _, s := range AllSpawns() {
		if h := Height(s.X, s.Z); h < 0.4 {
			t.Errorf("spawn %s in water (h=%.2f at %.0f,%.0f)", s.Mob, h, s.X, s.Z)
		}
	}
	for _, p := range POIs {
		if h := Height(p.X, p.Z); h < -1.2 {
			t.Errorf("POI %s deep under water (h=%.2f)", p.ID, h)
		}
	}
	for _, tw := range Towns {
		if h := Height(tw.X, tw.Z); h < 0.5 {
			t.Errorf("town %s under water", tw.ID)
		}
	}
}

func TestZonesOfPlaces(t *testing.T) {
	want := map[string]Biome{
		"eichenhain": BiomeForest, "wolfshuegel": BiomeForest, "keilerwiese": BiomeForest, "spinnenwald": BiomeForest,
		"banditenlager": BiomeForest, "baerenklamm": BiomeForest, "sonnenburg": BiomePlains, "weizenhof": BiomePlains,
		"jaegerhain": BiomePlains, "quellheiligtum": BiomePlains, "raeuberfestung": BiomePlains, "moorbruck": BiomeSwamp,
		"krokodilsuempfe": BiomeSwamp, "hexenhain": BiomeSwamp, "schiefer_turm": BiomeSwamp, "frostwacht": BiomeSnow,
		"yetihoehle": BiomeSnow, "rabenfels": BiomeSnow, "sandfurt": BiomeDesert, "skorpionduenen": BiomeDesert,
		"pharaonengrab": BiomeDesert, "windmuehlenhuegel": BiomePlains, "baerenklippe": BiomeSnow, "windpass": BiomeSnow,
	}
	for id, b := range want {
		p := POIByID(id)
		if p == nil {
			t.Errorf("missing POI %s", id)
			continue
		}
		if got := BiomeAt(p.X, p.Z); got != b {
			t.Errorf("%s is in biome %s, want %s", id, Zones[got].Key, Zones[b].Key)
		}
	}
}
