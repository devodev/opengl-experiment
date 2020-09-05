package components

import (
	"fmt"

	"github.com/devodev/opengl-experimentation/internal/opengl"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// SquareTexture .
type SquareTexture struct {
	vao     *opengl.VAO
	ibo     *opengl.IBO
	shader  *opengl.ShaderProgram
	texture *opengl.Texture
}

// NewSquareTexture .
func NewSquareTexture() (*SquareTexture, error) {
	shaderProgram, err := opengl.NewShaderProgram(
		"assets/shaders/vertexTexture.glsl",
		"assets/shaders/fragmentTexture.glsl",
	)
	if err != nil {
		return nil, err
	}

	square := []float32{
		// position(vec2), texCoord(vec2)
		-0.5, 0.5, 0, 1,
		-0.5, -0.5, 0, 0,
		0.5, -0.5, 1, 0,
		0.5, 0.5, 1, 1,
	}
	squareIndices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	vbo, err := opengl.NewVBO(square, opengl.FLOAT)
	if err != nil {
		return nil, err
	}
	vbo.AddElement(2, false)
	vbo.AddElement(2, false)

	vao := opengl.NewVAO()
	vao.AddVBO(vbo)

	ibo := opengl.NewIBO(squareIndices)

	texture, err := opengl.NewTexture("assets/textures/google_logo.png", 1)
	if err != nil {
		return nil, fmt.Errorf("error creating texture: %s", err)
	}

	component := &SquareTexture{
		vao:     vao,
		ibo:     ibo,
		shader:  shaderProgram,
		texture: texture,
	}
	return component, nil
}

// OnUpdate .
func (c *SquareTexture) OnUpdate() {
	c.shader.Bind()

	c.shader.SetUniform1i("tex", int32(c.texture.GetTextureUnit()-gl.TEXTURE0))

	c.texture.Bind()
	c.vao.Bind()
	c.ibo.Bind()
	defer func() {
		c.shader.Unbind()
		c.texture.Unbind()
		c.vao.Unbind()
		c.ibo.Unbind()
	}()

	gl.DrawElements(gl.TRIANGLES, c.ibo.GetCount(), gl.UNSIGNED_INT, nil)
}
