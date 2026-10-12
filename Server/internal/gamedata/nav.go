package gamedata

import (
	"math"
	"sync"
)

// Navigation grid for the open world: 1 m cells that are either free or blocked (deep water, lava, cliffs,
// trees, rocks, buildings - inflated by the body radius of a mob). It is derived from the same data as the
// collision (terrain + props.go + buildings), so there is nothing to maintain by hand. FindPath runs A* over
// it and returns smoothed waypoints.

const (
	navCell   = 1.0
	navAgentR = 0.55 // body radius the grid is inflated by
	navMaxSlp = 1.4  // steeper than this (m rise per m) is a cliff
)

var (
	navOnce sync.Once
	navNX   int
	navNZ   int
	navBlk  []uint8 // 1 = blocked

	// A* scratch buffers, reused between searches (the world goroutine is the only caller; the mutex keeps tests safe)
	navMu     sync.Mutex
	navStamp  []uint32
	navClosed []uint32
	navG      []float32
	navPar    []int32
	navGen    uint32
)

func resetNav() {
	navOnce = sync.Once{}
	navBlk = nil
	navStamp = nil
	navClosed = nil
	navG = nil
	navPar = nil
}

func initNav() {
	initColliders()
	navOnce.Do(func() {
		ensureWorld()
		navNX, navNZ = int((WorldMaxX-WorldMinX)/navCell), int((WorldMaxZ-WorldMinZ)/navCell)
		navBlk = make([]uint8, navNX*navNZ)
		for j := 0; j < navNZ; j++ {
			for i := 0; i < navNX; i++ {
				x, z := WorldMinX+(float64(i)+0.5)*navCell, WorldMinZ+(float64(j)+0.5)*navCell
				if !Walkable(x, z) {
					navBlk[j*navNX+i] = 1
					continue
				}
				h := Height(x, z)
				if math.Abs(Height(x+navCell, z)-h) > navMaxSlp*navCell || math.Abs(Height(x, z+navCell)-h) > navMaxSlp*navCell {
					navBlk[j*navNX+i] = 1
				}
			}
		}
		for _, cs := range circles {
			for _, c := range cs {
				navFillCircle(c.x, c.z, c.r+navAgentR)
			}
		}
		for _, b := range boxes {
			navFillBox(b.x-b.hw-navAgentR, b.z-b.hd-navAgentR, b.x+b.hw+navAgentR, b.z+b.hd+navAgentR)
		}
		navStamp = make([]uint32, navNX*navNZ)
		navClosed = make([]uint32, navNX*navNZ)
		navG = make([]float32, navNX*navNZ)
		navPar = make([]int32, navNX*navNZ)
		navGen = 0
	})
}

func navIdx(x, z float64) (int, int, bool) {
	i, j := int(math.Floor((x-WorldMinX)/navCell)), int(math.Floor((z-WorldMinZ)/navCell))
	return i, j, i >= 0 && j >= 0 && i < navNX && j < navNZ
}

func navFillCircle(cx, cz, r float64) {
	i0, j0, _ := navIdxClamp(cx-r, cz-r)
	i1, j1, _ := navIdxClamp(cx+r, cz+r)
	for j := j0; j <= j1; j++ {
		for i := i0; i <= i1; i++ {
			x, z := WorldMinX+(float64(i)+0.5)*navCell, WorldMinZ+(float64(j)+0.5)*navCell
			if math.Hypot(x-cx, z-cz) <= r {
				navBlk[j*navNX+i] = 1
			}
		}
	}
}

func navFillBox(x0, z0, x1, z1 float64) {
	i0, j0, _ := navIdxClamp(x0, z0)
	i1, j1, _ := navIdxClamp(x1, z1)
	for j := j0; j <= j1; j++ {
		for i := i0; i <= i1; i++ {
			x, z := WorldMinX+(float64(i)+0.5)*navCell, WorldMinZ+(float64(j)+0.5)*navCell
			if x >= x0 && x <= x1 && z >= z0 && z <= z1 {
				navBlk[j*navNX+i] = 1
			}
		}
	}
}

