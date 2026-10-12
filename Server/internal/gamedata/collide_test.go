package gamedata

import (
	"encoding/json"
	"os"
	"testing"
)

// TestVegetationCount prints the obstacle count; with DUMP_SCENERY=<dir> it also writes the scenery JSON and
// the solid list so the client layout can be compared (see Client dev notes).
func TestVegetationCount(t *testing.T) {
	n := VegetationCount()
	t.Logf("solid vegetation: %d, boxes: %d", n, len(boxes))
	if n < 1000 {
		t.Fatal("suspiciously few obstacles")
	}
	if dir := os.Getenv("DUMP_SCENERY"); dir != "" {
		b, _ := json.Marshal(BuildScenery())
		os.WriteFile(dir+"/scenery.json", b, 0o644)
		var out [][3]float64
		for _, cs := range circles {
			for _, c := range cs {
				out = append(out, [3]float64{c.x, c.z, c.r})
			}
		}
		b, _ = json.Marshal(out)
		os.WriteFile(dir+"/solids.json", b, 0o644)
	}
}
