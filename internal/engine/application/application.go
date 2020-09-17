package application

import (
	"errors"
	"fmt"

	"github.com/devodev/opengl-experimentation/internal/engine"
	"github.com/devodev/opengl-experimentation/internal/engine/renderer"
	"github.com/devodev/opengl-experimentation/internal/engine/window"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Errors
var (
	ErrAlreadyClosed = errors.New("application already closed")
)

// Application needs to be used exclusively
// on the main thread.
//
// Make sure you call runtime.LockOSThread()
// from the main function before initiliazing
// an application.
type Application struct {
	running bool

	window       *window.Window
	renderer     *renderer.Renderer
	logger       *engine.SimpleLogger
	frameCounter *FrameCounter

	layers []Layer
}

// New .
func New(options ...Option) (*Application, error) {
	window, err := window.New()
	if err != nil {
		return nil, err
	}
	renderer, err := renderer.New()
	if err != nil {
		return nil, err
	}
	app := &Application{
		window:       window,
		renderer:     renderer,
		logger:       engine.NewLogger(),
		frameCounter: NewFrameCounter(),
	}
	for _, option := range options {
		if err := option(app); err != nil {
			return nil, err
		}
	}
	if err := app.init(); err != nil {
		return nil, err
	}
	return app, nil
}

// Run .
func (a *Application) Run() error {
	a.running = true

	// init components
	for _, layer := range a.layers {
		layer.OnInit(a)
	}

	a.frameCounter.Init(glfw.GetTime())

	// main loop
	for a.running {
		a.frameCounter.OnUpdate(glfw.GetTime())
		deltaTime := a.frameCounter.GetDelta()

		a.renderer.Clear()
		for _, layer := range a.layers {
			layer.OnUpdate(a, deltaTime)
			layer.OnRender(a, deltaTime)
		}
		a.onUpdate()
	}
	return nil
}

// Close .
func (a *Application) Close() error {
	if !a.running {
		return ErrAlreadyClosed
	}
	a.RequestClose()
	if err := a.window.Close(); err != nil {
		return err
	}
	return nil
}

// RequestClose .
func (a *Application) RequestClose() {
	a.running = false
}

// AddLayer .
func (a *Application) AddLayer(l Layer) {
	a.layers = append(a.layers, l)
}

func (a *Application) init() error {
	if err := a.window.Init(); err != nil {
		return fmt.Errorf("error initializing window: %v", err)
	}
	a.logger.Printf("GLFW version: %s", glfw.GetVersionString())

	if err := a.renderer.Init(); err != nil {
		return fmt.Errorf("error initializing renderer: %v", err)
	}
	a.logger.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))
	return nil
}

// GetWindow .
func (a *Application) GetWindow() *window.Window {
	return a.window
}

// GetRenderer .
func (a *Application) GetRenderer() *renderer.Renderer {
	return a.renderer
}

func (a *Application) onUpdate() {
	if a.window.ShouldClose() {
		a.RequestClose()
		return
	}
	glfw.PollEvents()
	a.window.GetGLFWWindow().SwapBuffers()
}
