import { S } from '../state.js';
import { $, W, esc } from '../util.js';
import { BIOME_RGB } from './minimap.js';

const ZONE_BY_ID = ['sea', 'forest', 'plains', 'swamp', 'snow', 'desert'];
let hover = null;

function discovered(id) { return S.disc.includes(id); }

// Draws the whole continent. Zones you never visited stay hidden under fog.
export function drawWorldMap() {
  const cv = $('mapcv'), g = cv.getContext('2d');
  const size = cv.width;
  const k = size / (W.half * 2);
  const sc = S.scenery;
  g.clearRect(0, 0, size, size);
  g.fillStyle = '#10222e';
  g.fillRect(0, 0, size, size);
  if (!S.worldCanvas) return;
  g.imageSmoothingEnabled = true;
  g.drawImage(S.worldCanvas, 0, 0, size, size);

  // fog over zones that were never entered
  const fog = document.createElement('canvas');
  fog.width = fog.height = 150;
  const fg = fog.getContext('2d');
  const img = fg.createImageData(150, 150);
  for (let y = 0; y < 150; y++) {
    for (let x = 0; x < 150; x++) {
      const gx = Math.min(W.n - 1, Math.floor((x + 0.5) / 150 * W.n)), gz = Math.min(W.n - 1, Math.floor((y + 0.5) / 150 * W.n));
      const b = W.bm[gz * W.n + gx];
      const known = b === 0 || discovered('zone:' + ZONE_BY_ID[b]);
      const o = (y * 150 + x) * 4;
      img.data[o] = 14; img.data[o + 1] = 16; img.data[o + 2] = 22; img.data[o + 3] = known ? 0 : 215;
    }
  }
  fg.putImageData(img, 0, 0);
  g.imageSmoothingEnabled = true;
  g.drawImage(fog, 0, 0, size, size);

  // rivers and lakes already visible via heights; zone labels
  g.textAlign = 'center';
  g.font = 'bold 20px Georgia, serif';
  for (const z of sc.zones) {
    if (z.id === 0) continue;
    const cx = { 1: -40, 2: 170, 3: -190, 4: -60, 5: 180 }[z.id], cz = { 1: 75, 2: -225, 3: 120, 4: -180, 5: 235 }[z.id];
    const known = discovered('zone:' + z.key);
    const px = (cx + W.half) * k, py = (cz + W.half) * k;
    g.fillStyle = known ? 'rgba(255,240,200,0.92)' : 'rgba(180,180,190,0.55)';
    g.strokeStyle = 'rgba(0,0,0,0.8)'; g.lineWidth = 4;
    const name = known ? z.name : '???';
    g.strokeText(name, px, py); g.fillText(name, px, py);
    g.font = '13px Segoe UI, sans-serif';
    const lv = `Stufe ${z.minLvl}–${z.maxLvl}`;
    g.strokeText(lv, px, py + 16); g.fillText(lv, px, py + 16);
    g.font = 'bold 20px Georgia, serif';
  }

  const pts = [];
  // settlements
  for (const t of sc.towns) {
    if (t.kind === 'site') continue;
    const known = discovered(t.id);
    pts.push({ x: t.x, z: t.z, name: known ? t.name : '???', type: t.kind, known, id: t.id, fp: S.fp.includes(t.id) });
  }
  for (const p of sc.pois) {
    if (p.town) continue;
    if (!discovered(p.id)) continue;
    pts.push({ x: p.x, z: p.z, name: p.name, type: p.kind, known: true, id: p.id });
  }
  const dg = sc.pois.find((p) => p.id === 'rabenfels');
  if (dg && discovered('rabenfels')) pts.push({ x: dg.x, z: dg.z, name: 'Schwarzfelsmine (Dungeon)', type: 'dungeon', known: true, id: 'dungeon' });
  g.textAlign = 'center';
  for (const p of pts) {
    const px = (p.x + W.half) * k, py = (p.z + W.half) * k;
    g.fillStyle = p.known ? '#fff' : '#889';
    g.strokeStyle = '#000'; g.lineWidth = 3;
    const big = p.type === 'capital', town = p.type === 'village';
    if (big || town) {
      g.fillStyle = p.known ? (big ? '#ffd86a' : '#ffe9a8') : '#778';
      g.beginPath();
      if (big) { for (let i = 0; i < 5; i++) { const a = -Math.PI / 2 + i * 1.2566; g.lineTo(px + Math.cos(a) * 9, py + Math.sin(a) * 9); const a2 = a + 0.628; g.lineTo(px + Math.cos(a2) * 4, py + Math.sin(a2) * 4); } g.closePath(); }
      else g.rect(px - 5, py - 5, 10, 10);
      g.fill(); g.lineWidth = 1.5; g.stroke();
      if (p.fp) { g.font = '15px sans-serif'; g.fillText('🦅', px + 14, py + 5); }
      g.font = (big ? 'bold 15px' : 'bold 13px') + ' Georgia, serif';
      g.lineWidth = 3.5; g.fillStyle = p.known ? '#fff' : '#99a';
      g.strokeText(p.name, px, py - 12); g.fillText(p.name, px, py - 12);
    } else {
      g.font = '15px sans-serif';
      g.fillStyle = '#fff';
      const icon = { camp: '⛺', cave: '🕳️', lake: '💧', ruin: '🏛️', wild: '📍', dungeon: '💀' }[p.type] || '•';
      g.fillText(icon, px, py + 5);
      g.font = '11px Segoe UI, sans-serif';
      g.lineWidth = 3; g.fillStyle = '#eee';
      g.strokeText(p.name, px, py + 20); g.fillText(p.name, px, py + 20);
    }
  }
  S._mapPts = pts;

  // party members
  for (const st of S.party.states) {
    if (st.id === S.id || st.i !== 0) continue;
    const px = (st.x + W.half) * k, py = (st.z + W.half) * k;
    g.fillStyle = '#6aa8ff'; g.strokeStyle = '#000'; g.lineWidth = 2;
    g.beginPath(); g.arc(px, py, 5, 0, 6.3); g.fill(); g.stroke();
  }
  // quest givers with open quests nearby are shown on the minimap, here only the player
  const L = S.local;
  const px = (L.x + W.half) * k, py = (L.z + W.half) * k;
  if (!S.instance) {
    g.save(); g.translate(px, py); g.rotate(Math.PI - L.r);
    g.fillStyle = '#ff4a4a'; g.strokeStyle = '#fff'; g.lineWidth = 2;
    g.beginPath(); g.moveTo(0, -10); g.lineTo(7, 8); g.lineTo(0, 4); g.lineTo(-7, 8); g.closePath(); g.fill(); g.stroke();
    g.restore();
    const ring = (performance.now() / 600) % 1;
    g.strokeStyle = `rgba(255,80,80,${1 - ring})`; g.lineWidth = 2;
    g.beginPath(); g.arc(px, py, 8 + ring * 14, 0, 6.3); g.stroke();
  }
  g.strokeStyle = '#6d5a30'; g.lineWidth = 4; g.strokeRect(0, 0, size, size);
}

