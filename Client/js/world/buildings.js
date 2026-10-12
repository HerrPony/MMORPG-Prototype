import * as THREE from 'three';
import { heightAt } from '../util.js';
import { makeShip } from './ship.js';

// window glass: shared, so that all of them can light up at night (see lamps.js)
export const windowMats = [];
const WIN = new THREE.MeshLambertMaterial({ color: 0xffe9a0, emissive: 0xffc870, emissiveIntensity: 0.1 });
windowMats.push(WIN);

const M = (c, o = {}) => new THREE.MeshLambertMaterial({ color: c, ...o });
const box = (w, h, d, m) => { const x = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), m); x.castShadow = x.receiveShadow = true; return x; };
const cyl = (rt, rb, h, m, s = 12) => { const x = new THREE.Mesh(new THREE.CylinderGeometry(rt, rb, h, s), m); x.castShadow = true; return x; };

const STYLES = {
  timber: { wall: 0x9a6a3a, roof: 0x7a2f22, trim: 0x5a3a1c },
  plaster: { wall: 0xe8dcc0, roof: 0xb04a2a, trim: 0x6a4a2c },
  stilt: { wall: 0x5a4a32, roof: 0x6a7a3a, trim: 0x3a2e1e },
  lodge: { wall: 0x6a4a30, roof: 0xdde8f0, trim: 0x3a2a1c },
  adobe: { wall: 0xd8b878, roof: 0xc89858, trim: 0x8a6a3a },
  bamboo: { wall: 0xcdbc7a, roof: 0x9c8a3c, trim: 0x7a6a30 },
  jade: { wall: 0xb8362a, roof: 0x2f9a78, trim: 0xe8c860 },
  redrock: { wall: 0xc8693f, roof: 0xb4583a, trim: 0x8a4628 },
  obsidian: { wall: 0x2c2c38, roof: 0x14141c, trim: 0xff6a20 },
};

