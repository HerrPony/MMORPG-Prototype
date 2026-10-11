import * as THREE from 'three';
import { S } from './state.js';
import { $, esc, CLASS_COLORS, CLASS_ICONS, SKIN, HAIR, hex } from './util.js';
import { makePlayerModel } from './world/models.js';

// ---------- 3D character preview ----------
export class Preview {
  constructor() {
    this.scene = new THREE.Scene();
    this.scene.background = new THREE.Color(0x10202a);
    this.scene.fog = new THREE.Fog(0x10202a, 6, 22);
    this.camera = new THREE.PerspectiveCamera(35, 1, 0.1, 60);
    this.camera.position.set(0, 1.5, 4.6);
    this.camera.lookAt(0, 1.0, 0);
    this.scene.add(new THREE.HemisphereLight(0xbfd8ff, 0x2a3a22, 0.9));
    const key = new THREE.DirectionalLight(0xfff0d0, 1.6);
    key.position.set(3, 5, 4);
    this.scene.add(key);
    const rim = new THREE.DirectionalLight(0x6aa0ff, 0.8);
    rim.position.set(-4, 3, -3);
    this.scene.add(rim);
    const ground = new THREE.Mesh(new THREE.CircleGeometry(12, 40), new THREE.MeshLambertMaterial({ color: 0x27452a }));
    ground.rotation.x = -Math.PI / 2;
    this.scene.add(ground);
    const plate = new THREE.Mesh(new THREE.CylinderGeometry(1.1, 1.2, 0.1, 32), new THREE.MeshLambertMaterial({ color: 0x5a4a36 }));
    plate.position.y = -0.04;
    this.scene.add(plate);
    const trunk = new THREE.MeshLambertMaterial({ color: 0x4a3322 }), leaf = new THREE.MeshLambertMaterial({ color: 0x1f4a2a });
    for (let i = 0; i < 26; i++) {
      const a = (i / 26) * Math.PI * 2, d = 7 + (i % 3) * 2.2;
      const g = new THREE.Group();
      const t = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.3, 3, 6), trunk); t.position.y = 1.5; g.add(t);
      for (let k = 0; k < 3; k++) { const c = new THREE.Mesh(new THREE.ConeGeometry(1.8 - k * 0.45, 2.4, 7), leaf); c.position.y = 3.2 + k * 1.4; g.add(c); }
      g.position.set(Math.cos(a) * d, 0, Math.sin(a) * d - 2);
      this.scene.add(g);
    }
    this.model = null;
    this.yaw = 0.3;
    this.drag = false;
  }

  setStyle(cls, look) {
    if (this.model) this.scene.remove(this.model.obj);
    this.model = makePlayerModel(cls, look);
    this.model.obj.rotation.y = this.yaw;
    this.scene.add(this.model.obj);
  }

  clear() {
    if (this.model) this.scene.remove(this.model.obj);
    this.model = null;
  }

  update(dt, w, h) {
    this.camera.aspect = w / h;
    this.camera.updateProjectionMatrix();
    if (this.model) {
      if (!this.drag) this.yaw += dt * 0.5;
      this.model.obj.rotation.y = this.yaw;
      this._atk = (this._atk || 0) + dt; // show the attack animation now and then
      if (this._atk > 3.2) { this._atk = 0; this.model.attack?.(); }
      this.model.update(dt, { speed: 0, airborne: false, dead: false });
    }
  }
}

// ---------- screens ----------
const SCREENS = ['screen-login', 'screen-chars', 'screen-create'];
export function showScreen(id) {
  SCREENS.forEach((s) => $(s).classList.toggle('hidden', s !== id));
}

