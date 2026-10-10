import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import * as SkeletonUtils from 'three/addons/utils/SkeletonUtils.js';

// Humanoid characters use one rigged, animated model (Idle / Walk / Run) which is tinted and
// dressed per class, NPC and monster with props attached to its bones.

// Several rigged models: each has its own body (male soldier, female, robot/skeleton) and animations.
const RIG_FILES = { soldier: 'assets/Soldier.glb', michelle: 'assets/Michelle.glb', xbot: 'assets/Xbot.glb' };
const rigs = {}; // name -> { gltf, height, flip, clips }
let rigsReady = null;

function loadRig(name) {
  return new Promise((resolve) => {
    new GLTFLoader().load(RIG_FILES[name], (g) => {
      g.scene.updateMatrixWorld(true);
      const pos = (n) => { let o = null; g.scene.traverse((x) => { if (x.name === n) o = x; }); return o ? new THREE.Vector3().setFromMatrixPosition(o.matrixWorld) : null; };
      const head = pos('mixamorigHead'), foot = pos('mixamorigLeftFoot'), toe = pos('mixamorigLeftToeBase');
      let height = 1.8;
      if (head && foot) height = Math.max(0.3, (head.y - foot.y) * 1.12 + 0.12);
      // which way the model looks (+Z in the game): toes point forward
      const flip = toe && foot ? toe.z - foot.z < 0 : false;
      rigs[name] = { gltf: g, height, flip };
      console.log('rig', name, 'height', height.toFixed(2), 'flip', flip);
      resolve(true);
    }, undefined, () => resolve(false));
  });
}

export function loadHumanAssets() {
  if (!rigsReady) {
    rigsReady = loadRig('soldier').then(async (ok) => {
      await Promise.all([loadRig('michelle'), loadRig('xbot')]);
      return ok;
    });
  }
  return rigsReady;
}

export const humanAssetsLoaded = () => !!rigs.soldier;

// clips by role for each rig (Michelle has no walk/run of its own and re-uses the soldier clips: same bone names)
function clipsFor(name) {
  const r = rigs[name];
  const find = (arr, ...n) => arr.find((a) => n.includes(a.name));
  const own = r.gltf.animations;
  if (name === 'xbot') return { idle: find(own, 'idle'), walk: find(own, 'walk'), run: find(own, 'run') };
  const src = rigs.soldier.gltf.animations;
  // the hips track of the soldier clips is not compatible with this skeleton's rest pose
  const adapt = (c) => new THREE.AnimationClip(c.name, c.duration, c.tracks.filter((t) => !t.name.startsWith('mixamorigHips.')));
  return { idle: adapt(find(src, 'Idle')), walk: adapt(find(src, 'Walk')), run: adapt(find(src, 'Run')) };
}

const M = (c, o = {}) => new THREE.MeshLambertMaterial({ color: c, ...o });
const TIER_COLORS = [0x8a8a8a, 0x6a9a5a, 0x5a9ac8, 0x4a7ae0, 0x9a6ad8, 0xd8a040, 0xe08040, 0xe05050, 0xd04aa0, 0xffd860];

// item id suffix -> tier index (bf = dungeon set, shown as dark gold)
function tierOf(id) {
  if (!id) return 0;
  if (id.endsWith('_bf')) return 8;
  const t = parseInt(id.split('_').pop(), 10);
  return Number.isFinite(t) ? t : 0;
}
const metalColor = (id) => (id.endsWith('_bf') ? 0x6a5a30 : [0x9a9a9a, 0xa8b0b8, 0xb8c0cc, 0xc4ccd8, 0xaab4d0, 0xc8c8e0, 0xd8d0b0, 0xe8e0c0, 0xf0d890, 0xfff0b0][Math.min(9, tierOf(id))]);

