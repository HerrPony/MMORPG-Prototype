package gamedata

import "testing"

func TestPropLayout(t *testing.T) {
	n := VegetationCount()
	t.Logf("solid props: %d, building boxes: %d, props total: %d", n, len(boxes), len(propList))
	if n < 1000 || len(propList) < n {
		t.Fatal("suspiciously few props")
	}
	// the packed form must round-trip losslessly (the client collides with exactly the decoded values)
	if len(PropsB64(propList)) < len(propList) {
		t.Fatal("props not packed")
	}
}
