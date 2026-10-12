class_name Tooltip
extends PanelContainer
## Floating tooltip (BBCode) that sticks to the cursor and stays inside the window.

static var inst: Tooltip = null
var label: RichTextLabel

func _init() -> void:
	inst = self
	mouse_filter = Control.MOUSE_FILTER_IGNORE
	z_index = 100
	visible = false
	add_theme_stylebox_override("panel", UI.box(Color(0.04, 0.03, 0.02, 0.97), UI.GOLD_DIM, 1, 5, 9))
	label = RichTextLabel.new()
	label.bbcode_enabled = true
	label.fit_content = true
	label.scroll_active = false
	label.custom_minimum_size.x = 250
	label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	label.add_theme_font_size_override("normal_font_size", 13)
	label.add_theme_font_size_override("bold_font_size", 15)
	add_child(label)

static func show_at(bbcode: String, pos: Vector2, above := false) -> void:
	if inst == null:
		return
	var t := inst
	t.label.text = bbcode
	t.label.custom_minimum_size.x = 250
	t.size = Vector2.ZERO
	t.visible = true
	t.reset_size.call_deferred()
	t._place.call_deferred(pos, above)

func _place(pos: Vector2, above: bool) -> void:
	reset_size()
	var vs := get_viewport_rect().size
	var p := pos
	if above:
		p.y -= size.y + 8.0
	p.x = clampf(p.x, 6.0, vs.x - size.x - 6.0)
	p.y = clampf(p.y, 6.0, vs.y - size.y - 6.0)
	position = p

static func hide_tip() -> void:
	if inst:
		inst.visible = false

static func follow_mouse(bbcode: String, from: Control) -> void:
	var m := from.get_viewport().get_mouse_position()
	var vs := from.get_viewport_rect().size
	var p := m + Vector2(16, 14)
	if p.x > vs.x - 280.0:
		p.x = m.x - 280.0
	show_at(bbcode, p)
