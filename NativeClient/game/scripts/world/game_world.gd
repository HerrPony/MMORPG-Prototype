class_name GameWorld
extends Node3D
## The open world: terrain, water, sky, day/night, props and buildings. Everything is built from the server's scenery.

# per biome: fog colour, sky colours, fog range, light strength and the sky box parameters
const LOOK := {
	"forest": {"fog": 0xbcd6e6, "top": 0x4a86d6, "bot": 0xbcd6e6, "near": 45, "far": 170, "sun": 1.5, "hemi": 0.85, "cloud": 0.48, "ccol": 0xffffff, "speed": 0.004, "stars": 1.0, "haze": 0.15, "tw": 0xff9a62, "aurora": 0.0, "glow": 0.0},
	"plains": {"fog": 0xcfe3f0, "top": 0x5a9ae0, "bot": 0xd8e8f2, "near": 70, "far": 230, "sun": 1.7, "hemi": 0.95, "cloud": 0.58, "ccol": 0xffffff, "speed": 0.006, "stars": 1.0, "haze": 0.1, "tw": 0xff8a50, "aurora": 0.0, "glow": 0.0},
	"swamp": {"fog": 0x8a9a7a, "top": 0x6a7a68, "bot": 0x9aa890, "near": 22, "far": 110, "sun": 0.8, "hemi": 0.6, "cloud": 0.85, "ccol": 0x7d8a74, "speed": 0.002, "stars": 0.25, "haze": 0.55, "tw": 0x9aa05a, "aurora": 0.0, "glow": 0.0},
	"snow": {"fog": 0xdce6f0, "top": 0x8aa8cc, "bot": 0xe6eef6, "near": 35, "far": 150, "sun": 1.3, "hemi": 1.0, "cloud": 0.38, "ccol": 0xe8f0ff, "speed": 0.003, "stars": 1.4, "haze": 0.2, "tw": 0xf08aa8, "aurora": 1.0, "glow": 0.0},
	"desert": {"fog": 0xf0dcb0, "top": 0x6aa0d8, "bot": 0xf2e4c0, "near": 55, "far": 210, "sun": 1.9, "hemi": 0.95, "cloud": 0.12, "ccol": 0xfff0d0, "speed": 0.003, "stars": 1.6, "haze": 0.4, "tw": 0xff6a30, "aurora": 0.0, "glow": 0.0},
	"jungle": {"fog": 0x8fd8a8, "top": 0x1f9a8a, "bot": 0xa8e8c0, "near": 22, "far": 125, "sun": 1.2, "hemi": 0.8, "cloud": 0.72, "ccol": 0xd8f4e0, "speed": 0.003, "stars": 0.6, "haze": 0.45, "tw": 0xffb040, "aurora": 0.0, "glow": 0.0},
	"blossom": {"fog": 0xf4d6e4, "top": 0xb98ad0, "bot": 0xffdbe6, "near": 50, "far": 200, "sun": 1.6, "hemi": 1.05, "cloud": 0.5, "ccol": 0xffe8f2, "speed": 0.005, "stars": 1.1, "haze": 0.25, "tw": 0xff7aa8, "aurora": 0.0, "glow": 0.0},
	"mesa": {"fog": 0xe8b88c, "top": 0xd8804e, "bot": 0xf6d2a2, "near": 45, "far": 200, "sun": 2.0, "hemi": 0.95, "cloud": 0.26, "ccol": 0xffd8b0, "speed": 0.004, "stars": 1.3, "haze": 0.35, "tw": 0xff5a20, "aurora": 0.0, "glow": 0.0},
	"volcano": {"fog": 0x4c3034, "top": 0x2a161e, "bot": 0x92402e, "near": 18, "far": 115, "sun": 0.9, "hemi": 0.55, "cloud": 0.92, "ccol": 0x4a3034, "speed": 0.006, "stars": 0.12, "haze": 0.6, "tw": 0xff4a20, "aurora": 0.0, "glow": 1.0},
	"sea": {"fog": 0x9ab8d0, "top": 0x5a8ac0, "bot": 0xb0c8dc, "near": 60, "far": 220, "sun": 1.5, "hemi": 0.9, "cloud": 0.5, "ccol": 0xffffff, "speed": 0.005, "stars": 1.1, "haze": 0.25, "tw": 0xff8a60, "aurora": 0.0, "glow": 0.0},
}
const NUM_KEYS := ["near", "far", "sun", "hemi", "cloud", "speed", "stars", "haze", "aurora", "glow"]
const COL_KEYS := ["fog", "top", "bot", "ccol", "tw"]
const BIOME_COLORS := [
	Vector3(0.76, 0.7, 0.5), Vector3(0.2, 0.42, 0.13), Vector3(0.46, 0.58, 0.2), Vector3(0.27, 0.34, 0.17), Vector3(0.9, 0.94, 1.0), Vector3(0.82, 0.69, 0.4),
	Vector3(0.07, 0.5, 0.13), Vector3(0.74, 0.5, 0.58), Vector3(0.78, 0.34, 0.18), Vector3(0.17, 0.14, 0.15), Vector3(1.0, 0.36, 0.04),
]

