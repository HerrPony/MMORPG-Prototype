import * as THREE from 'three';
import { makeProps } from './vegetation.js';

const M = (c, o = {}) => new THREE.MeshLambertMaterial({ color: c, ...o });
const WALL_H = 9;

function rng(seed) {
  let s = seed >>> 0;
  return () => { s = (s * 1664525 + 1013904223) >>> 0; return s / 4294967296; };
}

// Builds the dark mine: rooms and corridors given in local coordinates around origin.
export function buildDungeon(info) {
  const scene = new THREE.Scene();
  scene.background = new THREE.Color(0x07070b);
  scene.fog = new THREE.Fog(0x07070b, 8, 55);
  scene.add(new THREE.AmbientLight(0x6a6a88, 0.9));
  scene.add(new THREE.HemisphereLight(0x8888aa, 0x221a14, 0.5));
  const ox = info.origin[0], oz = info.origin[1];
  const rooms = info.rooms.map((r) => ({ ...r, x: r.x + ox, z: r.z + oz }));
  const rock = M(0x4a443e), floorM = M(0x3a3632), dark = M(0x2a2622);
  const r = rng(77);
  const P = makeProps();

  const inside = (x, z, skip) => rooms.some((q, i) => i !== skip && x > q.x - q.w / 2 + 0.01 && x < q.x + q.w / 2 - 0.01 && z > q.z - q.d / 2 - 0.01 && z < q.z + q.d / 2 + 0.01);
  const wallPieces = [];
  const flames = [];

  rooms.forEach((q, i) => {
    const fl = new THREE.Mesh(new THREE.BoxGeometry(q.w, 0.4, q.d), q.theme === 'throne' ? M(0x2a2a3a) : q.theme === 'web' ? M(0x2e2e2a) : floorM);
    fl.position.set(q.x, -0.2, q.z); fl.receiveShadow = true; scene.add(fl);
    const ce = new THREE.Mesh(new THREE.PlaneGeometry(q.w, q.d), M(0x1c1a18, { side: THREE.BackSide }));
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
    c.setHSL(0.08, 0.1, 0.22 + r() * 0.1);
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
    const fl = new THREE.Mesh(new THREE.SphereGeometry(0.28, 6, 5), new THREE.MeshBasicMaterial({ color: 0xffaa33 })); fl.position.set(x, 2.4, z); scene.add(fl);
    flames.push(fl);
  };

  const lights = [];
  rooms.forEach((q, i) => {
    // torches along the long walls
    const n = Math.max(2, Math.floor(q.d / 14));
    for (let k = 0; k < n; k++) {
      const z = q.z - q.d / 2 + (k + 0.5) * (q.d / n);
      if (q.w > 8) { torch(q.x - q.w / 2 + 1.2, z); torch(q.x + q.w / 2 - 1.2, z); }
    }
    if (q.w > 8 && i % 2 === 0) {
      const l = new THREE.PointLight(0xff9a44, 1.3, 38); l.position.set(q.x, 5, q.z); scene.add(l); lights.push(l);
    }
    switch (q.theme) {
      case 'hall':
        for (const sx of [-1, 1]) for (const sz of [-1, 1]) {
          const p = new THREE.Mesh(new THREE.CylinderGeometry(0.7, 0.8, WALL_H, 8), M(0x6a645a)); p.position.set(q.x + sx * (q.w / 2 - 3), WALL_H / 2, q.z + sz * (q.d / 2 - 3)); scene.add(p);
        }
        break;
      case 'mine':
        for (let k = 0; k < 12; k++) prop('stalag', q.x + (r() - 0.5) * (q.w - 6), q.z + (r() - 0.5) * (q.d - 6), 0.6 + r() * 0.8, [0.08, 0.05, 0.3]);
        for (let k = 0; k < 6; k++) {
          const cr = new THREE.Mesh(new THREE.OctahedronGeometry(0.5 + r() * 0.5), new THREE.MeshBasicMaterial({ color: [0x4acfff, 0xff6aff, 0x6aff9a][k % 3] }));
          cr.position.set(q.x + (r() - 0.5) * (q.w - 4), 0.8, q.z + (r() - 0.5) * (q.d - 4)); scene.add(cr);
        }
        for (let k = 0; k < 4; k++) { const cart = new THREE.Mesh(new THREE.BoxGeometry(1.6, 0.9, 1.1), M(0x5a4030)); cart.position.set(q.x + (r() - 0.5) * (q.w - 6), 0.5, q.z + (r() - 0.5) * (q.d - 6)); cart.rotation.y = r() * 3; scene.add(cart); }
        break;
      case 'web':
        for (let k = 0; k < 14; k++) {
          const web = new THREE.Mesh(new THREE.PlaneGeometry(4 + r() * 4, 4 + r() * 3), new THREE.MeshBasicMaterial({ color: 0xdddde8, transparent: true, opacity: 0.35, side: THREE.DoubleSide, depthWrite: false }));
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
        const rune = new THREE.Mesh(new THREE.RingGeometry(8, 9, 40), new THREE.MeshBasicMaterial({ color: 0xc03030, side: THREE.DoubleSide, transparent: true, opacity: 0.6 }));
        rune.rotation.x = -Math.PI / 2; rune.position.set(q.x, 0.05, q.z + 2); scene.add(rune);
        for (const sx of [-1, 1]) for (let k = 0; k < 3; k++) {
          const p = new THREE.Mesh(new THREE.CylinderGeometry(0.9, 1, WALL_H, 8), M(0x4a4640)); p.position.set(q.x + sx * (q.w / 2 - 3), WALL_H / 2, q.z - q.d / 2 + 8 + k * 14); scene.add(p);
          torch(q.x + sx * (q.w / 2 - 4.2), q.z - q.d / 2 + 8 + k * 14);
        }
        const l = new THREE.PointLight(0xff5a30, 1.6, 42); l.position.set(q.x, 6, q.z); scene.add(l); lights.push(l);
        for (let k = 0; k < 14; k++) prop('bones', q.x + (r() - 0.5) * (q.w - 4), q.z + (r() - 0.5) * (q.d - 4), 1.2);
        break;
      }
    }
  });
  const playerLight = new THREE.PointLight(0xffd9a0, 1.1, 30);
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
      lights.forEach((l, i) => { l.intensity = 1.2 + Math.sin(t * 7 + i * 2) * 0.15; });
    },
  };
}
