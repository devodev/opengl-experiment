package main

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	glfwMajorVersion            = 4
	glfwMinorVersion            = 6
	glfwOpenGLCoreProfile       = glfw.OpenGLCoreProfile
	glfwOpenGLForwardCompatible = glfw.True
)

// WindowCleanupFn .
type WindowCleanupFn func()

// CreateWindow .
func CreateWindow(width, height int, title string) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("error initializing GLFW: %s", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, glfwMajorVersion)
	glfw.WindowHint(glfw.ContextVersionMinor, glfwMinorVersion)
	glfw.WindowHint(glfw.OpenGLProfile, glfwOpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfwOpenGLForwardCompatible)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	// create a window
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating window: %s", err)
	}
	window.MakeContextCurrent()

	return window, nil
}
