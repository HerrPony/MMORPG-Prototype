import { S } from '../state.js';
import { $, esc, fmtCopper, toast, CLASS_COLORS, contOf } from '../util.js';
import { renderTracker } from './hud.js';
import { drawWorldMap, setMapCont } from './worldmap.js';
import { renderPartyWindow } from './social.js';

// ---------- tooltip ----------
import { itemTip, renderBags, renderChar, openVendor, closeVendor, renderVendor, initItems } from './items.js';
export function showTip(html, x, y, above = false) {
  const t = $('tooltip');
  t.innerHTML = html;
  t.classList.remove('hidden');
  const w = t.offsetWidth, h = t.offsetHeight;
  let left = Math.min(innerWidth - w - 8, Math.max(8, x));
  let top = above ? y - h : y;
  top = Math.max(8, Math.min(innerHeight - h - 8, top));
  t.style.left = left + 'px';
  t.style.top = top + 'px';
}
export function hideTip() { $('tooltip').classList.add('hidden'); }


// ---------- draggable windows ----------
function makeDraggable(win) {
  const bar = win.querySelector('.wtitle');
  let ox = 0, oy = 0, drag = false;
  bar.addEventListener('mousedown', (e) => {
    if (e.target.closest('.x')) return;
    drag = true;
    const r = win.getBoundingClientRect();
    win.style.transform = 'none';
    win.style.left = r.left + 'px'; win.style.top = r.top + 'px';
    win.style.right = 'auto'; win.style.bottom = 'auto';
    ox = e.clientX - r.left; oy = e.clientY - r.top;
    e.preventDefault();
  });
  addEventListener('mousemove', (e) => {
    if (!drag) return;
    win.style.left = Math.max(0, Math.min(innerWidth - 80, e.clientX - ox)) + 'px';
    win.style.top = Math.max(0, Math.min(innerHeight - 40, e.clientY - oy)) + 'px';
  });
  addEventListener('mouseup', () => { drag = false; });
}

const WIN_IDS = ['dialog', 'lootwin', 'bags', 'questlog', 'charwin', 'menuwin', 'spellbook', 'vendorwin', 'mapwin', 'partywin'];
const show = (id) => $(id).classList.remove('hidden');
const hide = (id) => $(id).classList.add('hidden');
const visible = (id) => !$(id).classList.contains('hidden');

export function toggle(name) {
  const id = { menu: 'menuwin', char: 'charwin', quests: 'questlog', map: 'mapwin', party: 'partywin' }[name] || name;
  if (visible(id)) { hide(id); return; }
  show(id);
  if (id === 'bags') renderBags();
  if (id === 'questlog') renderQuestLog();
  if (id === 'charwin') renderChar();
  if (id === 'spellbook') renderSpellbook();
  if (id === 'mapwin') openMap();
  if (id === 'partywin') renderPartyWindow();
}

export function openMap() {
  const w = mapwin;
  w.style.left = Math.max(10, innerWidth / 2 - 380) + 'px';
  w.style.top = Math.max(10, innerHeight / 2 - 400) + 'px';
  w.style.transform = 'none';
  setMapCont(contOf(S.local.x));
  const sz = Math.max(300, Math.min(700, innerHeight - 190));
  $('mapcv').style.width = $('mapcv').style.height = sz + 'px';
  $('map-wrap').style.width = $('map-wrap').style.height = sz + 'px';
  drawWorldMap();
}

export function renderSpellbook() {
  const known = new Set(S.skillIds);
  const rows = S.cls.skills.map((id) => {
    const d = S.defs.skills[id];
    const meta = [];
    if (d.cost) meta.push(`${d.cost} ${{ rage: 'Wut', focus: 'Fokus', mana: 'Mana' }[S.cls.resource]}`);
    if (d.range) meta.push(`${d.range} m`);
    if (d.cast) meta.push(`${d.cast} Sek. Zauberzeit`);
    if (d.cooldown) meta.push(`${d.cooldown} Sek. Abklingzeit`);
    const state = known.has(id) ? '' : d.quest ? `<span class="tl">Klassenquest ab Stufe ${d.level}</span>` : `<span class="tl">Stufe ${d.level} · beim Klassenlehrer</span>`;
    return `<div class="trow ${known.has(id) ? 'known' : 'locked'} book"><div class="ti">${d.icon}</div>` +
      `<div class="tn"><b>${esc(d.name)}</b><small>${esc(meta.join(' · '))}</small><small>${esc(d.desc)}</small></div><div class="ta">${state}</div></div>`;
  });
  $('spell-body').innerHTML = rows.join('');
}

