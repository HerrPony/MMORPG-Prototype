// Zone banner, party, pet frames, invitations and the player context menu.
import { S } from '../state.js';
import { $, esc, CLASS_COLORS, CLASS_ICONS } from '../util.js';

let bannerT;

function showBanner(zone, sub, lo, hi, sea) {
  const b = $('zonebanner');
  $('zb-zone').textContent = zone;
  $('zb-sub').textContent = sub || (sea ? 'Offene See' : lo ? `Stufe ${lo}–${hi}` : '');
  b.classList.remove('hidden');
  b.style.animation = 'none'; void b.offsetWidth; b.style.animation = '';
  clearTimeout(bannerT);
  bannerT = setTimeout(() => b.classList.add('hidden'), 4200);
}

export function initSocial(net) {
  net.on('zone', (m) => {
    const zoneChanged = m.zone !== S.zone.name;
    S.zone = { name: m.zone, sub: m.sub, lo: m.lo, hi: m.hi };
    $('zone-name').textContent = m.sub ? `${m.zone} – ${m.sub}` : m.zone;
    if (zoneChanged) showBanner(m.zone, m.sub ? m.sub : '', m.lo, m.hi, m.sea);
    else if (m.sub) showBanner(m.sub, m.zone, 0, 0);
    S.onZone?.(m);
  });
  net.on('disc', (m) => {
    S.disc = m.ids;
    if (m.name && !m.new.startsWith('zone:')) { /* the banner already announced it */ }
  });
  net.on('fp', (m) => { S.fp = m.ids; });
  net.on('party', (m) => {
    S.party = { id: m.id, leader: m.leader, members: m.members, states: S.party.states };
    renderPartyWindow();
  });
  net.on('partyInvite', (m) => {
    $('inv-text').innerHTML = `<b>${esc(m.from)}</b> lädt dich in seine Gruppe ein.`;
    $('invite').classList.remove('hidden');
    clearTimeout(S._invT);
    S._invT = setTimeout(() => $('invite').classList.add('hidden'), 30000);
  });
  $('inv-yes').onclick = () => { net.send('inviteReply', { accept: true }); $('invite').classList.add('hidden'); };
  $('inv-no').onclick = () => { net.send('inviteReply', { accept: false }); $('invite').classList.add('hidden'); };
  net.on('flightStart', (m) => { $('fb-dest').textContent = (m.ship ? '⛵ ' : '') + m.dest; $('flightbar').classList.remove('hidden'); });
  net.on('flightEnd', () => $('flightbar').classList.add('hidden'));
  document.addEventListener('mousedown', (e) => { if (!e.target.closest('#playermenu')) $('playermenu').classList.add('hidden'); });
}

// right click on another player
export function playerMenu(e, x, y) {
  const menu = $('playermenu');
  const mine = S.party.members.some((m) => m.id === e.id);
  const leader = S.party.leader === S.id;
  menu.innerHTML = `<div class="pmn" style="color:${CLASS_COLORS[e.c] || '#fff'}">${esc(e.n)} <small>Stufe ${e.l}</small></div>`;
  const add = (label, fn) => { const b = document.createElement('button'); b.textContent = label; b.onclick = () => { fn(); menu.classList.add('hidden'); }; menu.appendChild(b); };
  if (!mine) add('👥 In Gruppe einladen', () => S.net.send('invite', { id: e.id }));
  if (mine && leader) { add('⭐ Zum Anführer machen', () => S.net.send('promote', { id: e.id })); add('🚪 Aus Gruppe werfen', () => S.net.send('kick', { id: e.id })); }
  add('💬 Anflüstern', () => S.openChat?.(`/p `));
  menu.style.left = Math.min(innerWidth - 190, x) + 'px';
  menu.style.top = Math.min(innerHeight - 140, y) + 'px';
  menu.classList.remove('hidden');
}

