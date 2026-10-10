// Items: tooltips, bags, paperdoll + equipment manager, vendor window.
import { S } from '../state.js';
import { $, esc, fmtCopper, CLASS_COLORS, toast } from '../util.js';
import { showTip, hideTip } from './windows.js';

export const SLOT_LABEL = {
  head: 'Kopf', neck: 'Hals', shoulder: 'Schultern', back: 'Rücken', chest: 'Brust', wrist: 'Handgelenke', hands: 'Hände', waist: 'Taille',
  legs: 'Beine', feet: 'Füße', finger: 'Finger', finger1: 'Finger', finger2: 'Finger', mainhand: 'Waffenhand', offhand: 'Nebenhand',
};
const SLOT_ICON = { head: '🪖', neck: '📿', shoulder: '🦺', back: '🧥', chest: '👕', wrist: '⌚', hands: '🧤', waist: '🎗️', legs: '👖', feet: '🥾', finger1: '💍', finger2: '💍', mainhand: '⚔️', offhand: '🛡️' };
const STAT_LABEL = { str: 'Stärke', agi: 'Beweglichkeit', int: 'Intelligenz', sta: 'Ausdauer' };
const RARITY_LABEL = { poor: 'Schlecht', common: 'Gewöhnlich', uncommon: 'Ungewöhnlich', rare: 'Selten', epic: 'Episch' };
export const RARITY_COLOR = { poor: '#9d9d9d', common: '#ffffff', uncommon: '#1eff00', rare: '#0070dd', epic: '#a335ee' };
const PRIMARY = { warrior: 'str', hunter: 'agi', mage: 'int' };

const def = (id) => S.defs.items[id];

export function itemTip(id, count, mode) {
  const it = def(id);
  if (!it) return '';
  let h = `<div class="tt r-${it.rarity}">${esc(it.name)}</div>`;
  if (it.slot) h += `<div class="tm"><span>${SLOT_LABEL[it.slot]}</span><span>${esc(it.type || '')}</span></div>`;
  else h += `<div class="tm">${RARITY_LABEL[it.rarity]}</div>`;
  if (it.armor) h += `<div>${it.armor} Rüstung</div>`;
  if (it.weapon) h += `<div>${it.weapon} Waffenschaden</div>`;
  for (const k of ['str', 'agi', 'int', 'sta']) if (it.stats?.[k]) h += `<div class="stat">+${it.stats[k]} ${STAT_LABEL[k]}</div>`;
  if (it.classes?.length) {
    const ok = it.classes.includes(S.self.class);
    h += `<div class="${ok ? '' : 'bad'}">Klasse: ${it.classes.map((c) => S.classes[c].name).join(', ')}</div>`;
  }
  if (it.req) h += `<div class="${S.me.lvl >= it.req ? '' : 'bad'}">Benötigt Stufe ${it.req}</div>`;
  if (it.desc) h += `<div class="desc">${esc(it.desc)}</div>`;
  if (mode === 'buy') h += `<div style="margin-top:4px">Preis: ${fmtCopper(it.buy)}</div>`;
  else if (it.sell) h += `<div style="margin-top:4px">Verkaufspreis: ${fmtCopper(it.sell * (count || 1))}</div>`;
  if (mode === 'bag') {
    if (S.vendorNpc) h += '<div class="hint">Rechtsklick: verkaufen · Doppelklick: anlegen/benutzen</div>';
    else if (it.slot) h += '<div class="hint">Rechtsklick: anlegen</div>';
    else if (it.heal) h += '<div class="hint">Rechtsklick: benutzen</div>';
  }
  return h;
}

function usable(it) {
  return (!it.classes?.length || it.classes.includes(S.self.class));
}

// ---------- bags ----------
export function useBagItem(idx, noSell = false) {
  const st = S.inv[idx];
  if (!st) return;
  const it = def(st.i);
  hideTip();
  if (S.vendorNpc && !noSell) S.net.send('sell', { npc: S.vendorNpc, idx });
  else if (it.slot) S.net.send('equip', { idx });
  else if (it.heal) S.net.send('use', { idx });
}

