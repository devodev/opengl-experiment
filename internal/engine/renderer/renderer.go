package renderer

import (
	"fmt"
	"image/color"

	"github.com/devodev/opengl-experimentation/internal/engine/window"
	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	defaultBackgroundColor = color.RGBA{51, 75, 75, 1}
)

// Renderer .
type Renderer struct {
	window *window.Window

	backgroundColor color.RGBA
}

// New .
func New(window *window.Window) (*Renderer, error) {
	r := &Renderer{
		window:          window,
		backgroundColor: defaultBackgroundColor,
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

	// set blending function
	r.setBlending()

	return nil
}

// Clear .
func (r *Renderer) Clear() {
	// clear buffers
	gl.ClearColor(
		float32(r.backgroundColor.R)/255,
		float32(r.backgroundColor.G)/255,
		float32(r.backgroundColor.B)/255,
		float32(r.backgroundColor.A)/255,
	)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *Renderer) setBlending() {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
