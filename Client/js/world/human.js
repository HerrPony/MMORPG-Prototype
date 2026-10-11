import * as THREE from 'three';

// Pixel-art humanoids: a blocky rig (head, torso, two arms, two legs) whose faces are painted
// with small procedurally drawn textures (nearest-neighbour filtered => crisp pixels).
// Clothes, armor, skin, hair and style (skeleton, zombie, yeti ...) are all painted, props are boxes.

export const loadHumanAssets = () => Promise.resolve(true);
export const humanAssetsLoaded = () => true;

// ---------- colours & random ----------
const hexCss = (c) => '#' + (c & 0xffffff).toString(16).padStart(6, '0');
export const shade = (c, f) => {
  const r = Math.min(255, Math.round(((c >> 16) & 255) * f)), g = Math.min(255, Math.round(((c >> 8) & 255) * f)), b = Math.min(255, Math.round((c & 255) * f));
  return (r << 16) | (g << 8) | b;
};
export const mix = (a, b, t) => {
  const ch = (s) => Math.round(((a >> s) & 255) * (1 - t) + ((b >> s) & 255) * t);
  return (ch(16) << 16) | (ch(8) << 8) | ch(0);
};
function rng(seed) {
  let s = seed >>> 0;
  return () => { s = (s + 0x6D2B79F5) >>> 0; let t = s; t = Math.imul(t ^ (t >>> 15), t | 1); t ^= t + Math.imul(t ^ (t >>> 7), t | 61); return ((t ^ (t >>> 14)) >>> 0) / 4294967296; };
}
function hashStr(s) { let h = 2166136261; for (let i = 0; i < s.length; i++) { h ^= s.charCodeAt(i); h = Math.imul(h, 16777619); } return h >>> 0; }

// ---------- pixel textures ----------
class Painter {
  constructor(g, w, h, seed) { this.g = g; this.w = w; this.h = h; this.r = rng(seed); }
  fill(c, amt = 0.12, x0 = 0, y0 = 0, w = this.w, h = this.h) {
    for (let y = y0; y < y0 + h; y++) for (let x = x0; x < x0 + w; x++) { this.g.fillStyle = hexCss(shade(c, 1 + (this.r() - 0.5) * amt)); this.g.fillRect(x, y, 1, 1); }
  }
  px(x, y, c) { if (x < 0 || y < 0 || x >= this.w || y >= this.h) return; this.g.fillStyle = hexCss(c); this.g.fillRect(x, y, 1, 1); }
  rect(x, y, w, h, c) { this.g.fillStyle = hexCss(c); this.g.fillRect(x, y, w, h); }
}

const texCache = new Map();
export function pixTex(key, w, h, draw) {
  let t = texCache.get(key);
  if (t) return t;
  const c = document.createElement('canvas');
  c.width = w; c.height = h;
  draw(new Painter(c.getContext('2d'), w, h, hashStr(key)));
  t = new THREE.CanvasTexture(c);
  t.magFilter = THREE.NearestFilter; t.minFilter = THREE.NearestFilter; t.generateMipmaps = false; t.colorSpace = THREE.SRGBColorSpace;
  texCache.set(key, t);
  return t;
}

const M = (c, o = {}) => new THREE.MeshLambertMaterial({ color: c, ...o });
const bx = (w, h, d, c, o = {}) => { const m = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), M(c, o)); m.castShadow = true; return m; };
const TIER_COLORS = [0x8a8a8a, 0x6a9a5a, 0x5a9ac8, 0x4a7ae0, 0x9a6ad8, 0xd8a040, 0xe08040, 0xe05050, 0xd04aa0, 0xffd860];

function tierOf(id) {
  if (!id) return 0;
  if (id.endsWith('_bf')) return 8;
  const t = parseInt(id.split('_').pop(), 10);
  return Number.isFinite(t) ? t : 0;
}
const metalColor = (id) => (id.endsWith('_bf') ? 0x6a5a30 : [0x9a9a9a, 0xa8b0b8, 0xb8c0cc, 0xc4ccd8, 0xaab4d0, 0xc8c8e0, 0xd8d0b0, 0xe8e0c0, 0xf0d890, 0xfff0b0][Math.min(9, tierOf(id))]);
const kindOfItem = (id) => (id || '').split('_')[0];

function armorColor(id) {
  const base = { cloth: 0x4a5ad0, leather: 0x7a5a38, mail: 0xa8aeb8 }[kindOfItem(id)] || 0x888888;
  const t = tierOf(id);
  return mix(base, TIER_COLORS[Math.min(9, t)], id.endsWith('_bf') ? 0.5 : Math.min(0.45, t * 0.05));
}

// ---------- props (all in local space, +Y along the item) ----------
function blade(len, w, color) {
  const g = new THREE.Group();
  const b = bx(w, len, w * 0.3, color, { emissive: 0x222222 }); b.position.y = len / 2 + 0.1; g.add(b);
  const guard = bx(w * 3.4, 0.05, 0.08, 0x8a6a28); guard.position.y = 0.09; g.add(guard);
  const grip = bx(0.05, 0.2, 0.05, 0x4a2a14); grip.position.y = -0.03; g.add(grip);
  return g;
}