// close top-most window; returns true if something was closed
export function closeTopWindow() {
  for (const id of ['lootwin', 'vendorwin', 'dialog', 'mapwin', 'partywin', 'spellbook', 'questlog', 'bags', 'charwin', 'menuwin']) {
    if (visible(id)) { hide(id); if (id === 'vendorwin') closeVendor(); if (id === 'dialog') S.dialogNpc = 0; if (id === 'lootwin') S.lootId = 0; return true; }
  }
  return false;
}

// ---------- dialog ----------
const replName = (t) => esc(t).replaceAll('%n', S.self.name);

function rewardHtml(q) {
  let h = '';
  if (q.reward.xp) h += `<span>🌟 ${q.reward.xp} Erfahrung</span>`;
  if (q.reward.copper) h += `<span>${fmtCopper(q.reward.copper)}</span>`;
  for (const it of q.reward.items || []) {
    const d = S.defs.items[it];
    h += `<span class="r-${d.rarity}" style="background:#0003;padding:2px 6px;border-radius:3px">${d.icon} ${esc(d.name)}</span>`;
  }
  return h;
}

function objectivesHtml(q, counts) {
  return q.objectives.map((o, i) => `<div>• ${esc(o.text)}${counts ? `: ${counts[i] || 0}/${o.count}` : ` (${o.count})`}</div>`).join('');
}

// An NPC dialog has up to four views: quests, trainer, griffin flights and (as a window) the vendor.
let dlgView = 'greet';
let dlgQuest = null; // quest id currently shown in the quest view

function dialogViews(m) {
  const v = [];
  if (m.quests?.length) v.push(['quest', m.quests.some((q) => q.state === 'ready') ? '❓ Quests' : '❗ Quests']);
  if (m.trainer) v.push(['trainer', '📚 Training']);
  if (m.flight) v.push(['flight', '🦅 Greifenflug']);
  if (m.ferry?.ship) v.push(['ferry', '⛵ Fähre']);
  return v;
}

function openDialog(m) {
  const prev = S.dlg && S.dlg.npc === m.npc && visible('dialog') ? dlgView : null;
  if (!(S.dlg && S.dlg.npc === m.npc && visible('dialog'))) dlgQuest = null;
  S.dlg = m;
  S.dialogNpc = m.npc;
  if (m.vendor && !m.quests?.length && !m.trainer && !m.flight && !m.ferry?.ship) { hide('dialog'); openVendor(m); return; }
  $('dlg-title').textContent = m.name + (m.title ? ` – ${m.title}` : '');
  const views = dialogViews(m);
  dlgView = views.find((v) => v[0] === prev) ? prev : (views[0]?.[0] || 'greet');
  if (dlgView === 'quest' && m.quests.length === 1) dlgQuest = m.quests[0].id;
  if (dlgQuest && !m.quests?.some((q) => q.id === dlgQuest)) dlgQuest = null;
  show('dialog');
  renderDialog();
}

const QMARK = { available: '❗', ready: '❓', progress: '❔' };

