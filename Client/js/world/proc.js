import * as THREE from 'three';
import { SKIN, HAIR } from '../util.js';
import { pixTex } from './human.js';

const mat = (c, o = {}) => new THREE.MeshLambertMaterial({ map: pixTex('n' + c, 4, 4, (p) => p.fill(c, 0.2)), ...o });
const box = (w, h, d, m) => { const x = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), m); x.castShadow = true; return x; };
const sph = (r, m) => { const x = new THREE.Mesh(new THREE.BoxGeometry(r * 1.7, r * 1.7, r * 1.7), m); x.castShadow = true; return x; };
const cyl = (rt, rb, h, m) => { const x = new THREE.Mesh(new THREE.BoxGeometry(Math.max(rt, rb) * 1.8, h, Math.max(rt, rb) * 1.8), m); x.castShadow = true; return x; };

function pivot(x, y, z) { const g = new THREE.Group(); g.position.set(x, y, z); return g; }

// Style for a player character derived from class + look.
export function playerStyle(cls, look) {
  const skin = SKIN[look?.skin ?? 0], hair = HAIR[look?.hair ?? 0];
  const base = { skin, hair, female: look?.gender === 'f', gear: cls };
  switch (cls) {
    case 'warrior': return { ...base, top: 0x8c8f99, bottom: 0x4a3b2c, accent: 0xb02a2a };
    case 'hunter': return { ...base, top: 0x4f7a3a, bottom: 0x6b4f2d, accent: 0x8a5a2b };
    case 'mage': return { ...base, top: 0x2a3f9a, bottom: 0x1d2a6b, accent: 0xd8b13a };
    default: return { ...base, top: 0x777777, bottom: 0x444444, accent: 0x999999 };
  }
}

export const FORESTER_STYLE = { skin: SKIN[1], hair: 0x8a8a8a, top: 0x3f6b35, bottom: 0x5b4426, accent: 0x7a4a1c, gear: 'forester', female: false };

