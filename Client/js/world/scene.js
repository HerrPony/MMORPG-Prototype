import * as THREE from 'three';
import { W, setWorldData, heightAt, biomeAt, distToPoly, smoothstep, BIOME_KEYS } from '../util.js';
import { buildVegetation } from './vegetation.js';
import { buildTowns, buildCabin, ships } from './buildings.js';
import { buildDungeon } from './dungeon.js';

function rng(seed) {
  let s = seed >>> 0;
  return () => { s = (s * 1664525 + 1013904223) >>> 0; return s / 4294967296; };
}

function noiseTexture() {
  const c = document.createElement('canvas');
  c.width = c.height = 256;
  const g = c.getContext('2d');
  g.fillStyle = '#b8b8b8';
  g.fillRect(0, 0, 256, 256);
  const r = rng(7);
  for (let i = 0; i < 5000; i++) {
    const v = 150 + Math.floor(r() * 100);
    g.fillStyle = `rgba(${v},${v},${v},0.5)`;
    g.fillRect(r() * 256, r() * 256, 1 + r() * 3, 1 + r() * 3);
  }
  const t = new THREE.CanvasTexture(c);
  t.wrapS = t.wrapT = THREE.RepeatWrapping;
  t.repeat.set(160, 160);
  t.colorSpace = THREE.SRGBColorSpace;
  return t;
}

// Base colors per biome (index = biome id) and the sky / fog look when standing in it.
const BIOME_COLOR = [
  [0.76, 0.7, 0.5],     // sea floor / beach
  [0.2, 0.42, 0.13],    // forest
  [0.46, 0.58, 0.2],    // plains
  [0.27, 0.34, 0.17],   // swamp
  [0.9, 0.94, 1.0],     // snow
  [0.82, 0.69, 0.4],    // desert
  [0.07, 0.5, 0.13],    // jungle
  [0.74, 0.5, 0.58],    // blossom
  [0.78, 0.34, 0.18],   // mesa
  [0.17, 0.14, 0.15],   // volcano
  [1.0, 0.36, 0.04],    // lava
];
export const ZONE_LOOK = {
  forest: { fog: 0xbcd6e6, top: 0x4a86d6, bot: 0xbcd6e6, near: 45, far: 170, sun: 1.5, hemi: 0.85 },
  plains: { fog: 0xcfe3f0, top: 0x5a9ae0, bot: 0xd8e8f2, near: 70, far: 230, sun: 1.7, hemi: 0.95 },
  swamp: { fog: 0x8a9a7a, top: 0x6a7a68, bot: 0x9aa890, near: 22, far: 110, sun: 0.8, hemi: 0.6 },
  snow: { fog: 0xdce6f0, top: 0x8aa8cc, bot: 0xe6eef6, near: 35, far: 150, sun: 1.3, hemi: 1.0 },
  desert: { fog: 0xf0dcb0, top: 0x6aa0d8, bot: 0xf2e4c0, near: 55, far: 210, sun: 1.9, hemi: 0.95 },
  jungle: { fog: 0x8fd8a8, top: 0x1f9a8a, bot: 0xa8e8c0, near: 22, far: 125, sun: 1.2, hemi: 0.8 },
  blossom: { fog: 0xf4d6e4, top: 0xb98ad0, bot: 0xffdbe6, near: 50, far: 200, sun: 1.6, hemi: 1.05 },
  mesa: { fog: 0xe8b88c, top: 0xd8804e, bot: 0xf6d2a2, near: 45, far: 200, sun: 2.0, hemi: 0.95 },
  volcano: { fog: 0x4c3034, top: 0x2a161e, bot: 0x92402e, near: 18, far: 115, sun: 0.9, hemi: 0.55 },
  sea: { fog: 0x9ab8d0, top: 0x5a8ac0, bot: 0xb0c8dc, near: 60, far: 220, sun: 1.5, hemi: 0.9 },
};


