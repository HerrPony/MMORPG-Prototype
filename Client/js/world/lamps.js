import * as THREE from 'three';
import { heightAt } from '../util.js';
import { windowMats } from './buildings.js';

// Street lamps and house lights that come on by themselves at night.
// Every lamp is a post with a glowing lantern and a soft halo; the nearest few light sources also get a real point
// light (the two nearest cast shadows) from a small pool, so the light falls on the ground and the walls.

const POOL = 7, SHADOW_LIGHTS = 2, REACH = 60, POST_H = 3.3;

function haloTexture() {
  const c = document.createElement('canvas'); c.width = c.height = 128;
  const g = c.getContext('2d');
  const grad = g.createRadialGradient(64, 64, 0, 64, 64, 64);
  grad.addColorStop(0, 'rgba(255,240,200,1)'); grad.addColorStop(0.12, 'rgba(255,200,120,0.85)');
  grad.addColorStop(0.4, 'rgba(255,150,70,0.28)'); grad.addColorStop(1, 'rgba(255,120,50,0)');
  g.fillStyle = grad; g.fillRect(0, 0, 128, 128);
  const t = new THREE.CanvasTexture(c); t.colorSpace = THREE.SRGBColorSpace;
  return t;
}

const styleColour = (s) => (s === 'jade' ? 0xff8d47 : s === 'obsidian' ? 0xff6624 : s === 'bamboo' ? 0xffb257 : 0xffb766);

