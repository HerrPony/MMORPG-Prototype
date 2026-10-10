import { makeHumanoid, loadHumanAssets as loadRigs, quiverProp } from './human.js';
import { loadBeastAssets, makePixelBeast } from './beasts.js';
import { makeSpider } from './proc.js';
import { SKIN, HAIR } from '../util.js';
import { makeCroc, makeScorpion, makeWisp, makeGriffin, makeObject } from './creatures2.js';

export { makeGriffin };
export const loadHumanAssets = () => Promise.all([loadRigs(), loadBeastAssets()]);
const PANTS = { warrior: 0x4a3b2c, hunter: 0x5a4a2d, mage: 0x23358a };

const CLASS_LOOK = {
  warrior: { pants: PANTS.warrior, tint: 0xc4c8d4, emissive: 0x0a0a10, head: { kind: 'helmet', id: 'mail_head_1' }, mh: 'sword', oh: 'shield', body: [1.14, 1, 1.12],
    base: { chest: 'mail_chest_1', legs: 'mail_legs_1', feet: 'mail_feet_1', hands: 'mail_hands_1', shoulder: 'mail_shoulder_1' } },
  hunter: { pants: PANTS.hunter, tint: 0x8aa66a, emissive: 0x08140a, hideVisor: true, head: { kind: 'hood', color: 0x2f4f2a }, mh: 'bow', oh: null, cloak: { color: 0x2f4f2a }, body: [0.97, 1.02, 0.97], quiver: true,
    base: { chest: 'leather_chest_0', legs: 'leather_legs_0', feet: 'leather_feet_0', hands: 'leather_hands_0' } },
  mage: { pants: PANTS.mage, tint: 0x8a96e8, emissive: 0x0e1440, hideVisor: true, head: { kind: 'wizardhat', color: 0x23358a }, mh: 'staff', oh: null, body: [0.92, 1.0, 0.92],
    base: { chest: 'cloth_chest_1', waist: 'cloth_waist_2', hands: 'cloth_hands_0' } },
};
// armor kind of an item id (cloth_head_3 -> cloth)
const kindOf = (id) => (id || '').split('_')[0];

// Builds the spec describing what a player model wears.
function playerSpec(cls, eq = {}, lk = {}) {
  const female = lk?.gender === 'f';
  const base = CLASS_LOOK[cls] || CLASS_LOOK.warrior;
  const spec = { tint: base.tint, pants: base.pants, props: [], scale: female ? 0.95 : 1, bodyScale: base.body, female, skin: SKIN[lk?.skin ?? 0], hair: HAIR[lk?.hair ?? 0] };
  const mh = eq.mainhand, oh = eq.offhand;
  const wk = (id, dflt) => (id ? id.split('_')[0] : dflt);
  const mhKind = wk(mh, base.mh);
  if (mhKind) spec.props.push({ kind: mhKind, id: mh || `${mhKind}_0`, slot: mhKind === 'bow' ? 'left' : 'right' });
  const ohKind = oh ? wk(oh, null) : base.oh;
  if (ohKind) spec.props.push({ kind: ohKind, id: oh || `${ohKind}_0`, slot: 'left' });
  if (eq.head) {
    const k = kindOf(eq.head);
    spec.head = { kind: k === 'mail' ? 'helmet' : k === 'leather' ? 'hood' : 'wizardhat', id: eq.head, color: k === 'cloth' ? 0x23358a : k === 'leather' ? 0x4a5a2a : undefined };
  } else if (base.head) spec.head = base.head;
  if (eq.shoulder || base.base?.shoulder) spec.shoulders = eq.shoulder || base.base.shoulder;
  if (eq.back) spec.cloak = { id: eq.back };
  else if (base.cloak) spec.cloak = { id: 'back_0', color: base.cloak.color };
  const bs = base.base || {};
  spec.chest = eq.chest || bs.chest; spec.legs = eq.legs || bs.legs; spec.feet = eq.feet || bs.feet; spec.hands = eq.hands || bs.hands; spec.waist = eq.waist || bs.waist;
  if (base.quiver) spec.extra = [{ bone: 'chest', obj: quiverProp(), off: [-0.08, 0.1, -0.2], rot: [0, 0, 0] }];
  if (eq.chest) {
    const t = parseInt(eq.chest.split('_').pop(), 10);
    if (eq.chest.endsWith('_bf') || t >= 7) spec.emissive = 0x3a2a08;
  }
  return spec;
}