export function makeHuman(style) {
  const root = new THREE.Group();
  const body = new THREE.Group();
  root.add(body);
  const skin = mat(style.skin), top = mat(style.top), bottom = mat(style.bottom), accent = mat(style.accent), hair = mat(style.hair);
  const metal = mat(0xb9bec8, { emissive: 0x111111 });
  const dark = mat(0x222222);

  const torso = box(0.5, 0.62, 0.28, top); torso.position.y = 1.12; body.add(torso);
  const hips = box(0.46, 0.2, 0.27, bottom); hips.position.y = 0.74; body.add(hips);
  const belt = box(0.5, 0.06, 0.3, accent); belt.position.y = 0.82; body.add(belt);

  const mkLeg = (x) => {
    const p = pivot(x, 0.76, 0);
    const l = box(0.2, 0.7, 0.22, bottom); l.position.y = -0.35; p.add(l);
    const b = box(0.22, 0.14, 0.3, mat(0x2e2218)); b.position.set(0, -0.66, 0.03); p.add(b);
    body.add(p); return p;
  };
  const legL = mkLeg(-0.12), legR = mkLeg(0.12);

  const mkArm = (x) => {
    const p = pivot(x, 1.4, 0);
    const a = box(0.15, 0.55, 0.16, top); a.position.y = -0.27; p.add(a);
    const h = sph(0.085, skin); h.position.y = -0.6; p.add(h);
    body.add(p); return p;
  };
  const armL = mkArm(-0.33), armR = mkArm(0.33);

  const neck = cyl(0.06, 0.06, 0.08, skin); neck.position.y = 1.46; body.add(neck);
  const head = new THREE.Group(); head.position.y = 1.62; body.add(head);
  const skull = sph(0.15, skin, 16, 12); skull.scale.set(1, 1.1, 1); head.add(skull);
  for (const sx of [-0.055, 0.055]) {
    const eye = sph(0.02, dark, 6, 6); eye.position.set(sx, 0.02, 0.135); head.add(eye);
  }
  const nose = box(0.03, 0.04, 0.04, skin); nose.position.set(0, -0.015, 0.15); head.add(nose);

  // hair (unless a helmet/hat replaces it)
  const noHair = style.gear === 'warrior' || style.gear === 'guard' || style.gear === 'bandit';
  if (!noHair) {
    const cap = new THREE.Mesh(new THREE.SphereGeometry(0.165, 16, 10, 0, Math.PI * 2, 0, Math.PI * 0.58), hair);
    cap.position.y = 0.03; cap.rotation.x = -0.25; cap.castShadow = true; head.add(cap);
    if (style.female) {
      const back = box(0.3, 0.42, 0.1, hair); back.position.set(0, -0.12, -0.12); head.add(back);
    }
  }

  const gear = style.gear;
  if (gear === 'warrior') {
    const helm = new THREE.Mesh(new THREE.SphereGeometry(0.17, 16, 10, 0, Math.PI * 2, 0, Math.PI * 0.62), metal);
    helm.position.y = 0.03; helm.castShadow = true; head.add(helm);
    const nasal = box(0.03, 0.12, 0.02, metal); nasal.position.set(0, 0, 0.16); head.add(nasal);
    for (const s of [-1, 1]) { const pad = sph(0.13, metal); pad.scale.set(1, 0.7, 1); pad.position.set(s * 0.34, 1.42, 0); body.add(pad); }
    const chest = box(0.54, 0.4, 0.32, metal); chest.position.y = 1.16; body.add(chest);
    // sword in right hand
    const sword = new THREE.Group();
    const blade = box(0.06, 0.85, 0.02, mat(0xdfe6ee, { emissive: 0x222222 })); blade.position.y = 0.45; sword.add(blade);
    const guard = box(0.22, 0.04, 0.05, mat(0x8a6a28)); guard.position.y = 0.02; sword.add(guard);
    sword.position.set(0, -0.6, 0.06); sword.rotation.x = Math.PI / 2; armR.add(sword);
    // shield on left arm
    const shield = cyl(0.27, 0.27, 0.05, mat(0x9a2a2a)); shield.rotation.z = Math.PI / 2; shield.position.set(-0.12, -0.3, 0.05);
    const rim = cyl(0.285, 0.285, 0.03, mat(0xc9a24a)); rim.rotation.z = Math.PI / 2; rim.position.set(-0.13, -0.3, 0.05);
    armL.add(rim, shield);
  } else if (gear === 'hunter') {
    const cloak = box(0.5, 0.75, 0.05, mat(0x2f4f2a)); cloak.position.set(0, 1.05, -0.17); body.add(cloak);
    const hood = new THREE.Mesh(new THREE.SphereGeometry(0.18, 12, 8, 0, Math.PI * 2, 0, Math.PI * 0.5), mat(0x2f4f2a));
    hood.position.set(0, 0.0, -0.02); hood.rotation.x = -0.3; head.add(hood);
    const quiver = cyl(0.06, 0.06, 0.45, mat(0x6b4423)); quiver.position.set(0.12, 1.2, -0.24); quiver.rotation.z = -0.25; body.add(quiver);
    for (let i = 0; i < 3; i++) { const ar = cyl(0.008, 0.008, 0.2, mat(0xeeeeee), 4); ar.position.set(0.14 + i * 0.015, 1.5, -0.24); ar.rotation.z = -0.25; body.add(ar); }
    // bow in left hand
    const bow = new THREE.Mesh(new THREE.TorusGeometry(0.42, 0.018, 6, 20, Math.PI), mat(0x6b4a22));
    bow.rotation.z = -Math.PI / 2; bow.position.set(0, -0.58, 0.1); bow.castShadow = true;
    const str = box(0.006, 0.84, 0.006, mat(0xdddddd)); str.position.set(0.0, -0.58, 0.1); str.rotation.z = 0;
    armL.add(bow, str);
  } else if (gear === 'mage') {
    const robe = cyl(0.25, 0.42, 0.85, mat(style.top), 14); robe.position.y = 0.5; body.add(robe);
    const trim = cyl(0.425, 0.425, 0.05, mat(style.accent), 14); trim.position.y = 0.1; body.add(trim);
    const hat = new THREE.Mesh(new THREE.ConeGeometry(0.2, 0.55, 12), mat(style.bottom)); hat.position.y = 0.4; hat.rotation.x = -0.15; hat.castShadow = true; head.add(hat);
    const brim = cyl(0.27, 0.27, 0.025, mat(style.bottom), 16); brim.position.y = 0.12; head.add(brim);
    const band = cyl(0.205, 0.205, 0.05, mat(style.accent), 14); band.position.y = 0.15; head.add(band);
    const staff = cyl(0.03, 0.03, 1.8, mat(0x6b4423), 6); staff.position.set(0, -0.2, 0.08); armR.add(staff);
    const orb = sph(0.1, new THREE.MeshBasicMaterial({ color: 0x66aaff })); orb.position.set(0, 0.74, 0.08); armR.add(orb);
  } else if (gear === 'merchant') {
    const apron = box(0.46, 0.55, 0.05, mat(style.accent)); apron.position.set(0, 0.95, 0.16); body.add(apron);
    const hat = cyl(0.19, 0.2, 0.1, mat(0x7a3a2a)); hat.position.y = 0.16; head.add(hat);
    const pack = box(0.4, 0.45, 0.22, mat(0x6b4423)); pack.position.set(0, 1.1, -0.25); body.add(pack);
    const roll = cyl(0.1, 0.1, 0.5, mat(0xc9b88a), 8); roll.rotation.z = Math.PI / 2; roll.position.set(0, 1.38, -0.28); body.add(roll);
  } else if (gear === 'farmer') {
    const hat = cyl(0.17, 0.2, 0.1, mat(0xe0c070)); hat.position.y = 0.16; head.add(hat);
    const brim = cyl(0.32, 0.32, 0.02, mat(0xe0c070), 16); brim.position.y = 0.1; head.add(brim);
    const skirt = cyl(0.26, 0.34, 0.5, mat(0x4a3a2a), 12); skirt.position.y = 0.55; body.add(skirt);
    const fork = new THREE.Group();
    const shaft = cyl(0.025, 0.025, 1.7, mat(0x6b4423), 6); shaft.position.y = 0.3; fork.add(shaft);
    for (const s of [-1, 0, 1]) { const tine = cyl(0.012, 0.012, 0.3, mat(0xaab0b8), 4); tine.position.set(s * 0.06, 1.2, 0); fork.add(tine); }
    fork.position.set(0, -0.6, 0.05); armR.add(fork);
  } else if (gear === 'guard') {
    const helm = new THREE.Mesh(new THREE.SphereGeometry(0.17, 12, 8, 0, Math.PI * 2, 0, Math.PI * 0.6), metal);
    helm.position.y = 0.03; head.add(helm);
    const plume = box(0.03, 0.12, 0.2, mat(style.accent)); plume.position.set(0, 0.22, -0.02); head.add(plume);
    const chest = box(0.54, 0.4, 0.32, mat(style.top)); chest.position.y = 1.16; body.add(chest);
    const cross = box(0.08, 0.3, 0.02, mat(style.accent)); cross.position.set(0, 1.16, 0.17); body.add(cross);
    const spear = new THREE.Group();
    const shaft = cyl(0.025, 0.025, 2.2, mat(0x6b4423), 6); shaft.position.y = 0.4; spear.add(shaft);
    const tip = new THREE.Mesh(new THREE.ConeGeometry(0.06, 0.28, 6), metal); tip.position.y = 1.64; spear.add(tip);
    spear.position.set(0, -0.6, 0.05); armR.add(spear);
  } else if (gear === 'huntsman') {
    const hat = cyl(0.18, 0.2, 0.14, mat(0x6b4a2a)); hat.position.y = 0.17; head.add(hat);
    const beard = box(0.22, 0.2, 0.1, mat(0x3a2a1a)); beard.position.set(0, -0.12, 0.09); head.add(beard);
    const fur = box(0.62, 0.14, 0.34, mat(0x5a3a22)); fur.position.y = 1.42; body.add(fur);
    const rifle = new THREE.Group();
    const bar = cyl(0.025, 0.025, 1.1, mat(0x333333), 6); bar.position.y = 0.3; rifle.add(bar);
    const stock = box(0.07, 0.3, 0.05, mat(0x6b4423)); stock.position.y = -0.2; rifle.add(stock);
    rifle.position.set(0, -0.6, 0.05); armR.add(rifle);
  } else if (gear === 'bandit') {
    const hood = new THREE.Mesh(new THREE.SphereGeometry(0.18, 12, 8, 0, Math.PI * 2, 0, Math.PI * 0.62), mat(0x2a2a2a));
    hood.position.set(0, 0.02, -0.01); hood.rotation.x = -0.2; head.add(hood);
    const mask = box(0.28, 0.1, 0.06, mat(0x8a1a1a)); mask.position.set(0, -0.06, 0.12); head.add(mask);
    const vest = box(0.52, 0.4, 0.3, mat(0x3a2a22)); vest.position.y = 1.1; body.add(vest);
    const dag = new THREE.Group();
    const blade = box(0.05, 0.4, 0.015, mat(0xb0b4bc)); blade.position.y = 0.25; dag.add(blade);
    dag.position.set(0, -0.6, 0.06); dag.rotation.x = Math.PI / 2; armR.add(dag);
  } else if (gear === 'forester') {
    const hat = cyl(0.17, 0.2, 0.14, mat(0x3a5a2a)); hat.position.y = 0.17; head.add(hat);
    const brim = cyl(0.27, 0.27, 0.02, mat(0x3a5a2a), 16); brim.position.y = 0.1; head.add(brim);
    const feather = box(0.02, 0.18, 0.04, mat(0xc03030)); feather.position.set(0.14, 0.27, 0); feather.rotation.z = -0.4; head.add(feather);
    const beard = box(0.2, 0.16, 0.1, mat(0x8a8a8a)); beard.position.set(0, -0.12, 0.09); head.add(beard);
    const vest = box(0.52, 0.4, 0.3, mat(0x5b4426)); vest.position.y = 1.08; body.add(vest);
    const axe = new THREE.Group();
    const shaft = cyl(0.025, 0.025, 0.9, mat(0x6b4423), 6); shaft.position.y = 0.2; axe.add(shaft);
    const head2 = box(0.22, 0.16, 0.04, mat(0xaab0b8)); head2.position.set(0.1, 0.6, 0); axe.add(head2);
    axe.position.set(0, -0.6, 0.05); armR.add(axe);
  }

  root.userData = { phase: 0 };
  const parts = { legL, legR, armL, armR, body, head };
  let attackT = 0;
  return {
    obj: root, height: gear === 'mage' ? 2.15 : 1.95, kind: 'human', parts,
    attack() { attackT = 0.01; },
    update(dt, st) {
      const u = root.userData;
      const moving = st.speed > 0.1 && !st.airborne;
      if (moving) u.phase += dt * 10 * Math.min(1.3, st.speed);
      const amp = moving ? 0.8 : 0;
      legL.rotation.x += ((Math.sin(u.phase) * amp) - legL.rotation.x) * Math.min(1, dt * 14);
      legR.rotation.x += ((-Math.sin(u.phase) * amp) - legR.rotation.x) * Math.min(1, dt * 14);
      const aTarget = moving ? 0.55 : 0;
      armL.rotation.x += ((-Math.sin(u.phase) * aTarget) - armL.rotation.x) * Math.min(1, dt * 14);
      if (st.airborne) { legL.rotation.x = -0.5; legR.rotation.x = 0.4; }
      if (attackT > 0) {
        attackT += dt;
        const k = Math.min(1, attackT / 0.35);
        armR.rotation.x = -Math.sin(k * Math.PI) * 2.2;
        if (k >= 1) attackT = 0;
      } else {
        armR.rotation.x += ((Math.sin(u.phase) * aTarget) - armR.rotation.x) * Math.min(1, dt * 14);
      }
      body.position.y = moving ? Math.abs(Math.sin(u.phase)) * 0.04 : Math.sin(performance.now() / 600) * 0.008;
      // death
      const dead = st.dead ? 1 : 0;
      u.dead = (u.dead ?? 0) + (dead - (u.dead ?? 0)) * Math.min(1, dt * 5);
      body.rotation.x = -u.dead * (Math.PI / 2);
      body.position.y += u.dead * 0.22;
    },
  };
}

