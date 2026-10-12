class_name GameWindow
extends PanelContainer
## A draggable window with a title bar and a close button. Content goes into `body`.

signal closed

var title_label: Label
var body: VBoxContainer
var _drag := false
var _drag_off := Vector2.ZERO
var win_id := ""
var default_pos := Vector2(200, 100)
var _placed := false

func _init(title := "", min_size := Vector2(300, 100)) -> void:
	custom_minimum_size = min_size
	visible = false
	mouse_filter = Control.MOUSE_FILTER_STOP
	var v := VBoxContainer.new()
	v.add_theme_constant_override("separation", 6)
	add_child(v)
	var bar := PanelContainer.new()
	bar.mouse_filter = Control.MOUSE_FILTER_STOP
	bar.add_theme_stylebox_override("panel", UI.box(Color(0.84, 0.7, 0.38, 0.1), Color(0, 0, 0, 0), 0, 4, 4))
	bar.mouse_default_cursor_shape = Control.CURSOR_MOVE
	var h := HBoxContainer.new()
	bar.add_child(h)
	title_label = UI.label(title, 16, UI.GOLD)
	title_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	title_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	h.add_child(title_label)
	var x := Button.new()
	x.text = "✕"
	x.flat = true
	x.focus_mode = Control.FOCUS_NONE
	x.pressed.connect(close)
	h.add_child(x)
	v.add_child(bar)
	body = VBoxContainer.new()
	body.add_theme_constant_override("separation", 6)
	body.size_flags_vertical = Control.SIZE_EXPAND_FILL
	v.add_child(body)
	bar.gui_input.connect(_bar_input)
	gui_input.connect(func(ev: InputEvent):
		if ev is InputEventMouseButton and ev.pressed:
			raise())

func raise() -> void:
	if get_parent():
		get_parent().move_child(self, -1)

func _bar_input(ev: InputEvent) -> void:
	if ev is InputEventMouseButton and ev.button_index == MOUSE_BUTTON_LEFT:
		_drag = ev.pressed
		_drag_off = get_global_mouse_position() - global_position
		raise()
	elif ev is InputEventMouseMotion and _drag:
		var vs := get_viewport_rect().size
		var p := get_global_mouse_position() - _drag_off
		position = Vector2(clampf(p.x, 0.0, vs.x - 80.0), clampf(p.y, 0.0, vs.y - 40.0))

func open(at := Vector2(-1, -1)) -> void:
	if not _placed or at.x >= 0.0:
		position = at if at.x >= 0.0 else default_pos
		_placed = true
	if not visible:
		Sound.sfx("ui_open", -10.0)
	visible = true
	UI.nofocus(self)
	raise()
	reset_size()
	_clamp.call_deferred()

func _clamp() -> void:
	var vs := get_viewport_rect().size
	position = Vector2(clampf(position.x, 0.0, maxf(0.0, vs.x - size.x)), clampf(position.y, 0.0, maxf(0.0, vs.y - size.y)))

func close() -> void:
	if visible:
		visible = false
		Sound.sfx("ui_close", -10.0)
		Tooltip.hide_tip()
		closed.emit()

func toggle() -> void:
	if visible:
		close()
	else:
		open()
