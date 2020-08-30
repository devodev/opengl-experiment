package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// shaders
var (
	vertexShaderSource = `
        #version 410
        in vec3 vp;
        void main() {
            gl_Position = vec4(vp, 1.0);
        }
    ` + "\x00"

	fragmentShaderSource = `
        #version 410
        out vec4 frag_colour;
        void main() {
            frag_colour = vec4(1, 1, 1, 1);
        }
    ` + "\x00"
)

// window attributes
var (
	windowWidth  = 640
	windowHeight = 480
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
	window, err := initGLFW()
	if err != nil {
		logger.Error(err)
		return
	}
	defer glfw.Terminate()

	logger.Printf("GLFW version: %s", glfw.GetVersionString())

	window.MakeContextCurrent()

	// initialize OpenGL
	// *always do this after a call to `window.MakeContextCurrent()`
	program, err := initOpenGL()
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	// create a triangle VAO
	triangle := []float32{
		// X, Y, Z
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}

	var components []*Component
	components = append(components, NewComponent(triangle))

	// start event loop
	for !window.ShouldClose() {
		draw(components, window, program)
	}
}

// Component .
type Component struct {
	points []float32
	vao    uint32
}

// NewComponent .
func NewComponent(points []float32) *Component {
	return &Component{points: points, vao: makeVao(points)}
}

// Draw .
func (c *Component) Draw() {
	gl.BindVertexArray(c.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(c.points)/3))
}

func draw(components []*Component, window *glfw.Window, program uint32) {
	// clear old frame and register the program to use
	// for subsequent draw calls
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// actual drawing
	for _, c := range components {
		c.Draw()
	}

	// events and buffer swap
	glfw.PollEvents()
	window.SwapBuffers()
}

func initGLFW() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize GLFW: %s", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitle, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create window: %s", err)
	}
	return window, nil
}

func initOpenGL() (uint32, error) {
	if err := gl.Init(); err != nil {
		return 0, fmt.Errorf("could not initialize OpenGL: %s", err)
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("could not compile vertex shader: %s", err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("could not compile fragment shader: %s", err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog, nil
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()

	gl.CompileShader(shader)

	// error handling
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
