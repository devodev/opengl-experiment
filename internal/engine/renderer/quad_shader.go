package renderer

const (
	quadVertexShader = `
#version 460 core
layout (location = 0) in vec4 position;
layout (location = 1) in vec2 texCoord;
layout (location = 2) in float texIndex;

out vec2 fragTexCoord;
out float fragTexIndex;

uniform mat4 vp;

void main() {
    fragTexCoord = texCoord;
    fragTexIndex = texIndex;
    gl_Position = vp * position;
}
    `

	quadFragmentShader = `
#version 460 core
layout (location = 0) out vec4 fragColor;

in vec2 fragTexCoord;
flat in float fragTexIndex;

uniform sampler2D tex[32];

void main() {
    // switch(int(fragTexIndex)) {
    //     case 0: fragColor = texture(tex[0], fragTexCoord); break;
    //     case 1: fragColor = texture(tex[1], fragTexCoord); break;
    //     case 2: fragColor = texture(tex[2], fragTexCoord); break;
    //     case 3: fragColor = texture(tex[3], fragTexCoord); break;
    //     case 4: fragColor = texture(tex[4], fragTexCoord); break;
    //     case 5: fragColor = texture(tex[5], fragTexCoord); break;
    //     case 6: fragColor = texture(tex[6], fragTexCoord); break;
    //     case 7: fragColor = texture(tex[7], fragTexCoord); break;
    //     case 8: fragColor = texture(tex[8], fragTexCoord); break;
    //     case 9: fragColor = texture(tex[9], fragTexCoord); break;
    //     case 10: fragColor = texture(tex[10], fragTexCoord); break;
    //     case 11: fragColor = texture(tex[11], fragTexCoord); break;
    //     case 12: fragColor = texture(tex[12], fragTexCoord); break;
    //     case 13: fragColor = texture(tex[13], fragTexCoord); break;
    // }
    fragColor = texture(tex[int(fragTexIndex)], fragTexCoord);
    //fragColor = texture(tex[15], fragTexCoord);
    //fragColor = vec4(1,1,1,1);
}
    `
)
