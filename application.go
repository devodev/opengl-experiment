package main

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	glfwMajorVersion            = 4
	glfwMinorVersion            = 6
	glfwOpenGLCoreProfile       = glfw.OpenGLCoreProfile
	glfwOpenGLForwardCompatible = glfw.True
)

var (
	defaultWindowWidth  = 1024
	defaultWindowHeight = 768
	defaultWindowTitle  = "Application"
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
	windowWidth  int
	windowHeight int
	windowTitle  string

	window *glfw.Window
	logger *SimpleLogger
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
		windowWidth:  defaultWindowWidth,
		windowHeight: defaultWindowHeight,
		windowTitle:  defaultWindowTitle,
		logger:       NewLogger(),
	}
	for _, option := range options {
		if err := option(app); err != nil {
			return nil, err
		}
	}
	return app, nil
}

// CloseFn .
type CloseFn func()

// Init .
func (a *Application) Init() (CloseFn, error) {
	close := func() {}
	// initialize GLFW
	if err := glfw.Init(); err != nil {
		return close, fmt.Errorf("error initializing GLFW: %s", err)
	}
	close = glfw.Terminate
	a.logger.Printf("GLFW version: %s", glfw.GetVersionString())

	// set glfw hints
	glfw.WindowHint(glfw.ContextVersionMajor, glfwMajorVersion)
	glfw.WindowHint(glfw.ContextVersionMinor, glfwMinorVersion)
	glfw.WindowHint(glfw.OpenGLProfile, glfwOpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfwOpenGLForwardCompatible)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	// create a window
	window, err := glfw.CreateWindow(a.windowWidth, a.windowHeight, a.windowTitle, nil, nil)
	if err != nil {
		return close, fmt.Errorf("error creating window: %s", err)
	}
	a.window = window
	a.window.MakeContextCurrent()
	// sync with monitor refresh rate
	glfw.SwapInterval(1)

	// initialize OpenGL
	// *always do this after a call to `window.MakeContextCurrent()`
	if err := gl.Init(); err != nil {
		return close, fmt.Errorf("error initializing OpenGL: %s", err)
	}
	a.logger.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	// set window resize callback
	a.window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	return close, nil
}

// Loop .
func (a *Application) Loop(fn func(*Application)) {
	for !a.window.ShouldClose() {
		a.processInput()

		fn(a)

		glfw.PollEvents()
		a.window.SwapBuffers()
	}
}

func (a *Application) draw(vao *VAO, ibo *IBO, shaderProgram *ShaderProgram) {
	// clear buffers
	gl.ClearColor(0.2, 0.3, 0.3, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// select shader program
	shaderProgram.Bind()

	// update uniform value
	currentTime := glfw.GetTime()
	greenValue := float32((math.Sin(currentTime) / 2.0) + 0.5)
	shaderProgram.SetUniform1f("variableColor", greenValue)

	// draw
	vao.Bind()
	// TODO: might want to benchmark this and maybe remove them
	// TODO: or do that only in debug mode or something..
	defer vao.Unbind()

	ibo.Bind()
	// TODO: might want to benchmark this and maybe remove them
	// TODO: or do that only in debug mode or something..
	defer ibo.Unbind()

	gl.DrawElements(gl.TRIANGLES, ibo.count, gl.UNSIGNED_INT, nil)
}

func (a *Application) processInput() {
	// close window
	if a.window.GetKey(glfw.KeyEscape) == glfw.Press {
		a.window.SetShouldClose(true)
	}

	// toggle wireframes
	keySpaceState := a.window.GetKey(glfw.KeySpace)
	if keySpaceState == glfw.Release && previousKeySpaceState == glfw.Press {
		if currentPolygonMode == gl.LINE {
			currentPolygonMode = gl.FILL
		} else {
			currentPolygonMode = gl.LINE
		}
		gl.PolygonMode(gl.FRONT_AND_BACK, uint32(currentPolygonMode))
	}
	previousKeySpaceState = keySpaceState
}
