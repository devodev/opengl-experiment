package main

import (
	"fmt"

	"github.com/devodev/opengl-experimentation/internal/engine/application"
	"github.com/devodev/opengl-experimentation/internal/opengl"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	deltaTime = float64(0)
	lastFrame = float64(0)
)

// camera movements
var (
	movingForward  = false
	movingBackward = false
	movingLeft     = false
	movingRight    = false
)

// SquareTexture .
type SquareTexture struct {
	vao     *opengl.VAO
	ibo     *opengl.IBO
	shader  *opengl.ShaderProgram
	texture *opengl.Texture
	camera  *Camera
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

	camera := &Camera{
		pos:   mgl32.Vec3{0, 0, 2},
		front: mgl32.Vec3{0, 0, -1},
		up:    mgl32.Vec3{0, 1, 0},
		speed: float32(0.05),
	}

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
	app.GetWindow().GetGLFWWindow().SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// we lost focus, dont process synthetic events
		if w.GetAttrib(glfw.Focused) == glfw.False {
			return
		}
		// close window
		if key == glfw.KeyEscape && action != glfw.Release {
			w.SetShouldClose(true)
		}
		// toggle wireframes
		if key == glfw.KeySpace && action == glfw.Press {
			var currentPolygonMode int32
			gl.GetIntegerv(gl.POLYGON_MODE, &currentPolygonMode)
			if currentPolygonMode == gl.LINE {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			} else {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			}
		}
		// move camera
		if key == glfw.KeyW {
			movingForward = !(action == glfw.Release)
		}
		if key == glfw.KeyS {
			movingBackward = !(action == glfw.Release)
		}
		if key == glfw.KeyA {
			movingLeft = !(action == glfw.Release)
		}
		if key == glfw.KeyD {
			movingRight = !(action == glfw.Release)
		}
	})
}

// OnUpdate .
func (c *SquareTexture) OnUpdate(app *application.Application) {
	currentTime := glfw.GetTime()
	deltaTime = currentTime - lastFrame
	lastFrame = currentTime

	c.camera.speed = float32(2 * deltaTime)

	if movingForward {
		c.camera.pos = c.camera.pos.Add(c.camera.front.Mul(c.camera.speed))
	}
	if movingBackward {
		c.camera.pos = c.camera.pos.Sub(c.camera.front.Mul(c.camera.speed))
	}
	if movingLeft {
		c.camera.pos = c.camera.pos.Sub(c.camera.front.Normalize().Cross(c.camera.up).Mul(c.camera.speed))
	}
	if movingRight {
		c.camera.pos = c.camera.pos.Add(c.camera.front.Normalize().Cross(c.camera.up).Mul(c.camera.speed))
	}

	width, height := app.GetWindow().GetGLFWWindow().GetSize()
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 10.0)
	camera := mgl32.LookAtV(c.camera.pos, c.camera.pos.Add(c.camera.front), c.camera.up)
	model := mgl32.Ident4()

	mvp := projection.Mul4(camera).Mul4(model)

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

// Camera .
type Camera struct {
	pos   mgl32.Vec3
	front mgl32.Vec3
	up    mgl32.Vec3
	speed float32
}