func navIdxClamp(x, z float64) (int, int, bool) {
	i, j := int(math.Floor((x-WorldMinX)/navCell)), int(math.Floor((z-WorldMinZ)/navCell))
	return max(0, min(navNX-1, i)), max(0, min(navNZ-1, j)), true
}

// NavBlocked reports whether a mob cannot stand at x,z.
func NavBlocked(x, z float64) bool {
	initNav()
	i, j, ok := navIdx(x, z)
	return !ok || navBlk[j*navNX+i] != 0
}

// NavLineFree reports whether a mob can walk the straight line without touching an obstacle.
func NavLineFree(x0, z0, x1, z1 float64) bool {
	if z0 > DungeonZ0-100 {
		return true
	}
	initNav()
	n := int(math.Ceil(math.Hypot(x1-x0, z1-z0) / (navCell * 0.5)))
	for k := 1; k <= n; k++ {
		t := float64(k) / float64(n)
		i, j, ok := navIdx(x0+(x1-x0)*t, z0+(z1-z0)*t)
		if !ok || navBlk[j*navNX+i] != 0 {
			return false
		}
	}
	return true
}

// nearestFree finds the free cell closest to a (possibly blocked) position, searching rings up to maxR cells.
func nearestFree(x, z float64, maxR int) (int, int, bool) {
	ci, cj, ok := navIdx(x, z)
	if !ok {
		return 0, 0, false
	}
	if navBlk[cj*navNX+ci] == 0 {
		return ci, cj, true
	}
	best, bi, bj := math.MaxFloat64, 0, 0
	for r := 1; r <= maxR; r++ {
		for dj := -r; dj <= r; dj++ {
			for di := -r; di <= r; di++ {
				if max(iabs(di), iabs(dj)) != r {
					continue
				}
				i, j := ci+di, cj+dj
				if i < 0 || j < 0 || i >= navNX || j >= navNZ || navBlk[j*navNX+i] != 0 {
					continue
				}
				if d := float64(di*di + dj*dj); d < best {
					best, bi, bj = d, i, j
				}
			}
		}
		if best < math.MaxFloat64 {
			return bi, bj, true
		}
	}
	return 0, 0, false
}

func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type navNode struct {
	f   float32
	idx int32
}

