// Procedural orchestral music and sound effects (WebAudio, no asset files).
//
// A small composition engine plays an endless score: every biome has its own theme (key, mode, tempo, chord
// progression, instrumentation, rhythm), the menu has a heroic main theme, dungeons a dark one. Melodies are built
// from a motif that is developed over 8-bar phrases (A A' B A'') instead of random notes. Time of day changes the
// orchestration (flute and strings by day, harp, bells and choir by night) and combat adds timpani/taiko, brass
// stabs and a driving string ostinato on top of the zone's theme.

const NOTE = (n) => 440 * Math.pow(2, (n - 69) / 12); // MIDI -> Hz

// ---- scales (semitone offsets) ----
const MODES = {
  major: [0, 2, 4, 5, 7, 9, 11], minor: [0, 2, 3, 5, 7, 8, 10], dorian: [0, 2, 3, 5, 7, 9, 10], lydian: [0, 2, 4, 6, 7, 9, 11],
  mixolydian: [0, 2, 4, 5, 7, 9, 10], phrygian: [0, 1, 3, 5, 7, 8, 10], harmonic: [0, 2, 3, 5, 7, 8, 11], phrygDom: [0, 1, 4, 5, 7, 8, 10],
  pentatonic: [0, 2, 4, 7, 9], minPent: [0, 3, 5, 7, 10], insen: [0, 1, 5, 7, 10], // insen = Japanese scale
};

// ---- themes ----
// root: MIDI tonic, mode, bpm, prog: chord roots as scale degrees (0-based) per 2 bars (4 entries = one 8-bar cycle),
// motif: [scale steps relative to the previous note, length in beats] ..., comp: accompaniment pattern, perc: drum style,
// lead / pad / pluck: instruments, night: how the night variant differs
const THEMES = {
  menu: { root: 50, mode: 'minor', bpm: 84, prog: [0, 5, 2, 6], motif: [[0, 1.5], [2, 0.5], [1, 1], [-1, 1], [3, 1.5], [-1, 0.5], [-1, 1], [-2, 1]], comp: 'arp8', perc: 'epic', lead: 'horn', pad: 'strings', pluck: 'harp', heroic: true },
  forest: { root: 50, mode: 'dorian', bpm: 78, prog: [0, 3, 6, 4], motif: [[0, 1], [2, 1], [1, 0.5], [-1, 0.5], [2, 1], [-2, 1], [1, 2], [0, 1]], comp: 'arp8', perc: 'frame', lead: 'flute', pad: 'strings', pluck: 'harp' },
  plains: { root: 55, mode: 'lydian', bpm: 92, prog: [0, 4, 5, 3], motif: [[0, 1], [1, 0.5], [1, 0.5], [2, 1], [-2, 1], [1, 1], [-1, 1], [0, 2]], comp: 'strum', perc: 'frame', lead: 'fiddle', pad: 'strings', pluck: 'guitar' },
  swamp: { root: 49, mode: 'phrygian', bpm: 58, prog: [0, 1, 0, 6], motif: [[0, 2], [1, 1], [-1, 1], [-1, 2], [2, 1], [-1, 1]], comp: 'drone', perc: 'heart', lead: 'flute', pad: 'choir', pluck: 'bell', eerie: true },
  snow: { root: 52, mode: 'minor', bpm: 66, prog: [0, 5, 2, 4], motif: [[0, 2], [4, 1], [-2, 1], [-1, 1], [-1, 1], [3, 2]], comp: 'bells', perc: 'none', lead: 'celesta', pad: 'choir', pluck: 'bell' },
  desert: { root: 50, mode: 'phrygDom', bpm: 96, prog: [0, 1, 0, 3], motif: [[0, 0.75], [1, 0.25], [1, 0.5], [-2, 0.5], [3, 1], [-1, 0.5], [-1, 0.5], [-1, 1]], comp: 'ostinato', perc: 'darbuka', lead: 'oud', pad: 'strings', pluck: 'oud' },
  jungle: { root: 53, mode: 'mixolydian', bpm: 104, prog: [0, 6, 3, 0], motif: [[0, 0.5], [2, 0.5], [-1, 0.5], [2, 0.5], [1, 1], [-3, 1], [2, 1]], comp: 'marimba', perc: 'tribal', lead: 'flute', pad: 'strings', pluck: 'marimba' },
  blossom: { root: 57, mode: 'insen', bpm: 72, prog: [0, 2, 1, 3], motif: [[0, 1], [1, 1], [1, 2], [-2, 1], [1, 1], [-1, 2]], comp: 'koto', perc: 'none', lead: 'flute', pad: 'strings', pluck: 'koto' },
  mesa: { root: 52, mode: 'phrygian', bpm: 88, prog: [0, 6, 5, 6], motif: [[0, 1], [0, 0.5], [1, 0.5], [-1, 1], [2, 1], [-3, 1], [1, 1]], comp: 'strum', perc: 'taiko', lead: 'horn', pad: 'strings', pluck: 'guitar' },
  volcano: { root: 45, mode: 'harmonic', bpm: 72, prog: [0, 5, 0, 4], motif: [[0, 1.5], [0, 0.5], [1, 1], [-1, 1], [3, 1], [-2, 1], [-1, 2]], comp: 'ostinato', perc: 'taiko', lead: 'horn', pad: 'choir', pluck: 'harp', dark: true },
  sea: { root: 55, mode: 'major', bpm: 84, prog: [0, 3, 4, 0], motif: [[0, 1], [2, 1], [1, 1], [-1, 1], [2, 2], [-2, 1], [1, 1]], comp: 'arp8', perc: 'frame', lead: 'flute', pad: 'strings', pluck: 'harp' },
  dungeon: { root: 43, mode: 'harmonic', bpm: 64, prog: [0, 1, 0, 5], motif: [[0, 2], [1, 1], [-1, 1], [2, 2]], comp: 'drone', perc: 'heart', lead: 'horn', pad: 'choir', pluck: 'bell', dark: true, eerie: true },
};

