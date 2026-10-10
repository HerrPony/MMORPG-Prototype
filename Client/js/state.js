// Shared client state (single source of truth for UI + world modules).
export const S = {
  net: null,
  user: '',
  // static data from the server
  classes: {}, classOrder: [], races: {}, raceOrder: [], allSkills: {},
  defs: null,            // { items, quests, skills }
  scenery: null,
  // in world
  id: 0,
  self: null,            // { name, race, class, look }
  cls: null,             // own class definition
  skillIds: [],
  me: { hp: 1, mh: 1, res: 0, rm: 1, xp: 0, xn: 100, lvl: 1, tg: 0, auto: false, d: false },
  ents: new Map(),       // id -> entity record (see world/entities.js)
  inv: [], copper: 0, invMax: 28,
  equip: {}, stats: {}, sets: {}, // equipment by slot, gear totals, saved sets
  vendorNpc: 0, vendorData: null,
  disc: [], fp: [],       // discovered places, unlocked griffin landing points
  party: { id: 0, leader: 0, members: [], states: [] },
  zone: { name: '', sub: '', lo: 1, hi: 6 },
  instance: null,         // dungeon info while inside one
  quests: [],            // [{id,status,counts}]
  marks: {},             // npcId -> '' | available | progress | ready
  cd: {},                // skill -> {end, dur}
  gcd: { end: 0, dur: 0 },
  cast: null,            // {skill, start, dur}
  local: { x: 0, y: 0, z: 0, r: 0, speed: 0, airborne: false }, // own predicted state
  typing: false,         // chat input focused
  inWorld: false,
};

export const now = () => performance.now() / 1000;
