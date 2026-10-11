import * as THREE from 'three';
import { S } from '../state.js';
import { $, angleLerp, CLASS_COLORS, esc } from '../util.js';
import { makeModel, applyEquipment, makeGriffin } from './models.js';

function glowTexture() {
  const c = document.createElement('canvas');
  c.width = c.height = 64;
  const g = c.getContext('2d');
  const grd = g.createRadialGradient(32, 32, 0, 32, 32, 32);
  grd.addColorStop(0, 'rgba(255,255,255,1)');
  grd.addColorStop(0.25, 'rgba(255,235,150,0.8)');
  grd.addColorStop(1, 'rgba(255,200,60,0)');
  g.fillStyle = grd;
  g.fillRect(0, 0, 64, 64);
  return new THREE.CanvasTexture(c);
}

function markTexture(ch, color) {
  const c = document.createElement('canvas');
  c.width = 64; c.height = 96;
  const g = c.getContext('2d');
  g.font = 'bold 84px Georgia, serif';
  g.textAlign = 'center'; g.textBaseline = 'middle';
  g.lineWidth = 9; g.strokeStyle = '#2a1500'; g.strokeText(ch, 32, 52);
  g.fillStyle = color; g.fillText(ch, 32, 52);
  const t = new THREE.CanvasTexture(c);
  t.colorSpace = THREE.SRGBColorSpace;
  return t;
}

export class Entities {
  constructor(scene) {
    this.scene = scene;
    this.glow = glowTexture();
    this.markTex = {
      available: markTexture('!', '#ffd01a'),
      ready: markTexture('?', '#ffd01a'),
      progress: markTexture('?', '#9a9a9a'),
    };
    this.ring = new THREE.Mesh(
      new THREE.RingGeometry(0.75, 0.9, 40),
      new THREE.MeshBasicMaterial({ color: 0xffff00, transparent: true, opacity: 0.85, depthTest: false, side: THREE.DoubleSide }),
    );
    this.ring.rotation.x = -Math.PI / 2;
    this.ring.renderOrder = 10;
    this.ring.visible = false;
    scene.add(this.ring);
    this.pickables = [];
    this.labels = $('labels');
    this.v = new THREE.Vector3();
  }

  // Replaces the scene the entities live in (open world <-> dungeon).
  setScene(scene) {
    this.scene = scene;
    scene.add(this.ring);
    for (const e of S.ents.values()) { e.group.removeFromParent(); scene.add(e.group); if (e.sparkle) { scene.add(e.sparkle.points); scene.add(e.sparkle.glow); } }
  }

  // ---- snapshot ----
  apply(list, lootable) {
    const seen = new Set();
    const lt = new Set(lootable || []);
    for (const v of list) {
      seen.add(v.id);
      let e = S.ents.get(v.id);
      if (!e) e = this.create(v);
      const wasDead = e.d, wasLt = e.lt;
      const modelChanged = v.m !== e.m;
      const eqKey = JSON.stringify(v.eq || null);
      // optional fields are omitted by the server when empty: reset them first
      Object.assign(e, { fl: false, d: false, tg: 0, st: undefined, el: false, bs: false, ow: 0, q: '', gn: false, eq: undefined, pt: false, ti: '', c: undefined }, v);
      e.d = !!v.d;
      e.lt = lt.has(v.id);
      e.tx = v.x; e.ty = v.y; e.tz = v.z; e.tr = v.r;
      if (modelChanged) this.rebuildModel(e);
      if (e.k === 'player' && eqKey !== e._eqKey) { e._eqKey = eqKey; applyEquipment(e.model, e); }
      this.syncMount(e);
      if (e.id === S.id && !e.init) { e.init = true; S.local.x = v.x; S.local.y = v.y; S.local.z = v.z; S.local.r = v.r; }
      if (e.lt !== wasLt) this.setSparkle(e, e.lt);
      if (e.k === 'npc') this.refreshMarker(e);
      if (e.k === 'obj') this.refreshObject(e);
      if (e.label) this.refreshLabel(e);
      if (e.d && !wasDead && e.k === 'mob') e.group.userData.justDied = true;
    }
    for (const [id, e] of S.ents) {
      if (!seen.has(id)) this.remove(e);
    }
  }

