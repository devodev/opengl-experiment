package window

import (
	"fmt"

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
	defaultWindowWidth     = 1024
	defaultWindowHeight    = 768
	defaultWindowTitle     = "Application"
	defaultWindowResizable = true
)

// Window .
type Window struct {
	width     int
	height    int
	title     string
	resizable bool

	window *glfw.Window
}

// New .
func New(options ...Option) (*Window, error) {
	window := &Window{
		width:     defaultWindowWidth,
		height:    defaultWindowHeight,
		title:     defaultWindowTitle,
		resizable: defaultWindowResizable,
	}

	for _, opt := range options {
		if err := opt(window); err != nil {
			return nil, err
		}
	}
	return window, nil
}

// Init .
func (w *Window) Init() error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("error initializing GLFW: %s", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, glfwMajorVersion)
	glfw.WindowHint(glfw.ContextVersionMinor, glfwMinorVersion)
	glfw.WindowHint(glfw.OpenGLProfile, glfwOpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfwOpenGLForwardCompatible)
	if w.resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	// create a window
	window, err := glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
	if err != nil {
		return fmt.Errorf("error creating window: %s", err)
	}
	w.window = window
	w.window.MakeContextCurrent()

	// set window resize callback
	if w.resizable {
		w.window.SetFramebufferSizeCallback(func(window *glfw.Window, width int, height int) {
			gl.Viewport(0, 0, int32(width), int32(height))
		})
	}

	// sync with monitor refresh rate
	glfw.SwapInterval(1)

	return nil
}

// GetSize .
func (w *Window) GetSize() (int, int) {
	return w.width, w.height
}

// ShouldClose .
func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}

// IsFocused .
func (w *Window) IsFocused() bool {
	return w.window.GetAttrib(glfw.Focused) == glfw.True
}

// IsKeyPressed .
func (w *Window) IsKeyPressed(key Key) bool {
	return w.window.GetKey(glfw.Key(key)) == glfw.Press
}

// IsMouseButtonPressed .
func (w *Window) IsMouseButtonPressed(m MouseButton) bool {
	return w.window.GetMouseButton(glfw.MouseButton(m)) == glfw.Press
}

// GetCursorPos .
func (w *Window) GetCursorPos() (float64, float64) {
	return w.window.GetCursorPos()
}

// Close .
func (w *Window) Close() error {
	glfw.Terminate()
	return nil
}

// GetGLFWWindow .
func (w *Window) GetGLFWWindow() *glfw.Window {
	return w.window
}
