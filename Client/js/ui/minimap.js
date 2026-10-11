import { S } from '../state.js';
import { $, W, distToPoly } from '../util.js';
import { getQuestAreas } from './qareas.js';

const PPM = 2;          // base map pixels per meter
const VIEW = 52;        // visible radius in meters
let base = null;

// index = biome id (sea, forest, plains, swamp, snow, desert, jungle, blossom, mesa, volcano, lava)
export const BIOME_RGB = [[40, 90, 140], [62, 122, 58], [170, 182, 84], [86, 106, 60], [226, 234, 244], [212, 182, 106],
  [22, 140, 72], [238, 170, 200], [200, 96, 56], [78, 62, 70], [255, 96, 30]];

// Renders a part of the world (x0,z0 .. x0+w,z0+h in meters) into an offscreen canvas (also used by the world map).
export function renderWorldCanvas(scale, sc, x0, z0, w, h) {
  const cw = Math.round(w * scale), ch = Math.round(h * scale);
  const cv = document.createElement('canvas');
  cv.width = cw; cv.height = ch;
  const g = cv.getContext('2d');
  const img = g.createImageData(cw, ch);
  for (let py = 0; py < ch; py++) {
    for (let px = 0; px < cw; px++) {
      const gx = Math.max(0, Math.min(W.nx - 1, Math.floor((x0 + px / scale - W.x0) / W.cell))), gz = Math.max(0, Math.min(W.nz - 1, Math.floor((z0 + py / scale - W.z0) / W.cell)));
      const idx = gz * W.nx + gx;
      const h0 = W.hm[idx];
      const b = W.bm[idx];
      const col = BIOME_RGB[b];
      const shade = 0.82 + Math.max(-0.2, Math.min(0.3, h0 * 0.012));
      let r = col[0] * shade, gg = col[1] * shade, bl = col[2] * shade;
      if (b === 10) { r = col[0]; gg = col[1]; bl = col[2]; }
      else if (h0 < 0 && b !== 0) { r = 50; gg = 110; bl = 160; } // rivers and lakes
      const o = (py * cw + px) * 4;
      img.data[o] = r; img.data[o + 1] = gg; img.data[o + 2] = bl; img.data[o + 3] = 255;
    }
  }
  g.putImageData(img, 0, 0);
  const X = (x) => (x - x0) * scale, Z = (z) => (z - z0) * scale;
  g.lineCap = 'round'; g.lineJoin = 'round';
  g.strokeStyle = 'rgba(150,115,70,0.95)'; g.lineWidth = Math.max(1.5, 2.2 * scale);
  for (const rd of sc.roads) {
    g.beginPath();
    rd.forEach(([x, z], i) => { i ? g.lineTo(X(x), Z(z)) : g.moveTo(X(x), Z(z)); });
    g.stroke();
  }
  // buildings
  g.fillStyle = 'rgba(120,50,40,0.9)';
  for (const b of sc.buildings) if (b.k === 'house' || b.k === 'hut' || b.k === 'keep' || b.k === 'pagoda') g.fillRect(X(b.x - b.w / 2), Z(b.z - b.d / 2), Math.max(1.5, b.w * scale), Math.max(1.5, b.d * scale));
  g.fillStyle = 'rgba(190,185,175,0.95)';
  for (const b of sc.buildings) if (b.k === 'wall') g.fillRect(X(b.x - b.w / 2), Z(b.z - b.d / 2), Math.max(1.2, b.w * scale), Math.max(1.2, b.d * scale));
  g.fillStyle = 'rgba(138,106,64,0.95)';
  for (const b of sc.buildings) if (b.k === 'dock') g.fillRect(X(b.x - b.w / 2), Z(b.z - b.d / 2), Math.max(1.5, b.w * scale), Math.max(1.5, b.d * scale));
  g.fillStyle = '#8a2f22';
  g.fillRect(X(-5.5), Z(-4.5), 11 * scale, 9 * scale);
  return cv;
}

export function buildMinimap(scenery) {
  const wpx = (W.nx - 1) * W.cell, hpx = (W.nz - 1) * W.cell;
  base = renderWorldCanvas(PPM, scenery, W.x0, W.z0, wpx, hpx);
  // one map image per continent for the world map
  S.contCanvas = (scenery.conts || []).map((c) => renderWorldCanvas(1.2, scenery, c.x0, W.z0, c.size, 600));
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
    g.drawImage(base, -(L.x - W.x0) * PPM, -(L.z - W.z0) * PPM);
    g.restore();
  }

  g.save();
  g.translate(c, c);
  const phi = camYaw - Math.PI;
  g.rotate(phi);

  // quest areas: a golden circle where the objective can be done, an arrow at the edge when it is out of view
  if (!S.instance) {
    const t = performance.now() / 1000;
    for (const a of getQuestAreas()) {
      const x = (a.x - L.x) * scale, y = (a.z - L.z) * scale, rr = a.r * scale;
      const d = Math.hypot(x, y);
      if (d < c + rr) {
        g.fillStyle = a.turnIn ? 'rgba(255,230,120,0.16)' : 'rgba(255,200,40,0.15)';
        g.strokeStyle = 'rgba(255,214,64,0.85)';
        g.lineWidth = 2;
        g.setLineDash([6, 5]);
        g.beginPath(); g.arc(x, y, rr, 0, 6.3); g.fill(); g.stroke();
        g.setLineDash([]);
        if (a.k === 'npc') { g.save(); g.translate(x, y); g.rotate(-phi); g.fillStyle = '#ffd01a'; g.strokeStyle = '#000'; g.lineWidth = 3; g.font = 'bold 18px serif'; g.textAlign = 'center'; g.strokeText('?', 0, 6); g.fillText('?', 0, 6); g.restore(); }
      }
      if (d > c - 6) {
        const ang = Math.atan2(y, x), rad = c - 9;
        g.save(); g.translate(Math.cos(ang) * rad, Math.sin(ang) * rad); g.rotate(ang + Math.PI / 2);
        g.globalAlpha = 0.65 + Math.sin(t * 4) * 0.3;
        g.fillStyle = '#ffd01a'; g.strokeStyle = '#000'; g.lineWidth = 1.5;
        g.beginPath(); g.moveTo(0, -7); g.lineTo(6, 5); g.lineTo(-6, 5); g.closePath(); g.fill(); g.stroke();
        g.restore();
      }
    }
  }

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
    } else if (e.k === 'ship') {
      g.font = '17px sans-serif'; g.textAlign = 'center'; g.fillText('⛵', 0, 6);
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
