package engine

import (
	"fmt"
	"image/color"

	"github.com/devodev/opengl-experimentation/internal/engine/components"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	defaultWindowWidth     = 1024
	defaultWindowHeight    = 768
	defaultWindowTitle     = "Application"
	defaultWindowResizable = true
	defaultBackgroundColor = color.RGBA{51, 75, 75, 1}
)

// processInput state.
// TODO: store these in a struct somewhere in the app
// TODO: or something..
var (
	previousKeySpaceState = glfw.Release
	currentPolygonMode    = gl.TRIANGLES
)

// Application needs to be used exclusively
// on the main thread.
//
// Make sure you call runtime.LockOSThread()
// from the main function before initiliazing
// an application.
type Application struct {
	windowWidth     int
	windowHeight    int
	windowTitle     string
	windowResizable bool

	running bool

	backgroundColor color.RGBA

	window *glfw.Window
	logger *SimpleLogger

	components []components.Component
}

// ApplicationOption .
type ApplicationOption func(*Application) error

// WithLoggerOption .
func WithLoggerOption(logger *SimpleLogger) ApplicationOption {
	return func(a *Application) error {
		a.logger = logger
		return nil
	}
}

// NewApplication .
func NewApplication(options ...ApplicationOption) (*Application, error) {
	app := &Application{
		windowWidth:     defaultWindowWidth,
		windowHeight:    defaultWindowHeight,
		windowTitle:     defaultWindowTitle,
		windowResizable: defaultWindowResizable,
		backgroundColor: defaultBackgroundColor,
		running:         true,
		logger:          NewLogger(),
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
	glfw.Terminate()
	a.running = false
	return nil
}

// Run .
func (a *Application) Run() error {
	if !a.running {
		return fmt.Errorf("application is closed")
	}
	// sync with monitor refresh rate
	glfw.SwapInterval(1)

	for _, c := range a.components {
		c.OnInit(a.window)
	}

	for !a.window.ShouldClose() {
		a.clear()

		for _, c := range a.components {
			c.OnUpdate(a.window)
		}

		glfw.PollEvents()
		a.window.SwapBuffers()
	}
	return nil
}

// AddComponent .
func (a *Application) AddComponent(c components.Component) {
	a.components = append(a.components, c)
}

func (a *Application) init() error {
	window, err := CreateWindow(a.windowWidth, a.windowHeight, a.windowTitle, a.windowResizable)
	if err != nil {
		return fmt.Errorf("error creating window: %s", err)
	}
	a.window = window
	a.logger.Printf("GLFW version: %s", glfw.GetVersionString())

	a.window.MakeContextCurrent()

	// initialize OpenGL
	// *always do this after a call to `window.MakeContextCurrent()`
	if err := gl.Init(); err != nil {
		return fmt.Errorf("error initializing OpenGL: %s", err)
	}
	a.logger.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	// set blending function
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// set window resize callback
	if a.windowResizable {
		a.window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
			gl.Viewport(0, 0, int32(width), int32(height))
		})
	}
	return nil
}

func (a *Application) clear() {
	// clear buffers
	gl.ClearColor(
		float32(a.backgroundColor.R)/255,
		float32(a.backgroundColor.G)/255,
		float32(a.backgroundColor.B)/255,
		float32(a.backgroundColor.A)/255,
	)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
