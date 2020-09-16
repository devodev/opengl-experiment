package renderer

import (
	"fmt"
	"image/color"
	"strings"
	"unsafe"

	"github.com/devodev/opengl-experimentation/internal/opengl"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	defaultBackgroundColor = color.RGBA{51, 75, 75, 1}
)

var (
	maxQuads    = 10000
	maxVertices = maxQuads * 4
	maxIndices  = maxQuads * 6
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Renderer .
type Renderer struct {
	backgroundColor color.RGBA

	quadVertexArray   *opengl.VAO
	quadVertexBuffer  *opengl.VBO
	quadShaderProgram *opengl.ShaderProgram

	quadData *QuadData
}

// New .
func New() (*Renderer, error) {
	r := &Renderer{
		backgroundColor: defaultBackgroundColor,
		quadData: &QuadData{
			Vertices: make([]QuadVertex, 0, maxVertices),
			Indices:  make([]uint32, 0, maxIndices),
		},
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

	r.setDebugging()
	r.setBlending()

	// initialize quad data
	quadVertexShaderSource := string(append([]byte(quadVertexShader), byte('\x00')))
	quadFragmentShaderSource := string(append([]byte(quadFragmentShader), byte('\x00')))
	shaderProgram, err := opengl.NewShaderProgram(quadVertexShaderSource, quadFragmentShaderSource)
	if err != nil {
		return err
	}

	vbo, err := opengl.NewVBO(maxVertices * r.quadData.GetVertexSize())
	if err != nil {
		return err
	}
	vbo.SetLayout(opengl.NewVBOLayout(
		opengl.VBOLayoutElement{Count: 2, Normalized: false, DataType: opengl.GLDataTypeFloat},
		opengl.VBOLayoutElement{Count: 2, Normalized: false, DataType: opengl.GLDataTypeFloat},
	))

	ibo := opengl.NewIBO(maxIndices)

	vao := opengl.NewVAO()
	vao.AddVBO(vbo)
	vao.SetIBO(ibo)

	r.quadVertexArray = vao
	r.quadVertexBuffer = vbo
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

// Begin .
func (r *Renderer) Begin(camera *Camera) {
	r.quadData = &QuadData{
		Vertices: make([]QuadVertex, 0, maxVertices),
		Indices:  make([]uint32, 0, maxIndices),
	}

	r.quadShaderProgram.Bind()
	defer r.quadShaderProgram.Unbind()

	vp := camera.GetViewProjectionMatrix()
	r.quadShaderProgram.SetUniformMatrix4fv("vp", 1, false, &vp[0])
}

// End .
func (r *Renderer) End() {
	r.quadVertexBuffer.SetData(r.quadData)
	r.quadVertexArray.GetIBO().SetData(r.quadData)

	r.quadShaderProgram.Bind()
	r.quadVertexArray.Bind()
	r.quadData.Texture.Bind()
	defer func() {
		r.quadVertexArray.Unbind()
		r.quadShaderProgram.Unbind()
		r.quadData.Texture.Unbind()
	}()

	r.quadShaderProgram.SetUniform1i("tex", int32(r.quadData.Texture.GetTextureUnit()-gl.TEXTURE0))

	count := r.quadVertexArray.GetIBO().GetCount()
	gl.DrawElements(gl.TRIANGLES, int32(count), gl.UNSIGNED_INT, nil)
}

// DrawTexturedQuad .
func (r *Renderer) DrawTexturedQuad(texture *opengl.Texture) {
	r.quadData.AddTexturedQuad(texture)
}

func (r *Renderer) setDebugging() {
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
			fmt.Printf("[OpenGL DEBUG]\t\t%v\n", strings.TrimLeft(message[msgIdx:min(msgIdx+lineLen, len(message))], "  "))
			msgIdx += lineLen
		}
		fmt.Println("[OpenGL DEBUG]")
	}, nil)
}

func (r *Renderer) setBlending() {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

// QuadVertex .
type QuadVertex struct {
	Position mgl32.Vec2
	TexCoord mgl32.Vec2
}

// QuadData .
type QuadData struct {
	Vertices []QuadVertex
	Indices  []uint32
	Texture  *opengl.Texture
}

// AddTexturedQuad .
func (d *QuadData) AddTexturedQuad(texture *opengl.Texture) {
	quad := []QuadVertex{
		QuadVertex{
			Position: mgl32.Vec2{-0.5, 0.5},
			TexCoord: mgl32.Vec2{0, 1},
		},
		QuadVertex{
			Position: mgl32.Vec2{-0.5, -0.5},
			TexCoord: mgl32.Vec2{0, 0},
		},
		QuadVertex{
			Position: mgl32.Vec2{0.5, -0.5},
			TexCoord: mgl32.Vec2{1, 0},
		},
		QuadVertex{
			Position: mgl32.Vec2{0.5, 0.5},
			TexCoord: mgl32.Vec2{1, 1},
		},
	}
	d.Vertices = append(d.Vertices, quad...)
	quadOffset := len(d.Vertices) - 4
	indices := []uint32{
		uint32(quadOffset + 0),
		uint32(quadOffset + 1),
		uint32(quadOffset + 2),
		uint32(quadOffset + 2),
		uint32(quadOffset + 3),
		uint32(quadOffset + 0),
	}
	d.Indices = append(d.Indices, indices...)
	d.Texture = texture
}

// GetVBOGLPtr .
func (d *QuadData) GetVBOGLPtr() unsafe.Pointer {
	return gl.Ptr(d.Vertices)
}

// GetVBOSize .
func (d *QuadData) GetVBOSize() int {
	return d.GetVertexSize() * len(d.Vertices)
}

// GetVertexSize .
func (d *QuadData) GetVertexSize() int {
	var quadVertex QuadVertex
	return int(unsafe.Sizeof(quadVertex))
}

// GetIBOGLPtr .
func (d *QuadData) GetIBOGLPtr() unsafe.Pointer {
	return gl.Ptr(d.Indices)
}

// GetIBOCount .
func (d *QuadData) GetIBOCount() int32 {
	return int32(len(d.Indices))
}