// small deterministic PRNG so that a given phrase always sounds the same
function mulberry(a) { return () => { a |= 0; a = (a + 0x6d2b79f5) | 0; let t = Math.imul(a ^ (a >>> 15), 1 | a); t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t; return ((t ^ (t >>> 14)) >>> 0) / 4294967296; }; }
const hashStr = (s) => { let h = 7; for (const c of s) h = (h * 31 + c.charCodeAt(0)) | 0; return h; };

export class Audio {
  constructor() {
    this.ctx = null;
    this.mode = 'menu';       // 'menu' | 'world' | 'dungeon'
    this.zone = 'forest';     // biome key while in the open world
    this.night = 0;           // 0 day .. 1 night
    this.intensity = 0;       // 0..1 combat level
    this.muted = false;
    this.vol = 0.6;
    try { const v = localStorage.getItem('vol'); if (v !== null) this.vol = +v; this.muted = localStorage.getItem('muted') === '1'; } catch { /* ignore */ }
    this.theme = null; this.themeKey = ''; this.barT = 0; this.bar = 0; this.phrase = 0;
  }

  // Must be called from a user gesture (browsers block audio until then).
  start() {
    if (this.ctx) { if (this.ctx.state === 'suspended') this.ctx.resume(); return; }
    const AC = window.AudioContext || window.webkitAudioContext;
    if (!AC) return;
    this.init(new AC());
    setInterval(() => this.schedule(), 100);
  }

  // builds the signal chain on any (also offline) context
  init(c) {
    this.ctx = c;
    this.master = c.createGain();
    this.master.gain.value = this.muted ? 0 : this.vol;
    const comp = c.createDynamicsCompressor();
    comp.threshold.value = -16; comp.knee.value = 14; comp.ratio.value = 4; comp.attack.value = 0.01; comp.release.value = 0.3;
    this.master.connect(comp); comp.connect(c.destination);
    // concert hall reverb: convolution with a synthetic impulse response
    const len = Math.floor(c.sampleRate * 3.2), ir = c.createBuffer(2, len, c.sampleRate);
    for (let ch = 0; ch < 2; ch++) { const d = ir.getChannelData(ch); for (let i = 0; i < len; i++) { const t = i / len; d[i] = (Math.random() * 2 - 1) * Math.pow(1 - t, 2.4) * (i < 600 ? i / 600 : 1); } }
    this.reverb = c.createConvolver(); this.reverb.buffer = ir;
    this.wet = c.createGain(); this.wet.gain.value = 0.42;
    this.reverb.connect(this.wet); this.wet.connect(this.master);
    this.musicBus = c.createGain(); this.musicBus.gain.value = 0.8;      // dry music
    this.musicBus.connect(this.master); this.musicBus.connect(this.reverb);
    this.drumBus = c.createGain(); this.drumBus.gain.value = 1;          // drums: mostly dry, a little hall
    const dsend = c.createGain(); dsend.gain.value = 0.25;
    this.drumBus.connect(this.master); this.drumBus.connect(dsend); dsend.connect(this.reverb);
    this.sfxBus = c.createGain(); this.sfxBus.gain.value = 0.7; this.sfxBus.connect(this.master);
    this.barT = c.currentTime + 0.3;
  }