export function initScreens(net, preview, enterHandlers) {
  // --- login / register ---
  let mode = 'login';
  const setMode = (m) => {
    mode = m;
    $('tab-login').classList.toggle('active', m === 'login');
    $('tab-register').classList.toggle('active', m === 'register');
    $('row-pass2').classList.toggle('hidden', m !== 'register');
    $('btn-login').textContent = m === 'login' ? 'Anmelden' : 'Account erstellen';
    $('login-error').textContent = '';
  };
  $('tab-login').onclick = () => setMode('login');
  $('tab-register').onclick = () => setMode('register');
  const submit = () => {
    const user = $('in-user').value.trim(), pass = $('in-pass').value;
    if (mode === 'register' && pass !== $('in-pass2').value) { $('login-error').textContent = 'Die Passwörter stimmen nicht überein.'; return; }
    $('login-error').textContent = '';
    net.send(mode, { user, pass });
  };
  $('btn-login').onclick = submit;
  for (const id of ['in-user', 'in-pass', 'in-pass2']) $(id).addEventListener('keydown', (e) => { if (e.key === 'Enter') submit(); });
  net.on('authError', (m) => { $('login-error').textContent = m.error; });

  net.on('auth', (m) => {
    S.user = m.user;
    S.classes = m.classes; S.classOrder = m.classOrder; S.races = m.races; S.raceOrder = m.raceOrder; S.allSkills = m.skills;
    $('login-error').textContent = '';
    $('in-pass').value = $('in-pass2').value = '';
  });

  // --- character list ---
  let chars = [], selected = null, afterCreate = null;
  net.on('chars', (m) => {
    chars = m.chars;
    if (enterHandlers.isInWorld()) enterHandlers.leaveWorld();
    if (afterCreate) { selected = afterCreate; afterCreate = null; }
    if (!chars.find((c) => c.id === selected)) selected = chars[0]?.id || null;
    if (!chars.length) { openCreate(); return; }
    showScreen('screen-chars');
    renderChars();
  });
  function renderChars() {
    const list = $('char-list');
    list.innerHTML = '';
    for (const c of chars) {
      const d = document.createElement('div');
      d.className = 'char-item' + (c.id === selected ? ' sel' : '');
      const cl = S.classes[c.class];
      d.innerHTML = `<div class="ico" style="border-color:${CLASS_COLORS[c.class]}">${CLASS_ICONS[c.class]}</div><div><div class="nm" style="color:${CLASS_COLORS[c.class]}">${esc(c.name)}</div><div class="sub">Stufe ${c.level} ${esc(S.races[c.race].name)} ${esc(cl.name)}</div></div>`;
      d.onclick = () => { selected = c.id; renderChars(); };
      d.ondblclick = () => net.send('enterWorld', { id: c.id });
      list.appendChild(d);
    }
    const c = chars.find((x) => x.id === selected);
    $('btn-new-char').classList.toggle('hidden', chars.length >= 5);
    $('btn-enter').disabled = !c; $('btn-delete').disabled = !c;
    if (c) {
      $('cd-name').textContent = c.name;
      $('cd-name').style.color = CLASS_COLORS[c.class];
      $('cd-info').textContent = `Stufe ${c.level} ${S.races[c.race].name} ${S.classes[c.class].name}`;
      preview.setStyle(c.class, c.look);
    } else { $('cd-name').textContent = '–'; $('cd-info').textContent = ''; preview.clear(); }
  }
  $('btn-enter').onclick = () => { if (selected) net.send('enterWorld', { id: selected }); };
  $('btn-delete').onclick = () => {
    const c = chars.find((x) => x.id === selected);
    if (c && confirm(`Charakter "${c.name}" wirklich endgültig löschen?`)) net.send('deleteChar', { id: c.id });
  };
  $('btn-new-char').onclick = () => openCreate();
  $('btn-account-out').onclick = () => location.reload();

  // --- creation ---
  const form = { race: 'human', cls: 'warrior', gender: 'm', hair: 0, skin: 0 };
  function openCreate() {
    $('create-error').textContent = '';
    $('in-charname').value = '';
    $('btn-create-back').classList.toggle('hidden', chars.length === 0);
    showScreen('screen-create');
    renderCreate();
  }
  function renderCreate() {
    $('race-list').innerHTML = '';
    for (const id of S.raceOrder) {
      const r = S.races[id];
      const d = document.createElement('div');
      d.className = 'card' + (form.race === id ? ' sel' : '');
      d.innerHTML = `<div class="big-ico">🧑</div><div class="cn">${esc(r.name)}</div><div class="cr">${esc(r.desc)}</div>`;
      d.onclick = () => { form.race = id; renderCreate(); };
      $('race-list').appendChild(d);
    }
    $('class-list').innerHTML = '';
    for (const id of S.classOrder) {
      const c = S.classes[id];
      const d = document.createElement('div');
      d.className = 'card' + (form.cls === id ? ' sel' : '');
      d.innerHTML = `<div class="big-ico">${CLASS_ICONS[id]}</div><div class="cn" style="color:${CLASS_COLORS[id]}">${esc(c.name)}</div><div class="cr">${esc(c.role)}</div>`;
      d.onclick = () => { form.cls = id; renderCreate(); };
      $('class-list').appendChild(d);
    }
    const sw = (arr, key, host) => {
      $(host).innerHTML = '';
      arr.forEach((col, i) => {
        const s = document.createElement('div');
        s.className = 'swatch' + (form[key] === i ? ' sel' : '');
        s.style.background = hex(col);
        s.onclick = () => { form[key] = i; renderCreate(); };
        $(host).appendChild(s);
      });
    };
    sw(HAIR, 'hair', 'hair-list'); sw(SKIN, 'skin', 'skin-list');
    $('g-m').classList.toggle('active', form.gender === 'm');
    $('g-f').classList.toggle('active', form.gender === 'f');
    // class detail
    const c = S.classes[form.cls];
    $('class-detail').innerHTML = `<h2 style="color:${CLASS_COLORS[form.cls]}">${esc(c.name)}</h2><div class="role">${esc(c.role)} · Ressource: ${{ rage: 'Wut', focus: 'Fokus', mana: 'Mana' }[c.resource]}</div>` +
      `<p style="font-size:13px;color:#cdbf9f">${esc(c.desc)}</p><h3>Fertigkeiten</h3>` +
      c.skills.map((id) => { const s = S.allSkills[id]; return `<div class="skill"><div class="si">${s.icon}</div><div><b>${esc(s.name)}</b><span class="lvl">${s.level > 1 ? `Stufe ${s.level} · beim Lehrer` : 'von Beginn an'}</span><div>${esc(s.desc)}</div></div></div>`; }).join('');
    preview.setStyle(form.cls, { gender: form.gender, hair: form.hair, skin: form.skin });
  }
  $('g-m').onclick = () => { form.gender = 'm'; renderCreate(); };
  $('g-f').onclick = () => { form.gender = 'f'; renderCreate(); };
  $('btn-create-back').onclick = () => { showScreen('screen-chars'); renderChars(); };
  $('btn-create').onclick = () => {
    net.send('createChar', { name: $('in-charname').value.trim(), race: form.race, class: form.cls, look: { gender: form.gender, hair: form.hair, skin: form.skin } });
  };
  $('in-charname').addEventListener('keydown', (e) => { e.stopPropagation(); if (e.key === 'Enter') $('btn-create').click(); });
  net.on('createError', (m) => { $('create-error').textContent = m.error; });
  net.on('created', (m) => { afterCreate = m.id; });
  net.on('enterError', (m) => { alert(m.error); });
}
