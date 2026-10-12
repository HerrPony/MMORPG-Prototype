import * as THREE from 'three';
import { Net } from './net.js';
import { audio } from './audio.js';
import { S } from './state.js';
import { $, W } from './util.js';
import { Preview, initScreens, showScreen } from './screens.js';
import { buildWorld } from './world/scene.js';
import { Entities } from './world/entities.js';
import { FX } from './world/fx.js';
import { loadHumanAssets } from './world/models.js';
import { Controller } from './player.js';
import { ui, initUI, enterHud, leaveHud, updateUI } from './ui/index.js';
import { renderTracker } from './ui/hud.js';
import { initTouch } from './touch.js';

const canvas = $('gl');
const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
renderer.setPixelRatio(Math.min(devicePixelRatio, 2));
const camera = new THREE.PerspectiveCamera(60, 1, 0.1, 600);
const preview = new Preview();
const net = new Net();
S.net = net;

let game = null;
let spawned = false;
let activeScene = null;

function resize() {
  renderer.setSize(innerWidth, innerHeight, false);
  camera.aspect = innerWidth / innerHeight;
  camera.updateProjectionMatrix();
}
addEventListener('resize', resize);
resize();

function ensureGame(scenery) {
  if (game) return;
  const world = buildWorld(scenery, renderer);
  const entities = new Entities(world.scene);
  const fx = new FX(world.scene, entities, () => game);
  const controller = new Controller({ canvas, camera, world, entities, ui });
  ui.setController(controller);
  game = { world, entities, fx, controller, camera };
  S.entities = entities;
  activeScene = world.scene;
  S.onZone = (m) => {
    if (S.instance || m.instance) return;
    const key = game.world.biomeKey(S.local.x, S.local.z);
    game.world.setZone(key);
  };
}

// switches between the open world and a dungeon scene
function setActiveScene(scene, dungeon) {
  activeScene = scene;
  game.entities.setScene(scene);
  game.fx.setScene(scene);
  game.controller.setDungeon(dungeon);
}

function leaveWorld() {
  S.inWorld = false;
  audio.setMode('menu');
  if (game) { game.entities.removeAll(); S.instance = null; if (game.world.dungeon) { game.world.leaveDungeon(); setActiveScene(game.world.scene, null); } }
  leaveHud();
  S.me.tg = 0;
}

// ---- network wiring ----
initUI(net);
initScreens(net, preview, { isInWorld: () => S.inWorld, leaveWorld });

net.on('world', (m) => {
  S.id = m.id; S.self = m.self; S.cls = m.cls; S.skillIds = m.skills; S.defs = m.defs; S.scenery = m.scenery;
  S.disc = m.disc || []; S.fp = m.fp || [];
  S.cd = {}; S.gcd = { end: 0, dur: 0 }; S.cast = null; S.quests = []; S.marks = {}; S.instance = null;
  S.me = { hp: 1, mh: 1, res: 0, rm: 1, xp: 0, xn: 100, lvl: 1, tg: 0, auto: false, d: false };
  ensureGame(m.scenery);
  if (game.world.dungeon) { game.world.leaveDungeon(); setActiveScene(game.world.scene, null); }
  game.entities.removeAll();
  spawned = false;
  preview.clear();
  showScreen('none');
  enterHud(m.scenery);
  renderTracker();
  S.inWorld = true;
  audio.setMode('world');
  S.addChat('sys', 'Tippe /hilfe für die Steuerung.');
});

net.on('s', (m) => {
  if (!S.inWorld) return;
  game.entities.apply(m.e, m.lt);
  const tg = S.me.tg;
  Object.assign(S.me, m.me);
  S.me.tg = tg;
  S.party.states = m.pt || [];
  const own = S.ents.get(S.id);
  if (own && !spawned) { spawned = true; game.controller.spawn(own.x, own.y, own.z, own.r); game.world.setZone(game.world.biomeKey(own.x, own.z)); }
});
net.on('correct', (m) => game?.controller.correct(m.x, m.y, m.z));

// dungeon instances
net.on('instance', (m) => {
  if (!game) return;
  if (m.on) {
    S.instance = { name: m.name, rooms: m.rooms, origin: m.origin };
    const d = game.world.enterDungeon({ rooms: m.rooms, origin: m.origin, style: m.style });
    setActiveScene(d.scene, d);
    audio.setMode('dungeon');
  } else {
    S.instance = null;
    game.world.leaveDungeon();
    setActiveScene(game.world.scene, null);
    audio.setMode('world');
  }
  game.entities.removeAll();
});

