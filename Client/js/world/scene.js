import * as THREE from 'three';
import { W, setWorldData, heightAt, SWIM_FLOOR, biomeAt, distToPoly, smoothstep, BIOME_KEYS } from '../util.js';
import { buildVegetation } from './vegetation.js';
import { buildTowns, buildCabin, ships } from './buildings.js';
import { buildLamps } from './lamps.js';
import { buildDungeon } from './dungeon.js';
import { sunInfo } from '../daynight.js';

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
// Per biome look. top/bot/fog are the daytime sky and fog colours; `sky` describes the sky box of the biome:
// cloud cover/colour/speed, star brightness, aurora, horizon haze/glow and the colour of the twilight band.
const SKY = (cloud, cloudCol, speed, stars, haze, tw, extra = {}) => ({ cloud, cloudCol, speed, stars, haze, tw, aurora: 0, glow: 0, ...extra });
export const ZONE_LOOK = {
  forest: { fog: 0xbcd6e6, top: 0x4a86d6, bot: 0xbcd6e6, near: 45, far: 170, sun: 1.5, hemi: 0.85, sky: SKY(0.48, 0xffffff, 0.004, 1, 0.15, 0xff9a62) },
  plains: { fog: 0xcfe3f0, top: 0x5a9ae0, bot: 0xd8e8f2, near: 70, far: 230, sun: 1.7, hemi: 0.95, sky: SKY(0.58, 0xffffff, 0.006, 1, 0.1, 0xff8a50) },
  swamp: { fog: 0x8a9a7a, top: 0x6a7a68, bot: 0x9aa890, near: 22, far: 110, sun: 0.8, hemi: 0.6, sky: SKY(0.85, 0x7d8a74, 0.002, 0.25, 0.55, 0x9aa05a) },
  snow: { fog: 0xdce6f0, top: 0x8aa8cc, bot: 0xe6eef6, near: 35, far: 150, sun: 1.3, hemi: 1.0, sky: SKY(0.38, 0xe8f0ff, 0.003, 1.4, 0.2, 0xf08aa8, { aurora: 1 }) },
  desert: { fog: 0xf0dcb0, top: 0x6aa0d8, bot: 0xf2e4c0, near: 55, far: 210, sun: 1.9, hemi: 0.95, sky: SKY(0.12, 0xfff0d0, 0.003, 1.6, 0.4, 0xff6a30) },
  jungle: { fog: 0x8fd8a8, top: 0x1f9a8a, bot: 0xa8e8c0, near: 22, far: 125, sun: 1.2, hemi: 0.8, sky: SKY(0.72, 0xd8f4e0, 0.003, 0.6, 0.45, 0xffb040) },
  blossom: { fog: 0xf4d6e4, top: 0xb98ad0, bot: 0xffdbe6, near: 50, far: 200, sun: 1.6, hemi: 1.05, sky: SKY(0.5, 0xffe8f2, 0.005, 1.1, 0.25, 0xff7aa8) },
  mesa: { fog: 0xe8b88c, top: 0xd8804e, bot: 0xf6d2a2, near: 45, far: 200, sun: 2.0, hemi: 0.95, sky: SKY(0.26, 0xffd8b0, 0.004, 1.3, 0.35, 0xff5a20) },
  volcano: { fog: 0x4c3034, top: 0x2a161e, bot: 0x92402e, near: 18, far: 115, sun: 0.9, hemi: 0.55, sky: SKY(0.92, 0x4a3034, 0.006, 0.12, 0.6, 0xff4a20, { glow: 1 }) },
  sea: { fog: 0x9ab8d0, top: 0x5a8ac0, bot: 0xb0c8dc, near: 60, far: 220, sun: 1.5, hemi: 0.9, sky: SKY(0.5, 0xffffff, 0.005, 1.1, 0.25, 0xff8a60) },
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
          if (Math.abs(x - t.x) < (t.style === 'jade' ? 52 : 78) && Math.abs(z - t.z) < (t.style === 'jade' ? 47 : 71)) {
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

export function buildSky() {
  const mat = new THREE.ShaderMaterial({
    side: THREE.BackSide, depthWrite: false, fog: false,
    uniforms: {
      top: { value: new THREE.Color(0x4a86d6) }, bot: { value: new THREE.Color(0xbcd6e6) },
      sunDir: { value: new THREE.Vector3(0, 1, 0) }, moonDir: { value: new THREE.Vector3(0, -1, 0) },
      sunCol: { value: new THREE.Color(0xfff0d0) }, cloudCol: { value: new THREE.Color(0xffffff) }, glowCol: { value: new THREE.Color(0xff4a20) },
      time: { value: 0 }, night: { value: 0 }, sunUp: { value: 1 }, cloud: { value: 0.5 }, speed: { value: 0.004 },
      stars: { value: 1 }, aurora: { value: 0 }, haze: { value: 0.2 }, glow: { value: 0 },
    },
    vertexShader: 'varying vec3 p; void main(){ p = position; gl_Position = projectionMatrix * modelViewMatrix * vec4(position,1.0); }',
    fragmentShader: `
      varying vec3 p;
      uniform vec3 top, bot, sunDir, moonDir, sunCol, cloudCol, glowCol;
      uniform float time, night, sunUp, cloud, speed, stars, aurora, haze, glow;
      float h21(vec2 q){ q = fract(q * vec2(123.34, 456.21)); q += dot(q, q + 45.32); return fract(q.x * q.y); }
      float h31(vec3 q){ q = fract(q * 0.3183099 + 0.1); q *= 17.0; return fract(q.x * q.y * q.z * (q.x + q.y + q.z)); }
      float vn(vec2 x){ vec2 i = floor(x), f = fract(x); f = f * f * (3.0 - 2.0 * f);
        return mix(mix(h21(i), h21(i + vec2(1,0)), f.x), mix(h21(i + vec2(0,1)), h21(i + vec2(1,1)), f.x), f.y); }
      float fbm(vec2 x){ float v = 0.0, a = 0.5; for (int i = 0; i < 5; i++) { v += a * vn(x); x = x * 2.03 + 17.0; a *= 0.5; } return v; }
      void main(){
        vec3 d = normalize(p);
        float hh = clamp(d.y * 1.6, 0.0, 1.0);
        vec3 col = mix(bot, top, pow(hh, 0.7));
        // clouds drifting over the biome
        float cm = 0.0; vec3 cc = cloudCol;
        if (d.y > 0.0) {
          vec2 uv = d.xz / (d.y + 0.22) * 0.85 + vec2(time * speed, time * speed * 0.35);
          float c = fbm(uv * 1.25);
          cm = smoothstep(1.0 - cloud, 1.0 - cloud + 0.3, c) * smoothstep(0.0, 0.22, d.y);
          float lit = 0.45 + 0.55 * smoothstep(0.35, 0.75, c);
          cc = cloudCol * mix(vec3(0.16, 0.2, 0.32), vec3(1.0), (1.0 - night)) * (0.7 + 0.3 * lit);
          cc += sunCol * pow(max(dot(d, sunDir), 0.0), 6.0) * 0.35 * sunUp;
          cc += glowCol * glow * 0.35 * (1.0 - hh);
        }
        // sun and moon
        float sd = max(dot(d, sunDir), 0.0), md = max(dot(d, moonDir), 0.0);
        col += sunCol * (pow(sd, 5.0) * 0.12 + pow(sd, 40.0) * 0.4) * sunUp * (1.0 - 0.5 * cm);
        col += sunCol * smoothstep(0.99935, 0.99965, sd) * 2.2 * sunUp * (1.0 - cm);
        float moonDisc = smoothstep(0.99905, 0.99925, md);
        col += (vec3(0.82, 0.88, 1.0) * (0.85 + 0.15 * vn(d.xz * 160.0)) * moonDisc * (1.0 - cm) + vec3(0.35, 0.45, 0.7) * pow(md, 60.0) * 0.5) * night;
        // stars
        if (night > 0.02 && d.y > 0.0) {
          vec3 q = d * 150.0; float s = h31(floor(q)); vec3 f = fract(q) - 0.5;
          float st = step(0.9955, s) * smoothstep(0.42, 0.0, length(f)) * (0.6 + 0.4 * sin(time * 2.0 + s * 80.0));
          vec3 q2 = d * 60.0; float s2 = h31(floor(q2) + 7.0); vec3 f2 = fract(q2) - 0.5;
          st += step(0.9975, s2) * smoothstep(0.4, 0.0, length(f2)) * 1.3;
          col += vec3(0.85, 0.9, 1.0) * st * stars * night * smoothstep(0.0, 0.18, d.y) * (1.0 - cm);
        }
        // aurora (snow lands)
        if (aurora > 0.01 && night > 0.05 && d.y > 0.0) {
          float band = sin(d.x * 5.0 + fbm(d.xz * 3.0 + time * 0.04) * 7.0) * 0.5 + 0.5;
          float yy = smoothstep(0.12, 0.38, d.y) * (1.0 - smoothstep(0.5, 0.9, d.y));
          col += mix(vec3(0.1, 0.95, 0.5), vec3(0.6, 0.3, 0.95), d.x * 0.5 + 0.5) * band * band * yy * aurora * night * 0.6;
        }
        col = mix(col, cc, cm * 0.92);
        // horizon haze and volcanic glow
        col = mix(col, bot, haze * pow(1.0 - hh, 2.0));
        col += glowCol * glow * pow(1.0 - hh, 3.0) * (0.35 + 0.65 * night) * 0.8;
        gl_FragColor = vec4(col, 1.0);
      }`,
  });
  const m = new THREE.Mesh(new THREE.SphereGeometry(450, 32, 20), mat);
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
  water.rotation.x = -Math.PI / 2; water.position.set(W.x0 + (W.nx - 1) * W.cell / 2, -0.12, 0);
  scene.add(water);
  const treeGrid = buildVegetation(scene, scenery);
  const solids = buildTowns(scene, scenery.buildings);
  buildCabin(scenery.house, scene);
  buildBridges(scene, scenery.bridges);
  const lamps = buildLamps(scene, scenery);
  const walls = scenery.house;
  const R = 0.4;

  // look: biome colours/sky blended over time, then graded by the time of day (sun, moon, twilight, night)
  const look = { ...ZONE_LOOK.forest };
  const skyP = { ...ZONE_LOOK.forest.sky };
  let target = ZONE_LOOK.forest;
  const cFog = new THREE.Color(look.fog), cTop = new THREE.Color(look.top), cBot = new THREE.Color(look.bot);
  const cCloud = new THREE.Color(skyP.cloudCol), cTw = new THREE.Color(skyP.tw);
  const tFog = new THREE.Color(), tTop = new THREE.Color(), tBot = new THREE.Color();
  const gTop = new THREE.Color(), gBot = new THREE.Color(), gFog = new THREE.Color(), tmp = new THREE.Color();
  const NIGHT_TOP = new THREE.Color(0x040a1c), NIGHT_BOT = new THREE.Color(0x101c3e), NIGHT_TINT = new THREE.Color(0.1, 0.14, 0.3);
  const SUN_DAY = new THREE.Color(0xfff0d8), SUN_LOW = new THREE.Color(0xff9a4a), MOON = new THREE.Color(0x8aa4ff);
  const HEMI_SKY_DAY = new THREE.Color(0xcfe6ff), HEMI_SKY_NIGHT = new THREE.Color(0x2c3c7a);
  const lightDir = new THREE.Vector3(0.5, 0.8, 0.3);
  let clockT = 0;

  function setZone(key) { target = ZONE_LOOK[key] || ZONE_LOOK.forest; }
  function updateLook(dt) {
    const k = 1 - Math.exp(-dt * 1.2);
    cFog.lerp(tFog.set(target.fog), k); cTop.lerp(tTop.set(target.top), k); cBot.lerp(tBot.set(target.bot), k);
    look.near += (target.near - look.near) * k; look.far += (target.far - look.far) * k;
    look.sun += (target.sun - look.sun) * k; look.hemi += (target.hemi - look.hemi) * k;
    const ts = target.sky;
    for (const n of ['cloud', 'speed', 'stars', 'haze', 'aurora', 'glow']) skyP[n] += (ts[n] - skyP[n]) * k;
    cCloud.lerp(tmp.set(ts.cloudCol), k); cTw.lerp(tmp.set(ts.tw), k);

    const t = sunInfo();
    clockT += dt;
    // sky: night colours, twilight band at the horizon
    gTop.copy(cTop).multiply(NIGHT_TINT).lerp(NIGHT_TOP, 0.55).lerp(cTop, t.day);
    gBot.copy(cBot).multiply(NIGHT_TINT).lerp(NIGHT_BOT, 0.6).lerp(cBot, t.day);
    gBot.lerp(cTw, t.twilight * 0.65);
    gTop.lerp(tmp.copy(cTw).multiplyScalar(0.55), t.twilight * 0.22);
    gFog.copy(gBot).multiplyScalar(0.55 + 0.45 * t.day);
    scene.fog.color.copy(gFog); scene.fog.near = look.near; scene.fog.far = look.far * (0.75 + 0.25 * t.day);
    const u = sky.material.uniforms;
    u.top.value.copy(gTop); u.bot.value.copy(gBot);
    u.sunDir.value.set(...t.sun); u.moonDir.value.set(...t.moon);
    u.sunCol.value.copy(SUN_LOW).lerp(SUN_DAY, Math.min(1, Math.max(0, t.elev * 2.2)));
    u.cloudCol.value.copy(cCloud); u.glowCol.value.copy(cTw);
    u.time.value = clockT; u.night.value = t.night; u.sunUp.value = Math.max(0, Math.min(1, (t.elev + 0.1) * 4));
    u.cloud.value = skyP.cloud; u.speed.value = skyP.speed; u.stars.value = skyP.stars; u.aurora.value = skyP.aurora;
    u.haze.value = skyP.haze; u.glow.value = skyP.glow;
    // lights: the sun by day, the moon by night (never lower than a minimum angle so shadows stay inside the shadow map)
    const dayLight = t.day, moonLight = t.night * 0.5;
    const src = t.elev > -0.04 ? t.sun : t.moon;
    lightDir.set(src[0], Math.max(0.3, src[1]), src[2]).normalize();
    sun.color.copy(MOON).lerp(SUN_LOW, Math.min(1, dayLight * 1.6)).lerp(SUN_DAY, Math.min(1, Math.max(0, (t.elev - 0.1) * 2)));
    sun.intensity = look.sun * dayLight * (0.35 + 0.65 * Math.min(1, t.elev * 3 + 0.2)) + moonLight * 0.9;
    hemi.intensity = look.hemi * (0.4 + 0.6 * dayLight);
    hemi.color.copy(HEMI_SKY_NIGHT).lerp(HEMI_SKY_DAY, dayLight);
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
    // the open sea is too deep to swim in, and slopes too steep to climb
    if (ox !== undefined) {
      const h = gh(x, z);
      if (h < SWIM_FLOOR) return [ox, oz]; // the open sea is closed: only the ferry crosses it
      const step = Math.hypot(x - ox, z - oz);
      if (step > 1e-4 && (h - gh(ox, oz)) / step > 1.35) {
        // try sliding along each axis
        if ((gh(x, oz) - gh(ox, oz)) / Math.max(Math.abs(x - ox), 1e-4) <= 1.35 && gh(x, oz) > SWIM_FLOOR) return [x, oz];
        if ((gh(ox, z) - gh(ox, oz)) / Math.max(Math.abs(z - oz), 1e-4) <= 1.35 && gh(ox, z) > SWIM_FLOOR) return [ox, z];
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
      sun.position.set(x + lightDir.x * 110, y + lightDir.y * 110, z + lightDir.z * 110);
    },
    update(dt, cam) { updateLook(dt); lamps.update(dt, cam || sky.position, sunInfo().day); for (const sh of ships) sh.update(dt, {}); water.position.y = -0.12 + Math.sin(performance.now() / 1400) * 0.04; },
    enterDungeon(info) { dungeon = buildDungeon(info); return dungeon; },
    leaveDungeon() { dungeon = null; },
    get dungeon() { return dungeon; },
    biomeKey: (x, z) => BIOME_KEYS[biomeAt(x, z)],
  };
}
