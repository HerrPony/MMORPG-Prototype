// WebSocket connection with a tiny event emitter.
export class Net {
  constructor() {
    this.handlers = {};
    this.ws = null;
  }

  on(type, fn) {
    (this.handlers[type] ||= []).push(fn);
    return this;
  }

  connect() {
    return new Promise((resolve, reject) => {
      const host = location.protocol === 'file:' ? 'localhost:8080' : location.host;
      const proto = location.protocol === 'https:' ? 'wss' : 'ws';
      const ws = new WebSocket(`${proto}://${host}/ws`);
      this.ws = ws;
      ws.onopen = () => resolve();
      ws.onerror = () => reject(new Error('Verbindung fehlgeschlagen'));
      ws.onclose = () => (this.handlers['_close'] || []).forEach((f) => f());
      ws.onmessage = (ev) => {
        let m;
        try { m = JSON.parse(ev.data); } catch { return; }
        (this.handlers[m.t] || []).forEach((f) => f(m));
      };
    });
  }

  send(t, data = {}) {
    if (this.ws && this.ws.readyState === 1) this.ws.send(JSON.stringify({ t, ...data }));
  }
}