// ---- props: Object3D whose local +Y is "along the held item" ----
function blade(len, w, color) {
  const g = new THREE.Group();
  const b = new THREE.Mesh(new THREE.BoxGeometry(w, len, w * 0.25), M(color, { emissive: 0x222222 }));
  b.position.y = len / 2 + 0.08; g.add(b);
  const guard = new THREE.Mesh(new THREE.BoxGeometry(w * 3.2, 0.05, 0.07), M(0x8a6a28)); guard.position.y = 0.08; g.add(guard);
  const grip = new THREE.Mesh(new THREE.CylinderGeometry(0.025, 0.025, 0.2, 6), M(0x4a2a14)); grip.position.y = -0.04; g.add(grip);
  return g;
}

function propFor(kind, id) {
  const tier = tierOf(id);
  const tc = TIER_COLORS[Math.min(9, tier)];
  const g = new THREE.Group();
  switch (kind) {
    case 'sword': return blade(0.5 + tier * 0.025, 0.07, metalColor(id || 'x_0'));
    case 'dagger': return blade(0.28, 0.05, 0xb0b4bc);
    case 'scimitar': { const b = blade(0.55, 0.08, 0xd8d8e0); b.rotation.z = 0.15; return b; }
    case 'axe': {
      const s = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.8, 6), M(0x6b4423)); s.position.y = 0.3; g.add(s);
      const h = new THREE.Mesh(new THREE.BoxGeometry(0.3, 0.2, 0.05), M(0xaab0b8)); h.position.set(0.12, 0.62, 0); g.add(h);
      return g;
    }
    case 'pick': {
      const s = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.7, 6), M(0x6b4423)); s.position.y = 0.25; g.add(s);
      const h = new THREE.Mesh(new THREE.ConeGeometry(0.05, 0.5, 5), M(0x8a8a90)); h.rotation.z = Math.PI / 2; h.position.set(0, 0.58, 0); g.add(h);
      return g;
    }
    case 'staff': {
      const s = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.035, 1.7, 6), M(0x6b4423)); s.position.y = 0.45; g.add(s);
      const orb = new THREE.Mesh(new THREE.SphereGeometry(0.1 + tier * 0.006, 10, 8), new THREE.MeshBasicMaterial({ color: tier > 0 ? tc : 0x66aaff })); orb.position.y = 1.35; g.add(orb);
      const halo = new THREE.Mesh(new THREE.TorusGeometry(0.14, 0.02, 5, 14), M(0xc8a24a)); halo.position.y = 1.3; halo.rotation.x = Math.PI / 2; g.add(halo);
      return g;
    }
    case 'witchstaff': {
      const s = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.04, 1.5, 5), M(0x3a2a1a)); s.position.y = 0.4; g.add(s);
      const orb = new THREE.Mesh(new THREE.SphereGeometry(0.1, 8, 6), new THREE.MeshBasicMaterial({ color: 0x9aff4a })); orb.position.y = 1.2; g.add(orb);
      return g;
    }
    case 'spear': {
      const s = new THREE.Mesh(new THREE.CylinderGeometry(0.025, 0.025, 2.1, 6), M(0x6b4423)); s.position.y = 0.5; g.add(s);
      const t = new THREE.Mesh(new THREE.ConeGeometry(0.06, 0.28, 6), M(0xc0c4cc)); t.position.y = 1.65; g.add(t);
      return g;
    }
    case 'bow': {
      const bow = new THREE.Mesh(new THREE.TorusGeometry(0.45 + tier * 0.01, 0.018, 6, 20, Math.PI), M(tier > 2 ? tc : 0x6b4a22));
      bow.rotation.z = -Math.PI / 2; g.add(bow);
      const str = new THREE.Mesh(new THREE.BoxGeometry(0.006, 0.9, 0.006), M(0xdddddd)); str.position.x = 0; g.add(str);
      return g;
    }
    case 'rifle': {
      const s = new THREE.Mesh(new THREE.CylinderGeometry(0.025, 0.025, 1.1, 6), M(0x333333)); s.position.y = 0.3; g.add(s);
      const st = new THREE.Mesh(new THREE.BoxGeometry(0.07, 0.35, 0.06), M(0x6b4423)); st.position.y = -0.2; g.add(st);
      return g;
    }
    case 'shield': {
      const sh = new THREE.Mesh(new THREE.CylinderGeometry(0.3, 0.3, 0.06, 16), M(tier > 0 ? tc : 0x9a2a2a)); sh.rotation.x = Math.PI / 2; g.add(sh);
      const rim = new THREE.Mesh(new THREE.TorusGeometry(0.3, 0.025, 5, 20), M(metalColor(id || 'x_0'))); g.add(rim);
      const boss = new THREE.Mesh(new THREE.SphereGeometry(0.07, 8, 6), M(0xc9a24a)); boss.position.z = 0.04; g.add(boss);
      return g;
    }
    case 'orb': {
      const o = new THREE.Mesh(new THREE.SphereGeometry(0.11, 10, 8), new THREE.MeshBasicMaterial({ color: tc })); g.add(o);
      return g;
    }
    case 'candle': {
      const c = new THREE.Mesh(new THREE.CylinderGeometry(0.03, 0.03, 0.15, 6), M(0xe8e0c0)); g.add(c);
      const f = new THREE.Mesh(new THREE.SphereGeometry(0.04, 6, 5), new THREE.MeshBasicMaterial({ color: 0xffaa33 })); f.position.y = 0.12; g.add(f);
      return g;
    }
  }
  return g;
}