function buildTerrain(sc) {
  const nx = W.nx, nz = W.nz, n = nx;
  const geo = new THREE.PlaneGeometry((nx - 1) * W.cell, (nz - 1) * W.cell, nx - 1, nz - 1);
  geo.rotateX(-Math.PI / 2);
  geo.translate(W.x0 + (nx - 1) * W.cell / 2, 0, W.z0 + (nz - 1) * W.cell / 2);
  const pos = geo.attributes.position;
  const colors = new Float32Array(pos.count * 3);
  const r = rng(99);
  const c = new THREE.Color(), tmp = new THREE.Color();
  // blurred biome colors
  const base = new Float32Array(nx * nz * 3);
  for (let j = 0; j < nz; j++) {
    for (let i = 0; i < nx; i++) {
      let R = 0, G = 0, B = 0, cnt = 0;
      for (let dj = -3; dj <= 3; dj += 3) for (let di = -3; di <= 3; di += 3) {
        const ii = Math.max(0, Math.min(nx - 1, i + di)), jj = Math.max(0, Math.min(nz - 1, j + dj));
        const col = BIOME_COLOR[W.bm[jj * nx + ii]];
        R += col[0]; G += col[1]; B += col[2]; cnt++;
      }
      const k = (j * nx + i) * 3;
      base[k] = R / cnt; base[k + 1] = G / cnt; base[k + 2] = B / cnt;
    }
  }
  const capital = sc.towns.find((t) => t.kind === 'capital');
  for (let j = 0; j < nz; j++) {
    for (let i = 0; i < nx; i++) {
      const idx = j * nx + i;
      const x = W.x0 + i * W.cell, z = W.z0 + j * W.cell;
      const h = W.hm[idx];
      pos.setY(idx, h);
      const k = idx * 3;
      const nz = Math.sin(x * 0.31) * Math.cos(z * 0.27) * 0.5 + r() * 0.5;
      c.setRGB(base[k] + nz * 0.04, base[k + 1] + nz * 0.05, base[k + 2] + nz * 0.025);
      const b = W.bm[idx];
      if (b === 4) { // snow: bare rock on steep slopes and high peaks stay white
        const sl = Math.abs(heightAt(x + 2, z) - h) + Math.abs(heightAt(x, z + 2) - h);
        c.lerp(tmp.setRGB(0.45, 0.45, 0.48), smoothstep(2.6, 5, sl) * 0.85);
      }
      if (b === 3 && h < 0.3) c.lerp(tmp.setRGB(0.12, 0.2, 0.12), 0.6); // swampy mud
      if (b === 7 && nz > 0.62) c.lerp(tmp.setRGB(0.98, 0.76, 0.84), 0.55); // drifts of petals
      if (b === 8) { const band = 0.88 + 0.12 * Math.sin(h * 2.6 + nz * 2); c.multiplyScalar(band); if (h > 9) c.lerp(tmp.setRGB(0.86, 0.5, 0.28), 0.25); } // rock strata
      if (b === 9) { const ember = Math.max(0, Math.sin(x * 0.21) * Math.cos(z * 0.17)); if (ember > 0.8) c.lerp(tmp.setRGB(0.9, 0.3, 0.05), (ember - 0.8) * 2.2); } // glowing cracks
      if (b === 6) c.offsetHSL(0, 0, Math.sin(x * 0.13 + z * 0.1) * 0.02);
      if (h < 0.6 && b !== 0 && b < 9) c.lerp(tmp.setRGB(0.76, 0.7, 0.5), smoothstep(0.6, 0.0, h) * 0.8); // beach
      if (h < -0.4 && b !== 10) c.lerp(tmp.setRGB(0.35, 0.34, 0.28), 0.7); // river / sea bed
      if (b === 10) c.setRGB(1.0, 0.3 + nz * 0.25, 0.03); // lava
      // roads
      let d = 1e9;
      for (const rd of sc.roads) d = Math.min(d, distToPoly(x, z, rd));
      const dirt = 1 - smoothstep(1.1, 2.6, d);
      if (dirt > 0) c.lerp(tmp.setRGB(0.5 + nz * 0.05, 0.4, 0.26), dirt * 0.9);
      // paving inside towns
      for (const t of sc.towns) {
        const dd = Math.hypot(x - t.x, z - t.z);
        if (t.kind === 'capital') {
          if (Math.abs(x - t.x) < 52 && Math.abs(z - t.z) < 47) {
            if (t.style === 'jade') c.lerp(tmp.setRGB(0.62, 0.5, 0.45), 0.7);
            else c.lerp(tmp.setRGB(0.62, 0.6, 0.55), 0.75);
          }
        } else if (t.kind === 'village' && dd < 15) c.lerp(tmp.setRGB(0.5, 0.42, 0.3), (1 - smoothstep(9, 15, dd)) * 0.8);
        else if (t.kind === 'site' && dd < 12) c.lerp(tmp.setRGB(0.4, 0.4, 0.42), 0.8);
      }
      colors[k] = c.r; colors[k + 1] = c.g; colors[k + 2] = c.b;
    }
  }
  geo.setAttribute('color', new THREE.BufferAttribute(colors, 3));
  geo.computeVertexNormals();
  const m = new THREE.Mesh(geo, new THREE.MeshLambertMaterial({ vertexColors: true, map: noiseTexture() }));
  m.receiveShadow = true;
  return m;
}

