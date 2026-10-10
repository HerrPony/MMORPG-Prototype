import * as THREE from 'three';

const M = (c, o = {}) => new THREE.MeshLambertMaterial({ color: c, ...o });
const B = (c) => new THREE.MeshBasicMaterial({ color: c });
const sph = (r, m, ws = 10, hs = 8) => { const x = new THREE.Mesh(new THREE.SphereGeometry(r, ws, hs), m); x.castShadow = true; return x; };
const box = (w, h, d, m) => { const x = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), m); x.castShadow = true; return x; };
const cone = (r, h, m, s = 6) => { const x = new THREE.Mesh(new THREE.ConeGeometry(r, h, s), m); x.castShadow = true; return x; };
const pivot = (x, y, z) => { const g = new THREE.Group(); g.position.set(x, y, z); return g; };

let glowTex = null;
function glow() {
  if (glowTex) return glowTex;
  const c = document.createElement('canvas');
  c.width = c.height = 64;
  const g = c.getContext('2d');
  const grd = g.createRadialGradient(32, 32, 0, 32, 32, 32);
  grd.addColorStop(0, 'rgba(255,255,255,1)');
  grd.addColorStop(0.4, 'rgba(255,255,255,0.5)');
  grd.addColorStop(1, 'rgba(255,255,255,0)');
  g.fillStyle = grd; g.fillRect(0, 0, 64, 64);
  glowTex = new THREE.CanvasTexture(c);
  return glowTex;
}

function deadFlip(root, u, st, dt, lift = 0) {
  u.dead += ((st.dead ? 1 : 0) - u.dead) * Math.min(1, dt * 6);
  root.rotation.z = u.dead * Math.PI;
  root.position.y = u.dead * lift;
}

export function makeCroc(seed = 0) {
  const root = new THREE.Group();
  const body = new THREE.Group(); root.add(body);
  const skin = M([0x4a6a2a, 0x3f5a28, 0x55702e][seed % 3]), belly = M(0xb8c080), dark = M(0x1a2a10);
  const torso = sph(0.45, skin, 12, 8); torso.scale.set(0.9, 0.55, 2.5); torso.position.set(0, 0.5, 0); body.add(torso);
  const under = sph(0.4, belly, 10, 6); under.scale.set(0.8, 0.35, 2.2); under.position.set(0, 0.34, 0); body.add(under);
  const head = new THREE.Group(); head.position.set(0, 0.55, 1.2); body.add(head);
  const upper = box(0.38, 0.16, 0.9, skin); upper.position.set(0, 0.05, 0.45); head.add(upper);
  const lower = box(0.34, 0.1, 0.8, skin); lower.position.set(0, -0.08, 0.4); head.add(lower);
  for (let i = 0; i < 5; i++) for (const s of [-1, 1]) { const t = cone(0.025, 0.08, M(0xf0ecd0), 4); t.rotation.x = Math.PI; t.position.set(s * 0.16, -0.02, 0.15 + i * 0.15); head.add(t); }
  for (const s of [-1, 1]) { const e = sph(0.06, B(0xffcc22), 6, 5); e.position.set(s * 0.14, 0.15, 0.1); head.add(e); }
  const tail = new THREE.Group(); tail.position.set(0, 0.5, -1.1); body.add(tail);
  const t1 = cone(0.3, 1.5, skin, 6); t1.rotation.x = -Math.PI / 2; t1.position.z = -0.7; tail.add(t1);
  const legs = [];
  for (const [x, z] of [[-0.4, 0.7], [0.4, 0.7], [-0.4, -0.6], [0.4, -0.6]]) {
    const p = pivot(x, 0.35, z); const l = box(0.16, 0.4, 0.2, skin); l.position.y = -0.15; p.add(l); body.add(p); legs.push(p);
  }
  for (let i = 0; i < 6; i++) { const sp = cone(0.05, 0.14, dark, 4); sp.position.set(0, 0.88, 0.9 - i * 0.35); body.add(sp); }
  const u = { phase: 0, dead: 0 };
  return {
    obj: root, height: 1.0, kind: 'croc',
    attack() { u.bite = 0.01; },
    update(dt, st) {
      const mv = st.speed > 0.1;
      if (mv) u.phase += dt * 8 * Math.min(1.4, st.speed);
      legs.forEach((l, i) => { l.rotation.x = Math.sin(u.phase + (i % 2 ? 0 : Math.PI)) * (mv ? 0.6 : 0); });
      tail.rotation.y = Math.sin(performance.now() / 500 + u.phase) * 0.3;
      if (u.bite > 0) { u.bite += dt; head.rotation.x = -Math.sin(Math.min(1, u.bite / 0.4) * Math.PI) * 0.5; if (u.bite > 0.4) { u.bite = 0; head.rotation.x = 0; } }
      deadFlip(root, u, st, dt, 0.7);
    },
  };
}