  setMode(m) { this.mode = m; }
  setVolume(v) {
    this.vol = v;
    try { localStorage.setItem('vol', String(v)); } catch { /* ignore */ }
    this.apply();
  }
  toggleMute() {
    this.muted = !this.muted;
    try { localStorage.setItem('muted', this.muted ? '1' : '0'); } catch { /* ignore */ }
    this.apply();
    return this.muted;
  }
  apply() {
    if (!this.ctx) return;
    this.master.gain.setTargetAtTime(this.muted ? 0 : this.vol, this.ctx.currentTime, 0.1);
  }

  wantedTheme() { return this.mode === 'menu' ? 'menu' : this.mode === 'dungeon' ? 'dungeon' : (THEMES[this.zone] ? this.zone : 'forest'); }

  // ---- composition: schedule whole bars ahead of time ----
  schedule() {
    const c = this.ctx;
    if (!c || c.state !== 'running') return;
    while (this.barT < c.currentTime + 1.2) this.playBar();
  }

  playBar() {
    const key = this.wantedTheme();
    if (key !== this.themeKey) { // theme change: start a fresh 8-bar cycle at this bar line
      this.themeKey = key; this.theme = THEMES[key]; this.bar = 0; this.phrase = 0;
      this.musicBus.gain.cancelScheduledValues(this.barT);
      this.musicBus.gain.setValueAtTime(0.35, this.barT); this.musicBus.gain.linearRampToValueAtTime(0.8, this.barT + 2.5);
    }
    const th = this.theme, beat = 60 / th.bpm, t0 = this.barT, barLen = beat * 4;
    const night = this.night > 0.5, fight = this.world() && this.intensity > 0.3, hot = this.intensity;
    const barIn = this.bar % 8, chordSlot = Math.floor(barIn / 2);
    const scale = MODES[th.mode], deg = th.prog[chordSlot];
    const rootDeg = deg;
    const pitch = (d, oct = 0) => { // scale degree (can be negative / > 6) -> MIDI
      const n = scale.length; const o = Math.floor(d / n), i = ((d % n) + n) % n;
      return th.root + scale[i] + 12 * (o + oct);
    };
    const chord = [rootDeg, rootDeg + 2, rootDeg + 4].map((d) => pitch(d)); // triad on the scale
    const rng = mulberry(hashStr(this.themeKey) + this.phrase * 977 + barIn * 31);

    // 1. strings / choir pad: swell on every chord change (2 bars)
    if (barIn % 2 === 0) {
      const len = barLen * 2;
      const voice = night && th.pad === 'strings' ? 'choir' : th.pad;
      for (const [i, n] of chord.entries()) this.padVoice(voice, NOTE(n + (i === 2 ? 0 : 0)), t0, len, 0.07 - (fight ? 0.01 : 0));
      this.padVoice(voice, NOTE(chord[0] - 12), t0, len, 0.05);
    }
    // 2. bass line: root and fifth, two notes per bar
    this.bassNote(NOTE(chord[0] - 12), t0, beat * 2, th.dark ? 0.2 : 0.15);
    this.bassNote(NOTE(chord[chordSlot % 2 ? 0 : 2] - 12), t0 + beat * 2, beat * 2, th.dark ? 0.15 : 0.1);
    // 3. accompaniment pattern
    this.accompany(th, chord, t0, beat, night, rng);
    // 4. melody: 8-bar phrase built from the motif (A A' B A''), rests in some bars to let the music breathe
    const form = this.phrase % 4, plays = (barIn % 2 === 0 || form !== 1) && !(barIn >= 6 && form === 3 && rng() < 0.4);
    if (plays) this.melody(th, scale, pitch, chordSlot, barIn, form, t0, beat, night, fight, rng);
    // 5. second voice (counter-melody) on later phrases, not at night during quiet phases
    if (form >= 2 && barIn % 2 === 1 && !fight) this.leadNote(night ? 'celesta' : th.lead, NOTE(pitch(rootDeg + 4, 2)), t0 + beat * 2, beat * 3, 0.07);
    // 6. percussion
    this.percussion(th, t0, beat, barIn, night, fight, hot);
    // 7. combat layer: brass stabs on the chord change and a driving string ostinato
    if (fight) this.combatLayer(th, chord, t0, beat, barIn, hot);

    this.barT += barLen; this.bar++;
    if (this.bar % 8 === 0) this.phrase++;
  }
  world() { return this.mode !== 'menu'; }

