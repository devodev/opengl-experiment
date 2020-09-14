package components

import (
	"math"

	"github.com/devodev/opengl-experimentation/internal/opengl"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Square .
type Square struct {
	vao    *opengl.VAO
	ibo    *opengl.IBO
	shader *opengl.ShaderProgram
}

// NewSquare .
func NewSquare() (*Square, error) {
	// create a program with default shaders
	shaderProgram, err := opengl.NewShaderProgram(
		"assets/shaders/vertexDefault.glsl",
		"assets/shaders/fragmentVariableColor.glsl",
	)
	if err != nil {
		return nil, err
	}

	// create vao
	square := []float32{
		-0.5, 0.5, 1, 0,
		-0.5, -0.5, 0, 1,
		0.5, -0.5, 0, 0,
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
	// This declares the interleaving layout of VBO.
	// Here, each `square` vertex contains:
	// 2 * vec2 float32 normalized values
	// *These map directly to vertex shader attributes
	vbo.AddElement(2, false)
	vbo.AddElement(2, false)
	vao := opengl.NewVAO()
	vao.AddVBO(vbo)
	ibo := opengl.NewIBO(squareIndices)

	component := &Square{
		vao:    vao,
		ibo:    ibo,
		shader: shaderProgram,
	}
	return component, nil
}

// OnInit .
func (c *Square) OnInit(window *glfw.Window) {
}

// OnUpdate .
func (c *Square) OnUpdate() {
	// select shader program
	c.shader.Bind()

	// update uniform value
	currentTime := glfw.GetTime()
	greenValue := float32((math.Sin(currentTime) / 2.0) + 0.5)
	c.shader.SetUniform1f("variableColor", greenValue)

	// draw
	c.vao.Bind()
	// TODO: might want to benchmark this and maybe remove them
	// TODO: or do that only in debug mode or something..
	defer c.vao.Unbind()

	c.ibo.Bind()
	// TODO: might want to benchmark this and maybe remove them
	// TODO: or do that only in debug mode or something..
	defer c.ibo.Unbind()

	gl.DrawElements(gl.TRIANGLES, c.ibo.GetCount(), gl.UNSIGNED_INT, nil)
}
