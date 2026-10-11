import * as THREE from 'three';
import { W, heightAt, biomeAt, distToPoly } from '../util.js';

function rng(seed) {
  let s = seed >>> 0;
  return () => { s = (s * 1664525 + 1013904223) >>> 0; return s / 4294967296; };
}

const mat = (c) => new THREE.MeshLambertMaterial({ color: c });
const white = new THREE.MeshLambertMaterial({ color: 0xffffff });

// Each kind is a list of {geo, mat, color?} parts sharing one transform per instance.
function parts() {
  const K = {};
  const trunk = (r1, r2, h, seg = 5) => new THREE.CylinderGeometry(r1, r2, h, seg).translate(0, h / 2, 0);
  K.pine = [
    { geo: trunk(0.18, 0.3, 3.2), mat: mat(0x5a3d22) },
    { geo: new THREE.ConeGeometry(2.0, 2.6, 6).translate(0, 3.4, 0), mat: white, tint: 'leaf' },
    { geo: new THREE.ConeGeometry(1.55, 2.3, 6).translate(0, 4.7, 0), mat: white, tint: 'leaf' },
    { geo: new THREE.ConeGeometry(1.05, 2.0, 6).translate(0, 5.9, 0), mat: white, tint: 'leaf' },
  ];
  K.oak = [
    { geo: trunk(0.25, 0.4, 2.6), mat: mat(0x6b4a2a) },
    { geo: new THREE.IcosahedronGeometry(2.0, 0).translate(0, 4.0, 0), mat: white, tint: 'leaf' },
    { geo: new THREE.IcosahedronGeometry(1.4, 0).translate(1.2, 3.2, 0.5), mat: white, tint: 'leaf' },
  ];
  K.dead = [
    { geo: trunk(0.12, 0.3, 4.2), mat: mat(0x3a3226) },
    { geo: new THREE.CylinderGeometry(0.05, 0.12, 2.0, 4).rotateZ(0.9).translate(0.9, 3.3, 0), mat: mat(0x3a3226) },
    { geo: new THREE.CylinderGeometry(0.04, 0.1, 1.6, 4).rotateZ(-1.0).translate(-0.7, 2.6, 0.2), mat: mat(0x3a3226) },
  ];
  K.snowpine = [
    { geo: trunk(0.18, 0.28, 2.4), mat: mat(0x4a3a2c) },
    { geo: new THREE.ConeGeometry(1.8, 2.6, 6).translate(0, 3.0, 0), mat: mat(0x2f5a46) },
    { geo: new THREE.ConeGeometry(1.45, 2.2, 6).translate(0, 4.4, 0), mat: mat(0x3a6a52) },
    { geo: new THREE.ConeGeometry(0.95, 1.8, 6).translate(0, 5.5, 0), mat: mat(0xe8f0f8) },
    { geo: new THREE.ConeGeometry(1.7, 0.8, 6).translate(0, 3.7, 0), mat: mat(0xeaf2fa) },
  ];
  K.cactus = [
    { geo: new THREE.CylinderGeometry(0.35, 0.4, 3.0, 7).translate(0, 1.5, 0), mat: mat(0x4f8a3a) },
    { geo: new THREE.CylinderGeometry(0.2, 0.22, 1.2, 6).rotateZ(1.2).translate(0.7, 1.7, 0), mat: mat(0x4f8a3a) },
    { geo: new THREE.CylinderGeometry(0.2, 0.22, 1.0, 6).rotateZ(-1.2).translate(-0.7, 2.1, 0), mat: mat(0x4f8a3a) },
  ];
  K.palm = [
    { geo: new THREE.CylinderGeometry(0.2, 0.3, 5.0, 6).translate(0, 2.5, 0), mat: mat(0x7a5a32) },
    ...[0, 1, 2, 3, 4, 5].map((i) => ({ geo: new THREE.ConeGeometry(0.5, 3.2, 4).rotateZ(1.3).rotateY((i / 6) * Math.PI * 2).translate(Math.cos((i / 6) * 6.283) * 1.4, 5.1, Math.sin((i / 6) * 6.283) * 1.4), mat: mat(0x3f8a3a) })),
  ];
  K.rock = [{ geo: new THREE.DodecahedronGeometry(0.8, 0), mat: white, tint: 'rock' }];
  K.bush = [{ geo: new THREE.SphereGeometry(0.8, 7, 5).translate(0, 0.45, 0), mat: white, tint: 'leaf' }];
  K.reed = [0, 1, 2, 3].map((i) => ({ geo: new THREE.ConeGeometry(0.06, 1.6, 4).rotateZ((i - 1.5) * 0.12).translate(i * 0.12 - 0.2, 0.8, (i % 2) * 0.1), mat: mat(0x6a7a3a) }));
  K.mushroom = [
    { geo: new THREE.CylinderGeometry(0.08, 0.1, 0.5, 5).translate(0, 0.25, 0), mat: mat(0xe8dcc0) },
    { geo: new THREE.SphereGeometry(0.4, 7, 4, 0, 6.283, 0, 1.4).translate(0, 0.5, 0), mat: mat(0xb03a30) },
  ];
  K.tuft = [{ geo: new THREE.ConeGeometry(0.07, 0.38, 4).translate(0, 0.19, 0), mat: white, tint: 'grass' }];
  K.flower = [{ geo: new THREE.SphereGeometry(0.09, 5, 4).translate(0, 0.35, 0), mat: white, tint: 'flower' }];
  K.bones = [
    { geo: new THREE.CylinderGeometry(0.05, 0.05, 0.9, 4).rotateZ(1.5).translate(0, 0.05, 0), mat: mat(0xe6dcc0) },
    { geo: new THREE.SphereGeometry(0.2, 6, 5).translate(0.6, 0.15, 0.1), mat: mat(0xe6dcc0) },
  ];
  K.stalag = [{ geo: new THREE.ConeGeometry(0.7, 3.5, 6).translate(0, 1.75, 0), mat: white, tint: 'rock' }];

  // ---- second continent ----
  const glowM = new THREE.MeshBasicMaterial({ color: 0xff6a20 });
  K.jtree = [
    { geo: trunk(0.3, 0.5, 7.5, 6), mat: mat(0x5a4028) },
    { geo: new THREE.IcosahedronGeometry(3.4, 0).scale(1, 0.45, 1).translate(0, 8.2, 0), mat: white, tint: 'leaf' },
    { geo: new THREE.IcosahedronGeometry(2.4, 0).scale(1, 0.5, 1).translate(2.2, 6.6, 1), mat: white, tint: 'leaf' },
    { geo: new THREE.CylinderGeometry(0.05, 0.05, 4, 4).translate(1.8, 5.8, 0.4), mat: mat(0x2f7a2a) },
  ];
  K.fern = [0, 1, 2, 3, 4].map((i) => ({ geo: new THREE.ConeGeometry(0.22, 1.7, 4).rotateZ(1.0).rotateY(i * 1.26).translate(Math.cos(i * 1.26) * 0.55, 0.5, Math.sin(i * 1.26) * 0.55), mat: white, tint: 'leaf' }));
  K.bamboo = [0, 1, 2, 3, 4, 5].map((i) => ({ geo: new THREE.CylinderGeometry(0.07, 0.09, 6 + (i % 3) * 1.4, 5).translate(Math.cos(i * 1.9) * 0.6, 3.2 + (i % 3) * 0.7, Math.sin(i * 1.9) * 0.6), mat: mat(0x8cb24a) }))
    .concat([0, 1, 2].map((i) => ({ geo: new THREE.ConeGeometry(0.6, 1.2, 5).translate(Math.cos(i * 2.1) * 0.5, 6.4 + i * 0.5, Math.sin(i * 2.1) * 0.5), mat: mat(0x5a9a30) })));
  K.sakura = [
    { geo: trunk(0.2, 0.36, 3.6, 6), mat: mat(0x4a3228) },
    { geo: new THREE.IcosahedronGeometry(2.4, 0).translate(0, 4.8, 0), mat: white, tint: 'leaf' },
    { geo: new THREE.IcosahedronGeometry(1.6, 0).translate(1.5, 4.0, 0.7), mat: white, tint: 'leaf' },
    { geo: new THREE.IcosahedronGeometry(1.4, 0).translate(-1.3, 4.3, -0.6), mat: white, tint: 'leaf' },
  ];
  K.mesapillar = [
    { geo: new THREE.CylinderGeometry(1.5, 2.0, 3.2, 7).translate(0, 1.6, 0), mat: mat(0xb0502e) },
    { geo: new THREE.CylinderGeometry(1.0, 1.4, 2.6, 7).translate(0.1, 4.4, 0), mat: mat(0xc86a3c) },
    { geo: new THREE.CylinderGeometry(1.6, 0.9, 1.2, 7).translate(0.15, 6.3, 0), mat: mat(0xd8803f) },
  ];
  K.drybush = [{ geo: new THREE.IcosahedronGeometry(0.7, 0).scale(1.2, 0.7, 1.2).translate(0, 0.4, 0), mat: mat(0x8a6a38) }];
  K.charred = [
    { geo: trunk(0.1, 0.28, 4.4, 5), mat: mat(0x1c1618) },
    { geo: new THREE.CylinderGeometry(0.04, 0.1, 1.8, 4).rotateZ(1.0).translate(0.8, 3.2, 0), mat: mat(0x1c1618) },
    { geo: new THREE.CylinderGeometry(0.04, 0.09, 1.4, 4).rotateZ(-0.9).translate(-0.6, 2.6, 0.2), mat: mat(0x1c1618) },
  ];
  K.obsidian = [0, 1, 2].map((i) => ({ geo: new THREE.ConeGeometry(0.45, 3.2 - i * 0.7, 5).rotateZ((i - 1) * 0.25).translate((i - 1) * 0.7, (3.2 - i * 0.7) / 2, i * 0.2), mat: mat(0x14141c) }));
  K.ember = [{ geo: new THREE.IcosahedronGeometry(0.25, 0).translate(0, 0.25, 0), mat: glowM }, { geo: new THREE.IcosahedronGeometry(0.18, 0).translate(0.4, 0.15, 0.2), mat: glowM }];
  return K;
}

