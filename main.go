package main

import (
	"runtime"

	"github.com/devodev/opengl-experimentation/internal/opengl"
)

func init() {
	// All calls to GLFW/OpenGL must happen on the main thread.
	// This locks the calling goroutine(main here) to
	// the current OS Thread(main here).
	runtime.LockOSThread()
}

func main() {
	logger := NewLogger()
	app, err := NewApplication(WithLoggerOption(logger))
	if err != nil {
		logger.Errorf("error creating application: %s", err)
		return
	}
	close, err := app.Init()
	if err != nil {
		logger.Errorf("error initializing application: %s", err)
		return
	}
	defer close()

	// create a program with default shaders
	shaderProgram, err := opengl.NewShaderProgram(
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

	vbo, err := opengl.NewVBO(square, opengl.FLOAT)
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

	ibo := opengl.NewIBO(squareIndices)

	vao := opengl.NewVAO()
	vao.AddVBO(vbo)

	app.Loop(func(a *Application) {
		a.draw(vao, ibo, shaderProgram)
	})
}