export function renderPartyWindow() {
  const body = $('party-body');
  if (!body) return;
  if (!S.party.members.length) {
    body.innerHTML = '<div class="setempty">Du bist in keiner Gruppe. Rechtsklick auf einen Spieler, um ihn einzuladen (oder <b>/einladen Name</b>). Gruppen haben bis zu 5 Spieler.</div>';
    return;
  }
  body.innerHTML = '';
  for (const m of S.party.members) {
    const row = document.createElement('div');
    row.className = 'prow';
    row.innerHTML = `<span class="pi" style="border-color:${CLASS_COLORS[m.c]}">${CLASS_ICONS[m.c] || '?'}</span><span class="pn">${m.id === S.party.leader ? '⭐ ' : ''}${esc(m.n)} <small>Stufe ${m.l}</small></span>`;
    if (S.party.leader === S.id && m.id !== S.id) {
      const k = document.createElement('button'); k.className = 'small'; k.textContent = '✕'; k.title = 'Aus der Gruppe werfen'; k.onclick = () => S.net.send('kick', { id: m.id });
      row.appendChild(k);
    }
    body.appendChild(row);
  }
  const leave = document.createElement('button');
  leave.className = 'small'; leave.textContent = 'Gruppe verlassen';
  leave.onclick = () => S.net.send('leaveParty');
  body.appendChild(leave);
}

// per frame: party member frames and pet frames
let partyKey = '', petKey = '';
export function updateSocial() {
  // party frames
  const box = $('partyframes');
  const others = S.party.members.filter((m) => m.id !== S.id);
  const key = others.map((m) => m.id).join(',');
  if (key !== partyKey) {
    partyKey = key;
    box.innerHTML = others.map((m) => `<div class="pf" data-id="${m.id}"><div class="pfn" style="color:${CLASS_COLORS[m.c]}">${m.id === S.party.leader ? '⭐' : ''}${esc(m.n)}</div><div class="pfb hp"><i></i></div><div class="pfb res"><i></i></div></div>`).join('');
    box.querySelectorAll('.pf').forEach((el) => { el.onclick = () => { S.net.send('target', { id: +el.dataset.id }); S.me.tg = +el.dataset.id; }; });
  }
  for (const st of S.party.states) {
    const el = box.querySelector(`.pf[data-id="${st.id}"]`);
    if (!el) continue;
    el.querySelector('.hp i').style.width = (st.mh ? st.hp / st.mh * 100 : 0) + '%';
    el.querySelector('.res i').style.width = (st.rm ? st.res / st.rm * 100 : 0) + '%';
    el.classList.toggle('dead', !!st.d);
    el.classList.toggle('far', st.i !== (S.instance ? 1 : 0));
  }
  // pet frames
  const pet = S.me.pet;
  const pk = pet ? [pet.id || 0, pet.eid || 0, pet.name].join() : '';
  const pbox = $('petframes');
  if (pk !== petKey) {
    petKey = pk;
    let h = '';
    if (pet && pet.name) h += `<div class="pf pet" data-id="${pet.id || 0}"><div class="pfn">🐾 ${esc(pet.name)}${pet.dead ? ' (gefallen)' : pet.id ? '' : ' (nicht gerufen)'}</div><div class="pfb hp"><i></i></div></div>`;
    if (pet && pet.eid) h += `<div class="pf pet" data-id="${pet.eid}"><div class="pfn">💧 Wasserelementar</div><div class="pfb hp"><i></i></div></div>`;
    pbox.innerHTML = h;
    pbox.querySelectorAll('.pf').forEach((el) => { el.onclick = () => { if (+el.dataset.id) { S.net.send('target', { id: +el.dataset.id }); S.me.tg = +el.dataset.id; } }; });
  }
  if (pet) {
    const els = pbox.querySelectorAll('.pf');
    let i = 0;
    if (pet.name) { const f = els[i++]; if (f) f.querySelector('i').style.width = (pet.mh ? pet.hp / pet.mh * 100 : 0) + '%'; }
    if (pet.eid) { const f = els[i++]; if (f) f.querySelector('i').style.width = (pet.emh ? pet.ehp / pet.emh * 100 : 0) + '%'; }
  }
}