const TINTS = {
  leaf: [() => new THREE.Color().setHSL(0.31 + Math.random() * 0.03, 0.45, 0.2 + Math.random() * 0.08)],
  rock: [() => new THREE.Color().setHSL(0.08, 0.05, 0.38 + Math.random() * 0.15)],
  grass: [() => new THREE.Color().setHSL(0.26 + Math.random() * 0.07, 0.5, 0.25 + Math.random() * 0.15)],
  flower: [() => new THREE.Color([0xff6aa0, 0xffe14a, 0xffffff, 0xa07aff][Math.floor(Math.random() * 4)])],
};

const BIOME_LEAF = {
  6: () => new THREE.Color().setHSL(0.33 + Math.random() * 0.05, 0.65, 0.2 + Math.random() * 0.09),
  7: () => new THREE.Color().setHSL(0.93 + Math.random() * 0.05, 0.6, 0.74 + Math.random() * 0.1),
  8: () => new THREE.Color().setHSL(0.1 + Math.random() * 0.03, 0.35, 0.28 + Math.random() * 0.06),
  9: () => new THREE.Color().setHSL(0.04, 0.2, 0.1 + Math.random() * 0.05),
  1: () => new THREE.Color().setHSL(0.31 + Math.random() * 0.03, 0.45, 0.17 + Math.random() * 0.08),
  2: () => new THREE.Color().setHSL(0.24 + Math.random() * 0.04, 0.5, 0.26 + Math.random() * 0.08),
  3: () => new THREE.Color().setHSL(0.22 + Math.random() * 0.04, 0.35, 0.16 + Math.random() * 0.06),
  4: () => new THREE.Color().setHSL(0.3, 0.3, 0.3),
  5: () => new THREE.Color().setHSL(0.18, 0.4, 0.3),
};

