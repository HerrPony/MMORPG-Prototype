shader_type spatial;
render_mode specular_schlick_ggx;

// Buildings: the blocky meshes keep their colours (vertex colour), but get a surface of their own - bricks, planks,
// stucco and roof shingles - with a bumped normal, so that sun, moon and lamp light play over the walls.
uniform float detail = 1.0;   // 0 = plain, 1 = full surface detail

varying vec3 wp;
varying vec3 wn;

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
vec3 srgb(vec3 c) {
	return mix(c / 12.92, pow((c + 0.055) / 1.055, vec3(2.4)), step(vec3(0.04045), c));
}

void vertex() {
	wp = (MODEL_MATRIX * vec4(VERTEX, 1.0)).xyz;
	wn = normalize((MODEL_MATRIX * vec4(NORMAL, 0.0)).xyz);
}

void fragment() {
	vec3 col = srgb(COLOR.rgb);
	float lum = dot(col, vec3(0.3, 0.59, 0.11));
	float sat = max(col.r, max(col.g, col.b)) - min(col.r, min(col.g, col.b));
	float dist = length(VERTEX);
	float fade = detail / (1.0 + dist * dist / 2400.0);

	float up = wn.y;
	float h = 0.0;      // height of the surface pattern (for the bump)
	float shade = 1.0;  // brightness variation of the pattern
	float rough = 0.9;
	if (up > 0.85 || up < -0.85) {
		// flat tops and undersides: weathered, with a few stains
		float n = vn(wp.xz * 0.6) * 0.6 + vn(wp.xz * 2.7) * 0.4;
		shade = 0.9 + n * 0.2;
		h = n * 0.3;
	} else if (up > 0.25) {
		// roof: rows of shingles, every other row shifted, with a rounded lower edge
		vec2 q = vec2((abs(wn.x) > abs(wn.z)) ? wp.z : wp.x, wp.y / max(0.35, sqrt(1.0 - up * up)));
		float row = floor(q.y / 0.34);
		float x = q.x / 0.5 + 0.5 * mod(row, 2.0);
		vec2 id = vec2(floor(x), row);
		float fy = fract(q.y / 0.34);
		float fx = fract(x);
		float edge = smoothstep(0.0, 0.28, fy) * (1.0 - smoothstep(0.92, 1.0, fy));
		float side = smoothstep(0.0, 0.07, fx) * smoothstep(0.0, 0.07, 1.0 - fx);
		float tile = h21(id);
		h = edge * side * (0.55 + 0.45 * fy) + tile * 0.1;
		shade = (0.78 + 0.3 * tile) * mix(0.62, 1.0, edge * side);
		rough = 0.62;
	} else if (lum > 0.62) {
		// plaster: fine stucco with a few cracks
		float n = vn(wp.xz * 5.0 + wp.y * 3.0) * 0.5 + vn(wp.xz * 17.0 + wp.y * 11.0) * 0.5;
		float crack = smoothstep(0.015, 0.0, abs(vn(vec2(wp.x + wp.z, wp.y) * 2.3) - 0.5));
		h = n * 0.35 - crack * 0.4;
		shade = 0.93 + n * 0.12 - crack * 0.25;
	} else if (sat < 0.14) {
		// stone: blocks in courses, every course shifted, each block its own tone
		float u = (abs(wn.x) > abs(wn.z)) ? wp.z : wp.x;
		float row = floor(wp.y / 0.45);
		float x = u / 0.9 + 0.5 * mod(row, 2.0) + h21(vec2(row, 3.0)) * 3.0;
		vec2 id = vec2(floor(x), row);
		float fx = fract(x);
		float fy = fract(wp.y / 0.45);
		float mortar = smoothstep(0.0, 0.05, fx) * smoothstep(0.0, 0.05, 1.0 - fx) * smoothstep(0.0, 0.07, fy) * smoothstep(0.0, 0.07, 1.0 - fy);
		float tone = h21(id);
		float grain = vn(wp.xz * 9.0 + wp.y * 9.0);
		h = mortar * (0.6 + 0.3 * tone) + grain * 0.12;
		shade = mix(0.55, 0.85 + 0.3 * tone, mortar) * (0.95 + 0.1 * grain);
	} else {
		// timber: vertical planks with grain and the odd knot
		float u = (abs(wn.x) > abs(wn.z)) ? wp.z : wp.x;
		float plank = floor(u / 0.3);
		float fx = fract(u / 0.3);
		float gap = smoothstep(0.0, 0.07, fx) * smoothstep(0.0, 0.07, 1.0 - fx);
		float grain = vn(vec2(u * 7.0 + plank * 5.3, wp.y * 1.1 + plank * 2.1)) * 0.6 + vn(vec2(u * 31.0, wp.y * 3.0)) * 0.4;
		float tone = h21(vec2(plank, 7.0));
		h = gap * 0.6 + grain * 0.15;
		shade = mix(0.5, 0.82 + 0.28 * tone, gap) * (0.9 + 0.2 * grain);
	}
	shade = mix(1.0, shade, fade);
	ALBEDO = col * shade;
	ROUGHNESS = rough;
	SPECULAR = 0.35;

	// bump: tilt the normal along the gradient of the pattern height
	vec3 dp1 = dFdx(VERTEX);
	vec3 dp2 = dFdy(VERTEX);
	vec3 r1 = cross(dp2, NORMAL);
	vec3 r2 = cross(NORMAL, dp1);
	float det = dot(dp1, r1);
	vec3 grad = sign(det) * (dFdx(h) * r1 + dFdy(h) * r2);
	NORMAL = normalize(abs(det) * NORMAL - grad * 0.5 * fade);
}
