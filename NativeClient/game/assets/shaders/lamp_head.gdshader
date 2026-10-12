shader_type spatial;
render_mode specular_schlick_ggx;

// Glass of the street lamps: a dim pane by day, a flickering flame behind frosted glass at night.
uniform float night = 0.0;
varying vec3 wp;

float h21(vec2 p) {
	p = fract(p * vec2(123.34, 456.21));
	p += dot(p, p + 45.32);
	return fract(p.x * p.y);
}

void vertex() {
	wp = (MODEL_MATRIX * vec4(VERTEX, 1.0)).xyz;
}

void fragment() {
	float id = h21(floor(wp.xz * 0.5));
	// each lamp lights up a little earlier or later than its neighbours
	float on = smoothstep(id * 0.3 + 0.1, id * 0.3 + 0.4, night);
	float flame = 0.92 + 0.08 * sin(TIME * 8.0 + id * 40.0) + 0.05 * sin(TIME * 21.0 + id * 90.0);
	vec3 col = COLOR.rgb;
	ALBEDO = col * 0.3;
	EMISSION = col * vec3(1.0, 0.8, 0.5) * (0.05 + 3.2 * on * flame);
	ROUGHNESS = 0.25;
	SPECULAR = 0.6;
}
