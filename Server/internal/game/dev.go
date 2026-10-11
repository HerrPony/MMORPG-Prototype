package game

import (
	"strconv"

	"mmorpg/internal/gamedata"
	"mmorpg/internal/progression"
	"mmorpg/internal/protocol"
)

// devCommand implements the test commands enabled with the -dev flag:
// /dev tp x z | level n | gold n | unlock | heal
func (w *World) devCommand(e *Entity, a []string) {
	if len(a) == 0 {
		return
	}
	num := func(i int) float64 {
		if i < len(a) {
			f, _ := strconv.ParseFloat(a[i], 64)
			return f
		}
		return 0
	}
	switch a[0] {
	case "tp":
		if e.Inst != 0 {
			return
		}
		e.P.Aboard = 0
		e.X, e.Z = num(1), num(2)
		w.correctPlayer(e)
	case "level":
		lvl := max(1, min(progression.MaxLevel, int(num(1))))
		e.Level, e.P.Char.Level, e.P.Char.XP = lvl, lvl, 0
		w.recalc(e)
		e.HP = e.MaxHP
		w.sendEquipment(e)
		w.send(e, protocol.M{"t": "lvl", "lvl": lvl})
		w.sendMarks(e)
	case "gold":
		e.P.Char.Copper += int(num(1)) * 100
		w.sendInventory(e)
	case "unlock":
		for _, t := range gamedata.Towns {
			if t.Kind != "site" {
				e.P.Char.FlightPoints = append(e.P.Char.FlightPoints, t.ID)
			}
		}
		w.send(e, protocol.M{"t": "fp", "ids": e.P.Char.FlightPoints})
	case "heal":
		e.HP = e.MaxHP
	}
}