// NPCs and monsters built from the same rig.
const HUMANOIDS = {
  forester: { tint: 0x5a8a4a, head: { kind: 'hood', color: 0x3a5a2a }, props: [{ kind: 'axe', slot: 'right' }] },
  trainer_warrior: { tint: 0x9a9aa8, props: [{ kind: 'sword', id: 'sword_2', slot: 'right' }, { kind: 'shield', id: 'shield_2', slot: 'left' }] },
  trainer_hunter: { tint: 0x5a8a45, hideVisor: true, head: { kind: 'hood', color: 0x2f4f2a }, props: [{ kind: 'bow', id: 'bow_2', slot: 'left' }] },
  trainer_mage: { tint: 0x7a3aa8, emissive: 0x180830, hideVisor: true, head: { kind: 'wizardhat', color: 0x4a2a7a }, props: [{ kind: 'staff', id: 'staff_2', slot: 'right' }] },
  merchant: { tint: 0xb8864a, hideVisor: true, head: { kind: 'straw' } },
  smith: { tint: 0x8a6a4a, props: [{ kind: 'pick', slot: 'right' }] },
  farmer: { tint: 0xb8704a, hideVisor: true, head: { kind: 'straw' }, props: [{ kind: 'pick', slot: 'right' }] },
  scholar: { tint: 0x3a9a6a, emissive: 0x08200f, hideVisor: true, head: { kind: 'wizardhat', color: 0x1f5a3a } },
  guard: { tint: 0x4a68a0, head: { kind: 'helmet' }, props: [{ kind: 'spear', slot: 'right' }] },
  huntsman: { tint: 0x7a5a3a, hideVisor: true, head: { kind: 'hood', color: 0x5a3a22 }, props: [{ kind: 'rifle', slot: 'right' }] },
  mayor: { tint: 0xc8a050, emissive: 0x2a1c04, hideVisor: true, head: { kind: 'crown' } },
  priest: { tint: 0xe8e8f4, emissive: 0x30304a, hideVisor: true, head: { kind: 'hood', color: 0xeeeeff }, props: [{ kind: 'staff', id: 'staff_1', slot: 'right' }] },
  elder: { tint: 0x8a8a7a, scale: 0.96, hideVisor: true, head: { kind: 'hood', color: 0x6a6a5a }, props: [{ kind: 'staff', id: 'staff_0', slot: 'right' }] },
  herbalist: { tint: 0x5a8a3a, hideVisor: true, head: { kind: 'hood', color: 0x3a6a2a } },
  fisher: { tint: 0x4a6a8a, hideVisor: true, head: { kind: 'straw' } },
  sheikh: { tint: 0xe8e0c8, emissive: 0x2a2418, hideVisor: true, head: { kind: 'turban', color: 0xe8e0c8 } },
  nomad: { tint: 0xc8a868, hideVisor: true, head: { kind: 'turban', color: 0xb89858 } },
  flightmaster: { tint: 0xc8a040, head: { kind: 'hood', color: 0x6a4a22 }, hideVisor: true, props: [{ kind: 'spear', slot: 'right' }] },
  // monsters
  bandit: { tint: 0x5a3a3a, hideVisor: true, head: { kind: 'mask' }, props: [{ kind: 'dagger', slot: 'right' }] },
  bandit_chief: { tint: 0x7a2a2a, scale: 1.05, hideVisor: true, head: { kind: 'mask', color: 0xaa2222 }, props: [{ kind: 'axe', slot: 'right' }] },
  desert_raider: { tint: 0xb89860, hideVisor: true, head: { kind: 'turban', color: 0x8a2a2a }, props: [{ kind: 'scimitar', slot: 'right' }] },
  skeleton: { tint: 0xe0dccc, emissive: 0x2a2a1a, hideVisor: true, props: [{ kind: 'sword', id: 'sword_0', slot: 'right' }] },
  zombie: { tint: 0x6a8a5a, hideVisor: true, slow: true, armsForward: true },
  zombie_giant: { tint: 0x5a7a4a, hideVisor: true, slow: true, armsForward: true },
  kobold: { tint: 0x9a6a3a, scale: 0.62, hideVisor: true, head: { kind: 'ears' }, props: [{ kind: 'pick', slot: 'right' }] },
  kobold_shaman: { tint: 0x7a3a9a, emissive: 0x200a30, scale: 0.62, hideVisor: true, head: { kind: 'ears', color: 0x7a3a9a }, props: [{ kind: 'witchstaff', slot: 'right' }] },
  kobold_boss: { tint: 0xb87a3a, scale: 0.7, hideVisor: true, head: { kind: 'horns' }, props: [{ kind: 'axe', slot: 'right' }] },
  kobold_king: { tint: 0xd8a838, emissive: 0x302200, scale: 0.7, hideVisor: true, head: { kind: 'crown' }, props: [{ kind: 'sword', id: 'sword_6', slot: 'right' }] },
  yeti: { tint: 0xe8f2ff, scale: 1.25, hideVisor: true, head: { kind: 'ears', color: 0xe8f2ff } },
  yeti_chief: { tint: 0xd8e8ff, emissive: 0x102030, scale: 1.25, hideVisor: true, head: { kind: 'horns' }, props: [{ kind: 'axe', slot: 'right', scale: 1.6 }] },
  witch: { tint: 0x6a3a8a, emissive: 0x150a20, hideVisor: true, head: { kind: 'witchhat' }, props: [{ kind: 'witchstaff', slot: 'right' }] },
  hag: { tint: 0x4a5a2a, emissive: 0x0a1a05, hideVisor: true, head: { kind: 'witchhat' }, props: [{ kind: 'witchstaff', slot: 'right' }] },
  mummy: { tint: 0xe8e0c0, hideVisor: true, slow: true, armsForward: true },
  tomb_guardian: { tint: 0xd8b038, emissive: 0x302200, hideVisor: true, head: { kind: 'pharaoh' }, props: [{ kind: 'spear', slot: 'right', scale: 1.1 }] },
  water_elem: { tint: 0x2a6ad8, emissive: 0x0a2a6a, opacity: 0.72, hideVisor: true, scale: 1.1 },
  water_elem_pet: { tint: 0x4aaaff, emissive: 0x1a5a9a, opacity: 0.72, hideVisor: true, scale: 0.95 },
};

