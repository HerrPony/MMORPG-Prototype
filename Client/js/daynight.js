// Game clock and sun/moon position. One game day lasts 6 real hours (a quarter of a real day), derived from the
// server's wall clock so that every player sees the same time of day: game midnight is at real 00:00, 06:00,
// 12:00 and 18:00 server time.
export const DAY_MS = 6 * 3600 * 1000;

let offset = 0;       // server clock minus local clock
let forced = null;    // ?hour=14 forces a game hour (testing)
try { const h = new URLSearchParams(location.search).get('hour'); if (h !== null && h !== '') forced = ((+h % 24) + 24) % 24; } catch { /* ignore */ }

export function setServerTime(serverNowMs) { offset = serverNowMs - Date.now(); }

// 0..1 through the game day (0 = midnight, 0.25 = sunrise, 0.5 = noon, 0.75 = sunset)
export function dayPhase() {
  if (forced !== null) return forced / 24;
  return ((((Date.now() + offset) % DAY_MS) + DAY_MS) % DAY_MS) / DAY_MS;
}
export const gameHours = () => dayPhase() * 24;

const smooth = (a, b, x) => { const t = Math.max(0, Math.min(1, (x - a) / (b - a))); return t * t * (3 - 2 * t); };

// Everything the sky, lights and music need to know about the current time.
export function sunInfo(ph = dayPhase()) {
  const a = (ph - 0.25) * Math.PI * 2;                 // 0 at sunrise
  const sx = Math.cos(a), sy = Math.sin(a), sz = 0.35; // sun arc, tilted a bit to the south
  const l = Math.hypot(sx, sy, sz);
  const sun = [sx / l, sy / l, sz / l];
  const day = smooth(-0.14, 0.22, sy);                 // 0 night .. 1 full day
  return {
    sun, moon: [-sun[0], -sun[1], -sun[2]], elev: sy,
    day, night: 1 - day,
    twilight: Math.exp(-Math.pow(sy / 0.2, 2)),        // peaks when the sun is at the horizon
    isMorning: ph < 0.5,
  };
}