var env: Environment
var world_env: WorldEnvironment
var sun: DirectionalLight3D
var sky_mat: ShaderMaterial
var terrain: MeshInstance3D
var water: MeshInstance3D
var props_root: Node3D
var buildings_root: Node3D
var look_cur := {}
var look_tgt := {}
var zone_key := "forest"
var clock := 0.0
var light_dir := Vector3(0.5, 0.8, 0.3)
var built := false
var active := true
var _last_prop_cfg := ""
var tex_hm: ImageTexture
var tex_bm: ImageTexture
var tex_map: ImageTexture      # 1 m/pixel map overlay: r road, g paving, b buildings, a walls

func _ready() -> void:
	Settings.changed.connect(apply_settings)

func build() -> void:
	WD.load_scenery(S.scenery)
	_build_environment()
	_build_terrain()
	_build_water()
	props_root = Node3D.new(); props_root.name = "Props"; add_child(props_root)
	buildings_root = Node3D.new(); buildings_root.name = "Buildings"; add_child(buildings_root)
	Props.build(props_root, float(Settings.v["prop_density"]), float(Settings.v["view_distance"]), Settings.v["shadows"] > 0)
	_last_prop_cfg = "%s|%s|%s" % [Settings.v["prop_density"], Settings.v["view_distance"], Settings.v["shadows"] > 0]
	Buildings.build(buildings_root)
	Buildings.build_cabin(buildings_root)
	Buildings.build_bridges(buildings_root)
	look_cur = _look_dict(LOOK["forest"])
	look_tgt = look_cur.duplicate()
	built = true
	apply_settings()

func _look_dict(l: Dictionary) -> Dictionary:
	var d := {}
	for k in NUM_KEYS:
		d[k] = float(l[k])
	for k in COL_KEYS:
		d[k] = Models.C(l[k])
	return d

# ---------------------------------------------------------------- environment

func _build_environment() -> void:
	world_env = WorldEnvironment.new()
	env = Environment.new()
	sky_mat = ShaderMaterial.new()
	sky_mat.shader = load("res://assets/shaders/sky.gdshader")
	var sky := Sky.new()
	sky.sky_material = sky_mat
	sky.radiance_size = Sky.RADIANCE_SIZE_32
	sky.process_mode = Sky.PROCESS_MODE_QUALITY
	env.background_mode = Environment.BG_SKY
	env.sky = sky
	env.ambient_light_source = Environment.AMBIENT_SOURCE_COLOR
	env.reflected_light_source = Environment.REFLECTION_SOURCE_DISABLED
	env.fog_enabled = true
	env.fog_mode = Environment.FOG_MODE_DEPTH
	env.fog_density = 1.0
	env.fog_depth_curve = 1.4
	env.fog_sky_affect = 0.0
	env.glow_enabled = true
	env.glow_hdr_threshold = 1.5
	env.glow_hdr_scale = 1.5
	env.glow_intensity = 0.7
	env.glow_bloom = 0.02
	world_env.environment = env
	add_child(world_env)
	sun = DirectionalLight3D.new()
	sun.name = "Sun"
	sun.shadow_enabled = true
	sun.directional_shadow_mode = DirectionalLight3D.SHADOW_PARALLEL_4_SPLITS
	sun.directional_shadow_blend_splits = true
	sun.shadow_bias = 0.04
	sun.shadow_normal_bias = 1.2
	sun.light_angular_distance = 0.4
	add_child(sun)

