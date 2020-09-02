package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// ShaderProgram .
type ShaderProgram struct {
	id               uint32
	uniformLocations map[string]int32
}

// NewShaderProgram .
func NewShaderProgram(vertexShaderFilepath, fragmentShaderFilepath string) (*ShaderProgram, error) {
	// TODO: we might want to only create an empty program on init
	// TODO: and provide methods for attaching shaders.

	// read files into memory
	// this might need to be put in a helper function
	// or provide another constructor using strings directly
	vertexShaderSource, err := ioutil.ReadFile(vertexShaderFilepath)
	if err != nil {
		return nil, err
	}
	fragmentShaderSource, err := ioutil.ReadFile(fragmentShaderFilepath)
	if err != nil {
		return nil, err
	}
	vertexShaderSource = append(vertexShaderSource, byte('\x00'))
	fragmentShaderSource = append(fragmentShaderSource, byte('\x00'))

	// compile shaders
	vertexShader, err := compileShader(string(vertexShaderSource), gl.VERTEX_SHADER)
	if err != nil {
		return nil, fmt.Errorf("could not compile vertex shader: %s", err)
	}
	fragmentShader, err := compileShader(string(fragmentShaderSource), gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, fmt.Errorf("could not compile fragment shader: %s", err)
	}

	// create shader program and link the shaders previously compiled
	shaderProgramID := gl.CreateProgram()
	gl.AttachShader(shaderProgramID, vertexShader)
	gl.AttachShader(shaderProgramID, fragmentShader)
	gl.LinkProgram(shaderProgramID)
	gl.ValidateProgram(shaderProgramID)

	defer func() {
		gl.DetachShader(shaderProgramID, vertexShader)
		gl.DetachShader(shaderProgramID, fragmentShader)
		gl.DeleteShader(vertexShader)
		gl.DeleteShader(fragmentShader)
	}()

	if err := retrieveProgramLinkError(shaderProgramID); err != nil {
		return nil, err
	}

	shaderProgram := &ShaderProgram{id: shaderProgramID}

	return shaderProgram, nil
}

// Bind .
func (s *ShaderProgram) Bind() {
	gl.UseProgram(s.id)
}

// Unbind .
func (s *ShaderProgram) Unbind() {
	gl.UseProgram(0)
}

func (s *ShaderProgram) getUniformLocation(name string) int32 {
	location, ok := s.uniformLocations[name]
	if !ok {
		location = gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	}
	return location
}

// SetUniform1f .
func (s *ShaderProgram) SetUniform1f(name string, v0 float32) {
	gl.Uniform1f(s.getUniformLocation(name), v0)
}

// SetUniform4f .
func (s *ShaderProgram) SetUniform4f(name string, v0, v1, v2, v3 float32) {
	gl.Uniform4f(s.getUniformLocation(name), v0, v1, v2, v3)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shaderProgramID := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shaderProgramID, 1, csources, nil)
	free()

	gl.CompileShader(shaderProgramID)
	if err := retrieveShaderCompileError(shaderProgramID); err != nil {
		return 0, err
	}
	return shaderProgramID, nil
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
