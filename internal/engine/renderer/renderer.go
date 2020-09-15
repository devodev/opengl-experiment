package renderer

import (
	"fmt"
	"image/color"

	"github.com/devodev/opengl-experimentation/internal/engine/window"
	"github.com/devodev/opengl-experimentation/internal/opengl"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	defaultBackgroundColor = color.RGBA{51, 75, 75, 1}
)

// Renderer .
type Renderer struct {
	window *window.Window

	backgroundColor color.RGBA

	quadVertexArray       *opengl.VAO
	quadIndexBufferObject *opengl.IBO
	quadShaderProgram     *opengl.ShaderProgram
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

	// initialize quad data
	shaderProgram, err := opengl.NewShaderProgram(
		"assets/shaders/vertexTexture.glsl",
		"assets/shaders/fragmentTexture.glsl",
	)
	if err != nil {
		return err
	}

	square := []float32{
		// position(vec2), texCoord(vec2)
		-0.5, 0.5, 0, 1,
		-0.5, -0.5, 0, 0,
		0.5, -0.5, 1, 0,
		0.5, 0.5, 1, 1,
	}
	squareIndices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	vbo, err := opengl.NewVBO(square, opengl.FLOAT)
	if err != nil {
		return err
	}
	vbo.AddElement(2, false)
	vbo.AddElement(2, false)

	vao := opengl.NewVAO()
	vao.AddVBO(vbo)

	ibo := opengl.NewIBO(squareIndices)

	r.quadVertexArray = vao
	r.quadIndexBufferObject = ibo
	r.quadShaderProgram = shaderProgram

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

// DrawQuad .
func (r *Renderer) DrawQuad(mvp mgl32.Mat4, texture *opengl.Texture) {
	r.quadShaderProgram.Bind()
	r.quadShaderProgram.SetUniform1i("tex", int32(texture.GetTextureUnit()-gl.TEXTURE0))
	r.quadShaderProgram.SetUniformMatrix4fv("mvp", 1, false, &mvp[0])

	texture.Bind()
	r.quadVertexArray.Bind()
	r.quadIndexBufferObject.Bind()
	defer func() {
		r.quadShaderProgram.Unbind()
		texture.Unbind()
		r.quadVertexArray.Unbind()
		r.quadIndexBufferObject.Unbind()
	}()

	gl.DrawElements(gl.TRIANGLES, r.quadIndexBufferObject.GetCount(), gl.UNSIGNED_INT, nil)
}

func (r *Renderer) setBlending() {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
