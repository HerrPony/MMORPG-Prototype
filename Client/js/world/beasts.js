import * as THREE from 'three';
import { pixTex, shade } from './human.js';

// Pixel-art four-legged animals (wolves, bears, boars, bulls ...): box body, head, snout, ears, tail and four
// swinging legs, all painted with small generated textures.

export const loadBeastAssets = () => Promise.resolve(true);
export const beastsLoaded = () => true;

// L body length, W width, H body height, G leg length, hs head size, sn snout length (metres, before scale s)
const KINDS = {
  wolf: { L: 1.0, W: 0.36, H: 0.4, G: 0.42, hs: 0.3, sn: 0.2, color: 0x6c7078, belly: 0xa8aaae, ears: 'point', tail: 0.4, eye: 0xffd040, s: 1 },
  wolf_alpha: { L: 1.0, W: 0.36, H: 0.4, G: 0.42, hs: 0.3, sn: 0.2, color: 0x2c2e36, belly: 0x56565f, ears: 'point', tail: 0.45, eye: 0xff3a20, s: 1.4 },
  icewolf: { L: 1.0, W: 0.36, H: 0.4, G: 0.42, hs: 0.3, sn: 0.2, color: 0xcfe2f4, belly: 0xf2f8ff, ears: 'point', tail: 0.4, eye: 0x40a0ff, s: 1.1 },
  hyena: { L: 1.0, W: 0.34, H: 0.4, G: 0.44, hs: 0.3, sn: 0.2, color: 0xa88a58, belly: 0xd8c090, pat: 'spots', mark: 0x4a3a20, ears: 'round', tail: 0.3, eye: 0xffd040, s: 1 },
  boar: { L: 0.95, W: 0.5, H: 0.5, G: 0.3, hs: 0.34, sn: 0.22, color: 0x5a4030, belly: 0x7a6048, pat: 'bristle', mark: 0x2a1c14, ears: 'round', tusks: true, tail: 0.15, eye: 0x201010, s: 1 },
  bear: { L: 1.2, W: 0.7, H: 0.7, G: 0.45, hs: 0.42, sn: 0.2, color: 0x5a3c24, belly: 0x7a5636, ears: 'round', tail: 0.12, eye: 0x201010, s: 1 },
  bear_big: { L: 1.2, W: 0.7, H: 0.7, G: 0.45, hs: 0.42, sn: 0.2, color: 0x3e2a1a, belly: 0x5a4028, ears: 'round', tail: 0.12, eye: 0xff5020, s: 1.35 },
  polarbear: { L: 1.2, W: 0.7, H: 0.7, G: 0.45, hs: 0.42, sn: 0.24, color: 0xf0f0ec, belly: 0xffffff, ears: 'round', tail: 0.12, eye: 0x202020, s: 1 },
  polarbear_big: { L: 1.2, W: 0.7, H: 0.7, G: 0.45, hs: 0.42, sn: 0.24, color: 0xe2eaf2, belly: 0xf6fbff, ears: 'round', tail: 0.12, eye: 0xff4020, s: 1.35 },
  bull: { L: 1.3, W: 0.6, H: 0.65, G: 0.55, hs: 0.36, sn: 0.22, color: 0x7a4a2a, belly: 0x8a5a38, pat: 'patch', mark: 0xf0e8d8, ears: 'point', horns: true, tail: 0.5, eye: 0x201010, s: 1 },
  sheep: { L: 0.8, W: 0.5, H: 0.5, G: 0.35, hs: 0.26, sn: 0.12, color: 0xf2f2ec, belly: 0xe4e4dc, pat: 'wool', head: 0x3a3a3a, ears: 'round', tail: 0.1, eye: 0xe0e0e0, s: 1 },
};

function bodyTex(K, c, key) {
  return pixTex(`bb${key}${c}`, 16, 8, (p) => {
    p.fill(c, 0.14);
    if (K.pat === 'wool') { for (let i = 0; i < 40; i++) p.rect((i * 7) % 15, (i * 5) % 7, 2, 2, shade(c, i % 2 ? 0.93 : 1.04)); return; }
    p.rect(0, 0, 16, 2, shade(c, 0.82)); // darker back
    p.fill(K.belly, 0.1, 0, 6, 16, 2);   // lighter belly
    if (K.pat === 'spots') for (let i = 0; i < 14; i++) p.rect((i * 5 + 2) % 15, 1 + (i * 3) % 5, 1 + (i % 2), 1 + (i % 2), K.mark);
    if (K.pat === 'bristle') { p.rect(0, 0, 16, 1, K.mark); for (let x = 0; x < 16; x += 2) p.px(x, 1, K.mark); }
    if (K.pat === 'patch') { p.rect(2, 2, 4, 3, K.mark); p.rect(9, 1, 5, 3, K.mark); p.rect(7, 4, 2, 2, K.mark); }
  });
}

function faceTex(K, c, key) {
  return pixTex(`bf${key}${c}`, 8, 8, (p) => {
    p.fill(c, 0.12);
    p.rect(0, 0, 8, 1, shade(c, 0.85));
    p.rect(1, 3, 2, 1, K.eye); p.rect(5, 3, 2, 1, K.eye);
    p.px(1, 3, 0x111111); p.px(6, 3, 0x111111); // pupils
  });
}

