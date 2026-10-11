export const $ = (id) => document.getElementById(id);

export function el(tag, cls, html) {
  const e = document.createElement(tag);
  if (cls) e.className = cls;
  if (html !== undefined) e.innerHTML = html;
  return e;
}

export function smoothstep(a, b, x) {
  const t = Math.max(0, Math.min(1, (x - a) / (b - a)));
  return t * t * (3 - 2 * t);
}

// ---- world grids (sent by the server) ----
// nx/nz grid vertices, x0/z0 world position of grid cell (0,0); two continents share the grid.
export const W = { nx: 0, nz: 0, cell: 2, x0: -300, z0: -300, hm: null, bm: null, ready: false, conts: [] };
export const BIOME_KEYS = ['sea', 'forest', 'plains', 'swamp', 'snow', 'desert', 'jungle', 'blossom', 'mesa', 'volcano', 'volcano'];
export const LAVA = 10;
export const contOf = (x) => (x > 420 ? 1 : 0);

export function setWorldData(sc) {
  W.nx = sc.nx; W.nz = sc.nz; W.cell = sc.cell; W.x0 = sc.x0; W.z0 = sc.z0; W.conts = sc.conts;
  const hb = Uint8Array.from(atob(sc.hm), (c) => c.charCodeAt(0));
  const i16 = new Int16Array(hb.buffer, hb.byteOffset, hb.byteLength / 2);
  W.hm = new Float32Array(i16.length);
  for (let i = 0; i < i16.length; i++) W.hm[i] = i16[i] / 20;
  W.bm = Uint8Array.from(atob(sc.bm), (c) => c.charCodeAt(0));
  W.ready = true;
}

function gridAt(i, j) {
  i = Math.max(0, Math.min(W.nx - 1, i));
  j = Math.max(0, Math.min(W.nz - 1, j));
  return W.hm[j * W.nx + i];
}

// Terrain height; dungeons (far south) are flat.
export function heightAt(x, z) {
  if (z > 500) return 0;
  if (!W.ready) return 0;
  const gx = (x - W.x0) / W.cell, gz = (z - W.z0) / W.cell;
  const i = Math.floor(gx), j = Math.floor(gz);
  const fx = gx - i, fz = gz - j;
  const a = gridAt(i, j), b = gridAt(i + 1, j), c = gridAt(i, j + 1), d = gridAt(i + 1, j + 1);
  return a + (b - a) * fx + (c - a) * fz + (a - b - c + d) * fx * fz;
}

// 0 sea, 1 forest, 2 plains, 3 swamp, 4 snow, 5 desert, 6 jungle, 7 blossom, 8 mesa, 9 volcano, 10 lava
export function biomeAt(x, z) {
  if (!W.ready || z > 500) return 1;
  const i = Math.max(0, Math.min(W.nx - 1, Math.round((x - W.x0) / W.cell)));
  const j = Math.max(0, Math.min(W.nz - 1, Math.round((z - W.z0) / W.cell)));
  return W.bm[j * W.nx + i];
}

export const CLASS_COLORS = { warrior: '#c79c6e', hunter: '#abd473', mage: '#69ccf0' };
export const CLASS_ICONS = { warrior: '⚔️', hunter: '🏹', mage: '🔮' };
export const SKIN = [0xf1c9a5, 0xe0ac82, 0xc68a5d, 0x8d5a3b, 0x5b3a29];
export const HAIR = [0x2b1b10, 0x6b4423, 0xc9a24a, 0xb0361d, 0x9aa0a6, 0x2c4ea8];
export const hex = (n) => '#' + n.toString(16).padStart(6, '0');

export function fmtCopper(c) {
  const g = Math.floor(c / 10000), s = Math.floor(c / 100) % 100, k = c % 100;
  let h = '';
  if (g) h += `<span class="g">${g}g</span> `;
  if (g || s) h += `<span class="s">${s}s</span> `;
  h += `<span class="c">${k}k</span>`;
  return `<span class="coin">${h}</span>`;
}

export function angleLerp(a, b, t) {
  let d = ((b - a + Math.PI) % (Math.PI * 2) + Math.PI * 2) % (Math.PI * 2) - Math.PI;
  return a + d * t;
}

export function esc(s) {
  return String(s).replace(/[&<>"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c]));
}

export function distToSeg(px, pz, ax, az, bx, bz) {
  const dx = bx - ax, dz = bz - az;
  const l2 = dx * dx + dz * dz;
  const t = l2 > 0 ? Math.max(0, Math.min(1, ((px - ax) * dx + (pz - az) * dz) / l2)) : 0;
  return Math.hypot(px - (ax + dx * t), pz - (az + dz * t));
}

export function distToPoly(x, z, pts) {
  let d = Infinity;
  for (let i = 0; i + 1 < pts.length; i++) d = Math.min(d, distToSeg(x, z, pts[i][0], pts[i][1], pts[i + 1][0], pts[i + 1][1]));
  return d;
}

let toastTimer;
export function toast(msg) {
  const t = $('toast');
  t.textContent = msg;
  t.classList.add('show');
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => t.classList.remove('show'), 2500);
}
