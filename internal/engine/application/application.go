package application

import (
	"fmt"

	"github.com/devodev/opengl-experimentation/internal/engine/components"
	"github.com/devodev/opengl-experimentation/internal/engine/window"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Application needs to be used exclusively
// on the main thread.
//
// Make sure you call runtime.LockOSThread()
// from the main function before initiliazing
// an application.
type Application struct {
	running bool

	window *window.Window
	logger *SimpleLogger

	components []components.Component
}

// New .
func New(options ...Option) (*Application, error) {
	window, err := window.New()
	if err != nil {
		return nil, err
	}
	app := &Application{
		running: true,
		window:  window,
		logger:  NewLogger(),
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

// Close .
func (a *Application) Close() error {
	if !a.running {
		return fmt.Errorf("already closed")
	}
	a.running = false
	if err := a.window.Close(); err != nil {
		return err
	}
	return nil
}

// Run .
func (a *Application) Run() error {
	if !a.running {
		return fmt.Errorf("application is closed")
	}

	// init components
	for _, c := range a.components {
		c.OnInit(a.window.GetGLFWWindow())
	}

	// main loop
	for !a.isAlive() {
		// TODO: move to renderer
		a.window.Clear()
		for _, c := range a.components {
			// TODO: pass application instead of window
			// TODO: add method on app to get Window
			c.OnUpdate(a.window.GetGLFWWindow())
		}
		a.onUpdate()
	}
	return nil
}

// AddComponent .
func (a *Application) AddComponent(c components.Component) {
	a.components = append(a.components, c)
}

func (a *Application) init() error {
	if err := a.window.Init(); err != nil {
		return fmt.Errorf("error initializing window: %v", err)
	}
	a.logger.Printf("GLFW version: %s", glfw.GetVersionString())

	// TODO: v move to renderer v
	// initialize OpenGL
	// *always do this after a call to `window.MakeContextCurrent()`
	if err := gl.Init(); err != nil {
		return fmt.Errorf("error initializing OpenGL: %s", err)
	}
	a.logger.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	// set blending function
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// TODO: ^ move to renderer ^
	return nil
}

func (a *Application) isAlive() bool {
	return a.window.ShouldClose()
}

func (a *Application) onUpdate() {
	glfw.PollEvents()
	a.window.GetGLFWWindow().SwapBuffers()
}