## Applies the graphics settings that live in the environment, lights and viewport.
func apply_settings() -> void:
	if not built:
		return
	var v := Settings.v
	env.tonemap_mode = {"filmic": Environment.TONE_MAPPER_FILMIC, "aces": Environment.TONE_MAPPER_ACES, "reinhard": Environment.TONE_MAPPER_REINHARDT, "linear": Environment.TONE_MAPPER_LINEAR}.get(v["tonemap"], Environment.TONE_MAPPER_FILMIC)
	env.adjustment_enabled = true
	env.adjustment_brightness = float(v["brightness"])
	env.adjustment_contrast = float(v["contrast"])
	env.adjustment_saturation = float(v["saturation"])
	env.glow_enabled = bool(v["bloom"])
	env.ssao_enabled = bool(v["ssao"])
	env.ssao_radius = 1.6
	env.ssao_intensity = 2.2
	env.ssil_enabled = bool(v["ssil"])
	env.volumetric_fog_enabled = bool(v["volumetric_fog"])
	env.volumetric_fog_density = 0.012
	env.volumetric_fog_gi_inject = 0.0
	env.volumetric_fog_anisotropy = 0.6
	env.fog_enabled = bool(v["fog"])
	sun.shadow_enabled = int(v["shadows"]) > 0
	sun.directional_shadow_max_distance = float(v["shadow_distance"])
	var atlas: int = [1024, 1024, 2048, 4096, 8192][clampi(int(v["shadows"]), 0, 4)]
	RenderingServer.directional_shadow_atlas_set_size(atlas, true)
	RenderingServer.directional_soft_shadow_filter_set_quality([RenderingServer.SHADOW_QUALITY_HARD, RenderingServer.SHADOW_QUALITY_SOFT_LOW, RenderingServer.SHADOW_QUALITY_SOFT_MEDIUM, RenderingServer.SHADOW_QUALITY_SOFT_HIGH, RenderingServer.SHADOW_QUALITY_SOFT_ULTRA][clampi(int(v["shadows"]), 0, 4)])
	var vp := get_viewport()
	vp.msaa_3d = [Viewport.MSAA_DISABLED, Viewport.MSAA_2X, Viewport.MSAA_4X, Viewport.MSAA_8X][clampi(int(v["msaa"]), 0, 3)]
	vp.use_taa = bool(v["taa"])
	vp.screen_space_aa = Viewport.SCREEN_SPACE_AA_FXAA if v["fxaa"] else Viewport.SCREEN_SPACE_AA_DISABLED
	vp.scaling_3d_mode = Viewport.SCALING_3D_MODE_FSR if v["fsr"] else Viewport.SCALING_3D_MODE_BILINEAR
	vp.scaling_3d_scale = clampf(float(v["res_scale"]), 0.5, 1.0)
	if water and water.material_override:
		(water.material_override as ShaderMaterial).set_shader_parameter("quality", float(v["water_quality"]))
	rebuild_props.call_deferred()
	if props_root:
		for c in props_root.get_children():
			var mi := c as GeometryInstance3D
			if mi:
				mi.visibility_range_end = float(v["view_distance"]) if mi.visibility_range_end > 80.0 else mi.visibility_range_end

## Hides the open world while the player is inside a dungeon instance.
func set_active(on: bool) -> void:
	active = on
	visible = on
	if sun:
		sun.visible = on

func rebuild_props() -> void:
	var cfg := "%s|%s|%s" % [Settings.v["prop_density"], Settings.v["view_distance"], Settings.v["shadows"] > 0]
	if cfg == _last_prop_cfg and props_root and props_root.get_child_count() > 0:
		return
	_last_prop_cfg = cfg
	if props_root:
		Props.build(props_root, float(Settings.v["prop_density"]), float(Settings.v["view_distance"]), Settings.v["shadows"] > 0)

