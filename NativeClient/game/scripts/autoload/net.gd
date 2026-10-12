extends Node
## WebSocket connection to the game server. Messages are JSON objects with a type field "t".
## Use Net.on("type", callable) to receive them and Net.send("type", {...}) to send.

signal connected
signal closed

var peer: WebSocketPeer = null
var is_open := false
var _handlers: Dictionary = {}
var url := ""

func on(type: String, fn: Callable) -> void:
	if not _handlers.has(type):
		_handlers[type] = []
	_handlers[type].append(fn)

func off_all(owner: Object) -> void:
	for t in _handlers:
		_handlers[t] = _handlers[t].filter(func(c: Callable): return c.get_object() != owner)

func connect_to(u: String) -> int:
	close()
	url = u
	peer = WebSocketPeer.new()
	peer.inbound_buffer_size = 1 << 25 # the world message holds the whole map
	peer.outbound_buffer_size = 1 << 20
	peer.max_queued_packets = 8192
	return peer.connect_to_url(u)

func close() -> void:
	if peer:
		peer.close()
	peer = null
	is_open = false

func send(t: String, data: Dictionary = {}) -> void:
	if peer == null or peer.get_ready_state() != WebSocketPeer.STATE_OPEN:
		return
	var m := _intify(data.duplicate(true)) as Dictionary
	m["t"] = t
	peer.send_text(JSON.stringify(m))

## Godot parses every JSON number as float and would write 100.0 back; the Go server wants 100 for integer fields.
static func _intify(v: Variant) -> Variant:
	if v is float:
		if absf(v) < 9.0e15 and v == floorf(v):
			return int(v)
		return v
	if v is Dictionary:
		for k in v:
			v[k] = _intify(v[k])
		return v
	if v is Array:
		for i in v.size():
			v[i] = _intify(v[i])
		return v
	return v

func _process(_dt: float) -> void:
	if peer == null:
		return
	peer.poll()
	match peer.get_ready_state():
		WebSocketPeer.STATE_OPEN:
			if not is_open:
				is_open = true
				connected.emit()
			while peer.get_available_packet_count() > 0:
				var txt := peer.get_packet().get_string_from_utf8()
				var m = JSON.parse_string(txt)
				if m is Dictionary and m.has("t"):
					_dispatch(m["t"], m)
		WebSocketPeer.STATE_CLOSED:
			var was_open := is_open
			peer = null
			is_open = false
			if was_open:
				closed.emit()
			else:
				closed.emit()

## Go sends null for empty slices and maps: drop those keys so that `m.get(key, default)` gives the default.
static func strip_nulls(d: Dictionary) -> void:
	for k in d.keys():
		var v = d[k]
		if v == null:
			d.erase(k)
		elif v is Dictionary:
			strip_nulls(v)
		elif v is Array:
			for it in v:
				if it is Dictionary:
					strip_nulls(it)

func _dispatch(t: String, m: Dictionary) -> void:
	if t != "s":
		strip_nulls(m)
	if _handlers.has(t):
		for fn in _handlers[t].duplicate():
			if fn.is_valid():
				fn.call(m)