export function initWorldMap() {
  const cv = $('mapcv');
  cv.addEventListener('mousemove', (e) => {
    const r = cv.getBoundingClientRect();
    const mx = (e.clientX - r.left) / r.width * W.half * 2 - W.half, mz = (e.clientY - r.top) / r.height * W.half * 2 - W.half;
    let best = null, bd = 18;
    for (const p of S._mapPts || []) { const d = Math.hypot(p.x - mx, p.z - mz); if (d < bd) { bd = d; best = p; } }
    const tip = $('map-tip');
    if (best) {
      tip.textContent = best.name + (best.fp ? ' · Flugpunkt' : '');
      tip.style.left = (e.clientX - r.left + 12) + 'px'; tip.style.top = (e.clientY - r.top + 12) + 'px';
      tip.style.display = 'block';
    } else tip.style.display = 'none';
  });
  cv.addEventListener('mouseleave', () => { $('map-tip').style.display = 'none'; });
  $('map-legend').innerHTML = '<span>★ Hauptstadt</span><span>■ Dorf</span><span>🦅 Flugpunkt</span><span>⛺ Lager</span><span>🕳️ Höhle</span><span>💀 Dungeon</span><span style="color:#ff6a6a">▲ Du</span><span style="color:#6aa8ff">● Gruppe</span>';
}
