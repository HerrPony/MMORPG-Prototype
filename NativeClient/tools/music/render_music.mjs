// Renders the web client's procedural score (Client/js/audio.js) offline into OGG files for the native client.
//   node render_music.mjs <outDir> [--sfx-only] [--music-only] [--only name,name]
// Needs: node + playwright (chromium) and ffmpeg with libvorbis. The WebAudio engine is the single source of the
// music, so web and native client sound the same.
import { createRequire } from 'module';
import http from 'http';
import fs from 'fs';
import path from 'path';
import { execFileSync } from 'child_process';
import { fileURLToPath } from 'url';

const require = createRequire(import.meta.url);
let chromium;
try { ({ chromium } = require('playwright')); } catch { ({ chromium } = require('/opt/node-tools/node_modules/playwright')); }

const here = path.dirname(fileURLToPath(import.meta.url));
const clientDir = path.resolve(here, '../../../Client');
const out = path.resolve(process.argv[2] || path.join(here, 'out'));
const sfxOnly = process.argv.includes('--sfx-only');
const musicOnly = process.argv.includes('--music-only');
const onlyArg = process.argv.indexOf('--only');
const only = onlyArg > 0 ? process.argv[onlyArg + 1].split(',') : null;
fs.mkdirSync(out, { recursive: true });

const server = http.createServer((req, res) => {
  const p = new URL(req.url, 'http://x').pathname;
  if (p === '/') { res.setHeader('content-type', 'text/html'); res.end('<html><body></body></html>'); return; }
  const f = path.join(clientDir, p);
  if (!f.startsWith(clientDir) || !fs.existsSync(f)) { res.statusCode = 404; res.end(); return; }
  res.setHeader('content-type', f.endsWith('.js') ? 'text/javascript' : 'application/octet-stream');
  res.end(fs.readFileSync(f));
}).listen(0);
const port = server.address().port;

const browser = await chromium.launch({ args: ['--autoplay-policy=no-user-gesture-required'] });
const page = await browser.newPage();
page.on('console', (m) => { if (m.type() === 'error') console.log('page:', m.text()); });
await page.goto(`http://localhost:${port}/`);
await page.addScriptTag({ type: 'module', content: `
import { Audio, THEMES } from '/js/audio.js';
window.THEMES = THEMES;
function wav(buf) {
  const ch = buf.numberOfChannels, n = buf.length, sr = buf.sampleRate;
  const bytes = new Uint8Array(44 + n * ch * 2), dv = new DataView(bytes.buffer);
  const w = (o, s) => { for (let i = 0; i < s.length; i++) dv.setUint8(o + i, s.charCodeAt(i)); };
  w(0, 'RIFF'); dv.setUint32(4, 36 + n * ch * 2, true); w(8, 'WAVE'); w(12, 'fmt '); dv.setUint32(16, 16, true); dv.setUint16(20, 1, true);
  dv.setUint16(22, ch, true); dv.setUint32(24, sr, true); dv.setUint32(28, sr * ch * 2, true); dv.setUint16(32, ch * 2, true); dv.setUint16(34, 16, true);
  w(36, 'data'); dv.setUint32(40, n * ch * 2, true);
  const data = [...Array(ch).keys()].map((c) => buf.getChannelData(c));
  let o = 44;
  for (let i = 0; i < n; i++) for (let c = 0; c < ch; c++) { const v = Math.max(-1, Math.min(1, data[c][i])); dv.setInt16(o, v < 0 ? v * 32768 : v * 32767, true); o += 2; }
  let s = ''; const CH = 0x8000;
  for (let i = 0; i < bytes.length; i += CH) s += String.fromCharCode.apply(null, bytes.subarray(i, i + CH));
  return btoa(s);
}
window.renderMusic = async ({ theme, mode, zone, night, intensity, bars, phrase }) => {
  const sr = 44100, th = THEMES[theme];
  const barLen = 4 * 60 / th.bpm, dur = bars * barLen + 3.5;
  const ctx = new OfflineAudioContext(2, Math.ceil(dur * sr), sr);
  const a = new Audio();
  a.muted = false; a.vol = 0.9;
  a.init(ctx);
  a.mode = mode; a.zone = zone; a.night = night; a.intensity = intensity; a.barT = 0.05; a.phrase = phrase || 0;
  while (a.bar < bars) a.playBar();
  return wav(await ctx.startRendering());
};
window.renderSfx = async (name, dur) => {
  const sr = 44100;
  const ctx = new OfflineAudioContext(2, Math.ceil(dur * sr), sr);
  Object.defineProperty(ctx, 'state', { value: 'running' });
  const a = new Audio();
  a.muted = false; a.vol = 1;
  a.init(ctx);
  if (name.startsWith('x_')) extra(a, ctx, name);
  else a.sfx(name);
  return wav(await ctx.startRendering());
};
// effects that exist only in the native client
function extra(a, c, name) {
  switch (name) {
    case 'x_step': a.noise(0.09, 700, 0.35, 'lowpass'); a.tone(90, 0.07, 0.18, 'sine', 55); break;
    case 'x_step_snow': a.noise(0.14, 1800, 0.22, 'lowpass'); break;
    case 'x_step_sand': a.noise(0.12, 2500, 0.22, 'bandpass'); break;
    case 'x_step_stone': a.noise(0.05, 3000, 0.3, 'highpass'); a.tone(220, 0.05, 0.12, 'square', 120); break;
    case 'x_jump': a.tone(220, 0.18, 0.2, 'sine', 420); a.noise(0.1, 1200, 0.15); break;
    case 'x_land': a.noise(0.12, 500, 0.4, 'lowpass'); a.tone(80, 0.12, 0.3, 'sine', 40); break;
    case 'x_splash': a.noise(0.5, 1500, 0.4, 'bandpass'); a.noise(0.35, 700, 0.3, 'lowpass'); break;
    case 'x_ui_open': a.tone(520, 0.08, 0.12, 'triangle', 780); break;
    case 'x_ui_close': a.tone(640, 0.08, 0.12, 'triangle', 420); break;
    case 'x_ui_click': a.tone(900, 0.04, 0.1, 'square'); break;
    case 'x_error': a.tone(180, 0.18, 0.2, 'sawtooth', 120); break;
    case 'x_cast': a.tone(300, 0.6, 0.12, 'sine', 700); a.tone(450, 0.6, 0.08, 'triangle', 900); break;
    case 'x_equip': a.noise(0.08, 2500, 0.3, 'highpass'); a.tone(300, 0.1, 0.15, 'square', 200); break;
    case 'x_sell': a.tone(1976, 0.08, 0.15, 'square'); a.tone(2637, 0.1, 0.15, 'square', 0, 0.06); a.tone(3136, 0.18, 0.12, 'square', 0, 0.12); break;
    case 'x_death': a.tone(200, 0.9, 0.3, 'sawtooth', 40); a.noise(0.6, 400, 0.3, 'lowpass'); break;
    case 'x_wolf': a.tone(300, 0.7, 0.2, 'sawtooth', 500); a.tone(500, 0.6, 0.12, 'sine', 380, 0.2); break;
    case 'x_boss': a.tone(60, 1.6, 0.4, 'sawtooth', 35); a.noise(1.2, 300, 0.4, 'lowpass'); break;
    case 'x_fly': a.noise(1.4, 900, 0.25, 'bandpass'); a.tone(120, 1.4, 0.1, 'sine', 160); break;
  }
}
window.ready = true;
` });
await page.waitForFunction('window.ready === true');

