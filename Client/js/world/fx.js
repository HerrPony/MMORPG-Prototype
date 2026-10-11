import * as THREE from 'three';
import { S } from '../state.js';
import { $ } from '../util.js';

function radialTex() {
  const c = document.createElement('canvas');
  c.width = c.height = 64;
  const g = c.getContext('2d');
  const grd = g.createRadialGradient(32, 32, 0, 32, 32, 32);
  grd.addColorStop(0, 'rgba(255,255,255,1)');
  grd.addColorStop(0.4, 'rgba(255,255,255,0.55)');
  grd.addColorStop(1, 'rgba(255,255,255,0)');
  g.fillStyle = grd; g.fillRect(0, 0, 64, 64);
  return new THREE.CanvasTexture(c);
}

export class FX {
  constructor(scene, entities, getCamera) {
    this.ents = entities;
    this.getCamera = getCamera;
    this.setScene(scene);
    this.tex = radialTex();
    this.fct = $('fct');
    this.charges = new Map(); // caster id -> sprite
  }

  // all effects live in one group so that a scene change (preview -> world) cannot leave sprites behind
  setScene(scene) {
    if (this.root) this.root.removeFromParent();
    this.scene = scene; this.list = [];
    this.root = new THREE.Group();
    scene.add(this.root);
  }

  pos(id, f = 0.6) {
    const e = S.ents.get(id);
    if (!e) return null;
    return new THREE.Vector3(e.rx, e.ry + e.model.height * f, e.rz);
  }

  sprite(color, size, opacity = 1) {
    const s = new THREE.Sprite(new THREE.SpriteMaterial({
      map: this.tex, color, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending, opacity,
    }));
    s.scale.setScalar(size);
    this.root.add(s);
    return s;
  }

  burst(at, color, size = 1.4, dur = 0.35) {
    const s = this.sprite(color, size * 0.3);
    s.position.copy(at);
    let t = 0;
    this.list.push({
      update: (dt) => {
        t += dt;
        const k = t / dur;
        s.scale.setScalar(size * (0.3 + k));
        s.material.opacity = Math.max(0, 1 - k);
        if (k >= 1) { this.root.remove(s); return false; }
        return true;
      },
    });
  }

  ring(at, color, r0, r1, dur, vertical = 0) {
    const m = new THREE.Mesh(
      new THREE.RingGeometry(0.85, 1, 40),
      new THREE.MeshBasicMaterial({ color, transparent: true, depthWrite: false, side: THREE.DoubleSide, blending: THREE.AdditiveBlending, opacity: 0.9 }),
    );
    m.rotation.x = -Math.PI / 2;
    m.position.copy(at);
    this.root.add(m);
    let t = 0;
    this.list.push({
      update: (dt) => {
        t += dt;
        const k = Math.min(1, t / dur);
        m.scale.setScalar(r0 + (r1 - r0) * k);
        m.position.y = at.y + vertical * k;
        m.material.opacity = 0.9 * (1 - k);
        if (k >= 1) { this.root.remove(m); return false; }
        return true;
      },
    });
  }

  projectile(from, dstId, style, color) {
    const col = new THREE.Color(color);
    let mesh;
    if (style === 'arrow') {
      mesh = new THREE.Mesh(new THREE.CylinderGeometry(0.02, 0.02, 0.7, 5), new THREE.MeshBasicMaterial({ color: 0xe8d8a8 }));
      mesh.geometry.rotateX(Math.PI / 2);
      const tip = new THREE.Mesh(new THREE.ConeGeometry(0.05, 0.14, 5), new THREE.MeshBasicMaterial({ color: 0xcccccc }));
      tip.rotation.x = Math.PI / 2; tip.position.z = 0.4; mesh.add(tip);
    } else {
      mesh = new THREE.Mesh(new THREE.SphereGeometry(style === 'fire' ? 0.3 : 0.2, 10, 8), new THREE.MeshBasicMaterial({ color: col }));
      const halo = new THREE.Sprite(new THREE.SpriteMaterial({ map: this.tex, color: col, transparent: true, depthWrite: false, blending: THREE.AdditiveBlending }));
      halo.scale.setScalar(style === 'fire' ? 1.6 : 1.1);
      mesh.add(halo);
    }
    mesh.position.copy(from);
    this.root.add(mesh);
    const speed = style === 'arrow' ? 42 : style === 'fire' ? 22 : 28;
    this.list.push({
      update: (dt) => {
        const to = this.pos(dstId, 0.6);
        if (!to) { this.root.remove(mesh); return false; }
        const d = to.clone().sub(mesh.position);
        const step = speed * dt;
        if (d.length() <= step + 0.3) {
          this.root.remove(mesh);
          this.burst(to, col, style === 'fire' ? 2.2 : 1.4);
          return false;
        }
        d.normalize();
        mesh.position.addScaledVector(d, step);
        if (style === 'arrow') mesh.lookAt(to);
        return true;
      },
    });
  }