export function renderBags() {
  const g = $('bag-grid');
  g.innerHTML = '';
  for (let i = 0; i < S.invMax; i++) {
    const st = S.inv[i];
    const s = document.createElement('div');
    s.className = 'bslot';
    if (st) {
      const d = def(st.i);
      s.innerHTML = `${d.icon}${st.c > 1 ? `<span class="cnt">${st.c}</span>` : ''}`;
      s.style.borderColor = RARITY_COLOR[d.rarity];
      if (d.slot && d.req > S.me.lvl) s.classList.add('locked');
      else if (d.slot && !usable(d)) s.classList.add('locked');
      s.onmouseenter = (e) => showTip(itemTip(st.i, st.c, 'bag'), e.clientX + 14, e.clientY + 8);
      s.onmouseleave = hideTip;
      s.oncontextmenu = (e) => { e.preventDefault(); useBagItem(i); };
      s.ondblclick = () => useBagItem(i, true);
      s.draggable = true;
      s.ondragstart = (e) => { e.dataTransfer.setData('text/plain', 'bag:' + i); hideTip(); };
    }
    g.appendChild(s);
  }
  g.ondragover = (e) => e.preventDefault();
  g.ondrop = (e) => {
    const d = e.dataTransfer.getData('text/plain');
    if (d.startsWith('eq:')) S.net.send('unequip', { slot: d.slice(3) });
  };
  $('bag-money').innerHTML = fmtCopper(S.copper);
}

// ---------- character window ----------
const LEFT = ['head', 'neck', 'shoulder', 'back', 'chest', 'wrist'];
const RIGHT = ['hands', 'waist', 'legs', 'feet', 'finger1', 'finger2'];
const WEAP = ['mainhand', 'offhand'];

function slotEl(slot) {
  const id = S.equip[slot];
  const d = id && def(id);
  const el = document.createElement('div');
  el.className = 'eslot' + (d ? ' filled' : '');
  el.innerHTML = d ? d.icon : `<span class="ph">${SLOT_ICON[slot]}</span>`;
  el.title = '';
  if (d) {
    el.style.borderColor = RARITY_COLOR[d.rarity];
    el.onmouseenter = (e) => showTip(itemTip(id, 1, 'eq'), e.clientX + 14, e.clientY + 8);
    el.onmouseleave = hideTip;
    el.oncontextmenu = (e) => { e.preventDefault(); hideTip(); S.net.send('unequip', { slot }); };
    el.draggable = true;
    el.ondragstart = (e) => { e.dataTransfer.setData('text/plain', 'eq:' + slot); hideTip(); };
  } else {
    el.onmouseenter = (e) => showTip(`<div class="tt" style="color:#999">${SLOT_LABEL[slot]}</div>`, e.clientX + 14, e.clientY + 8);
    el.onmouseleave = hideTip;
  }
  el.ondragover = (e) => e.preventDefault();
  el.ondrop = (e) => {
    e.preventDefault();
    const dd = e.dataTransfer.getData('text/plain');
    if (dd.startsWith('bag:')) S.net.send('equip', { idx: +dd.slice(4), slot });
  };
  return el;
}

export function renderChar() {
  const me = S.me, c = S.cls, st = S.stats, body = $('char-body');
  const prim = PRIMARY[S.self.class];
  body.innerHTML = '';
  const head = document.createElement('div');
  head.innerHTML = `<div class="cname" style="color:${CLASS_COLORS[S.self.class]}">${esc(S.self.name)}</div>` +
    `<div class="csub">Stufe ${me.lvl} ${esc(S.races[S.self.race]?.name || '')} ${esc(c.name)}</div>`;
  body.appendChild(head);

  const paper = document.createElement('div');
  paper.className = 'paper';
  const col = (slots) => { const d = document.createElement('div'); d.className = 'pcol'; slots.forEach((s) => d.appendChild(slotEl(s))); return d; };
  const mid = document.createElement('div');
  mid.className = 'pmid';
  mid.innerHTML = `<div class="doll">${S.self.look?.gender === 'f' ? '👩' : '🧑'}</div>`;
  const weps = document.createElement('div');
  weps.className = 'pweps';
  WEAP.forEach((s) => weps.appendChild(slotEl(s)));
  mid.appendChild(weps);
  paper.append(col(LEFT), mid, col(RIGHT));
  body.appendChild(paper);

  const row = (a, b, cls = '') => `<div class="row2 ${cls}"><span>${a}</span><b>${b}</b></div>`;
  const stats = document.createElement('div');
  stats.className = 'statsbox';
  stats.innerHTML = row('Lebenspunkte', `${me.hp} / ${me.mh}`) + row({ rage: 'Wut', focus: 'Fokus', mana: 'Mana' }[c.resource], `${me.res} / ${me.rm}`) +
    ['str', 'agi', 'int', 'sta'].map((k) => row(STAT_LABEL[k], st[k] || 0, k === prim ? 'prim' : '')).join('') +
    row('Rüstung', `${st.armor || 0} <small>(−${st.red || 0}% Schaden)</small>`) +
    row('Waffenschaden', st.weapon || 0) + row('Schadensbonus', `+${st.dmg || 0}%`) +
    row('Erfahrung', me.xn ? `${me.xp} / ${me.xn}` : 'Max') + row('Geld', fmtCopper(S.copper));
  body.appendChild(stats);

  // equipment manager
  const mgr = document.createElement('div');
  mgr.className = 'mgr';
  mgr.innerHTML = '<div class="mh">Ausrüstungsmanager</div>';
  const names = Object.keys(S.sets).sort();
  const list = document.createElement('div');
  for (const n of names) {
    const set = S.sets[n];
    const r = document.createElement('div');
    r.className = 'setrow';
    const icons = Object.values(set).slice(0, 8).map((id) => def(id)?.icon || '').join('');
    r.innerHTML = `<span class="sn">${esc(n)}</span><span class="si">${icons}</span>`;
    const eq = document.createElement('button'); eq.className = 'small'; eq.textContent = 'Anlegen';
    eq.onclick = () => S.net.send('equipSet', { name: n });
    const sv = document.createElement('button'); sv.className = 'small'; sv.textContent = 'Überschreiben'; sv.title = 'Aktuelle Ausrüstung in diesem Set speichern';
    sv.onclick = () => S.net.send('saveSet', { name: n });
    const del = document.createElement('button'); del.className = 'small'; del.textContent = '✕';
    del.onclick = () => S.net.send('deleteSet', { name: n });
    r.append(eq, sv, del);
    list.appendChild(r);
  }
  if (!names.length) list.innerHTML = '<div class="setempty">Noch keine Sets. Lege deine Ausrüstung an und speichere sie.</div>';
  mgr.appendChild(list);
  const add = document.createElement('div');
  add.className = 'setadd';
  add.innerHTML = '<input id="set-name" maxlength="16" placeholder="Setname"><button class="small" id="set-save">Aktuelle Ausrüstung speichern</button>';
  mgr.appendChild(add);
  body.appendChild(mgr);
  const inp = add.querySelector('input');
  inp.addEventListener('keydown', (e) => { e.stopPropagation(); if (e.key === 'Enter') add.querySelector('button').click(); });
  inp.addEventListener('focus', () => { S.typing = true; });
  inp.addEventListener('blur', () => { S.typing = false; });
  add.querySelector('button').onclick = () => {
    const n = inp.value.trim();
    if (!n) { toast('Bitte einen Setnamen eingeben.'); return; }
    S.net.send('saveSet', { name: n });
  };
}

