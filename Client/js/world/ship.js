import * as THREE from 'three';
import { pixTex } from './human.js';

// The ferry: a wooden pixel ship with a symmetric hull (it sails bow or stern first and never turns),
// a mast with a sail across the middle and an open deck you can walk on. The deck is 3 m wide and 12 m long
// (z = -6 .. 6, top at y = 1.2) so that passengers can step on from the pier.
const wood = (c) => new THREE.MeshLambertMaterial({ map: pixTex('sw' + c, 8, 4, (p) => { p.fill(c, 0.22); for (let y = 1; y < 4; y += 2) p.rect(0, y, 8, 1, 0x000000 + ((c >> 1) & 0x7f7f7f)); }) });
const cloth = (c) => new THREE.MeshLambertMaterial({ map: pixTex('sc' + c, 8, 8, (p) => { p.fill(c, 0.1); for (let y = 2; y < 8; y += 3) p.rect(0, y, 8, 1, 0x000000 + ((c >> 1) & 0x7f7f7f)); }), side: THREE.DoubleSide });

export function makeShip(style = 'cog') {
  const g = new THREE.Group();
  const junk = style === 'junk';
  const hullC = junk ? 0x8a2f22 : 0x6a4a2c, deckC = 0xa07a48, darkC = 0x4a3220;
  const add = (w, h, d, m, x, y, z) => { const b = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), m); b.position.set(x, y, z); b.castShadow = true; b.receiveShadow = true; g.add(b); return b; };
  // hull: stepped boxes narrowing towards both ends
  add(3.6, 1.2, 8, wood(hullC), 0, 0.55, 0);
  for (const s of [-1, 1]) {
    add(3.1, 1.0, 2.2, wood(hullC), 0, 0.65, s * 5.0);
    add(2.2, 0.9, 1.6, wood(hullC), 0, 0.85, s * 6.5);
    add(1.1, 0.7, 1.2, wood(hullC), 0, 1.05, s * 7.6);
    add(0.4, 1.4, 0.4, wood(darkC), 0, 1.6, s * 8.1); // stem / stern post
  }
  add(3.3, 0.16, 10, wood(deckC), 0, 1.16, 0);
  for (const s of [-1, 1]) add(2.7, 0.16, 3.6, wood(deckC), 0, 1.16, s * 5.6);
  // rails only at both ends: the middle stays open for boarding
  for (const sx of [-1, 1]) for (const sz of [-1, 1]) add(0.14, 0.6, 3.0, wood(hullC), sx * 1.6, 1.6, sz * 5.0);
  // mast in the middle and a long sail
  add(0.3, 9.5, 0.3, wood(darkC), 0, 5.9, 0);
  const sail = new THREE.Mesh(new THREE.PlaneGeometry(junk ? 7 : 6, 6.6), cloth(junk ? 0xc8402a : 0xf0ead8));
  sail.position.set(0, 6.3, 0); sail.rotation.y = Math.PI / 2; sail.castShadow = true; g.add(sail);
  add(0.14, 0.14, junk ? 7 : 6, wood(darkC), 0, 9.5, 0);
  add(0.14, 0.14, junk ? 7 : 6, wood(darkC), 0, 3.1, 0);
  const flag = add(0.05, 0.9, 1.4, cloth(junk ? 0xe8c040 : 0xb02a2a), 0, 10.9, -0.7);
  // crates, barrels and ropes as decoration at the sides (deck stays free in the middle)
  add(0.8, 0.7, 0.8, wood(0x7a5a38), 1.2, 1.6, 4.4);
  add(0.7, 0.6, 0.7, wood(0x7a5a38), -1.2, 1.55, -4.6);
  add(0.6, 0.8, 0.6, wood(0x5a3a22), -1.2, 1.65, 4.6);
  return {
    obj: g, height: 4, kind: 'ship', hitR: 7,
    attack() {},
    update() {
      const t = performance.now() / 1000;
      g.position.y = Math.sin(t * 1.4) * 0.04;
      g.rotation.z = Math.sin(t * 1.1) * 0.012;
      flag.rotation.y = Math.sin(t * 6) * 0.3;
    },
  };
}
