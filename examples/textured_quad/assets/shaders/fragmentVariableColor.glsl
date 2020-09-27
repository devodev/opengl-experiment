#version 460
layout (location = 0) out vec4 frag_color;

in vec2 out_color;

uniform float variableColor;

void main() {
    frag_color = vec4(out_color, variableColor, 1.0);
}