function headPropFor(kind, id, color) {
  const tier = tierOf(id);
  const mc = id ? metalColor(id) : 0xb9bec8;
  const g = new THREE.Group();
  switch (kind) {
    case 'helmet': {
      const h = new THREE.Mesh(new THREE.SphereGeometry(0.17, 14, 10, 0, 6.283, 0, 1.75), M(mc, { emissive: 0x111111 })); h.position.y = 0.02; g.add(h);
      const nasal = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.12, 0.02), M(mc)); nasal.position.set(0, -0.02, 0.17); g.add(nasal);
      if (tier >= 3) { const crest = new THREE.Mesh(new THREE.BoxGeometry(0.03, 0.12, 0.3), M(TIER_COLORS[Math.min(9, tier)])); crest.position.set(0, 0.2, 0); g.add(crest); }
      return g;
    }
    case 'hood': {
      const h = new THREE.Mesh(new THREE.SphereGeometry(0.19, 12, 8, 0, 6.283, 0, 1.9), M(color || 0x2f4f2a)); h.position.set(0, 0.0, -0.02); g.add(h);
      return g;
    }
    case 'wizardhat': {
      const cone = new THREE.Mesh(new THREE.ConeGeometry(0.2, 0.55, 12), M(color || 0x23358a)); cone.position.y = 0.4; cone.rotation.x = -0.15; g.add(cone);
      const brim = new THREE.Mesh(new THREE.CylinderGeometry(0.3, 0.3, 0.025, 16), M(color || 0x23358a)); brim.position.y = 0.14; g.add(brim);
      const band = new THREE.Mesh(new THREE.CylinderGeometry(0.205, 0.205, 0.05, 14), M(0xd8b13a)); band.position.y = 0.17; g.add(band);
      return g;
    }
    case 'straw': {
      const c = new THREE.Mesh(new THREE.CylinderGeometry(0.15, 0.18, 0.1, 12), M(0xe0c070)); c.position.y = 0.15; g.add(c);
      const b = new THREE.Mesh(new THREE.CylinderGeometry(0.34, 0.34, 0.02, 16), M(0xe0c070)); b.position.y = 0.1; g.add(b);
      return g;
    }
    case 'turban': {
      const t = new THREE.Mesh(new THREE.SphereGeometry(0.2, 12, 8, 0, 6.283, 0, 1.6), M(color || 0xe8e0c8)); t.position.y = 0.07; g.add(t);
      return g;
    }
    case 'crown': {
      const base = new THREE.Mesh(new THREE.CylinderGeometry(0.17, 0.19, 0.12, 8), M(0xe8c040, { emissive: 0x443300 })); base.position.y = 0.15; g.add(base);
      for (let i = 0; i < 6; i++) { const p = new THREE.Mesh(new THREE.ConeGeometry(0.035, 0.12, 4), M(0xe8c040, { emissive: 0x443300 })); const a = (i / 6) * 6.283; p.position.set(Math.cos(a) * 0.17, 0.27, Math.sin(a) * 0.17); g.add(p); }
      return g;
    }
    case 'witchhat': {
      const cone = new THREE.Mesh(new THREE.ConeGeometry(0.22, 0.6, 10), M(0x2a1a3a)); cone.position.y = 0.42; cone.rotation.z = 0.2; g.add(cone);
      const brim = new THREE.Mesh(new THREE.CylinderGeometry(0.34, 0.34, 0.02, 14), M(0x2a1a3a)); brim.position.y = 0.14; g.add(brim);
      return g;
    }
    case 'mask': {
      const m = new THREE.Mesh(new THREE.BoxGeometry(0.26, 0.1, 0.06), M(color || 0x8a1a1a)); m.position.set(0, -0.04, 0.15); g.add(m);
      const cap = new THREE.Mesh(new THREE.SphereGeometry(0.175, 12, 8, 0, 6.283, 0, 1.4), M(0x2a2a2a)); cap.position.y = 0.03; g.add(cap);
      return g;
    }
    case 'horns': {
      for (const s of [-1, 1]) { const h = new THREE.Mesh(new THREE.ConeGeometry(0.05, 0.3, 5), M(0xe8e0c8)); h.position.set(s * 0.15, 0.2, 0); h.rotation.z = -s * 0.5; g.add(h); }
      return g;
    }
    case 'pharaoh': {
      const h = new THREE.Mesh(new THREE.CylinderGeometry(0.16, 0.22, 0.34, 4), M(0xe8c040, { emissive: 0x443300 })); h.position.y = 0.18; h.rotation.y = Math.PI / 4; g.add(h);
      return g;
    }
    case 'ears': {
      for (const s of [-1, 1]) { const e = new THREE.Mesh(new THREE.ConeGeometry(0.07, 0.35, 4), M(color || 0x8a6a3a)); e.position.set(s * 0.22, 0.02, 0); e.rotation.z = -s * 1.2; g.add(e); }
      return g;
    }
  }
  return g;
}

