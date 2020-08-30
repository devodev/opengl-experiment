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
	program, err := createProgram(string(vertexShaderSource)+"\x00", string(fragmentShaderSource)+"\x00")
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
		0, 3, 2,
	}
	vao := makeVao(square, squareIndices)

	// draw wireframes
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// start event loop
	for !window.ShouldClose() {
		draw(vao, window, program)
	}
}
