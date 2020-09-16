package main

import (
	"fmt"

	"github.com/devodev/opengl-experimentation/internal/engine/application"
	"github.com/devodev/opengl-experimentation/internal/engine/renderer"
	"github.com/devodev/opengl-experimentation/internal/opengl"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// SquareTextureLayer .
type SquareTextureLayer struct {
	texture *opengl.Texture
	camera  *renderer.Camera
}

// NewSquareTextureLayer .
func NewSquareTextureLayer(app *application.Application) (*SquareTextureLayer, error) {
	texture, err := opengl.NewTexture("assets/textures/google_logo.png", 1)
	if err != nil {
		return nil, fmt.Errorf("error creating texture: %s", err)
	}

	width, height := app.GetWindow().GetGLFWWindow().GetSize()
	camera := renderer.NewCamera(width, height)

	component := &SquareTextureLayer{
		texture: texture,
		camera:  camera,
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
	app.GetRenderer().Begin(c.camera)
	app.GetRenderer().DrawTexturedQuad(c.texture)
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