  slash(at, color) {
    const m = new THREE.Mesh(
      new THREE.RingGeometry(0.5, 0.7, 20, 1, 0.2, 2.2),
      new THREE.MeshBasicMaterial({ color, transparent: true, side: THREE.DoubleSide, depthWrite: false, blending: THREE.AdditiveBlending }),
    );
    m.position.copy(at);
    m.rotation.set(Math.random() * 0.8, Math.random() * 6, Math.random() * 3);
    this.root.add(m);
    let t = 0;
    this.list.push({ update: (dt) => { t += dt; m.scale.setScalar(1 + t * 4); m.material.opacity = 1 - t / 0.25; if (t > 0.25) { this.root.remove(m); return false; } return true; } });
    this.burst(at, color, 1.0, 0.25);
  }

  rising(id, color, dur = 1.0) {
    const base = this.pos(id, 0);
    if (!base) return;
    for (let i = 0; i < 10; i++) {
      const s = this.sprite(color, 0.4, 0.9);
      const a = (i / 10) * 6.28, r = 0.5;
      let t = -i * 0.05;
      this.list.push({
        update: (dt) => {
          t += dt;
          if (t < 0) return true;
          const p = this.pos(id, 0);
          if (!p) { this.root.remove(s); return false; }
          s.position.set(p.x + Math.cos(a + t * 4) * r, p.y + t * 2.2, p.z + Math.sin(a + t * 4) * r);
          s.material.opacity = 0.9 * Math.max(0, 1 - t / dur);
          if (t > dur) { this.root.remove(s); return false; }
          return true;
        },
      });
    }
  }

  levelUp(id) {
    const p = this.pos(id, 0);
    if (!p) return;
    const beam = new THREE.Mesh(
      new THREE.CylinderGeometry(0.9, 0.9, 7, 20, 1, true),
      new THREE.MeshBasicMaterial({ color: 0xffdd55, transparent: true, depthWrite: false, side: THREE.DoubleSide, blending: THREE.AdditiveBlending, opacity: 0.7 }),
    );
    this.root.add(beam);
    let t = 0;
    this.list.push({
      update: (dt) => {
        t += dt;
        const q = this.pos(id, 0);
        if (q) beam.position.set(q.x, q.y + 3.5, q.z);
        beam.material.opacity = 0.7 * Math.max(0, 1 - t / 1.8);
        beam.scale.set(1 + t * 0.2, 1, 1 + t * 0.2);
        if (t > 1.8) { this.root.remove(beam); return false; }
        return true;
      },
    });
    this.rising(id, 0xffe27a, 1.6);
    this.ring(p.clone().setY(p.y + 0.1), 0xffe27a, 0.5, 3.5, 1.0, 1.5);
  }