function renderDialog() {
  const m = S.dlg;
  if (!m) return;
  const views = dialogViews(m);
  const tabs = $('dlg-tabs');
  tabs.innerHTML = '';
  if (views.length > 1 || (m.vendor && views.length)) {
    if (m.vendor) {
      const vb = document.createElement('button');
      vb.textContent = '💰 Handel';
      vb.onclick = () => openVendor(m);
      tabs.appendChild(vb);
    }
    for (const [id, label] of views) {
      const b = document.createElement('button');
      b.className = id === dlgView ? 'active' : '';
      b.textContent = label;
      b.onclick = () => { dlgView = id; renderDialog(); };
      tabs.appendChild(b);
    }
  }
  const main = $('dlg-main');
  main.onclick = null;
  main.classList.add('hidden');
  $('dlg-quest').textContent = ''; $('dlg-text').innerHTML = ''; $('dlg-obj').innerHTML = ''; $('dlg-reward').innerHTML = ''; $('dlg-list').innerHTML = '';

  if (dlgView === 'quest') {
    if (!dlgQuest) { // several quests: pick one
      $('dlg-text').textContent = m.greet || '';
      const list = $('dlg-list');
      for (const o of m.quests) {
        const q = S.defs.quests[o.id];
        const row = document.createElement('div');
        row.className = 'trow sell';
        row.innerHTML = `<div class="ti">${QMARK[o.state]}</div><div class="tn"><b>${esc(q.title)}</b><small>Stufe ${q.level}${q.tag ? ' · ' + esc(q.tag) : ''}</small></div>`;
        row.onclick = () => { dlgQuest = o.id; renderDialog(); };
        list.appendChild(row);
      }
      return;
    }
    const offer = m.quests.find((x) => x.id === dlgQuest);
    const q = S.defs.quests[dlgQuest];
    const st = S.quests.find((x) => x.id === dlgQuest);
    $('dlg-quest').innerHTML = `${esc(q.title)} ${q.tag ? `<span class="qtag">${esc(q.tag)}</span>` : ''}`;
    const meta = `<div class="qmeta">Stufe ${q.level}${q.timeLimit ? ` · Zeitlimit: ${q.timeLimit} Sek.` : ''}</div>`;
    if (offer.state === 'available') {
      $('dlg-text').innerHTML = meta + replName(q.text);
      $('dlg-obj').innerHTML = `<div class="ob">Aufgabe</div>${objectivesHtml(q)}`;
      $('dlg-reward').innerHTML = `<div class="ob">Belohnung</div><div class="reward">${rewardHtml(q)}</div>`;
      main.textContent = 'Quest annehmen';
      main.onclick = () => { S.net.send('acceptQuest', { quest: q.id }); hide('dialog'); S.dialogNpc = 0; };
      main.classList.remove('hidden');
    } else if (offer.state === 'progress') {
      $('dlg-text').innerHTML = replName(q.progress);
      $('dlg-obj').innerHTML = `<div class="ob">Aufgabe</div>${objectivesHtml(q, st?.counts)}`;
    } else {
      $('dlg-text').innerHTML = replName(q.complete);
      $('dlg-reward').innerHTML = `<div class="ob">Belohnung</div><div class="reward">${rewardHtml(q)}</div>`;
      main.textContent = 'Quest abschließen';
      main.onclick = () => { S.net.send('turnInQuest', { quest: q.id }); hide('dialog'); S.dialogNpc = 0; };
      main.classList.remove('hidden');
    }
    if (m.quests.length > 1) {
      const back = document.createElement('button');
      back.className = 'small'; back.textContent = '‹ Alle Quests'; back.onclick = () => { dlgQuest = null; renderDialog(); };
      $('dlg-list').appendChild(back);
    }
  } else if (dlgView === 'ferry') {
    const f = m.ferry;
    $('dlg-text').innerHTML = `<b>${esc(f.status)}</b><br><br>Die Fähre pendelt im festen Takt zwischen Sturmhafen und Lotushafen und liegt dabei längsseits am Steg. Geh einfach über den Steg auf das Schiff, bewege dich frei an Deck und steige im Zielhafen wieder über den Steg aus.`;
  } else if (dlgView === 'flight') {
    const f = m.flight;
    $('dlg-text').innerHTML = f.ship
      ? ''
      : `Du bist in <b>${esc(f.hereName)}</b>. Dein Flugpunkt hier ist freigeschaltet. Dein Geld: ${fmtCopper(S.copper)}`;
    const list = $('dlg-list');
    if (!f.points.length) list.innerHTML = '<div style="padding:8px;color:#6a5a3a">Du kennst noch keinen anderen Flugpunkt. Sprich mit den Greifenmeistern in anderen Städten, um sie freizuschalten. Greifen fliegen nur auf demselben Kontinent.</div>';
    for (const p of f.points) {
      const row = document.createElement('div');
      row.className = 'trow sell';
      row.innerHTML = `<div class="ti">${f.ship ? '⛵' : '🦅'}</div><div class="tn"><b>${esc(p.name)}</b></div><div class="tp">${fmtCopper(p.price)}</div>`;
      row.onclick = () => { if (S.copper >= p.price) { S.net.send('flight', { npc: m.npc, dest: p.id }); hide('dialog'); S.dialogNpc = 0; } else toast('Du hast nicht genug Geld.'); };
      list.appendChild(row);
    }
  } else if (dlgView === 'trainer') {
    const t = m.trainer;
    if (!t.mine) {
      $('dlg-text').textContent = `Ich bilde nur ${t.className} aus. Such dir einen Lehrer deiner eigenen Klasse.`;
    } else {
      $('dlg-text').innerHTML = `Dein Geld: ${fmtCopper(S.copper)}`;
      const list = $('dlg-list');
      for (const s of t.skills) {
        const d = S.defs.skills[s.id];
        const row = document.createElement('div');
        row.className = 'trow ' + s.state;
        const action = s.state === 'known' ? '<span class="tk">Gelernt</span>'
          : s.state === 'locked' ? `<span class="tl">Ab Stufe ${s.level}</span>`
            : `<button class="small learn" ${S.copper < s.price ? 'disabled' : ''}>Lernen</button>`;
        row.innerHTML = `<div class="ti">${d.icon}</div><div class="tn"><b>${esc(d.name)}</b><small>Stufe ${s.level} · ${esc(d.desc)}</small></div>` +
          `<div class="tp">${s.state === 'known' ? '' : fmtCopper(s.price)}</div><div class="ta">${action}</div>`;
        const btn = row.querySelector('.learn');
        if (btn) btn.onclick = () => S.net.send('train', { npc: m.npc, skill: s.id });
        list.appendChild(row);
      }
      if (S.cls.skills.some((id) => S.defs.skills[id].quest)) {
        const note = document.createElement('div');
        note.className = 'qmeta'; note.style.padding = '8px 4px';
        note.textContent = 'Besondere Fertigkeiten (Stufe 10) erlernst du durch deine Klassenquest in Sonnenburg.';
        list.appendChild(note);
      }
    }
  } else {
    $('dlg-text').textContent = m.greet || '…';
  }
}