  melody(th, scale, pitch, chordSlot, barIn, form, t0, beat, night, fight, rng) {
    // the motif is played over bars 0-1 of every 2-bar slot and transposed with the chord (a "sequence"); in form 2 it is inverted
    const motif = th.motif, rootDeg = th.prog[chordSlot];
    let d = rootDeg + 2 + (form === 2 ? 2 : 0), t = 0;
    const lead = night ? (th.lead === 'horn' ? 'flute' : th.lead === 'flute' ? 'celesta' : th.lead) : th.lead;
    const vol = (night ? 0.1 : 0.14) * (th.heroic ? 1.15 : 1);
    const halfBar = barIn % 2; // each motif spans two bars; play the half that belongs to this bar
    const total = motif.reduce((s, m) => s + m[1], 0), scaleLen = 8 / total; // fit the motif into 8 beats
    let acc = 0;
    for (const [step, len] of motif) {
      const dur = len * scaleLen, start = acc; acc += dur;
      d += form === 2 ? -step : step;
      if (Math.floor(start / 4) !== halfBar) continue;
      if (rng() < 0.08 && start % 4 !== 0) continue; // occasional rests
      // strong beats land on chord tones: snap to the nearest of root/third/fifth
      let dd = d;
      if (start % 2 === 0) { const ct = [rootDeg, rootDeg + 2, rootDeg + 4, rootDeg + 7]; dd = ct.reduce((b, x) => (Math.abs(x - d) < Math.abs(b - d) ? x : b), ct[0]); }
      // phrase ending: last bar resolves to the tonic
      if (barIn === 7 && start + dur >= 7.9) dd = form === 3 ? 0 : 4;
      const f = NOTE(pitch(dd, 1));
      this.leadNote(lead, f, t0 + (start % 4) * beat, dur * beat * (lead === 'celesta' || lead === 'oud' ? 1.2 : 0.95), vol * (0.8 + rng() * 0.3));
    }
  }

  accompany(th, chord, t0, beat, night, rng) {
    const notes = [chord[0], chord[1], chord[2], chord[0] + 12, chord[1] + 12, chord[2] + 12];
    const pl = night && th.pluck === 'guitar' ? 'harp' : th.pluck;
    const v = night ? 0.06 : 0.08;
    switch (th.comp) {
      case 'arp8': for (let i = 0; i < 8; i++) this.pluckNote(pl, NOTE(notes[[0, 1, 2, 3, 2, 1, 4, 1][i]]), t0 + i * beat / 2, beat, v * (i % 2 ? 0.7 : 1)); break;
      case 'strum': for (const [i, off] of [0, 1, 1.5, 2.5, 3].entries()) for (let k = 0; k < 3; k++) this.pluckNote(pl, NOTE(chord[k] + (i % 2 ? 12 : 0)), t0 + off * beat + k * 0.018, beat * 1.2, v * 0.75); break;
      case 'ostinato': for (let i = 0; i < 8; i++) this.pluckNote(pl, NOTE(notes[[0, 2, 0, 2, 1, 2, 0, 2][i]] - (i % 4 === 0 ? 12 : 0)), t0 + i * beat / 2, beat * 0.6, v * (i % 4 === 0 ? 1.2 : 0.7)); break;
      case 'marimba': for (let i = 0; i < 8; i++) if (i % 3 !== 2 || rng() < 0.5) this.pluckNote('marimba', NOTE(notes[[0, 2, 1, 3, 2, 4, 1, 3][i]]), t0 + i * beat / 2, beat * 0.5, v); break;
      case 'bells': for (const off of [0, 1.5, 2.5]) this.pluckNote('bell', NOTE(notes[[3, 4, 5][Math.floor(rng() * 3)]]), t0 + off * beat, beat * 3, v * 0.8); break;
      case 'koto': for (const off of [0, 1, 2, 3]) if (rng() < 0.75) this.pluckNote('koto', NOTE(notes[[0, 1, 2, 3, 4][Math.floor(rng() * 5)]]), t0 + off * beat + (rng() < 0.3 ? beat / 2 : 0), beat * 2, v); break;
      case 'drone': if (this.bar % 2 === 0) this.pluckNote('bell', NOTE(chord[0] + 12), t0 + beat * 1.5, beat * 6, v * 0.7); break;
      default: break;
    }
  }