export function buildLamps(scene, scenery) {
  const lamps = scenery.lamps || [];
  const sources = []; // { x, y, z, col, e, r, ph, halo }
  const root = new THREE.Group();
  scene.add(root);

  // ---- posts (instanced) ----
  const n = lamps.length;
  const poleG = new THREE.CylinderGeometry(0.08, 0.1, POST_H, 6); poleG.translate(0, POST_H / 2, 0);
  const baseG = new THREE.BoxGeometry(0.55, 0.4, 0.55); baseG.translate(0, 0.2, 0);
  const capG = new THREE.ConeGeometry(0.62, 0.45, 4); capG.rotateY(Math.PI / 4); capG.translate(0, POST_H + 0.82, 0);
  const headG = new THREE.BoxGeometry(0.5, 0.6, 0.5); headG.translate(0, POST_H + 0.3, 0);
  const iron = new THREE.MeshLambertMaterial({ color: 0x2e2a28 });
  const stone = new THREE.MeshLambertMaterial({ color: 0x5a5650 });
  const cap = new THREE.MeshLambertMaterial({ color: 0x2a2420 });
  const capJade = new THREE.MeshLambertMaterial({ color: 0xc0302a });
  const headMat = new THREE.MeshLambertMaterial({ color: 0x6a4a28, emissive: 0xffa040, emissiveIntensity: 0.08 });
  const pole = new THREE.InstancedMesh(poleG, iron, n), base = new THREE.InstancedMesh(baseG, stone, n);
  const head = new THREE.InstancedMesh(headG, headMat, n), caps = new THREE.InstancedMesh(capG, cap, n);
  for (const m of [pole, base, caps]) m.castShadow = true;
  const mat = new THREE.Matrix4();
  lamps.forEach((l, i) => {
    const y = heightAt(l.x, l.z) - 0.15;
    mat.makeTranslation(l.x, y, l.z);
    pole.setMatrixAt(i, mat); base.setMatrixAt(i, mat); head.setMatrixAt(i, mat); caps.setMatrixAt(i, mat);
    sources.push({ x: l.x, y: y + POST_H + 0.3, z: l.z, col: styleColour(l.s), e: 55, r: 17, ph: Math.random() * 6.28, halo: 1.6 });
  });
  for (const m of [pole, base, head, caps]) { m.instanceMatrix.needsUpdate = true; m.frustumCulled = false; root.add(m); }

  // ---- other sources: house windows, lanterns, tower windows ----
  for (const b of scenery.buildings) {
    const h = heightAt(b.x, b.z) - 0.3;
    if (b.k === 'house' || b.k === 'hut') {
      const lift = b.k === 'hut' ? 1.3 : 0;
      sources.push({ x: b.x, y: h + lift + b.h * 0.62, z: b.z + b.d / 2 + 1.2, col: b.s === 'obsidian' ? 0xff6a24 : 0xffcc80, e: 38, r: 12, ph: Math.random() * 6.28, halo: 0 });
    } else if (b.k === 'lantern') {
      sources.push({ x: b.x, y: h + b.h + 0.3, z: b.z, col: 0xffa047, e: 50, r: 15, ph: Math.random() * 6.28, halo: 1.6 });
    } else if (b.k === 'tower') {
      sources.push({ x: b.x, y: h + b.h * 0.6, z: b.z + b.w / 2 + 1, col: 0xffc880, e: 40, r: 14, ph: Math.random() * 6.28, halo: 0 });
    }
  }

  // ---- halos: additive sprites at the lamps ----
  const haloMat = new THREE.SpriteMaterial({ map: haloTexture(), blending: THREE.AdditiveBlending, depthWrite: false, transparent: true, opacity: 0, fog: false });
  const halos = [];
  for (const s of sources) {
    if (!s.halo) continue;
    const sp = new THREE.Sprite(haloMat);
    sp.position.set(s.x, s.y, s.z); sp.scale.setScalar(4.2);
    sp.userData.ph = s.ph;
    root.add(sp); halos.push(sp);
  }

  // ---- light pool ----
  const lights = [];
  for (let i = 0; i < POOL; i++) {
    const l = new THREE.PointLight(0xffb766, 0, 16, 2);
    if (i < SHADOW_LIGHTS) {
      l.castShadow = true; l.shadow.mapSize.set(512, 512); l.shadow.bias = -0.002; l.shadow.normalBias = 0.05; l.shadow.camera.near = 0.3;
    }
    l.userData.src = -1;
    root.add(l); lights.push(l);
  }

  let t = 0, timer = 0, near = [];
  function assign(cam) {
    near = [];
    for (let i = 0; i < sources.length; i++) {
      const s = sources[i];
      const d = Math.hypot(s.x - cam.x, s.y - cam.y, s.z - cam.z);
      if (d < REACH + 8) near.push([d, i]);
    }
    near.sort((a, b) => a[0] - b[0]);
    const take = near.slice(0, POOL);
    const info = new Map(take.map((e, r) => [e[1], e[0]]));   // source -> distance
    const taken = new Set(), free = [];
    lights.forEach((l, i) => {
      if (l.userData.src >= 0 && info.has(l.userData.src)) taken.add(l.userData.src);
      else { l.userData.src = -1; free.push(i); }
    });
    for (const [, si] of take) {
      if (taken.has(si) || !free.length) continue;
      const l = lights[free.pop()], s = sources[si];
      l.userData.src = si; l.position.set(s.x, s.y, s.z); l.color.setHex(s.col); l.distance = s.r;
      taken.add(si);
    }
    for (const l of lights) l.userData.fade = l.userData.src >= 0 ? 1 - smooth(REACH - 22, REACH + 8, info.get(l.userData.src)) : 0;
  }
  const smooth = (a, b, x) => { const k = Math.max(0, Math.min(1, (x - a) / (b - a))); return k * k * (3 - 2 * k); };

  return {
    // day: 0 night .. 1 day
    update(dt, cam, day) {
      t += dt;
      const on = Math.max(0, Math.min(1, (0.62 - day) / 0.4));
      headMat.emissiveIntensity = 0.08 + 2.4 * on;
      haloMat.opacity = 0.9 * on;
      for (const m of windowMats) m.emissiveIntensity = 0.1 + 1.6 * on;
      if (on <= 0.01) { for (const l of lights) l.intensity = 0; timer = 0; return; }
      timer -= dt;
      if (timer <= 0) { timer = 0.25; assign(cam); }
      for (const l of lights) {
        const si = l.userData.src;
        if (si < 0) { l.intensity = 0; continue; }
        const s = sources[si];
        l.intensity = s.e * on * (l.userData.fade ?? 1) * (1 + 0.07 * Math.sin(t * 9 + s.ph) + 0.04 * Math.sin(t * 23 + s.ph * 2));
      }
    },
  };
}
