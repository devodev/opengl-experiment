package renderer

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"unsafe"

	"github.com/devodev/opengl-experiment/internal/opengl"
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

	quadVertexSize = int(unsafe.Sizeof(QuadVertex{}))

	quadVertices = []mgl32.Vec4{
		{-0.5, 0.5, 0.0, 1.0},
		{-0.5, -0.5, 0.0, 1.0},
		{0.5, -0.5, 0.0, 1.0},
		{0.5, 0.5, 0.0, 1.0},
	}
	quadTexCoords = []mgl32.Vec2{
		{0, 1},
		{0, 0},
		{1, 0},
		{1, 1},
	}
	quadLayout = opengl.NewVBOLayout(
		opengl.VBOLayoutElement{Count: 4, Normalized: false, DataType: opengl.GLDataTypeFloat},
		opengl.VBOLayoutElement{Count: 2, Normalized: false, DataType: opengl.GLDataTypeFloat},
		opengl.VBOLayoutElement{Count: 1, Normalized: false, DataType: opengl.GLDataTypeFloat},
	)
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

	// quad-related rendering primitives
	quadVertexArray   *opengl.VAO
	quadVertexBuffer  *opengl.VBO
	quadShaderProgram *opengl.ShaderProgram
	// quad-related batch rendering data
	quadData *QuadData
}

// New .
func New() (*Renderer, error) {
	r := &Renderer{
		bgColor:  defaultBackgroundColor,
		quadData: NewQuadData(),
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
	quadVertexShaderSource := string(append([]byte(quadVertexShader), byte('\x00')))
	quadFragmentShaderSource := string(append([]byte(quadFragmentShader), byte('\x00')))
	shaderProgram, err := opengl.NewShaderProgram(quadVertexShaderSource, quadFragmentShaderSource)
	if err != nil {
		return err
	}

	// create VBO
	vbo, err := opengl.NewVBO(maxVertices * quadVertexSize)
	if err != nil {
		return err
	}
	vbo.SetLayout(quadLayout)

	// create IBO
	ibo := opengl.NewIBO(maxIndices)

	// create VAO and set buffers on it
	vao := opengl.NewVAO()
	vao.AddVBO(vbo)
	vao.SetIBO(ibo)

	//
	samplers := make([]int32, maxTextures)
	for i := 0; i < maxTextures; i++ {
		samplers[i] = int32(i)
	}

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
		float32(r.bgColor.R)/255,
		float32(r.bgColor.G)/255,
		float32(r.bgColor.B)/255,
		float32(r.bgColor.A)/255,
	)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Begin .
func (r *Renderer) Begin(cameraController *CameraController) {
	// reset data each frame
	r.quadData = NewQuadData()

	// set view-projection matrix
	r.quadShaderProgram.Bind()
	vp := cameraController.GetViewProjectionMatrix()
	r.quadShaderProgram.SetUniformMatrix4fv("vp", 1, false, &vp[0])
	r.quadShaderProgram.Unbind()
}

// End .
func (r *Renderer) End() {
	// fmt.Printf("> End\n")
	// fmt.Printf("\tVertices: %v\n", r.quadData.Vertices)
	// fmt.Printf("\tIndices: %v\n", r.quadData.Indices)
	// fmt.Printf("\tTextures: %v\n", r.quadData.Textures)
	r.quadVertexBuffer.SetData(r.quadData)
	r.quadVertexArray.IBO().SetData(r.quadData)

	for _, t := range r.quadData.Textures {
		t.Bind()
	}
	r.quadShaderProgram.Bind()
	r.quadVertexArray.Bind()

	// actual draw call
	count := r.quadVertexArray.IBO().Count()
	gl.DrawElements(gl.TRIANGLES, int32(count), gl.UNSIGNED_INT, nil)

	for _, t := range r.quadData.Textures {
		t.Unbind()
	}
	r.quadVertexArray.Unbind()
	r.quadShaderProgram.Unbind()
	// fmt.Printf("< End\n")
}

// DrawTexturedQuad .
func (r *Renderer) DrawTexturedQuad(quad *TexturedQuad) {
	if err := r.quadData.AddTexturedQuad(quad); err != nil {
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

// QuadVertex .
type QuadVertex struct {
	Position mgl32.Vec4
	TexCoord mgl32.Vec2
	TexIndex float32
}

// QuadData .
type QuadData struct {
	Textures map[int]opengl.Texture
	Vertices []QuadVertex
	Indices  []uint32
}

func NewQuadData() *QuadData {
	return &QuadData{
		Textures: make(map[int]opengl.Texture),
		Vertices: make([]QuadVertex, 0, maxVertices),
		Indices:  make([]uint32, 0, maxIndices),
	}
}

// AddTexturedQuad .
func (d *QuadData) AddTexturedQuad(quad *TexturedQuad) error {
	if err := d.addTexture(quad.Texture); err != nil {
		return err
	}

	// add indices
	quadOffset := len(d.Vertices)
	indices := []uint32{
		uint32(quadOffset),
		uint32(quadOffset + 1),
		uint32(quadOffset + 2),
		uint32(quadOffset + 2),
		uint32(quadOffset + 3),
		uint32(quadOffset + 0),
	}
	d.Indices = append(d.Indices, indices...)

	// add vertices
	for i := 0; i < len(quadVertices); i++ {
		vertex := QuadVertex{
			Position: quad.Transform.Mul4x1(quadVertices[i]),
			TexCoord: quadTexCoords[i],
			TexIndex: float32(quad.Texture.Index()),
		}
		d.Vertices = append(d.Vertices, vertex)
	}

	return nil
}

func (d *QuadData) addTexture(texture opengl.Texture) error {
	// noop if already registered
	if _, ok := d.Textures[texture.Index()]; ok {
		return nil
	}
	d.Textures[texture.Index()] = texture
	return nil
}

// VBOGLPtr .
func (d *QuadData) VBOGLPtr() unsafe.Pointer {
	return gl.Ptr(d.Vertices)
}

// VBOSize .
func (d *QuadData) VBOSize() int {
	return quadVertexSize * len(d.Vertices)
}

// IBOGLPtr .
func (d *QuadData) IBOGLPtr() unsafe.Pointer {
	return gl.Ptr(d.Indices)
}

// IBOCount .
func (d *QuadData) IBOCount() int32 {
	return int32(len(d.Indices))
}