  // ---- network events ----
  onFx(m) {
    const skill = S.defs?.skills?.[m.k.replace(':cast', '')];
    const src = S.ents.get(m.s);
    if (m.k === 'levelup') { this.levelUp(m.s); return; }
    if (m.k === 'mobshot') {
      this.ents.attack(m.s);
      const from = this.pos(m.s, 0.9);
      if (from) this.projectile(from, m.d, m.style === 'frost' ? 'frost' : m.style === 'fire' ? 'fire' : 'bolt',
        { frost: '#7ad8ff', fire: '#ff7a1a', bolt: '#9aff4a' }[m.style] || '#c070ff');
      return;
    }
    if (m.k === 'pet_summon') { const p = this.pos(m.s, 0.3); if (p) { this.ring(p, 0x66ccff, 0.5, 2.2, 0.6); this.rising(m.s, 0x66ccff, 0.8); } return; }
    if (m.k === 'heal') { this.rising(m.s, 0x66ff88, 1.0); return; }
    if (m.k === 'bossheavy') { const p = this.pos(m.d, 0); if (p) { this.ents.attack(m.s); this.ring(p, 0xff4020, 0.5, 3, 0.45); this.burst(p.clone().setY(p.y + 1), 0xff6030, 2.6, 0.4); } return; }
    if (m.k === 'bossnova') { const p = this.pos(m.s, 0); if (p) this.ring(p.clone().setY(p.y + 0.1), 0xff5030, 0.5, 14, 0.8); return; }
    if (m.k === 'ground') {
      const base = new THREE.Vector3(m.x, 0, m.z);
      base.y = (S.ents.get(m.s)?.ry ?? 0);
      const col = new THREE.Color(m.color || '#aee4ff');
      const disc = new THREE.Mesh(new THREE.CircleGeometry(m.r, 32), new THREE.MeshBasicMaterial({ color: col, transparent: true, opacity: 0.25, depthWrite: false, side: THREE.DoubleSide }));
      disc.rotation.x = -Math.PI / 2; disc.position.set(m.x, base.y + 0.15, m.z); this.root.add(disc);
      let t = 0;
      this.list.push({ update: (dt) => {
        t += dt;
        for (let i = 0; i < 3; i++) { const q = this.sprite(col, 0.5, 0.8); q.position.set(m.x + (Math.random() - 0.5) * m.r * 2, base.y + 4 + Math.random() * 2, m.z + (Math.random() - 0.5) * m.r * 2); let tt = 0; this.list.push({ update: (d2) => { tt += d2; q.position.y -= d2 * 5; q.material.opacity = 0.8 * (1 - tt / 0.9); if (tt > 0.9) { this.root.remove(q); return false; } return true; } }); }
        disc.material.opacity = 0.25 * Math.max(0, 1 - t / m.dur);
        if (t > m.dur) { this.root.remove(disc); return false; }
        return true;
      } });
      return;
    }
    if (m.k === 'bite') {
      this.ents.attack(m.s);
      const to = this.pos(m.d, 0.6);
      if (to) setTimeout(() => this.burst(to, 0xff3030, 1.1, 0.3), 200);
      return;
    }
    if (!skill) return;
    const color = new THREE.Color(skill.color || '#ffffff');
    if (m.k.endsWith(':cast')) {
      this.ents.castAnim(m.s, true);
      setTimeout(() => this.ents.castAnim(m.s, false), (skill.cast || 1) * 1000);
      const p = this.pos(m.s, 0.9);
      if (!p) return;
      const s = this.sprite(color, 0.2, 0.9);
      let t = 0;
      this.list.push({
        update: (dt) => {
          t += dt;
          const q = this.pos(m.s, 0.9);
          if (!q) { this.root.remove(s); return false; }
          s.position.copy(q).add(new THREE.Vector3(Math.sin(src.rr) * 0.6, 0, Math.cos(src.rr) * 0.6));
          s.scale.setScalar(0.3 + Math.min(t, 1.5) * 0.6);
          if (t > Math.max(0.8, skill.cast || 1)) { this.root.remove(s); return false; }
          return true;
        },
      });
      return;
    }
    if (src) {
      this.ents.attack(m.s);
      if (src.id === S.id || true) {
        // face the target briefly is handled client-side for the own player
      }
    }
    const from = this.pos(m.s, 0.85);
    const to = this.pos(m.d, 0.6);
    switch (skill.style) {
      case 'melee':
        if (to) setTimeout(() => this.slash(to, color), 120);
        break;
      case 'nova':
        if (from) this.ring(from.clone().setY(from.y - 0.7), color, 0.5, skill.radius || 5, 0.5);
        break;
      case 'buff':
        this.rising(m.s, color, 1.1);
        break;
      default:
        if (from && to) {
          const delay = skill.style === 'arrow' && src?.model?.bowDelay ? src.model.bowDelay * 1000 : 0;
          if (delay) setTimeout(() => { const f2 = this.pos(m.s, 0.85); if (f2) this.projectile(f2, m.d, skill.style, skill.color); }, delay);
          else this.projectile(from, m.d, skill.style, skill.color);
        }
    }
  }

  onDmg(m) {
    const cam = this.getCamera();
    const w = innerWidth, h = innerHeight;
    const sp = this.ents.screenPos(m.id, 0.85, cam.camera, w, h);
    if (!sp) return;
    const cls = m.s === S.id ? 'out' : m.id === S.id ? 'in' : 'other';
    const mk = (text, c) => {
      const d = document.createElement('div');
      d.className = 'fct ' + c;
      d.textContent = text;
      d.style.left = (sp.x + (Math.random() - 0.5) * 50) + 'px';
      d.style.top = sp.y + 'px';
      this.fct.appendChild(d);
      setTimeout(() => d.remove(), 1250);
    };
    if (m.a > 0) mk(m.a + (m.c ? '!' : ''), cls + (m.c ? ' crit' : ''));
    if (m.ab > 0) mk(m.a > 0 ? `(${m.ab} absorbiert)` : 'Absorbiert', 'abs');
  }

  update(dt) {
    this.list = this.list.filter((f) => f.update(dt));
  }
}