// ---------- loot ----------
function openLoot(m) {
  S.lootId = m.id;
  $('loot-title').textContent = m.name;
  const body = $('loot-body');
  body.innerHTML = '';
  m.items.forEach((it, idx) => {
    const d = S.defs.items[it.i];
    const row = document.createElement('div');
    row.className = 'lootrow';
    row.innerHTML = `<div class="ii r-${d.rarity}">${d.icon}</div><div class="r-${d.rarity}">${esc(d.name)}${it.c > 1 ? ' x' + it.c : ''}</div>`;
    row.onclick = () => S.net.send('loot', { id: m.id, idx });
    row.onmouseenter = (e) => showTip(itemTip(it.i, it.c), e.clientX + 14, e.clientY + 8);
    row.onmouseleave = hideTip;
    body.appendChild(row);
  });
  if (m.copper > 0) {
    const row = document.createElement('div');
    row.className = 'lootrow';
    row.innerHTML = `<div class="ii" style="border-color:#c90">💰</div><div>${fmtCopper(m.copper)}</div>`;
    row.onclick = () => S.net.send('loot', { id: m.id, idx: -2 });
    body.appendChild(row);
  }
  $('loot-all').onclick = () => S.net.send('loot', { id: m.id, idx: -1 });
  const w = $('lootwin');
  if (!visible('lootwin')) {
    w.style.transform = 'none';
    w.style.left = Math.round(innerWidth * 0.42) + 'px';
    w.style.top = Math.round(innerHeight * 0.3) + 'px';
  }
  show('lootwin');
}

let selQuest = null;
export function renderQuestLog() {
  const list = $('ql-list'), det = $('ql-detail');
  list.innerHTML = '';
  if (!S.quests.length) { det.innerHTML = 'Du hast keine aktiven Quests.'; return; }
  if (!S.quests.find((q) => q.id === selQuest)) selQuest = S.quests[0].id;
  for (const q of S.quests) {
    const d = S.defs.quests[q.id];
    const row = document.createElement('div');
    row.className = 'qi' + (q.id === selQuest ? ' sel' : '');
    row.textContent = d.title;
    row.onclick = () => { selQuest = q.id; renderQuestLog(); };
    list.appendChild(row);
  }
  const q = S.quests.find((x) => x.id === selQuest), d = S.defs.quests[selQuest];
  det.innerHTML = `<h4>${esc(d.title)} ${d.tag ? `<span class="qtag">${esc(d.tag)}</span>` : ''}</h4><div class="qmeta">Stufe ${d.level}${q.end ? ` · Zeit übrig: ${Math.max(0, Math.ceil(q.end - performance.now() / 1000))} Sek.` : ''}</div><div style="white-space:pre-wrap">${replName(q.status === 'complete' ? d.complete : d.text)}</div>` +
    `<div class="ob">Ziele</div>${objectivesHtml(d, q.counts)}` +
    `<div class="ob">Belohnung</div><div class="reward">${rewardHtml(d)}</div>` +
    `<button class="small" id="ql-abandon">Quest abbrechen</button>`;
  $('ql-abandon').onclick = () => S.net.send('abandonQuest', { quest: selQuest });
}