  create(v) {
    const model = makeModel(v);
    const group = new THREE.Group();
    group.add(model.obj);
    group.position.set(v.x, v.y, v.z);
    group.rotation.y = v.r;
    const sc = v.sc || 1;
    group.scale.setScalar(sc);
    const hitR = model.kind === 'ship' ? 6 : 0.7;
    const hit = new THREE.Mesh(
      new THREE.CylinderGeometry(hitR, hitR, model.height, 8),
      new THREE.MeshBasicMaterial({ visible: false }),
    );
    hit.position.y = model.height / 2;
    hit.userData.entId = v.id;
    group.add(hit);
    this.scene.add(group);
    this.pickables.push(hit);

    const label = document.createElement('div');
    label.className = 'nameplate';
    this.labels.appendChild(label);

    const e = {
      id: v.id, model, group, hit, label, init: false, sc,
      rx: v.x, ry: v.y, rz: v.z, rr: v.r, speed: 0, airborne: false,
    };
    S.ents.set(v.id, e);
    if (v.m === 'flightmaster') { // a griffin waits next to the griffin master
      const gr = makeGriffin();
      gr.obj.position.set(2.8, 0, 0.5);
      gr.obj.rotation.y = -0.6;
      group.add(gr.obj);
      e.pad = gr;
    }
    return e;
  }

  rebuildModel(e) {
    e.model.obj.removeFromParent();
    e.model = makeModel(e);
    e.group.add(e.model.obj);
    this.pickables.splice(this.pickables.indexOf(e.hit), 1);
    e.hit.geometry.dispose();
    e.hit.geometry = new THREE.CylinderGeometry(0.7, 0.7, e.model.height, 8);
    e.hit.position.y = e.model.height / 2;
    this.pickables.push(e.hit);
    e._lk = null;
  }

  // Players riding a griffin during a flight.
  syncMount(e) {
    if (e.k !== 'player') return;
    if (e.fl && !e.mount) {
      e.mount = makeGriffin();
      e.group.add(e.mount.obj);
      e.model.obj.position.set(e.mount.seatX || 0, e.mount.seatY, e.mount.seatZ);
    } else if (!e.fl && e.mount) {
      e.mount.obj.removeFromParent();
      e.mount = null;
      e.model.obj.position.set(0, 0, 0);
    }
  }

  remove(e) {
    e.group.removeFromParent();
    const i = this.pickables.indexOf(e.hit);
    if (i >= 0) this.pickables.splice(i, 1);
    e.label.remove();
    if (e.sparkle) { this.scene.remove(e.sparkle.points); this.scene.remove(e.sparkle.glow); }
    S.ents.delete(e.id);
    if (S.me.tg === e.id) S.me.tg = 0;
  }

  removeAll() { for (const e of [...S.ents.values()]) this.remove(e); }

  refreshLabel(e) {
    let color = '#fff', sub = '', cls = 'nameplate';
    if (e.k === 'player') color = CLASS_COLORS[e.c] || '#fff';
    else if (e.k === 'npc') { color = '#ffe27a'; sub = `<small>&lt;${esc(e.ti || '')}&gt;</small>`; }
    else if (e.k === 'pet') { color = '#9fe89f'; sub = `<small>Begleiter</small>`; }
    else if (e.k === 'ship') { color = '#9fd8ff'; sub = `<small>${esc(e.ti || '')}</small>`; }
    else if (e.k === 'obj') { color = e.q && S.quests.some((q) => q.status === 'active' && this.questNeeds(q.id, e)) ? '#ffe14a' : '#cccccc'; }
    else { color = e.d ? '#888' : '#ffe14a'; if (e.el) { cls += ' elite'; sub = `<small class="tag">${e.bs ? 'Boss' : 'Elite'}</small>`; } }
    const isT = S.me.tg === e.id;
    let hp = '';
    if (isT && !e.d && (e.k === 'mob' || e.k === 'pet')) hp = `<div class="hpb"><i style="width:${Math.max(0, e.hp / e.mh * 100)}%"></i></div>`;
    const key = color + e.n + sub + hp + e.l + cls + e.q;
    if (e._lk === key) return;
    e._lk = key;
    e.label.className = cls;
    const lvl = (e.k === 'mob' && !e.d) || e.k === 'pet' ? ` <small style="display:inline;color:${e.l > S.me.lvl + 4 ? '#ff5a4a' : '#ccc'}">${e.l}</small>` : '';
    e.label.innerHTML = `<span style="color:${color}">${esc(e.n)}</span>${lvl}${sub}${hp}`;
  }

