package renderer

import (
	"github.com/devodev/opengl-experiment/internal/opengl"
	"github.com/go-gl/mathgl/mgl32"
)

type TexturedQuad struct {
	Transform mgl32.Mat4
	Texture   opengl.Texture
}