# ---------------------------------------------------------------- terrain and water

func _build_terrain() -> void:
	var sh := Shader.new()
	sh.code = FileAccess.get_file_as_string("res://assets/shaders/terrain.gdshader")
	var mat := ShaderMaterial.new()
	mat.shader = sh
	var himg := Image.create_from_data(WD.nx, WD.nz, false, Image.FORMAT_RF, WD.hm.to_byte_array())
	tex_hm = ImageTexture.create_from_image(himg)
	mat.set_shader_parameter("hm_tex", tex_hm)
	var bimg := Image.create_from_data(WD.nx, WD.nz, false, Image.FORMAT_R8, WD.bm)
	tex_bm = ImageTexture.create_from_image(bimg)
	mat.set_shader_parameter("bm_tex", tex_bm)
	var mask := _bake_mask()
	mat.set_shader_parameter("mask_tex", ImageTexture.create_from_image(mask))
	tex_map = ImageTexture.create_from_image(_bake_map_overlay(mask))
	mat.set_shader_parameter("grid", Vector2(WD.nx, WD.nz))
	mat.set_shader_parameter("origin", Vector2(WD.x0, WD.z0))
	mat.set_shader_parameter("cell", WD.cell)
	var pal := PackedVector3Array()
	for c in BIOME_COLORS:
		pal.append(c)
	mat.set_shader_parameter("palette", pal)
	var plane := PlaneMesh.new()
	plane.size = Vector2((WD.nx - 1) * WD.cell, (WD.nz - 1) * WD.cell)
	plane.subdivide_width = WD.nx - 2
	plane.subdivide_depth = WD.nz - 2
	plane.custom_aabb = AABB(Vector3(-plane.size.x / 2.0, -10.0, -plane.size.y / 2.0), Vector3(plane.size.x, 80.0, plane.size.y))
	terrain = MeshInstance3D.new()
	terrain.name = "Terrain"
	terrain.mesh = plane
	terrain.material_override = mat
	terrain.position = Vector3(WD.x0 + plane.size.x / 2.0, 0, WD.z0 + plane.size.y / 2.0)
	terrain.extra_cull_margin = 60.0
	terrain.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_ON
	add_child(terrain)

## Roads (dirt) and paving of the settlements as a 1 m/pixel mask: r = dirt, g = paving, b = paving type.
func _bake_mask() -> Image:
	var w := int((WD.nx - 1) * WD.cell)
	var h := int((WD.nz - 1) * WD.cell)
	var data := PackedByteArray()
	data.resize(w * h * 4)
	var put := func(px: int, pz: int, r: float, g: float, b: float) -> void:
		if px < 0 or pz < 0 or px >= w or pz >= h:
			return
		var o := (pz * w + px) * 4
		if r > 0.0:
			data[o] = maxi(data[o], int(r * 255.0))
		if g > 0.0 and int(g * 255.0) >= data[o + 1]:
			data[o + 1] = int(g * 255.0)
			data[o + 2] = int(b * 255.0)
	for rd in WD.roads:
		for i in range(rd.size() - 1):
			var a := Vector2(rd[i][0], rd[i][1])
			var b := Vector2(rd[i + 1][0], rd[i + 1][1])
			var steps := int(a.distance_to(b) / 0.6) + 1
			for s in steps + 1:
				var p := a.lerp(b, float(s) / steps)
				var cx := int(p.x - WD.x0)
				var cz := int(p.y - WD.z0)
				for dz in range(-3, 4):
					for dx in range(-3, 4):
						var d := Vector2(cx + dx + 0.5 - (p.x - WD.x0), cz + dz + 0.5 - (p.y - WD.z0)).length()
						var dirt := 1.0 - DayNight.smooth(1.1, 2.6, d)
						if dirt > 0.0:
							put.call(cx + dx, cz + dz, dirt, 0.0, 0.0)
	for t in WD.towns:
		var tx: float = t.x
		var tz: float = t.z
		match t.kind:
			"capital":
				var jade: bool = t.get("style", "") == "jade"
				for z in range(int(tz - WD.z0 - 47), int(tz - WD.z0 + 47)):
					for x in range(int(tx - WD.x0 - 52), int(tx - WD.x0 + 52)):
						put.call(x, z, 0.0, 0.7 if jade else 0.75, 1.0 if jade else 0.6)
			"village":
				for z in range(int(tz - WD.z0 - 15), int(tz - WD.z0 + 15)):
					for x in range(int(tx - WD.x0 - 15), int(tx - WD.x0 + 15)):
						var dd := Vector2(x + 0.5 - (tx - WD.x0), z + 0.5 - (tz - WD.z0)).length()
						if dd < 15.0:
							put.call(x, z, 0.0, (1.0 - DayNight.smooth(9.0, 15.0, dd)) * 0.8, 0.0)
			"site":
				for z in range(int(tz - WD.z0 - 12), int(tz - WD.z0 + 12)):
					for x in range(int(tx - WD.x0 - 12), int(tx - WD.x0 + 12)):
						if Vector2(x + 0.5 - (tx - WD.x0), z + 0.5 - (tz - WD.z0)).length() < 12.0:
							put.call(x, z, 0.0, 0.8, 0.3)
	return Image.create_from_data(w, h, false, Image.FORMAT_RGBA8, data)