function propFor(kind, id) {
  const tier = tierOf(id);
  const tc = TIER_COLORS[Math.min(9, tier)];
  const g = new THREE.Group();
  const glow = (s, c) => { const m = new THREE.Mesh(new THREE.BoxGeometry(s, s, s), new THREE.MeshBasicMaterial({ color: c })); m.rotation.set(0.6, 0.6, 0); return m; };
  switch (kind) {
    case 'sword': return blade(0.5 + tier * 0.025, 0.08, metalColor(id || 'x_0'));
    case 'dagger': return blade(0.28, 0.06, 0xb0b4bc);
    case 'scimitar': { const b = blade(0.55, 0.09, 0xd8d8e0); b.rotation.z = 0.15; return b; }
    case 'axe': {
      const s = bx(0.05, 0.8, 0.05, 0x6b4423); s.position.y = 0.3; g.add(s);
      const h = bx(0.3, 0.22, 0.05, 0xaab0b8); h.position.set(0.13, 0.62, 0); g.add(h);
      return g;
    }
    case 'pick': {
      const s = bx(0.05, 0.7, 0.05, 0x6b4423); s.position.y = 0.25; g.add(s);
      const h = bx(0.5, 0.07, 0.07, 0x8a8a90); h.position.set(0, 0.6, 0); g.add(h);
      return g;
    }
    case 'staff': {
      const s = bx(0.06, 1.7, 0.06, 0x6b4423); s.position.y = 0.45; g.add(s);
      const orb = glow(0.18 + tier * 0.01, tier > 0 ? tc : 0x66aaff); orb.position.y = 1.4; g.add(orb);
      const halo = bx(0.26, 0.04, 0.26, 0xc8a24a); halo.position.y = 1.28; g.add(halo);
      return g;
    }
    case 'witchstaff': {
      const s = bx(0.06, 1.5, 0.06, 0x3a2a1a); s.position.y = 0.4; g.add(s);
      const orb = glow(0.16, 0x9aff4a); orb.position.y = 1.2; g.add(orb);
      return g;
    }
    case 'spear': {
      const s = bx(0.05, 2.1, 0.05, 0x6b4423); s.position.y = 0.5; g.add(s);
      const t = bx(0.1, 0.3, 0.04, 0xc0c4cc); t.position.y = 1.67; g.add(t);
      return g;
    }
    case 'bow': {
      // pixel bow: a connected arc along Y that bulges towards +X (tips at x = 0)
      const col = tier > 2 ? tc : 0x6b4a22;
      const arc = [];
      for (let k = 0; k <= 8; k++) { const y = -0.5 + k / 8; arc.push([0.17 * (1 - Math.pow(y / 0.5, 2)), y]); }
      for (let k = 0; k < 8; k++) {
        const a = arc[k], b2 = arc[k + 1];
        const seg = bx(0.07, Math.hypot(b2[0] - a[0], b2[1] - a[1]) + 0.03, 0.07, col);
        seg.position.set((a[0] + b2[0]) / 2, (a[1] + b2[1]) / 2, 0);
        seg.rotation.z = -Math.atan2(b2[0] - a[0], b2[1] - a[1]);
        g.add(seg);
      }
      const grip = bx(0.1, 0.2, 0.1, 0x3a2412); grip.position.set(0.17, 0, 0); g.add(grip);      // the string is two thin boxes from the tips to the nock; the nock slides back when the bow is drawn
      const halves = [1, -1].map(() => { const h = bx(0.012, 1, 0.012, 0xdddddd); g.add(h); return h; });
      const nock = new THREE.Group(); g.add(nock);
      const shaft = bx(0.7, 0.014, 0.014, 0xd8d0b8); shaft.position.x = 0.35; nock.add(shaft);
      const head = bx(0.08, 0.05, 0.03, 0xc0c4cc); head.position.x = 0.72; nock.add(head);
      const fl = bx(0.1, 0.05, 0.012, 0xc83030); fl.position.x = 0.04; nock.add(fl);
      // draw: how far the nock is pulled back (0..0.38); arrow: whether an arrow is nocked
      g.userData.setDraw = (d, arrow) => {
        const nx = -0.03 - d;
        nock.position.x = nx; nock.visible = arrow;
        halves.forEach((h, i) => {
          const s = i === 0 ? 1 : -1;
          const dx = nx - 0.0, dy = -s * 0.5; // from the tip (0, s*0.5) to the nock
          const len = Math.hypot(dx, dy);
          h.scale.y = len; h.position.set(dx / 2, s * 0.5 + dy / 2, 0);
          h.rotation.z = Math.atan2(dy, dx) - Math.PI / 2;
        });
      };
      g.userData.setDraw(0, true);
      return g;
    }
    case 'rifle': {
      const s = bx(0.05, 1.1, 0.05, 0x333333); s.position.y = 0.3; g.add(s);
      const st = bx(0.08, 0.35, 0.07, 0x6b4423); st.position.y = -0.2; g.add(st);
      return g;
    }
    case 'shield': {
      const body = bx(0.5, 0.62, 0.06, tier > 0 ? tc : 0x9a2a2a); g.add(body);
      const rim = bx(0.54, 0.66, 0.04, metalColor(id || 'x_0')); rim.position.z = -0.01; g.add(rim);
      const boss = bx(0.14, 0.14, 0.05, 0xc9a24a); boss.position.z = 0.05; g.add(boss);
      return g;
    }
    case 'orb': { g.add(bx(0.2, 0.2, 0.2, tc)); return g; }
    case 'candle': {
      g.add(bx(0.06, 0.16, 0.06, 0xe8e0c0));
      const f = new THREE.Mesh(new THREE.BoxGeometry(0.06, 0.07, 0.06), new THREE.MeshBasicMaterial({ color: 0xffaa33 })); f.position.y = 0.12; g.add(f);
      return g;
    }
  }
  return g;
}