const QUADS = {
  wolf: { shades: [0x6e6e72, 0x5b5b60, 0x7a7468, 0x4d4d52], belly: 0x9c9a95, scale: 1, torsoW: 0.95, snout: 0.3, tail: 0.6, earR: 0.065, eye: 0xffcc22, tusks: false },
  boar: { shades: [0x6b4a32, 0x5a3d28, 0x7a5a3c], belly: 0x8a6a4a, scale: 0.85, torsoW: 1.25, snout: 0.22, tail: 0.18, earR: 0.08, eye: 0x331100, tusks: true },
  bear: { shades: [0x5a3a22, 0x3d2a1a, 0x6b4a2c], belly: 0x6b4f33, scale: 1.55, torsoW: 1.15, snout: 0.2, tail: 0.15, earR: 0.1, eye: 0x221100, tusks: false },
  hyena: { shades: [0xa88a58, 0x987a4a, 0xb89a68], belly: 0xd8c898, scale: 0.9, torsoW: 0.85, snout: 0.26, tail: 0.4, earR: 0.075, eye: 0xffaa22, tusks: false },
  icewolf: { shades: [0xdfe8f2, 0xc8d8e8, 0xe8eef6], belly: 0xffffff, scale: 1.15, torsoW: 0.95, snout: 0.3, tail: 0.6, earR: 0.065, eye: 0x66ccff, tusks: false },
  polarbear: { shades: [0xf0f4f8, 0xe0e8f0, 0xe8edf4], belly: 0xffffff, scale: 1.65, torsoW: 1.15, snout: 0.22, tail: 0.12, earR: 0.1, eye: 0x112233, tusks: false },
  bull: { shades: [0x6a4a2a, 0x4a3420, 0x8a6a40], belly: 0x8a7a5a, scale: 1.25, torsoW: 1.2, snout: 0.22, tail: 0.5, earR: 0.08, eye: 0x331100, tusks: false, horns: true },
  sheep: { shades: [0xf4f4f0], belly: 0xeeeeea, scale: 0.8, torsoW: 1.35, snout: 0.16, tail: 0.15, earR: 0.09, eye: 0x222222, tusks: false, wool: true },
};

