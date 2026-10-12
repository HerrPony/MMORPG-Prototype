import { S, now } from '../state.js';
import { $, esc, CLASS_COLORS, CLASS_ICONS } from '../util.js';
import { gameHours, sunInfo } from '../daynight.js';

const STATUS_ICONS = { slow: ['🐌', 'Verlangsamt', true], dot: ['🩸', 'Schaden über Zeit', true], buff: ['⬆️', 'Gestärkt', false], shield: ['🛡️', 'Schild', false] };

function setBar(bar, cur, max, text) {
  const f = bar.querySelector('.fill');
  const pct = max > 0 ? Math.max(0, Math.min(100, (cur / max) * 100)) : 0;
  f.style.width = pct + '%';
  bar.querySelector('.txt').textContent = text;
}

function portraitEmoji(e) {
  if (e.k === 'player') return e.lk?.gender === 'f' ? '👩' : '🧑';
  if (e.m === 'smith') return '🔨';
  const mobs = { wolf: '🐺', boar: '🐗', spider: '🕷️', bear: '🐻', bandit: '🥷', hyena: '🐕', plains_bandit: '🥷', bull: '🐂', wolf_alpha: '🐺', bandit_chief: '🥷', water_elem: '💧', croc: '🐊', zombie: '🧟', wisp: '✨', witch: '🧙', hag: '🧙', zombie_giant: '🧟', icewolf: '🐺', yeti: '❄️', yeti_chief: '❄️', polarbear: '🐻‍❄️', polarbear_big: '🐻‍❄️', bear_big: '🐻', frostspirit: '👻', scorpion: '🦂', scorpion_king: '🦂', desert_raider: '🥷', mummy: '🧟', tomb_guardian: '🗿', kobold: '👹', kobold_shaman: '👹', kobold_boss: '👹', kobold_king: '👑', skeleton: '💀', cavespider: '🕷️', spider_queen: '🕷️', sheep: '🐑', water_elem_pet: '💧' };
  if (mobs[e.m]) return mobs[e.m];
  const npcs = { forester: '🧔', farmer: '👩‍🌾', scholar: '👩‍🔬', guard: '💂', huntsman: '🧔🏿', merchant: '👩‍💼', trainer_warrior: '🛡️', trainer_hunter: '🏹', trainer_mage: '🔮', mayor: '👑', priest: '🙏', elder: '👴', herbalist: '🌿', fisher: '🎣', sheikh: '🕌', nomad: '🐪', flightmaster: '🦅' };
  return npcs[e.m] || '❓';
}

function renderBuffs(box, st) {
  const key = (st || []).join(',');
  if (box._k === key) return;
  box._k = key;
  box.innerHTML = (st || []).map((s) => {
    const d = STATUS_ICONS[s]; if (!d) return '';
    return `<i class="${d[2] ? 'bad' : ''}" title="${d[1]}">${d[0]}</i>`;
  }).join('');
}

