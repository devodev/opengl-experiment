package opengl

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// ShaderProgram .
type ShaderProgram struct {
	id               uint32
	uniformLocations map[string]int32
}

// NewShaderProgram requires that both vertex and fragment shader
// sources be null terminated strings.
func NewShaderProgram(vertexShaderSource, fragmentShaderSource string) (*ShaderProgram, error) {
	// compile shaders
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return nil, fmt.Errorf("could not compile vertex shader: %s", err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, fmt.Errorf("could not compile fragment shader: %s", err)
	}

	// create shader program
	shaderProgramID := gl.CreateProgram()
	gl.AttachShader(shaderProgramID, vertexShader)
	gl.AttachShader(shaderProgramID, fragmentShader)

	// link shaders previously compiled
	gl.LinkProgram(shaderProgramID)
	gl.ValidateProgram(shaderProgramID)

	// free memory
	gl.DetachShader(shaderProgramID, vertexShader)
	gl.DetachShader(shaderProgramID, fragmentShader)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

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
	s.Bind()
	gl.Uniform1f(s.getUniformLocation(name), v0)
	s.Unbind()
}

// SetUniform1i .
func (s *ShaderProgram) SetUniform1i(name string, v0 int32) {
	s.Bind()
	gl.Uniform1i(s.getUniformLocation(name), v0)
	s.Unbind()
}

// SetUniform1iv .
func (s *ShaderProgram) SetUniform1iv(name string, count int32, value *int32) {
	s.Bind()
	gl.Uniform1iv(s.getUniformLocation(name), count, value)
	s.Unbind()
}

// SetUniform4f .
func (s *ShaderProgram) SetUniform4f(name string, v0, v1, v2, v3 float32) {
	s.Bind()
	gl.Uniform4f(s.getUniformLocation(name), v0, v1, v2, v3)
	s.Unbind()
}

// SetUniformMatrix4fv .
func (s *ShaderProgram) SetUniformMatrix4fv(name string, count int32, transpose bool, value *float32) {
	s.Bind()
	gl.UniformMatrix4fv(s.getUniformLocation(name), count, transpose, value)
	s.Unbind()
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