export function makeScorpion(seed = 0) {
  const root = new THREE.Group();
  const body = new THREE.Group(); root.add(body);
  const shell = M([0x8a6a30, 0x6a4a28, 0x2a2a2a][seed % 3]), dark = M(0x3a2a14);
  const torso = sph(0.4, shell, 10, 8); torso.scale.set(1, 0.55, 1.3); torso.position.set(0, 0.45, 0); body.add(torso);
  const head = sph(0.22, shell, 8, 6); head.position.set(0, 0.42, 0.55); body.add(head);
  for (const s of [-1, 1]) {
    const e = sph(0.04, B(0xff3030), 5, 4); e.position.set(s * 0.1, 0.52, 0.7); body.add(e);
    const arm = new THREE.Group(); arm.position.set(s * 0.25, 0.4, 0.7);
    const upper = box(0.1, 0.1, 0.5, shell); upper.position.z = 0.25; arm.add(upper);
    const c1 = cone(0.14, 0.35, shell, 4); c1.rotation.x = Math.PI / 2; c1.position.z = 0.6; arm.add(c1);
    arm.rotation.y = -s * 0.35; body.add(arm);
  }
  const tail = new THREE.Group(); tail.position.set(0, 0.5, -0.5); body.add(tail);
  let py = 0, pz = 0;
  for (let i = 0; i < 5; i++) {
    const seg = sph(0.15 - i * 0.012, shell, 6, 5); const a = 0.3 + i * 0.42;
    pz -= Math.cos(a) * 0.2; py += Math.sin(a) * 0.2 + 0.04;
    seg.position.set(0, py, pz); tail.add(seg);
  }
  const sting = cone(0.07, 0.3, M(0xe8e0a0), 5); sting.position.set(0, py + 0.05, pz + 0.25); sting.rotation.x = 2.2; tail.add(sting);
  const legs = [];
  for (let i = 0; i < 8; i++) {
    const side = i < 4 ? -1 : 1, k = i % 4;
    const p = pivot(side * 0.3, 0.4, 0.25 - k * 0.2);
    const l = box(0.45, 0.05, 0.05, dark); l.position.set(side * 0.22, -0.12, 0); l.rotation.z = side * -0.5; p.add(l);
    body.add(p); legs.push(p);
  }
  const u = { phase: 0, dead: 0 };
  return {
    obj: root, height: 0.9, kind: 'scorpion',
    attack() { u.sting = 0.01; },
    update(dt, st) {
      const mv = st.speed > 0.1;
      if (mv) u.phase += dt * 14 * Math.min(1.4, st.speed);
      legs.forEach((l, i) => { l.rotation.y = Math.sin(u.phase + i * 1.2) * (mv ? 0.35 : 0.03); });
      tail.rotation.x = -0.1 + Math.sin(performance.now() / 700) * 0.06;
      if (u.sting > 0) { u.sting += dt; tail.rotation.x = Math.sin(Math.min(1, u.sting / 0.35) * Math.PI) * 0.8; if (u.sting > 0.35) u.sting = 0; }
      deadFlip(root, u, st, dt, 0.8);
    },
  };
}

export function makeWisp(color = 0x9aff4a, face = true) {
  const root = new THREE.Group();
  const core = sph(0.28, B(color), 12, 10); core.position.y = 1.4; root.add(core);
  const halo = new THREE.Sprite(new THREE.SpriteMaterial({ map: glow(), color, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending }));
  halo.scale.set(1.8, 1.8, 1); halo.position.y = 1.4; root.add(halo);
  if (face) for (const s of [-1, 1]) { const e = sph(0.045, B(0x112211), 5, 4); e.position.set(s * 0.1, 1.45, 0.25); root.add(e); }
  const wisps = [];
  for (let i = 0; i < 6; i++) { const p = new THREE.Sprite(new THREE.SpriteMaterial({ map: glow(), color, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending, opacity: 0.6 })); p.scale.set(0.35, 0.35, 1); root.add(p); wisps.push(p); }
  const u = { dead: 0 };
  return {
    obj: root, height: 1.9, kind: 'wisp',
    attack() { u.punch = 0.01; },
    update(dt, st) {
      const t = performance.now() / 1000;
      const bob = Math.sin(t * 2.2) * 0.12;
      u.dead += ((st.dead ? 1 : 0) - u.dead) * Math.min(1, dt * 5);
      core.position.y = 1.4 + bob - u.dead * 1.2;
      halo.position.y = core.position.y;
      halo.scale.setScalar((1.8 + Math.sin(t * 5) * 0.2) * (1 - u.dead * 0.8));
      wisps.forEach((p, i) => { const a = t * 1.5 + i; p.position.set(Math.cos(a) * 0.5, core.position.y - 0.1 + Math.sin(a * 1.3) * 0.3, Math.sin(a) * 0.5); });
      core.scale.setScalar(1 - u.dead * 0.7);
    },
  };
}

