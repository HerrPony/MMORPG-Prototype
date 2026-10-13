package gamedata

import (
	"math"
	"sync"
)

// Static obstacles of the open world (trees, rocks, buildings), used by the server to keep mobs from walking
// through them and to block shots without line of sight.
//
// The vegetation layout (props.go) is generated here, once, and sent to the clients, which render and collide
// with exactly this list. There is no second copy of the placement logic.

type solidCircle struct {
	x, z, r float64
	los     bool // thick enough to block line of sight (trees, rocks), not bushes
}
type solidBox struct{ x, z, hw, hd float64 }

var (
	collOnce  sync.Once
	circles   map[[2]int][]solidCircle // 8x8 cells
	boxes     []solidBox
	boxCircle map[[2]int][]int // 8x8 cells -> indices into boxes
)

func cellKey(x, z float64) [2]int { return [2]int{int(math.Floor(x / 8)), int(math.Floor(z / 8))} }

func initColliders() {
	collOnce.Do(func() {
		sc := baseScenery()
		circles = map[[2]int][]solidCircle{}
		boxCircle = map[[2]int][]int{}
		for _, b := range sc.House {
			boxes = append(boxes, solidBox{b.X, b.Z, b.W / 2, b.D / 2})
		}
		for _, b := range sc.Buildings {
			switch b.Kind {
			case "pad", "dock", "ship", "lantern", "pond", "gate", "archgate":
				continue
			case "tent":
				b.W, b.D = b.W*0.9, b.D*0.9
			}
			boxes = append(boxes, solidBox{b.X, b.Z, b.W / 2, b.D / 2})
		}
		for i, b := range boxes {
			for cx := int(math.Floor((b.x - b.hw - 1) / 8)); cx <= int(math.Floor((b.x+b.hw+1)/8)); cx++ {
				for cz := int(math.Floor((b.z - b.hd - 1) / 8)); cz <= int(math.Floor((b.z+b.hd+1)/8)); cz++ {
					boxCircle[[2]int{cx, cz}] = append(boxCircle[[2]int{cx, cz}], i)
				}
			}
		}
		for _, p := range worldProps(sc) {
			if p.Solid > 0 {
				k := cellKey(p.X, p.Z)
				circles[k] = append(circles[k], solidCircle{p.X, p.Z, p.Solid * p.S, p.Solid >= 0.5})
			}
		}
	})
}

// ---- queries ----

// Collide pushes a body of radius rad out of trees, rocks and buildings (open world only).
func Collide(x, z, rad float64) (float64, float64) {
	if z > DungeonZ0-100 {
		return x, z
	}
	initColliders()
	gx, gz := int(math.Floor(x/8)), int(math.Floor(z/8))
	for ix := gx - 1; ix <= gx+1; ix++ {
		for iz := gz - 1; iz <= gz+1; iz++ {
			for _, c := range circles[[2]int{ix, iz}] {
				dx, dz := x-c.x, z-c.z
				d := math.Hypot(dx, dz)
				if m := c.r + rad; d < m && d > 1e-4 {
					x, z = c.x+dx/d*m, c.z+dz/d*m
				}
			}
		}
	}
	for _, i := range boxCircle[[2]int{gx, gz}] {
		b := boxes[i]
		hx, hz := b.hw+rad, b.hd+rad
		dx, dz := x-b.x, z-b.z
		if math.Abs(dx) < hx && math.Abs(dz) < hz {
			px, pz := hx-math.Abs(dx), hz-math.Abs(dz)
			if px < pz {
				x = b.x + math.Copysign(hx, dx+1e-9)
			} else {
				z = b.z + math.Copysign(hz, dz+1e-9)
			}
		}
	}
	return x, z
}

// LineBlocked reports whether the straight line between two points is cut by a building (and, with trees,
// also by tree trunks and rocks). Used for line of sight.
func LineBlocked(x0, z0, x1, z1 float64, trees bool) bool {
	if z0 > DungeonZ0-100 {
		return false
	}
	initColliders()
	n := int(math.Ceil(math.Hypot(x1-x0, z1-z0) / 0.7))
	for i := 1; i < n; i++ { // end points themselves are never blocked
		t := float64(i) / float64(n)
		x, z := x0+(x1-x0)*t, z0+(z1-z0)*t
		k := cellKey(x, z)
		for _, bi := range boxCircle[k] {
			b := boxes[bi]
			if math.Abs(x-b.x) < b.hw && math.Abs(z-b.z) < b.hd {
				return true
			}
		}
		if !trees {
			continue
		}
		for cx := k[0] - 1; cx <= k[0]+1; cx++ {
			for cz := k[1] - 1; cz <= k[1]+1; cz++ {
				for _, c := range circles[[2]int{cx, cz}] {
					if c.los && math.Hypot(x-c.x, z-c.z) < c.r*0.6 {
						return true
					}
				}
			}
		}
	}
	return false
}

// VegetationCount is used by tests to compare against the client's layout.
func VegetationCount() (n int) {
	initColliders()
	for _, c := range circles {
		n += len(c)
	}
	return
}
