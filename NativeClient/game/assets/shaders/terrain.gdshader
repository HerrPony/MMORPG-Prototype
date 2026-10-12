shader_type spatial;
render_mode specular_disabled, cull_back;

// Terrain: displaced by the server's height grid, coloured per biome on the GPU, flat shaded (low poly look).
uniform sampler2D hm_tex : filter_nearest, repeat_disable;
uniform sampler2D bm_tex : filter_nearest, repeat_disable;
uniform sampler2D mask_tex : filter_linear, repeat_disable;
uniform vec2 grid = vec2(721.0, 301.0);
uniform vec2 origin = vec2(-300.0, -300.0);
uniform float cell = 2.0;
uniform vec3 palette[11];
uniform float mosaic = 1.0;

varying vec3 wp;

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
int biome_at(ivec2 q) {
	ivec2 gmax = ivec2(grid) - 1;
	return int(texelFetch(bm_tex, clamp(q, ivec2(0), gmax), 0).r * 255.0 + 0.5);
}
float height_at(ivec2 q) {
	ivec2 gmax = ivec2(grid) - 1;
	return texelFetch(hm_tex, clamp(q, ivec2(0), gmax), 0).r;
}

void vertex() {
	vec2 g = round(UV * (grid - 1.0));
	float h = height_at(ivec2(g));
	VERTEX.y = h;
	wp = vec3(origin.x + g.x * cell, h, origin.y + g.y * cell);
}

void fragment() {
	vec2 gf = (wp.xz - origin) / cell;
	ivec2 gi = ivec2(round(gf));
	vec3 base = vec3(0.0);
	for (int j = -1; j <= 1; j++) {
		for (int i = -1; i <= 1; i++) {
			base += palette[biome_at(gi + ivec2(i * 3, j * 3))];
		}
	}
	base /= 9.0;
	int b = biome_at(gi);
	float h = wp.y;
	vec2 mc = floor(wp.xz * 2.0);
	float n = vn(wp.xz * 0.31) * 0.5 + h21(mc) * 0.5;
	vec3 c = base + vec3(n * 0.04, n * 0.05, n * 0.025);
	vec3 emis = vec3(0.0);
	if (b == 4) {
		float sl = abs(height_at(gi + ivec2(1, 0)) - h) + abs(height_at(gi + ivec2(0, 1)) - h);
		c = mix(c, vec3(0.45, 0.45, 0.48), smoothstep(2.6, 5.0, sl) * 0.85);
	}
	if (b == 3 && h < 0.3) c = mix(c, vec3(0.12, 0.2, 0.12), 0.6);
	if (b == 7 && n > 0.62) c = mix(c, vec3(0.98, 0.76, 0.84), 0.55);
	if (b == 8) {
		c *= 0.88 + 0.12 * sin(h * 2.6 + n * 2.0);
		if (h > 9.0) c = mix(c, vec3(0.86, 0.5, 0.28), 0.25);
	}
	if (b == 9) {
		float ember = max(0.0, sin(wp.x * 0.21) * cos(wp.z * 0.17));
		if (ember > 0.8) {
			float k = (ember - 0.8) * 2.2;
			c = mix(c, vec3(0.9, 0.3, 0.05), k);
			emis = vec3(0.9, 0.25, 0.04) * k * 0.8;
		}
	}
	if (b == 6) c += sin(wp.x * 0.13 + wp.z * 0.1) * 0.02;
	if (h < 0.6 && b != 0 && b < 9) c = mix(c, vec3(0.76, 0.7, 0.5), smoothstep(0.6, 0.0, h) * 0.8);
	if (h < -0.4 && b != 10) c = mix(c, vec3(0.35, 0.34, 0.28), 0.7);
	if (b == 10) {
		c = vec3(1.0, 0.3 + n * 0.25, 0.03);
		emis = vec3(1.0, 0.35, 0.05) * 1.4;
	}
	// roads and paving baked into the mask (r = dirt, g = paving strength, b = paving type)
	vec2 muv = (wp.xz - origin) / ((grid - 1.0) * cell);
	vec4 m = texture(mask_tex, muv);
	c = mix(c, vec3(0.5 + n * 0.05, 0.4, 0.26), m.r * 0.9);
	if (m.g > 0.01) {
		vec3 pave = m.b < 0.2 ? vec3(0.5, 0.42, 0.3) : (m.b < 0.4 ? vec3(0.4, 0.4, 0.42) : (m.b < 0.75 ? vec3(0.62, 0.6, 0.55) : vec3(0.62, 0.5, 0.45)));
		c = mix(c, pave * (0.95 + 0.1 * h21(floor(wp.xz))), m.g);
	}
	// 0.5 m mosaic: every block of ground gets its own slight brightness, like voxel terrain
	c *= 1.0 + (h21(mc + 3.7) - 0.5) * 0.14 * mosaic;
	ALBEDO = c;
	EMISSION = emis;
	ROUGHNESS = 1.0;
	NORMAL = normalize(cross(dFdx(VERTEX), dFdy(VERTEX)));
}