function house(b, hh) {
  const st = STYLES[b.s] || STYLES.timber;
  const g = new THREE.Group();
  const lift = b.k === 'hut' ? 1.3 : 0;
  const walls = box(b.w, b.h, b.d, M(st.wall));
  walls.position.y = lift + b.h / 2;
  g.add(walls);
  // wooden corner posts and a base
  const trim = M(st.trim);
  for (const sx of [-1, 1]) for (const sz of [-1, 1]) {
    const p = box(0.35, b.h, 0.35, trim); p.position.set(sx * b.w / 2, lift + b.h / 2, sz * b.d / 2); g.add(p);
  }
  if (b.k === 'hut') for (const sx of [-1, 1]) for (const sz of [-1, 1]) {
    const leg = cyl(0.2, 0.25, lift + 0.8, trim, 6); leg.position.set(sx * b.w / 2.3, (lift + 0.8) / 2 - 0.4, sz * b.d / 2.3); g.add(leg);
  }
  if (b.s === 'adobe') {
    const top = box(b.w + 0.4, 0.5, b.d + 0.4, M(st.roof)); top.position.y = lift + b.h + 0.25; g.add(top);
    const dome = new THREE.Mesh(new THREE.SphereGeometry(Math.min(b.w, b.d) * 0.28, 12, 8, 0, 6.283, 0, 1.57), M(0xe8d8a8));
    dome.position.y = lift + b.h + 0.5; g.add(dome);
  } else if (b.s === 'redrock') { // flat clay roof with beam ends
    const top = box(b.w + 0.6, 0.6, b.d + 0.6, M(st.roof)); top.position.y = lift + b.h + 0.3; g.add(top);
    const rim = box(b.w + 0.9, 0.35, b.d + 0.9, M(st.trim)); rim.position.y = lift + b.h + 0.75; g.add(rim);
    for (let i = -2; i <= 2; i++) { const beam = cyl(0.14, 0.14, b.d + 1.6, M(0x5a3a22), 6); beam.rotation.x = Math.PI / 2; beam.position.set(i * b.w / 5, lift + b.h - 0.2, 0); g.add(beam); }
    const step = box(b.w * 0.45, 0.4, b.d * 0.4, M(st.wall)); step.position.set(0, lift + b.h + 0.95, 0); g.add(step);
  } else if (b.s === 'jade') { // two-tier roof with upturned corners
    const r1 = new THREE.Mesh(new THREE.ConeGeometry(1, 1.5, 4), M(st.roof));
    r1.rotation.y = Math.PI / 4; r1.scale.set((b.w / 2 + 1.6) * Math.SQRT2, 1, (b.d / 2 + 1.6) * Math.SQRT2); r1.position.y = lift + b.h + 0.75; r1.castShadow = true; g.add(r1);
    const r2 = new THREE.Mesh(new THREE.ConeGeometry(1, 1.5, 4), M(st.roof));
    r2.rotation.y = Math.PI / 4; r2.scale.set((b.w / 2 + 0.2) * Math.SQRT2 * 0.7, 1, (b.d / 2 + 0.2) * Math.SQRT2 * 0.7); r2.position.y = lift + b.h + 2.0; g.add(r2);
    for (const sx of [-1, 1]) for (const sz of [-1, 1]) { const tip = box(0.4, 0.4, 0.4, M(st.trim)); tip.position.set(sx * (b.w / 2 + 1.5), lift + b.h + 0.55, sz * (b.d / 2 + 1.5)); tip.rotation.z = sx * 0.6; g.add(tip); }
    const fin = box(0.2, 0.9, 0.2, M(st.trim)); fin.position.y = lift + b.h + 3.0; g.add(fin);
  } else if (b.s === 'obsidian') { // sharp black roof with glowing seams
    const rf = new THREE.Mesh(new THREE.ConeGeometry(1, 3.6, 5), M(st.roof));
    rf.scale.set(b.w / 2 + 0.6, 1, b.d / 2 + 0.6); rf.position.y = lift + b.h + 1.8; rf.castShadow = true; g.add(rf);
    const seam = box(b.w + 0.2, 0.2, b.d + 0.2, new THREE.MeshBasicMaterial({ color: 0xff6a20 })); seam.position.y = lift + b.h + 0.1; g.add(seam);
  } else if (b.s === 'bamboo') { // thatched cone roof
    const rf = new THREE.Mesh(new THREE.ConeGeometry(1, 3.2, 8), M(st.roof));
    rf.scale.set((Math.max(b.w, b.d) / 2 + 1.2), 1, (Math.max(b.w, b.d) / 2 + 1.2)); rf.position.y = lift + b.h + 1.6; rf.castShadow = true; g.add(rf);
  } else {
    const long = b.w >= b.d;
    const roof = new THREE.Mesh(new THREE.ConeGeometry(1, 2.4 + (b.h > 5 ? 0.5 : 0), 4), M(st.roof));
    roof.rotation.y = Math.PI / 4;
    roof.scale.set((b.w / 2 + 0.7) * Math.SQRT2, 1, (b.d / 2 + 0.7) * Math.SQRT2);
    roof.position.y = lift + b.h + 1.2;
    roof.castShadow = true;
    g.add(roof);
    if (long) roof.scale.z *= 1; else roof.scale.x *= 1;
  }
  // door and windows on the south side
  const door = box(1.3, 2.2, 0.2, M(0x3a2410)); door.position.set(0, lift + 1.1, b.d / 2 + 0.05); g.add(door);
  const wm = b.s === 'obsidian' ? new THREE.MeshBasicMaterial({ color: 0xff7a30 }) : WIN;
  for (const sx of [-1, 1]) { const w = box(0.9, 0.9, 0.2, wm); w.position.set(sx * b.w * 0.28, lift + b.h * 0.6, b.d / 2 + 0.05); g.add(w); }
  if (b.k === 'hut') { const st2 = box(2.2, 0.3, 2.0, trim); st2.position.set(0, 0.5, b.d / 2 + 1.3); g.add(st2); }
  g.position.set(b.x, hh, b.z);
  return g;
}

function tower(b, hh) {
  const g = new THREE.Group();
  const stone = M(0x9a9690);
  const body = cyl(b.w / 2, b.w / 2 + 0.5, b.h, stone, 14); body.position.y = b.h / 2; g.add(body);
  const roof = new THREE.Mesh(new THREE.ConeGeometry(b.w / 2 + 0.9, 4, 14), M(0x3a5a9a)); roof.position.y = b.h + 2; roof.castShadow = true; g.add(roof);
  const crown = cyl(b.w / 2 + 0.4, b.w / 2 + 0.4, 0.8, stone, 14); crown.position.y = b.h - 0.4; g.add(crown);
  g.position.set(b.x, hh, b.z);
  return g;
}

