class_name UI
extends RefCounted
## Theme and small factory helpers for the fantasy-styled interface (dark glass panels with gold borders).

const GOLD := Color("f0cf7a")
const GOLD_DIM := Color("a98a45")
const EDGE := Color(0.84, 0.7, 0.38, 0.55)
const EDGE_SOFT := Color(0.84, 0.7, 0.38, 0.25)
const TEXT := Color("ece3cf")
const MUTED := Color("a89a7e")
const GREEN := Color("3ccf4e")
const RED := Color("e24a3a")
const BG := Color(0.07, 0.065, 0.055, 0.94)
const BG_FLAT := Color(0.055, 0.052, 0.047, 0.8)
const CLASS_COLORS := {"warrior": Color("c79c6e"), "hunter": Color("abd473"), "mage": Color("69ccf0")}
const CLASS_ICONS := {"warrior": "⚔️", "hunter": "🏹", "mage": "🔮"}
const RARITY := {"poor": Color("9d9d9d"), "common": Color("ffffff"), "uncommon": Color("1eff00"), "rare": Color("4a9bff"), "epic": Color("c26bff")}

static var _theme: Theme = null
static var emoji_font: Font = null

static func box(bg: Color, border := Color(0, 0, 0, 0), bw := 0, radius := 6, pad := 0) -> StyleBoxFlat:
	var s := StyleBoxFlat.new()
	s.bg_color = bg
	s.border_color = border
	s.set_border_width_all(bw)
	s.set_corner_radius_all(radius)
	if pad > 0:
		s.set_content_margin_all(pad)
	return s

