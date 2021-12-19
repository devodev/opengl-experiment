package renderer

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"unsafe"

	"github.com/devodev/opengl-experiment/internal/opengl"
	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	defaultBackgroundColor = color.RGBA{51, 75, 75, 1}
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Renderer .
type Renderer struct {
	bgColor color.RGBA

	quadProgram *Quad
}

// New .
func New() (*Renderer, error) {
	r := &Renderer{
		bgColor:     defaultBackgroundColor,
		quadProgram: &Quad{},
	}
	return r, nil
}

// Init .
func (r *Renderer) Init() error {
	// initialize OpenGL
	// *always do this after a call to `window.MakeContextCurrent()`
	if err := gl.Init(); err != nil {
		return fmt.Errorf("error initializing OpenGL: %s", err)
	}

	if os.Getenv("DEBUG") == "true" {
		r.enableDebugging()
	}

	r.enableBlending()

	// initialize quad-related rendering primitives
	if err := r.quadProgram.Init(); err != nil {
		return err
	}

	return nil
}

// Clear .
func (r *Renderer) Clear() {
	// clear buffers
	gl.ClearColor(
		float32(r.bgColor.R)/255,
		float32(r.bgColor.G)/255,
		float32(r.bgColor.B)/255,
		float32(r.bgColor.A)/255,
	)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// BeginQuad .
func (r *Renderer) BeginQuad(cameraController *CameraController) {
	vp := cameraController.GetViewProjectionMatrix()
	r.quadProgram.Begin(vp)
}

// EndQuad .
func (r *Renderer) EndQuad() {
	r.quadProgram.End()
}

// DrawTexturedQuad .
func (r *Renderer) DrawTexturedQuad(quad *TexturedQuad) {
	if err := r.quadProgram.AddTextured(quad); err != nil {
		panic(err)
	}
}

func (r *Renderer) enableDebugging() {
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(func(
		source uint32,
		gltype uint32,
		id uint32,
		severity uint32,
		length int32,
		message string,
		userParam unsafe.Pointer) {

		fmt.Println("[OpenGL DEBUG]")
		fmt.Printf("[OpenGL DEBUG]\tsource (0x%x): %v\n", source, strings.Join(opengl.GlEnums[source], ", "))
		fmt.Printf("[OpenGL DEBUG]\tgltype (0x%x): %v\n", gltype, strings.Join(opengl.GlEnums[gltype], ", "))
		fmt.Printf("[OpenGL DEBUG]\tseverity (0x%x): %v\n", severity, strings.Join(opengl.GlEnums[severity], ", "))
		fmt.Printf("[OpenGL DEBUG]\tid: %v\n", id)
		fmt.Println("[OpenGL DEBUG]\tmessage:")
		lineLen := 90
		msgIdx := 0
		for msgIdx < len(message) {
			fmt.Printf("[OpenGL DEBUG]\t\t%v\n", strings.TrimLeft(message[msgIdx:min(msgIdx+lineLen, len(message))], " "))
			msgIdx += lineLen
		}
		fmt.Println("[OpenGL DEBUG]")
	}, nil)
}

func (r *Renderer) enableBlending() {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