function wall(b, hh) {
  const g = new THREE.Group();
  const m = box(b.w, b.h, b.d, M(0xaaa69e)); m.position.y = b.h / 2; g.add(m);
  // crenellations
  const long = b.w > b.d;
  const n = Math.floor((long ? b.w : b.d) / 3);
  for (let i = 0; i < n; i++) {
    const c = box(long ? 1.4 : b.d, 1.2, long ? b.d : 1.4, M(0xaaa69e));
    const off = -((long ? b.w : b.d) / 2) + 1.5 + i * 3;
    c.position.set(long ? off : 0, b.h + 0.6, long ? 0 : off);
    g.add(c);
  }
  g.position.set(b.x, hh, b.z);
  return g;
}

function keep(b, hh) {
  const g = new THREE.Group();
  const stone = M(b.s === 'adobe' ? 0xd8b878 : b.s === 'redrock' ? 0xc8693f : 0xcfc9bc);
  const body = box(b.w, b.h, b.d, stone); body.position.y = b.h / 2; g.add(body);
  const roof = new THREE.Mesh(new THREE.ConeGeometry(1, 5, 4), M(b.s === 'adobe' ? 0xc89858 : b.s === 'redrock' ? 0xa84e2e : 0x3a5a9a));
  roof.rotation.y = Math.PI / 4; roof.scale.set(b.w / 2 * Math.SQRT2, 1, b.d / 2 * Math.SQRT2); roof.position.y = b.h + 2.5; roof.castShadow = true; g.add(roof);
  if (b.s !== 'adobe' && b.s !== 'redrock') for (const sx of [-1, 1]) for (const sz of [-1, 1]) {
    const t = cyl(2.2, 2.4, b.h + 4, stone, 12); t.position.set(sx * b.w / 2, (b.h + 4) / 2, sz * b.d / 2); g.add(t);
    const tr = new THREE.Mesh(new THREE.ConeGeometry(3, 4, 12), M(0xc8a24a)); tr.position.set(sx * b.w / 2, b.h + 6, sz * b.d / 2); g.add(tr);
  }
  const gate = box(4, 6, 0.4, M(0x3a2410)); gate.position.set(0, 3, b.d / 2 + 0.1); g.add(gate);
  for (let i = -2; i <= 2; i++) {
    const f = box(1.2, 3, 0.1, M(i % 2 ? 0xc8a24a : 0x3a5a9a)); f.position.set(i * 4, b.h - 3, b.d / 2 + 0.12); g.add(f);
  }
  g.position.set(b.x, hh, b.z);
  return g;
}

function fountain(b, hh) {
  const g = new THREE.Group();
  const stone = M(0xc4bfb4);
  const basin = cyl(4, 4.4, 1.0, stone, 20); basin.position.y = 0.5; g.add(basin);
  const water = cyl(3.6, 3.6, 0.1, new THREE.MeshPhongMaterial({ color: 0x4aa8d8, transparent: true, opacity: 0.8 }), 20); water.position.y = 1.0; g.add(water);
  const col = cyl(0.5, 0.7, 3, stone, 10); col.position.y = 1.5; g.add(col);
  const top = cyl(1.6, 0.4, 0.6, stone, 12); top.position.y = 3.1; g.add(top);
  g.position.set(b.x, hh, b.z);
  return g;
}

function pad(b, hh) {
  const g = new THREE.Group();
  const stone = M(0xa8a298);
  const disc = cyl(b.w / 2, b.w / 2 + 0.3, 0.5, stone, 24); disc.position.y = 0.25; g.add(disc);
  const ring = new THREE.Mesh(new THREE.TorusGeometry(b.w / 2 - 0.8, 0.12, 6, 24), M(0xc8a24a)); ring.rotation.x = Math.PI / 2; ring.position.y = 0.55; g.add(ring);
  for (const a of [0.8, 2.4]) {
    const p = cyl(0.25, 0.3, 4, M(0x6a4a2a), 6); p.position.set(Math.cos(a) * 4, 2, Math.sin(a) * 4); g.add(p);
    const f = box(1.4, 0.9, 0.05, M(0xc8a24a)); f.position.set(Math.cos(a) * 4 + 0.7, 3.4, Math.sin(a) * 4); g.add(f);
  }
  const nest = new THREE.Mesh(new THREE.TorusGeometry(1.6, 0.5, 6, 12), M(0x8a6a3a)); nest.rotation.x = Math.PI / 2; nest.position.y = 0.7; g.add(nest);
  g.position.set(b.x, hh, b.z);
  return g;
}