static func theme() -> Theme:
	if _theme != null:
		return _theme
	var t := Theme.new()
	# serif like the reference interface (Windows: Palatino/Georgia); falls back to the built-in font elsewhere
	var base := SystemFont.new()
	base.font_names = PackedStringArray(["Palatino Linotype", "Book Antiqua", "Georgia", "Cambria", "Times New Roman", "DejaVu Serif", "Liberation Serif", "serif"])
	base.allow_system_fallback = true
	var fv := FontVariation.new()
	fv.base_font = base
	if ResourceLoader.exists("res://assets/fonts/NotoColorEmoji.ttf"):
		emoji_font = load("res://assets/fonts/NotoColorEmoji.ttf")
		fv.fallbacks = [emoji_font]
	t.default_font = fv
	t.default_font_size = 14
	t.set_color("font_color", "Label", TEXT)
	t.set_color("font_shadow_color", "Label", Color(0, 0, 0, 0.7))
	t.set_constant("shadow_offset_x", "Label", 1)
	t.set_constant("shadow_offset_y", "Label", 1)
	t.set_color("default_color", "RichTextLabel", TEXT)
	t.set_font_size("normal_font_size", "RichTextLabel", 14)
	# panels
	var glass := box(BG, EDGE, 1, 6, 10)
	glass.shadow_color = Color(0, 0, 0, 0.6)
	glass.shadow_size = 12
	t.set_stylebox("panel", "PanelContainer", glass)
	t.set_stylebox("panel", "Panel", glass)
	var flat := box(BG_FLAT, EDGE_SOFT, 1, 4, 6)
	t.add_type("FlatPanel")
	t.set_type_variation("FlatPanel", "PanelContainer")
	t.set_stylebox("panel", "FlatPanel", flat)
	# buttons
	var bn := box(Color(0.05, 0.045, 0.04, 0.62), EDGE, 1, 4, 6)
	bn.content_margin_left = 12
	bn.content_margin_right = 12
	var bh := box(Color(0.2, 0.16, 0.09, 0.85), GOLD, 1, 4, 6)
	bh.content_margin_left = 12
	bh.content_margin_right = 12
	var bp := box(Color(0.42, 0.33, 0.14, 0.92), GOLD, 1, 4, 6)
	bp.content_margin_left = 12
	bp.content_margin_right = 12
	var bd := box(Color(0.05, 0.045, 0.04, 0.4), EDGE_SOFT, 1, 4, 6)
	t.set_stylebox("normal", "Button", bn)
	t.set_stylebox("hover", "Button", bh)
	t.set_stylebox("pressed", "Button", bp)
	t.set_stylebox("disabled", "Button", bd)
	t.set_stylebox("focus", "Button", StyleBoxEmpty.new())
	t.set_color("font_color", "Button", TEXT)
	t.set_color("font_hover_color", "Button", Color.WHITE)
	t.set_color("font_disabled_color", "Button", Color(1, 1, 1, 0.3))
	t.add_type("BigButton")
	t.set_type_variation("BigButton", "Button")
	var gn := box(Color("cf9a33"), Color("6a4a10"), 1, 4, 10)
	var gh := box(Color("e4b04a"), Color("6a4a10"), 1, 4, 10)
	var gp := box(Color("b88626"), Color("6a4a10"), 1, 4, 10)
	t.set_stylebox("normal", "BigButton", gn)
	t.set_stylebox("hover", "BigButton", gh)
	t.set_stylebox("pressed", "BigButton", gp)
	t.set_stylebox("disabled", "BigButton", box(Color(0.3, 0.27, 0.2), Color(0.1, 0.1, 0.1), 1, 4, 10))
	t.set_color("font_color", "BigButton", Color("2b1a05"))
	t.set_color("font_hover_color", "BigButton", Color("2b1a05"))
	t.set_color("font_pressed_color", "BigButton", Color("2b1a05"))
	t.set_color("font_disabled_color", "BigButton", Color(0.55, 0.5, 0.4))
	t.set_font_size("font_size", "BigButton", 16)
	# inputs
	var le := box(Color(0, 0, 0, 0.45), EDGE_SOFT, 1, 4, 8)
	t.set_stylebox("normal", "LineEdit", le)
	t.set_stylebox("focus", "LineEdit", box(Color(0, 0, 0, 0.55), GOLD, 1, 4, 8))
	t.set_color("font_color", "LineEdit", Color.WHITE)
	t.set_color("caret_color", "LineEdit", GOLD)
	t.set_color("font_placeholder_color", "LineEdit", Color(1, 1, 1, 0.3))
	# check boxes / sliders
	t.set_color("font_color", "CheckBox", TEXT)
	t.set_color("font_color", "CheckButton", TEXT)
	t.set_color("font_color", "OptionButton", TEXT)
	t.set_stylebox("normal", "OptionButton", bn)
	t.set_stylebox("hover", "OptionButton", bh)
	t.set_stylebox("pressed", "OptionButton", bp)
	var slider := box(Color(0, 0, 0, 0.5), EDGE_SOFT, 1, 3)
	slider.content_margin_top = 4
	slider.content_margin_bottom = 4
	t.set_stylebox("slider", "HSlider", slider)
	t.set_stylebox("grabber_area", "HSlider", box(GOLD_DIM, Color(0, 0, 0, 0), 0, 3))
	t.set_stylebox("grabber_area_highlight", "HSlider", box(GOLD, Color(0, 0, 0, 0), 0, 3))
	# scrollbars
	var sb := box(Color("6d5a30"), Color(0, 0, 0, 0), 0, 4)
	t.set_stylebox("grabber", "VScrollBar", sb)
	t.set_stylebox("grabber_highlight", "VScrollBar", box(Color("8a7240"), Color(0, 0, 0, 0), 0, 4))
	t.set_stylebox("grabber_pressed", "VScrollBar", box(GOLD_DIM, Color(0, 0, 0, 0), 0, 4))
	t.set_stylebox("scroll", "VScrollBar", box(Color(0, 0, 0, 0.2), Color(0, 0, 0, 0), 0, 4))
	# tabs
	t.set_stylebox("tab_selected", "TabContainer", box(Color(0.94, 0.81, 0.48, 0.14), GOLD, 1, 4, 6))
	t.set_stylebox("tab_unselected", "TabContainer", box(Color(1, 1, 1, 0.04), EDGE_SOFT, 1, 4, 6))
	t.set_stylebox("tab_hovered", "TabContainer", box(Color(1, 1, 1, 0.09), EDGE, 1, 4, 6))
	t.set_stylebox("panel", "TabContainer", box(Color(0, 0, 0, 0), Color(0, 0, 0, 0), 0, 0, 4))
	t.set_color("font_selected_color", "TabContainer", GOLD)
	t.set_color("font_unselected_color", "TabContainer", MUTED)
	t.set_color("font_hovered_color", "TabContainer", Color.WHITE)
	# tooltip
	t.set_stylebox("panel", "TooltipPanel", box(Color(0.04, 0.03, 0.02, 0.97), GOLD_DIM, 1, 5, 8))
	_theme = t
	return t

