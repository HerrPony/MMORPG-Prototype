shader_type spatial;
render_mode specular_schlick_ggx;

// Windows and lanterns of the buildings: dim glass by day; at dusk every window lights up at its own moment and
// glows warm, with a faint candle flicker.
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
	float id = h21(floor(wp.xz * 0.7) + floor(wp.y));
	float on = smoothstep(id * 0.3 + 0.1, id * 0.3 + 0.4, night);
	float flame = 0.93 + 0.07 * sin(TIME * 6.0 + id * 50.0) + 0.04 * sin(TIME * 17.0 + id * 80.0);
	ALBEDO = COLOR.rgb * (0.3 + 0.1 * (1.0 - night));
	EMISSION = COLOR.rgb * vec3(1.0, 0.8, 0.5) * (0.05 + 2.1 * on * flame);
	ROUGHNESS = 0.3;
	SPECULAR = 0.7;
}