function tent(b, hh) {
  const g = new THREE.Group();
  const cone = new THREE.Mesh(new THREE.ConeGeometry(b.w * 0.7, b.h, 6), M([0x8a3a2a, 0x6a5a3a, 0x4a4a5a][Math.floor(Math.abs(b.x * 7 + b.z) % 3)]));
  cone.position.y = b.h / 2; cone.castShadow = true; g.add(cone);
  const flap = box(1.2, 2, 0.1, M(0x2a1a10)); flap.position.set(0, 1, b.w * 0.45); flap.rotation.x = -0.3; g.add(flap);
  g.position.set(b.x, hh, b.z);
  return g;
}

function cave(b, hh) {
  const g = new THREE.Group();
  const kind = b.s || 'mine';
  const rockC = kind === 'fire' ? 0x24181c : kind === 'temple' ? 0x6a8a78 : 0x5a5650;
  const rock = new THREE.Mesh(new THREE.DodecahedronGeometry(1, 1), M(rockC));
  rock.scale.set(b.w * 0.6, b.h * 0.8, b.d * 0.9); rock.position.set(0, b.h * 0.3, -b.d * 0.3); rock.castShadow = true; g.add(rock);
  const front = b.d * 0.62;
  const dark = box(b.w * 0.34, b.h * 0.55, 0.3, M(0x050505)); dark.position.set(0, b.h * 0.28, front); g.add(dark);
  for (const sx of [-1, 1]) { const p = box(0.9, b.h * 0.62, 0.9, M(0x7a746a)); p.position.set(sx * b.w * 0.2, b.h * 0.31, front + 0.1); g.add(p); }
  const lintel = box(b.w * 0.5, 0.9, 1.0, M(0x7a746a)); lintel.position.set(0, b.h * 0.62, front + 0.1); g.add(lintel);
  const torch = new THREE.PointLight(kind === 'temple' ? 0x60ffb0 : 0xff9040, 1.2, 14); torch.position.set(0, 3, front + 1.5); g.add(torch);
  if (kind === 'temple') { // carved columns and a jade lintel
    for (const sx of [-1, 1]) { const col = cyl(0.55, 0.65, b.h * 0.8, M(0x9ac8b0), 8); col.position.set(sx * b.w * 0.3, b.h * 0.4, front + 0.6); g.add(col); }
    const arch = box(b.w * 0.8, 0.8, 1.2, M(0x2f9a78)); arch.position.set(0, b.h * 0.82, front + 0.6); g.add(arch);
  } else if (kind === 'fire') { // glowing maw
    const glow = box(b.w * 0.3, b.h * 0.45, 0.2, new THREE.MeshBasicMaterial({ color: 0xff5a1a })); glow.position.set(0, b.h * 0.25, front + 0.2); g.add(glow);
    for (let i = 0; i < 5; i++) { const sp = new THREE.Mesh(new THREE.ConeGeometry(0.5, 2.4, 5), M(0x14141c)); sp.position.set((i - 2) * b.w * 0.16, b.h * 0.85, front - 0.5); g.add(sp); }
  }
  g.position.set(b.x, hh, b.z);
  return g;
}

// ---- Jadereich: pagodas, gates, ponds, lanterns, harbor ----
function pagoda(b, hh) {
  const g = new THREE.Group();
  const st = STYLES[b.s] || STYLES.jade;
  const base = box(b.w + 5, 1.2, b.d + 5, M(0xb8a898)); base.position.y = 0.6; g.add(base);
  const steps = box(b.w * 0.5, 0.6, 3, M(0xc8b8a8)); steps.position.set(0, 0.3, b.d / 2 + 3.4); g.add(steps);
  const tiers = Math.max(3, Math.round(b.h / 4.4));
  let y = 1.2, w = b.w, d = b.d;
  for (let i = 0; i < tiers; i++) {
    const hgt = i === 0 ? 4.6 : 3.4;
    const body = box(w * 0.82, hgt, d * 0.82, M(st.wall)); body.position.y = y + hgt / 2; g.add(body);
    for (const sx of [-1, 1]) for (const sz of [-1, 1]) { const pil = cyl(0.3, 0.3, hgt, M(st.trim), 6); pil.position.set(sx * w * 0.41, y + hgt / 2, sz * d * 0.41); g.add(pil); }
    const win = box(w * 0.3, 1.1, 0.2, M(0xffe9a0, { emissive: 0x6a5a20 })); win.position.set(0, y + hgt * 0.55, d * 0.41 + 0.05); g.add(win);
    const roof = new THREE.Mesh(new THREE.ConeGeometry(1, 1.8, 4), M(st.roof));
    roof.rotation.y = Math.PI / 4; roof.scale.set((w / 2 + 1.6) * Math.SQRT2, 1, (d / 2 + 1.6) * Math.SQRT2); roof.position.y = y + hgt + 0.9; roof.castShadow = true; g.add(roof);
    for (const sx of [-1, 1]) for (const sz of [-1, 1]) { const tip = box(0.5, 0.5, 0.5, M(st.trim)); tip.position.set(sx * (w / 2 + 1.5), y + hgt + 0.6, sz * (d / 2 + 1.5)); tip.rotation.z = sx * 0.7; g.add(tip); }
    y += hgt + 1.1; w *= 0.8; d *= 0.8;
  }
  const spire = cyl(0.12, 0.3, 3.6, M(st.trim), 6); spire.position.y = y + 1.4; g.add(spire);
  const orb = new THREE.Mesh(new THREE.SphereGeometry(0.4, 8, 6), M(0xe8c040, { emissive: 0x443300 })); orb.position.y = y + 3.4; g.add(orb);
  g.position.set(b.x, hh, b.z);
  return g;
}

