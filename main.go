package main

import (
	"fmt"
	"math"
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
	shaderProgram, err := NewShaderProgram(
		"assets/shaders/vertexDefault.glsl",
		"assets/shaders/fragmentVariableColor.glsl",
	)
	if err != nil {
		logger.Error(err)
		return
	}

	// create vao
	square := []float32{
		// X, Y, Z
		-0.5, 0.5, 1, 0, // top-left
		-0.5, -0.5, 0, 1, // bottom-left
		0.5, -0.5, 0, 0, // bottom-right
		0.5, 0.5, 1, 1, // top-right
	}
	squareIndices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	vbo, err := NewVBO(square, gl.FLOAT)
	if err != nil {
		logger.Error(err)
		return
	}
	// This declares the interleaving layout of VBO.
	// Here, each `square` vertex contains:
	// 2 * vec2 float32 normalized values
	// *These map directly to vertex shader attributes
	vbo.AddElement(2, false)
	vbo.AddElement(2, false)

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

func createWindow(width, height int, title string) (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create window: %s", err)
	}
	return window, nil
}

func draw(vao *VAO, ibo *IBO, shaderProgram *ShaderProgram) {
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