// Armor (visible parts only: shoulder pads / cloak).
function shoulderProp(id, side) {
  const mc = metalColor(id);
  const g = new THREE.Group();
  const pad = new THREE.Mesh(new THREE.SphereGeometry(0.13, 10, 8), M(new THREE.Color(mc).multiplyScalar(0.75), { emissive: 0x111111 }));
  pad.scale.set(1, 0.7, 1);
  g.add(pad);
  if (tierOf(id) >= 3) { const sp = new THREE.Mesh(new THREE.ConeGeometry(0.04, 0.16, 5), M(TIER_COLORS[Math.min(9, tierOf(id))])); sp.position.y = 0.12; g.add(sp); }
  g.userData.side = side;
  return g;
}

export function quiverProp() {
  const g = new THREE.Group();
  const q = new THREE.Mesh(new THREE.CylinderGeometry(0.07, 0.06, 0.5, 8), M(0x6b4a28)); q.rotation.z = 0.25; g.add(q);
  for (let i = 0; i < 3; i++) { const a = new THREE.Mesh(new THREE.CylinderGeometry(0.008, 0.008, 0.18, 4), M(0xd8d0b8)); a.position.set(-0.07 + i * 0.03, 0.3, 0); a.rotation.z = 0.25; g.add(a); }
  return g;
}

function cloakProp(id, color) {
  const g = new THREE.Group();
  const c = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.95, 0.04), M(color || [0x6a3a2a, 0x3a5a3a, 0x3a3a7a, 0x7a2a2a][tierOf(id) % 4]));
  c.position.set(0, -0.45, 0);
  g.add(c);
  return g;
}