function buildSky() {
  const mat = new THREE.ShaderMaterial({
    side: THREE.BackSide, depthWrite: false, fog: false,
    uniforms: { top: { value: new THREE.Color(0x4a86d6) }, bot: { value: new THREE.Color(0xbcd6e6) } },
    vertexShader: 'varying vec3 p; void main(){ p = position; gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.0); }',
    fragmentShader: 'varying vec3 p; uniform vec3 top; uniform vec3 bot; void main(){ float h = clamp(normalize(p).y*1.6, 0.0, 1.0); gl_FragColor = vec4(mix(bot, top, pow(h, 0.7)), 1.0); }',
  });
  const m = new THREE.Mesh(new THREE.SphereGeometry(450, 24, 16), mat);
  m.frustumCulled = false;
  return m;
}

function buildBridges(scene, bridges) {
  const wood = new THREE.MeshLambertMaterial({ color: 0x7a5a32 });
  for (const b of bridges) {
    const g = new THREE.Group();
    const deck = new THREE.Mesh(new THREE.BoxGeometry(4, 0.3, b.len), wood);
    deck.position.y = 0.55; deck.castShadow = deck.receiveShadow = true; g.add(deck);
    for (const sx of [-1, 1]) {
      const rail = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.9, b.len), wood); rail.position.set(sx * 1.9, 1.1, 0); g.add(rail);
      for (const sz of [-1, 1]) { const post = new THREE.Mesh(new THREE.BoxGeometry(0.3, 1.8, 0.3), wood); post.position.set(sx * 1.9, 0.9, sz * b.len / 2.2); g.add(post); }
    }
    g.position.set(b.x, Math.max(0.1, heightAt(b.x, b.z)), b.z);
    g.rotation.y = b.rot;
    scene.add(g);
  }
}