// painted style, gender and proportions (bodyScale x,y,z ; headScale) of NPCs and monsters
const STYLE_OF = {
  skeleton: 'skeleton', zombie: 'zombie', zombie_giant: 'zombie', mummy: 'bandage', tomb_guardian: 'gold', yeti: 'fur', yeti_chief: 'fur',
  kobold: 'kobold', kobold_shaman: 'kobold', kobold_boss: 'kobold', kobold_king: 'kobold', water_elem: 'water', water_elem_pet: 'water',
};
const FEMALE_NPC = new Set(['priest', 'herbalist', 'witch', 'hag', 'scholar', 'merchant', 'trainer_mage', 'fisher']);
const BODY_OF = {
  guard: { b: [1.12, 1, 1.1] }, trainer_warrior: { b: [1.14, 1, 1.12] }, smith: { b: [1.2, 0.97, 1.2] }, mayor: { b: [1.25, 0.95, 1.3] }, elder: { b: [0.95, 0.92, 0.95] },
  bandit_chief: { b: [1.2, 1, 1.15] }, farmer: { b: [1.05, 1, 1.1] }, merchant: { b: [1.0, 1, 1.0] },
  kobold: { b: [1.1, 0.95, 1.1], h: 1.3 }, kobold_shaman: { b: [1.0, 0.95, 1.0], h: 1.3 }, kobold_boss: { b: [1.35, 0.95, 1.3], h: 1.25 }, kobold_king: { b: [1.4, 0.98, 1.35], h: 1.2 },
  yeti: { b: [1.55, 1.0, 1.35], h: 0.9 }, yeti_chief: { b: [1.7, 1.0, 1.45], h: 0.9 },
  skeleton: { b: [0.8, 1.0, 0.8] }, zombie_giant: { b: [1.6, 1.0, 1.45] },
  tomb_guardian: { b: [1.35, 1.0, 1.3] }, water_elem: { b: [1.3, 1.0, 1.3] }, water_elem_pet: { b: [1.1, 1.0, 1.1] }, hag: { b: [0.9, 0.92, 0.9] },
};
const QUAD = new Set(['wolf', 'wolf_alpha', 'icewolf', 'hyena', 'boar', 'bear', 'bear_big', 'polarbear', 'polarbear_big', 'bull', 'sheep']);

export function makeModel(ent) {
  const k = ent.k, m = ent.m;
  if (k === 'obj') return makeObject(m);
  if (QUAD.has(m)) return makePixelBeast(m, ent.id);
  switch (m) {
    case 'spider': return makeSpider(ent.id);
    case 'cavespider': return makeSpider(ent.id, 0x3a3a52);
    case 'spider_queen': return makeSpider(ent.id, 0x5a1a2e);
    case 'scorpion': return makeScorpion(ent.id);
    case 'scorpion_king': return makeScorpion(2);
    case 'croc': return makeCroc(ent.id);
    case 'wisp': return makeWisp(0x9aff4a);
    case 'frostspirit': return makeWisp(0x7ad8ff);
    case 'griffin': return makeGriffin();
  }
  if (k === 'player') {
    const model = makeHumanoid(playerSpec(ent.c, ent.eq || {}, ent.lk));
    model.playerClass = ent.c;
    return model;
  }
  const spec = HUMANOIDS[m] || { tint: 0x888888 };
  const bo = BODY_OF[m];
  const id = ent.id || 0;
  return makeHumanoid({
    ...spec, style: STYLE_OF[m], female: FEMALE_NPC.has(m), bodyScale: bo?.b, headScale: bo?.h,
    skin: SKIN[id % 4], hair: HAIR[(id >> 1) % 6],
  });
}
// Re-dresses a player model when the equipment changed.
export function applyEquipment(model, ent) {
  if (model.kind !== 'human' || !model.dress || !model.playerClass) return;
  const spec = playerSpec(ent.c, ent.eq || {}, ent.lk);
  model.setTint(spec.tint);
  model.dress(spec);
}

export function makePlayerModel(cls, look, eq = {}) {
  return makeModel({ k: 'player', c: cls, lk: look, eq, id: 1, m: 'human' });
}