const encode = (b64, file, q) => {
  const wavPath = file.replace(/\.ogg$/, '.wav');
  fs.writeFileSync(wavPath, Buffer.from(b64, 'base64'));
  execFileSync('ffmpeg', ['-y', '-loglevel', 'error', '-i', wavPath, '-c:a', 'libvorbis', '-q:a', String(q), file]);
  fs.unlinkSync(wavPath);
  return fs.statSync(file).size;
};

const themes = await page.evaluate('Object.keys(window.THEMES)');
const wanted = (name) => !only || only.includes(name);
if (!sfxOnly) {
  fs.mkdirSync(path.join(out, 'music'), { recursive: true });
  const jobs = [];
  for (const t of themes) {
    if (t === 'menu') jobs.push({ name: 'menu', theme: 'menu', mode: 'menu', zone: 'forest', night: 0, intensity: 0, bars: 24 });
    else if (t === 'dungeon') {
      jobs.push({ name: 'dungeon', theme: 'dungeon', mode: 'dungeon', zone: 'forest', night: 0, intensity: 0, bars: 24 });
      jobs.push({ name: 'dungeon_fight', theme: 'dungeon', mode: 'dungeon', zone: 'forest', night: 0, intensity: 1, bars: 16, phrase: 1 });
    } else {
      jobs.push({ name: `${t}_day`, theme: t, mode: 'world', zone: t, night: 0, intensity: 0, bars: 24 });
      jobs.push({ name: `${t}_night`, theme: t, mode: 'world', zone: t, night: 1, intensity: 0, bars: 16, phrase: 2 });
      jobs.push({ name: `${t}_fight`, theme: t, mode: 'world', zone: t, night: 0, intensity: 1, bars: 16, phrase: 1 });
    }
  }
  for (const j of jobs) {
    if (!wanted(j.name)) continue;
    const t0 = Date.now();
    const b64 = await page.evaluate((o) => window.renderMusic(o), j);
    const size = encode(b64, path.join(out, 'music', j.name + '.ogg'), 1.5);
    console.log(`music ${j.name}: ${(size / 1e6).toFixed(2)} MB in ${((Date.now() - t0) / 1000).toFixed(1)} s`);
  }
}
if (!musicOnly) {
  fs.mkdirSync(path.join(out, 'sfx'), { recursive: true });
  const sfx = { hit: 0.4, hurt: 0.5, melee: 0.4, arrow: 0.5, bolt: 0.5, fire: 0.8, frost: 0.7, nova: 0.9, buff: 0.8, loot: 0.5, quest: 0.9, levelup: 1.8, learn: 0.9, coin: 0.4, click: 0.2,
    x_step: 0.2, x_step_snow: 0.25, x_step_sand: 0.25, x_step_stone: 0.15, x_jump: 0.4, x_land: 0.3, x_splash: 0.8, x_ui_open: 0.25, x_ui_close: 0.25, x_ui_click: 0.15, x_error: 0.4,
    x_cast: 1.0, x_equip: 0.3, x_sell: 0.5, x_death: 1.5, x_wolf: 1.2, x_boss: 2.2, x_fly: 1.8 };
  for (const [name, dur] of Object.entries(sfx)) {
    if (!wanted(name)) continue;
    const b64 = await page.evaluate(([n, d]) => window.renderSfx(n, d), [name, dur]);
    const size = encode(b64, path.join(out, 'sfx', name.replace(/^x_/, '') + '.ogg'), 3);
    console.log(`sfx ${name}: ${(size / 1000).toFixed(0)} KB`);
  }
}
await browser.close();
server.close();
