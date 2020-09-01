package main

import (
	"io/ioutil"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// window attributes
var (
	windowWidth  = 1024
	windowHeight = 768
	windowTitle  = "Hello World"
)

var (
	previousKeySpaceState = glfw.Release
	currentPolygonMode    = gl.TRIANGLES
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
	if err := glfw.Init(); err != nil {
		logger.Errorf("could not initialize GLFW: %s", err)
		return
	}
	defer glfw.Terminate()
	logger.Printf("GLFW version: %s", glfw.GetVersionString())

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// create a window
	window, err := createWindow(windowWidth, windowHeight, windowTitle)
	if err != nil {
		logger.Error(err)
		return
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	// initialize OpenGL
	// *always do this after a call to `window.MakeContextCurrent()`
	if err := gl.Init(); err != nil {
		logger.Errorf("could not initialize OpenGL: %s", err)
		return
	}
	logger.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	// set window resize callback
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// create a program with default shaders
	vertexShaderSource, err := ioutil.ReadFile("assets/shaders/vertexDefault.glsl")
	if err != nil {
		logger.Error(err)
		return
	}
	fragmentShaderSource, err := ioutil.ReadFile("assets/shaders/fragmentVariableColor.glsl")
	if err != nil {
		logger.Error(err)
		return
	}
	vertexShaderSource = append(vertexShaderSource, byte('\x00'))
	fragmentShaderSource = append(fragmentShaderSource, byte('\x00'))

	shaderProgram, err := createProgram(string(vertexShaderSource), string(fragmentShaderSource))
	if err != nil {
		logger.Error(err)
		return
	}

	// create vao
	square := []float32{
		// X, Y, Z
		-0.5, 0.5, 0, // top-left
		-0.5, -0.5, 0, // bottom-left
		0.5, -0.5, 0, // bottom-right
		0.5, 0.5, 0, // top-right
	}
	squareIndices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	vbo := NewVBO()
	vbo.AddElement(square, 3, gl.FLOAT, true)

	ibo := NewIBO(squareIndices)

	vao := NewVAO()
	vao.AddVBO(vbo)

	// start event loop
	for !window.ShouldClose() {
		processInput(window)

		draw(vao, ibo, shaderProgram)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func processInput(w *glfw.Window) {
	// close window
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}

	// toggle wireframes
	keySpaceState := w.GetKey(glfw.KeySpace)
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