export function makeWolf(seed = 0, kind = 'wolf') {
  const V = QUADS[kind];
  const root = new THREE.Group();
  const body = new THREE.Group();
  root.add(body);
  root.scale.setScalar(V.scale);
  const fur = mat(V.shades[seed % V.shades.length]);
  const belly = mat(V.belly), dark = mat(0x1d1d1f);

  const torso = sph(0.38, fur, 14, 10); torso.scale.set(V.torsoW, 0.85, 1.9); torso.position.set(0, 0.78, 0); body.add(torso);
  const chest = sph(0.34, fur, 12, 8); chest.position.set(0, 0.84, 0.35); body.add(chest);
  const under = sph(0.3, belly, 10, 8); under.scale.set(0.9, 0.6, 1.5); under.position.set(0, 0.62, 0.05); body.add(under);
  const head = new THREE.Group(); head.position.set(0, 0.98, 0.78); body.add(head);
  const skull = sph(0.2, fur, 12, 10); skull.scale.set(1, 0.95, 1.1); head.add(skull);
  const snout = box(0.15, 0.13, V.snout, kind === 'wolf' ? mat(0x8a8780) : fur); snout.position.set(0, -0.06, 0.1 + V.snout * 0.45); head.add(snout);
  const nose = box(0.06, 0.05, 0.05, dark); nose.position.set(0, -0.02, 0.1 + V.snout * 0.9); head.add(nose);
  for (const s of [-1, 1]) {
    const ear = new THREE.Mesh(kind === 'wolf' ? new THREE.ConeGeometry(V.earR, 0.2, 5) : new THREE.SphereGeometry(V.earR, 8, 6), fur);
    ear.position.set(s * 0.11, kind === 'wolf' ? 0.2 : 0.17, -0.02); ear.castShadow = true; head.add(ear);
    const eye = sph(0.03, new THREE.MeshBasicMaterial({ color: V.eye }), 6, 6); eye.position.set(s * 0.09, 0.05, 0.16); head.add(eye);
    if (V.horns) {
      const horn = new THREE.Mesh(new THREE.ConeGeometry(0.045, 0.32, 5), mat(0xe8e0c8));
      horn.position.set(s * 0.17, 0.2, 0.0); horn.rotation.z = -s * 0.9; head.add(horn);
    }
    if (V.tusks) {
      const tusk = new THREE.Mesh(new THREE.ConeGeometry(0.03, 0.18, 5), mat(0xf0e6c8));
      tusk.position.set(s * 0.07, -0.06, 0.26); tusk.rotation.x = -0.5; head.add(tusk);
    }
  }
  if (V.wool) { const fluff = sph(0.5, mat(0xf4f4ee), 10, 8); fluff.scale.set(1.15, 0.95, 1.9); fluff.position.set(0, 0.88, 0); body.add(fluff); }
  const tail = new THREE.Mesh(new THREE.ConeGeometry(0.1, V.tail, 6), fur); tail.position.set(0, 0.85, -0.8); tail.rotation.x = -2.0; tail.castShadow = true; body.add(tail);

  const legs = [];
  for (const [x, z] of [[-0.2, 0.5], [0.2, 0.5], [-0.2, -0.5], [0.2, -0.5]]) {
    const p = pivot(x, 0.6, z);
    const l = box(0.13, 0.6, 0.13, fur); l.position.y = -0.3; p.add(l);
    const paw = box(0.15, 0.06, 0.2, dark); paw.position.set(0, -0.6, 0.03); p.add(paw);
    body.add(p); legs.push(p);
  }
  let attackT = 0;
  const u = { phase: 0, dead: 0 };
  return {
    obj: root, height: 1.5 * V.scale, kind: 'wolf', parts: { head, body },
    attack() { attackT = 0.01; },
    update(dt, st) {
      const moving = st.speed > 0.1;
      if (moving) u.phase += dt * 11 * Math.min(1.5, st.speed);
      const a = moving ? 0.9 : 0;
      legs[0].rotation.x += (Math.sin(u.phase) * a - legs[0].rotation.x) * Math.min(1, dt * 15);
      legs[3].rotation.x += (Math.sin(u.phase) * a - legs[3].rotation.x) * Math.min(1, dt * 15);
      legs[1].rotation.x += (-Math.sin(u.phase) * a - legs[1].rotation.x) * Math.min(1, dt * 15);
      legs[2].rotation.x += (-Math.sin(u.phase) * a - legs[2].rotation.x) * Math.min(1, dt * 15);
      tail.rotation.z = Math.sin(performance.now() / 300 + u.phase) * 0.25;
      let lunge = 0;
      if (attackT > 0) {
        attackT += dt;
        const k = Math.min(1, attackT / 0.4);
        lunge = Math.sin(k * Math.PI) * 0.45;
        head.rotation.x = Math.sin(k * Math.PI) * 0.5;
        if (k >= 1) { attackT = 0; head.rotation.x = 0; }
      }
      body.position.z = lunge;
      body.position.y = moving ? Math.abs(Math.sin(u.phase)) * 0.05 : 0;
      u.dead += ((st.dead ? 1 : 0) - u.dead) * Math.min(1, dt * 6);
      root.rotation.z = u.dead * (Math.PI / 2);
      // lying on its side: local +x becomes world +y, lift so it rests on the ground
      body.position.x = u.dead * 0.4;
    },
  };
}

