shader_type spatial;
render_mode specular_schlick_ggx, cull_back;

// Terrain: displaced by the server's height grid and coloured per biome on the GPU. The ground gets a surface of its own:
// blades and tufts of grass, rippled sand, pebbles on the paths, cobblestones, rock strata on steep slopes and sparkling snow,
// lit with a smooth normal plus a fine bump.
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

float fbm(vec2 p) {
	float a = 0.5;
	float s = 0.0;
	for (int i = 0; i < 4; i++) {
		s += a * vn(p);
		p = p * 2.03 + vec2(17.1, 9.7);
		a *= 0.5;
	}
	return s / 0.9375;
}
// bilinear terrain height at a fractional grid position
float hb(vec2 g) {
	vec2 f = fract(g);
	ivec2 i = ivec2(floor(g));
	return mix(mix(height_at(i), height_at(i + ivec2(1, 0)), f.x), mix(height_at(i + ivec2(0, 1)), height_at(i + ivec2(1, 1)), f.x), f.y);
}
vec3 flower(vec2 p, float id) {
	if (id < 0.25) return vec3(0.95, 0.95, 0.9);
	if (id < 0.5) return vec3(0.98, 0.85, 0.2);
	if (id < 0.75) return vec3(0.9, 0.35, 0.55);
	return vec3(0.4, 0.5, 0.95);
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
	vec2 p = wp.xz;
	float dist = length(VERTEX);
	float lod = 1.0 / (1.0 + dist * dist / 9000.0);       // fine detail fades out in the distance (no shimmering)
	float macro = fbm(p * 0.045);
	float mid = fbm(p * 0.33 + vec2(5.0, 2.0));
	float fine = (vn(p * 3.1) * 0.55 + vn(p * 8.3) * 0.45);

	// smooth terrain normal from the height grid (world space)
	vec3 N = normalize(vec3(hb(gf - vec2(1.0, 0.0)) - hb(gf + vec2(1.0, 0.0)), 2.0 * cell, hb(gf - vec2(0.0, 1.0)) - hb(gf + vec2(0.0, 1.0))));
	float steep = 1.0 - N.y;

	vec3 c = base * (0.9 + 0.2 * macro);
	vec3 emis = vec3(0.0);
	float rough = 0.95;
	float bump = fine * 0.5 + mid * 0.5;   // height of the small scale surface, for the bump
	float bump_k = 0.35;

	bool grass = (b == 1 || b == 2 || b == 3 || b == 6 || b == 7);
	if (grass) {
		// blades and tufts: dark clumps, light tips, patches of dry grass
		float blades = vn(p * vec2(23.0, 17.0)) * 0.6 + vn(p * vec2(9.0, 13.0)) * 0.4;
		float tuft = smoothstep(0.55, 0.8, vn(p * 2.3 + 4.0));
		float dry = smoothstep(0.58, 0.78, macro) * (b == 2 ? 1.0 : 0.5);
		vec3 dryc = vec3(0.62, 0.56, 0.22);
		c = mix(c, c * dryc / max(base, vec3(0.05)), dry * 0.45);
		float gv = mix(1.0, 0.72 + blades * 0.6, lod);
		c *= gv * (1.0 - tuft * 0.18 * lod) * (0.92 + 0.16 * mid);
		bump = blades * 0.5 + tuft * 0.5;
		bump_k = 0.3;
		rough = 0.92;
		// flowers
		if (b == 2 || b == 1 || b == 7) {
			vec2 fc = floor(p * 2.4);
			float fid = h21(fc);
			float thr = b == 2 ? 0.965 : 0.985;
			if (fid > thr && h < 6.0) {
				vec2 centre = (fc + vec2(0.3 + 0.4 * h21(fc + 7.0), 0.3 + 0.4 * h21(fc + 3.0))) / 2.4;
				float fd = length(p - centre);
				float on = smoothstep(0.09, 0.05, fd) * lod;
				c = mix(c, flower(p, h21(fc + 11.0)), on);
			}
		}
	}
	if (b == 4) {
		// snow: soft drifts with blue shadows and glittering crystals
		float drift = fbm(p * 0.18);
		c = mix(c, c * vec3(0.84, 0.9, 1.0), smoothstep(0.35, 0.7, drift) * 0.45);
		float spark = pow(vn(p * 45.0 + floor(TIME * 0.5)), 14.0) * smoothstep(0.2, 0.0, steep) * lod;
		c += vec3(spark * 0.9);
		float sl = abs(height_at(gi + ivec2(1, 0)) - h) + abs(height_at(gi + ivec2(0, 1)) - h);
		c = mix(c, vec3(0.45, 0.45, 0.48), smoothstep(2.6, 5.0, sl) * 0.85);
		bump = drift;
		bump_k = 0.2;
		rough = 0.55;
	}
	if (b == 3 && h < 0.3) {
		c = mix(c, vec3(0.12, 0.2, 0.12), 0.6);
		rough = 0.35;   // wet mud shines
	}
	if (b == 7 && vn(p * 1.2) > 0.62) c = mix(c, vec3(0.98, 0.76, 0.84), 0.55 * lod + 0.3);   // fallen blossoms
	if (b == 5 || (h < 0.9 && b != 0 && b < 9 && b != 4)) {
		// sand: wind ripples, grains and darker troughs
		float rip = sin(dot(p, vec2(0.62, 0.42)) * 2.6 + fbm(p * 0.18) * 9.0);
		float ripple = smoothstep(-0.2, 0.9, rip);
		float grain = vn(p * 26.0);
		float sand_k = b == 5 ? 1.0 : smoothstep(0.9, 0.2, h);
		vec3 sandc = base * 0.0 + (b == 5 ? c : vec3(0.78, 0.7, 0.5));
		sandc *= (0.9 + 0.12 * ripple) * (0.94 + 0.1 * grain * lod);
		c = mix(c, sandc, sand_k * 0.92);
		if (sand_k > 0.5) {
			bump = ripple * 0.7 + grain * 0.15;
			bump_k = 0.55;
			rough = 0.97;
		}
	}
	if (b == 8) {
		// mesa: coloured strata
		float strata = sin(h * 2.6 + fbm(p * 0.07) * 3.0) * 0.5 + 0.5;
		c *= 0.82 + 0.28 * strata;
		if (h > 9.0) c = mix(c, vec3(0.86, 0.5, 0.28), 0.25);
		bump = strata * 0.6 + fine * 0.3;
	}
	if (b == 9) {
		// ash and cooled lava: dark crust, glowing cracks
		float ember = max(0.0, sin(wp.x * 0.21) * cos(wp.z * 0.17));
		float crack = smoothstep(0.045, 0.0, abs(fbm(p * 0.22) - 0.5)) * 0.8;
		c *= 0.8 + 0.4 * fine;
		if (ember > 0.8 || crack > 0.1) {
			float k = max((ember - 0.8) * 2.2, crack);
			c = mix(c, vec3(0.9, 0.3, 0.05), k);
			emis = vec3(0.9, 0.25, 0.04) * k * 0.9;
		}
		bump = fine;
		bump_k = 0.6;
	}
	if (b == 6) c += sin(wp.x * 0.13 + wp.z * 0.1) * 0.02;
	if (h < -0.4 && b != 10) c = mix(c, vec3(0.35, 0.34, 0.28), 0.7);
	if (b == 10) {
		c = vec3(1.0, 0.3 + fine * 0.25, 0.03);
		emis = vec3(1.0, 0.35, 0.05) * 1.4;
		rough = 0.4;
	}

	// rock: steep slopes show bare stone in strata
	if (b != 10 && b != 0 && h > 0.3) {
		float rk = smoothstep(0.26, 0.5, steep);
		float strata2 = fbm(vec2(p.x * 0.5 + p.y * 0.35, h * 1.8));
		vec3 rock = mix(vec3(0.38, 0.36, 0.34), vec3(0.55, 0.5, 0.45), strata2);
		rock = mix(rock, base * 0.75, b == 8 || b == 5 ? 0.6 : 0.12);
		rock *= 0.8 + 0.35 * vn(p * 4.0 + h * 3.0) * lod + 0.2 * fine * lod;
		c = mix(c, rock, rk * (b == 4 ? 0.7 : 0.9));
		bump = mix(bump, strata2 * 0.7 + fine * 0.5, rk);
		bump_k = mix(bump_k, 0.9, rk);
		rough = mix(rough, 0.85, rk);
	}

	// roads and paving baked into the mask (r = dirt, g = paving strength, b = paving type)
	vec2 muv = (wp.xz - origin) / ((grid - 1.0) * cell);
	vec4 m = texture(mask_tex, muv);
	if (m.r > 0.01) {
		// dirt track: packed earth with pebbles and worn ruts
		float peb = smoothstep(0.68, 0.92, vn(p * 9.0)) * 0.8;
		vec3 dirt = vec3(0.5 + fine * 0.06, 0.4, 0.27) * (0.88 + 0.2 * mid);
		dirt = mix(dirt, vec3(0.6, 0.57, 0.52) * (0.8 + 0.3 * vn(p * 31.0)), peb * 0.35 * lod);
		c = mix(c, dirt, m.r * 0.92);
		bump = mix(bump, mid * 0.6 + peb * 0.15, m.r);
		bump_k = mix(bump_k, 0.6, m.r);
		rough = mix(rough, 0.97, m.r);
	}
	if (m.g > 0.01) {
		// cobblestones: every stone has its own tone, the joints are dark
		vec3 pave = m.b < 0.2 ? vec3(0.5, 0.42, 0.3) : (m.b < 0.4 ? vec3(0.4, 0.4, 0.42) : (m.b < 0.75 ? vec3(0.62, 0.6, 0.55) : vec3(0.62, 0.5, 0.45)));
		vec2 q = p / vec2(0.44, 0.34);
		float row = floor(q.y);
		q.x += 0.5 * mod(row, 2.0) + h21(vec2(row, 1.0)) * 2.0;
		vec2 id = vec2(floor(q.x), row);
		vec2 d = abs(fract(q) - 0.5) * 2.0;
		float e = pow(pow(d.x, 5.0) + pow(d.y, 5.0), 0.2);       // rounded stone outline
		float stone = 1.0 - smoothstep(0.74, 0.95, e);
		float tone = h21(id);
		vec3 cob = pave * (0.74 + 0.44 * tone) * (0.9 + 0.2 * vn(p * 14.0));
		cob *= mix(0.4, 1.0, mix(1.0, stone, lod));
		c = mix(c, cob, m.g);
		bump = mix(bump, stone * (1.0 - e * 0.5) + tone * 0.1, m.g);
		bump_k = mix(bump_k, 0.8, m.g);
		rough = mix(rough, 0.8, m.g);
	}
	// a light mosaic of the old voxel look remains, much weaker than before
	vec2 mc = floor(p * 2.0);
	c *= 1.0 + (h21(mc + 3.7) - 0.5) * 0.05 * mosaic;
	ALBEDO = c;
	EMISSION = emis;
	ROUGHNESS = rough;
	SPECULAR = 0.3;

	// normal: smooth terrain normal plus the bump of the surface
	vec3 Nv = normalize((VIEW_MATRIX * vec4(N, 0.0)).xyz);
	vec3 dp1 = dFdx(VERTEX);
	vec3 dp2 = dFdy(VERTEX);
	vec3 r1 = cross(dp2, Nv);
	vec3 r2 = cross(Nv, dp1);
	float det = dot(dp1, r1);
	float hh = bump * bump_k * lod;
	vec3 grad = sign(det) * (dFdx(hh) * r1 + dFdy(hh) * r2);
	NORMAL = normalize(abs(det) * Nv - grad * 0.6);
}
