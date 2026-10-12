package gamedata

import (
	"math"
	"testing"
	"time"
)

func TestPathAroundCabin(t *testing.T) {
	// the cabin (10x8 m) at the world origin lies between the two points
	if NavLineFree(-14.5, -3.5, 14, 0) {
		t.Fatal("the straight line through the cabin must be blocked")
	}
	start := time.Now()
	p, partial, ok := FindPath(-14.5, -3.5, 14, 0, 20000)
	if !ok || partial || len(p) < 2 {
		t.Fatalf("no path around the cabin: ok=%v partial=%v %v", ok, partial, p)
	}
	t.Logf("path with %d waypoints in %v: %v", len(p), time.Since(start), p)
	x, z, length := -14.5, -3.5, 0.0
	for _, wp := range p {
		if !NavLineFree(x, z, wp[0], wp[1]) {
			t.Fatalf("segment %.1f,%.1f -> %.1f,%.1f cuts through an obstacle", x, z, wp[0], wp[1])
		}
		length += math.Hypot(wp[0]-x, wp[1]-z)
		x, z = wp[0], wp[1]
	}
	if math.Hypot(x-14, z) > 1.5 {
		t.Fatalf("path ends at %.1f,%.1f", x, z)
	}
	if length < 28 || length > 50 {
		t.Fatalf("implausible path length %.1f (straight line is 30)", length)
	}
}

func TestPathIsFastAndHandlesBlockedGoal(t *testing.T) {
	FindPath(-14.5, -3.5, 14, 0, 20000) // warm up (builds the grid)
	start := time.Now()
	for i := 0; i < 50; i++ {
		FindPath(-40, -30+float64(i), 60, 30, 20000)
	}
	if d := time.Since(start) / 50; d > 20*time.Millisecond {
		t.Fatalf("path search too slow: %v per search", d)
	}
	// a goal inside the cabin walls is moved to the nearest free spot
	if p, _, ok := FindPath(-15, 0, 0, 0, 20000); !ok || len(p) == 0 {
		t.Fatal("goal inside an obstacle should lead to its edge")
	}
}