export function makeSpider(seed = 0, col0) {
  const root = new THREE.Group();
  const body = new THREE.Group();
  root.add(body);
  const col = col0 ?? [0x2a2024, 0x3a2a1e, 0x24202e][seed % 3];
  const m = mat(col), mark = mat(0x9a2a1a), eyeM = new THREE.MeshBasicMaterial({ color: 0xff3030 });
  const abdomen = sph(0.45, m, 12, 10); abdomen.scale.set(1, 0.85, 1.25); abdomen.position.set(0, 0.62, -0.5); body.add(abdomen);
  const stripe = sph(0.2, mark, 8, 6); stripe.scale.set(0.8, 0.3, 1.2); stripe.position.set(0, 0.98, -0.5); body.add(stripe);
  const head = sph(0.27, m, 10, 8); head.position.set(0, 0.58, 0.25); body.add(head);
  for (let i = 0; i < 4; i++) { const e = sph(0.04, eyeM, 5, 5); e.position.set((i % 2 ? 1 : -1) * (0.07 + (i > 1 ? 0.07 : 0)), 0.65 + (i > 1 ? 0.05 : 0), 0.48); body.add(e); }
  for (const s of [-1, 1]) { const f = cyl(0.025, 0.04, 0.22, mat(0xcccc99), 5); f.position.set(s * 0.08, 0.45, 0.5); f.rotation.x = 0.4; body.add(f); }
  const legs = [];
  for (let i = 0; i < 8; i++) {
    const side = i < 4 ? -1 : 1, k = i % 4;
    const p = pivot(side * 0.2, 0.62, 0.35 - k * 0.22);
    p.rotation.y = side * (0.5 + k * 0.28 - 0.6);
    const up = cyl(0.03, 0.03, 0.55, m, 5); up.rotation.z = side * -1.0; up.position.set(side * 0.24, 0.1, 0); p.add(up);
    const lo = cyl(0.025, 0.02, 0.6, m, 5); lo.rotation.z = side * 0.35; lo.position.set(side * 0.55, -0.15, 0); p.add(lo);
    body.add(p); legs.push(p);
  }
  let attackT = 0; const u = { phase: 0, dead: 0 };
  return {
    obj: root, height: 1.0, kind: 'spider', parts: { body },
    attack() { attackT = 0.01; },
    update(dt, st) {
      const moving = st.speed > 0.1;
      if (moving) u.phase += dt * 16 * Math.min(1.4, st.speed);
      legs.forEach((l, i) => { const t = Math.sin(u.phase + i * 1.3) * (moving ? 0.35 : 0.03); l.rotation.z += (t - l.rotation.z) * Math.min(1, dt * 18); });
      let lunge = 0;
      if (attackT > 0) { attackT += dt; const k = Math.min(1, attackT / 0.35); lunge = Math.sin(k * Math.PI) * 0.4; if (k >= 1) attackT = 0; }
      body.position.z = lunge;
      u.dead += ((st.dead ? 1 : 0) - u.dead) * Math.min(1, dt * 6);
      root.rotation.z = u.dead * Math.PI;
      root.position.y = u.dead * 1.0;
    },
  };
}

