shader_type spatial;
render_mode blend_mix, depth_draw_never, cull_back, specular_schlick_ggx;

// Realistic water: a sum of travelling waves (long swell in the vertices, fine ripples in the normal), depth absorption,
// fresnel, sky/sun reflection through the environment, subsurface glow on the crests and foam on the shore and on steep crests.
// quality: 0 = calm flat plane, 1 = waves + depth colour, 2 = additionally foam and crest highlights.

uniform vec3 shallow_col = vec3(0.14, 0.55, 0.56);
uniform vec3 deep_col = vec3(0.02, 0.17, 0.27);
uniform vec3 scatter_col = vec3(0.10, 0.45, 0.42);
uniform vec3 sky_col = vec3(0.45, 0.62, 0.78);
uniform vec3 foam_col = vec3(0.93, 0.97, 1.0);
uniform sampler2D depth_tex : hint_depth_texture, repeat_disable, filter_nearest;
uniform float quality = 2.0;

// xy = direction, z = wavelength (m), w = slope (steepness of the wave)
const int NW = 8;
const vec4 WV[8] = vec4[8](
	vec4(1.0, 0.2, 80.0, 0.020),
	vec4(-0.6, 0.8, 47.0, 0.026),
	vec4(0.3, 1.0, 21.0, 0.036),
	vec4(-0.9, 0.4, 11.0, 0.050),
	vec4(0.7, -0.7, 5.7, 0.060),
	vec4(-0.2, -1.0, 3.1, 0.065),
	vec4(0.95, 0.3, 1.7, 0.055),
	vec4(-0.5, 0.85, 0.9, 0.040)
);

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

// returns (height, d/dx, d/dz) of waves [from, to)
vec3 waves(vec2 p, float t, int from, int to) {
	vec3 r = vec3(0.0);
	for (int i = from; i < to; i++) {
		vec4 w = WV[i];
		vec2 d = normalize(w.xy);
		float k = 6.2831853 / w.z;
		float ph = k * dot(d, p) - sqrt(9.81 * k) * t * 0.55;
		r.x += w.w / k * sin(ph);
		r.yz += w.w * cos(ph) * d;
	}
	return r;
}

varying vec3 wpos;
varying float swell;

void vertex() {
	wpos = (MODEL_MATRIX * vec4(VERTEX, 1.0)).xyz;
	swell = 0.0;
	if (quality > 0.5) {
		vec2 p = wpos.xz + (vec2(vn(wpos.xz * 0.02 + 3.1), vn(wpos.xz * 0.02 - 7.7)) - 0.5) * 10.0;
		float h = waves(p, TIME, 0, 3).x;
		VERTEX.y += h;
		swell = h;
	}
}

void fragment() {
	float dist = length(VERTEX);
	// depth of the water above the terrain behind it
	float d = 6.0;
	if (quality > 0.5) {
		float depth = texture(depth_tex, SCREEN_UV).r;
#if CURRENT_RENDERER == RENDERER_COMPATIBILITY
		vec3 ndc = vec3(SCREEN_UV * 2.0 - 1.0, depth * 2.0 - 1.0); // OpenGL depth runs 0..1, its clip space -1..1
#else
		vec3 ndc = vec3(SCREEN_UV * 2.0 - 1.0, depth);
#endif
		vec4 view = INV_PROJECTION_MATRIX * vec4(ndc, 1.0);
		view.xyz /= view.w;
		d = max(0.0, -view.z + VERTEX.z);
	}

	// surface normal from the travelling waves; far away the fine detail fades out (no shimmering)
	vec2 p = wpos.xz + (vec2(vn(wpos.xz * 0.02 + 3.1), vn(wpos.xz * 0.02 - 7.7)) - 0.5) * 10.0;
	vec3 wv = vec3(0.0);
	float detail = 1.0 / (1.0 + dist * dist / 5000.0);
	if (quality > 0.5) {
		wv = waves(p, TIME, 0, 4);
		vec3 fine = waves(p, TIME, 4, NW);
		wv.x += fine.x * detail;
		wv.yz += fine.yz * detail;
		// irregular micro ripples so the pattern never looks tiled
		vec2 q = wpos.xz * 1.9 + vec2(TIME * 0.35, -TIME * 0.22);
		vec2 g = vec2(vn(q + vec2(0.7, 0.0)) - vn(q - vec2(0.7, 0.0)), vn(q + vec2(0.0, 0.7)) - vn(q - vec2(0.0, 0.7)));
		wv.yz += g * 0.12 * detail;
	}
	vec3 wn = normalize(vec3(-wv.y, 1.0, -wv.z));
	NORMAL = normalize((VIEW_MATRIX * vec4(wn, 0.0)).xyz);

	// colour: shallow turquoise over sand, deep blue further out (light is absorbed with depth)
	float absorb = 1.0 - exp(-d * 0.16);
	vec3 col = mix(shallow_col, deep_col, absorb);
	// subsurface scattering: light shining through the thin crests
	float crest = clamp(swell * 6.0 + wv.x * 3.0, 0.0, 1.0);
	col += scatter_col * crest * 0.35 * (1.0 - absorb * 0.5);

	// foam: lapping edge at the shore and whitecaps on steep crests
	float foam = 0.0;
	if (quality > 1.5) {
		float n1 = vn(wpos.xz * 1.6 + vec2(TIME * 0.25, TIME * 0.17));
		float n2 = vn(wpos.xz * 4.1 - vec2(TIME * 0.31, TIME * 0.2));
		float lap = 0.5 + 0.5 * sin(TIME * 0.9 + wpos.x * 0.21 + wpos.z * 0.17);
		float edge = 1.0 - smoothstep(0.0, 0.22 + 0.5 * lap, d);
		float band = (1.0 - smoothstep(0.35, 1.4 + 0.6 * lap, d)) * smoothstep(0.52, 0.72, n1 * 0.6 + n2 * 0.5);
		float steep = smoothstep(0.30, 0.50, length(wv.yz)) * smoothstep(0.45, 0.75, n2) * detail;
		foam = clamp(edge * (0.55 + 0.45 * n1) + band * 0.8 + steep * 0.7, 0.0, 1.0);
	}
	// fresnel: nearly transparent when looking down, a mirror at grazing angles
	float fres = pow(1.0 - clamp(dot(NORMAL, VIEW), 0.0, 1.0), 4.0);
	col = mix(col, sky_col, fres * 0.55); // the bright horizon is mirrored in the water
	ALBEDO = mix(col, foam_col, foam);
	float alpha = clamp(0.45 + absorb * 0.5, 0.45, 0.95);
	ALPHA = clamp(mix(alpha, 1.0, fres * 0.7) + foam * 0.5, 0.0, 1.0);
	ROUGHNESS = mix(0.04, 0.8, foam);
	SPECULAR = mix(0.62, 0.15, foam);
}