export function initHud(net) {
  // errors
  let errT;
  net.on('err', (m) => showError(m.text));
  function showError(text) {
    const e = $('errmsg');
    e.textContent = text;
    e.classList.add('show');
    clearTimeout(errT);
    errT = setTimeout(() => e.classList.remove('show'), 1600);
  }

  // chat
  const log = $('chat-log'), input = $('chat-input'), chat = $('chat');
  function addChat(ch, html) {
    const d = document.createElement('div');
    d.className = 'c-' + ch;
    d.innerHTML = html;
    log.appendChild(d);
    while (log.children.length > 150) log.firstChild.remove();
    log.scrollTop = log.scrollHeight;
  }
  S.addChat = addChat;
  net.on('chat', (m) => {
    if (m.ch === 'say') addChat('say', `<b style="color:#ffd">[${esc(m.from)}]</b>: ${esc(m.text)}`);
    else if (m.ch === 'party') addChat('party', `<b>[Gruppe] ${esc(m.from)}</b>: ${esc(m.text)}`);
    else addChat(m.ch, esc(m.text));
  });
  document.querySelectorAll('#chat-tabs button').forEach((b) => (b.onclick = () => {
    document.querySelectorAll('#chat-tabs button').forEach((x) => x.classList.toggle('active', x === b));
    chat.dataset.f = b.dataset.f;
    log.scrollTop = log.scrollHeight;
  }));
  input.addEventListener('focus', () => { S.typing = true; chat.classList.add('typing'); });
  input.addEventListener('blur', () => { S.typing = false; chat.classList.remove('typing'); });
  input.addEventListener('keydown', (e) => {
    e.stopPropagation();
    if (e.key === 'Enter') {
      const t = input.value.trim();
      input.value = '';
      input.blur();
      if (!t) return;
      if (t === '/hilfe' || t === '/help') {
        addChat('sys', 'Steuerung: WASD bewegen · Leertaste springen · Rechtsklick halten: Kamera · Linksklick: auswählen · Rechtsklick: interagieren/angreifen/looten · Tab: Ziel wechseln · 1-5: Fertigkeiten · B/L/C/P: Fenster · R: Autorun · F: mit Ziel interagieren');
        return;
      }
      net.send('chat', { text: t.startsWith('/s ') ? t.slice(3) : t });
    } else if (e.key === 'Escape') { input.value = ''; input.blur(); }
  });
  S.openChat = (prefix) => { input.value = prefix; input.focus(); };

  // revive
  $('btn-revive').onclick = () => net.send('revive');

  // level up banner
  net.on('lvl', (m) => {
    $('levelup-n').textContent = m.lvl;
    const b = $('levelup');
    b.classList.remove('hidden');
    b.style.animation = 'none'; void b.offsetWidth; b.style.animation = '';
    setTimeout(() => b.classList.add('hidden'), 3000);
  });

  // cast bar
  net.on('cast', (m) => { S.cast = { skill: m.skill, start: now(), dur: m.dur }; });
  net.on('castStop', () => { S.cast = null; });

  const pf = $('pframe'), tf = $('tframe');
  const pHp = pf.querySelector('.bar.hp'), pRes = pf.querySelector('.bar.res'), tHp = tf.querySelector('.bar.hp');
  let lastClock = '';

  return {
    showError,
    update() {
      const me = S.me, c = S.cls, ent = S.ents.get(S.id);
      // player frame
      $('pf-name').textContent = S.self.name;
      $('pf-name').style.color = CLASS_COLORS[S.self.class] || '#fff';
      $('pf-lvl').textContent = me.lvl;
      pf.querySelector('.pemoji').textContent = S.self.look?.gender === 'f' ? '👩' : '🧑';
      $('pf-portrait').style.background = `radial-gradient(circle at 35% 30%, ${CLASS_COLORS[S.self.class]}88, #1b2530)`;
      setBar(pHp, me.hp, me.mh, `${me.hp} / ${me.mh}`);
      pRes.querySelector('.fill').style.background = c.resColor;
      setBar(pRes, me.res, me.rm, `${me.res} / ${me.rm}`);
      renderBuffs($('pf-buffs'), ent?.st);

      // target frame
      const t = S.ents.get(me.tg);
      if (t) {
        tf.classList.remove('hidden');
        $('tf-name').textContent = t.n;
        $('tf-name').style.color = t.k === 'player' ? CLASS_COLORS[t.c] : t.k === 'npc' ? '#ffe27a' : '#ffe14a';
        $('tf-lvl').textContent = t.k === 'npc' ? '!' : t.l;
        tf.querySelector('.pemoji').textContent = portraitEmoji(t);
        tHp.classList.toggle('hostile', t.k === 'mob');
        tHp.classList.toggle('dead', !!t.d);
        setBar(tHp, t.hp, t.mh, t.k === 'npc' ? '' : t.d ? 'Tot' : `${t.hp} / ${t.mh}`);
        renderBuffs($('tf-buffs'), t.st);
      } else tf.classList.add('hidden');

      // stance, elite tag, quest timers
      const stEl = $('pf-stance');
      if (S.cls.id === 'warrior') {
        stEl.classList.remove('hidden');
        stEl.textContent = me.stance === 'defensive' ? '🛡️ Verteidigungshaltung' : '⚔️ Kampfhaltung';
      } else stEl.classList.add('hidden');
      $('tf-tag').textContent = t ? (t.bs ? 'Boss' : t.el ? 'Elite' : '') : '';
      tf.classList.toggle('elite', !!(t && t.el));
      const nowS = performance.now() / 1000;
      document.querySelectorAll('.qtimer').forEach((el) => {
        const left = Math.max(0, Math.ceil(+el.dataset.end - nowS));
        el.textContent = `⏱ ${left} Sek.`;
        el.style.color = left < 10 ? '#ff5a4a' : '#ffd650';
      });

      // xp
      $('xpfill').style.width = (me.xn > 0 ? Math.min(100, (me.xp / me.xn) * 100) : 100) + '%';
      $('xptext').textContent = me.xn > 0 ? `Erfahrung: ${me.xp} / ${me.xn}` : 'Höchststufe';

      // cast bar
      const cb = $('castbar');
      if (S.cast) {
        const k = (now() - S.cast.start) / S.cast.dur;
        if (k >= 1) S.cast = null;
        else {
          cb.classList.remove('hidden');
          cb.querySelector('.fill').style.width = (k * 100) + '%';
          cb.querySelector('span').textContent = S.defs.skills[S.cast.skill]?.name || '';
        }
      }
      if (!S.cast) cb.classList.add('hidden');

      // death
      $('deathbox').classList.toggle('hidden', !me.d);

      const gh = gameHours(), hh = Math.floor(gh), mm = Math.floor((gh - hh) * 60);
      const clock = (sunInfo().elev > -0.05 ? '☀ ' : '☾ ') + hh.toString().padStart(2, '0') + ':' + mm.toString().padStart(2, '0');
      if (clock !== lastClock) { lastClock = clock; $('clock').textContent = clock; }
    },
  };
}

// Quest tracker (right side under the minimap)
export function renderTracker() {
  const body = $('tracker-body');
  let h = '';
  for (const q of S.quests) {
    const def = S.defs.quests[q.id];
    if (!def) continue;
    h += `<div class="tq">${esc(def.title)}${q.status === 'complete' ? ' <span style="color:#7e7">(fertig)</span>' : ''}</div>`;
    if (q.end && q.status === 'active') h += `<div class="to qtimer" data-end="${q.end}">⏱</div>`;
    def.objectives.forEach((o, i) => {
      const n = q.counts[i] || 0;
      h += `<div class="to ${n >= o.count ? 'done' : ''}">- ${esc(o.text)}: ${n}/${o.count}</div>`;
    });
  }
  body.innerHTML = h || '<div class="to" style="color:#999">Keine aktiven Quests</div>';
}

export { CLASS_ICONS };