const NPC_STYLES = {
  forester: FORESTER_STYLE,
  trainer_warrior: { skin: SKIN[2], hair: 0x555555, top: 0x7a7e88, bottom: 0x3a2f24, accent: 0x9a2a2a, gear: 'warrior', female: false },
  trainer_hunter: { skin: SKIN[1], hair: 0x8a4a1a, top: 0x3f6a32, bottom: 0x5a4126, accent: 0x7a4a1c, gear: 'hunter', female: true },
  trainer_mage: { skin: SKIN[0], hair: 0xcccccc, top: 0x5a2a7a, bottom: 0x3a1a5a, accent: 0xe0c040, gear: 'mage', female: false },
  smith: { skin: SKIN[2], hair: 0x3a2a1a, top: 0x6a4a3a, bottom: 0x3a3a3a, accent: 0x2a2a2a, gear: 'merchant', female: false },
  merchant: { skin: SKIN[1], hair: 0x6b4423, top: 0xa86a2a, bottom: 0x5a3a1c, accent: 0xe8e0c8, gear: 'merchant', female: true },
  farmer: { skin: SKIN[0], hair: 0xc9a24a, top: 0xb0603a, bottom: 0x4a3a2a, accent: 0xe8d8a0, gear: 'farmer', female: true },
  scholar: { skin: SKIN[0], hair: 0x6b4423, top: 0x2f7a5a, bottom: 0x1f5a3a, accent: 0xd8c070, gear: 'mage', female: true },
  guard: { skin: SKIN[2], hair: 0x2b1b10, top: 0x35507a, bottom: 0x2a2a3a, accent: 0xc9a24a, gear: 'guard', female: false },
  huntsman: { skin: SKIN[3], hair: 0x2b1b10, top: 0x6b4423, bottom: 0x3a2a1a, accent: 0x9a7a3a, gear: 'huntsman', female: false },
  bandit: { skin: SKIN[2], hair: 0x1b1b1b, top: 0x4a2a2a, bottom: 0x2a2a2a, accent: 0x888888, gear: 'bandit', female: false },
};

export function makeModel(ent) {
  switch (ent.m) {
    case 'wolf': case 'boar': case 'bear': return makeWolf(ent.id, ent.m);
    case 'spider': return makeSpider(ent.id);
    default:
      if (NPC_STYLES[ent.m]) return makeHuman(NPC_STYLES[ent.m]);
      return makeHuman(playerStyle(ent.c, ent.lk));
  }
}