export function quiverProp() {
  const g = new THREE.Group();
  const q = bx(0.12, 0.5, 0.1, 0x6b4a28); q.rotation.z = 0.25; g.add(q);
  for (let i = 0; i < 3; i++) { const a = bx(0.02, 0.18, 0.02, 0xd8d0b8); a.position.set(-0.07 + i * 0.04, 0.32, 0); a.rotation.z = 0.25; g.add(a); }
  return g;
}

// head props are centered on the middle of the (0.42 m) head
function headPropFor(kind, id, color) {
  const tier = tierOf(id);
  const mc = kind === 'helmet' && color ? color : id ? metalColor(id) : 0xb9bec8;
  const g = new THREE.Group();
  const add = (m, x = 0, y = 0, z = 0) => { m.position.set(x, y, z); g.add(m); return m; };
  switch (kind) {
    case 'captainhat': {
      add(bx(0.62, 0.14, 0.62, 0x1a1a24), 0, 0.22, 0);
      add(bx(0.46, 0.18, 0.46, 0x1a1a24), 0, 0.34, 0);
      add(bx(0.48, 0.05, 0.48, 0xe8c040), 0, 0.27, 0);
      return g;
    }
    case 'helmet': {
      add(bx(0.48, 0.18, 0.48, mc, { emissive: 0x111111 }), 0, 0.15, 0);
      add(bx(0.06, 0.22, 0.48, mc), 0.21, 0.0, 0); add(bx(0.06, 0.22, 0.48, mc), -0.21, 0.0, 0);
      add(bx(0.48, 0.22, 0.06, mc), 0, 0.0, -0.21);
      add(bx(0.05, 0.16, 0.03, mc), 0, 0.0, 0.23);
      if (tier >= 3) add(bx(0.05, 0.14, 0.36, TIER_COLORS[Math.min(9, tier)]), 0, 0.3, 0);
      return g;
    }
    case 'hood': {
      const c = color || 0x2f4f2a;
      add(bx(0.48, 0.1, 0.48, c), 0, 0.2, 0);
      add(bx(0.48, 0.44, 0.1, c), 0, 0.0, -0.19);
      add(bx(0.07, 0.4, 0.4, c), 0.21, 0.0, -0.02); add(bx(0.07, 0.4, 0.4, c), -0.21, 0.0, -0.02);
      add(bx(0.48, 0.08, 0.06, shade(c, 0.8)), 0, 0.13, 0.2);
      return g;
    }
    case 'wizardhat': {
      const c = color || 0x23358a;
      add(bx(0.72, 0.04, 0.72, c), 0, 0.22, 0);
      add(bx(0.44, 0.04, 0.44, 0xd8b13a), 0, 0.26, 0);
      add(bx(0.4, 0.2, 0.4, c), 0, 0.36, 0);
      add(bx(0.3, 0.2, 0.3, c), 0, 0.55, -0.02);
      add(bx(0.2, 0.18, 0.2, c), 0, 0.73, -0.06);
      add(bx(0.1, 0.14, 0.1, c), 0, 0.88, -0.12);
      return g;
    }
    case 'straw': {
      add(bx(0.78, 0.03, 0.78, 0xe0c070), 0, 0.2, 0);
      add(bx(0.42, 0.12, 0.42, 0xd0b060), 0, 0.27, 0);
      return g;
    }
    case 'turban': {
      const c = color || 0xe8e0c8;
      add(bx(0.5, 0.2, 0.5, c), 0, 0.2, 0); add(bx(0.4, 0.1, 0.4, shade(c, 0.9)), 0, 0.35, 0);
      return g;
    }
    case 'crown': {
      add(bx(0.46, 0.1, 0.46, 0xe8c040, { emissive: 0x443300 }), 0, 0.25, 0);
      for (const [x, z] of [[-0.18, -0.18], [0.18, -0.18], [-0.18, 0.18], [0.18, 0.18], [0, 0.2], [0, -0.2]]) add(bx(0.06, 0.12, 0.06, 0xe8c040, { emissive: 0x443300 }), x, 0.36, z);
      return g;
    }
    case 'witchhat': {
      const c = 0x2a1a3a;
      add(bx(0.8, 0.03, 0.8, c), 0, 0.22, 0);
      add(bx(0.42, 0.2, 0.42, c), 0, 0.34, 0);
      add(bx(0.3, 0.2, 0.3, c), 0.02, 0.53, 0);
      add(bx(0.2, 0.18, 0.2, c), 0.07, 0.71, 0);
      add(bx(0.1, 0.14, 0.1, c), 0.14, 0.86, 0);
      add(bx(0.43, 0.05, 0.43, 0x8a3ad8), 0, 0.27, 0);
      return g;
    }
    case 'mask': {
      add(bx(0.44, 0.14, 0.04, color || 0x8a1a1a), 0, -0.04, 0.22);
      add(bx(0.46, 0.14, 0.46, 0x2a2a2a), 0, 0.18, 0);
      return g;
    }
    case 'horns': {
      for (const s of [-1, 1]) { add(bx(0.07, 0.26, 0.07, 0xe8e0c8), s * 0.26, 0.28, 0).rotation.z = -s * 0.5; add(bx(0.06, 0.12, 0.06, 0xe8e0c8), s * 0.37, 0.43, 0).rotation.z = -s * 0.9; }
      return g;
    }
    case 'pharaoh': {
      add(bx(0.5, 0.3, 0.5, 0xe8c040, { emissive: 0x443300 }), 0, 0.17, 0);
      for (const s of [-1, 1]) add(bx(0.08, 0.5, 0.2, 0x2a4aa8), s * 0.28, -0.05, 0.0);
      return g;
    }
    case 'ears': {
      for (const s of [-1, 1]) add(bx(0.24, 0.07, 0.05, color || 0x8a6a3a), s * 0.32, 0.06, 0).rotation.z = s * 0.35;
      return g;
    }
  }
  return g;
}