  percussion(th, t0, beat, barIn, night, fight, hot) {
    const p = th.perc, soft = night ? 0.5 : 1;
    const hit = (kind, off, vol) => this.perc(kind, t0 + off * beat, vol);
    if (fight) return; // combat has its own drums
    switch (p) {
      case 'frame': if (!night) { hit('frame', 0, 0.22); hit('frame', 2, 0.14); hit('shaker', 1, 0.1); hit('shaker', 3, 0.1); } else if (barIn % 2 === 0) hit('frame', 0, 0.12); break;
      case 'darbuka': for (const [o, k, v] of [[0, 'dum', 0.3], [1, 'tek', 0.18], [1.5, 'tek', 0.14], [2, 'dum', 0.24], [3, 'tek', 0.18], [3.5, 'tek', 0.12]]) hit(k, o, v * soft); break;
      case 'tribal': for (const [o, k, v] of [[0, 'dum', 0.3], [0.75, 'tom', 0.18], [1.5, 'tom', 0.2], [2, 'dum', 0.26], [2.75, 'tom', 0.16], [3.5, 'tom', 0.2]]) hit(k, o, v * soft); break;
      case 'taiko': if (barIn % 2 === 0) { hit('taiko', 0, 0.45 * soft); hit('taiko', 2.5, 0.28 * soft); } break;
      case 'heart': hit('dum', 0, 0.28); hit('dum', 0.6, 0.18); break;
      case 'epic': if (barIn % 4 === 0) { hit('timp', 0, 0.5); hit('timp', 3, 0.3); } else if (barIn % 2 === 0) hit('timp', 0, 0.3); if (barIn === 7) for (let i = 0; i < 8; i++) hit('timp', 2 + i * 0.25, 0.1 + i * 0.03); break;
      default: break;
    }
  }

  combatLayer(th, chord, t0, beat, barIn, hot) {
    const v = 0.5 + hot * 0.5;
    // war drums: driving eighths with accents, a timpani on the downbeat
    for (let i = 0; i < 8; i++) this.perc(i % 4 === 0 ? 'taiko' : i % 2 ? 'tom' : 'dum', t0 + i * beat / 2, (i % 4 === 0 ? 0.55 : 0.2) * v);
    if (hot > 0.6) for (let i = 0; i < 4; i++) this.perc('tom', t0 + beat * 3.5 + i * beat / 8, 0.12 * v);
    // string ostinato (spiccato sixteenth pattern on root and fifth)
    for (let i = 0; i < 16; i++) this.stringStab(NOTE(chord[i % 4 === 3 ? 2 : 0] + 12), t0 + i * beat / 4, beat / 5, 0.045 * v * (i % 4 === 0 ? 1.4 : 1));
    // brass stabs on the chord change
    if (barIn % 2 === 0) { for (const n of chord) this.brass(NOTE(n), t0, beat * 1.4, 0.07 * v); this.brass(NOTE(chord[0] - 12), t0, beat * 2, 0.1 * v); }
    if (barIn % 2 === 1) this.brass(NOTE(chord[2]), t0 + beat * 2, beat * 0.6, 0.06 * v);
  }

  // ---- instruments ----
  osc(type, freq, t, dur, vol, attack, release, dest, detune = 0) {
    const c = this.ctx, o = c.createOscillator(), g = c.createGain();
    o.type = type; o.frequency.value = freq; o.detune.value = detune;
    g.gain.setValueAtTime(0, t); g.gain.linearRampToValueAtTime(vol, t + attack);
    g.gain.setValueAtTime(vol, Math.max(t + attack, t + dur - release)); g.gain.linearRampToValueAtTime(0, t + dur + release);
    o.connect(g); g.connect(dest);
    o.start(t); o.stop(t + dur + release + 0.05);
    return { o, g };
  }

