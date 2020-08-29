package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	windowWidth  = 640
	windowHeight = 480
	windowTitle  = "Hello World"
)

func main() {
	// Create logger
	logger := NewLogger()

	// all calls to GLFW must be run on main thread
	runtime.LockOSThread()

	// initialize GLFW
	if err := glfw.Init(); err != nil {
		logger.FatalError(GLFWInitError, err)
	}
	defer glfw.Terminate()
	glfwVersion := glfw.GetVersionString()
	logger.Printf("GLFW version: %s", glfwVersion)

	// create a window
	window, err := createWindow(windowWidth, windowHeight, windowTitle)
	if err != nil {
		logger.FatalError(GLFWCreateWindowError, err)
	}
	window.MakeContextCurrent()

	// initialize OpenGL
	// *always do this after a call to `window.MakeContextCurrent()`
	if err := gl.Init(); err != nil {
		logger.FatalError(GLInitError, err)
	}
	glVersion := gl.GoStr(gl.GetString(gl.VERSION))
	logger.Printf("OpenGL version: %s", glVersion)

	// create a program
	prog := gl.CreateProgram()
	gl.LinkProgram(prog)

	// start event loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func createWindow(width, height int, title string) (*glfw.Window, error) {
	// Not working on wsl2
	//glfw.WindowHint(glfw.ContextVersionMajor, 4)
	//glfw.WindowHint(glfw.ContextVersionMinor, 1)
	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	return window, nil
}