// ---------- vendor ----------
let onlyUsable = true;

export function openVendor(m) {
  S.vendorNpc = m.npc;
  S.vendorData = m;
  $('vendor-title').textContent = m.name;
  $('vendorwin').classList.remove('hidden');
  $('bags').classList.remove('hidden');
  renderBags();
  renderVendor();
}

export function closeVendor() {
  S.vendorNpc = 0; S.vendorData = null;
  $('vendorwin').classList.add('hidden');
}

export function renderVendor() {
  const m = S.vendorData;
  if (!m) return;
  const box = $('vendor-list');
  box.innerHTML = '';
  const ids = (m.stock || []).filter((id) => !onlyUsable || usable(def(id)));
  for (const id of ids) {
    const it = def(id);
    const r = document.createElement('div');
    r.className = 'vrow' + (S.copper < it.buy ? ' poor' : '') + (it.req > S.me.lvl ? ' lock' : '');
    r.innerHTML = `<div class="ii r-${it.rarity}">${it.icon}</div><div class="vn"><span class="r-${it.rarity}">${esc(it.name)}</span>` +
      `<small>${it.slot ? SLOT_LABEL[it.slot] + ' · ' + esc(it.type || '') : ''}${it.req ? ` · Stufe ${it.req}` : ''}</small></div><div class="vp">${fmtCopper(it.buy)}</div>`;
    r.onmouseenter = (e) => showTip(itemTip(id, 1, 'buy'), e.clientX - 280, e.clientY + 8);
    r.onmousemove = (e) => showTip(itemTip(id, 1, 'buy'), e.clientX - 280, e.clientY + 8);
    r.onmouseleave = hideTip;
    r.onclick = (e) => S.net.send('buy', { npc: m.npc, item: id, count: e.shiftKey && it.stack > 1 ? 5 : 1 });
    box.appendChild(r);
  }
  if (!ids.length) box.innerHTML = '<div class="setempty">Nichts im Angebot.</div>';
  $('vendor-money').innerHTML = fmtCopper(S.copper);
}

export function initItems(net) {
  $('vendor-only').onchange = (e) => { onlyUsable = e.target.checked; renderVendor(); };
  $('vendor-junk').onclick = () => { if (S.vendorNpc) net.send('sellJunk', { npc: S.vendorNpc }); };
  net.on('equip', (m) => { S.equip = m.eq || {}; if (!$('charwin').classList.contains('hidden')) renderChar(); });
  net.on('stats', (m) => { S.stats = m; if (!$('charwin').classList.contains('hidden')) renderChar(); });
  net.on('sets', (m) => { S.sets = m.sets || {}; if (!$('charwin').classList.contains('hidden')) renderChar(); });
}
