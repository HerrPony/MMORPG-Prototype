package gamedata

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGeoDump(t *testing.T) {
	if os.Getenv("GEO") == "" {
		t.Skip("set GEO=1 to print the world map")
	}
	ch := map[Biome]byte{BiomeSea: '~', BiomeForest: 'f', BiomePlains: 'p', BiomeSwamp: 's', BiomeSnow: 'n', BiomeDesert: 'd', BiomeJungle: 'J', BiomeBlossom: 'B', BiomeMesa: 'M', BiomeVolcano: 'V', BiomeLava: 'L'}
	var sb strings.Builder
	for z := -296.0; z < 300; z += 10 {
		for x := -300.0; x < 1140; x += 10 {
			c := ch[BiomeAt(x, z)]
			if c != '~' && !Walkable(x, z) {
				c = '#'
			}
			for _, tw := range Towns {
				if abs(tw.X-x) < 5 && abs(tw.Z-z) < 5 {
					c = '@'
				}
			}
			sb.WriteByte(c)
		}
		sb.WriteByte('\n')
	}
	fmt.Println(sb.String())
	for _, tw := range Towns {
		if tw.Harbor {
			bx, bz := Berth(tw)
			px, pz := PadPos(tw)
			fmt.Printf("HARBOR %s berth=(%.0f,%.0f) h=%.1f pad=(%.0f,%.0f) h=%.1f\n", tw.ID, bx, bz, Height(bx, bz), px, pz, Height(px, pz))
			for x := px; x < px+60 && x > px-60; x += map[bool]float64{true: 4, false: -4}[tw.X < 400] {
				fmt.Printf(" %.0f:%.1f", x, Height(x, pz))
			}
			fmt.Println()
		}
		px, pz := PadPos(tw)
		fmt.Printf("%-14s h=%.1f biome=%d walk=%v pad=(%.0f,%.0f) walk=%v\n", tw.ID, Height(tw.X, tw.Z), BiomeAt(tw.X, tw.Z), Walkable(tw.X, tw.Z), px, pz, Walkable(px, pz))
	}
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