function gate(b, hh) { // torii
  const g = new THREE.Group();
  const red = M(0xc0302a), black = M(0x201a1a);
  for (const s of [-1, 1]) { const p = cyl(0.5, 0.6, b.h, red, 8); p.position.set(s * b.w / 2.4, b.h / 2, 0); g.add(p); }
  const beam1 = box(b.w + 4, 0.7, 1.0, black); beam1.position.y = b.h + 0.3; g.add(beam1);
  const beam2 = box(b.w + 1.2, 0.5, 0.8, red); beam2.position.y = b.h - 1.4; g.add(beam2);
  for (const s of [-1, 1]) { const tip = box(1.4, 0.5, 1.0, black); tip.position.set(s * (b.w / 2 + 2.2), b.h + 0.7, 0); tip.rotation.z = s * 0.35; g.add(tip); }
  g.rotation.y = b.r ? 0 : Math.PI / 2; // the beam spans the street
  g.position.set(b.x, hh, b.z);
  return g;
}

function pond(b, hh) {
  const g = new THREE.Group();
  const rim = box(b.w + 1.2, 0.6, b.d + 1.2, M(0xb8a898)); rim.position.y = 0.2; g.add(rim);
  const water = box(b.w, 0.2, b.d, new THREE.MeshPhongMaterial({ color: 0x3aa8b8, transparent: true, opacity: 0.85, shininess: 90 })); water.position.y = 0.55; g.add(water);
  for (let i = 0; i < 7; i++) { const pad = cyl(0.55, 0.55, 0.06, M(0x3f8a3a), 8); pad.position.set(((i * 37) % 17 - 8) / 8 * b.w * 0.4, 0.7, ((i * 53) % 13 - 6) / 6 * b.d * 0.38); g.add(pad); }
  // arched red bridge
  const deck = box(3, 0.4, b.d + 4, M(0xc0302a)); deck.position.y = 1.2; g.add(deck);
  const hump = box(3, 0.4, b.d * 0.5, M(0xc0302a)); hump.position.y = 1.7; g.add(hump);
  for (const s of [-1, 1]) { const rail = box(0.2, 0.8, b.d + 4, M(0xe8c860)); rail.position.set(s * 1.4, 1.8, 0); g.add(rail); }
  g.position.set(b.x, hh, b.z);
  return g;
}

function lantern(b, hh) {
  const g = new THREE.Group();
  const post = cyl(0.14, 0.18, b.h, M(0x3a2a22), 6); post.position.y = b.h / 2; g.add(post);
  const lamp = box(0.7, 0.8, 0.7, new THREE.MeshBasicMaterial({ color: 0xffb040 })); lamp.position.y = b.h + 0.3; g.add(lamp);
  const cap = new THREE.Mesh(new THREE.ConeGeometry(0.7, 0.5, 4), M(0xc0302a)); cap.rotation.y = Math.PI / 4; cap.position.y = b.h + 0.95; g.add(cap);
  g.position.set(b.x, hh, b.z);
  return g;
}

function dock(b, hh) {
  const g = new THREE.Group();
  const plank = M(0x8a6a40);
  const deck = box(b.w, 0.35, b.d, plank); deck.position.y = 1.03; g.add(deck);
  for (let x = -b.w / 2 + 1.5; x <= b.w / 2; x += 4) for (const z of [-1, 1]) { const pile = cyl(0.25, 0.3, 4.2, M(0x4a3220), 6); pile.position.set(x, -0.6, z * (b.d / 2 - 0.3)); g.add(pile); }
  g.position.set(b.x, 0, b.z);
  return g;
}

