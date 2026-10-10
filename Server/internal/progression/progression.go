// Package progression contains XP curves and stat scaling.
package progression

import "mmorpg/internal/gamedata"

const MaxLevel = gamedata.MaxLevel

// XPToNext is the XP needed to advance from the given level.
func XPToNext(level int) int { return gamedata.XPToNext(level) }

// MaxHP of a class at a level.
func MaxHP(c *gamedata.Class, level int) int {
	return c.HP + c.HPPerLevel*(level-1)
}

// MaxResource of a class at a level (mana/focus grow a little, rage is fixed).
func MaxResource(c *gamedata.Class, level int) int {
	if c.Resource == "rage" {
		return c.ResMax
	}
	return c.ResMax + 8*(level-1)
}

// DamageScale multiplies skill power with the character level.
func DamageScale(level int) float64 { return 1 + 0.15*float64(level-1) }

// KillXP is the XP a mob of mobLevel grants a player of playerLevel.
func KillXP(base, mobLevel, playerLevel int) int {
	xp := float64(base + 10*mobLevel)
	f := 1 + 0.1*float64(mobLevel-playerLevel)
	if f < 0.1 {
		f = 0.1
	}
	if f > 1.5 {
		f = 1.5
	}
	return int(xp * f)
}

// AddXP adds xp and returns the new level, remaining xp and number of level-ups.
func AddXP(level, xp, gain int) (newLevel, newXP, ups int) {
	xp += gain
	for level < MaxLevel {
		need := XPToNext(level)
		if xp < need {
			break
		}
		xp -= need
		level++
		ups++
	}
	if level >= MaxLevel {
		xp = 0
	}
	return level, xp, ups
}