net.on('dmg', (m) => {
  game?.fx.onDmg(m);
  if (m.id === S.id) { lastCombat = performance.now() / 1000; audio.sfx('hurt'); }
  else if (m.s === S.id) { lastCombat = performance.now() / 1000; audio.sfx('hit'); }
});
net.on('fx', (m) => {
  game?.fx.onFx(m);
  if (m.s !== S.id) return;
  if (m.k === 'levelup') { audio.sfx('levelup'); return; }
  const st = S.defs?.skills?.[m.k.replace(':cast', '')]?.style;
  if (st && !m.k.endsWith(':cast')) audio.sfx(st === 'melee' ? 'melee' : st);
});
net.on('questDone', () => audio.sfx('quest'));
net.on('skills', () => audio.sfx('learn'));
net.on('loot', () => audio.sfx('loot'));
// admin mode: accounts on the server's admin list get the map editor button; edits reload everyone's client
net.on('auth', (m) => { S.admin = !!m.admin; $('btn-editor').classList.toggle('hidden', !S.admin); });
$('btn-editor').onclick = () => {
  const q = S.inWorld ? `?x=${S.local.x.toFixed(0)}&z=${S.local.z.toFixed(0)}` : '';
  window.open('editor.html' + q, '_blank');
};
net.on('worldChanged', () => { alert('Ein Admin hat die Welt verändert. Das Spiel wird neu geladen.'); location.reload(); });
net.on('kicked', (m) => { alert(m.text); location.reload(); });
net.on('_close', () => { alert('Verbindung zum Server getrennt.'); location.reload(); });

initTouch(() => game?.controller, ui);

// ---- audio + fullscreen ----
const unlock = () => audio.start();
addEventListener('pointerdown', unlock);
addEventListener('keydown', unlock);
function toggleFullscreen() {
  if (document.fullscreenElement) document.exitFullscreen?.();
  else document.documentElement.requestFullscreen?.().catch(() => {});
}
function refreshMuteBtn() { $('btn-mute').textContent = audio.muted ? '🔇' : '🔊'; }
$('btn-fs').onclick = toggleFullscreen;
$('m-fs').onclick = toggleFullscreen;
$('btn-mute').onclick = () => { audio.start(); audio.toggleMute(); refreshMuteBtn(); };
$('vol').value = audio.vol;
$('vol').oninput = (e) => { audio.setVolume(+e.target.value); if (audio.muted) { audio.toggleMute(); refreshMuteBtn(); } };
refreshMuteBtn();
addEventListener('keydown', (e) => {
  if (e.key === 'F11') { e.preventDefault(); toggleFullscreen(); return; }
  if (S.typing || e.target.tagName === 'INPUT') return;
  if (e.key.toLowerCase() === 'n' && S.inWorld) { audio.toggleMute(); refreshMuteBtn(); }
});
let lastCombat = -99, intensity = 0;

// ---- main loop ----
let last = performance.now();
function frame(t) {
  requestAnimationFrame(frame);
  step(t);
}
// keep simulating when the tab is in the background (rAF is paused there)
setInterval(() => { const t = performance.now(); if (t - last > 150) step(t); }, 50);
function step(t) {
  const dt = Math.min(0.05, (t - last) / 1000);
  last = t;
  const target = t / 1000 - lastCombat < 6 && !S.me.d ? 1 : 0;
  intensity += (target - intensity) * Math.min(1, dt * 0.8);
  audio.intensity = S.inWorld ? intensity : 0;
  if (S.inWorld && game && spawned) {
    const g = game;
    g.controller.update(dt);
    g.entities.update(dt, camera, innerWidth, innerHeight);
    g.fx.update(dt);
    if (!S.instance) g.world.update(dt);
    // keep the sky matched to the biome we are standing in
    if (!S.instance && ((t / 1000) | 0) % 2 === 0) g.world.setZone(g.world.biomeKey(S.local.x, S.local.z));
    updateUI(dt, g.controller.camYaw);
    renderer.render(activeScene, camera);
  } else {
    preview.update(dt, innerWidth, innerHeight);
    renderer.render(preview.scene, preview.camera);
  }
}
requestAnimationFrame(frame);

// drag to rotate the preview model
addEventListener('mousemove', (e) => { if (preview.drag) preview.yaw += e.movementX * 0.01; });
addEventListener('mouseup', () => { preview.drag = false; });
addEventListener('mousedown', (e) => { if (!S.inWorld && e.target === canvas) preview.drag = true; });

// ---- connect (after the character rig has been loaded) ----
$('login-status').textContent = 'Lade Spielwelt…';
loadHumanAssets().then(() => net.connect()).then(
  () => { $('login-status').textContent = 'Verbunden. Bitte anmelden oder einen Account erstellen.'; },
  () => { $('login-status').textContent = 'Server nicht erreichbar. Läuft der Go-Server?'; },
);