// Builds all vegetation (instanced, in chunks) and returns the tree collision grid.
export function buildVegetation(scene, sc) {
  const K = parts();
  const r = rng(2024);
  const chunks = new Map(); // key -> kind -> [{x,y,z,s,ry,biome}]
  const grid = {};          // collision cells
  const clears = sc.clears;
  const roads = sc.roads;
  const rivers = sc.rivers;
  const lakes = sc.lakes;
  const STEP = 3;

  const blocked = (x, z) => {
    for (const c of clears) if (Math.hypot(x - c.x, z - c.z) < c.r + 1.2) return true;
    for (const rd of roads) if (distToPoly(x, z, rd) < 3.2) return true;
    for (const b of sc.buildings) if (b.k !== 'wall' && Math.abs(x - b.x) < b.w / 2 + 2 && Math.abs(z - b.z) < b.d / 2 + 2) return true;
    for (const b of sc.house) if (Math.abs(x - b.x) < b.w / 2 + 3 && Math.abs(z - b.z) < b.d / 2 + 3) return true;
    return false;
  };

  const add = (kind, x, z, s, biome, solid) => {
    const key = Math.floor(x / 64) + ',' + Math.floor(z / 64);
    let c = chunks.get(key);
    if (!c) { c = {}; chunks.set(key, c); }
    (c[kind] ||= []).push({ x, z, s, ry: r() * 6.28, biome });
    if (solid) {
      const gk = Math.floor(x / 8) + ',' + Math.floor(z / 8);
      (grid[gk] ||= []).push({ x, z, r: solid * s });
    }
  };

  const xEnd = W.x0 + (W.nx - 1) * W.cell - 4, zEnd = W.z0 + (W.nz - 1) * W.cell - 4;
  for (let z = W.z0 + 4; z < zEnd; z += STEP) {
    for (let x = W.x0 + 4; x < xEnd; x += STEP) {
      const px = x + (r() - 0.5) * STEP, pz = z + (r() - 0.5) * STEP;
      const b = biomeAt(px, pz);
      if (b === 0 || b === 10) continue;
      const h = heightAt(px, pz);
      if (h < 0.35) {
        if (b === 3 && h > -1.0 && h < 0.35 && r() < 0.12 && !blocked(px, pz)) add('reed', px, pz, 0.9 + r() * 0.8, b, 0);
        continue;
      }
      const q = r();
      if (q > 0.5) continue; // quick reject: most cells stay empty
      if (blocked(px, pz)) continue;
      const slope = Math.abs(heightAt(px + 2, pz) - h) + Math.abs(heightAt(px, pz + 2) - h);
      const s = 0.8 + r() * 0.9;
      const u = r();
      switch (b) {
        case 1: // forest
          if (u < 0.55) add(r() < 0.72 ? 'pine' : 'oak', px, pz, s, b, 0.5);
          else if (u < 0.62) add('bush', px, pz, s, b, 0.4);
          else if (u < 0.66) add('rock', px, pz, s * 0.8, b, 0.6);
          else if (u < 0.90) add('tuft', px, pz, s, b, 0);
          else add('flower', px, pz, s, b, 0);
          break;
        case 2: // plains
          if (u < 0.08) add('oak', px, pz, s, b, 0.5);
          else if (u < 0.12) add('bush', px, pz, s, b, 0.4);
          else if (u < 0.40) add('tuft', px, pz, s, b, 0);
          else if (u < 0.55) add('flower', px, pz, s, b, 0);
          break;
        case 3: // swamp
          if (u < 0.22) add('dead', px, pz, s, b, 0.4);
          else if (u < 0.32) add('mushroom', px, pz, s, b, 0);
          else if (u < 0.38) add('bush', px, pz, s, b, 0.4);
          else if (u < 0.55) add('tuft', px, pz, s, b, 0);
          break;
        case 4: // snow
          if (h > 24 || slope > 3.2) { if (u < 0.25) add('rock', px, pz, s * 1.4, b, 0.9); }
          else if (u < 0.30) add('snowpine', px, pz, s, b, 0.5);
          else if (u < 0.34) add('rock', px, pz, s * 1.1, b, 0.8);
          break;
        case 5: // desert
          if (u < 0.06) add('cactus', px, pz, s, b, 0.5);
          else if (u < 0.085) add('rock', px, pz, s, b, 0.7);
          else if (u < 0.10 && h < 1.8) add('palm', px, pz, s, b, 0.4);
          else if (u < 0.105) add('bones', px, pz, s, b, 0);
          else if (u < 0.15) add('bush', px, pz, s * 0.7, b, 0);
          break;
        case 6: // jungle: dense and tall
          if (u < 0.46) add('jtree', px, pz, s * 1.05, b, 0.7);
          else if (u < 0.51 && h < 2.2) add('palm', px, pz, s, b, 0.4);
          else if (u < 0.68) add('fern', px, pz, s, b, 0);
          else if (u < 0.76) add('bamboo', px, pz, s * 0.9, b, 0.4);
          else if (u < 0.82) add('bush', px, pz, s, b, 0.4);
          else if (u < 0.88) add('flower', px, pz, s, b, 0);
          else if (u < 0.9) add('rock', px, pz, s, b, 0.6);
          break;
        case 7: // cherry blossom groves
          if (u < 0.2) add('sakura', px, pz, s, b, 0.5);
          else if (u < 0.28) add('bamboo', px, pz, s, b, 0.4);
          else if (u < 0.32) add('bush', px, pz, s * 0.8, b, 0.4);
          else if (u < 0.34) add('rock', px, pz, s * 0.8, b, 0.6);
          else if (u < 0.58) add('tuft', px, pz, s, b, 0);
          else if (u < 0.72) add('flower', px, pz, s, b, 0);
          break;
        case 8: // red canyons
          if (u < 0.045 && slope < 2.5) add('mesapillar', px, pz, s * 1.2, b, 1.3);
          else if (u < 0.12) add('rock', px, pz, s * 1.1, b, 0.8);
          else if (u < 0.16) add('cactus', px, pz, s, b, 0.5);
          else if (u < 0.26) add('drybush', px, pz, s, b, 0);
          else if (u < 0.28) add('bones', px, pz, s, b, 0);
          break;
        case 9: // volcanic ash
          if (u < 0.12 && slope < 3) add('charred', px, pz, s, b, 0.4);
          else if (u < 0.25) add('rock', px, pz, s * 1.2, b, 0.8);
          else if (u < 0.32) add('obsidian', px, pz, s, b, 0.7);
          else if (u < 0.36) add('ember', px, pz, s, b, 0);
          break;
      }
    }
  }

  const dummy = new THREE.Object3D();
  const col = new THREE.Color();
  const group = new THREE.Group();
  for (const [key, kinds] of chunks) {
    for (const [kind, list] of Object.entries(kinds)) {
      for (const part of K[kind]) {
        const im = new THREE.InstancedMesh(part.geo, part.mat, list.length);
        list.forEach((it, i) => {
          dummy.position.set(it.x, heightAt(it.x, it.z) - 0.1, it.z);
          dummy.rotation.set(0, it.ry, 0);
          dummy.scale.set(it.s, it.s * (0.9 + (i % 3) * 0.1), it.s);
          dummy.updateMatrix();
          im.setMatrixAt(i, dummy.matrix);
          if (part.tint) {
            if (part.tint === 'leaf') col.copy(BIOME_LEAF[it.biome]());
            else col.copy(TINTS[part.tint][0]());
            if (part.tint === 'rock' && it.biome === 4) col.offsetHSL(0, 0, 0.25);
            if (part.tint === 'rock' && it.biome === 5) col.setHSL(0.1, 0.35, 0.5);
            if (part.tint === 'rock' && it.biome === 8) col.setHSL(0.04, 0.55, 0.3 + Math.random() * 0.12);
            if (part.tint === 'rock' && it.biome === 9) col.setHSL(0.02, 0.2, 0.08 + Math.random() * 0.06);
            if (part.tint === 'rock' && it.biome === 7) col.setHSL(0.08, 0.1, 0.5);
            if (part.tint === 'grass' && it.biome === 5) col.setHSL(0.14, 0.4, 0.45);
            im.setColorAt(i, col);
          }
        });
        im.castShadow = kind !== 'tuft' && kind !== 'flower';
        im.computeBoundingSphere();
        group.add(im);
      }
    }
  }
  scene.add(group);
  return grid;
}

// Dungeon decoration reuses the same geometry kinds.
export function makeProps() { return parts(); }
