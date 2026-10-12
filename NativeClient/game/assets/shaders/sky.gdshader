shader_type sky;

// Procedural sky box with sun, moon, stars, drifting clouds, aurora and volcanic glow. All parameters are set per biome
// and time of day by the world (see game_world.gd), the shader itself has no state.
uniform vec3 top_col = vec3(0.2, 0.4, 0.8);
uniform vec3 bot_col = vec3(0.7, 0.8, 0.9);
uniform vec3 sun_dir = vec3(0.0, 1.0, 0.0);
uniform vec3 moon_dir = vec3(0.0, -1.0, 0.0);
uniform vec3 sun_col = vec3(1.0, 0.94, 0.8);
uniform vec3 cloud_col = vec3(1.0);
uniform vec3 glow_col = vec3(1.0, 0.3, 0.1);
uniform float clock = 0.0;
uniform float night = 0.0;
uniform float sun_up = 1.0;
uniform float cloud = 0.5;
uniform float cspeed = 0.004;
uniform float stars = 1.0;
uniform float aurora = 0.0;
uniform float haze = 0.2;
uniform float glow = 0.0;

float h21(vec2 q) {
	q = fract(q * vec2(123.34, 456.21));
	q += dot(q, q + 45.32);
	return fract(q.x * q.y);
}
float h31(vec3 q) {
	q = fract(q * 0.3183099 + 0.1);
	q *= 17.0;
	return fract(q.x * q.y * q.z * (q.x + q.y + q.z));
}
float vn(vec2 x) {
	vec2 i = floor(x);
	vec2 f = fract(x);
	f = f * f * (3.0 - 2.0 * f);
	return mix(mix(h21(i), h21(i + vec2(1.0, 0.0)), f.x), mix(h21(i + vec2(0.0, 1.0)), h21(i + vec2(1.0, 1.0)), f.x), f.y);
}
float fbm(vec2 x) {
	float v = 0.0;
	float a = 0.5;
	for (int i = 0; i < 5; i++) {
		v += a * vn(x);
		x = x * 2.03 + 17.0;
		a *= 0.5;
	}
	return v;
}

void sky() {
	vec3 d = normalize(EYEDIR);
	float hh = clamp(d.y * 1.6, 0.0, 1.0);
	vec3 col = mix(bot_col, top_col, pow(hh, 0.7));
	float cm = 0.0;
	vec3 cc = cloud_col;
	if (d.y > 0.0) {
		vec2 uv = d.xz / (d.y + 0.22) * 0.85 + vec2(clock * cspeed, clock * cspeed * 0.35);
		// blocky clouds: quantise the sample position
		uv = floor(uv * 14.0) / 14.0 * 0.5 + uv * 0.5;
		float c = fbm(uv * 1.25);
		cm = smoothstep(1.0 - cloud, 1.0 - cloud + 0.3, c) * smoothstep(0.0, 0.22, d.y);
		float lit = 0.45 + 0.55 * smoothstep(0.35, 0.75, c);
		cc = cloud_col * mix(vec3(0.16, 0.2, 0.32), vec3(1.0), 1.0 - night) * (0.7 + 0.3 * lit);
		cc += sun_col * pow(max(dot(d, sun_dir), 0.0), 6.0) * 0.35 * sun_up;
		cc += glow_col * glow * 0.35 * (1.0 - hh);
	}
	float sd = max(dot(d, sun_dir), 0.0);
	float md = max(dot(d, moon_dir), 0.0);
	col += sun_col * (pow(sd, 5.0) * 0.12 + pow(sd, 40.0) * 0.4) * sun_up * (1.0 - 0.5 * cm);
	// square sun and moon (it is a voxel world)
	vec3 sr = normalize(cross(sun_dir, vec3(0.0, 0.0, 1.0)));
	vec3 su = cross(sr, sun_dir);
	float sq = max(abs(dot(d, sr)), abs(dot(d, su)));
	col += sun_col * (1.0 - smoothstep(0.034, 0.038, sq)) * step(0.9, sd) * 2.4 * sun_up * (1.0 - cm);
	vec3 mr = normalize(cross(moon_dir, vec3(0.0, 0.0, 1.0)));
	vec3 mu = cross(mr, moon_dir);
	float mq = max(abs(dot(d, mr)), abs(dot(d, mu)));
	float moon_disc = (1.0 - smoothstep(0.03, 0.033, mq)) * step(0.9, md);
	col += (vec3(0.82, 0.88, 1.0) * (0.85 + 0.15 * vn(d.xz * 160.0)) * moon_disc * (1.0 - cm) + vec3(0.35, 0.45, 0.7) * pow(md, 60.0) * 0.5) * night;
	if (night > 0.02 && d.y > 0.0) {
		vec3 q = d * 150.0;
		float s = h31(floor(q));
		vec3 f = fract(q) - 0.5;
		float st = step(0.9955, s) * smoothstep(0.42, 0.0, length(f)) * (0.6 + 0.4 * sin(clock * 2.0 + s * 80.0));
		vec3 q2 = d * 60.0;
		float s2 = h31(floor(q2) + 7.0);
		vec3 f2 = fract(q2) - 0.5;
		st += step(0.9975, s2) * smoothstep(0.4, 0.0, length(f2)) * 1.3;
		col += vec3(0.85, 0.9, 1.0) * st * stars * night * smoothstep(0.0, 0.18, d.y) * (1.0 - cm);
	}
	if (aurora > 0.01 && night > 0.05 && d.y > 0.0) {
		float band = sin(d.x * 5.0 + fbm(d.xz * 3.0 + clock * 0.04) * 7.0) * 0.5 + 0.5;
		float yy = smoothstep(0.12, 0.38, d.y) * (1.0 - smoothstep(0.5, 0.9, d.y));
		col += mix(vec3(0.1, 0.95, 0.5), vec3(0.6, 0.3, 0.95), d.x * 0.5 + 0.5) * band * band * yy * aurora * night * 0.6;
	}
	col = mix(col, cc, cm * 0.92);
	col = mix(col, bot_col, haze * pow(1.0 - hh, 2.0));
	col += glow_col * glow * pow(1.0 - hh, 3.0) * (0.35 + 0.65 * night) * 0.8;
	COLOR = col;
}
