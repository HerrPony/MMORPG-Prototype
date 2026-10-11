import * as THREE from 'three';
import { S } from './state.js';
import { heightAt } from './util.js';

const SLOT_CODES = ['Digit1', 'Digit2', 'Digit3', 'Digit4', 'Digit5', 'Digit6', 'Digit7', 'Digit8', 'Digit9', 'Digit0', 'Minus', 'Equal', 'BracketLeft', 'BracketRight'];

const RUN = 7, BACK = 4.2, WALK = 3.2, JUMP = 8.2, GRAVITY = 24;

// Local player: movement, camera, mouse + keyboard controls.
export class Controller {
  constructor({ canvas, camera, world, entities, ui }) {
    this.canvas = canvas;
    this.camera = camera;
    this.world = world;
    this.entities = entities;
    this.ui = ui;
    this.keys = {};
    this.camYaw = 0;
    this.camPitch = 0.35;
    this.camDist = 9;
    this.vy = 0;
    this.grounded = true;
    this.autorun = false;
    this.rmb = false; this.lmb = false;
    this.down = { x: 0, y: 0, t: 0, acc: 0, btn: -1 };
    this.lastSend = 0; this.lastSent = { x: 0, y: 0, z: 0, r: 0 };
    this.ray = new THREE.Raycaster();
    this.hoverT = 0;
    this.bind();
  }

  // Place the player and align the camera behind them.
  spawn(x, y, z, r) {
    Object.assign(S.local, { x, y, z, r, speed: 0, airborne: false });
    this.camYaw = r;
  }

  bind() {
    const c = this.canvas;
    addEventListener('keydown', (e) => this.onKey(e, true));
    addEventListener('keyup', (e) => this.onKey(e, false));
    addEventListener('blur', () => { this.keys = {}; this.rmb = this.lmb = false; });
    c.addEventListener('contextmenu', (e) => e.preventDefault());
    c.addEventListener('mousedown', (e) => this.onDown(e));
    addEventListener('mouseup', (e) => this.onUp(e));
    addEventListener('mousemove', (e) => this.onMove(e));
    c.addEventListener('wheel', (e) => {
      e.preventDefault();
      this.camDist = Math.max(2, Math.min(22, this.camDist * (1 + Math.sign(e.deltaY) * 0.1)));
    }, { passive: false });
  }

  // ---------- keyboard ----------
  onKey(e, down) {
    if (!S.inWorld) return;
    if (S.typing) return; // chat handles its own keys
    const k = e.key.length === 1 ? e.key.toLowerCase() : e.key;
    if (e.repeat && down) { if (k === 'Tab' || k === ' ') e.preventDefault(); return; }
    const move = { w: 'w', a: 'a', s: 's', d: 'd', ArrowUp: 'w', ArrowDown: 's', ArrowLeft: 'a', ArrowRight: 'd', ' ': 'jump', Shift: 'shift' }[k];
    if (move) { this.keys[move] = down; if (down) { e.preventDefault(); if (move === 'w' || move === 's') this.autorun = false; } }
    if (!down) return;
    if (e.ctrlKey || e.metaKey || e.altKey) return;
    switch (k) {
      case 'Tab': e.preventDefault(); this.targetNext(e.shiftKey); break;
      case 'Escape': this.ui.escape(); break;
      case 'Enter': e.preventDefault(); this.ui.openChat(''); break;
      case '/': e.preventDefault(); this.ui.openChat('/'); break;
      case 'b': this.ui.toggle('bags'); break;
      case 'l': this.ui.toggle('questlog'); break;
      case 'c': this.ui.toggle('charwin'); break;
      case 'p': this.ui.toggle('spellbook'); break;
      case 'm': this.ui.toggle('map'); break;
      case 'o': this.ui.toggle('party'); break;
      case 'r': this.autorun = !this.autorun; break;
      case 'f': this.interactTarget(); break;
      default: {
        const i = SLOT_CODES.indexOf(e.code);
        if (i >= 0) this.ui.castSlot(i);
      }
    }
  }

  // ---------- mouse ----------
  onDown(e) {
    if (!S.inWorld) return;
    if (e.button === 0) this.lmb = true;
    if (e.button === 2) this.rmb = true;
    this.down = { x: e.clientX, y: e.clientY, t: performance.now(), acc: 0, btn: e.button };
    if (e.button === 2 && !S.me.d) {
      try { const p = this.canvas.requestPointerLock?.(); if (p && p.catch) p.catch(() => {}); } catch { /* ignore */ }
    }
    this.ui.closeTooltip?.();
  }

