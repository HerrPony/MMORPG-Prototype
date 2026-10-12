class_name DayNight
extends RefCounted
## Game clock and sun/moon: one game day = 6 real hours, taken from the server's wall clock.

const DAY_MS := 6.0 * 3600.0 * 1000.0
static var forced := -1.0  # forced game hour (debug), -1 = off

static func now_ms() -> float:
	return Time.get_unix_time_from_system() * 1000.0 + S.server_offset_ms

static func phase() -> float:
	if forced >= 0.0:
		return forced / 24.0
	return fposmod(now_ms(), DAY_MS) / DAY_MS

static func hours() -> float:
	return phase() * 24.0

static func smooth(a: float, b: float, x: float) -> float:
	var t := clampf((x - a) / (b - a), 0.0, 1.0)
	return t * t * (3.0 - 2.0 * t)

## Returns sun direction (towards the sun), moon direction, elevation, day (0..1), night, twilight.
static func info(ph: float = -1.0) -> Dictionary:
	if ph < 0.0:
		ph = phase()
	var a := (ph - 0.25) * TAU
	var sun := Vector3(cos(a), sin(a), 0.35).normalized()
	var sy := sin(a)
	var day := smooth(-0.14, 0.22, sy)
	return {"sun": sun, "moon": -sun, "elev": sy, "day": day, "night": 1.0 - day, "twilight": exp(-pow(sy / 0.2, 2.0))}
