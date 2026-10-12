import { S } from './state.js';

// Touch controls for phones/tablets in landscape: floating joystick (left), camera drag + tap-to-target (right),
// action buttons, fullscreen + landscape lock.
const isTouch = matchMedia('(pointer: coarse)').matches || (navigator.maxTouchPoints > 0 && !matchMedia('(hover: hover)').matches);
const STICK = 52; // joystick radius in px

export function initTouch(getCtl, ui) {
  if (!isTouch) return;
  document.body.classList.add('touch');
  const canvas = document.getElementById('gl');
  const mk = (cls, html, parent = document.body) => {
    const el = document.createElement('div');
    el.className = cls; el.innerHTML = html; parent.appendChild(el); return el;
  };

  // ---- fullscreen / landscape ----
  const goFullscreen = async () => {
    if (!document.fullscreenElement) { try { await document.documentElement.requestFullscreen({ navigationUI: 'hide' }); } catch { /* iOS Safari: not supported */ } }
    try { await screen.orientation?.lock?.('landscape'); } catch { /* needs fullscreen on Android */ }
  };
  document.getElementById('rotate-fs').onclick = goFullscreen;
  let firstTouch = true;
  addEventListener('touchend', () => { if (firstTouch) { firstTouch = false; goFullscreen(); } }, { passive: true });
  addEventListener('touchend', () => setTimeout(() => ui.closeTooltip?.(), 2500), { passive: true });

  // ---- HUD controls ----
  const root = mk('', '', document.body);
  root.id = 'touch-ui';
  const stick = mk('', '<div class="knob"></div>', root);
  stick.id = 'joy';
  const knob = stick.firstChild;
  const btn = (id, label, cls = '') => { const b = mk('tbtn ' + cls, label, root); b.id = id; return b; };
  const bAtk = btn('t-atk', '⚔️', 'big'), bJump = btn('t-jump', '⤒'), bTgt = btn('t-tgt', '🎯'), bUse = btn('t-use', '✋'), bChat = btn('t-chat', '💬');
  const hold = (el, down, up) => {
    el.addEventListener('touchstart', (e) => { e.preventDefault(); el.classList.add('on'); down(); }, { passive: false });
    const end = (e) => { e.preventDefault(); el.classList.remove('on'); up?.(); };
    el.addEventListener('touchend', end, { passive: false });
    el.addEventListener('touchcancel', end, { passive: false });
  };
  hold(bJump, () => { const c = getCtl(); if (c) c.keys.jump = true; }, () => { const c = getCtl(); if (c) c.keys.jump = false; });
  hold(bAtk, () => {
    const c = getCtl(); if (!c) return;
    const t = S.ents.get(S.me.tg);
    if (!t || t.k !== 'mob' || t.d) c.targetNext(false);
    ui.castSlot(0);
  });
  hold(bTgt, () => getCtl()?.targetNext(false));
  hold(bUse, () => {
    const c = getCtl(); if (!c) return;
    let t = S.ents.get(S.me.tg);
    const usable = (e) => e && (e.k === 'npc' || e.k === 'obj' || (e.k === 'mob' && e.d && e.lt));
    if (!usable(t)) { // pick the nearest usable thing
      let best = null, bd = 8;
      for (const e of S.ents.values()) {
        if (!usable(e)) continue;
        const d = Math.hypot(e.rx - S.local.x, e.rz - S.local.z);
        if (d < bd) { bd = d; best = e; }
      }
      if (!best) return;
      c.setTarget(best.id);
    }
    c.interactTarget();
  });
  hold(bChat, () => {
    const chat = document.getElementById('chat');
    const open = chat.classList.toggle('open');
    if (open) S.openChat?.('');
  });

  // ---- canvas gestures ----
  let joy = null;           // { id, ox, oy }
  let cam = null;           // { id, x, y, sx, sy, t, moved }
  let pinch = null;         // { a, b, dist }
  const byId = (e, id) => [...e.touches].find((t) => t.identifier === id);
  const setJoy = (x, y) => { const c = getCtl(); if (c) { c.joy.x = x; c.joy.y = y; } };
  const showStick = (on) => stick.classList.toggle('active', on);

  canvas.addEventListener('touchstart', (e) => {
    if (!S.inWorld) return;
    e.preventDefault();
    ui.closeTooltip?.();
    for (const t of e.changedTouches) {
      if (!joy && t.clientX < innerWidth * 0.42 && t.clientY > innerHeight * 0.3) {
        joy = { id: t.identifier, ox: t.clientX, oy: t.clientY };
        stick.style.left = (t.clientX - 70) + 'px'; stick.style.top = (t.clientY - 70) + 'px';
        knob.style.transform = 'translate(0,0)';
        showStick(true);
      } else if (!cam) {
        cam = { id: t.identifier, x: t.clientX, y: t.clientY, sx: t.clientX, sy: t.clientY, t: performance.now(), moved: 0 };
      } else if (!pinch) {
        const a = byId(e, cam.id);
        if (a) { pinch = { id2: t.identifier, dist: Math.hypot(a.clientX - t.clientX, a.clientY - t.clientY) }; cam.moved = 99; }
      }
    }
  }, { passive: false });

  canvas.addEventListener('touchmove', (e) => {
    if (!S.inWorld) return;
    e.preventDefault();
    const c = getCtl(); if (!c) return;
    if (joy) {
      const t = byId(e, joy.id);
      if (t) {
        let dx = t.clientX - joy.ox, dy = t.clientY - joy.oy;
        const len = Math.hypot(dx, dy);
        if (len > STICK) { dx = dx / len * STICK; dy = dy / len * STICK; }
        knob.style.transform = `translate(${dx}px,${dy}px)`;
        const m = Math.min(1, len / STICK);
        if (m < 0.18) setJoy(0, 0); else setJoy(dx / Math.max(len, 1e-3) * m, -dy / Math.max(len, 1e-3) * m);
      }
    }
    if (cam) {
      const t = byId(e, cam.id);
      if (t) {
        const dx = t.clientX - cam.x, dy = t.clientY - cam.y;
        cam.x = t.clientX; cam.y = t.clientY;
        cam.moved += Math.abs(dx) + Math.abs(dy);
        if (cam.moved > 10) {
          c.camYaw -= dx * 0.006;
          c.camPitch = Math.max(-0.35, Math.min(1.45, c.camPitch + dy * 0.006));
        }
        if (pinch) {
          const t2 = byId(e, pinch.id2);
          if (t2) {
            const d = Math.hypot(t.clientX - t2.clientX, t.clientY - t2.clientY);
            c.camDist = Math.max(2, Math.min(22, c.camDist * (pinch.dist / Math.max(d, 1))));
            pinch.dist = d;
          }
        }
      }
    }
  }, { passive: false });

  const end = (e) => {
    e.preventDefault();
    for (const t of e.changedTouches) {
      if (joy && t.identifier === joy.id) { joy = null; setJoy(0, 0); showStick(false); }
      else if (pinch && t.identifier === pinch.id2) pinch = null;
      else if (cam && t.identifier === cam.id) {
        const tap = cam.moved < 10 && performance.now() - cam.t < 350;
        if (tap) onTap(t.clientX, t.clientY);
        cam = null; pinch = null;
      }
    }
  };
  canvas.addEventListener('touchend', end, { passive: false });
  canvas.addEventListener('touchcancel', end, { passive: false });

  // tap: NPCs, objects and lootable corpses are used directly; hostile mobs / players need a second tap
  function onTap(x, y) {
    const c = getCtl(); if (!c) return;
    const id = c.pick(x, y);
    if (!id) return;
    const e = S.ents.get(id);
    if (!e) return;
    const direct = e.k === 'npc' || e.k === 'obj' || (e.k === 'mob' && e.d);
    c.click(x, y, direct || S.me.tg === id);
  }

  // long press on item slots = right click (use / equip / unequip)
  let lp = null;
  document.addEventListener('touchstart', (e) => {
    const el = e.target.closest?.('.bslot, .eslot');
    if (!el) return;
    const t = e.touches[0];
    lp = setTimeout(() => { el.dispatchEvent(new MouseEvent('contextmenu', { bubbles: true, cancelable: true, clientX: t.clientX, clientY: t.clientY })); lp = null; }, 550);
  }, { passive: true });
  const cancelLp = () => { if (lp) { clearTimeout(lp); lp = null; } };
  document.addEventListener('touchend', cancelLp, { passive: true });
  document.addEventListener('touchmove', cancelLp, { passive: true });

  // joystick must not stay stuck if the window loses focus
  addEventListener('blur', () => { joy = null; cam = null; setJoy(0, 0); showStick(false); });
  addEventListener('contextmenu', (e) => e.preventDefault());
}
