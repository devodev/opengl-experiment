package main

import (
	"fmt"

	"github.com/devodev/opengl-experimentation/internal/engine/application"
	"github.com/devodev/opengl-experimentation/internal/engine/renderer"
	"github.com/devodev/opengl-experimentation/internal/opengl"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// SquareTexture .
type SquareTexture struct {
	texture *opengl.Texture
	camera  *renderer.Camera
}

// NewSquareTexture .
func NewSquareTexture(app *application.Application) (*SquareTexture, error) {
	texture, err := opengl.NewTexture("assets/textures/google_logo.png", 1)
	if err != nil {
		return nil, fmt.Errorf("error creating texture: %s", err)
	}

	width, height := app.GetWindow().GetGLFWWindow().GetSize()
	camera := renderer.NewCamera(width, height)

	component := &SquareTexture{
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
	c.camera.OnUpdate(app.GetWindow().GetGLFWWindow(), deltaTime)
}

// OnRender .
func (c *SquareTexture) OnRender(app *application.Application, deltaTime float64) {
	app.GetRenderer().Begin(c.camera)
	app.GetRenderer().DrawTexturedQuad(c.texture)
	app.GetRenderer().End()
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