## Copy of the road/paving mask with buildings (blue) and walls (alpha) stamped in, for the minimap and world map.
func _bake_map_overlay(mask: Image) -> Image:
	var img := mask.duplicate() as Image
	var w := img.get_width()
	var h := img.get_height()
	var data := img.get_data()
	for i in range(0, data.size(), 4):
		data[i + 2] = 0
		data[i + 3] = 0
	for b in WD.buildings:
		var k: String = b.k
		var chan := -1
		var val := 255
		if k == "house" or k == "hut" or k == "keep" or k == "pagoda":
			chan = 2
		elif k == "dock":
			chan = 2
			val = 120
		elif k == "wall":
			chan = 3
		if chan < 0:
			continue
		var x0 := maxi(0, int(b.x - b.w / 2.0 - WD.x0))
		var x1 := mini(w - 1, int(b.x + b.w / 2.0 - WD.x0))
		var z0 := maxi(0, int(b.z - b.d / 2.0 - WD.z0))
		var z1 := mini(h - 1, int(b.z + b.d / 2.0 - WD.z0))
		for z in range(z0, z1 + 1):
			for x in range(x0, x1 + 1):
				data[(z * w + x) * 4 + chan] = val
	return Image.create_from_data(w, h, false, Image.FORMAT_RGBA8, data)

func _build_water() -> void:
	var plane := PlaneMesh.new()
	plane.size = Vector2(3400, 2400)
	plane.subdivide_width = 160
	plane.subdivide_depth = 120
	water = MeshInstance3D.new()
	water.name = "Water"
	water.mesh = plane
	var mat := ShaderMaterial.new()
	mat.shader = load("res://assets/shaders/water.gdshader")
	water.material_override = mat
	water.position = Vector3(420, -0.12, 0)
	water.cast_shadow = GeometryInstance3D.SHADOW_CASTING_SETTING_OFF
	add_child(water)

# ---------------------------------------------------------------- per frame

func set_zone(key: String) -> void:
	zone_key = key if LOOK.has(key) else "forest"
	look_tgt = _look_dict(LOOK[zone_key])

