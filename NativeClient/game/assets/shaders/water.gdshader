shader_type spatial;
render_mode blend_mix, depth_draw_never, cull_back, specular_schlick_ggx;

uniform vec3 shallow_col = vec3(0.18, 0.55, 0.62);
uniform vec3 deep_col = vec3(0.05, 0.2, 0.38);
uniform sampler2D depth_tex : hint_depth_texture, repeat_disable, filter_nearest;
uniform float quality = 2.0;
uniform float wave = 0.06;

float h21(vec2 p) {
	p = fract(p * vec2(123.34, 456.21));
	p += dot(p, p + 45.32);
	return fract(p.x * p.y);
}
float vn(vec2 x) {
	vec2 i = floor(x);
	vec2 f = fract(x);
	f = f * f * (3.0 - 2.0 * f);
	return mix(mix(h21(i), h21(i + vec2(1.0, 0.0)), f.x), mix(h21(i + vec2(0.0, 1.0)), h21(i + vec2(1.0, 1.0)), f.x), f.y);
}
varying vec3 wpos;

void vertex() {
	wpos = (MODEL_MATRIX * vec4(VERTEX, 1.0)).xyz;
	if (quality > 0.5) VERTEX.y += (vn(wpos.xz * 0.35 + TIME * 0.4) - 0.5) * wave;
}

void fragment() {
	float d = 1.0;
	if (quality > 0.5) {
		float depth = texture(depth_tex, SCREEN_UV).r;
		vec3 ndc = vec3(SCREEN_UV * 2.0 - 1.0, depth);
		vec4 view = INV_PROJECTION_MATRIX * vec4(ndc, 1.0);
		view.xyz /= view.w;
		d = max(0.0, -view.z + VERTEX.z);
	}
	vec3 col = mix(shallow_col, deep_col, clamp(d * 0.28, 0.0, 1.0));
	float foam = quality > 1.5 ? (1.0 - smoothstep(0.0, 0.45, d)) * (0.55 + 0.45 * vn(wpos.xz * 1.8 + TIME * 0.6)) : 0.0;
	ALBEDO = mix(col, vec3(0.95), foam * 0.8);
	ALPHA = clamp(0.5 + d * 0.35, 0.5, 0.93);
	ROUGHNESS = 0.08;
	SPECULAR = 0.7;
	// blocky ripples: normal tilt per 1.5 m block
	vec2 bq = floor(wpos.xz / 1.5);
	float t1 = sin(TIME * 1.1 + h21(bq) * 6.28) * 0.07;
	float t2 = cos(TIME * 0.9 + h21(bq + 9.1) * 6.28) * 0.07;
	vec3 wn = normalize(vec3(t1, 1.0, t2));
	NORMAL = normalize((VIEW_MATRIX * vec4(wn, 0.0)).xyz);
}