  onUp(e) {
    if (!S.inWorld) return;
    const d = this.down;
    if (e.button === 0 && this.lmb) {
      this.lmb = false;
      if (d.btn === 0 && d.acc < 5) this.click(d.x, d.y, false);
    }
    if (e.button === 2 && this.rmb) {
      this.rmb = false;
      if (document.pointerLockElement) document.exitPointerLock();
      if (d.btn === 2 && d.acc < 6 && performance.now() - d.t < 500) this.click(d.x, d.y, true);
    }
  }

  onMove(e) {
    if (!S.inWorld) return;
    const dragging = this.lmb || this.rmb;
    if (dragging) {
      this.down.acc += Math.abs(e.movementX) + Math.abs(e.movementY);
      if (this.rmb || this.down.acc > 4) {
        this.camYaw -= e.movementX * 0.0045;
        this.camPitch = Math.max(-0.35, Math.min(1.45, this.camPitch + e.movementY * 0.0045));
      }
      return;
    }
    const t = performance.now();
    if (t - this.hoverT > 80 && e.target === this.canvas) {
      this.hoverT = t;
      const id = this.pick(e.clientX, e.clientY);
      const ent = S.ents.get(id);
      let cur = 'default';
      if (ent) {
        if (ent.k === 'npc' || ent.k === 'obj' || (ent.k === 'mob' && ent.d && ent.lt)) cur = 'pointer';
        else if (ent.k === 'mob' && !ent.d) cur = 'crosshair';
      }
      this.canvas.style.cursor = cur;
    }
  }

  pick(cx, cy) {
    const v = new THREE.Vector2((cx / innerWidth) * 2 - 1, -(cy / innerHeight) * 2 + 1);
    this.ray.setFromCamera(v, this.camera);
    const hits = this.ray.intersectObjects(this.entities.pickables, false);
    return hits.length ? hits[0].object.userData.entId : 0;
  }

  click(x, y, right) {
    const id = this.pick(x, y);
    if (!id) return;
    const e = S.ents.get(id);
    if (!e) return;
    if (S.me.tg !== id) this.setTarget(id);
    if (!right) return;
    if (e.k === 'npc' || e.k === 'obj' || (e.k === 'mob' && e.d)) {
      S.net.send('interact', { id });
    } else if (e.k === 'player' && e.id !== S.id) {
      this.ui.playerMenu?.(e, x, y);
    } else if (e.k === 'pet') {
      // pets cannot be attacked
    } else if (e.k === 'mob' && !e.d) {
      this.ui.castSlot(0); // basic attack, starts auto-attack
    }
  }

  interactTarget() {
    const e = S.ents.get(S.me.tg);
    if (e && (e.k === 'npc' || e.k === 'obj' || (e.k === 'mob' && e.d))) S.net.send('interact', { id: e.id });
  }

  setTarget(id) {
    S.me.tg = id;
    S.net.send('target', { id });
  }

  targetNext(back) {
    const me = S.local;
    const list = [...S.ents.values()]
      .filter((e) => e.k === 'mob' && !e.d && Math.hypot(e.rx - me.x, e.rz - me.z) < 45)
      .sort((a, b) => Math.hypot(a.rx - me.x, a.rz - me.z) - Math.hypot(b.rx - me.x, b.rz - me.z));
    if (!list.length) return;
    let i = list.findIndex((e) => e.id === S.me.tg);
    i = back ? (i <= 0 ? list.length - 1 : i - 1) : (i + 1) % list.length;
    this.setTarget(list[i].id);
  }

  faceTarget() {
    const t = S.ents.get(S.me.tg);
    if (t && t.k === 'mob') S.local.r = Math.atan2(t.rx - S.local.x, t.rz - S.local.z);
  }