const CONTROLS = [
  ['W A S D', 'Bewegen'], ['Leertaste', 'Springen'], ['Rechtsklick halten', 'Kamera drehen (Charakter blickt mit)'],
  ['Linksklick halten', 'Kamera drehen'], ['Linksklick', 'NPC / Spieler / Gegner auswählen'],
  ['Rechtsklick', 'NPC ansprechen · Gegner angreifen · Leiche looten'], ['Beide Maustasten', 'Vorwärts laufen'],
  ['Mausrad', 'Zoom'], ['Tab', 'Nächsten Gegner anvisieren'], ['1 – 5', 'Fertigkeiten'], ['F', 'Mit Ziel interagieren'],
  ['R', 'Automatisch laufen'], ['Shift', 'Gehen (gedrückt halten)'], ['Enter  /  /', 'Chat'], ['B  L  C  P', 'Taschen · Questlog · Charakter · Zauberbuch'], ['Esc', 'Fenster schließen / Ziel abwählen / Menü'],
];

// ---------- init ----------
export function initWindows(net) {
  WIN_IDS.forEach((id) => makeDraggable($(id)));
  document.querySelectorAll('[data-close]').forEach((b) => (b.onclick = () => {
    const id = b.dataset.close;
    hide(id);
    if (id === 'dialog') S.dialogNpc = 0;
    if (id === 'lootwin') S.lootId = 0;
    if (id === 'vendorwin') closeVendor();
  }));
  document.querySelectorAll('#micro [data-win]').forEach((b) => (b.onclick = () => toggle(b.dataset.win)));

  $('controls-help').innerHTML = CONTROLS.map(([k, v]) => `<div><b>${k}</b>${v}</div>`).join('');
  $('m-controls').onclick = () => $('controls-help').classList.toggle('hidden');
  $('m-logout').onclick = () => { hide('menuwin'); net.send('logout'); };

  initItems(net);
  net.on('dialog', openDialog);
  net.on('loot', openLoot);
  net.on('lootClose', (m) => { if (S.lootId === m.id) { hide('lootwin'); S.lootId = 0; } hideTip(); });
  net.on('inv', (m) => {
    S.inv = m.items || []; S.copper = m.copper; S.invMax = m.max;
    if (visible('bags')) renderBags();
    if (visible('charwin')) renderChar();
    if (visible('dialog') && S.dlg) renderDialog();
    if (visible('vendorwin')) renderVendor();
  });
  net.on('skills', (m) => {
    S.skillIds = m.skills;
    if (visible('spellbook')) renderSpellbook();
    const d = S.defs.skills[m.learned];
    if (d) toast(`Neue Fertigkeit: ${d.name}`);
  });
  net.on('quests', (m) => {
    const now = performance.now() / 1000;
    S.quests = m.list.map((q) => ({ ...q, end: q.left !== undefined ? now + q.left : 0 }));
    renderTracker();
    if (visible('questlog')) renderQuestLog();
  });
  net.on('questDone', () => { hide('dialog'); S.dialogNpc = 0; });
}

// close dialogs/loot when walking away
export function updateWindows() {
  const check = (id, key, maxD) => {
    if (!visible(id) || !S[key]) return;
    const e = S.ents.get(S[key]);
    if (!e || Math.hypot(e.rx - S.local.x, e.rz - S.local.z) > maxD) { hide(id); S[key] = 0; }
  };
  check('dialog', 'dialogNpc', 12);
  check('lootwin', 'lootId', 12);
  if (S.vendorNpc) {
    const e = S.ents.get(S.vendorNpc);
    if (!e || Math.hypot(e.rx - S.local.x, e.rz - S.local.z) > 12) { hide('vendorwin'); closeVendor(); }
  }
}

export function resetWindows() {
  WIN_IDS.forEach(hide);
  closeVendor();
  S.dialogNpc = 0; S.lootId = 0;
}