export function makeGriffin() {
  const root = new THREE.Group();
  const body = new THREE.Group(); root.add(body);
  const fur = M(0xc89a58), feather = M(0xf2ece0), dark = M(0x6a4a2a);
  const torso = sph(0.55, fur, 12, 10); torso.scale.set(0.9, 0.8, 1.9); torso.position.set(0, 1.0, -0.3); body.add(torso);
  const chest = sph(0.5, feather, 10, 8); chest.position.set(0, 1.15, 0.65); body.add(chest);
  const neck = new THREE.Mesh(new THREE.CylinderGeometry(0.22, 0.3, 0.7, 8), feather); neck.position.set(0, 1.6, 0.9); neck.rotation.x = 0.4; body.add(neck);
  const head = new THREE.Group(); head.position.set(0, 1.98, 1.15); body.add(head);
  const skull = sph(0.26, feather, 10, 8); head.add(skull);
  const beak = cone(0.12, 0.5, M(0xe8b030), 5); beak.rotation.x = Math.PI / 2 + 0.2; beak.position.set(0, -0.08, 0.38); head.add(beak);
  for (const s of [-1, 1]) {
    const e = sph(0.05, B(0xffd040), 6, 5); e.position.set(s * 0.15, 0.06, 0.16); head.add(e);
    const ear = cone(0.06, 0.28, feather, 4); ear.position.set(s * 0.13, 0.28, -0.05); head.add(ear);
  }
  const tail = cone(0.12, 1.0, dark, 6); tail.rotation.x = -2.0; tail.position.set(0, 1.05, -1.4); body.add(tail);
  const legs = [];
  for (const [x, z, f] of [[-0.32, 0.5, true], [0.32, 0.5, true], [-0.3, -0.7, false], [0.3, -0.7, false]]) {
    const p = pivot(x, 0.75, z); const l = box(0.14, 0.7, 0.16, f ? feather : fur); l.position.y = -0.3; p.add(l);
    const claw = cone(0.09, 0.2, M(0xe8b030), 4); claw.position.set(0, -0.7, 0.08); claw.rotation.x = 1.57; p.add(claw);
    body.add(p); legs.push(p);
  }
  const wings = [];
  for (const s of [-1, 1]) {
    const w = new THREE.Group(); w.position.set(s * 0.5, 1.45, 0.1);
    for (let i = 0; i < 6; i++) {
      const f = box(0.35, 0.03, 1.6 - i * 0.12, feather); f.position.set(s * (0.45 + i * 0.5), 0.0, -0.3 - i * 0.05); f.rotation.y = s * i * 0.04; w.add(f);
    }
    const arm = box(2.9, 0.1, 0.4, fur); arm.position.set(s * 1.5, 0.05, 0.2); w.add(arm);
    body.add(w); wings.push({ g: w, s });
  }
  const u = { t: 0 };
  return {
    obj: root, height: 2.4, kind: 'griffin', seatY: 1.55, seatZ: -0.05,
    attack() {},
    update(dt, st) {
      u.t += dt * (st.flying ? 4.2 : 1.2);
      const flap = Math.sin(u.t * 3);
      wings.forEach(({ g, s }) => { g.rotation.z = s * (st.flying ? flap * 0.7 : -0.9 + flap * 0.05); });
      body.position.y = st.flying ? Math.sin(u.t * 3) * 0.1 : 0;
      head.rotation.y = Math.sin(u.t * 0.7) * 0.25;
      legs.forEach((l) => { l.rotation.x = st.flying ? -0.6 : 0; });
    },
  };
}