function ship(b, hh) {
  const s = makeShip(b.s === 'bamboo' ? 'junk' : 'cog');
  const g = new THREE.Group();
  g.add(s.obj);
  g.rotation.y = b.r === 1 ? Math.PI / 2 : -Math.PI / 2; // sideways along the pier
  g.position.set(b.x, -0.1, b.z);
  g.userData.ship = s;
  return g;
}

const BUILD = { house, hut: house, tower, wall, keep, fountain, pad, tent, cave, pagoda, gate, pond, lantern, dock, ship };
const NO_SOLID = new Set(['pad', 'dock', 'ship', 'lantern', 'pond', 'gate']);

// Builds the cabin at the world origin with its door and furniture (from the first server release).
export function buildCabin(boxes, scene) {
  const g = new THREE.Group();
  const logs = M(0x8a5a2e), stone = M(0x8a8a84), wood = M(0x5a3a1c);
  const h0 = heightAt(0, 0);
  for (const b of boxes) {
    const m = box(b.w, b.h, b.d, logs);
    m.position.set(b.x, h0 + b.h / 2, b.z);
    g.add(m);
  }
  const base = box(10.6, 0.5, 8.6, stone); base.position.set(0, h0 + 0.1, 0); g.add(base);
  const floor = box(9.6, 0.1, 7.6, M(0x6b4a28)); floor.position.set(0, h0 + 0.36, 0); g.add(floor);
  const roof = new THREE.Mesh(new THREE.ConeGeometry(1, 2.6, 4), M(0x8a2f22));
  roof.rotation.y = Math.PI / 4;
  roof.scale.set(6.4 / Math.SQRT1_2, 1, 5.3 / Math.SQRT1_2);
  roof.position.set(0, h0 + 3 + 1.3, 0); roof.castShadow = true; g.add(roof);
  const chim = box(0.9, 2.6, 0.9, stone); chim.position.set(-2.8, h0 + 4.4, -1.6); g.add(chim);
  const frame = box(2.1, 0.2, 0.5, wood); frame.position.set(0, h0 + 2.5, 4); g.add(frame);
  for (const sx of [-1, 1]) { const post = box(0.2, 2.5, 0.5, wood); post.position.set(sx * 1.0, h0 + 1.25, 4); g.add(post); }
  const winM = WIN;
  for (const [x, z, ry] of [[-3.2, 4.0, 0], [3.2, 4.0, 0], [-5, 0, Math.PI / 2], [5, 0, Math.PI / 2]]) {
    const w = box(1.1, 0.9, 0.5, winM); w.position.set(x, h0 + 1.8, z); w.rotation.y = ry; g.add(w);
  }
  const bed = box(1.2, 0.5, 2.2, M(0x9a3a3a)); bed.position.set(3.6, h0 + 0.7, -2.6); g.add(bed);
  const table = box(1.6, 0.9, 1.0, wood); table.position.set(-2.4, h0 + 0.85, 1.4); g.add(table);
  const lamp = new THREE.PointLight(0xffb060, 1.2, 12); lamp.position.set(0, h0 + 2.4, 0); g.add(lamp);
  for (let i = 0; i < 6; i++) {
    const lg = cyl(0.16, 0.16, 1.2, wood, 8);
    lg.rotation.z = Math.PI / 2; lg.position.set(6.2, h0 + 0.2 + Math.floor(i / 3) * 0.3, 1 + (i % 3) * 0.34); g.add(lg);
  }
  scene.add(g);
}

// Builds every town building and returns collision boxes.
export const ships = [];
export function buildTowns(scene, buildings) {
  const solids = [];
  const g = new THREE.Group();
  for (const b of buildings) {
    const fn = BUILD[b.k];
    if (!fn) continue;
    const hh = heightAt(b.x, b.z);
    const obj = fn(b, b.k === 'pad' ? hh : hh - 0.3);
    g.add(obj);
    if (obj.userData.ship) ships.push(obj.userData.ship);
    if (!NO_SOLID.has(b.k)) solids.push({ x: b.x, z: b.z, w: b.k === 'tent' ? b.w * 0.9 : b.w, d: b.k === 'tent' ? b.d * 0.9 : b.d, tower: b.k === 'tower' });
  }
  scene.add(g);
  return solids;
}