export function buildWorld(scenery, renderer) {
  setWorldData(scenery);
  const scene = new THREE.Scene();
  scene.fog = new THREE.Fog(0xbcd6e6, 45, 170);
  const sky = buildSky();
  scene.add(sky);
  const hemi = new THREE.HemisphereLight(0xcfe6ff, 0x4a5a36, 0.85);
  scene.add(hemi);
  const sun = new THREE.DirectionalLight(0xfff0d0, 1.5);
  sun.position.set(60, 90, 40);
  sun.castShadow = true;
  sun.shadow.mapSize.set(2048, 2048);
  const sc = sun.shadow.camera;
  sc.left = -45; sc.right = 45; sc.top = 45; sc.bottom = -45; sc.near = 10; sc.far = 260;
  sun.shadow.bias = -0.0006;
  sun.shadow.normalBias = 0.04;
  scene.add(sun, sun.target);
  renderer.shadowMap.enabled = true;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;

  scene.add(buildTerrain(scenery));
  const water = new THREE.Mesh(new THREE.PlaneGeometry(3400, 2400), new THREE.MeshPhongMaterial({ color: 0x2f7fb0, transparent: true, opacity: 0.78, shininess: 80, specular: 0x88aacc }));
  water.rotation.x = -Math.PI / 2; water.position.set(420, -0.12, 0);
  scene.add(water);
  const treeGrid = buildVegetation(scene, scenery);
  const solids = buildTowns(scene, scenery.buildings);
  buildCabin(scenery.house, scene);
  buildBridges(scene, scenery.bridges);
  const walls = scenery.house;
  const R = 0.4;

  // look transition
  const look = { ...ZONE_LOOK.forest };
  let target = ZONE_LOOK.forest;
  const cFog = new THREE.Color(look.fog), cTop = new THREE.Color(look.top), cBot = new THREE.Color(look.bot);
  const tFog = new THREE.Color(), tTop = new THREE.Color(), tBot = new THREE.Color();

  function setZone(key) { target = ZONE_LOOK[key] || ZONE_LOOK.forest; }
  function updateLook(dt) {
    const k = 1 - Math.exp(-dt * 1.2);
    cFog.lerp(tFog.set(target.fog), k); cTop.lerp(tTop.set(target.top), k); cBot.lerp(tBot.set(target.bot), k);
    look.near += (target.near - look.near) * k; look.far += (target.far - look.far) * k;
    look.sun += (target.sun - look.sun) * k; look.hemi += (target.hemi - look.hemi) * k;
    scene.fog.color.copy(cFog); scene.fog.near = look.near; scene.fog.far = look.far;
    sky.material.uniforms.top.value.copy(cTop); sky.material.uniforms.bot.value.copy(cBot);
    sun.intensity = look.sun; hemi.intensity = look.hemi;
  }

  const docks = scenery.buildings.filter((b) => b.k === 'dock');
  const gh = (x, z) => (onDock(x, z) ? 1.2 : heightAt(x, z)); // the pier is a walkway above the water
  const onDock = (x, z) => docks.some((b) => Math.abs(x - b.x) <= b.w / 2 && Math.abs(z - b.z) <= b.d / 2);
  function collide(x, z, ox, oz) {
    // trees and rocks
    const gx = Math.floor(x / 8), gz = Math.floor(z / 8);
    for (let ix = gx - 1; ix <= gx + 1; ix++) {
      for (let iz = gz - 1; iz <= gz + 1; iz++) {
        const cell = treeGrid[ix + ',' + iz];
        if (!cell) continue;
        for (const t of cell) {
          const dx = x - t.x, dz = z - t.z, d = Math.hypot(dx, dz), min = t.r + R;
          if (d < min && d > 1e-4) { x = t.x + (dx / d) * min; z = t.z + (dz / d) * min; }
        }
      }
    }
    // cabin walls and town buildings (axis aligned boxes)
    const boxes = [];
    for (const b of walls) boxes.push(b);
    for (const b of solids) {
      if (Math.abs(x - b.x) < b.w / 2 + 3 && Math.abs(z - b.z) < b.d / 2 + 3) boxes.push(b);
    }
    for (const b of boxes) {
      const hx = b.w / 2 + R, hz = b.d / 2 + R;
      const dx = x - b.x, dz = z - b.z;
      if (Math.abs(dx) < hx && Math.abs(dz) < hz) {
        const px = hx - Math.abs(dx), pz = hz - Math.abs(dz);
        if (px < pz) x = b.x + Math.sign(dx || 1) * hx; else z = b.z + Math.sign(dz || 1) * hz;
      }
    }
    const xMax = W.x0 + (W.nx - 1) * W.cell - 6, zMax = W.z0 + (W.nz - 1) * W.cell - 6;
    x = Math.max(W.x0 + 6, Math.min(xMax, x)); z = Math.max(W.z0 + 6, Math.min(zMax, z));
    // water too deep to wade through, and slopes too steep to climb
    if (ox !== undefined) {
      const h = gh(x, z);
      if (h < -1.5) return [ox, oz];
      const step = Math.hypot(x - ox, z - oz);
      if (step > 1e-4 && (h - gh(ox, oz)) / step > 1.35) {
        // try sliding along each axis
        if ((gh(x, oz) - gh(ox, oz)) / Math.max(Math.abs(x - ox), 1e-4) <= 1.35 && gh(x, oz) > -1.5) return [x, oz];
        if ((gh(ox, z) - gh(ox, oz)) / Math.max(Math.abs(z - oz), 1e-4) <= 1.35 && gh(ox, z) > -1.5) return [ox, z];
        return [ox, oz];
      }
    }
    return [x, z];
  }

  let dungeon = null;
  return {
    scene, collide, setZone, onDock,
    followSun(x, y, z) {
      sky.position.set(x, 0, z);
      sun.target.position.set(x, y, z);
      sun.position.set(x + 60, y + 90, z + 40);
    },
    update(dt) { updateLook(dt); for (const sh of ships) sh.update(dt, {}); water.position.y = -0.12 + Math.sin(performance.now() / 1400) * 0.04; },
    enterDungeon(info) { dungeon = buildDungeon(info); return dungeon; },
    leaveDungeon() { dungeon = null; },
    get dungeon() { return dungeon; },
    biomeKey: (x, z) => BIOME_KEYS[biomeAt(x, z)],
  };
}