  questNeeds(qid) { return !!S.defs.quests[qid]; }

  refreshMarker(e) {
    const st = S.marks[e.id] || '';
    if (e.markState === st) return;
    e.markState = st;
    if (e.marker) { e.group.remove(e.marker); e.marker = null; }
    if (!st) return;
    const sp = new THREE.Sprite(new THREE.SpriteMaterial({ map: this.markTex[st], transparent: true, depthWrite: false }));
    sp.scale.set(0.75, 1.1, 1);
    sp.position.y = e.model.height + 0.75;
    e.marker = sp;
    e.group.add(sp);
  }

  refreshAllMarkers() {
    for (const e of S.ents.values()) if (e.k === 'npc') this.refreshMarker(e);
  }

  // quest objects sparkle while they are useful to the player
  refreshObject(e) {
    const want = !!e.q && S.quests.some((q) => q.id === e.q && q.status === 'active');
    if (!!e.objGlow === want) return;
    e.objGlow = want;
    this.setSparkle(e, want || e.lt);
  }

  setSparkle(e, on) {
    if (!on) {
      if (e.sparkle) { this.scene.remove(e.sparkle.points); this.scene.remove(e.sparkle.glow); e.sparkle = null; }
      return;
    }
    if (e.sparkle) return;
    const n = 18;
    const geo = new THREE.BufferGeometry();
    const arr = new Float32Array(n * 3);
    geo.setAttribute('position', new THREE.BufferAttribute(arr, 3));
    const pts = new THREE.Points(geo, new THREE.PointsMaterial({
      map: this.glow, size: 0.45, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending, color: 0xffeeaa,
    }));
    pts.frustumCulled = false;
    const glow = new THREE.Sprite(new THREE.SpriteMaterial({ map: this.glow, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending, opacity: 0.6 }));
    glow.scale.set(2.2, 2.2, 1);
    this.scene.add(pts, glow);
    e.sparkle = { points: pts, glow, seeds: Array.from({ length: n }, (_, i) => ({ a: Math.random() * 6.28, r: 0.2 + Math.random() * 0.8, t: Math.random() * 2, s: 0.4 + Math.random() * 0.6, i })) };
  }

  // ---- animation / rendering ----
  attack(id) { const e = S.ents.get(id); if (e) e.model.attack(); }
  castAnim(id, on) { const e = S.ents.get(id); if (e && e.model.cast) e.model.cast(on); }

