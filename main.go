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
        #version 460
        layout (location = 0) in vec3 vp;
        void main() {
            gl_Position = vec4(vp, 1.0);
        }
    ` + "\x00"

	fragmentShaderSource = `
        #version 460
        out vec4 frag_colour;
        void main() {
            frag_colour = vec4(1, 0.5, 0.2, 1);
        }
    ` + "\x00"
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
	window, err := initGLFW(windowWidth, windowHeight, windowTitle)
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

	// set window resize callback
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// set background color
	gl.ClearColor(0.2, 0.3, 0.3, 1)

	// create a component
	square := []float32{
		// X, Y, Z
		-0.5, 0.5, 0, // top-left
		-0.5, -0.5, 0, // bottom-left
		0.5, -0.5, 0, // bottom-right

		-0.5, 0.5, 0, // top-left
		0.5, 0.5, 0, // top-right
		0.5, -0.5, 0, // bottom-right
	}

	var components []*Component
	components = append(components, NewComponent(square))

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

	// Unbinding is optional if we always bind a VAO before a draw call
	// Also, would like to benchmark this
	// It is still safer to unbind so that if someone tries
	// to draw without binding a VAO prior, it fails right away
	gl.BindVertexArray(0)
}

func draw(components []*Component, window *glfw.Window, program uint32) {
	processInput(window)

	// clear buffers
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// actual drawing
	gl.UseProgram(program)
	for _, c := range components {
		c.Draw()
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func processInput(w *glfw.Window) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
}

func initGLFW(width, height int, title string) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize GLFW: %s", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
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

	// error handling
	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return prog, nil
}

func makeVao(points []float32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	// VertexAttribPointer index refers to `layout (location = 0) ` in the vertex shader
	// stride can be set to 0 when the values are tightly packed
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*4), nil)
	gl.EnableVertexAttribArray(0)

	// unbind objects
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

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
