#version 460
out vec4 frag_color;

uniform vec4 variableColor;

void main() {
    frag_color = variableColor;
}
