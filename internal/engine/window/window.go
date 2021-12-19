package window

import (
	"fmt"
	"os"

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

	keyPressed  map[Key]bool
	keyReleased map[Key]bool

	window *glfw.Window
}

// New .
func New(options ...Option) (*Window, error) {
	window := &Window{
		width:       defaultWindowWidth,
		height:      defaultWindowHeight,
		title:       defaultWindowTitle,
		resizable:   defaultWindowResizable,
		keyPressed:  make(map[Key]bool),
		keyReleased: make(map[Key]bool),
	}

	for _, opt := range options {
		if err := opt(window); err != nil {
			return nil, err
		}
	}
	return window, nil
}

func (w *Window) initialized() bool {
	return w.window != nil
}

// Init .
func (w *Window) Init() error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("error initializing GLFW: %s", err)
	}

	// set opengl hints on the window
	glfw.WindowHint(glfw.OpenGLProfile, glfwOpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfwOpenGLForwardCompatible)
	glfw.WindowHint(glfw.ContextVersionMajor, glfwMajorVersion)
	glfw.WindowHint(glfw.ContextVersionMinor, glfwMinorVersion)

	if w.resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	// create a window
	window, err := glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
	if err != nil {
		w.Close()
		return fmt.Errorf("error creating window: %s", err)
	}
	w.window = window
	w.window.MakeContextCurrent()

	w.window.SetKeyCallback(func(ww *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		mKey := Key(key)
		switch action {
		case glfw.Press:
			delete(w.keyReleased, mKey)
			if _, ok := w.keyPressed[mKey]; !ok {
				w.keyPressed[mKey] = false
			}
		case glfw.Release:
			delete(w.keyPressed, mKey)
			if _, ok := w.keyReleased[mKey]; !ok {
				w.keyReleased[mKey] = false
			}
		}
	})

	// set window resize callback
	if w.resizable {
		w.window.SetFramebufferSizeCallback(func(window *glfw.Window, width int, height int) {
			gl.Viewport(0, 0, int32(width), int32(height))
			w.width = width
			w.height = height
		})
	}

	// set vsync (synchronize buffer swap with monitor refresh rate)
	if os.Getenv("VSYNC") == "true" {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}

	return nil
}

// GetSize .
func (w *Window) GetSize() (int, int) {
	return w.width, w.height
}

// SetSize .
func (w *Window) SetSize(width, height int) {
	if w.initialized() {
		panic("cant set window size: already initialized")
	}
	w.width = width
	w.height = height
}

// ShouldClose .
func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}

// IsFocused .
func (w *Window) IsFocused() bool {
	return w.window.GetAttrib(glfw.Focused) == glfw.True
}

func (w *Window) GetKeyDown(key Key) bool {
	requested, ok := w.keyPressed[key]
	if ok && !requested {
		w.keyPressed[key] = true
		return true
	}
	return false
}

func (w *Window) GetKeyUp(key Key) bool {
	requested, ok := w.keyReleased[key]
	if ok && !requested {
		w.keyReleased[key] = true
		return true
	}
	return false
}

// IsKeyPressed .
func (w *Window) IsKeyPressed(key Key) bool {
	_, ok := w.keyPressed[key]
	return ok
}

// IsKeyReleased .
func (w *Window) IsKeyReleased(key Key) bool {
	_, ok := w.keyReleased[key]
	return ok
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
	w.window = nil
	return nil
}

// GetGLFWWindow .
func (w *Window) GetGLFWWindow() *glfw.Window {
	return w.window
}