// worn armor pieces drawn on top of the rig
function armorColor(id) {
  const k = (id || '').split('_')[0];
  const base = { cloth: 0x4a5ad0, leather: 0x7a5a38, mail: 0xa8aeb8 }[k] || 0x888888;
  const c = new THREE.Color(base);
  const t = tierOf(id);
  c.lerp(new THREE.Color(TIER_COLORS[Math.min(9, t)]), id.endsWith('_bf') ? 0.5 : Math.min(0.45, t * 0.05));
  return c.getHex();
}

function chestPlate(id) {
  const g = new THREE.Group();
  const col = armorColor(id), k = id.split('_')[0];
  const body = new THREE.Mesh(new THREE.BoxGeometry(k === 'cloth' ? 0.36 : 0.42, 0.42, k === 'cloth' ? 0.24 : 0.3), M(col, { emissive: tierOf(id) >= 7 ? 0x332200 : 0x000000 }));
  g.add(body);
  if (k === 'cloth') { const skirt = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.34, 0.5, 10), M(col)); skirt.position.y = -0.45; g.add(skirt); }
  const trim = new THREE.Mesh(new THREE.BoxGeometry(0.44, 0.05, 0.32), M(TIER_COLORS[Math.min(9, tierOf(id))])); trim.position.y = -0.2; g.add(trim);
  return g;
}

function limbPiece(id, len, r) {
  const m = new THREE.Mesh(new THREE.CylinderGeometry(r, r * 0.85, len, 8), M(armorColor(id)));
  m.position.y = -len / 2 + 0.02;
  const g = new THREE.Group(); g.add(m); return g;
}

const SLOT_BONES = {
  right: 'RightHand', left: 'LeftHand', head: 'Head', spine: 'Spine2', rshoulder: 'RightShoulder', lshoulder: 'LeftShoulder',
};

function findBone(root, name) {
  let found = null;
  root.traverse((o) => { if (!found && o.isBone && (o.name === 'mixamorig' + name || o.name === 'mixamorig:' + name)) found = o; });
  return found;
}

// Where props sit on the bones (offset in meters in bone-local axes, tuned for this rig).
const HOLD = {
  right: { pos: [0, 0.06, 0], rot: [Math.PI / 2, 0, 0] },
  left: { pos: [0, 0.06, 0], rot: [0, 0, 0] },
  head: { pos: [0, 0.2, 0.03], rot: [0, 0, 0] },
};

function attach(bone, obj, worldOffset, rot, scale = 1) {
  if (!bone) return;
  const ws = new THREE.Vector3();
  bone.getWorldScale(ws);
  const s = 1 / (ws.x || 1);
  obj.position.set(worldOffset[0] * s, worldOffset[1] * s, worldOffset[2] * s);
  obj.rotation.set(rot[0], rot[1], rot[2]);
  obj.scale.setScalar(s * scale);
  bone.add(obj);
}

/**
 * Creates an animated humanoid.
 * opts: { tint, emissive, scale, hideVisor, props: [{kind, id, slot}], head: {kind, id, color}, cloak, shoulders, transparent, opacity }
 */
