// Procedural music and sound effects (WebAudio, no asset files needed).
// Music: a slow, folky A-minor pad + pentatonic melody; a drum pulse fades in during combat.

const NOTE = (n) => 440 * Math.pow(2, (n - 69) / 12); // MIDI -> Hz
// chord progression (MIDI roots): Am - F - C - G, each 8 seconds
const PROG = [[57, 60, 64], [53, 57, 60], [48, 52, 55], [55, 59, 62]];
const SCALE = [69, 72, 74, 76, 79, 81, 84]; // A minor pentatonic-ish (A C D E G A C)

export class Audio {
  constructor() {
    this.ctx = null;
    this.mode = 'menu';       // 'menu' | 'world'
    this.intensity = 0;       // 0..1 combat level
    this.muted = false;
    this.vol = 0.6;
    try { const v = localStorage.getItem('vol'); if (v !== null) this.vol = +v; this.muted = localStorage.getItem('muted') === '1'; } catch { /* ignore */ }
    this.nextChord = 0; this.chordIdx = 0;
    this.nextMel = 0; this.nextBeat = 0; this.beat = 0;
    this.lastNote = 3;
  }

  // Must be called from a user gesture (browsers block audio until then).
  start() {
    if (this.ctx) { if (this.ctx.state === 'suspended') this.ctx.resume(); return; }
    const AC = window.AudioContext || window.webkitAudioContext;
    if (!AC) return;
    const c = this.ctx = new AC();
    this.master = c.createGain();
    this.master.gain.value = this.muted ? 0 : this.vol;
    this.master.connect(c.destination);
    // simple "reverb": filtered feedback delay
    this.bus = c.createGain();
    const delay = c.createDelay(1.0); delay.delayTime.value = 0.37;
    const fb = c.createGain(); fb.gain.value = 0.38;
    const lp = c.createBiquadFilter(); lp.type = 'lowpass'; lp.frequency.value = 2200;
    this.bus.connect(this.master);
    this.bus.connect(delay); delay.connect(lp); lp.connect(fb); fb.connect(delay); lp.connect(this.master);
    this.sfxBus = c.createGain(); this.sfxBus.gain.value = 0.7; this.sfxBus.connect(this.master);
    this.nextChord = this.nextMel = this.nextBeat = c.currentTime + 0.2;
    setInterval(() => this.schedule(), 120);
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

  // ---- music scheduling ----
  schedule() {
    const c = this.ctx;
    if (!c || c.state !== 'running') return;
    const ahead = c.currentTime + 0.5;
    const world = this.mode === 'world';
    while (this.nextChord < ahead) {
      this.pad(PROG[this.chordIdx % PROG.length], this.nextChord, 8);
      this.chordIdx++;
      this.nextChord += 8;
    }
    const step = world ? 0.5 : 1.0;
    while (this.nextMel < ahead) {
      const dense = world ? 0.55 + this.intensity * 0.3 : 0.4;
      if (Math.random() < dense) {
        // random walk on the scale
        this.lastNote = Math.max(0, Math.min(SCALE.length - 1, this.lastNote + Math.floor(Math.random() * 5) - 2));
        const dur = Math.random() < 0.3 ? 1.6 : 0.8;
        this.pluck(NOTE(SCALE[this.lastNote]), this.nextMel, dur, world ? 0.16 : 0.1);
        if (Math.random() < 0.12) this.pluck(NOTE(SCALE[this.lastNote] - 12), this.nextMel, dur, 0.08);
      }
      this.nextMel += step * (this.intensity > 0.5 ? 0.5 : 1);
    }
    // combat drums
    while (this.nextBeat < ahead) {
      if (this.intensity > 0.3 && world) {
        const b = this.beat % 8;
        if (b % 4 === 0) this.drum(this.nextBeat, 70, 0.5 * this.intensity, 0.2);
        else if (b % 2 === 0) this.drum(this.nextBeat, 180, 0.22 * this.intensity, 0.08);
      }
      this.beat++;
      this.nextBeat += 0.25;
    }
  }

  pad(chord, t, dur) {
    const c = this.ctx;
    const world = this.mode === 'world';
    for (const [i, n] of [...chord, chord[0] - 12].entries()) {
      for (const det of [-6, 6]) {
        const o = c.createOscillator();
        o.type = i === 3 ? 'sine' : 'triangle';
        o.frequency.value = NOTE(n);
        o.detune.value = det;
        const g = c.createGain();
        const peak = (i === 3 ? 0.11 : 0.05) * (world ? 1 : 0.8);
        g.gain.setValueAtTime(0, t);
        g.gain.linearRampToValueAtTime(peak, t + 2.5);
        g.gain.setValueAtTime(peak, t + dur - 2.5);
        g.gain.linearRampToValueAtTime(0, t + dur + 1.5);
        const f = c.createBiquadFilter(); f.type = 'lowpass'; f.frequency.value = 900;
        o.connect(f); f.connect(g); g.connect(this.bus);
        o.start(t); o.stop(t + dur + 1.6);
      }
    }
  }

  pluck(freq, t, dur, vol) {
    const c = this.ctx;
    const o = c.createOscillator(); o.type = 'triangle'; o.frequency.value = freq;
    const o2 = c.createOscillator(); o2.type = 'sine'; o2.frequency.value = freq * 2;
    const g = c.createGain();
    g.gain.setValueAtTime(0, t);
    g.gain.linearRampToValueAtTime(vol, t + 0.015);
    g.gain.exponentialRampToValueAtTime(0.0008, t + dur + 0.6);
    const g2 = c.createGain(); g2.gain.value = 0.25;
    o.connect(g); o2.connect(g2); g2.connect(g); g.connect(this.bus);
    o.start(t); o2.start(t); o.stop(t + dur + 0.7); o2.stop(t + dur + 0.7);
  }

  drum(t, freq, vol, dur) {
    const c = this.ctx;
    const o = c.createOscillator(); o.type = 'sine';
    o.frequency.setValueAtTime(freq * 1.8, t);
    o.frequency.exponentialRampToValueAtTime(freq * 0.6, t + dur);
    const g = c.createGain();
    g.gain.setValueAtTime(vol * 0.5, t);
    g.gain.exponentialRampToValueAtTime(0.001, t + dur + 0.1);
    o.connect(g); g.connect(this.master);
    o.start(t); o.stop(t + dur + 0.15);
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
