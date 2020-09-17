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

	maxTextures = 32

	quadVertices = []mgl32.Vec4{
		mgl32.Vec4{-0.5, 0.5, 0.0, 1.0},
		mgl32.Vec4{-0.5, -0.5, 0.0, 1.0},
		mgl32.Vec4{0.5, -0.5, 0.0, 1.0},
		mgl32.Vec4{0.5, 0.5, 0.0, 1.0},
	}
	quadTexCoords = []mgl32.Vec2{
		mgl32.Vec2{0, 1},
		mgl32.Vec2{0, 0},
		mgl32.Vec2{1, 0},
		mgl32.Vec2{1, 1},
	}
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
			Textures: make(map[int]*opengl.Texture),
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
		opengl.VBOLayoutElement{Count: 4, Normalized: false, DataType: opengl.GLDataTypeFloat},
		opengl.VBOLayoutElement{Count: 2, Normalized: false, DataType: opengl.GLDataTypeFloat},
		opengl.VBOLayoutElement{Count: 1, Normalized: false, DataType: opengl.GLDataTypeFloat},
	))

	ibo := opengl.NewIBO(maxIndices)

	vao := opengl.NewVAO()
	vao.AddVBO(vbo)
	vao.SetIBO(ibo)

	samplers := make([]int32, maxTextures)
	for i := 0; i < maxTextures; i++ {
		samplers[i] = int32(i)
	}
	shaderProgram.Bind()
	defer shaderProgram.Unbind()
	shaderProgram.SetUniform1iv("tex", int32(len(samplers)), &samplers[0])

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
		Textures: make(map[int]*opengl.Texture),
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
	// fmt.Printf("> End\n")
	// fmt.Printf("\tVertices: %v\n", r.quadData.Vertices)
	// fmt.Printf("\tIndices: %v\n", r.quadData.Indices)
	// fmt.Printf("\tTextures: %v\n", r.quadData.Textures)
	r.quadVertexBuffer.SetData(r.quadData)
	r.quadVertexArray.GetIBO().SetData(r.quadData)

	for _, t := range r.quadData.Textures {
		t.Bind()
	}
	r.quadShaderProgram.Bind()
	r.quadVertexArray.Bind()

	defer func() {
		for _, t := range r.quadData.Textures {
			t.Unbind()
		}
		r.quadVertexArray.Unbind()
		r.quadShaderProgram.Unbind()
	}()

	count := r.quadVertexArray.GetIBO().GetCount()
	gl.DrawElements(gl.TRIANGLES, int32(count), gl.UNSIGNED_INT, nil)
	// fmt.Printf("< End\n")
}

// DrawTexturedQuad .
func (r *Renderer) DrawTexturedQuad(transform mgl32.Mat4, texture *opengl.Texture) {
	if err := r.quadData.AddTexturedQuad(transform, texture); err != nil {
		panic(err)
	}
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
	Position mgl32.Vec4
	TexCoord mgl32.Vec2
	TexIndex float32
}

// QuadData .
type QuadData struct {
	Textures map[int]*opengl.Texture
	Vertices []QuadVertex
	Indices  []uint32
}

// AddTexturedQuad .
func (d *QuadData) AddTexturedQuad(transform mgl32.Mat4, texture *opengl.Texture) error {
	if err := d.addTexture(texture); err != nil {
		return err
	}
	quadOffset := len(d.Vertices)
	quad := []QuadVertex{
		QuadVertex{
			Position: transform.Mul4x1(quadVertices[0]),
			TexCoord: quadTexCoords[0],
			TexIndex: float32(texture.GetIndex()),
		},
		QuadVertex{
			Position: transform.Mul4x1(quadVertices[1]),
			TexCoord: quadTexCoords[1],
			TexIndex: float32(texture.GetIndex()),
		},
		QuadVertex{
			Position: transform.Mul4x1(quadVertices[2]),
			TexCoord: quadTexCoords[2],
			TexIndex: float32(texture.GetIndex()),
		},
		QuadVertex{
			Position: transform.Mul4x1(quadVertices[3]),
			TexCoord: quadTexCoords[3],
			TexIndex: float32(texture.GetIndex()),
		},
	}
	d.Vertices = append(d.Vertices, quad...)
	indices := []uint32{
		uint32(quadOffset + 0),
		uint32(quadOffset + 1),
		uint32(quadOffset + 2),
		uint32(quadOffset + 2),
		uint32(quadOffset + 3),
		uint32(quadOffset + 0),
	}
	d.Indices = append(d.Indices, indices...)
	return nil
}

func (d *QuadData) addTexture(texture *opengl.Texture) error {
	// already registered
	if _, ok := d.Textures[texture.GetIndex()]; ok {
		return nil
	}
	if len(d.Textures) == maxTextures {
		return fmt.Errorf("maximum texture count reached")
	}
	d.Textures[texture.GetIndex()] = texture
	return nil
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