  // ---------- per frame ----------
  update(dt) {
    const L = S.local;
    const k = this.keys;
    if (S.me.fl) { // griffin flight: the server moves us, only the camera is ours
      this.updateCamera();
      return;
    }
    const dead = S.me.d;
    const both = this.lmb && this.rmb;
    let fw = (k.w ? 1 : 0) - (k.s ? 1 : 0) + (both || this.autorun ? 1 : 0);
    let rt = (k.d ? 1 : 0) - (k.a ? 1 : 0);
    if (dead || S.typing) { fw = 0; rt = 0; }
    fw = Math.max(-1, Math.min(1, fw));

    const fx = Math.sin(this.camYaw), fz = Math.cos(this.camYaw);
    let dx = fx * fw + -fz * rt, dz = fz * fw + fx * rt;
    const len = Math.hypot(dx, dz);
    let speed = 0;
    if (len > 0.001) {
      dx /= len; dz /= len;
      speed = fw < 0 && rt === 0 ? BACK : k.shift ? WALK : RUN;
      if (fw < 0 && rt !== 0) speed = RUN * 0.8;
    }
    // facing
    if (!dead) {
      // with the camera held (RMB) the character normally faces the camera direction; strafing forward/sideways
      // (A/D, W+A, W+D) still turns it smoothly towards the direction it actually runs in
      const strafe = this.rmb && len > 0.001 && rt !== 0 && fw >= 0;
      if (this.rmb && !strafe) L.r = this.camYaw;
      else if (len > 0.001) {
        const want = Math.atan2(dx, dz);
        let d = ((want - L.r + Math.PI) % (Math.PI * 2) + Math.PI * 2) % (Math.PI * 2) - Math.PI;
        L.r += d * Math.min(1, dt * 14);
      }
    }
    let nx = L.x + dx * speed * dt, nz = L.z + dz * speed * dt;
    if (speed > 0) [nx, nz] = this.dungeon ? this.dungeon.collide(nx, nz) : this.world.collide(nx, nz, L.x, L.z);
    const moved = Math.hypot(nx - L.x, nz - L.z) / Math.max(dt, 1e-4);
    L.x = nx; L.z = nz;
    L.speed = moved / RUN;

    // vertical
    const ground = heightAt(L.x, L.z);
    if (this.grounded) {
      if (k.jump && !dead && !S.typing) { this.vy = JUMP; this.grounded = false; }
      else if (L.y - ground > 0.35) { this.grounded = false; this.vy = 0; }
      else L.y = ground;
    }
    if (!this.grounded) {
      this.vy -= GRAVITY * dt;
      L.y += this.vy * dt;
      if (L.y <= ground && this.vy <= 0) { L.y = ground; this.vy = 0; this.grounded = true; }
    }
    L.airborne = !this.grounded;

    this.sendMove();
    this.updateCamera();
  }

  sendMove() {
    const L = S.local, t = performance.now();
    if (t - this.lastSend < 50) return;
    const s = this.lastSent;
    const changed = Math.abs(L.x - s.x) > 0.01 || Math.abs(L.z - s.z) > 0.01 || Math.abs(L.y - s.y) > 0.01 || Math.abs(L.r - s.r) > 0.02;
    if (!changed && t - this.lastSend < 600) return;
    this.lastSend = t;
    Object.assign(s, { x: L.x, y: L.y, z: L.z, r: L.r });
    S.net.send('move', { x: +L.x.toFixed(2), y: +L.y.toFixed(2), z: +L.z.toFixed(2), r: +L.r.toFixed(3) });
  }

  setDungeon(d) { this.dungeon = d; }

  correct(x, y, z) {
    Object.assign(S.local, { x, y, z });
    this.vy = 0;
    this.grounded = true;
  }

  updateCamera() {
    const L = S.local;
    const tx = L.x, ty = L.y + 1.7, tz = L.z;
    const cp = Math.cos(this.camPitch), sp = Math.sin(this.camPitch);
    let cx = tx - Math.sin(this.camYaw) * cp * this.camDist;
    let cz = tz - Math.cos(this.camYaw) * cp * this.camDist;
    let cy = ty + sp * this.camDist;
    if (this.dungeon) [cx, cz] = this.dungeon.collide(cx, cz); // keep the camera inside the walls
    const minY = heightAt(cx, cz) + 0.5;
    if (cy < minY) cy = minY;
    this.camera.position.set(cx, cy, cz);
    this.camera.lookAt(tx, ty + 0.1, tz);
    if (this.dungeon) this.dungeon.update(performance.now() / 1000, L.x, L.z);
    else this.world.followSun(tx, L.y, tz);
  }
}
