package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/devodev/opengl-experiment/internal/engine"
	"github.com/devodev/opengl-experiment/internal/engine/application"
	"github.com/devodev/opengl-experiment/internal/engine/renderer"
	"github.com/devodev/opengl-experiment/internal/opengl"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	logger := engine.NewLogger()
	application.SetLogger(logger)

	layer := &SquareTextureLayer{}
	application.AddLayer(layer)

	if os.Getenv("PPROF") == "true" {
		application.EnableProfiling()
	}

	if err := application.Run(); err != nil {
		logger.Errorf("error running application: %s", err)
		return
	}
}

var (
	defaultWidth  = 1024
	defaultHeight = 768
)

// SquareTextureLayer .
type SquareTextureLayer struct {
	quads []*renderer.TexturedQuad

	cameraController *renderer.CameraController
}

// OnInit .
func (c *SquareTextureLayer) OnInit() error {
	texture1, err := opengl.NewNRGBATexture("assets/textures/google_logo.png")
	if err != nil {
		return fmt.Errorf("error creating texture1: %s", err)
	}
	texture2, err := opengl.NewNRGBATexture("assets/textures/facebook_logo.png")
	if err != nil {
		return fmt.Errorf("error creating texture2: %s", err)
	}
	texture3, err := opengl.NewNRGBATexture("assets/textures/instagram_logo.png")
	if err != nil {
		return fmt.Errorf("error creating texture3: %s", err)
	}

	c.quads = []*renderer.TexturedQuad{
		{Texture: texture1, Transform: mgl32.Translate3D(-0.5, 0, 2)},
		{Texture: texture2, Transform: mgl32.Translate3D(0.5, 0, 1)},
		{Texture: texture3, Transform: mgl32.Translate3D(0, 0.5, 0.5)},
	}

	application.SetWindowSize(defaultWidth, defaultHeight)

	// cameraController := renderer.NewCameraController(renderer.NewCameraPerspective(defaultWidth, defaultHeight))
	c.cameraController = renderer.NewCameraController(renderer.NewCameraOrthographic(defaultWidth, defaultHeight))

	return nil
}

// OnUpdate .
func (c *SquareTextureLayer) OnUpdate(deltaTime float64) {
	c.processInput()
	c.cameraController.OnUpdate(application.GetWindow(), deltaTime)
}

// OnRender .
func (c *SquareTextureLayer) OnRender(deltaTime float64) {
	application.GetRenderer().Begin(c.cameraController)
	for _, q := range c.quads {
		application.GetRenderer().DrawTexturedQuad(q)
	}
	application.GetRenderer().End()
}

func (c *SquareTextureLayer) processInput() {
	glfwWindow := application.GetWindow().GetGLFWWindow()

	// we lost focus, dont process synthetic events
	if glfwWindow.GetAttrib(glfw.Focused) == glfw.False {
		return
	}

	// close window
	if glfwWindow.GetKey(glfw.KeyEscape) != glfw.Release {
		application.Close()
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
