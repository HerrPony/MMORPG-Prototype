shader_type canvas_item;
// Top-down map of the whole world drawn on the GPU from the terrain textures (minimap and world map).
uniform sampler2D hm_tex : filter_nearest, repeat_disable;
uniform sampler2D bm_tex : filter_nearest, repeat_disable;
uniform sampler2D map_tex : filter_linear, repeat_disable;
uniform vec2 grid = vec2(721.0, 301.0);
uniform vec2 origin = vec2(-300.0, -300.0);
uniform float cell = 2.0;
uniform vec2 center = vec2(0.0);
uniform vec2 view = vec2(60.0);
uniform vec2 right = vec2(1.0, 0.0);
uniform vec2 fwd = vec2(0.0, -1.0);
uniform bool circle = true;
uniform bool use_fog = false;
uniform float known[11];

const vec3 PAL[11] = vec3[11](
	vec3(40.0, 90.0, 140.0), vec3(62.0, 122.0, 58.0), vec3(170.0, 182.0, 84.0), vec3(86.0, 106.0, 60.0), vec3(226.0, 234.0, 244.0), vec3(212.0, 182.0, 106.0),
	vec3(22.0, 140.0, 72.0), vec3(238.0, 170.0, 200.0), vec3(200.0, 96.0, 56.0), vec3(78.0, 62.0, 70.0), vec3(255.0, 96.0, 30.0));

float hgt(vec2 gi) {
	return texture(hm_tex, (gi + 0.5) / grid).r;
}

void fragment() {
	vec2 p = (UV - 0.5) * 2.0;
	if (circle && length(p) > 1.0) {
		discard;
	}
	vec2 off = p * view;
	vec2 w = center + right * off.x - fwd * off.y;
	vec2 g = floor((w - origin) / cell);
	if (g.x < 0.0 || g.y < 0.0 || g.x >= grid.x || g.y >= grid.y) {
		COLOR = vec4(0.05, 0.08, 0.11, 1.0);
	} else {
		float b = floor(texture(bm_tex, (g + 0.5) / grid).r * 255.0 + 0.5);
		float h0 = hgt(g);
		vec3 col = PAL[int(clamp(b, 0.0, 10.0))] / 255.0;
		float slope = (hgt(g + vec2(-1.0, -1.0)) - hgt(g + vec2(1.0, 1.0))) * 0.05;
		float shade = 0.82 + clamp(h0 * 0.012, -0.2, 0.3) + clamp(slope, -0.2, 0.2);
		col *= shade;
		if (b > 9.5) {
			col = PAL[10] / 255.0;
		} else if (h0 < 0.0 && b > 0.5) {
			col = vec3(50.0, 110.0, 160.0) / 255.0;
		}
		vec2 muv = (w - origin) / ((grid - 1.0) * cell);
		vec4 m = texture(map_tex, muv);
		if (b > 0.5 && h0 >= 0.0) {
			col = mix(col, vec3(0.59, 0.45, 0.27), smoothstep(0.25, 0.7, m.r) * 0.95);
			col = mix(col, vec3(0.72, 0.69, 0.62), smoothstep(0.1, 0.5, m.g) * 0.9);
		}
		if (m.b > 0.7) {
			col = vec3(0.47, 0.2, 0.16);
		} else if (m.b > 0.3) {
			col = vec3(0.54, 0.42, 0.25);
		}
		if (m.a > 0.5) {
			col = vec3(0.75, 0.73, 0.69);
		}
		if (use_fog) {
			float kb = b > 9.5 ? 9.0 : b;
			float kn = known[int(clamp(kb, 0.0, 10.0))];
			col = mix(vec3(0.055, 0.063, 0.086), col, mix(0.16, 1.0, kn));
		}
		COLOR = vec4(col, 1.0);
	}
}
