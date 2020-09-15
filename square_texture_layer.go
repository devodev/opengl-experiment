package main

import (
	"fmt"

	"github.com/devodev/opengl-experimentation/internal/engine/application"
	"github.com/devodev/opengl-experimentation/internal/opengl"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// SquareTexture .
type SquareTexture struct {
	vao     *opengl.VAO
	ibo     *opengl.IBO
	shader  *opengl.ShaderProgram
	texture *opengl.Texture
	camera  *application.Camera
}

// NewSquareTexture .
func NewSquareTexture(app *application.Application) (*SquareTexture, error) {
	shaderProgram, err := opengl.NewShaderProgram(
		"assets/shaders/vertexTexture.glsl",
		"assets/shaders/fragmentTexture.glsl",
	)
	if err != nil {
		return nil, err
	}

	square := []float32{
		// position(vec2), texCoord(vec2)
		-0.5, 0, 0, 1,
		-0.5, -0.5, 0, 0,
		0, -0.5, 1, 0,
		0, 0, 1, 1,

		0, 0.5, 0, 1,
		0, 0, 0, 0,
		0.5, 0, 1, 0,
		0.5, 0.5, 1, 1,
	}
	squareIndices := []uint32{
		0, 1, 2, 2, 3, 0,
		4, 5, 6, 6, 7, 4,
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

	width, height := app.GetWindow().GetGLFWWindow().GetSize()
	camera := application.NewCamera(width, height)

	component := &SquareTexture{
		vao:     vao,
		ibo:     ibo,
		shader:  shaderProgram,
		texture: texture,
		camera:  camera,
	}
	return component, nil
}

// OnInit .
func (c *SquareTexture) OnInit(app *application.Application) {
}

// OnUpdate .
func (c *SquareTexture) OnUpdate(app *application.Application, deltaTime float64) {
	c.processInput(app)
	c.camera.OnUpdate(app, deltaTime)
}

// OnRender .
func (c *SquareTexture) OnRender(app *application.Application, deltaTime float64) {
	vp := c.camera.GetViewProjectionMatrix()
	model := mgl32.Ident4()
	mvp := vp.Mul4(model)

	c.shader.Bind()
	c.shader.SetUniform1i("tex", int32(c.texture.GetTextureUnit()-gl.TEXTURE0))
	c.shader.SetUniformMatrix4fv("mvp", 1, false, &mvp[0])

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

func (c *SquareTexture) processInput(app *application.Application) {
	glfwWindow := app.GetWindow().GetGLFWWindow()
	// we lost focus, dont process synthetic events
	if glfwWindow.GetAttrib(glfw.Focused) == glfw.False {
		return
	}

	// close window
	if glfwWindow.GetKey(glfw.KeyEscape) != glfw.Release {
		app.RequestClose()
		return
	}
	// toggle wireframes
	if glfwWindow.GetKey(glfw.KeySpace) == glfw.Press {
		var currentPolygonMode int32
		gl.GetIntegerv(gl.POLYGON_MODE, &currentPolygonMode)
		if currentPolygonMode == gl.LINE {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		}
	}
}
