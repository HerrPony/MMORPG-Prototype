// Aggregates all HUD modules behind one small facade used by the controller and main loop.
import { S } from '../state.js';
import { $ } from '../util.js';
import { initHud } from './hud.js';
import { initActionbar, castSlot, updateActionbar } from './actionbar.js';
import { initWindows, toggle, closeTopWindow, updateWindows, hideTip, resetWindows } from './windows.js';
import { buildMinimap, drawMinimap } from './minimap.js';
import { initSocial, updateSocial, playerMenu } from './social.js';
import { initWorldMap, drawWorldMap } from './worldmap.js';

let hud = null;
let controller = null;

export const ui = {
  setController(c) { controller = c; },
  castSlot,
  toggle,
  playerMenu,
  closeTooltip: hideTip,
  openChat: (p) => S.openChat(p),
  escape() {
    if (closeTopWindow()) return;
    if (S.me.tg) { controller.setTarget(0); return; }
    toggle('menu');
  },
};

let initialized = false;
export function initUI(net) {
  if (!initialized) {
    initialized = true;
    hud = initHud(net);
    initWindows(net);
    initActionbar(net, () => controller);
    initSocial(net);
    initWorldMap();
    net.on('marks', (m) => { S.marks = m.marks; S.entities?.refreshAllMarkers(); });
  }
}

export function enterHud(scenery) {
  resetWindows();
  buildMinimap(scenery);
  $('hud').classList.remove('hidden');
  $('zone-name').textContent = 'Erlenwald';
}

export function leaveHud() {
  $('hud').classList.add('hidden');
  $('labels').innerHTML = '';
  $('fct').innerHTML = '';
  $('chat-log').innerHTML = '';
  S.party = { id: 0, leader: 0, members: [], states: [] };
  resetWindows();
}

let mmT = 0, mapT = 0;
export function updateUI(dt, camYaw) {
  hud.update();
  updateActionbar();
  updateWindows();
  updateSocial();
  mmT += dt;
  if (mmT > 0.05) { mmT = 0; drawMinimap(camYaw); }
  if (!$('mapwin').classList.contains('hidden')) { mapT += dt; if (mapT > 0.1) { mapT = 0; drawWorldMap(); } }
}

export function showError(t) { hud?.showError(t); }
