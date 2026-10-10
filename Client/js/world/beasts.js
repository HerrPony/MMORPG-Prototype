import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import * as SkeletonUtils from 'three/addons/utils/SkeletonUtils.js';

// Four-legged canines use the animated Fox model (Khronos glTF sample, CC-BY 4.0): Survey / Walk / Run.

let fox = null, ready = null;

export function loadBeastAssets() {
  if (!ready) {
    ready = new Promise((resolve) => {
      new GLTFLoader().load('assets/Fox.glb', (g) => {
        fox = { gltf: g, len: 1, dims: new THREE.Vector3(1, 1, 1) };
        // body length of the bind pose (skinned bounding boxes ignore the node scale, so measure the geometry)
        const box = new THREE.Box3();
        g.scene.traverse((o) => { if (o.isMesh) { o.geometry.computeBoundingBox(); box.union(o.geometry.boundingBox); } });
        box.getSize(fox.dims);
        fox.len = Math.max(fox.dims.x, fox.dims.y, fox.dims.z);
        fox.up = fox.dims.y >= fox.dims.x && fox.dims.y >= fox.dims.z;
        resolve(true);
      }, undefined, () => resolve(false));
    });
  }
  return ready;
}

export const beastsLoaded = () => !!fox;

// Wolf-like variants: colour, body length (m)
const VARIANTS = {
  wolf: { color: 0x6c7078, len: 1.5 },
  wolf_alpha: { color: 0x2c2e36, len: 2.0, eye: 0xff4020 },
  icewolf: { color: 0xcfe2f4, len: 1.65, emissive: 0x20303c },
  hyena: { color: 0xa88a58, len: 1.35, emissive: 0x181008 },
};

export function makeFoxBeast(kind = 'wolf') {
  const V = VARIANTS[kind] || VARIANTS.wolf;
  const root = new THREE.Group();
  const inst = SkeletonUtils.clone(fox.gltf.scene);
  const sc = (V.len / fox.len) * 2.2; // tuned by eye: the skinned fox renders smaller than its geometry box suggests
  inst.scale.setScalar(sc);
  const pivot = new THREE.Group();
  pivot.add(inst); root.add(pivot);
  inst.traverse((o) => {
    if (o.isMesh) {
      o.castShadow = true; o.frustumCulled = false;
      o.material = o.material.clone();
      o.material.map = null;
      o.material.color.set(V.color);
      if (V.emissive) o.material.emissive = new THREE.Color(V.emissive);
      o.material.needsUpdate = true;
    }
  });
  const mixer = new THREE.AnimationMixer(inst);
  const find = (n) => fox.gltf.animations.find((a) => a.name === n);
  const act = { idle: mixer.clipAction(find('Survey')), walk: mixer.clipAction(find('Walk')), run: mixer.clipAction(find('Run')) };
  Object.values(act).forEach((a) => { a.play(); a.setEffectiveWeight(0); });
  act.idle.setEffectiveWeight(1);
  const w = { idle: 1, walk: 0, run: 0 };
  let attackT = 0, dead = 0;
  const h = V.len * 0.55;
  return {
    obj: root, height: h + 0.3, kind: 'wolf',
    attack() { attackT = 0.01; },
    update(dt, st) {
      const sp = st.speed;
      const tIdle = sp < 0.08 ? 1 : 0, tWalk = sp >= 0.08 && sp < 0.45 ? 1 : 0, tRun = sp >= 0.45 ? 1 : 0;
      const k = Math.min(1, dt * 10);
      w.idle += (tIdle - w.idle) * k; w.walk += (tWalk - w.walk) * k; w.run += (tRun - w.run) * k;
      act.idle.setEffectiveWeight(w.idle); act.walk.setEffectiveWeight(w.walk); act.run.setEffectiveWeight(w.run);
      act.walk.setEffectiveTimeScale(Math.max(0.6, sp / 0.3));
      act.run.setEffectiveTimeScale(Math.max(0.8, sp / 0.9));
      mixer.update(dt);
      let lunge = 0;
      if (attackT > 0) {
        attackT += dt;
        const t = Math.min(1, attackT / 0.4);
        lunge = Math.sin(t * Math.PI) * 0.5;
        if (t >= 1) attackT = 0;
      }
      pivot.position.z = lunge;
      dead += ((st.dead ? 1 : 0) - dead) * Math.min(1, dt * 6);
      root.rotation.z = dead * (Math.PI / 2);
      pivot.position.x = dead * 0.3;
      mixer.timeScale = st.dead && w.idle > 0.9 ? 0 : 1;
    },
  };
}
