package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// shaders
var (
	vertexShaderSource = `
        #version 460
        layout (location = 0) in vec3 vp;
        void main() {
            gl_Position = vec4(vp, 1.0);
        }
    ` + "\x00"

	fragmentShaderSource = `
        #version 460
        out vec4 frag_colour;
        void main() {
            frag_colour = vec4(1, 0.5, 0.2, 1);
        }
    ` + "\x00"
)

// window attributes
var (
	windowWidth  = 1024
	windowHeight = 768
	windowTitle  = "Hello World"
)

func init() {
	// All calls to GLFW must be run on main thread
	// This locks the calling goroutine(main here) to the current OS Thread
	runtime.LockOSThread()
}

func main() {
	// Create logger
	logger := NewLogger()

	// initialize GLFW
	window, err := initGLFW(windowWidth, windowHeight, windowTitle)
	if err != nil {
		logger.Error(err)
		return
	}
	defer glfw.Terminate()
	logger.Printf("GLFW version: %s", glfw.GetVersionString())

	window.MakeContextCurrent()

	// initialize OpenGL
	// *always do this after a call to `window.MakeContextCurrent()`
	program, err := initOpenGL(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	// set window resize callback
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// set background color
	gl.ClearColor(0.2, 0.3, 0.3, 1)

	// create a component
	square := []float32{
		// X, Y, Z
		-0.5, 0.5, 0, // top-left
		-0.5, -0.5, 0, // bottom-left
		0.5, -0.5, 0, // bottom-right

		-0.5, 0.5, 0, // top-left
		0.5, 0.5, 0, // top-right
		0.5, -0.5, 0, // bottom-right
	}

	var components []*Component
	components = append(components, NewComponent(square))

	// start event loop
	for !window.ShouldClose() {
		draw(components, window, program)
	}
}
