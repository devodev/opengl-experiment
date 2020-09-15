package renderer

const (
	quadVertexShader = `
#version 460
layout (location = 0) in vec2 position;
layout (location = 1) in vec2 texCoord;

out vec2 fragTexCoord;

uniform mat4 mvp;

void main() {
    gl_Position = mvp * vec4(position, 0.0, 1.0);
    fragTexCoord = texCoord;
}
    `

	quadFragmentShader = `
#version 460
layout (location = 0) out vec4 fragColor;

in vec2 fragTexCoord;

uniform sampler2D tex;

void main() {
    fragColor = texture(tex, fragTexCoord);
}
    `
)
