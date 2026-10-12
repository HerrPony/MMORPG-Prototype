import * as THREE from 'three';
import { W, heightAt, biomeAt } from '../util.js';

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
// Draws the props the server placed (sc.pk = kind names, sc.props = packed list, see Server gamedata/props.go) and
// returns the collision grid (8x8 m cells of {x,z,r}) built from the same list. The server decides what stands
// where; nothing is generated here.
export function buildVegetation(scene, sc) {
  const K = parts();
  const chunks = new Map(); // key -> kind -> [{x,y,z,s,ry,biome}]
  const grid = {};          // collision cells
  const bin = Uint8Array.from(atob(sc.props), (c) => c.charCodeAt(0));
  for (let o = 0; o + 8 <= bin.length; o += 8) {
    const kind = sc.pk[bin[o]];
    const x = W.x0 + (bin[o + 1] | (bin[o + 2] << 8)) / 40, z = W.z0 + (bin[o + 3] | (bin[o + 4] << 8)) / 40;
    const s = bin[o + 5] / 100, ry = bin[o + 6] / 255 * 6.28, solid = bin[o + 7] / 100;
    const key = Math.floor(x / 64) + ',' + Math.floor(z / 64);
    let c = chunks.get(key);
    if (!c) { c = {}; chunks.set(key, c); }
    (c[kind] ||= []).push({ x, z, s, ry, biome: biomeAt(x, z) });
    if (solid > 0) (grid[Math.floor(x / 8) + ',' + Math.floor(z / 8)] ||= []).push({ x, z, r: solid * s });
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
