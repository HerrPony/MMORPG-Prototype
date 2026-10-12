package gamedata

import "testing"

// Swimming must never connect the two continents: from the first continent's coast the swimmable
// water must not reach the second one (the ferry is the only way across).
func TestSwimDoesNotCrossSea(t *testing.T) {
	ensureWorld()
	idx := func(i, j int) int { return j*GridNX + i }
	ok := func(i, j int) bool {
		x, z := gridXZ(i, j)
		return Swimmable(x, z)
	}
	seen := make([]bool, GridNX*GridNZ)
	var q [][2]int
	for j := 0; j < GridNZ; j++ {
		for i := 0; i < GridNX; i++ {
			x, z := gridXZ(i, j)
			if x < ContSplitX && Walkable(x, z) && !seen[idx(i, j)] && x > -50 && x < 50 && z > -50 && z < 50 {
				seen[idx(i, j)] = true
				q = append(q, [2]int{i, j})
			}
		}
	}
	minX := 1e9
	for len(q) > 0 {
		c := q[0]
		q = q[1:]
		x, _ := gridXZ(c[0], c[1])
		if x > ContSplitX {
			t.Fatalf("swimming leads to the second continent at x=%.0f", x)
		}
		_ = minX
		for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			i, j := c[0]+d[0], c[1]+d[1]
			if i < 0 || j < 0 || i >= GridNX || j >= GridNZ || seen[idx(i, j)] || !ok(i, j) {
				continue
			}
			seen[idx(i, j)] = true
			q = append(q, [2]int{i, j})
		}
	}
	// report how far one can get from the coast of the first continent
	maxX := -1e9
	for j := 0; j < GridNZ; j++ {
		for i := 0; i < GridNX; i++ {
			if seen[idx(i, j)] {
				if x, _ := gridXZ(i, j); x > maxX {
					maxX = x
				}
			}
		}
	}
	t.Logf("easternmost swimmable x of continent 1: %.0f (split at %.0f)", maxX, ContSplitX)
}

// Every harbor must offer swimmable water in front of the pier.
func TestHarborHasSwimWater(t *testing.T) {
	for _, tw := range Towns {
		if !tw.Harbor {
			continue
		}
		bx, bz := Berth(tw)
		if !Swimmable(bx, bz) {
			t.Errorf("berth of %s is not swimmable", tw.ID)
		}
	}
}
