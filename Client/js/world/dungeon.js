import * as THREE from 'three';
import { makeProps } from './vegetation.js';

const M = (c, o = {}) => new THREE.MeshLambertMaterial({ color: c, ...o });
const B = (c, o = {}) => new THREE.MeshBasicMaterial({ color: c, ...o });
const WALL_H = 9;

function rng(seed) {
  let s = seed >>> 0;
  return () => { s = (s * 1664525 + 1013904223) >>> 0; return s / 4294967296; };
}

// Look of each dungeon: colours of rock and floor, fog, torch flames and the lights.
const STYLES = {
  mine: { bg: 0x07070b, amb: 0x6a6a88, hemi: [0x8888aa, 0x221a14], rock: 0x4a443e, floor: 0x3a3632, ceil: 0x1c1a18, wallH: [0.08, 0.1], flame: 0xffaa33, light: 0xff9a44, player: 0xffd9a0 },
  temple: { bg: 0x04100e, amb: 0x5a8a7a, hemi: [0x78b8a0, 0x10241c], rock: 0x4a6a5a, floor: 0x2e463c, ceil: 0x142a22, wallH: [0.42, 0.14], flame: 0x60ffb8, light: 0x50ffc0, player: 0xc8fff0 },
  fire: { bg: 0x120504, amb: 0x8a5a4a, hemi: [0xb06a50, 0x2a0e08], rock: 0x2c1e1e, floor: 0x241618, ceil: 0x120a0a, wallH: [0.0, 0.1], flame: 0xff5a1a, light: 0xff5a20, player: 0xffb888 },
};

