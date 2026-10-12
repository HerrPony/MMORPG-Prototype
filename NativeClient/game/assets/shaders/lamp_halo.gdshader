shader_type spatial;
render_mode blend_add, unshaded, depth_draw_never, cull_disabled, fog_disabled;

// Soft glow around street lamps and lanterns: a camera facing quad with a bright core and a wide, soft falloff.
// Per instance (custom data): x = flicker phase, y = size in metres, z = strength.
uniform float on = 0.0;
uniform vec3 tint = vec3(1.0, 0.62, 0.28);

varying float k;
varying float fl;

void vertex() {
	mat4 mv = VIEW_MATRIX * mat4(INV_VIEW_MATRIX[0], INV_VIEW_MATRIX[1], INV_VIEW_MATRIX[2], MODEL_MATRIX[3]);
	vec3 centre = (VIEW_MATRIX * vec4(MODEL_MATRIX[3].xyz, 1.0)).xyz;
	float dist = length(centre);
	// fades out in the distance and when the camera is right next to the lamp
	k = INSTANCE_CUSTOM.z * smoothstep(190.0, 80.0, dist) * smoothstep(0.8, 4.0, dist);
	fl = 1.0 + 0.08 * sin(TIME * 9.0 + INSTANCE_CUSTOM.x * 6.28) + 0.05 * sin(TIME * 23.0 + INSTANCE_CUSTOM.x * 17.0);
	// a bit bigger when seen from afar so that it stays visible
	VERTEX *= INSTANCE_CUSTOM.y * (1.0 + dist * 0.004);
	MODELVIEW_MATRIX = mv;
	MODELVIEW_NORMAL_MATRIX = mat3(mv);
}

void fragment() {
	float r = length(UV - vec2(0.5)) * 2.0;
	float core = exp(-r * r * 22.0);
	float halo = exp(-r * r * 4.5) * 0.55;
	float rays = pow(max(0.0, 1.0 - r), 3.0) * 0.25;
	float g = (core * 1.6 + halo + rays) * step(r, 1.0) * (1.0 - smoothstep(0.82, 1.0, r));
	ALBEDO = mix(tint, vec3(1.0, 0.92, 0.75), core) * g * on * k * fl * 0.9;
	ALPHA = 1.0;
}
