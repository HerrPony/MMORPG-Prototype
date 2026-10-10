package gamedata

const MaxLevel = 20

// XPToNext is the XP needed to advance from the given level.
func XPToNext(level int) int {
	if level >= MaxLevel {
		return 0
	}
	return 80*level + 20*level*level
}

// questXP is a quest reward: share of one level, scaled by weight.
func questXP(level int, weight float64) int {
	return int(float64(XPToNext(level)) * 0.075 * weight)
}

// questCopper is a quest money reward.
func questCopper(level int, weight float64) int {
	return int(float64(10+level*level*2+level*12) * weight)
}

// Mob formulas for the generated mobs.
func mobHP(level int) float64  { return float64(40 + 22*level + level*level) }
func mobDmg(level int) float64 { return float64(2) + 1.5*float64(level) }
