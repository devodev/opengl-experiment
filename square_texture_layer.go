package main

import (
	"fmt"

	"github.com/devodev/opengl-experimentation/internal/engine/application"
	"github.com/devodev/opengl-experimentation/internal/engine/renderer"
	"github.com/devodev/opengl-experimentation/internal/opengl"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// SquareTextureLayer .
type SquareTextureLayer struct {
	texture1 *opengl.Texture
	texture2 *opengl.Texture
	texture3 *opengl.Texture
	camera   *renderer.Camera
}

// NewSquareTextureLayer .
func NewSquareTextureLayer(app *application.Application) (*SquareTextureLayer, error) {
	texture1, err := opengl.NewTexture("assets/textures/google_logo.png", 0)
	if err != nil {
		return nil, fmt.Errorf("error creating texture1: %s", err)
	}
	texture2, err := opengl.NewTexture("assets/textures/facebook_logo.png", 10)
	if err != nil {
		return nil, fmt.Errorf("error creating texture2: %s", err)
	}
	texture3, err := opengl.NewTexture("assets/textures/instagram_logo.png", 15)
	if err != nil {
		return nil, fmt.Errorf("error creating texture3: %s", err)
	}

	width, height := app.GetWindow().GetGLFWWindow().GetSize()
	camera := renderer.NewCamera(width, height)

	component := &SquareTextureLayer{
		texture1: texture1,
		texture2: texture2,
		texture3: texture3,
		camera:   camera,
	}
	return component, nil
}

// OnInit .
func (c *SquareTextureLayer) OnInit(app *application.Application) {
}

// OnUpdate .
func (c *SquareTextureLayer) OnUpdate(app *application.Application, deltaTime float64) {
	c.processInput(app)
	c.camera.OnUpdate(app.GetWindow().GetGLFWWindow(), deltaTime)
}

// OnRender .
func (c *SquareTextureLayer) OnRender(app *application.Application, deltaTime float64) {
	pos1 := mgl32.Translate3D(-0.5, 0, -0.5)
	pos2 := mgl32.Translate3D(0.5, 0, 1)
	pos3 := mgl32.Translate3D(0, 0.5, 0.5)
	pos4 := mgl32.Translate3D(0, -0.5, -1)
	pos5 := mgl32.Translate3D(0, 0.5, -1)

	app.GetRenderer().Begin(c.camera)
	app.GetRenderer().DrawTexturedQuad(pos1, c.texture1)
	app.GetRenderer().DrawTexturedQuad(pos2, c.texture1)
	app.GetRenderer().DrawTexturedQuad(pos5, c.texture2)
	app.GetRenderer().DrawTexturedQuad(pos3, c.texture2)
	app.GetRenderer().DrawTexturedQuad(pos4, c.texture2)
	app.GetRenderer().DrawTexturedQuad(pos3, c.texture3)
	app.GetRenderer().DrawTexturedQuad(pos4, c.texture3)
	app.GetRenderer().End()
}

func (c *SquareTextureLayer) processInput(app *application.Application) {
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
