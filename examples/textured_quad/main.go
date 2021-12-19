package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/devodev/opengl-experiment/internal/engine"
	"github.com/devodev/opengl-experiment/internal/engine/application"
	"github.com/devodev/opengl-experiment/internal/engine/renderer"
	"github.com/devodev/opengl-experiment/internal/opengl"
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

	application.SetWindowSize(1024, 768)
	if err := application.Run(); err != nil {
		logger.Errorf("error running application: %s", err)
		return
	}
}

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

	w, h := application.GetWindow().GetSize()
	// cameraController := renderer.NewCameraController(renderer.NewCameraPerspective(w, h))
	c.cameraController = renderer.NewCameraController(renderer.NewCameraOrthographic(w, h))

	return nil
}

// OnUpdate .
func (c *SquareTextureLayer) OnUpdate(deltaTime float64) {
	c.cameraController.OnUpdate(application.GetWindow(), deltaTime)
}

// OnRender .
func (c *SquareTextureLayer) OnRender(deltaTime float64) {
	application.GetRenderer().BeginQuad(c.cameraController)
	for _, q := range c.quads {
		application.GetRenderer().DrawTexturedQuad(q)
	}
	application.GetRenderer().EndQuad()
}