// FindPath returns smoothed waypoints from (x0,z0) to (x1,z1) around obstacles. ok is false if no route exists
// within maxExpand search steps; when the search is cut short the path leads to the closest point reached
// (partial = true), which is the right thing to follow while chasing.
func FindPath(x0, z0, x1, z1 float64, maxExpand int) (path [][2]float64, partial, ok bool) {
	if z0 > DungeonZ0-100 {
		return nil, false, false
	}
	initNav()
	navMu.Lock()
	defer navMu.Unlock()
	si, sj, ok1 := nearestFree(x0, z0, 4)
	gi, gj, ok2 := nearestFree(x1, z1, 6)
	if !ok1 || !ok2 {
		return nil, false, false
	}
	start, goal := int32(sj*navNX+si), int32(gj*navNX+gi)
	navGen++
	if navGen == 0 { // wrapped: clear the stamps
		for i := range navStamp {
			navStamp[i], navClosed[i] = 0, 0
		}
		navGen = 1
	}
	heur := func(i, j int) float32 {
		dx, dz := float32(iabs(i-gi)), float32(iabs(j-gj))
		if dx < dz {
			dx, dz = dz, dx
		}
		return dx + 0.41421356*dz
	}
	var open []navNode
	push := func(n navNode) {
		open = append(open, n)
		for c := len(open) - 1; c > 0; {
			p := (c - 1) / 2
			if open[p].f <= open[c].f {
				break
			}
			open[p], open[c] = open[c], open[p]
			c = p
		}
	}
	pop := func() navNode {
		top := open[0]
		last := len(open) - 1
		open[0] = open[last]
		open = open[:last]
		for c := 0; ; {
			l, r, m := 2*c+1, 2*c+2, c
			if l < len(open) && open[l].f < open[m].f {
				m = l
			}
			if r < len(open) && open[r].f < open[m].f {
				m = r
			}
			if m == c {
				break
			}
			open[m], open[c] = open[c], open[m]
			c = m
		}
		return top
	}
	navStamp[start], navG[start], navPar[start] = navGen, 0, -1
	push(navNode{heur(si, sj), start})
	best, bestH := start, heur(si, sj)
	found := false
	for expanded := 0; len(open) > 0 && expanded < maxExpand; {
		cur := pop()
		if navClosed[cur.idx] == navGen {
			continue
		}
		navClosed[cur.idx] = navGen
		expanded++
		if cur.idx == goal {
			found = true
			break
		}
		ci, cj := int(cur.idx)%navNX, int(cur.idx)/navNX
		if h := heur(ci, cj); h < bestH {
			best, bestH = cur.idx, h
		}
		for dj := -1; dj <= 1; dj++ {
			for di := -1; di <= 1; di++ {
				if di == 0 && dj == 0 {
					continue
				}
				ni, nj := ci+di, cj+dj
				if ni < 0 || nj < 0 || ni >= navNX || nj >= navNZ || navBlk[nj*navNX+ni] != 0 {
					continue
				}
				cost := float32(1)
				if di != 0 && dj != 0 { // no cutting corners
					if navBlk[cj*navNX+ni] != 0 || navBlk[nj*navNX+ci] != 0 {
						continue
					}
					cost = 1.41421356
				}
				ni32 := int32(nj*navNX + ni)
				g := navG[cur.idx] + cost
				if navStamp[ni32] == navGen && navG[ni32] <= g {
					continue
				}
				navStamp[ni32], navG[ni32], navPar[ni32] = navGen, g, cur.idx
				push(navNode{g + heur(ni, nj), ni32})
			}
		}
	}
	end := goal
	if !found {
		if best == start {
			return nil, false, false
		}
		end, partial = best, true
	}
	var cells [][2]float64
	for c := end; c != -1; c = navPar[c] {
		cells = append(cells, [2]float64{WorldMinX + (float64(int(c)%navNX)+0.5)*navCell, WorldMinZ + (float64(int(c)/navNX)+0.5)*navCell})
	}
	for a, b := 0, len(cells)-1; a < b; a, b = a+1, b-1 {
		cells[a], cells[b] = cells[b], cells[a]
	}
	if !partial { // finish exactly on the target when it is reachable in a straight line from the last cell
		if i, j, ok := navIdx(x1, z1); ok && navBlk[j*navNX+i] == 0 {
			cells = append(cells, [2]float64{x1, z1})
		}
	}
	// string pulling: skip waypoints that can be bypassed in a straight line
	out := [][2]float64{}
	from := [2]float64{x0, z0}
	for k := 0; k < len(cells); {
		far := k
		for m := len(cells) - 1; m > k; m-- {
			if navLineFreeLocked(from[0], from[1], cells[m][0], cells[m][1]) {
				far = m
				break
			}
		}
		out = append(out, cells[far])
		from = cells[far]
		k = far + 1
	}
	return out, partial, true
}

func navLineFreeLocked(x0, z0, x1, z1 float64) bool {
	n := int(math.Ceil(math.Hypot(x1-x0, z1-z0) / (navCell * 0.5)))
	for k := 1; k <= n; k++ {
		t := float64(k) / float64(n)
		i, j, ok := navIdx(x0+(x1-x0)*t, z0+(z1-z0)*t)
		if !ok || navBlk[j*navNX+i] != 0 {
			return false
		}
	}
	return true
}