  update(dt, camera, w, h) {
    const k = 1 - Math.exp(-dt * 11);
    const tNow = performance.now() / 1000;
    for (const e of S.ents.values()) {
      const self = e.id === S.id;
      if ((S.me.tg === e.id) !== !!e._wasT) { e._wasT = S.me.tg === e.id; this.refreshLabel(e); }
      if (self && !S.me.fl) {
        e.speed = S.local.speed; e.airborne = S.local.airborne;
        e.rx = S.local.x; e.ry = S.local.y; e.rz = S.local.z; e.rr = angleLerp(e.rr, S.local.r, Math.min(1, dt * 14));
      } else {
        const ox = e.rx, oz = e.rz;
        e.rx += (e.tx - e.rx) * k; e.rz += (e.tz - e.rz) * k; e.ry += (e.ty - e.ry) * k;
        const sp = Math.hypot(e.rx - ox, e.rz - oz) / Math.max(dt, 1e-3);
        e.speed += (sp / (e.k === 'mob' ? 5 : 7) - e.speed) * Math.min(1, dt * 8);
        e.airborne = e.ty - e.ry > 0.6 && !e.fl;
        e.rr = angleLerp(e.rr, e.tr, Math.min(1, dt * 10));
        if (self && S.me.fl) { S.local.x = e.rx; S.local.y = e.ry; S.local.z = e.rz; S.local.r = e.rr; S.local.speed = 0; }
      }
      e.group.position.set(e.rx, e.ry, e.rz);
      e.group.rotation.y = e.rr;
      const st = { speed: e.fl ? 0 : e.speed, airborne: e.airborne, dead: e.d, flying: e.fl };
      e.model.update(dt, st);
      if (e.mount) e.mount.update(dt, st);
      if (e.pad) e.pad.update(dt, { speed: 0, flying: false });
      if (e.marker) e.marker.position.y = e.model.height + 0.75 + Math.sin(tNow * 3) * 0.08;

      if (e.sparkle) {
        const sp = e.sparkle, pos = sp.points.geometry.attributes.position;
        sp.seeds.forEach((s, i) => {
          s.t += dt * s.s;
          if (s.t > 1.6) { s.t = 0; s.a = Math.random() * 6.28; s.r = 0.2 + Math.random() * 0.8; }
          const rr = s.r * (1 - s.t * 0.3);
          pos.setXYZ(i, e.rx + Math.cos(s.a + s.t * 2) * rr, e.ry + 0.2 + s.t * 1.4, e.rz + Math.sin(s.a + s.t * 2) * rr);
        });
        pos.needsUpdate = true;
        sp.points.material.opacity = 0.65 + Math.sin(tNow * 7) * 0.35;
        sp.glow.position.set(e.rx, e.ry + 0.4, e.rz);
        sp.glow.material.opacity = 0.35 + Math.sin(tNow * 4) * 0.2;
      }

      // nameplate
      const lab = e.label;
      const far = Math.hypot(e.rx - S.local.x, e.rz - S.local.z) > (e.k === 'mob' ? 32 : e.k === 'obj' ? 14 : 55);
      if (far || (e.k === 'mob' && e.d) || (self && S.me.fl)) { lab.style.display = 'none'; continue; }
      this.v.set(e.rx, e.ry + (e.model.height + (e.marker ? 0.2 : 0.25)) * e.sc + (e.fl ? 2 : 0), e.rz).project(camera);
      if (this.v.z > 1 || this.v.z < -1) { lab.style.display = 'none'; continue; }
      lab.style.display = '';
      lab.style.left = ((this.v.x * 0.5 + 0.5) * w) + 'px';
      lab.style.top = ((-this.v.y * 0.5 + 0.5) * h) + 'px';
    }

    // selection ring
    const t = S.ents.get(S.me.tg);
    if (t && t.group) {
      this.ring.visible = true;
      this.ring.position.set(t.rx, t.ry + 0.08, t.rz);
      const col = t.k === 'mob' ? (t.d ? 0x888888 : t.el ? 0xff9a20 : 0xffe000) : t.k === 'obj' ? 0x40c8ff : 0x40ff40;
      this.ring.material.color.setHex(col);
      this.ring.rotation.z += dt * 0.8;
      this.ring.scale.setScalar((t.model.kind === 'human' ? 0.85 : Math.max(0.9, t.model.height * 0.75)) * t.sc);
    } else this.ring.visible = false;
  }

  screenPos(id, heightFactor, camera, w, h) {
    const e = S.ents.get(id);
    if (!e) return null;
    this.v.set(e.rx, e.ry + e.model.height * e.sc * heightFactor, e.rz).project(camera);
    if (this.v.z > 1) return null;
    return { x: (this.v.x * 0.5 + 0.5) * w, y: (-this.v.y * 0.5 + 0.5) * h };
  }
}