export function makeHumanoid(opts = {}) {
  const root = new THREE.Group();
  const rigName = rigs[opts.rig] ? opts.rig : 'soldier';
  const rig = rigs[rigName];
  const inst = SkeletonUtils.clone(rig.gltf.scene);
  const target = 1.8 * (opts.scale || 1);
  const sc = target / rig.height;
  const bs = opts.bodyScale || [1, 1, 1];
  inst.scale.set(sc * bs[0], sc * bs[1], sc * bs[2]);
  const pivotG = new THREE.Group(); // turns the model so that it looks along +Z
  if (rig.flip) pivotG.rotation.y = Math.PI;
  pivotG.add(inst);
  root.add(pivotG);
  const mats = [];
  inst.traverse((o) => {
    if (o.isMesh) {
      o.castShadow = true;
      o.frustumCulled = false;
      o.material = o.material.clone();
      if (opts.tint != null && o.material.color) o.material.color.set(opts.tint);
      if (opts.emissive != null) { o.material.emissive = new THREE.Color(opts.emissive); }
      if (opts.opacity != null) { o.material.transparent = true; o.material.opacity = opts.opacity; o.castShadow = false; }
      if (o.name.includes('visor') && opts.hideVisor) o.visible = false;
      mats.push(o.material);
    }
  });

  const mixer = new THREE.AnimationMixer(inst);
  const cl = clipsFor(rigName);
  const act = { idle: mixer.clipAction(cl.idle), walk: mixer.clipAction(cl.walk), run: mixer.clipAction(cl.run) };
  Object.values(act).forEach((a) => { a.play(); a.setEffectiveWeight(0); });
  act.idle.setEffectiveWeight(1);
  const w = { idle: 1, walk: 0, run: 0 };

  const bones = {
    rightHand: findBone(inst, 'RightHand'), leftHand: findBone(inst, 'LeftHand'), head: findBone(inst, 'Head'),
    spine: findBone(inst, 'Spine2'), rArm: findBone(inst, 'RightArm'), lArm: findBone(inst, 'LeftArm'),
    rFore: findBone(inst, 'RightForeArm'), lFore: findBone(inst, 'LeftForeArm'), spine1: findBone(inst, 'Spine1'),
    rShoulder: findBone(inst, 'RightShoulder'), lShoulder: findBone(inst, 'LeftShoulder'),
    hips: findBone(inst, 'Hips'), lUpLeg: findBone(inst, 'LeftUpLeg'), rUpLeg: findBone(inst, 'RightUpLeg'),
    lFoot: findBone(inst, 'LeftFoot'), rFoot: findBone(inst, 'RightFoot'),
  };
  if (opts.headScale && bones.head) bones.head.scale.multiplyScalar(opts.headScale);
  const attached = [];
  const clear = () => { attached.forEach((o) => o.removeFromParent()); attached.length = 0; };
  const put = (bone, obj, off, rot, scale) => { attach(bone, obj, off, rot, scale); attached.push(obj); };

  // props & worn items; can be rebuilt when the equipment changes
  function dress(spec) {
    clear();
    for (const p of spec.props || []) {
      const obj = propFor(p.kind, p.id);
      const rot = p.slot === 'left' ? (p.kind === 'shield' ? [0, Math.PI / 2, Math.PI / 2] : p.kind === 'bow' ? [0, 0, 0] : [Math.PI / 2, 0, 0]) : [Math.PI / 2, 0, 0];
      const off = p.slot === 'left' ? [p.kind === 'shield' ? 0.14 : 0.0, 0.05, 0.02] : [0, 0.05, 0];
      put(p.slot === 'left' ? bones.leftHand : bones.rightHand, obj, off, rot, p.scale || 1);
    }
    if (spec.head) put(bones.head, headPropFor(spec.head.kind, spec.head.id, spec.head.color), [0, 0.1, 0.02], [0, 0, 0], 1);
    if (spec.shoulders) {
      put(bones.rShoulder, shoulderProp(spec.shoulders, 'r'), [0, 0.13, 0], [0, 0, 0], 0.8);
      put(bones.lShoulder, shoulderProp(spec.shoulders, 'l'), [0, 0.13, 0], [0, 0, 0], 0.8);
    }
    if (spec.chest) put(bones.spine, chestPlate(spec.chest), [0, 0.02, 0.0], [0, 0, 0], 1);
    if (spec.legs) { put(bones.lUpLeg, limbPiece(spec.legs, 0.42, 0.1), [0, 0, 0], [0, 0, 0], 1); put(bones.rUpLeg, limbPiece(spec.legs, 0.42, 0.1), [0, 0, 0], [0, 0, 0], 1); }
    if (spec.feet) for (const b of [bones.lFoot, bones.rFoot]) { const bt = new THREE.Mesh(new THREE.BoxGeometry(0.13, 0.2, 0.28), M(armorColor(spec.feet))); put(b, bt, [0, 0.0, 0.04], [0, 0, 0], 1); }
    if (spec.hands) for (const b of [bones.leftHand, bones.rightHand]) { const gl = new THREE.Mesh(new THREE.SphereGeometry(0.075, 8, 6), M(armorColor(spec.hands))); put(b, gl, [0, 0.07, 0], [0, 0, 0], 1); }
    if (spec.waist) { const belt = new THREE.Mesh(new THREE.TorusGeometry(0.2, 0.04, 6, 14), M(armorColor(spec.waist))); put(bones.hips, belt, [0, 0.12, 0], [Math.PI / 2, 0, 0], 1); }
    if (spec.cloak) put(bones.spine, cloakProp(spec.cloak.id, spec.cloak.color), [0, 0.22, -0.18], [0, 0, 0], 1);
    if (spec.extra) for (const ex of spec.extra) put(bones[ex.bone] || bones.spine, ex.obj, ex.off || [0, 0, 0], ex.rot || [0, 0, 0], ex.scale || 1);
  }
  dress(opts);

  let attackT = 0, castT = 0, hurt = 0, dead = 0;
  const qa = new THREE.Quaternion(), ax = new THREE.Vector3(1, 0, 0), az = new THREE.Vector3(0, 0, 1);
  const state = { phase: 0 };
  return {
    obj: root, height: target + 0.1, kind: 'human', tintMats: mats, setTint(c) { mats.forEach((m) => m.color.set(c)); },
    dress,
    attack() { attackT = 0.01; },
    cast(on) { castT = on ? 1 : 0; },
    update(dt, st) {
      // blend idle / walk / run
      const sp = st.speed;
      let tIdle = 0, tWalk = 0, tRun = 0;
      if (st.airborne) tIdle = 1;
      else if (sp < 0.08) tIdle = 1;
      else if (sp < 0.5) { tWalk = 1; } else tRun = 1;
      const k = Math.min(1, dt * 10);
      w.idle += (tIdle - w.idle) * k; w.walk += (tWalk - w.walk) * k; w.run += (tRun - w.run) * k;
      act.idle.setEffectiveWeight(w.idle); act.walk.setEffectiveWeight(w.walk); act.run.setEffectiveWeight(w.run);
      act.walk.setEffectiveTimeScale(Math.max(0.5, sp / 0.35) * (opts.slow ? 0.6 : 1));
      act.run.setEffectiveTimeScale(Math.max(0.7, sp / 0.95) * (opts.slow ? 0.6 : 1));
      mixer.update(dt);

      // overlays on top of the animation
      if (attackT > 0) {
        attackT += dt;
        const t = Math.min(1, attackT / 0.38);
        const swing = Math.sin(t * Math.PI);
        if (bones.rArm) { qa.setFromAxisAngle(ax, -swing * 1.9); bones.rArm.quaternion.premultiply(qa); }
        if (bones.rFore) { qa.setFromAxisAngle(ax, swing * 0.6); bones.rFore.quaternion.premultiply(qa); }
        if (bones.spine1) { qa.setFromAxisAngle(new THREE.Vector3(0, 1, 0), swing * 0.35); bones.spine1.quaternion.premultiply(qa); }
        if (t >= 1) attackT = 0;
      }
      if (castT > 0 || opts.armsForward) {
        for (const b of [bones.lArm, bones.rArm]) if (b) { qa.setFromAxisAngle(ax, -1.4); b.quaternion.premultiply(qa); }
      }
      // death: lie down
      dead += ((st.dead ? 1 : 0) - dead) * Math.min(1, dt * 6);
      inst.rotation.x = -dead * (Math.PI / 2);
      inst.position.y = dead * 0.22;
      if (st.dead && w.idle > 0.9) mixer.timeScale = 0; else mixer.timeScale = 1;
    },
  };
}
