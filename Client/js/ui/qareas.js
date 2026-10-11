import { S } from '../state.js';

// Map areas of the active quests: where to go to finish the open objectives, or where to hand the quest in.
// Returns [{ x, z, r, k, qid, title, text }]
export function getQuestAreas() {
  const out = [];
  for (const q of S.quests) {
    if (q.status !== 'active' && q.status !== 'complete') continue;
    const def = S.defs?.quests?.[q.id];
    if (!def) continue;
    if (q.status === 'complete') {
      if (def.handIn) out.push({ ...def.handIn, qid: q.id, title: def.title, text: 'Quest abgeben', turnIn: true });
      continue;
    }
    (def.objectives || []).forEach((o, i) => {
      if ((q.counts?.[i] ?? 0) >= o.count) return;
      for (const a of def.areas?.[i] || []) out.push({ ...a, qid: q.id, title: def.title, text: o.text });
    });
  }
  return out;
}
