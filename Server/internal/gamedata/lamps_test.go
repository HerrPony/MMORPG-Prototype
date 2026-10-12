package gamedata

import (
	"math"
	"testing"
)

func TestLamps(t *testing.T) {
	l := Lamps()
	t.Logf("%d street lamps", len(l))
	if len(l) < 40 || len(l) > 900 {
		t.Fatalf("unexpected number of lamps: %d", len(l))
	}
	for i, a := range l {
		if Height(a.X, a.Z) < 0.9 {
			t.Errorf("lamp %d stands in the water (%.0f,%.0f)", i, a.X, a.Z)
		}
		for _, b := range l[i+1:] {
			if math.Hypot(a.X-b.X, a.Z-b.Z) < lampMinGap-0.5 {
				t.Errorf("lamps too close: %.0f,%.0f and %.0f,%.0f", a.X, a.Z, b.X, b.Z)
			}
		}
	}
}