func update(dt: float) -> void:
	if not built or not active:
		return
	clock += dt
	var k := 1.0 - exp(-dt * 1.2)
	for n in NUM_KEYS:
		look_cur[n] = lerpf(look_cur[n], look_tgt[n], k)
	for c in COL_KEYS:
		look_cur[c] = (look_cur[c] as Color).lerp(look_tgt[c], k)
	var t := DayNight.info()
	var day: float = t.day
	var night: float = t.night
	var tw: float = t.twilight
	var night_tint := Color(0.1, 0.14, 0.3)
	var g_top: Color = (look_cur["top"] * night_tint).lerp(Color(0.016, 0.04, 0.11), 0.55).lerp(look_cur["top"], day)
	var g_bot: Color = (look_cur["bot"] * night_tint).lerp(Color(0.06, 0.11, 0.24), 0.6).lerp(look_cur["bot"], day)
	g_bot = g_bot.lerp(look_cur["tw"], tw * 0.65)
	g_top = g_top.lerp(look_cur["tw"] * 0.55, tw * 0.22)
	var g_fog: Color = g_bot * (0.55 + 0.45 * day)
	var vd := float(Settings.v["view_distance"])
	var span: float = look_cur["far"] * (0.75 + 0.25 * day) * (vd / 170.0)
	env.fog_light_color = g_fog
	env.fog_depth_begin = look_cur["near"] * (vd / 170.0) * 0.8
	env.fog_depth_end = span
	var cam := get_viewport().get_camera_3d()
	if cam:
		cam.far = maxf(vd * 1.6, 300.0)
	var u := sky_mat
	u.set_shader_parameter("top_col", Vector3(g_top.r, g_top.g, g_top.b))
	u.set_shader_parameter("bot_col", Vector3(g_bot.r, g_bot.g, g_bot.b))
	u.set_shader_parameter("sun_dir", t.sun)
	u.set_shader_parameter("moon_dir", t.moon)
	var sc_low := Color(1.0, 0.6, 0.29)
	var sc_day := Color(1.0, 0.94, 0.85)
	var sun_c := sc_low.lerp(sc_day, clampf(float(t.elev) * 2.2, 0.0, 1.0))
	u.set_shader_parameter("sun_col", Vector3(sun_c.r, sun_c.g, sun_c.b))
	var cc: Color = look_cur["ccol"]
	u.set_shader_parameter("cloud_col", Vector3(cc.r, cc.g, cc.b))
	var gc: Color = look_cur["tw"]
	u.set_shader_parameter("glow_col", Vector3(gc.r, gc.g, gc.b))
	u.set_shader_parameter("clock", clock)
	u.set_shader_parameter("night", night)
	u.set_shader_parameter("sun_up", clampf((float(t.elev) + 0.1) * 4.0, 0.0, 1.0))
	u.set_shader_parameter("cloud", look_cur["cloud"])
	u.set_shader_parameter("cspeed", look_cur["speed"])
	u.set_shader_parameter("stars", look_cur["stars"])
	u.set_shader_parameter("aurora", look_cur["aurora"])
	u.set_shader_parameter("haze", look_cur["haze"])
	u.set_shader_parameter("glow", look_cur["glow"])
	# light: sun by day, moon by night, never lower than a minimum angle so shadows stay inside the shadow map
	var src: Vector3 = t.sun if float(t.elev) > -0.04 else t.moon
	light_dir = Vector3(src.x, maxf(0.3, src.y), src.z).normalized()
	var moon_c := Color(0.54, 0.64, 1.0)
	var lc := moon_c.lerp(sc_low, minf(1.0, day * 1.6)).lerp(Color(1.0, 0.94, 0.85), clampf((float(t.elev) - 0.1) * 2.0, 0.0, 1.0))
	sun.light_color = lc
	sun.light_energy = look_cur["sun"] * day * (0.35 + 0.65 * minf(1.0, float(t.elev) * 3.0 + 0.2)) * 0.9 + night * 0.5 * 0.75
	var amb_day := Color(0.8, 0.9, 1.0)
	var amb_night := Color(0.3, 0.38, 0.64)
	env.ambient_light_color = amb_night.lerp(amb_day, day)
	env.ambient_light_energy = maxf(0.45, look_cur["hemi"] * (0.5 + 0.5 * day) * 0.8)
	# windows and lamps glow at night
	Buildings.glow_material().set_shader_parameter("night", night)

## Called every frame with the player's position: the light and its shadows follow.
func follow(p: Vector3) -> void:
	if sun:
		sun.global_transform = Transform3D(Basis.looking_at(-light_dir, Vector3.UP), p + light_dir * 110.0)