static func label(text: String, size := 14, color := TEXT, align := HORIZONTAL_ALIGNMENT_LEFT) -> Label:
	var l := Label.new()
	l.text = text
	l.horizontal_alignment = align
	l.add_theme_font_size_override("font_size", size)
	l.add_theme_color_override("font_color", color)
	return l

static func button(text: String, cb: Callable = Callable(), big := false) -> Button:
	var b := Button.new()
	b.text = text
	b.focus_mode = Control.FOCUS_NONE   # a focused button would be pressed again by the jump key
	if big:
		b.theme_type_variation = "BigButton"
		b.custom_minimum_size.y = 40
	if cb.is_valid():
		b.pressed.connect(cb)
	return b

static func spacer(h := 8) -> Control:
	var c := Control.new()
	c.custom_minimum_size = Vector2(0, h)
	return c

static func fill() -> Control:
	var c := Control.new()
	c.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	c.size_flags_vertical = Control.SIZE_EXPAND_FILL
	return c

static func hbox(sep := 8) -> HBoxContainer:
	var h := HBoxContainer.new()
	h.add_theme_constant_override("separation", sep)
	return h

static func vbox(sep := 6) -> VBoxContainer:
	var v := VBoxContainer.new()
	v.add_theme_constant_override("separation", sep)
	return v

static func rich(bbcode: String, size := 14) -> RichTextLabel:
	var r := RichTextLabel.new()
	r.bbcode_enabled = true
	r.fit_content = true
	r.scroll_active = false
	r.text = bbcode
	r.add_theme_font_size_override("normal_font_size", size)
	r.add_theme_font_size_override("bold_font_size", size)
	return r

## A clickable card (panel with child controls) used for selections; call .set_meta("sel", bool) via set_selected.
static func card(on_click: Callable, min_size := Vector2(0, 0)) -> PanelContainer:
	var p := PanelContainer.new()
	p.custom_minimum_size = min_size
	p.mouse_filter = Control.MOUSE_FILTER_STOP
	p.mouse_default_cursor_shape = Control.CURSOR_POINTING_HAND
	set_selected(p, false)
	p.gui_input.connect(func(ev: InputEvent):
		if ev is InputEventMouseButton and ev.pressed and ev.button_index == MOUSE_BUTTON_LEFT:
			on_click.call())
	return p

static func set_selected(p: PanelContainer, sel: bool) -> void:
	var s := box(Color(0.94, 0.81, 0.48, 0.14) if sel else Color(1, 1, 1, 0.04), GOLD if sel else EDGE_SOFT, 2 if sel else 1, 5, 8)
	p.add_theme_stylebox_override("panel", s)

## Buttons and sliders must not keep keyboard focus: space (jump) and the arrow keys would operate them.
static func nofocus(n: Node) -> void:
	if n is BaseButton or n is Slider or n is TabBar or n is TabContainer:
		(n as Control).focus_mode = Control.FOCUS_NONE
	for c in n.get_children():
		nofocus(c)

static func clear(n: Node) -> void:
	for c in n.get_children():
		n.remove_child(c)
		c.queue_free()

static func hex(c: int) -> Color:
	return Models.C(c)

static func copper_text(c: int) -> String:
	var g := c / 10000
	var s := (c / 100) % 100
	var k := c % 100
	var out := ""
	if g > 0:
		out += "[color=#ffd700]%dg[/color] " % g
	if g > 0 or s > 0:
		out += "[color=#ccccdd]%ds[/color] " % s
	out += "[color=#d9803d]%dk[/color]" % k
	return out

static func rarity_color(r: String) -> Color:
	return RARITY.get(r, Color.WHITE)