export function makePixelBeast(kind = 'wolf', seed = 0) {
  const K = KINDS[kind] || KINDS.wolf;
  const S = K.s;
  const col = shade(K.color, 1 + ((seed % 3) - 1) * 0.07);
  const root = new THREE.Group();
  const inner = new THREE.Group(); root.add(inner);
  const body = new THREE.Group(); inner.add(body);
  const L = K.L * S, W = K.W * S, H = K.H * S, G = K.G * S, hs = K.hs * S, sn = K.sn * S;
  const mat = (t) => new THREE.MeshLambertMaterial({ map: t });
  const box = (w, h, d, m) => { const x = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), m); x.castShadow = true; return x; };
  const flat = (c) => mat(pixTex('bs' + c, 4, 4, (p) => p.fill(c, 0.16)));
  const headCol = K.head ?? col;

  const torso = box(W, H, L, mat(bodyTex(K, col, kind))); torso.position.set(0, G + H / 2, 0); body.add(torso);
  // head, snout, ears
  const head = new THREE.Group(); head.position.set(0, G + H * 0.85, L / 2 + hs * 0.2); body.add(head);
  const hm = mat(pixTex('bh' + headCol, 4, 4, (p) => p.fill(headCol, 0.14)));
  const face = mat(faceTex(K, headCol, kind));
  const skull = box(hs, hs, hs, [hm, hm, hm, hm, face, hm]); skull.position.set(0, 0, hs * 0.3); head.add(skull);
  const snout = box(hs * 0.62, hs * 0.5, sn, flat(shade(headCol, 1.1))); snout.position.set(0, -hs * 0.18, hs * 0.8 + sn / 2); head.add(snout);
  const nose = box(hs * 0.22, hs * 0.16, 0.04, flat(0x111111)); nose.position.set(0, -hs * 0.06, hs * 0.8 + sn); head.add(nose);
  for (const s of [-1, 1]) {
    if (K.ears === 'point') { const e = box(hs * 0.22, hs * 0.38, hs * 0.1, flat(shade(headCol, 0.9))); e.position.set(s * hs * 0.34, hs * 0.5, hs * 0.2); head.add(e); }
    else if (K.ears === 'round') { const e = box(hs * 0.24, hs * 0.24, hs * 0.14, flat(shade(headCol, 0.9))); e.position.set(s * hs * 0.4, hs * 0.42, hs * 0.2); head.add(e); }
    if (K.tusks) { const t = box(0.05, 0.12, 0.05, flat(0xf0e6c8)); t.position.set(s * hs * 0.24, -hs * 0.36, hs * 0.8 + sn * 0.8); head.add(t); }
    if (K.horns) { const h1 = box(0.07, 0.07, 0.2, flat(0xe8e0c8)); h1.position.set(s * hs * 0.55, hs * 0.38, hs * 0.2); head.add(h1); const h2 = box(0.07, 0.2, 0.07, flat(0xe8e0c8)); h2.position.set(s * hs * 0.7, hs * 0.5, hs * 0.25); head.add(h2); }
  }
  // tail
  const tail = new THREE.Group(); tail.position.set(0, G + H * 0.85, -L / 2); body.add(tail);
  const tl = K.tail * S;
  const tb = box(W * 0.25, W * 0.25, tl, flat(shade(col, 0.9))); tb.position.set(0, -tl * 0.2, -tl / 2); tail.add(tb); tail.rotation.x = -0.5;
  // legs
  const legs = [];
  const lm = mat(pixTex('bl' + col, 4, 8, (p) => { p.fill(col, 0.14); p.fill(shade(col, 0.6), 0.14, 0, 6, 4, 2); }));
  for (const [x, z] of [[-1, 1], [1, 1], [-1, -1], [1, -1]]) {
    const pv = new THREE.Group(); pv.position.set(x * (W / 2 - W * 0.14), G, z * (L / 2 - W * 0.18)); body.add(pv);
    const lw = Math.max(0.1, W * 0.3);
    const leg = box(lw, G, lw, lm); leg.position.y = -G / 2; pv.add(leg);
    legs.push(pv);
  }

  let attackT = 0, dead = 0, phase = 0, sw = 0, t = Math.random() * 10;
  return {
    obj: root, height: G + H + hs * 0.6, kind: 'wolf', parts: { head, body },
    attack() { attackT = 0.01; },
    update(dt, st) {
      t += dt;
      const moving = st.speed > 0.08;
      if (moving) phase += dt * Math.max(5, st.speed * 13);
      sw += ((moving ? Math.min(1, 0.45 + st.speed * 0.6) : 0) - sw) * Math.min(1, dt * 12);
      const a = Math.sin(phase) * sw * 0.8;
      legs[0].rotation.x = a; legs[3].rotation.x = a; legs[1].rotation.x = -a; legs[2].rotation.x = -a;
      tail.rotation.y = Math.sin(t * 6 + phase) * 0.3;
      head.rotation.x = Math.sin(t * 1.5) * 0.04;
      let lunge = 0;
      if (attackT > 0) {
        attackT += dt;
        const k = Math.min(1, attackT / 0.4);
        lunge = Math.sin(k * Math.PI) * 0.4;
        head.rotation.x = Math.sin(k * Math.PI) * 0.5;
        if (k >= 1) attackT = 0;
      }
      body.position.z = lunge;
      body.position.y = moving ? Math.abs(Math.sin(phase)) * 0.04 : 0;
      dead += ((st.dead ? 1 : 0) - dead) * Math.min(1, dt * 6);
      root.rotation.z = dead * (Math.PI / 2);
      inner.position.x = dead * (W / 2 + 0.05);
    },
  };
}