  filtered(type, freq, q, dest) { const f = this.ctx.createBiquadFilter(); f.type = type; f.frequency.value = freq; f.Q.value = q; f.connect(dest); return f; }

  vibrato(o, t, depth = 6, rate = 5.2, delay = 0.35) {
    const c = this.ctx, l = c.createOscillator(), g = c.createGain();
    l.frequency.value = rate; g.gain.setValueAtTime(0, t); g.gain.linearRampToValueAtTime(depth, t + delay + 0.5);
    l.connect(g); g.connect(o.detune); l.start(t); l.stop(t + 12);
  }

  padVoice(kind, f, t, dur, vol) {
    const c = this.ctx;
    if (kind === 'choir') { // sung "ooh/aah": sawtooth through two vowel formants
      const out = c.createGain(); out.gain.value = 1; out.connect(this.musicBus);
      const f1 = this.filtered('bandpass', 650, 5, out), f2 = this.filtered('bandpass', 1080, 6, out);
      for (const det of [-9, 0, 9]) { const o = this.osc('sawtooth', f, t, dur, vol * 0.9, 1.6, 1.8, f1, det); const g2 = c.createGain(); g2.gain.value = 0.6; o.g.connect(g2); g2.connect(f2); this.vibrato(o.o, t, 7, 5.5); }
      return;
    }
    // string section: detuned sawtooth ensemble, low-passed, slow bow attack
    const lp = this.filtered('lowpass', 1500, 0.4, this.musicBus);
    lp.frequency.setValueAtTime(500, t); lp.frequency.linearRampToValueAtTime(1700, t + Math.min(2.5, dur * 0.5));
    for (const det of [-11, 0, 11]) { const o = this.osc('sawtooth', f, t, dur, vol * 0.5, 1.3, 1.6, lp, det); this.vibrato(o.o, t, 5, 5); }
  }

  bassNote(f, t, dur, vol) { this.osc('sine', f, t, dur, vol, 0.04, 0.4, this.musicBus); this.osc('triangle', f * 2, t, dur, vol * 0.35, 0.04, 0.3, this.musicBus); }

  leadNote(kind, f, t, dur, vol) {
    const c = this.ctx;
    switch (kind) {
      case 'flute': { // breathy sine with vibrato and a puff of noise
        const o = this.osc('sine', f, t, dur, vol, 0.09, 0.18, this.musicBus); this.vibrato(o.o, t, 7, 5.4, 0.25);
        this.osc('triangle', f * 2, t, dur, vol * 0.18, 0.1, 0.15, this.musicBus);
        this.breath(t, Math.min(dur, 0.35), vol * 0.25, f * 2.5); break;
      }
      case 'horn': { // french horn: saw through an opening low-pass
        const lp = this.filtered('lowpass', 700, 0.8, this.musicBus); lp.frequency.setValueAtTime(380, t); lp.frequency.linearRampToValueAtTime(1500, t + 0.25); lp.frequency.linearRampToValueAtTime(900, t + dur);
        const o = this.osc('sawtooth', f, t, dur, vol * 0.9, 0.12, 0.25, lp); this.vibrato(o.o, t, 5, 5, 0.4);
        this.osc('sawtooth', f, t, dur, vol * 0.5, 0.12, 0.25, lp, 8); break;
      }
      case 'fiddle': { const lp = this.filtered('lowpass', 2600, 0.7, this.musicBus); const o = this.osc('sawtooth', f, t, dur, vol * 0.7, 0.07, 0.15, lp); this.vibrato(o.o, t, 12, 6, 0.2); this.osc('sawtooth', f, t, dur, vol * 0.4, 0.07, 0.15, lp, 7); break; }
      case 'celesta': this.pluckNote('bell', f * 2, t, dur, vol * 1.2); break;
      case 'oud': this.pluckNote('oud', f, t, dur, vol * 1.3); break;
      default: this.osc('triangle', f, t, dur, vol, 0.05, 0.2, this.musicBus);
    }
    void c;
  }

