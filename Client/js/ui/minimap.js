import { S } from '../state.js';
import { $, W, distToPoly } from '../util.js';

const PPM = 2;          // base map pixels per meter
const VIEW = 52;        // visible radius in meters
let base = null, half = 300;

export const BIOME_RGB = [[40, 90, 140], [62, 122, 58], [170, 182, 84], [86, 106, 60], [226, 234, 244], [212, 182, 106]];

// Renders the whole continent into an offscreen canvas (also used by the world map).
export function renderWorldCanvas(scale, sc) {
  const size = Math.round(W.half * 2 * scale);
  const cv = document.createElement('canvas');
  cv.width = cv.height = size;
  const g = cv.getContext('2d');
  const img = g.createImageData(size, size);
  for (let py = 0; py < size; py++) {
    for (let px = 0; px < size; px++) {
      const gx = Math.min(W.n - 1, Math.floor(px / scale / W.cell)), gz = Math.min(W.n - 1, Math.floor(py / scale / W.cell));
      const idx = gz * W.n + gx;
      const h = W.hm[idx];
      let col = BIOME_RGB[W.bm[idx]];
      let shade = 0.82 + Math.max(-0.2, Math.min(0.3, h * 0.012));
      let r = col[0] * shade, gg = col[1] * shade, b = col[2] * shade;
      if (h < 0 && W.bm[idx] !== 0) { r = 50; gg = 110; b = 160; } // rivers and lakes
      const o = (py * size + px) * 4;
      img.data[o] = r; img.data[o + 1] = gg; img.data[o + 2] = b; img.data[o + 3] = 255;
    }
  }
  g.putImageData(img, 0, 0);
  g.lineCap = 'round'; g.lineJoin = 'round';
  g.strokeStyle = 'rgba(150,115,70,0.95)'; g.lineWidth = Math.max(1.5, 2.2 * scale);
  for (const rd of sc.roads) {
    g.beginPath();
    rd.forEach(([x, z], i) => { const px = (x + W.half) * scale, py = (z + W.half) * scale; i ? g.lineTo(px, py) : g.moveTo(px, py); });
    g.stroke();
  }
  // buildings
  g.fillStyle = 'rgba(120,50,40,0.9)';
  for (const b of sc.buildings) if (b.k === 'house' || b.k === 'hut' || b.k === 'keep') g.fillRect((b.x - b.w / 2 + W.half) * scale, (b.z - b.d / 2 + W.half) * scale, Math.max(1.5, b.w * scale), Math.max(1.5, b.d * scale));
  g.fillStyle = 'rgba(190,185,175,0.95)';
  for (const b of sc.buildings) if (b.k === 'wall') g.fillRect((b.x - b.w / 2 + W.half) * scale, (b.z - b.d / 2 + W.half) * scale, Math.max(1.2, b.w * scale), Math.max(1.2, b.d * scale));
  g.fillStyle = '#8a2f22';
  g.fillRect((-5.5 + W.half) * scale, (-4.5 + W.half) * scale, 11 * scale, 9 * scale);
  return cv;
}

export function buildMinimap(scenery) {
  half = scenery.half;
  base = renderWorldCanvas(PPM, scenery);
  S.worldCanvas = renderWorldCanvas(1.2, scenery);
}

export function drawMinimap(camYaw) {
  const cv = $('minimap'), g = cv.getContext('2d');
  const Wd = cv.width, c = Wd / 2, scale = c / VIEW; // canvas px per meter
  const L = S.local;
  g.clearRect(0, 0, Wd, Wd);
  if (S.instance) {
    g.fillStyle = '#16141a'; g.fillRect(0, 0, Wd, Wd);
    g.save(); g.translate(c, c); g.rotate(camYaw - Math.PI); g.scale(scale, scale);
    g.fillStyle = '#4a443e';
    const ox = S.instance.origin[0], oz = S.instance.origin[1];
    for (const r of S.instance.rooms) g.fillRect(r.x + ox - r.w / 2 - L.x, r.z + oz - r.d / 2 - L.z, r.w, r.d);
    g.restore();
  } else if (base) {
    g.save();
    g.translate(c, c);
    g.rotate(camYaw - Math.PI);
    g.scale(scale / PPM, scale / PPM);
    g.drawImage(base, -(L.x + half) * PPM, -(L.z + half) * PPM);
    g.restore();
  }

  g.save();
  g.translate(c, c);
  const phi = camYaw - Math.PI;
  g.rotate(phi);
  for (const e of S.ents.values()) {
    if (e.id === S.id || e.k === 'obj') continue;
    const x = (e.rx - L.x) * scale, y = (e.rz - L.z) * scale;
    if (Math.hypot(x, y) > c + 6) continue;
    g.save();
    g.translate(x, y);
    g.rotate(-phi);
    if (e.k === 'npc') {
      const m = S.marks[e.id];
      g.fillStyle = '#ffe14a';
      g.beginPath(); g.arc(0, 0, 4, 0, 6.3); g.fill();
      if (m) { g.fillStyle = m === 'progress' ? '#aaa' : '#ffd01a'; g.strokeStyle = '#000'; g.lineWidth = 3; g.font = 'bold 20px serif'; g.textAlign = 'center'; g.strokeText(m === 'available' ? '!' : '?', 0, -8); g.fillText(m === 'available' ? '!' : '?', 0, -8); }
    } else if (e.k === 'mob') {
      g.fillStyle = e.d ? (e.lt ? '#ffd24a' : '#666') : e.el ? '#ff8a20' : '#ffd23a';
      g.strokeStyle = '#000'; g.lineWidth = 1.5;
      g.beginPath(); g.arc(0, 0, e.d ? 3 : e.el ? 5.5 : 4, 0, 6.3); g.fill(); g.stroke();
    } else if (e.k === 'pet') {
      g.fillStyle = '#7ad87a'; g.beginPath(); g.arc(0, 0, 3, 0, 6.3); g.fill();
    } else {
      const inParty = S.party.members.some((m) => m.id === e.id);
      g.fillStyle = inParty ? '#6aa8ff' : '#e8e8e8';
      g.beginPath(); g.arc(0, 0, 4.5, 0, 6.3); g.fill(); g.strokeStyle = '#000'; g.lineWidth = 1.5; g.stroke();
    }
    g.restore();
  }
  g.restore();

  g.save();
  g.translate(c, c);
  const hx = Math.sin(L.r), hz = Math.cos(L.r);
  const ax = hx * Math.cos(phi) - hz * Math.sin(phi), az = hx * Math.sin(phi) + hz * Math.cos(phi);
  g.rotate(Math.atan2(ax, -az));
  g.fillStyle = '#fff'; g.strokeStyle = '#000'; g.lineWidth = 2;
  g.beginPath(); g.moveTo(0, -11); g.lineTo(8, 9); g.lineTo(0, 4); g.lineTo(-8, 9); g.closePath(); g.fill(); g.stroke();
  g.restore();

  const n = $('mm-north'), R = 80;
  n.style.left = (88 + Math.sin(phi) * R) + 'px';
  n.style.top = (88 - Math.cos(phi) * R) + 'px';
  n.style.transform = 'translate(-50%,-50%)';
}
