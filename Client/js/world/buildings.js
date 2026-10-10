import * as THREE from 'three';
import { heightAt } from '../util.js';

const M = (c, o = {}) => new THREE.MeshLambertMaterial({ color: c, ...o });
const box = (w, h, d, m) => { const x = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), m); x.castShadow = x.receiveShadow = true; return x; };
const cyl = (rt, rb, h, m, s = 12) => { const x = new THREE.Mesh(new THREE.CylinderGeometry(rt, rb, h, s), m); x.castShadow = true; return x; };

const STYLES = {
  timber: { wall: 0x9a6a3a, roof: 0x7a2f22, trim: 0x5a3a1c },
  plaster: { wall: 0xe8dcc0, roof: 0xb04a2a, trim: 0x6a4a2c },
  stilt: { wall: 0x5a4a32, roof: 0x6a7a3a, trim: 0x3a2e1e },
  lodge: { wall: 0x6a4a30, roof: 0xdde8f0, trim: 0x3a2a1c },
  adobe: { wall: 0xd8b878, roof: 0xc89858, trim: 0x8a6a3a },
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
  const wm = M(0xffe9a0, { emissive: 0x5a4a1a });
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
  const stone = M(b.s === 'adobe' ? 0xd8b878 : 0xcfc9bc);
  const body = box(b.w, b.h, b.d, stone); body.position.y = b.h / 2; g.add(body);
  const roof = new THREE.Mesh(new THREE.ConeGeometry(1, 5, 4), M(b.s === 'adobe' ? 0xc89858 : 0x3a5a9a));
  roof.rotation.y = Math.PI / 4; roof.scale.set(b.w / 2 * Math.SQRT2, 1, b.d / 2 * Math.SQRT2); roof.position.y = b.h + 2.5; roof.castShadow = true; g.add(roof);
  if (b.s !== 'adobe') for (const sx of [-1, 1]) for (const sz of [-1, 1]) {
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
  const rock = new THREE.Mesh(new THREE.DodecahedronGeometry(1, 1), M(0x5a5650));
  rock.scale.set(b.w * 0.6, b.h * 0.8, b.d * 0.9); rock.position.set(0, b.h * 0.3, -b.d * 0.3); rock.castShadow = true; g.add(rock);
  const front = b.d * 0.62;
  const dark = box(b.w * 0.34, b.h * 0.55, 0.3, M(0x050505)); dark.position.set(0, b.h * 0.28, front); g.add(dark);
  for (const sx of [-1, 1]) { const p = box(0.9, b.h * 0.62, 0.9, M(0x7a746a)); p.position.set(sx * b.w * 0.2, b.h * 0.31, front + 0.1); g.add(p); }
  const lintel = box(b.w * 0.5, 0.9, 1.0, M(0x7a746a)); lintel.position.set(0, b.h * 0.62, front + 0.1); g.add(lintel);
  const torch = new THREE.PointLight(0xff9040, 1.2, 14); torch.position.set(0, 3, front + 1.5); g.add(torch);
  g.position.set(b.x, hh, b.z);
  return g;
}

const BUILD = { house, hut: house, tower, wall, keep, fountain, pad, tent, cave };

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
  const winM = M(0xffe9a0, { emissive: 0x6a5a20 });
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
export function buildTowns(scene, buildings) {
  const solids = [];
  const g = new THREE.Group();
  for (const b of buildings) {
    const fn = BUILD[b.k];
    if (!fn) continue;
    const hh = heightAt(b.x, b.z);
    g.add(fn(b, b.k === 'pad' ? hh : hh - 0.3));
    if (b.k !== 'pad') solids.push({ x: b.x, z: b.z, w: b.k === 'tent' ? b.w * 0.9 : b.w, d: b.k === 'tent' ? b.d * 0.9 : b.d, tower: b.k === 'tower' });
  }
  scene.add(g);
  return solids;
}