function shoulderProp(id) {
  const mc = shade(metalColor(id), 0.85);
  const g = new THREE.Group();
  g.add(bx(0.27, 0.12, 0.28, mc, { emissive: 0x111111 }));
  const side = bx(0.07, 0.16, 0.26, shade(mc, 0.85)); side.position.set(0.0, -0.08, 0); g.add(side);
  if (tierOf(id) >= 3) { const sp = bx(0.05, 0.12, 0.05, TIER_COLORS[Math.min(9, tierOf(id))]); sp.position.y = 0.1; g.add(sp); }
  return g;
}

function cloakProp(id, color) {
  const g = new THREE.Group();
  const c = bx(0.5, 0.95, 0.04, color || [0x6a3a2a, 0x3a5a3a, 0x3a3a7a, 0x7a2a2a][tierOf(id) % 4]);
  c.position.set(0, -0.2, 0); g.add(c);
  return g;
}

// ---------- appearance ----------
// style: normal | skeleton | zombie | bandage | gold | fur | kobold | water
function paletteOf(spec) {
  const style = spec.style || 'normal';
  const P = { style, skin: spec.skin ?? 0xe0ac82, hair: spec.hair ?? 0x4a3020, hairOn: spec.hair != null && !spec.bald, eyes: 0x1a1a22, shirt: spec.tint ?? 0x777777, female: !!spec.female };
  P.pants = spec.pants ?? shade(P.shirt, 0.55);
  P.boots = spec.boots ?? 0x2e2218;
  switch (style) {
    case 'skeleton': P.skin = 0xe6e2d0; P.shirt = 0xdcd8c4; P.pants = 0xcfcbb8; P.boots = 0xcfcbb8; P.hairOn = false; break;
    case 'zombie': P.skin = 0x6f8f5a; P.hairOn = false; P.eyes = 0xcc2222; P.pants = 0x4a4a3a; break;
    case 'bandage': P.skin = 0xe0d8b8; P.shirt = 0xe0d8b8; P.pants = 0xd4ccac; P.boots = 0xd4ccac; P.hairOn = false; break;
    case 'gold': P.skin = 0xd8b038; P.shirt = 0xe0b840; P.pants = 0xc89c30; P.boots = 0xb08828; P.hairOn = false; break;
    case 'fur': P.skin = spec.tint ?? 0xe8f2ff; P.shirt = P.skin; P.pants = shade(P.skin, 0.94); P.boots = shade(P.skin, 0.85); P.hairOn = false; break;
    case 'kobold': P.skin = 0xa8663a; P.hairOn = false; P.eyes = 0xffd020; break;
    case 'lizard': P.skin = 0x3f9a5a; P.hairOn = false; P.eyes = 0xffd020; break;
    case 'ghost': P.skin = 0xbcd8f0; P.shirt = spec.tint ?? 0x9ab8e0; P.pants = shade(P.shirt, 0.8); P.boots = P.pants; P.hairOn = false; P.eyes = 0xffffff; break;
    case 'stone': P.skin = spec.tint ?? 0x8c8c86; P.shirt = P.skin; P.pants = shade(P.skin, 0.92); P.boots = shade(P.skin, 0.8); P.hairOn = false; P.eyes = 0xffa030; break;
    case 'lava': P.skin = 0x2e1c1a; P.shirt = P.skin; P.pants = P.skin; P.boots = 0x1a1010; P.hairOn = false; P.eyes = 0xffe040; break;
    case 'obsidian': P.skin = 0x1c1c26; P.shirt = 0x242432; P.pants = 0x1a1a24; P.boots = 0x101018; P.hairOn = false; P.eyes = 0xff7a30; break;
    case 'demon': P.skin = 0xa02c24; P.hairOn = false; P.eyes = 0xffe040; P.pants = 0x3a1410; break;
    case 'water': P.skin = spec.tint ?? 0x2a6ad8; P.shirt = P.skin; P.pants = P.skin; P.boots = P.skin; P.hairOn = false; P.eyes = 0xffffff; break;
  }
  return P;
}

const SPEC_KEY = (P, extra = '') => `${P.style}|${P.skin}|${P.hair}|${P.hairOn}|${P.eyes}|${P.shirt}|${P.pants}|${P.boots}|${P.female}|${extra}`;