// Builds the dungeon: rooms and corridors given in local coordinates around origin.
export function buildDungeon(info) {
  const ST = STYLES[info.style] || STYLES.mine;
  const scene = new THREE.Scene();
  scene.background = new THREE.Color(ST.bg);
  scene.fog = new THREE.Fog(ST.bg, 8, 55);
  scene.add(new THREE.AmbientLight(ST.amb, 0.9));
  scene.add(new THREE.HemisphereLight(ST.hemi[0], ST.hemi[1], 0.5));
  const ox = info.origin[0], oz = info.origin[1];
  const rooms = info.rooms.map((r) => ({ ...r, x: r.x + ox, z: r.z + oz }));
  const rock = M(ST.rock), floorM = M(ST.floor);
  const r = rng(77);
  const P = makeProps();

  const inside = (x, z, skip) => rooms.some((q, i) => i !== skip && x > q.x - q.w / 2 + 0.01 && x < q.x + q.w / 2 - 0.01 && z > q.z - q.d / 2 - 0.01 && z < q.z + q.d / 2 + 0.01);
  const wallPieces = [];
  const flames = [];
  const anims = [];

  const floorFor = (theme) => {
    switch (theme) {
      case 'throne': return M(0x2a2a3a);
      case 'web': return M(0x2e2e2a);
      case 'shrine': return M(0x24483c);
      case 'flooded': return M(0x28403a);
      case 'inferno': return M(0x1c1012);
      case 'obsidian': return M(0x16161e);
      case 'forge': return M(0x2a1c1c);
      default: return floorM;
    }
  };

  rooms.forEach((q, i) => {
    const fl = new THREE.Mesh(new THREE.BoxGeometry(q.w, 0.4, q.d), floorFor(q.theme));
    fl.position.set(q.x, -0.2, q.z); fl.receiveShadow = true; scene.add(fl);
    const ce = new THREE.Mesh(new THREE.PlaneGeometry(q.w, q.d), M(ST.ceil, { side: THREE.BackSide }));
    ce.rotation.x = Math.PI / 2; ce.position.set(q.x, WALL_H, q.z); scene.add(ce);
    // perimeter walls in 2 m pieces, skipping those that open into another room
    const x0 = q.x - q.w / 2, x1 = q.x + q.w / 2, z0 = q.z - q.d / 2, z1 = q.z + q.d / 2;
    for (let x = x0 + 1; x < x1; x += 2) {
      for (const [z, nz] of [[z0, -1], [z1, 1]]) {
        if (!inside(x, z + nz * 0.5, i) && !inside(x, z - nz * 0.5, i)) wallPieces.push({ x, z: z + nz * 0.6, w: 2.1, d: 1.2 });
      }
    }
    for (let z = z0 + 1; z < z1; z += 2) {
      for (const [x, nx] of [[x0, -1], [x1, 1]]) {
        if (!inside(x + nx * 0.5, z, i) && !inside(x - nx * 0.5, z, i)) wallPieces.push({ x: x + nx * 0.6, z, w: 1.2, d: 2.1 });
      }
    }
  });
  const wallGeo = new THREE.BoxGeometry(1, WALL_H, 1);
  const walls = new THREE.InstancedMesh(wallGeo, rock, wallPieces.length);
  const d = new THREE.Object3D();
  const c = new THREE.Color();
  wallPieces.forEach((p, i) => {
    d.position.set(p.x, WALL_H / 2, p.z);
    d.scale.set(p.w, 0.8 + r() * 0.3, p.d);
    d.updateMatrix();
    walls.setMatrixAt(i, d.matrix);
    c.setHSL(ST.wallH[0], ST.wallH[1], 0.22 + r() * 0.1);
    walls.setColorAt(i, c);
  });
  walls.castShadow = false;
  scene.add(walls);

  const prop = (kind, x, z, s = 1, tintHSL) => {
    const g = new THREE.Group();
    for (const part of P[kind]) {
      const mesh = new THREE.Mesh(part.geo, part.tint && tintHSL ? M(new THREE.Color().setHSL(...tintHSL)) : part.mat);
      mesh.castShadow = false;
      g.add(mesh);
    }
    g.position.set(x, 0, z); g.scale.setScalar(s); g.rotation.y = r() * 6;
    scene.add(g);
  };
  const torch = (x, z) => {
    const post = new THREE.Mesh(new THREE.CylinderGeometry(0.1, 0.12, 2.2, 5), M(0x4a3020)); post.position.set(x, 1.1, z); scene.add(post);
    const fl = new THREE.Mesh(new THREE.SphereGeometry(0.28, 6, 5), B(ST.flame)); fl.position.set(x, 2.4, z); scene.add(fl);
    flames.push(fl);
  };
  const lights = [];
  const addLight = (color, intensity, dist, x, y, z) => { const l = new THREE.PointLight(color, intensity, dist); l.position.set(x, y, z); scene.add(l); lights.push(l); return l; };
  const pillar = (x, z, rad, color) => { const p = new THREE.Mesh(new THREE.CylinderGeometry(rad * 0.9, rad, WALL_H, 8), M(color)); p.position.set(x, WALL_H / 2, z); scene.add(p); };
  const pool = (x, z, rad, color, light = false) => {
    const m = new THREE.Mesh(new THREE.CircleGeometry(rad, 20), B(color));
    m.rotation.x = -Math.PI / 2; m.position.set(x, 0.04, z); scene.add(m);
    const rim = new THREE.Mesh(new THREE.RingGeometry(rad, rad + 0.7, 20), M(0x1a1214)); rim.rotation.x = -Math.PI / 2; rim.position.set(x, 0.03, z); scene.add(rim);
    anims.push((t) => { m.material.color.setHex(color).multiplyScalar(0.85 + Math.sin(t * 2 + x) * 0.15); });
    if (light) addLight(color, 1.0, 16, x, 2, z);
  };

  rooms.forEach((q, i) => {
    // torches along the long walls
    const n = Math.max(2, Math.floor(q.d / 14));
    for (let k = 0; k < n; k++) {
      const z = q.z - q.d / 2 + (k + 0.5) * (q.d / n);
      if (q.w > 8) { torch(q.x - q.w / 2 + 1.2, z); torch(q.x + q.w / 2 - 1.2, z); }
    }
    if (q.w > 8 && i % 2 === 0) addLight(ST.light, 1.3, 38, q.x, 5, q.z);
    switch (q.theme) {
      case 'hall':
        for (const sx of [-1, 1]) for (const sz of [-1, 1]) pillar(q.x + sx * (q.w / 2 - 3), q.z + sz * (q.d / 2 - 3), 0.8, 0x6a645a);
        break;
      case 'mine':
        for (let k = 0; k < 12; k++) prop('stalag', q.x + (r() - 0.5) * (q.w - 6), q.z + (r() - 0.5) * (q.d - 6), 0.6 + r() * 0.8, [0.08, 0.05, 0.3]);
        for (let k = 0; k < 6; k++) {
          const cr = new THREE.Mesh(new THREE.OctahedronGeometry(0.5 + r() * 0.5), B([0x4acfff, 0xff6aff, 0x6aff9a][k % 3]));
          cr.position.set(q.x + (r() - 0.5) * (q.w - 4), 0.8, q.z + (r() - 0.5) * (q.d - 4)); scene.add(cr);
        }
        for (let k = 0; k < 4; k++) { const cart = new THREE.Mesh(new THREE.BoxGeometry(1.6, 0.9, 1.1), M(0x5a4030)); cart.position.set(q.x + (r() - 0.5) * (q.w - 6), 0.5, q.z + (r() - 0.5) * (q.d - 6)); cart.rotation.y = r() * 3; scene.add(cart); }
        break;
      case 'web':
        for (let k = 0; k < 14; k++) {
          const web = new THREE.Mesh(new THREE.PlaneGeometry(4 + r() * 4, 4 + r() * 3), B(0xdddde8, { transparent: true, opacity: 0.35, side: THREE.DoubleSide, depthWrite: false }));
          web.position.set(q.x + (r() - 0.5) * (q.w - 2), 3 + r() * 4, q.z + (r() - 0.5) * (q.d - 2)); web.rotation.set(r() * 1.2, r() * 3, r() * 1.2); scene.add(web);
        }
        for (let k = 0; k < 8; k++) { const co = new THREE.Mesh(new THREE.SphereGeometry(0.7, 8, 6), M(0xc8c8b8)); co.scale.y = 1.6; co.position.set(q.x + (r() - 0.5) * (q.w - 6), 1, q.z + (r() - 0.5) * (q.d - 6)); scene.add(co); }
        for (let k = 0; k < 8; k++) prop('bones', q.x + (r() - 0.5) * (q.w - 4), q.z + (r() - 0.5) * (q.d - 4), 1);
        break;
      case 'throne': {
        const th = new THREE.Group();
        const seat = new THREE.Mesh(new THREE.BoxGeometry(3, 1.4, 3), M(0x6a5a2a)); seat.position.y = 0.7; th.add(seat);
        const back = new THREE.Mesh(new THREE.BoxGeometry(3, 5, 0.8), M(0x8a7a3a)); back.position.set(0, 3.5, -1.4); th.add(back);
        th.position.set(q.x, 0, q.z - q.d / 2 + 4); scene.add(th);
        const rune = new THREE.Mesh(new THREE.RingGeometry(8, 9, 40), B(0xc03030, { side: THREE.DoubleSide, transparent: true, opacity: 0.6 }));
        rune.rotation.x = -Math.PI / 2; rune.position.set(q.x, 0.05, q.z + 2); scene.add(rune);
        for (const sx of [-1, 1]) for (let k = 0; k < 3; k++) {
          pillar(q.x + sx * (q.w / 2 - 3), q.z - q.d / 2 + 8 + k * 14, 0.95, 0x4a4640);
          torch(q.x + sx * (q.w / 2 - 4.2), q.z - q.d / 2 + 8 + k * 14);
        }
        addLight(0xff5a30, 1.6, 42, q.x, 6, q.z);
        for (let k = 0; k < 14; k++) prop('bones', q.x + (r() - 0.5) * (q.w - 4), q.z + (r() - 0.5) * (q.d - 4), 1.2);
        break;
      }

      // ---- Nebeltempel ----
      case 'temple':
        for (const sx of [-1, 1]) for (const sz of [-1, 1]) pillar(q.x + sx * (q.w / 2 - 3), q.z + sz * (q.d / 2 - 3), 0.9, 0x78b8a0);
        for (let k = 0; k < 4; k++) { const rn = new THREE.Mesh(new THREE.RingGeometry(1.4, 1.7, 6), B(0x60ffb8, { side: THREE.DoubleSide, transparent: true, opacity: 0.5 })); rn.rotation.x = -Math.PI / 2; rn.position.set(q.x + (r() - 0.5) * (q.w - 8), 0.05, q.z + (r() - 0.5) * (q.d - 8)); scene.add(rn); }
        break;
      case 'flooded': {
        const w = new THREE.Mesh(new THREE.PlaneGeometry(q.w - 1.5, q.d - 1.5), new THREE.MeshPhongMaterial({ color: 0x2a9a9a, transparent: true, opacity: 0.55, shininess: 90 }));
        w.rotation.x = -Math.PI / 2; w.position.set(q.x, 0.12, q.z); scene.add(w);
        anims.push((t) => { w.position.y = 0.12 + Math.sin(t * 1.3) * 0.03; });
        for (let k = 0; k < 7; k++) { const col = new THREE.Mesh(new THREE.CylinderGeometry(0.7, 0.8, 1.5 + r() * 4, 7), M(0x78a898)); col.position.set(q.x + (r() - 0.5) * (q.w - 6), 1, q.z + (r() - 0.5) * (q.d - 6)); scene.add(col); }
        for (let k = 0; k < 6; k++) { const cr = new THREE.Mesh(new THREE.OctahedronGeometry(0.5 + r() * 0.4), B(0x50ffd0)); cr.position.set(q.x + (r() - 0.5) * (q.w - 4), 0.9, q.z + (r() - 0.5) * (q.d - 4)); scene.add(cr); anims.push((t) => { cr.rotation.y = t + k; }); }
        for (let k = 0; k < 10; k++) { const pad = new THREE.Mesh(new THREE.CircleGeometry(0.6, 8), M(0x3f8a3a)); pad.rotation.x = -Math.PI / 2; pad.position.set(q.x + (r() - 0.5) * (q.w - 4), 0.16, q.z + (r() - 0.5) * (q.d - 4)); scene.add(pad); }
        addLight(0x40ffd0, 1.0, 30, q.x, 4, q.z);
        break;
      }
      case 'shrine': {
        const th = new THREE.Group();
        const seat = new THREE.Mesh(new THREE.BoxGeometry(4, 1.2, 4), M(0x2f9a78)); seat.position.y = 0.6; th.add(seat);
        const back = new THREE.Mesh(new THREE.BoxGeometry(5, 6, 0.8), M(0x3fb890)); back.position.set(0, 3.7, -1.8); th.add(back);
        for (const sx of [-1, 1]) { const cobra = new THREE.Mesh(new THREE.ConeGeometry(0.8, 4.5, 6), M(0x2a8a60)); cobra.position.set(sx * 3.4, 2.3, -1.5); th.add(cobra); }
        th.position.set(q.x, 0, q.z - q.d / 2 + 5); scene.add(th);
        const rune = new THREE.Mesh(new THREE.RingGeometry(8, 9, 40), B(0x30ffa0, { side: THREE.DoubleSide, transparent: true, opacity: 0.6 }));
        rune.rotation.x = -Math.PI / 2; rune.position.set(q.x, 0.05, q.z + 2); scene.add(rune);
        anims.push((t) => { rune.material.opacity = 0.45 + Math.sin(t * 2) * 0.2; });
        for (const sx of [-1, 1]) for (let k = 0; k < 3; k++) { pillar(q.x + sx * (q.w / 2 - 3), q.z - q.d / 2 + 8 + k * 14, 0.95, 0x78b8a0); torch(q.x + sx * (q.w / 2 - 4.2), q.z - q.d / 2 + 8 + k * 14); }
        addLight(0x30ffa0, 1.6, 42, q.x, 6, q.z);
        break;
      }

      // ---- Feuerschlund ----
      case 'forge': {
        for (const sx of [-1, 1]) for (const sz of [-1, 1]) pillar(q.x + sx * (q.w / 2 - 3), q.z + sz * (q.d / 2 - 3), 0.9, 0x3a2a2a);
        for (const sx of [-1, 1]) { const ch = new THREE.Mesh(new THREE.BoxGeometry(2, 0.12, q.d - 8), B(0xff5a1a)); ch.position.set(q.x + sx * (q.w / 2 - 5.5), 0.06, q.z); scene.add(ch); }
        for (let k = 0; k < 3; k++) { const an = new THREE.Group(); const base = new THREE.Mesh(new THREE.BoxGeometry(1.4, 1.1, 2), M(0x2a2a30)); base.position.y = 0.55; an.add(base); const top = new THREE.Mesh(new THREE.BoxGeometry(2.2, 0.4, 1.0), M(0x3a3a44)); top.position.y = 1.3; an.add(top); an.position.set(q.x + (k - 1) * 7, 0, q.z + (r() - 0.5) * (q.d - 12)); scene.add(an); }
        addLight(0xff5a20, 1.4, 34, q.x, 4, q.z);
        break;
      }
      case 'lava':
        pool(q.x - q.w / 2 + 5, q.z - q.d / 4, 4.2, 0xff5a1a);
        pool(q.x + q.w / 2 - 5, q.z + q.d / 4, 4.8, 0xff5a1a);
        pool(q.x + q.w / 2 - 6, q.z - q.d / 3, 3.2, 0xff7a20);
        for (let k = 0; k < 8; k++) { const sp = new THREE.Mesh(new THREE.ConeGeometry(0.5, 2 + r() * 2.4, 5), M(0x14141c)); sp.position.set(q.x + (r() - 0.5) * (q.w - 6), 1.2, q.z + (r() - 0.5) * (q.d - 6)); scene.add(sp); }
        break;
      case 'obsidian':
        for (let k = 0; k < 12; k++) { const cr = new THREE.Mesh(new THREE.OctahedronGeometry(0.7 + r() * 0.8), M(0x14141c)); cr.scale.y = 2; cr.position.set(q.x + (r() - 0.5) * (q.w - 4), 1.2, q.z + (r() - 0.5) * (q.d - 4)); scene.add(cr); }
        for (let k = 0; k < 6; k++) { const cr = new THREE.Mesh(new THREE.OctahedronGeometry(0.4), B(0xff6a20)); cr.position.set(q.x + (r() - 0.5) * (q.w - 6), 0.7, q.z + (r() - 0.5) * (q.d - 6)); scene.add(cr); anims.push((t) => { cr.rotation.y = t * 1.5 + k; }); }
        for (const sx of [-1, 1]) for (let k = 0; k < 2; k++) pillar(q.x + sx * (q.w / 2 - 3), q.z - q.d / 4 + k * 16, 0.9, 0x1a1a24);
        pool(q.x, q.z + q.d / 2 - 5, 3.4, 0xff5a1a);
        break;
      case 'inferno': {
        const th = new THREE.Group();
        const seat = new THREE.Mesh(new THREE.BoxGeometry(3.4, 1.4, 3.4), M(0x1a1a24)); seat.position.y = 0.7; th.add(seat);
        const back = new THREE.Mesh(new THREE.BoxGeometry(3.4, 5.5, 0.8), M(0x20202c)); back.position.set(0, 3.8, -1.6); th.add(back);
        for (const sx of [-1, 1]) { const horn = new THREE.Mesh(new THREE.ConeGeometry(0.45, 3.2, 5), M(0x14141c)); horn.position.set(sx * 1.7, 7.2, -1.6); horn.rotation.z = -sx * 0.4; th.add(horn); }
        th.position.set(q.x, 0, q.z - q.d / 2 + 4); scene.add(th);
        const rune = new THREE.Mesh(new THREE.RingGeometry(8, 9, 40), B(0xff3a10, { side: THREE.DoubleSide, transparent: true, opacity: 0.7 }));
        rune.rotation.x = -Math.PI / 2; rune.position.set(q.x, 0.05, q.z + 2); scene.add(rune);
        anims.push((t) => { rune.material.opacity = 0.5 + Math.sin(t * 3) * 0.25; });
        for (const sx of [-1, 1]) { pool(q.x + sx * (q.w / 2 - 4), q.z - q.d / 6, 4.2, 0xff5a1a); pool(q.x + sx * (q.w / 2 - 4), q.z + q.d / 4, 3.6, 0xff6a20); }
        for (const sx of [-1, 1]) for (let k = 0; k < 3; k++) { pillar(q.x + sx * (q.w / 2 - 3), q.z - q.d / 2 + 8 + k * 14, 0.95, 0x1a1a24); torch(q.x + sx * (q.w / 2 - 4.2), q.z - q.d / 2 + 8 + k * 14); }
        addLight(0xff4a20, 1.8, 46, q.x, 6, q.z);
        break;
      }
    }
  });
  const playerLight = new THREE.PointLight(ST.player, 1.1, 30);
  scene.add(playerLight);

  // walkable rectangles in world coordinates, for collision
  const rects = rooms.map((q) => ({ x0: q.x - q.w / 2 + 0.6, x1: q.x + q.w / 2 - 0.6, z0: q.z - q.d / 2 - 0.4, z1: q.z + q.d / 2 + 0.4 }));
  function collide(x, z) {
    for (const q of rects) if (x >= q.x0 && x <= q.x1 && z >= q.z0 && z <= q.z1) return [x, z];
    let best = [x, z], bd = Infinity;
    for (const q of rects) {
      const cx = Math.max(q.x0, Math.min(q.x1, x)), cz = Math.max(q.z0, Math.min(q.z1, z));
      const dd = (cx - x) ** 2 + (cz - z) ** 2;
      if (dd < bd) { bd = dd; best = [cx, cz]; }
    }
    return best;
  }

  return {
    scene, collide, playerLight,
    update(t, px, pz) {
      playerLight.position.set(px, 3.5, pz);
      flames.forEach((f, i) => f.scale.setScalar(0.85 + Math.sin(t * 9 + i) * 0.15));
      lights.forEach((l, i) => { l.intensity = (l.userData.base ??= l.intensity) * (0.9 + Math.sin(t * 7 + i * 2) * 0.1); });
      for (const a of anims) a(t);
    },
  };
}