  breath(t, dur, vol, freq) {
    const c = this.ctx, buf = c.createBuffer(1, Math.ceil(c.sampleRate * dur), c.sampleRate), d = buf.getChannelData(0);
    for (let i = 0; i < d.length; i++) d[i] = (Math.random() * 2 - 1) * Math.sin(Math.PI * i / d.length);
    const s = c.createBufferSource(); s.buffer = buf; const g = c.createGain(); g.gain.value = vol;
    s.connect(this.filtered('bandpass', freq, 1.2, g)); g.connect(this.musicBus); s.start(t);
  }

  pluckNote(kind, f, t, dur, vol) {
    const decay = { harp: 1.6, guitar: 1.0, bell: 2.6, marimba: 0.4, koto: 1.4, oud: 0.9 }[kind] || 1;
    const c = this.ctx;
    const o = c.createOscillator(), g = c.createGain();
    o.type = kind === 'bell' ? 'sine' : kind === 'marimba' ? 'sine' : 'triangle'; o.frequency.value = f;
    g.gain.setValueAtTime(0, t); g.gain.linearRampToValueAtTime(vol, t + 0.006); g.gain.exponentialRampToValueAtTime(0.0008, t + Math.min(decay, dur + 0.5));
    o.connect(g);
    const lp = this.filtered('lowpass', kind === 'guitar' || kind === 'oud' ? 2400 : 4800, 0.5, this.musicBus); g.connect(lp);
    o.start(t); o.stop(t + decay + 0.1);
    // overtones: harp/koto bright octave, bell inharmonic partials, marimba a 4th-octave click
    const parts = { harp: [[2, 0.3]], guitar: [[2, 0.35], [3, 0.15]], koto: [[2, 0.4], [3, 0.25]], oud: [[2, 0.45], [3, 0.2]], bell: [[2.76, 0.35], [5.4, 0.15], [8.9, 0.07]], marimba: [[4, 0.25]] }[kind] || [];
    for (const [m, a] of parts) { const o2 = c.createOscillator(), g2 = c.createGain(); o2.type = 'sine'; o2.frequency.value = f * m; g2.gain.setValueAtTime(vol * a, t); g2.gain.exponentialRampToValueAtTime(0.0005, t + decay * (kind === 'marimba' ? 0.15 : 0.6)); o2.connect(g2); g2.connect(lp); o2.start(t); o2.stop(t + decay); }
  }

  stringStab(f, t, dur, vol) { const lp = this.filtered('lowpass', 2400, 0.5, this.musicBus); this.osc('sawtooth', f, t, dur, vol, 0.008, 0.05, lp); this.osc('sawtooth', f, t, dur, vol * 0.8, 0.008, 0.05, lp, 9); }
  brass(f, t, dur, vol) {
    const lp = this.filtered('lowpass', 800, 1, this.musicBus); lp.frequency.setValueAtTime(500, t); lp.frequency.linearRampToValueAtTime(2200, t + 0.08); lp.frequency.exponentialRampToValueAtTime(900, t + dur);
    this.osc('sawtooth', f, t, dur, vol, 0.03, 0.2, lp); this.osc('square', f, t, dur, vol * 0.4, 0.03, 0.2, lp, -6);
  }

  // drums
  perc(kind, t, vol) {
    const c = this.ctx;
    const drum = (f0, f1, dur, v, bus = this.drumBus) => { const o = c.createOscillator(), g = c.createGain(); o.type = 'sine'; o.frequency.setValueAtTime(f0, t); o.frequency.exponentialRampToValueAtTime(f1, t + dur * 0.8); g.gain.setValueAtTime(v, t); g.gain.exponentialRampToValueAtTime(0.001, t + dur); o.connect(g); g.connect(bus); o.start(t); o.stop(t + dur + 0.05); };
    const noise = (dur, freq, v, type = 'bandpass', q = 1) => { const buf = c.createBuffer(1, Math.ceil(c.sampleRate * dur), c.sampleRate), d = buf.getChannelData(0); for (let i = 0; i < d.length; i++) d[i] = (Math.random() * 2 - 1) * (1 - i / d.length); const s = c.createBufferSource(); s.buffer = buf; const g = c.createGain(); g.gain.value = v; s.connect(this.filtered(type, freq, q, g)); g.connect(this.drumBus); s.start(t); };
    switch (kind) {
      case 'timp': drum(110, 62, 1.1, vol * 0.8); drum(220, 110, 0.4, vol * 0.25); noise(0.06, 900, vol * 0.2); break;
      case 'taiko': drum(150, 48, 0.8, vol); noise(0.05, 500, vol * 0.3, 'lowpass'); break;
      case 'dum': drum(120, 55, 0.35, vol); break;
      case 'tom': drum(200, 110, 0.22, vol * 0.8); break;
      case 'tek': noise(0.07, 3800, vol, 'bandpass', 2); break;
      case 'frame': drum(180, 130, 0.25, vol * 0.6); noise(0.05, 1800, vol * 0.3); break;
      case 'shaker': noise(0.09, 7000, vol, 'highpass'); break;
      default: break;
    }
  }

