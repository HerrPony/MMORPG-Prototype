package gamedata

import (
	"math"
	"sync"
)

// Lamp is a street lamp: a post with a lantern that lights up at night. The server places them once (along the
// roads and in front of houses) so that every client shows the same lamps.
type Lamp struct {
	X     float64 `json:"x"`
	Z     float64 `json:"z"`
	Style string  `json:"s,omitempty"` // building style of the nearest town (paper lanterns in jade towns ...)
}

const (
	lampSpacing = 26.0 // along roads
	lampOffset  = 2.8  // distance from the road centre
	lampMinGap  = 11.0 // no two lamps closer than this
)

var (
	lampsOnce sync.Once
	lampList  []Lamp
)

// Lamps returns all street lamps (generated once, deterministic).
func Lamps() []Lamp {
	lampsOnce.Do(func() {
		blds := Buildings()
		free := func(x, z float64) bool {
			if Height(x, z) < 0.9 || BiomeAt(x, z) == BiomeLava || BiomeAt(x, z) == BiomeSea {
				return false
			}
			for _, b := range blds {
				if math.Abs(x-b.X) < b.W/2+1.2 && math.Abs(z-b.Z) < b.D/2+1.2 {
					return false
				}
			}
			for _, l := range lampList {
				if math.Hypot(x-l.X, z-l.Z) < lampMinGap {
					return false
				}
			}
			return true
		}
		add := func(x, z float64) {
			if !free(x, z) {
				return
			}
			style := ""
			best := 80.0
			for _, t := range Towns {
				if d := math.Hypot(x-t.X, z-t.Z); d < best {
					best, style = d, t.Style
				}
			}
			lampList = append(lampList, Lamp{X: math.Round(x*10) / 10, Z: math.Round(z*10) / 10, Style: style})
		}
		// along the roads, alternating sides
		side := 1.0
		for _, r := range Roads {
			next := lampSpacing / 2
			walked := 0.0
			for i := 0; i+1 < len(r.Pts); i++ {
				a, b := r.Pts[i], r.Pts[i+1]
				seg := math.Hypot(b[0]-a[0], b[1]-a[1])
				if seg < 1e-6 {
					continue
				}
				dx, dz := (b[0]-a[0])/seg, (b[1]-a[1])/seg
				for walked+seg >= next {
					t := next - walked
					add(a[0]+dx*t-dz*lampOffset*side, a[1]+dz*t+dx*lampOffset*side)
					side = -side
					next += lampSpacing
				}
				walked += seg
			}
		}
		// in front of the houses (the door is on the +z side)
		for _, b := range blds {
			if b.Kind == "house" || b.Kind == "hut" {
				add(b.X+b.W/2+1.1, b.Z+b.D/2+1.5)
			}
		}
	})
	return lampList
}
