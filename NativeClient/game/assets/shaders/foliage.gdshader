shader_type spatial;
render_mode cull_back, specular_disabled;

// Instanced props (trees, rocks, bushes ...). Vertex colour alpha 0.5 marks the parts that take the per-instance tint
// (leaves, grass, rock); wind sways everything above the ground.
uniform float wind = 1.0;
uniform float fade_start = 150.0;
varying vec3 vcol;
varying float vtint;
varying vec3 itint;

void vertex() {
	vec3 wpos = (MODEL_MATRIX * vec4(VERTEX, 1.0)).xyz;
	float sway = sin(TIME * 1.5 + wpos.x * 0.33 + wpos.z * 0.27) * 0.03 * wind * max(VERTEX.y, 0.0);
	VERTEX.x += sway;
	VERTEX.z += sway * 0.6;
	vcol = COLOR.rgb;
	vtint = COLOR.a < 0.75 ? 1.0 : 0.0;
	itint = INSTANCE_CUSTOM.rgb;
}

void fragment() {
	vec3 c = vcol;
	if (vtint > 0.5) c = vcol * itint;
	ALBEDO = c;
	ROUGHNESS = 0.95;
}