function headTex(P, face) { // face: front | back | top | bottom | sideR | sideL
  return pixTex('h' + face + SPEC_KEY(P), 8, 8, (p) => {
    p.fill(P.skin, P.style !== 'normal' && P.style !== 'fur' ? 0.2 : 0.1);
    const hl = P.female ? 7 : 3; // how far the hair reaches down on the sides / back
    if (face === 'top') { if (P.hairOn) p.fill(P.hair, 0.3); return; }
    if (face === 'bottom') return;
    if (face === 'back') { if (P.hairOn) p.fill(P.hair, 0.3, 0, 0, 8, hl + 1); return; }
    if (face === 'sideR' || face === 'sideL') {
      if (P.hairOn) {
        p.fill(P.hair, 0.3, 0, 0, 8, 2);
        p.fill(P.hair, 0.3, face === 'sideR' ? 4 : 0, 2, 4, hl - 1); // the back of the head
      }
      return;
    }
    // front
    if (P.hairOn) {
      p.fill(P.hair, 0.3, 0, 0, 8, 2);
      p.fill(P.hair, 0.3, 0, 2, 1, hl - 1); p.fill(P.hair, 0.3, 7, 2, 1, hl - 1);
      p.px(1, 2, P.hair); p.px(6, 2, P.hair); p.px(3, 2, P.hair);
    }
    if (P.style === 'skeleton') {
      p.rect(1, 3, 2, 2, 0x111111); p.rect(5, 3, 2, 2, 0x111111); p.px(3, 5, shade(P.skin, 0.6));
      for (let x = 2; x <= 5; x++) p.px(x, 6, x % 2 ? 0x222222 : 0xf0f0e8);
      return;
    }
    if (P.style === 'bandage') {
      for (let y = 1; y < 8; y += 2) p.rect(0, y, 8, 1, 0xb8b090);
      p.rect(1, 4, 2, 1, 0x1a1a1a); p.rect(5, 4, 2, 1, 0x1a1a1a); p.px(1, 4, 0xcc3322); p.px(6, 4, 0xcc3322);
      return;
    }
    if (P.style === 'gold') { p.rect(1, 4, 2, 1, 0x101018); p.rect(5, 4, 2, 1, 0x101018); p.rect(0, 3, 8, 1, 0x2a4aa8); p.rect(3, 6, 2, 1, shade(P.skin, 0.6)); return; }
    if (P.style === 'fur') { p.rect(1, 3, 6, 4, shade(P.skin, 0.82)); p.px(2, 4, 0x2a4a8a); p.px(5, 4, 0x2a4a8a); p.rect(2, 6, 4, 1, 0x555a66); p.px(2, 6, 0xffffff); p.px(5, 6, 0xffffff); return; }
    if (P.style === 'water') { for (let i = 0; i < 10; i++) p.px((i * 5) % 8, (i * 3) % 8, shade(P.skin, 1.5)); p.rect(1, 3, 2, 2, 0xffffff); p.rect(5, 3, 2, 2, 0xffffff); return; }
    if (P.style === 'lizard') { p.rect(0, 3, 8, 1, shade(P.skin, 0.75)); p.px(1, 4, P.eyes); p.px(2, 4, 0x101008); p.px(5, 4, 0x101008); p.px(6, 4, P.eyes); p.rect(2, 6, 4, 1, 0x1a3a20); p.px(2, 7, 0xf0f0e0); p.px(5, 7, 0xf0f0e0); return; }
    if (P.style === 'stone') { p.rect(1, 3, 2, 1, P.eyes); p.rect(5, 3, 2, 1, P.eyes); p.rect(2, 6, 4, 1, 0x1c1c1c); p.px(0, 1, shade(P.skin, 0.6)); p.px(1, 2, shade(P.skin, 0.6)); p.px(7, 5, 0x5a7a3a); p.px(6, 6, 0x5a7a3a); return; }
    if (P.style === 'lava') { for (let i = 0; i < 8; i++) p.px((i * 3) % 8, (i * 5 + 1) % 8, 0xff6a10); p.rect(1, 3, 2, 1, P.eyes); p.rect(5, 3, 2, 1, P.eyes); p.rect(2, 6, 4, 1, 0xff8a20); return; }
    if (P.style === 'obsidian') { p.rect(0, 0, 8, 1, 0x3a3a50); p.rect(1, 4, 6, 1, P.eyes); p.px(0, 4, 0x3a3a50); p.px(7, 4, 0x3a3a50); p.rect(2, 7, 4, 1, 0x3a3a50); return; }
    if (P.style === 'demon') { p.rect(1, 3, 2, 1, 0x101008); p.rect(5, 3, 2, 1, 0x101008); p.px(1, 4, P.eyes); p.px(2, 4, P.eyes); p.px(5, 4, P.eyes); p.px(6, 4, P.eyes); p.rect(2, 6, 4, 1, 0x2a0a08); p.px(2, 7, 0xf0f0e0); p.px(5, 7, 0xf0f0e0); return; }
    if (P.style === 'ghost') { p.rect(1, 3, 2, 2, 0xffffff); p.rect(5, 3, 2, 2, 0xffffff); p.rect(3, 6, 2, 2, 0x203050); return; }
    if (P.style === 'kobold') { p.px(1, 4, P.eyes); p.px(2, 4, 0x201008); p.px(5, 4, 0x201008); p.px(6, 4, P.eyes); p.rect(3, 5, 2, 1, shade(P.skin, 0.75)); p.rect(2, 6, 4, 1, 0x5a2a14); p.px(3, 6, 0xffffff); p.px(5, 6, 0xffffff); return; }
    if (P.style === 'zombie') { p.px(1, 4, 0xddddaa); p.px(2, 4, P.eyes); p.px(5, 4, P.eyes); p.px(6, 4, 0xddddaa); p.rect(2, 6, 4, 1, 0x2a1a1a); p.px(0, 5, 0x4a5a3a); p.px(7, 3, 0x4a5a3a); return; }
    // normal
    p.px(1, 4, 0xffffff); p.px(2, 4, P.eyes); p.px(5, 4, P.eyes); p.px(6, 4, 0xffffff);
    if (P.hairOn) for (const x of [1, 2, 5, 6]) p.px(x, 3, shade(P.hair, 0.8));
    p.rect(3, 5, 2, 1, shade(P.skin, 0.85));
    p.rect(3, 6, 2, 1, mix(P.skin, 0xa04040, 0.45));
  });
}

