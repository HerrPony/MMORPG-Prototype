import { S } from '../state.js';
import { $, W, contOf, BIOME_KEYS } from '../util.js';
import { getQuestAreas } from './qareas.js';

let centroids = null;

function discovered(id) { return S.disc.includes(id); }

// Centre of every zone per continent (used for the zone labels).
function zoneCentroids() {
  if (centroids) return centroids;
  const acc = {};
  for (let j = 0; j < W.nz; j += 3) {
    for (let i = 0; i < W.nx; i += 3) {
      let b = W.bm[j * W.nx + i];
      if (b === 0) continue;
      if (b === 10) b = 9;
      const x = W.x0 + i * W.cell;
      const a = (acc[b] ||= { x: 0, z: 0, n: 0 });
      a.x += x; a.z += W.z0 + j * W.cell; a.n++;
    }
  }
  centroids = {};
  for (const [b, a] of Object.entries(acc)) centroids[b] = { x: a.x / a.n, z: a.z / a.n };
  return centroids;
}

export function setMapCont(c) { S.mapCont = c; syncTabs(); }
function syncTabs() {
  document.querySelectorAll('#map-tabs button').forEach((b) => b.classList.toggle('on', +b.dataset.c === (S.mapCont ?? 0)));
}

// Draws one continent. Zones you never visited stay hidden under fog.
export function drawWorldMap() {
  const cv = $('mapcv'), g = cv.getContext('2d');
  const size = cv.width;
  const c = S.mapCont ?? contOf(S.local.x);
  const cont = S.scenery.conts[c];
  const k = size / cont.size;
  const X = (x) => (x - cont.x0) * k, Z = (z) => (z - W.z0) * k;
  const sc = S.scenery;
  g.clearRect(0, 0, size, size);
  g.fillStyle = '#10222e';
  g.fillRect(0, 0, size, size);
  if (!S.contCanvas?.[c]) return;
  g.imageSmoothingEnabled = true;
  g.drawImage(S.contCanvas[c], 0, 0, size, size);

  // fog over zones that were never entered
  const FN = 150;
  const fog = document.createElement('canvas');
  fog.width = fog.height = FN;
  const fg = fog.getContext('2d');
  const img = fg.createImageData(FN, FN);
  for (let y = 0; y < FN; y++) {
    for (let x = 0; x < FN; x++) {
      const wx = cont.x0 + (x + 0.5) / FN * cont.size, wz = W.z0 + (y + 0.5) / FN * 600;
      const gx = Math.max(0, Math.min(W.nx - 1, Math.floor((wx - W.x0) / W.cell))), gz = Math.max(0, Math.min(W.nz - 1, Math.floor((wz - W.z0) / W.cell)));
      const b = W.bm[gz * W.nx + gx];
      const known = b === 0 || discovered('zone:' + BIOME_KEYS[b]);
      const o = (y * FN + x) * 4;
      img.data[o] = 14; img.data[o + 1] = 16; img.data[o + 2] = 22; img.data[o + 3] = known ? 0 : 215;
    }
  }
  fg.putImageData(img, 0, 0);
  g.drawImage(fog, 0, 0, size, size);

  // zone labels
  g.textAlign = 'center';
  g.font = 'bold 20px Georgia, serif';
  const cen = zoneCentroids();
  for (const z of sc.zones) {
    if (z.id === 0 || z.id === 10 || z.cont !== c) continue;
    const ce = cen[z.id];
    if (!ce) continue;
    const known = discovered('zone:' + z.key);
    const px = X(ce.x), py = Z(ce.z);
    g.fillStyle = known ? 'rgba(255,240,200,0.92)' : 'rgba(180,180,190,0.55)';
    g.strokeStyle = 'rgba(0,0,0,0.8)'; g.lineWidth = 4;
    const name = known ? z.name : '???';
    g.strokeText(name, px, py); g.fillText(name, px, py);
    g.font = '13px Segoe UI, sans-serif';
    const lv = `Stufe ${z.minLvl}–${z.maxLvl}`;
    g.strokeText(lv, px, py + 16); g.fillText(lv, px, py + 16);
    g.font = 'bold 20px Georgia, serif';
  }

  // quest areas
  const areas = getQuestAreas().filter((a) => contOf(a.x) === c);
  const tt = performance.now() / 1000;
  for (const a of areas) {
    const px = X(a.x), py = Z(a.z), rr = Math.max(8, a.r * k);
    g.fillStyle = 'rgba(255,200,40,0.2)';
    g.strokeStyle = `rgba(255,214,64,${0.7 + Math.sin(tt * 3) * 0.2})`;
    g.lineWidth = 2.5; g.setLineDash([7, 5]);
    g.beginPath(); g.arc(px, py, rr, 0, 6.3); g.fill(); g.stroke();
    g.setLineDash([]);
    g.font = 'bold 18px serif'; g.textAlign = 'center'; g.fillStyle = '#ffd01a'; g.strokeStyle = '#000'; g.lineWidth = 3.5;
    g.strokeText(a.turnIn || a.k === 'npc' ? '?' : '!', px, py + 6); g.fillText(a.turnIn || a.k === 'npc' ? '?' : '!', px, py + 6);
  }
  S._mapAreas = areas.map((a) => ({ ...a, cont: c }));

  const pts = [];
  for (const t of sc.towns) {
    if (t.kind === 'site' || contOf(t.x) !== c) continue;
    const known = discovered(t.id);
    pts.push({ x: t.x, z: t.z, name: known ? t.name : '???', type: t.harbor ? 'harbor' : t.kind, known, id: t.id, fp: S.fp.includes(t.id) });
  }
  for (const p of sc.pois) {
    if (p.town || contOf(p.x) !== c) continue;
    if (!discovered(p.id)) continue;
    pts.push({ x: p.x, z: p.z, name: p.name, type: p.kind, known: true, id: p.id });
  }
  for (const [id, name] of [['rabenfels', 'Schwarzfelsmine (Dungeon)'], ['nebeltempel', 'Nebeltempel (Dungeon)'], ['feuerschlund', 'Feuerschlund (Dungeon)']]) {
    const dg = sc.pois.find((p) => p.id === id);
    if (dg && contOf(dg.x) === c && discovered(id)) pts.push({ x: dg.x, z: dg.z, name, type: 'dungeon', known: true, id: 'dungeon' });
  }
  g.textAlign = 'center';
  for (const p of pts) {
    const px = X(p.x), py = Z(p.z);
    g.fillStyle = p.known ? '#fff' : '#889';
    g.strokeStyle = '#000'; g.lineWidth = 3;
    const big = p.type === 'capital', town = p.type === 'village' || p.type === 'harbor';
    if (big || town) {
      g.fillStyle = p.known ? (big ? '#ffd86a' : '#ffe9a8') : '#778';
      g.beginPath();
      if (big) { for (let i = 0; i < 5; i++) { const a = -Math.PI / 2 + i * 1.2566; g.lineTo(px + Math.cos(a) * 9, py + Math.sin(a) * 9); const a2 = a + 0.628; g.lineTo(px + Math.cos(a2) * 4, py + Math.sin(a2) * 4); } g.closePath(); }
      else g.rect(px - 5, py - 5, 10, 10);
      g.fill(); g.lineWidth = 1.5; g.stroke();
      if (p.fp) { g.font = '15px sans-serif'; g.fillText('🦅', px + 14, py + 5); }
      if (p.type === 'harbor') { g.font = '15px sans-serif'; g.fillText('⚓', px + (px < size / 2 ? 14 : -14), py + 5); }
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
  // the ferry route
  const hb = sc.towns.filter((t) => t.harbor);
  const here = hb.find((t) => contOf(t.x) === c), other = hb.find((t) => contOf(t.x) !== c);
  if (here && other && discovered(here.id)) {
    g.setLineDash([3, 6]); g.strokeStyle = 'rgba(255,255,255,0.65)'; g.lineWidth = 2;
    g.beginPath(); g.moveTo(X(here.x), Z(here.z)); g.lineTo(c === 0 ? size : 0, Z(here.z) + (Z(other.z) - Z(here.z)) * 0.2); g.stroke(); g.setLineDash([]);
    g.font = 'italic 12px Segoe UI, sans-serif'; g.fillStyle = '#cfe6ff'; g.lineWidth = 3; g.strokeStyle = '#000';
    const lx = c === 0 ? size - 70 : 70, ly = Z(here.z) + 28;
    g.strokeText('⛵ Überfahrt', lx, ly); g.fillText('⛵ Überfahrt', lx, ly);
  }
  S._mapPts = pts.map((p) => ({ ...p, cont: c }));

  // party members
  for (const st of S.party.states) {
    if (st.id === S.id || st.i !== 0 || contOf(st.x) !== c) continue;
    const px = X(st.x), py = Z(st.z);
    g.fillStyle = '#6aa8ff'; g.strokeStyle = '#000'; g.lineWidth = 2;
    g.beginPath(); g.arc(px, py, 5, 0, 6.3); g.fill(); g.stroke();
  }
  const L = S.local;
  if (!S.instance && contOf(L.x) === c) {
    const px = X(L.x), py = Z(L.z);
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
    const c = S.mapCont ?? contOf(S.local.x);
    const cont = S.scenery.conts[c];
    const r = cv.getBoundingClientRect();
    const mx = cont.x0 + (e.clientX - r.left) / r.width * cont.size, mz = W.z0 + (e.clientY - r.top) / r.height * 600;
    let best = null, bd = 18, text = '';
    for (const p of S._mapPts || []) { const d = Math.hypot(p.x - mx, p.z - mz); if (d < bd) { bd = d; best = p; } }
    if (best) text = best.name + (best.fp ? ' · Flugpunkt' : best.type === 'harbor' ? ' · Hafen' : '');
    const qa = (S._mapAreas || []).find((a) => Math.hypot(a.x - mx, a.z - mz) <= Math.max(10, a.r));
    if (qa) text += (text ? '\n' : '') + `Quest: ${qa.title} – ${qa.text}`;
    const tip = $('map-tip');
    if (text) {
      tip.textContent = text;
      tip.style.whiteSpace = 'pre';
      tip.style.left = (e.clientX - r.left + 12) + 'px'; tip.style.top = (e.clientY - r.top + 12) + 'px';
      tip.style.display = 'block';
    } else tip.style.display = 'none';
  });
  cv.addEventListener('mouseleave', () => { $('map-tip').style.display = 'none'; });
  document.querySelectorAll('#map-tabs button').forEach((b) => b.addEventListener('click', () => { S.mapCont = +b.dataset.c; syncTabs(); drawWorldMap(); }));
  $('map-legend').innerHTML = '<span>★ Hauptstadt</span><span>■ Dorf</span><span>⚓ Hafen</span><span>🦅 Flugpunkt</span><span>⛺ Lager</span><span>🕳️ Höhle</span><span>💀 Dungeon</span><span style="color:#ffd01a">◎ Questgebiet</span><span style="color:#ff6a6a">▲ Du</span><span style="color:#6aa8ff">● Gruppe</span>';
}