// ---- interactive objects ----
export function makeObject(kind) {
  const root = new THREE.Group();
  let anim = () => {};
  let height = 0.8;
  switch (kind) {
    case 'pumpkin': {
      const p = sph(0.38, M(0xe8801a), 10, 8); p.scale.y = 0.8; p.position.y = 0.3; root.add(p);
      const st = box(0.08, 0.15, 0.08, M(0x4a6a2a)); st.position.y = 0.65; root.add(st);
      break;
    }
    case 'wheat': {
      for (let i = 0; i < 7; i++) { const s = cone(0.05, 0.9, M(0xe0b838), 4); s.position.set(Math.cos(i) * 0.15, 0.45, Math.sin(i) * 0.15); s.rotation.z = Math.cos(i) * 0.15; root.add(s); }
      const tie = new THREE.Mesh(new THREE.TorusGeometry(0.15, 0.02, 4, 8), M(0x6a4a2a)); tie.rotation.x = Math.PI / 2; tie.position.y = 0.3; root.add(tie);
      break;
    }
    case 'spring': {
      const d = new THREE.Mesh(new THREE.CylinderGeometry(0.5, 0.55, 0.12, 14), new THREE.MeshPhongMaterial({ color: 0x4ac8ff, transparent: true, opacity: 0.85, emissive: 0x1a5a8a })); d.position.y = 0.1; root.add(d);
      const h = new THREE.Sprite(new THREE.SpriteMaterial({ map: glow(), color: 0x66ccff, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending })); h.scale.set(1.8, 1.8, 1); h.position.y = 0.4; root.add(h);
      anim = (t) => { h.material.opacity = 0.6 + Math.sin(t * 3) * 0.3; };
      break;
    }
    case 'tracks': {
      for (let i = 0; i < 4; i++) { const p = new THREE.Mesh(new THREE.CircleGeometry(0.14, 8), M(0x3a2a1a)); p.rotation.x = -Math.PI / 2; p.position.set((i % 2 ? 0.2 : -0.2), 0.03, i * 0.35 - 0.5); root.add(p); }
      height = 0.5;
      break;
    }
    case 'crystal': {
      const c = new THREE.Mesh(new THREE.OctahedronGeometry(0.4), new THREE.MeshBasicMaterial({ color: 0x5ad0ff })); c.scale.y = 1.6; c.position.y = 0.8; root.add(c);
      const h = new THREE.Sprite(new THREE.SpriteMaterial({ map: glow(), color: 0x66ccff, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending })); h.scale.set(2, 2, 1); h.position.y = 0.8; root.add(h);
      anim = (t) => { c.rotation.y = t * 1.2; c.position.y = 0.8 + Math.sin(t * 2) * 0.08; };
      height = 1.4;
      break;
    }
    case 'lotus': {
      const pad = new THREE.Mesh(new THREE.CircleGeometry(0.5, 10), M(0x3f7a3a)); pad.rotation.x = -Math.PI / 2; pad.position.y = 0.05; root.add(pad);
      for (let i = 0; i < 6; i++) { const pe = cone(0.12, 0.35, M(0xffa0d0), 4); const a = (i / 6) * 6.28; pe.position.set(Math.cos(a) * 0.12, 0.25, Math.sin(a) * 0.12); pe.rotation.z = Math.cos(a) * 0.5; pe.rotation.x = -Math.sin(a) * 0.5; root.add(pe); }
      break;
    }
    case 'cactus': {
      const c = sph(0.3, M(0xe83a6a), 8, 6); c.scale.y = 0.8; c.position.y = 0.3; root.add(c);
      for (let i = 0; i < 4; i++) { const s = cone(0.03, 0.2, M(0xe8e0a0), 3); s.position.set(Math.cos(i * 1.6) * 0.25, 0.4, Math.sin(i * 1.6) * 0.25); root.add(s); }
      break;
    }
    case 'urn': {
      const u = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.28, 0.55, 10), M(0xb87a48)); u.position.y = 0.3; root.add(u);
      const n = new THREE.Mesh(new THREE.CylinderGeometry(0.14, 0.2, 0.18, 10), M(0xb87a48)); n.position.y = 0.65; root.add(n);
      break;
    }
    case 'portal': {
      const ring = new THREE.Mesh(new THREE.TorusGeometry(1.5, 0.15, 8, 32), M(0x8a60d8, { emissive: 0x4a20a0 })); ring.position.y = 1.9; root.add(ring);
      const swirl = new THREE.Mesh(new THREE.CircleGeometry(1.4, 24), new THREE.MeshBasicMaterial({ color: 0xa070ff, transparent: true, opacity: 0.55, side: THREE.DoubleSide })); swirl.position.y = 1.9; root.add(swirl);
      const h = new THREE.Sprite(new THREE.SpriteMaterial({ map: glow(), color: 0xb080ff, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending })); h.scale.set(4.5, 4.5, 1); h.position.y = 1.9; root.add(h);
      const base = new THREE.Mesh(new THREE.CylinderGeometry(1.8, 2, 0.3, 16), M(0x5a5a64)); base.position.y = 0.15; root.add(base);
      anim = (t) => { swirl.rotation.z = t * 1.5; swirl.material.opacity = 0.45 + Math.sin(t * 3) * 0.15; h.material.opacity = 0.7 + Math.sin(t * 2) * 0.2; };
      height = 3.6;
      break;
    }
  }
  return { obj: root, height, kind: 'object', attack() {}, update(dt, st) { anim(performance.now() / 1000); } };
}
