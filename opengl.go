package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func makeVaoAndIbo(points []float32, indices []uint32) (uint32, uint32) {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	_ = NewVBO(points)
	ibo := NewIBO(indices)

	// unbind objects
	gl.BindVertexArray(0)

	return vao, ibo.id
}

func draw(vao, ebo uint32, program uint32) {
	// clear buffers
	gl.ClearColor(0.2, 0.3, 0.3, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// select shader program
	gl.UseProgram(program)

	// update uniform value
	timeValue := glfw.GetTime()
	greenValue := float32((math.Sin(timeValue) / 2.0) + 0.5)
	vertexColorLocation := gl.GetUniformLocation(program, gl.Str("variableColor\x00"))
	gl.Uniform4f(vertexColorLocation, 0.0, greenValue, 0.0, 1.0)

	// draw vao
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

	// Unbinding is optional if we always bind a VAO before a draw call
	// Also, would like to benchmark with and without
	// It is still safer to unbind so that if someone tries
	// to draw without binding a VAO prior, it fails right away
	gl.BindVertexArray(0)
}

func createWindow(width, height int, title string) (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create window: %s", err)
	}
	return window, nil
}

func createProgram(vsSource, fsSource string) (uint32, error) {
	vertexShader, err := compileShader(vsSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("could not compile vertex shader: %s", err)
	}
	fragmentShader, err := compileShader(fsSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("could not compile fragment shader: %s", err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	gl.ValidateProgram(prog)

	// free memory once attached to a program
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// dont do this if we need to debug
	// the shaders in the GPU
	gl.DetachShader(prog, vertexShader)
	gl.DetachShader(prog, fragmentShader)

	if err := retrieveProgramLinkError(prog); err != nil {
		return 0, err
	}
	return prog, nil
}

func retrieveProgramLinkError(program uint32) error {
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %v", log)
	}
	return nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()

	gl.CompileShader(shader)
	if err := retrieveShaderCompileError(shader); err != nil {
		return 0, err
	}
	return shader, nil
}

func retrieveShaderCompileError(shader uint32) error {
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to compile shader: %v", log)
	}
	return nil
}