// worn armor panel with a pattern depending on the armor kind
function armorPanel(p, w, h, base, kind, tier) {
  p.fill(base, 0.12);
  if (kind === 'mail') {
    for (let y = 0; y < h; y++) for (let x = 0; x < w; x++) if ((x + y) % 2 === 0) p.px(x, y, shade(base, 1.16));
    p.rect(0, 0, w, 1, shade(base, 1.3));
  } else if (kind === 'leather') {
    for (let y = 2; y < h; y += 4) p.rect(0, y, w, 1, shade(base, 0.75));
    p.rect(w >> 1, 0, 1, h, shade(base, 0.7));
  } else if (kind === 'cloth') {
    p.rect(0, h - 1, w, 1, 0xd8b13a); if (w > 4) p.rect(w >> 1, 0, 1, h, 0xd8b13a);
  }
  if (tier >= 5) for (let i = 0; i < w; i += 2) p.px(i, 1, TIER_COLORS[Math.min(9, tier)]);
}

// style specific speckles (scales, cracks, lava veins ...)
function decorate(p, P, w, h) {
  const speck = (c, n) => { for (let i = 0; i < n; i++) p.px(Math.floor(p.r() * w), Math.floor(p.r() * h), c); };
  const n = Math.floor(w * h * 0.16);
  switch (P.style) {
    case 'lizard': speck(shade(P.skin, 0.7), n); speck(shade(P.skin, 1.2), n >> 1); break;
    case 'stone': speck(shade(P.skin, 0.62), n); speck(0x5a7a3a, n >> 2); speck(shade(P.skin, 1.15), n >> 1); break;
    case 'lava': speck(0xff6a10, n); speck(0xffb030, n >> 2); break;
    case 'obsidian': speck(0x4a4a66, n); speck(0x6a3a8a, n >> 3); break;
    case 'ghost': speck(0xe8f4ff, n); break;
    case 'demon': speck(0x6a1814, n); break;
  }
}

function torsoTex(P, A, face) { // face: front | back | side
  const chest = A.chest;
  const w = face === 'side' ? 4 : 8;
  return pixTex('t' + face + SPEC_KEY(P, `${chest}|${A.waist}`), w, 12, (p) => {
    if (chest) armorPanel(p, w, 12, armorColor(chest), kindOfItem(chest), tierOf(chest));
    else p.fill(P.shirt, P.style === 'normal' ? 0.14 : 0.2);
    if (!chest) decorate(p, P, w, 12);
    if (!chest && P.style === 'skeleton') { for (const y of [2, 4, 6]) p.rect(face === 'side' ? 0 : 1, y, face === 'side' ? 4 : 6, 1, 0x4a4638); if (face !== 'side') p.rect(3, 1, 2, 7, 0xb8b4a0); }
    if (!chest && P.style === 'bandage') for (let y = 1; y < 12; y += 2) p.rect(0, y, w, 1, 0xb8b090);
    if (!chest && P.style === 'zombie' && face === 'front') { p.rect(1, 3, 3, 2, 0x4a3a2a); p.rect(5, 6, 2, 3, 0x4a3a2a); }
    if (!chest && P.style === 'water') for (let i = 0; i < 8; i++) p.px((i * 3) % w, (i * 5) % 12, shade(P.shirt, 1.5));
    if (!chest && P.style === 'fur') for (let i = 0; i < 14; i++) p.px((i * 3) % w, (i * 5) % 12, shade(P.shirt, 0.88));
    if (!chest && P.style === 'normal' && face === 'front') { p.px(3, 0, P.skin); p.px(4, 0, P.skin); p.px(3, 1, shade(P.shirt, 0.8)); p.px(4, 1, shade(P.shirt, 0.8)); }
    const clothed = P.style === 'normal' || P.style === 'zombie' || P.style === 'kobold' || P.style === 'lizard' || P.style === 'demon';
    if (A.waist) { p.rect(0, 9, w, 1, armorColor(A.waist)); if (face === 'front') p.px(3, 9, 0xe8c040); }
    else if (clothed) { p.rect(0, 9, w, 1, 0x4a3520); if (face === 'front') p.px(3, 9, 0xc8a040); }
    if (!chest && clothed) p.fill(P.pants, 0.12, 0, 10, w, 2);
  });
}

function armTex(P, A) {
  const chest = A.chest, hands = A.hands;
  return pixTex('a' + SPEC_KEY(P, `${chest}|${hands}`), 4, 12, (p) => {
    const clothed = P.style === 'normal' || P.style === 'kobold' || P.style === 'zombie' || P.style === 'lizard' || P.style === 'demon';
    const sleeve = chest ? 9 : (clothed ? 5 : 12);
    p.fill(P.skin, 0.12);
    if (chest) armorPanel(p, 4, sleeve, armorColor(chest), kindOfItem(chest), tierOf(chest));
    else p.fill(P.shirt, 0.14, 0, 0, 4, sleeve);
    if (!chest) decorate(p, P, 4, 12);
    if (!chest && P.style === 'skeleton') { p.fill(0xe6e2d0, 0.1); p.rect(1, 0, 2, 12, 0xc8c4b0); }
    if (!chest && P.style === 'bandage') for (let y = 1; y < 12; y += 2) p.rect(0, y, 4, 1, 0xb8b090);
    if (hands) { p.fill(armorColor(hands), 0.12, 0, 9, 4, 3); p.rect(0, 9, 4, 1, shade(armorColor(hands), 1.25)); }
    else if (clothed) p.fill(P.skin, 0.1, 0, 10, 4, 2);
  });
}

function legTex(P, A) {
  const legs = A.legs, feet = A.feet;
  return pixTex('l' + SPEC_KEY(P, `${legs}|${feet}`), 4, 12, (p) => {
    if (legs) armorPanel(p, 4, 9, armorColor(legs), kindOfItem(legs), tierOf(legs)); else { p.fill(P.pants, 0.14, 0, 0, 4, 9); decorate(p, P, 4, 9); }
    if (!legs && P.style === 'skeleton') { p.fill(0xe6e2d0, 0.1, 0, 0, 4, 9); p.rect(1, 0, 2, 9, 0xc8c4b0); }
    if (!legs && P.style === 'bandage') for (let y = 1; y < 9; y += 2) p.rect(0, y, 4, 1, 0xb8b090);
    if (!legs && P.style === 'zombie') p.rect(0, 4, 2, 2, 0x3a3a2a);
    const boot = feet ? armorColor(feet) : P.boots;
    p.fill(boot, 0.14, 0, 9, 4, 3); p.rect(0, 9, 4, 1, shade(boot, 1.3));
  });
}