  // ---- sound effects ----
  noise(dur, freq, vol, type = 'bandpass') {
    const c = this.ctx, t = c.currentTime;
    const buf = c.createBuffer(1, Math.ceil(c.sampleRate * dur), c.sampleRate);
    const d = buf.getChannelData(0);
    for (let i = 0; i < d.length; i++) d[i] = (Math.random() * 2 - 1) * (1 - i / d.length);
    const s = c.createBufferSource(); s.buffer = buf;
    const f = c.createBiquadFilter(); f.type = type; f.frequency.value = freq;
    const g = c.createGain(); g.gain.value = vol;
    s.connect(f); f.connect(g); g.connect(this.sfxBus);
    s.start(t);
  }

  tone(freq, dur, vol, type = 'sine', slideTo = 0, delay = 0) {
    const c = this.ctx, t = c.currentTime + delay;
    const o = c.createOscillator(); o.type = type; o.frequency.setValueAtTime(freq, t);
    if (slideTo) o.frequency.exponentialRampToValueAtTime(slideTo, t + dur);
    const g = c.createGain();
    g.gain.setValueAtTime(vol, t);
    g.gain.exponentialRampToValueAtTime(0.001, t + dur);
    o.connect(g); g.connect(this.sfxBus);
    o.start(t); o.stop(t + dur + 0.05);
  }

  sfx(name) {
    if (!this.ctx || this.muted || this.ctx.state !== 'running') return;
    switch (name) {
      case 'hit': this.noise(0.12, 900, 0.5); this.tone(150, 0.1, 0.3, 'square', 60); break;
      case 'hurt': this.noise(0.18, 500, 0.6); this.tone(110, 0.18, 0.35, 'sawtooth', 50); break;
      case 'melee': this.noise(0.15, 2500, 0.35, 'highpass'); break;
      case 'arrow': this.tone(300, 0.25, 0.25, 'triangle', 900); this.noise(0.1, 3000, 0.2); break;
      case 'bolt': this.tone(500, 0.3, 0.2, 'sine', 1200); break;
      case 'fire': this.noise(0.45, 600, 0.45, 'lowpass'); this.tone(200, 0.4, 0.2, 'sawtooth', 500); break;
      case 'frost': this.tone(1400, 0.4, 0.16, 'sine', 700); this.tone(2100, 0.3, 0.1, 'sine', 900); break;
      case 'nova': this.tone(180, 0.5, 0.3, 'sine', 700); this.noise(0.3, 1200, 0.3); break;
      case 'buff': [523, 659, 784].forEach((f, i) => this.tone(f, 0.35, 0.14, 'triangle', 0, i * 0.07)); break;
      case 'loot': this.tone(1318, 0.12, 0.2, 'triangle'); this.tone(1760, 0.2, 0.2, 'triangle', 0, 0.08); break;
      case 'quest': [392, 523, 659, 784].forEach((f, i) => this.tone(f, 0.4, 0.18, 'triangle', 0, i * 0.1)); break;
      case 'levelup': [523, 659, 784, 1047, 1319].forEach((f, i) => this.tone(f, 0.7, 0.2, 'triangle', 0, i * 0.12)); break;
      case 'learn': [659, 880, 1175].forEach((f, i) => this.tone(f, 0.45, 0.2, 'sine', 0, i * 0.09)); break;
      case 'coin': this.tone(1976, 0.08, 0.15, 'square'); this.tone(2637, 0.18, 0.15, 'square', 0, 0.06); break;
      case 'click': this.tone(800, 0.05, 0.1, 'square'); break;
    }
  }
}

export const audio = new Audio();
export { THEMES };
