#version 460
layout (location = 0) out vec4 fragColor;

in vec2 fragTexCoord;

uniform sampler2D tex;

void main() {
    fragColor = texture(tex, fragTexCoord);
}