const solidTex = (c) => pixTex('s' + c, 4, 4, (p) => p.fill(c, 0.16));

/**
 * Creates an animated pixel humanoid.
 * spec: { tint (shirt), skin, hair, bald, female, style, pants, boots, emissive, scale, bodyScale, headScale, opacity,
 *         props, head, cloak, shoulders, chest, legs, feet, hands, waist, extra, slow, armsForward }
 */
export function makeHumanoid(spec0 = {}) {
  let spec = spec0;
  const root = new THREE.Group();
  const scale = spec0.scale || 1;
  const inner = new THREE.Group(); inner.scale.setScalar(scale); root.add(inner);
  const body = new THREE.Group(); inner.add(body);
  const bs = spec0.bodyScale || [1, 1, 1];
  body.scale.set(bs[0], bs[1], bs[2]);

  const torsoW = spec0.female ? 0.42 : 0.5;
  const mk = (t) => {
    const m = new THREE.MeshLambertMaterial({ map: t });
    if (spec.emissive != null) m.emissive = new THREE.Color(spec.emissive);
    if (spec.opacity != null) { m.transparent = true; m.opacity = spec.opacity; }
    return m;
  };
  const part = (w, h, d) => { const m = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), []); m.castShadow = spec0.opacity == null; return m; };

  // legs
  const legs = [], feetMeshes = [];
  for (const x of [-0.125, 0.125]) {
    const pv = new THREE.Group(); pv.position.set(x, 0.72, 0); body.add(pv);
    const l = part(0.24, 0.72, 0.26); l.position.y = -0.36; pv.add(l);
    const f = part(0.25, 0.1, 0.36); f.position.set(0, -0.67, 0.05); pv.add(f); feetMeshes.push(f);
    legs.push({ pv, mesh: l });
  }
  const torso = part(torsoW, 0.66, spec0.female ? 0.26 : 0.3); torso.position.y = 1.05; body.add(torso);
  const chestA = new THREE.Group(); chestA.position.y = 1.05; body.add(chestA);
  const hipsA = new THREE.Group(); hipsA.position.y = 0.72; body.add(hipsA);
  // arms (the character's right arm is on -x because the model looks along +z)
  const arms = [];
  for (const s of [-1, 1]) {
    const pv = new THREE.Group(); pv.position.set(s * (torsoW / 2 + 0.09), 1.35, 0); body.add(pv);
    const a = part(0.18, 0.66, 0.2); a.position.y = -0.33; pv.add(a);
    const hand = new THREE.Group(); hand.position.set(0, -0.64, 0.02); pv.add(hand);
    const sh = new THREE.Group(); pv.add(sh);
    arms.push({ pv, mesh: a, hand, sh, s });
  }
  const armR = arms[0], armL = arms[1];
  const headPv = new THREE.Group(); headPv.position.y = 1.38; body.add(headPv);
  if (spec0.headScale) headPv.scale.setScalar(spec0.headScale);
  const headMesh = part(0.42, 0.42, 0.42); headMesh.position.y = 0.21; headPv.add(headMesh);
  const headA = new THREE.Group(); headA.position.y = 0.21; headPv.add(headA);

  const attached = [];
  const clear = () => { attached.forEach((o) => o.removeFromParent()); attached.length = 0; };
  const put = (parent, obj, pos = [0, 0, 0], rot = [0, 0, 0], sc = 1) => { obj.position.set(pos[0], pos[1], pos[2]); obj.rotation.set(rot[0], rot[1], rot[2]); obj.scale.setScalar(sc); parent.add(obj); attached.push(obj); };
  const POINT = { sword: 1.2, dagger: 1.2, scimitar: 1.2, axe: 0.9, pick: 0.9, rifle: 0.9 };

  let bowObj = null;
  function dress(next) {
    if (next) spec = next;
    clear();
    bowObj = null;
    const P = paletteOf(spec);
    const A = { chest: spec.chest, legs: spec.legs, feet: spec.feet, hands: spec.hands, waist: spec.waist };
    const T = (f) => mk(torsoTex(P, A, f));
    torso.material = [T('side'), T('side'), T('side'), T('side'), T('front'), T('back')];
    const hm = (f) => mk(headTex(P, f));
    headMesh.material = [hm('sideR'), hm('sideL'), hm('top'), hm('bottom'), hm('front'), hm('back')];
    const am = mk(armTex(P, A));
    for (const a of arms) a.mesh.material = am;
    const lm = mk(legTex(P, A));
    for (const l of legs) l.mesh.material = lm;
    const fm = mk(solidTex(spec.feet ? armorColor(spec.feet) : P.boots));
    for (const f of feetMeshes) f.material = fm;

    for (const p of spec.props || []) {
      const obj = propFor(p.kind, p.id);
      const left = p.slot === 'left';
      const arm = left ? armL : armR;
      if (left && p.kind === 'shield') put(arm.hand, obj, [0.13, 0.0, 0.05], [0, Math.PI / 2, 0], p.scale || 1);
      else if (left && p.kind === 'bow') {
        // the left arm points forward: hand-local +Z is up, -Y is forward; the bow stands upright with its belly facing the target
        put(arm.hand, obj, [0, 0, 0.02], [0, 0, 0], (p.scale || 1) * 1.25);
        obj.quaternion.setFromRotationMatrix(new THREE.Matrix4().makeBasis(new THREE.Vector3(0, -1, 0), new THREE.Vector3(0, 0, 1), new THREE.Vector3(-1, 0, 0)));
        bowObj = obj;
      }
      else put(arm.hand, obj, [0, 0, 0.04], [POINT[p.kind] ?? 0, 0, 0], p.scale || 1);
    }
    if (spec.head) put(headA, headPropFor(spec.head.kind, spec.head.id, spec.head.color));
    if (spec.shoulders) for (const a of arms) put(a.sh, shoulderProp(spec.shoulders), [a.s * 0.02, 0.05, 0]);
    if (spec.chest && kindOfItem(spec.chest) === 'cloth') { // robe: a skirt over the legs
      const c = armorColor(spec.chest);
      const sk = bx(torsoW + 0.1, 0.5, 0.38, c);
      sk.material = mk(pixTex('skirt' + c, 8, 8, (p) => armorPanel(p, 8, 8, c, 'cloth', tierOf(spec.chest))));
      put(hipsA, sk, [0, -0.12, 0]);
    }
    if (spec.cloak) put(chestA, cloakProp(spec.cloak.id, spec.cloak.color), [0, 0.0, -0.2]);
    if (spec.extra) for (const ex of spec.extra) put(ex.bone === 'head' ? headA : chestA, ex.obj, ex.off || [0, 0, 0], ex.rot || [0, 0, 0], ex.scale || 1);
  }
  dress(spec0);

  let attackT = 0, castT = 0, dead = 0, phase = 0, t = Math.random() * 10, sw = 0;
  return {
    obj: root, height: 1.8 * scale * (bs[1] || 1) + 0.1, kind: 'human',
    setTint(c) { spec = { ...spec, tint: c }; },
    dress,
    bowDelay: spec0.props?.some((p) => p.kind === 'bow') ? 0.3 : 0, // the arrow leaves the string this long after the shot started
    attack() { attackT = 0.01; },
    cast(on) { castT = on ? 1 : 0; },
    update(dt, st) {
      t += dt;
      const sp = st.speed;
      const moving = sp > 0.08 && !st.airborne;
      const slow = spec.slow ? 0.6 : 1;
      if (moving) phase += dt * Math.max(4, sp * 13) * slow;
      sw += ((moving ? Math.min(1, 0.4 + sp * 0.6) * slow : 0) - sw) * Math.min(1, dt * 12);
      const s = Math.sin(phase) * sw;
      legs[0].pv.rotation.x = s * 0.9; legs[1].pv.rotation.x = -s * 0.9;
      if (st.airborne) { legs[0].pv.rotation.x = -0.5; legs[1].pv.rotation.x = 0.4; }
      const breathe = Math.sin(t * 1.8) * 0.03;
      let rArm = -s * 0.8 + breathe, lArm = s * 0.8 - breathe;
      if (spec.props?.some((p) => p.slot !== 'left')) rArm = rArm * 0.5 - 0.15; // holding a weapon
      if (spec.armsForward) { rArm = -1.45 + s * 0.1; lArm = -1.45 - s * 0.1; }
      if (castT > 0) { rArm = -1.7; lArm = -1.7; }
      let rYaw = 0;
      if (bowObj) {
        // archer stance: bow arm forward, string hand near the chest; attack = draw (0-0.3 s), release, recoil
        lArm = -1.4; rArm = -1.3; rYaw = 0.55;
        let draw = 0.04, arrow = true;
        if (attackT > 0) {
          attackT += dt;
          const k = attackT / 0.3; // 0..1 = drawing
          if (k < 1) { const e = k * k * (3 - 2 * k); draw = 0.04 + 0.34 * e; rYaw = 0.55 + 0.45 * e; }
          else if (attackT < 0.42) { draw = 0.02; rYaw = 0.6 - ((attackT - 0.3) / 0.12) * 0.15; arrow = false; lArm = -1.4 + 0.08 * Math.sin((attackT - 0.3) * 60); }
          else if (attackT < 0.65) { draw = 0.04; arrow = ((attackT - 0.42) / 0.23) > 0.4; }
          else attackT = 0;
          body.rotation.y = k < 1 ? -0.25 * k : -0.25 * Math.max(0, 1 - (attackT - 0.3) * 4);
          if (!attackT) body.rotation.y = 0;
        }
        bowObj.userData.setDraw?.(draw, arrow);
      } else if (attackT > 0) {
        attackT += dt;
        const k = Math.min(1, attackT / 0.4);
        rArm = -0.4 - Math.sin(k * Math.PI) * 2.0;
        body.rotation.y = Math.sin(k * Math.PI) * 0.3;
        if (k >= 1) { attackT = 0; body.rotation.y = 0; }
      }
      armR.pv.rotation.x = rArm; armL.pv.rotation.x = lArm;
      armR.pv.rotation.y = rYaw; armL.pv.rotation.y = 0;
      armR.pv.rotation.z = 0.04; armL.pv.rotation.z = -0.04;
      body.position.y = moving ? Math.abs(Math.sin(phase)) * 0.04 * sw : 0;
      headPv.rotation.y = moving ? 0 : Math.sin(t * 0.6) * 0.12;
      dead += ((st.dead ? 1 : 0) - dead) * Math.min(1, dt * 6); // death: lie on the back
      inner.rotation.x = -dead * (Math.PI / 2);
      inner.position.y = dead * 0.22 * scale;
    },
  };
}
