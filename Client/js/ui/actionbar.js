import { S, now } from '../state.js';
import { $, esc } from '../util.js';
import { showTip, hideTip } from './windows.js';

const KEYS = ['1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '=', '[', ']'];
let slots = [];
let faceTarget = () => {};

export function initActionbar(net, controllerRef) {
  const bar = $('actionbar');
  bar.innerHTML = '';
  slots = [];
  KEYS.forEach((k, i) => {
    const el = document.createElement('div');
    el.className = 'slot';
    el.innerHTML = `<span class="ico"></span><span class="key">${k}</span><span class="cost"></span><div class="cdo"></div><div class="cdt"></div>`;
    el.onclick = () => castSlot(i);
    el.onmouseenter = () => {
      const sk = skillAt(i);
      if (!sk) return;
      const r = el.getBoundingClientRect();
      let h = `<div class="tt">${esc(sk.name)}</div>`;
      const meta = [];
      if (sk.cost) meta.push(`${sk.cost} ${resName()}`);
      if (sk.range) meta.push(`${sk.range} m Reichweite`);
      if (sk.cast) meta.push(`${sk.cast} Sek. Zauberzeit`); else if (!sk.basic) meta.push('Sofort');
      if (sk.cooldown) meta.push(`${sk.cooldown} Sek. Abklingzeit`);
      h += `<div class="tm">${meta.join(' · ')}</div><div style="margin-top:4px">${esc(sk.desc)}</div>`;
      showTip(h, r.left, r.top - 8, true);
    };
    el.onmouseleave = hideTip;
    bar.appendChild(el);
    slots.push({ el, ico: el.querySelector('.ico'), cdo: el.querySelector('.cdo'), cdt: el.querySelector('.cdt'), cost: el.querySelector('.cost') });
  });
  faceTarget = () => controllerRef().faceTarget();

  net.on('skills', (m) => {
    const s = slots[m.skills.indexOf(m.learned)];
    if (s) { s.el.classList.remove('flash'); void s.el.offsetWidth; s.el.classList.add('flash'); }
  });
  net.on('cd', (m) => {
    const t = now();
    if (m.dur > 0) S.cd[m.skill] = { end: t + m.dur, dur: m.dur };
    if (m.gcd > 0) S.gcd = { end: t + m.gcd, dur: m.gcd };
  });
}

const resName = () => ({ rage: 'Wut', focus: 'Fokus', mana: 'Mana' }[S.cls.resource]);
const skillAt = (i) => S.defs.skills[S.skillIds[i]];

export function castSlot(i) {
  const sk = skillAt(i);
  if (!sk || S.me.d) return;
  if (['damage', 'dot', 'aoetarget', 'ground', 'polymorph', 'tame', 'taunt', 'charge'].includes(sk.kind)) faceTarget();
  S.net.send('cast', { skill: sk.id });
  const s = slots[i];
  s.el.classList.add('press');
  setTimeout(() => s.el.classList.remove('press'), 100);
}

export function updateActionbar() {
  const t = now();
  slots.forEach((s, i) => {
    const sk = skillAt(i);
    if (!sk) { s.el.classList.remove('has'); s.ico.textContent = ''; return; }
    s.el.classList.add('has');
    if (s.ico.textContent !== sk.icon) s.ico.textContent = sk.icon;
    s.cost.textContent = sk.cost || '';
    let rem = 0, dur = 1;
    const cd = S.cd[sk.id];
    if (cd && cd.end > t) { rem = cd.end - t; dur = cd.dur; }
    if (!sk.basic && S.gcd.end > t && S.gcd.end - t > rem) { rem = S.gcd.end - t; dur = S.gcd.dur; }
    if (rem > 0) {
      s.cdo.style.display = 'block';
      s.cdo.style.background = `conic-gradient(#0000 ${(1 - rem / dur) * 360}deg, #000b 0)`;
      s.cdt.style.display = rem > 1.0 && dur > 1.2 ? 'grid' : 'none';
      s.cdt.textContent = rem.toFixed(rem < 10 ? 0 : 0);
    } else { s.cdo.style.display = 'none'; s.cdt.style.display = 'none'; }
    s.el.classList.toggle('nores', S.me.res < sk.cost);
    s.el.classList.toggle('active', (sk.basic && S.me.auto) || (sk.kind === 'stance' && S.me.stance === sk.stance));
  });
}
